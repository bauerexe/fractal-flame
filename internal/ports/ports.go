package ports

import "gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"

type AffineRepository interface {
	AddAffine(affine domain.Affine)
	Apply(i int, p domain.Point) (domain.Point, error)
	ColorAt(i int) (r, g, b float64)
	Len() int
}

type TransformRepository interface {
	AddTransform(transform domain.FuncTransform)
	Apply(p domain.Point) domain.Point
}

type SampleSink interface {
	Hit(p domain.Point)
	CloneEmpty() SampleSink
	MergeFrom(other SampleSink)
	Bounds() (minX float64, maxX float64, minY float64, maxY float64, ok bool)
}
