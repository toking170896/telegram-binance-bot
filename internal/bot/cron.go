package bot

import (
	"fmt"
	"github.com/adshao/go-binance/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron/v3"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"telegram-signals-bot/internal/api"
	"telegram-signals-bot/internal/db"
	"telegram-signals-bot/internal/sheets"
	"time"
)

func (s *Svc) StartCronJobs() {
	logger := cron.VerbosePrintfLogger(log.New(os.Stdout, "", log.LstdFlags))
	cronJob := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(logger),
			cron.SkipIfStillRunning(logger),
		),
	)

	//every wednesday at 12.00
	cronJob.AddFunc("0 0 12 * * 3", func() {
		s.processUsers()
	})

	//every friday at 12.00
	cronJob.AddFunc("0 0 12 * * 5", func() {
		s.remindAboutThePayment()
	})

	//ping service
	cronJob.AddFunc("@every 15m", func() {
		PingService()
	})
	cronJob.Start()
}

func (s *Svc) processUsers() {
	now := time.Now()
	startTime := now.Add(-7*24*time.Hour).UnixNano() / int64(time.Millisecond)
	endTime := now.UnixNano() / int64(time.Millisecond)

	users, err := s.DbSvc.GetUsers()
	if err != nil {
		log.Println(fmt.Sprintf("Error appered while trying to get users in sendPaymentReport(), Error: %s", err.Error()))
		return
	}

	signals, err := s.DbSvc.GetSignalsSince(startTime)
	if err != nil {
		log.Println(fmt.Sprintf("Error appered while trying to get signals in sendPaymentReport(), Error: %s", err.Error()))
		return
	}

	var (
		wg         = &sync.WaitGroup{}
		reports    []*sheets.Report
		reportChan = make(chan *sheets.Report)
		quit       = make(chan bool)
	)

	for _, u := range users {
		if u.RegistrationTimestamp.String != "" && !u.Blocked.Bool {
			wg.Add(1)
			go s.processUser(u, signals, startTime, endTime, wg, reportChan)
		}
	}

	go func() {
		for {
			select {
			case <-quit:
				return
			case report := <-reportChan:
				reports = append(reports, report)
			}
		}
	}()

	wg.Wait()
	quit <- true
	close(reportChan)
	close(quit)

	if len(reports) == 0 {
		return
	}

	err = s.GoogleCli.InsertNewRows(reports)
	if err != nil {
		log.Println(fmt.Sprintf("Error appeared while trying to insert new rows to sheets, Error: %s", err.Error()))
	}
}

func (s *Svc) processUser(user db.User, signals []db.Signal, startTime, endTime int64, wg *sync.WaitGroup, reportChan chan *sheets.Report) {
	defer wg.Done()

	binanceSvc := api.NewBinanceSvc(user.BinanceApiKey.String, user.BinanceApiSecret.String)
	processedSymbols := make(map[string]string)

	var feeSum float64
	var totalTrades int
	var report, notFilledOrderIDs string
	for _, signal := range signals {
		symbol := signal.Symbol.String
		if symbol == "" {
			continue
		}

		if _, found := processedSymbols[symbol]; found {
			continue
		}

		//user trades
		fee, tradesNum, reportLine := s.processTrades(binanceSvc, startTime, endTime, symbol, user)
		totalTrades += tradesNum
		feeSum += fee
		report += reportLine

		//spot trades
		fee, tradesNum, reportLine, ids := s.processSpotTrades(binanceSvc, startTime, endTime, symbol, user)
		totalTrades += tradesNum
		feeSum += fee
		report += reportLine
		notFilledOrderIDs += ids

		processedSymbols[symbol] = symbol
	}

	err := s.DbSvc.UpdateUserTimestamp(user.UserID.String, endTime)
	if err != nil {
		log.Println(err)
	}

	if totalTrades > 0 && feeSum >= 0.01 {
		feeSum = math.Round(feeSum*100) / 100
		reportUuid := uuid.NewV4()
		paymentLink, err := binanceSvc.GetPaymentLink(feeSum, user.UserID.String, reportUuid.String())
		if err != nil {
			log.Println(err)
		}

		report = s.addStartAndEndToReport(report, paymentLink, totalTrades, feeSum)
		if notFilledOrderIDs != "" {
			err = s.DbSvc.UpdateUserNotFilledOrderIDs(user.UserID.String, notFilledOrderIDs)
			if err != nil {
				log.Println(err)
			}
		}

		err = s.DbSvc.InsertReport(user.UserID.String, feeSum, user.Username.String, report, reportUuid.String())
		if err != nil {
			log.Println(err)
		}

		reportObj := &sheets.Report{
			UserID:     user.UserID.String,
			Username:   user.Username.String,
			Fees:       feeSum,
			ReportPaid: "false",
			UUID:       reportUuid.String(),
		}
		reportChan <- reportObj

		id, err := strconv.Atoi(user.UserID.String)
		if err != nil {
			log.Println(err)
		}

		message := tgbotapi.NewMessage(int64(id), report)
		message.ReplyMarkup = GenerateNewLinkKeyboard()
		_, err = s.Bot.Send(message)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Svc) processTrades(binanceSvc *api.BinanceSvc, startTime, endTime int64, symbol string, user db.User) (float64, int, string) {
	var feeSum float64
	var amountOfTrades int
	var report string

	trades, err := binanceSvc.GetUserTrades(symbol, startTime, endTime)
	if err != nil {
		log.Println(err)
	}

	for _, t := range trades {
		realizedPnl, err := strconv.ParseFloat(t.RealizedPnl, 64)
		if err != nil {
			log.Println(err)
		}
		if realizedPnl > 0 {
			amountOfTrades++
			fee := s.addFee(user, realizedPnl)
			feeSum += fee
			report += s.addReportLine(realizedPnl, fee, symbol, t.Time)
		}
	}

	return feeSum, amountOfTrades, report
}

func (s *Svc) processSpotTrades(binanceSvc *api.BinanceSvc, startTime, endTime int64, symbol string, user db.User) (float64, int, string, string) {
	var notFilledOrderIDs, report string
	var totalFee float64
	var amountOfTrades int

	var prevSyncOrders []*binance.TradeV3
	if user.NotFilledOrderIDs.String != "" {
		ids := strings.Split(user.NotFilledOrderIDs.String, ",")
		for _, i := range ids {
			id, err := strconv.Atoi(i)
			if err != nil {
				log.Print(err)
				continue
			}
			prevSyncOrders = append(prevSyncOrders, &binance.TradeV3{OrderID: int64(id)})
		}
	}

	allTrades, err := binanceSvc.ListAllMyTradesForTheWeek(symbol, startTime, endTime)
	if err != nil {
		log.Println(err)
	}

	allTrades = append(allTrades, prevSyncOrders...)

	var orders []*binance.Order
	for _, trade := range allTrades {
		order, err := binanceSvc.GetOrderByID(trade.OrderID, symbol)
		if err != nil {
			log.Println(err)
		}
		if order != nil {
			if order.Status == binance.OrderStatusTypeFilled {
				orders = append(orders, order)
			} else if order.Status == binance.OrderStatusTypeNew || order.Status == binance.OrderStatusTypePartiallyFilled {
				notFilledOrderIDs += fmt.Sprintf("%s_%d,", symbol, int(order.OrderID))
			}
		}
	}

	for _, o := range orders {
		symbolPrice, err := binanceSvc.GetSymbolPrice(o.Symbol)
		if err != nil {
			log.Print(err)
			continue
		}

		if symbolPrice == nil {
			continue
		}
		currentPrice, err := strconv.ParseFloat(symbolPrice.Price, 64)
		if err != nil {
			log.Print(err)
			continue
		}

		executedQty, err := strconv.ParseFloat(o.ExecutedQuantity, 64)
		if err != nil {
			log.Print(err)
			continue
		}

		cumulativeQuoteQty, err := strconv.ParseFloat(o.CummulativeQuoteQuantity, 64)
		if err != nil {
			log.Print(err)
			continue
		}

		profit := currentPrice*executedQty - cumulativeQuoteQty
		if profit > 0 {
			amountOfTrades++
			fee := s.addFee(user, profit)
			totalFee += fee
			report += s.addReportLine(profit, fee, symbol, o.UpdateTime)
		}
	}

	return totalFee, amountOfTrades, report, notFilledOrderIDs
}

func (s *Svc) addFee(user db.User, realizedPnl float64) float64 {
	var from0To100, from100To500, from500To1000, from1000To5000, result float64
	from0To100 = 20
	from100To500 = 10
	from500To1000 = 5
	from1000To5000 = 4

	if user.FeesPercentage.String != "" {
		percentages := strings.Split(strings.TrimSpace(user.FeesPercentage.String), ",")
		for position, percent := range percentages {
			value, err := strconv.ParseFloat(percent, 64)
			if err != nil {
				log.Println(err)
				break
			}

			switch position {
			case 0:
				from0To100 = value
			case 1:
				from100To500 = value
			case 2:
				from500To1000 = value
			case 3:
				from1000To5000 = value
			}
		}
	}

	switch {
	case realizedPnl < 100:
		result = realizedPnl * from0To100 / 100
	case realizedPnl > 100 && realizedPnl < 500:
		result = realizedPnl * from100To500 / 100
	case realizedPnl > 500 && realizedPnl < 1000:
		result = realizedPnl * from500To1000 / 100
	case realizedPnl > 1000 && realizedPnl < 5000:
		result = realizedPnl * from1000To5000 / 100
	}

	return result
}

func (s *Svc) addReportLine(realizedPnl, fee float64, symbol string, closedDate int64) string {
	date := time.Unix(0, closedDate*int64(time.Millisecond))
	return fmt.Sprintf(FeeLineMsgStructure, symbol, date.Format("2006-01-02"), realizedPnl, fee)
}

func (s *Svc) addStartAndEndToReport(reportLines, paymentLink string, amountOfTrades int, feeSum float64) string {
	fromDate := time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	reportStart := fmt.Sprintf(ReportStartMsg, fromDate, toDate)
	reportEnd := fmt.Sprintf(ReportEndMsg, amountOfTrades, feeSum, paymentLink)

	return reportStart + reportLines + reportEnd
}

func PingService() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://prx-bot.herokuapp.com", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(string(bodyBytes))
		fmt.Println(err)
	}
	return
}
