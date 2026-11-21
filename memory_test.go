package main

import "testing"

func TestMemoryOperations(t *testing.T) {
	var m Memory

	if got := m.Recall(); got != 0 {
		t.Fatalf("initial Recall() = %v; want 0", got)
	}

	if got := m.Add(10); got != 10 {
		t.Errorf("Add(10) = %v; want 10", got)
	}

	if got := m.Add(5.5); got != 15.5 {
		t.Errorf("Add(5.5) = %v; want 15.5", got)
	}

	if got := m.Sub(3.5); got != 12 {
		t.Errorf("Sub(3.5) = %v; want 12", got)
	}

	if got := m.Recall(); got != 12 {
		t.Errorf("Recall() = %v; want 12", got)
	}

	m.Clear()
	if got := m.Recall(); got != 0 {
		t.Errorf("after Clear(), Recall() = %v; want 0", got)
	}
}
