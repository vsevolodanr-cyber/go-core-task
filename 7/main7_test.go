package main

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"
)

func channel(values ...int) <-chan int {
	ch := make(chan int, len(values))
	for _, value := range values { ch <- value }
	close(ch)
	return ch
}

func collect(t *testing.T, ch <-chan int) []int {
	t.Helper()
	result := make([]int, 0)
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()
	for {
		select {
		case value, ok := <-ch:
			if !ok { return result }
			result = append(result, value)
		case <-timer.C:
			t.Fatal("output did not close")
		}
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name string
		inputs []<-chan int
		want []int
	}{
		{"none", nil, []int{}},
		{"nil", []<-chan int{nil, nil}, []int{}},
		{"closed", []<-chan int{channel(), channel()}, []int{}},
		{"one", []<-chan int{channel(1, 2, 3)}, []int{1, 2, 3}},
		{"many", []<-chan int{channel(1, 2), nil, channel(2, 3), channel(0)}, []int{0, 1, 2, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			got := collect(t, merge(ctx, tt.inputs...))
			sort.Ints(got)
			if !reflect.DeepEqual(got, tt.want) { t.Fatalf("got %v, want %v", got, tt.want) }
		})
	}
}

func TestMergePreservesEachInputOrder(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	got := collect(t, merge(ctx, channel(1, 2, 3), channel(10, 20, 30)))
	var low, high []int
	for _, value := range got {
		if value < 10 { low = append(low, value) } else { high = append(high, value) }
	}
	if !reflect.DeepEqual(low, []int{1, 2, 3}) || !reflect.DeepEqual(high, []int{10, 20, 30}) {
		t.Fatalf("order lost: %v", got)
	}
}

func TestMergeCancelWaitingInput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := merge(ctx, make(chan int))
	cancel()
	if got := collect(t, out); len(got) != 0 { t.Fatal(got) }
}

func TestMergeCancelBlockedOutput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	in := make(chan int)
	out := merge(ctx, in)
	select {
	case in <- 42: // Пересылатель принял число, но выход пока никто не читает.
	case <-time.After(2 * time.Second): t.Fatal("input not consumed")
	}
	cancel()
	collect(t, out)
}

func TestMergeManyInputs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inputs := make([]<-chan int, 100)
	for i := range inputs { inputs[i] = channel(i) }
	got := collect(t, merge(ctx, inputs...))
	sort.Ints(got)
	if len(got) != 100 { t.Fatalf("got %d items", len(got)) }
	for i, value := range got { if value != i { t.Fatalf("item %d = %d", i, value) } }
}
