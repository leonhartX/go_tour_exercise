package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z, prev := 1.0, 0.0
	diff := 0.01
	for math.Abs(z-prev) > diff {
		prev = z
		z = z - (z*z-x)/2*z
		fmt.Println(z, prev)
	}
	return z
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(math.Sqrt(2))
}
