package statistics

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type StatisticsRepository interface {
	GetAll() ([]ClientStatistics, error)
	CreateBatch(stats []ClientStatistics) error
	GetByDate(month, year int) ([]ClientStatistics, error)
}

type statisticsImpl struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewStatisticsRepository(db *gorm.DB, log *zap.Logger) StatisticsRepository {
	return &statisticsImpl{db: db, log: log}
}

func (repo *statisticsImpl) GetByDate(month, year int) ([]ClientStatistics, error) {
	var clientStatistics []ClientStatistics
	err := repo.db.Where("month = ? AND year = ?", month, year).Find(&clientStatistics).Error
	return clientStatistics, err
}

func (repo *statisticsImpl) GetAll() ([]ClientStatistics, error) {
	var result []ClientStatistics
	err := repo.db.Find(&result).Error
	return result, err
}

func (repo *statisticsImpl) CreateBatch(stats []ClientStatistics) error {
	err := repo.db.Create(&stats).Error
	return err
}
