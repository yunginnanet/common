package list

import (
	"strings"
	"testing"
)

type testCommon interface {
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
}

func addGarbo(ll *LockingList, amount int, t testCommon) {
	t.Helper()
	for i := 0; i < amount; i++ {
		if i%2 == 0 {
			if err := ll.Push(strings.Repeat("e", i)); err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			continue
		}
		if err := ll.Push(i); err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
}

func TestContains(t *testing.T) {
	ll := New()
	if ll.Contains(1) {
		t.Error("expected 1 not to be in an empty list")
	}
	ll.PushFront(2)
	if ll.Contains(1) {
		t.Error("expected 1 not to be in the list")
	}
	if !ll.Contains(2) {
		t.Error("expected 2 to be in the list")
	}
	_ = ll.Pop()

	addGarbo(ll, 1000, t)

	for i := 0; i < 1000; i++ {
		switch i % 2 {
		case 0:
			if ll.Contains(i) {
				t.Fatalf("%d should not be in list", i)
			}
			if !ll.Contains(strings.Repeat("e", i)) {
				t.Fatal("missing string value")
			}
		default:
			if !ll.Contains(i) {
				t.Fatalf("%d should be in list", i)
			}
		}
	}

	t.Run("deep", func(t *testing.T) {
		type mcgee int

		type structy struct {
			yeeterson any
			yeetsalot string
		}

		if err := ll.Push(&structy{yeeterson: mcgee(1), yeetsalot: "yeet"}); err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		if !ll.ContainsDeep(&structy{yeeterson: mcgee(1), yeetsalot: "yeet"}) {
			t.Fatal("reflect.DeepEqual should have caught this")
		}

	})

	t.Run("list not initialized", func(t *testing.T) {
		var l = &LockingList{}
		if l.Contains(1) {
			t.Error("expected 1 not to be in an empty list")
		}
	})
}

/*

Wed Jul 17 04:20:42 PM PDT 2024

goos: linux
goarch: amd64
pkg: github.com/yunginnanet/common/list
cpu: 13th Gen Intel(R) Core(TM) i9-13900K
BenchmarkLockingList_Contains-32        	   6034	   196409 ns/op	 120801 B/op	   7550 allocs/op
BenchmarkLockingList_ContainsDeep-32    	   4717	   245886 ns/op	 120800 B/op	   7550 allocs/op
PASS
ok  	github.com/yunginnanet/common/list	16.810s

*/

func BenchmarkLockingList_Contains(b *testing.B) {
	b.StopTimer()
	ll := New()
	addGarbo(ll, 100, b)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for ii := 0; ii < 100; ii++ {
			b.StartTimer()
			_ = ll.Contains(ii)
			b.StopTimer()
		}
	}
}

func BenchmarkLockingList_ContainsDeep(b *testing.B) {
	b.StopTimer()
	ll := New()
	addGarbo(ll, 100, b)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for ii := 0; ii < 100; ii++ {
			b.StartTimer()
			_ = ll.ContainsDeep(ii)
			b.StopTimer()
		}
	}
}
