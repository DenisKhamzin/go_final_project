package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
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
func Init(dbName string) error {
	// булевая переменная для необходимости создания таблицы и индекса
	install := false
	// проверка существования файла в корне проекта
	_, err := os.Stat(dbName)
	// в случае, если файл не найден, начение переменной install меняется на true
	if errors.Is(err, fs.ErrNotExist) {
		install = true
	}
	// подготовка структуры для работы с БД
	DB, err = sql.Open("sqlite", dbName)
	if err != nil {
		return fmt.Errorf("Ошибка открытия базы данных: %w", err)
	}
	// установка соединения (и создание файла БД, если его не было)
	err = DB.Ping()
	if err != nil {
		// в случае ошибки соединение принудительно закрывается
		DB.Close()
		return fmt.Errorf("Ошибка подключения к базе данных: %w", err)
	}
	// Если файла раньше не было - создаем таблицу и индекс
	if install == true {
		// попытка создания таблицы
		_, err = DB.Exec(schema)
		if err != nil {
			// в случае ошибки соединение принудительно закрывается
			DB.Close()
			return fmt.Errorf("Ошибка создания таблицы: %w", err)
		}
		// попытка создания индекса
		_, err = DB.Exec(index)
		if err != nil {
			// в случае ошибки соединение принудительно закрывется
			DB.Close()
			return fmt.Errorf("Ошибка создания индекса: %w", err)
		}
	}
	// если ошибок при инициализации не возникает, соединение остается открытым
	return nil
}
