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
	conn, err := GetDatabaseConnection()
	if err != nil {
		log.Fatal("Error connection to db")
	}
	conn.startup()
	defer conn.db.Close()

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	userId := os.Getenv("USERID")
	password := os.Getenv("PASSWORD")

	worker := webWorker.NewWebWorker(userId, password)
	defer worker.Close()
	worker.Login()

	log.Println("logged into ******** successfully")

	portfolio := worker.GetPortfolio()
	log.Println("got portfolio for ******** successfully")

	//for _, item := range portfolio.Positions {
	//	printPosition(item)
	//}
	gain := ((portfolio.EquityValue / portfolio.EquityIssuePrice) - 1) * 100
	value := math.Floor(portfolio.EquityValue*100) / 100
	gain = math.Floor(gain*100) / 100
	log.Println("current value of portfolio:", value, "€")
	log.Println("with a total gain of", gain, "%")
	_, err = conn.InsertPortfolio(portfolio)
	if err != nil {
		log.Fatal("Error inserting portfolio", err)
	}
	time.Sleep(10 * time.Second)
	result, err := conn.SelectAll()
	if err != nil {
		log.Fatal("Error selecting All", err)
	}
	log.Println(result)
	log.Println("webWorker completed")
}

func printPosition(position webWorker.Position) {
	log.Println("name:", position.Name)
	log.Println("amount:", position.Amount)
	log.Println("currentValue:", position.CurrentValue)
	log.Println("currentPrice:", position.CurrentPrice)
	log.Println("issueValue", position.IssueValue)
	log.Println("issuePrice", position.IssuePrice)
	log.Println("closingYesterday", position.ClosingYesterday)
	log.Println("developmentToday", position.DevelopmentToday)
	log.Println("developmentAbsolute", position.DevelopmentAbsolutePercent)
}
