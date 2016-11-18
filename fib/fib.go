package main

import "fmt"

func fib(x int) int {
	switch x {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return fib(x-1) + fib(x-2)
	}
}

func main() {
	x := 10
	fmt.Println(fib(10))
	fmt.Println(&x)
}
