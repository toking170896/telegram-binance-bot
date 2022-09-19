package bot

import (
	"fmt"
	"log"
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

		s.sendMsg(user.UserID.String, fmt.Sprintf(PaymentReminderMsg, paymentLink))
	}
}
