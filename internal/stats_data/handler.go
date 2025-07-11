package stats_data

import (
	"fmt"
	"github.com/crafty-ezhik/carbonstats/internal/carbon"
	"github.com/crafty-ezhik/carbonstats/internal/service_description"
	"github.com/crafty-ezhik/carbonstats/internal/statistics"
	"github.com/crafty-ezhik/carbonstats/pkg/res"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type StatsDataHandler interface {
	GetStats() http.HandlerFunc
}

type statsDataHandler struct {
	StatisticsRepo statistics.StatisticsRepository
	ServiceDesc    service_description.ServiceDescriptionRepository
	Billing        carbon.CarbonBilling
	log            *zap.Logger
}

func NewStatsDataHandler(StatisticsRepo statistics.StatisticsRepository, ServiceDesc service_description.ServiceDescriptionRepository, billing carbon.CarbonBilling, log *zap.Logger) StatsDataHandler {
	return &statsDataHandler{
		StatisticsRepo: StatisticsRepo,
		ServiceDesc:    ServiceDesc,
		Billing:        billing,
		log:            log,
	}
}

func (h *statsDataHandler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		month, err := strconv.Atoi(r.URL.Query().Get("month"))
		if err != nil {
			res.JSON(w, err.Error(), http.StatusBadRequest)
			return
		}
		year, err := strconv.Atoi(r.URL.Query().Get("year"))
		if err != nil {
			res.JSON(w, err.Error(), http.StatusBadRequest)
			return
		}
		// TODO: Сделать выбор по месяцам
		h.log.Info(fmt.Sprintf("Запрос на получение статистики за %d-%d", month, year))

		res.JSON(w, "Coming soon...", http.StatusTeapot)
	}
}
