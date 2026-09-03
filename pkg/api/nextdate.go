package api

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

func nextDate(now time.Time, dstart string, repeat string) (string, error) {
	startDate, err := time.Parse("20060102", dstart)
	if err != nil {
		log.Printf("Неверный формат dstart: %v", err)
		return "", errors.New("неверный формат даты")
	}

	if repeat == "" {
		return "", errors.New("repeat не может быть пустым")
	}

	slice := strings.Split(repeat, " ")
	if len(slice) == 0 {
		return "", errors.New("пустой repeat")
	}

	if slice[0] != "d" && slice[0] != "y" {
		log.Printf("Неверный тип в repeat: %s", slice[0])
		return "", errors.New("неверный тип в repeat")
	}

	if slice[0] == "y" {
		if len(slice) != 1 {
			log.Printf("Неверное значение после y в repeat: %v", slice)
			return "", errors.New("неверное значение после y")
		}

		// ВАЖНО: сначала прибавляем год, потом проверяем
		startDate = startDate.AddDate(1, 0, 0)

		// Теперь проверяем, что дата > now
		for !startDate.After(now) {
			startDate = startDate.AddDate(1, 0, 0)
		}

		nextDate := startDate.Format("20060102")
		return nextDate, nil
	}

	if slice[0] == "d" {
		if len(slice) != 2 {
			log.Printf("Неверное значение после d в repeat: %v", slice)
			return "", errors.New("неверное значение после d")
		}

		num, err := strconv.Atoi(slice[1])
		if err != nil {
			log.Printf("Невозможно конвертировать в число: %v", err)
			return "", errors.New("неверное число в repeat")
		}

		if num < 1 || num > 400 {
			log.Printf("Некорректное число в repeat: %d", num)
			return "", errors.New("число должно быть от 1 до 400")
		}

		// ВАЖНО: сначала прибавляем дни, потом проверяем
		startDate = startDate.AddDate(0, 0, num)

		// Теперь проверяем, что дата > now
		for !startDate.After(now) {
			startDate = startDate.AddDate(0, 0, num)
		}

		nextDate := startDate.Format("20060102")
		return nextDate, nil
	}

	return "", errors.New("неизвестный тип repeat")
}
