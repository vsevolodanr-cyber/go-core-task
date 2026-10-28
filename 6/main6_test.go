package main

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func nextNumber(t *testing.T, ch <-chan int) int {
	t.Helper()
	select {
	case value, ok := <-ch:
		if !ok {
			t.Fatal("channel closed too early")
		}
		return value
	case <-time.After(2 * time.Second):
		t.Fatal("generator did not send")
		return 0
	}
}

func awaitClosed(t *testing.T, ch <-chan int) {
	t.Helper()
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-timer.C:
			t.Fatal("generator did not close after cancellation")
		}
	}
}

func TestRandomGenerator(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := randomGenerator(ctx, 42)
	if cap(ch) != 0 {
		t.Fatal("channel must be unbuffered")
	}
	expected := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		got := nextNumber(t, ch)
		if got < 0 || got >= 100 {
			t.Fatalf("out of range: %d", got)
		}
		if want := expected.Intn(100); got != want {
			t.Fatalf("item %d: got %d, want %d", i, got, want)
		}
	}
	cancel()
	awaitClosed(t, ch)
}

func TestCancelWithoutConsumer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := randomGenerator(ctx, 1)
	cancel()
	awaitClosed(t, ch)
}

func TestAlreadyCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ch := randomGenerator(ctx, 1)
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("already canceled generator produced a value")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("generator did not close")
	}
}
