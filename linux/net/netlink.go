//go:build linux

package net

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/vishvananda/netlink"
)

var ErrNoDefaultRoute = errors.New("no default route")

func getDefaultIPv4RouteString() (string, error) {
	routes, err := netlink.RouteListFiltered(syscall.AF_INET, &netlink.Route{Dst: nil}, netlink.RT_FILTER_DST)
	if err != nil {
		return "", fmt.Errorf("failed to get default route with error: %w", err)
	}
	for _, route := range routes {
		if route.Dst == nil {
			return route.Gw.String(), nil
		}
	}
	return "", ErrNoDefaultRoute
}
