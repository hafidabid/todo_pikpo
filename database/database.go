package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"time"
	"todo_pikpo/config"
	model "todo_pikpo/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Postgres *gorm.DB
	Redis    *redis.Client
}

func (db *Database) Migrate() error {
	err := db.Postgres.AutoMigrate(&model.TodoModel{})
	return err
}

func (db *Database) Flush() error {
	err := db.Postgres.Where("id is not null").Delete(&model.TodoModel{}).Error
	return err
}

func (db *Database) AddRedis(key string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error(time.Now().Format("2006-01-02 15:04:05"), " Database AddRedis ", err)
		return err
	}

	err = db.Redis.Set("pikpo-"+key, jsonData, time.Duration(1200)*time.Second).Err()

	log.Info(time.Now().Format("2006-01-02 15:04:05"), " Database AddRedis ", data)

	return err
}

func (db *Database) GetRedis(key string, res interface{}) error {
	val := db.Redis.Get("pikpo-" + key)
	if val == nil {
		log.Error(time.Now().Format("2006-01-02 15:04:05"), "Database GetRedis Key Not Found")
		return errors.New("key not found")
	}
	err := json.Unmarshal([]byte(val.Val()), res)

	log.Info(time.Now().Format("2006-01-02 15:04:05"), " Database GetRedis ", res)
	return err
}

func (db *Database) RedisRemove(addPrefix string) error {
	keys, err := db.Redis.Keys("pikpo-" + addPrefix + "*").Result()
	if err != nil {
		log.Error(time.Now().Format("2006-01-02 15:04:05"), " Database RedisRemove ", err)
		return err
	}

	if len(keys) > 0 {
		err = db.Redis.Del(keys...).Err()
		if err != nil {
			log.Error(time.Now().Format("2006-01-02 15:04:05"), " Database RedisRemove -> Key deletion ", err)
			return err
		}
	}

	log.Info(time.Now().Format("2006-01-02 15:04:05"), " Database RedisRemove ", keys)
	return nil
}

func NewDatabase(conf config.ConfigApp) (db Database, err error) {
	var newDatabase Database
	conn, err := gorm.Open(postgres.Open(
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			conf.DbUsername,
			conf.DbPassword,
			conf.DbHost,
			conf.DbPort,
			conf.DbName),
	), &gorm.Config{})
	if err != nil {
		return newDatabase, err
	}

	newDatabase.Postgres = conn

	newDatabase.Redis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort),
	})
	return newDatabase, nil
}
