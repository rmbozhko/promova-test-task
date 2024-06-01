package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"promova-test-task/util"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	testDatabase := SetupTestDatabase()
	config, err := util.LoadConfig("../..")
	testDB, err = sql.Open(config.DBDriver, testDatabase.DbSource)
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	testQueries = New(testDB)

	defer testDatabase.TearDown()
	os.Exit(m.Run())
}
