package nextdate

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// проверка на корректное входное значение dstart
	startDate, err := time.Parse("20060102", dstart)
	if err != nil {
		log.Fatal("Неверный формат dstart:", err)
		return "", err
	}
	// проверка на корректное значение repeat
	slice := strings.Split(repeat, " ")
	if slice[0] != "d" && slice[0] != "y" {
		err = errors.New("Ошибка в repeat")
		log.Fatal("Неверный индекс в repeat:", err)
		return "", err
	}
	// проверка на корректность в случае значения `y`
	if slice[0] == "y" && len(slice) != 1 {
		err = errors.New("Значение после `y` в repeat")
		log.Fatal("Неверное значение после `y` в repeat:", err)
		return "", err
	}
	// проверка на корректность в случае значения `d`
	if slice[0] == "d" && len(slice) != 2 {
		err = errors.New("Значение после `d` в repeat")
		log.Fatal("Неверное значение после `d` в repeat:", err)
		return "", err
	}
	// проверка на корректность числа после `d` в repeat
	num, err := strconv.Atoi(slice[1])
	if err != nil {
		log.Fatal("невозможно конвертировать в число:", err)
		return "", err
	}
	// проверка чила на <= 400 по условию задачи
	if num < 1 || num > 400 {
		err = errors.New("Некорректное число в repeat после `d`")
		log.Fatal("Некорректный repeat:", err)
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
	return !date.Before(now)
}
