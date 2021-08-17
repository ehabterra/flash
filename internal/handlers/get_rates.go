package handlers

import (
	"fmt"

	"github.com/ehabterra/flash_api/api/models"
	"github.com/ehabterra/flash_api/api/restapi/operations/rates"
	"github.com/go-openapi/runtime/middleware"
)

type GetRatesUserService interface {
	GetRates(base string, target string) (float64, error)
}

type getRatesHandler struct {
	service GetRatesUserService
}

func NewRatesGetRatesHandler(service GetRatesUserService) rates.GetRatesHandler {
	return &getRatesHandler{service}
}

func (h *getRatesHandler) Handle(params rates.GetRatesParams, principle *models.Principle) middleware.Responder {
	base := params.Base
	target := params.Target
	rate, err := h.service.GetRates(base, target)
	if err != nil {
		fmt.Println(err.Error())
		msg := err.Error()
		return rates.NewGetRatesDefault(500).WithPayload(&models.Error{Message: &msg})
	}
	return rates.NewGetRatesOK().WithPayload(rate)
}
