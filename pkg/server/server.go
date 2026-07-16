package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/api"
)

const (
	defaultPort = "7540"
	webDir      = "./web"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()

	// Сначала регистрируем API.
	api.Init(mux)

	// Обработчик "/" должен регистрироваться после API.
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	address := ":" + port

	fmt.Printf("Scheduler server started: http://localhost:%s\n", port)

	return http.ListenAndServe(address, mux)
}
