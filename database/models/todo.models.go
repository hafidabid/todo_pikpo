package models

import "time"

type TodoModel struct {
	Id          string
	Author      string
	Title       string
	Description string
	IsDone      bool
	StartDate   time.Time
	EndDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
