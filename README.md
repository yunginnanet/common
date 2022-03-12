# common
[![GoDoc](https://godoc.org/git.tcp.direct/kayos/common?status.svg)](https://pkg.go.dev/git.tcp.direct/kayos/common) [![codecov](https://codecov.io/gh/yunginnanet/common/branch/master/graph/badge.svg?token=vk5frSGqhq)](https://codecov.io/gh/yunginnanet/common)


Welcome to things. Here are some of the aforementioned:

* [hash](https://pkg.go.dev/git.tcp.direct/kayos/common/hash)

* [linux](https://pkg.go.dev/git.tcp.direct/kayos/common/linux)

* [squish](https://pkg.go.dev/git.tcp.direct/kayos/common/squish)

* [entropy](https://pkg.go.dev/git.tcp.direct/kayos/common/entropy)

## base

    `import "git.tcp.direct/kayos/common"`

## Base

#### func  Abs

```go
func Abs(n int) int
```
Abs will give you the positive version of a negative integer, quickly.

#### func  Fprint

```go
func Fprint(w io.Writer, s string)
```
Fprint is fmt.Fprint with error handling.
