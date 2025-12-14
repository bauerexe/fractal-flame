package flag_parse

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

// MyInt16 - type clone int16, to implement interface Value
type MyInt16 int16

func (i *MyInt16) String() string {
	return strconv.Itoa(int(*i))
}

func (i *MyInt16) Set(s string) error {
	val, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid int value %q: %w", s, err)
	}
	*i = MyInt16(val)
	return nil
}

type AffineParams []domain.Affine

func (a *AffineParams) String() string {
	if a == nil || len(*a) == 0 {
		return ""
	}
	parts := make([]string, 0, len(*a))
	for _, tr := range *a {
		parts = append(parts,
			fmt.Sprintf("%g,%g,%g,%g,%g,%g",
				tr.A, tr.B, tr.C, tr.D, tr.E, tr.F,
			),
		)
	}
	return strings.Join(parts, "/")
}

func (a *AffineParams) Set(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("affine params string is empty")
	}

	blocks := strings.Split(s, "/")
	result := make([]domain.Affine, 0, len(blocks))

	for _, block := range blocks {
		block = strings.TrimSpace(block)

		parts := strings.Split(block, ",")
		if len(parts) != 6 {
			return fmt.Errorf("invalid affine block %q: expected 6 values, got %d", block, len(parts))
		}

		vals := make([]float64, 6)
		for i, p := range parts {
			p = strings.TrimSpace(p)
			f, err := strconv.ParseFloat(p, 64)
			if err != nil {
				return fmt.Errorf("invalid float %q in affine block %q: %w", p, block, err)
			}
			vals[i] = f
		}

		result = append(result, domain.Affine{
			A: vals[0],
			B: vals[1],
			C: vals[2],
			D: vals[3],
			E: vals[4],
			F: vals[5],
		})
	}

	if len(result) == 0 {
		return fmt.Errorf("no affine transforms parsed from %q", s)
	}

	*a = result
	return nil
}

type FuncConfig struct {
	Name   string
	Weight float64
}

type FuncsParams []FuncConfig

func (f *FuncsParams) String() string {
	if f == nil || len(*f) == 0 {
		return ""
	}
	parts := make([]string, 0, len(*f))
	for _, fn := range *f {
		parts = append(parts, fmt.Sprintf("%s:%g", fn.Name, fn.Weight))
	}
	return strings.Join(parts, ",")
}

func (f *FuncsParams) Set(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("functions string is empty")
	}

	chunks := strings.Split(s, ",")
	result := make([]FuncConfig, 0, len(chunks))

	for _, ch := range chunks {
		ch = strings.TrimSpace(ch)
		if ch == "" {
			continue
		}

		parts := strings.SplitN(ch, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid function spec %q: expected name:weight", ch)
		}

		name := strings.TrimSpace(parts[0])
		if name == "" {
			return fmt.Errorf("empty function name in %q", ch)
		}

		wStr := strings.TrimSpace(parts[1])
		w, err := strconv.ParseFloat(wStr, 64)
		if err != nil {
			return fmt.Errorf("invalid weight %q for function %q: %w", wStr, name, err)
		}

		result = append(result, FuncConfig{
			Name:   name,
			Weight: w,
		})
	}

	if len(result) == 0 {
		return fmt.Errorf("no functions parsed from %q", s)
	}

	*f = result
	return nil
}
