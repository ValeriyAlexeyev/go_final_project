package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
)

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)

		writeError(
			w,
			http.StatusMethodNotAllowed,
			fmt.Sprintf(
				"метод %s не поддерживается",
				r.Method,
			),
		)
		return
	}

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		log.Printf("ошибка получения списка задач: %v", err)

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
		TasksResponse{
			Tasks: tasks,
		},
	)
}
