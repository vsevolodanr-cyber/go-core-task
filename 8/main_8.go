package main

import "fmt"

type WaitGroup struct {
	semaphore chan struct{}
	count     int
	done      chan struct{}
}

func NewWaitGroup() *WaitGroup {
	done := make(chan struct{})
	close(done)
	return &WaitGroup{
		semaphore: make(chan struct{}, 1),
		done:      done,
	}
}

func (wg *WaitGroup) Add(delta int) {
	wg.semaphore <- struct{}{}
	defer func() { <-wg.semaphore }()

	next := wg.count + delta
	if next < 0 || (delta > 0 && next < wg.count) {
		panic("WaitGroup: invalid counter")
	}
	if wg.count == 0 && next > 0 {
		wg.done = make(chan struct{})
	}
	if wg.count > 0 && next == 0 {
		close(wg.done)
	}
	wg.count = next
}

func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

func (wg *WaitGroup) Wait() {
	wg.semaphore <- struct{}{}
	done := wg.done
	<-wg.semaphore
	<-done
}

func main() {
	wg := NewWaitGroup()
	results := make([]int, 5)
	wg.Add(len(results))
	for i := range results {
		go func(index int) {
			defer wg.Done()
			results[index] = (index + 1) * (index + 1)
		}(i)
	}
	wg.Wait()
	fmt.Println("Все горутины завершены:", results)
}
