package model

import (
	"reflect"
	"time"
	"todo_pikpo/database"
)
import "todo_pikpo/interface"

type TodoModel struct {
	Id          string    `json:"id" gorm:"primary_key"`
	Author      string    `json:"author" gorm:"not_null"`
	Title       string    `json:"title" gorm:"not_null"`
	Description string    `json:"description" gorm:"type:text"`
	IsDone      bool      `json:"isDone" gorm:"default:false"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TodoDTO struct {
	_interface.DtoInterface[TodoModel]
	db *database.Database
}

func (td *TodoDTO) SetDb(db *database.Database) {
	td.db = db
}

func (td *TodoDTO) GetMany(filter map[string]interface{}, page uint, pageSize uint) ([]TodoModel, error) {
	var data []TodoModel
	err := td.db.Postgres.Limit(int(pageSize)).Offset(int(page*pageSize)).Find(&data, filter).Error
	if err != nil {
		return []TodoModel{}, err
	}

	return data, nil
}

func (td *TodoDTO) GetSingle(id string) (TodoModel, error) {
	var data TodoModel
	err := td.db.Postgres.First(&data, "id = ?", id).Error
	if err != nil {
		return TodoModel{}, err
	}
	return data, nil
}

func (td *TodoDTO) Create(data TodoModel) (TodoModel, error) {
	err := td.db.Postgres.Create(&data).Error
	if err != nil {
		return TodoModel{}, err
	}
	return data, nil
}

func (td *TodoDTO) Update(id string, data TodoModel) (TodoModel, error) {
	var ret TodoModel
	err := td.db.Postgres.First(&ret, "id = ?", id).Error
	if err != nil {
		return TodoModel{}, err
	}

	// use reflection to update important data / selected data only
	v := reflect.ValueOf(data)
	z := reflect.ValueOf(&ret)
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i)
		value := v.Field(i)

		z.FieldByName(fieldName.Name).Set(value)
	}

	ret.UpdatedAt = time.Now()
	err = td.db.Postgres.Save(&ret).Error

	return ret, err
}

func (td *TodoDTO) Delete(id string) (TodoModel, error) {
	var data TodoModel
	err := td.db.Postgres.Where("id = ?", id).Delete(&data).Error
	if err != nil {
		return TodoModel{}, err
	}

	return data, nil
}
