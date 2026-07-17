package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"не указан идентификатор",
		)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		log.Printf("ошибка удаления задачи %q: %v", id, err)

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
		map[string]any{},
	)
}
