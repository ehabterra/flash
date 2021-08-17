package handlers

import (
	"fmt"

	"github.com/ehabterra/flash_api/api/models"
	"github.com/ehabterra/flash_api/api/restapi/operations/users"
	"github.com/go-openapi/runtime/middleware"
)

type SendUserService interface {
	Send(id string, usernameOrEmail string, amount int64) error
}

type sendHandler struct {
	service SendUserService
}

func NewUsersSendHandler(service SendUserService) users.SendHandler {
	return &sendHandler{service}
}

func (h *sendHandler) Handle(params users.SendParams, principle *models.Principle) middleware.Responder {
	amount := params.Amount
	usernameOrEmail := params.UsernameOrEmail
	err := h.service.Send(*principle.ID, usernameOrEmail, amount)
	if err != nil {
		fmt.Println(err.Error())
		msg := err.Error()
		return users.NewSendDefault(500).WithPayload(&models.Error{Message: &msg})
	}

	return users.NewSendOK()
}
