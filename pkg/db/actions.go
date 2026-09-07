package db

import (
	"errors"
	"strconv"
	"time"
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
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTask(id int64) (Task, error) {
	var task Task
	var id64 int64
	err := DB.QueryRow("SELECT * FROM scheduler WHERE id = ?", id).Scan(&id64, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}
	task.ID = strconv.FormatInt(id64, 10)
	return task, nil
}

func UpdateTask(task Task) error {
	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return errors.New("Неверный ID")
	}

	if task.Date == "" || task.Date == "today" {
		task.Date = time.Now().Format(DateFormat)
	} else {
		_, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			return errors.New("Неверный формат даты")
		}
	}

	today := time.Now().Format(DateFormat)
	if task.Date < today {
		task.Date = today
	}

	_, err = DB.Exec(
		"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, id,
	)
	if err != nil {
		return errors.New("Ошибка обновления задачи")
	}
	return nil
}

func DeleteTask(id int64) error {
	_, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return errors.New("Ошибка удваления задачи")
	}
	return nil
}
