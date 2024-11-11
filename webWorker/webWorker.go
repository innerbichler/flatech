package webWorker

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

const (
	seleniumPath    = ""                      // Path to Selenium server jar (if using a standalone server, otherwise leave empty)
	geckoDriverPath = "/snap/bin/geckodriver" // Path to geckodriver (adjust if necessary)
	port            = 4444                    // Port for the WebDriver server
)

func Run(userId string, password string) {
	// Start a WebDriver server instance
	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(geckoDriverPath), // Specify GeckoDriver path
	}
	service, err := selenium.NewGeckoDriverService(geckoDriverPath, port, opts...)
	if err != nil {
		log.Fatalf("Error starting the Geckodriver service: %v", err)
	}
	defer service.Stop()

	// Connect to the WebDriver instance
	caps := selenium.Capabilities{"browserName": "firefox"}
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		log.Fatalf("Error connecting to the WebDriver: %v", err)
	}
	defer driver.Quit()

	login(driver, userId, password)
	log.Println("logged into " + userId + " successfully")

	time.Sleep(1 * time.Second)
	getAll(driver)
	log.Println("got all for " + userId + " successfully")
	time.Sleep(100 * time.Second)

	log.Println("webWorker completed")
}

func login(driver selenium.WebDriver, userId string, password string) {
	err := driver.Get("https://konto.flatex.at/login.at/loginIFrameFormAction.do")
	if err != nil {
		log.Fatalf("Failed to load website: %v", err)
	}
	// acceptCookies(driver)

	// wait for page to be loaded fully !!!
	time.Sleep(2 * time.Second)

	userField, err := driver.FindElement(selenium.ByName, "userId")
	if err != nil {
		log.Fatalf("Failed to find the userField: %v", err)
	}
	userField.SendKeys(userId)

	passwordField, err := driver.FindElement(selenium.ByName, "password")
	if err != nil {
		log.Fatalf("Failed to find the passwordField: %v", err)
	}
	passwordField.SendKeys(password)

	time.Sleep(5 * time.Second)

	loginButton, err := driver.FindElement(selenium.ByID, "btnSubmitForm")
	if err != nil {
		log.Fatalf("Failed to find the loginButton: %v", err)
	}
	loginButton.Click()
}

func getAll(driver selenium.WebDriver) {
	err := driver.Get("https://konto.flatex.at/next-desktop.at/overviewFormAction.do")
	if err != nil {
		log.Fatalf("Failed to load website: %v", err)
	}
	time.Sleep(5 * time.Second)
	allButton, err := driver.FindElement(selenium.ByID, "__1551384479")
	if err != nil {
		log.Fatalf("Failed to find the element: %v", err)
	}
	allButton.Click()
	time.Sleep(5 * time.Second)

	// page, err := driver.PageSource()
	// log.Println(page)

	names, err := driver.FindElements(selenium.ByCSSSelector, "[class=Ellipsis]")
	if err != nil {
		log.Fatalf("Failed to find the element: %v", err)
	}
	amount, err := driver.FindElements(selenium.ByCSSSelector, "[class=PositiveAmount]")
	if err != nil {
		log.Fatalf("Failed to find the element: %v", err)
	}
	for _, item := range names {
		log.Println(item.Text())
	}
	for _, item := range amount {
		log.Println(item.Text())
	}

	time.Sleep(1 * time.Second)
}

func acceptCookies(driver selenium.WebDriver) {
	cookieButton, err := driver.FindElement(selenium.ByID, "CybotCookiebotDialogBodyLevelButtonLevelOptinAllowAll")
	if err != nil {
		log.Fatalf("Failed to find the cookieButton: %v", err)
	}
	cookieButton.Click()
}
