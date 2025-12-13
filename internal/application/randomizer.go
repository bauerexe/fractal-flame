package application

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/flag_parse"
)

type RandomizeOptions struct {
	Count    int
	Override bool
}

func RandomizeAffineParamsWithOptions(opts *flag_parse.Opts, rng *rand.Rand, options RandomizeOptions) {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if len(opts.AffineParams) == 0 && options.Count > 0 {
		opts.AffineParams = make([]domain.Affine, options.Count)
	}

	for i := range opts.AffineParams {
		if options.Override || (opts.AffineParams[i] == (domain.Affine{})) {
			opts.AffineParams[i] = randomGoodAffine(rng)
		}
	}

	randomizeFunctionsSmart(opts, rng, options.Override)
}

func randomGoodAffine(rng *rand.Rand) domain.Affine {
	const (
		minScale   = 0.3
		maxScale   = 0.9
		minShift   = -1.0
		maxShift   = 1.0
		maxRetries = 100
	)

	for attempt := 0; attempt < maxRetries; attempt++ {
		scale := minScale + rng.Float64()*(maxScale-minScale)
		angle := rng.Float64() * 2 * math.Pi

		cosA := math.Cos(angle)
		sinA := math.Sin(angle)

		a := scale * cosA
		b := -scale * sinA
		d := scale * sinA
		e := scale * cosA

		c := minShift + (maxShift-minShift)*rng.Float64()
		f := minShift + (maxShift-minShift)*rng.Float64()

		if isGoodLinearPart(a, b, d, e) {
			return domain.Affine{A: a, B: b, C: c, D: d, E: e, F: f}
		}
	}

	return domain.Affine{A: 0.5, E: 0.5}
}

func randomizeFunctionsSmart(opts *flag_parse.Opts, rng *rand.Rand, override bool) {
	allVariations := []string{"linear", "swirl", "horseshoe", "spherical", "sinusoidal", "popcorn", "pdj", "blob", "cosine"}

	if override || len(opts.Functions) == 0 {
		n := 1 + rng.Intn(4)
		selected := allVariations[:n]

		fns := make([]flag_parse.FuncConfig, 0, n)
		var sumW float64

		for _, name := range selected {
			var w float64
			switch name {
			case "pdj", "cosine":
				w = 0.1 + 0.3*rng.Float64()
			default:
				w = 0.5 + 0.5*rng.Float64()
			}

			sumW += w
			fns = append(fns, flag_parse.FuncConfig{Name: name, Weight: w})
		}

		if sumW > 0 {
			for i := range fns {
				fns[i].Weight /= sumW
			}
		}

		opts.Functions = fns
		return
	}

	var sumW float64
	for i := range opts.Functions {
		name := strings.ToLower(strings.TrimSpace(opts.Functions[i].Name))
		var w float64

		switch name {
		case "pdj", "cosine":
			w = 0.1 + 0.3*rng.Float64()
		default:
			w = 0.5 + 0.5*rng.Float64()
		}

		opts.Functions[i].Weight = w
		sumW += w
	}

	if sumW > 0 {
		for i := range opts.Functions {
			opts.Functions[i].Weight /= sumW
		}
	}
}

func isGoodLinearPart(a, b, d, e float64) bool {
	norm := math.Sqrt(a*a + b*b + d*d + e*e)
	return !(norm < 0.25 || norm > 0.95)
}
