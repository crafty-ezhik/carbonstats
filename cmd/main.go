package main

import (
	"fmt"
	"github.com/crafty-ezhik/carbonstats/config"
	"github.com/crafty-ezhik/carbonstats/internal/carbon"
	"github.com/crafty-ezhik/carbonstats/internal/db"
	"github.com/crafty-ezhik/carbonstats/internal/routes"
	"github.com/crafty-ezhik/carbonstats/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func main() {
	myLogger := logger.NewLogger(true)
	cfg := config.LoadConfig()
	_ = carbon.NewCarbonBilling(&cfg.Carbon, myLogger)

	database := db.GetConnection(&cfg.DB)
	db.GoMigrate(database)

	//billing.Run()

	// Инициализация роутера, middlewares, маршрутов
	router := chi.NewRouter()

	routes.InitRoutes(router)
	routes.InitMiddleware(router, cfg.Server.Timeout)

	// Кофигурирование сервера
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Старт сервера
	myLogger.Info("Starting proxy server on port: " + strconv.Itoa(cfg.Server.Port))
	err := server.ListenAndServe()
	if err != nil {
		myLogger.Error("Error starting server.")
		panic(err)
	}
}
