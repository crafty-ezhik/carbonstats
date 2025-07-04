package carbon

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crafty-ezhik/carbonstats/config"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CarbonBilling struct {
	abonentsList *AbonentsInfoList
	servAddr     string
	pastDate     time.Time
	currentDate  time.Time
	client       *http.Client
	log          *zap.Logger
}

// TODO: После реализации методов, сделать отдачу интерфейса, а не структуры
func NewCarbonBilling(carbonCfg *config.CarbonConfig, logger *zap.Logger) *CarbonBilling {
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	carbon := &CarbonBilling{
		servAddr: fmt.Sprintf("%s:%d", carbonCfg.Host, carbonCfg.Port),
		client:   client,
		log:      logger,
	}

	abonList, err := carbon.getAbonentsList(carbonCfg.Parents)
	if err != nil {
		panic("Failed to create abonents list")
	}
	carbon.abonentsList = abonList

	return carbon
}

func (c *CarbonBilling) getAbonentDocument() {

}

func (c *CarbonBilling) getAbonentsList(parents []string) (*AbonentsInfoList, error) {
	data := RequestParams{
		Method1: ObjFilter,
		Arg1: []Pair{
			{
				Key:   "parent__range",
				Value: parents,
			},
		},
		Fields: []string{FieldName, FiledEmail, FieldOperatorID, FieldParentID, FieldContractNumber},
	}

	formData, err := c.buildFormData(data)
	if err != nil {
		//c.log.Error("Error creating formData", zap.Error(err))
		return nil, err
	}

	resp, err := c.callApi(ModelAbonents, []byte(formData.Encode()))
	if err != nil {
		//c.log.Error("Error sent request to CarbonBilling", zap.Error(err))
		return nil, err
	}

	var respData ResponseWithManyRes
	err = json.Unmarshal(resp, &respData)
	if err != nil {
		//c.log.Error("Error unmarshalling response", zap.Error(err))
		return nil, err
	}
	if respData.Error != "" {
		return nil, errors.New(respData.Error)
	}

	var output AbonentsInfoList
	for _, abonent := range respData.Result {
		fields := abonent["fields"].(map[string]any)
		abonInfo := AbonentInfo{
			PK:             int(abonent["pk"].(float64)),
			Name:           fields["name"].(string),
			ContractNumber: fields["contract_number"].(string),
			Email:          fields["email"].(string),
			OperatorID:     int(fields["operator_id"].(float64)),
			ParentID:       int(fields["parent_id"].(float64)),
		}
		output.Abonents = append(output.Abonents, abonInfo)
	}

	return &output, nil
}

func (c *CarbonBilling) callApi(model string, params []byte) ([]byte, error) {
	apiUrl := fmt.Sprintf("http://%s/rest_api/v2/%s/", c.servAddr, model)

	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(params))
	if err != nil {
		//c.log.Error("Error creating request", zap.Error(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		//c.log.Error("Error sending request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		//c.log.Error("Bad response status", zap.Int("status", resp.StatusCode))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		//c.log.Error("Error reading response", zap.Error(err))
		return nil, err
	}

	fmt.Println(string(respBody))
	return respBody, nil
}

func (c *CarbonBilling) buildFormData(params RequestParams) (url.Values, error) {
	formData := url.Values{}

	formData.Add(MethodOne, params.Method1)
	if err := c.addArgs(&formData, params.Arg1, ArgOne); err != nil {
		return nil, err
	}
	if params.Method2 != "" {
		formData.Add(MethodTwo, params.Method2)
		if err := c.addArgs(&formData, params.Arg2, ArgTwo); err != nil {
			return nil, err
		}
	}
	if params.Method3 != "" {
		formData.Add(MethodThree, params.Method3)
		if err := c.addArgs(&formData, params.Arg3, ArgThree); err != nil {
			return nil, err
		}
	}
	if params.Fields != nil {
		formData.Add(Fields, "["+strings.Join(params.Fields, ",")+"]")
	}

	return formData, nil
}

func (c *CarbonBilling) addArgs(formData *url.Values, args []Pair, argsNumber string) error {
	if len(args) < 1 {
		return nil
	}
	tempArgs := make(AnyMap)
	for _, pair := range args {
		tempArgs[pair.Key] = pair.Value
	}
	argsJson, err := json.Marshal(tempArgs)
	if err != nil {
		return err
	}
	formData.Add(argsNumber, string(argsJson))
	return nil
}

func (c *CarbonBilling) PrintAbonentsList() {
	for _, item := range c.abonentsList.Abonents {
		fmt.Println(item)
	}
}
