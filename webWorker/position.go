package main

import (
	"fmt"
)

type Position struct {
	Name                       string
	Amount                     int64
	CurrentValue               float64
	CurrentPrice               float64
	IssueValue                 float64
	IssuePrice                 float64
	DevelopmentAbsolutePercent float64
	ClosingYesterday           float64
	DevelopmentToday           float64
}

func (p Position) AsDBString(timestamp int64) string {
	return fmt.Sprintf("'%s',%d,%f,%f,%f,%f,%f,%f,%f,%d",
		p.Name,
		p.Amount,
		p.CurrentValue,
		p.CurrentPrice,
		p.IssueValue,
		p.IssuePrice,
		p.DevelopmentAbsolutePercent,
		p.ClosingYesterday,
		p.DevelopmentToday,
		timestamp)
}
