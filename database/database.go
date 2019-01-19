package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DB_USER     = "hibk"
	DB_PASSWORD = "password"
	DB_NAME     = "hibk"
	DB_ADDR     = "127.0.0.1:3306"
)

var ErrDatabaseNotResponding = errors.New("database: database not responding")

// Opens the db and returns a db pointer
func open() (*sql.DB, error) {
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		DB_USER, DB_PASSWORD, DB_ADDR, DB_NAME)

	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return nil, err
	}
	// Check if the database is connected
	if err := db.Ping(); err != nil {
		log.Println(err)
		return nil, ErrDatabaseNotResponding
	}
	return db, nil
}

// Calls all init functions
func Init() {
	db, err := open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	initMusic(db) // musicQueries.go
	initUser(db)  //userQueries.go
}

// Handles errors
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
