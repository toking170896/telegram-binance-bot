package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"telegram-signals-bot/internal/config"
)

type (
	Svc struct {
		Db     *sql.DB
	}
)

func OpenDatabase(c *config.Config) (*Svc, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", c.Username, c.Password, c.Host, c.Db))
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (" +
		"id int primary key auto_increment, " +
		"username text, " +
		"userID int, " +
		"rulesAccepted boolean, " +
		"blocked boolean, " +
		"licenseKey text, " +
		"binanceApiKey text, " +
		"binanceApiSecret text, " +
		"feesPercentage text, " +
		"notFilledOrderIDs text, " +
		"lastPaymentLinkReceivedTimestamp text, " +
		"registrationTimestamp text)")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS signals (" +
		"id int primary key auto_increment, " +
		"symbol text, " +
		"timestamp text)")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS reports (" +
		"id int primary key auto_increment, " +
		"uuid text, " +
		"username text, " +
		"userID int, " +
		"reportInfo text, " +
		"paid boolean, " +
		"fees text, " +
		"timestamp text)")
	if err != nil {
		return nil, err
	}

	return &Svc{db}, nil
}

//func AddWallet(db *sql.DB, walletAddress string, userID, chatID int, transactionID string) error {
//	unix := time.Now().Unix()
//	_, err := db.Exec("INSERT INTO wallets (walletAddress, userID, chatID, transactionID, unixtimestamp) values (?, ?, ?, ?, ?)", walletAddress, userID, chatID, transactionID, unix)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func UpdateWalletWithTransaction(db *sql.DB, walletAddress, transactionID string) error {
//	unix := time.Now().Unix()
//	sqlStatement := `UPDATE wallets SET transactionID = ?, unixtimestamp = ? WHERE walletAddress = ?;`
//	_, err := db.Exec(sqlStatement, transactionID, unix, walletAddress)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func RemoveWallet(db *sql.DB, walletAddress string, userID int) error {
//	_, err := db.Exec("DELETE FROM wallets WHERE walletAddress = ? AND userID = ?", walletAddress, userID)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func GetUserWalletsAddresses(db *sql.DB, userID int) ([]string, error) {
//	var wallets []string
//	var wallet string
//	rows, err := db.Query("SELECT walletAddress FROM wallets where userID = ?", userID)
//	//rows, err := db.Query("SELECT walletAddress FROM wallets WHERE userID = $1;", userID)
//	defer rows.Close()
//
//	for rows.Next() {
//		err := rows.Scan(&wallet)
//		if err != nil {
//			return nil, err
//		}
//		wallets = append(wallets, wallet)
//	}
//
//	if err != nil {
//		return nil, err
//	}
//	return wallets, nil
//}
//
//func GetAllWalletsAddresses(db *sql.DB) ([]string, error) {
//	var wallets []string
//	var wallet string
//	rows, err := db.Query("SELECT walletAddress FROM wallets")
//
//	defer rows.Close()
//
//	for rows.Next() {
//		err := rows.Scan(&wallet)
//		if err != nil {
//			return nil, err
//		}
//		wallets = append(wallets, wallet)
//	}
//
//	if err != nil {
//		return nil, err
//	}
//	return wallets, nil
//}
//
//func GetAllWallets(db *sql.DB) ([]DbWallet, error) {
//	var wallets []DbWallet
//	var wallet DbWallet
//	rows, err := db.Query("SELECT * FROM wallets")
//
//	defer rows.Close()
//
//	for rows.Next() {
//		err := rows.Scan(&wallet.ID, &wallet.Address, &wallet.UserID, &wallet.ChatID, &wallet.TransactionID, &wallet.UnixTimestamp)
//		if err != nil {
//			return nil, err
//		}
//		wallets = append(wallets, wallet)
//	}
//
//	if err != nil {
//		return nil, err
//	}
//	return wallets, nil
//}
