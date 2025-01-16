package main

import (
	"errors"
	"sync"
)

type fn func() error

func main() {
	expErr := errors.New("error")
	funcs := []fn{
		func() error { return nil },
		func() error { return nil },
		func() error { return expErr },
		func() error { return nil },
	}
	if err := Run(funcs...); !errors.Is(err, expErr) {
		panic("wrong code")
	}
}

func Run(fs ...fn) error {
	errCh := make(chan error)
	var wg sync.WaitGroup

	for _, f := range fs {
		wg.Add(1)
		go func(f fn) {
			defer wg.Done()
			errCh <- f()
		}(f)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
