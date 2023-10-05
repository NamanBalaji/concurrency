package main

import (
	"fmt"
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

func packItems(done <-chan bool, items <-chan item) <-chan int {
	packages := make(chan int)
	go func() {
		for item := range items {
			select {
			case <-done:
				return
			case packages <- item.id:
				time.Sleep(item.effort)
				fmt.Printf("Shipping package no. %d\n", item.id)
			}
		}
		close(packages)
	}()

	return packages
}

func main() {
	done := make(chan bool)
	defer close(done)

	start := time.Now()

	packages := packItems(done, prepareItems(done))
	numPackages := 0
	for range packages {
		numPackages++
	}

	fmt.Printf("Took %fs to ship %d packages\n", time.Since(start).Seconds(), numPackages)
}
