/*
Package cleaningRegister provides a map that periodically cleans up entries.

*/

package cleaningRegister

import (
	"runtime"
	"sync"
	"time"
)

type Removed func(interface{}, interface{})
type ShouldRemove func(interface{}, interface{}) bool

// GC trick
// https://groups.google.com/forum/#!topic/golang-nuts/1ItNOOj8yW8/discussion
type cleaningRegister struct {
	register     map[interface{}]interface{}
	shouldRemove ShouldRemove
	removed      Removed
	mutex        sync.RWMutex
}

type CleaningRegister struct {
	*cleaningRegister
}

func New(cleanupInterval time.Duration, shouldRemove ShouldRemove, removed Removed) *CleaningRegister {
	reg := &cleaningRegister{
		register:     make(map[interface{}]interface{}),
		shouldRemove: shouldRemove,
		removed:      removed,
	}

	register := CleaningRegister{reg}

	cleanUp := make(chan bool)

	go func() {
		c := time.Tick(cleanupInterval)
		for {
			select {
			case _ = <-c:
				// Reference the internal struct
				reg.cleanUp()
			case _ = <-cleanUp:
				return
			}

		}
	}()

	runtime.SetFinalizer(&register, func(reg interface{}) { close(cleanUp) })

	return &register
}

func (r cleaningRegister) cleanUp() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for k, v := range r.register {
		if r.shouldRemove == nil || r.shouldRemove(k, v) {
			delete(r.register, k)
			if r.removed != nil {
				r.removed(k, v)
			}
		}
	}

}

func (r cleaningRegister) Get(key interface{}) (interface{}, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	val, ok := r.register[key]
	return val, ok
}

func (r cleaningRegister) Copy() map[interface{}]interface{} {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	values := make(map[interface{}]interface{})
	for k, v := range r.register {
		values[k] = v
	}

	return values
}

func (r cleaningRegister) Put(key interface{}, value interface{}) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.register[key] = value
}

func (r cleaningRegister) Pop(key interface{}) (interface{}, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	val, ok := r.register[key]
	if ok {
		delete(r.register, key)
	}
	return val, ok
}
