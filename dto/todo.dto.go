package dto

import (
	"time"
	"todo_pikpo/database"
	model "todo_pikpo/database/models"
	_interface "todo_pikpo/interface"

	log "github.com/sirupsen/logrus"
)

type TodoDTO struct {
	_interface.DtoInterface[model.TodoModel]
	Db *database.Database
}

func (td *TodoDTO) SetDb(db *database.Database) {
	td.Db = db
}

func (td *TodoDTO) GetMany(filter map[string]interface{}, page uint, pageSize uint) ([]model.TodoModel, error) {
	var data []model.TodoModel

	var err error
	if len(filter) >= 1 {
		err = td.Db.Postgres.Limit(int(pageSize)).Offset(int(page*pageSize)).Find(&data, filter).Error
	} else {
		err = td.Db.Postgres.Limit(int(pageSize)).Offset(int(page * pageSize)).Find(&data).Error
	}

	if err != nil {
		log.Error(err)
		return []model.TodoModel{}, err
	}

	return data, nil
}

func (td *TodoDTO) GetSingle(id string) (model.TodoModel, error) {
	var data model.TodoModel
	err := td.Db.Postgres.First(&data, "id = ?", id).Error
	if err != nil {
		return model.TodoModel{}, err
	}
	return data, nil
}

func (td *TodoDTO) Create(data model.TodoModel) (model.TodoModel, error) {
	err := td.Db.Postgres.Create(&data).Error
	if err != nil {
		return model.TodoModel{}, err
	}
	return data, nil
}

func (td *TodoDTO) Update(id string, data model.TodoModel) (model.TodoModel, error) {
	var ret model.TodoModel
	err := td.Db.Postgres.First(&ret, "id = ?", id).Error
	if err != nil {
		return model.TodoModel{}, err
	}

	ret.IsDone = data.IsDone
	ret.Author = data.Author
	ret.Description = data.Description
	ret.Title = data.Title
	ret.StartDate = data.StartDate
	ret.EndDate = data.EndDate

	ret.UpdatedAt = time.Now()
	err = td.Db.Postgres.Save(&ret).Error

	return ret, err
}

func (td *TodoDTO) Delete(id string) (model.TodoModel, error) {
	var data model.TodoModel
	err := td.Db.Postgres.Where("id = ?", id).Delete(&data).Error
	if err != nil {
		return model.TodoModel{}, err
	}

	return data, nil
}
