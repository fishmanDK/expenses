package telegram

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func ifAddProduct(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Добавить товар"{		
		continuationAddProduct(bt, update)
	}

}

func continuationAddProduct(bt BotConfig, update tgbotapi.Update){

	categoryKeyboard, keyboard := newKeyboardForCategoryes(bt, update)

	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбирете категорию")
	message.ReplyMarkup = keyboard
	_, err := bt.Bot.Send(message)
	if err != nil {
	    panic(err)
	}

	workWithAddproduct(bt, categoryKeyboard)
}

func workWithAddproduct(bt BotConfig, categoryKeyboard map[string]int){
	for update := range bt.Updates {
		if update.Message != nil{	
			if exitBool := exit(bt, update); exitBool{
				break
			}
			nameProduct, price := getInfAboutAddProduct(bt, update, categoryKeyboard)
			
			addProductSQL(bt, categoryKeyboard[update.Message.Text], update.Message.Chat.ID, nameProduct, price)

			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Товар успешно добавлен")
			message.ReplyMarkup = bt.Keyboard
			_, err := bt.Bot.Send(message)
			if err != nil {
				panic(err)
			}
			break
		}
	}
}

func getInfAboutAddProduct(bt BotConfig, update tgbotapi.Update, categoryKeyboard map[string]int) (string, string){
	nameProduct, err := getNameProduct(bt, update)
	if err != nil{
		log.Fatalf("erron in 'getNameProduct': %s", err.Error())
	}

	price, err := getPrice(bt, update)
	if err != nil{
		log.Fatalf("erron in 'getPrice': %s", err.Error())
	}

	return nameProduct, price
}

func getPrice(bt BotConfig, update tgbotapi.Update) (string, error){
    var price string
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Теперь его цену")
	_, err := bt.Bot.Send(message)
	if err != nil {
		log.Fatalf("error in continuationCategory 'Введите название товара': %s", err.Error())
	}
    for update := range bt.Updates {
        if update.Message != nil {
            price += update.Message.Text
            return price, nil
        }
    }
    return "", err
}

func getNameProduct(bt BotConfig, update tgbotapi.Update) (string, error){	
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название товара")
	message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	
	_, err := bt.Bot.Send(message)
	if err != nil {
		log.Fatalf("error in getNameProduct 'Введите название товара': %s", err.Error())
	}
	for update_forChatId := range bt.Updates {
		if update_forChatId.Message != nil{
			nameProduct := update_forChatId.Message.Text
			return nameProduct, nil
			
		}
	}
	return "", err
}
//----------------------------------------------------------------------------------------------------------------------------------------------------------

func ifAddPurchase(bt BotConfig, update tgbotapi.Update) {
    if update.Message.Text == "Добавить покупку" {
		continuationPurchase(bt, update)
    }
}

func continuationPurchase(bt BotConfig, update tgbotapi.Update){
	var nameCategory string
	categoryKeyboard, keyboardCategory := newKeyboardForCategoryes(bt, update)

	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбирете категорию")
	message.ReplyMarkup = keyboardCategory
	_, err := bt.Bot.Send(message)
	if err != nil {
	    panic(err)
	}
	for updateNameCategory := range bt.Updates{
		if updateNameCategory.Message != nil{
			if exitBool := exit(bt, update); exitBool{
				break
			}
			nameCategory = updateNameCategory.Message.Text
			break
		}
	}
	fmt.Println(nameCategory)
	if nameCategory != "" && nameCategory != "Выйти"{
		nameProduct, countProduct := getInfoAboutPurchase(bt, update, categoryKeyboard, nameCategory)
		if nameProduct != ""{
			slqCountProductbt(bt, countProduct, update.Message.Chat.ID, nameProduct)

			message = tgbotapi.NewMessage(update.Message.Chat.ID, "Добавление прошло успешно")
			message.ReplyMarkup = bt.Keyboard		
			_, err = bt.Bot.Send(message)
			if err != nil {
				log.Fatalf("error in continuationCategory 'Введите название товара': %s", err.Error())
			}
		} 
	} else if nameCategory == "Выйти"{
		message = tgbotapi.NewMessage(update.Message.Chat.ID, "Возвращаемся в главное меню")
			message.ReplyMarkup = bt.Keyboard		
			_, err = bt.Bot.Send(message)
			if err != nil {
				log.Fatalf("error in continuationPurchase: %s", err.Error())
			}
	}
	
}

func getInfoAboutPurchase(bt BotConfig, update tgbotapi.Update, categoryKeyboard map[string]int, nameCategory string) (string, int){
	if update.Message != nil{
		keyboard := newKeyboardForProductsForPurchase(bt, update, categoryKeyboard, nameCategory)
		nameProduct, err := selectProduct(bt, update, keyboard)
		if nameProduct == ""{
			return "", 0
		}
		fmt.Println(nameProduct)
		if err != nil{
			log.Fatalf("error in selectProduct: %s", err.Error())
		}
		
		countProduct, err := getCountProduct(bt, update)
		if err != nil{
			log.Fatalf("error in selectProduct: %s", err.Error()) /////////
		}
		return nameProduct, countProduct
	}
	
	return "", 0
}

func selectProduct(bt BotConfig, update tgbotapi.Update, keyboard tgbotapi.ReplyKeyboardMarkup) (string, error){
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

func getCountProduct(bt BotConfig, update tgbotapi.Update) (int, error){
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Теперь количество")
	message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err := bt.Bot.Send(message)
	if err != nil {
		log.Fatalf("error in getCountProduct: %s", err.Error())
	}
	for update_forChatId := range bt.Updates {
		if update_forChatId.Message != nil{
			if exitBool := exit(bt, update); exitBool{
				break
			}
			countProduct, _ := strconv.Atoi(update_forChatId.Message.Text)
			return countProduct, nil
			
		}
	}
	return 0, err
}

func newKeyboardForProductsForPurchase(bt BotConfig, update tgbotapi.Update, categoryKeyboard map[string]int, nameCategory string) tgbotapi.ReplyKeyboardMarkup{
	var keyboard tgbotapi.ReplyKeyboardMarkup
	_, nameProducts, _ := getProductsSQL(bt, update, categoryKeyboard[nameCategory])
	
	for _, name := range nameProducts {
	    keyboard.Keyboard = append(keyboard.Keyboard, []tgbotapi.KeyboardButton{
	        tgbotapi.NewKeyboardButton(name),
	    })
	}
	keyboard.Keyboard = append(keyboard.Keyboard, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Выйти"),
	})

	return keyboard
}
//----------------------------------------------------------------------------------------------------------------------------------------------------------

func ifAddCatogory(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Добавить категорию"{
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите имя категории")
		message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Выйти"),
			))
		message.ReplyMarkup = keyboard

		exit(bt, update)
		
		_, err := bt.Bot.Send(message)
		if err != nil {
		    panic(err)
		}
		
		workWithAddCategory(bt, update)
	}
}

func workWithAddCategory(bt BotConfig, update tgbotapi.Update){
	for update := range bt.Updates {
		if update.Message != nil{	
			if exitBool := exit(bt, update); exitBool{
				break
			}

			addCategorySQL(bt, update.Message.Chat.ID, update.Message.Text, update)
			
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Категория успешно добавленна")
			message.ReplyMarkup = bt.Keyboard
			_, err := bt.Bot.Send(message)
			if err != nil {
				panic(err)
			}
			break
		}
	}
}
//----------------------------------------------------------------------------------------------------------------------------------------------------------