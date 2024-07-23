package main

import (
	"github.com/TauAdam/timer-bot/internal/inmemdb"
	"github.com/TauAdam/timer-bot/internal/storage"
	"github.com/TauAdam/timer-bot/internal/timer"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func setTimer(bot *tgbotapi.BotAPI, update tgbotapi.Update, seconds time.Duration, db storage.Storage) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Timer set for "+seconds.String()+" seconds")
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}

	t := timer.Timer{Duration: seconds, StartTime: time.Now()}
	if err := db.AddTimer(update.Message.From.UserName, t); err != nil {
		log.Panic(err)
	}

	time.Sleep(seconds)

	if err := db.DeleteTimer(update.Message.From.UserName); err != nil {
		log.Panic(err)
	}

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

	db := inmemdb.NewInMemoryDB()

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
				if len(args) == 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /timer <seconds> or /timer <option>")
					msg.ReplyToMessageID = update.Message.MessageID
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
					continue
				}

				var seconds time.Duration
				switch args {
				case "10", "30", "60", "90":
					seconds, _ = time.ParseDuration(args + "s")
				default:
					seconds, err = time.ParseDuration(args + "s")
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid duration")
						msg.ReplyToMessageID = update.Message.MessageID
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
						continue
					}
				}

				go setTimer(bot, update, seconds, db)
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
