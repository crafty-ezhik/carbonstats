package service_description

import (
	"github.com/crafty-ezhik/carbonstats/pkg/req"
	"github.com/crafty-ezhik/carbonstats/pkg/res"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ServiceDescriptionHandler interface {
	GetByCarbonPK() http.HandlerFunc
	GetAll() http.HandlerFunc
	Create() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
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
		carbonPkStr := chi.URLParam(r, "carbon_pk")
		if carbonPkStr == "" {
			res.JSON(w, "Invalid carbon pk", http.StatusBadRequest)
			return
		}
		carbonPk, err := strconv.Atoi(carbonPkStr)
		result, err := h.SerDescRepo.Get(uint(carbonPk))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.JSON(w, result, http.StatusOK)
	}
}

func (h *serviceDescriptionHandlerImpl) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := h.SerDescRepo.List()
		if err != nil {
			res.JSON(w, err, http.StatusBadRequest)
		}
		res.JSON(w, result, http.StatusOK)
	}
}

func (h *serviceDescriptionHandlerImpl) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("batch") == "true" {
			body, err := req.HandleBody[CreateBatchRequest](w, r)
			if err != nil {
				return
			}

			var data []*ServiceDescription
			for _, item := range body.Data {
				amount, err := decimal.NewFromString(item.VPBXAmount)
				if err != nil {
					res.JSON(w, "Некорректный vpbx amount", http.StatusBadRequest)
				}

				temp := &ServiceDescription{
					CarbonPK:     item.CarbonPK,
					NumbersCount: item.NumbersCount,
					VPBXAmount:   amount,
					ServiceDesc:  item.ServiceDesc,
				}

				data = append(data, temp)
			}
			err = h.SerDescRepo.CreateBatch(data)
			if err != nil {
				res.JSON(w, err, http.StatusInternalServerError)
				return
			}

			res.JSON(w, "Данные успешно добавлены", http.StatusOK)

		} else {
			body, err := req.HandleBody[CreateRequest](w, r)
			if err != nil {
				return
			}
			amount, err := decimal.NewFromString(body.VPBXAmount)
			if err != nil {
				res.JSON(w, "Некорректный vpbx amount", http.StatusBadRequest)
			}

			data := &ServiceDescription{
				CarbonPK:     body.CarbonPK,
				NumbersCount: body.NumbersCount,
				VPBXAmount:   amount,
				ServiceDesc:  body.ServiceDesc,
			}
			if err = h.SerDescRepo.Create(data); err != nil {
				res.JSON(w, err.Error(), http.StatusBadRequest)
			}
			res.JSON(w, "Запись добавлена", http.StatusCreated)
		}
	}
}

func (h *serviceDescriptionHandlerImpl) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		carbonPkStr := chi.URLParam(r, "carbon_pk")
		if carbonPkStr == "" {
			res.JSON(w, "Invalid carbon pk", http.StatusBadRequest)
			return
		}
		carbonPk, err := strconv.Atoi(carbonPkStr)
		if err != nil {
			res.JSON(w, err, http.StatusBadRequest)
		}

		body, err := req.HandleBody[UpdateRequest](w, r)
		if err != nil {
			return
		}
		amount, err := decimal.NewFromString(body.VPBXAmount)
		if err != nil {
			res.JSON(w, "Некорректный vpbx amount", http.StatusBadRequest)
		}

		data := &ServiceDescription{
			NumbersCount: body.NumbersCount,
			VPBXAmount:   amount,
			ServiceDesc:  body.ServiceDesc,
		}

		err = h.SerDescRepo.Update(uint(carbonPk), data)
		if err != nil {
			res.JSON(w, err, http.StatusBadRequest)
			return
		}
		res.JSON(w, body, http.StatusOK)
	}
}

func (h *serviceDescriptionHandlerImpl) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		carbonPkStr := chi.URLParam(r, "carbon_pk")
		if carbonPkStr == "" {
			res.JSON(w, "Invalid carbon pk", http.StatusBadRequest)
			return
		}
		carbonPk, err := strconv.Atoi(carbonPkStr)
		if err != nil {
			res.JSON(w, err, http.StatusBadRequest)
			return
		}
		err = h.SerDescRepo.Delete(uint(carbonPk))
		if err != nil {
			res.JSON(w, err, http.StatusBadRequest)
			return
		}
		res.JSON(w, nil, http.StatusNoContent)
	}
}
