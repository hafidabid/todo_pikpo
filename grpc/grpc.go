package grpc

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
	pb "todo_pikpo/grpc/proto"
)

type GrpcServer struct {
	pb.TodoServiceServer
	pb.StreamServiceServer
}

func (gs *GrpcServer) GetTodo(ctx context.Context, filter *pb.FilterRequest) (*pb.ArrResponse, error) {

	return &pb.ArrResponse{Result: []*pb.Response{
		&pb.Response{
			Author:      "james",
			Title:       "halo halo bandung",
			Description: "afasdfasdasdf",
			IsDone:      false,
			StartDate:   0,
			EndDate:     0,
			CreatedAt:   0,
			UpdatedAt:   0,
			Id:          "1",
		},
		&pb.Response{
			Author:      "akbar",
			Title:       "bandung panas",
			Description: "afasdfasdasdf",
			IsDone:      true,
			StartDate:   0,
			EndDate:     0,
			CreatedAt:   0,
			UpdatedAt:   0,
			Id:          "2",
		},
	}}, nil
}

func (gs *GrpcServer) GetOneTodo(ctx context.Context, id *pb.IdQuery) (*pb.Response, error) {

	return &pb.Response{
		Author:      "akbar",
		Title:       "bandung panas",
		Description: "afasdfasdasdf",
		IsDone:      true,
		StartDate:   0,
		EndDate:     0,
		CreatedAt:   0,
		UpdatedAt:   0,
		Id:          "1",
	}, nil
}

func (gs *GrpcServer) GetStreamingTodo(
	filter *pb.FilterRequest,
	stream pb.StreamService_GetStreamingTodoServer,
) error {
	respList := []*pb.Response{
		&pb.Response{
			Author:      "james",
			Title:       "halo halo bandung",
			Description: "afasdfasdasdf",
			IsDone:      false,
			StartDate:   0,
			EndDate:     0,
			CreatedAt:   0,
			UpdatedAt:   0,
			Id:          "1",
		},
		&pb.Response{
			Author:      "akbar",
			Title:       "bandung panas",
			Description: "afasdfasdasdf",
			IsDone:      true,
			StartDate:   0,
			EndDate:     0,
			CreatedAt:   0,
			UpdatedAt:   0,
			Id:          "2",
		},
		&pb.Response{
			Author:      "rojak",
			Title:       "bandung panas",
			Description: "afasdfasdasdf",
			IsDone:      true,
			StartDate:   0,
			EndDate:     0,
			CreatedAt:   0,
			UpdatedAt:   0,
			Id:          "3",
		},
		&pb.Response{
			Author:      "malik",
			Title:       "bandung panas",
			Description: "afasdfasdasdf",
			IsDone:      true,
			StartDate:   0,
			EndDate:     0,
			CreatedAt:   0,
			UpdatedAt:   0,
			Id:          "4",
		},
	}

	for _, d := range respList {
		if e := stream.Send(d); e != nil {
			log.Error(e)
			return e
		}
		time.Sleep(1500 * time.Millisecond)
	}
	return nil
}

func (gs *GrpcServer) AddTodo(ctx context.Context, data *pb.AddRequest) (*pb.Response, error) {

	return &pb.Response{
		Author:      "akbar",
		Title:       "bandung panas",
		Description: "afasdfasdasdf",
		IsDone:      true,
		StartDate:   0,
		EndDate:     0,
		CreatedAt:   0,
		UpdatedAt:   0,
		Id:          "1",
	}, nil
}

func (gs *GrpcServer) EditTodo(ctx context.Context, data *pb.EditRequest) (*pb.Response, error) {

	return &pb.Response{
		Author:      "akbar",
		Title:       "bandung panas",
		Description: "afasdfasdasdf",
		IsDone:      true,
		StartDate:   0,
		EndDate:     0,
		CreatedAt:   0,
		UpdatedAt:   0,
		Id:          "1",
	}, nil
}

func (gs *GrpcServer) DeleteTodo(ctx context.Context, id *pb.IdQuery) (*pb.Response, error) {

	return &pb.Response{
		Author:      "akbar",
		Title:       "bandung panas",
		Description: "afasdfasdasdf",
		IsDone:      true,
		StartDate:   0,
		EndDate:     0,
		CreatedAt:   0,
		UpdatedAt:   0,
		Id:          "1",
	}, nil
}
