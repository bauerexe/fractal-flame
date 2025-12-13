package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

func TestTransformRepository(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		transforms []domain.FuncTransform
		point      domain.Point
		expected   domain.Point
	}{
		{
			name: "averages weighted transforms",
			transforms: []domain.FuncTransform{
				{Variation: func(p *domain.Point) { p.X += 1; p.Y += 1 }, Weight: 1},
				{Variation: func(p *domain.Point) { p.X *= 2; p.Y *= 2 }, Weight: 2},
			},
			point:    domain.Point{X: 1, Y: 2},
			expected: domain.Point{X: 2, Y: 11.0 / 3.0},
		},
		{
			name: "ignores zero or negative weights",
			transforms: []domain.FuncTransform{
				{Variation: func(p *domain.Point) { p.X += 10; p.Y += 10 }, Weight: 0},
				{Variation: func(p *domain.Point) { p.X -= 5; p.Y -= 5 }, Weight: -1},
			},
			point:    domain.Point{X: 3, Y: -2},
			expected: domain.Point{},
		},
		{
			name:       "single transform",
			transforms: []domain.FuncTransform{{Variation: func(p *domain.Point) { p.X *= -1; p.Y *= 0.5 }, Weight: 1}},
			point:      domain.Point{X: 2, Y: 4},
			expected:   domain.Point{X: -2, Y: 2},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewTransformRepository()
			for _, tr := range tt.transforms {
				repo.AddTransform(tr)
			}

			got := repo.Apply(tt.point)
			assert.InDelta(t, tt.expected.X, got.X, 1e-9)
			assert.InDelta(t, tt.expected.Y, got.Y, 1e-9)
		})
	}
}
