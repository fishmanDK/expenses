package telegram

import (
	"fmt"
	"log"

	// "strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	// "github.com/fishmanDK/expenses/pkg/telegram/DB"
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

			ifAddCatogory(bt, update)

			ifPrintProductsINCategory(bt, update)

			ifDelCategory(bt, update)

			ifPrintAllCategoryes(bt, update)

			ifDelProduct(bt, update)
		}
	}
}


func ifCommandStart(bt BotConfig, update tgbotapi.Update){
	if update.Message.Command() == "start"{
		examUser(bt, update)
		startMessage := fmt.Sprintf("Привет %s, это бот для оптимизации твоох расходов", update.Message.Chat.UserName)
		message := tgbotapi.NewMessage(update.Message.Chat.ID, startMessage)
		message.ReplyMarkup = bt.Keyboard
		
		_, err := bt.Bot.Send(message)
		if err != nil {
			panic(err)
		}		
	}
}

func examUser(bt BotConfig, update tgbotapi.Update){
	fmt.Println(1)
	statusExam := examUserSQL(bt, int(update.Message.Chat.ID))
	fmt.Println(statusExam)
	if statusExam{
		fmt.Println(3)
		fmt.Println("Регистрация")
		statusRegistration := registrationSQL(bt, update.Message.Chat.FirstName + " " + update.Message.Chat.LastName, int(update.Message.Chat.ID))
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Зарегистрироваться"),
			))
		
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Вероятно вы впервые зашли сюда, для продолжения вам необходимо зарегистрироваться")
		message.ReplyMarkup = keyboard		
		_, err := bt.Bot.Send(message)
		if err != nil {
			log.Fatalf("error in examUser 'Вероятно вы впервые зашли сюда, для продолжения вам необходимо зарегистрироваться': %s", err.Error())
		}
		for update := range bt.Updates{
			if update.Message != nil && update.Message.Text == "Зарегистрироваться"{
				if !statusRegistration{
					message := tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, произошли неполадки. Попробуйте позже зарегистрироваться")
					message.ReplyMarkup = bt.Keyboard		
					_, err := bt.Bot.Send(message)
					if err != nil {
						log.Fatalf("error in examUser 'Извините, произошли неполадки. Попробуйте позже зарегистрироваться': %s", err.Error())
					}
				}
				break
			}
		}


	}
}

func exit(bt BotConfig, update tgbotapi.Update) bool{
	if update.Message.Text == "Выйти"{
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Возвращаемся в главное меню")
		message.ReplyMarkup = bt.Keyboard

		_, err := bt.Bot.Send(message)
		if err != nil {
			panic(err)
		}
		return true
	}
	return false
}
 
