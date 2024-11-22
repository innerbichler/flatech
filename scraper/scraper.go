package main

import (
	"flag"
	"log"
	"os"
	"time"
	_ "time/tzdata"

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
		if !isMarketOpen() {
			continue
		}
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

func isMarketOpen() bool {
	loc, _ := time.LoadLocation("Europe/Vienna")
	t := time.Now().In(loc)
	switch t.Weekday() {
	case time.Saturday:
		h, m, _ := t.Clock()
		now_minutes := (h * 60) + m
		waiting_minutes := (24 * 60) - now_minutes - 5
		log.Println("todays saturday, do something else - waiting", waiting_minutes, "minutes")
		time.Sleep(time.Duration(waiting_minutes) * time.Hour)
		return false
	case time.Sunday:
		h, m, _ := t.Clock()
		now_minutes := (h * 60) + m
		waiting_minutes := (24 * 60) - now_minutes - 5
		log.Println("todays sunday, do something else - waiting", waiting_minutes, "minutes")
		time.Sleep(time.Duration(waiting_minutes) * time.Hour)
		return false
	case time.Friday:
		h, m, _ := t.Clock()
		if h > 22 {
			log.Println("market closed for the week - waiting for monday")
			time.Sleep(2 * time.Hour)
			return false
		}
		if h < 9 {
			now_minutes := (h * 60) + m
			waiting_minutes := (9 * 60) - now_minutes - 1
			log.Println("market not open yet - waiting for", waiting_minutes, "minutes")
			time.Sleep(time.Duration(waiting_minutes) * time.Hour)
			return false
		}
	default:
		h, m, _ := t.Clock()
		if h > 22 {
			log.Println("market closed for the week - waiting for monday")
			time.Sleep(2 * time.Hour)
			return false
		}
		if h < 9 {
			now_minutes := (h * 60) + m
			waiting_minutes := (9 * 60) - now_minutes - 1
			log.Println("market not open yet - waiting for", waiting_minutes, "minutes")
			time.Sleep(time.Duration(waiting_minutes) * time.Hour)
			return false
		}
	}
	return true
}

func scrapePortfolio(userId string, password string) webWorker.Portfolio {
	worker := webWorker.NewWebWorker(userId, password)
	defer worker.Close()
	worker.Login()
	portfolio := worker.GetPortfolio()
	return portfolio
}
