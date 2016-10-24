package main

import (
	"fmt"
)

type test_struct struct {
	Test string
}

func main() {
	var k []interface{}
	fmt.Println(k) //result : null interface
	s := make([]interface{}, 3)
	s[0] = 5
	s[1] = false
	s[2] = "c"
	k = append(k, s...)
	fmt.Println(k)

}
