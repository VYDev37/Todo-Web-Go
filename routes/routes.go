package routes

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "todo-grpc-go/protobuf"
	todo "todo-grpc-go/todo"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedTaskServiceServer
	Manager *todo.TaskManager
}

func WritePb(task todo.Task) *pb.Task {
	return &pb.Task{
		ID:   int32(task.ID),
		Name: task.Name,
		Due:  task.Due,
		Done: task.Done,
	}
}

func CreateGRPCService(manager *todo.TaskManager) *GRPCServer {
	return &GRPCServer{Manager: manager}
}

func (server *GRPCServer) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	taskList := server.Manager.Get()

	var latestId int32 = 1
	if len(taskList) > 0 {
		latestId = int32(taskList[len(taskList)-1].ID)
	}

	var tasks []*pb.Task
	for _, task := range taskList {
		tasks = append(tasks, WritePb(task))
	}

	return &pb.GetTasksResponse{
		Id:      latestId,
		Message: "Parsed data.",
		Tasks:   tasks,
	}, nil
}

func (server *GRPCServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	name := req.GetName() // body
	due := req.GetDue()   // body

	latestId, err := server.Manager.Add(name, due)
	if err != nil {
		return &pb.CreateTaskResponse{
			Id:      -1,
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.CreateTaskResponse{
		Id:      int32(latestId),
		Success: true,
		Message: "Task added.",
	}, nil
}

func (server *GRPCServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	id := req.GetID() // params

	if err := server.Manager.Remove(int16(id)); err != nil {
		return &pb.DeleteTaskResponse{
			Id:      -1,
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteTaskResponse{
		Id:      id,
		Success: true,
		Message: fmt.Sprintf("Task #%d removed from the list.", id),
	}, nil
}

func (server *GRPCServer) DeleteTasks(ctx context.Context, req *pb.DeleteTasksRequest) (*pb.DeleteTasksResponse, error) {
	if err := server.Manager.RemoveAll(); err != nil {
		return &pb.DeleteTasksResponse{
			Id:      -1,
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteTasksResponse{
		Id:      0,
		Success: true,
		Message: "Successfully purged tasks.",
	}, nil
}

func (server *GRPCServer) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	id := req.GetID()     // params
	name := req.GetName() // body
	due := req.GetDue()   // body
	done := req.GetDone() // body

	if err := server.Manager.Update(int16(id), name, due, done); err != nil {
		return &pb.UpdateTaskResponse{
			Id:      -1,
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.UpdateTaskResponse{
		Id:      -1,
		Success: true,
		Message: fmt.Sprintf("Updated task #%d.", id),
	}, nil
}

func (server *GRPCServer) Run(manager *todo.TaskManager, serverPort uint16) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		log.Fatalf("Failed to start server in port %d: %v\n", serverPort, err)
		return err
	}

	gRpcServer := grpc.NewServer()
	todoService := CreateGRPCService(manager)

	pb.RegisterTaskServiceServer(gRpcServer, todoService)
	fmt.Printf("Server is running in port %d.\n", serverPort)

	if err := gRpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to run server in port %d: %v\n", serverPort, err)
		return err
	}

	return nil
}
