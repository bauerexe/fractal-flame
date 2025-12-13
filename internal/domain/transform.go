package domain

type FuncTransform struct {
	Variation func(p *Point)
	Weight    float64
}
