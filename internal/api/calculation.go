package api

import (
	"fmt"
	"github.com/adshao/go-binance/v2"
	uuid "github.com/satori/go.uuid"
	"log"
	"math"
	"strconv"
	"time"
)

type DailyReport struct {
	UUID           string
	UserID         string
	Username       string
	Fees           float64
	ClosedDate     string
	Symbol         string
	BuyPrice       float64
	SellPrice      float64
	BuyCCQ         float64
	Commission     float64
	ProfitWithDust float64
	Profit         float64
	TodayDate      string
}

func (s *BinanceSvc) CalculateProfit(symbol string, startTime, endTime int64) []DailyReport {
	results, err := s.profitCalculationUpdated(symbol, startTime, endTime)
	if err != nil {
		log.Println(err)
		return nil
	}

	return results
}

func (s *BinanceSvc) getPrevDayLastBuyOrderAndTrades(symbol string, startTime int64) ([]*binance.Order, []*binance.TradeV3, error) {
	var order *binance.Order
	var orderTrades []*binance.TradeV3
	endTime := startTime
	startTime = time.Unix(0, startTime*int64(time.Millisecond)).Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)

	var allOrders []*binance.Order
	var err error
	for {
		allOrders, err = s.ListMyOrders(symbol, startTime, endTime)
		if err != nil {
			return nil, nil, err
		}

		for _, o := range allOrders {
			if o.Status == binance.OrderStatusTypeFilled && o.Side == binance.SideTypeBuy {
				order = o
			}
		}

		if order != nil {
			break
		}
		endTime = startTime
		startTime = time.Unix(0, startTime*int64(time.Millisecond)).Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)
	}

	allTrades, err := s.ListMyTrades(symbol, startTime, endTime)
	if err != nil {
		return nil, nil, err
	}

	for _, trade := range allTrades {
		if trade.OrderID == order.OrderID {
			orderTrades = append(orderTrades, trade)
		}
	}

	if order == nil {
		return nil, orderTrades, nil
	}

	return []*binance.Order{order}, orderTrades, nil
}

func (s *BinanceSvc) profitCalculationUpdated(symbol string, startTime, endTime int64) ([]DailyReport, error) {
	var allOrders []*binance.Order
	orders, err := s.ListMyOrders(symbol, startTime, endTime)
	if err != nil {
		return nil, err
	}

	if len(orders) > 0 && orders[len(orders)-1].Side == binance.SideTypeBuy {
		for key, o := range orders {
			if key == len(orders)-1 {
				continue
			}
			allOrders = append(allOrders, o)
		}
	} else {
		allOrders = orders
	}

	allTrades, err := s.ListMyTrades(symbol, startTime, endTime)
	if err != nil {
		return nil, err
	}

	if len(allOrders) > 0 && allOrders[0].Side == binance.SideTypeSell {
		prevDayLastBuy, prevDayBuyTrades, err := s.getPrevDayLastBuyOrderAndTrades(symbol, startTime)
		if err != nil {
			return nil, err
		}
		allTrades = append(allTrades, prevDayBuyTrades...)
		allOrders = append(prevDayLastBuy, allOrders...)
	}

	return matchOrders(allOrders, allTrades), nil
}

func matchOrders(allOrders []*binance.Order, allTrades []*binance.TradeV3) []DailyReport {
	var results []DailyReport
	var buyPrice, buyTradesCount, sellPrice, sellTradesCount, tradesCommission, buyCCQ float64
	var closingTime int64

	location, _ := time.LoadLocation("Europe/Rome")

	for _, o := range allOrders {
		if buyPrice != 0 && sellPrice != 0 && buyCCQ != 0 {
			sellPrice = sellPrice / sellTradesCount
			if buyTradesCount != 0 {
				buyPrice = buyPrice / buyTradesCount
			}
			orderProfitWithDust := (sellPrice-buyPrice)/buyPrice*buyCCQ - tradesCommission
			orderProfit := orderProfitWithDust - (0.01 * orderProfitWithDust)

			result := DailyReport{
				UUID:           uuid.NewV4().String(),
				BuyPrice:       math.Round(buyPrice*10000) / 10000,
				SellPrice:      math.Round(sellPrice*10000) / 10000,
				BuyCCQ:         math.Round(buyCCQ*10000) / 10000,
				Commission:     math.Round(tradesCommission*10000) / 10000,
				ProfitWithDust: math.Round(orderProfitWithDust*10000) / 10000,
				Profit:         math.Round(orderProfit*10000) / 10000,
			}
			if location == nil {
				result.TodayDate = time.Now().Format("2006-01-02")
				result.ClosedDate = time.Unix(0, closingTime*int64(time.Millisecond)).Format(time.RFC3339)
			} else {
				result.TodayDate = time.Now().In(location).Format("2006-01-02")
				result.ClosedDate = time.Unix(0, closingTime*int64(time.Millisecond)).In(location).Format(time.RFC3339)
			}
			results = append(results, result)

			buyPrice = 0
			buyTradesCount = 0
			sellPrice = 0
			sellTradesCount = 0
			buyCCQ = 0
			tradesCommission = 0
		}

		if o.Side == binance.SideTypeBuy {
			if ccq, err := strconv.ParseFloat(o.CummulativeQuoteQuantity, 64); err == nil {
				if p, err := strconv.ParseFloat(o.Price, 64); err == nil {
					buyCCQ = ccq
					buyPrice = p
				}
			}
			for _, t := range allTrades {
				if o.OrderID == t.OrderID {
					if commission, err := strconv.ParseFloat(t.Commission, 64); err == nil {
						tradesCommission += commission
						if buyPrice == 0.0 {
							if p, err := strconv.ParseFloat(t.Price, 64); err == nil {
								buyPrice += p
								buyTradesCount++
							}
						}
					}
				}
			}
		}

		if o.Side == binance.SideTypeSell {
			closingTime = o.UpdateTime
			for _, t := range allTrades {
				if o.OrderID == t.OrderID {
					if commission, err := strconv.ParseFloat(t.Commission, 64); err == nil {
						if p, err := strconv.ParseFloat(t.Price, 64); err == nil {
							sellTradesCount++
							sellPrice += p
							tradesCommission += commission
						}
					}
				}
			}
		}
	}

	if buyPrice != 0 && sellPrice != 0 && buyCCQ != 0 {
		sellPrice = sellPrice / sellTradesCount
		if buyTradesCount != 0 {
			buyPrice = buyPrice / buyTradesCount
		}
		orderProfitWithDust := (sellPrice-buyPrice)/buyPrice*buyCCQ - tradesCommission
		orderProfit := orderProfitWithDust - (0.01 * orderProfitWithDust)
		result := DailyReport{
			UUID:           uuid.NewV4().String(),
			BuyPrice:       math.Round(buyPrice*10000) / 10000,
			SellPrice:      math.Round(sellPrice*10000) / 10000,
			BuyCCQ:         math.Round(buyCCQ*10000) / 10000,
			Commission:     math.Round(tradesCommission*10000) / 10000,
			ProfitWithDust: math.Round(orderProfitWithDust*10000) / 10000,
			Profit:         math.Round(orderProfit*10000) / 10000,
		}
		if location == nil {
			result.TodayDate = time.Now().Format("2006-01-02")
			result.ClosedDate = time.Unix(0, closingTime*int64(time.Millisecond)).Format(time.RFC3339)
		} else {
			result.TodayDate = time.Now().In(location).Format("2006-01-02")
			result.ClosedDate = time.Unix(0, closingTime*int64(time.Millisecond)).In(location).Format(time.RFC3339)
		}
		results = append(results, result)
	}

	return results
}

func (s *BinanceSvc) profitCalculation(symbol string, startTime, endTime int64) (float64, int64, error) {
	allOrders, err := s.ListMyOrders(symbol, startTime, endTime)
	if err != nil {
		return 0, 0, err
	}
	allTrades, err := s.ListMyTrades(symbol, startTime, endTime)
	if err != nil {
		return 0, 0, err
	}

	trackData := map[string]float64{
		"TBOUGHT":    0,
		"TSOLD":      0,
		"TCOMMISION": 0,
		"TLASTBUY":   0,
	}

	arrayOfBuyID := make([]int64, 0)
	arrayOfSellID := make([]int64, 0)
	var lastBuyID int64 = 0

	if len(allOrders) > 0 && allOrders[0].Side == binance.SideTypeSell {
		prevDayLastBuy, prevDayBuyTrades, err := s.getPrevDayLastBuyOrderAndTrades(symbol, startTime)
		if err != nil {
			return 0, 0, err
		}
		allTrades = append(allTrades, prevDayBuyTrades...)
		allOrders = append(prevDayLastBuy, allOrders...)
	}

	for _, o := range allOrders {
		if o.Status == "FILLED" {
			if o.Side == "BUY" {
				arrayOfBuyID = append(arrayOfBuyID, o.OrderID)
			} else if o.Side == "SELL" {
				arrayOfSellID = append(arrayOfSellID, o.OrderID)
			}
		}
	}

	for _, t := range allTrades {
		for _, id := range arrayOfBuyID {
			if id == t.OrderID {
				if id != lastBuyID {
					lastBuyID = id
					trackData["TLASTBUY"] = 0
				}
				tempPrice := t.Price
				tempExecQty := t.Quantity
				tempCmsn := t.Commission
				if s, err := strconv.ParseFloat(tempPrice, 64); err == nil {
					if d, err := strconv.ParseFloat(tempExecQty, 64); err == nil {
						if c, err := strconv.ParseFloat(tempCmsn, 64); err == nil {
							trackData["TBOUGHT"] += s * d
							trackData["TCOMMISION"] += c
							trackData["TLASTBUY"] += s * d
							lastBuyID = id
						}
					}
				}
			}
		}
		for _, id := range arrayOfSellID {
			if id == t.OrderID {
				tempPrice := t.Price
				tempExecQty := t.Quantity
				tempCmsn := t.Commission
				if s, err := strconv.ParseFloat(tempPrice, 64); err == nil {
					if d, err := strconv.ParseFloat(tempExecQty, 64); err == nil {
						if c, err := strconv.ParseFloat(tempCmsn, 64); err == nil {
							trackData["TSOLD"] += s * d
							trackData["TCOMMISION"] += c
						}
					}
				}
			}
		}
	}

	totalBought := trackData["TBOUGHT"]
	totalSold := trackData["TSOLD"]
	totalPnL := totalSold - totalBought

	if math.Mod(totalSold, totalPnL) >= 1000 {
		totalPnL += trackData["TLASTBUY"]
		totalBought -= trackData["TLASTBUY"]
	}

	cornixMethod := ((totalSold - totalBought) / totalBought)
	cornixMethodExCommision := ((totalSold - (totalBought + trackData["TCOMMISION"])) / (totalBought + trackData["TCOMMISION"]))
	cornixMethodExCommisionExDust := (((totalSold - (totalBought + trackData["TCOMMISION"])) / (totalBought + trackData["TCOMMISION"])) * 0.90)
	cornixMethodPercentage := ((cornixMethod * totalBought) - trackData["TCOMMISION"]) * 0.90
	// cornixMethodExCommision := cornixMethod - trackData["TCOMMISION"]

	closedTime := allOrders[len(allOrders)-1].UpdateTime
	if lastBuyID != 0 {
		closedTime = allOrders[len(allOrders)-2].UpdateTime
	}

	fmt.Println("Profit in ")
	fmt.Println(cornixMethod)

	fmt.Println("TBOUGHT: ")
	fmt.Println(totalBought)
	fmt.Println("=========================")
	fmt.Println("TSOLD:")
	fmt.Println(totalSold)
	fmt.Println("=========================")
	fmt.Println("Total Last Buy: ")
	fmt.Println(trackData["TLASTBUY"])
	fmt.Println("=========================")
	fmt.Println("Total PnL:")
	fmt.Println(totalPnL)
	fmt.Println("Total PnL in %:")
	fmt.Println(cornixMethod * 100)
	fmt.Println("=========================")
	fmt.Println("Total PnL in % Ex-Commision:")
	fmt.Println(cornixMethodExCommision * 100)
	fmt.Println("=========================")
	fmt.Println("Total PnL Ex-Commission and Ex Dust:")
	fmt.Println(cornixMethodExCommisionExDust * 100)
	fmt.Println("=========================")
	fmt.Println("Total PnL Ex-Commission:")
	fmt.Println(totalPnL - trackData["TCOMMISION"])
	fmt.Println("=========================")
	fmt.Println("Total PnL Ex-Commission and Ex Dust:")
	fmt.Println(cornixMethodPercentage)
	fmt.Println("=========================")

	fmt.Println("Total Commision Paid: ")
	fmt.Println(trackData["TCOMMISION"])
	fmt.Println("=========================")

	return cornixMethodPercentage, closedTime, nil
}
