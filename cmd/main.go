package main

import (
	"fmt"
	"github.com/crafty-ezhik/carbonstats/config"
	"github.com/crafty-ezhik/carbonstats/internal/carbon"
)

func main() {
	//logger := logger2.NewLogger(true)
	cfg := config.LoadConfig()
	billing := carbon.NewCarbonBilling(cfg.Carbon.Host, cfg.Carbon.Port)

	data, err := billing.GetAbonentsList([]int{1455, 1456})
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range data.Abonents {
		fmt.Println(item)
	}

}
