//go:build linux

package net

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func Test_getDefaultIPv4Route(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("skipping test; must be root to run")
	}
	// use cmd exec to parse and check `ip route show	 default` for a control value
	cmd := exec.Command("ip", "route", "show", "default")
	out, err := cmd.Output()
	if err != nil {
		t.Skipf("failed to get default route with error: %e", err)
	}
	if len(out) < 1 {
		t.Skip("failed to get default route control value")
	}
	fields := bytes.Fields(out)
	if len(fields) < 3 {
		t.Skip("failed to get default route control value")
	}
	value := fields[2]
	if len(value) < 1 {
		t.Skip("failed to get default route control value")
	}
	t.Logf("control value: %s", string(value))
	var res string
	if res, err = getDefaultIPv4RouteString(); err != nil {
		t.Fatalf("failed to get default route with error: %e", err)
	}
	if !strings.EqualFold(res, string(value)) {
		t.Fatalf("failed to get default route with error: %e", err)
	}
	t.Logf("default route: %s", res)
}
