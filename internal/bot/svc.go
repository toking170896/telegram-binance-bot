package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	cmap "github.com/orcaman/concurrent-map"
	"log"
	"telegram-signals-bot/internal/config"
	db2 "telegram-signals-bot/internal/db"
	"telegram-signals-bot/internal/sheets"
)

type Svc struct {
	Bot    *tgbotapi.BotAPI
	States cmap.ConcurrentMap
	DbSvc  *db2.Svc
	ChatID int64
	GoogleCli *sheets.GoogleCli
}

func NewSvc(c *config.Config) (*Svc, error) {
	bot, err := tgbotapi.NewBotAPI(c.Token)
	if err != nil {
		log.Println(err)
	}

	db, err := db2.OpenDatabase(c)
	if err != nil {
		return nil, err
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	googleCli, err := sheets.NewGoogleClient()
	if err != nil {
		return nil, err
	}

	s := &Svc{
		Bot: bot,
		DbSvc: db,
		States: cmap.New(),
		ChatID: c.ChatID,
		GoogleCli: googleCli,
	}

	return s, nil
}
