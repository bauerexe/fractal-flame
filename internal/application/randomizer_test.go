package application

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/flag_parse"
)

func TestRandomizeAffineParamsWithOptions(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {
		opts := flag_parse.Opts{}
		rng := rand.New(rand.NewSource(int64(1 + i)))

		RandomizeAffineParamsWithOptions(&opts, rng, RandomizeOptions{Count: 2, Override: true})

		require.Len(t, opts.AffineParams, 2)
		assert.NotZero(t, opts.AffineParams[0].A)
		assert.NotZero(t, opts.AffineParams[1].B)
		assert.NotEmpty(t, opts.Functions)
		for _, fn := range opts.Functions {
			assert.Greater(t, fn.Weight, 0.0)
			assert.LessOrEqual(t, fn.Weight, 1.0)
		}
	}
}

func TestRandomizeFunctionsOverrideWeights(t *testing.T) {
	t.Parallel()

	for i := 0; i < 3; i++ {
		opts := flag_parse.Opts{Functions: flag_parse.FuncsParams{{Name: "pdj", Weight: 10}}}
		rng := rand.New(rand.NewSource(int64(2 + i)))

		RandomizeAffineParamsWithOptions(&opts, rng, RandomizeOptions{Count: 0, Override: false})

		assert.Greater(t, opts.Functions[0].Weight, 0.0)
		assert.LessOrEqual(t, opts.Functions[0].Weight, 1.0)
		assert.Equal(t, "pdj", opts.Functions[0].Name)
	}
}
