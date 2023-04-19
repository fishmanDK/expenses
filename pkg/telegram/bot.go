package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func MainReply(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, keyboard tgbotapi.ReplyKeyboardMarkup){
	for update := range updates {
		if update.Message != nil{
			if update.Message.Command() == "start"{
				startMessage := fmt.Sprintf("Привет %s, это бот для оптимизации твоох расходов", update.Message.Chat.UserName)
				message := tgbotapi.NewMessage(update.Message.Chat.ID, startMessage)
				message.ReplyMarkup = keyboard
				
				_, err := bot.Send(message)
				if err != nil {
					panic(err)
				}
			}

			if update.Message.Text == "Добавить категорию"{
				message := tgbotapi.NewMessage(update.Message.Chat.ID, "Как она называется")
				_, err := bot.Send(message)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}