package main

import (
	"log"
	"net/http"

	"github.com/egor_lukyanovich/avito/internal/routing"
	"github.com/egor_lukyanovich/avito/pkg/app"
)

func main() {
	storage, err := app.InitDB()
	if err != nil {
		log.Fatalf("DB initialization failed: %v", err)
	}

	defer storage.DB.Close()

	router := routing.NewRouter(*storage)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
