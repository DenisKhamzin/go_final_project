package db

import (
	"strconv"
)

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetTasks() ([]Task, error) {
	query := `SELECT * FROM scheduler ORDER BY date`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTask(id int64) (Task, error) {
	var task Task
	var id64 int64
	err := DB.QueryRow("SELECT * FROM scheduler WHERE id = ?", id).Scan(&id64, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err // возвращаем пустую задачу и ошибку
	}
	task.ID = strconv.FormatInt(id64, 10)
	return task, nil
}
