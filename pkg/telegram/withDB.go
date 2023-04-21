package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// "github.com/fishmanDK/expenses"
// "github.com/fishmanDK/expenses/pkg/repository"


type Product struct{ 
	Product_name string  
	Price		 float32 
}

type User struct{
	FirstName_lastNAme string 
	ChatID             string	  
}

func ifPrintCategoryes(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Вывести список продуктов"{

		allRows, _ := printAllProductsWithSelectDB(bt, update)

		message := tgbotapi.NewMessage(update.Message.Chat.ID, allRows)
		_, err := bt.Bot.Send(message)
		if err != nil {
			panic(err)
		}

	}

}

func printAllProductsWithSelectDB(bt BotConfig, update tgbotapi.Update) (string, error){
	var (
		allRows string
		sum float32
	)


	allRows += fmt.Sprintf("Список товаров пользователя: %s\n\n\n", update.Message.Chat.UserName)
	allRows += "Категория: КАТЕГОРИЯ\n\n"
	product := Product{}
	rows, _ := bt.ConfigDB.Queryx("SELECT product_name, price FROM product")
	for rows.Next() {
		err := rows.StructScan(&product)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
			return " ", err
		}

		fmt.Printf("%#v\n", product)
		allRows += fmt.Sprintf("%s: %.2fp.\n", product.Product_name, product.Price)

		sum += sumPrice(product.Price)
	}
	allRows += fmt.Sprintf("\nОбщая сумма: %.2f", sum)
	return allRows, nil
}

func sumPrice(price float32) float32{
	var sum float32
	sum += price

	return sum
}