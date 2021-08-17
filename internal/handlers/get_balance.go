package handlers

import (
	"fmt"

	"github.com/ehabterra/flash_api/api/models"
	"github.com/ehabterra/flash_api/api/restapi/operations/users"
	"github.com/go-openapi/runtime/middleware"
)

type GetBalanceUserService interface {
	GetUserBalance(id string) (float64, error)
}

type getBalanceHandler struct {
	service GetBalanceUserService
}

func NewUsersGetBalanceHandler(service GetBalanceUserService) users.GetBalanceHandler {
	return &getBalanceHandler{service}
}

func (h *getBalanceHandler) Handle(params users.GetBalanceParams, principle *models.Principle) middleware.Responder {
	balance, err := h.service.GetUserBalance(*principle.ID)
	if err != nil {
		fmt.Println(err.Error())
		msg := err.Error()
		return users.NewGetBalanceDefault(500).WithPayload(&models.Error{Message: &msg})
	}
	return users.NewGetBalanceOK().WithPayload(balance)
}
