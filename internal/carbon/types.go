package carbon

import "github.com/shopspring/decimal"

const (
	MethodOne   = "method1"
	ArgOne      = "arg1"
	MethodTwo   = "method2"
	ArgTwo      = "arg2"
	MethodThree = "method3"
	ArgThree    = "arg3"
	Fields      = "fields"
)

const (
	MethodGetDetails = "get_details"
)

const (
	ObjFilter = "objects.filter"
	ObjAll    = "objects.all"
	ObjGet    = "objects.get"
)

const (
	ModelAbonents         = "Abonents"
	ModelFinanceOperation = "FinanceOperations"
)

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
)

type AnyMap map[string]interface{}

type Pair struct {
	Key   string
	Value any
}

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
	Error  string          `json:"error"`
}

type ResponseWithManyRes2 struct {
	Call   string          `json:"call"`
	Result []ResultRequest `json:"result"`
	Error  string          `json:"error"`
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
