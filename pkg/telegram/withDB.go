package telegram

import (
	"fmt"
	"log"
	"strconv"

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
			log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
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
			log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
		}
	}
}


func addCategorySQL(bt BotConfig, user_id int64, category_name string){

	category := new(Category)
	rows, _ := bt.ConfigDB.Queryx(fmt.Sprintf("insert into category (user_id, category_name) values (%d, '%s')", user_id, category_name))

	for rows.Next() {
		err := rows.StructScan(&category)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
		}
	}
}




func slqCountProductbt(bt BotConfig, newValue int, user_id int64, product_name string){
// 		rows, _ := bt.ConfigDB.Queryx(fmt.Sprintf("UPDATE product SET count = count + %d WHERE user_id = %d, product_name = '%s'", countProduct, user_id, product_name))
// 		for rows.Next() {
// 			err := rows.StructScan(&category)
// 			if err != nil {
// 				log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
// 			}
// 		}
// }
	category := new(Category)
	fmt.Println("========")
	// newValueInt, _ := strconv.Atoi(newValue)
	rows, _ := bt.ConfigDB.Queryx(fmt.Sprintf("UPDATE product SET count = count + %d WHERE user_id = %d and product_name = '%s'", newValue, user_id, product_name))
	fmt.Println("========")
	for rows.Next() {
		err := rows.StructScan(&category)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`slqCountProductbt`: %s", err.Error())
		}
	}
}

func printCategoryes(bt BotConfig, userName string, update tgbotapi.Update)(map[string]int, error){

	cagegoryes := make(map[string]int)
	category := new(Category)

	rows, err := bt.ConfigDB.Queryx(fmt.Sprintf("SELECT category_name, id FROM category WHERE user_id = %d", update.Message.Chat.ID))
	if err != nil {
		log.Fatalf("Error in SELECT in table:`category`: %s", err.Error())
	}

	for rows.Next() {
		err := rows.StructScan(&category)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
		}

		cagegoryes[category.Category_name] = category.ID

	}
	return cagegoryes, nil

}

func printCategoryesInMainReply(bt BotConfig, userName string, update tgbotapi.Update) string{

	// cagegoryes := make(map[string]int)

	var (
		products []string
		price []string
	)
	category := new(Category)
	product := Product{}

	var message string
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
		
	}
	rows, _ = bt.ConfigDB.Queryx(fmt.Sprintf("SELECT * FROM product WHERE product.category_id = %d", category.ID))
	for rows.Next() {
		err := rows.StructScan(&product)
		if err != nil {
			log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
		}
		floatCount, _ := strconv.ParseFloat(product.Count, 32)
		floatPrice := strconv.FormatFloat(float64(product.Price * float32(floatCount)), 'f', -1, 32)
		price = append(price, floatPrice)
	}

	for _, product := range products{
		message += product
	}
	return message
}
// func printProducts(bt BotConfig, userName string, update tgbotapi.Update, category_id int) (string, error){
// 	var (
// 		message string
// 	)
// 	cagegoryes := make(map[string]int)
// 	category := new(Category)

// 	rows, err := bt.ConfigDB.Queryx(fmt.Sprintf("SELECT product_name FROM category WHERE category_id = %d", category_id))
// 	if err != nil {
// 		log.Fatalf("Error in SELECT in table:`category`: %s", err.Error())
// 	}

// 	message += fmt.Sprintf("Список категорий пользователя: %s\n\n\n", userName)
// 	for rows.Next() {
// 		err := rows.StructScan(&category)
// 		if err != nil {
// 			log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
// 		}

// 		fmt.Printf("%#v\n", category)
// 		cagegoryes[category.Category_name] = category.ID

// 	}
// 	return cagegoryes, nil

// }

// func getProductsSQL(bt BotConfig, update tgbotapi.Update, category_id int) (string, error){
// 	var (
// 		allRows string
// 		sum float32
// 	)
// 	product := Product{}


// 	allRows += fmt.Sprintf("Список товаров пользователя: %s\n\n\n", update.Message.Chat.UserName)
// 	allRows += fmt.Sprintf("Категория: %s\n\n", update.Message.Text)

// 	rows, _ := bt.ConfigDB.Queryx(fmt.Sprintf("SELECT * FROM product WHERE product.category_id = %d", category_id))
// 	for rows.Next() {
// 		err := rows.StructScan(&product)
// 		if err != nil {
// 			log.Fatalf("Error in SELECT in table:`product`: %s", err.Error())
// 		}
// 		fmt.Printf("%#v\n", product)
// 		allRows += fmt.Sprintf("%s: %.2fp.\n", product.Product_name, product.Price)
// 		fmt.Println("===", product.Product_name)
// 		sum += sumPrice(product.Price)
// 	}
// 	allRows += fmt.Sprintf("\nОбщая сумма: %.2f", sum)
// 	return allRows, nil
// }