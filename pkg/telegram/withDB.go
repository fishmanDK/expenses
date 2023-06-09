package telegram

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type Category struct{
	ID 			  int     `db:"id"`
	Category_name string  `db:"category_name"`
}


type Product struct{
	ID 			 string  `db:"id"`
	Category_id  int  `db:"category_id"`
	User_id 	 string  `db:"user_id"`
	Product_name string  `db:"product_name"`
	Price		 float32 `db:"price"`
	Count 		 string  `db:"count"`
}

type User struct{
	FirstName_lastNAme string `db:"firstName_lastNAme"`
	ChatID             string `db:"id"`
}

func examUserSQL(bt BotConfig, chatID int) bool{
	user := new(User)
	rows, _ := bt.ConfigDB.Queryx(fmt.Sprintf("select id from users_list where chatID = %d", chatID))
	for rows.Next() {
		err := rows.StructScan(&user)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`product9`: %s", err.Error())
		}
	}
	// fmt.Println(rows)

	return user.ChatID == ""
}

func registrationSQL(bt BotConfig, firstName_lastNAme string, chatID int) bool{
	rows, _ := bt.ConfigDB.Queryx(fmt.Sprintf("insert into users_list (firstName_lastNAme, chatID) values ('%s', %d)", firstName_lastNAme, chatID))
	fmt.Println(rows.Close())

	return true
}

func getProductsSQL(bt BotConfig, update tgbotapi.Update, category_id int) (string, []string, error){
	var (
		allRows string
		sum float32
		products []string
	)
	product := Product{}


	allRows += fmt.Sprintf("Список товаров пользователя: %s\n\n\n", update.Message.Chat.UserName)
	allRows += fmt.Sprintf("Категория: %s\n\n", update.Message.Text)

	rows, _ := bt.ConfigDB.Queryx(fmt.Sprintf("SELECT * FROM product WHERE product.category_id = %d", category_id))
	for rows.Next() {
		err := rows.StructScan(&product)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`product2`: %s", err.Error())
		}
		allRows += fmt.Sprintf("%s: %.2fp.  x   %sшт.\n", product.Product_name, product.Price, product.Count)
		floatCount, _ := strconv.Atoi(product.Count)
		sum += sumPrice(product.Price * float32(floatCount))

		products = append(products, product.Product_name)
	}
	allRows += fmt.Sprintf("\nОбщая сумма: %.2f", sum)
	return allRows, products, nil
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
			log.Fatalf("Error in SELECT in table:`product3`: %s", err.Error())
		}
	}
}


func addCategorySQL(bt BotConfig, user_id int64, category_name string,  update tgbotapi.Update){
	category := new(Category)
	rows, _ := bt.ConfigDB.Queryx(fmt.Sprintf("insert into category (user_id, category_name) values (%d, '%s')", user_id, category_name))

	for rows.Next() {
		err := rows.StructScan(&category)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`product4`: %s", err.Error())
		}
	}
}

func delCategorySQL(bt BotConfig, chatID int, nameProduct string) bool{
	category := new(Category)
	rows, err := bt.ConfigDB.Queryx(fmt.Sprintf("delete from category where user_id = %d and category_name = '%s'", chatID, nameProduct))
	if err != nil{
		return false
	}
	for rows.Next() {
		err := rows.StructScan(&category)
		if err != nil {
			log.Fatalf("error in ifDelCategorySQL: %s", err.Error())
		}
	}

	return true
}

func delProductSQL(bt BotConfig, chatID int, categoryID int, productName string) bool{
	category := new(Category)
	rows, err := bt.ConfigDB.Queryx(fmt.Sprintf("delete from product where user_id = %d and category_id = %d and product_name = '%s'", chatID, categoryID, productName))
	if err != nil{
		return false
	}
	for rows.Next() {
		err := rows.StructScan(&category)
		if err != nil {
			log.Fatalf("error in ifDelCategorySQL: %s", err.Error())
		}
	}

	return true
}


func slqCountProductbt(bt BotConfig, newValue int, user_id int64, product_name string){
	category := new(Category)
	rows, _ := bt.ConfigDB.Queryx(fmt.Sprintf("UPDATE product SET count = count + %d WHERE user_id = %d and product_name = '%s'", newValue, user_id, product_name))
	for rows.Next() {
		err := rows.StructScan(&category)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`slqCountProductbt`: %s", err.Error())
		}
	}
}

func printCategoryesSQL(bt BotConfig, userName string, update tgbotapi.Update)(map[string]int, error){

	cagegoryes := make(map[string]int)
	category := new(Category)

	rows, err := bt.ConfigDB.Queryx(fmt.Sprintf("SELECT category_name, id FROM category WHERE (select chatID from users_list where chatID = %d) = %d", update.Message.Chat.ID, update.Message.Chat.ID))
	if err != nil {
		log.Fatalf("Error in SELECT in table:`category`: %s", err.Error())
	}

	for rows.Next() {
		err := rows.StructScan(&category)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`product5`: %s", err.Error())
		}

		cagegoryes[category.Category_name] = category.ID

	}
	return cagegoryes, nil

}

func printCategoryesInMainReply(bt BotConfig, update tgbotapi.Update) ([]string, []string){
	var (
		products []string
		price []string
		floatPrice float32
	)
	category := new(Category)
	product := Product{}

	rows, err := bt.ConfigDB.Queryx(fmt.Sprintf("SELECT category_name, id FROM category WHERE user_id = %d", update.Message.Chat.ID))
	if err != nil {
		log.Fatalf("Error in SELECT in table:`category`: %s", err.Error())
	}

	for rows.Next() {
		err := rows.StructScan(&category)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
		}
		products = append(products, category.Category_name)
		rows, _ = bt.ConfigDB.Queryx(fmt.Sprintf("SELECT * FROM product WHERE user_id = %d", update.Message.Chat.ID))
		for rows.Next() {
			err := rows.StructScan(&product)
			if err != nil {
				log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
			}
			if category.ID == product.Category_id{
				floatCount, _ := strconv.ParseFloat(product.Count, 32)
				floatPrice += float32(product.Price * float32(floatCount))
			}
		}
		price = append(price, strconv.FormatFloat(float64(floatPrice), 'f', -1, 32))
	}

	return products, price
}
