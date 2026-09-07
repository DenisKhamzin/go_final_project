package api

import "net/http"

// регистрация хендлеров для каждого эндпоинта
func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDayHandler)   // ожидаем метод GET
	mux.HandleFunc("/api/task", taskHandler)          // ожидаем метод GET
	mux.HandleFunc("/api/tasks", getTasksHandler)     // ожидаем метод GET
	mux.HandleFunc("/api/task/done", doneTaskHandler) // ожидаем метод POST
}

// регистрация хендлеров для эндпоинта "/api/task" в зависимости от метода
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		taskAddHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}
