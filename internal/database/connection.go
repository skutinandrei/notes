package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func InitDB() error {

	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, используем значения по умолчанию")
	}

	dbPath := os.Getenv("DB_PATH")

	if dbPath == "" {
		dbPath = "./note.db"
		log.Printf("DB_PATH не указан, используем: %s", dbPath)
	}

	DB, err = sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("ошибка открытия подключения к БД: %w", err)
	}

	_, err = DB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return fmt.Errorf("ошибка включения foreign keys: %w", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("ошибка проверки подключения к БД: %w", err)
	}

	log.Printf("Подключение к базе данных успешно: %s", dbPath)
	return nil
}

func GetDB() *sqlx.DB {
	return DB
}

func CloseDB() error {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			return fmt.Errorf("ошибка закрытия подключения к БД: %w", err)
		}
		log.Println("Подключение к базе данных закрыто")
	}
	return nil
}
