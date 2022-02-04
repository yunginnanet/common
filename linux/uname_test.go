package linux

import "testing"

func TestGetUname(t *testing.T) {
	uname, err := GetUname("smrdnv")
	if err != nil {
		t.Fatalf("failed to get uname with error: %e", err)
	}
	if len(uname) < 1 {
		t.Fatalf("failed to get uname")
	} else {
		t.Logf("%s", uname)
	}
}
