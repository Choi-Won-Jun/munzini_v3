// main.go
package main

import (
	"fmt"
)

func main() {
	a := []int{1, 2, 3}
	var b []int
	b = append(b, a...)
	fmt.Println(b)
}
