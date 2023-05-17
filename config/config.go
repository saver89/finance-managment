package config

import (
	"os"
	"regexp"

	"github.com/spf13/viper"
)

type Config struct {
	AppVersion string
	Server     Server
	PostgresDB PostgresDB
	Logger     Logger
}

type Server struct {
	Port        string
	Address     string
	Development bool
}

type PostgresDB struct {
	Host     string
	Port     string
	User     string
	Password string
	SslMode  string
	Driver   string
	DbName   string
}

type Logger struct {
	Level string
}

const projectDirName = "finance-managment"

func LoadConfig() (config Config, err error) {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	viper.AddConfigPath(string(rootPath) + `/config`)
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
