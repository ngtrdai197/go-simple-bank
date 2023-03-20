package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/ngtrdai197/simple-bank/api"
	db "github.com/ngtrdai197/simple-bank/db/sqlc"
	"github.com/ngtrdai197/simple-bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect db", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
