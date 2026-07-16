package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, fmt.Errorf("не указан идентификатор"))
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]any{})
}
