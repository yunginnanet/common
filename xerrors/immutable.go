package xerrors

// ErrorsImmutable is a stack of multiple errors popped from [errstack.Errors].
// Internally it contains a private [errstack.Errors] with the immutable flag set to true.
//
// It's public methods are a partial subset of [errstack.Errors] methods operating as pass-throughs.
// Consequently, more information on the contained methods can be found in the [errstack.Errors] documentation.
type ErrorsImmutable struct {
	e Errors // immutable
	i int
}

func (e *Errors) immutableCopy() *ErrorsImmutable {
	errs := make([]error, len(e.errs))
	copy(errs, e.errs)
	return &ErrorsImmutable{
		e: Errors{errs: errs, immutable: true},
	}
}

// Copy returns an immutable copy of the error stack. It does not clear the original stack.
// Copy is safe to call concurrently.
func (e *Errors) Copy() *ErrorsImmutable {
	e.mu.RLock()
	ec := e.immutableCopy()
	e.mu.RUnlock()
	return ec
}

// PopAllImmutable returns an immutable copy of the error stack, and clears the original stack, leaving it empty.
func (e *Errors) PopAllImmutable() *ErrorsImmutable {
	if e.immutable {
		panic("PopAllImmutable called on immutable error stack")
	}
	e.mu.Lock()
	ec := e.immutableCopy()
	e.clear()
	e.mu.Unlock()
	return ec
}

func (e *ErrorsImmutable) Len() int {
	return e.e.Len()
}

func (e *ErrorsImmutable) Is(err error) bool {
	return e.e.Is(err)
}

func (e *ErrorsImmutable) As(i interface{}) bool {
	return e.e.As(i)
}

func (e *ErrorsImmutable) Error() string {
	return e.e.Error()
}

func (e *ErrorsImmutable) Errors() []error {
	return e.e.Errors()
}

func (e *ErrorsImmutable) Next() error {
	return e.e.next(&e.i)
}

func (e *ErrorsImmutable) Seek(offset int64, whence int) (int64, error) {
	return e.e.seek(&e.i, offset, whence)
}
