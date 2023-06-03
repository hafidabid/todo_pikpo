package database

import (
	"testing"
	"todo_pikpo/config"

	"github.com/stretchr/testify/suite"
)

type DbTestSuite struct {
	suite.Suite
	config config.ConfigApp
	db     Database
}

func (s *DbTestSuite) SetupSuite() {
	c, err := config.NewAppConfig("../.env.test")
	if err != nil {
		s.T().Error("Failed to create config:", err)
		return
	}
	s.config = c
	db, err := NewDatabase(CreateURI(c))
	if err != nil {
		s.T().Error("Failed to create database:", err)
		return
	}
	s.db = db
}

func (s *DbTestSuite) SetupTest() {
}

func (s *DbTestSuite) TearDownSuite() {
	s.db.Flush()
}

func (s *DbTestSuite) TearDownTest() {
	s.db.Flush()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

func (s *DbTestSuite) TestMigrate() {
	s.db.Migrate()

	err := s.db.Postgres.Exec("SELECT * FROM todo_models").Error
	s.Equal(err, nil)
}

func (s *DbTestSuite) TestConnect() {
	err := s.db.Postgres.Exec("SELECT 1").Error
	s.Equal(err, nil)
}
