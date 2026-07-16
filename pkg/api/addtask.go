package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ValeriyAlexeyev/go_final_project/pkg/db"
)

type addTaskResponse struct {
	ID string `json:"id"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		writeError(w, fmt.Errorf("ошибка десериализации JSON: %w", err))
		return
	}

	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		writeError(w, fmt.Errorf("не указан заголовок задачи"))
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, addTaskResponse{
		ID: fmt.Sprintf("%d", id),
	})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := normalizeDate(now)

	if task.Date == "" {
		task.Date = today.Format(DateFormat)
	}

	taskDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректная дата: %w", err)
	}

	taskDate = normalizeDate(taskDate)
	task.Repeat = strings.TrimSpace(task.Repeat)

	var next string

	if task.Repeat != "" {
		next, err = NextDate(today, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("некорректное правило повторения: %w", err)
		}
	}

	if taskDate.Before(today) {
		if task.Repeat == "" {
			task.Date = today.Format(DateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}
