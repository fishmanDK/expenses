package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"

	"github.com/fishmanDK/expenses/pkg/repository"
	telegram "github.com/fishmanDK/expenses/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	if err := initConfig(); err != nil{
		log.Fatalf("problem with config: %s", err.Error())
	}

	if err := gotenv.Load(); err != nil{
		log.Fatalf("problem with '.env': %s", err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("tele_token"))
	if err != nil {
		log.Panicf("error in main.go tgbotapi.NewBotAPI: %s", err.Error())
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить категорию"),
			tgbotapi.NewKeyboardButton("Добавить товар"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить покупку"),
			tgbotapi.NewKeyboardButton("Вывести все категории"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Удалить продукт"),
			tgbotapi.NewKeyboardButton("Удалить категорию"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вывести список продуктов (в категории)"),
		),
	)
	updates := bot.GetUpdatesChan(u)

	configDB, err := repository.NewConfigDB(repository.Config{
		Host: 	  viper.GetString("db.host"),
		Port: 	  viper.GetString("db.port"),
		UserName: viper.GetString("db.userName"),
		Password: os.Getenv("DB_Password"),
		DBName:   viper.GetString("db.dbName"),
		SSLMode:  viper.GetString("db.sslMode"),
	}) 

	if err != nil{
		log.Fatal("No conection DB")
	}
	
	fmt.Println("БД подключенна")

	fmt.Println(configDB)

	telegram.MainReply(*telegram.NewBotConfig(bot, updates, keyboard, configDB))
	
}


func initConfig() error{
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}