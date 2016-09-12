package main

import (
	"fmt"
	"strings"
)

func main() {
	someString := "one                           two   three four "

	words := strings.Fields(someString)

	fmt.Println(words, len(words)) // [one two three four] 4
}
