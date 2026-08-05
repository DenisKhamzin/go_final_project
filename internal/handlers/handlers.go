package handlers

import (
	//"fmt"
	//"io"
	"net/http"
	//"os"
	//"time"
	//"github.com/deniskhamzin/go_final_project/web"
)

func SimpeGetHandler(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./web/pep.html")
}
