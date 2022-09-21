package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"telegram-signals-bot/internal/api"
	"telegram-signals-bot/internal/db"
)

func (s *Svc) remindAboutThePayment() {
	users, err := s.DbSvc.GetUsers()
	if err != nil {
		log.Println(fmt.Sprintf("Error appered while trying to get users in sendPaymentReport(), Error: %s", err.Error()))
		return
	}

	for _, u := range users {
		go s.remindUser(u)
	}
}

func (s *Svc) remindUser(user db.User) {
	report, err := s.DbSvc.GetLastUserReport(user.UserID.String)
	if err != nil {
		log.Println(err)
		return
	}

	if !report.Paid.Bool {
		binanceSvc := api.NewBinanceSvc(user.BinanceApiKey.String, user.BinanceApiSecret.String)
		paymentLink, err := binanceSvc.GetPaymentLink(report.Fees.Float64, user.UserID.String, report.UUID.String)
		if err != nil {
			log.Println(err)
			return
		}

		id, err := strconv.Atoi(user.UserID.String)
		if err != nil {
			log.Println(err)
		}

		message := tgbotapi.NewMessage(int64(id), fmt.Sprintf(PaymentReminderMsg, report.ReportInfo.String))
		message.ReplyMarkup = GenerateNewLinkKeyboard(paymentLink)
		_, err = s.Bot.Send(message)
		if err != nil {
			log.Println(err)
		}
	}
}
