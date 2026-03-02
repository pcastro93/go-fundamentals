package main

import "fmt"

func main() {
	l := NewLinkedList[int]()
	for i := range 5 {
		l.Append(i * 2)
	}
	data, ok := l.Get(3)
	if ok {
		fmt.Println("Element at idx 3:", data)
	}
	fmt.Println("Size:", l.Size())
}
