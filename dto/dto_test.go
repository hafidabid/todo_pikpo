package dto

import (
	"fmt"
	"testing"
	"time"
	"todo_pikpo/config"
	"todo_pikpo/database"
	model "todo_pikpo/database/models"

	"github.com/stretchr/testify/suite"
)

type DtoTestSuite struct {
	suite.Suite
	dto TodoDTO
}

func (s *DtoTestSuite) SetupSuite() {
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
	s.dto.db = &db
}

func (s *DtoTestSuite) SetupTest() {
	s.dto.db.Migrate()
}

func (s *DtoTestSuite) TearDownSuite() {

}

func (s *DtoTestSuite) TearDownTest() {
	s.dto.db.Flush()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(DtoTestSuite))
}

func (s *DtoTestSuite) TestAdd() {
	s.dto.Create(model.TodoModel{
		Id:          "1",
		Author:      "-",
		Title:       "test",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	var res model.TodoModel
	err := s.dto.db.Postgres.Where("id = ?", "1").First(&res).Error

	a := s.Suite.Assert()

	a.Equal(err, nil)
	a.Equal(res.Id, "1")
	a.Equal(res.Author, "-")

	err = s.dto.db.Postgres.Where("id = ?", "10").First(&res).Error
	a.NotEqual(err, nil)

	// Test for duplicate ID
	_, err = s.dto.Create(model.TodoModel{
		Id:          "1",
		Author:      "-",
		Title:       "test",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	a.NotEqual(err, nil)
}

func (s *DtoTestSuite) TestGetMany() {
	s.dto.Create(model.TodoModel{
		Id:          "1",
		Author:      "-",
		Title:       "test",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	s.dto.Create(model.TodoModel{
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

	data, err := s.dto.GetMany(map[string]interface{}{}, 0, 10)
	a := s.Suite.Assert()

	a.Equal(err, nil)
	a.Equal(len(data), 2)

	s.dto.Create(model.TodoModel{
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

	data, err = s.dto.GetMany(map[string]interface{}{}, 0, 10)
	a.Equal(err, nil)
	a.Equal(len(data), 3)
}

func (s *DtoTestSuite) TestGetManyPagination() {
	a := s.Suite.Assert()
	s.dto.Create(model.TodoModel{
		Id:          "1",
		Author:      "-",
		Title:       "test",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	s.dto.Create(model.TodoModel{
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

	s.dto.Create(model.TodoModel{
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

	// Test for pagination
	data, err := s.dto.GetMany(map[string]interface{}{}, 0, 1)
	fmt.Println(len(data))
	a.Equal(err, nil)
	a.Equal(len(data), 1)
	a.Equal(data[0].Id, "1")

	data, err = s.dto.GetMany(map[string]interface{}{}, 1, 1)
	a.Equal(err, nil)
	a.Equal(len(data), 1)
	a.Equal(data[0].Id, "2")

	data, err = s.dto.GetMany(map[string]interface{}{}, 3, 2)
	a.Equal(err, nil)
	a.Equal(len(data), 0)
}

func (s *DtoTestSuite) TestGetManyQuery() {
	a := s.Suite.Assert()
	s.dto.Create(model.TodoModel{
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

	s.dto.Create(model.TodoModel{
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

	s.dto.Create(model.TodoModel{
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

	// Test for filter
	data, err := s.dto.GetMany(map[string]interface{}{
		"title": "test",
	}, 0, 10)
	a.Equal(err, nil)
	a.Equal(len(data), 2)
	a.Equal(data[0].Id, "2")

	data, err = s.dto.GetMany(map[string]interface{}{
		"author": "James",
	}, 0, 10)
	a.Equal(err, nil)
	a.Equal(len(data), 1)
	a.Equal(data[0].Id, "3")

	data, err = s.dto.GetMany(map[string]interface{}{
		"title": "will not found there",
	}, 0, 10)
	a.Equal(err, nil)
	a.Equal(len(data), 0)
}

func (s *DtoTestSuite) TestGetOne() {
	a := s.Suite.Assert()
	s.dto.Create(model.TodoModel{
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

	data, err := s.dto.GetSingle("1")
	a.Equal(err, nil)
	a.Equal(data.Id, "1")
	a.Equal(data.Title, "test for make sure")

	data, err = s.dto.GetSingle("2")
	a.NotEqual(err, nil)
}

func (s *DtoTestSuite) TestUpdate() {
	a := s.Suite.Assert()
	s.dto.Create(model.TodoModel{
		Id:          "1",
		Author:      "-",
		Title:       "test for make sure",
		Description: "-",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
	})

	data, err := s.dto.Update("1", model.TodoModel{
		Id:          "1",
		Author:      "James",
		Title:       "changed",
		Description: "this is changed too",
		IsDone:      false,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
	})
	a.Equal(err, nil)
	a.Equal(data.Id, "1")
	a.Equal(data.Title, "changed")
	a.Equal(data.Description, "this is changed too")
	a.Equal(data.Author, "James")

	_, err = s.dto.GetSingle("2")
	a.NotEqual(err, nil)
}

func (s *DtoTestSuite) TestDelete() {
	a := s.Suite.Assert()
	s.dto.Create(model.TodoModel{
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

	s.dto.Create(model.TodoModel{
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

	var data []model.TodoModel
	err := s.dto.db.Postgres.Find(&data).Error
	a.Equal(err, nil)
	a.Equal(len(data), 2)

	_, err = s.dto.Delete("1")
	a.Equal(err, nil)

	err = s.dto.db.Postgres.Find(&data).Error
	a.Equal(err, nil)
	a.Equal(len(data), 1)

	_, err = s.dto.Delete("2")
	a.Equal(err, nil)

	err = s.dto.db.Postgres.Find(&data).Error
	a.Equal(err, nil)
	a.Equal(len(data), 0)
}
