package statistics

import "go.uber.org/zap"

type StatisticsService interface {
}

type statisticsServiceImpl struct {
	StatsRepo StatisticsRepository
	log       *zap.Logger
}

func NewStatisticsService(statsRepo StatisticsRepository, logger *zap.Logger) StatisticsService {
	return &statisticsServiceImpl{
		StatsRepo: statsRepo,
		log:       logger,
	}
}
