package main

import (
	"errors"
	"fmt"
	"math"

	"fixed_point_iteration/pars"
)

func Fix(x_first, epsilon float64, max_interation int, list []pars.Querry) (float64, error) {
	var x_1, x_2 float64
	var i int

	x_1 = x_first
	for i = 0; i < max_interation; i++ {
		x_2 = pars.Function(x_1, list)
		var result = math.Abs(x_1 - x_2)
		if result < epsilon {
			return x_2, nil
		}
		x_1 = x_2
	}

	return 0, errors.New("Cannot find result")
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
	var x_1, epsilon float64
	var max_interation int

	x_1 = 10
	epsilon = 0.0000001
	max_interation = 100

	equation := pars.ReadFunction()
	list := pars.Make_list(equation)
	list = pars.Insert(list)

	var found_values = []float64{}

	for {
		var result, e = Fix(x_1, epsilon, max_interation, list)
		if e != nil {
			fmt.Printf("No root!\n")
			return
		}
		for _, s := range found_values {
			if math.Abs(s-result) < epsilon {
				fmt.Printf("No more root!\n")
				return
			}
		}
		found_values = append(found_values, result)
		fmt.Printf("%d%s root is %f\n", len(found_values), getOrdinalSuffix(len(found_values)), result)
		x_1 = result + 10
	}
}
