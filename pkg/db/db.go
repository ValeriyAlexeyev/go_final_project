package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(256) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX scheduler_date ON scheduler(date);
`

// Init открывает SQLite-базу и создаёт таблицу scheduler,
// если файл базы данных ещё не существует.
func Init(dbFile string) error {
	if dbFile == "" {
		return errors.New("database file path is empty")
	}

	_, err := os.Stat(dbFile)

	install := false

	switch {
	case err == nil:
		// Файл существует, создавать схему не требуется.
	case errors.Is(err, os.ErrNotExist):
		install = true
	default:
		return fmt.Errorf("check database file: %w", err)
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	// Проверяем, что соединение действительно работает.
	if err = DB.Ping(); err != nil {
		_ = DB.Close()
		DB = nil
		return fmt.Errorf("ping database: %w", err)
	}

	if install {
		if _, err = DB.Exec(schema); err != nil {
			_ = DB.Close()
			DB = nil
			return fmt.Errorf("create database schema: %w", err)
		}
	}

	return nil
}

// Close закрывает соединение с базой данных.
func Close() error {
	if DB == nil {
		return nil
	}

	return DB.Close()
}
