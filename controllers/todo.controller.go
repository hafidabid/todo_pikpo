package controllers

import (
	"errors"
	"time"
	"todo_pikpo/database"
	model "todo_pikpo/database/models"
)

type TodoController struct {
	dto model.TodoDTO
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
		return model.TodoModel{}, code, err
	}

	data.IsDone = false
	data.UpdatedAt = time.Now()
	data.CreatedAt = time.Now()
	res, err := tc.dto.Create(data)
	if err != nil {
		return model.TodoModel{}, 500, err
	}
	return res, 200, nil
}

func (tc TodoController) GetTodos(filter map[string]interface{}, page uint, limit uint) ([]model.TodoModel, int, error) {
	data, err := tc.dto.GetMany(filter, page, limit)
	if err != nil {
		return []model.TodoModel{}, 500, err
	}
	return data, 200, nil
}

func (tc TodoController) GetTodo(id string) (model.TodoModel, int, error) {
	data, err := tc.dto.GetSingle(id)
	if err != nil {
		return model.TodoModel{}, 404, err
	}
	return data, 200, nil
}

func (tc TodoController) EditTodo(id string, data model.TodoModel) (model.TodoModel, int, error) {
	if code, err := tc.verify(&data); err != nil {
		return model.TodoModel{}, code, err
	}

	if _, err := tc.dto.GetSingle(id); err != nil {
		return model.TodoModel{}, 404, err
	}

	data.UpdatedAt = time.Now()
	data.Id = id

	result, err := tc.dto.Update(id, data)
	if err != nil {
		return model.TodoModel{}, 500, err
	}

	return result, 200, nil
}

func (tc TodoController) DeleteTodo(id string) (model.TodoModel, int, error) {
	if _, err := tc.dto.GetSingle(id); err != nil {
		return model.TodoModel{}, 404, err
	}

	result, err := tc.dto.Delete(id)
	if err != nil {
		return model.TodoModel{}, 500, err
	}

	return result, 200, nil
}

func CreateTodoController(db *database.Database) (TodoController, error) {
	var res TodoController

	return res, nil
}
