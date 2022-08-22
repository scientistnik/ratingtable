package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"ratingtable/internal/app"
	"ratingtable/internal/storage"
	"ratingtable/internal/telegram"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	telegramToken := os.Getenv("TELEGRAM_TOKEN")

	store, err := storage.NewSQLiteStorage("database.db")
	if err != nil {
		fmt.Printf("error in create store instance: %s", err.Error())
		return
	}

	err = store.Connect()
	if err != nil {
		fmt.Printf("error in sqlite connection: %s", err.Error())
		return
	}

	application := app.NewApp(store)

	links := map[string]string{
		"telegram": "1",
	}

	user, err := application.GetOrCreateUser(links)
	if err != nil {
		fmt.Printf("error in create user: %s", err.Error())
		return
	}

	fmt.Printf("%#v", user)

	telegram.NewTelegram(telegramToken).Launch(ctx, application)

	cancelCannel := make(chan os.Signal, 1)
	signal.Notify(cancelCannel, os.Interrupt)

	<-cancelCannel

	cancel()
}
