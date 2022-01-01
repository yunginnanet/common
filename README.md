# common
--
    import "git.tcp.direct/kayos/common"

stuff I use a lot.


#### func  Abs

```go
func Abs(n int) int
```

#### func  CompareChecksums

```go
func CompareChecksums(a []byte, b []byte) bool
```

#### func  Fprint

```go
func Fprint(w io.Writer, s string)
```
Fprint is fmt.Fprint with error handling.

#### func  RNG

```go
func RNG(n int) int
```

Random Number Generator (uses a combo of crypto/rand and math/rand for better performance)

#### func  RangeIterate

```go
func RangeIterate(ips interface{}) chan *ipa.IP
```

IP Address iteration

#### func  SnoozeMS

```go
func SnoozeMS(n int)
```

Random sleep, max of n seconds.
