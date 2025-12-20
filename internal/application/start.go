package application

import (
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/spf13/afero"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/flag_parse"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/ports"
)

func (c *Conversion) Start(opts flag_parse.Opts, fs afero.Fs, logger *slog.Logger) error {
	affRepo := infrastructure.NewAffineRepository()

	if opts.Symmetry < 1 {
		opts.Symmetry = 1
	}

	baseSeed := int64(opts.Seed)
	boundsRnd := rand.New(rand.NewSource(baseSeed))
	if opts.Randomize || len(opts.AffineParams) == 0 || len(opts.Functions) == 0 {
		RandomizeAffineParamsWithOptions(&opts, rand.New(rand.NewSource(time.Now().UnixNano())), RandomizeOptions{
			Count:    2,
			Override: true,
		})
	}

	for i, a := range opts.AffineParams {
		if a.ColorR == 0 && a.ColorG == 0 && a.ColorB == 0 {
			a.ColorR = boundsRnd.Float64()
			a.ColorG = boundsRnd.Float64()
			a.ColorB = boundsRnd.Float64()
		}

		affRepo.AddAffine(a)
		opts.AffineParams[i] = a
	}

	trRepo := infrastructure.NewTransformRepository()
	for _, fnCfg := range opts.Functions {
		v := VariationByName(fnCfg.Name)
		trRepo.AddTransform(domain.FuncTransform{Variation: v, Weight: fnCfg.Weight})
	}

	boundsAcc := infrastructure.NewBoundsAccumulator()

	boundsConv := NewConversion(affRepo, trRepo, boundsRnd, boundsAcc)
	boundsConv.logger = logger
	boundsConv.worker = "bounds"

	if err := boundsConv.Iterate(domain.Point{}, opts.IterationCount); err != nil {
		return err
	}

	minX, maxX, minY, maxY, ok := boundsAcc.Bounds()
	if !ok {
		return fmt.Errorf("fractal produced no points")
	}

	minX, maxX, minY, maxY = infrastructure.AdjustBoundsToAspect(
		minX, maxX, minY, maxY, int(opts.Width), int(opts.Height), 0.1,
	)

	imageAcc := infrastructure.NewImageAccumulator(
		fs,
		int(opts.Width),
		int(opts.Height),
		minX,
		maxX,
		minY-0.15,
		maxY-0.15,
		opts.GammaCorrection,
		opts.Gamma,
		opts.Symmetry,
	)

	if logger != nil {
		logger.Info("rendering image", "width", opts.Width, "height", opts.Height, "iterations", opts.IterationCount, "threads", opts.Threads)
	}

	if err := iterateWithThreads(affRepo, trRepo, imageAcc, baseSeed+1, opts.Threads, opts.IterationCount, logger); err != nil {
		return err
	}

	if logger != nil {
		logger.Info("saving image", "path", opts.OutputPath)
	}

	return imageAcc.RenderPNG(opts.OutputPath)
}

func iterateWithThreads(
	affRepo ports.AffineRepository,
	trRepo ports.TransformRepository,
	sink ports.SampleSink,
	seed int64,
	threads int,
	iterations int,
	logger *slog.Logger,
) error {
	if threads <= 1 {
		rnd := rand.New(rand.NewSource(seed))
		conv := NewConversion(affRepo, trRepo, rnd, sink)
		conv.worker = "single"
		if logger != nil {
			conv.logger = logger.With("worker", conv.worker)
		}
		return conv.Iterate(domain.Point{}, iterations)
	}

	if iterations < threads {
		threads = iterations
	}

	base := iterations / threads
	rem := iterations % threads

	sinks := make([]ports.SampleSink, threads)
	errCh := make(chan error, threads)
	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		idx := i
		iterCount := base
		if idx < rem {
			iterCount++
		}

		if iterCount == 0 {
			continue
		}

		sinks[idx] = sink.CloneEmpty()
		rnd := rand.New(rand.NewSource(seed + int64(idx+1)))
		conv := NewConversion(affRepo, trRepo, rnd, sinks[idx])
		if logger != nil {
			conv.logger = logger.With("worker", idx)
			conv.worker = fmt.Sprintf("worker-%d", idx)
		}

		wg.Add(1)
		go func(iterations int) {
			defer wg.Done()
			if err := conv.Iterate(domain.Point{}, iterations); err != nil {
				errCh <- err
			}
		}(iterCount)
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	for _, s := range sinks {
		sink.MergeFrom(s)
	}
	return nil
}
