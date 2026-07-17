package api

import (
	"encoding/json"
	"fmt"
	"log"
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
		writeJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": fmt.Sprintf(
					"ошибка десериализации JSON: %v",
					err,
				),
			},
		)
		return
	}

	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		writeJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "не указан заголовок задачи",
			},
		)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": err.Error(),
			},
		)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("ошибка при добавлении задачи: %v", err)

		writeJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "внутренняя ошибка сервера",
			},
		)
		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		addTaskResponse{
			ID: id,
		},
	)
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
			return fmt.Errorf(
				"некорректное правило повторения: %w",
				err,
			)
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
