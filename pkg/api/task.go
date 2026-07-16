package api

import (
	"fmt"
	"net/http"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
)

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, fmt.Errorf("метод %s не поддерживается", r.Method))
		return
	}

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, TasksResponse{
		Tasks: tasks,
	})
}
