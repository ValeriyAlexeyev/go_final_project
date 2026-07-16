package main

import (
	"log"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
