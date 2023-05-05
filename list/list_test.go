package list

import (
	"container/list"
	"errors"
	"reflect"
	"sync"
	"testing"
	"unsafe"
)

// Adapted from golang source
func checkListLen(t *testing.T, l *LockingList, length int) bool {
	t.Helper()
	if n := l.Len(); n != length {
		t.Errorf("l.Len() = %d, want %d", n, length)
		return false
	}
	return true
}

// Adapted from golang source
func checkList(t *testing.T, l *LockingList, es []any) {
	t.Helper()
	if !checkListLen(t, l, len(es)) {
		return
	}

	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		le := e.Value()
		if le != es[i] {
			t.Errorf("\uF630 elt[%d].Value = %v, want %v", i, le, es[i])
		}
		// t.Logf("\uF634 elt[%d].Value = %v, want %v", i, le, es[i])
		i++
	}
}

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func checkListPointers(t *testing.T, ll *LockingList, es []*list.Element) {
	if !checkListLen(t, ll, len(es)) {
		return
	}

	newList := reflect.New(reflect.TypeOf(*ll.l)).Elem()
	root := newList.FieldByName("root")

	// zero length lists must be the zero value or properly initialized (sentinel circle)
	if len(es) == 0 {
		next := root.FieldByName("next")
		prev := root.FieldByName("prev")
		if !next.IsNil() && next != root || !prev.IsNil() && prev != root {
			t.Errorf("l.root.next = %v; should be nil or %v",
				next, root.Type(),
			)
		}
		if !prev.IsNil() && next != root || !prev.IsNil() && prev != root {
			t.Errorf("l.root.prev = %p; should be nil or %v",
				prev.Type(), root.Type(),
			)
		}
		return
	}
}

// Adapted from golang source
func TestExtending(t *testing.T) { //nolint:funlen,gocyclo
	t.Parallel()
	l1 := New()
	l2 := New()

	l1.PushBack(1)
	l1.PushBack(2)
	l1.PushBack(3)

	l2.PushBack(4)
	l2.PushBack(5)

	l3 := New()
	if err := l3.PushBackList(l1); err != nil {
		t.Error(errors.New("PushBackList failed"))
	}
	checkList(t, l3, []any{1, 2, 3})
	if err := l3.PushBackList(l2); err != nil {
		t.Error(errors.New("PushBackList failed"))
	}
	checkList(t, l3, []any{1, 2, 3, 4, 5})

	l3 = New()
	if err := l3.PushFrontList(l2); err != nil {
		t.Error(errors.New("PushFrontList failed"))
	}
	checkList(t, l3, []any{4, 5})
	if err := l3.PushFrontList(l1); err != nil {
		t.Error(errors.New("PushFrontList failed"))
	}
	checkList(t, l3, []any{1, 2, 3, 4, 5})

	checkList(t, l1, []any{1, 2, 3})
	checkList(t, l2, []any{4, 5})

	l3 = New()
	if err := l3.PushBackList(l1); err != nil {
		t.Error(errors.New("PushBackList failed"))
	}
	checkList(t, l3, []any{1, 2, 3})
	if err := l3.PushBackList(l3); err != nil {
		t.Error(errors.New("PushBackList failed"))
	}
	checkList(t, l3, []any{1, 2, 3, 1, 2, 3})

	l3 = New()
	if err := l3.PushFrontList(l1); err != nil {
		return
	}
	checkList(t, l3, []any{1, 2, 3})
	if err := l3.PushFrontList(l3); err != nil {
		t.Error(errors.New("PushFrontList failed"))
	}
	checkList(t, l3, []any{1, 2, 3, 1, 2, 3})

	l3 = New()
	if err := l1.PushBackList(l3); err != nil {
		t.Error(errors.New("PushBackList failed"))
	}
	checkList(t, l1, []any{1, 2, 3})
	if err := l1.PushFrontList(l3); err != nil {
		t.Error(errors.New("PushFrontList failed"))
	}
	checkList(t, l1, []any{1, 2, 3})
}

// Adapted from golang source
func TestIssue4103(t *testing.T) {
	t.Parallel()
	l1 := New()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := New()
	l2.PushBack(3)
	l2.PushBack(4)

	e := l1.Front()
	err := l2.Remove(e)
	if err == nil {
		t.Errorf("l2.Remove(e) = %v, want ErrElementNotInList", err)
	}
	if !errors.Is(err, ErrElementNotInList) {
		t.Errorf("l2.Remove(e) = %v, want ErrElementNotInList", err)
	}

	// l2 should not change because e is not an element of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	var ne *Element

	if ne, err = l1.InsertBefore(8, e); err != nil {
		t.Errorf("l1.InsertBefore(8, e) = %v, want nil", err)
	}

	//goland:noinspection GoCommentLeadingSpace (lol)
	if ne == nil { // nolint:SA5011
		t.Fatalf("l1.InsertBefore(8, e) = nil, want non-nil")
	}

	//goland:noinspection GoCommentLeadingSpace (lol)
	if ne.Element == nil { // nolint:SA5011
		t.Errorf("l1.InsertBefore(8, e) = nil, want non-nil")
	}
	if ne.Value() != 8 {
		t.Errorf("l1.InsertBefore(8, e) = %v, want 8", ne.Value())
	}

	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

// Adapted from golang source
func TestIssue6349(t *testing.T) {
	t.Parallel()
	l := New()
	l.PushBack(1)
	l.PushBack(2)

	e := l.Front()
	if e.Value() != 1 {
		t.Errorf("e.value = %d, want 1", e.Value())
	}
	if err := l.Remove(e); err != nil {
		t.Errorf("l.Remove(e) = %v, want nil", err)
	}
	if e.Next() != nil {
		t.Errorf("e.Next() != nil")
	}
	if e.Prev() != nil {
		t.Errorf("e.Prev() != nil")
	}
}

// Test PushFront, PushBack, PushFrontList, PushBackList with uninitialized List
// Adapted from golang source.
func TestZeroList(t *testing.T) {
	t.Parallel()
	var l1 = new(LockingList)
	l1.PushFront(1)
	checkList(t, l1, []any{1})

	var l2 = new(LockingList)
	l2.PushBack(1)
	checkList(t, l2, []any{1})

	var l3 = new(LockingList)
	if err := l3.PushFrontList(l1); err != nil {
		t.Errorf("l3.PushFrontList(l1) = %v, want nil", err)
	}
	checkList(t, l3, []any{1})

	var l4 = new(LockingList)
	if err := l4.PushBackList(l2); err != nil {
		t.Errorf("l4.PushBackList(l2) = %v, want nil", err)
	}
	checkList(t, l4, []any{1})
}

// Test that a list l is not modified when calling InsertBefore with a mark that is not an element of l.
// Adapted from golang source.
func TestInsertBeforeUnknownMark(t *testing.T) {
	t.Parallel()
	var l = New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	_, err := l.InsertBefore(1, new(Element))
	if err == nil || !errors.Is(err, ErrMarkNotInList) {
		t.Errorf("l.InsertBefore(1, new(Element)) = %v, want ErrMarkNotInList", err)
	}
	checkList(t, l, []any{1, 2, 3})
}

// TestList tests the list implementation.
// Mostly adapted from golang source
func TestList(t *testing.T) {
	l := New()
	checkListPointers(t, l, []*list.Element{})

	// Single element list
	e := l.PushFront("a")
	checkListPointers(t, l, []*list.Element{e.Element})
	if err := l.MoveToFront(e); err != nil {
		t.Errorf("MoveToFront(e) = %v, want nil", err)
	}
	if err := l.MoveAfter(e, e); err != nil {
		t.Errorf("MoveAfter(e, e) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e.Element})
	if err := l.MoveToBack(e); err != nil {
		t.Errorf("MoveToBack(e) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e.Element})
	if err := l.Remove(e); err != nil {
		t.Errorf("Remove(e) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{})

	// Bigger list
	e2 := l.PushFront(2)
	e1 := l.PushFront(1)
	e3 := l.PushBack(3)
	e4 := l.PushBack("banana")
	checkListPointers(t, l, []*list.Element{e1.Element, e2.Element, e3.Element, e4.Element})

	var err error

	if err = l.Remove(e2); err != nil {
		t.Errorf("Remove(e2) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e1.Element, e3.Element, e4.Element})

	// move from middle
	if err = l.MoveToFront(e3); err != nil {
		t.Errorf("MoveToFront(e3) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e3.Element, e1.Element, e4.Element})

	if err = l.MoveToFront(e1); err != nil {
		t.Errorf("MoveToFront(e1) = %v, want nil", err)
	}

	// move from middle
	if err = l.MoveToBack(e3); err != nil {
		t.Errorf("MoveToBack(e3) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e1.Element, e4.Element, e3.Element})

	// move from back
	if err = l.MoveToFront(e3); err != nil {
		t.Errorf("MoveToFront(e3) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e3.Element, e1.Element, e4.Element})

	// should be no-op
	if err = l.MoveToFront(e3); err != nil {
		t.Errorf("MoveToFront(e3) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e3.Element, e1.Element, e4.Element})

	// move from front
	if err = l.MoveToBack(e3); err != nil {
		t.Errorf("MoveToBack(e3) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e1.Element, e4.Element, e3.Element})

	// should be no-op
	if err = l.MoveToBack(e3); err != nil {
		t.Errorf("MoveToBack(e3) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e1.Element, e4.Element, e3.Element})

	// insert before front
	if e2, err = l.InsertBefore(2, e1); err != nil {
		t.Errorf("InsertBefore(2, e1) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e2.Element, e1.Element, e4.Element, e3.Element})
	if err = l.Remove(e2); err != nil {
		t.Errorf("Remove(e2) = %v, want nil", err)
	}

	// insert before middle
	if e2, err = l.InsertBefore(2, e4); err != nil {
		t.Errorf("InsertBefore(2, e4) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e1.Element, e2.Element, e4.Element, e3.Element})
	if err = l.Remove(e2); err != nil {
		t.Errorf("Remove(e2) = %v, want nil", err)
	}

	// insert before back
	if e2, err = l.InsertBefore(2, e3); err != nil {
		t.Errorf("InsertBefore(2, e3) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e1.Element, e4.Element, e2.Element, e3.Element})
	if err = l.Remove(e2); err != nil {
		t.Errorf("Remove(e2) = %v, want nil", err)
	}

	// insert after front
	if e2, err = l.InsertAfter(2, e1); err != nil {
		t.Errorf("InsertAfter(2, e1) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e1.Element, e2.Element, e4.Element, e3.Element})
	if err = l.Remove(e2); err != nil {
		t.Errorf("Remove(e2) = %v, want nil", err)
	}

	// insert after middle
	if e2, err = l.InsertAfter(2, e4); err != nil {
		t.Errorf("InsertAfter(2, e4) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e1.Element, e4.Element, e2.Element, e3.Element})
	if err = l.Remove(e2); err != nil {
		t.Errorf("Remove(e2) = %v, want nil", err)
	}

	// insert after back
	if e2, err = l.InsertAfter(2, e3); err != nil {
		t.Errorf("InsertAfter(2, e3) = %v, want nil", err)
	}
	checkListPointers(t, l, []*list.Element{e1.Element, e4.Element, e3.Element, e2.Element})
	if err = l.Remove(e2); err != nil {
		t.Errorf("Remove(e2) = %v, want nil", err)
	}

	// Check standard iteration.
	sum := 0
	for e = l.Front(); e != nil; e = e.Next() {
		if i, ok := e.Element.Value.(int); ok {
			sum += i
		}
	}
	if sum != 4 {
		t.Errorf("sum over l = %d, want 4", sum)
	}
}

func TestThePlanet(t *testing.T) {
	t.Parallel()
	var err error
	var nl = New()

	if nl.Len() != 0 {
		t.Errorf("Init() failed to reset list length to 0")
	}

	nl.PushFront(1)
	if nl.Pop() != 1 {
		t.Errorf("Pop() failed to return first element")
	}
	if nl.Len() != 0 {
		t.Errorf("Pop() failed to remove first element")
	}

	if err = nl.Push(1); err != nil {
		t.Errorf("Push(1) = %v, want nil", err)
	}

	nl.l = nil
	if err = nl.Push(1); err != nil {
		t.Errorf("Push(1) = %v, want %v", err, nil)
	}
	if nl.Pop() != 1 {
		t.Errorf("Pop() = %v, want %v", err, 1)
	}
	if nl.Pop() != nil {
		t.Errorf("Pop() = %v, want %v", err, nil)
	}

	t.Run("PushPop", func(t *testing.T) {
		t.Parallel()
		pl := New()
		for i := 1; i < 5; i++ {
			if err = pl.Push(i); err != nil {
				t.Errorf("Push(%d) = %v, want nil", i, err)
			}
		}
		for i := 1; i < 5; i++ {
			if got := pl.Pop(); got != i {
				t.Errorf("Pop() = %d, want %d", got, i)
			}
		}
		for i := 1; i < 5; i++ {
			if err = pl.Push(i); err != nil {
				t.Errorf("Push(%d) = %v, want nil", i, err)
			}
		}
	})

	t.Run("Concurrent", func(t *testing.T) {
		t.Parallel()
		cl := New()
		wg := &sync.WaitGroup{}
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(n int) {
				if cerr := cl.Push(n); cerr != nil {
					t.Errorf("Push(%d) err = %v, want nil", n, cerr)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		if cl.Len() != 1000 {
			t.Errorf("Len() = %d, want 1000", cl.Len())
		}
		for i := 0; i < 1000; i++ {
			if cl.Pop() == nil {
				t.Errorf("Pop() = %d, want %d", cl.Pop(), i)
			}
			if cl.Len() != 999-i {
				t.Errorf("Len() = %d, want %d", cl.Len(), 999-i)
			}
		}
	})

	nl.Init()
	for i := 1; i < 5; i++ {
		if err = nl.Push(i); err != nil {
			t.Errorf("Push(%d) = %v, want nil", i, err)
		}
	}

	for i := 1; i < 100; i++ {
		for j := 1; j < 5; j++ {
			res := nl.Rotate()
			t.Logf("Rotate() = %d", res.Value())
			if res.Value() != j {
				t.Errorf("Rotate() = %d, want %d", res.Value(), j)
			}
			if j > 4 && res.Next() == nil {
				t.Fatalf("%d.Next is nil", res.Value())
			}
			if j < 2 && res.Prev() == nil {
				t.Errorf("Prev is nil")
			}
			if res.Next() != nil && res.Next().Value() != j+1 {
				t.Errorf("Next = %d, want %d", res.Next().Value(), j+1)
			}
			if res.Prev() != nil && j-1 != 0 && res.Prev().Value() != j-1 {
				t.Errorf("Prev = %d, want %d", res.Prev().Value(), j-1)
			}
		}
	}

	t.Run("CoverageStuff", func(t *testing.T) {
		nl.RWMutex = nil
		if !errors.Is(nl.Lock(), ErrUninitialized) {
			t.Errorf("Lock() = %v, want %v", err, ErrUninitialized)
		}
		var yeet *Element
		if yeet, err = nl.InsertAfter(1, nil); !errors.Is(err, ErrUninitialized) {
			t.Errorf("InsertAfter(1, nil) = %v, want %v", err, ErrUninitialized)
		}
		if yeet.Next() != nil {
			t.Errorf("Next() = %v, want %v", yeet.Next(), nil)
		}
		if yeet.Prev() != nil {
			t.Errorf("Prev() = %v, want %v", yeet.Prev(), nil)
		}
		if _, err = nl.InsertBefore(1, nil); !errors.Is(err, ErrUninitialized) {
			t.Errorf("InsertAfter(1, nil) = %v, want %v", err, ErrUninitialized)
		}
		if err = nl.Remove(nil); !errors.Is(err, ErrUninitialized) {
			t.Errorf("Remove(nil) = %v, want %v", err, ErrUninitialized)
		}
		if err = nl.Push(1); !errors.Is(err, ErrUninitialized) {
			t.Errorf("Push(1) = %v, want %v", err, ErrUninitialized)
		}
		if e := nl.PushFront(1); e != nil {
			t.Errorf("PushFront(1) = %v, want %v", err, ErrUninitialized)
		}
		nl.l = nil
		nl.Rotate()
		if nl.wrapElement(nil) != nil {
			t.Errorf("wrapElement(e) = %v, want nil", err)
		}
		if nl.wrapElement(nil).Value() != nil {
			t.Errorf("wrapElement(e).Value() = %v, want nil", err)
		}
		el := &Element{}
		if el.Value() != nil {
			t.Errorf("el.Value() = %v, want nil", err)
		}
		if _, err = nl.InsertAfter(1, nil); err == nil {
			t.Errorf("InsertAfter(1, nil) = %v, want error", err)
		}
		nl.Init()
		nl.Init()
		nl.l = nil
		if !errors.Is(nl.Remove(nil), ErrUninitialized) {
			t.Errorf("Remove(nil) = %v, want %v", err, ErrUninitialized)
		}
		nl.Init()
		_ = nl.Push(1)
		_ = nl.Push(2)
		if nl.Back().Value() == 1 {
			t.Errorf("Back() = %v, want %v", err, 1)
		}
	})

	// Clear all elements by iterating
	for nl.Len() > 0 {
		nl.Pop()
	}
	checkListPointers(t, nl, []*list.Element{})
}
