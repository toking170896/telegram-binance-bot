package db

import (
	"database/sql"
	"strconv"
	"time"
)

type (
	Report struct {
		ID         int             `db:"id,omitempty"`
		UserID     sql.NullString  `db:"userID,omitempty"`
		Username   sql.NullString  `db:"username,omitempty"`
		Fees       sql.NullFloat64 `db:"fees,omitempty"`
		Timestamp  sql.NullString  `db:"timestamp,omitempty"`
		ReportInfo sql.NullString  `db:"reportInfo,omitempty"`
		Paid       sql.NullBool    `db:"paid,omitempty"`
		UUID       sql.NullString  `db:"uuid,omitempty"`
	}
)

func (s *Svc) InsertReport(userID string, fees float64, username, reportInfo, uuid string) error {
	unix := strconv.Itoa(int(time.Now().UnixNano()))
	_, err := s.Db.Exec("INSERT INTO reports (username, userID, timestamp, reportInfo, fees, uuid) values (?, ?, ?, ?, ?, ?)", username, userID, unix, reportInfo, fees, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svc) GetLastUserReport(userID string) (*Report, error) {
	var report *Report
	rows, err := s.Db.Query("SELECT * FROM reports where userID = ? ORDER BY timestamp DESC LIMIT 1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := Report{}
		err = rows.Scan(&r.ID, &r.UUID, &r.Username, &r.UserID, &r.ReportInfo, &r.Paid, &r.Fees, &r.Timestamp)
		if err != nil {
			return nil, err
		}
		report = &r
		break
	}

	return report, nil
}

func (s *Svc) UpdateReportPaymentStatus(uuid string) error {
	sqlStatement := `UPDATE reports SET paid = ? WHERE uuid = ?;`
	_, err := s.Db.Exec(sqlStatement, true, uuid)
	if err != nil {
		return err
	}
	return nil
}
