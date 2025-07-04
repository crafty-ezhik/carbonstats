package main

import (
	"github.com/crafty-ezhik/carbonstats/config"
	"github.com/crafty-ezhik/carbonstats/internal/carbon"
	"github.com/crafty-ezhik/carbonstats/logger"
)

func main() {
	myLogger := logger.NewLogger(true)
	cfg := config.LoadConfig()
	billing := carbon.NewCarbonBilling(&cfg.Carbon, myLogger)

	billing.Run()

}
