package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"strings"
)

func (s *Svc) StartBot() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := s.Bot.GetUpdatesChan(u)
	if err != nil {
		return
	}
	for update := range updates {
		go s.processUpdate(update)
	}
}

func (s *Svc) processUpdate(update tgbotapi.Update) {
	if update.ChannelPost != nil && update.ChannelPost.Chat != nil && update.ChannelPost.Chat.ID == s.ChatID {
		s.processSignal(update.ChannelPost.Text)
		return
	}

	//update is a msg
	if update.Message != nil {
		s.processMsg(update)
		return
	}

	if update.CallbackQuery != nil {
		s.processCallbackQuery(update)
	}
}

func (s *Svc) processMsg(update tgbotapi.Update) {
	if update.Message.Chat == nil {
		return
	}

	msg := update.Message.Text
	if update.Message.Chat.IsChannel() || update.Message.Chat.IsGroup() {
		if update.Message.Chat.ID == s.ChatID {
			s.processSignal(msg)
			return
		}
	}

	userID := strconv.Itoa(int(update.Message.Chat.ID))
	//if username is empty ask user to set it up and repeat
	if update.Message.From.UserName == "" {
		s.sendMsg(userID, EmptyUsernameErr)
		return
	}

	//Check if user already registered
	user, err := s.DbSvc.GetUserByUserID(userID)
	if err != nil {
		log.Println(err.Error())
	}
	if user != nil && user.RegistrationTimestamp.String != "" {
		s.sendMsg(userID, fmt.Sprintf("You are already registered with userName: %s", user.Username.String))
		return
	}

	state := s.getUserState(userID)
	switch state {
	case DenyTermsState:
		//if user denied terms and tries to send a msg
		s.termsDeniedByUser(userID)

	case InsertLicenseKeyState, InvalidLicenseKey:
		s.ValidateLicenseKey(msg, update.Message.From.UserName, userID)

	case ValidLicenseKey:
		s.acceptFeesRetry(userID)

	case InsertBinanceKeysState, InvalidBinanceKeys:
		s.ValidateBinanceApiKeys(msg, userID)

	default:
		s.acceptTerms(userID)
	}
}

func (s *Svc) processCallbackQuery(update tgbotapi.Update) {
	userID := strconv.Itoa(update.CallbackQuery.From.ID)

	switch update.CallbackQuery.Data {
	case AcceptTermsState:
		//if terms are accepted, ask for license
		s.insertLicenseKey(userID)
	case DenyTermsState:
		// if terms are denied, show msg and ask again
		s.termsDeniedByUser(userID)
	case AcceptFeesState:
		err := s.DbSvc.UpdateUserWithRulesAcceptedByUserID(userID, true)
		if err != nil {
			log.Println(err.Error())
			s.sendMsg(userID, err.Error())
			s.acceptFees(userID)
			return
		}
		s.insertBinanceKeys(userID)
	case DenyFeesState:
		s.acceptFeesRetry(userID)
		//err := s.DbSvc.UpdateUserWithRulesAcceptedByUserID(userID, false)
		//if err != nil {
		//	log.Println(err.Error())
		//	s.sendMsg(userID, err.Error())
		//	s.acceptFees(userID)
		//	return
		//}
		//s.feesDeniedByUser(userID)
	case GeneratePaymentLinkState:
		s.generatePaymentLink(userID)
	}
}

func (s *Svc) processSignal(msg string) {
	var prevLine, symbol string
	var symbolLine bool

	log.Println("Caught a signal")
	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		if symbolLine {
			symbol = strings.TrimSpace(line)
			break
		}
		if strings.Contains(line, "*") && strings.Contains(prevLine, "Trading Alert") {
			symbolLine = true
		}
		prevLine = line
	}

	if symbol != "" {
		log.Println(fmt.Sprintf("Adding signal for %s", symbol))
		err := s.DbSvc.InsertSignal(symbol)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (s *Svc) sendMsg(userID, text string) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %s got such an error: %s", userID, err.Error()))
	}
	message := tgbotapi.NewMessage(int64(id), text)
	_, err = s.Bot.Send(message)
	if err != nil {
		log.Print(fmt.Sprintf("UserID %s got such an error: %s", userID, err.Error()))
	}
}
