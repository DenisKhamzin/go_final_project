package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/deniskhamzin/go_final_project/pkg/db"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("POST /api/task", taskAddHandler)
}

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

func taskAddHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic в taskAdder: %v", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			//resp := db.ErrorAdderResponse{Error: "Внутренняя ошибка сервера"}
			writeError(w, "Внутренняя ошибка сервера")
			return
		}
	}()

	var task db.Task

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	// проверяем, проситался ли JSON
	err := decoder.Decode(&task)
	if err != nil {
		//resp := db.ErrorAdderResponse{Error: "ошибка декодирования входящего json"}
		//err = writeError(w, resp)
		//return
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

	//	// проверяем, что дата задачи позже сегодняшнего числа
	//	if dateTask.Before(time.Now()) {
	//		if task.Repeat == "" {
	//			task.Date = time.Now().Format("20060102")
	//		} else {
	//			task.Date, err = nextDate(time.Now(), task.Date, task.Repeat)
	//			if err != nil {
	//				resp := db.ErrorAdderResponse{Error: "Некорректный repeat"}
	//				err = writeError(w, resp)
	//				return
	//			}
	//		}
	//	}

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
		log.Printf("Ошибка добавления задачи: %w", err)
	}
	if id == 0 {
		writeError(w, "Ошибка сохранения задачи")
		return
	}
	//resp := db.IdAdderResponse{ID: id}
	writeID(w, id)
}

func writeID(w http.ResponseWriter, id int64) {
	// 1. Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// 2. Устанавливаем HTTP статус
	w.WriteHeader(http.StatusOK)

	// 3. Сериализуем и отправляем JSON
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func writeError(w http.ResponseWriter, err string) {
	// 1. Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 2. Устанавливаем HTTP статус
	w.WriteHeader(http.StatusBadRequest)
	// 3. Сериализуем и отправляем JSON
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}
