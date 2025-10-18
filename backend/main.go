package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type Todo struct {
	Id        string `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var (
	todos     []Todo
	todoMutex sync.Mutex
)

func generateRandomId() string {
	return uuid.New().String()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{
		Status:  "Ok",
		Message: "API health is ok",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.Method)
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	case "POST":
		var newTodo Todo
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Unable to read from request body", http.StatusBadRequest)
			return
		}

		// fmt.Println(body)
		err = json.Unmarshal(body, &newTodo)
		if err != nil || newTodo.Task == "" {
			http.Error(w, "No inputs found", http.StatusBadRequest)
			return
		}

		newTodo.Id = generateRandomId()

		todoMutex.Lock()
		todos = append(todos, newTodo)
		todoMutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTodo)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func todoByIdHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/todos", todosHandler)
	http.HandleFunc("/todos/", todoByIdHandler)

	fmt.Println("App is running in PORT 3000")

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println("Error in starting the application", err)
	}
}
