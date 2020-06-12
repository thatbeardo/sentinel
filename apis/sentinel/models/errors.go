package errors

import "errors"

var (
	// ErrDatabase represents a database error during one of the CRUD operations
	ErrDatabase = errors.New("Database Error")

	// ErrNotFound represents that a particular element could not found
	ErrNotFound = errors.New("Data not found")

	// ErrNoContent represents a condition when a particular document isn't found
	ErrNoContent = errors.New("Document not found")

	// ErrPathNotFound represents a condition when a particular route isn't found
	ErrPathNotFound = errors.New("Path not found")
)
