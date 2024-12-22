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

	// Стек для чисел и операторов
	var nums []int
	var ops []byte
	currentNum := 0

	for i := 0; i < len(expression); i++ {
		char := expression[i]

		// Если символ - цифра, формируем текущее число
		if char >= '0' && char <= '9' {
			currentNum = currentNum*10 + int(char-'0')
		}

		// Если символ - оператор или последний символ
		if char == '+' || char == '-' || char == '*' || char == '/' || char == '(' || char == ')' || i == len(expression)-1 {
			if char == ')' {
				// Обрабатываем все операции до открывающей скобки
				for len(ops) > 0 && ops[len(ops)-1] != '(' {
					currentNum = applyOperation(nums, ops, currentNum)
				}
				// Удаляем открывающую скобку
				ops = ops[:len(ops)-1]
			} else {
				// Обрабатываем операции с высоким приоритетом
				for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(char) {
					currentNum = applyOperation(nums, ops, currentNum)
				}
				// Добавляем текущее число и оператор в списки
				nums = append(nums, currentNum)
				ops = append(ops, char)
				currentNum = 0
			}
		}
	}

	// Обрабатываем оставшиеся операции
	nums = append(nums, currentNum)
	for len(ops) > 0 {
		currentNum = applyOperation(nums, ops, currentNum)
	}

	return nums[0], nil
}

// Функция для применения операции
func applyOperation(nums []int, ops []byte, currentNum int) int {
	lastNum := nums[len(nums)-1]
	nums = nums[:len(nums)-1]
	lastOp := ops[len(ops)-1]
	ops = ops[:len(ops)-1]

	switch lastOp {
	case '+':
		return lastNum + currentNum
	case '-':
		return lastNum - currentNum
	case '*':
		return lastNum * currentNum
	case '/':
		if currentNum == 0 {
			panic("division by zero")
		}
		return lastNum / currentNum
	}
	return currentNum
}

// Функция для определения приоритета операций
func precedence(op byte) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	case '(':
		return 0
	}
	return -1
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
