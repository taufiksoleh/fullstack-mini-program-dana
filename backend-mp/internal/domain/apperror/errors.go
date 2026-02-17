package apperror

import "errors"

var ErrNotFound   = errors.New("not found")
var ErrConflict   = errors.New("already exists")
var ErrValidation = errors.New("validation error")
