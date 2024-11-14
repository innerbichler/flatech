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
	caps := selenium.Capabilities{
		"browserName": "firefox",
		"moz:firefoxOptions": map[string]interface{}{
			"args": []string{"--headless", "--disable-gpu"},
		},
	}
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

func (w WebWorker) GetAll() []Position {
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

func (w WebWorker) GetPortfolio() {}
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
