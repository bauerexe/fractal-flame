package application

import (
	"log/slog"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/ports"
)

type Conversion struct {
	affRepo ports.AffineRepository
	trRepo  ports.TransformRepository
	rnd     ports.RandomSource
	sink    ports.SampleSink
	logger  *slog.Logger
	worker  string
}

func NewConversion(aff ports.AffineRepository, tr ports.TransformRepository, rnd ports.RandomSource, sink ports.SampleSink) *Conversion {
	return &Conversion{
		affRepo: aff,
		trRepo:  tr,
		rnd:     rnd,
		sink:    sink,
	}
}

func (c *Conversion) Step(save bool, p domain.Point) (domain.Point, error) {
	n := c.affRepo.Len()
	if n <= 0 {
		return domain.Point{}, ErrGenerate
	}

	i := c.rnd.Intn(n)
	p, err := c.affRepo.Apply(i, p)
	if err != nil {
		return domain.Point{}, err
	}

	p = c.trRepo.Apply(p)

	cr, cg, cb := c.affRepo.ColorAt(i)

	const colorAlpha = 0.5

	p.R = (1.0-colorAlpha)*p.R + colorAlpha*cr
	p.G = (1.0-colorAlpha)*p.G + colorAlpha*cg
	p.B = (1.0-colorAlpha)*p.B + colorAlpha*cb

	if c.sink != nil && save {
		c.sink.Hit(p)
	}

	return p, nil
}

func (c *Conversion) Iterate(p domain.Point, iterations int) error {
	if iterations <= 0 {
		return ErrGenerate
	}

	var err error

	p.R = c.rnd.Float64()
	p.G = c.rnd.Float64()
	p.B = c.rnd.Float64()

	warmup := iterations / 10
	if warmup > 20_000 {
		warmup = 20_000
	}
	for i := 0; i < warmup; i++ {
		if p, err = c.Step(false, p); err != nil {
			return err
		}
	}

	progressEvery := iterations / 10
	if progressEvery == 0 {
		progressEvery = iterations
	}

	for i := 0; i < iterations; i++ {
		if p, err = c.Step(true, p); err != nil {
			return err
		}

		if c.logger != nil && (i+1)%progressEvery == 0 {
			c.logger.Info("iteration progress", "completed", i+1, "total", iterations, "worker", c.worker)
		}
	}

	return nil
}
