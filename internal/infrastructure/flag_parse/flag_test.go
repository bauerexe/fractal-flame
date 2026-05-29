package flag_parse

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		name          string
		args          []string
		want          *Opts
		wantAffineStr string
		wantFuncStr   string
		wantErr       string
		wantRandomize *bool
	}

	defaultAffine := "1,0,0,1,0,0"
	defaultFuncs := "linear:1"

	tests := []TestCase{
		{
			name: "ok: minimal required flags with long names",
			args: []string{
				"--affine_params", defaultAffine,
				"--functions", defaultFuncs,
			},
			want: &Opts{
				Width:           MyInt16(1920),
				Height:          MyInt16(1080),
				Seed:            5,
				IterationCount:  2500,
				OutputPath:      "result.png",
				Threads:         1,
				Config:          "",
				GammaCorrection: false,
				Gamma:           2.2,
				Symmetry:        1,
			},
			wantAffineStr: defaultAffine,
			wantFuncStr:   defaultFuncs,
		},
		{
			name: "ok: all options with short aliases",
			args: []string{
				"-w", "800",
				"-h", "600",
				"-seed", "42",
				"-i", "1000",
				"-o", "out.png",
				"-t", "4",
				"-g",
				"-gamma", "2.4",
				"-s", "3",
				"-ap", "1,0,0,1,0,0/0.5,0,0,0.5,10,10",
				"-f", "sinus:0.3,swirl:0.7",
				"-config", "config.yaml",
			},
			want: &Opts{
				Width:           MyInt16(800),
				Height:          MyInt16(600),
				Seed:            42,
				IterationCount:  1000,
				OutputPath:      "out.png",
				Threads:         4,
				Config:          "config.yaml",
				GammaCorrection: true,
				Gamma:           2.4,
				Symmetry:        3,
			},
			wantAffineStr: "1,0,0,1,0,0/0.5,0,0,0.5,10,10",
			wantFuncStr:   "sinus:0.3,swirl:0.7",
		},
		{
			name: "ok: symmetry less than 1 is clamped to 1",
			args: []string{
				"--affine_params", defaultAffine,
				"--functions", defaultFuncs,
				"--symmetry_level", "0",
			},
			want: &Opts{
				Width:           MyInt16(1920),
				Height:          MyInt16(1080),
				Seed:            5,
				IterationCount:  2500,
				OutputPath:      "result.png",
				Threads:         1,
				Config:          "",
				GammaCorrection: false,
				Gamma:           2.2,
				Symmetry:        1,
			},
			wantAffineStr: defaultAffine,
			wantFuncStr:   defaultFuncs,
		},

		{
			name: "ok: missing affine params triggers randomize",
			args: []string{
				"--functions", defaultFuncs,
			},
			wantRandomize: ptr(true),
		},
		{
			name: "ok: missing functions triggers randomize",
			args: []string{
				"--affine_params", defaultAffine,
			},
			wantRandomize: ptr(true),
		},
		{
			name: "error: width must be > 0",
			args: []string{
				"--affine_params", defaultAffine,
				"--functions", defaultFuncs,
				"--width", "0",
			},
			wantErr: "width must be > 0",
		},
		{
			name: "error: height must be > 0",
			args: []string{
				"--affine_params", defaultAffine,
				"--functions", defaultFuncs,
				"--height", "0",
			},
			wantErr: "height must be > 0",
		},
		{
			name: "error: iteration-count must be > 0",
			args: []string{
				"--affine_params", defaultAffine,
				"--functions", defaultFuncs,
				"--iteration-count", "0",
			},
			wantErr: "iteration-count must be > 0",
		},
		{
			name: "error: threads must be > 0",
			args: []string{
				"--affine_params", defaultAffine,
				"--functions", defaultFuncs,
				"--threads", "0",
			},
			wantErr: "threads must be > 0",
		},
		{
			name: "error: empty output-path",
			args: []string{
				"--affine_params", defaultAffine,
				"--functions", defaultFuncs,
				"--output-path", "",
			},
			wantErr: "output-path must not be empty",
		},
		{
			name: "error: invalid affine params format (parse error inside AffineParams.Set)",
			args: []string{
				"--affine_params", "invalid_format",
				"--functions", defaultFuncs,
			},

			wantErr: "invalid",
		},
		{
			name: "error: invalid functions format (parse error inside FuncsParams.Set)",
			args: []string{
				"--affine_params", defaultAffine,
				"--functions", "bad_format",
			},
			wantErr: "invalid",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			args := make([]string, len(tc.args))
			copy(args, tc.args)
			for i := 0; i < len(args); i++ {
				if args[i] == "-config" || args[i] == "--config" {
					if i+1 < len(args) {
						cfgPath := writeConfig(t, `{}`)
						args[i+1] = cfgPath
						if tc.want != nil {
							tc.want.Config = cfgPath
						}
					}
				}
			}

			fs := flag.NewFlagSet("fractal flame test", flag.ContinueOnError)
			opts, err := parseFlags(fs, args)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, opts)

			if tc.want != nil {
				assert.Equal(t, tc.want.Width, opts.Width)
				assert.Equal(t, tc.want.Height, opts.Height)
				assert.Equal(t, tc.want.Seed, opts.Seed)
				assert.Equal(t, tc.want.IterationCount, opts.IterationCount)
				assert.Equal(t, tc.want.OutputPath, opts.OutputPath)
				assert.Equal(t, tc.want.Threads, opts.Threads)
				assert.Equal(t, tc.want.Config, opts.Config)
				assert.Equal(t, tc.want.GammaCorrection, opts.GammaCorrection)
				assert.InDelta(t, tc.want.Gamma, opts.Gamma, 1e-9)
				assert.Equal(t, tc.want.Symmetry, opts.Symmetry)
			}

			if tc.wantAffineStr != "" {
				assert.Equal(t, tc.wantAffineStr, opts.AffineParams.String())
			}

			if tc.wantFuncStr != "" {
				assert.Equal(t, tc.wantFuncStr, opts.Functions.String())
			}

			if tc.wantRandomize != nil {
				assert.Equal(t, *tc.wantRandomize, opts.Randomize)
			}
		})
	}

}

func ptr[T any](v T) *T {
	return &v
}
