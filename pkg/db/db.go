package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var schema string = `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR DAFAULT "",
		comment TEXT DEFAULT "",
		repeat VARCHAR(128) DEFAULT "");`

var index string = `CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);`
var DB *sql.DB

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func Init(dbName string) error {
	_, err := os.Stat(dbName)
	install := false
	if err != nil {
		install = true
	}

	DB, err = sql.Open("sqlite", dbName)
	if err != nil {
		log.Fatal("Ошибка создания (открытия) базы данных:", err)
		return err
	}
	//defer DB.Close()

	err = DB.Ping()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
		return err
	}

	if install == true {
		_, err = DB.Exec(schema)
		if err != nil {
			log.Fatal("Ошибка создания таблицы:", err)
			return err
		}
		_, err = DB.Exec(index)
		if err != nil {
			log.Fatal("Ошибка создания индекса:", err)
			return err
		}
	}
	return err
}
