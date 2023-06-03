package model

import (
	"time"
)

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
