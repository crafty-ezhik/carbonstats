package carbon

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type CarbonBilling struct {
	servAddr string
	client   *http.Client
	log      *zap.Logger
}

// TODO: После реализации методов, сделать отдачу интерфейса, а не структуры
func NewCarbonBilling(host string, port int, logger *zap.Logger) *CarbonBilling {
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	carbon := &CarbonBilling{
		servAddr: fmt.Sprintf("%s:%d", host, port),
		client:   client,
		log:      logger,
	}
	return carbon
}

func (c *CarbonBilling) CallApi(model string, params []byte) (map[string]any, error) {
	apiUrl := fmt.Sprintf("http://%s/rest_api/v2/%s/", c.servAddr, model)

	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(params))
	if err != nil {
		c.log.Error("Error creating request", zap.Error(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("Error sending request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.log.Error("Bad response status", zap.Int("status", resp.StatusCode))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error("Error reading response", zap.Error(err))
		return nil, err
	}

	fmt.Println(string(respBody))
	return nil, nil
}
