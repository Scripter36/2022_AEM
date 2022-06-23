package main // 실제 모듈로 사용하려면 변경이 필요합니다.

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type EquationElementType int

const (
	VAR EquationElementType = iota
	NUM
	ADD
	MULT
	DIV
	POW
	UNKNOWN
	ANY
)

type EquationElement struct {
	Type  EquationElementType
	Left  *EquationElement
	Right *EquationElement
	Value string
	Num   float64
}

type Equation struct {
	RawValue string
	Value    EquationElement
	Parsed   bool
}

func strToOperator(operator byte) EquationElementType {
	switch operator {
	case '+':
		return ADD
	case '*':
		return MULT
	case '/':
		return DIV
	case '^':
		return POW
	default:
		return UNKNOWN
	}
}

func findMatchingParenthesisEnd(str string, start int) int {
	count := 1
	length := len(str)
	for i := start; i < length; i++ {
		if str[i] == '(' {
			count++
		} else if str[i] == ')' {
			count--
			if count == 0 {
				return i
			}
		}
	}
	return -1
}

func appendMinus(element EquationElement) (EquationElement, bool) {
	switch element.Type {
	case NUM:
		return EquationElement{Type: NUM, Num: element.Num * -1, Value: fmt.Sprintf("%f", element.Num*-1)}, true
	case MULT:
		left_result, success := appendMinus(*element.Left)
		if success {
			return EquationElement{Type: MULT, Left: &left_result, Right: element.Right}, true
		} else {
			right_result, success := appendMinus(*element.Right)
			if success {
				return EquationElement{Type: MULT, Left: element.Left, Right: &right_result}, true
			} else {
				return EquationElement{Type: MULT, Left: &EquationElement{Type: NUM, Value: "-1", Num: -1}, Right: &element}, false
			}
		}
	default:
		return EquationElement{Type: MULT, Left: &EquationElement{Type: NUM, Value: "-1", Num: -1}, Right: &element}, false
	}
}

func parseEquationElement(element EquationElement) (EquationElement, error) {
	fmt.Println(element)
	switch element.Type {
	case UNKNOWN:
		parenthesis_start := strings.Index(element.Value, "(")
		if parenthesis_start != -1 {
			parenthesis_end := findMatchingParenthesisEnd(element.Value, parenthesis_start+1)

			left := element.Value[:parenthesis_start]
			parenthesis := element.Value[parenthesis_start+1 : parenthesis_end]
			right := element.Value[parenthesis_end+1:]

			has_left := len(left) > 0
			has_right := len(right) > 0

			var left_operator, right_operator EquationElementType
			var left_element, right_element EquationElement
			if has_left {
				left_operator = strToOperator(left[len(left)-1])

				new_value := left
				if left_operator == UNKNOWN {
					new_value = new_value + "*"
				}
				new_value = new_value + "t"

				var err error
				left_element, err = parseEquationElement(EquationElement{
					Type:  UNKNOWN,
					Value: new_value,
				})
				if err != nil {
					return EquationElement{}, err
				}
			}
			if has_right {
				right_operator = strToOperator(right[0])

				new_value := right
				if right_operator == UNKNOWN {
					new_value = "*" + new_value
				}
				new_value = "t" + new_value

				var err error
				right_element, err = parseEquationElement(EquationElement{
					Type:  UNKNOWN,
					Value: new_value,
				})
				if err != nil {
					return EquationElement{}, err
				}
			}

			parenthesisElement, err := parseEquationElement(EquationElement{Type: UNKNOWN, Value: parenthesis})
			if err != nil {
				return EquationElement{}, err
			}

			resultElement := &parenthesisElement

			if has_left {
				prev := &left_element
				goal := &left_element
				for goal.Right != nil {
					prev = goal
					goal = goal.Right
				}
				prev.Right = resultElement

				resultElement = &left_element
			}

			if has_right {
				prev := &right_element
				goal := &right_element
				for goal.Left != nil {
					prev = goal
					goal = goal.Left
				}
				prev.Left = resultElement
				resultElement = &right_element
			}
			return *resultElement, nil
		} else {
			operator_strings := []string{"+", "-", "*", "/", "^"}
			operator_types := []EquationElementType{ADD, ADD, MULT, DIV, POW}
			operator_names := []string{"ADD", "SUB", "MULT", "DIV", "POW"}
			for index := range operator_strings {
				operator_index := strings.Index(element.Value, operator_strings[index])
				if operator_index != -1 {
					if operator_strings[index] != "-" && operator_index == 0 {
						return EquationElement{}, errors.New(fmt.Sprintf("SyntaxError: %s operator should have left operand.", operator_names[index]))
					}
					if operator_index == len(element.Value)-1 {
						return EquationElement{}, errors.New(fmt.Sprintf("SyntaxError: %s operator should have right operand.", operator_names[index]))
					}
					left := element.Value[:operator_index]
					if operator_types[index] == POW {
						onevar_reg := regexp.MustCompile("(\\d*(\\.\\d*)?|.)$")
						str := onevar_reg.FindString(element.Value[:operator_index])
						left = string(element.Value[operator_index-len(str) : operator_index])
					}
					var left_element EquationElement
					if operator_strings[index] == "-" && len(left) == 0 {
						left_element = EquationElement{Type: NUM, Num: 0, Value: "0"}
					} else {
						var err error
						left_element, err = parseEquationElement(EquationElement{
							Type:  UNKNOWN,
							Value: left,
						})
						if err != nil {
							return EquationElement{}, nil
						}
					}

					right := element.Value[operator_index+1:]

					right_element, err := parseEquationElement(EquationElement{
						Type:  UNKNOWN,
						Value: right,
					})
					if err != nil {
						return EquationElement{}, nil
					}

					if operator_types[index] == POW {
						lefter := element.Value[:operator_index-len(left)]
						if len(lefter) > 0 {
							lefter_element, err := parseEquationElement(EquationElement{
								Type:  UNKNOWN,
								Value: lefter,
							})
							if err != nil {
								return EquationElement{}, nil
							}
							return EquationElement{
								Type: MULT,
								Left: &lefter_element,
								Right: &EquationElement{
									Type:  operator_types[index],
									Left:  &left_element,
									Right: &right_element,
								},
							}, nil
						}
					}

					if operator_strings[index] == "-" {
						new_right, _ := appendMinus(right_element)
						return EquationElement{
							Type:  operator_types[index],
							Left:  &left_element,
							Right: &new_right,
						}, nil
					}

					return EquationElement{
						Type:  operator_types[index],
						Left:  &left_element,
						Right: &right_element,
					}, nil
				}
			}
			// 연산자가 없으면 VAR
			return parseEquationElement(EquationElement{
				Type:  VAR,
				Value: element.Value,
			})
		}
	case VAR:
		// 사이에 생략된 곱하기를 추가한다.
		elements := make([](*EquationElement), 0)
		number_reg := regexp.MustCompile("^\\d*(\\.\\d*)?")
		number_str := number_reg.FindString(element.Value)
		if len(number_str) > 0 {
			num, err := strconv.ParseFloat(number_str, 64)
			if err != nil {
				return EquationElement{}, err
			}
			number_element := EquationElement{
				Type: NUM,
				Num:  num,
			}
			elements = append(elements, &number_element)
		}
		variables_str := element.Value[len(number_str):]
		for _, s := range variables_str {
			var_element := EquationElement{
				Type:  VAR,
				Value: string(s),
			}
			elements = append(elements, &var_element)
		}

		if len(elements) == 1 {
			return *elements[0], nil
		} else {
			result := EquationElement{
				Type:  MULT,
				Left:  elements[len(elements)-2],
				Right: elements[len(elements)-1],
			}
			for i := len(elements) - 3; i >= 0; i-- {
				prev := result
				result = EquationElement{
					Type:  MULT,
					Left:  elements[i],
					Right: &prev,
				}
			}
			return result, nil
		}
	case ADD:
		fallthrough
	case MULT:
		fallthrough
	case DIV:
		fallthrough
	case POW:
		fallthrough
	case ANY:
		return element, nil
	}
	return element, nil
}

func (e Equation) Parse() (Equation, error) {
	value, err := parseEquationElement(EquationElement{Type: UNKNOWN, Value: e.RawValue})
	if err != nil {
		return Equation{}, err
	}
	e.Value = value
	return e, nil
}

func (element EquationElement) String() string {
	switch element.Type {
	case UNKNOWN:
		return "UNKNOWN(" + element.Value + ")"
	case VAR:
		return "VAR(" + element.Value + ")"
	case NUM:
		return "NUM(" + fmt.Sprintf("%v", element.Num) + ")"
	case ADD:
		return "ADD(" + element.Left.String() + ", " + element.Right.String() + ")"
	case MULT:
		return "MULT(" + element.Left.String() + ", " + element.Right.String() + ")"
	case DIV:
		return "DIV(" + element.Left.String() + ", " + element.Right.String() + ")"
	case POW:
		return "POW(" + element.Left.String() + ", " + element.Right.String() + ")"
	case ANY:
		return "ANY(" + element.Value + ")"
	}
	return ""
}

func (element EquationElement) Reduce(variableValues map[string]float64) EquationElement {
	switch element.Type {
	case UNKNOWN:
		fallthrough
	case ANY:
		fallthrough
	case NUM:
		return element
	case VAR:
		val, exists := variableValues[element.Value]
		if exists {
			return EquationElement{Type: NUM, Num: val}
		}
		return element
	case ADD:
		left := element.Left.Reduce(variableValues)
		right := element.Right.Reduce(variableValues)
		if left.Type == NUM && right.Type == NUM {
			return EquationElement{Type: NUM, Num: left.Num + right.Num}
		}
		element.Left = &left
		element.Right = &right
		return element
	case MULT:
		left := element.Left.Reduce(variableValues)
		right := element.Right.Reduce(variableValues)
		if left.Type == NUM && right.Type == NUM {
			return EquationElement{Type: NUM, Num: left.Num * right.Num}
		}
		element.Left = &left
		element.Right = &right
		return element
	case DIV:
		left := element.Left.Reduce(variableValues)
		right := element.Right.Reduce(variableValues)
		if left.Type == NUM && right.Type == NUM {
			return EquationElement{Type: NUM, Num: left.Num / right.Num}
		}
		element.Left = &left
		element.Right = &right
		return element
	case POW:
		left := element.Left.Reduce(variableValues)
		right := element.Right.Reduce(variableValues)
		if left.Type == NUM && right.Type == NUM {
			return EquationElement{Type: NUM, Num: math.Pow(left.Num, right.Num)}
		}
		element.Left = &left
		element.Right = &right
		return element
	}
	return element
}

func (element EquationElement) isEqual(target EquationElement) bool {
	if element.Type == ANY || target.Type == ANY {
		return true
	}
	if element.Type != target.Type {
		return false
	}
	switch element.Type {
	case UNKNOWN:
		return true
	case NUM:
		return element.Num == target.Num
	case VAR:
		return element.Value == target.Value
	case ADD:
		fallthrough
	case MULT:
		return (element.Left.isEqual(*target.Left) && element.Right.isEqual(*target.Right)) || element.Right.isEqual(*target.Left) && element.Left.isEqual(*target.Right)
	case DIV:
		fallthrough
	case POW:
		return element.Left.isEqual(*target.Left) && element.Right.isEqual(*target.Right)
	}
	return false
}

func (element EquationElement) FindPattern(pattern Equation) (EquationElement, error) {
	if !pattern.Parsed {
		pattern, err := pattern.Parse()
		if err != nil {
			return EquationElement{}, err
		}
		return element.FindPatternByElement(pattern.Value), nil
	}
	return element.FindPatternByElement(pattern.Value), nil
}

func (element EquationElement) FindPatternByElement(pattern EquationElement) EquationElement {
	if element.isEqual(pattern) {
		return element
	}
	switch element.Type {
	case UNKNOWN:
		fallthrough
	case NUM:
		fallthrough
	case VAR:
		return EquationElement{Type: UNKNOWN}
	case ADD:
		fallthrough
	case MULT:
		fallthrough
	case DIV:
		fallthrough
	case POW:
		left_result := element.Left.FindPatternByElement(pattern)
		if left_result.Type != UNKNOWN {
			return left_result
		}
		return element.Right.FindPatternByElement(pattern)
	}
	return EquationElement{Type: UNKNOWN}
}

func main() {
	eq := Equation{RawValue: "x^2-2x+3"}
	var err error
	eq, err = eq.Parse()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(eq.Value)
	m := make(map[string]float64)
	m["x"] = 10
	fmt.Println(eq.Value.Reduce(m))
	pattern := EquationElement{Type: MULT, Left: &EquationElement{Type: ANY}, Right: &EquationElement{Type: VAR, Value: "x"}}
	fmt.Println(pattern)
	result := eq.Value.FindPatternByElement(pattern)
	if result.Left.Type == NUM {
		fmt.Println(result.Left.Num)
	} else {
		fmt.Println(result.Right.Num)
	}
}
