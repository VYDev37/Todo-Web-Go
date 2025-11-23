package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	pb "todo-grpc-go/protobuf"
)

type HTTPServer struct {
	Client pb.TaskServiceClient
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

func WriteJSON(res http.ResponseWriter, status int, v any) error {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	return json.NewEncoder(res).Encode(v) // json structure for message
}

func SendMessage(res http.ResponseWriter, message string, status int) error {
	return WriteJSON(res, status, map[string]string{"message": message}) // json structure for message
}

func (server *HTTPServer) HandleGetTodo(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // lewat timeout = cancel

	result, err := server.Client.GetTasks(ctx, &pb.GetTasksRequest{})
	if err != nil {
		SendMessage(res, fmt.Sprintf("Error when trying to get todos: %v.\n", err), http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	WriteJSON(res, http.StatusOK, result.Tasks)
}

func (server *HTTPServer) HandleAddTodo(res http.ResponseWriter, req *http.Request) {
	var data struct {
		Name string `json:"name"`
		Due  string `json:"due"`
	}

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		SendMessage(res, fmt.Sprintf("Failed to retrieve data: %v.", err), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := server.Client.CreateTask(ctx, &pb.CreateTaskRequest{
		Name: data.Name,
		Due:  data.Due,
	})
	if err != nil {
		SendMessage(res, fmt.Sprintf("Failed to add task: %v.", err), http.StatusInternalServerError)
		return
	}

	SendMessage(res, "Added task to the list.", http.StatusCreated)
}

func (server *HTTPServer) HandleDeleteTodo(res http.ResponseWriter, req *http.Request) {
	rawId := req.PathValue("id")
	id, err := strconv.Atoi(rawId)

	if err != nil {
		SendMessage(res, "ID must be number.", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err2 := server.Client.DeleteTask(ctx, &pb.DeleteTaskRequest{ID: int32(id)})
	if err2 != nil {
		SendMessage(res, fmt.Sprintf("Failed to delete task: %v.", err), http.StatusBadRequest)
		return
	}

	SendMessage(res, "Task removed from the list.", http.StatusOK)
}

func (server *HTTPServer) HandleDeleteAll(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := server.Client.DeleteTasks(ctx, &pb.DeleteTasksRequest{})
	if err != nil {
		SendMessage(res, fmt.Sprintf("Failed to purge task: %v.", err), http.StatusBadRequest)
		return
	}

	SendMessage(res, "Task purged.", http.StatusAccepted)
}

func (server *HTTPServer) HandleUpdateTodo(res http.ResponseWriter, req *http.Request) {
	var data struct {
		Name string `json:"name"`
		Due  string `json:"due"`
		Done bool   `json:"done"`
	}

	rawId := req.PathValue("id")
	id, err := strconv.Atoi(rawId)

	if err != nil {
		SendMessage(res, "ID must be number.", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		SendMessage(res, fmt.Sprintf("Failed to read data: %v.", err), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err2 := server.Client.UpdateTask(ctx, &pb.UpdateTaskRequest{
		ID:   int32(id),
		Name: data.Name,
		Due:  data.Due,
		Done: data.Done,
	})
	if err2 != nil {
		SendMessage(res, fmt.Sprintf("Failed to update task: %v.", err), http.StatusBadRequest)
		return
	}

	SendMessage(res, "Task updated.", http.StatusAccepted)
}

func (server *HTTPServer) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "Hello world!")
	})

	mux.HandleFunc("GET /todos", server.HandleGetTodo)
	mux.HandleFunc("POST /add-todo", server.HandleAddTodo)
	mux.HandleFunc("DELETE /todo/{id}", server.HandleDeleteTodo)
	mux.HandleFunc("DELETE /todos", server.HandleDeleteAll)
	mux.HandleFunc("PUT /todo/{id}", server.HandleUpdateTodo)

	fmt.Println("Routes registered.")
}
