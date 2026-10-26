package main

import (
	"reflect"
	"testing"
)

func TestIntersection(t *testing.T) {
	tests := []struct {
		name          string
		first, second []int
		wantFound     bool
		want          []int
	}{
		{"example", []int{65, 3, 58, 678, 64}, []int{64, 2, 3, 43}, true, []int{64, 3}},
		{"no overlap", []int{1, 2}, []int{3, 4}, false, []int{}},
		{"one match", []int{1, 2}, []int{2}, true, []int{2}},
		{"duplicates", []int{2, 2, 1, 1}, []int{1, 1, 2, 2}, true, []int{1, 2}},
		{"zero and negative", []int{-3, 0, 2}, []int{0, -3, -1}, true, []int{0, -3}},
		{"same values", []int{1, 2, 3}, []int{3, 2, 1}, true, []int{3, 2, 1}},
		{"nil first", nil, []int{1}, false, []int{}},
		{"nil second", []int{1}, nil, false, []int{}},
		{"both nil", nil, nil, false, []int{}},
		{"both empty", []int{}, []int{}, false, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, got := intersection(tt.first, tt.second)
			if found != tt.wantFound || !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got (%v, %#v), want (%v, %#v)", found, got, tt.wantFound, tt.want)
			}
		})
	}
}

func TestIntersectionIndependent(t *testing.T) {
	first := []int{3, 2, 1, 99}
	second := []int{2, 3, 88}
	_, result := intersection(first[:3], second[:2])
	if !reflect.DeepEqual(first, []int{3, 2, 1, 99}) ||
		!reflect.DeepEqual(second, []int{2, 3, 88}) {
		t.Fatal("inputs were modified")
	}
	first[0], second[0] = -1, -2
	if !reflect.DeepEqual(result, []int{2, 3}) {
		t.Fatal("result shares input storage")
	}
	result[0] = -3
	if first[0] != -1 || first[1] != 2 || second[0] != -2 {
		t.Error("changing result affected input")
	}
}
