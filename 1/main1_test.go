package main

import (
	"reflect"
	"testing"
)

func TestTypeOfAndToString(t *testing.T) {
	tests := []struct {
		value    any
		typeName string
		text     string
	}{
		{42, "int", "42"},
		{052, "int", "42"},
		{0x2A, "int", "42"},
		{3.14, "float64", "3.14"},
		{"Golang", "string", "Golang"},
		{true, "bool", "true"},
		{complex64(1 + 2i), "complex64", "(1+2i)"},
	}
	for _, tt := range tests {
		if got := TypeOf(tt.value); got != tt.typeName {
			t.Errorf("TypeOf(%v) = %q, want %q", tt.value, got, tt.typeName)
		}
		if got := ToString(tt.value); got != tt.text {
			t.Errorf("ToString(%v) = %q, want %q", tt.value, got, tt.text)
		}
	}
}

func TestJoinValues(t *testing.T) {
	got := JoinValues(42, 052, 0x2A, 3.14, "Golang", true, complex64(1+2i))
	if want := "4242423.14Golangtrue(1+2i)"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got := JoinValues(); got != "" {
		t.Errorf("empty input: got %q", got)
	}
}

func TestToRunes(t *testing.T) {
	got := ToRunes("Я🙂a")
	want := []rune{'Я', '🙂', 'a'}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if len(ToRunes("")) != 0 {
		t.Error("empty string must produce zero runes")
	}
}

func TestInsertSalt(t *testing.T) {
	tests := []struct{ input, salt, want string }{
		{"abcd", "go-2024", "abgo-2024cd"},
		{"abcde", "go-2024", "abgo-2024cde"},
		{"", "go-2024", "go-2024"},
		{"a", "go-2024", "go-2024a"},
		{"Я🙂界", "go-2024", "Яgo-2024🙂界"},
		{"abcd", "", "abcd"},
		{"ab", "Ж🙂", "aЖ🙂b"},
	}
	for _, tt := range tests {
		t.Run(tt.input+"/"+tt.salt, func(t *testing.T) {
			input := make([]rune, len([]rune(tt.input)), 100)
			copy(input, []rune(tt.input))
			got := InsertSalt(input, tt.salt)
			if string(got) != tt.want {
				t.Errorf("got %q, want %q", string(got), tt.want)
			}
			if string(input) != tt.input {
				t.Error("input was modified")
			}
			if len(got) > 0 && len(input) > 0 {
				got[0] = '!'
				if string(input) != tt.input {
					t.Error("result shares input storage")
				}
			}
		})
	}
}

func TestHashRunes(t *testing.T) {
	tests := []struct{ input, salt, want string }{
		{"", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"abc", "", "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
		{"", "go-2024", "66802df107aace17871a5b610ff9eb11706e13477bb24e93966ca80671c0fac6"},
		{"Я🙂界", "go-2024", "ec1447767b2f0c35e9a70e8f11d60a867694a8cdf01310b1973716c60b4a97f5"},
		{"4242423.14Golangtrue(1+2i)", "go-2024", "53f2f60ac6c41389d3ed3d84d88d8c2860bf8981c677be18243a6f35a6b6a1b3"},
	}
	for _, tt := range tests {
		if got := HashRunes([]rune(tt.input), tt.salt); got != tt.want {
			t.Errorf("HashRunes(%q, %q) = %q, want %q", tt.input, tt.salt, got, tt.want)
		}
	}
}
