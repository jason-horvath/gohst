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
	defer func() {
        if r := recover(); r != nil {
            log.Println("Recovered from panic:", r)
        }
    }()

	config.InitConfig()
	session.Init()
	db.InitDB()
	defer db.CloseDB()

	if config.App.IsDevelopment() {
		log.Println("config.App:", config.App)
		log.Println("config.Vite:", config.Vite)
		log.Println("config.DB:", config.DB)
	}

	mux := routes.SetupRoutes()
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
