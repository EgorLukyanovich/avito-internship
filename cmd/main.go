package main

import (
	"log"
	"net/http"

	team "github.com/egor_lukyanovich/avito/internal/handlers"
	"github.com/egor_lukyanovich/avito/internal/routing"
	"github.com/egor_lukyanovich/avito/pkg/app"
)

/*
TODO:

	Подумать как автоматизировать поднятие зависимостей
	Перед сдачей не забудь поменять url в goose с localhost на db инчае будут ошибки
*/
func main() {
	storage, port, err := app.InitDB()
	if err != nil {
		log.Fatalf("DB initialization failed: %v", err)
	}

	defer storage.DB.Close()

	teamHandlers := team.NewTeamHandlers(storage.Queries)
	router := routing.NewRouter(*storage)

	router.Post("/team/add", teamHandlers.AddTeam)
	router.Get("/team/get", teamHandlers.GetTeam)

	server := &http.Server{
		Handler: router,
		Addr:    port,
	}

	log.Printf("Starting server on :%s", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to listen to server: %v", err)
	}
}
