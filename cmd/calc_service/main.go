package main

import (
	"calc_service/internal/calculator"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/calculate", calculator.CalculateHandler)
	http.ListenAndServe(":8080", nil)
}
