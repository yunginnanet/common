# xerrors

Package xerrors provides a stack of multiple errors that can be pushed to,
popped from, and concatenated.


#### type ErrorStack

```go
type ErrorStack interface {
	Len() int
	Is(error) bool
	As(interface{}) bool
	Errors() []error
	Next() error
	io.Seeker
	error
}
```


#### type Errors

```go
type Errors struct {
}
```

Errors is a stack of multiple errors that can be pushed to, popped from, or
concatenated. It is safe for concurrent use, and can return an immutable copy of
itself as an [errstack.ErrorsImmutable].

#### func  NewErrors

```go
func NewErrors() *Errors
```
NewErrors returns a new [Errors] stack.

#### func (*Errors) As

```go
func (e *Errors) As(target interface{}) bool
```

#### func (*Errors) Clear

```go
func (e *Errors) Clear()
```
Clear clears the error stack.

#### func (*Errors) Concat

```go
func (e *Errors) Concat() error
```
Concat concatenates all errors in the stack into a single error. It does not
clear the original stack.

#### func (*Errors) Copy

```go
func (e *Errors) Copy() *ErrorsImmutable
```
Copy returns an immutable copy of the error stack. It does not clear the
original stack. Copy is safe to call concurrently.

#### func (*Errors) Error

```go
func (e *Errors) Error() string
```
Error implements the error interface. Internally it uses [Errors.Concat].

#### func (*Errors) Errors

```go
func (e *Errors) Errors() []error
```
Errors returns a slice containing a copy of all errors in the stack. It does not
clear the original stack.

#### func (*Errors) Is

```go
func (e *Errors) Is(sought error) bool
```

#### func (*Errors) Len

```go
func (e *Errors) Len() int
```
Len returns the number of errors in the stack.

#### func (*Errors) Next

```go
func (e *Errors) Next() error
```
Next returns the next error in the stack, incrementing the internal index. If
we've reached the end of the stack, it returns nil.

Use [errstack.Errors.Seek] to rewind the internal index if needed.

#### func (*Errors) Pop

```go
func (e *Errors) Pop() error
```
Pop pops one error from the stack, removing it from the stack and returning it.

#### func (*Errors) PopAll

```go
func (e *Errors) PopAll() []error
```

#### func (*Errors) PopAllImmutable

```go
func (e *Errors) PopAllImmutable() *ErrorsImmutable
```
PopAllImmutable returns an immutable copy of the error stack, and clears the
original stack, leaving it empty.

#### func (*Errors) Push

```go
func (e *Errors) Push(err error)
```
Push adds an error to the stack. It is safe for concurrent use.

#### func (*Errors) Seek

```go
func (e *Errors) Seek(offset int64, whence int) (int64, error)
```
Seek implements an [io.Seeker] for the purposes of controlling
[errstack.Errors.Next] output.

#### type ErrorsImmutable

```go
type ErrorsImmutable struct {
}
```

ErrorsImmutable is a stack of multiple errors popped from [errstack.Errors].
Internally it contains a private [errstack.Errors] with the immutable flag set
to true.

It's public methods are a partial subset of [errstack.Errors] methods operating
as pass-throughs. Consequently, more information on the contained methods can be
found in the [errstack.Errors] documentation.

#### func (*ErrorsImmutable) As

```go
func (e *ErrorsImmutable) As(i interface{}) bool
```

#### func (*ErrorsImmutable) Error

```go
func (e *ErrorsImmutable) Error() string
```

#### func (*ErrorsImmutable) Errors

```go
func (e *ErrorsImmutable) Errors() []error
```

#### func (*ErrorsImmutable) Is

```go
func (e *ErrorsImmutable) Is(err error) bool
```

#### func (*ErrorsImmutable) Len

```go
func (e *ErrorsImmutable) Len() int
```

#### func (*ErrorsImmutable) Next

```go
func (e *ErrorsImmutable) Next() error
```

#### func (*ErrorsImmutable) Seek

```go
func (e *ErrorsImmutable) Seek(offset int64, whence int) (int64, error)
```

---
