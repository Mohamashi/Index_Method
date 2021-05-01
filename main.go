package main

import (
	"fmt"
	"math"
)

func f1(x float64) float64 {
	return math.Exp(-0.5*x) * math.Sin(6*x-1.5)
}

func f2(x float64) float64 {
	return math.Abs(x) * math.Sin(2*math.Pi*x-0.5)
}

func fi(x float64) float64 {
	return math.Cos(18*x-3)*math.Sin(10*x-7) + 1.5
}

func main() {
	var a, b, eps, r float64

	fmt.Printf("a = ")
	fmt.Scanln(&a)
	fmt.Printf("b = ")
	fmt.Scanln(&b)
	fmt.Printf("eps = ")
	fmt.Scanln(&eps)
	fmt.Printf("r = ")
	fmt.Scanln(&r)

}
