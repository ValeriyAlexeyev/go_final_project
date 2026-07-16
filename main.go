package main

import (
	"log"
	"os"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
	"github.com/ValeriyAlexeyev/go_final_project/pkg/server"
)

const defaultDBFile = "scheduler.db"

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("database initialization error: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("database close error: %v", err)
		}
	}()

	if err := server.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
