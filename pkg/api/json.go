package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(
	w http.ResponseWriter,
	statusCode int,
	data any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("ошибка записи JSON-ответа: %v", err)
	}
}

func writeError(
	w http.ResponseWriter,
	statusCode int,
	message string,
) {
	writeJSON(
		w,
		statusCode,
		map[string]string{
			"error": message,
		},
	)
}

func writeInternalError(
	w http.ResponseWriter,
	err error,
) {
	log.Printf("внутренняя ошибка: %v", err)

	writeError(
		w,
		http.StatusInternalServerError,
		"внутренняя ошибка сервера",
	)
}
