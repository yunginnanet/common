package pool

type Pool[T any] interface {
	Get() T
	Put(T)
}

type WithPutError[T any] interface {
	Get() T
	Put(T) error
}

type BufferFactoryInterfaceCompat struct {
	BufferFactory
}

func (b BufferFactoryInterfaceCompat) Put(buf *Buffer) {
	_ = b.BufferFactory.Put(buf)
}
