package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"
	"github.com/saver89/finance-management/internal/server"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/finances?sslmode=disable"
	serverAddress = "0.0.0.0:8099"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := server.NewServer(store)

	err = server.Run(serverAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
}
