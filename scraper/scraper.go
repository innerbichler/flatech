package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/innerbichler/flatech/webWorker"
	"github.com/joho/godotenv"
)

func main() {
	filePtr := flag.String("file", ".secrets", "path to your .env file, that holds USERID and PASSWORD")
	dbPtr := flag.String("database", "./test", "path to your database file")
	minutePtr := flag.Int("time", 10, "time between scrapes in minutes")

	flag.Parse()

	log.Println("+++++++++++++++++++++++++++++++")
	log.Println("Starting scraper for flatech with options:")
	log.Println("file", *filePtr)
	log.Println("database", *dbPtr)
	log.Println("time", *minutePtr, "min")
	log.Println("+++++++++++++++++++++++++++++++")
	webWorker.SWAG()

	err := godotenv.Load(*filePtr)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	userId := os.Getenv("USERID")
	password := os.Getenv("PASSWORD")

	conn, err := webWorker.GetDatabaseConnection(*dbPtr)
	if err != nil {
		log.Fatal("Error connecting to db")
	}

	conn.Startup()
	defer conn.Connection.Close()

	for {
		portfolio := scrapePortfolio(userId, password)
		conn, err := webWorker.GetDatabaseConnection(*dbPtr)
		if err != nil {
			log.Fatal("Error connecting to db")
		}

		conn.Startup()
		conn.InsertPortfolio(portfolio)
		log.Println("scraped portfolio for ******** successfully:", portfolio.Value, "€")
		time.Sleep(time.Duration(*minutePtr) * time.Minute)
	}
}

func scrapePortfolio(userId string, password string) webWorker.Portfolio {
	worker := webWorker.NewWebWorker(userId, password)
	defer worker.Close()
	worker.Login()
	portfolio := worker.GetPortfolio()
	return portfolio
}
