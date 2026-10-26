package main

import "fmt"

type StringIntMap struct {
	data map[string]int
}

func (m *StringIntMap) Add(key string, value int) {
	if m.data == nil {
		m.data = make(map[string]int)
	}
	m.data[key] = value
}

func (m *StringIntMap) Remove(key string) {
	delete(m.data, key)
}

func (m *StringIntMap) Copy() map[string]int {
	result := make(map[string]int, len(m.data))
	for key, value := range m.data {
		result[key] = value
	}
	return result
}

func (m *StringIntMap) Exists(key string) bool {
	_, ok := m.data[key]
	return ok
}

func (m *StringIntMap) Get(key string) (int, bool) {
	value, ok := m.data[key]
	return value, ok
}

func main() {
	var m StringIntMap
	m.Add("apple", 10)
	m.Add("banana", 20)
	m.Add("apple", 15)
	fmt.Println("Карта:", m.Copy())
	fmt.Println("Есть apple:", m.Exists("apple"))
	value, ok := m.Get("apple")
	fmt.Println("Значение apple:", value, "найдено:", ok)

	copied := m.Copy()
	m.Add("apple", 99)
	m.Remove("banana")
	fmt.Println("Оригинал после изменений:", m.Copy())
	fmt.Println("Независимая копия:", copied)
	fmt.Println("Есть banana:", m.Exists("banana"))
	value, ok = m.Get("missing")
	fmt.Println("Отсутствующий ключ:", value, "найдено:", ok)
}
