package infrastructure

import (
	"math"
	"sync"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/ports"
)

type BoundsAccumulator struct {
	minX, maxX  float64
	minY, maxY  float64
	initialized bool

	n     int64
	sumX  float64
	sumY  float64
	sumXX float64
	sumYY float64

	sigmaK float64

	mu sync.Mutex
}

func NewBoundsAccumulator() *BoundsAccumulator {
	return &BoundsAccumulator{
		minX:   math.Inf(1),
		maxX:   math.Inf(-1),
		minY:   math.Inf(1),
		maxY:   math.Inf(-1),
		sigmaK: 3.0,
	}
}

func (b *BoundsAccumulator) CloneEmpty() ports.SampleSink {
	return &BoundsAccumulator{
		minX:   math.Inf(1),
		maxX:   math.Inf(-1),
		minY:   math.Inf(1),
		maxY:   math.Inf(-1),
		sigmaK: b.sigmaK,
	}
}

func (b *BoundsAccumulator) MergeFrom(other ports.SampleSink) {
	o, ok := other.(*BoundsAccumulator)
	if !ok || o == nil {
		return
	}

	b.mu.Lock()
	o.mu.Lock()
	defer b.mu.Unlock()
	defer o.mu.Unlock()

	b.n += o.n
	b.sumX += o.sumX
	b.sumY += o.sumY
	b.sumXX += o.sumXX
	b.sumYY += o.sumYY

	if !o.initialized {
		return
	}

	if !b.initialized {
		b.minX, b.maxX = o.minX, o.maxX
		b.minY, b.maxY = o.minY, o.maxY
		b.initialized = true
		return
	}

	if o.minX < b.minX {
		b.minX = o.minX
	}
	if o.maxX > b.maxX {
		b.maxX = o.maxX
	}
	if o.minY < b.minY {
		b.minY = o.minY
	}
	if o.maxY > b.maxY {
		b.maxY = o.maxY
	}
}

func (b *BoundsAccumulator) Hit(p domain.Point) {
	x := p.X
	y := p.Y

	b.mu.Lock()
	defer b.mu.Unlock()

	b.n++
	b.sumX += x
	b.sumY += y
	b.sumXX += x * x
	b.sumYY += y * y

	if !b.initialized {
		b.minX, b.maxX = x, x
		b.minY, b.maxY = y, y
		b.initialized = true
		return
	}

	if x < b.minX {
		b.minX = x
	}
	if x > b.maxX {
		b.maxX = x
	}
	if y < b.minY {
		b.minY = y
	}
	if y > b.maxY {
		b.maxY = y
	}
}

func (b *BoundsAccumulator) Bounds() (minX, maxX, minY, maxY float64, ok bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.initialized || b.n == 0 {
		return 0, 0, 0, 0, false
	}

	meanX := b.sumX / float64(b.n)
	meanY := b.sumY / float64(b.n)

	varX := b.sumXX/float64(b.n) - meanX*meanX
	varY := b.sumYY/float64(b.n) - meanY*meanY
	if varX < 0 {
		varX = 0
	}
	if varY < 0 {
		varY = 0
	}

	stdX := math.Sqrt(varX)
	stdY := math.Sqrt(varY)

	if stdX == 0 || stdY == 0 {
		return b.minX, b.maxX, b.minY, b.maxY, true
	}

	loX := meanX - b.sigmaK*stdX
	hiX := meanX + b.sigmaK*stdX
	loY := meanY - b.sigmaK*stdY
	hiY := meanY + b.sigmaK*stdY

	minX = math.Max(b.minX, loX)
	maxX = math.Min(b.maxX, hiX)
	minY = math.Max(b.minY, loY)
	maxY = math.Min(b.maxY, hiY)

	if minX >= maxX || minY >= maxY {
		return b.minX, b.maxX, b.minY, b.maxY, true
	}

	return minX, maxX, minY, maxY, true
}

func AdjustBoundsToAspect(
	minX, maxX, minY, maxY float64,
	width, height int,
	padding float64,
) (float64, float64, float64, float64) {
	w := maxX - minX
	h := maxY - minY
	if w <= 0 || h <= 0 {
		return minX, maxX, minY, maxY
	}

	imgAspect := float64(width) / float64(height)
	boundsAspect := w / h

	if boundsAspect > imgAspect {
		newH := w / imgAspect
		centerY := (minY + maxY) / 2
		minY = centerY - newH/2
		maxY = centerY + newH/2
	} else {
		newW := h * imgAspect
		centerX := (minX + maxX) / 2
		minX = centerX - newW/2
		maxX = centerX + newW/2
	}

	w = maxX - minX
	h = maxY - minY
	minX -= w * padding
	maxX += w * padding
	minY -= h * padding
	maxY += h * padding

	return minX, maxX, minY, maxY
}
