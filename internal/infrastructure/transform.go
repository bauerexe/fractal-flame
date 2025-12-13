package infrastructure

import "gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"

type TransformRepository struct {
	transforms []domain.FuncTransform
	sumWeight  float64
}

func NewTransformRepository() *TransformRepository {
	return &TransformRepository{
		transforms: make([]domain.FuncTransform, 0),
	}
}

func (t *TransformRepository) AddTransform(transform domain.FuncTransform) {
	t.transforms = append(t.transforms, transform)
	if transform.Weight > 0 {
		t.sumWeight += transform.Weight
	}
}

func (t *TransformRepository) Apply(p domain.Point) domain.Point {
	if t.sumWeight == 0 {
		return domain.Point{}
	}

	var accX, accY float64

	for _, tr := range t.transforms {
		if tr.Weight <= 0 {
			continue
		}

		tmp := p
		tr.Variation(&tmp)

		accX += tr.Weight * tmp.X
		accY += tr.Weight * tmp.Y
	}

	inv := 1.0 / t.sumWeight
	p.X = accX * inv
	p.Y = accY * inv
	return p
}
