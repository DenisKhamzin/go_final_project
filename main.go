package main

import (
	"log"
	"os"

	"github.com/deniskhamzin/go_final_project/internal/server"
)

func main() {
	// creating new logger
	logger := log.New(os.Stdout, "http-server", log.LstdFlags|log.Lshortfile)
	// creating new http-server using current logger
	Server := server.NewServer(logger)

	// server starts listening to port 7540
	err := Server.HttpServer.ListenAndServe()
	// logging errors
	if err != nil {
		logger.Fatal("HTTP server didn't start: ", err)
	}
}