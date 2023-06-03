package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ConfigApp struct {
	DbUsername string `mapstructure:"PG_USERNAME"`
	DbPassword string `mapstructure:"PG_PASS"`
	DbName     string `mapstructure:"PG_DB"`
	DbPort     uint16 `mapstructure:"PG_PORT"`
	DbHost     string `mapstructure:"PG_HOST"`
	RedisHost  string `mapstructure:"REDIS_HOST"`
	RedisPort  string `mapstructure:"REDIS_PORT"`
	RedisUsn   string `mapstructure:"REDIS_USN"`
	RedisPass  string `mapstructure:"REDIS_PASS"`
	EncryptKey string `mapstructure:"KEY"`
	Port       uint16 `mapstructure:"PORT"`
}

func NewAppConfig(filePath string) (c ConfigApp, e error) {
	var confResult ConfigApp
	viper.SetConfigFile(filePath)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if e := viper.ReadInConfig(); e != nil {
		log.Error("error in creating NewAppConfig with error ", e)
	}
	e = viper.Unmarshal(&confResult)
	return confResult, e
}
