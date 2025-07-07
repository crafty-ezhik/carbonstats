package carbon

import "github.com/shopspring/decimal"

// Параметры для формирования запроса
const (
	MethodOne   = "method1"
	ArgOne      = "arg1"
	MethodTwo   = "method2"
	ArgTwo      = "arg2"
	MethodThree = "method3"
	ArgThree    = "arg3"
	Fields      = "fields"
)

// Дополнительные методы получения данных
const (
	MethodGetDetails = "get_details"
)

// Методы получения данных
const (
	ObjFilter = "objects.filter"
	ObjAll    = "objects.all"
	ObjGet    = "objects.get"
)

// Модели в Carbon Billing к котором делаются запросы
const (
	ModelAbonents         = "Abonents"
	ModelFinanceOperation = "FinanceOperations"
	ModelVoipCounters     = "VoipCounters"
	ModelVoipLog          = "VoipLog"
)

// Поля, которые можно выбрать при запросе
const (
	FieldName           = "\"name\""
	FiledEmail          = "\"email\""
	FieldContractNumber = "\"contract_number\""
	FieldOperatorID     = "\"operator_id\""
	FieldParentID       = "\"parent_id\""

	FieldOpSumma = "\"op_summa\""
	FieldNumber  = "\"number\""

	FieldPrice  = "\"price\""
	FieldVolume = "\"vv\""

	FieldMonth           = "\"month_number\""
	FieldYear            = "\"year_number\""
	FieldOutgoingTraffic = "\"v_out\""
	FieldAmount          = "\"summa\""

	FieldBilled = "\"billed\""
)

type AnyMap map[string]interface{}

// Pair - структура для описания аргументов запроса
type Pair struct {
	Key   string
	Value any
}

// RequestParams - структура для формирования тела запроса
type RequestParams struct {
	Method1 string
	Arg1    []Pair
	Method2 string
	Arg2    []Pair
	Method3 string
	Arg3    []Pair
	Fields  []string
}

type Response struct {
	Call   string        `json:"call"`
	Result ResultRequest `json:"result"`
	Error  string        `json:"error"`
}

type ResponseWithManyRes struct {
	Call   string          `json:"call"`
	Result []ResultRequest `json:"result"`
	Error  []string        `json:"error"`
}

type ResponseWithManyRes2 struct {
	Call   string          `json:"call"`
	Result []ResultRequest `json:"result"`
	Error  []string        `json:"error"`
}

type ResultRequest struct {
	PK     int            `json:"pk"`
	Model  string         `json:"model"`
	Fields map[string]any `json:"fields"`
}

type AbonentsInfoList struct {
	Abonents []AbonentInfo
}

type AbonentInfo struct {
	PK             int    `json:"pk"`
	Name           string `json:"name"`
	ContractNumber string `json:"contract_number"`
	Email          string `json:"email"`
	OperatorID     int    `json:"operator_id"`
	ParentID       int    `json:"parent_id"`
}

type DocumentInfo struct {
	Number string          `json:"number"`
	Amount decimal.Decimal `json:"op_summa"`
}

type MinutesInfo struct {
	Count  decimal.Decimal `json:"count"`
	Amount decimal.Decimal `json:"op_summa"`
}
