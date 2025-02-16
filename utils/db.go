package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Проверка существования базы данных и её создание, если её нет
func EnsureDatabaseExists(dbUser, dbPassword, dbName string) error {
	// Подключение к PostgreSQL без указания конкретной базы данных
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		dbUser,
		dbPassword,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к PostgreSQL: %w", err)
	}
	defer db.Close()

	// Проверяем существование базы данных
	var exists bool
	err = db.QueryRow(`SELECT 1 FROM pg_database WHERE datname = $1`, dbName).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("ошибка проверки существования базы данных: %w", err)
	}

	// Если база данных не существует, создаём её
	if !exists {
		log.Printf("Базы данных %s не существует. Создание базы данных...", dbName)
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("ошибка создания базы данных %s: %w", dbName, err)
		}
		log.Printf("База данных %s создана.", dbName)
	} else {
		log.Printf("База данных %s уже существует.", dbName)
	}

	return nil
}

// ConnectDB устанавливает соединение с базой данных
func ConnectDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Успешное подключение к базе данных")

	// Создание таблицы wallets, если её нет
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS wallets (
            id TEXT PRIMARY KEY,
            balance REAL NOT NULL,
            currency TEXT NOT NULL
        );
    `)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	log.Println("Таблица wallets создана или уже существует")

	return db, nil
}
