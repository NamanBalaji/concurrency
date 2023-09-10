package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Throttle allows a limited number of workers to run at a time. It also
// provides a mechanism to check for errors encountered by workers and wait for
// them to finish.
//
// How to user throttle?
// -> Call throttle.Do before starting the goroutine.
// -> Call throttle.Done in the goroutine.
// -> Call throttle.Finish to wait for all routines to finish.
type Throttle struct {
	once      sync.Once
	wg        sync.WaitGroup
	ch        chan struct{}
	errCh     chan error
	finishErr error
}

// NewThrottle creates a new throttle with a max number of workers.
func NewThrottle(max int) *Throttle {
	return &Throttle{
		ch:    make(chan struct{}, max),
		errCh: make(chan error, max),
	}
}

// Do should be called by workers before they start working. It blocks if there
// are already maximum number of workers working. If it detects an error from
// previously Done workers, it would return it.
func (t *Throttle) Do() error {
	for {
		select {
		case t.ch <- struct{}{}:
			t.wg.Add(1)
			return nil
		case err := <-t.errCh:
			if err != nil {
				return err
			}
		}
	}
}

// Done should be called by workers when they finish working. They can also
// pass the error status of work done.
func (t *Throttle) Done(err error) {
	if err != nil {
		t.errCh <- err
	}
	select {
	case <-t.ch:
	default:
		panic("Throttle Do Done mismatch")
	}

	t.wg.Done()
}

// Finish waits until all workers have finished working. It would return any error passed by Done.
// If Finish is called multiple time, it will wait for workers to finish only once(first time).
// From next calls, it will return same error as found on first call.
func (t *Throttle) Finish() error {
	t.once.Do(func() {
		t.wg.Done()
		close(t.ch)
		close(t.errCh)
		for err := range t.errCh {
			if err != nil {
				t.finishErr = err
				return
			}
		}
	})

	return t.finishErr
}

func doWork(th *Throttle) {
	defer th.Done(nil)

	time.Sleep(time.Second)
}

func printStats() {
	for {
		time.Sleep(time.Millisecond * 500)
		fmt.Printf("Active number of goroutines = %+v\n", runtime.NumGoroutine())
	}
}

func main() {
	go printStats()

	th := NewThrottle(5)

	for i := 0; i < 100; i++ {
		th.Do()
		fmt.Printf("Processing item number = %+v\n", i)
		go doWork(th)
	}

	th.Finish()
}