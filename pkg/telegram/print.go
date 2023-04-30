package telegram

import (

	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func ifPrintProductsINCategory(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Вывести список продуктов (в категории)"{


		categoryKeyboard, keyboard := newKeyboardForCategoryes(bt, update)

		getInfAboutProductsInCategory(bt, update, categoryKeyboard, keyboard)
	}
}

func getInfAboutProductsInCategory(bt BotConfig, update tgbotapi.Update, categoryKeyboard map[string]int, keyboard tgbotapi.ReplyKeyboardMarkup) {
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбирете категорию")
	message.ReplyMarkup = keyboard
	_, err := bt.Bot.Send(message)
	if err != nil {
	    panic(err)
	}

	for update := range bt.Updates {
		if update.Message != nil{	
			if exitBool := exit(bt, update); exitBool{
				break
			}

			nameCategory, _, _ := getProductsSQL(bt, update, categoryKeyboard[update.Message.Text])
			
			message := tgbotapi.NewMessage(update.Message.Chat.ID, nameCategory)
			_, err := bt.Bot.Send(message)
			if err != nil {
				log.Fatalf("error in ifPrintProductsINCategory: %s", err.Error())
			}


		}
	}
}
//----------------------------------------------------------------------------------------------------------------------------------------------------------

func ifPrintAllCategoryes(bt BotConfig, update tgbotapi.Update){
	// if update.Message != nil && update.Message.Text == "Вывести все категории"{
	// 	products, prices := printCategoryesInMainReply(bt, update)
	// }
}
//----------------------------------------------------------------------------------------------------------------------------------------------------------
