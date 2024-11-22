package webWorker

import (
	"database/sql"
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
	timestamp INTEGER NOT NULL,
	portfolio INTEGER NOT NULL,
	FOREIGN KEY(portfolio) REFERENCES portfolio_snapshots(id)
    );`
	return con.Connection.Exec(sql)
}

func (con DBConnection) InsertPortfolio(portfolio Portfolio) (sql.Result, error) {
	timestamp := time.Now().Unix()
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
	portfolio := []Portfolio{}
	rows, err := con.Connection.Query(sql)
	if err != nil {
		return portfolio, err
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
		portfolio = append(portfolio, *c)
	}
	return portfolio, nil
}
