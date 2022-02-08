# common

Welcome to things. Here are some of the aforementioned:

* [hash](./hash)

* [linux](./linux)

* [squish](./squish)

* [entropy](./entropy)

## base

    `import "git.tcp.direct/kayos/common"`

### Usage

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

#### func  OneInA

```go
func OneInA(million int) bool
```

#### func  RNG

```go
func RNG(n int) int
```
RNG is a Random Number Generator

#### func  RandSleepMS

```go
func RandSleepMS(n int)
```
RandSleepMS sleeps for a random period of time with a maximum of n milliseconds.

#### func  RandStr

```go
func RandStr(size int) string
```
RandStr generates a random alphanumeric string with a max length of size.

#### func  RangeIterate

```go
func RangeIterate(ips interface{}) chan *ipa.IP
```
