package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, task)
}
