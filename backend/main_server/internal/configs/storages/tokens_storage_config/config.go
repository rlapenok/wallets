package tokens_storage_config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type MongoConfigInterface interface {
	GetUri() string
	GetTimeOut() uint8
	GetPoolSize() uint8
	GetCollName() string
}

type MongoConfig struct {
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	DbName   string `yaml:"db_name"`
	CollName string `yaml:"coll_name"`
	Port     int    `yaml:"port"`
	Pool     int    `yaml:"pool"`
	Time     int    `yaml:"time"`
}

func New() MongoConfig {
	env := os.Getenv("mongo_config")
	if env == "" {
		log.Fatal("Empty env")
		//add normal loggining
	}
	var config MongoConfig
	err := cleanenv.ReadConfig(env, &config)
	if err != nil {
		log.Fatal(err)
		//add normal logginig
	}
	return config

}

func (cfg *MongoConfig) GetUri() string {
	uri := fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	return uri
}
