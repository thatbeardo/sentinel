package errors

import "errors"

var (
	// ErrDatabase represents a database error during one of the CRUD operations
	ErrDatabase = errors.New("Database Error")

	// ErrNoContent represents a condition when a particular document isn't found
	ErrNoContent = errors.New("Error: Document not found")
)
