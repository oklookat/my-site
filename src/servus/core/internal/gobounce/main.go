package gobounce

// fork of https://github.com/bep/debounce (MIT)

import (
	"sync"
	"time"
)

// create new debouncer.
//
// withSkips = true: run now() -> wait then run() -> run now() -> wait then run()
//
// withSkips = false: wait then run() -> wait then run() -> wait then run() -> wait then run()
func New(after time.Duration, withSkips bool) func(f func()) {
	d := &debouncer{after: after}

	if withSkips {
		return func(f func()) {
			d.addWithSkip(f)
		}
	}

	return func(f func()) {
		d.add(f)
	}
}

type debouncer struct {
	mu             sync.Mutex
	after          time.Duration
	timer          *time.Timer
	isCalledBefore bool
}

func (d *debouncer) add(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		return
	}

	d.timer = time.AfterFunc(d.after, func() {
		f()
		d.timer = nil
	})
}

func (d *debouncer) addWithSkip(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		return
	}

	if !d.isCalledBefore {
		f()
		d.isCalledBefore = true
		return
	}

	d.timer = time.AfterFunc(d.after, func() {
		f()
		d.timer = nil
		d.isCalledBefore = false
	})
}
