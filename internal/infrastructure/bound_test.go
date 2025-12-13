package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

func TestBoundsAccumulator(t *testing.T) {
	t.Parallel()

	acc := NewBoundsAccumulator()
	points := []domain.Point{{X: -1, Y: 0}, {X: 2, Y: 3}, {X: 0, Y: -2}}
	for _, p := range points {
		acc.Hit(p)
	}

	minX, maxX, minY, maxY, ok := acc.Bounds()
	require.True(t, ok)
	assert.Greater(t, maxX, minX)
	assert.Greater(t, maxY, minY)
}

func TestAdjustBoundsToAspect(t *testing.T) {
	t.Parallel()

	minX, maxX, minY, maxY := AdjustBoundsToAspect(-1, 1, -1, 1, 200, 100, 0)
	assert.InDelta(t, 4.0, maxX-minX, 1e-9)
	assert.InDelta(t, 2.0, maxY-minY, 1e-9)
}
