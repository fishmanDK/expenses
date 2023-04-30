package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ifDelCategory(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Удалить категорию"{
		_, keyboard := newKeyboardForCategoryes(bt, update)

		nameCategory := selectCategoryForDelCategory(bt, update, keyboard)

		if nameCategory != "" && nameCategory != "Выйти"{
			status := delCategorySQL(bt, int(update.Message.Chat.ID), nameCategory)

			if !status{
				message := tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, категорию удалить не получилось.\nПопробуйте позже.")
				message.ReplyMarkup = bt.Keyboard
				_, err := bt.Bot.Send(message)
				if err != nil {
				    panic(err)
				}
			} else{
				message := tgbotapi.NewMessage(update.Message.Chat.ID, "Удаление прошло успешно")
				message.ReplyMarkup = bt.Keyboard
				_, err := bt.Bot.Send(message)
				if err != nil {
				    panic(err)
				}
			}
		}
	
	}
}

func selectCategoryForDelCategory(bt BotConfig, update tgbotapi.Update, keyboard tgbotapi.ReplyKeyboardMarkup) string{
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбирете категорию")
	message.ReplyMarkup = keyboard
	_, err := bt.Bot.Send(message)
	if err != nil {
	    panic(err)
	}

	for update := range bt.Updates {
		if update.Message != nil{	
			if exitBool := exit(bt, update); exitBool{
				return "Выйти"
			}
			return update.Message.Text
		}
	}
	return ""
}
//----------------------------------------------------------------------------------------------------------------------------------------------------------

func ifDelProduct(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Удалить продукт"{
		cagegoryes, categoryKeyboard := newKeyboardForCategoryes(bt, update)

		updateNameCategory, statusSelectCategory := selectCategoryForDelProduct(bt, update, categoryKeyboard)
		if statusSelectCategory{

			productKeyboardd := newKeyboardForProducts(bt, updateNameCategory, cagegoryes)

			nameProduct, _ := selectProductForDelProduct(bt, update, productKeyboardd)
			fmt.Println(nameProduct)
			if nameProduct != ""{
				fmt.Println("suka blyat 2")
				status := delProductSQL(bt, int(update.Message.Chat.ID), cagegoryes[updateNameCategory.Message.Text], nameProduct)
	
				if !status{
					message := tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, категорию удалить не получилось.\nПопробуйте позже.")
					message.ReplyMarkup = bt.Keyboard
					_, err := bt.Bot.Send(message)
					if err != nil {
						panic(err)
					}
				} else{
					message := tgbotapi.NewMessage(update.Message.Chat.ID, "Удаление прошло успешно")
					message.ReplyMarkup = bt.Keyboard
					_, err := bt.Bot.Send(message)
					if err != nil {
						panic(err)
					}
				}
			}

		}
		
	}
}

func selectProductForDelProduct(bt BotConfig, update tgbotapi.Update, keyboard tgbotapi.ReplyKeyboardMarkup) (string, error){
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбирете продукт")
	message.ReplyMarkup = keyboard
	_, err := bt.Bot.Send(message)
	if err != nil {
		panic(err)
	}
	for updateNameProduct := range bt.Updates {
		if updateNameProduct.Message != nil{
			if exitBool := exit(bt, updateNameProduct); exitBool{
				break
			}
			nameProduct := updateNameProduct.Message.Text
			return nameProduct, nil
		}
	}
	return "", err
}


func selectCategoryForDelProduct(bt BotConfig, update tgbotapi.Update, keyboard tgbotapi.ReplyKeyboardMarkup) (tgbotapi.Update, bool){
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
			return update, true
		}
	}
	return update, false
}
//----------------------------------------------------------------------------------------------------------------------------------------------------------