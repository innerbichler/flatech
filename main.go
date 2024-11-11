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

	worker.GetAll()
	log.Println("got all for " + userId + " successfully")

	time.Sleep(100 * time.Second)
	log.Println("webWorker completed")
}
