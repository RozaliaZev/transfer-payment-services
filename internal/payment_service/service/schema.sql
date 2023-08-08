-- Создание таблицы с балансами пользователей
CREATE TABLE IF NOT EXISTS balances (
		sender_id VARCHAR(20) PRIMARY KEY,
		balance FLOAT(2) NOT NULL
	);

-- Создание таблицы с переводами
CREATE TABLE IF NOT EXISTS transfers (
        request_id VARCHAR(20) NOT NULL,
        sender_id VARCHAR(20) NOT NULL,
        amount FLOAT(2) NOT NULL,
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        status VARCHAR(20) DEFAULT 'successful',
        FOREIGN KEY (sender_id) REFERENCES balances_test (sender_id) ON DELETE CASCADE
    );