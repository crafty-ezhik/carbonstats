package carbon

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crafty-ezhik/carbonstats/config"
	"github.com/crafty-ezhik/carbonstats/internal/statistics"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CarbonBilling interface {
	calculatingCostAdditionalServices(docInfo *DocumentInfo, minInfo *MinutesInfo) decimal.Decimal
	getOutgoingCallsOverPastPeriod(abonent *AbonentInfo) (int, error)
	getMinutesForPastPeriod(monthNumber, year, clientPk int) (*MinutesInfo, error)
	getAbonentDocument(abonent *AbonentInfo) (*DocumentInfo, error)
	getDocumentAmount(operationPk, clientPk int) (decimal.Decimal, error)
	getAbonentsList(parents []string) (*AbonentsInfoList, error)
	callApi(model string, params []byte) ([]byte, error)
	buildFormData(params RequestParams) (url.Values, error)
	addArgs(formData *url.Values, args []Pair, argsNumber string) error
	StartStatisticsCollection()
}

type CarbonBillingImpl struct {
	abonentsList  *AbonentsInfoList
	servAddr      string
	carbonParents []string
	pastDate      time.Time
	currentDate   time.Time
	client        *http.Client

	log *zap.Logger
}

// TODO: После реализации методов, сделать отдачу интерфейса, а не структуры
func NewCarbonBilling(carbonCfg *config.CarbonConfig, logger *zap.Logger) CarbonBilling {
	client := &http.Client{
		Timeout: time.Second * 60,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	carbon := &CarbonBillingImpl{
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
	carbon.carbonParents = carbonCfg.Parents

	return carbon
}

func (c *CarbonBillingImpl) StartStatisticsCollection() {
	wg := &sync.WaitGroup{}
	resChan := make(chan statistics.ClientStatistics)
	numWorkers := len(c.abonentsList.Abonents)
	month := int(c.pastDate.Month())
	year := c.pastDate.Year()

	// Обновление списка клиентов
	abonList, err := c.getAbonentsList(c.carbonParents)
	if err != nil {
		return
	}
	c.abonentsList = abonList

	// Асинхронно для каждого клиента необходимо получить все данные, а затем отдать готовую структуру с данными
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			abonent := c.abonentsList.Abonents[i]
			docInfo, err := c.getAbonentDocument(&abonent)
			if err != nil {
				return
			}
			minInfo, err := c.getMinutesForPastPeriod(month, year, abonent.PK)
			if err != nil {
				return
			}
			additionalCost := c.calculatingCostAdditionalServices(docInfo, minInfo)
			callsCount, err := c.getOutgoingCallsOverPastPeriod(&abonent)
			if err != nil {
				return
			}

			var company string
			if abonent.OperatorID == 1450 {
				company = "БЛ"
			} else {
				company = "БИ"
			}

			data := statistics.ClientStatistics{
				Month:               uint(month),
				Year:                uint(year),
				CarbonPK:            uint(abonent.PK),
				ClientName:          abonent.Name,
				MinutesCount:        minInfo.Count,
				MinutesAmountWoTax:  minInfo.Amount,
				ServicesAmountWoTaz: additionalCost,
				TotalAmountWoTax:    docInfo.Amount,
				DocNumber:           docInfo.Number,
				CompanyAffiliation:  company,
				CallsCount:          callsCount,
			}
			resChan <- data
		}()
	}

	go func() {
		wg.Wait()
		close(resChan)
	}()

	for data := range resChan {
		fmt.Println(data)
	}
}

func (c *CarbonBillingImpl) calculatingCostAdditionalServices(docInfo *DocumentInfo, minInfo *MinutesInfo) decimal.Decimal {
	return docInfo.Amount.Sub(minInfo.Amount)

}

func (c *CarbonBillingImpl) getOutgoingCallsOverPastPeriod(abonent *AbonentInfo) (int, error) {
	startPastDate := c.pastDate.AddDate(0, 0, -c.pastDate.Day())
	data := RequestParams{
		Method1: ObjFilter,
		Arg1: []Pair{
			{
				Key:   "s_time__range",
				Value: []string{startPastDate.Format(time.DateOnly), c.pastDate.Format(time.DateOnly)},
			},
			{
				Key:   "abonent_id",
				Value: abonent.PK,
			},
			{
				Key:   "billed",
				Value: 1,
			},
			{
				Key:   "duration__range",
				Value: []int{1, 3600},
			},
		},
		Fields: []string{FieldBilled},
	}

	formData, err := c.buildFormData(data)
	if err != nil {
		c.log.Error("Error creating formData", zap.Error(err))
		return 0, err
	}

	resp, err := c.callApi(ModelVoipLog, []byte(formData.Encode()))
	if err != nil {
		c.log.Error("Error sent request to CarbonBillingImpl", zap.Error(err))
		return 0, err
	}

	var respData ResponseWithManyRes
	err = json.Unmarshal(resp, &respData)
	if err != nil {
		c.log.Error("Error unmarshalling response", zap.Error(err))
		return 0, nil
	}
	if len(respData.Error) > 0 {
		c.log.Error("Error in the CarbonBillingImpl response ", zap.String("Error", respData.Error[len(respData.Error)-1]))
		return 0, errors.New(respData.Error[len(respData.Error)-1])
	}

	if len(respData.Result) < 1 {
		c.log.Debug("There are no elements in the response")
		return 0, nil
	}

	return len(respData.Result), nil
}

func (c *CarbonBillingImpl) getMinutesForPastPeriod(monthNumber, year, clientPk int) (*MinutesInfo, error) {
	data := RequestParams{
		Method1: ObjFilter,
		Arg1: []Pair{
			{
				Key:   "abonent_id",
				Value: clientPk,
			},
			{
				Key:   "month_number",
				Value: monthNumber,
			},
			{
				Key:   "year_number",
				Value: year,
			},
		},
		Fields: []string{FieldMonth, FieldYear, FieldAmount, FieldOutgoingTraffic},
	}

	formData, err := c.buildFormData(data)
	if err != nil {
		c.log.Error("Error creating formData", zap.Error(err))
		return &MinutesInfo{
			Count:  decimal.Zero,
			Amount: decimal.Zero,
		}, nil
	}

	resp, err := c.callApi(ModelVoipCounters, []byte(formData.Encode()))
	if err != nil {
		c.log.Error("Error sent request to CarbonBillingImpl", zap.Error(err))
		return &MinutesInfo{
			Count:  decimal.Zero,
			Amount: decimal.Zero,
		}, nil
	}

	var respData ResponseWithManyRes
	err = json.Unmarshal(resp, &respData)
	if err != nil {
		c.log.Error("Error unmarshalling response", zap.Error(err))
		return &MinutesInfo{
			Count:  decimal.Zero,
			Amount: decimal.Zero,
		}, nil
	}
	if len(respData.Error) > 0 {
		c.log.Error("Error in the CarbonBillingImpl response ", zap.String("Error", respData.Error[len(respData.Error)-1]))
		return &MinutesInfo{
			Count:  decimal.Zero,
			Amount: decimal.Zero,
		}, errors.New(respData.Error[len(respData.Error)-1])
	}

	if len(respData.Result) < 1 {
		c.log.Debug("There are no elements in the response")
		return &MinutesInfo{
			Count:  decimal.Zero,
			Amount: decimal.Zero,
		}, nil
	}

	output := MinutesInfo{}
	if len(respData.Result) > 1 {
		for _, item := range respData.Result {
			field := item.Fields
			amount, _ := decimal.NewFromString(field["summa"].(string))
			count, _ := decimal.NewFromString(field["v_out"].(string))

			output.Count = output.Count.Add(count)
			output.Amount = output.Amount.Add(amount)
		}
	} else {
		field := respData.Result[0].Fields
		amount, _ := decimal.NewFromString(field["summa"].(string))
		count, _ := decimal.NewFromString(field["v_out"].(string))

		output.Count = output.Count.Add(count)
		output.Amount = output.Amount.Add(amount)
	}

	output.Count = output.Count.Round(2)
	output.Amount = output.Amount.Round(2)

	return &output, nil
}

func (c *CarbonBillingImpl) getAbonentDocument(abonent *AbonentInfo) (*DocumentInfo, error) {
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
		c.log.Error("Error sent request to CarbonBillingImpl", zap.Error(err))
		return nil, err
	}

	var respData ResponseWithManyRes
	err = json.Unmarshal(resp, &respData)
	if err != nil {
		c.log.Error("Error unmarshalling response", zap.Error(err))
		return nil, err
	}
	if len(respData.Error) > 0 {
		c.log.Error("Error in the CarbonBillingImpl response ", zap.String("Error", respData.Error[len(respData.Error)-1]))
		return nil, errors.New(respData.Error[len(respData.Error)-1])
	}

	if len(respData.Result) < 1 {
		c.log.Debug("There are no elements in the response")
		return &DocumentInfo{
			Number: 0,
			Amount: decimal.NewFromFloat(0),
		}, nil
	}

	docInfo := respData.Result[0]
	opSumma, err := c.getDocumentAmount(docInfo.PK, abonent.PK)
	if err != nil {
		c.log.Error("Error getting document amount", zap.Error(err))
		return nil, err
	}

	docNumber, err := strconv.Atoi(docInfo.Fields["number"].(string))
	if err != nil {
		c.log.Error("Error parsing document number", zap.Error(err))
		return nil, err
	}

	output := &DocumentInfo{
		Number: docNumber,
		Amount: opSumma,
	}

	return output, nil
}

func (c *CarbonBillingImpl) getDocumentAmount(operationPk, clientPk int) (decimal.Decimal, error) {
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
		c.log.Error("Error sent request to CarbonBillingImpl", zap.Error(err))
		return decimal.Zero, err
	}

	var respData ResponseWithManyRes
	err = json.Unmarshal(resp, &respData)
	if err != nil {
		c.log.Error("Error unmarshalling response", zap.Error(err))
		return decimal.Zero, err
	}
	if len(respData.Error) > 0 {
		c.log.Error("Error in the CarbonBillingImpl response ", zap.String("Error", respData.Error[len(respData.Error)-1]))
		return decimal.Zero, errors.New(respData.Error[len(respData.Error)-1])
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

func (c *CarbonBillingImpl) getAbonentsList(parents []string) (*AbonentsInfoList, error) {
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
		c.log.Error("Error sent request to CarbonBillingImpl", zap.Error(err))
		return nil, err
	}

	var respData ResponseWithManyRes
	err = json.Unmarshal(resp, &respData)
	if err != nil {
		c.log.Error("Error unmarshalling response", zap.Error(err))
		return nil, err
	}
	if len(respData.Error) > 0 {
		c.log.Error("Error in the CarbonBillingImpl response ", zap.String("Error", respData.Error[len(respData.Error)-1]))
		return nil, errors.New(respData.Error[len(respData.Error)-1])
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

func (c *CarbonBillingImpl) callApi(model string, params []byte) ([]byte, error) {
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
	c.log.Debug(fmt.Sprintf("Запрос к модели %s, выполнене успешно", model))
	return respBody, nil
}

func (c *CarbonBillingImpl) buildFormData(params RequestParams) (url.Values, error) {
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

func (c *CarbonBillingImpl) addArgs(formData *url.Values, args []Pair, argsNumber string) error {
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

func (c *CarbonBillingImpl) PrintAbonentsList() {
	for _, item := range c.abonentsList.Abonents {
		fmt.Println(item)
	}
}
