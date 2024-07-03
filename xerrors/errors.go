// Package xerrors provides a stack of multiple errors that can be pushed to, popped from, and concatenated.
package xerrors

import (
	"errors"
	"io"
	"slices"
	"sync"
)

type ErrorStack interface {
	Len() int
	Is(error) bool
	As(interface{}) bool
	Errors() []error
	Next() error
	io.Seeker
	error
}

// Errors is a stack of multiple errors that can be pushed to, popped from, or concatenated.
// It is safe for concurrent use, and can return an immutable copy of itself as an [errstack.ErrorsImmutable].
type Errors struct {
	mu        sync.RWMutex
	errs      []error
	i         int
	immutable bool
}

// NewErrors returns a new [Errors] stack.
func NewErrors() *Errors {
	return &Errors{
		errs:      make([]error, 0),
		immutable: false,
	}
}

// Next returns the next error in the stack, incrementing the internal index.
// If we've reached the end of the stack, it returns nil.
//
// Use [errstack.Errors.Seek] to rewind the internal index if needed.
func (e *Errors) Next() error {
	e.mu.RLock()
	err := e.next(&e.i)
	e.mu.RUnlock()
	return err
}

// Push adds an error to the stack.
// It is safe for concurrent use.
func (e *Errors) Push(err error) {
	if e.immutable {
		panic("Add called on immutable error stack")
	}
	e.mu.Lock()
	if len(e.errs) > 0 {
		newStack := []error{err}
		e.errs = append(newStack, e.errs...)
	} else {
		e.errs = append(e.errs, err)
	}
	e.mu.Unlock()
}

// Pop pops one error from the stack, removing it from the stack and returning it.
func (e *Errors) Pop() error {
	if e.immutable {
		panic("Pop called on immutable error stack")
	}
	e.mu.Lock()
	if len(e.errs) == 0 {
		e.mu.Unlock()
		return nil
	}
	err := e.errs[0]
	e.errs = e.errs[1:]
	e.mu.Unlock()
	return err
}

// Len returns the number of errors in the stack.
func (e *Errors) Len() int {
	e.mu.RLock()
	l := len(e.errs)
	e.mu.RUnlock()
	return l
}

func (e *Errors) clear() {
	e.errs = make([]error, 0)
}

// Clear clears the error stack.
func (e *Errors) Clear() {
	if e.immutable {
		panic("clear called on immutable error stack")
	}
	e.mu.Lock()
	e.clear()
	e.mu.Unlock()
}

func (e *Errors) PopAll() []error {
	if e.immutable {
		panic("PopAll called on immutable error stack")
	}
	e.mu.RLock()
	retErrs := make([]error, len(e.errs))
	copy(retErrs, e.errs)
	e.clear()
	e.mu.RUnlock()
	return retErrs
}

func (e *Errors) Is(sought error) bool {
	if e.Len() == 0 {
		return sought == nil
	}
	e.mu.RLock()
	for _, err := range e.errs {
		if errors.Is(err, sought) {
			e.mu.RUnlock()
			return true
		}
	}
	e.mu.RUnlock()
	return false
}

func (e *Errors) As(target interface{}) bool {
	e.mu.RLock()
	for _, err := range e.errs {
		if //goland:noinspection GoErrorsAs
		errors.As(err, target) {
			e.mu.RUnlock()
			return true
		}
	}
	e.mu.RUnlock()
	return false
}

// Concat concatenates all errors in the stack into a single error. It does not clear the original stack.
func (e *Errors) Concat() error {
	e.mu.RLock()
	// errors.Join handles nil checks
	errStack := append(make([]error, 0), e.errs...)
	slices.Reverse[[]error](errStack)
	concat := errors.Join(errStack...)
	e.mu.RUnlock()
	return concat
}

// Errors returns a slice containing a copy of all errors in the stack. It does not clear the original stack.
func (e *Errors) Errors() []error {
	e.mu.RLock()
	errs := make([]error, len(e.errs))
	copy(errs, e.errs)
	e.mu.RUnlock()
	return errs
}

// Error implements the error interface. Internally it uses [Errors.Concat].
func (e *Errors) Error() string {
	return e.Concat().Error()
}
