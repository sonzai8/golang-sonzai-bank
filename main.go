package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sonzai8/golang-sonzai-bank/api"
	db "github.com/sonzai8/golang-sonzai-bank/db/sqlc"
	"github.com/sonzai8/golang-sonzai-bank/gapi"
	"github.com/sonzai8/golang-sonzai-bank/pb"
	"github.com/sonzai8/golang-sonzai-bank/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
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
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg(fmt.Sprintf("environment: %s", config.AppConfig.Environment))

	if config.AppConfig.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	pg := config.DbDriver
	dns := "postgresql://%s:%s@%s:%s/%s?sslmode=%s"
	var s = fmt.Sprintf(dns, pg.User, pg.Pass, pg.Host, pg.Port, pg.Name, pg.SSLMode)
	conf, err := pgxpool.ParseConfig(s)
	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conf.ConnConfig.Tracer = &SimpleSQLTracer{}

	pgPool, err = pgxpool.NewWithConfig(ctx, conf)

	// run db migration
	runDBMigration(config.AppConfig.MigrationURL, s)
	store := db.NewStore(pgPool)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)

}

func runDBMigration(migrateURL string, dbSource string) {
	log.Print("Running migrations...", dbSource)
	mgr, err := migrate.New(migrateURL, dbSource)
	if err != nil {
		//log.Fatal().Msg(migrateURL)
		//log.Fatal().Msg(dbSource)
		log.Fatal().Msg("cannot create new migrate instance")
	}
	err = mgr.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Msg("failed to run migrate up: ")
	}
}

func runGrpcServer(config utils.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("can not create grpc server:")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSonZaiBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.AppConfig.GrpcPort)
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)

	}
	log.Printf("grpc server listening at %v", listener.Addr())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}

func runGatewayServer(config utils.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("can not create grpc server:")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	//grpcServer := grpc.NewServer()
	//pb.RegisterSonZaiBankServer(grpcServer, server)
	//reflection.Register(grpcServer)

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSonZaiBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msg("can not register grpc gateway server:")
	}
	mux := http.NewServeMux()

	mux.Handle("/", grpcMux)
	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	addr := fmt.Sprintf(":%s", config.AppConfig.HttpPort)
	fmt.Printf("grpc gateway server listening at %v \n", addr)
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal().Msgf("failed to listen gateway: %v \n")

	}

	log.Printf("http gateway server listening at %v \n", listener.Addr())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msgf("failed to start http gateway serve: %v", err)
	}
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	err = server.Start(config.AppConfig.HttpPort)
	if err != nil {
		log.Fatal().Msg("cannot start server:")
	}
}
