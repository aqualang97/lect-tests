package server

import (
	"auth/config"
	"auth/server/handlers"
	"log"
	"net/http"
)

func Start(cfg *config.Config) {
	authHandler := handlers.NewAuthHandler(cfg)
	userHandler := handlers.NewUserHandler(cfg)

	http.HandleFunc("/login", authHandler.Login)
	http.HandleFunc("/profile", userHandler.GetProfile)

	log.Fatal(http.ListenAndServe(cfg.Port, nil))
}
