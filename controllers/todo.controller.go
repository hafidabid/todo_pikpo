package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
	"todo_pikpo/database"
	model "todo_pikpo/database/models"
	"todo_pikpo/dto"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
		log.Error(time.Now().Format("2006-01-02 15:04:05"), " AddTodo controller ", err)

		return model.TodoModel{}, code, err
	}

	res, err := tc.dto.Create(model.TodoModel{
		Id:          uuid.New().String(),
		Author:      data.Author,
		Title:       data.Title,
		Description: data.Description,
		IsDone:      false,
		StartDate:   data.StartDate,
		EndDate:     data.EndDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		log.Error(time.Now().Format("2006-01-02 15:04:05"), " AddTodo controller ", err)
		return model.TodoModel{}, 500, err
	}

	//Revoke data from redis too
	_ = tc.dto.Db.RedisRemove("list-")

	return res, 200, nil
}

func (tc TodoController) GetTodos(filter map[string]interface{}, page uint, limit uint) ([]model.TodoModel, int, error) {
	// Get data from redis first
	var data []model.TodoModel
	jd, e0 := json.Marshal(filter)
	if e0 == nil {
		md5hash := md5.Sum(jd)
		hashed := hex.EncodeToString(md5hash[:])
		eRedis := tc.dto.Db.GetRedis("list-"+hashed, &data)
		if eRedis == nil {
			return data, 200, nil
		}
	}

	// Get data from postgres
	data, err := tc.dto.GetMany(filter, page, limit)
	if err != nil {
		log.Error(time.Now().Format("2006-01-02 15:04:05"), " GetTodos controller ", err)

		return []model.TodoModel{}, 500, err
	}

	// Insert data into redis
	if e0 == nil {
		md5hash := md5.Sum(jd)
		hashed := hex.EncodeToString(md5hash[:])
		_ = tc.dto.Db.AddRedis("list-"+hashed, data)
	}

	return data, 200, nil

}

func (tc TodoController) GetTodo(id string) (model.TodoModel, int, error) {

	// Get data from redis first
	var data model.TodoModel
	err := tc.dto.Db.GetRedis(id, &data)
	if err == nil {
		return data, 200, nil
	}

	// Get data from postgres
	data, err = tc.dto.GetSingle(id)

	if err != nil {
		log.Error(time.Now().Format("2006-01-02 15:04:05"), " GetTodo controller ", err)

		return model.TodoModel{}, 404, err
	}

	// Insert data into redis
	_ = tc.dto.Db.AddRedis(id, data)

	return data, 200, nil
}

func (tc TodoController) EditTodo(id string, data model.TodoModel) (model.TodoModel, int, error) {
	if code, err := tc.verify(&data); err != nil {
		log.Error(time.Now().Format("2006-01-02 15:04:05"), " EditTodo controller ", err)

		return model.TodoModel{}, code, err
	}

	if _, err := tc.dto.GetSingle(id); err != nil {
		log.Error(time.Now().Format("2006-01-02 15:04:05"), " EditTodo controller ", err)

		return model.TodoModel{}, 404, err
	}

	data.UpdatedAt = time.Now()
	data.Id = id

	result, err := tc.dto.Update(id, data)
	if err != nil {
		log.Error(time.Now(), " EditTodo controller ", err)

		return model.TodoModel{}, 500, err
	}

	//Revoke data from redis too
	_ = tc.dto.Db.RedisRemove(id)
	_ = tc.dto.Db.RedisRemove("list-")

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

	//Revoke data from redis too
	_ = tc.dto.Db.RedisRemove(id)
	_ = tc.dto.Db.RedisRemove("list-")

	return result, 200, nil
}

func CreateTodoController(db *database.Database) (TodoController, error) {
	var res TodoController
	res.dto = dto.TodoDTO{}
	res.dto.SetDb(db)
	return res, nil
}
