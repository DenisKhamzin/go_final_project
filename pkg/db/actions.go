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

func UpdateTask(task Task) error {
	// 1. Проверка ID
	if task.ID == "" {
		return errors.New("ID не указан")
	}
	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return errors.New("Неверный ID")
	}

	// 2. Проверка существования задачи
	var exists bool
	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM scheduler WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		return errors.New("Задача не найдена")
	}

	// 3. Валидация заголовка
	if task.Title == "" {
		return errors.New("Заголовок не может быть пустым")
	}

	// 4. Обработка даты
	if task.Date == "" || task.Date == "today" {
		task.Date = time.Now().Format("20060102")
	} else {
		if _, err := time.Parse("20060102", task.Date); err != nil {
			return errors.New("Неверный формат даты")
		}
	}

	// 6. Автокоррекция даты в прошлом
	today := time.Now().Format("20060102")
	if task.Date < today {
		task.Date = today
	}

	// 7. Выполнение UPDATE
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
