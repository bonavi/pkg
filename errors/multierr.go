package errors

import "sync"

type MultiError struct {
	mu     sync.Mutex
	errors []error
}

func NewMultiError() *MultiError {
	return &MultiError{
		errors: make([]error, 0),
		mu:     sync.Mutex{},
	}
}

func (e *MultiError) Append(errs ...error) *MultiError {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.errors = append(e.errors, errs...)

	return e
}

func (e *MultiError) Get() []error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.errors
}
