package main

import (
	"bytes"
	"calc_service/internal/calculator"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		expectedResult string
		expectedError  string
		expectedCode   int
	}{
		{
			name:           "Valid Expression",
			expression:     "2+2*2",
			expectedResult: "6.000000",
			expectedError:  "",
			expectedCode:   http.StatusOK,
		},
		{
			name:           "Valid Expression",
			expression:     "(2+2)*2",
			expectedResult: "8.000000",
			expectedError:  "",
			expectedCode:   http.StatusOK,
		},
		{
			name:           "Valid Expression",
			expression:     "2+2",
			expectedResult: "4.000000",
			expectedError:  "",
			expectedCode:   http.StatusOK,
		},
		{
			name:           "Invalid Expression",
			expression:     "2+2*2a",
			expectedResult: "",
			expectedError:  "Expression is not valid",
			expectedCode:   http.StatusUnprocessableEntity,
		},
		{
			name:           "Division by Zero",
			expression:     "2/0",
			expectedResult: "",
			expectedError:  "Internal server error",
			expectedCode:   http.StatusInternalServerError,
		},
		{
			name:           "Mismatched Parentheses",
			expression:     "(2+3",
			expectedResult: "",
			expectedError:  "mismatched parentheses",
			expectedCode:   http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(calculator.RequestBody{Expression: tt.expression})
			req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(calculator.CalculateHandler)
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}
			var response calculator.ResponseBody
			json.NewDecoder(rr.Body).Decode(&response)
			if response.Error != "" && response.Error != tt.expectedError {
				t.Errorf("handler returned unexpected error: got %v want %v", response.Error, tt.expectedError)
			}
			if response.Result != tt.expectedResult {
				t.Errorf("handler returned unexpected result: got %v want %v", response.Result, tt.expectedResult)
			}
		})
	}
}
