package domain

// Affine - struct with coefs
type Affine struct {
	A, B, C, D, E, F       float64
	ColorR, ColorG, ColorB float64
}

func (a *Affine) String() {
	println("a: ", a.A, ",b: ", a.B, " c: ", a.C, ",d: ", a.D, ",e: ", a.E, ",f: ", a.F)
	println("R: ", a.ColorR, ",B: ", a.ColorB, ",G: ", a.ColorG)
}
