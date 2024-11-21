package main

import (
	"database/sql"
	"flatech/webWorker"
	"log"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type DBConnection struct {
	db *sql.DB
}

func (con DBConnection) startup() {
	_, err := createPortfolioSnapshotsTable(con.db)
	if err != nil {
		log.Println("error in creating PortfolioSnapshots table")
		log.Println(err)
		return
	}
	_, err = createPositionSnapshotsTable(con.db)
	if err != nil {
		log.Println("error in creating PositionsSnapshot table")
		log.Println(err)
		return
	}
}

func GetDatabaseConnection() (DBConnection, error) {
	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		log.Println("error in openning db")
		log.Println(err)
		return DBConnection{}, err
	}
	log.Println("Connected to the SQLite database successfully")
	return DBConnection{db}, nil
}

func createPortfolioSnapshotsTable(db *sql.DB) (sql.Result, error) {
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
	return db.Exec(sql)
}

func createPositionSnapshotsTable(db *sql.DB) (sql.Result, error) {
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
	return db.Exec(sql)
}

func (con DBConnection) InsertPortfolio(portfolio webWorker.Portfolio) (sql.Result, error) {
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

	return con.db.Exec(sql)
}

func (con DBConnection) SelectAll() ([]webWorker.Portfolio, error) {
	sql := `SELECT * FROM portfolio_snapshots;`
	portfolio := []webWorker.Portfolio{}
	rows, err := con.db.Query(sql)
	if err != nil {
		return portfolio, err
	}
	defer rows.Close()

	for rows.Next() {
		c := &webWorker.Portfolio{}
		b := &webWorker.CurrentAccount{}
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
