package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func DBSet() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:0987@tcp(localhost:3306)/clothing")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// set how much min connection
	db.SetMaxIdleConns(10)
	// set how mmuch max connetion
	db.SetMaxOpenConns(100)
	// set how long connection will lose
	db.SetConnMaxIdleTime(10 * time.Minute)
	// set how long can use
	db.SetConnMaxLifetime(60 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Println("failed to connection to mysql:", err)
		return nil, err
	}

	fmt.Println("successfuly connected to Mysql")
	return db, nil
}

var DB *sql.DB

func init() {
	var err error
	DB, err = DBSet()
	if err != nil {
		log.Fatal(err)
	}
}

func UserData(User string) (*sql.Rows, error) {
	query := "SELECT * FROM " + User

	rows, err := DB.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer DB.Close()
	return rows, nil
}

func ProductData(Product string) (*sql.Rows, error) {
	query := "SELECT * FROM " + Product

	rows, err := DB.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return rows, nil
}
