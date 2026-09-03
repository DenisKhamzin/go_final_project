package db

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type ErrorAdderResponse struct {
	Error string `json:"error"`
}

type IdAdderResponse struct {
	ID int64 `json:"id"`
}
