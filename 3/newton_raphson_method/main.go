package main

import (
	"fmt"
	"math"
	"newton_raphson_method/pars"
)

func Newton(x_first, delta float64, list []pars.Querry) float64 {
	x_2 := x_first + delta
	f_2 := pars.Function(x_2, list)
	f_1 := (f_2 - pars.Function(x_first, list)) / delta

	x_2 = x_first - f_2/f_1
	return x_2
}

func getOrdinalSuffix(num int) string {
	var first = num % 10
	if first == 1 {
		return "st"
	} else if first == 2 {
		return "nd"
	} else if first == 3 {
		return "rd"
	} else {
		return "th"
	}
}

func main() {
	var x [10]float64

	var delta, epsilon, x_1, result float64
	var max_iteration, i int
	delta = 0.0000001
	epsilon = 0.0000001
	x_1 = 0.0
	result = 1.0
	max_iteration = 100
	i = 0

	equation := pars.ReadFunction()
	list := pars.Make_list(equation)
	list = pars.Insert(list)

	var found_values = []float64{}

	x[0] = x_1
	for i = 0; i < max_iteration; i++ {
		result = Newton(x_1, delta, list)
		if math.Abs(result-x_1) < epsilon {
			for _, v := range found_values {
				if math.Abs(v-result) < epsilon {
					fmt.Printf("No more Root!\n")
					return
				}
			}
			found_values = append(found_values, result)
			fmt.Printf("%d%s root is %f\n", len(found_values), getOrdinalSuffix(len(found_values)), result)
			x_1 = result + 10*math.Pow(-1, float64(len(found_values)))
			i = 0
		} else {
			x_1 = result
		}
	}

	fmt.Printf("No Root!\n")
}
