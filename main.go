package main

import (
	"fmt"
	"net"
	"todo_pikpo/config"
	"todo_pikpo/controllers"
	"todo_pikpo/database"
	myGrpc "todo_pikpo/grpc"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "todo_pikpo/grpc/proto"
	midw "todo_pikpo/middleware"
)

func main() {

	conf, err := config.NewAppConfig(".env")
	if err != nil {
		log.Error("something wrong while loading app config -> ", err)
		panic(err)
	}

	db, err := database.NewDatabase(conf)
	if err != nil {
		log.Error("something wrong while loading app database -> ", err)
		panic(err)
	}

	if err = db.Migrate(); err != nil {
		log.Error("something wrong while loading app database migration -> ", err)
		panic(err)
	}

	ctrl, err := controllers.CreateTodoController(&db)
	if err != nil {
		log.Error("something wrong while creating app controller -> ", err)
		panic(err)
	}

	gService := myGrpc.StartGrpc(&ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		panic(err)
	}

	mdl := midw.NewMiddleware(conf)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(mdl.UnaryAuth),
		grpc.StreamInterceptor(mdl.StreamAuth),
	)
	pb.RegisterTodoServiceServer(s, &gService)
	pb.RegisterStreamServiceServer(s, &gService)

	log.Printf("ToDo Service started with gRPC on port %d\n", conf.Port)
	if err = s.Serve(lis); err != nil {
		panic(err)
	}

	defer lis.Close()
}
