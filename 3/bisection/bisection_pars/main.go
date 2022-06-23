package main

import (
	"fmt"

	"bisec/pars"
)

func Bisec(LIST []pars.Querry, x_first, x_last, Delta, epsilon float64) (float64, float64) {
	var result, result1, result2 float64
	delta := Delta
	x_1 := x_first
	x_2 := x_1 + delta
	result1 = pars.Function(x_1, LIST)

	for {
		result2 = pars.Function(x_2, LIST)
		result = result1 * result2
		if result1 == 0.0 || result2 == 0.0 {
			return 0.0, x_1
		}

		if result < 0.0 {
			delta = 0.5 * delta
			if -result < epsilon {
				return result1, x_1
			}
			if delta <= epsilon {
				return result1, x_1
			}
			Bisec(LIST, x_1, x_2, delta, epsilon)
		} else {
			x_1 = x_2
			result1 = result2
			x_2 += delta
			if x_1 >= x_last {
				return 99.0, (x_last + 1.0)
			}
		}
	}
}

func main() {
	var start, finish, delta, epsilon, f_x float64 = -100.0, 100.0, 0.001, 0.0001, 0.0
	No_Root := 1
	x_1 := start

	Equation := pars.ReadFunction()
	LIST := pars.MakeList(Equation)
	LIST = pars.Insert(LIST)

	for {
		f_x, x_1 = Bisec(LIST, x_1, finish, delta, epsilon)
		if x_1 > finish {
			fmt.Printf("\n ----> No more Real ROOT FROM %v to %v!! \n\n", start, finish)
			return
		} else {
			fmt.Printf("The %3d th Real ROOT ( x=%6.3f, f(x)= %10.7f)\n", No_Root, x_1, f_x)
			No_Root++
			x_1 += delta
		}
	}
}
