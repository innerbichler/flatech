package webWorker

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type DBConnection struct {
	Connection *sql.DB
}

func (con DBConnection) Startup() {
	_, err := con.createPortfolioSnapshotsTable()
	if err != nil {
		log.Println("error in creating PortfolioSnapshots table")
		log.Println(err)
		return
	}
	_, err = con.createPositionSnapshotsTable()
	if err != nil {
		log.Println("error in creating PositionsSnapshot table")
		log.Println(err)
		return
	}
}

func GetDatabaseConnection(dbPath string) (DBConnection, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Println("error in openning db")
		log.Println(err)
		return DBConnection{}, err
	}
	return DBConnection{db}, nil
}

func (con DBConnection) createPortfolioSnapshotsTable() (sql.Result, error) {
	// table for all positions -> use the id of portfolio to get the correct one
	sql := `CREATE TABLE IF NOT EXISTS portfolio_snapshots (
        id INTEGER PRIMARY KEY,
        accountName TEXT NOT NULL,
        balance REAL NOT NULL,
        available REAL NOT NULL,
        availableCredit REAL NOT NULL,
        value REAL NOT NULL,
        equityValue REAL NOT NULL,
        equityIssuePrice REAL NOT NULL,
	timestamp INTEGER NOT NULL
    );`
	return con.Connection.Exec(sql)
}

func (con DBConnection) createPositionSnapshotsTable() (sql.Result, error) {
	// table for all positions -> use the id of portfolio to get the correct one
	sql := `CREATE TABLE IF NOT EXISTS position_snapshots (
        id INTEGER PRIMARY KEY,
        name     TEXT NOT NULL,
        amount INTEGER NOT NULL,
        currentValue REAL NOT NULL,
        currentPrice REAL NOT NULL,
        issueValue REAL NOT NULL,
        issuePrice REAL NOT NULL,
	developmentAbsolutePercent REAL NOT NULL,
	closingYesterday REAL NOT NULL,
	developmentToday REAL NOT NULL,
	timestamp INTEGER NOT NULL
    );`
	return con.Connection.Exec(sql)
}

func (con DBConnection) InsertPosition(position Position, timestamp int64) (sql.Result, error) {
	sql := `INSERT INTO position_snapshots (
	name,
        amount,
        currentValue,
        currentPrice,
        issueValue,
        issuePrice,
	developmentAbsolutePercent,
	closingYesterday,
	developmentToday,
	timestamp
    ) VALUES (` + position.AsDBString(timestamp) + `);`

	return con.Connection.Exec(sql)
}

func (con DBConnection) InsertPortfolio(portfolio Portfolio) (sql.Result, error) {
	timestamp := time.Now().Unix()
	for _, position := range portfolio.Positions {
		_, err := con.InsertPosition(position, timestamp)
		if err != nil {
			log.Fatalln(err)
		}
	}
	sql := `INSERT INTO portfolio_snapshots (
	accountName,
        balance,
        available,
        availableCredit,
        value,
        equityValue,
        equityIssuePrice,
	timestamp
    ) VALUES (` + portfolio.AsDBString(timestamp) + `);`

	return con.Connection.Exec(sql)
}

func (con DBConnection) SelectAll() ([]Portfolio, error) {
	sql := `SELECT * FROM portfolio_snapshots;`
	portfolios := []Portfolio{}
	rows, err := con.Connection.Query(sql)
	if err != nil {
		return portfolios, err
	}
	defer rows.Close()

	for rows.Next() {
		c := &Portfolio{}
		b := &CurrentAccount{}
		var d *int64
		err := rows.Scan(&d, &c.AccountName, &b.Balance, &b.Available, &b.AvailableCredit, &c.Value, &c.EquityValue, &c.EquityIssuePrice, &c.Timestamp)
		if err != nil {
			return nil, err
		}

		c.Balance = *b
		portfolios = append(portfolios, *c)
	}
	return portfolios, nil
}

func (con DBConnection) SelectPositionsFromPortfolio(portfolio Portfolio) ([]Position, error) {
	timestamp := fmt.Sprintf("%d", portfolio.Timestamp)
	sql := `SELECT * FROM position_snapshots WHERE timestamp=` + timestamp + `;`
	rows, err := con.Connection.Query(sql)
	positions := []Position{}
	if err != nil {
		return positions, err
	}
	defer rows.Close()

	for rows.Next() {
		p := &Position{}
		var d *int64
		var t *int64
		err := rows.Scan(&d, &p.Name, &p.Amount, &p.CurrentValue, &p.CurrentPrice, &p.IssueValue, &p.IssuePrice, &p.DevelopmentAbsolutePercent, &p.ClosingYesterday, &p.DevelopmentToday, &t)
		if err != nil {
			return nil, err
		}
		positions = append(positions, *p)
	}
	return positions, nil
}
