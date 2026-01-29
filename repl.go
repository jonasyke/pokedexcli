package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	words := strings.Fields(lowercase)
	return words
}

func main() {
	input := "  HeLLo   WoRLD  This    iS    a   TeSt  "
	result := cleanInput(input)
	fmt.Printf("%#v\n", result)
}
