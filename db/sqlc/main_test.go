package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
	"time"
)

var testQueries *Queries
var testDB *pgxpool.Pool

const (
	dbDriver = "postgres"
	connStr  = "postgresql://root:sonzai@123456@localhost:5433/sonzai-bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	//conn, err := sql.Open(dbDriver, dbSource)
	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	testDB, err = pgxpool.NewWithConfig(ctx, conf)
	testQueries = New(testDB)
	os.Exit(m.Run())
}
