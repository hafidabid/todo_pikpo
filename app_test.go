package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"testing"
	"time"
	"todo_pikpo/config"
	"todo_pikpo/controllers"
	"todo_pikpo/database"
	model "todo_pikpo/database/models"
	myGrpc "todo_pikpo/grpc"

	pb "todo_pikpo/grpc/proto"

	midw "todo_pikpo/middleware"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AppTest struct {
	suite.Suite
	grpc myGrpc.GrpcServer
	db   *database.Database
	cnt  controllers.TodoController
	conf config.ConfigApp
}

func (s *AppTest) grpcRunner() (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.conf.Port))
	if err != nil {
		panic(err)
	}
	mdl := midw.NewMiddleware(s.conf)
	sr := grpc.NewServer(
		grpc.UnaryInterceptor(mdl.UnaryAuth),
		grpc.StreamInterceptor(mdl.StreamAuth),
	)
	pb.RegisterTodoServiceServer(sr, &s.grpc)
	pb.RegisterStreamServiceServer(sr, &s.grpc)

	log.Printf("ToDo Service started with gRPC on port %d\n", s.conf.Port)
	if err = sr.Serve(lis); err != nil {
		panic(err)
	}

	return lis, err
}

func (s *AppTest) SetupSuite() {
	c, err := config.NewAppConfig(".env.test")
	if err != nil {
		s.T().Error("Failed to create config:", err)
		return
	}

	s.conf = c

	db, err := database.NewDatabase(c)
	if err != nil {
		s.T().Error("Failed to create database:", err)
		return
	}
	s.db = &db

	cnt, err := controllers.CreateTodoController(&db)
	if err != nil {
		s.T().Error("Failed to create controller:", err)
		return
	}
	s.cnt = cnt

	s.grpc = myGrpc.StartGrpc(&cnt)
}

func (s *AppTest) createDummyData() {
	s.cnt.AddTodo(model.TodoModel{
		Author:      "james",
		Title:       "test this is title",
		Description: "lorem ipsom dolom amet",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(1 * time.Hour),
	})
	s.cnt.AddTodo(model.TodoModel{
		Author:      "robert",
		Title:       "jakarta unit test",
		Description: "lorem ipsom dolom amet",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(1 * time.Hour),
	})
	s.cnt.AddTodo(model.TodoModel{
		Author:      "ali",
		Title:       "singapore is awesome",
		Description: "lorem ipsom dolom amet",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(1 * time.Hour),
	})
}

func (s *AppTest) TearDownSuite() {
	s.db.Flush()
}

func (s *AppTest) SetupTest() {
	s.conf.Port += uint16(rand.Intn(500))
	s.db.Migrate()

}

func (s *AppTest) TearDownTest() {
	s.db.Flush()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(AppTest))
}

func (s *AppTest) TestRPC1() {
	a := s.Suite.Assert()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+s.conf.EncryptKey)
	s.createDummyData()
	go func() {
		l, e := s.grpcRunner()
		if e != nil {
			s.Suite.T().Error()
		}
		defer l.Close()
	}()

	cc, err := grpc.Dial(fmt.Sprintf(":%d", s.conf.Port), grpc.WithInsecure())
	if err != nil {
		s.T().Error(err)
	}
	defer cc.Close()

	client := pb.NewTodoServiceClient(cc)

	resp, err := client.GetTodo(ctx, &pb.FilterRequest{})
	a.Equal(err, nil)
	a.Equal(resp.GetIsOk(), true)
	a.Equal(len(resp.GetValue()), 3)

	// Test with pagination
	resp, err = client.GetTodo(ctx, &pb.FilterRequest{
		Limit: 2,
	})
	a.Equal(err, nil)
	a.Equal(resp.GetIsOk(), true)
	a.Equal(len(resp.GetValue()), 2)

	resp, err = client.GetTodo(ctx, &pb.FilterRequest{
		Limit: 2,
		Page:  1,
	})
	a.Equal(err, nil)
	a.Equal(resp.GetIsOk(), true)
	a.Equal(len(resp.GetValue()), 1)

	// Test with filter

	resp, err = client.GetTodo(ctx, &pb.FilterRequest{
		Title: "jakarta unit test",
	})
	a.Equal(err, nil)
	a.Equal(resp.GetIsOk(), true)
	a.Equal(len(resp.GetValue()), 1)

	resp, err = client.GetTodo(ctx, &pb.FilterRequest{
		Author: "james",
	})
	a.Equal(err, nil)
	a.Equal(resp.GetIsOk(), true)
	a.Equal(len(resp.GetValue()), 1)

}

func (s *AppTest) TestRPC2() {
	a := s.Suite.Assert()

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+s.conf.EncryptKey)
	s.createDummyData()
	go func() {
		l, e := s.grpcRunner()
		if e != nil {
			s.Suite.T().Error()
		}
		defer l.Close()
	}()

	cc, err := grpc.Dial(fmt.Sprintf(":%d", s.conf.Port), grpc.WithInsecure())
	if err != nil {
		s.T().Error(err)
	}
	defer cc.Close()

	client := pb.NewTodoServiceClient(cc)
	resp, err := client.GetTodo(ctx, &pb.FilterRequest{})
	a.Equal(err, nil)
	a.Equal(resp.GetIsOk(), true)

	resp2, err := client.GetOneTodo(ctx, &pb.IdQuery{Id: "1"})
	a.Equal(err, nil)
	a.Equal(resp2.GetIsOk(), false)
	a.Equal(int(resp2.GetError().GetCode()), 404)

	resp2, err = client.GetOneTodo(ctx, &pb.IdQuery{Id: resp.GetValue()[0].GetId()})
	a.Equal(err, nil)
	a.Equal(resp2.GetIsOk(), true)
	a.Equal(resp2.GetValue().GetId(), resp.GetValue()[0].GetId())

}

func (s *AppTest) TestRPC3() {
	a := s.Suite.Assert()

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+s.conf.EncryptKey)
	s.createDummyData()
	go func() {
		l, e := s.grpcRunner()
		if e != nil {
			s.Suite.T().Error()
		}
		defer l.Close()
	}()

	cc, err := grpc.Dial(fmt.Sprintf(":%d", s.conf.Port), grpc.WithInsecure())
	if err != nil {
		s.T().Error(err)
	}
	defer cc.Close()

	client := pb.NewTodoServiceClient(cc)

	data, err := client.GetTodo(ctx, &pb.FilterRequest{})
	a.Equal(err, nil)
	a.Equal(data.GetIsOk(), true)

	resp, err := client.DeleteTodo(ctx, &pb.IdQuery{
		Id: data.GetValue()[0].GetId(),
	})
	a.Equal(err, nil)
	a.Equal(resp.GetIsOk(), true)

	data, err = client.GetTodo(ctx, &pb.FilterRequest{})
	a.Equal(err, nil)
	a.Equal(data.GetIsOk(), true)
	a.Equal(len(data.GetValue()), 2)

	resp, err = client.DeleteTodo(ctx, &pb.IdQuery{
		Id: "fasfad",
	})
	a.Equal(resp.GetIsOk(), false)
	a.Equal(int(resp.GetError().GetCode()), 404)
}

func (s *AppTest) TestRPC4() {
	a := s.Suite.Assert()

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+s.conf.EncryptKey)
	s.createDummyData()
	go func() {
		l, e := s.grpcRunner()
		if e != nil {
			s.Suite.T().Error()
		}
		defer l.Close()
	}()

	cc, err := grpc.Dial(fmt.Sprintf(":%d", s.conf.Port), grpc.WithInsecure())
	if err != nil {
		s.T().Error(err)
	}
	defer cc.Close()

	client := pb.NewTodoServiceClient(cc)

	resp, err := client.AddTodo(ctx, &pb.AddRequest{
		Title:       "jakarta unit test2",
		Author:      "james roberto",
		Description: "just description",
		IsDone:      true,
		StartDate:   uint64(time.Now().Unix()),
		EndDate:     uint64(time.Now().Add(24 * time.Hour).Unix()),
	})
	a.Equal(resp.GetIsOk(), true)
	a.Equal(resp.GetValue().GetIsDone(), false)
	a.Equal(resp.GetValue().GetTitle(), "jakarta unit test2")
	a.Equal(resp.GetValue().GetAuthor(), "james roberto")

	data, err := client.GetTodo(ctx, &pb.FilterRequest{})
	a.Equal(err, nil)
	a.Equal(data.GetIsOk(), true)
	a.Equal(len(data.GetValue()), 4)

	resp, err = client.AddTodo(ctx, &pb.AddRequest{
		Title:       "jakarta unit test3",
		Author:      "james roberto",
		Description: "just description",
		IsDone:      true,
		StartDate:   uint64(time.Now().Unix()),
		EndDate:     uint64(time.Now().Unix()),
	})
	a.Equal(resp.GetIsOk(), false)
	a.Equal(int(resp.GetError().GetCode()), 400)
}

func (s *AppTest) TestRPC5() {
	a := s.Suite.Assert()

	s.createDummyData()
	go func() {
		l, e := s.grpcRunner()
		if e != nil {
			s.Suite.T().Error()
		}
		defer l.Close()
	}()

	cc, err := grpc.Dial(fmt.Sprintf(":%d", s.conf.Port), grpc.WithInsecure())
	if err != nil {
		s.T().Error(err)
	}
	defer cc.Close()

	client := pb.NewTodoServiceClient(cc)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+s.conf.EncryptKey)
	data, err := client.GetTodo(ctx, &pb.FilterRequest{})
	a.Equal(err, nil)
	a.Equal(data.GetIsOk(), true)
	a.Equal(len(data.GetValue()), 3)

	for _, v := range data.GetValue() {
		resp, err := client.EditTodo(ctx, &pb.EditRequest{
			Id: &pb.IdQuery{Id: v.GetId()},
			Data: &pb.AddRequest{
				Title:       "changed title",
				Author:      "changed author",
				Description: "changed description",
				IsDone:      true,
				StartDate:   v.GetStartDate(),
				EndDate:     v.GetEndDate(),
			},
		})
		a.Equal(err, nil)
		a.Equal(resp.GetIsOk(), true)
		a.Equal(resp.GetValue().GetTitle(), "changed title")
		a.Equal(resp.GetValue().GetAuthor(), "changed author")
		a.Equal(resp.GetValue().GetDescription(), "changed description")
		a.Equal(resp.GetValue().GetIsDone(), true)
	}

	for _, v := range data.GetValue() {
		resp, err := client.EditTodo(ctx, &pb.EditRequest{
			Id: &pb.IdQuery{Id: v.GetId()},
			Data: &pb.AddRequest{
				Title:       "changed title",
				Author:      "changed author",
				Description: "changed description",
				IsDone:      true,
				StartDate:   v.GetEndDate(),
				EndDate:     v.GetStartDate(),
			},
		})
		a.Equal(err, nil)
		a.Equal(resp.GetIsOk(), false)
		a.Equal(int(resp.GetError().GetCode()), 400)
	}
}

func (s *AppTest) TestRPCStream() {
	a := s.Suite.Assert()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+s.conf.EncryptKey)
	s.createDummyData()
	go func() {
		l, e := s.grpcRunner()
		if e != nil {
			s.Suite.T().Error()
		}
		defer l.Close()
	}()

	cc, err := grpc.Dial(fmt.Sprintf(":%d", s.conf.Port), grpc.WithInsecure())
	if err != nil {
		s.T().Error(err)
	}
	defer cc.Close()

	client := pb.NewStreamServiceClient(cc)
	stream, err := client.GetStreamingTodo(ctx, &pb.FilterRequest{})
	a.Equal(err, nil)
	var c = 0
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		a.Equal(err, nil)
		a.Greater(len(msg.GetId()), 1)
		c += 1
	}
	a.Equal(c, 3)

}

func (s *AppTest) TestRPCNonAuth() {
	a := s.Suite.Assert()

	s.createDummyData()
	go func() {
		l, e := s.grpcRunner()
		if e != nil {
			s.Suite.T().Error()
		}
		defer l.Close()
	}()

	cc, err := grpc.Dial(fmt.Sprintf(":%d", s.conf.Port), grpc.WithInsecure())
	if err != nil {
		s.T().Error(err)
	}
	defer cc.Close()

	client := pb.NewTodoServiceClient(cc)
	_, err = client.GetTodo(context.Background(), &pb.FilterRequest{})
	a.NotEqual(err, nil)
}
