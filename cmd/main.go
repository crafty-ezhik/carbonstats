package main

import (
	"fmt"
	"github.com/crafty-ezhik/carbonstats/config"
	"github.com/crafty-ezhik/carbonstats/internal/carbon"
	"github.com/crafty-ezhik/carbonstats/internal/db"
	"github.com/crafty-ezhik/carbonstats/internal/excel"
	"github.com/crafty-ezhik/carbonstats/internal/routes"
	"github.com/crafty-ezhik/carbonstats/internal/service_description"
	"github.com/crafty-ezhik/carbonstats/internal/statistics"
	"github.com/crafty-ezhik/carbonstats/internal/stats_data"
	"github.com/crafty-ezhik/carbonstats/internal/utils"
	"github.com/crafty-ezhik/carbonstats/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func main() {
	myLogger := logger.NewLogger(true)
	cfg := config.LoadConfig()
	billing := carbon.NewCarbonBilling(&cfg.Carbon, myLogger)

	database := db.GetConnection(&cfg.DB)
	db.GoMigrate(database)

	// Инициализация репозиториев
	servDescRepo := service_description.NewServiceDescriptionRepository(database, myLogger)
	statsRepo := statistics.NewStatisticsRepository(database, myLogger)

	// Инициализация обработчиков
	servDescHandler := service_description.NewServiceDescriptionHandler(myLogger, servDescRepo)
	statsHandler := statistics.NewStatisticsHandler(statsRepo, myLogger)
	statsDataHandler := stats_data.NewStatsDataHandler(statsRepo, servDescRepo, billing, myLogger)

	// Инициализация роутера, middlewares, маршрутов
	router := chi.NewRouter()

	routes.InitMiddleware(router, cfg.Server.Timeout)
	routes.InitRoutes(router, servDescHandler, statsHandler, statsDataHandler)

	// TODO: Delete later
	excelFile := excel.New(myLogger, "test")
	excelFile.AddData(utils.DataPreparation(billing, servDescRepo, statsRepo, myLogger))
	if err := excelFile.Save(); err != nil {
		myLogger.Error(err.Error())
	}

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
