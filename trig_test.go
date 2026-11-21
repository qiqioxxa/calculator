package main

import (
	"math"
	"testing"
)

func TestSin(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		want float64
	}{
		{"zero", 0, 0},
		{"pi/2", math.Pi / 2, 1},
		{"pi", math.Pi, 0},
		{"-pi/2", -math.Pi / 2, -1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Sin(tc.a)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Sin(%v) = %v; want %v", tc.a, got, tc.want)
			}
		})
	}
}

func TestCos(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		want float64
	}{
		{"zero", 0, 1},
		{"pi/2", math.Pi / 2, 0},
		{"pi", math.Pi, -1},
		{"2pi", 2 * math.Pi, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Cos(tc.a)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Cos(%v) = %v; want %v", tc.a, got, tc.want)
			}
		})
	}
}
