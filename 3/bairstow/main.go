package main

import (
	"bairstow_method/pars"
	"fmt"
	"math"
	"os"
)

type QuadRoot struct {
	real float64
	imag float64
}

func Make_Coeff(LL []pars.Querry) []float64 {

	var O, N string
	var V float64

	// Input as "1x^5-3.5x^4+2.75x^3+2.125x^2-3.875x+1.25"--
	F := make([]float64, 0)
	// get coefficient data from after inserting LIST -----------
	for i := 0; i < len(LL); i++ {
		N = pars.Get_name(LL, i)
		if N == "x" {
			if i == 0 {
				V = 1.0
			} else {
				V = pars.Get_value(LL, i-2)
				if i > 2 {
					O = pars.Get_operator(LL, i-3)
				}
				if O == "-" {
					V = -V
				}
			}
			F = append(F, V)
		}
	}
	// add last data as x^0 coefficient -------------------------
	V = pars.Get_value(LL, len(LL)-1)
	if V != 0.0 {
		O = pars.Get_operator(LL, len(LL)-2)
		if O == "-" {
			V = -V
		}
		F = append(F, V)
	}
	// Make reverse Coefficient array ---------------------------
	// Cieff array => { 1.25, -3.875, 2.125, 2.75, -3.5, 1 } ----
	Coeff := make([]float64, len(F))
	rank := len(F) - 1
	for i := rank; i >= 0; i-- {
		Coeff[rank-i] = F[i]
	}
	return Coeff
}

func Quad(rs []float64) (QuadRoot, QuadRoot) {

	var x1, x2 QuadRoot
	d := rs[0]*rs[0] + 4*rs[1]

	switch {
	case d > 0:
		x1.real = 0.5 * (rs[0] + math.Sqrt(d))
		x1.imag = 0.0
		x2.real = 0.5 * (rs[0] - math.Sqrt(d))
		x2.imag = 0.0
	case d == 0:
		x1.real = 0.5 * rs[0]
		x1.imag = 0.0
		x2.real = 0.5 * rs[0]
		x2.imag = 0.0
	default:
		x1.real = 0.5 * rs[0]
		x1.imag = 0.5 * math.Sqrt(-d)
		x2.real = 0.5 * rs[0]
		x2.imag = -0.5 * math.Sqrt(-d)
	}
	return x1, x2
}

func errFunc(rank int, Coef, b, rs []float64) []float64 {

	c := make([]float64, len(Coef))
	e := make([]float64, 2)

	// generate assumed coefficient after devide x^2-rx-s ------
	tem := rank - 1
	b[tem] = Coef[tem]
	b[tem-1] = Coef[tem-1] + rs[0]*b[tem]
	c[tem] = b[tem]
	c[tem-1] = b[tem-1] + rs[0]*c[tem]
	for k := tem - 2; k >= 0; k-- {
		b[k] = Coef[k] + rs[0]*b[k+1] + rs[1]*b[k+2]
		c[k] = b[k] + rs[0]*c[k+1] + rs[1]*c[k+2]
	}

	// Crammar rule -> x1= |D1|/|D|, x2=|D2|/|D| --------------
	det := c[2]*c[2] - c[3]*c[1]
	e[0] = (-c[2]*b[1] + c[3]*b[0]) / det
	e[1] = (-c[2]*b[0] + c[1]*b[1]) / det

	// arrange r, s with error for iteration ------------------
	rs[0] += e[0]
	rs[1] += e[1]
	e[0] /= rs[0]
	e[1] /= rs[1]

	// make error as absolute value ---------------------------
	if e[0] < 0.0 {
		e[0] = -e[0]
	}
	if e[1] < 0.0 {
		e[1] = -e[1]
	}
	return e
}

func Bairstow(epsilon float64, Coeff []float64, Max_i int) []QuadRoot {

	rank := len(Coeff)
	x := make([]QuadRoot, rank)
	b := make([]float64, rank)
	rs := []float64{-1.0, -1.0} // start values ----
	i := 0

	for {
		if rank <= 3 {
			break
		}
		i++
		switch {
		case i > Max_i:
			fmt.Printf(" It doesen't converge !!!\n")
			os.Exit(0)
		default:
			err := errFunc(rank, Coeff, b, rs)
			if err[0] <= epsilon && err[1] <= epsilon {
				// if b1, b0 is 0, then roots of x^2-rx-s are..
				x[rank-1], x[rank-2] = Quad(rs)
				rank -= 2 // for another iteration --------
				for j := 0; j < rank; j++ {
					Coeff[j] = b[j+2]
				}
			}
		}
	}
	if rank == 3 { // x^2 -rx -s ----------------------
		rs[0] = -Coeff[1] / Coeff[2]
		rs[1] = -Coeff[0] / Coeff[2]
		x[rank-1], x[rank-2] = Quad(rs)
	} else { // rank=2, ax +b =0 ----------------
		x[rank-1].real = -Coeff[0] / Coeff[1]
		x[rank-1].imag = 0.0
	}
	return x
}

func main() {

	var epsilon float64
	var Max_i int

	epsilon = 0.0001
	Max_i = 1000

	Eq := pars.ReadFunction()
	LIST := pars.Make_list(Eq)
	LIST = pars.Insert(LIST)

	Coef := Make_Coeff(LIST)
	result := Bairstow(epsilon, Coef, Max_i)

	fmt.Printf("\n ---------------- The Roots Are ----------------\n")
	for i := 1; i < len(result); i++ {
		fmt.Printf("   x[%1d].real = %7.3f,    x[%1d].imag = %7.3f\n",
			i, result[i].real, i, result[i].imag)
	}
	fmt.Printf(" ---------------- The Roots Are ----------------\n\n")
}
