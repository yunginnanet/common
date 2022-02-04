## linux

    import "git.tcp.direct/kayos/common/linux"

### Usage

```go
const (
	// UnameOS is "sysname" or "uname -s".
	UnameOS
	// UnameArch is "machine" or "uname -m".
	UnameArch
	// UnameRelease is "release" or "uname -r".
	UnameRelease
	// UnameDomain is "domainname", the kernel domain name.
	UnameDomain
	// UnameVersion is "version", or "uname -v".
	UnameVersion
	// UnameHostname is "Nodename" or "uname -n"
	UnameHostname
)
```

#### func  GetUname

```go
func GetUname(unameFlags string) (un string, err error)
```
GetUname uses system calls to retrieve the same values as the uname linux
command
