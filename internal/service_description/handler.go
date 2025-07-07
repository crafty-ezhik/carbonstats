package service_description

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type ServiceDescriptionHandler interface {
	GetByCarbonPK() http.HandlerFunc
}

type serviceDescriptionHandlerImpl struct {
	log         *zap.Logger
	SerDescRepo ServiceDescriptionRepository
}

func NewServiceDescriptionHandler(logger *zap.Logger, servDescRepo ServiceDescriptionRepository) ServiceDescriptionHandler {
	return &serviceDescriptionHandlerImpl{
		log:         logger,
		SerDescRepo: servDescRepo,
	}
}

func (h *serviceDescriptionHandlerImpl) GetByCarbonPK() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		result, err := h.SerDescRepo.Get(1345)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonData, err := json.Marshal(result)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
