package db

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

//type ErrorAdderResponse struct {
//	Error string `json:"error"`
//}

//type IdAdderResponse struct {
//	ID int64 `json:"id"`
//}
