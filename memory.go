package main

type Memory struct {
	value float64
}

func (m *Memory) Add(a float64) float64 {
	m.value += a
	return m.value
}

func (m *Memory) Sub(a float64) float64 {
	m.value -= a
	return m.value
}

func (m *Memory) Recall() float64 {
	return m.value
}
