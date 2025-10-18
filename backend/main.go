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
	id := r.URL.Path[len("/todos/"):]

	todoMutex.Lock()
	defer todoMutex.Unlock()

	for i, todo := range todos {
		if todo.Id == id {
			switch r.Method {
			case "GET":
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(todo)
			case "PUT":
				var updatedTodo Todo
				body, err := ioutil.ReadAll(r.Body)

				if err != nil {
					http.Error(w, "Unable to read from request body", http.StatusBadRequest)
					return
				}

				err = json.Unmarshal(body, &updatedTodo)

				if err != nil || updatedTodo.Task == "" {
					http.Error(w, "Invalid input", http.StatusBadRequest)
					return
				}

				todos[i].Task = updatedTodo.Task
				todos[i].Completed = updatedTodo.Completed

				json.NewEncoder(w).Encode(todos[i])
			case "DELETE":
				deletedTodoId := todo.Id
				todos = append(todos[:i], todos[i+1:]...)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"message": "Todo with Id" + deletedTodoId + "is successfully deleted"})

			default:
				http.Error(w, "Method not allowed", http.StatusBadRequest)
			}
		}
	}

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
