package main

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
			// do some magic string stuff
			splitText := strings.SplitAfter(text, "Stück")
			firstPart := strings.Split(splitText[0], "\n")
			name := firstPart[0]
			amountString := strings.Split(firstPart[1], " ")
			amount := amountString[len(amountString)-2]
			test := []string{name, amount}
			secondPart := strings.Split(splitText[1], "\n")

			correctedText := append(test, secondPart...)
			positions = append(positions, createNewPositionHelper(correctedText))

		}
	}
	return positions
}

func createNewPositionHelper(data []string) Position {
	amount := formatAmount(data[1])
	currentValue := formatCurrentPrice(data[3])
	currentPrice := currentValue / float64(amount)

	// data[2] is a space

	issueValue := formatCurrentPrice(data[4])
	issuePrice := issueValue / float64(amount)

	name := strings.ReplaceAll(data[0], " ", "-")
	return Position{
		Name:                       name,
		Amount:                     amount,
		CurrentValue:               currentValue,
		CurrentPrice:               currentPrice,
		IssueValue:                 issueValue,
		IssuePrice:                 issuePrice,
		DevelopmentAbsolutePercent: formatCurrentPrice(data[5]),
		ClosingYesterday:           formatCurrentPrice(data[6]),
		DevelopmentToday:           formatCurrentPrice(data[7]),
	}
}

func formatAmount(amount string) int64 {
	finished, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		panic(err)
	}
	return finished
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
