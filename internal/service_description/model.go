package service_description

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type ServiceDescription struct {
	ID           uint            `gorm:"primarykey"`
	CarbonPK     uint            `json:"carbon_pk"`
	NumbersCount uint            `json:"numbers_count"`
	VPBXAmount   decimal.Decimal `gorm:"type:numeric(10,2)" json:"vpbx_amount"`
	ServiceDesc  string          `json:"service_desc"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"deleted_at"`
}
