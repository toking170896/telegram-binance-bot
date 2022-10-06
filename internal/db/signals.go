package db

import (
	"database/sql"
	"strconv"
	"time"
)

type (
	Signal struct {
		ID        int            `db:"id,omitempty"`
		Symbol    sql.NullString `db:"symbol,omitempty"`
		Timestamp sql.NullString `db:"timestamp,omitempty"`
		TradeTime sql.NullString `db:"tradeTime,omitempty"`
	}
)

func (s *Svc) InsertSignal(symbol string) error {
	tradeTime := ""
	unixNano := time.Now().UnixNano()
	location, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		tradeTime = time.Unix(0, unixNano).Format(time.RFC3339)
	} else {
		tradeTime = time.Unix(0, unixNano).In(location).Format(time.RFC3339)
	}

	_, err = s.Db.Exec("INSERT INTO signals (symbol, timestamp, tradeTime) values (?, ?, ?)", symbol, strconv.Itoa(int(unixNano)), tradeTime)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svc) GetSignalsSince(timestamp int64) ([]Signal, error) {
	var signals []Signal
	rows, err := s.Db.Query("SELECT * FROM signals WHERE timestamp >= ?", timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		signal := Signal{}
		err = rows.Scan(&signal.ID, &signal.Symbol, &signal.Timestamp, &signal.TradeTime)
		if err != nil {
			return nil, err
		}
		signals = append(signals, signal)
	}

	return signals, nil
}
