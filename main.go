package main

import (
	"log"
	"net/http"

	"wallet/handlers"
	"wallet/utils"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки переменных окружения: %v", err)
	}

	// Подключаемся к базе данных
	db, err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Организуем роутер
	r := chi.NewRouter()

	// Определяем пути
	r.Post("/api/v1/wallet", handlers.HandleWalletOperation)
	r.Get("/api/v1/wallets/{walletId}", handlers.GetWalletBalance)

	// Запускаем сервер
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
