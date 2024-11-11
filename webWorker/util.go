package webWorker

import (
	"strings"

	"github.com/tebeka/selenium"
)

func newPositionFromList(data []selenium.WebElement) []Position {
	positions := []Position{}
	for _, item := range data {
		text, err := item.Text()
		if err == nil {
			splitText := strings.Split(text, "\n")
			positions = append(positions, createNewPositionHelper(splitText))

		}
	}
	return positions
}

func createNewPositionHelper(data []string) Position {
	return Position{
		data[0],
		data[1],
		data[2],
		data[3],
		data[4],
		data[5],
		data[6],
	}
}
