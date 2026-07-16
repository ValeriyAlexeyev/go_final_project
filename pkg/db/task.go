package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("database is not initialized")
	}

	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`

	result, err := DB.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get inserted task id: %w", err)
	}

	return id, nil
}
func Tasks(limit int, search string) ([]*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	if limit <= 0 {
		limit = 50
	}

	tasks := make([]*Task, 0)

	search = strings.TrimSpace(search)

	var (
		rows *sql.Rows
		err  error
	)

	if search == "" {
		query := `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			ORDER BY date, id
			LIMIT ?
		`

		rows, err = DB.Query(query, limit)
	} else if date, parseErr := time.Parse("02.01.2006", search); parseErr == nil {
		query := `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE date = ?
			ORDER BY date, id
			LIMIT ?
		`

		rows, err = DB.Query(query, date.Format("20060102"), limit)
	} else {
		query := `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE title LIKE ? OR comment LIKE ?
			ORDER BY date, id
			LIMIT ?
		`

		pattern := "%" + search + "%"
		rows, err = DB.Query(query, pattern, pattern, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}

		if err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}
