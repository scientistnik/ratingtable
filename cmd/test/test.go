package main

import (
	"fmt"
	"ratingtable/internal/app"
	"ratingtable/internal/storage"
)

func main() {
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
}
