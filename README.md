# Simple Todo apps with go and gRPC

This guide will help you set up a  and run this project.

## About this project
this is simple todo service app project using golang, postgres, redis, and gRPC as part of pikpo software engineer test
for candidate Hafid Abi D.

## Prerequisites

Before proceeding, make sure you have the following installed on your system:

- Go: [Official Installation Guide](https://golang.org/doc/install)
- PostgreSQL: [Official Downloads](https://www.postgresql.org/download)
- Redis: [Official Downloads](https://redis.io/download)

## Features
- gRPC CRUD
- gRPC authentication with bearer token
- gRPC stream
- GORM implementation
- Dependency injection (not yet implement wire for enhance Dependency Injecton process)
- Redis caching
- Unit testing
## Setup Steps

1. Clone the repository:

   ```bash
    git clone https://github.com/hafidabid/todo_pikpo.git
   ```
   
2. Change into the project directory:
    ```bash
    cd todo_pikpo
    ```
3. Install Go dependencies:
    ```bash
    go mod download
   ```     
4. Set up dot env:
   - copy .env.example into .env for running in dev/prod mode and .env.test for testing purpose
5. Run the generate_proto.sh script:
    - Make sure you have the necessary permissions to execute the script.
    - Open a terminal and navigate to the project root directory.
    - Execute the following command:
   ```bash
    ./generate_proto.sh
    ```
   This script generates Go code from your protocol buffer definition files (.proto) and places the generated code in the appropriate Go package.
    
    if you don't have protoc installed, please follow there
   https://grpc.io/docs/languages/go/quickstart/

6. Build and run the project:
    ```bash
    go run main.go
    ```
   or you can use docker to run the app
    ```bash
    docker-compose up
    ```
