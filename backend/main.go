package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{
		Status:  "Ok",
		Message: "API health is ok",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func main() {
	http.HandleFunc("/health", healthHandler)

	fmt.Println("App is running in PORT 3000")

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println("Error in starting the application", err)
	}
}
