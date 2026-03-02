package main

import (
	"fmt"
	"runtime"
)

func getFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func zeroValues() {
	fmt.Printf("==========\n")
	fmt.Printf("%s\n", getFunctionName())
	var i int
	var s string
	var b bool
	var f float64
	var sl []int

	fmt.Printf("int: %d\n", i)
	fmt.Printf("string: %s\n", s)
	fmt.Printf("bool: %t\n", b)
	fmt.Printf("float: %f\n", f)
	fmt.Printf("slice: %v\n", sl)
	fmt.Printf("==========\n\n")
}

func typeInference() {
	fmt.Printf("==========\n")
	fmt.Printf("%s\n", getFunctionName())

	i := 4
	l := int64(1125899906842624)
	s := "something"
	b := true
	f := 2.23
	sl := []int{1, 2, 3}

	fmt.Printf("int32: %T\n", i)
	fmt.Printf("int64: %T\n", l)
	fmt.Printf("string: %T\n", s)
	fmt.Printf("bool: %T\n", b)
	fmt.Printf("float: %T\n", f)
	fmt.Printf("slice: %T\n", sl)
	fmt.Printf("==========\n\n")
}

func multiReturnFunc() (int, string, bool) {
	return 1, "2", true
}

func blankIdentifier() {
	fmt.Printf("==========\n")
	fmt.Printf("%s\n", getFunctionName())
	a, b, c := multiReturnFunc()
	fmt.Printf("%d %s %t\n", a, b, c)
	fmt.Printf("==========\n\n")
}

func main() {
	zeroValues()
	typeInference()
	blankIdentifier()
}
