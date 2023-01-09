package main

import (
	"github.com/el27egs/twittor-tilotta/db"
	"github.com/el27egs/twittor-tilotta/handlers"
	"log"
)

func main() {
	if db.CheckConnection() == false {
		log.Fatal("Database connection failed")
		return
	}
	handlers.StartServer()
}
