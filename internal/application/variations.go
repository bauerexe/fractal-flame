package application

import (
	"strings"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/variations"
)

func VariationByName(name string) func(p *domain.Point) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "swirl":
		return variations.Swirl
	case "horseshoe":
		return variations.Horseshoe
	case "linear":
		return variations.Linear
	case "spherical":
		return variations.Spherical
	case "popcorn":
		return variations.Popcorn
	case "pdj":
		return variations.PDJ
	case "sinusoidal":
		return variations.Sinusoidal
	case "blob":
		return variations.Blob
	case "cosine":
		return variations.Cosine
	default:
		return variations.Linear
	}
}
