package pool

import (
	"bytes"
	"sync"
	"testing"
)

type othaBuffa struct {
	ByteBuffer
}

// ensure compatibility with interface
func TestInterfaces(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			t.Error("interface not implemented")
		}
	}()
	var (
		bf       any = NewBufferFactory()
		bfCompat any = BufferFactoryInterfaceCompat{NewBufferFactory()}
		sPool    any = &sync.Pool{
			New: func() any { return new(bytes.Buffer) },
		}
	)
	if _, ok := sPool.(Pool[any]); !ok {
		t.Fatal("Pool[any] not implemented by sync.Pool")
	}
	testMe1, ok1 := bfCompat.(Pool[*Buffer])
	if !ok1 {
		t.Fatal("Pool[*Buffer] not implemented")
	}

	t.Run("Pool", func(t *testing.T) {
		t.Parallel()
		b := testMe1.Get()
		if _, err := b.WriteString("test"); err != nil {
			t.Fatal(err)
		}
		testMe1.Put(b)
		b = testMe1.Get()
		if b.Len() != 0 {
			t.Fatal("buffer not reset")
		}
		testMe1.Put(b)
	})

	t.Run("PoolWithPutError", func(t *testing.T) {
		t.Parallel()
		testMe2, ok2 := bf.(WithPutError[*Buffer])
		if !ok2 {
			t.Error("PoolWithPutError[*Buffer] not implemented")
		}
		b := testMe2.Get()
		if _, err := b.WriteString("test"); err != nil {
			t.Fatal(err)
		}
		if err := testMe2.Put(b); err != nil {
			t.Fatal(err)
		}
		b = testMe2.Get()
		if b.Len() != 0 {
			t.Fatal("buffer not reset")
		}
		if err := testMe2.Put(b); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("BufferFactoryByteBufferCompat", func(t *testing.T) {
		t.Parallel()
		bf := BufferFactoryByteBufferCompat{NewBufferFactory()}
		b := bf.Get()
		if _, err := b.WriteString("test"); err != nil {
			t.Fatal(err)
		}
		bf.Put(b)
		b = bf.Get()
		if b.Len() != 0 {
			t.Fatal("buffer not reset")
		}
		foreign := &bytes.Buffer{}
		foreign.WriteString("test")
		bf.Put(foreign)
		if foreign.Len() != 0 {
			t.Fatal("buffer not reset")
		}
		foreignGot := bf.Get()
		if foreignGot.Len() != 0 {
			t.Fatal("buffer not reset")
		}
		bf.Put(foreignGot)
		t.Run("must panic after wrapped and put twice", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("panic expected")
				}
			}()
			bf.Put(foreignGot)
		})
		t.Run("must panic on invalid type", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("panic expected")
				}
			}()
			bf.Put(&othaBuffa{})
		})

	})

}
