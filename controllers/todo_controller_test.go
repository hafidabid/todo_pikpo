package controllers

import (
	"testing"
	"time"
	"todo_pikpo/config"
	"todo_pikpo/database"
	model "todo_pikpo/database/models"

	"github.com/stretchr/testify/suite"
)

type ControllerTest struct {
	suite.Suite
	controller TodoController
	db         *database.Database
}

func (s *ControllerTest) SetupTest() {
	s.db.Migrate()
}

func (s *ControllerTest) TearDownTest() {
	s.db.Flush()
}

func (s *ControllerTest) SetupSuite() {
	c, err := config.NewAppConfig("../.env.test")
	if err != nil {
		s.T().Error("Failed to create config:", err)
		return
	}

	db, err := database.NewDatabase(database.CreateURI(c))
	if err != nil {
		s.T().Error("Failed to create database:", err)
		return
	}

	nc, err := CreateTodoController(&db)
	if err != nil {
		s.T().Error("Failed to create controller:", err)
		return
	}
	s.controller = nc
	s.db = &db
}

func (s *ControllerTest) TearDownSuite() {
	s.db.Flush()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTest))
}

func (s *ControllerTest) TestAdd() {
	a := s.Suite.Assert()

	_, code, err := s.controller.AddTodo(model.TodoModel{
		Id:          "1",
		Author:      "-",
		Title:       "test",
		Description: "-",
		IsDone:      true,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	a.Equal(code, 400)
	a.NotEqual(err, nil)

	res, code, err := s.controller.AddTodo(model.TodoModel{
		Id:          "1",
		Author:      "james",
		Title:       "test this is title",
		Description: "lorem ipsom dolom amet",
		IsDone:      true,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(1 * time.Hour),
		CreatedAt:   time.Now().Add(-300 * time.Hour),
		UpdatedAt:   time.Now().Add(-400 * time.Hour),
	})
	a.Equal(code, 200)
	a.Equal(err, nil)
	a.NotEqual(res.Id, "1")
	a.Equal(res.Author, "james")
	a.Equal(res.IsDone, false)
	a.Equal(res.CreatedAt.Unix(), res.UpdatedAt.Unix())
}

func (s *ControllerTest) TestDelete() {
	a := s.Suite.Assert()

	res, code, err := s.controller.AddTodo(model.TodoModel{
		Id:          "1",
		Author:      "james",
		Title:       "test this is title",
		Description: "lorem ipsom dolom amet",
		IsDone:      true,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(72 * time.Hour),
		CreatedAt:   time.Now().Add(-300 * time.Hour),
		UpdatedAt:   time.Now().Add(-400 * time.Hour),
	})
	a.Equal(code, 200)
	a.Equal(err, nil)
	a.NotEqual(res.Id, "1")

	_, code, err = s.controller.DeleteTodo(res.Id)
	a.Equal(code, 200)

	_, code, err = s.controller.GetTodo(res.Id)
	a.Equal(code, 404)
}

func (s *ControllerTest) TestUpdate() {
	a := s.Suite.Assert()

	var tempData = model.TodoModel{
		Id:          "1",
		Author:      "james",
		Title:       "test this is title",
		Description: "lorem ipsom dolom amet",
		IsDone:      true,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(72 * time.Hour),
	}

	res, code, err := s.controller.AddTodo(tempData)
	a.Equal(code, 200)
	a.Equal(err, nil)
	a.NotEqual(res.Id, "1")

	_, code, err = s.controller.EditTodo(res.Id, model.TodoModel{
		Author:      "james",
		Title:       "test this is title",
		Description: "lorem ipsom dolom amet",
		IsDone:      true,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(-1 * time.Hour),
	})
	a.Equal(code, 400)
	a.NotEqual(err, nil)

	time.Sleep(2 * time.Second)

	res2, code, err := s.controller.EditTodo(res.Id, model.TodoModel{
		Author:      "james",
		Title:       "test this is title",
		Description: "lorem ipsom dolom amet",
		IsDone:      true,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(10 * time.Hour),
	})
	a.Equal(code, 200)
	a.Equal(err, nil)
	a.Greater(res2.UpdatedAt.Unix(), res2.CreatedAt.Unix())
	a.Greater(res2.UpdatedAt.Unix(), res.UpdatedAt.Unix())
	a.Equal(res2.IsDone, true)
	a.Equal(res2.Author, "james")
}

func (s *ControllerTest) TestList() {
	a := s.Suite.Assert()

	s.controller.dto.Create(model.TodoModel{
		Id:          "1",
		Author:      "-",
		Title:       "test for make sure",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	s.controller.dto.Create(model.TodoModel{
		Id:          "2",
		Author:      "-",
		Title:       "test",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	s.controller.dto.Create(model.TodoModel{
		Id:          "3",
		Author:      "James",
		Title:       "test",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	data, code, err := s.controller.GetTodos(map[string]interface{}{}, 0, 10)
	a.Equal(code, 200)
	a.Equal(err, nil)
	a.Equal(len(data), 3)
}

func (s *ControllerTest) TestListPagination() {
	a := s.Suite.Assert()

	s.controller.dto.Create(model.TodoModel{
		Id:          "1",
		Author:      "-",
		Title:       "test for make sure",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	s.controller.dto.Create(model.TodoModel{
		Id:          "2",
		Author:      "-",
		Title:       "test",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	s.controller.dto.Create(model.TodoModel{
		Id:          "3",
		Author:      "James",
		Title:       "test",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	data, code, err := s.controller.GetTodos(map[string]interface{}{}, 1, 1)
	a.Equal(code, 200)
	a.Equal(err, nil)
	a.Equal(len(data), 1)
	a.Equal(data[0].Id, "2")

	data, code, err = s.controller.GetTodos(map[string]interface{}{}, 2, 1)
	a.Equal(code, 200)
	a.Equal(err, nil)
	a.Equal(len(data), 1)
	a.Equal(data[0].Id, "3")

	data, code, err = s.controller.GetTodos(map[string]interface{}{}, 2, 10)
	a.Equal(code, 200)
	a.Equal(err, nil)
	a.Equal(len(data), 0)
}

func (s *ControllerTest) TestListQuery() {
	a := s.Suite.Assert()

	s.controller.dto.Create(model.TodoModel{
		Id:          "1",
		Author:      "-",
		Title:       "test for make sure",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	s.controller.dto.Create(model.TodoModel{
		Id:          "2",
		Author:      "-",
		Title:       "test",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	s.controller.dto.Create(model.TodoModel{
		Id:          "3",
		Author:      "James",
		Title:       "test",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	data, code, err := s.controller.GetTodos(map[string]interface{}{
		"author": "James",
	}, 0, 10)
	a.Equal(code, 200)
	a.Equal(err, nil)
	a.Equal(len(data), 1)

	data, code, err = s.controller.GetTodos(map[string]interface{}{
		"title": "James",
	}, 0, 10)
	a.Equal(code, 200)
	a.Equal(err, nil)
	a.Equal(len(data), 0)
}

func (s *ControllerTest) TestGet() {
	a := s.Suite.Assert()

	_, code, err := s.controller.GetTodo("test")
	a.Equal(code, 404)
	a.NotEqual(err, nil)

	s.controller.dto.Create(model.TodoModel{
		Id:          "1",
		Author:      "-",
		Title:       "test",
		Description: "-",
		IsDone:      false,
	})

	data, code, err := s.controller.GetTodo("1")
	a.Equal(code, 200)
	a.Equal(err, nil)
	a.Equal(data.Id, "1")
}
