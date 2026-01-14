package application

import (
	"fmt"
	"math/rand"
	"testing"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/ports"
)

type benchSink struct{ hits int }

func (s *benchSink) Bounds() (minX float64, maxX float64, minY float64, maxY float64, ok bool) {
	panic("implement me")
}

func (s *benchSink) Hit(p domain.Point)           { s.hits++ }
func (s *benchSink) CloneEmpty() ports.SampleSink { return &benchSink{} }
func (s *benchSink) MergeFrom(other ports.SampleSink) {
	if o, ok := other.(*benchSink); ok {
		s.hits += o.hits
	}
}

type benchAffineRepo struct{}

func (r benchAffineRepo) AddAffine(affine domain.Affine) {
	panic("implement me")
}

func (benchAffineRepo) Apply(i int, p domain.Point) (domain.Point, error) { return p, nil }
func (benchAffineRepo) ColorAt(i int) (r, g, b float64)                   { return 0, 0, 0 }
func (benchAffineRepo) Len() int                                          { return 1 }

type benchTransformRepo struct{}

func (r benchTransformRepo) AddTransform(transform domain.FuncTransform) {
	panic("implement me")
}

func (benchTransformRepo) Apply(p domain.Point) domain.Point { return p }

func BenchmarkConversionIterateSingleThread(b *testing.B) {
	sink := &benchSink{}
	conv := NewConversion(benchAffineRepo{}, benchTransformRepo{}, rand.New(rand.NewSource(1)), sink, 1)

	for i := 0; i < b.N; i++ {
		_ = conv.Iterate(domain.Point{}, 1000)
	}
}

func BenchmarkConversionIterateThreads(b *testing.B) {
	affRepo := infrastructure.NewAffineRepository()
	affRepo.AddAffine(domain.Affine{A: 0.5, E: 0.5, ColorR: 0.2, ColorG: 0.5, ColorB: 0.7})

	trRepo := infrastructure.NewTransformRepository()
	trRepo.AddTransform(domain.FuncTransform{Variation: VariationByName("linear"), Weight: 1.0})

	threadCases := []int{1, 2, 4, 8}

	type testCase struct {
		name       string
		iterations int
	}
	testCases := []testCase{
		{
			name:       "4_000 iterations",
			iterations: 4_000,
		},
		{
			name:       "16_000 iterations",
			iterations: 16_000,
		},
		{
			name:       "32_000 iterations",
			iterations: 32_000,
		},
		{
			name:       "100_000 iterations",
			iterations: 100_000,
		},
		{
			name:       "1_000_000 iterations",
			iterations: 1_000_000,
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for _, threads := range threadCases {
				b.Run(fmt.Sprintf("threads-%d", threads), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						sink := &benchSink{}
						_ = iterateWithThreads(affRepo, trRepo, sink, int64(i+1), threads, tc.iterations, nil)
					}
				})
			}
		})
	}
}
