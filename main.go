package main

import (
	"database/sql"
	"net"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/ngtrdai197/simple-bank/cmd/api"
	"github.com/ngtrdai197/simple-bank/cmd/gapi"
	db "github.com/ngtrdai197/simple-bank/db/sqlc"
	"github.com/ngtrdai197/simple-bank/pb"
	"github.com/ngtrdai197/simple-bank/util"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect db")
	}
	runMigration(config)

	store := db.NewStore(conn)

	// go runGrpcServer(config, store)
	runGinServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create gRPC server")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create listener")
	}

	log.Printf("Start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start gRPC server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server")
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server")
	}
}

func runMigration(config util.Config) {
	migration, err := migrate.New(
		config.MigrationUrl,
		config.DBSource)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new migrate instance")
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("Failed to run migrate up")
	}

	log.Info().Msg("DB migrated successfully")
}
