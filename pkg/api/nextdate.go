package api

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/deniskhamzin/go_final_project/pkg/db"
)

// функция для вычисления следующей даты согласно правилам repeat
func nextDate(now time.Time, dstart string, repeat string) (string, error) {
	// проверяем на корректность входящую строку dstart
	startDate, err := time.Parse(db.DateFormat, dstart)
	if err != nil {
		log.Printf("Неверный формат dstart: %v", err)
		return "", errors.New("неверный формат даты")
	}
	// проверяем, что repeat не пустой
	if repeat == "" {
		return "", errors.New("repeat не может быть пустым")
	}
	// обработка repeat
	slice := strings.Split(repeat, " ")
	if len(slice) == 0 {
		return "", errors.New("неверный тип repeat")
	}
	// хендлер обрабатывает правила только для "y" и "d"
	if slice[0] != "d" && slice[0] != "y" {
		log.Printf("Неверный тип в repeat: %s", slice[0])
		return "", errors.New("неверный тип в repeat")
	}
	// вычисление следующей даты по правилу "y"
	if slice[0] == "y" {
		if len(slice) != 1 {
			log.Printf("Неверное значение после y в repeat: %v", slice)
			return "", errors.New("неверное значение после y")
		}
		// сначала прибавляем год, потом проверяем
		startDate = startDate.AddDate(1, 0, 0)
		// теперь проверяем, что дата > now
		for !startDate.After(now) {
			startDate = startDate.AddDate(1, 0, 0)
		}
		// возврат next date в виде строки
		nextDate := startDate.Format(db.DateFormat)
		return nextDate, nil
	}
	// вычисление следующей даты по правилу для "d"
	if slice[0] == "d" {
		if len(slice) != 2 {
			log.Printf("Неверное значение после d в repeat: %v", slice)
			return "", errors.New("неверное значение после d")
		}
		// проверка на коректность второго значения slice
		num, err := strconv.Atoi(slice[1])
		if err != nil {
			log.Printf("Невозможно конвертировать в число: %v", err)
			return "", errors.New("неверное число в repeat")
		}
		// проверка максимально возможного значения "d"
		if num < 1 || num > 400 {
			log.Printf("Некорректное число в repeat: %d", num)
			return "", errors.New("число должно быть от 1 до 400")
		}
		// сначала прибавляем дни, потом проверяем
		startDate = startDate.AddDate(0, 0, num)
		// теперь проверяем, что дата > now
		for !startDate.After(now) {
			startDate = startDate.AddDate(0, 0, num)
		}
		// возврат next date в виде строки
		nextDate := startDate.Format(db.DateFormat)
		return nextDate, nil
	}
	// возврат ошибки в случае некорретного repeat
	return "", errors.New("неизвестный тип repeat")
}
