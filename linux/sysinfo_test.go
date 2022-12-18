//go:build linux

package linux

import (
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// getUptimeControlValue from /proc/uptime to compare against our syscall value.
func getUptimeControlValue(t *testing.T) time.Duration {
	f, err := os.Open("/proc/uptime")
	if err != nil {
		t.Fatalf("failed to open /proc/uptime with error: %v", err)
	}
	buf, err := io.ReadAll(f)

	if err != nil {
		t.Fatalf("failed to read /proc/uptime with error: %v", err)
	}

	t.Logf("read %d bytes from /proc/uptime: %s", len(buf), buf)

	controlSeconds, err := strconv.ParseInt(string(buf[:strings.IndexByte(string(buf), '.')]), 10, 64)
	if err != nil {
		t.Fatalf("failed to parse /proc/uptime with error: %e", err)
	}

	if controlSeconds < 1 {
		t.Fatalf("failed to parse /proc/uptime (zero value)")
	}

	return time.Duration(controlSeconds) * time.Second
}

func TestUptime(t *testing.T) {
	uptimeCtrl := getUptimeControlValue(t)
	t.Logf("control uptime: %v", uptimeCtrl)

	uptime, err := Uptime()
	if err != nil {
		t.Fatalf("failed to get uptime with error: %e", err)
	}

	if uptime < 1 {
		t.Fatalf("failed to get uptime (zero value)")
	}

	t.Logf("uptime: %v", uptime)

	matching := uptime == uptimeCtrl

	// if somehow the uptime is less than the control, which was called first
	// then this has failed terribly.
	// If it's greater, then it's possible we are within an acceptable tolerance.

	if !matching && uptime < uptimeCtrl {
		t.Fatalf("uptime does not match control uptime (uptime < uptimeCtrl)!!")
	}

	if !matching {
		t.Logf("no match, allowing for a 1 second tolerance...")
		matching = uptime == uptimeCtrl+time.Second
		if matching {
			t.Logf("success! (uptime == uptimeCtrl+time.Second)")
		}
	}

	if !matching {
		t.Errorf("uptime does not match control uptime: %v != %v", uptime, uptimeCtrl)
	}
}

func TestSysinfo(t *testing.T) {
	t.Parallel()
	si, err := Sysinfo()
	if err != nil {
		t.Fatalf("failed to get sysinfo with error: %e", err)
	}
	if si == nil {
		t.Fatalf("failed to get sysinfo (nil)")
	}

	t.Logf("sysinfo: %s", spew.Sdump(si))
}
