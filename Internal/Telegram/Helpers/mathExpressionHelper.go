package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

func EvaluateExpression(expr string) (float64, error) {
	parts := strings.Fields(expr)
	if len(parts) != 3 {
		return 0, fmt.Errorf("Expression must have exactly 3 parts (e.g. '2+2')")
	}

	a, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("First number is invalid.")
	}

	b, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, fmt.Errorf("Second number is invalid.")
	}

	op := parts[1]
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("Division by zero.")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("Unsupported operation '%s' (use +, -, *, /)", op)
	}
}
