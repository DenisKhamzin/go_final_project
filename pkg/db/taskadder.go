package db

import (
	"database/sql"
	"log"
)

func AddTask(task *Task) (int64, error) {
	db, err := sql.Open("sqlite", "scheduler")
	var id int64
	if err != nil {
		log.Fatal("Ошибка создания (открытия) базы данных:", err)
		return id, err
	}
	defer db.Close()
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		id, err = res.LastInsertId()
	}
	return id, err
}
