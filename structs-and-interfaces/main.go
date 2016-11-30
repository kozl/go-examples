package main

import (
	"fmt"
	"math"
)

type Shape interface {
	area() float64
}

type MultiShape struct {
	shapes []Shape
}

type Circle struct {
	x, y, r float64
}

type Rectangle struct {
	x1, x2, x3, x4, y1, y2, y3, y4 float64
}

func distance(x1, y1, x2, y2 float64) float64 {
	a := x2 - x1
	b := y2 - y1
	return math.Sqrt(a*a + b*b)
}

func (m *MultiShape) area() float64 {
	var total float64
	for _, v := range m.shapes {
		total += v.area()
	}
	return total
}

func (c *Circle) area() float64 {
	return math.Pi * c.r * c.r
}

func (r *Rectangle) area() float64 {
	l := distance(r.x1, r.y1, r.x2, r.y2)
	w := distance(r.x1, r.y1, r.x3, r.y3)
	return w * l
}

func main() {
	mShape := MultiShape{
		shapes: []Shape{
			&Rectangle{x1: 0, y1: 0, x2: 0, y2: 2, x3: 5, y3: 0, x4: 5, y4: 2},
			&Circle{x: 0, y: 0, r: 1},
		},
	}
	fmt.Println(mShape.area())
}
