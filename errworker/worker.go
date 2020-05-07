package errworker

import (
	"sync"
)

type ErrWorkgroup struct {
	limiterc      chan struct{}
	wg            sync.WaitGroup
	errMutex      sync.RWMutex
	err           error
	skipWhenError bool
}

func NewErrWorkgroup(size int, skipWhenError bool) ErrWorkgroup {
	if size < 1 {
		size = 1
	}
	return ErrWorkgroup{
		limiterc:      make(chan struct{}, size),
		skipWhenError: skipWhenError,
	}
}

// Wait waits till all current jobs finish and returns first occurred error
// in case something went wrong.
func (w *ErrWorkgroup) Wait() error {
	w.wg.Wait()
	return w.err
}

// Go adds work func with error to the ErrWorkgroup. If err occurred other jobs won't proceed.
func (w *ErrWorkgroup) Go(work func() error) {
	w.wg.Add(1)
	go func(fn func() error) {
		w.limiterc <- struct{}{}

		if w.skipWhenError {
			// if ErrWorkgroup corrupted -> skip work execution
			w.errMutex.RLock()
			if w.err == nil {
				w.errMutex.RUnlock()
				w.execute(fn)
			} else {
				w.errMutex.RUnlock()
			}
		} else {
			w.execute(fn)
		}

		w.wg.Done()
		<-w.limiterc
	}(work)
}

func (w *ErrWorkgroup) execute(work func() error) {
	if err := work(); err != nil {
		w.errMutex.Lock()
		w.err = err
		w.errMutex.Unlock()
	}
}
