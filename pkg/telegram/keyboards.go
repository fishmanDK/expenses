package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func newKeyboardForCategoryes(bt BotConfig, update tgbotapi.Update) (map[string]int, tgbotapi.ReplyKeyboardMarkup){
	var keyboard tgbotapi.ReplyKeyboardMarkup
	categoryKeyboard, _ := printCategoryesSQL(bt, update.Message.Chat.UserName, update)

	for category_name := range categoryKeyboard {
	    keyboard.Keyboard = append(keyboard.Keyboard, []tgbotapi.KeyboardButton{
	        tgbotapi.NewKeyboardButton(category_name),
	    })
	}
	keyboard.Keyboard = append(keyboard.Keyboard, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Выйти"),
	})

	return categoryKeyboard, keyboard
	
}

func newKeyboardForProducts(bt BotConfig, update tgbotapi.Update, categoryKeyboard map[string]int) tgbotapi.ReplyKeyboardMarkup{
	var keyboard tgbotapi.ReplyKeyboardMarkup
	_, nameProducts, _ := getProductsSQL(bt, update, categoryKeyboard[update.Message.Text])
	

	for _, name := range nameProducts {
	    keyboard.Keyboard = append(keyboard.Keyboard, []tgbotapi.KeyboardButton{
	        tgbotapi.NewKeyboardButton(name),
	    })
	}
	keyboard.Keyboard = append(keyboard.Keyboard, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Выйти"),
	})

	return keyboard
	// message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбирете продукт")
	// message.ReplyMarkup = keyboard
	// _, err := bt.Bot.Send(message)
	// if err != nil {
	//     panic(err)
	// }
}