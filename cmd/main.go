package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/crafty-ezhik/carbonstats/config"
	"github.com/crafty-ezhik/carbonstats/internal/carbon"
	"github.com/crafty-ezhik/carbonstats/internal/db"
	"github.com/crafty-ezhik/carbonstats/internal/periodic_tasks"
	"github.com/crafty-ezhik/carbonstats/internal/routes"
	"github.com/crafty-ezhik/carbonstats/internal/service_description"
	"github.com/crafty-ezhik/carbonstats/internal/statistics"
	"github.com/crafty-ezhik/carbonstats/internal/stats_data"
	"github.com/crafty-ezhik/carbonstats/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
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

	// Кофигурирование сервера
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Объявление контекста и WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Запуск задачи по созданию отчетов
	wg.Add(1)
	go func() {
		err := periodic_tasks.RunMonthlyTask(ctx, &wg, myLogger, &periodic_tasks.CustomerReport{
			CarbonCfg:    &cfg.Carbon,
			ServDescRepo: servDescRepo,
			StatsRepo:    statsRepo,
			Log:          myLogger,
		})
		if err != nil && !errors.Is(err, context.Canceled) {
			myLogger.Warn("periodic_tasks.RunMonthlyTask", zap.Error(err))
		}
	}()

	// Запуск сервера
	wg.Add(1)
	go func() {
		defer wg.Done()
		myLogger.Info("Starting proxy server on port: " + strconv.Itoa(cfg.Server.Port))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			myLogger.Error("Error starting server.", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		myLogger.Info("Received shutdown signal")
		cancel()

		// Отдельный контекст, чтобы http сервер отдал ответы и корректно завершил работу
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutDownTimeout)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			myLogger.Error("Error shutting down server.", zap.Error(err))
		}
		sqlDb, err := database.DB()
		if err != nil {
			myLogger.Error("Error connecting to database.", zap.Error(err))
		}
		if err := sqlDb.Close(); err != nil {
			myLogger.Error("Error closing database connection.", zap.Error(err))
		}
	}()
	wg.Wait()
	myLogger.Info("Server stopped")
}
