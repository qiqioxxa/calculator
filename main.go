package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Calculator struct {
	Display          string
	Expr             string
	Mem              Memory
	firstOperand     float64
	hasFirstOperand  bool
	operator         string
	waitingForSecond bool
	hasError         bool
}

func NewCalculator() *Calculator {
	return &Calculator{Display: "0"}
}

func formatFloat(v float64) string {
	if v == 0 {
		v = 0
	}
	return strconv.FormatFloat(v, 'g', 12, 64)
}

func (c *Calculator) Clear() {
	c.Display = "0"
	c.Expr = ""
	c.firstOperand = 0
	c.hasFirstOperand = false
	c.operator = ""
	c.waitingForSecond = false
	c.hasError = false
}

func (c *Calculator) InputDigit(d string) {
	if c.hasError {
		c.Clear()
	}
	if c.waitingForSecond {
		c.Display = d
		c.waitingForSecond = false
		return
	}
	if c.Display == "0" {
		c.Display = d
	} else {
		c.Display += d
	}
}

func (c *Calculator) InputDot() {
	if c.hasError {
		c.Clear()
	}
	if c.waitingForSecond {
		c.Display = "0."
		c.waitingForSecond = false
		return
	}
	if !strings.Contains(c.Display, ".") {
		c.Display += "."
	}
}

func (c *Calculator) ToggleSign() {
	if c.hasError || c.Display == "0" {
		return
	}
	if strings.HasPrefix(c.Display, "-") {
		c.Display = strings.TrimPrefix(c.Display, "-")
	} else {
		c.Display = "-" + c.Display
	}
}

func (c *Calculator) Backspace() {
	if c.hasError {
		c.Clear()
		return
	}
	if c.waitingForSecond {
		return
	}
	if len(c.Display) > 1 {
		c.Display = c.Display[:len(c.Display)-1]
	} else {
		c.Display = "0"
	}
	if c.Display == "-" {
		c.Display = "0"
	}
}

func (c *Calculator) ChooseOp(op, symbol string) {
	if c.hasError {
		return
	}
	if c.operator != "" && !c.waitingForSecond {
		c.Calculate()
		if c.hasError {
			return
		}
	}
	val, err := strconv.ParseFloat(c.Display, 64)
	if err != nil {
		return
	}
	c.firstOperand = val
	c.hasFirstOperand = true
	c.operator = op
	c.waitingForSecond = true
	c.Expr = fmt.Sprintf("%s %s", formatFloat(val), symbol)
}

func (c *Calculator) Calculate() {
	if c.hasError || !c.hasFirstOperand || c.operator == "" {
		return
	}
	second, err := strconv.ParseFloat(c.Display, 64)
	if err != nil {
		return
	}

	var res float64
	switch c.operator {
	case "add":
		res, err = Add(c.firstOperand, second)
	case "sub":
		res, err = Subtract(c.firstOperand, second)
	case "mul":
		res, err = Multiply(c.firstOperand, second)
	case "div":
		res, err = Divide(c.firstOperand, second)
	case "mod":
		res, err = Mod(c.firstOperand, second)
	case "pow":
		res, err = Power(c.firstOperand, second)
	default:
		err = fmt.Errorf("unknown op")
	}

	if err != nil {
		c.Display = "Err: " + err.Error()
		c.hasError = true
	} else if math.IsNaN(res) || math.IsInf(res, 0) {
		c.Display = "Err: math error"
		c.hasError = true
	} else {
		c.Display = formatFloat(res)
	}

	c.Expr = ""
	c.hasFirstOperand = false
	c.operator = ""
	c.waitingForSecond = true
}

func (c *Calculator) UnaryOp(op, label string) {
	if c.hasError {
		return
	}
	val, err := strconv.ParseFloat(c.Display, 64)
	if err != nil {
		return
	}

	var res float64
	switch op {
	case "sin":
		res, err = Sin(val)
	case "cos":
		res, err = Cos(val)
	case "sqrt":
		res, err = Sqrt(val)
	case "floor":
		res, err = Floor(val)
	case "ceil":
		res, err = Ceil(val)
	default:
		err = fmt.Errorf("unknown op")
	}

	if err != nil {
		c.Display = "Err: " + err.Error()
		c.hasError = true
	} else if math.IsNaN(res) || math.IsInf(res, 0) {
		c.Display = "Err: math error"
		c.hasError = true
	} else {
		c.Expr = fmt.Sprintf("%s(%s)", label, formatFloat(val))
		c.Display = formatFloat(res)
		c.waitingForSecond = true
	}
}

func (c *Calculator) MemOp(op string) {
	if c.hasError && op != "mc" && op != "mr" {
		return
	}
	val, _ := strconv.ParseFloat(c.Display, 64)
	switch op {
	case "m+":
		c.Mem.Add(val)
	case "m-":
		c.Mem.Sub(val)
	case "mr":
		c.Display = formatFloat(c.Mem.Recall())
		c.hasError = false
		c.waitingForSecond = false
	case "mc":
		c.Mem.Clear()
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Калькулятор")

	calc := NewCalculator()

	memLabel := widget.NewLabel("M: 0")
	exprLabel := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})
	displayLabel := widget.NewLabelWithStyle("0", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	refresh := func() {
		displayLabel.SetText(calc.Display)
		exprLabel.SetText(calc.Expr)
		memLabel.SetText("M: " + formatFloat(calc.Mem.Recall()))
	}

	btn := func(label string, fn func()) *widget.Button {
		return widget.NewButton(label, func() {
			fn()
			refresh()
		})
	}

	clearBtn := btn("C", func() { calc.Clear() })
	clearBtn.Importance = widget.DangerImportance

	eqBtn := btn("=", func() { calc.Calculate() })
	eqBtn.Importance = widget.HighImportance

	grid := container.NewGridWithColumns(4,
		btn("MC", func() { calc.MemOp("mc") }),
		btn("MR", func() { calc.MemOp("mr") }),
		btn("M+", func() { calc.MemOp("m+") }),
		btn("M-", func() { calc.MemOp("m-") }),

		btn("sin", func() { calc.UnaryOp("sin", "sin") }),
		btn("cos", func() { calc.UnaryOp("cos", "cos") }),
		btn("√", func() { calc.UnaryOp("sqrt", "√") }),
		btn("xʸ", func() { calc.ChooseOp("pow", "^") }),

		btn("floor", func() { calc.UnaryOp("floor", "floor") }),
		btn("ceil", func() { calc.UnaryOp("ceil", "ceil") }),
		btn("%", func() { calc.ChooseOp("mod", "%") }),
		clearBtn,

		btn("7", func() { calc.InputDigit("7") }),
		btn("8", func() { calc.InputDigit("8") }),
		btn("9", func() { calc.InputDigit("9") }),
		btn("÷", func() { calc.ChooseOp("div", "÷") }),

		btn("4", func() { calc.InputDigit("4") }),
		btn("5", func() { calc.InputDigit("5") }),
		btn("6", func() { calc.InputDigit("6") }),
		btn("×", func() { calc.ChooseOp("mul", "×") }),

		btn("1", func() { calc.InputDigit("1") }),
		btn("2", func() { calc.InputDigit("2") }),
		btn("3", func() { calc.InputDigit("3") }),
		btn("−", func() { calc.ChooseOp("sub", "-") }),

		btn("±", func() { calc.ToggleSign() }),
		btn("0", func() { calc.InputDigit("0") }),
		btn(".", func() { calc.InputDot() }),
		btn("+", func() { calc.ChooseOp("add", "+") }),
	)

	bottomRow := container.NewGridWithColumns(2,
		btn("⌫", func() { calc.Backspace() }),
		eqBtn,
	)

	statusRow := container.NewBorder(nil, nil, memLabel, exprLabel)
	content := container.NewVBox(
		statusRow,
		displayLabel,
		widget.NewSeparator(),
		grid,
		bottomRow,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(340, 440))
	w.ShowAndRun()
}
