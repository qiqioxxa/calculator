package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
)

type calcRequest struct {
	Op string  `json:"op"`
	A  float64 `json:"a"`
	B  float64 `json:"b"`
}

type calcResponse struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

type memoryRequest struct {
	Op    string  `json:"op"`
	Value float64 `json:"value"`
}

type memoryResponse struct {
	Memory float64 `json:"memory"`
	Error  string  `json:"error,omitempty"`
}

var (
	calcMem Memory
	memMu   sync.Mutex
)

func handleCalc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req calcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var (
		res float64
		err error
	)

	switch req.Op {
	case "add":
		res, err = Add(req.A, req.B)
	case "sub":
		res, err = Subtract(req.A, req.B)
	case "mul":
		res, err = Multiply(req.A, req.B)
	case "div":
		res, err = Divide(req.A, req.B)
	case "mod":
		res, err = Mod(req.A, req.B)
	case "pow":
		res, err = Power(req.A, req.B)
	case "sin":
		res, err = Sin(req.A)
	case "cos":
		res, err = Cos(req.A)
	case "sqrt":
		res, err = Sqrt(req.A)
	case "floor":
		res, err = Floor(req.A)
	case "ceil":
		res, err = Ceil(req.A)
	default:
		err = fmt.Errorf("unknown operation: %s", req.Op)
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(calcResponse{Error: err.Error()})
		return
	}
	if math.IsNaN(res) || math.IsInf(res, 0) {
		json.NewEncoder(w).Encode(calcResponse{Error: "math error"})
		return
	}

	json.NewEncoder(w).Encode(calcResponse{Result: res})
}

func handleMemory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req memoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	memMu.Lock()
	defer memMu.Unlock()

	var val float64
	var errMsg string

	switch req.Op {
	case "m+":
		val = calcMem.Add(req.Value)
	case "m-":
		val = calcMem.Sub(req.Value)
	case "mr":
		val = calcMem.Recall()
	case "mc":
		calcMem.Clear()
		val = calcMem.Recall()
	default:
		errMsg = "unknown memory operation"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(memoryResponse{Memory: val, Error: errMsg})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHTML)
	})
	http.HandleFunc("/api/calc", handleCalc)
	http.HandleFunc("/api/memory", handleMemory)

	fmt.Println("Calculator running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

const indexHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Калькулятор</title>
<style>
  * { box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
  body {
    margin: 0; min-height: 100vh; display: flex; align-items: center; justify-content: center;
    background: #121218; color: #fff;
  }
  .calc {
    width: 340px; background: #1e1e2a; border-radius: 18px; padding: 20px;
    box-shadow: 0 20px 50px rgba(0,0,0,0.5);
  }
  .status {
    display: flex; justify-content: space-between; font-size: 13px; color: #8c8ca1; min-height: 20px; margin-bottom: 4px;
  }
  .display {
    width: 100%; background: #14141d; border-radius: 12px; padding: 14px 16px;
    text-align: right; font-size: 32px; font-weight: 600; overflow-x: auto; white-space: nowrap;
    margin-bottom: 16px;
  }
  .grid {
    display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px;
  }
  button {
    border: none; border-radius: 10px; padding: 14px 6px; font-size: 16px; font-weight: 500;
    cursor: pointer; background: #2c2c3d; color: #fff; transition: background 0.15s, transform 0.05s;
  }
  button:hover { background: #3a3a50; }
  button:active { transform: scale(0.96); }
  .mem { background: #263248; color: #8ab4f8; font-size: 14px; }
  .sci { background: #2b2740; color: #c5b4ff; font-size: 14px; }
  .op  { background: #3f365c; color: #e2d9ff; font-size: 18px; }
  .eq  { background: #4f7cff; color: #fff; grid-column: span 2; font-size: 20px; }
  .eq:hover { background: #648cff; }
  .clr { background: #5a2a33; color: #ffb3be; }
</style>
</head>
<body>
<div class="calc">
  <div class="status">
    <span id="mem">M: 0</span>
    <span id="expr"></span>
  </div>
  <div class="display" id="display">0</div>
  <div class="grid">
    <button class="mem" onclick="memOp('mc')">MC</button>
    <button class="mem" onclick="memOp('mr')">MR</button>
    <button class="mem" onclick="memOp('m+')">M+</button>
    <button class="mem" onclick="memOp('m-')">M-</button>

    <button class="sci" onclick="unaryOp('sin', 'sin')">sin</button>
    <button class="sci" onclick="unaryOp('cos', 'cos')">cos</button>
    <button class="sci" onclick="unaryOp('sqrt', '√')">√</button>
    <button class="sci" onclick="chooseOp('pow', '^')">xʸ</button>

    <button class="sci" onclick="unaryOp('floor', 'floor')">floor</button>
    <button class="sci" onclick="unaryOp('ceil', 'ceil')">ceil</button>
    <button class="op"  onclick="chooseOp('mod', '%')">%</button>
    <button class="clr" onclick="clearAll()">C</button>

    <button onclick="inputDigit('7')">7</button>
    <button onclick="inputDigit('8')">8</button>
    <button onclick="inputDigit('9')">9</button>
    <button class="op" onclick="chooseOp('div', '÷')">÷</button>

    <button onclick="inputDigit('4')">4</button>
    <button onclick="inputDigit('5')">5</button>
    <button onclick="inputDigit('6')">6</button>
    <button class="op" onclick="chooseOp('mul', '×')">×</button>

    <button onclick="inputDigit('1')">1</button>
    <button onclick="inputDigit('2')">2</button>
    <button onclick="inputDigit('3')">3</button>
    <button class="op" onclick="chooseOp('sub', '-')">−</button>

    <button onclick="toggleSign()">±</button>
    <button onclick="inputDigit('0')">0</button>
    <button onclick="inputDot()">.</button>
    <button class="op" onclick="chooseOp('add', '+')">+</button>

    <button class="clr" onclick="backspace()" style="grid-column: span 2;">⌫</button>
    <button class="eq" onclick="calculate()">=</button>
  </div>
</div>

<script>
  let current = '0';
  let firstOperand = null;
  let operator = null;
  let waitingForSecond = false;
  let hasError = false;

  const displayEl = document.getElementById('display');
  const exprEl = document.getElementById('expr');
  const memEl = document.getElementById('mem');

  function formatNum(n) {
    return Number(n.toPrecision(12)).toString();
  }

  function updateDisplay() {
    displayEl.textContent = current;
  }

  function clearAll() {
    current = '0';
    firstOperand = null;
    operator = null;
    waitingForSecond = false;
    hasError = false;
    exprEl.textContent = '';
    updateDisplay();
  }

  function inputDigit(d) {
    if (hasError) clearAll();
    if (waitingForSecond) {
      current = d;
      waitingForSecond = false;
    } else {
      current = current === '0' ? d : current + d;
    }
    updateDisplay();
  }

  function inputDot() {
    if (hasError) clearAll();
    if (waitingForSecond) {
      current = '0.';
      waitingForSecond = false;
    } else if (!current.includes('.')) {
      current += '.';
    }
    updateDisplay();
  }

  function toggleSign() {
    if (hasError || current === '0') return;
    current = current.startsWith('-') ? current.slice(1) : '-' + current;
    updateDisplay();
  }

  function backspace() {
    if (hasError) { clearAll(); return; }
    if (waitingForSecond) return;
    current = current.slice(0, -1);
    if (current === '' || current === '-') current = '0';
    updateDisplay();
  }

  async function chooseOp(op, symbol) {
    if (hasError) return;
    if (operator && !waitingForSecond) {
      await calculate();
      if (hasError) return;
    }
    firstOperand = parseFloat(current);
    operator = op;
    waitingForSecond = true;
    exprEl.textContent = firstOperand + ' ' + symbol;
  }

  async function calculate() {
    if (hasError || operator === null || firstOperand === null) return;
    const second = parseFloat(current);
    const resp = await fetch('/api/calc', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ op: operator, a: firstOperand, b: second })
    });
    const data = await resp.json();
    if (data.error) {
      current = 'Err: ' + data.error;
      hasError = true;
    } else {
      current = formatNum(data.result);
    }
    exprEl.textContent = '';
    firstOperand = null;
    operator = null;
    waitingForSecond = true;
    updateDisplay();
  }

  async function unaryOp(op, label) {
    if (hasError) return;
    const val = parseFloat(current);
    const resp = await fetch('/api/calc', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ op: op, a: val, b: 0 })
    });
    const data = await resp.json();
    if (data.error) {
      current = 'Err: ' + data.error;
      hasError = true;
    } else {
      exprEl.textContent = label + '(' + val + ')';
      current = formatNum(data.result);
      waitingForSecond = true;
    }
    updateDisplay();
  }

  async function memOp(op) {
    if (hasError && op !== 'mc' && op !== 'mr') return;
    const val = hasError ? 0 : parseFloat(current);
    const resp = await fetch('/api/memory', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ op: op, value: val })
    });
    const data = await resp.json();
    memEl.textContent = 'M: ' + formatNum(data.memory);
    if (op === 'mr') {
      current = formatNum(data.memory);
      hasError = false;
      waitingForSecond = false;
      updateDisplay();
    }
  }
</script>
</body>
</html>
`
