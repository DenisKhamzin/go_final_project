package db

import (
	"errors"
	"strconv"
)

// функция для добавления задачи
func AddTask(task *Task) (int64, error) {
	var id int64
	// SQL-запрос в БД
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	// выполнение SQL-запроса
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	// возврат id для созданной записи в БД
	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// функция для возврата всех предстоящих задач
func GetTasks() ([]*Task, error) {
	// SQL-запрос в БД
	query := `SELECT id, title, comment, date, repeat FROM scheduler ORDER BY date LIMIT 50`
	// выполнение SQL-запроса
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	// закрытие rows после обработки
	defer rows.Close()
	// список задач для возврата из функции
	tasks := make([]*Task, 0)
	// перебор полученых строк, создание объекта Task для каждой строки
	for rows.Next() {
		var task Task
		var id int64
		// заполнение полей задачи
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		// перевод id в строковый формат
		task.ID = strconv.FormatInt(id, 10)
		// добавление задачи в результирующий слайс
		tasks = append(tasks, &task)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	// возврат результирующего слайса
	return tasks, nil
}

// функция для добавления задачи
func GetTask(id int64) (Task, error) {
	var task Task
	// выполнение запроса в БД по полю id
	err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}
	// перевод int64 в строковый формат для записи в task
	task.ID = strconv.FormatInt(id, 10)
	return task, nil
}

// функция для изменения задачи
func UpdateTask(task Task) error {
	// перевод id в формат int64 для запроса в БД
	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return errors.New("Неверный ID")
	}
	// выполенение SQL-запроса на изменение текущей задачи
	_, err = DB.Exec(
		"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, id,
	)
	if err != nil {
		return errors.New("Ошибка обновления задачи")
	}
	return nil
}

// функция для удаления записи из БД
func DeleteTask(id int64) error {
	// проверка существования задачи с данным id
	_, err := GetTask(id)
	if err != nil {
		return errors.New("Задача с данным id отсутсвует")
	}
	// выполнение SQL-запроса на удаление
	_, err = DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return errors.New("Ошибка удваления задачи")
	}
	return nil
}
