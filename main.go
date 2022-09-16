package main

import (
	"fmt"
	"telegram-signals-bot/internal/bot"
	"telegram-signals-bot/internal/config"
)

func main() {
	c, err := config.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	svc, err := bot.NewSvc(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	//start cronjobs
	go svc.StartCronJobs()

	////start server
	go svc.StartServer()

	svc.StartBot()
}
