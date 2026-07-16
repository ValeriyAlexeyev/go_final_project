package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, fmt.Errorf("ошибка десериализации JSON: %w", err))
		return
	}

	task.ID = strings.TrimSpace(task.ID)
	task.Title = strings.TrimSpace(task.Title)

	if task.ID == "" {
		writeError(w, fmt.Errorf("не указан идентификатор"))
		return
	}

	if task.Title == "" {
		writeError(w, fmt.Errorf("не указан заголовок задачи"))
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]any{})
}
