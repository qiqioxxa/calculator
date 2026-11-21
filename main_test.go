package main

import "testing"

func TestCalculatorBasicFlow(t *testing.T) {
	c := NewCalculator()

	c.InputDigit("1")
	c.InputDigit("2")
	c.ChooseOp("add", "+")
	c.InputDigit("8")
	c.Calculate()

	if c.Display != "20" {
		t.Fatalf("12 + 8: got %s; want 20", c.Display)
	}

	c.ChooseOp("div", "÷")
	c.InputDigit("0")
	c.Calculate()
	if !c.hasError {
		t.Errorf("expected error on division by zero")
	}

	c.Clear()
	if c.Display != "0" || c.hasError {
		t.Errorf("Clear() failed: Display=%s, hasError=%v", c.Display, c.hasError)
	}
}

func TestCalculatorUnaryAndMemory(t *testing.T) {
	c := NewCalculator()

	c.InputDigit("9")
	c.UnaryOp("sqrt", "√")
	if c.Display != "3" {
		t.Fatalf("sqrt(9): got %s; want 3", c.Display)
	}

	c.MemOp("m+")
	c.Clear()
	c.MemOp("mr")
	if c.Display != "3" {
		t.Fatalf("MR: got %s; want 3", c.Display)
	}

	c.MemOp("mc")
	if c.Mem.Recall() != 0 {
		t.Errorf("MC: memory = %v; want 0", c.Mem.Recall())
	}
}
