package main

import (
	"errors"
	"fmt"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func multiple_returns() {
	fmt.Printf("==========\n")
	r, e := divide(1, 2)
	fmt.Printf("%d %v\n", r, e)
	r, e = divide(2, 1)
	fmt.Printf("%d %v\n", r, e)
	r, e = divide(2, 0)
	fmt.Printf("%d %v\n", r, e)
}

func named_returns() {
	fmt.Printf("==========\n")
	swap := func(a, b string) (first, second string) {
		first, second = b, a
		return
	}
	f, s := swap("first", "2")
	fmt.Printf("%s %s\n", f, s)
}

func variadic_functions() {
	fmt.Printf("==========\n")
	var_func := func(to_print ...int) {
		for i, v := range to_print {
			fmt.Printf("pos: %d - elment: %d\n", i, v)
		}
	}
	var_func(1)
	var_func(1, 2)
	sl := []int{1, 2, 3}
	var_func(sl...)
}

func main() {
	multiple_returns()
	named_returns()
	variadic_functions()
}
