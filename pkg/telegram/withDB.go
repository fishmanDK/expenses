package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// "github.com/fishmanDK/expenses"
// "github.com/fishmanDK/expenses/pkg/repository"
type Category struct{
	ID 			  int  `db:"id"`
	Category_name string  `db:"category_name"`
}


type Product struct{
	ID 			 string  `db:"id"`
	Category_id  string  `db:"category_id"`
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

func printAllProductsWithSelectDB(bt BotConfig, update tgbotapi.Update) (string, error){
	var (
		allRows string
		sum float32
	)


	allRows += fmt.Sprintf("Список товаров пользователя: %s\n\n\n", update.Message.Chat.UserName)
	allRows += "Категория: КАТЕГОРИЯ\n\n"
	product := Product{}
	rows, err := bt.ConfigDB.Queryx(fmt.Sprintf("SELECT * FROM product WHERE user_id = %d", update.Message.Chat.ID))
	if err != nil{
		log.Fatalf("error in printAllProductsWithSelectDB(rows): %s", err.Error())
	}
	fmt.Println(rows)
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

func addProductSQL(bt BotConfig, category_id int, user_id int64, product_name string, price string){

	category := new(Category)
	rows, _ := bt.ConfigDB.Queryx(fmt.Sprintf("insert into product (category_id, user_id, product_name, price, count) values (%d, %d, '%s', %s, 1)", category_id, user_id, product_name, price))

	for rows.Next() {
		err := rows.StructScan(&category)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
		}
	}
}



func printProducts(bt BotConfig, userName string, update tgbotapi.Update)(map[string]int, error){
	var (
		message string
	)
	cagegoryes := make(map[string]int)
	category := new(Category)

	rows, err := bt.ConfigDB.Queryx(fmt.Sprintf("SELECT category_name, id FROM category WHERE user_id = %d", update.Message.Chat.ID))
	if err != nil {
		log.Fatalf("Error in SELECT in table:`category`: %s", err.Error())
	}

	message += fmt.Sprintf("Список категорий пользователя: %s\n\n\n", userName)
	for rows.Next() {
		err := rows.StructScan(&category)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
		}

		fmt.Printf("%#v\n", category)
		cagegoryes[category.Category_name] = category.ID

	}
	return cagegoryes, nil

}