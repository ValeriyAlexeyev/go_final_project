package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		writeError(
			w,
			http.StatusMethodNotAllowed,
			fmt.Sprintf("метод %s не поддерживается", r.Method),
		)
		return
	}

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

	if strings.TrimSpace(task.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			log.Printf("ошибка удаления выполненной задачи %q: %v", id, err)

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
		return
	}

	nextDate, err := NextDate(
		time.Now(),
		task.Date,
		task.Repeat,
	)
	if err != nil {
		log.Printf(
			"ошибка расчёта следующей даты задачи %q: %v",
			id,
			err,
		)

		writeError(
			w,
			http.StatusInternalServerError,
			"внутренняя ошибка сервера",
		)
		return
	}

	if err := db.UpdateDate(nextDate, id); err != nil {
		log.Printf("ошибка обновления даты задачи %q: %v", id, err)

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
