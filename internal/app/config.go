package app

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
	Type     string
}

func (dbs *DBConfig) URL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		dbs.Type, dbs.User, dbs.Password, dbs.Host, dbs.Port, dbs.Name)
}

func LoadConfig(path string) Config {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		log.Print(err)
	}

	var config Config

	if str, ok := os.LookupEnv("DB_HOST"); ok {
		config.DB.Host = str
	} else {
		config.DB.Host = viper.GetString("DB.HOST")
	}

	if str, ok := os.LookupEnv("DB_NAME"); ok {
		config.DB.Name = str
	} else {
		config.DB.Name = viper.GetString("DB.NAME")
	}

	if str, ok := os.LookupEnv("DB_PASSWORD"); ok {
		config.DB.Password = str
	} else {
		config.DB.Password = viper.GetString("DB.PASSWORD")
	}

	if str, ok := os.LookupEnv("DB_PORT"); ok {
		config.DB.Port, err = strconv.Atoi(str)
		if err != nil {
			log.Print(err)
		}
	} else {
		config.DB.Port = viper.GetInt("DB.PORT")
	}

	if str, ok := os.LookupEnv("DB_TYPE"); ok {
		config.DB.Type = str
	} else {
		config.DB.Type = viper.GetString("DB.type")
	}

	if str, ok := os.LookupEnv("DB_USER"); ok {
		config.DB.User = str
	} else {
		config.DB.User = viper.GetString("DB.USER")
	}

	return config
}
