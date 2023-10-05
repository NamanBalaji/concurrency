package main

import (
	"fmt"
	"sync"
	"time"
)

var itemsToShip = []item{
	item{0, "Shirt", 1 * time.Second},
	item{1, "Legos", 1 * time.Second},
	item{2, "TV", 5 * time.Second},
	item{3, "Bananas", 2 * time.Second},
	item{4, "Hat", 1 * time.Second},
	item{5, "Phone", 2 * time.Second},
	item{6, "Plates", 3 * time.Second},
	item{7, "Computer", 5 * time.Second},
	item{8, "Pint Glass", 3 * time.Second},
	item{9, "Watch", 2 * time.Second},
}

type item struct {
	id     int
	name   string
	effort time.Duration
}

func prepareItems(done <-chan bool) <-chan item {
	items := make(chan item)

	go func() {
		for _, item := range itemsToShip {
			select {
			case <-done:
				return
			case items <- item:
			}
		}
		close(items)
	}()

	return items
}

func packItems(done <-chan bool, items <-chan item, workerId int) <-chan int {
	packages := make(chan int)
	go func() {
		for item := range items {
			select {
			case <-done:
				return
			case packages <- item.id:
				time.Sleep(item.effort)
				fmt.Printf("Worker #%d: Shipping package no. %d, took %ds to pack\n", workerId, item.id, item.effort/time.Second)
			}
		}
		close(packages)
	}()

	return packages
}

func merge(done <-chan bool, chans ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	wg.Add(len(chans))

	outgoingPackages := make(chan int)
	multiplex := func(c <-chan int) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case outgoingPackages <- i:
			}
		}
	}

	for _, c := range chans {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(outgoingPackages)
	}()

	return outgoingPackages
}

func main() {
	done := make(chan bool)
	defer close(done)

	start := time.Now()

	items := prepareItems(done)

	workers := make([]<-chan int, 4)

	for i := 0; i < 4; i++ {
		workers[i] = packItems(done, items, i)
	}

	numPackages := 0
	for range merge(done, workers...) {
		numPackages++
	}

	fmt.Printf("Took %fs to ship %d packages\n", time.Since(start).Seconds(), numPackages)
}
