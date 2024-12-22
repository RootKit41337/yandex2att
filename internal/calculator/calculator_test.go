package calculator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		expression string
		expected   string
		statusCode int
	}{
		{"2 + 3 * (4 - 1)", "11", http.StatusOK},
		{"2 + 2 * a", "Expression is not valid", http.StatusUnprocessableEntity},
		{"2 + (3 * (4 - 1)", "Mismatched parentheses", http.StatusUnprocessableEntity},
		{"2 / 0", "Internal server error", http.StatusInternalServerError},
	}

	for _, test := range tests {
		reqBody, _ := json.Marshal(Request{Expression: test.expression})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CalculateHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != test.statusCode {
			t.Errorf("handler returned wrong status code: got %v want %v", status, test.statusCode)
		}

		var response Response
		json.NewDecoder(rr.Body).Decode(&response)

		if response.Error != "" {
			if response.Error != test.expected {
				t.Errorf("handler returned unexpected body: got %v want %v", response.Error, test.expected)
			}
		} else {
			if response.Result != test.expected {
				t.Errorf("handler returned unexpected body: got %v want %v", response.Result, test.expected)
			}
		}
	}
}
