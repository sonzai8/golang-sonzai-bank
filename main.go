package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonzai8/golang-sonzai-bank/api"
	db "github.com/sonzai8/golang-sonzai-bank/db/sqlc"
	"github.com/sonzai8/golang-sonzai-bank/gapi"
	"github.com/sonzai8/golang-sonzai-bank/pb"
	"github.com/sonzai8/golang-sonzai-bank/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
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
	go runGatewayServer(config, store)
	runGrpcServer(config, store)

}

func runGrpcServer(config utils.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("can not create grpc server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSonZaiBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.AppConfig.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)

	}
	log.Printf("grpc server listening at %v", listener.Addr())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runGatewayServer(config utils.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("can not create grpc server:", err)
	}

	//grpcServer := grpc.NewServer()
	//pb.RegisterSonZaiBankServer(grpcServer, server)
	//reflection.Register(grpcServer)

	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSonZaiBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("can not register grpc gateway server:", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	addr := fmt.Sprintf(":%s", config.AppConfig.HttpPort)
	fmt.Printf("grpc gateway server listening at %v \n", addr)
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("failed to listen gateway: %v \n", err)

	}

	log.Printf("http gateway server listening at %v \n", listener.Addr())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalf("failed to start http gateway serve: %v", err)
	}
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}
	err = server.Start(config.AppConfig.HttpPort)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
