package main

import (
	"fmt"
	"github.com/crafty-ezhik/carbonstats/internal/carbon"
)

func main() {
	//logger := logger2.NewLogger(true)
	billing := carbon.NewCarbonBilling("78.155.208.228", 8082)

	data, err := billing.GetAbonentsList([]int{1455, 1456})
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range data.Abonents {
		fmt.Println(item)
	}

}
