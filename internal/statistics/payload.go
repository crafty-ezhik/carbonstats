package statistics

type CreateData struct {
	Month               uint   `json:"month" validate:"required,gte=1,lte=12"`
	Year                uint   `json:"year" validate:"required,gte=2025,lt=2100"`
	CarbonPK            uint   `json:"carbon_pk" validate:"required,gt=0"`
	ClientName          string `json:"client_name" validate:"required,gt=0"`
	MinutesCount        string `json:"minutes_count" validate:"required,gt=0"`
	MinutesAmountWoTax  string `json:"minutes_amount_wo_tax" validate:"required,gt=0"`
	ServicesAmountWoTaz string `json:"services_amount_wo_tax" validate:"required,gt=0"`
	TotalAmountWoTax    string `json:"total_amount_wo_tax" validate:"required,gt=0"`
	AmountFromBLToBI    string `json:"amount_from_bl_to_bi" validate:"required,gt=0"`
	DocNumber           int    `json:"doc_number" validate:"required,gt=0"`
	CompanyAffiliation  string `json:"company_affiliation" validate:"len=2,required"`
	CallsCount          int    `json:"calls_count" validate:"required,gt=0"`
}

type CreateRequest struct {
	Data []CreateData `json:"data" validate:"required,gt=0,dive,required"`
}
