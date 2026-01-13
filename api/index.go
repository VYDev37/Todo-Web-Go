package handler

import (
	"fmt"
	"net/http"
	"sync"
	routes "todo-rest-go/routes"
	model "todo-rest-go/todo"
)

var (
	manager *model.TaskManager
	once    sync.Once
	mux     *http.ServeMux
)

// Serverless handler
func Handler(res http.ResponseWriter, req *http.Request) {
	once.Do(func() {
		mux = http.NewServeMux()

		manager = &model.TaskManager{}
		if err := manager.Load(); err != nil {
			fmt.Printf("Warning: Failed to load data with error: %v", err)
		}

		server := routes.CreateAPIServer(manager)
		server.RegisterRoutes(mux)
	})

	// cors handler
	corsHandler := routes.AllowCORS(mux)
	corsHandler.ServeHTTP(res, req)

	if manager == nil {
		http.Error(res, "Failed to initialize manager and database", 500)
		return
	}
}
