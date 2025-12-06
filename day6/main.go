package main

import (
	"advent-of-code-2025/support"
	"fmt"
	"regexp"
	"strings"
)

func main() {
	input := support.LoadInput()

	fmt.Println(partOne(input))
	fmt.Println(partTwo(input))
}

func partOne(input string) int {
	lines := strings.Split(input, "\n")
	components := make([][]string, len(lines))

	whitespace := regexp.MustCompile(`\s+`)

	// Break the numbers and operators into cells...
	for i, line := range lines {
		components[i] = whitespace.Split(strings.TrimSpace(line), -1)
	}

	// Then transpose the cells such that we have a list of operands terminated by an operation
	problems := support.Map(support.Transpose(components), func(problem []string) Problem {
		return Problem{
			operator: stringOpToFuncOp(problem[len(problem)-1]),
			operands: support.SliceOfNumericStringsToSliceOfInts(problem[:len(problem)-1]),
		}
	})

	sum := 0

	for _, problem := range problems {
		sum += problem.evaluate()
	}

	return sum
}

func partTwo(input string) int {
	lines := strings.Split(input, "\n")

	problems := make([]Problem, 0)
	operands := make([]string, 0)

	// Process the input right-to-left, column-wise
	for i := len(lines[0]) - 1; i >= 0; i-- {
		number := ""

		for j := range lines {
			digit := lines[j][i]

			// If the final row contains an operator, we're at the end of the problem
			if j == len(lines)-1 {
				// We're not at the end of the problem yet
				if digit == ' ' {
					continue
				}

				problems = append(problems, Problem{
					operator: stringOpToFuncOp(string(digit)),
					operands: support.SliceOfNumericStringsToSliceOfInts(operands),
				})

				// Reset operands as we're about to start a new problem
				operands = make([]string, 0)

				// Skip the empty column separating problems
				i--

				break
			}

			// If we're not on the last row, add the next digit if we have one...
			if digit != ' ' {
				number += string(digit)
			}

			// ...then if we're at the end of a number, add it to the list of operands
			if j == len(lines)-2 {
				operands = append(operands, number)
				continue
			}
		}
	}

	sum := 0

	for _, problem := range problems {
		sum += problem.evaluate()
	}

	return sum
}

type Operator func(...int) int

type Problem struct {
	operator Operator
	operands []int
}

func (p *Problem) evaluate() int {
	return p.operator(p.operands...)
}

func mult(in ...int) int {
	out := 1
	for _, v := range in {
		out *= v
	}

	return out
}

func plus(in ...int) int {
	out := 0
	for _, v := range in {
		out += v
	}

	return out
}

func stringOpToFuncOp(op string) Operator {
	var fnOp Operator

	switch op {
	case "*":
		fnOp = mult
	case "+":
		fnOp = plus
	default:
		panic("Unimplemented operator")
	}

	return fnOp
}
