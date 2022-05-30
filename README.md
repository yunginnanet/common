# common
[![GoDoc](https://godoc.org/git.tcp.direct/kayos/common?status.svg)](https://pkg.go.dev/git.tcp.direct/kayos/common) [![codecov](https://codecov.io/gh/yunginnanet/common/branch/master/graph/badge.svg?token=vk5frSGqhq)](https://codecov.io/gh/yunginnanet/common)


Welcome to things. Here are some of the aforementioned:

* [hash](https://pkg.go.dev/git.tcp.direct/kayos/common/hash)

* [linux](https://pkg.go.dev/git.tcp.direct/kayos/common/linux)

* [squish](https://pkg.go.dev/git.tcp.direct/kayos/common/squish)

* [entropy](https://pkg.go.dev/git.tcp.direct/kayos/common/entropy)

* [network](https://pkg.go.dev/git.tcp.direct/kayos/common/network)

## base

`import "git.tcp.direct/kayos/common"`

## Base Module

#### func  Abs

```go
func Abs(n int) int
```
Abs will give you the positive version of a negative integer, quickly.

#### func  BytesToFloat64

```go
func BytesToFloat64(bytes []byte) float64
```
BytesToFloat64 will take a slice of bytes and convert it to a float64 type.

#### func  Float64ToBytes

```go
func Float64ToBytes(f float64) []byte
```
Float64ToBytes will take a float64 type and convert it to a slice of bytes.

#### func  Fprint

```go
func Fprint(w io.Writer, s string)
```
Fprint is fmt.Fprint with error handling.
