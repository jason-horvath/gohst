package main

import (
	"log"
	"net/http"
	"strconv"

	"gohst/app/config"
	appHelpers "gohst/app/helpers"
	appRoutes "gohst/app/routes"
	"gohst/internal/render"
	"gohst/internal/routes"

	coreConfig "gohst/internal/config"
	"gohst/internal/db"
	"gohst/internal/session"
)

func main() {
	defer func() {
        if r := recover(); r != nil {
            log.Println("Recovered from panic:", r)
        }
    }()

	coreConfig.RegisterAppConfig(config.InitAppConfig())
	coreConfig.InitConfig()    // Initialize app-specific config

	dbConfigs := config.CreateDBConfigs()   // Initialize database configurations
	session.Init()
	db.InitDBPool(dbConfigs) // Initialize database connections
	defer db.CloseDBPool()

	// App-specific setup
    render.RegisterTemplateFuncs(appHelpers.AppTemplateFuncs())

	if config.App.IsDevelopment() {

		log.Println("appConfig.App:", config.App)
		log.Println("config.Vite:", coreConfig.Vite)
		log.Println("config.DB:", coreConfig.DB)
	}

	appRouter := appRoutes.NewAppRouter()
	mux := routes.RegisterRouter(appRouter)
	port := strconv.Itoa(config.App.Port)
	server := http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Println("Starting server on port:" , port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
