package main

import (
	"encoding/json"
	"fmt"
	"github.com/crafty-ezhik/carbonstats/internal/carbon"
	logger2 "github.com/crafty-ezhik/carbonstats/logger"
	"net/url"
)

type Response struct {
	Call   string           `json:"call"`
	Result []map[string]any `json:"result"`
}

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
	fmt.Println(string(resp))

	var data Response
	err := json.Unmarshal(resp, &data)
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range data.Result[0] {
		fmt.Println(item)
		if items, ok := item.(map[string]any); ok {
			for key, value := range items {
				fmt.Printf("Поле: %s, Значение: %v\n", key, value)
			}
		}
	}
}
