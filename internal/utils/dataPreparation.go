package utils

import (
	"github.com/crafty-ezhik/carbonstats/internal/excel"
	"github.com/crafty-ezhik/carbonstats/internal/service_description"
	"github.com/crafty-ezhik/carbonstats/internal/statistics"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"time"
)

func DataPreparation(
	sd service_description.ServiceDescriptionRepository,
	carbonData []statistics.ClientStatistics,
	logger *zap.Logger) *excel.Rows {

	tax := decimal.NewFromFloat(1.2)

	output := excel.Rows{
		BI: excel.CompanyData{
			Data:                  nil,
			SumMinutesCount:       decimal.NewFromFloat(0.00),
			SumNumbersCount:       0,
			SumAdditionalServices: decimal.NewFromFloat(0.00),
			SumTotalAmountWithTax: decimal.NewFromFloat(0.00),
			SumCallsCount:         0,
		},
		BL: excel.CompanyData{
			Data:                  nil,
			SumMinutesCount:       decimal.NewFromFloat(0.00),
			SumNumbersCount:       0,
			SumAdditionalServices: decimal.NewFromFloat(0.00),
			SumTotalAmountWithTax: decimal.NewFromFloat(0.00),
			SumCallsCount:         0,
		},
		Month: (time.Now().AddDate(0, 0, -time.Now().Day())).Month(),
		Year:  (time.Now().AddDate(0, 0, -time.Now().Day())).Year(),
	}

	for _, client := range carbonData {
		servDesc, err := sd.Get(client.CarbonPK)
		if err != nil {
			if err.Error() == "record not found" {
				logger.Debug("Запись не найдена")
				servDesc = service_description.ServiceDescription{
					NumbersCount: 0,
					VPBXAmount:   decimal.NewFromFloat(0.00),
					ServiceDesc:  "Запись в базе не найдена",
				}
			} else {
				logger.Error("Ошибка получения данных", zap.Error(err))
				servDesc = service_description.ServiceDescription{
					NumbersCount: 0,
					VPBXAmount:   decimal.NewFromFloat(0.00),
					ServiceDesc:  "Ошибка получения данных в базе",
				}
			}
		}

		row := excel.Row{
			ClientName:               client.ClientName,
			MinutesCount:             client.MinutesCount,
			MinutesAmountWoTax:       client.MinutesAmountWoTax,
			NumbersCount:             servDesc.NumbersCount,
			ServiceDescription:       servDesc.ServiceDesc,
			ServicesAmountWithoutTax: client.ServicesAmountWoTaz,
			ServicesAmountWithTax:    client.ServicesAmountWoTaz.Mul(tax).Round(2),
			TotalAmountWithoutTax:    client.TotalAmountWoTax,
			TotalAmountWithTax:       client.TotalAmountWoTax.Mul(tax).Round(2),
			CompanyAffiliation:       client.CompanyAffiliation,
			DocNumber:                client.DocNumber,
			VPBXAmountWithTax:        servDesc.VPBXAmount,
			AmountFromBLToBI:         (client.TotalAmountWoTax.Mul(tax)).Sub(servDesc.VPBXAmount).Round(2),
			CallsCount:               client.CallsCount,
		}

		var companyInfo *excel.CompanyData
		switch client.CompanyAffiliation {
		case "БЛ":
			companyInfo = &output.BL
		case "БИ":
			companyInfo = &output.BI
		default:
			companyInfo = &output.BL
		}

		// Суммирование значений для поля Итого
		companyInfo.SumAdditionalServices = companyInfo.SumAdditionalServices.Add(row.ServicesAmountWithTax)
		companyInfo.SumMinutesCount = companyInfo.SumMinutesCount.Add(row.MinutesCount)
		companyInfo.SumTotalAmountWithTax = companyInfo.SumTotalAmountWithTax.Add(row.TotalAmountWithTax)
		companyInfo.SumNumbersCount += row.NumbersCount
		companyInfo.SumCallsCount += row.CallsCount

		// Добавление готовой строки
		companyInfo.Data = append(companyInfo.Data, row)
	}

	return &output
}
