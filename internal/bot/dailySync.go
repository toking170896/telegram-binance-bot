package bot

import (
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/common"
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
	startTime := now.Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)
	endTime := now.UnixNano() / int64(time.Millisecond)

	users, err := s.DbSvc.GetUsers()
	if err != nil {
		log.Println(fmt.Sprintf("Error appeared while trying to get users in DAILY sendPaymentReport(), Error: %s", err.Error()))
		return
	}

	timestamp := startTime * int64(time.Millisecond)
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
			if realizedPnl > 0 {
				fee := s.addFee(user, realizedPnl)
				orderRealizedPnl += realizedPnl
				orderFee += fee
				closedTime = t.Time
			}
		}

		if orderFee > 0 {
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
	allTrades, err := binanceSvc.ListMyTrades(symbol, startTime, endTime)
	if err != nil {
		if common.IsAPIError(err) {
			apiErr := err.(*common.APIError)
			if apiErr.Code == -1003 {
				log.Println("Reached binance api limit, cool down for 60 sec")
				time.Sleep(1 * time.Minute)
				allTrades, _ = binanceSvc.ListMyTrades(symbol, startTime, endTime)
			}
		}
		log.Println(fmt.Sprintf("Error: %s, Username: %s", err.Error(), user.Username.String))
	}

	orders := make(map[int64]*binance.Order)
	for _, trade := range allTrades {
		if _, found := orders[trade.OrderID]; found {
			continue
		}
		order, err := binanceSvc.GetOrderByID(trade.OrderID, symbol)
		if err != nil {
			log.Println(fmt.Sprintf("Error: %s, Username: %s", err.Error(), user.Username.String))
		}
		if order != nil {
			if order.Status == binance.OrderStatusTypeFilled {
				orders[trade.OrderID] = order
			}
		}
	}

	for _, o := range orders {
		symbolPrice, err := binanceSvc.GetSymbolPrice(o.Symbol)
		if err != nil {
			log.Println(fmt.Sprintf("Error: %s, Username: %s", err.Error(), user.Username.String))
			continue
		}

		if symbolPrice == nil {
			continue
		}
		currentPrice, err := strconv.ParseFloat(symbolPrice.Price, 64)
		if err != nil {
			log.Println(fmt.Sprintf("Error: %s, Username: %s", err.Error(), user.Username.String))
			continue
		}

		executedQty, err := strconv.ParseFloat(o.ExecutedQuantity, 64)
		if err != nil {
			log.Println(fmt.Sprintf("Error: %s, Username: %s", err.Error(), user.Username.String))
			continue
		}

		cumulativeQuoteQty, err := strconv.ParseFloat(o.CummulativeQuoteQuantity, 64)
		if err != nil {
			log.Println(fmt.Sprintf("Error: %s, Username: %s", err.Error(), user.Username.String))
			continue
		}

		profit := currentPrice*executedQty - cumulativeQuoteQty
		if profit > 0 {
			fee := s.addFee(user, profit)
			closedDate := time.Unix(0, o.UpdateTime*int64(time.Millisecond)).Format("2006-01-02")
			reportUuid := uuid.NewV4()
			err = s.DbSvc.InsertDailyReport(user.UserID.String, user.Username.String, symbol, reportUuid.String(), closedDate, fee, profit)
			if err != nil {
				log.Println(fmt.Sprintf("Error: %s, Username: %s", err.Error(), user.Username.String))
			}
		}
	}
}
