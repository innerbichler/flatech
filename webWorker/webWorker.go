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

type WebWorker struct {
	userId   string
	password string
	driver   selenium.WebDriver
	service  *selenium.Service
}

func NewWebWorker(userId string, password string) WebWorker {
	// starts webdriver and logs in -> returns the driver instance
	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(geckoDriverPath), // Specify GeckoDriver path
	}
	service, err := selenium.NewGeckoDriverService(geckoDriverPath, port, opts...)
	if err != nil {
		log.Fatalf("Error starting the Geckodriver service: %v", err)
	}

	// Connect to the WebDriver instance
	caps := selenium.Capabilities{"browserName": "firefox"}
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		log.Fatalf("Error connecting to the WebDriver: %v", err)
	}

	return WebWorker{
		userId,
		password,
		driver,
		service,
	}
}

func (w WebWorker) Close() {
	w.service.Stop()
	w.driver.Quit()
}

func (w WebWorker) Login() {
	err := w.driver.Get("https://konto.flatex.at/login.at/loginIFrameFormAction.do")
	if err != nil {
		log.Fatalf("Failed to load website: %v", err)
	}
	// acceptCookies(driver)

	// wait for page to be loaded fully !!!
	time.Sleep(2 * time.Second)

	userField, err := w.driver.FindElement(selenium.ByName, "userId")
	if err != nil {
		log.Fatalf("Failed to find the userField: %v", err)
	}
	userField.SendKeys(w.userId)

	passwordField, err := w.driver.FindElement(selenium.ByName, "password")
	if err != nil {
		log.Fatalf("Failed to find the passwordField: %v", err)
	}
	passwordField.SendKeys(w.password)

	time.Sleep(5 * time.Second)

	loginButton, err := w.driver.FindElement(selenium.ByID, "btnSubmitForm")
	if err != nil {
		log.Fatalf("Failed to find the loginButton: %v", err)
	}
	loginButton.Click()
	time.Sleep(1 * time.Second)
}

func (w WebWorker) GetAll() {
	err := w.driver.Get("https://konto.flatex.at/next-desktop.at/overviewFormAction.do")
	if err != nil {
		log.Fatalf("Failed to load website: %v", err)
	}
	time.Sleep(5 * time.Second)
	allButton, err := w.driver.FindElement(selenium.ByID, "__1551384479")
	if err != nil {
		log.Fatalf("Failed to find the element: %v", err)
	}
	allButton.Click()
	time.Sleep(5 * time.Second)

	// page, err := driver.PageSource()
	// log.Println(page)

	names, err := w.driver.FindElements(selenium.ByCSSSelector, "[class=Ellipsis]")
	if err != nil {
		log.Fatalf("Failed to find the element: %v", err)
	}
	amount, err := w.driver.FindElements(selenium.ByCSSSelector, "[class=PositiveAmount]")
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
