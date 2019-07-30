package test

import (
	"fmt"
	// "log"
)

func sayHello() int {
	fmt.Println("Hello")
	return 3
}

var x int = sayHello()

func GetX() int {
	return x
}
