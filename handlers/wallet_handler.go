package handlers

import (
	"encoding/json"
	"net/http"

	"wallet/models"
	"wallet/services"
	"wallet/utils"

	"github.com/go-chi/chi/v5"
)

// Обработчик операций с кошельком
func HandleWalletOperation(w http.ResponseWriter, r *http.Request) {
	var req models.OperationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	db, err := utils.ConnectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = services.PerformOperation(db, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// Обработчик получения баланса
func GetWalletBalance(w http.ResponseWriter, r *http.Request) {

	walletID := chi.URLParam(r, "walletId")

	if walletID == "" {
		http.Error(w, "Неверный ID кошелька", http.StatusBadRequest)
		return
	}

	db, err := utils.ConnectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	balance, err := services.GetBalance(db, walletID)
	if err != nil {
		if err.Error() == "кошелек не найден" {
			http.Error(w, "кошелек не найден", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]float64{"Баланс": balance})
}
