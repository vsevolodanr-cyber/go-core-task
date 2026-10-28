package main

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func collect(t *testing.T, ch <-chan float64) []float64 {
	t.Helper()
	result := make([]float64, 0)
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()
	for {
		select {
		case value, ok := <-ch:
			if !ok { return result }
			result = append(result, value)
		case <-timer.C: t.Fatal("pipeline did not close")
		}
	}
}

func TestCubePipeline(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	input := make(chan uint8, 6)
	for _, value := range []uint8{0, 1, 2, 3, 10, 255} { input <- value }
	close(input)
	got := collect(t, cubePipeline(ctx, input))
	want := []float64{0, 1, 8, 27, 1000, 16581375}
	if !reflect.DeepEqual(got, want) { t.Fatalf("got %v, want %v", got, want) }
}

func TestEmptyInput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	input := make(chan uint8)
	close(input)
	if got := collect(t, cubePipeline(ctx, input)); len(got) != 0 { t.Fatal(got) }
}

func TestStreamingBeforeInputCloses(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	input := make(chan uint8)
	out := cubePipeline(ctx, input)
	for _, n := range []uint8{2, 4, 2} {
		select {
		case input <- n:
		case <-time.After(2 * time.Second): t.Fatal("input was not consumed")
		}
		select {
		case got, ok := <-out:
			v := float64(n)
			if !ok || got != v*v*v { t.Fatalf("unexpected result: %v, %v", got, ok) }
		case <-time.After(2 * time.Second): t.Fatal("pipeline waited for input to close")
		}
	}
	close(input)
	if tail := collect(t, out); len(tail) != 0 { t.Fatal(tail) }
}

func TestCancelWaitingInput(t *testing.T) {
	for _, input := range []<-chan uint8{nil, make(chan uint8)} {
		ctx, cancel := context.WithCancel(context.Background())
		out := cubePipeline(ctx, input)
		cancel()
		if got := collect(t, out); len(got) != 0 { t.Fatal(got) }
	}
}

func TestCancelBlockedOutput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	input := make(chan uint8)
	out := cubePipeline(ctx, input)
	select {
	case input <- 255:
	case <-time.After(2 * time.Second): t.Fatal("input not consumed")
	}
	cancel()
	collect(t, out)
}
