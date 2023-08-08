-- Создание таблицы с балансами пользователей
CREATE TABLE IF NOT EXISTS balances (
    sender_id SERIAL PRIMARY KEY,
    balance FLOAT(2) NOT NULL
);

-- Создание таблицы с переводами
CREATE TABLE IF NOT EXISTS transfers (
    request_id UUID PRIMARY KEY,
    sender_id INT NOT NULL,
    amount FLOAT(2) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'successful',
    FOREIGN KEY (sender_id) REFERENCES balances (sender_id) ON DELETE CASCADE
);