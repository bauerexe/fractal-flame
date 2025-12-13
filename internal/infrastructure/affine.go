package infrastructure

import (
	"errors"
	"sync"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

var ErrInvalidId = errors.New("invalid id")

type AffineRepository struct {
	mu  sync.RWMutex
	aff []domain.Affine
	len int
}

func NewAffineRepository() *AffineRepository {
	return &AffineRepository{
		aff: make([]domain.Affine, 0),
		len: 0,
	}
}

func (a *AffineRepository) AddAffine(affine domain.Affine) {
	//a.mu.Lock()
	//defer a.mu.Unlock()

	a.aff = append(a.aff, affine)
	a.len++
}

func (a *AffineRepository) Apply(i int, p domain.Point) (domain.Point, error) {
	if i >= a.len || a.len < 0 {
		return domain.Point{}, ErrInvalidId
	}

	aff := a.aff[i]
	x := aff.A*p.X + aff.B*p.Y + aff.C
	y := aff.D*p.X + aff.E*p.Y + aff.F
	p.X = x
	p.Y = y

	return p, nil
}

func (a *AffineRepository) Len() int {
	//a.mu.RLock()
	//defer a.mu.RUnlock()

	return a.len
}

func (a *AffineRepository) ColorAt(i int) (r, g, b float64) {
	//a.mu.RLock()
	//defer a.mu.RUnlock()

	if i < 0 || i >= a.len {
		return 0, 0, 0
	}
	aff := a.aff[i]
	return aff.ColorR, aff.ColorG, aff.ColorB
}
