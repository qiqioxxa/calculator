package main

import (
	"fmt"
	"math"
)

func Power(a, b float64) (float64, error) {
	return math.Pow(a, b), nil
}

func Sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, fmt.Errorf("square root of negative number")
	}
	return math.Sqrt(a), nil
}
