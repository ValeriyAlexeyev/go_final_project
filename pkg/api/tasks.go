package api

import (
	"fmt"
	"net/http"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPut:
		updateTaskHandler(w, r)

	case http.MethodDelete:
		deleteTaskHandler(w, r)

	default:
		w.Header().Set(
			"Allow",
			http.MethodGet+", "+
				http.MethodPost+", "+
				http.MethodPut+", "+
				http.MethodDelete,
		)

		writeError(w, fmt.Errorf("метод %s не поддерживается", r.Method))
	}
}
