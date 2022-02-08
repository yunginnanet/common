# hash

    `import "git.tcp.direct/kayos/common/hash"`


## Usage

#### func  Blake2bSum

```go
func Blake2bSum(b []byte) []byte
```
Blake2bSum ignores all errors and gives you a blakae2b 64 hash value as a byte
slice. (or panics somehow)

#### func  BlakeEqual

```go
func BlakeEqual(a []byte, b []byte) bool
```
BlakeEqual will take in two byte slices, hash them with blake2b, and tell you if
the resulting checksums match.
