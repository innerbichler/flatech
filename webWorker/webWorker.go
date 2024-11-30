package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

const (
	seleniumPath = ""
	// geckoDriverPath = "/usr/bin/geckodriver" // Path to geckodriver (adjust if necessary)
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

	caps := selenium.Capabilities{
		"browserName": "firefox",
		"moz:firefoxOptions": map[string]interface{}{
			"args": []string{"--headless", "--disable-gpu"},
		},
	}
	// uncomment below to start in gui mode
	// caps = selenium.Capabilities{"browserName": "firefox"}

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
	w.driver.Quit()
	w.service.Stop()
}

func (w WebWorker) Login() {
	err := w.driver.Get("https://konto.flatex.at/login.at/loginIFrameFormAction.do")
	if err != nil {
		log.Fatalf("Failed to load website: %v", err)
	}
	// acceptCookies(driver)

	// wait for page to be loaded fully !!!
	time.Sleep(5 * time.Second)

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

	time.Sleep(3 * time.Second)

	loginButton, err := w.driver.FindElement(selenium.ByID, "btnSubmitForm")
	if err != nil {
		log.Fatalf("Failed to find the loginButton: %v", err)
	}

	loginButton.Click()
	time.Sleep(4 * time.Second)

	handles, err := w.driver.WindowHandles()

	// close login window and switch to new one that flatex automatically opens
	w.driver.SwitchWindow(handles[len(handles)-1])
	time.Sleep(3 * time.Second)
	w.driver.CloseWindow(handles[0])
}

func (w WebWorker) GetPositions() []Position {
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

	evenDetailsRow, err := w.driver.FindElements(selenium.ByCSSSelector, "[class='I1 Even DetailsAvailable']")
	if err != nil {
		log.Fatalf("Failed to find the element: %v", err)
	}
	oddDetailsRow, err := w.driver.FindElements(selenium.ByCSSSelector, "[class='I1 Odd DetailsAvailable']")
	if err != nil {
		log.Fatalf("Failed to find the element: %v", err)
	}
	oddLastDetailsRow, err := w.driver.FindElements(selenium.ByCSSSelector, "[class='I1 Odd LastRow DetailsAvailable']")
	if err != nil {
		log.Fatalf("Failed to find the element: %v", err)
	}
	positions := []Position{}
	positions = append(positions, newPositionFromList(evenDetailsRow)...)
	positions = append(positions, newPositionFromList(oddDetailsRow)...)
	positions = append(positions, newPositionFromList(oddLastDetailsRow)...)

	time.Sleep(1 * time.Second)
	return positions
}

func (w WebWorker) GetCurrentAccount() CurrentAccount {
	err := w.driver.Get("https://konto.flatex.at/next-desktop.at/overviewFormAction.do")
	if err != nil {
		log.Fatalf("Failed to load website: %v", err)
	}
	time.Sleep(5 * time.Second)
	accountButton, err := w.driver.FindElement(selenium.ByID, "__611878645")
	if err != nil {
		log.Fatalf("Failed to find the element: %v", err)
	}
	accountButton.Click()
	time.Sleep(5 * time.Second)
	balanceCreditCard, err := w.driver.FindElements(selenium.ByCSSSelector, "[class='BalanceCreditAreaEntryValue']")
	if err != nil {
		log.Fatalf("Failed to find the element: %v", err)
	}

	numbers := []float64{}
	for _, item := range balanceCreditCard {
		text, err := item.Text()
		if err == nil {
			splitText := strings.Split(text, "\n")
			numbers = append(numbers, formatCurrentPrice(splitText[0]))
		}
	}

	time.Sleep(2 * time.Second)
	return CurrentAccount{
		Balance:         numbers[0],
		Available:       numbers[1],
		AvailableCredit: numbers[2],
	}
}

func (w WebWorker) GetPortfolio() Portfolio {
	positions := w.GetPositions()
	currentValue := 0.0
	issueValue := 0.0
	for _, pos := range positions {
		currentValue += pos.CurrentValue
		issueValue += pos.IssueValue
	}

	account := w.GetCurrentAccount()

	return Portfolio{
		Timestamp:        time.Now().Unix(),
		AccountName:      w.userId,
		Positions:        positions,
		Balance:          account,
		Value:            currentValue + account.Balance,
		EquityValue:      currentValue,
		EquityIssuePrice: issueValue,
	}
}

func (w WebWorker) GetAccountAssetNames() []string {
	/* returns the names of all the assets in your Account
	 */
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

	nameElements, err := w.driver.FindElements(selenium.ByCSSSelector, "[class=Ellipsis]")
	if err != nil {
		log.Fatalf("Failed to find the element: %v", err)
	}

	names := []string{}
	for _, item := range nameElements {
		name, err := item.Text()
		if err == nil {
			names = append(names, name)
		}
	}
	return names
}

func acceptCookies(driver selenium.WebDriver) {
	cookieButton, err := driver.FindElement(selenium.ByID, "CybotCookiebotDialogBodyLevelButtonLevelOptinAllowAll")
	if err != nil {
		log.Fatalf("Failed to find the cookieButton: %v", err)
	}
	cookieButton.Click()
}
