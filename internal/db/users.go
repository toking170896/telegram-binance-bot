package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

type (
	User struct {
		ID                               int            `db:"id,omitempty"`
		Username                         sql.NullString `db:"username,omitempty"`
		UserID                           sql.NullString `db:"userID,omitempty"`
		RulesAccepted                    sql.NullBool   `db:"rulesAccepted,omitempty"`
		Blocked                          sql.NullBool   `db:"blocked,omitempty"`
		LicenseKey                       sql.NullString `db:"licenseKey,omitempty"`
		BinanceApiKey                    sql.NullString `db:"binanceApiKey,omitempty"`
		BinanceApiSecret                 sql.NullString `db:"binanceApiSecret,omitempty"`
		FeesPercentage                   sql.NullString `db:"feesPercentage,omitempty"`
		NotFilledOrderIDs                sql.NullString `db:"notFilledOrderIDs,omitempty"`
		RegistrationTimestamp            sql.NullString `db:"registrationTimestamp,omitempty"`
		LastPaymentLinkReceivedTimestamp sql.NullString `db:"lastPaymentLinkReceivedTimestamp,omitempty"`
	}
)

func (s *Svc) GetUserByLicenseKey(licenseKey string) (*User, error) {
	var user *User
	rows, err := s.Db.Query("SELECT * FROM users where licenseKey = ?", licenseKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.ID, &u.Username, &u.UserID, &u.RulesAccepted, &u.Blocked, &u.LicenseKey, &u.BinanceApiKey,
			&u.BinanceApiSecret, &u.FeesPercentage, &u.NotFilledOrderIDs, &u.LastPaymentLinkReceivedTimestamp, &u.RegistrationTimestamp)
		if err != nil {
			return nil, err
		}
		user = &u
		break
	}

	return user, nil
}

func (s *Svc) GetUserByUserID(userID string) (*User, error) {
	var user *User
	rows, err := s.Db.Query("SELECT * FROM users where userID = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.ID, &u.Username, &u.UserID, &u.RulesAccepted, &u.Blocked, &u.LicenseKey, &u.BinanceApiKey,
			&u.BinanceApiSecret, &u.FeesPercentage, &u.NotFilledOrderIDs, &u.LastPaymentLinkReceivedTimestamp, &u.RegistrationTimestamp)
		if err != nil {
			return nil, err
		}
		user = &u
		break
	}

	return user, nil
}

func (s *Svc) GetUsers() ([]User, error) {
	var users []User
	rows, err := s.Db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.ID, &u.Username, &u.UserID, &u.RulesAccepted, &u.Blocked, &u.LicenseKey, &u.BinanceApiKey,
			&u.BinanceApiSecret, &u.FeesPercentage, &u.NotFilledOrderIDs, &u.LastPaymentLinkReceivedTimestamp, &u.RegistrationTimestamp)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (s *Svc) UpdateUserWithIDAndUsernameByLicenseKey(userID string, userName, licenseKey string) error {
	sqlStatement := `UPDATE users SET userID = ?, username = ? WHERE licenseKey = ?;`
	_, err := s.Db.Exec(sqlStatement, userID, userName, licenseKey)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svc) UpdateUserWithRulesAcceptedByUserID(userID string, rulesAccepted bool) error {
	sqlStatement := `UPDATE users SET rulesAccepted = ? WHERE userID = ?;`
	_, err := s.Db.Exec(sqlStatement, rulesAccepted, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svc) UpdateUserBinanceKeysAndTimestamp(userID string, apiKey, apiSecret string) error {
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))
	sqlStatement := `UPDATE users SET binanceApiKey = ?, binanceApiSecret = ?, registrationTimestamp = ? WHERE userID = ?;`
	_, err := s.Db.Exec(sqlStatement, apiKey, apiSecret, timestamp, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svc) UpdateUserNotFilledOrderIDs(userID string, orderIDs string) error {
	sqlStatement := `UPDATE users SET notFilledOrderIDs = ? WHERE userID = ?;`
	_, err := s.Db.Exec(sqlStatement, orderIDs, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svc) UpdateUserTimestamp(userID string, unix int64) error {
	timestamp := strconv.Itoa(int(unix))
	sqlStatement := `UPDATE users SET lastPaymentLinkReceivedTimestamp = ? WHERE userID = ?;`
	_, err := s.Db.Exec(sqlStatement, timestamp, userID)
	if err != nil {
		return err
	}
	return nil
}
