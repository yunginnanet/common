# pool

`import git.tcp.direct/kayos/common/pool`

## Overview

pool contains two components, both of which are forms of buffer pools.

  ### *BufferFactory*

  BufferFactory is a potentially safer sync.Pool of [bytes.Buffer](https://pkg.go.dev/bytes#Buffer) types that will not allow you to accidentally re-use buffers after you return them to the pool.

  - `func NewBufferFactory() BufferFactory`
  - `func NewSizedBufferFactory(size int) BufferFactory`
  - `func (cf BufferFactory) Get() *Buffer`
  - `func (cf BufferFactory) MustPut(buf *Buffer)`
  - `func (cf BufferFactory) Put(buf *Buffer) error`

  ### *StringFactory*

  StringFactory is very much like BufferFactory, except for it's a pool of [strings.Builder](https://pkg.go.dev/strings#Builder) types instead.

  - `func NewStringFactory() StringFactory`
  - `func (sf StringFactory) Get() *String`
  - `func (sf StringFactory) MustPut(buf *String)`
  - `func (sf StringFactory) Put(buf *String) error`

## Benchmarks

#### foreword

In some usecases, this package will actually allocate slightly more than if one were not using it. That is because each Buffer or String type has the additional header of a [sync.Once](https://pkg.go.dev/sync#Once). The overhead of this is minimal, and in the benchmarks I've put together that attempt to emulate a more realistic use case, the benefits are easy to see. This is comparing it to use without buffer pools at all. Using [sync.Pool](https://pkg.go.dev/sync#Pool) alone without this package will always be ever-so-slightly more efficient, but the consequences of mis-using it could be catastrophic. That is the trade-off we make.

_Note: "bytes" here refers to the size of buffer pre-allocation. Either with NewSizedBufferFactory or bytes.NewBuffer_

---

### Using this package

```
BenchmarkBufferFactory/SingleProc-64-bytes-24         	 410792	     3467 ns/op	  24685 B/op	      5 allocs/op
BenchmarkBufferFactory/SingleProc-1024-bytes-24       	 376539	     3668 ns/op	  24687 B/op	      5 allocs/op
BenchmarkBufferFactory/SingleProc-4096-bytes-24       	 281402	     3599 ns/op	  24688 B/op	      5 allocs/op
BenchmarkBufferFactory/SingleProc-65536-bytes-24      	 340872	     3591 ns/op	  24830 B/op	      5 allocs/op
```

```
BenchmarkBufferFactory/Concurrent-x2-64-bytes-24      	 498634	     2076 ns/op	  24678 B/op	      5 allocs/op
BenchmarkBufferFactory/Concurrent-x2-1024-bytes-24    	 569001	     2128 ns/op	  24684 B/op	      5 allocs/op
BenchmarkBufferFactory/Concurrent-x2-4096-bytes-24    	 602946	     2131 ns/op	  24688 B/op	      5 allocs/op
BenchmarkBufferFactory/Concurrent-x2-65536-bytes-24   	 779770	     1522 ns/op	  24747 B/op	      5 allocs/op
```

```
BenchmarkBufferFactory/Concurrent-x4-64-bytes-24      	 319677	     3271 ns/op	  24695 B/op	      5 allocs/op
BenchmarkBufferFactory/Concurrent-x4-1024-bytes-24    	 597859	     2005 ns/op	  24677 B/op	      5 allocs/op
BenchmarkBufferFactory/Concurrent-x4-4096-bytes-24    	 586940	     2092 ns/op	  24685 B/op	      5 allocs/op
BenchmarkBufferFactory/Concurrent-x4-65536-bytes-24   	 861326	     1420 ns/op	  24719 B/op	      5 allocs/op
```

```
BenchmarkBufferFactory/Concurrent-x8-64-bytes-24      	 621253	     1893 ns/op	  24672 B/op	      5 allocs/op
BenchmarkBufferFactory/Concurrent-x8-1024-bytes-24    	 626067	     2155 ns/op	  24678 B/op	      5 allocs/op
BenchmarkBufferFactory/Concurrent-x8-4096-bytes-24    	 644078	     2030 ns/op	  24682 B/op	      5 allocs/op
BenchmarkBufferFactory/Concurrent-x8-65536-bytes-24   	 878794	     1382 ns/op	  24718 B/op	      5 allocs/op
```

---

### Allocating new bytes.Buffers

```
BenchmarkNotUsingPackage/SingleProc-64-bytes-24       	 121978	     9401 ns/op	  52000 B/op	      5 allocs/op
BenchmarkNotUsingPackage/SingleProc-1024-bytes-24     	 129132	     9004 ns/op	  52000 B/op	      5 allocs/op
BenchmarkNotUsingPackage/SingleProc-4096-bytes-24     	 138590	     8894 ns/op	  52000 B/op	      5 allocs/op
BenchmarkNotUsingPackage/SingleProc-65536-bytes-24    	 134730	     8891 ns/op	  52000 B/op	      5 allocs/op
```

```
BenchmarkNotUsingPackage/Concurrent-x2-64-bytes-24    	 150890	     8488 ns/op	  52225 B/op	      8 allocs/op
BenchmarkNotUsingPackage/Concurrent-x2-1024-bytes-24  	 174734	     8924 ns/op	  55937 B/op	      7 allocs/op
BenchmarkNotUsingPackage/Concurrent-x2-4096-bytes-24  	 101695	    11032 ns/op	  65537 B/op	      6 allocs/op
BenchmarkNotUsingPackage/Concurrent-x2-65536-bytes-24 	  43708	    28977 ns/op	 286724 B/op	      5 allocs/op
```

```
BenchmarkNotUsingPackage/Concurrent-x4-64-bytes-24    	 189124	     6314 ns/op	  52224 B/op	      8 allocs/op
BenchmarkNotUsingPackage/Concurrent-x4-1024-bytes-24  	 144324	     7010 ns/op	  55937 B/op	      7 allocs/op
BenchmarkNotUsingPackage/Concurrent-x4-4096-bytes-24  	 129194	     9325 ns/op	  65537 B/op	      6 allocs/op
BenchmarkNotUsingPackage/Concurrent-x4-65536-bytes-24 	  43641	    28730 ns/op	 286723 B/op	      5 allocs/op
```

```
BenchmarkNotUsingPackage/Concurrent-x8-64-bytes-24    	 184123	     5911 ns/op	  52224 B/op	      8 allocs/op
BenchmarkNotUsingPackage/Concurrent-x8-1024-bytes-24  	 172016	     6602 ns/op	  55937 B/op	      7 allocs/op
BenchmarkNotUsingPackage/Concurrent-x8-4096-bytes-24  	 129916	     9030 ns/op	  65537 B/op	      6 allocs/op
BenchmarkNotUsingPackage/Concurrent-x8-65536-bytes-24 	  45729	    26344 ns/op	 286723 B/op	      5 allocs/op
```
