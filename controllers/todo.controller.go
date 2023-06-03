package controllers

import (
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"time"
	"todo_pikpo/database"
	model "todo_pikpo/database/models"
	"todo_pikpo/dto"
)

type TodoController struct {
	dto dto.TodoDTO
}

func (tc TodoController) verify(data *model.TodoModel) (int, error) {
	now := time.Now()
	if len(data.Author) < 3 {
		return 400, errors.New("author column should be filled with minimum 3 characters")
	}
	if len(data.Title) < 5 {
		return 400, errors.New("title column should be filled with minimum of 5 characters")
	}
	if data.EndDate.Unix() <= data.StartDate.Unix() {
		return 400, errors.New("EndDate should be greater than StartDate")
	}
	if data.EndDate.Unix() <= now.Unix() {
		return 400, errors.New("EndDate should be greater than now")
	}

	return 200, nil
}

func (tc TodoController) AddTodo(data model.TodoModel) (model.TodoModel, int, error) {
	if code, err := tc.verify(&data); err != nil {
		log.Error(time.Now(), " AddTodo controller ", err)

		return model.TodoModel{}, code, err
	}

	data.Id = uuid.New().String()
	data.IsDone = false
	data.UpdatedAt = time.Now()
	data.CreatedAt = time.Now()
	res, err := tc.dto.Create(data)
	if err != nil {
		log.Error(time.Now(), " AddTodo controller ", err)
		return model.TodoModel{}, 500, err
	}
	return res, 200, nil
}

func (tc TodoController) GetTodos(filter map[string]interface{}, page uint, limit uint) ([]model.TodoModel, int, error) {
	data, err := tc.dto.GetMany(filter, page, limit)
	if err != nil {
		log.Error(time.Now(), " GetTodos controller ", err)

		return []model.TodoModel{}, 500, err
	}
	return data, 200, nil
}

func (tc TodoController) GetTodo(id string) (model.TodoModel, int, error) {
	data, err := tc.dto.GetSingle(id)
	if err != nil {
		log.Error(time.Now(), " GetTodo controller ", err)

		return model.TodoModel{}, 404, err
	}
	return data, 200, nil
}

func (tc TodoController) EditTodo(id string, data model.TodoModel) (model.TodoModel, int, error) {
	if code, err := tc.verify(&data); err != nil {
		log.Error(time.Now(), " EditTodo controller ", err)

		return model.TodoModel{}, code, err
	}

	if _, err := tc.dto.GetSingle(id); err != nil {
		log.Error(time.Now(), " EditTodo controller ", err)

		return model.TodoModel{}, 404, err
	}

	data.UpdatedAt = time.Now()
	data.Id = id

	result, err := tc.dto.Update(id, data)
	if err != nil {
		log.Error(time.Now(), " EditTodo controller ", err)

		return model.TodoModel{}, 500, err
	}

	return result, 200, nil
}

func (tc TodoController) DeleteTodo(id string) (model.TodoModel, int, error) {
	if _, err := tc.dto.GetSingle(id); err != nil {
		log.Error(time.Now(), " DeleteTodo controller ", err)

		return model.TodoModel{}, 404, err
	}

	result, err := tc.dto.Delete(id)
	if err != nil {
		log.Error(time.Now(), " DeleteTodo controller ", err)

		return model.TodoModel{}, 500, err
	}

	return result, 200, nil
}

func CreateTodoController(db *database.Database) (TodoController, error) {
	var res TodoController
	res.dto = dto.TodoDTO{}
	res.dto.SetDb(db)
	return res, nil
}
