package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func TypeOf(value any) string {
	return fmt.Sprintf("%T", value)
}

func ToString(value any) string {
	return fmt.Sprint(value)
}

func JoinValues(values ...any) string {
	var result strings.Builder
	for _, value := range values {
		result.WriteString(ToString(value))
	}
	return result.String()
}

func ToRunes(value string) []rune {
	return []rune(value)
}

func InsertSalt(runes []rune, salt string) []rune {
	middle := len(runes) / 2
	saltRunes := []rune(salt)
	result := make([]rune, 0, len(runes)+len(saltRunes))
	result = append(result, runes[:middle]...)
	result = append(result, saltRunes...)
	result = append(result, runes[middle:]...)
	return result
}

func HashRunes(runes []rune, salt string) string {
	salted := InsertSalt(runes, salt)
	hash := sha256.Sum256([]byte(string(salted)))
	return fmt.Sprintf("%x", hash)
}

func main() {
	var numDecimal int = 42
	var numOctal int = 052
	var numHexadecimal int = 0x2A
	var pi float64 = 3.14
	var name string = "Golang"
	var isActive bool = true
	var complexNum complex64 = 1 + 2i

	values := []any{numDecimal, numOctal, numHexadecimal, pi, name, isActive, complexNum}
	for _, value := range values {
		fmt.Printf("Значение: %v, тип: %s\n", value, TypeOf(value))
	}

	joined := JoinValues(values...)
	runes := ToRunes(joined)
	const salt = "go-2024"
	fmt.Println("Объединённая строка:", joined)
	fmt.Println("Срез рун:", runes)
	fmt.Println("SHA-256:", HashRunes(runes, salt))
}
