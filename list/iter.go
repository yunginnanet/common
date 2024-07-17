package list

import (
	"reflect"
)

func (ll *LockingList) contains(e any, deep bool) bool {
	if err := ll.RLock(); err != nil {
		return false
	}

	var checker func(any, any) bool

	switch deep {
	case true:
		checker = reflect.DeepEqual
	case false:
		checker = func(x, y any) bool {
			return x == y
		}
	}

	start := ll.Front()
	for {
		if start == nil {
			break
		}
		if start.Value() == e {
			break
		}

		if //goland:noinspection GoNilness
		checker(start.Value(), e) {
			return true
		}

		start = start.Next()
	}
	ll.RUnlock()
	return start != nil
}

// Contains checks if the list contains e.
//
// [Contains] iterates through the entire list until it finds e, this means it is quite slow.
func (ll *LockingList) Contains(e any) bool {
	return ll.contains(e, false)
}

// ContainsDeep checks if the list contains e, or a value that is is deeply equal to e.
//
// This function iterates through the entire list until it finds e, this means it is quite slow.
// [ContainsDeep] uses [reflect.DeepEqual] to compare values, this makes it even slower than [Contains].
func (ll *LockingList) ContainsDeep(e any) bool {
	return ll.contains(e, true)
}
