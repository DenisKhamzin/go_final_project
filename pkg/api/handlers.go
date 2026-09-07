package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/deniskhamzin/go_final_project/pkg/db"
)

// этот хендлер не нужно переносить
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic в nextDayHandler: %v", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
	}()

	nowStr := r.URL.Query().Get("now")
	var nowTime time.Time
	if nowStr == "" {
		nowTime = time.Now()
	} else {
		var err error
		nowTime, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, "Неверный формат параметра now, ожидается YYYYMMDD", http.StatusBadRequest)
			return
		}
	}

	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if date == "" {
		http.Error(w, "Параметр date обязателен", http.StatusBadRequest)
		return
	}

	if repeat == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
		return
	}

	result, err := nextDate(nowTime, date, repeat)
	if err != nil {
		log.Printf("Ошибка в nextDate: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

// этот хэндлер перенесен
func taskAddHandler(w http.ResponseWriter, r *http.Request) {

	var task db.Task

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	// проверяем, проситался ли JSON
	err := decoder.Decode(&task)
	if err != nil {
		writeError(w, "Неверный формат json")
		return
	}
	// проверяем, что title не пустой
	if task.Title == "" {
		writeError(w, "Title не может быть пустым")
		return
	}
	// если строка date не указана или это пустая строка - берем сегодняшнее число
	if task.Date == "" || task.Date == "today" {
		task.Date = time.Now().Format("20060102")
	}
	// вычисляем сегодняшнее число в формате time.Time (дата корректно распознается)
	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		writeError(w, "Неверный формат даты")
		return
	}
	// Автокоррекция даты в прошлом на сегодняшнюю
	today := time.Now().Format("20060102")
	if task.Date < today {
		task.Date = today
	}

	// проверяем правило repeat
	if task.Repeat != "" {
		if _, err := nextDate(time.Now(), task.Date, task.Repeat); err != nil {
			writeError(w, "Неверное правило повторения")
			return
		}
	}
	// отправляем запрос в бд
	id, err := db.AddTask(&task)

	if err != nil {
		log.Printf("Ошибка добавления задачи: %v", err)
	}
	if id == 0 {
		writeError(w, "Ошибка сохранения задачи")
		return
	}
	// instead of writeID()
	// 1. Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 2. Устанавливаем HTTP статус
	w.WriteHeader(http.StatusOK)
	// 3. Сериализуем и отправляем JSON
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// это не хэндлер
func writeError(w http.ResponseWriter, err string) {
	// 1. Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 2. Устанавливаем HTTP статус
	w.WriteHeader(http.StatusBadRequest)
	// 3. Сериализуем и отправляем JSON
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}

// хендлер перенесен
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.GetTasks()
	if err != nil {
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"tasks": tasks}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// хендлер перенесен
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода (добавлена для полноты)
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Парсинг ID из запроса
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

	// Вызов функции БД
	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, "Задача не найдена")
		} else {
			writeError(w, "Ошибка получения задачи")
		}
		return
	}

	// Отправка успешного ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// хендлер перенесен
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// (Опционально) Проверка метода
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Декодирование JSON
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, "Неверный формат JSON")
		return
	}

	if task.Repeat != "" {
		if _, err := nextDate(time.Now(), task.Date, task.Repeat); err != nil {
			writeError(w, "Неверное правило повторения")
			return
		}
	}

	// Вызов функции обновления
	err = db.UpdateTask(task)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

// HTTP-обработчик, который вызывает функцию БД
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}
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

	task, err := db.GetTask(id)

	if task.Repeat == "" {
		db.DeleteTask(id)
	} else {
		nextDate, err := nextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Ошибка вычисления следующей даты")
			return
		}
		task.Date = nextDate
		err = db.UpdateTask(task)
		if err != nil {
			writeError(w, "Ошибка обновления даты")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	err = db.DeleteTask(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
