//go:build linux

package linux

import (
	"errors"
	"strings"
	"syscall"
)

// unameFlags are the different uname switches for our uname syscall function.
type unameFlags uint8

var unameMap = map[[1]string]unameFlags{
	{"s"}: UnameOS,
	{"m"}: UnameArch,
	{"r"}: UnameRelease,
	{"d"}: UnameDomain,
	{"n"}: UnameHostname,
	{"v"}: UnameVersion,
}

func unameFlag(flag [1]string) (res unameFlags) {
	var ok bool
	res, ok = unameMap[flag]
	if !ok {
		return 0
	}
	return
}

const (
	// UnameNull is a placeholder.
	UnameNull unameFlags = iota
	// UnameOS is "sysname" or "uname -s".
	UnameOS
	// UnameArch is "machine" or "uname -m".
	UnameArch
	// UnameRelease is "release" or "uname -r".
	UnameRelease
	// UnameDomain is "domainname", the kernel domain name.
	UnameDomain
	// UnameHostname is "Nodename" or "uname -n"
	UnameHostname
	// UnameVersion is "version", or "uname -v".
	UnameVersion
)

// GetUname uses system calls to retrieve the same values as the uname linux command
func GetUname(unameFlags string) (un string, err error) {
	ub := &syscall.Utsname{}
	_ = syscall.Uname(ub)
	var targets []*[65]int8
	for _, n := range unameFlags {
		var flag = [1]string{string(n)}
		p := unameFlag(flag)
		switch p {
		case UnameNull:
			continue
		case UnameDomain:
			targets = append(targets, &ub.Domainname)
		case UnameHostname:
			targets = append(targets, &ub.Nodename)
		case UnameRelease:
			targets = append(targets, &ub.Release)
		case UnameArch:
			targets = append(targets, &ub.Machine)
		case UnameOS:
			targets = append(targets, &ub.Sysname)
		case UnameVersion:
			targets = append(targets, &ub.Version)
		}
	}

	if len(targets) < 2 {
		return "", errors.New("no valid uname targets in string")
	}

	var uns []string
	for _, target := range targets {
		var sub []string
		for _, r := range target {
			sub = append(sub, string(rune(r)))
		}
		uns = append(uns, strings.Join(sub, ""))
	}

	un = strings.Join(uns, " ")

	return
}
