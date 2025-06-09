package pool

import "errors"

var (
	ErrBufferReturned   = errors.New("buffer already returned")
	ErrBufferHasNoMutex = errors.New("buffer has no mutex, use WithMutex method to acquire a mutex for the buffer")
)
