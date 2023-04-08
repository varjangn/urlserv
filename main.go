package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/varjangn/urlserv/api"
	"github.com/varjangn/urlserv/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbPath := os.Getenv("SQLITE_DB_PATH")
	store, err := storage.NewSqliteStorage(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if err = store.Init(); err != nil {
		log.Fatal(err)
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	server := api.NewAPIServer(listenAddr, store)
	if err = server.Run(); err != nil {
		log.Fatal(err)
	}

}
