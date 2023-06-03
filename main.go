package main

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	myGrpc "todo_pikpo/grpc"
)
import pb "todo_pikpo/grpc/proto"

func main() {
	lis, err := net.Listen("tcp", ":14049")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, &myGrpc.GrpcServer{})
	pb.RegisterStreamServiceServer(s, &myGrpc.GrpcServer{})

	log.Println("ToDo Service started with gRPC on port " + "14049")
	if err = s.Serve(lis); err != nil {
		panic(err)
	}

	defer lis.Close()
}
