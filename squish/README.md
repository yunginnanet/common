# squish
--
    import "git.tcp.direct/kayos/common/squish"

## Usage

#### func  B64d

```go
func B64d(str string) (data []byte)
```
B64d decodes the given string into the original slice of bytes. Do note that
this is for non critical tasks, it has no error handling for purposes of clean
code.

#### func  B64e

```go
func B64e(cytes []byte) (data string)
```
B64e encodes the given slice of bytes into base64 standard encoding.

#### func  Gunzip

```go
func Gunzip(data []byte) ([]byte, error)
```
Gunzip decompresses a gzip compressed slice of bytes.

#### func  Gzip

```go
func Gzip(data []byte) ([]byte, error)
```
Gzip compresses as slice of bytes using gzip compression.

#### func  UnpackStr

```go
func UnpackStr(encoded string) string
```
UnpackStr UNsafely unpacks (usually banners) that have been base64'd and then
gzip'd.
