package main

import (
	"flatech/webWorker"
	"log"
	"math"
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

	value := 0.0
	for _, pos := range positions {
		value += pos.CurrentPrice
	}
	value = math.Floor(value)
	log.Println("current value of portfolio:", value, "€")
	time.Sleep(10 * time.Second)
	log.Println("webWorker completed")
}
