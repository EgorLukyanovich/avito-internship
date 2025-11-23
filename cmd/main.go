package main

import (
	"log"
	"net/http"
	"time"

	handlers "github.com/egor_lukyanovich/avito/internal/handlers"
	"github.com/egor_lukyanovich/avito/internal/routing"
	"github.com/egor_lukyanovich/avito/pkg/app"
)

func main() {
	storage, port, err := app.InitDB()
	if err != nil {
		log.Fatalf("DB initialization failed: %v", err)
	}

	defer storage.DB.Close()

	teamHandlers := handlers.NewTeamHandlers(storage.Queries)
	userHandlers := handlers.NewUserHandlers(storage.Queries)
	pullReqHadnlers := handlers.NewPullRequestHandlers(storage.Queries)
	statsHandler := handlers.NewStatsHandler(storage.Queries)
	router := routing.NewRouter(*storage)

	router.Post("/team/add", teamHandlers.AddTeam)
	router.Get("/team/get", teamHandlers.GetTeam)

	router.Get("/users/get", userHandlers.GetUser)
	router.Get("/users/getReview", userHandlers.GetPRsForReview)
	router.Delete("/users/delete", userHandlers.DeleteUser)
	router.Post("/users/setIsActive", userHandlers.SetUserActive)
	router.Post("/users/upsertUser", userHandlers.UpsertUser)

	router.Post("/pullRequest/create", pullReqHadnlers.CreatePullRequest)
	router.Post("/pullRequest/merge", pullReqHadnlers.MergePullRequest)
	router.Post("/pullRequest/reassign", pullReqHadnlers.ReassignReviewer)

	router.Get("/stats", statsHandler.GetStats)

	server := &http.Server{
		Handler:           router,
		Addr:              port,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on :%s", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to listen to server: %v", err)
	}
}
