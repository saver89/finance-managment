package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	"github.com/saver89/finance-management/config"
	db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"
	"github.com/saver89/finance-management/internal/server"
	"github.com/saver89/finance-management/pkg/logger"
)

func main() {
	log.Println("Starting Finance Management APP...")
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.PostgresDB.Host,
		config.PostgresDB.Port,
		config.PostgresDB.User,
		config.PostgresDB.Password,
		config.PostgresDB.DbName,
		config.PostgresDB.SslMode)
	conn, err := sql.Open(config.PostgresDB.Driver, connString)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	appLogger := logger.NewApiLogger(&config)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, DevelopmentMode: %s",
		config.AppVersion,
		config.Logger.Level,
		config.Server.Development,
	)
	appLogger.Infof("Success parsed config: %#v", config.AppVersion)

	store := db.NewStore(conn)
	server := server.NewServer(store, appLogger)

	addressUri := fmt.Sprintf("%s:%s", config.Server.Address, config.Server.Port)
	err = server.Run(addressUri)
	if err != nil {
		log.Fatal(err.Error())
	}
}
