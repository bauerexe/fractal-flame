package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

func TestAffineRepository(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		affine  domain.Affine
		point   domain.Point
		want    domain.Point
		wantRGB [3]float64
	}{
		{
			name:    "translation",
			affine:  domain.Affine{A: 1, B: 0, C: 2, D: 0, E: 1, F: -1, ColorR: 0.2, ColorG: 0.3, ColorB: 0.4},
			point:   domain.Point{X: 3, Y: 4},
			want:    domain.Point{X: 5, Y: 3},
			wantRGB: [3]float64{0.2, 0.3, 0.4},
		},
		{
			name:    "scaling and shear",
			affine:  domain.Affine{A: 0.5, B: 0.2, C: -1, D: 0.1, E: 0.5, F: 2, ColorR: 0.1, ColorG: 0.9, ColorB: 0.6},
			point:   domain.Point{X: 4, Y: -2},
			want:    domain.Point{X: 0.6, Y: 1.4},
			wantRGB: [3]float64{0.1, 0.9, 0.6},
		},
		{
			name:    "identity",
			affine:  domain.Affine{A: 1, B: 0, C: 0, D: 0, E: 1, F: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5},
			point:   domain.Point{X: -2, Y: 3},
			want:    domain.Point{X: -2, Y: 3},
			wantRGB: [3]float64{0.5, 0.5, 0.5},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewAffineRepository()
			repo.AddAffine(tt.affine)

			require.Equal(t, 1, repo.Len())

			got, err := repo.Apply(0, tt.point)
			require.NoError(t, err)
			assert.InDelta(t, tt.want.X, got.X, 1e-9)
			assert.InDelta(t, tt.want.Y, got.Y, 1e-9)

			r, g, b := repo.ColorAt(0)
			assert.InDelta(t, tt.wantRGB[0], r, 1e-9)
			assert.InDelta(t, tt.wantRGB[1], g, 1e-9)
			assert.InDelta(t, tt.wantRGB[2], b, 1e-9)
		})
	}

	t.Run("invalid id returns error", func(t *testing.T) {
		t.Parallel()

		repo := NewAffineRepository()

		_, err := repo.Apply(len(tests)+1, domain.Point{})
		assert.ErrorIs(t, err, ErrInvalidId)
	})
}
