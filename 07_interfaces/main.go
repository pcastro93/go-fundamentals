package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

type Rectangle struct {
	width  float64
	height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

func main() {
	var shape Shape
	var c = Circle{radius: 10}
	var r = Rectangle{2, 3}
	shape = c
	fmt.Printf("Area: %f\n", shape.Area())
	fmt.Printf("Perimeter: %f\n", shape.Perimeter())
	shape = r
	fmt.Printf("Area: %f\n", shape.Area())
	fmt.Printf("Perimeter: %f\n", shape.Perimeter())
}
