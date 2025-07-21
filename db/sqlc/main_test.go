package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonzai8/golang-sonzai-bank/utils"
	"log"
	"os"
	"testing"
	"time"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {

	config, err := utils.LoadConfig("../../.github/workflows/")
	if err != nil {
		log.Fatal(err)
	}
	pg := config.DbDriver
	dns := "postgresql://%s:%s@%s:%s/%s?sslmode=%s"
	var s = fmt.Sprintf(dns, pg.User, pg.Pass, pg.Host, pg.Port, pg.Name, pg.SSLMode)
	fmt.Println(s)
	//conn, err := sql.Open(dbDriver, dbSource)
	conf, err := pgxpool.ParseConfig(s)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	testDB, err = pgxpool.NewWithConfig(ctx, conf)
	testQueries = New(testDB)
	os.Exit(m.Run())
}
