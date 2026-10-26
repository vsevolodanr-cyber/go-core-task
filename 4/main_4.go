package main

import "fmt"

func difference(first, second []string) []string {
	excluded := make(map[string]struct{}, len(second))
	for _, value := range second {
		excluded[value] = struct{}{}
	}

	result := make([]string, 0)
	for _, value := range first {
		if _, found := excluded[value]; !found {
			result = append(result, value)
		}
	}
	return result
}

func main() {
	slice1 := []string{"apple", "banana", "cherry", "date", "43", "lead", "gno1"}
	slice2 := []string{"banana", "date", "fig"}
	fmt.Println("Первый слайс:", slice1)
	fmt.Println("Второй слайс:", slice2)
	fmt.Println("Только в первом:", difference(slice1, slice2))
}
