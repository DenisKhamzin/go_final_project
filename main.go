package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deniskhamzin/go_final_project/pkg/api"
	"github.com/deniskhamzin/go_final_project/pkg/db"
)

func main() {
	// объявление переменные для директории со статическими объектами и названием БД
	webDir := "./web"
	dbFile := "scheduler.db"
	// инициация БД
	err := db.Init(dbFile)
	defer db.DB.Close()
	if err != nil {
		// здесь и далее на уровне main.go ошибки логируются в консоли
		log.Fatal("Ошибка создания (открытия) Базы Данных:", err)
	}
	// регистрация роутера
	mux := http.NewServeMux()
	// регистрация хендлеров
	api.Init(mux)
	// регистрация папки со статикой
	fileServer := http.FileServer(http.Dir(webDir))
	mux.Handle("/", fileServer)
	// сообщение об успешном запуске сервера
	fmt.Println("Сервер запущен на http://localhost:7540")
	// запуск сервера на порту 7540, логирование ошибки в случае неудачи
	if err := http.ListenAndServe(":7540", mux); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
