package scraper

import (
	"log"
	"math"
	"os"

	"github.com/innerbichler/flatech/webWorker"

	"github.com/joho/godotenv"
)

func main() {
}

func scrapePortfolio() webWorker.Portfolio {
	err := godotenv.Load()
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

	gain := ((portfolio.EquityValue / portfolio.EquityIssuePrice) - 1) * 100
	value := math.Floor(portfolio.EquityValue*100) / 100
	gain = math.Floor(gain*100) / 100

	log.Println("current value of portfolio:", value, "€")
	log.Println("with a total gain of", gain, "%")
	log.Println("webWorker completed")
	return portfolio
}
