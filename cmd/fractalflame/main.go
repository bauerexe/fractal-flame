package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/afero"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/application"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/flag_parse"
)

const help = `fractalflame \
  --width 1920 --height 1080 --seed 10 --iteration-count 2000000 \
  --output-path result.png --threads 1 \
  --affine_params "0.878017,-0.443449,-0.415128,0.768696,0.469698,0.330358/0.734865,0.105944,0.106802,0.384882,0.113686,-0.478264/0.615335,0.223012,0.49482,0.436119,0.629154,0.061848" \
  --functions "linear:1,popcorn:0.2" \
  --s 3 --gamma_correction=true --gamma 2.2
`

func main() {
	usage := func(full bool) {
		_, _ = fmt.Fprintln(os.Stdout, "Usage: fractalflame [flags]")
		if full {
			_, _ = fmt.Fprintln(os.Stdout, "\nExample:\n"+help)
		} else {
			_, _ = fmt.Fprintln(os.Stdout, "Run with -h for full help.")
		}
	}

	opts, err := flag_parse.ParseFlags()
	if errors.Is(err, flag.ErrHelp) {
		usage(true)
		return
	}
	if err != nil {
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})).
			Error("failed to parse flags", "err", err)
		usage(false)
		os.Exit(2)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	if err := (&application.Conversion{}).Start(*opts, afero.NewOsFs(), logger); err != nil {
		logger.Error("failed to generate fractal", "err", err)
		os.Exit(1)
	}
}
