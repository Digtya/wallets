package tests

import (
	"os"
	"testing"

	"database/sql"
	"wallet/models"
	"wallet/services"
	"wallet/utils"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Настройка тестовой базы данных
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	// Загрузка переменных окружения
	err := godotenv.Load("../config.env")
	require.NoError(t, err, "Ошибка загрузки переменных окружения")

	// Проверяем и создаём базу данных, если её нет
	err = utils.EnsureDatabaseExists(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	require.NoError(t, err, "Ошибка проверки и создания базы данных")

	// Подключение к базе данных
	db, err := utils.ConnectDB()
	require.NoError(t, err, "Ошибка подключения к базе данных")

	// Создаем таблицу wallets
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS wallets (
            id TEXT PRIMARY KEY,
            balance REAL NOT NULL,
            currency TEXT NOT NULL
        );
    `)
	require.NoError(t, err, "Ошибка создания таблицы wallets")

	// Очищаем таблицу перед тестами
	_, err = db.Exec("DELETE FROM wallets")
	require.NoError(t, err, "Ошибка очистки таблицы wallets")

	// Функция для закрытия соединения с базой данных
	teardown := func() {
		err := db.Close()
		require.NoError(t, err, "Ошибка закрытия соединения с базой данных")
	}

	return db, teardown
}

// Тест для PerformOperation
func TestWalletOperation(t *testing.T) {
	// Настройка тестовой базы данных
	db, teardown := setupTestDB(t)
	defer teardown()

	// Добавляем тестовый кошелек в базу данных
	_, err := db.Exec("INSERT INTO wallets (id, balance, currency) VALUES ($1, $2, $3)", "test-wallet-id", 500.0, "USD")
	require.NoError(t, err, "Ошибка добавления кошелька в базу данных")

	// Создаем запрос на депозит
	req := models.OperationRequest{
		WalletID:      "test-wallet-id",
		OperationType: "DEPOSIT",
		Amount:        100.0,
	}

	// Вызываем тестируемую функцию
	err = services.PerformOperation(db, req)
	require.NoError(t, err, "PerformOperation failed")

	// Проверяем, что баланс обновился корректно
	var balance float64
	err = db.QueryRow("SELECT balance FROM wallets WHERE id = $1", "test-wallet-id").Scan(&balance)
	require.NoError(t, err, "Ошибка получения баланса кошелька")
	assert.Equal(t, 600.0, balance, "Неверный баланс кошелька")

	// Создаем запрос на снятие средств
	req = models.OperationRequest{
		WalletID:      "test-wallet-id",
		OperationType: "WITHDRAW",
		Amount:        200.0,
	}

	// Вызываем тестируемую функцию
	err = services.PerformOperation(db, req)
	require.NoError(t, err, "Функция вернула ошибку")

	// Проверяем, что баланс обновился корректно
	err = db.QueryRow("SELECT balance FROM wallets WHERE id = $1", "test-wallet-id").Scan(&balance)
	require.NoError(t, err, "Failed to query wallet balance")
	assert.Equal(t, 400.0, balance, "Неверный баланс кошелька")
}
