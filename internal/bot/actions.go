package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"telegram-signals-bot/internal/api"
)

func (s *Svc) acceptTerms(userID int64) {
	message := tgbotapi.NewMessage(userID, AcceptTermsMsg)
	message.ReplyMarkup = AcceptTermsKeyboard()

	_, err := s.Bot.Send(message)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %d got such an error: %s", userID, err.Error()))
	}
}

func (s *Svc) termsDeniedByUser(userID int64) {
	s.sendMsg(userID, DeniedTermsErr)
	s.acceptTerms(userID)
	s.updateUserState(userID, DenyTermsState)
}

func (s *Svc) insertLicenseKey(userID int64) {
	s.sendMsg(userID, InsertLicenseKeyMsg)
	s.updateUserState(userID, InsertLicenseKeyState)
}

func (s *Svc) invalidLicenseKey(userID int64) {
	s.sendMsg(userID, InvalidLicenseKeyMsg)
	s.insertLicenseKey(userID)
	s.updateUserState(userID, InvalidLicenseKey)
}

func (s *Svc) acceptFees(userID int64) {
	message := tgbotapi.NewMessage(userID, ValidLicenseKeyMsg)
	message.ReplyMarkup = AcceptFeesKeyboard()
	_, err := s.Bot.Send(message)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %d got such an error: %s", userID, err.Error()))
	}

	s.updateUserState(userID, ValidLicenseKey)
}

func (s *Svc) acceptFeesRetry(userID int64) {
	message := tgbotapi.NewMessage(userID, AcceptFeesRetryMsg)
	message.ReplyMarkup = AcceptFeesKeyboard()
	_, err := s.Bot.Send(message)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %d got such an error: %s", userID, err.Error()))
	}

	s.updateUserState(userID, ValidLicenseKey)
}

func (s *Svc) feesDeniedByUser(userID int64) {
	s.sendMsg(userID, DeniedFeesMsg)
	s.insertBinanceKeys(userID)
}

func (s *Svc) insertBinanceKeys(userID int64) {
	s.sendMsg(userID, InsertBinanceKeyMsg)
	s.updateUserState(userID, InsertBinanceKeysState)
}

func (s *Svc) validBinanceKeys(userID int64, inviteUrl string) {
	s.sendMsg(userID, fmt.Sprintf(ValidBinanceKeysMsg, inviteUrl))
	s.States.Remove(strconv.Itoa(int(userID)))
}

func (s *Svc) invalidBinanceKeys(userID int64) {
	s.sendMsg(userID, InvalidBinanceKeysMsg)
	s.insertBinanceKeys(userID)
	s.updateUserState(userID, InvalidBinanceKeys)
}

func (s *Svc) generatePaymentLink(userID int64) {
	report, err := s.DbSvc.GetLastUserReport(userID)
	if err != nil {
		log.Println(err)
		s.sendMsg(userID, err.Error())
	}

	user, err := s.DbSvc.GetUserByUserID(userID)
	if err != nil {
		log.Println(err)
		s.sendMsg(userID, err.Error())
	}

	binanceSvc := api.NewBinanceSvc(user.BinanceApiKey.String, user.BinanceApiSecret.String)
	paymentLink, err := binanceSvc.GetPaymentLink(report.Fees.Float64, userID, report.UUID.String)
	if err != nil {
		log.Println(err)
		s.sendMsg(userID, err.Error())
	}

	s.sendMsg(userID, fmt.Sprintf(NewlyGeneratedPaymentLinkMsg, paymentLink))
}