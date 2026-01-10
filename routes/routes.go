package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	todo "todo-rest-go/todo"
)

type AddTodoBody struct {
	Name string `json:"name"`
	Due  string `json:"due"`
}

type UpdateTodoBody struct {
	AddTodoBody      // embedded
	Done        bool `json:"done"`
}

type APIServer struct {
	Manager *todo.TaskManager
}

func AllowCORS(next http.Handler) http.Handler {
	allowedOrigins := map[string]bool{ // domain_name: can access / not
		"http://localhost:5173": true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Lanjut ke handler asli
		next.ServeHTTP(w, r)
	})
}

func CreateAPIServer(manager *todo.TaskManager) *APIServer {
	return &APIServer{Manager: manager}
}

func SendMessage(res http.ResponseWriter, status int, message string) error {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	return json.NewEncoder(res).Encode(map[string]string{"message": message}) // json structure for message
}

func (server *APIServer) HandleGetTodo(res http.ResponseWriter, req *http.Request) {
	taskList := server.Manager.Get()

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(res).Encode(taskList); err != nil {
		http.Error(res, fmt.Sprintf("Failed to retrieve data: %v.", err), http.StatusInternalServerError)
	}
}

func (server *APIServer) HandleAddTodo(res http.ResponseWriter, req *http.Request) {
	var data AddTodoBody

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		http.Error(res, fmt.Sprintf("Failed to retrieve data: %v.", err), http.StatusBadRequest)
		return
	}

	if err := server.Manager.Add(data.Name, data.Due); err != nil {
		http.Error(res, fmt.Sprintf("Failed to add task: %v.", err), http.StatusInternalServerError)
		return
	}

	SendMessage(res, http.StatusAccepted, "Added task to the list.")
}

func (server *APIServer) HandleDeleteTodo(res http.ResponseWriter, req *http.Request) {
	rawId := req.PathValue("id")
	id, err := strconv.Atoi(rawId)

	if err != nil {
		http.Error(res, "ID must be number.", http.StatusBadRequest)
		return
	}

	if err := server.Manager.Remove(int16(id)); err != nil {
		http.Error(res, fmt.Sprintf("Failed to delete task: %v.", err), http.StatusBadRequest)
		return
	}

	SendMessage(res, http.StatusAccepted, "Task removed from the list.")
}

func (server *APIServer) HandleDeleteAll(res http.ResponseWriter, req *http.Request) {
	if err := server.Manager.RemoveAll(); err != nil {
		http.Error(res, fmt.Sprintf("Failed to purge task: %v.", err), http.StatusBadRequest)
		return
	}

	SendMessage(res, http.StatusAccepted, "Task purged.")
}

func (server *APIServer) HandleUpdateTodo(res http.ResponseWriter, req *http.Request) {
	var data UpdateTodoBody

	rawId := req.PathValue("id")
	id, err := strconv.Atoi(rawId)

	if err != nil {
		http.Error(res, "ID must be number.", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		http.Error(res, fmt.Sprintf("Failed to read data: %v.", err), http.StatusBadRequest)
		return
	}

	if err := server.Manager.Update(int16(id), data.Name, data.Due, data.Done); err != nil {
		http.Error(res, fmt.Sprintf("Failed to update task: %v.", err), http.StatusBadRequest)
		return
	}

	SendMessage(res, http.StatusAccepted, "Task updated.")
}

func (server *APIServer) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "Hello world!")
	})

	mux.HandleFunc("GET /todos", server.HandleGetTodo)
	mux.HandleFunc("POST /add-todo", server.HandleAddTodo)
	mux.HandleFunc("DELETE /todo/{id}", server.HandleDeleteTodo)
	mux.HandleFunc("DELETE /todos", server.HandleDeleteAll)
	mux.HandleFunc("PUT /todo/{id}", server.HandleUpdateTodo)
}
