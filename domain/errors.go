package domain

import "errors"

// Domain layer errors
var (
	ErrValidation          = errors.New("invalid import status")
	ErrImportAlreadyExists = errors.New("import ID already exists")
	ErrUserAlreadyExists   = errors.New("import ID already exists")
	ErrInvalidInput        = errors.New("invalid input")
)
