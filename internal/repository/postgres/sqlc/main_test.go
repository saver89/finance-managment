package db

import (
	"database/sql"
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

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
