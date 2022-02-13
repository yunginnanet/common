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

func TestGetUnameFailure(t *testing.T) {
	uname, err := GetUname("frickhole johnson")
	if err == nil {
		t.Fatalf("[FAIL] We failed to fail. Wanted an error. %e", err)
	}
	if len(uname) > 1 {
		t.Fatalf("[FAIL] Despite erroring out we still received a value: %v", uname)
	}
}
