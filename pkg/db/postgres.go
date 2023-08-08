package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

/*
var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var sslmode = os.Getenv("SSLMODE")
var dbname = os.Getenv("DBNAME")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
*/
var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", "localhost", "5432", "postgres", "Password10", "base", "disable")

type TransferPayment struct {
	SenderId  string
	RequestId string
	Amount    float64
}

func CreateTableBalances() error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `CREATE TABLE IF NOT EXISTS balances_test (
		sender_id VARCHAR(20) PRIMARY KEY,
		balance FLOAT(2) NOT NULL
	)`

	_, err = db.Exec(query)
	if err != nil {
		log.Println("error creating table:", err)
		return err
	}

	return nil
}

func CreateTableTransfers(senderId string) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	tableName := fmt.Sprintf("transfers_test_%v", senderId)

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        request_id VARCHAR(20) NOT NULL,
        sender_id VARCHAR(20) NOT NULL,
        amount FLOAT(2) NOT NULL,
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        status VARCHAR(20) DEFAULT 'successful',
        FOREIGN KEY (sender_id) REFERENCES balances_test (sender_id) ON DELETE CASCADE
    )`, tableName)

	_, err = db.Exec(query)
	if err != nil {
		log.Println("error creating table:", err)
		return err
	}

	return nil
}

func CheckIdRepeatition(senderId string, requestId string) (bool, error) {
	err := CreateTableTransfers(senderId)
	if err != nil {
		return false, err
	}

	conn, err := pgx.Connect(context.Background(), dbInfo)
	if err != nil {
		return false, err
	}
	defer conn.Close(context.Background())

	tableName := fmt.Sprintf("transfers_test_%v", senderId)

	var count int
	err = conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM "+tableName+" WHERE request_id = $1", requestId).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func CheckAndChangeBalance(senderId string, amount float64) (bool, error) {
	db, err := pgx.Connect(context.Background(), dbInfo)
	if err != nil {
		return false, err
	}
	defer db.Close(context.Background())

	query := `SELECT balance FROM balances_test WHERE sender_id = $1`
	var balance float64
	err = db.QueryRow(context.Background(), query, senderId).Scan(&balance)
	if err == pgx.ErrNoRows {
		err = createNewBalance(db, senderId)
		if err != nil {
			return false, err
		}
		return false, nil
	} else if err != nil {
		return false, err
	}

	if balance >= amount {
		query := `UPDATE balances_test SET balance = balance - $1 WHERE sender_id = $2`
		_, err = db.Exec(context.Background(), query, amount, senderId)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func createNewBalance(db *pgx.Conn, senderId string) error {
	query := `INSERT INTO balances_test (sender_id, balance) VALUES ($1, 0)`
	_, err := db.Exec(context.Background(), query, senderId)
	if err != nil {
		log.Println("error creating balance:", err)
		return err
	}

	return nil
}

func AddRequestTransferPayment(transferPayment *TransferPayment) error {
	db, err := pgx.Connect(context.Background(), dbInfo)
	if err != nil {
		return err
	}
	defer db.Close(context.Background())

	tableName := fmt.Sprintf("transfers_test_%v", transferPayment.SenderId)

	query := fmt.Sprintf(`INSERT INTO %s (request_id, sender_id, amount) VALUES ($1, $2, $3)`, tableName)
	_, err = db.Exec(context.Background(), query, transferPayment.RequestId, transferPayment.SenderId, transferPayment.Amount)
	if err != nil {
		log.Println("error creating transfer:", err)
		return err
	}

	return nil
}

func SetStatusRequestTransferPayment(transferPayment *TransferPayment) error {
	db, err := pgx.Connect(context.Background(), dbInfo)
	if err != nil {
		return err
	}
	defer db.Close(context.Background())

	tableName := fmt.Sprintf("transfers_test_%v", transferPayment.SenderId)

	query := fmt.Sprintf(`UPDATE %s SET status = $1 WHERE request_id = $2`, tableName)
	_, err = db.Exec(context.Background(), query, "not successful", transferPayment.RequestId)
	if err != nil {
		return err
	}
	return nil
}
