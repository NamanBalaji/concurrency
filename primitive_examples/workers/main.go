package main

import (
	"fmt"
	"math/rand"
)

type Work struct {
	value int
}

type Result struct {
	value int
}

func main() {
	workChan := make(chan Work)
	resultChan := make(chan Result)
	done := make(chan bool)

	workQueue := make([]Work, 100)
	for i := range workQueue {
		workQueue[i].value = rand.Int()
	}

	// create 10 worker go routines
	for i := 0; i < 10; i++ {
		go func() {
			for {
				// get work from worker chan
				work := <-workChan
				// compute result
				result := Result{
					value: work.value * 2,
				}
				// send teh result via result chan
				resultChan <- result
			}
		}()
	}

	results := make([]Result, 0)
	go func() {
		// colect all the results
		for i := 0; i < len(workQueue); i++ {
			results = append(results, <-resultChan)
		}
		// when all the results are collected, notify the done channel
		done <- true
	}()

	// send all the work to the workers
	for _, work := range workQueue {
		workChan <- work
	}
	// wait until everything is done
	<-done
	fmt.Println(results)
}
