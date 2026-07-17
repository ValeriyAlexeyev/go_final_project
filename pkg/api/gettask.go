package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"не указан идентификатор",
		)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("ошибка получения задачи %q: %v", id, err)

		writeError(
			w,
			http.StatusInternalServerError,
			"внутренняя ошибка сервера",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		task,
	)
}
