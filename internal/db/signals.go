package db

import (
	"database/sql"
	"strconv"
	"time"
)

type(
	Signal struct {
		ID               int    `db:"id,omitempty"`
		Symbol           sql.NullString `db:"symbol,omitempty"`
		Timestamp        sql.NullString  `db:"timestamp,omitempty"`
	}
)

func (s *Svc) InsertSignal(symbol string) error {
	unix := strconv.Itoa(int(time.Now().UnixNano()))
	_, err := s.Db.Exec("INSERT INTO signals (symbol, timestamp) values (?, ?)", symbol, unix)
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
		err = rows.Scan(&signal.ID, &signal.Symbol, &signal.Timestamp)
		if err != nil {
			return nil, err
		}
		signals = append(signals, signal)
	}

	return signals, nil
}
