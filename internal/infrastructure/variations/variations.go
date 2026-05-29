package variations

import (
	"math"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

func Linear(p *domain.Point) {
	x := p.X
	y := p.Y

	p.X, p.Y = x, y
}

func Swirl(p *domain.Point) {
	r2 := p.X*p.X + p.Y*p.Y
	s := math.Sin(r2)
	c := math.Cos(r2)

	x := p.X*s - p.Y*c
	y := p.X*c + p.Y*s

	p.X, p.Y = x, y
}

func Horseshoe(p *domain.Point) {
	r := math.Hypot(p.X, p.Y)
	if r == 0 {
		return
	}

	x := (p.X*p.X - p.Y*p.Y) / r
	y := 2 * p.X * p.Y / r

	p.X, p.Y = x, y
}

func Sinusoidal(p *domain.Point) {
	p.X = math.Sin(p.X)
	p.Y = math.Sin(p.Y)
}

func Spherical(p *domain.Point) {
	r2 := p.X*p.X + p.Y*p.Y
	if r2 == 0 {
		return
	}
	k := 1.0 / r2

	p.X *= k
	p.Y *= k
}

const (
	popcornC = 0.02
	popcornF = 0.1
)

func Popcorn(p *domain.Point) {
	x := p.X + popcornC*math.Sin(math.Tan(3*p.Y))
	y := p.Y + popcornF*math.Sin(math.Tan(3*p.X))

	p.X, p.Y = x, y
}

const (
	pdjP1    = 1.0
	pdjP2    = 1.0
	pdjP3    = 1.0
	pdjP4    = 1.0
	pdjScale = 1.2
)

func PDJ(p *domain.Point) {
	x := math.Sin(pdjP1*p.Y) - math.Cos(pdjP2*p.X)
	y := math.Sin(pdjP3*p.X) - math.Cos(pdjP4*p.Y)

	p.X = pdjScale * x
	p.Y = pdjScale * y
}

const (
	blobLow   = 0.2
	blobHigh  = 1.2
	blobWaves = 3.0
)

func Blob(p *domain.Point) {
	if p.X == 0 && p.Y == 0 {
		return
	}

	r := math.Hypot(p.X, p.Y)
	theta := math.Atan2(p.Y, p.X)

	newR := r * (blobLow + (blobHigh-blobLow)*(math.Sin(blobWaves*theta)+1)/2)

	p.X = newR * math.Cos(theta)
	p.Y = newR * math.Sin(theta)
}

func Cosine(p *domain.Point) {
	x := p.X
	y := p.Y

	p.X = math.Cos(math.Pi*x) * math.Cosh(y)
	p.Y = -math.Sin(math.Pi*x) * math.Sinh(y)
}
