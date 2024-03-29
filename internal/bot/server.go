package bot

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

type (
	PaymentWebhook struct {
		BizType   string `json:"bizType"`
		Data      string `json:"data"`
		BizStatus string `json:"bizStatus"`
	}

	WebhookResp struct {
		ReturnCode string `json:"returnCode"`
	}

	Data struct {
		MerchantTradeNo string `json:"merchantTradeNo"`
		PassThroughInfo string `json:"passThroughInfo"`
	}
)

func (s *Svc) StartServer() {
	engine := gin.Default()
	engine.POST("/payment", s.handlePayment)
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Service is running!")
		return
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Starting http server on port :" + port)
	err := engine.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Svc) handlePayment(c *gin.Context) {
	payload := &PaymentWebhook{}
	if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if payload.BizStatus == "PAY_CLOSED" {
		c.JSON(http.StatusOK, WebhookResp{ReturnCode: "SUCCESS"})
		return
	}

	e, err := json.Marshal(payload)
	log.Print(string(e))

	var (
		data     *Data
		reportID string
	)
	err = json.Unmarshal([]byte(payload.Data), &data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if strings.Contains(data.PassThroughInfo, "_") {
		position := strings.LastIndex(data.PassThroughInfo, "_")
		reportID = strings.TrimSpace(data.PassThroughInfo[position+1:])
	}

	go func(reportID string) {
		err = s.DbSvc.UpdateReportPaymentStatus(reportID)
		if err != nil {
			log.Println(err.Error())
			return
		}

		err = s.GoogleCli.UpdateRowWithPaymentDate(reportID)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}(reportID)

	c.JSON(http.StatusOK, WebhookResp{ReturnCode: "SUCCESS"})
	return
}
