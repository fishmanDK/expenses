package telegram

import (
	"fmt"
	"log"
	"strconv"

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
			// m := printCategoryesInMainReply(bt, update.Message.Chat.UserName, update)
			// message := tgbotapi.NewMessage(update.Message.Chat.ID, m)
			// _, err := bt.Bot.Send(message)
			// if err != nil {
			// 	log.Fatalf("error in mainMessage: %s", err.Error())
			// }
			
			ifCommandStart(bt, update)

			ifAddProduct(bt, update)

			ifAddPurchase(bt, update)

			ifAddCatogory(bt, update)

			ifPrintProductsINCategory(bt, update)
		}
	}
}

func ifAddPurchase(bt BotConfig, update tgbotapi.Update) {
    if update.Message.Text == "Добавить покупку" {
		continuationPurchase(bt, update)
    }
}
func continuationPurchase(bt BotConfig, update tgbotapi.Update){
	categoryKeyboard := newKeyboardForCategoryes(bt, update)
	fmt.Println(categoryKeyboard)


	for update := range bt.Updates {
		if update.Message != nil{	
			if exitBool := exit(bt, update); exitBool == true{
				break
			}
			newKeyboardForProducts(bt, update, categoryKeyboard)
			var (
					nameProduct string
					countProduct int
			)
			for updateNameProduct := range bt.Updates {
				if updateNameProduct.Message != nil{
					nameProduct = updateNameProduct.Message.Text
					break
					
				}
			}
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Теперь количество")
			message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			
			_, err := bt.Bot.Send(message)
			if err != nil {
				log.Fatalf("error in continuationCategory 'Введите название товара': %s", err.Error())
			}
			for update_forChatId := range bt.Updates {
				if update_forChatId.Message != nil{
					countProduct, _ = strconv.Atoi(update_forChatId.Message.Text)
					break
					
				}
			}
			fmt.Println(nameProduct,countProduct)
			slqCountProductbt(bt, countProduct, update.Message.Chat.ID, nameProduct)
			
			message = tgbotapi.NewMessage(update.Message.Chat.ID, "Добавление прошло успешно")
			message.ReplyMarkup = bt.Keyboard		
			_, err = bt.Bot.Send(message)
			if err != nil {
				log.Fatalf("error in continuationCategory 'Введите название товара': %s", err.Error())
			}
			break
			
		}
	}
}


func getCount(bt BotConfig) int {
    var price int
    for update := range bt.Updates {
        if update.Message != nil {
            price, _ = strconv.Atoi(update.Message.Text)
            return price
        }
    }
    return 0
}



func ifPrintProductsINCategory(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Вывести список продуктов (в категории)"{
		categoryKeyboard := newKeyboardForCategoryes(bt, update)
		for update := range bt.Updates {
			if update.Message != nil{	
				if exitBool := exit(bt, update); exitBool == true{
					break
				}
				fmt.Println(categoryKeyboard[update.Message.Text])
				str, _, _ := getProductsSQL(bt, update, categoryKeyboard[update.Message.Text])
			    
				message := tgbotapi.NewMessage(update.Message.Chat.ID, str)
				_, err := bt.Bot.Send(message)
				if err != nil {
				    panic(err)
				}


			}
		}
	}
}


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
		
		printUserDoc(update)
		workWithAddCategory(bt, update)
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
		
		continuationAddProduct(bt, update)
	}

}

func continuationAddProduct(bt BotConfig, update tgbotapi.Update){

	categoryKeyboard := newKeyboardForCategoryes(bt, update)

	workWithAddproduct(bt, categoryKeyboard)
}


func newKeyboardForCategoryes(bt BotConfig, update tgbotapi.Update) map[string]int{
	var keyboard tgbotapi.ReplyKeyboardMarkup
	categoryKeyboard, _ := printCategoryes(bt, update.Message.Chat.UserName, update)
	for category_name, _ := range categoryKeyboard {
		fmt.Println(category_name)
	    keyboard.Keyboard = append(keyboard.Keyboard, []tgbotapi.KeyboardButton{
	        tgbotapi.NewKeyboardButton(category_name),
	    })
	}
	keyboard.Keyboard = append(keyboard.Keyboard, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Выйти"),
	})
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбирете категорию")
	message.ReplyMarkup = keyboard
	_, err := bt.Bot.Send(message)
	if err != nil {
	    panic(err)
	}

	return categoryKeyboard
	
}

func newKeyboardForProducts(bt BotConfig, update tgbotapi.Update, categoryKeyboard map[string]int){
	var keyboard tgbotapi.ReplyKeyboardMarkup
	_, nameProducts, _ := getProductsSQL(bt, update, categoryKeyboard[update.Message.Text])
	for _, name := range nameProducts {
		fmt.Println(name)
	    keyboard.Keyboard = append(keyboard.Keyboard, []tgbotapi.KeyboardButton{
	        tgbotapi.NewKeyboardButton(name),
	    })
	}
	keyboard.Keyboard = append(keyboard.Keyboard, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Выйти"),
	})
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбирете продукт1")
	message.ReplyMarkup = keyboard
	_, err := bt.Bot.Send(message)
	if err != nil {
	    panic(err)
	}
}


func workWithAddCategory(bt BotConfig, update tgbotapi.Update){
	for update := range bt.Updates {
		if update.Message != nil{	
			if exitBool := exit(bt, update); exitBool == true{
				break
			}

			addCategorySQL(bt, update.Message.Chat.ID, update.Message.Text)
			
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

func workWithAddproduct(bt BotConfig, categoryKeyboard map[string]int){
	for update := range bt.Updates {
		if update.Message != nil{	
			if exitBool := exit(bt, update); exitBool == true{
				break
			}

			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название товара")
			message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			
			_, err := bt.Bot.Send(message)
			if err != nil {
				log.Fatalf("error in continuationCategory 'Введите название товара': %s", err.Error())
			}

			addProduct(bt, update, categoryKeyboard)
			
			message = tgbotapi.NewMessage(update.Message.Chat.ID, "Товар успешно добавлен")
			message.ReplyMarkup = bt.Keyboard
			_, err = bt.Bot.Send(message)
			if err != nil {
				panic(err)
			}
			break


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

func getPrice(bt BotConfig) string {
    var price string
    for update := range bt.Updates {
        if update.Message != nil {
            price += update.Message.Text
            return price
        }
    }
    return ""
}

func addProduct(bt BotConfig, update tgbotapi.Update, categoryKeyboard map[string]int){
	var nameProduct string
	for update_forChatId := range bt.Updates {
		if update_forChatId.Message != nil{
			nameProduct = update_forChatId.Message.Text
			break
			
		}
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Теперь его цену")
	_, err := bt.Bot.Send(message)
	if err != nil {
		log.Fatalf("error in continuationCategory 'Введите название товара': %s", err.Error())
	}
	price := getPrice(bt)

	addProductSQL(bt, categoryKeyboard[update.Message.Text], update.Message.Chat.ID, nameProduct, price)

}


 
func printUserDoc(update tgbotapi.Update) {
    fmt.Println(update.Message.Chat.UserName, update.Message.Chat.FirstName+" "+update.Message.Chat.LastName, update.Message.Chat.ID, update.Message.Text)
}

