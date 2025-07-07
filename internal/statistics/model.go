package statistics

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Statistics struct {
	ID                  uint            `gorm:"primarykey"`
	CarbonPK            uint            `json:"carbon_pk"`
	ClientName          string          `gorm:"unique" json:"client_name"`
	MinutesCount        decimal.Decimal `gorm:"type:numeric(10,2)" json:"minutes_count"`
	MinutesAmountWoTax  decimal.Decimal `gorm:"type:numeric(10,2)" json:"minutes_amount_wo_tax"`
	ServicesAmountWoTaz decimal.Decimal `gorm:"type:numeric(10,2)" json:"services_amount_wo_tax"`
	TotalAmountWoTax    decimal.Decimal `gorm:"type:numeric(10,2)" json:"total_amount_wo_tax"`
	AmountFromBLToBI    decimal.Decimal `gorm:"type:numeric(10,2)" json:"amount_from_bl_to_bi"`
	DocNumber           int             `json:"doc_number"`
	CompanyAffiliation  string          `json:"company_affiliation"`
	CallsCount          int             `json:"calls_count"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
	DeletedAt           gorm.DeletedAt  `gorm:"index" json:"deleted_at"`
}
