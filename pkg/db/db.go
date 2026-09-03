package db

import (
	"database/sql"
	// "fmt"
	// "time"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var schema string = `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR DAFAULT "",
		comment TEXT DEFAULT "",
		repeat VARCHAR(128) DEFAULT "");`

var index string = `CREATE INDEX IF NOT EXIST idx_scheduler_date ON scheduler (date);`

func Init(dbName string) error {
	_, err := os.Stat(dbName)
	install := false
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		log.Fatal("Ошибка создания (открытия) базы данных:", err)
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
		return err
	}

	if install == true {
		_, err = db.Exec(schema)
		if err != nil {
			log.Fatal("Ошибка создания таблицы:", err)
			return err
		}
		_, err = db.Exec(index)
		if err != nil {
			log.Fatal("Ошибка создания индекса:", err)
			return err
		}
	}
	return nil
}

func TaskAder(title, date, comment, repeat string) (int64, error) {
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		log.Fatal("Ошибка создания (открытия) базы данных:", err)
		return 0, err
	}
	defer db.Close()

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, date, title, comment, repeat)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
