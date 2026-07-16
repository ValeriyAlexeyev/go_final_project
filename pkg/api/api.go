package api

import "net/http"

// Init регистрирует все API-обработчики приложения.
func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
}
