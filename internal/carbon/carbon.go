package carbon

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crafty-ezhik/carbonstats/config"
	"github.com/shopspring/decimal"
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
		pastDate: time.Now().AddDate(0, 0, -time.Now().Day()),
	}

	abonList, err := carbon.getAbonentsList(carbonCfg.Parents)
	if err != nil {
		panic("Failed to create abonents list")
	}
	carbon.abonentsList = abonList

	return carbon
}

func (c *CarbonBilling) Run() {
	for _, abonent := range c.abonentsList.Abonents {
		result, err := c.getAbonentDocument(&abonent)
		if err != nil {
			c.log.Error("Failed to get abonent document", zap.Error(err))
		}

		fmt.Printf("Абонент: %s. Документ:%v\n\n", abonent.Name, result) // TODO: Убрать
		// TODO: Убрать
	}
}

func (c *CarbonBilling) getAbonentDocument(abonent *AbonentInfo) (*DocumentInfo, error) {
	data := RequestParams{
		Method1: ObjFilter,
		Arg1: []Pair{
			{
				Key:   "abonent",
				Value: abonent.PK,
			},
			{
				Key:   "op_type",
				Value: 1,
			},
			{
				Key:   "period_end_date",
				Value: c.pastDate.Format(time.DateOnly),
			},
		},
		Fields: []string{FieldOpSumma, FieldNumber},
	}

	formData, err := c.buildFormData(data)
	if err != nil {
		c.log.Error("Error creating formData", zap.Error(err))
		return nil, err
	}

	resp, err := c.callApi(ModelFinanceOperation, []byte(formData.Encode()))
	if err != nil {
		c.log.Error("Error sent request to CarbonBilling", zap.Error(err))
		return nil, err
	}

	var respData ResponseWithManyRes
	err = json.Unmarshal(resp, &respData)
	if err != nil {
		c.log.Error("Error unmarshalling response", zap.Error(err))
		return nil, err
	}
	if respData.Error != "" {
		c.log.Error("Error in the CarbonBilling response ", zap.String("Error", respData.Error))
		return nil, errors.New(respData.Error)
	}

	if len(respData.Result) < 1 {
		c.log.Debug("There are no elements in the response")
		return &DocumentInfo{
			Number: "Нет документа за данный период",
			Amount: decimal.NewFromFloat(0),
		}, nil
	}

	docInfo := respData.Result[0]
	opSumma, err := c.getDocumentAmount(docInfo.PK, abonent.PK)
	if err != nil {
		c.log.Error("Error getting document amount", zap.Error(err))
		return nil, err
	}

	output := &DocumentInfo{
		Number: docInfo.Fields["number"].(string),
		Amount: opSumma,
	}

	return output, nil
}

func (c *CarbonBilling) getDocumentAmount(operationPk, clientPk int) (decimal.Decimal, error) {
	data := RequestParams{
		Method1: ObjGet,
		Arg1: []Pair{
			{
				Key:   "abonent",
				Value: clientPk,
			},
			{
				Key:   "op_type",
				Value: 1,
			},
			{
				Key:   "period_end_date",
				Value: c.pastDate.Format(time.DateOnly),
			},
			{
				Key:   "op_id",
				Value: operationPk,
			},
		},
		Method2: MethodGetDetails,
		Fields:  []string{FieldVolume, FieldPrice},
	}

	formData, err := c.buildFormData(data)
	if err != nil {
		c.log.Error("Error creating formData", zap.Error(err))
		return decimal.Zero, err
	}

	resp, err := c.callApi(ModelFinanceOperation, []byte(formData.Encode()))
	if err != nil {
		c.log.Error("Error sent request to CarbonBilling", zap.Error(err))
		return decimal.Zero, err
	}

	var respData ResponseWithManyRes
	err = json.Unmarshal(resp, &respData)
	if err != nil {
		c.log.Error("Error unmarshalling response", zap.Error(err))
		return decimal.Zero, err
	}
	if respData.Error != "" {
		c.log.Error("Error in the CarbonBilling response ", zap.String("Error", respData.Error))
		return decimal.Zero, errors.New(respData.Error)
	}
	if len(respData.Result) < 1 {
		c.log.Debug("There are no elements in the response")
		return decimal.Zero, nil
	}

	summa := decimal.NewFromFloat(0.0)
	for _, operation := range respData.Result {
		fields := operation.Fields
		price := decimal.NewFromFloat(fields["price"].(float64))
		volume := decimal.NewFromFloat(fields["vv"].(float64))

		summa = summa.Add(price.Mul(volume))
	}
	return summa.Round(2), nil
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
		c.log.Error("Error creating formData", zap.Error(err))
		return nil, err
	}

	resp, err := c.callApi(ModelAbonents, []byte(formData.Encode()))
	if err != nil {
		c.log.Error("Error sent request to CarbonBilling", zap.Error(err))
		return nil, err
	}

	var respData ResponseWithManyRes
	err = json.Unmarshal(resp, &respData)
	if err != nil {
		c.log.Error("Error unmarshalling response", zap.Error(err))
		return nil, err
	}
	if respData.Error != "" {
		return nil, errors.New(respData.Error)
	}

	var output AbonentsInfoList
	for _, abonent := range respData.Result {
		fields := abonent.Fields
		abonInfo := AbonentInfo{
			PK:             abonent.PK,
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

	c.log.Debug("Response from Carbon Billing", zap.String("ResponseBody", string(respBody))) // TODO: Убрать
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
