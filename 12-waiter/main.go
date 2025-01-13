package main

import (
	"context"
	"errors"
	"sync"
)

type waiter interface {
	wait() error
	run(ctx context.Context, f func(ctx context.Context) error)
}

type waitGroup struct {
	mu     sync.Mutex
	wg     sync.WaitGroup
	sem    chan struct{}
	errors []error
}

func (g *waitGroup) wait() error {
	g.wg.Wait()
	return errors.Join(g.errors...)
}

func (g *waitGroup) run(ctx context.Context, fn func(ctx context.Context) error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		g.sem <- struct{}{}
		defer func() { <-g.sem }()

		if err := fn(ctx); err != nil {
			g.mu.Lock()
			g.errors = append(g.errors, err)
			g.mu.Unlock()
		}
	}()
}

func newGroupWait(maxParallel int) waiter {
	g := &waitGroup{
		sem: make(chan struct{}, maxParallel),
	}

	return g
}

func main() {
	g := newGroupWait(2)
	ctx := context.Background()
	expErr1 := errors.New("got error 1")
	expErr2 := errors.New("got error 2")

	g.run(ctx, func(ctx context.Context) error {
		return nil
	})

	g.run(ctx, func(ctx context.Context) error {
		return expErr2
	})

	g.run(ctx, func(ctx context.Context) error {
		return expErr1
	})

	err := g.wait()
	if !errors.Is(err, expErr1) || !errors.Is(err, expErr2) {
		panic("wrong code")
	}
}
