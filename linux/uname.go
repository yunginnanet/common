//go:build linux

package linux

import (
	"strings"
	"syscall"
)

// unameFlags are the different uname switches for our uname syscall function.
type unameFlags uint8

var unameMap = map[[1]string]unameFlags{
	[1]string{"s"}: UnameOS,
	[1]string{"m"}: UnameArch,
	[1]string{"r"}: UnameRelease,
	[1]string{"d"}: UnameDomain,
	[1]string{"n"}: UnameHostname,
	[1]string{"v"}: UnameVersion,
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
	// UnameVersion is "version", or "uname -v".
	UnameVersion
	// UnameHostname is "Nodename" or "uname -n"
	UnameHostname
)

func GetUname(unameFlags string) (un string, err error) {
	ub := &syscall.Utsname{}
	err = syscall.Uname(ub)
	if err != nil {
		return
	}
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
		case UnameArch:
			targets = append(targets, &ub.Machine)
		case UnameOS:
			targets = append(targets, &ub.Sysname)
		case UnameVersion:
			targets = append(targets, &ub.Version)
		}
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
