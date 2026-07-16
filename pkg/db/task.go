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
func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	if id == "" {
		return nil, fmt.Errorf("не указан идентификатор")
	}

	task := &Task{}

	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`

	err := DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("get task: %w", err)
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	if DB == nil {
		return fmt.Errorf("database is not initialized")
	}

	if task.ID == "" {
		return fmt.Errorf("не указан идентификатор")
	}

	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`

	result, err := DB.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
func DeleteTask(id string) error {
	if DB == nil {
		return fmt.Errorf("database is not initialized")
	}

	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("не указан идентификатор")
	}

	result, err := DB.Exec(
		`DELETE FROM scheduler WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func UpdateDate(nextDate, id string) error {
	if DB == nil {
		return fmt.Errorf("database is not initialized")
	}

	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("не указан идентификатор")
	}

	result, err := DB.Exec(
		`UPDATE scheduler SET date = ? WHERE id = ?`,
		nextDate,
		id,
	)
	if err != nil {
		return fmt.Errorf("update task date: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
