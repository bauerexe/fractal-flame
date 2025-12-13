package infrastructure

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"runtime"
	"sync"

	"github.com/spf13/afero"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/ports"
)

func (acc *ImageAccumulator) CloneEmpty() ports.SampleSink {
	n := acc.width * acc.height
	return &ImageAccumulator{
		fs:              acc.fs,
		width:           acc.width,
		height:          acc.height,
		minX:            acc.minX,
		maxX:            acc.maxX,
		minY:            acc.minY,
		maxY:            acc.maxY,
		scaleX:          acc.scaleX,
		scaleY:          acc.scaleY,
		hits:            make([]float64, n),
		sumR:            make([]float64, n),
		sumG:            make([]float64, n),
		sumB:            make([]float64, n),
		gammaCorrection: acc.gammaCorrection,
		gamma:           acc.gamma,
		symmetry:        acc.symmetry,
		symAngles:       acc.symAngles,
	}
}

func (acc *ImageAccumulator) MergeFrom(other ports.SampleSink) {
	o, ok := other.(*ImageAccumulator)
	if !ok || o == nil {
		return
	}

	if acc.width != o.width || acc.height != o.height {
		return
	}

	n := acc.width * acc.height
	workers := runtime.NumCPU()
	if workers < 1 {
		workers = 1
	}
	step := (n + workers - 1) / workers

	var wg sync.WaitGroup
	for start := 0; start < n; start += step {
		end := start + step
		if end > n {
			end = n
		}

		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()
			for i := s; i < e; i++ {
				acc.hits[i] += o.hits[i]
				acc.sumR[i] += o.sumR[i]
				acc.sumG[i] += o.sumG[i]
				acc.sumB[i] += o.sumB[i]
			}
		}(start, end)
	}

	wg.Wait()
}

type ImageAccumulator struct {
	fs afero.Fs

	width  int
	height int

	minX, maxX float64
	minY, maxY float64

	scaleX float64
	scaleY float64

	hits []float64
	sumR []float64
	sumG []float64
	sumB []float64

	gammaCorrection bool
	gamma           float64
	symmetry        int
	symAngles       []cosSin
}

type cosSin struct {
	c float64
	s float64
}

func NewImageAccumulator(fs afero.Fs, width, height int,
	minX, maxX, minY, maxY float64, gammaCorrection bool,
	gamma float64, symmetry int,
) *ImageAccumulator {
	if gamma == 0 {
		gamma = 2.2
	}

	invRangeX := maxX - minX
	if invRangeX != 0 {
		invRangeX = float64(width) / invRangeX
	}

	invRangeY := maxY - minY
	if invRangeY != 0 {
		invRangeY = float64(height) / invRangeY
	}

	angles := make([]cosSin, 0, symmetry-1)
	if symmetry > 1 {
		angleStep := 2 * math.Pi / float64(symmetry)
		for i := 1; i < symmetry; i++ {
			a := angleStep * float64(i)
			angles = append(angles, cosSin{c: math.Cos(a), s: math.Sin(a)})
		}
	}
	n := width * height
	return &ImageAccumulator{
		fs:              fs,
		width:           width,
		height:          height,
		minX:            minX,
		maxX:            maxX,
		minY:            minY,
		maxY:            maxY,
		scaleX:          invRangeX,
		scaleY:          invRangeY,
		hits:            make([]float64, n),
		sumR:            make([]float64, n),
		sumG:            make([]float64, n),
		sumB:            make([]float64, n),
		gammaCorrection: gammaCorrection,
		gamma:           gamma,
		symmetry:        symmetry,
		symAngles:       angles,
	}
}
func (acc *ImageAccumulator) hitSingle(p domain.Point) {
	ix := int((p.X - acc.minX) * acc.scaleX)
	iy := int((acc.maxY - p.Y) * acc.scaleY)

	if ix < 0 || ix >= acc.width || iy < 0 || iy >= acc.height {
		return
	}

	idx := iy*acc.width + ix
	acc.hits[idx]++
	acc.sumR[idx] += p.R
	acc.sumG[idx] += p.G
	acc.sumB[idx] += p.B
}

func (acc *ImageAccumulator) Hit(p domain.Point) {
	n := acc.symmetry
	if n <= 1 {
		acc.hitSingle(p)
		return
	}

	acc.hitSingle(p)
	for _, ang := range acc.symAngles {
		x := p.X*ang.c - p.Y*ang.s
		y := p.X*ang.s + p.Y*ang.c

		tmp := p
		tmp.X = x
		tmp.Y = y

		acc.hitSingle(tmp)
	}
}

func (acc *ImageAccumulator) RenderPNG(path string) error {

	img := image.NewRGBA(image.Rect(0, 0, acc.width, acc.height))

	var maxHits float64
	for _, h := range acc.hits {
		if h > maxHits {
			maxHits = h
		}
	}
	if maxHits <= 0 {
		maxHits = 1
	}

	const exposure = 8.0

	const noiseCutoff = 0.0001

	clamp01 := func(v float64) float64 {
		if v < 0 {
			return 0
		}
		if v > 1 {
			return 1
		}
		return v
	}

	for y := 0; y < acc.height; y++ {
		for x := 0; x < acc.width; x++ {
			idx := y*acc.width + x
			h := acc.hits[idx]

			if h == 0 {
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
				continue
			}

			norm := h / maxHits

			if norm < noiseCutoff {
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
				continue
			}

			norm = (norm - noiseCutoff) / (1.0 - noiseCutoff)
			if norm < 0 {
				norm = 0
			}

			r := acc.sumR[idx] / h
			g := acc.sumG[idx] / h
			b := acc.sumB[idx] / h

			L := norm * exposure

			L = L / (1.0 + L)

			if acc.gammaCorrection && acc.gamma > 0 {
				L = math.Pow(L, 1.0/acc.gamma)
			}

			L = clamp01(L)

			Rf := clamp01(r * L)
			Gf := clamp01(g * L)
			Bf := clamp01(b * L)

			img.Set(x, y, color.RGBA{
				R: uint8(Rf * 255),
				G: uint8(Gf * 255),
				B: uint8(Bf * 255),
				A: 255,
			})
		}
	}

	if acc.fs == nil {
		acc.fs = afero.NewOsFs()
	}

	f, err := acc.fs.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
