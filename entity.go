package main

import (
	"math"
)

type Geometry interface {
	area() float64
}

type GeometryObject struct {
	Id int
}

type Rect struct {
	GeometryObject
	Width  int
	Length int
}

func (rect Rect) area() float64 {
	return float64(rect.Width * rect.Length)
}

type Circle struct {
	GeometryObject
	radius float64
}

func (c Circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
