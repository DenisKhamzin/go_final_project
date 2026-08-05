package server

import (
	"log"
	"net/http"
	"time"

	"github.com/deniskhamzin/go_final_project/internal/handlers"
)

// http server struct
type Server struct {
	Logger *log.Logger
	HttpServer *http.Server
}

// function for creating server instance
func NewServer(logger *log.Logger) *Server {
	// creating router
	router := http.NewServeMux()
	// register handlers in current router
	router.HandleFunc("/", handlers.SimpeGetHandler)

	// creating instance if server with requirment params
	serverHttp := &http.Server {
		Addr: ":7540",
		Handler: router,
		ErrorLog: logger,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 15 * time.Second,
	}

	// returning created instance of  http-server
	return &Server {
		Logger: logger,
		HttpServer: serverHttp,
	}
}