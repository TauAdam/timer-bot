package main

import (
	"github.com/TauAdam/timer-bot/internal/inmemdb"
	"github.com/TauAdam/timer-bot/internal/storage"
	"github.com/TauAdam/timer-bot/internal/timer"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendInitialTimerSetMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, seconds time.Duration, label string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Timer "+label+" set for "+seconds.String()+" seconds")
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
}

func addTimerToDB(db storage.Storage, userName string, t timer.Timer) {
	if err := db.AddTimer(userName, t); err != nil {
		log.Fatalf("Error adding timer: %v", err)
	}
}

func resetTimerInDB(db storage.Storage, userName string) {
	if err := db.ResetTimer(userName); err != nil {
		log.Fatalf("Error resetting timer: %v", err)
	}
}

func sendTimerFinishedMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, label string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Timer "+label+" finished!")
	if _, err := bot.Send(msg); err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
}

func setTimer(bot *tgbotapi.BotAPI, update tgbotapi.Update, seconds time.Duration, label string, db storage.Storage) {
	sendInitialTimerSetMessage(bot, update, seconds, label)
	t := timer.Timer{Duration: seconds, StartTime: time.Now(), Label: label}
	addTimerToDB(db, update.Message.From.UserName, t)
	time.Sleep(seconds)
	resetTimerInDB(db, update.Message.From.UserName)
	sendTimerFinishedMessage(bot, update, label)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Bot token error: %v", err)
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
				argsParts := strings.SplitN(args, " ", 2)
				if len(argsParts) < 2 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /timer <seconds> <label>")
					msg.ReplyToMessageID = update.Message.MessageID
					if _, err := bot.Send(msg); err != nil {
						log.Fatalf("Error sending message: %v", err)
					}
					continue
				}

				seconds, err := time.ParseDuration(argsParts[0] + "s")
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid duration")
					msg.ReplyToMessageID = update.Message.MessageID
					if _, err := bot.Send(msg); err != nil {
						log.Fatalf("Error sending message: %v", err)
					}
					continue
				}

				label := argsParts[1]
				go setTimer(bot, update, seconds, label, db)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				log.Fatalf("Error sending message: %v", err)
			}
		}
	}
}
