package main

import (
	"github.com/TauAdam/timer-bot/internal/inmemdb"
	"github.com/TauAdam/timer-bot/internal/storage"
	"github.com/TauAdam/timer-bot/internal/timer"
	"github.com/joho/godotenv"
	"log"
	"os"
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
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose a timer duration:")
				msg.ReplyMarkup = createTimerOptionsKeyboard()
				if _, err := bot.Send(msg); err != nil {
					log.Fatalf("Error sending message: %v", err)
				}
				if update.CallbackQuery != nil {
					callbackData := update.CallbackQuery.Data
					seconds, err := time.ParseDuration(callbackData)
					if err != nil {
						log.Fatalf("Invalid duration: %v", err)
					}
					label := "Predefined timer"
					callbackUpdate := tgbotapi.Update{
						Message: update.CallbackQuery.Message,
					}
					go setTimer(bot, callbackUpdate, seconds, label, db)
					callbackConfig := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
					if _, err := bot.Request(callbackConfig); err != nil {
						log.Fatalf("Error acknowledging callback query: %v", err)
					}
				}
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

func createTimerOptionsKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1 minute", "1m"),
			tgbotapi.NewInlineKeyboardButtonData("5 minutes", "5m"),
			tgbotapi.NewInlineKeyboardButtonData("10 minutes", "10m"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("15 minutes", "15m"),
			tgbotapi.NewInlineKeyboardButtonData("30 minutes", "30m"),
			tgbotapi.NewInlineKeyboardButtonData("60 minutes", "60m"),
		),
	)
	return keyboard
}
