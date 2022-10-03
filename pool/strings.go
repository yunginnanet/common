package pool

import (
	"strings"
	"sync"
)

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
	buf.o.Do(func() {
		_ = buf.Reset()
		sf.pool.Put(buf.Builder)
		buf.Builder = nil
		err = nil
	})
	return err
}

func (sf StringFactory) MustPut(buf *String) {
	if err := sf.Put(buf); err != nil {
		panic(err)
	}
}

// Get returns a strings.Builder from the pool.
func (sf StringFactory) Get() *String {
	return &String{
		sf.pool.Get().(*strings.Builder),
		&sync.Once{},
	}
}

type String struct {
	*strings.Builder
	o *sync.Once
}

func (s String) String() string {
	if s.Builder == nil {
		return ""
	}
	return s.Builder.String()
}

func (s String) MustString() string {
	if s.Builder == nil {
		panic(ErrBufferReturned)
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

func (s String) MustReset() {
	if err := s.Reset(); err != nil {
		panic(err)
	}
	s.Builder.Reset()
}

func (s String) WriteString(str string) (int, error) {
	if s.Builder == nil {
		return 0, ErrBufferReturned
	}
	return s.Builder.WriteString(str)
}

// MustWriteString means Must Write String, like WriteString but will panic on error.
func (s String) MustWriteString(str string) {
	if s.Builder == nil {
		panic(ErrBufferReturned)
	}
	if str == "" {
		panic("nil string")
	}
	_, _ = s.Builder.WriteString(str)
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

func (s String) Len() int {
	if s.Builder == nil {
		return 0
	}
	return s.Builder.Len()
}

func (s String) MustLen() int {
	if s.Builder == nil {
		panic(ErrBufferReturned)
	}
	return s.Builder.Len()
}

func (s String) Grow(n int) error {
	if s.Builder == nil {
		return ErrBufferReturned
	}
	s.Builder.Grow(n)
	return nil
}

func (s String) MustGrow(n int) {
	if s.Builder == nil {
		panic(ErrBufferReturned)
	}
	s.Builder.Grow(n)
}

func (s String) Cap() int {
	if s.Builder == nil {
		return 0
	}
	return s.Builder.Cap()
}
