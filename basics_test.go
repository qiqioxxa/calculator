package main

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 2, 3, 5},
		{"negative and positive", -2, 5, 3},
		{"zeros", 0, 0, 0},
		{"floats", 1.5, 2.25, 3.75},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Add(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Add(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive result", 10, 4, 6},
		{"negative result", 4, 10, -6},
		{"zeros", 0, 0, 0},
		{"floats", 5.5, 2.2, 3.3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Subtract(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Subtract(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 4, 5, 20},
		{"negative and positive", -3, 5, -15},
		{"multiply by zero", 123.45, 0, 0},
		{"floats", 2.5, 4, 10},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Multiply(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Multiply(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr bool
	}{
		{"exact division", 10, 2, 5, false},
		{"float result", 7, 2, 3.5, false},
		{"negative division", -9, 3, -3, false},
		{"division by zero", 10, 0, 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Divide(tc.a, tc.b)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Divide(%v, %v): expected error, got nil", tc.a, tc.b)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Divide(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestMod(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr bool
	}{
		{"remainder exists", 10, 3, 1, false},
		{"no remainder", 9, 3, 0, false},
		{"float remainder", 5.5, 2, 1.5, false},
		{"mod by zero", 10, 0, 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Mod(tc.a, tc.b)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Mod(%v, %v): expected error, got nil", tc.a, tc.b)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Mod(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}
