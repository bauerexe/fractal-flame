package infrastructure

import (
	"image/png"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

func TestImageAccumulator_RenderPNG(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	acc := NewImageAccumulator(fs, 10, 10, -1, 1, -1, 1, true, 2.2, 1)
	acc.Hit(domain.Point{X: 0, Y: 0, R: 1, G: 0.5, B: 0.25})

	err := acc.RenderPNG("test.png")
	require.NoError(t, err)

	info, err := fs.Stat("test.png")
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))

	f, err := fs.Open("test.png")
	require.NoError(t, err)
	defer f.Close()

	_, err = png.Decode(f)
	assert.NoError(t, err)
}

func TestImageAccumulator_Symmetry(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	acc := NewImageAccumulator(fs, 5, 5, -1, 1, -1, 1, false, 0, 4)

	acc.Hit(domain.Point{X: 0.5, Y: 0, R: 1})

	totalHits := 0.0
	for _, h := range acc.hits {
		totalHits += h
	}

	assert.GreaterOrEqual(t, totalHits, 4.0)
}
