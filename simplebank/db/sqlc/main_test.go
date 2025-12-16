package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbSource := config.DBSource

	testDB, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
