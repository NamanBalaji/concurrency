package main

import "sync"

type WaitGroup struct {
	sync.Mutex
	n    int
	lock sync.Mutex
}

func (wg *WaitGroup) Add(delta int) {
	wg.Lock()
	defer wg.Unlock()

	if delta == 0 {
		return
	}

	if delta < 0 {
		panic("invalid delta")
	}

	if wg.n == 0 {
		wg.lock.Lock()
	}

	wg.n += delta
}

func (wg *WaitGroup) Done() {
	wg.Lock()
	defer wg.Unlock()

	if wg.n == 0 {
		panic("nagative n")
	}

	wg.n--
	if wg.n == 0 {
		wg.lock.Unlock()
	}
}

func (wg *WaitGroup) Wait() {
	wg.lock.Lock()
	wg.lock.Unlock()
}
