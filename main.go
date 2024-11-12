package main

import (
	"flatech/webWorker"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	userId := os.Getenv("USERID")
	password := os.Getenv("PASSWORD")

	worker := webWorker.NewWebWorker(userId, password)
	defer worker.Close()
	worker.Login()
	log.Println("logged into " + userId + " successfully")

	positions := worker.GetAll()
	log.Println("got all for " + userId + " successfully")

	for _, pos := range positions {
		log.Println(pos)
	}

	time.Sleep(10 * time.Second)
	log.Println("webWorker completed")
}
