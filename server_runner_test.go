package main

import (
	"testing"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
	"github.com/ValeriyAlexeyev/go_final_project/pkg/server"
)

// TestServerRunner временно используется для запуска сервера,
// когда выполнение локального scheduler.exe блокирует Windows.
func TestServerRunner(t *testing.T) {
	if err := db.Init("scheduler.db"); err != nil {
		t.Fatalf("database initialization error: %v", err)
	}
	defer db.Close()

	if err := server.Run(); err != nil {
		t.Fatalf("server error: %v", err)
	}
}
