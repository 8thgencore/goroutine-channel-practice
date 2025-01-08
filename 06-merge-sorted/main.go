package main

import "fmt"

func mergeSorted(a, b <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		var valA, okA = <-a
		var valB, okB = <-b

		for okA || okB {
			// Если один из каналов закрыт, берем значения только из другого
			if !okA {
				out <- valB
				valB, okB = <-b
				continue
			}
			if !okB {
				out <- valA
				valA, okA = <-a
				continue
			}

			// Если оба канала открыты, берем меньший элемент
			if valA < valB {
				out <- valA
				valA, okA = <-a
			} else {
				out <- valB
				valB, okB = <-b
			}
		}
	}()

	return out
}

func fillChanA(c chan int) {
	c <- 1
	c <- 2
	c <- 4
	close(c)
}

func fillChanB(c chan int) {
	c <- -1
	c <- 4
	c <- 5
	close(c)
}

func main() {
	a, b := make(chan int), make(chan int)
	go fillChanA(a)
	go fillChanB(b)
	c := mergeSorted(a, b)

	for val := range c {
		fmt.Printf("%d ", val)
	}
}
