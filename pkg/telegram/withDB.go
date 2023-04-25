package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// "github.com/fishmanDK/expenses"
// "github.com/fishmanDK/expenses/pkg/repository"
type Category struct{
	Category_name string  `db:"category_name"`
}


type Product struct{
	ID 			 string  `db:"id"`
	Categoru_id  string  `db:"category_id"`
	User_id 	 string  `db:"user_id"`
	Product_name string  `db:"product_name"`
	Price		 float32 `db:"price"`
	Count 		 string  `db:"count"`
}

type User struct{
	FirstName_lastNAme string `db:"firstName_lastNAme"`
	ChatID             string `db:"id"`
}

func ifPrintProductsINCategory(bt BotConfig, update tgbotapi.Update){
	if update.Message.Text == "Вывести список продуктов"{

		allRows, err := printAllProductsWithSelectDB(bt, update)
		if err != nil{
			log.Fatalf("error with ifPrintProductsINCategory: %s", err.Error())
		}

		message := tgbotapi.NewMessage(update.Message.Chat.ID, allRows)
		_, err = bt.Bot.Send(message)
		if err != nil {
			panic(err)
		}

	}

}

// func printAllProductsWithSelectDB(bt BotConfig, update tgbotapi.Update) (string, error){
	// var (
	// 	allRows string
	// 	sum float32
	// )



	// allRows += fmt.Sprintf("Список товаров пользователя: %s\n\n\n", update.Message.Chat.UserName)
	// allRows += "Категория: КАТЕГОРИЯ\n\n"
	// product := Product{}
	// rows, _ := bt.ConfigDB.Queryx(fmt.Sprintf("SELECT * FROM product WHERE user_id = (SELECT id FROM users_list WHERE chatID = %d)", update.Message.Chat.ID))

	// allRows += fmt.Sprintf("Список товаров пользователя: %s\n\n\n", update.Message.Chat.UserName)
	// for rows.Next() {
	// 	err := rows.StructScan(&product)
	// 	if err != nil {
	// 		log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
	// 		return " ", err
	// 	}

	// 	fmt.Printf("%#v\n", product)
	// 	allRows += fmt.Sprintf("%s: %.2fp.\n", product.Product_name, product.Price)

	// 	sum += sumPrice(product.Price)
	// }
	// allRows += fmt.Sprintf("\nОбщая сумма: %.2f", sum)
	// return allRows, nil
// }

func printAllProductsWithSelectDB(bt BotConfig, update tgbotapi.Update) (string, error){
	var (
		allRows string
		sum float32
	)


	allRows += fmt.Sprintf("Список товаров пользователя: %s\n\n\n", update.Message.Chat.UserName)
	allRows += "Категория: КАТЕГОРИЯ\n\n"
	product := Product{}
	rows, err := bt.ConfigDB.Queryx(fmt.Sprintf("SELECT * FROM product WHERE user_id = (SELECT id FROM users_list WHERE chatID = %d)", update.Message.Chat.ID))
	if err != nil{
		log.Fatalf("error in printAllProductsWithSelectDB(rows): %s", err.Error())
	}
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


func printProducts(bt BotConfig, userName string, update tgbotapi.Update)(string, error){
	var (
		allRows string
	)
	category := new(Category)

	rows, _ := bt.ConfigDB.Queryx(fmt.Sprintf("SELECT category_name FROM category WHERE user_id = (SELECT id FROM users_list WHERE chatID = %d)", update.Message.Chat.ID))


	allRows += fmt.Sprintf("Список категорий пользователя: %s\n\n\n", userName)
	// allRows += fmt.Sprintf("Категория: %s\n\n", category.Category_name)
	for rows.Next() {
		err := rows.StructScan(&category)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
			return " ", err
		}

		fmt.Printf("%#v\n", category)
		allRows += fmt.Sprintf("%s\n", category.Category_name)
	}
	return allRows, nil

}