package periodic_tasks

import (
	"github.com/crafty-ezhik/carbonstats/config"
	"github.com/crafty-ezhik/carbonstats/internal/carbon"
	"github.com/crafty-ezhik/carbonstats/internal/excel"
	"github.com/crafty-ezhik/carbonstats/internal/service_description"
	"github.com/crafty-ezhik/carbonstats/internal/statistics"
	"github.com/crafty-ezhik/carbonstats/internal/utils"
	"go.uber.org/zap"
	"time"
)

type CustomerReport struct {
	CarbonCfg    *config.CarbonConfig
	ServDescRepo service_description.ServiceDescriptionRepository
	StatsRepo    statistics.StatisticsRepository
	Log          *zap.Logger
	Filename     string
	Date         time.Time
}

func (cr *CustomerReport) RunTask() error {
	cr.Log.Info("Starting GetMonthlyClientStatistics task")
	billing := carbon.NewCarbonBilling(cr.CarbonCfg, cr.Log)

	cr.Log.Info("Initializing Excel module")
	ex := excel.New(cr.Log, cr.Filename)
	cr.Log.Info("Excel module has been successfully initialized")

	cr.Log.Info("Starting receiving data from billing...")
	data, err := billing.StartStatisticsCollection(cr.Date)
	if err != nil {
		cr.Log.Fatal("Error getting billing data", zap.Error(err))
		return err
	}
	cr.Log.Info("Billing data has been successfully received")

	cr.Log.Info("Starting the data conversion process")
	formatData := utils.DataPreparation(cr.ServDescRepo, data, cr.Log)
	cr.Log.Info("The data has been successfully converted to the required format")

	cr.Log.Info("Adding the received statistics to the database...")
	err = cr.StatsRepo.CreateBatch(data)
	if err != nil {
		cr.Log.Fatal("Error adding statistics to the database", zap.Error(err))
		return err
	}
	cr.Log.Info("Statistics have been successfully added to the database!")

	cr.Log.Info("Starting the formation of the excel file...")
	err = ex.AddData(formatData)
	if err != nil {
		cr.Log.Fatal("Error processing excel file", zap.Error(err))
		return err
	}
	cr.Log.Info("Finished processing excel file")

	cr.Log.Info("Save excel file")
	err = ex.Save()
	if err != nil {
		cr.Log.Fatal("Error saving excel file", zap.Error(err))
		return err
	}
	cr.Log.Info("Finished saving excel file")
	cr.Log.Info("GetMonthlyClientStatistics task has been successfully completed")
	return nil
}
