package main

import "math"

func Floor(a float64) (float64, error) {
	return math.Floor(a), nil
}

func Ceil(a float64) (float64, error) {
	return math.Ceil(a), nil
}
