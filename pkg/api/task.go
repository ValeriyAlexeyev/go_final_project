package api

import (
	"fmt"
	"net/http"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	default:
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, fmt.Errorf("метод %s не поддерживается", r.Method))
	}
}
