package api

import "net/http"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("POST /api/task", taskAddHandler)
	mux.HandleFunc("GET /api/tasks", getTasksHandler)
	mux.HandleFunc("GET /api/task", getTaskHandler)
	mux.HandleFunc("PUT /api/task", updateTaskHandler)
	mux.HandleFunc("POST /api/task/done", doneTaskHandler)
	mux.HandleFunc("DELETE /api/task", deleteTaskHandler)
}

// go test -run ^TestApp$ ./tests
// go test -run ^TestDB$ ./tests
// go test -run ^TestNextDate$ ./tests
// go test -run ^TestAddTask$ ./tests
// go test -run ^TestTasks$ ./tests
// go test -run ^TestTask$ ./tests
// go test -run ^TestDone$ ./tests
// go test -run ^TestDelTask$ ./tests
