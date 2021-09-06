package main

import (
	"fmt"
	"milliard-easy/daaz_simbank/context"
	"milliard-easy/daaz_simbank/daaz"
	"milliard-easy/daaz_simbank/internal/ports"
	"milliard-easy/daaz_simbank/internal/trade"
	"milliard-easy/daaz_simbank/log"
	"milliard-easy/daaz_simbank/telegram"
	"strings"

	"github.com/sirupsen/logrus"
)

func main() {
	log.InitLogger()

	logrus.Infof(log.InfoColor, "Running simbank. ©milliardeasy ©daazweb")

	ports, err := ports.ScanUSB()
	if err != nil {
		logrus.Fatalf(log.ErrorColor, "Cannot scan usb ports: %v", err)
	}

	logrus.WithField("count", len(ports)).Infof(log.InfoColor, "Port scan successfully")

	if len(ports) == 0 {
		logrus.Fatalf(log.ErrorColor, "Ports not found")
	}

	cfg, err := context.InitConfig()
	if err != nil {
		logrus.Fatalf(fmt.Sprintf(log.ErrorColor, "Cannot init config: %v"), err)
	}

	logrus.WithFields(logrus.Fields{
		"token":      cfg.Daazweb.Token,
		"operator":   cfg.Daazweb.Operator,
		"limit":      cfg.Daazweb.Limit,
		"email from": cfg.SMTP.From,
		"email to":   cfg.SMTP.To,
	}).Infof(log.InfoColor, "Configuration read successfully")

	bot := telegram.NewBot(cfg)
	iam, err := bot.Start()
	if err != nil {
		logrus.WithError(err).Fatalf(log.ErrorColor, "Cannot connect telegram bot")
	}
	logrus.WithField("username", iam.UserName).Infof(log.InfoColor, "Telegram bot successfully connected")

	bufMsgController := telegram.NewBufMsgController(bot)
	bufMsgController.Start()

	tradeCont := trade.NewController(cfg, ports, bot, daaz.NewAPI(cfg), bufMsgController)

	bot.Send(strings.Join([]string{
		fmt.Sprintf("<b>▶ Система активирована</b>"),
		fmt.Sprintf(""),
		fmt.Sprintf("<b>Токен: <code>%s</code></b>", cfg.Daazweb.Token),
		fmt.Sprintf("<b>Количество портов: <code>%d</code></b>", len(ports)),
		fmt.Sprintf("<b>Оператор: <code>%s</code></b>", cfg.Daazweb.Operator),
		fmt.Sprintf("<b>Лимит: <code>%s₽</code></b>", cfg.Daazweb.Limit),
		fmt.Sprintf("<b></b>"),
		fmt.Sprintf("<b>Да прибудет с Вами вечный разблок! 🙏</b>"),
	}, "\n"))

	tradeCont.Start()
}
