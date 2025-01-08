package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	pipeline := squarer(ctx, generator(ctx, 1, 2, 3))
	for x := range pipeline {
		fmt.Println(x)
	}
}

func generator(ctx context.Context, in ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, num := range in {
			select {
			case <-ctx.Done():
				return
			case out <- num:
			}
		}
	}()

	return out
}

func squarer(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for num := range in {
			select {
			case <-ctx.Done():
				return
			case out <- num * num:
			}
		}
	}()

	return out
}
