package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/saver89/finance-management/config"
	db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"
	"github.com/saver89/finance-management/internal/server"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := server.NewServer(store)

	err = server.Run(config.ServerAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
}
