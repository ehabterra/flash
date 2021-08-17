package handlers

import (
	"fmt"

	"github.com/ehabterra/flash_api/api/models"
	"github.com/ehabterra/flash_api/api/restapi/operations/users"
	"github.com/go-openapi/runtime/middleware"
)

type UploadUserService interface {
	Upload(id string, number string, amount int64) error
}

type uploadHandler struct {
	service UploadUserService
}

func NewUsersUploadHandler(service UploadUserService) users.UploadHandler {
	return &uploadHandler{service}
}

func (h *uploadHandler) Handle(params users.UploadParams, principle *models.Principle) middleware.Responder {
	amount := params.Amount
	accountNumber := params.AccountNumber
	err := h.service.Upload(*principle.ID, accountNumber, amount)
	if err != nil {
		fmt.Println(err.Error())
		msg := err.Error()
		return users.NewUploadDefault(500).WithPayload(&models.Error{Message: &msg})
	}

	return users.NewUploadOK()
}
