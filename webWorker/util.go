package webWorker

import (
	"strconv"
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
		formatAmount(data[1]),
		formatCurrentPrice(data[2]),
		formatCurrentPrice(data[3]),
		formatCurrentPrice(data[4]),
		formatCurrentPrice(data[5]),
		formatCurrentPrice(data[6]),
	}
}

func formatAmount(data string) string {
	return strings.Split(data, " ")[1]
}

func formatCurrentPrice(data string) float64 {
	number := strings.Split(data, " ")[0]

	noDot := strings.Replace(number, ".", "", 1)
	noComa := strings.Replace(noDot, ",", ".", 1)
	finished, err := strconv.ParseFloat(noComa, 64)
	if err != nil {
		panic(err)
	}
	return finished
}
