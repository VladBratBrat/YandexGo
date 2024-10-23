package main

import (
	"errors"
	"strings"
)

var (
	unclosedBracket = errors.New("There is unclosed bracket")
	weirdSymbols    = errors.New("There are odd symbols")
	weirdPosition   = errors.New("Possition of symbols is weird")
	zeroInBgein     = errors.New("A number starts with a 0")
	divisionByZero  = errors.New("Division by zero")
	operands        = map[rune]bool{'+': true, '-': true, '*': true, '/': true}
	digit           = map[rune]bool{'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true, '7': true, '8': true, '9': true}
)

func SkipUnderscore(expression string) string {
	return strings.ReplaceAll(expression, " ", "")
}

func CheckBalance(expression string) error {
	balance := 0
	for _, val := range expression {
		switch val {
		case '(':
			balance += 1
		case ')':
			balance -= 1
		}
		if balance < 0 {
			return unclosedBracket
		}
	}
	if balance != 0 {
		return unclosedBracket
	}
	return nil
}

func CheckSymbols(expression string) error {
	for _, val := range expression {
		if !(digit[val] || operands[val] || val == '(' || val == ')') {
			return weirdSymbols
		}
	}
	return nil
}

func CheckPos(expression string) error {
	// operands in odd places
	for i, val := range expression {
		if operands[val] && (i == 0 || i == len(expression)-1 || expression[i+1] == ')' || expression[i-1] == '(' || operands[rune(expression[i-1])] || operands[rune(expression[i+1])]) {
			return weirdPosition
		}
		if digit[rune(expression[i])] && i > 0 && expression[i-1] == ')' {
			return weirdPosition
		}
		if i > 0 && expression[i] == ')' && expression[i-1] == '(' {
			return weirdPosition
		}
	}
	// zero in begin
	for i, val := range expression {
		if i < len(expression)-1 && digit[rune(expression[i+1])] && val == '0' && (i == 0 || !digit[rune(expression[i-1])]) {
			return zeroInBgein
		}
	}
	return nil
}

var exp string
var DivisionByZero error

func RecursFormula(l int) (int, float64) {
	if DivisionByZero != nil {
		return -1, -1
	}
	i := l
	vals := make([]float64, 0)
	oper := make([]int, 0)                           // 1 <=> +, 2 <=> -
	value, lst, ind, filled := float64(0.0), 0, 0, 0 // 3 <=> *, 4 <=> /
	for ; i < len(exp); i += 1 {
		if digit[rune(exp[i])] {
			value = value*10 + float64(rune(exp[i])-rune('0'))
			filled = 1
		} else {
			if filled == 1 && ind >= len(vals) {
				vals = append(vals, value)
			}
			filled = 0
			if exp[i] != '(' {
				if lst == 3 {
					vals[ind] *= value
					value = 0
				} else if lst == 4 {
					if value == 0 {
						DivisionByZero = divisionByZero
						return -1, -1
					}
					vals[ind] /= value
					value = 0
				}
			}
			if exp[i] == '\n' || exp[i] == ')' {
				break
			}
			if operands[rune(exp[i])] {
				if exp[i] == '+' {
					oper = append(oper, 1)
					ind += 1
					lst = 0
				} else if exp[i] == '-' {
					oper = append(oper, 2)
					ind += 1
					lst = 0
				} else if exp[i] == '*' {
					lst = 3
				} else {
					lst = 4
				}
				value = 0
			} else {
				r, res := RecursFormula(i + 1)
				filled = 1
				i = r
				value = res
				if lst == 3 {
					vals[ind] *= value
				} else if lst == 4 {
					if value == 0 {
						DivisionByZero = divisionByZero
						return -1, -1
					}
					vals[ind] /= value
				}
				lst = 0
			}
		}
	}
	ans := 0.0
	for j := 0; j <= ind; j += 1 {
		if j == 0 {
			ans += vals[j]
		} else {
			if oper[j-1] == 1 {
				ans += vals[j]
			} else {
				ans -= vals[j]
			}
		}
	}
	return i, ans
}

func Calc(expression string) (float64, error) {
	input := expression
	if len(expression) == 0 {
		return 0, errors.New("Empty statement")
	}
	text := ""
	for i := 0; i < len(input); i += 1 {
		if i == 0 && input[i] == '(' {
			text = text + "1*("
		} else if (i > 0 && input[i] == '-' && input[i-1] == '(') || (i == 0 && input[i] == '-') {
			text = text + "0-" // zero before statement
		} else if i > 0 && input[i] == '(' && input[i-1] == '(' {
			text = text + "1*(" // multiplication before statemnt
		} else if i > 0 && input[i] == '(' && digit[rune(input[i-1])] {
			text = text + "*(" // multiplication before statemnt
		} else {
			text = text + string(input[i])
		}
	}
	expression = text

	expression = SkipUnderscore(expression)

	err := CheckBalance(expression)
	if err != nil {
		return 0, err
	}
	err = CheckSymbols(expression)
	if err != nil {
		return 0, err
	}
	err = CheckPos(expression)
	if err != nil {
		return 0, err
	}
	exp = expression + "\n"
	_, ans := RecursFormula(0)
	if DivisionByZero != nil {
		return 0, DivisionByZero
	}
	return ans, nil
}
