package flag_parse

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

// Opts - the struct with all params to work with fractal flame
type Opts struct {
	AffineParams   AffineParams
	Functions      FuncsParams
	Seed           int
	IterationCount int
	OutputPath     string
	Threads        int
	Config         string
	Gamma          float64
	Symmetry       int
	Randomize      bool

	Width           MyInt16
	Height          MyInt16
	GammaCorrection bool
}

func parseFlags(fs *flag.FlagSet, args []string) (*Opts, error) {
	var opts Opts

	opts.Width = MyInt16(1920)
	opts.Height = MyInt16(1080)

	fs.Var(&opts.Width, "width", "width of image")
	fs.Var(&opts.Width, "w", "width of image")

	fs.Var(&opts.Height, "height", "height of image")
	fs.Var(&opts.Height, "h", "height of image")

	fs.IntVar(&opts.Seed, "seed", 5, "seed of image")

	fs.IntVar(&opts.IterationCount, "iteration-count", 2500, "iteration for generation")
	fs.IntVar(&opts.IterationCount, "i", 2500, "iteration for generation")

	fs.StringVar(&opts.OutputPath, "output-path", "result.png", "path to result image")
	fs.StringVar(&opts.OutputPath, "o", "result.png", "path to result image")

	fs.IntVar(&opts.Threads, "threads", 1, "threads count")
	fs.IntVar(&opts.Threads, "t", 1, "threads count")

	fs.BoolVar(&opts.GammaCorrection, "gamma_correction", false, "enable gamma correction")
	fs.BoolVar(&opts.GammaCorrection, "g", false, "enable gamma correction")
	fs.Float64Var(&opts.Gamma, "gamma", 2.2, "gamma for brightness correction")

	fs.IntVar(&opts.Symmetry, "symmetry_level", 1, "symmetry level")
	fs.IntVar(&opts.Symmetry, "s", 1, "symmetry level")

	fs.BoolVar(&opts.Randomize, "randomize", false, "enable random generation of affine params and functions")
	fs.BoolVar(&opts.Randomize, "r", false, "enable random generation of affine params and functions")

	fs.Var(&opts.AffineParams, "affine_params",
		"affine params in format <a_1>,<b_1>,<c_1>,<d_1>,<e_1>,<f_1>/<a_N>,<b_N>,<c_N>,<d_N>,<e_N>,<f_N>")
	fs.Var(&opts.AffineParams, "ap",
		"affine params in format <a_1>,<b_1>,<c_1>,<d_1>,<e_1>,<f_1>/<a_N>,<b_N>,<c_N>,<d_N>,<e_N>,<f_N>")

	fs.Var(&opts.Functions, "functions",
		"functions config in format <func_N>:<weigh_func>,<func_N>:<weigh_func>,...")
	fs.Var(&opts.Functions, "f",
		"functions config in format <func_N>:<weigh_func>,<func_N>:<weigh_func>,...")

	fs.StringVar(&opts.Config, "config", "", "optional path to config file")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	visited := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})

	if opts.Config != "" {
		cfg, err := loadConfig(opts.Config)
		if err != nil {
			return nil, err
		}

		if err := applyConfig(&opts, cfg, func(names ...string) bool {
			for _, n := range names {
				if visited[n] {
					return true
				}
			}
			return false
		}); err != nil {
			return nil, err
		}
	}

	if opts.Symmetry < 1 {
		opts.Symmetry = 1
	}
	if opts.Width <= 0 {
		return nil, fmt.Errorf("width must be > 0, got %d", opts.Width)
	}
	if opts.Height <= 0 {
		return nil, fmt.Errorf("height must be > 0, got %d", opts.Height)
	}
	if opts.IterationCount <= 0 {
		return nil, fmt.Errorf("iteration-count must be > 0, got %d", opts.IterationCount)
	}
	if opts.Threads <= 0 {
		return nil, fmt.Errorf("threads must be > 0, got %d", opts.Threads)
	}
	if strings.TrimSpace(opts.OutputPath) == "" {
		return nil, fmt.Errorf("output-path must not be empty")
	}

	if len(opts.AffineParams) == 0 || len(opts.Functions) == 0 {
		opts.Randomize = true
	}

	return &opts, nil
}

func ParseFlags() (*Opts, error) {
	fs := flag.NewFlagSet("fractal flame", flag.ContinueOnError)
	return parseFlags(fs, os.Args[1:])
}

type jsonConfig struct {
	Size struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"size"`
	IterationCount  int             `json:"iteration_count"`
	OutputPath      string          `json:"output_path"`
	Threads         int             `json:"threads"`
	Seed            json.Number     `json:"seed"`
	Functions       []FuncConfig    `json:"functions"`
	AffineParams    []domain.Affine `json:"affine_params"`
	Gamma           float64         `json:"gamma"`
	GammaCorrection *bool           `json:"gamma_correction"`
	Symmetry        int             `json:"symmetry"`
	Randomize       *bool           `json:"randomize"`
}

func loadConfig(path string) (*jsonConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config %q: %w", path, err)
	}

	var cfg jsonConfig
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.UseNumber()
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config %q: %w", path, err)
	}

	return &cfg, nil
}

func applyConfig(opts *Opts, cfg *jsonConfig, wasSet func(...string) bool) error {
	if cfg == nil {
		return nil
	}

	if cfg.Size.Width > 0 && !wasSet("width", "w") {
		opts.Width = MyInt16(cfg.Size.Width)
	}
	if cfg.Size.Height > 0 && !wasSet("height", "h") {
		opts.Height = MyInt16(cfg.Size.Height)
	}
	if cfg.IterationCount > 0 && !wasSet("iteration-count", "i") {
		opts.IterationCount = cfg.IterationCount
	}
	if cfg.OutputPath != "" && !wasSet("output-path", "o") {
		opts.OutputPath = cfg.OutputPath
	}
	if cfg.Threads > 0 && !wasSet("threads", "t") {
		opts.Threads = cfg.Threads
	}
	if cfg.Symmetry > 0 && !wasSet("symmetry_level", "s") {
		opts.Symmetry = cfg.Symmetry
	}
	if cfg.Randomize != nil && !wasSet("randomize", "r") {
		opts.Randomize = *cfg.Randomize
	}
	if cfg.Gamma > 0 && !wasSet("gamma") {
		opts.Gamma = cfg.Gamma
	}
	if cfg.GammaCorrection != nil && !wasSet("gamma_correction", "g") {
		opts.GammaCorrection = *cfg.GammaCorrection
	}
	if cfg.Seed != "" && !wasSet("seed") {
		if i, err := cfg.Seed.Int64(); err == nil {
			opts.Seed = int(i)
		} else if f, err := cfg.Seed.Float64(); err == nil {
			opts.Seed = int(f)
		} else {
			return fmt.Errorf("invalid seed in config: %w", err)
		}
	}
	if len(cfg.Functions) > 0 && !wasSet("functions", "f") {
		opts.Functions = cfg.Functions
	}
	if len(cfg.AffineParams) > 0 && !wasSet("affine_params", "ap") {
		opts.AffineParams = cfg.AffineParams
	}

	return nil
}
