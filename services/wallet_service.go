package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"wallet/models"
)

// Функция для выполнения операции
func PerformOperation(db *sql.DB, req models.OperationRequest) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var currentBalance float64
	row := tx.QueryRow("SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", req.WalletID)
	err = row.Scan(&currentBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("кошелек не найден")
		}
		return err
	}

	switch req.OperationType {
	case "DEPOSIT":
		currentBalance += req.Amount
	case "WITHDRAW":
		if currentBalance < req.Amount {
			return fmt.Errorf("недостаточно средств для снятия")
		}
		currentBalance -= req.Amount
	default:
		return fmt.Errorf("неправильный тип операции")
	}

	_, err = tx.Exec("UPDATE wallets SET balance = $1 WHERE id = $2", currentBalance, req.WalletID)
	if err != nil {
		return err
	}

	log.Printf("Баланс кошелька %s изменен. Новый баланс: %.2f", req.WalletID, currentBalance)
	return nil
}

// Функция для получения баланса
func GetBalance(db *sql.DB, walletID string) (float64, error) {

	var balance float64
	row := db.QueryRow("SELECT balance FROM wallets WHERE id = $1", walletID)

	err := row.Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("кошелек не найден")
		}
		return 0, fmt.Errorf("ошибка получения баланса: %w", err)
	}

	log.Printf("Получен баланс кошелька %s: %.2f", walletID, balance)
	return balance, nil
}
