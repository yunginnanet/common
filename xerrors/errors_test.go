package xerrors

import (
	"errors"
	"io"
	"io/fs"
	"sync"
	"testing"
)

func TestNewErrorsIsEmpty(t *testing.T) {
	e := NewErrors()
	if len(e.errs) != 0 {
		t.Errorf("expected empty error stack, got %d", len(e.errs))
	}
}

func TestAddErrorIncreasesLength(t *testing.T) {
	e := NewErrors()
	err := errors.New("test error")
	e.Push(err)
	if len(e.errs) != 1 {
		t.Errorf("expected error stack length to be 1, got %d", len(e.errs))
		for _, existing := range e.errs {
			println(existing.Error())
		}
	}
}

func TestPopOneReturnsLastAddedError(t *testing.T) {
	e := NewErrors()
	firstErr := errors.New("first error")
	secondErr := errors.New("second error")
	e.Push(nil)
	e.Push(firstErr)
	e.Push(secondErr)
	e.Push(nil)
	poppedErr := e.Pop()
	if !errors.Is(poppedErr, secondErr) {
		t.Log(poppedErr.Error())
		t.Errorf("expected to pop the first added error")
	}
}

func TestPopAllReturnsAllErrors(t *testing.T) {
	e := NewErrors()
	e.Push(errors.New("first error"))
	e.Push(errors.New("second error"))
	errs := e.PopAll()
	if len(errs) != 2 {
		t.Errorf("expected to pop all errors, got %d", len(errs))
	}
}

func TestImmutableWriteOpShouldPanic(t *testing.T) {
	t.Run("Push", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Push should panic when called on an immutable error stack")
			}
		}()
		e := NewErrors()
		e.Push(errors.New("test error"))
		ec := e.Copy()
		ec.e.Push(errors.New("new error"))
	})
	t.Run("Pop", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Pop should panic when called on an immutable error stack")
			}
		}()
		e := NewErrors()
		e.Push(errors.New("test error"))
		ec := e.Copy()
		_ = ec.e.Pop()
	})
	t.Run("PopAll", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("PopAll should panic when called on an immutable error stack")
			}
		}()
		e := NewErrors()
		e.Push(errors.New("test error"))
		ec := e.Copy()
		_ = ec.e.PopAll()
	})
	t.Run("Clear", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Clear should panic when called on an immutable error stack")
			}
		}()
		e := NewErrors()
		e.Push(errors.New("test error"))
		ec := e.Copy()
		ec.e.Clear()
	})
	t.Run("PopAllImmutable", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("PopAllImmutable should panic when called on an immutable error stack")
			}
		}()
		e := NewErrors()
		e.Push(errors.New("test error"))
		ec := e.Copy()
		_ = ec.e.PopAllImmutable()
	})
}

func TestClearEmptiesErrorStack(t *testing.T) {
	e := NewErrors()
	e.Push(errors.New("test error"))
	e.Clear()
	if len(e.errs) != 0 {
		t.Errorf("expected error stack to be empty after clear, got %d", len(e.errs))
	}
}

func TestLenReturnsCorrectNumberOfErrors(t *testing.T) {
	e := NewErrors()
	e.Push(errors.New("first error"))
	e.Push(errors.New("second error"))
	if e.Len() != 2 {
		t.Errorf("expected length to be 2, got %d", e.Len())
	}
}

func TestIsFindsError(t *testing.T) {
	targetErr := errors.New("target error")
	e := NewErrors()
	e.Push(targetErr)
	if !e.Is(targetErr) {
		t.Errorf("expected to find the target error")
	}
	if !errors.Is(e.Pop(), targetErr) {
		t.Fatal("expected to find the target error on pop")
	}
	if !e.Is(nil) {
		t.Errorf("expected Is(nil) to be true after popping only error from stack")
	}
	if e.Is(targetErr) {
		t.Errorf("expected to not find the target error after popping")
	}
	newErr := errors.New("new error")
	e.Push(newErr)
	if e.Is(targetErr) {
		t.Errorf("expected to not find the target error after pushing a new error")
	}
	if !e.Is(newErr) {
		t.Errorf("expected to find the new error")
	}
}

func TestAsFindsErrorType(t *testing.T) {
	targetErr := &fs.PathError{
		Op:   "yeet",
		Path: "/yeet",
		Err:  errors.New("yeet error"),
	}
	e := NewErrors()
	e.Push(targetErr)
	var result *fs.PathError
	if !errors.As(e.errs[0], &result) {
		t.Fatal("expected to find the target error")
	}
	if result.Op != targetErr.Op || result.Path != targetErr.Path || result.Err.Error() != targetErr.Err.Error() {
		t.Errorf("errors.As did not fill the target variable with the correct error")
	}
	result = &fs.PathError{}
	if !e.As(&result) {
		t.Fatal("expected to find the target error")
	}
	if result.Op != targetErr.Op || result.Path != targetErr.Path || result.Err.Error() != targetErr.Err.Error() {
		t.Errorf("e.As did not fill the target variable with the correct error")
	}
	e2 := e.PopAllImmutable()
	if !e2.As(&result) {
		t.Fatal("expected to find the target error")
	}
	if result.Op != targetErr.Op || result.Path != targetErr.Path || result.Err.Error() != targetErr.Err.Error() {
		t.Errorf("e2.As did not fill the target variable with the correct error")
	}
	if e.As(&result) {
		t.Fatal("expected not to find the target error after PopAllImmutable")
	}

}

func TestConcatCombinesErrors(t *testing.T) {
	e := NewErrors()
	e.Push(errors.New("first error"))
	e.Push(errors.New("second error"))
	concatErr := e.Concat()
	if concatErr == nil {
		t.Fatal("expected concatenated error to be non-nil")
	}
	expected := "first error\nsecond error"
	if concatErr.Error() != expected {
		t.Errorf("expected concatenated error to be '%s', got '%s'", expected, concatErr.Error())
	}
}

func TestPopReturnsAndClearsErrors(t *testing.T) {
	e := NewErrors()
	e.Push(errors.New("error"))
	popped := e.PopAllImmutable()
	if popped.Len() != 1 || len(e.errs) != 0 {
		t.Errorf("PopAllImmutable should return errors and clear the original stack")
	}
}

func TestPushAddsErrorToNonEmptyStack(t *testing.T) {
	e := NewErrors()
	e.Push(errors.New("existing error"))
	e.Push(errors.New("new error"))
	if len(e.errs) != 2 {
		t.Errorf("Expected stack length to be 2, got %d", len(e.errs))
	}
}

func TestPushAddsErrorToEmptyStack(t *testing.T) {
	e := NewErrors()
	e.Push(errors.New("new error"))
	if len(e.errs) != 1 {
		t.Errorf("Expected stack length to be 1, got %d", len(e.errs))
	}
}

func TestPopFromNonEmptyStack(t *testing.T) {
	e := NewErrors()
	err := errors.New("test error")
	e.Push(err)
	poppedErr := e.Pop()
	if !errors.Is(poppedErr, err) {
		t.Errorf("Expected to pop the pushed error")
	}
}

func TestPopFromEmptyStack(t *testing.T) {
	e := NewErrors()
	poppedErr := e.Pop()
	if poppedErr != nil {
		t.Errorf("Expected to pop nil from an empty stack")
	}
}

func TestImmutableCopyPreventsModification(t *testing.T) {
	e := NewErrors()
	e.Push(errors.New("test error"))
	ec := e.Copy()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when modifying immutable copy")
		}
	}()
	ec.e.Push(errors.New("should not add"))
}

func TestPopAllRemovesAllErrors(t *testing.T) {
	e := NewErrors()
	e.Push(errors.New("first error"))
	e.Push(errors.New("second error"))
	poppedErrs := e.PopAll()
	if len(poppedErrs) != 2 || len(e.errs) != 0 {
		t.Errorf("Expected to pop all errors and empty the stack")
	}
}

func TestConcatWithMultipleErrors(t *testing.T) {
	e := NewErrors()
	e.Push(errors.New("first error"))
	e.Push(errors.New("second error"))
	concatErr := e.Concat()
	if concatErr == nil {
		t.Fatal("Expected concatenated error to be non-nil")
	}
	expected := "first error\nsecond error" // Assuming errors.Join concatenates in reverse order
	if concatErr.Error() != expected {
		t.Errorf("Expected concatenated error to be '%s', got '%s'", expected, concatErr.Error())
	}
}

func TestErrorsReturnsCopyOfAllErrors(t *testing.T) {
	e := NewErrors()
	firstErr := errors.New("first error")
	secondErr := errors.New("second error")
	e.Push(firstErr)
	e.Push(secondErr)
	errs := e.Errors()
	if len(errs) != 2 || !errors.Is(errs[0], secondErr) || !errors.Is(errs[1], firstErr) {
		t.Errorf("Expected Errors to return a copy of all errors in reverse order")
	}
	t.Run("ErrorsImmutable", func(t *testing.T) {
		immutable := e.Copy()
		if len(immutable.Errors()) != 2 {
			t.Fatal("Expected immutable copy to return all errors")
		}
		if !immutable.Is(firstErr) || !immutable.Is(secondErr) {
			t.Fatal("Expected immutable copy to contain all errors")
		}
		if immutable.Error() == "" {
			t.Fatal("Expected immutable copy to return non-empty string")
		}
		if immutable.Error() != e.Error() {
			t.Fatal("Expected immutable Error method to return the same string as the original")
		}
	})
}

func TestErrorImplementsErrorInterface(t *testing.T) {
	e := NewErrors()
	e.Push(errors.New("error"))
	if e.Error() == "" {
		t.Errorf("Expected Error to implement error interface and return non-empty string")
	}
}

func testSeeking(e ErrorStack, t *testing.T) {
	seekOnce := &sync.Once{}

	doGoto := false

iter:
	i := 1

	name := "First Iteration"

	if doGoto {
		name = "Second Iteration"
	}

	t.Run(name, func(t *testing.T) {
		for ee := e.Next(); ee != nil; ee = e.Next() {
			switch i {
			case 1:
				if ee.Error() != "first error" {
					t.Errorf("Expected first error, got %s", ee.Error())
				}
			case 2:
				if ee.Error() != "second error" {
					t.Errorf("Expected second error, got %s", ee.Error())
				}
			case 3:
				if ee.Error() != "third error" {
					t.Errorf("Expected third error, got %s", ee.Error())
				}
			case 4:
				if ee.Error() != "fourth error" {
					t.Errorf("Expected fourth error, got %s", ee.Error())
				}
			case 5:
				if ee.Error() != "fifth error" {
					t.Errorf("Expected fifth error, got %s", ee.Error())
				}
			case 6:
				t.Fatal("Expected to break after fifth error")
			}
			i++
		}
	})

	doGoto = false
	seekOnce.Do(func() {
		if n, err := e.Seek(0, io.SeekStart); n != 0 || err != nil {
			t.Fatalf("Expected to seek to beginning of stack, got %d, %v", n, err)
		}
		doGoto = true
	})
	if doGoto {
		goto iter
	}

	t.Run("Seek(n,io.SeekStart)", func(t *testing.T) {
		if n, err := e.Seek(1, io.SeekStart); n != 1 || err != nil {
			t.Fatalf("Expected to seek to second error, got %d, %v", n, err)
		}
	})

	t.Run("Seek(n,io.SeekCurrent)", func(t *testing.T) {
		if n, err := e.Seek(1, io.SeekCurrent); n != 2 || err != nil {
			t.Fatalf("Expected to seek to third error, got %d, %v", n, err)
		}
	})

	t.Run("Seek(n,io.SeekEnd)", func(t *testing.T) {
		if n, err := e.Seek(-1, io.SeekEnd); n != 4 || err != nil {
			t.Fatalf("Expected to seek to fourth error, got %d, %v", n, err)
		}
		if n, err := e.Seek(1, io.SeekEnd); n != 0 || err == nil {
			if err == nil {
				t.Fatalf("Expected seeking past end to return EOF, got %d, %v", n, err)
			}
			t.Fatalf("Expected seeking past end to return 0, got %d, %v", n, err)
		}
	})

}

func TestErrorsNextAndSeek(t *testing.T) {
	e := NewErrors()
	e.Push(errors.New("first error"))
	e.Push(errors.New("second error"))
	e.Push(errors.New("third error"))
	e.Push(errors.New("fourth error"))
	e.Push(errors.New("fifth error"))
	testSeeking(e, t)

	ec := e.PopAllImmutable()
	testSeeking(ec, t)
}
