package webWorker

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
