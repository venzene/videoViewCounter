package database

import (
	"database/sql"
	"fmt"
	"log"

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

func CreateDB(dbName string) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%d sslmode=disable", dbUser, dbPassword, dbHost, dbPort))
	if err != nil {
		log.Fatal("Error connecting to PostgreSQL server: ", err)
	}
	defer db.Close()
	// _, err = db.Exec("CREATE DATABASE view_count")
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(dbName)))
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "42P04" {
			log.Println("Database already exists.")
		} else {
			log.Fatalf("Error creating database: %v", err)
		}
	} else {
		log.Println("Database created successfully.")
	}
	dbNew, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort))
	if err != nil {
		log.Fatal("Error connecting to the new PostgreSQL database: ", err)
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
		log.Fatalf("Error creating table schema: %v", err)
	} else {
		log.Println("Table 'videos' created successfully.")
	}
}
