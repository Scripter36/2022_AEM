package main

import (
	"fmt"
	"math/big"
	pars_v2 "runge_kutta/pars"
)

func factorial(n int) int {
	fac := new(big.Int)
	fac.MulRange(1, int64(n))
	return int(fac.Int64())
}

func RungeKutta(equation pars_v2.Equation, n int, start, end, start_value, step float64, a []float64) float64 {
	current_value := start_value
	for x := start; x < end; x += step {
		kValues := make([]float64, 1)
		m := make(map[string]float64)
		m["x"] = x
		m["y"] = current_value
		kValues[0] = equation.Value.Reduce(m).Num
		for j := 1; j < n; j++ { // n차에 대해...
			kx := x + 1.0/float64(factorial(j+1))/a[j]*step
			ky := current_value
			for l := 1; l < j; l++ {
				ky += 1.0 / float64(factorial(j+1)) / a[j] * step * float64(kValues[l-1])
			}
			m := make(map[string]float64)
			m["x"] = kx
			m["y"] = ky
			kValues = append(kValues, equation.Value.Reduce(m).Num)
		}
		for j := 0; j < n; j++ {
			current_value += a[j] * step * kValues[j]
		}
	}

	return current_value
}

func main() {
	// y'
	equation, err := pars_v2.Equation{RawValue: "0-2x^3+12x^2-20x+8.5"}.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 초기값
	start := 0.0
	start_value := 1.0
	// 목표
	end := 4.0
	// step
	step := 0.5
	// 계수들
	a := []float64{1}
	fmt.Println("Euler:", RungeKutta(equation, 1, start, end, start_value, step, a))
	a = []float64{0.5, 0.5}
	fmt.Println("HEUN:", RungeKutta(equation, 2, start, end, start_value, step, a))
	a = []float64{0, 1}
	fmt.Println("Midpoint:", RungeKutta(equation, 2, start, end, start_value, step, a))
	a = []float64{1.0 / 3.0, 2.0 / 3.0}
	fmt.Println("2nd-order:", RungeKutta(equation, 2, start, end, start_value, step, a))
}
