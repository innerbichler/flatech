package webWorker

import (
	"log"
	"math"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	conn, err := GetDatabaseConnection("./test.db")
	if err != nil {
		log.Fatal("Error connection to db")
	}
	conn.Startup()
	defer conn.Connection.Close()

	err = godotenv.Load("/home/alex/tmp/.secrets")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	result, err := conn.SelectAll()
	if err != nil {
		log.Fatal("Error selecting All", err)
	}
	log.Println(result)
	portfolio := scrapePortfolio()
	for _, item := range portfolio.Positions {
		printPosition(item)
	}
	conn.InsertPortfolio(portfolio)
}

func scrapePortfolio() Portfolio {
	err := godotenv.Load("/home/alex/tmp/.secrets")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	userId := os.Getenv("USERID")
	password := os.Getenv("PASSWORD")

	worker := NewWebWorker(userId, password)
	defer worker.Close()
	worker.Login()

	log.Println("logged into ******** successfully")

	portfolio := worker.GetPortfolio()
	log.Println("got portfolio for ******** successfully")

	gain := ((portfolio.EquityValue / portfolio.EquityIssuePrice) - 1) * 100
	value := math.Floor(portfolio.EquityValue*100) / 100
	gain = math.Floor(gain*100) / 100

	log.Println("current value of portfolio:", value, "€")
	log.Println("with a total gain of", gain, "%")
	log.Println("webWorker completed")
	return portfolio
}

func printPosition(position Position) {
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
