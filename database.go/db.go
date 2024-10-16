package database

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

const (
	dbUser     = "postgres"
	dbPassword = "vishal"
	// dbName     = "view_count"
	dbHost = "localhost"
	dbPort = 5433
)

func Connect(dbName string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	// remove all log.Println from code. use only mw
	return db, err
}

func CreateDB(dbName string) error {
	db, err := Connect("")
	if err != nil {
		return err
	}
	defer db.Close()
	// _, err = db.Exec("CREATE DATABASE view_count")
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(dbName)))
	if err != nil && err.(*pq.Error).Code != "42P04" {
		return err
	}

	dbNew, err := Connect(dbName)
	if err != nil {
		return err
	}
	defer dbNew.Close()
	_, err = dbNew.Exec(`
        CREATE TABLE IF NOT EXISTS videos (
            id TEXT PRIMARY KEY,
            views INT NOT NULL,
            last_updated TIMESTAMP NOT NULL
        );
    `)
	if err != nil {
		return err
	}
	return nil
}
