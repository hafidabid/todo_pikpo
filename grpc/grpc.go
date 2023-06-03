package grpc

import (
	"context"
	"time"
	"todo_pikpo/controllers"
	model "todo_pikpo/database/models"
	pb "todo_pikpo/grpc/proto"
)

type GrpcServer struct {
	pb.TodoServiceServer
	pb.StreamServiceServer
	controller *controllers.TodoController
}

func (gs *GrpcServer) todoGetter(filter *pb.FilterRequest) ([]*pb.DataResponse, error) {
	var query = map[string]interface{}{}
	pg := 0
	limit := 10
	if len(filter.GetAuthor()) > 0 {
		query["Author"] = filter.GetAuthor()
	}
	if len(filter.GetTitle()) > 0 {
		query["Title"] = filter.GetTitle()
	}
	if filter.GetIsDone() == true && filter.GetIsDone() == false {
		query["IsDone"] = filter.GetIsDone()
	}
	if filter.GetPage() > 0 {
		pg = int(filter.GetPage())
	}
	if filter.GetLimit() > 0 {
		limit = int(filter.GetLimit())
	}

	var listOfData []*pb.DataResponse

	res, _, err := gs.controller.GetTodos(query, uint(pg), uint(limit))
	if res != nil {
		return nil, err
	}
	for _, d := range res {
		listOfData = append(listOfData, &pb.DataResponse{
			Author:      d.Author,
			Title:       d.Title,
			Description: d.Description,
			IsDone:      d.IsDone,
			StartDate:   uint64(d.StartDate.Unix()),
			EndDate:     uint64(d.EndDate.Unix()),
			CreatedAt:   uint64(d.CreatedAt.Unix()),
			UpdatedAt:   uint64(d.UpdatedAt.Unix()),
			Id:          d.Id,
		})
	}
	return listOfData, nil
}

func (gs *GrpcServer) GetTodo(ctx context.Context, filter *pb.FilterRequest) (*pb.ArrResponse, error) {
	lData, err := gs.todoGetter(filter)
	return &pb.ArrResponse{
		IsOk:  err == nil,
		Value: lData,
		Error: &pb.ErrorResponse{
			Code:    uint32(500),
			Message: err.Error(),
			Details: nil,
		},
	}, nil
}

func (gs *GrpcServer) GetOneTodo(ctx context.Context, id *pb.IdQuery) (*pb.Response, error) {
	resp, status, e := gs.controller.GetTodo(id.GetId())

	return &pb.Response{
		IsOk: e == nil,
		Value: &pb.DataResponse{
			Author:      resp.Author,
			Title:       resp.Title,
			Description: resp.Description,
			IsDone:      resp.IsDone,
			StartDate:   uint64(resp.StartDate.Unix()),
			EndDate:     uint64(resp.EndDate.Unix()),
			CreatedAt:   uint64(resp.CreatedAt.Unix()),
			UpdatedAt:   uint64(resp.UpdatedAt.Unix()),
			Id:          resp.Id,
		},
		Error: &pb.ErrorResponse{
			Code:    uint32(status),
			Message: e.Error(),
			Details: nil,
		},
	}, nil
}

func (gs *GrpcServer) GetStreamingTodo(
	filter *pb.FilterRequest,
	stream pb.StreamService_GetStreamingTodoServer,
) error {
	lData, err := gs.todoGetter(filter)
	if err != nil {
		return err
	}

	for _, d := range lData {
		if e := stream.Send(d); e != nil {
			return e
		}
	}

	return nil
}

func (gs *GrpcServer) AddTodo(ctx context.Context, data *pb.AddRequest) (*pb.Response, error) {
	res, code, err := gs.controller.AddTodo(model.TodoModel{
		Author:      data.GetAuthor(),
		Title:       data.GetTitle(),
		Description: data.GetDescription(),
		StartDate:   time.Unix(int64(data.GetStartDate()), 0),
		EndDate:     time.Unix(int64(data.GetEndDate()), 0),
	})
	return &pb.Response{
		IsOk: err == nil,
		Value: &pb.DataResponse{
			Author:      res.Author,
			Title:       res.Title,
			Description: res.Description,
			IsDone:      res.IsDone,
			StartDate:   uint64(res.StartDate.Unix()),
			EndDate:     uint64(res.EndDate.Unix()),
			CreatedAt:   uint64(res.CreatedAt.Unix()),
			UpdatedAt:   uint64(res.UpdatedAt.Unix()),
			Id:          res.Id,
		},
		Error: &pb.ErrorResponse{
			Code:    uint32(code),
			Message: err.Error(),
		},
	}, nil
}

func (gs *GrpcServer) EditTodo(ctx context.Context, data *pb.EditRequest) (*pb.Response, error) {
	res, code, err := gs.controller.EditTodo(data.GetId().GetId(), model.TodoModel{
		Author:      data.GetData().GetAuthor(),
		Title:       data.GetData().GetTitle(),
		Description: data.GetData().GetDescription(),
		IsDone:      data.GetData().GetIsDone(),
		StartDate:   time.Unix(int64(data.GetData().GetStartDate()), 0),
		EndDate:     time.Unix(int64(data.GetData().GetEndDate()), 0),
	})
	return &pb.Response{
		IsOk: err == nil,
		Value: &pb.DataResponse{
			Author:      res.Author,
			Title:       res.Title,
			Description: res.Description,
			IsDone:      res.IsDone,
			StartDate:   uint64(res.StartDate.Unix()),
			EndDate:     uint64(res.EndDate.Unix()),
			CreatedAt:   uint64(res.CreatedAt.Unix()),
			UpdatedAt:   uint64(res.UpdatedAt.Unix()),
			Id:          res.Id,
		},
		Error: &pb.ErrorResponse{
			Code:    uint32(code),
			Message: err.Error(),
		},
	}, nil
}

func (gs *GrpcServer) DeleteTodo(ctx context.Context, id *pb.IdQuery) (*pb.Response, error) {
	res, code, err := gs.controller.DeleteTodo(id.GetId())
	return &pb.Response{
		IsOk: err == nil,
		Value: &pb.DataResponse{
			Author:      res.Author,
			Title:       res.Title,
			Description: res.Description,
			IsDone:      res.IsDone,
			StartDate:   uint64(res.StartDate.Unix()),
			EndDate:     uint64(res.EndDate.Unix()),
			CreatedAt:   uint64(res.CreatedAt.Unix()),
			UpdatedAt:   uint64(res.UpdatedAt.Unix()),
			Id:          res.Id,
		},
		Error: &pb.ErrorResponse{
			Code:    uint32(code),
			Message: err.Error(),
		},
	}, nil
}

func StartGrpc(controller *controllers.TodoController) GrpcServer {
	g := GrpcServer{controller: controller}
	return g
}
