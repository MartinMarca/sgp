package authz

import "errors"

// Errores de autorización (evitan ciclo de imports con services).
var (
	ErrForbidden = errors.New("forbidden")
	ErrNotFound  = errors.New("not found")
)
