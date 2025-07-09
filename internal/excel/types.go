package excel

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	ColorBlack     = "#000000"
	ColorLightGray = "D3D3D3"
)

const (
	BorderLeft   = "left"
	BorderRight  = "right"
	BorderBottom = "bottom"
	BorderTop    = "top"
)

const (
	AlignmentCenter = "center"
	AlignmentLeft   = "left"
	AlignmentRight  = "right"
	AlignmentTop    = "top"
)

type Row struct {
	ClientName               string          `json:"client_name"`
	MinutesCount             decimal.Decimal `json:"minutes_count"`
	MinutesAmountWoTax       decimal.Decimal `json:"minutes_amount_wo_tax"`
	NumbersCount             uint            `json:"numbers_count"`
	ServiceDescription       string          `json:"service_description"`
	ServicesAmountWithoutTax decimal.Decimal `json:"services_amount_without_tax"`
	ServicesAmountWithTax    decimal.Decimal `json:"services_amount_with_tax"`
	TotalAmountWithoutTax    decimal.Decimal `json:"total_amount_without_tax"`
	TotalAmountWithTax       decimal.Decimal `json:"total_amount_with_tax"`
	CompanyAffiliation       string          `json:"company_affiliation"`
	DocNumber                int             `json:"doc_number"`
	VPBXAmountWithTax        decimal.Decimal `json:"vpbx_amount_with_tax"`
	AmountFromBLToBI         decimal.Decimal `json:"amount_from_bl_to_bi"`
	CallsCount               int             `json:"calls_count"`
}

type CompanyData struct {
	Data                  []Row           `json:"data"`
	SumMinutesCount       decimal.Decimal `json:"sum_minutes_count"`
	SumNumbersCount       uint            `json:"sum_numbers_count"`
	SumAdditionalServices decimal.Decimal `json:"sum_additional_services"`
	SumTotalAmountWithTax decimal.Decimal `json:"sum_total_amount_with_tax"`
	SumCallsCount         int             `json:"sum_calls_count"`
}

type Rows struct {
	BL    CompanyData `json:"bl"`
	BI    CompanyData `json:"bi"`
	Month time.Month  `json:"month"`
	Year  int         `json:"year"`
}

// Flatten - возвращает структуру в виде []any в том же порядке, что и определены поля
//
//	Используется при добавлении строки в Excel
func (r Row) Flatten() []any {
	return []any{
		r.ClientName,
		r.MinutesCount,
		r.MinutesAmountWoTax,
		r.NumbersCount,
		r.ServiceDescription,
		r.ServicesAmountWithoutTax,
		r.ServicesAmountWithTax,
		r.TotalAmountWithoutTax,
		r.TotalAmountWithTax,
		r.CompanyAffiliation,
		r.DocNumber,
		r.VPBXAmountWithTax,
		r.AmountFromBLToBI,
		r.CallsCount,
	}
}
