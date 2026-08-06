package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	webDir := "./web"

	// Создаем файловый сервер
	fileServer := http.FileServer(http.Dir(webDir))

	// Обработчик для всех запросов
	http.Handle("/", fileServer)

	// Запускаем сервер на порту 7540
	fmt.Println("Сервер запущен на http://localhost:7540")

	if err := http.ListenAndServe(":7540", nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
