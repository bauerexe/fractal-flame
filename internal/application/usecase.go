package application

import (
	"errors"
)

var (
	ErrGenerate  = errors.New("error generate")
	ErrInvalidId = errors.New("invalid id")
)
