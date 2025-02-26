package main

import (
	"log"
	"net/http"

	"gohst/internal/config"
	"gohst/internal/db"
	"gohst/internal/routes"
	"gohst/internal/session"
)

func main() {
	config.InitConfig()
	sm := session.NewSessionManager()
	db.InitDB()
	defer db.CloseDB()

	log.Println("From Config APP_ENV_KEY:", config.App.EnvKey)
	log.Println("config.App:", config.App)
	log.Println("config.Vite:", config.Vite)
	log.Println("config.DB:", config.DB)

	rc := routes.RouteConfig{
		SessionManager: sm,
	}
	mux := routes.SetupRoutes(rc)
	port := config.App.PortStr()
	server := http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Println("Starting server on port:" , port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
