package pars

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Querry struct {
	identity string  // Type of value as N,P,Q,V,O,F
	name     string  // Function and variable name
	value    float64 // evaluated value
	operator string  // (, ), +,-,*,/,^
}

func ReadFunction() string {
	fmt.Printf("\n\nInput Function : ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	fmt.Printf("\n\n     Your Input: %v", string([]byte(input)))
	input = strings.Replace(input, " ", "", -1)

	fmt.Printf("Condensed Func.: %v", string([]byte(input)))
	fmt.Printf("\n")
	return string([]byte(input))
}

func MakeList(fomular string) []Querry {
	var Q Querry
	var i int

	LIST := make([]Querry, 0)

	for i = 0; i < len(fomular); i++ {
		c := string([]byte(fomular)[i])
		switch {
		case c >= "0" && c <= "9":
			fallthrough
		case c == ".":
			{
				number := c
				for {
					i++
					c := string([]byte(fomular)[i])
					if c >= "0" && c <= "9" || c == "." {
						number = number + c // concaternate string
					} else {
						Q.value, _ = strconv.ParseFloat(number, 64)
						Q.identity, Q.name, Q.operator = "V", "", ""
						LIST = append(LIST, Q)
						i--
						break
					}
				}
			}
		case c >= "a" && c <= "z":
			fallthrough
		case c >= "A" && c <= "Z":
			{
				number := c
				for {
					i++
					c := string([]byte(fomular)[i])
					if c >= "a" && c <= "z" || c >= "A" && c <= "Z" {
						number = number + c // concaternate string
					} else {
						Q.identity, Q.value, Q.name, Q.operator = "F", 0.0, number, ""
						LIST = append(LIST, Q)
						i--
						break
					}
				}
			}
		case c == "+":
			fallthrough
		case c == "-":
			fallthrough
		case c == "*":
			fallthrough
		case c == "/":
			fallthrough
		case c == "^":
			Q.identity, Q.value, Q.name, Q.operator = "O", 0.0, "", c
			if i == 0 {
				Q.identity = "N"
			}
			LIST = append(LIST, Q)
		case c == "(":
			Q.identity, Q.value, Q.name, Q.operator = "P", 0.0, "", c
			LIST = append(LIST, Q)
		case c == ")":
			Q.identity, Q.value, Q.name, Q.operator = "Q", 0.0, "", c
			LIST = append(LIST, Q)
		}
	}
	return LIST
}

func deleteN(L []Querry, ID int, DN int) []Querry {
	copy(L[ID:], L[ID+DN:]) // copy i+1, i+2, .. --> i, i+1, ..
	L = L[:len(L)-DN]       // truncate last DN slot
	return L
}

func Insert(L []Querry) []Querry {
	var P Querry

	Q := make([]Querry, 0) // assign pointer Querry structure
	l := len(L)            // length of L structure array

	for i := 0; i < l-1; i++ { // unnary operator
		if (L[i].identity == "N") ||
			(L[i].operator == "-" && L[i+1].identity == "O") {
			if L[i+1].operator == "-" {
				L[i+1].identity, L[i+1].value = "V", -1.0
			} else {
				L[i].identity, L[i].value = "V", -1.0
			}
		} else if L[i].identity == "P" && L[i+1].operator == "-" {
			L[i+1].identity, L[i+1].value = "V", -1.0
		} else if L[i].operator == "^" && L[i+1].operator == "-" {
			L[i+2].value = -1.0 * L[i+2].value
			L = deleteN(L, i+1, 1)
		} else if L[i].identity == "O" && L[i+1].operator == "-" {
			L[i+1].identity, L[i+1].value = "V", -1.0
		}
	}

	l = len(L) // length of L structure array
	for i := 0; i < l-1; i++ {
		if (L[i].identity == "V" && L[i+1].identity == "F") ||
			(L[i].identity == "V" && L[i+1].identity == "V") ||
			(L[i].identity == "V" && L[i+1].identity == "P") ||
			(L[i].identity == "F" && L[i+1].identity == "V") {
			P.identity, P.operator = "O", "*"
			Q = append(Q, L[i])
			Q = append(Q, P) // Insertion of * operator
		} else {
			Q = append(Q, L[i])
		}
	}
	Q = append(Q, L[l-1]) // The last character
	return Q
}

func Print_List(Q []Querry) {

	for i := 0; i < len(Q); i++ {
		fmt.Printf("%v : ", Q[i].identity)
		switch Q[i].identity {
		case "V":
			{
				fmt.Printf("%v\n", Q[i].value)
			}
		case "N":
			fallthrough
		case "O":
			fallthrough
		case "P":
			fallthrough
		case "Q":
			{
				fmt.Printf("%v\n", Q[i].operator)
			}
		case "F":
			{
				fmt.Printf("%v\n", Q[i].name)
			}
		}
	}
}

func reduce(L []Querry, NO int) []Querry {
	P := Querry{identity: "V"}
	j := L[NO].operator
	switch j {
	case "+":
		P.value = L[NO-1].value + L[NO+1].value // plus
	case "-":
		P.value = L[NO-1].value - L[NO+1].value // minus
	case "*":
		P.value = L[NO-1].value * L[NO+1].value // times
	case "/":
		P.value = L[NO-1].value / L[NO+1].value // devide
	case "^":
		P.value = math.Pow(L[NO-1].value, L[NO+1].value) // exponential
	}
	// delete NO-1,NO,NO+1 and assign computed value into NO-1 array
	L = deleteN(L, NO, 2)
	L[NO-1] = P // copy P into L
	return L
}

func reduce_parenthesis(L []Querry, start, finish int) []Querry {
	for i := start; i <= finish; i++ {
		if L[i].identity == "F" { // reduce 1 slot in reduc_function
			L = reduce_function(L, i)
			i, finish = start, finish-1
		}
	}
	for i := start; i <= finish; i++ {
		if L[i].operator == "^" { // reduce 2 slot in reducing
			L = reduce(L, i)
			i, finish = start, finish-2
		}
	}
	for i := start; i <= finish; i++ {
		if L[i].operator == "*" { // reduce 2 slot in reducing
			L = reduce(L, i)
			i, finish = start, finish-2
		}
		if L[i].operator == "/" { // reduce 2 slot in reducing
			L = reduce(L, i)
			i, finish = start, finish-2
		}
	}
	for i := start; i <= finish; i++ {
		if L[i].operator == "+" { // reduce 2 slot in reducing
			L = reduce(L, i)
			i, finish = start, finish-2
		}
		if L[i].operator == "-" { // reduce   slot in reducing
			L = reduce(L, i)
			i, finish = start, finish-2
		}
	}
	L = deleteN(L, finish, 1) // you should delete finish first
	L = deleteN(L, start, 1)
	return L
}

func reduce_function(L []Querry, NO int) []Querry {
	switch L[NO].name {
	case "sin":
		L[NO].value, L[NO].identity = math.Sin(L[NO+1].value), "V"
	case "cos":
		L[NO].value, L[NO].identity = math.Cos(L[NO+1].value), "V"
	case "tan":
		L[NO].value, L[NO].identity = math.Tan(L[NO+1].value), "V"
	case "abs":
		L[NO].value, L[NO].identity = math.Abs(L[NO+1].value), "V"
	}

	// ------------- delete f->next ----------------- */
	L = deleteN(L, NO+1, 1)
	return L
}

func evaluate(v float64, L []Querry) float64 {
	// substitution x value -----------------------------------------------
	for i := 0; i < len(L); i++ {
		if L[i].name == "x" {
			L[i].value, L[i].identity = v, "V"
		}
	}

	// Reduce for parenthesis ----------------------------------------------
	var start, finish int
	var stack [10]int
	stack_counter := 0

	l := len(L)
	for i := 0; i < l; i++ {
		switch L[i].operator {
		case "(":
			{
				stack[stack_counter] = i
				stack_counter++
			}
		case ")":
			{
				stack_counter--
				start = stack[stack_counter]
				finish = i
				L = reduce_parenthesis(L, start, finish)
				l = len(L)
				i = start // adjust loop length
			}
		}
	}

	// compute functions such as sin(), cos().------------------------------
	for i := 0; i < l-1; i++ {
		if L[i].identity == "F" {
			L = reduce_function(L, i)
			l = len(L)
		}
	}

	// compute and reducing list by priority ===============================
	// firstly exponential -----------------------------------------------
	for i := 0; i < l-1; i++ {
		if L[i].operator == "^" {
			L = reduce(L, i)
			i, l = 0, len(L)
		}
	}

	// secondly times & devide -------------------------------------------
	for i := 0; i < l; i++ {
		if L[i].operator == "*" {
			L = reduce(L, i)
			i, l = 0, len(L)
		}
		if L[i].operator == "/" {
			L = reduce(L, i)
			i, l = 0, len(L)
		}
	}

	// lastly plus & minus -----------------------------------------------
	for i := 0; i < l; i++ {
		if L[i].operator == "+" {
			L = reduce(L, i)
			i, l = 0, len(L)
		}
		if L[i].operator == "-" {
			L = reduce(L, i)
			i, l = 0, len(L)
		}
	}

	// return function value ---------------------------------------------
	return L[0].value
}

func Function(x float64, LIST []Querry) float64 {
	T_LIST := make([]Querry, len(LIST)) // Prepare to Copy LIST
	copy(T_LIST, LIST)                  // copy LIST into T_LIST
	return (evaluate(x, T_LIST))
}
