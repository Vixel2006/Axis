package services

import "fmt"

// NotFoundError is returned when a requested resource is not found.
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Message)
}

// ConflictError is returned when a request cannot be completed due to a conflict with the current state of the resource.
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Message)
}

// ForbiddenError is returned when a user is not authorized to perform an action.
type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("forbidden: %s", e.Message)
}
