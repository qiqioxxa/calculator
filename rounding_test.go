package main

import "testing"

func TestFloor(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		want float64
	}{
		{"positive float", 3.7, 3},
		{"integer", 3.0, 3},
		{"negative float", -2.3, -3},
		{"zero", 0, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Floor(tc.a)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Floor(%v) = %v; want %v", tc.a, got, tc.want)
			}
		})
	}
}

func TestCeil(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		want float64
	}{
		{"positive float", 3.2, 4},
		{"integer", 3.0, 3},
		{"negative float", -2.7, -2},
		{"zero", 0, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Ceil(tc.a)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Ceil(%v) = %v; want %v", tc.a, got, tc.want)
			}
		})
	}
}
