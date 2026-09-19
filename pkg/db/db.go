package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// SQL-запрос для добавления новой таблицы
var schema string = `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR DEFAULT "",
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
	// создание таблицы  и индекса в БД в случае их отсутствия
	if install == true {
		_, err = DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("Ошибка создания таблицы: %w", err)
		}
		_, err = DB.Exec(index)
		if err != nil {
			return fmt.Errorf("Ошибка создания индекса: %w", err)
		}
		// закрытие соединения с созданной БД
		DB.Close()
	}

	DB, err = sql.Open("sqlite", dbName)
	if err != nil {
		return fmt.Errorf("Ошибка создания (открытия) базы данных: %w", err)
	}
	// проверка подлючения к БД
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("Ошибка подключения к БД: %w", err)
	}
	// все ошибки уже обработаны
	return nil
}
