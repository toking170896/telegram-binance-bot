package api

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type (
	BinanceSvc struct {
		ApiKey    string
		ApiSecret string
		BaseUrl   string
		Cli       *binance.Client
		httpCli   *http.Client
	}

	Trade struct {
		Buyer           bool   `json:"buyer"`
		Commission      string `json:"commission"`
		CommissionAsset string `json:"commissionAsset"`
		ID              int    `json:"id"`
		Maker           bool   `json:"maker"`
		OrderID         int    `json:"orderId"`
		Price           string `json:"price"`
		Qty             string `json:"qty"`
		QuoteQty        string `json:"quoteQty"`
		RealizedPnl     string `json:"realizedPnl"`
		Side            string `json:"side"`
		PositionSide    string `json:"positionSide"`
		Symbol          string `json:"symbol"`
		Time            int64  `json:"time"`
	}

	CreateOrderReq struct {
		PassThroughInfo string `json:"passThroughInfo"`
		WebhookUrl      string `json:"webhookUrl"`
		Env             struct {
			TerminalType string `json:"terminalType"`
		} `json:"env"`
		MerchantTradeNo string  `json:"merchantTradeNo"`
		OrderAmount     float64 `json:"orderAmount"`
		Currency        string  `json:"currency"`
		Goods           struct {
			GoodsType        string `json:"goodsType"`
			GoodsCategory    string `json:"goodsCategory"`
			ReferenceGoodsID string `json:"referenceGoodsId"`
			GoodsName        string `json:"goodsName"`
			GoodsDetail      string `json:"goodsDetail"`
		} `json:"goods"`
	}

	CreateOrderRes struct {
		Status string `json:"status"`
		Code   string `json:"code"`
		Data   struct {
			PrepayID     string `json:"prepayId"`
			TerminalType string `json:"terminalType"`
			ExpireTime   int64  `json:"expireTime"`
			QrcodeLink   string `json:"qrcodeLink"`
			QrContent    string `json:"qrContent"`
			CheckoutURL  string `json:"checkoutUrl"`
			Deeplink     string `json:"deeplink"`
			UniversalURL string `json:"universalUrl"`
		} `json:"data"`
	}
)

const (
	MerchantApiKey    = "1q9xg3o42s0qu1gotxkr290za1yqzj4xah4usueic5sfmh8nqm7crgzakify8dzi"
	MerchantApiSecret = "tu6bzs9bomvbdzcmvlpwgx13yzqmyiefrh9xwjfz62x6odrl8mr1pl8ck6r6xrix"
	BinanceBaseUrl    = "https://fapi.binance.com"
)

func NewBinanceSvc(apiKey, apiSecret string) *BinanceSvc {
	cli := binance.NewClient(apiKey, apiSecret)
	httpCli := &http.Client{}

	return &BinanceSvc{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		BaseUrl:   BinanceBaseUrl,
		Cli:       cli,
		httpCli:   httpCli,
	}
}

func (s *BinanceSvc) ListAllMyTradesForTheWeek(symbol string, startTime, endTime int64) ([]*binance.TradeV3, error) {
	var allTrades []*binance.TradeV3

	dayEnd := time.Unix(0, startTime*int64(time.Millisecond)).Add(24*time.Hour).UnixNano() / int64(time.Millisecond)
	day := 1
	for {
		if day == 8 {
			break
		}
		trades, err := s.ListMyTrades(symbol, startTime, dayEnd)
		if err != nil {
			return nil, err
		}
		allTrades = append(allTrades, trades...)
		day++
		startTime = dayEnd
		dayEnd = time.Unix(0, startTime*int64(time.Millisecond)).Add(24*time.Hour).UnixNano() / int64(time.Millisecond)
	}

	return allTrades, nil
}

func (s *BinanceSvc) GetSymbolPrice(symbol string) (*binance.SymbolPrice, error) {
	price, err := s.Cli.NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return nil, err
	}

	if len(price) == 0 {
		return nil, nil
	}

	if price[0] != nil {
		return price[0], nil
	}
	return nil, nil
}

func (s *BinanceSvc) ListMyTrades(symbol string, startTime, endTime int64) ([]*binance.TradeV3, error) {
	trades, err := s.Cli.NewListTradesService().Symbol(symbol).StartTime(startTime).EndTime(endTime).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return trades, nil
}

func (s *BinanceSvc) GetOrderByID(id int64, symbol string) (*binance.Order, error) {
	order, err := s.Cli.NewGetOrderService().OrderID(id).Symbol(symbol).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *BinanceSvc) ValidApiKeys() bool {
	_, err := s.Cli.NewListTradesService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}

func (s *BinanceSvc) GetUserTrades(symbol string, startTime, endTime int64) ([]Trade, error) {
	path := s.BaseUrl + "/fapi/v1/userTrades"

	q := url.Values{}
	q.Add("symbol", symbol)
	q.Add("startTime", strconv.Itoa(int(startTime)))
	q.Add("endTime", strconv.Itoa(int(endTime)))
	q.Add("timestamp", getTimestamp())
	q.Add("signature", signature(q.Encode(), s.ApiSecret))

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("X-MBX-APIKEY", s.ApiKey)

	resp, err := s.httpCli.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res []Trade
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (s *BinanceSvc) GetPaymentLink(amount float64, userID int64, reportUuid string) (string, error) {
	fromDate := time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")
	path := "https://bpay.binanceapi.com/binancepay/openapi/v2/order"
	timestamp := getTimestamp()
	nonce := RandStringBytes(32)

	body := &CreateOrderReq{
		Env: struct {
			TerminalType string `json:"terminalType"`
		}{},
		MerchantTradeNo: "",
		OrderAmount:     0,
		Currency:        "",
		Goods: struct {
			GoodsType        string `json:"goodsType"`
			GoodsCategory    string `json:"goodsCategory"`
			ReferenceGoodsID string `json:"referenceGoodsId"`
			GoodsName        string `json:"goodsName"`
			GoodsDetail      string `json:"goodsDetail"`
		}{},
	}

	tradeNo := uuid.NewV4()
	body.Env.TerminalType = "WEB"
	body.MerchantTradeNo = hex.EncodeToString(tradeNo.Bytes())
	body.Currency = "BUSD"
	body.OrderAmount = amount
	body.Goods.GoodsType = "02"
	body.Goods.GoodsCategory = "0000"
	body.Goods.ReferenceGoodsID = tradeNo.String()
	body.Goods.GoodsName = "Trading Fees"
	body.Goods.GoodsDetail = fmt.Sprintf("Trading Fees %s to %s", fromDate, toDate)
	body.WebhookUrl = "https://p.grz.media/payment"
	body.PassThroughInfo = fmt.Sprintf("%d_%s", userID, reportUuid)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Add("BinancePay-Timestamp", timestamp)
	req.Header.Add("BinancePay-Nonce", nonce)
	req.Header.Add("BinancePay-Certificate-SN", MerchantApiKey)
	req.Header.Add("BinancePay-Signature", binancePaySignature(timestamp, nonce, string(jsonBody), MerchantApiSecret))
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.httpCli.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res *CreateOrderRes
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return "", err
	}

	if res == nil {
		return "", nil
	}

	return res.Data.CheckoutURL, nil
}

func getTimestamp() string {
	return strconv.Itoa(int(time.Now().UnixNano() / int64(time.Millisecond)))
}

func signature(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	signingKey := fmt.Sprintf("%x", mac.Sum(nil))
	return signingKey
}

func binancePaySignature(timestamp, nonce, body, secret string) string {
	payload := timestamp + "\n" + nonce + "\n" + body + "\n"
	s := hmac.New(sha512.New, []byte(secret))
	s.Write([]byte(payload))

	return strings.ToUpper(fmt.Sprintf("%x", s.Sum(nil)))
}

func RandStringBytes(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
