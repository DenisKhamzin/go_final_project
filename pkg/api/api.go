package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	//"github.com/deniskhamzin/go_final_project/pkg/nextdate"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDayHandler)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем и парсим параметр now
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

	// 2. Получаем остальные параметры
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	// 3. Вызываем бизнес-логику
	result, err := nextDate(nowTime, date, repeat)
	if err != nil {
		// Логируем ошибку (но НЕ убиваем сервер!)
		log.Printf("Ошибка в nextDate: %v", err)
		http.Error(w, "Невозможно вычислить следующую дату", http.StatusBadRequest)
		return
	}

	// 4. Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func nextDate(now time.Time, dstart string, repeat string) (string, error) {
	// проверка на корректное входное значение dstart
	startDate, err := time.Parse("20060102", dstart)
	if err != nil {
		log.Printf("Неверный формат dstart: %v", err)
		return "", err
	}
	// проверка на корректное значение repeat
	slice := strings.Split(repeat, " ")
	if slice[0] != "d" && slice[0] != "y" {
		err = errors.New("Ошибка в repeat")
		log.Printf("Неверный индекс в repeat: %v", err)
		return "", err
	}
	// проверка на корректность в случае значения `y`
	if slice[0] == "y" && len(slice) != 1 {
		err = errors.New("Значение после `y` в repeat")
		log.Printf("Неверное значение после `y` в repeat: %v", err)
		return "", err
	}
	// проверка на корректность в случае значения `d`
	if slice[0] == "d" && len(slice) != 2 {
		err = errors.New("Значение после `d` в repeat")
		log.Printf("Неверное значение после `d` в repeat: %v", err)
		return "", err
	}
	// проверка на корректность числа после `d` в repeat
	num, err := strconv.Atoi(slice[1])
	if err != nil {
		log.Printf("невозможно конвертировать в число: %v", err)
		return "", err
	}
	// проверка чила на <= 400 по условию задачи
	if num < 1 || num > 400 {
		err = errors.New("Некорректное число в repeat после `d`")
		log.Printf("Некорректный repeat: %v", err)
		return "", err
	}
	// вычисляем следующую дату для `y` и для `d`
	if slice[0] == "y" {
		for afterNow(startDate, now) {
			startDate = startDate.AddDate(1, 0, 0)
		}
	} else {
		for afterNow(startDate, now) {
			startDate = startDate.AddDate(0, 0, num)
		}
	}
	nextDate := startDate.Format("20060102")
	return nextDate, nil
}

func afterNow(date, now time.Time) bool {
	return date.Before(now)
}

//func afterNow(date, now time.Time) bool {
//    return date.Before(now) || date.Equal(now) // нужно добавлять пока дата <= now
//}
