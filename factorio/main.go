package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Inserter1(belt chan int) {
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		belt <- i
	}

	fmt.Println("[Inserter1] My work is done. Closing the channel")
	close(belt)
}

func Inserter2(ctx context.Context, belt chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case i, ok := <-belt:
			if ok {
				fmt.Printf("[Inserter2] I got %d\n", i)
			} else {
				fmt.Println("[Inserter2] The belt is gone. I'm done")
				return
			}
		case <-ctx.Done():
			fmt.Printf("[Inserter2] Cancelling my work because : %q\n", ctx.Err())
			return
		}
	}
}

func main() {
	bgCtx := context.Background()
	ctx, _ := context.WithTimeout(bgCtx, 3*time.Second)

	belt := make(chan int)

	var i2Wg sync.WaitGroup

	go Inserter1(belt)

	for i := 0; i < 3; i++ {
		i2Wg.Add(1)
		go Inserter2(ctx, belt, &i2Wg)
	}

	i2Wg.Wait()

	fmt.Println("[Main Factory] All done! Closing up shop")
}

func KillSwitch(cancel context.CancelFunc) {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)
	<-sig
	fmt.Println("[Kill Switch] Ctrl-C pressed. Cancelling everything")
	cancel()
}
