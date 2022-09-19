package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"telegram-signals-bot/internal/api"
)

func (s *Svc) ValidateBinanceApiKeys(msg string, userID string) {
	var apiKey, apiSecret string
	if strings.Contains(msg, "_") {
		position := strings.LastIndex(msg, "_")
		apiKey = strings.TrimSpace(msg[:position])
		apiSecret = strings.TrimSpace(msg[position+1:])
	} else {
		s.invalidBinanceKeys(userID)
		return
	}

	binanceSvc := api.NewBinanceSvc(apiKey, apiSecret)
	valid := binanceSvc.ValidApiKeys()
	if valid {
		//update db record with keys and reg. timestamp
		err := s.DbSvc.UpdateUserBinanceKeysAndTimestamp(userID, apiKey, apiSecret)
		if err != nil {
			log.Println(err.Error())
			s.sendMsg(userID, err.Error())
			s.invalidBinanceKeys(userID)
			return
		}

		inviteLink, err := s.Bot.GetInviteLink(tgbotapi.ChatConfig{ChatID: s.ChatID})
		if err != nil {
			log.Println(err.Error())
			s.sendMsg(userID, err.Error())
			s.invalidBinanceKeys(userID)
			return
		}

		//send msg
		s.validBinanceKeys(userID, inviteLink)
	} else {
		s.invalidBinanceKeys(userID)
	}
}

func (s *Svc) ValidateLicenseKey(msg, username string, userID string) {
	//check if license key record exists in DB
	user, err := s.DbSvc.GetUserByLicenseKey(msg)
	if err != nil {
		log.Println(err.Error())
	}
	if user != nil && user.RegistrationTimestamp.String == "" {
		err = s.DbSvc.UpdateUserWithIDAndUsernameByLicenseKey(userID, username, msg)
		if err != nil {
			log.Println(err.Error())
			s.sendMsg(userID, err.Error())
			s.insertLicenseKey(userID)
			return
		}
		s.acceptFees(userID)
	} else {
		s.invalidLicenseKey(userID)
	}
}
