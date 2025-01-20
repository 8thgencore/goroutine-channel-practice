package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	produceCount = 3
	produceStop  = 10
)

func produce(ctx context.Context, pipe chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		time.Sleep(time.Second * 3)
		fmt.Println("produce finished")
	}()

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		case pipe <- i:
		}
	}
}

func main() {
	pipe := make(chan int)
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < produceCount; i++ {
		wg.Add(1)
		go produce(ctx, pipe, &wg)
	}

	for i := range pipe {
		fmt.Println(i)

		if i == produceStop {
			cancel()
			break
		}
	}

	wg.Wait()
	close(pipe)

	fmt.Println("main finished")
}
