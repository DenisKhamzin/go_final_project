package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/deniskhamzin/go_final_project/pkg/db"
)

// хендлер для вычисления следующей даты согласно правилам repeat
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
	nowString := r.URL.Query().Get("now")
	// переменная для отправки в функцию nextDate()
	var nowTime time.Time
	if nowString == "" {
		nowTime = time.Now()
	} else {
		var err error
		nowTime, err = time.Parse(DateFormat, nowString)
		if err != nil {
			writeError(w, "Неверный формат параметра now, ожидается YYYYMMDD", http.StatusBadRequest)
			return
		}
	}
	// переменные для отправки в функцию nextDate()
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	if date == "" {
		writeError(w, "Параметр date обязателен", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		writeError(w, "Не указан repeat", http.StatusBadRequest)
		return
	}
	// вызов функции nextDate() с полученными аргументами
	result, err := nextDate(nowTime, date, repeat)
	if err != nil {
		writeError(w, "Ошибка в вычислении next date", http.StatusInternalServerError)
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
		writeError(w, "Неверный формат json", http.StatusBadRequest)
		return
	}
	// проверка Title на заполненность
	if task.Title == "" {
		writeError(w, "Title не может быть пустым", http.StatusBadRequest)
		return
	}
	// автозаполнение Date в случаях, если указано сегодняшнее число или оно не указано
	if task.Date == "" || task.Date == "today" {
		task.Date = time.Now().Format(DateFormat)
	}
	// проверка формата присланной даты
	_, err = time.Parse(DateFormat, task.Date)
	if err != nil {
		writeError(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}
	// проверка корректности полученного правила repeat функцией nextDate()
	if task.Repeat != "" {
		_, err := nextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Неверное правило повторения", http.StatusBadRequest)
			return
		}
	}
	// если присланная дата уже прошла, происходит подстановка сегодняшней даты
	today := time.Now().Format(DateFormat)
	if task.Date < today {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date, err = nextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				writeError(w, "Ошибка вычисления даты", http.StatusInternalServerError)
			}
		}
	}
	// вызов функции для добавления задачи в БД
	id, err := db.AddTask(&task)
	// запись ответа в случая ошибки при помощи функции writeError()
	if err != nil || id == 0 {
		writeError(w, "Ошибка добавления задачи", http.StatusInternalServerError)
		return
	}
	idString := strconv.Itoa(int(id))
	// запись json с id добавленной задачи в случае успешнгого добавления
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// так как при записи http-response заголовки и статус-код уже записаны, ошибка только логируется
	err = json.NewEncoder(w).Encode(map[string]string{"id": idString})
	if err != nil {
		log.Println("Ошибка сериализации json для ответа", err)
	}

}

// вспомогательная функция для записи в ответ json с ошибкой
func writeError(w http.ResponseWriter, err string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	// так как при записи http-response заголовок и статус-код уже записаны, ошибка только логируется
	errJson := json.NewEncoder(w).Encode(map[string]string{"error": err})
	if errJson != nil {
		log.Println("Ошибка записи json в тело ответа при вызове writeError: ", errJson)
	}
}

// хендлер для получения списка предстоящих задач
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	// вызов функции для получения слайса с объектами db.Task
	tasks, err := db.GetTasks()
	if err != nil {
		writeError(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}
	// результрующая строковая переменная в формате json
	response := map[string][]*db.Task{"tasks": tasks}
	w.Header().Set("Content-Type", "application/json")
	// так как при записи http-response заголовки и статус-код уже записаны, возможная ошибка только логируется
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Ошибка сериализации json для ответа: ", err)
	}
}

// хендлер для получения задачи по id
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	// проверка на отсутствие значения
	if idString == "" {
		writeError(w, "ID не указан", http.StatusBadRequest)
		return
	}
	// перевод строки с id в формат int64
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		writeError(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	// вызов функции для обращения к БД
	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, "Задача не найдена", http.StatusInternalServerError)
			return
		} else {
			writeError(w, "Ошибка получения задачи", http.StatusInternalServerError)
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
		writeError(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}
	// проверка Title на заполненность
	if task.Title == "" {
		writeError(w, "Title не может быть пустым", http.StatusBadRequest)
		return
	}
	// автозаполнение Date в случаях, если указано сегодняшнее число или оно не указано
	if task.Date == "" || task.Date == "today" {
		task.Date = time.Now().Format(DateFormat)
	}
	// проверка формата присланной даты
	_, err = time.Parse(DateFormat, task.Date)
	if err != nil {
		writeError(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}
	// проверка корректности полученного правила repeat функцией nextDate()
	if task.Repeat != "" {
		_, err := nextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Неверное правило повторения", http.StatusBadRequest)
			return
		}
	}
	// если присланная дата уже прошла, происходит подстановка сегодняшней даты
	today := time.Now().Format(DateFormat)
	if task.Date < today {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date, err = nextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				writeError(w, "Ошибка вычисления даты", http.StatusInternalServerError)
			}
		}
	}
	// перевод строки с id в формат int64
	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		writeError(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	// после проверки всех полей необходимо проверить, есть ли в БД задача с соотвествующим id
	_, err = db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, "Задача не найдена", http.StatusInternalServerError)
			return
		} else {
			writeError(w, "Ошибка получения задачи", http.StatusInternalServerError)
			return
		}
	}
	// вызов функции для измениния задачи в БД
	err = db.UpdateTask(task)
	if err != nil {
		writeError(w, "Ошибка измения задачи в БД", http.StatusInternalServerError)
		return
	}
	// запись пустого json в случае успешного изменения задачи
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

// хендлер для отметки задачи как выполненной
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// проверка http-метода
	if r.Method != http.MethodPost {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	// проверка поля id на отсутствие значения
	idString := r.URL.Query().Get("id")
	if idString == "" {
		writeError(w, "ID не указан", http.StatusBadRequest)
		return
	}
	// перевод строкового id вы формат int64
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		writeError(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	// получение задачи по id для проверки параметров repeat
	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, "Задача не найдена", http.StatusInternalServerError)
			return
		} else {
			writeError(w, "Ошибка получения задачи", http.StatusInternalServerError)
			return
		}
	}
	// в случае пустого поля repeat задача удаляется
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, "Ошибка удаления задачи", http.StatusInternalServerError)
		}
		// в случае, если repeat указан, вычисляется следующая дата задачи путем вызова nextDate()
	} else {
		now := time.Now()
		nextDate, err := nextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Ошибка вычисления следующей даты", http.StatusInternalServerError)
			return
		}
		// обновление поля Date у задачи
		task.Date = nextDate
		// вызов функции для изменения задачи в БД
		err = db.UpdateTask(task)
		if err != nil {
			writeError(w, "Ошибка обновления даты", http.StatusInternalServerError)
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
		writeError(w, "ID не указан", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	// вызоы функции для удаления задачи из БД
	err = db.DeleteTask(id)
	if err != nil {
		writeError(w, "Ошибка удаления задачи", http.StatusInternalServerError)
		return
	}
	// запись в тело ответа пустого json в случе успешного удаления задачи
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
