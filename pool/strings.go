package pool

import (
	"errors"
	"strings"
	"sync"
)

var ErrBufferReturned = errors.New("buffer already returned")

type String struct {
	*strings.Builder
	*sync.Once
}

func (s String) String() string {
	if s.Builder == nil {
		return ""
	}
	return s.Builder.String()
}

func (s String) Reset() error {
	if s.Builder == nil {
		return ErrBufferReturned
	}
	s.Builder.Reset()
	return nil
}

func (s String) WriteString(str string) (int, error) {
	if s.Builder == nil {
		return 0, ErrBufferReturned
	}
	return s.Builder.WriteString(str)
}

func (s String) Len() int {
	if s.Builder == nil {
		return 0
	}
	return s.Builder.Len()
}

func (s String) Write(p []byte) (int, error) {
	if s.Builder == nil {
		return 0, ErrBufferReturned
	}
	return s.Builder.Write(p)
}

func (s String) WriteRune(r rune) (int, error) {
	if s.Builder == nil {
		return 0, ErrBufferReturned
	}
	return s.Builder.WriteRune(r)
}

func (s String) WriteByte(c byte) error {
	if s.Builder == nil {
		return ErrBufferReturned
	}
	return s.Builder.WriteByte(c)
}

func (s String) Grow(n int) error {
	if s.Builder == nil {
		return ErrBufferReturned
	}
	s.Builder.Grow(n)
	return nil
}

func (s String) Cap() int {
	if s.Builder == nil {
		return 0
	}
	return s.Builder.Cap()
}

type StringFactory struct {
	pool *sync.Pool
}

// NewStringFactory creates a new strings.Builder pool.
func NewStringFactory() StringFactory {
	return StringFactory{
		pool: &sync.Pool{
			New: func() any { return new(strings.Builder) },
		},
	}
}

// Put returns a strings.Builder back into to the pool after resetting it.
func (sf StringFactory) Put(buf *String) error {
	var err = ErrBufferReturned
	buf.Do(func() {
		_ = buf.Reset()
		sf.pool.Put(buf.Builder)
		buf.Builder = nil
		err = nil
	})
	return err
}

// Get returns a strings.Builder from the pool.
func (sf StringFactory) Get() *String {
	return &String{
		sf.pool.Get().(*strings.Builder),
		&sync.Once{},
	}
}
