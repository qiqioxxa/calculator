package main

import (
	"fmt"
	"math"
)

func Add(a, b float64) (float64, error) {
	return a + b, nil
}

func Subtract(a, b float64) (float64, error) {
	return a - b, nil
}

func Multiply(a, b float64) (float64, error) {
	return a * b, nil
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func Mod(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("mod by zero")
	}
	return float64(math.Mod(a, b)), nil
}
