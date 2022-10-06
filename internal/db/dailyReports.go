package db

import (
	"database/sql"
	"math"
	"telegram-signals-bot/internal/api"
	"time"
)

type (
	DailyReport struct {
		ID             int             `db:"id,omitempty"`
		UUID           sql.NullString  `db:"uuid,omitempty"`
		UserID         sql.NullString  `db:"userID,omitempty"`
		Username       sql.NullString  `db:"username,omitempty"`
		Fees           sql.NullFloat64 `db:"fees,omitempty"`
		ClosedDate     sql.NullString  `db:"closedDate,omitempty"`
		Symbol         sql.NullString  `db:"symbol,omitempty"`
		BuyPrice       sql.NullFloat64 `db:"buyPrice,omitempty"`
		SellPrice      sql.NullFloat64 `db:"sellPrice,omitempty"`
		BuyCCQ         sql.NullFloat64 `db:"buyccq,omitempty"`
		Commission     sql.NullFloat64 `db:"commission,omitempty"`
		ProfitWithDust sql.NullFloat64 `db:"profitWithDust,omitempty"`
		Profit         sql.NullFloat64 `db:"profit,omitempty"`
		TodayDate      sql.NullString  `db:"todayDate,omitempty"`
	}
)

func (s *Svc) InsertDailyReport(userID, username, symbol, uuid, closedDate string, fees, profit float64) error {
	fees = math.Round(fees*100) / 100
	profit = math.Round(profit*100) / 100

	location, _ := time.LoadLocation("Europe/Rome")
	todayDate := time.Now().In(location).Format("2006-01-02")
	_, err := s.Db.Exec("INSERT INTO dailyReports (username, userID, closedDate, symbol, fees, profit, todayDate, uuid) values (?, ?, ?, ?, ?, ?, ?, ?)", username, userID, closedDate, symbol, fees, profit, todayDate, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svc) InsertDailyReportForSpot(userID, username, symbol string, report api.DailyReport) error {
	_, err := s.Db.Exec("INSERT INTO dailyReports (username, userID, closedDate, symbol, buyPrice, sellPrice, buyccq, commission, profitWithDust, profit, fees, todayDate, uuid) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		username, userID, report.ClosedDate, symbol, report.BuyPrice, report.SellPrice, report.BuyCCQ, report.Commission, report.ProfitWithDust, report.Profit, report.Fees, report.TodayDate, report.UUID)
	if err != nil {
		return err
	}
	return nil
}
