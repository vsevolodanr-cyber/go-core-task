package main

import (
	"reflect"
	"testing"
)

func TestAddAndGet(t *testing.T) {
	var m StringIntMap
	for _, tt := range []struct {
		key   string
		value int
	}{
		{"apple", 10},
		{"apple", 20},
		{"zero", 0},
		{"", -5},
		{"ключ🙂", 42},
	} {
		m.Add(tt.key, tt.value)
		if got, ok := m.Get(tt.key); !ok || got != tt.value {
			t.Errorf("Get(%q) = (%d, %v), want (%d, true)", tt.key, got, ok, tt.value)
		}
	}
	if got, ok := m.Get("missing"); ok || got != 0 {
		t.Errorf("missing key: got (%d, %v), want (0, false)", got, ok)
	}
}

func TestExistsAndRemove(t *testing.T) {
	var m StringIntMap
	if m.Exists("zero") {
		t.Fatal("empty map must not contain key")
	}
	m.Remove("missing")
	m.Add("zero", 0)
	m.Add("keep", 5)
	if !m.Exists("zero") {
		t.Fatal("key with zero value must exist")
	}
	m.Remove("zero")
	m.Remove("zero")
	if m.Exists("zero") {
		t.Error("removed key still exists")
	}
	if value, ok := m.Get("zero"); ok || value != 0 {
		t.Error("Get must report removed key as absent")
	}
	if value, ok := m.Get("keep"); !ok || value != 5 {
		t.Error("Remove changed unrelated key")
	}
	m.Add("zero", 7)
	if value, ok := m.Get("zero"); !ok || value != 7 {
		t.Error("cannot re-add removed key")
	}
}

func TestZeroValue(t *testing.T) {
	var m StringIntMap
	if value, ok := m.Get("x"); value != 0 || ok {
		t.Error("Get on zero value must return (0, false)")
	}
	copy := m.Copy()
	if copy == nil || len(copy) != 0 {
		t.Fatal("Copy must return a writable empty map")
	}
	copy["x"] = 1
	if m.Exists("x") {
		t.Error("copy affected original")
	}
}

func TestCopyIndependent(t *testing.T) {
	var m StringIntMap
	m.Add("a", 1)
	m.Add("b", 2)
	copied := m.Copy()
	if !reflect.DeepEqual(copied, map[string]int{"a": 1, "b": 2}) {
		t.Fatalf("unexpected copy: %v", copied)
	}
	m.Add("a", 99)
	m.Remove("b")
	m.Add("c", 3)
	if !reflect.DeepEqual(copied, map[string]int{"a": 1, "b": 2}) {
		t.Error("original changes affected copy")
	}
	copied["a"] = -1
	copied["d"] = 4
	delete(copied, "b")
	if !reflect.DeepEqual(m.Copy(), map[string]int{"a": 99, "c": 3}) {
		t.Error("copy changes affected original")
	}
}
