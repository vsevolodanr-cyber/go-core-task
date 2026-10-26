package main

import "fmt"

func intersection(first, second []int) (bool, []int) {
	available := make(map[int]struct{}, len(first))
	for _, value := range first {
		available[value] = struct{}{}
	}

	result := make([]int, 0)
	for _, value := range second {
		if _, found := available[value]; found {
			result = append(result, value)
			delete(available, value) // Исключаем повторную выдачу этого значения.
		}
	}
	return len(result) > 0, result
}

func main() {
	a := []int{65, 3, 58, 678, 64}
	b := []int{64, 2, 3, 43}
	found, common := intersection(a, b)
	fmt.Println("Первый слайс:", a)
	fmt.Println("Второй слайс:", b)
	fmt.Println("Есть пересечения:", found)
	fmt.Println("Общие значения:", common)
}
