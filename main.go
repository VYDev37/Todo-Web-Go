package main

import (
	"fmt"
	"log"
	"net/http"
	pb "todo-grpc-go/protobuf"
	"todo-grpc-go/routes"
	model "todo-grpc-go/todo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverPort uint16 = 8080
var grpcPort uint16 = 50051

func main() {
	manager := model.TaskManager{}
	manager.SetFile("tasks.json")

	if err := manager.Load(); err != nil {
		fmt.Printf("Warning: Failed to load data with error: %v\n", err)
		return
	}

	go func() {
		server := routes.GRPCServer{}
		err := server.Run(&manager, grpcPort)

		if err != nil {
			log.Fatalf("An error occured when trying to start GRPC Server: %v\n", err)
			return
		}
	}()

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials())) // client port
	if err != nil {
		log.Fatalf("An error occured when trying to make new grpc client: %v\n", err)
		return
	}

	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	httpServer := &routes.HTTPServer{Client: client}

	mux := http.NewServeMux()
	httpServer.RegisterRoutes(mux)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), routes.AllowCORS(mux)); err != nil {
		log.Fatalf("An error occured when trying to run http: %v\n", err)
		return
	}

	fmt.Printf("HTTP is running in port %d.\n", serverPort)
}
