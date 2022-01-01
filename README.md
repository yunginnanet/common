# common
--
    import "git.tcp.direct/kayos/common"

stuff I use a lot.

#### func  Abs

```go
func Abs(n int) int
```
Abs will give you the positive version of a negative integer, quickly.

#### func  BytesToBlake2b

```go
func BytesToBlake2b(b []byte) []byte
```
BytesToBlake2b ignores all errors and gives you a blakae2b 64 hash value as a
byte slice. (or panics somehow)

#### func  CompareChecksums

```go
func CompareChecksums(a []byte, b []byte) bool
```
CompareChecksums will take in two byte slices, hash them with blake2b, and tell
you if the resulting checksums match.

#### func  Fprint

```go
func Fprint(w io.Writer, s string)
```
Fprint is fmt.Fprint with error handling.

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
