package webWorker

import (
	"fmt"
)

type Portfolio struct {
	AccountName      string
	Timestamp        int64 // unix timestamp
	Positions        []Position
	Balance          CurrentAccount
	Value            float64
	EquityValue      float64
	EquityIssuePrice float64
}

func (p Portfolio) AsDBString(timestamp int64) string {
	return fmt.Sprintf("%s,%f,%f,%f,%f,%f,%f,%d",
		p.AccountName,
		p.Balance.Balance,
		p.Balance.Available,
		p.Balance.AvailableCredit,
		p.Value,
		p.EquityValue,
		p.EquityIssuePrice,
		timestamp)
}
