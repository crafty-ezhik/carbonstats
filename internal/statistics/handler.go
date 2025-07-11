package statistics

import (
	"fmt"
	"github.com/crafty-ezhik/carbonstats/pkg/req"
	"github.com/crafty-ezhik/carbonstats/pkg/res"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type StatisticsHandler interface {
	GetAll() http.HandlerFunc
	GetByDate() http.HandlerFunc
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

func (h *statisticsHandlerImpl) GetByDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		month, err := strconv.Atoi(r.URL.Query().Get("month"))
		if err != nil {
			h.log.Warn("month", zap.Error(err))
			res.JSON(w, "Month must be positive integer", http.StatusBadRequest)
			return
		}

		year, err := strconv.Atoi(r.URL.Query().Get("year"))
		if err != nil {
			h.log.Warn("month", zap.Error(err))
			res.JSON(w, "Year must be positive integer", http.StatusBadRequest)
			return
		}

		result, err := h.StatisticsRepo.GetByDate(month, year)
		if len(result) == 0 {
			h.log.Info(fmt.Sprintf("Not found records for period month: %v, year: %v", month, year))
			res.JSON(w, fmt.Sprintf("Not found records for period month: %v, year: %v", month, year), http.StatusNotFound)
			return
		}
		if err != nil {
			h.log.Warn(fmt.Sprintf("error getting statistics for month %d, year %d", month, year), zap.Error(err))
			res.JSON(w, "Error getting statistics for month %d, year", http.StatusInternalServerError)
			return
		}
		res.JSON(w, result, http.StatusOK)
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
