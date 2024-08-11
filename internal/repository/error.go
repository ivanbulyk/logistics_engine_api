package repository

import "errors"

// ErrNotFound is returned when a requested metrics report
// record is not found
var ErrNotFound = errors.New("metrics report not found")

// ErrAlreadyExists is returned when a requested metrics report
// record already exists
var ErrAlreadyExists = errors.New("metrics report already exists")
