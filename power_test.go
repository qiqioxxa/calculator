package main

import (
	"math"
	"testing"
)

func TestPower(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive exponent", 2, 3, 8},
		{"zero exponent", 5, 0, 1},
		{"negative exponent", 2, -1, 0.5},
		{"fractional exponent", 9, 0.5, 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Power(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Power(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		want    float64
		wantErr bool
	}{
		{"perfect square", 16, 4, false},
		{"zero", 0, 0, false},
		{"float", 2.25, 1.5, false},
		{"negative number", -9, 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Sqrt(tc.a)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Sqrt(%v): expected error, got nil", tc.a)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Sqrt(%v) = %v; want %v", tc.a, got, tc.want)
			}
		})
	}
}
