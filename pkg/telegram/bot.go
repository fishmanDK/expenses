package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotConfig struct{
	Bot      *tgbotapi.BotAPI
	Updates  tgbotapi.UpdatesChannel
	Keyboard tgbotapi.ReplyKeyboardMarkup
	// Update   tgbotapi.Update
}

func NewBotConfig(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, keyboard tgbotapi.ReplyKeyboardMarkup) *BotConfig{
	return &BotConfig{
		Bot: bot,
		Updates: updates,
		Keyboard: keyboard,
	}
}

func MainReply(bt BotConfig){
	for update := range bt.Updates {
		if update.Message != nil{
			ifCommandStart(bt, update)

			ifAddCategory(bt, update)

			ifAddPurchase(bt, update)
		}
	}
}


func ifCommandStart(bt BotConfig, update tgbotapi.Update){
	if update.Message.Command() == "start"{
		startMessage := fmt.Sprintf("Привет %s, это бот для оптимизации твоох расходов", update.Message.Chat.UserName)
		message := tgbotapi.NewMessage(update.Message.Chat.ID, startMessage)
		message.ReplyMarkup = bt.Keyboard
		
		_, err := bt.Bot.Send(message)
		if err != nil {
			panic(err)
		}
		fmt.Println(update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.Text)
	}
}

func ifAddCategory(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Добавить товар"{
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Как он называется")
		_, err := bt.Bot.Send(message)
		if err != nil {
			panic(err)
		}
		fmt.Println(update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.Text)
		
		continuationCategory(bt)
	}
}

func ifAddPurchase(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Добавить покупку"{
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название товара")
		_, err := bt.Bot.Send(message)
		if err != nil {
			panic(err)
		}
		fmt.Println(update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.Text)
	}
}


func continuationCategory(bt BotConfig){
	for update := range bt.Updates {
		if update.Message != nil{
			if update.Message.Text != ""{
				message := tgbotapi.NewMessage(update.Message.Chat.ID, "Товар успешно добавлен")
				_, err := bt.Bot.Send(message)
				if err != nil {
					panic(err)
				}
				fmt.Println(update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.Text)
			}
		break
		}
	}
}