package application

import (
	"errors"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

var (
	ErrGenerate  = errors.New("error generate")
	ErrInvalidId = errors.New("invalid id")
)

type AffineRepository interface {
	Apply(i int, p domain.Point) (domain.Point, error)
	ColorAt(i int) (r, g, b float64)
	Len() int
}

type TransformRepository interface {
	Apply(p domain.Point) domain.Point
}

type RandomSource interface {
	Intn(n int) int
	Float64() float64
}

type SampleSink interface {
	Hit(p domain.Point)
	CloneEmpty() SampleSink
	MergeFrom(other SampleSink)
}
