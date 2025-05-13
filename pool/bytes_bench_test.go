package pool

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

// =========================================================================

// BenchmarkBufferFactory tries to emulate real world usage of a buffer pool.
// It creates a buffer, writes to it, and then returns it to the pool.
//
// Then it repeats this process b.N times.
// This should be a decent way to test the performance of a buffer pool.
//
// See bytes_bench.go for more information.
func BenchmarkBufferFactory(b *testing.B) {
	benchmarkBufferFactory(b)
}

// BenchmarkNotUsingPackage is a benchmark that does not use github.com/yunginnanet/common/pool.
// It mimics the behavior of the BufferFactory benchmark, but does not use a buffer pool.
//
// See bytes_test.go for more information.
func BenchmarkNotUsingPackage(b *testing.B) {
	benchmarkNewBytesBuffer(b)
}

// =========================================================================

const (
	hello64 = `its me ur new best friend tell me a buf i tell you where it ends`
	hello   = `hello world, it's me, your new best friend. tell me the buffer and i'll tell you where it ends.`
	lip     = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed a ante sit amet purus blandit auctor. Nullam ornare enim sed nibh consequat molestie. Duis est lectus, vestibulum vel felis vel, convallis cursus ex. Morbi nec placerat orci. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Praesent a erat sit amet libero convallis ornare a venenatis dolor. Pellentesque euismod risus et metus porttitor, vel consectetur lacus tempus. Integer elit arcu, condimentum quis nisi eget, dapibus imperdiet nulla. Cras sit amet ante in urna varius tempus. Integer tristique sagittis nunc vel tincidunt. Integer non suscipit ligula, et fermentum sem. Duis id odio lorem. Sed id placerat urna, eu vehicula risus. Duis porttitor hendrerit risus. Curabitur id tellus ac arcu aliquet finibus. Pellentesque et nisl ante. Mauris sapien nisl, pretium in ligula tempus, posuere mattis turpis. Proin et tempus enim. Nullam at diam est. Vivamus ut lectus hendrerit, interdum ex id, ultricies sapien. Praesent rhoncus turpis dolor, quis lobortis tortor pellentesque id. Pellentesque eget nisi laoreet, fringilla augue eu, cursus risus. Integer consectetur ornare laoreet. Praesent ligula sem, tincidunt at ligula at, condimentum venenatis tortor. Nam laoreet enim leo, sed finibus lorem egestas vel. Maecenas varius a leo non placerat. Donec scelerisque, risus vel finibus ornare, arcu ligula interdum justo, in ultricies urna mi et neque. Curabitur sed sem dui. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Maecenas eget laoreet nisi. Nam rhoncus sapien ac interdum sagittis. Nulla fermentum sem nec tellus dignissim lacinia. Curabitur ornare lectus non dictum laoreet. Praesent tempor risus at tortor tempor finibus. Cras id dolor mi. Mauris ut mi quis est vehicula molestie. Mauris eu varius urna. Integer sodales nunc at risus rutrum eleifend. In sed bibendum lectus. Morbi ipsum sapien, blandit in dignissim eu, ultrices non odio. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Nulla eget volutpat ligula, at elementum dui. Aliquam sed enim scelerisque, facilisis magna vitae, dignissim enim. Pellentesque non ultricies urna. Proin fermentum erat semper efficitur auctor. Vestibulum posuere non tortor vitae tincidunt.`
)

var lipTenIcedTea = strings.Repeat(lip, 10)

func poolbench(f BufferFactory) {
	buf := f.Get()
	buf.MustWrite([]byte(hello64))
	f.MustPut(buf)
	buf = f.Get()
	buf.MustWrite([]byte(hello))
	f.MustPut(buf)
	buf = f.Get()
	buf.MustWrite([]byte(lip))
	f.MustPut(buf)
	buf = f.Get()
	buf.MustWrite([]byte(lipTenIcedTea))
	f.MustPut(buf)
}

func parabench(pb *testing.PB, f BufferFactory) {
	for pb.Next() {
		poolbench(f)
	}
}

func sized(b *testing.B, size int, para int) {
	b.ReportAllocs()

	f := NewBufferFactory()
	if size != 0 {
		f = NewSizedBufferFactory(size)
	}

	if para != 0 {
		b.SetParallelism(para)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) { parabench(pb, f) })
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolbench(f)
	}
}

// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------

func sizedbytesbench(initial func() []byte) {
	buf := bytes.NewBuffer(initial())
	buf.Write([]byte(hello64))
	buf = bytes.NewBuffer(initial())
	buf.Write([]byte(hello))
	buf = bytes.NewBuffer(initial())
	buf.Write([]byte(lip))
	buf = bytes.NewBuffer(initial())
	buf.Write([]byte(lipTenIcedTea))
}

func parabytesbench(pb *testing.PB, size int) {
	for pb.Next() {
		if size != 0 {
			sizedbytesbench(func() []byte { return make([]byte, 0, size) })
		} else {
			sizedbytesbench(func() []byte { return nil })
		}
	}
}

func bytessized(b *testing.B, size int, para int) {
	b.ReportAllocs()
	if para != 0 {
		b.SetParallelism(para)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) { parabytesbench(pb, size) })
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sizedbytesbench(func() []byte { return nil })
	}
}

var divier = []byte(strings.Repeat("-", 130) + "\n")

func bigDivide(n int) {
	if n > 1 {
		_, _ = os.Stdout.Write([]byte{'\n'})
	}
	for i := 0; i < n; i++ {
		_, _ = os.Stdout.Write(divier)
	}

	_, _ = os.Stdout.Write([]byte{'\n'})
}

func benchmarkBufferFactory(b *testing.B) {
	b.ReportAllocs()
	concurrency := []int{0, 2, 4, 8}
	size := []int{64, 1024, 4096, 65536}

	defer bigDivide(2)

	for _, c := range concurrency {
		for _, s := range size {
			label := fmt.Sprintf("Concurrent-x%d-%d-bytes", c, s)
			if c == 0 {
				label = fmt.Sprintf("SingleProc-%d-bytes", s)
			}
			b.Run(label, func(b *testing.B) { sized(b, s, c) })
		}
		_, _ = os.Stdout.Write(divier)
	}
}

func benchmarkNewBytesBuffer(b *testing.B) {
	b.ReportAllocs()
	concurrency := []int{0, 2, 4, 8}
	size := []int{64, 1024, 4096, 65536}

	defer bigDivide(1)

	for _, c := range concurrency {
		for _, s := range size {
			label := fmt.Sprintf("Concurrent-x%d-%d-bytes", c, s)
			if c == 0 {
				label = fmt.Sprintf("SingleProc-%d-bytes", s)
			}
			b.Run(label, func(b *testing.B) { bytessized(b, s, c) })
		}
		_, _ = os.Stdout.Write(divier)
	}
}
