package server

import (
	"fmt"
	"net/http"
	"os"
)

const (
	defaultPort = "7540"
	webDir      = "./web"
)

// Run настраивает маршруты и запускает HTTP-сервер.
func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()

	// Раздача статических файлов из каталога web.
	fileServer := http.FileServer(http.Dir(webDir))
	mux.Handle("/", fileServer)

	address := ":" + port

	fmt.Printf("Scheduler server started: http://localhost:%s\n", port)

	return http.ListenAndServe(address, mux)
}
