package main

import (
	"encoding/json"
	"fmt"
	"github.com/crafty-ezhik/carbonstats/internal/carbon"
	logger2 "github.com/crafty-ezhik/carbonstats/logger"
	"net/url"
)

func main() {
	logger := logger2.NewLogger(true)
	billing := carbon.NewCarbonBilling("host", 0, logger)

	params := map[string]any{
		"parent__range": []int{0, 0},
	}

	paramsJson, _ := json.Marshal(params)

	formData := url.Values{}
	formData.Add("method1", "objects.filter")
	formData.Add("arg1", string(paramsJson))

	resp, _ := billing.CallApi("Abonents", []byte(formData.Encode()))
	fmt.Println(resp)
}
