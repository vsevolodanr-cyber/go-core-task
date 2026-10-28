package main

import (
	"context"
	"fmt"
)

func cubePipeline(ctx context.Context, input <-chan uint8) <-chan float64 {
	out := make(chan float64)
	go func() {
		defer close(out)
		for {
			if ctx.Err() != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case number, ok := <-input:
				if !ok {
					return
				}
				value := float64(number)
				select {
				case <-ctx.Done():
					return
				case out <- value * value * value:
				}
			}
		}
	}()
	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	input := make(chan uint8)
	output := cubePipeline(ctx, input)
	go func() {
		defer close(input)
		for _, value := range []uint8{0, 1, 2, 3, 10, 255} {
			select {
			case <-ctx.Done():
				return
			case input <- value:
			}
		}
	}()
	fmt.Println("Кубы чисел:")
	for value := range output {
		fmt.Printf("%.0f\n", value)
	}
}
