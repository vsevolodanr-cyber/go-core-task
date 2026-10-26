package main

import (
	"fmt"
	"math/rand"
	"time"
)

func sliceExample(input []int) []int {
	result := make([]int, 0)
	for _, value := range input {
		if value%2 == 0 {
			result = append(result, value)
		}
	}
	return result
}

func addElements(input []int, value int) []int {
	result := make([]int, len(input)+1)
	copy(result, input)
	result[len(input)] = value
	return result
}

func copySlice(input []int) []int {
	if input == nil {
		return nil
	}
	result := make([]int, len(input))
	copy(result, input)
	return result
}

func removeElement(input []int, index int) ([]int, error) {
	if index < 0 || index >= len(input) {
		return nil, fmt.Errorf("индекс %d вне границ слайса длиной %d", index, len(input))
	}

	result := make([]int, len(input)-1)
	copy(result, input[:index])
	copy(result[index:], input[index+1:])
	return result, nil
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	originalSlice := make([]int, 10)
	for i := range originalSlice {
		originalSlice[i] = rng.Intn(100)
	}
	fmt.Println("Исходный слайс:", originalSlice)
	fmt.Println("Чётные числа:", sliceExample(originalSlice))
	fmt.Println("Добавлено число:", addElements(originalSlice, 42))

	copied := copySlice(originalSlice)
	originalSlice[3] = 111
	fmt.Println("Копия:", copied)
	fmt.Println("Оригинал:", originalSlice)

	removed, err := removeElement(originalSlice, 3)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Удалён элемент:", removed)
	fmt.Println("Исходный слайс после операций:", originalSlice)
}
