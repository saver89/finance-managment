package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/saver89/finance-management/config"
)

var (
	roundPrecision float64 = 1000000
	testDB         *sql.DB
	testQueries    *Queries
)

func TestMain(m *testing.M) {
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
	testDB, err = sql.Open(config.PostgresDB.Driver, connString)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
