package infra

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

const dbAddress = "./db/bid_data.sql"

func OpenDatabase() (*sql.DB, error) {
	exists := isDBExists(dbAddress)
	var firstConnection bool
	if !exists {
		firstConnection = true
		log.Println("Creating database")
		err := createDatabase()
		if err != nil {
			return nil, err
		}
	}
	log.Println("Opening database")
	sqlite3, err := open(firstConnection)
	if err != nil {
		return nil, err
	}
	return sqlite3, err
}

func open(firstConnection bool) (*sql.DB, error) {
	sqlite3, err := sql.Open("sqlite3", dbAddress)
	if err != nil {
		return nil, err
	}
	if firstConnection {
		log.Println("Creating BID table")
		err = createTable(sqlite3)
		if err != nil {
			return nil, err
		}
	}
	return sqlite3, nil
}

func createDatabase() error {
	dbFile, err := os.Create(dbAddress)
	if err != nil {
		return err
	}
	dbFile.Close()

	return nil
}

func createTable(db *sql.DB) error {
	bidTableSQL := `CREATE TABLE bid ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "value" TEXT);`
	stmt, err := db.Prepare(bidTableSQL)
	if err != nil {

	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func isDBExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
