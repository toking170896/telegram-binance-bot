package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"telegram-signals-bot/internal/api"
)

func (s *Svc) acceptTerms(userID string) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %s got such an error: %s", userID, err.Error()))
	}
	message := tgbotapi.NewMessage(int64(id), AcceptTermsMsg)
	message.ReplyMarkup = AcceptTermsKeyboard()

	_, err = s.Bot.Send(message)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %s got such an error: %s", userID, err.Error()))
	}
}

func (s *Svc) termsDeniedByUser(userID string) {
	s.sendMsg(userID, DeniedTermsErr)
	s.acceptTerms(userID)
	s.updateUserState(userID, DenyTermsState)
}

func (s *Svc) insertLicenseKey(userID string) {
	s.sendMsg(userID, InsertLicenseKeyMsg)
	s.updateUserState(userID, InsertLicenseKeyState)
}

func (s *Svc) invalidLicenseKey(userID string) {
	s.sendMsg(userID, InvalidLicenseKeyMsg)
	s.insertLicenseKey(userID)
	s.updateUserState(userID, InvalidLicenseKey)
}

func (s *Svc) acceptFees(userID string) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %s got such an error: %s", userID, err.Error()))
	}
	message := tgbotapi.NewMessage(int64(id), ValidLicenseKeyMsg)
	message.ReplyMarkup = AcceptFeesKeyboard()
	_, err = s.Bot.Send(message)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %s got such an error: %s", userID, err.Error()))
	}

	s.updateUserState(userID, ValidLicenseKey)
}

func (s *Svc) acceptFeesRetry(userID string) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %s got such an error: %s", userID, err.Error()))
	}
	message := tgbotapi.NewMessage(int64(id), AcceptFeesRetryMsg)
	message.ReplyMarkup = AcceptFeesKeyboard()
	_, err = s.Bot.Send(message)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %s got such an error: %s", userID, err.Error()))
	}

	s.updateUserState(userID, ValidLicenseKey)
}

func (s *Svc) feesDeniedByUser(userID string) {
	s.sendMsg(userID, DeniedFeesMsg)
	s.insertBinanceKeys(userID)
}

func (s *Svc) insertBinanceKeys(userID string) {
	s.sendMsg(userID, InsertBinanceKeyMsg)
	s.updateUserState(userID, InsertBinanceKeysState)
}

func (s *Svc) validBinanceKeys(userID string, inviteUrl string) {
	s.sendMsg(userID, fmt.Sprintf(ValidBinanceKeysMsg, inviteUrl))
	s.States.Remove(userID)
}

func (s *Svc) invalidBinanceKeys(userID string) {
	s.sendMsg(userID, InvalidBinanceKeysMsg)
	s.insertBinanceKeys(userID)
	s.updateUserState(userID, InvalidBinanceKeys)
}

func (s *Svc) generatePaymentLink(userID string) {
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

	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %s got such an error: %s", userID, err.Error()))
	}
	message := tgbotapi.NewMessage(int64(id), NewlyGeneratedPaymentLinkMsg)
	message.ReplyMarkup = GenerateNewLinkKeyboard(paymentLink)
	_, err = s.Bot.Send(message)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %s got such an error: %s", userID, err.Error()))
	}
}
