package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deniskhamzin/go_final_project/pkg/api"
	"github.com/deniskhamzin/go_final_project/pkg/db"
	//"github.com/deniskhamzin/go_final_project/pkg/nextdate"
)

func main() {
	webDir := "./web"
	dbFile := "scheduler.db"

	err := db.Init(dbFile)
	if err != nil {
		log.Fatal("Ошибка создания (открытия) Базы Данных:", err)
	}

	mux := http.NewServeMux()

	api.Init(mux)

	// Создаем файловый сервер
	fileServer := http.FileServer(http.Dir(webDir))

	// Обработчик для всех запросов
	mux.Handle("/", fileServer)

	// Запускаем сервер на порту 7540
	fmt.Println("Сервер запущен на http://localhost:7540")

	if err := http.ListenAndServe(":7540", mux); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}

}
