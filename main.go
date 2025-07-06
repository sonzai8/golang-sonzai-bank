package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonzai8/golang-sonzai-bank/api"
	db "github.com/sonzai8/golang-sonzai-bank/db/sqlc"
	"github.com/sonzai8/golang-sonzai-bank/utils"
	"log"
	"time"
)

// SimpleSQLTracer is a basic implementation of the pgx.Tracer interface.
type SimpleSQLTracer struct{}

// TraceQueryStart is called before a query is executed.
func (t *SimpleSQLTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	log.Printf("Executing query: %s with args: %v", data.SQL, data.Args)
	return ctx
}

// TraceQueryEnd is called after a query has finished.
func (t *SimpleSQLTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		log.Printf("Query failed: %v", data.Err)
	} else {
		log.Printf("Query successful, time taken: %s", data.CommandTag)
	}
}

var testQueries *db.Queries
var pgPool *pgxpool.Pool

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	pg := config.DbDriver
	dns := "postgresql://%s:%s@%s:%s/%s?sslmode=%s"
	var s = fmt.Sprintf(dns, pg.User, pg.Pass, pg.Host, pg.Port, pg.Name, pg.SSLMode)
	conf, err := pgxpool.ParseConfig(s)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conf.ConnConfig.Tracer = &SimpleSQLTracer{}

	pgPool, err = pgxpool.NewWithConfig(ctx, conf)

	store := db.NewStore(pgPool)
	server := api.NewServer(store)
	err = server.Start(config.AppConfig.Port)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
