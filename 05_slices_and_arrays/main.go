package main

import "fmt"

func creation() {
	arr := [5]int{}
	sl := make([]int, 5)
	fmt.Printf("%#v\n", arr)
	fmt.Printf("%#v\n", sl)
}

func appending() {
	arr := [5]int{1, 2, 3, 4, 5}
	sl := make([]int, 5)
	// arr = append(arr, 6) // does not compile
	sl = append(sl, 6)
	fmt.Printf("%#v\n", arr)
	fmt.Printf("%#v\n", sl)
}

func copying() {
	sl := make([]int, 5)
	sl = append(sl, 1, 2, 3)
	sl2 := make([]int, len(sl))
	copy(sl2, sl)
	sl2 = append(sl2, 7)
	fmt.Printf("sl: %#v\n", sl)
	fmt.Printf("sl2: %#v\n", sl2)
}

func main() {
	creation()
	appending()
	copying()
}
