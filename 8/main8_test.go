package main

import (
	"testing"
	"time"
)

func waitResult(t *testing.T, done <-chan struct{}) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return")
	}
}

func waitAsync(wg *WaitGroup) <-chan struct{} {
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	return done
}

func TestEmptyWait(t *testing.T) {
	wg := NewWaitGroup()
	wg.Add(0)
	waitResult(t, waitAsync(wg))
}

func TestWaitForAllTasks(t *testing.T) {
	wg := NewWaitGroup()
	const count = 128
	results := make([]int, count)
	wg.Add(count)
	for i := range results {
		go func(index int) {
			defer wg.Done()
			results[index] = index + 1
		}(i)
	}
	waitResult(t, waitAsync(wg))
	for i, value := range results {
		if value != i+1 {
			t.Fatalf("task %d unfinished: %d", i, value)
		}
	}
}

func TestWaitBlocksUntilZero(t *testing.T) {
	wg := NewWaitGroup()
	wg.Add(2)
	done := waitAsync(wg)
	wg.Done()
	select {
	case <-done:
		t.Fatal("Wait returned before all tasks finished")
	default:
	}
	wg.Done()
	waitResult(t, done)
}

func TestMultipleWaitersAndReuse(t *testing.T) {
	wg := NewWaitGroup()
	for round := 0; round < 10; round++ {
		wg.Add(1)
		waiters := make([]<-chan struct{}, 10)
		for i := range waiters {
			waiters[i] = waitAsync(wg)
		}
		wg.Done()
		for _, waiter := range waiters {
			waitResult(t, waiter)
		}
	}
}

func TestAddDuringActiveGroup(t *testing.T) {
	wg := NewWaitGroup()
	wg.Add(1)
	wg.Add(2)
	wg.Add(-2)
	wg.Done()
	waitResult(t, waitAsync(wg))
}

func TestInvalidCounter(t *testing.T) {
	wg := NewWaitGroup()
	func() {
		defer func() {
			if recover() == nil {
				t.Error("expected panic on negative count")
			}
		}()
		wg.Done()
	}()

	wg.Add(1)
	wg.Done()
	waitResult(t, waitAsync(wg))
}
