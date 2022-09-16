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

	//start cronjob, which send fees reports to clients once per week
	//go svc.StartCronJobs()

	////start server
	go svc.StartServer()

	svc.StartBot()
}
//ADP1VDs71kkuc0DV9sz4B1Nao37QYiy3eS7Fx3DIM2g3IfjPIALND8VwZO7jmCFP_eUt1nIS6jgsZlABVnykoA0xqKrT9M5E9yZEIKscOoWnNCIrMPC2mqiHwgsxtW2b5
//ADAUSDT
//BTCUSDT
//XRPBUSD