package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func setTimer(bot *tgbotapi.BotAPI, update tgbotapi.Update, seconds time.Duration) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Timer set for "+seconds.String()+" seconds")
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}

	time.Sleep(seconds)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Timer finished!")
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI("some-token")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "timer":
				args := update.Message.CommandArguments()
				if len(args) != 1 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /timer <seconds>")
					msg.ReplyToMessageID = update.Message.MessageID
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
					continue
				}

				seconds, err := time.ParseDuration(args + "s")
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid duration")
					msg.ReplyToMessageID = update.Message.MessageID
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
					continue
				}

				go setTimer(bot, update, seconds)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}
