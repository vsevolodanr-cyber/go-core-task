package main

import (
	"context"
	"fmt"
	"sync"
)

func merge(ctx context.Context, inputs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, input := range inputs {
		if input == nil {
			continue
		}
		wg.Add(1)
		go func(in <-chan int) {
			defer wg.Done()
			for {
				if ctx.Err() != nil {
					return
				}
				select {
				case <-ctx.Done():
					return
				case value, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case out <- value:
					}
				}
			}
		}(input)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	first := make(chan int, 3)
	second := make(chan int, 2)
	third := make(chan int, 1)
	first <- 1
	first <- 2
	first <- 3
	second <- 10
	second <- 20
	third <- 100
	close(first)
	close(second)
	close(third)
	fmt.Println("Объединённый поток:")
	for value := range merge(context.Background(), first, second, third) {
		fmt.Println(value)
	}
}
