package cli

import (
	"errors"
)

var (
	ErrValidation = errors.New("validation error")
	ErrProgram    = errors.New("program error")
)
