package main

import (
	"log"
	"net/http"

	appConfig "gohst/app/config"
	appHelpers "gohst/app/helpers"
	appRoutes "gohst/app/routes"
	"gohst/internal/render"
	"gohst/internal/routes"

	"gohst/internal/config"
	"gohst/internal/db"
	"gohst/internal/session"
)

func main() {
	defer func() {
        if r := recover(); r != nil {
            log.Println("Recovered from panic:", r)
        }
    }()

	config.InitConfig()
	appConfig.Initialize()      // Initialize app-specific config
	dbConfigs := appConfig.CreateDBConfigs()   // Initialize database configurations
	session.Init()
	db.InitDBPool(dbConfigs) // Initialize database connections
	defer db.CloseDBPool()

	// App-specific setup
    render.RegisterTemplateFuncs(appHelpers.AppTemplateFuncs())

	if config.App.IsDevelopment() {
		log.Println("config.App:", config.App)
		log.Println("appConfig.App:", appConfig.App)
		log.Println("config.Vite:", config.Vite)
		log.Println("config.DB:", config.DB)
	}

	appRouter := appRoutes.NewAppRouter()
	mux := routes.RegisterRouter(appRouter)
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
