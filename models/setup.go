package models

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func ConnectDatabase(user, password, dbname string) {
	connectionString := fmt.Sprintf("dbname=%s user=%s password=%s host=localhost sslmode=disable", dbname, user, password)

	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
}
