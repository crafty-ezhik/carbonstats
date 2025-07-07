package service_description

// CreateRequest - Структура для описания принимаемых параметров при запросе создания
type CreateRequest struct {
	CarbonPK     uint   `json:"carbon_pk" validate:"min=1,gt=0,required"`
	NumbersCount uint   `json:"numbers_count" validate:"gt=0,required"`
	VPBXAmount   string `json:"vpbx_amount" validate:"gt=0,required"`
	ServiceDesc  string `json:"service_desc"`
}

// CreateBatchRequest - Структура для описания принимаемых параметров при запросе создания нескольких записей
type CreateBatchRequest struct {
	Data []CreateRequest `json:"data" validate:"required,gt=0,dive,required"`
}

type UpdateRequest struct {
	NumbersCount uint   `json:"numbers_count" validate:"gt=0,required"`
	VPBXAmount   string `json:"vpbx_amount" validate:"gt=0,required"`
	ServiceDesc  string `json:"service_desc"`
}
