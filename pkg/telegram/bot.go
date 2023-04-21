package telegram

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

type BotConfig struct{
	Bot      *tgbotapi.BotAPI
	Updates  tgbotapi.UpdatesChannel
	Keyboard tgbotapi.ReplyKeyboardMarkup
	ConfigDB *sqlx.DB
}

func NewBotConfig(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, keyboard tgbotapi.ReplyKeyboardMarkup, configDB *sqlx.DB) *BotConfig{
	return &BotConfig{
		Bot: bot,
		Updates: updates,
		Keyboard: keyboard,
		ConfigDB: configDB,
	}
}

func MainReply(bt BotConfig){
	for update := range bt.Updates {
		if update.Message != nil{
			ifCommandStart(bt, update)

			ifAddCategory(bt, update)

			ifAddPurchase(bt, update)

			ifPrintCategoryes(bt, update)
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
		printUserDoc(update)
	}
}

func ifAddCategory(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Добавить товар"{
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Как он называется")
		_, err := bt.Bot.Send(message)
		if err != nil {
			panic(err)
		}
		printUserDoc(update)
		
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
		printUserDoc(update)
		continuationPurchase(bt)
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
				fmt.Println(update.Message.Chat.UserName, update.Message.Chat.FirstName + " " + update.Message.Chat.LastName, update.Message.Chat.ID, update.Message.Text)			}
		break
		}
	}
}

func continuationPurchase(bt BotConfig){
	for update := range bt.Updates {
		if update.Message != nil{
			if update.Message.Text != ""{
				message := tgbotapi.NewMessage(update.Message.Chat.ID, "Теперь его цену")
				_, err := bt.Bot.Send(message)
				if err != nil {
					panic(err)
				}
				printUserDoc(update)
				continuationPrise(bt)
			}
		break
		}
	}
}

func continuationPrise(bt BotConfig){
	for update := range bt.Updates {
		if update.Message != nil{
			prise := update.Message.Text
			if floatPrise, _ := strconv.ParseFloat(prise, 32); floatPrise > 0{
				message := tgbotapi.NewMessage(update.Message.Chat.ID, "Покупка успешно добавленна")
				_, err := bt.Bot.Send(message)
				if err != nil {
					panic(err)
				}
				printUserDoc(update)
				break

			} else{
				message := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите правально цену")
				_, err := bt.Bot.Send(message)
				if err != nil {
					panic(err)
				}
				printUserDoc(update)
			}
		}
	}
}
 
func printUserDoc(update tgbotapi.Update) {
    fmt.Println(update.Message.Chat.UserName, update.Message.Chat.FirstName+" "+update.Message.Chat.LastName, update.Message.Chat.ID, update.Message.Text)
}