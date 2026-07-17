package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"ошибка десериализации JSON",
		)
		return
	}

	task.ID = strings.TrimSpace(task.ID)
	task.Title = strings.TrimSpace(task.Title)

	if task.ID == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"не указан идентификатор",
		)
		return
	}

	if task.Title == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"не указан заголовок задачи",
		)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		log.Printf(
			"ошибка обновления задачи %q: %v",
			task.ID,
			err,
		)

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
