package telegram

import (
	"fmt"
	"log"
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

			ifAddProduct(bt, update)

			ifAddPurchase(bt, update)

			ifPrintProductsINCategory(bt, update)
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

func ifAddProduct(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Добавить товар"{
		
		printUserDoc(update)
		
		continuationCategory(bt, update)
	}
}

func ifAddPurchase(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Добавить покупку"{
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название товара")
		_, err := bt.Bot.Send(message)
		if err != nil {
			log.Fatalf("error in ifAddPurchase: %s", err.Error())
		}
		printUserDoc(update)
		continuationPurchase(bt)
	}
}


func continuationCategory(bt BotConfig, update tgbotapi.Update){

	var keyboard tgbotapi.ReplyKeyboardMarkup
	categoryKeyboard, _ := printProducts(bt, update.Message.Chat.UserName, update)
	for category_name, _ := range categoryKeyboard {
		fmt.Println(category_name)
	    keyboard.Keyboard = append(keyboard.Keyboard, []tgbotapi.KeyboardButton{
	        tgbotapi.NewKeyboardButton(category_name),
	    })
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбирете категорию")
	message.ReplyMarkup = keyboard
	_, err := bt.Bot.Send(message)
	if err != nil {
	    panic(err)
	}
	
	for update := range bt.Updates {
		if update.Message != nil{
			if update.Message.Text != ""{
				message = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название товара")
				_, err = bt.Bot.Send(message)
				if err != nil {
					log.Fatalf("error in continuationCategory 'Введите название товара': %s", err.Error())
				}

				for update_forChatId := range bt.Updates {
					if update_forChatId.Message != nil{
						if update_forChatId.Message.Text != ""{
							fmt.Println(categoryKeyboard[update.Message.Text])
							AddProduct(bt, categoryKeyboard[update.Message.Text], update.Message.Chat.ID, update_forChatId.Message.Text)
						}
					}
					break
				}
				
				message := tgbotapi.NewMessage(update.Message.Chat.ID, "Товар успешно добавлен")
				_, err := bt.Bot.Send(message)
				if err != nil {
					panic(err)
				}
				fmt.Println(update.Message.Chat.UserName, update.Message.Chat.FirstName + " " + update.Message.Chat.LastName, update.Message.Chat.ID, update.Message.Text)
			}
		
		}
		break
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
					log.Panicf("error in continuationPrise: %s", err.Error())
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