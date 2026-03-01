// Package init provides sync.Once alternatives for initialisers.
//
// Once is similar to sync.Once except it returns an error natively so that there's no need to separately set an error.
//
// SuccessOnce allows the initialiser to be run multiple times (but not concurrently) until the first success.
package init

import (
	"sync"
	"sync/atomic"
)

// Once is similar to sync.Once: the initializer is run exactly once, and any error, nil or non-nil, is always returned
// in subsequent invocations.
type Once struct {
	_    noCopy
	once sync.Once
	err  error
}

// Do calls the function f exactly once as if sync.Once has been called.
func (o *Once) Do(f func() error) error {
	o.once.Do(func() {
		o.err = f()
	})
	return o.err
}

// SuccessOnce allows the initializer to run as many times as it needs until it returns a non-nil error.
type SuccessOnce struct {
	_ noCopy
	// done and m copies logic of sync.Once except done is only set on success.
	done atomic.Bool
	m    sync.Mutex
}

// Do calls f as many times as it needs until the first invocation returns nil error.
func (o *SuccessOnce) Do(f func() error) error {
	if !o.done.Load() {
		return o.doSlow(f)
	}

	return nil
}

func (o *SuccessOnce) doSlow(f func() error) (err error) {
	o.m.Lock()
	defer o.m.Unlock()
	if !o.done.Load() {
		if err = f(); err == nil {
			o.done.Store(true)
		}
	}

	return
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
