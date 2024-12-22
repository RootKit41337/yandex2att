package main

import (
	"log"
	"net/http"
	"calc_service/internal/calculator"
)

func main() {
	http.HandleFunc("/api/v1/calculate", calculator.CalculateHandler)
	log.Println("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}