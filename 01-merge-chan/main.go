package main

import (
	"fmt"
	"sync"
)

// merge — соединяет каналы в один
func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup

	for _, c := range cs {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(c)
	}

	// закрываем выходной канал, когда все горутины завершены
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// fillChan — заполняет канал числами от 0 до n-1
func fillChan(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // закрываем канал после записи всех значений
		for i := 0; i < n; i++ {
			out <- i
		}
	}()

	return out
}

func main() {
	a := fillChan(2)
	b := fillChan(3)
	c := fillChan(4)
	res := merge(a, b, c)
	for v := range res {
		fmt.Println(v)
	}
}
