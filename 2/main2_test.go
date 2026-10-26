package main

import (
	"reflect"
	"testing"
)

func TestSliceExample(t *testing.T) {
	tests := []struct {
		name        string
		input, want []int
	}{
		{"mixed", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []int{2, 4, 6, 8, 10}},
		{"negative and zero", []int{-3, -2, 0, 1, 2}, []int{-2, 0, 2}},
		{"duplicates", []int{4, 2, 4}, []int{4, 2, 4}},
		{"odd only", []int{1, 3, 5}, []int{}},
		{"empty", []int{}, []int{}},
		{"nil", nil, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceExample(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddElements(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		value int
		want  []int
	}{
		{"normal", []int{1, 2}, 3, []int{1, 2, 3}},
		{"negative", []int{1}, -5, []int{1, -5}},
		{"empty", []int{}, 0, []int{0}},
		{"nil", nil, 7, []int{7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addElements(tt.input, tt.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopySlice(t *testing.T) {
	for _, input := range [][]int{nil, {}, {1, 2, 3}} {
		if got := copySlice(input); !reflect.DeepEqual(got, input) {
			t.Errorf("copySlice(%#v) = %#v", input, got)
		}
	}
}

func TestRemoveElement(t *testing.T) {
	tests := []struct {
		name    string
		input   []int
		index   int
		want    []int
		wantErr bool
	}{
		{"first", []int{1, 2, 3}, 0, []int{2, 3}, false},
		{"middle", []int{1, 2, 3}, 1, []int{1, 3}, false},
		{"last", []int{1, 2, 3}, 2, []int{1, 2}, false},
		{"only", []int{5}, 0, []int{}, false},
		{"negative index", []int{1, 2}, -1, nil, true},
		{"index equals length", []int{1, 2}, 2, nil, true},
		{"index beyond length", []int{1, 2}, 100, nil, true},
		{"empty", []int{}, 0, nil, true},
		{"nil", nil, 0, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := removeElement(tt.input, tt.index)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

// Проверяем независимость памяти, даже когда append мог бы использовать
// свободную ёмкость исходного слайса.
func TestIndependentResults(t *testing.T) {
	tests := []struct {
		name  string
		apply func([]int) []int
	}{
		{"filter", sliceExample},
		{"add", func(s []int) []int { return addElements(s, 8) }},
		{"copy", copySlice},
		{"remove", func(s []int) []int {
			result, err := removeElement(s, 1)
			if err != nil {
				t.Fatal(err)
			}
			return result
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backing := []int{2, 4, 6, 99, 100}
			input := backing[:3]
			got := tt.apply(input)
			if !reflect.DeepEqual(backing, []int{2, 4, 6, 99, 100}) {
				t.Fatalf("function modified original storage: %v", backing)
			}
			if len(got) == 0 {
				t.Fatal("unexpected empty result")
			}
			first := got[0]
			input[0] = -10
			if got[0] != first {
				t.Error("changing original changed result")
			}
			got[0] = -20
			if input[0] != -10 {
				t.Error("changing result changed original")
			}
		})
	}
}
