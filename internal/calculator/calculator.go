package calculator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// RequestBody структура для входящего запроса
type RequestBody struct {
	Expression string `json:"expression"`
}

// ResponseBody структура для исходящего ответа
type ResponseBody struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// CalculateHandler обрабатывает запросы на вычисление выражений
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := calculate(reqBody.Expression)
	if err != nil {
		if err.Error() == "invalid expression" {
			http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		} else {
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	response := ResponseBody{Result: fmt.Sprintf("%f", result)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// calculate выполняет вычисление арифметического выражения
func calculate(expression string) (float64, error) {
	// Удаляем пробелы
	expression = strings.ReplaceAll(expression, " ", "")
	if strings.ContainsAny(expression, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return 0, fmt.Errorf("invalid expression")
	}
	

	postfix, err := infixToPostfix(expression)
	if err != nil {
		return 0, err
	}

	return evaluatePostfix(postfix)
}

func infixToPostfix(expression string) ([]string, error) {
	var output []string
	var stack []string

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
		"(": 0,
	}

	for i := 0; i < len(expression); i++ {
		token := string(expression[i])

		if isNumber(token) {
			// Если токен - число, добавляем его в выходной массив
			num, j := readNumber(expression, i)
			output = append(output, num)
			i = j - 1
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			stack = stack[:len(stack)-1] 
		} else if precedence[token] > 0 {
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[token] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else {
			return nil, fmt.Errorf("invalid token: %s", token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func evaluatePostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if isNumber(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				result = a / b
			default:
				return 0, fmt.Errorf("invalid operator: %s", token)
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}
	return stack[0], nil
}

// isNumber проверяет, является ли строка числом
func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// readNumber считывает число из строки и возвращает его вместе с индексом следующего символа
func readNumber(expression string, start int) (string, int) {
	end := start
	for end < len(expression) && (isNumber(string(expression[end])) || expression[end] == '.') {
		end++
	}
	return expression[start:end], end
}
