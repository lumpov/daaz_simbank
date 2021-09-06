package telegram

import (
	"fmt"
	"milliard-easy/daaz_simbank/context"
	"milliard-easy/daaz_simbank/log"
	"strings"

	"github.com/sirupsen/logrus"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Bot for telegram
type Bot struct {
	cfg *context.Config
	api *tgbot.BotAPI
}

// NewBot instance
func NewBot(cfg *context.Config) *Bot {
	return &Bot{
		cfg: cfg,
	}
}

// Start listening msgs
func (b *Bot) Start() (*tgbot.User, error) {
	var err error
	b.api, err = tgbot.NewBotAPI(b.cfg.Telegram.Token)
	if err != nil {
		return nil, err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}

	go func() {
		for i := range updates {
			if b.checkAuth(i) {
				logrus.WithFields(logrus.Fields{
					"from":    i.Message.Chat.UserName,
					"chat id": i.Message.Chat.ID,
					"message": i.Message.Text,
				}).Debugf(log.DebugColor, "New telegram message")

				b.send(i.Message.Chat.ID, "<b>⛔ Неизвестная команда</b>")
			}
		}
	}()

	return &b.api.Self, nil
}

func (b *Bot) checkAuth(update tgbotapi.Update) bool {
	for i := range b.cfg.Telegram.Chats {
		if update.Message.Chat.ID == b.cfg.Telegram.Chats[i] {
			return true
		}
	}

	l := logrus.WithField("id", update.Message.Chat.ID)
	if update.Message.Command() == "connect" {
		token := strings.Split(update.Message.Text, " ")[1]
		if b.cfg.Telegram.AccessToken == token {
			b.cfg.Telegram.Chats = append(b.cfg.Telegram.Chats, update.Message.Chat.ID)
			if err := b.cfg.Save(); err != nil {
				b.send(update.Message.Chat.ID, strings.Join([]string{
					fmt.Sprintf("<b>❌ Не удалось сохранить чат</b>"),
					fmt.Sprintf("<b>Подробнее в журналах приложения</b>"),
				}, "\n"))
				l.WithError(err).Errorf(log.ErrorColor, "cannot save telegram chat")
			} else {
				b.send(update.Message.Chat.ID, strings.Join([]string{
					fmt.Sprintf("<b>✅ Чат успешно добавлен</b>"),
					fmt.Sprintf("<b>Ваш ID: <code>%d</code></b>", update.Message.Chat.ID),
					fmt.Sprintf("<b>Всего чатов: <code>%d</code></b>", len(b.cfg.Telegram.Chats)),
					fmt.Sprintf(""),
					fmt.Sprintf("<b>Для поддержания эстетического кайфа,</b>"),
					fmt.Sprintf("<b>рекомендуем очистить переписку сейчас.</b>"),
					fmt.Sprintf("<b>Чат выше - бесполезен</b>"),
				}, "\n"))
				l.Infof(log.InfoColor, "Telegram chat connected")
			}
		} else {
			b.send(update.Message.Chat.ID, "<b>❌ Неверный токен доступа</b>")
			l.WithFields(logrus.Fields{
				"token":        token,
				"access token": b.cfg.Telegram.AccessToken,
			}).Warnf(log.WarningColor, "Invalid access token")
		}
	} else {
		b.send(update.Message.Chat.ID, strings.Join([]string{
			fmt.Sprintf("<b>⛔ Ваш чат не зарегистрирован</b>"),
			fmt.Sprintf(""),
			fmt.Sprintf("<b>Используйте команду: </b>"),
			fmt.Sprintf("<code>/connect (токен - доступа)</code>"),
			fmt.Sprintf("<b>для регистрации чата</b>"),
		}, "\n"))
	}
	return false
}

func (b *Bot) send(id int64, str string) {
	msg := tgbot.NewMessage(id, str)
	msg.ParseMode = "HTML"

	l := logrus.WithFields(logrus.Fields{
		"id":      id,
		"message": str,
	})
	l.Debugf(log.DebugColor, "Sending telegram message")
	_, err := b.api.Send(msg)
	if err != nil {
		l.Errorf(log.ErrorColor, fmt.Errorf("Cannot send telegram message: %v", err))
	}
}

// Send message to all connected telegram accounts
func (b *Bot) Send(str string) {
	for i := range b.cfg.Telegram.Chats {
		b.send(b.cfg.Telegram.Chats[i], str)
	}
}
