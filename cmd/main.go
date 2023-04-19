package main

import (
	"log"
	"github.com/fishmanDK/expenses/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	bot, err := tgbotapi.NewBotAPI("5700090858:AAHjE0mZFZR6z-H9S2pxLu5WhhXPknLDkFA")
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить категорию"),
			tgbotapi.NewKeyboardButton("Добавить покупку"),
		),
	)
	updates := bot.GetUpdatesChan(u)

	
	telegram.MainReply(bot, updates, keyboard)


	
}
