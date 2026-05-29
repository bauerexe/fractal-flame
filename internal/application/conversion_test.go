package application

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/ports"
)

type stubSink struct {
	hits  int
	store []domain.Point
}

func (s *stubSink) Bounds() (minX, maxX, minY, maxY float64, ok bool) {
	return 0, 0, 0, 0, false
}

func (s *stubSink) Hit(p domain.Point) {
	s.hits++
	s.store = append(s.store, p)
}

func (s *stubSink) CloneEmpty() ports.SampleSink {
	return &stubSink{}
}

func (s *stubSink) MergeFrom(other ports.SampleSink) {
	o, ok := other.(*stubSink)
	if !ok || o == nil {
		return
	}
	s.hits += o.hits
	s.store = append(s.store, o.store...)
}

type stubAffineRepo struct{}

func (stubAffineRepo) AddAffine(domain.Affine) {
}

func (stubAffineRepo) Apply(i int, p domain.Point) (domain.Point, error) { return p, nil }
func (stubAffineRepo) ColorAt(i int) (r, g, b float64)                   { return 0.1, 0.2, 0.3 }
func (stubAffineRepo) Len() int                                          { return 1 }

type stubTransformRepo struct{}

func (stubTransformRepo) AddTransform(domain.FuncTransform) {
}

func (stubTransformRepo) Apply(p domain.Point) domain.Point { return p }

func TestConversionIterate(t *testing.T) {
	t.Parallel()

	sink := &stubSink{}
	rnd := rand.New(rand.NewSource(1))

	conv := NewConversion(stubAffineRepo{}, stubTransformRepo{}, rnd, sink, 1)

	err := conv.Iterate(domain.Point{}, 10)
	require.NoError(t, err)
	assert.Equal(t, 10, sink.hits)
}

func TestConversionStepValidation(t *testing.T) {
	t.Parallel()

	rnd := rand.New(rand.NewSource(1))

	conv := NewConversion(stubAffineRepo{}, stubTransformRepo{}, rnd, nil, 1)
	_, err := conv.Step(false, domain.Point{})
	require.NoError(t, err)

	convEmpty := NewConversion(stubAffineRepo{}, stubTransformRepo{}, rnd, nil, 1)
	convEmpty.affRepo = stubAffineRepoLenZero{}

	_, err = convEmpty.Step(true, domain.Point{})
	require.ErrorContains(t, err, domain.ErrGenerate.Error())
}

type stubAffineRepoLenZero struct{}

func (stubAffineRepoLenZero) AddAffine(domain.Affine) {
}

func (stubAffineRepoLenZero) Apply(i int, p domain.Point) (domain.Point, error) {
	return p, domain.ErrGenerate
}

func (stubAffineRepoLenZero) ColorAt(i int) (r, g, b float64) { return 0, 0, 0 }
func (stubAffineRepoLenZero) Len() int                        { return 0 }
