package flag_parse

import (
	"flag"
	"os"
	"testing"
)

func TestParseFlagsEnablesRandomizeWhenMissingParams(t *testing.T) {
	t.Parallel()

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	opts, err := parseFlags(fs, []string{"-w", "10", "-h", "10", "-i", "100"})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if !opts.Randomize {
		t.Fatalf("expected randomize to be enabled when params are missing")
	}
}

func TestParseFlagsUsesConfigWhenFlagsAbsent(t *testing.T) {
	t.Parallel()

	cfgPath := writeConfig(t, `{
  "size": {"width": 128, "height": 96},
  "iteration_count": 300,
  "output_path": "cfg.png",
  "threads": 3,
  "seed": 42,
  "functions": [{"name": "linear", "weight": 1.0}],
  "affine_params": [{"a": 1, "b": 0, "c": 0, "d": 0, "e": 1, "f": 0}]
}`)

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	opts, err := parseFlags(fs, []string{"--config", cfgPath})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if opts.Width != 128 || opts.Height != 96 || opts.IterationCount != 300 || opts.OutputPath != "cfg.png" || opts.Threads != 3 || opts.Seed != 42 {
		t.Fatalf("config values not applied: %+v", opts)
	}
	if opts.Randomize {
		t.Fatalf("randomize should remain false when config provides params")
	}
}

func TestParseFlagsPrioritizesCLIOverConfig(t *testing.T) {
	t.Parallel()

	cfgPath := writeConfig(t, `{"size": {"width": 500, "height": 500}, "threads": 4}`)

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	opts, err := parseFlags(fs, []string{"--config", cfgPath, "-w", "64", "-h", "32", "-t", "1"})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if opts.Width != 64 || opts.Height != 32 || opts.Threads != 1 {
		t.Fatalf("cli flags should override config: %+v", opts)
	}
}

func writeConfig(t *testing.T, data string) string {
	t.Helper()

	path := t.TempDir() + "/config.json"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	return path
}
