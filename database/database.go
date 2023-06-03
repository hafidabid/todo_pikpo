package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"todo_pikpo/config"
	model "todo_pikpo/database/models"
)

type Database struct {
	Postgres *gorm.DB
}

func (db *Database) Migrate() error {
	err := db.Postgres.AutoMigrate(&model.TodoModel{})
	return err
}

func (db *Database) Flush() error {
	err := db.Postgres.Delete(&model.TodoModel{}).Error
	return err
}

func NewDatabase(connURI string) (db Database, err error) {
	var newDatabase Database

	conn, err := gorm.Open(postgres.Open(connURI), &gorm.Config{})
	if err != nil {
		return newDatabase, err
	}

	newDatabase.Postgres = conn

	return newDatabase, nil
}

func CreateURI(config config.ConfigApp) string {
	return ""
}
