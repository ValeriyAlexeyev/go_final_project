package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, fmt.Errorf("метод %s не поддерживается", r.Method))
		return
	}

	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, fmt.Errorf("не указан идентификатор"))
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err)
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, map[string]any{})
		return
	}

	nextDate, err := NextDate(
		time.Now(),
		task.Date,
		task.Repeat,
	)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := db.UpdateDate(nextDate, id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]any{})
}
