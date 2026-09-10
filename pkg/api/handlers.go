package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/deniskhamzin/go_final_project/pkg/db"
)

// хендлер для вычисления следующей даты согласно правилам repeat
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowString := r.URL.Query().Get("now")
	// переменная для отправки в функцию nextDate()
	var nowTime time.Time
	if nowString == "" {
		nowTime = time.Now()
	} else {
		var err error
		nowTime, err = time.Parse(db.DateFormat, nowString)
		if err != nil {
			writeError(w, "Неверный формат параметра now, ожидается YYYYMMDD")
			return
		}
	}
	// переменные для отправки в функцию nextDate()
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	if date == "" {
		writeError(w, "Параметр date обязателен")
		return
	}
	if repeat == "" {
		w.WriteHeader(http.StatusOK)
		return
	}
	// вызов функции nextDate() с полученными аргументами
	result, err := nextDate(nowTime, date, repeat)
	if err != nil {
		writeError(w, "Ошибка в вычислении next date")
		return
	}
	// запись результата в тело ответа
	w.Write([]byte(result))
}

// хэндлер для добавления задачи
func taskAddHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	// заполнение структуры db.Task значениями из request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&task)
	// проверка корректности присланного json
	if err != nil {
		writeError(w, "Неверный формат json")
		return
	}
	// проверка Title на заполненность
	if task.Title == "" {
		writeError(w, "Title не может быть пустым")
		return
	}
	// автозаполнение Date в случаях, если указано сегодняшнее число или оно не указано
	if task.Date == "" || task.Date == "today" {
		task.Date = time.Now().Format(db.DateFormat)
	}
	// проверка формата присланной даты
	_, err = time.Parse(db.DateFormat, task.Date)
	if err != nil {
		writeError(w, "Неверный формат даты")
		return
	}
	// если присланная дата уже прошла, происходит подстановка сегодняшней даты
	today := time.Now().Format(db.DateFormat)
	if task.Date < today {
		task.Date = today
	}
	// проверка корректности полученного правила repeat функцией nextDate()
	if task.Repeat != "" {
		now := time.Now()
		_, err := nextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Неверное правило повторения")
			return
		}
	}
	// вызов функции для добавления задачи в БД
	id, err := db.AddTask(&task)
	// запись ответа в случая ошибки при помощи функции writeError()
	if err != nil || id == 0 {
		writeError(w, "Ошибка добавления задачи")
		return
	}
	// запись json с id добавленной задачи в случае успешнгого добавления
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// вспомогательная функция для записи в ответ json с ошибкой
func writeError(w http.ResponseWriter, err string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}

// хендлер для получения списка предстоящих задач
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	// вызов функции для получения слайса с объектами db.Task
	tasks, err := db.GetTasks()
	if err != nil {
		writeError(w, "Ошибка запроса к БД")
		return
	}
	// результрующая строковая переменная в формате json
	response := map[string][]db.Task{"tasks": tasks}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// хендлер для получения задачи по id
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	// проверка на отсутствие значения
	if idString == "" {
		writeError(w, "ID не указан")
		return
	}
	// перевод строки с id в формат int64
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		writeError(w, "Неверный ID")
		return
	}
	// вызов функции для обращения к БД
	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, "Задача не найдена")
			return
		} else {
			writeError(w, "Ошибка получения задачи")
			return
		}
	}
	// запись json в тело ответа в случае успешного получения задачи
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// хендлер для изменения задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	// заполнение структуры db.Task значениями из request
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, "Неверный формат JSON")
		return
	}
	// проверка repeat на соответствие правилам повторения
	if task.Repeat != "" {
		now := time.Now()
		_, err := nextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Неверное правило повторения")
			return
		}
	}
	// перевод строки с id в формат int64
	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		writeError(w, "Неверный ID")
		return
	}
	// проверка поля title нулевое значение
	if task.Title == "" {
		writeError(w, "Поле title не может быть пустым")
		return
	}
	// после проверки всех полей необходимо проверить, есть ли в БД задача с соотвествующим id
	_, err = db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, "Задача не найдена")
			return
		} else {
			writeError(w, "Ошибка получения задачи")
			return
		}
	}
	// вызов функции для измениния задачи в БД
	err = db.UpdateTask(task)
	if err != nil {
		writeError(w, "Ошибка измения задачи в БД")
		return
	}
	// запись пустого json в случае успешного изменения задачи
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

// хендлер для
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// проверка http-метода
	if r.Method != http.MethodPost {
		writeError(w, "Поддерживается только метод POST")
		return
	}
	// проверка поля id на отсутствие значения
	idString := r.URL.Query().Get("id")
	if idString == "" {
		writeError(w, "ID не указан")
		return
	}
	// перевод строкового id вы формат int64
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		writeError(w, "Неверный ID")
		return
	}
	// получение задачи по id для проверки параметров repeat
	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, "Задача не найдена")
			return
		} else {
			writeError(w, "Ошибка получения задачи")
			return
		}
	}
	// в случае пустого поля repeat задача удаляется
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, "Ошибка удаления задачи")
		}
		// в случае, если repeat указан, вычисляется следующая дата задачи путем вызова nextDate()
	} else {
		now := time.Now()
		nextDate, err := nextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Ошибка вычисления следующей даты")
			return
		}
		// обновление поля Date у задачи
		task.Date = nextDate
		// вызов функции для изменения задачи в БД
		err = db.UpdateTask(task)
		if err != nil {
			writeError(w, "Ошибка обновления даты")
			return
		}
	}
	// возврат пустого json в теле ответа в случае успешного изменения статуса задачи
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

// хендлер для удаления задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// получение id из теля запроса, проверки на пустое значение и перевод в формат int64
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, "ID не указан")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, "Неверный ID")
		return
	}
	// вызоы функции для удаления задачи из БД
	err = db.DeleteTask(id)
	if err != nil {
		writeError(w, "Ошибка удаления задачи")
	}
	// запись в тело ответа пустого json в случе успешного удаления задачи
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
