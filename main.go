package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bayudha2/go-test-0/app"
	"github.com/bayudha2/go-test-0/models"
)

func main() {
	models.ConnectDatabase(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	app.Initialize()
	log.Fatal(http.ListenAndServe(":8010", app.R))
}
