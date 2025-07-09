package statistics

import (
	"github.com/crafty-ezhik/carbonstats/pkg/req"
	"github.com/crafty-ezhik/carbonstats/pkg/res"
	"go.uber.org/zap"
	"net/http"
)

type StatisticsHandler interface {
	GetAll() http.HandlerFunc
	Create() http.HandlerFunc
}

type statisticsHandlerImpl struct {
	StatisticsRepo StatisticsRepository

	log            *zap.Logger
}

func NewStatisticsHandler(statsRepo StatisticsRepository, logger *zap.Logger) StatisticsHandler {
	return &statisticsHandlerImpl{
		StatisticsRepo: statsRepo,
		log:            logger,
	}
}

func (h *statisticsHandlerImpl) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := h.StatisticsRepo.GetAll()
		if err != nil {
			h.log.Error("Error getting statistics", zap.Error(err))
			res.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.JSON(w, result, http.StatusOK)
	}
}

func (h *statisticsHandlerImpl) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[CreateRequest](w, r)
		if err != nil {
			return
		}
		res.JSON(w, body, http.StatusOK)
	}
}


