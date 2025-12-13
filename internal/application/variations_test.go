package application

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/variations"
)

func TestVariationByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		fnName     string
		input      domain.Point
		expectedFn func(p *domain.Point)
	}{
		{
			name:       "defaults to linear",
			fnName:     "unknown",
			input:      domain.Point{X: 1, Y: 2},
			expectedFn: variations.Linear,
		},
		{
			name:       "swirl",
			fnName:     "swirl",
			input:      domain.Point{X: 1, Y: 1},
			expectedFn: variations.Swirl,
		},
		{
			name:       "horseshoe",
			fnName:     "horseshoe",
			input:      domain.Point{X: 1, Y: 2},
			expectedFn: variations.Horseshoe,
		},
		{
			name:       "sinusoidal",
			fnName:     "sinusoidal",
			input:      domain.Point{X: 1, Y: -1},
			expectedFn: variations.Sinusoidal,
		},
		{
			name:       "spherical",
			fnName:     "spherical",
			input:      domain.Point{X: 1, Y: 2},
			expectedFn: variations.Spherical,
		},
		{
			name:       "popcorn",
			fnName:     "popcorn",
			input:      domain.Point{X: 1, Y: 0.5},
			expectedFn: variations.Popcorn,
		},
		{
			name:       "pdj",
			fnName:     "pdj",
			input:      domain.Point{X: 0.3, Y: -0.7},
			expectedFn: variations.PDJ,
		},
		{
			name:       "blob",
			fnName:     "blob",
			input:      domain.Point{X: 0.5, Y: 0.5},
			expectedFn: variations.Blob,
		},
		{
			name:       "cosine",
			fnName:     "cosine",
			input:      domain.Point{X: 0.25, Y: -0.5},
			expectedFn: variations.Cosine,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.input
			VariationByName(tt.fnName)(&got)

			want := tt.input
			tt.expectedFn(&want)

			assert.InDelta(t, want.X, got.X, 1e-9)
			assert.InDelta(t, want.Y, got.Y, 1e-9)
		})
	}
}
