package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/deniskhamzin/go_final_project/pkg/db"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	var nowTime time.Time
	if nowStr == "" {
		nowTime = time.Now()
	} else {
		var err error
		nowTime, err = time.Parse(db.DateFormat, nowStr)
		if err != nil {
			writeError(w, "Неверный формат параметра now, ожидается YYYYMMDD")
			return
		}
	}
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	if date == "" {
		writeError(w, "Параметр date обязателен")
		return
	}
	if repeat == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
		return
	}
	result, err := nextDate(nowTime, date, repeat)
	if err != nil {
		writeError(w, "Ошибка в вычислении next date")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func taskAddHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&task)
	if err != nil {
		writeError(w, "Неверный формат json")
		return
	}
	if task.Title == "" {
		writeError(w, "Title не может быть пустым")
		return
	}
	if task.Date == "" || task.Date == "today" {
		task.Date = time.Now().Format(db.DateFormat)
	}
	_, err = time.Parse(db.DateFormat, task.Date)
	if err != nil {
		writeError(w, "Неверный формат даты")
		return
	}
	today := time.Now().Format(db.DateFormat)
	if task.Date < today {
		task.Date = today
	}
	if task.Repeat != "" {
		now := time.Now()
		_, err := nextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Неверное правило повторения")
			return
		}
	}

	id, err := db.AddTask(&task)

	if err != nil {
		writeError(w, "Ошибка добавления задачи")
		return
	}
	if id == 0 {
		writeError(w, "Ошибка сохранения задачи")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func writeError(w http.ResponseWriter, err string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetTasks()
	if err != nil {
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}
	response := map[string][]db.Task{"tasks": tasks}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, "Задача не найдена")
			return
		} else {
			writeError(w, "Ошибка получения задачи")
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, "Неверный формат JSON")
		return
	}
	if task.Repeat != "" {
		now := time.Now()
		_, err := nextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Неверное правило повторения")
			return
		}
	}
	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		writeError(w, "Неверный ID")
		return
	}
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
	if task.Title == "" {
		writeError(w, "Поле title не может быть пустым")
		return
	}
	err = db.UpdateTask(task)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Поддерживается только метод POST")
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
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, "Ошибка удаления задачи")
		}
	} else {
		now := time.Now()
		nextDate, err := nextDate(now, task.Date, task.Repeat)
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
	w.Write([]byte("{}"))
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
	if err != nil {
		writeError(w, "Ошибка удаления задачи")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
