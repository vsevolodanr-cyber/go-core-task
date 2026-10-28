package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func randomGenerator(ctx context.Context, seed int64) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		rng := rand.New(rand.NewSource(seed))
		for {
			if ctx.Err() != nil {
				return
			}
			value := rng.Intn(100)
			select {
			case <-ctx.Done():
				return
			case out <- value:
			}
		}
	}()
	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	numbers := randomGenerator(ctx, time.Now().UnixNano())
	fmt.Println("10 случайных чисел:")
	for i := 0; i < 10; i++ {
		fmt.Println(<-numbers)
	}
	cancel()

	for range numbers {
	}
	fmt.Println("Генератор остановлен")
}
