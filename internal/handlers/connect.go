package handlers

import (
	"fmt"

	"github.com/ehabterra/flash_api/internal/models"

	apiModels "github.com/ehabterra/flash_api/api/models"
	"github.com/ehabterra/flash_api/api/restapi/operations/users"
	"github.com/go-openapi/runtime/middleware"
)

type ConnectUserService interface {
	Connect(account *models.Account) error
}

type connectHandler struct {
	service ConnectUserService
}

func NewUsersConnectHandler(service ConnectUserService) users.ConnectHandler {
	return &connectHandler{service}
}

func (h *connectHandler) Handle(params users.ConnectParams, principle *apiModels.Principle) middleware.Responder {
	err := h.service.Connect(&models.Account{
		AccountNumber: params.AccountNumber,
		UserID:        *principle.ID,
		BankID:        *params.Body.BankID,
		BranchNumber:  *params.Body.BranchNumber,
		HolderName:    *params.Body.HolderName,
		Reference:     params.Body.Reference,
	})
	if err != nil {
		return h.showError(err)
	}

	return users.NewConnectOK()
}

func (h *connectHandler) showError(err error) middleware.Responder {
	fmt.Println(err.Error())
	msg := err.Error()
	return users.NewConnectDefault(500).WithPayload(&apiModels.Error{Message: &msg})
}
