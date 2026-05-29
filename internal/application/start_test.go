package application

import (
	"io"
	"log/slog"
	"math/rand"
	"testing"

	"github.com/spf13/afero"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/flag_parse"
)

func TestConversionStart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts flag_parse.Opts
	}{
		{
			name: "randomizes missing config",
			opts: flag_parse.Opts{
				Width:          flag_parse.MyInt16(64),
				Height:         flag_parse.MyInt16(64),
				IterationCount: 200,
				OutputPath:     "out_random.png",
				Seed:           1,
			},
		},
		{
			name: "uses provided config",
			opts: flag_parse.Opts{
				Width:          flag_parse.MyInt16(32),
				Height:         flag_parse.MyInt16(32),
				IterationCount: 150,
				OutputPath:     "out_provided.png",
				Seed:           1,
				AffineParams: flag_parse.AffineParams{
					{A: 1, B: 0, C: 0, D: 0, E: 1, F: 0, ColorR: 0.1, ColorG: 0.2, ColorB: 0.3},
				},
				Functions: flag_parse.FuncsParams{{Name: "linear", Weight: 1}},
				Randomize: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()

			affRepo := infrastructure.NewAffineRepository()
			trRepo := infrastructure.NewTransformRepository()
			boundsAcc := infrastructure.NewBoundsAccumulator()
			boundsRnd := rand.New(rand.NewSource(tt.opts.Seed))

			conv := NewConversion(affRepo, trRepo, boundsRnd, boundsAcc, tt.opts.Seed)

			logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelWarn}))

			if err := conv.Start(tt.opts, fs, logger); err != nil {
				t.Fatalf("start failed: %v", err)
			}

			if _, err := fs.Stat(tt.opts.OutputPath); err != nil {
				t.Fatalf("expected output file %s to be created: %v", tt.opts.OutputPath, err)
			}
		})
	}
}
