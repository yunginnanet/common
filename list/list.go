// Package list implements a locking list.l
package list

import (
	"container/list"
	"errors"
	"sync"
)

var (
	ErrElementNotInList = errors.New("element not in list")
	ErrMarkNotInList    = errors.New("mark not in list")
	ErrNilValue         = errors.New("nil element")
	ErrUninitialized    = errors.New("uninitialized list")
)

type LockingList struct {
	l *list.List
	*sync.RWMutex
}

func (ll *LockingList) Lock() error {
	if ll == nil || ll.RWMutex == nil {
		return ErrUninitialized
	}
	ll.RWMutex.Lock()
	return nil
}

func (ll *LockingList) Unlock() {
	if ll.RWMutex != nil {
		ll.RWMutex.Unlock()
	}
}

func (ll *LockingList) RLock() error {
	switch {
	case
		ll == nil,
		ll.RWMutex == nil,
		ll.l == nil:
		return ErrUninitialized
	default:
		//
	}
	ll.RWMutex.RLock()
	return nil
}

func (ll *LockingList) RUnlock() {
	if ll.RWMutex != nil {
		ll.RWMutex.RUnlock()
	}
}

func New() *LockingList {
	ll := &LockingList{
		l:       list.New(),
		RWMutex: &sync.RWMutex{},
	}
	return ll
}

func (ll *LockingList) wrapElement(e *list.Element) *Element {
	if e == nil {
		return nil
	}
	return &Element{
		list:    ll,
		Element: e,
	}
}

// Init initializes or clears list l.
func (ll *LockingList) Init() *LockingList {
	if ll.l != nil {
		_ = ll.Lock()
		ll.l.Init()
		ll.Unlock()
		return ll
	}
	ll.l = list.New()
	ll.RWMutex = &sync.RWMutex{}
	return ll
}

func (ll *LockingList) InsertAfter(v any, mark *Element) (*Element, error) {
	if err := ll.check(v, mark, true); err != nil {
		return nil, err
	}
	_ = ll.Lock()
	res := ll.l.InsertAfter(v, mark.Element)
	ll.Unlock()
	return ll.wrapElement(res), nil
}

func (ll *LockingList) InsertBefore(v any, mark *Element) (*Element, error) {
	if err := ll.check(v, mark, true); err != nil {
		return nil, err
	}
	_ = ll.Lock()
	res := ll.wrapElement(
		ll.l.InsertBefore(v, mark.Element),
	)
	ll.Unlock()
	return res, nil
}

func (ll *LockingList) Len() int {
	_ = ll.RLock()
	l := ll.l.Len()
	ll.RUnlock()
	return l
}

func (ll *LockingList) MoveAfter(e, mark *Element) error {
	_ = ll.Lock()
	ll.l.MoveAfter(e.Element, mark.Element)
	ll.Unlock()
	return nil
}

func (ll *LockingList) Front() *Element {
	_ = ll.RLock()
	e := ll.l.Front()
	ll.RUnlock()
	return ll.wrapElement(e)
}

func (ll *LockingList) Back() *Element {
	_ = ll.RLock()
	e := ll.l.Back()
	ll.RUnlock()
	return ll.wrapElement(e)
}

func (ll *LockingList) MoveToBack(e *Element) error {
	_ = ll.Lock()
	ll.l.MoveToBack(e.Element)
	ll.Unlock()
	return nil
}

func (ll *LockingList) MoveToFront(e *Element) error {
	_ = ll.Lock()
	ll.l.MoveToFront(e.Element)
	ll.Unlock()
	return nil
}

func (ll *LockingList) PushBack(v any) *Element {
	if ll.l == nil {
		ll.Init()
	}
	_ = ll.Lock()
	e := ll.l.PushBack(v)
	ll.Unlock()
	return ll.wrapElement(e)
}

func (ll *LockingList) PushFront(v any) *Element {
	if ll.l == nil {
		ll.Init()
	}
	if err := ll.Lock(); err != nil {
		return nil
	}
	e := ll.l.PushFront(v)
	ll.Unlock()
	return ll.wrapElement(e)
}

func (ll *LockingList) Remove(elm *Element) error {
	if ll.l == nil {
		return ErrUninitialized
	}
	if err := ll.check(elm, nil, false); err != nil {
		return err
	}
	_ = ll.Lock()
	_ = ll.l.Remove(elm.Element)
	elm.list = nil    // avoid memory leaks
	elm.Element = nil // avoid memory leaks
	ll.Unlock()
	return nil
}

// Rotate moves the first element to the back of the list and returns it.
func (ll *LockingList) Rotate() *Element {
	if ll.l == nil {
		ll.Init()
	}
	if ll.Len() < 1 {
		return nil
	}
	_ = ll.Lock()
	e := ll.l.Front()
	ll.l.MoveToBack(e)
	ll.Unlock()
	return ll.wrapElement(e)
}

func (ll *LockingList) Push(item any) (err error) {
	if ll.l == nil {
		ll.Init()
	}
	if err = ll.Lock(); err != nil {
		return err
	}
	ll.l.PushBack(item)
	ll.Unlock()
	return nil
}

func (ll *LockingList) Pop() any {
	if ll.Len() < 1 {
		return nil
	}
	_ = ll.Lock()
	e := ll.l.Front()
	ll.l.Remove(e)
	ll.Unlock()
	return e.Value
}

func (ll *LockingList) PushBackList(other *LockingList) error {
	if ll.l == nil {
		ll.Init()
	}
	_ = ll.Lock()
	ll.l.PushBackList(other.l)
	ll.Unlock()
	return nil
}

func (ll *LockingList) PushFrontList(other *LockingList) error {
	if ll.l == nil {
		ll.Init()
	}
	_ = ll.Lock()
	ll.l.PushFrontList(other.l)
	ll.Unlock()
	return nil
}

type Element struct {
	*list.Element
	list *LockingList
}

func (e *Element) Value() any {
	if e == nil {
		return nil
	}
	if e.Element == nil {
		return nil
	}
	return e.Element.Value
}

func (e *Element) Next() *Element {
	if e == nil {
		return nil
	}
	if err := e.list.RLock(); err != nil {
		return nil
	}
	ne := e.list.wrapElement(e.Element.Next())
	e.list.RUnlock()
	return ne
}

func (e *Element) Prev() *Element {
	if e == nil {
		return nil
	}
	if err := e.list.RLock(); err != nil {
		return nil
	}
	pe := e.list.wrapElement(e.Element.Prev())
	e.list.RUnlock()
	return pe
}

func (ll *LockingList) check(item any, mark *Element, needsMark bool) (err error) {
	var elm *Element
	var isElement bool

	elm, isElement = item.(*Element)

	if err = ll.RLock(); err != nil {
		return err
	}
	switch {
	case
		needsMark && mark == nil,
		needsMark && mark.list != ll:
		err = ErrMarkNotInList

	case
		isElement && elm.Element == nil,
		isElement && elm.list != ll:
		err = ErrElementNotInList

	default:
		err = nil
	}
	ll.RUnlock()
	return
}
