package repository

import "errors"

// ErrNotFound é retornado quando um recurso não existe.
// Os handlers mapeiam esse erro para HTTP 404.
var ErrNotFound = errors.New("recurso não encontrado")

// ErrInvalidID é retornado quando um id não está no formato esperado.
var ErrInvalidID = errors.New("id inválido")
