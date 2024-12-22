package calculator

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func Calculate(expression string) (string, error) {
	// Проверка на наличие недопустимых символов
	if strings.ContainsAny(expression, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return "", errors.New("invalid expression")
	}

	// Вычисление выражения
	result, err := eval(expression)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(result), nil
}

func eval(expression string) (int, error) {
	// Удаляем пробелы
	expression = strings.ReplaceAll(expression, " ", "")

	var result int
	var currentNum int
	var operation byte = '+'

	for i := 0; i < len(expression); i++ {
		char := expression[i]

		if char >= '0' && char <= '9' {
			currentNum = currentNum*10 + int(char-'0')
		}

		if char == '+' || char == '-' || char == '*' || char == '/' || i == len(expression)-1 {
			switch operation {
			case '+':
				result += currentNum
			case '-':
				result -= currentNum
			case '*':
				result *= currentNum
			case '/':
				if currentNum == 0 {
					return 0, errors.New("division by zero")
				}
				result /= currentNum
			}

			operation = char
			currentNum = 0
		}
	}

	return result, nil
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := Calculate(req.Expression)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		} else {
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	response := Response{Result: result}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
