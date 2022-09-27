package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"telegram-signals-bot/internal/config"
)

type (
	Svc struct {
		Db *sql.DB
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
		"userID text, " +
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
		"userID text, " +
		"reportInfo text, " +
		"paid boolean, " +
		"fees text, " +
		"timestamp text)")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS dailyReports (" +
		"id int primary key auto_increment, " +
		"uuid text, " +
		"todayDate text, " +
		"username text, " +
		"userID text, " +
		"closedDate text, " +
		"symbol text, " +
		"profit text, " +
		"fees text)")
	if err != nil {
		return nil, err
	}

	return &Svc{db}, nil
}
