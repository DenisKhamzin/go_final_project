package db

import (
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

// константа для перевода формата времени в строку и в time.Time
const DateFormat string = "20060102"

// SQL-запрос для добавления новой таблицы
var schema string = `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR DAFAULT "",
		comment TEXT DEFAULT "",
		repeat VARCHAR(128) DEFAULT "");`

// SQL-запрос для создания индекса по полю date
var index string = `CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);`

// переменная для хранения подключения к БД
var DB *sql.DB

// структура для маршализации и анмаршализации json с полями задачи
type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// функция для инициации БД
// на этом уровне
func Init(dbName string) error {
	// проверяем наличие файла БД
	_, err := os.Stat(dbName)
	install := false
	if err != nil {
		install = true
	}
	DB, err = sql.Open("sqlite", dbName)
	if err != nil {
		return errors.New("Ошибка создания (открытия) базы данных:")
	}
	// проверка подлючения к БД
	err = DB.Ping()
	if err != nil {
		return errors.New("Ошибка подключения к БД:")
	}
	// создание таблицы и индекса в случае, если install == true
	if install == true {
		_, err = DB.Exec(schema)
		if err != nil {
			return errors.New("Ошибка создания таблицы:")
		}
		_, err = DB.Exec(index)
		if err != nil {
			return errors.New("Ошибка создания индекса:")
		}
	}
	// все возможные ошибки уже обработаны
	return nil
}
