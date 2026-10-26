package main

import (
	"reflect"
	"testing"
)

func TestDifference(t *testing.T) {
	tests := []struct {
		name                string
		first, second, want []string
	}{
		{"example", []string{"apple", "banana", "cherry", "date", "43", "lead", "gno1"}, []string{"banana", "date", "fig"}, []string{"apple", "cherry", "43", "lead", "gno1"}},
		{"no overlap", []string{"b", "a"}, []string{"c"}, []string{"b", "a"}},
		{"all excluded", []string{"a", "b"}, []string{"b", "a"}, []string{}},
		{"duplicates", []string{"a", "b", "a", "b"}, []string{"b", "b"}, []string{"a", "a"}},
		{"unicode and empty string", []string{"", "ёж", "🙂", "Apple", "apple"}, []string{"", "ёж", "apple"}, []string{"🙂", "Apple"}},
		{"nil first", nil, []string{"a"}, []string{}},
		{"nil second", []string{"b", "a", "b"}, nil, []string{"b", "a", "b"}},
		{"both nil", nil, nil, []string{}},
		{"both empty", []string{}, []string{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := difference(tt.first, tt.second); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestDifferenceIndependent(t *testing.T) {
	first := []string{"a", "b", "c", "sentinel"}
	second := []string{"b", "other"}
	result := difference(first[:3], second[:1])
	if !reflect.DeepEqual(first, []string{"a", "b", "c", "sentinel"}) ||
		!reflect.DeepEqual(second, []string{"b", "other"}) {
		t.Fatal("inputs were modified")
	}
	first[0] = "changed"
	if !reflect.DeepEqual(result, []string{"a", "c"}) {
		t.Fatal("result shares input storage")
	}
	result[0] = "result changed"
	if first[0] != "changed" || second[0] != "b" {
		t.Error("changing result affected input")
	}
}
