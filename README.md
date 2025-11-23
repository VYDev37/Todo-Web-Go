# todo-grpc-go

## Description: Only a simple Todo REST API
## Project only for my Go learning purpose.

## Progression
- [x] Basic Server 
- [x] GET Todo Router
- [x] POST Todo Router
- [x] PUT Todo Router
- [x] DELETE Todo Router
- [x] Massive DELETE Todo
- [x] The frontend (for the use of API)

## Setup
1. Open terminal in current working directory.
2. Download the required modules with `go mod download` or `go mod tidy`.
3. Done.

## In case you'd like to modify the structure
1. Download [protoc](https://github.com/protocolbuffers/protobuf/releases) and introduce it to environment variable.
2. Modify the message structure in `protobuf/tasks.proto`.
3. Copy and run this command to generate the protobuf output: `protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/todo.proto`.
4. Done.

## Running
- Run the server with `go run .` command.

## Changed
- PATCH Router + Mark as done -> Massive Delete + another DELETE route