package bot

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"strconv"
	"sync"
	"telegram-signals-bot/internal/api"
	"telegram-signals-bot/internal/db"
	"time"
)

func (s *Svc) CreateDailyReports() {
	now := time.Now()
	signalsTimestamp := now.Add(-24*8*time.Hour).UnixNano() / int64(time.Millisecond)
	startTime := now.Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)
	endTime := now.UnixNano() / int64(time.Millisecond)
	//endTime = time.Unix(1664257335, 0).UnixNano()/ int64(time.Millisecond)
	//startTime = time.Unix(1664191980, 0).UnixNano()/ int64(time.Millisecond)

	users, err := s.DbSvc.GetUsers()
	if err != nil {
		log.Println(fmt.Sprintf("Error appeared while trying to get users in DAILY sendPaymentReport(), Error: %s", err.Error()))
		return
	}

	timestamp := signalsTimestamp * int64(time.Millisecond)
	signals, err := s.DbSvc.GetSignalsSince(timestamp)
	if err != nil {
		log.Println(fmt.Sprintf("Error appered while trying to get signals in DAILY sendPaymentReport(), Error: %s", err.Error()))
		return
	}

	var wg = &sync.WaitGroup{}
	for _, u := range users {
		if u.RegistrationTimestamp.String != "" && !u.Blocked.Bool {
			wg.Add(1)
			go s.createReportsForUser(u, signals, startTime, endTime, wg)
		}
	}

	wg.Wait()
}

func (s *Svc) createReportsForUser(user db.User, signals []db.Signal, startTime, endTime int64, wg *sync.WaitGroup) {
	defer wg.Done()

	binanceSvc := api.NewBinanceSvc(user.BinanceApiKey.String, user.BinanceApiSecret.String)
	processedSymbols := make(map[string]string)

	for _, signal := range signals {
		symbol := signal.Symbol.String
		if symbol == "" {
			continue
		}

		if _, found := processedSymbols[symbol]; found {
			continue
		}

		//user trades
		s.processDayTrades(binanceSvc, startTime, endTime, symbol, user)

		//spot trades
		s.processSpotDayTrades(binanceSvc, startTime, endTime, symbol, user)

		processedSymbols[symbol] = symbol
	}
}

func (s *Svc) processDayTrades(binanceSvc *api.BinanceSvc, startTime, endTime int64, symbol string, user db.User) {
	trades, err := binanceSvc.GetUserTrades(symbol, startTime, endTime)
	if err != nil {
		log.Println(fmt.Sprintf("Error: %s, Username: %s", err.Error(), user.Username.String))
	}

	orders := make(map[int]float64)
	for _, t := range trades {
		if _, found := orders[t.OrderID]; !found {
			orders[t.OrderID] = 0
		}
	}

	for orderID, _ := range orders {
		var closedTime int64
		var orderFee, orderRealizedPnl float64
		for _, t := range trades {
			if t.OrderID != orderID {
				continue
			}
			realizedPnl, err := strconv.ParseFloat(t.RealizedPnl, 64)
			if err != nil {
				log.Println(fmt.Sprintf("Error: %s, Username: %s", err.Error(), user.Username.String))
			}
			if realizedPnl != 0 {
				fee := s.addFee(user, realizedPnl)
				orderRealizedPnl += realizedPnl
				orderFee += fee
				closedTime = t.Time
			}
		}

		log.Println(fmt.Sprintf("OrderID: %d, profit: %.2f", orderID, orderFee))
		if orderFee != 0 {
			closedDate := time.Unix(0, closedTime*int64(time.Millisecond)).Format("2006-01-02")
			reportUuid := uuid.NewV4()
			err = s.DbSvc.InsertDailyReport(user.UserID.String, user.Username.String, symbol, reportUuid.String(), closedDate, orderFee, orderRealizedPnl)
			if err != nil {
				log.Println(fmt.Sprintf("Error: %s, Username: %s", err.Error(), user.Username.String))
			}
		}
	}
}

func (s *Svc) processSpotDayTrades(binanceSvc *api.BinanceSvc, startTime, endTime int64, symbol string, user db.User) {
	results := binanceSvc.CalculateProfit(symbol, startTime, endTime)
	for _, r := range results {
		if r.Profit != 0 {
			r.Fees = s.addFee(user, r.Profit)
			err := s.DbSvc.InsertDailyReportForSpot(user.UserID.String, user.Username.String, symbol, r)
			if err != nil {
				log.Println(fmt.Sprintf("Error: %s, Username: %s", err.Error(), user.Username.String))
			}
		}
	}
}
