-- Шаг 1: Создание таблицы wallets, если её нет
CREATE TABLE IF NOT EXISTS wallets (
    id TEXT PRIMARY KEY,
    balance REAL NOT NULL,
    currency TEXT NOT NULL
);

-- Шаг 2: Добавление записей, если их ещё нет
DO $$
BEGIN
    -- Проверка наличия записи с ID 'user1-wallet'
    IF NOT EXISTS (SELECT 1 FROM wallets WHERE id = 'user1-wallet') THEN
        INSERT INTO wallets (id, balance, currency) VALUES ('user1-wallet', 1000.0, 'USD');
    END IF;

    -- Проверка наличия записи с ID 'user2-wallet'
    IF NOT EXISTS (SELECT 1 FROM wallets WHERE id = 'user2-wallet') THEN
        INSERT INTO wallets (id, balance, currency) VALUES ('user2-wallet', 500.0, 'EUR');
    END IF;

    -- Добавьте больше записей по необходимости
END $$;