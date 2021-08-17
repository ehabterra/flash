package handlers

import (
	"fmt"
	"time"

	models "github.com/ehabterra/flash_api/internal/models"

	apiModels "github.com/ehabterra/flash_api/api/models"
	"github.com/ehabterra/flash_api/api/restapi/operations/users"
	"github.com/ehabterra/flash_api/internal/middlewares"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserService interface {
	GetUserByUsernameOrEmail(usernameOrEmail string) (*models.User, error)
}

type loginHandler struct {
	service LoginUserService
}

func NewUsersLoginHandler(service LoginUserService) users.LoginHandler {
	return &loginHandler{service}
}

func (h *loginHandler) Handle(params users.LoginParams) middleware.Responder {
	usernameOrEmail := params.Body.UsernameOrEmail
	userInfo, err := h.service.GetUserByUsernameOrEmail(*usernameOrEmail)
	responder := errorHandling(err)
	if responder != nil {
		return responder
	}
	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(*params.Body.Password))
	responder = errorHandling(err)
	if responder != nil {
		return responder
	}
	token, err := middlewares.GenerateJWT(userInfo.ID, userInfo.Email, userInfo.Username)
	responder = errorHandling(err)
	if responder != nil {
		return responder
	}

	now := strfmt.DateTime(time.Now())
	return users.NewLoginOK().WithPayload(&apiModels.LoginResponse{ExpireDate: &now, Token: &token})
}

func errorHandling(err error) middleware.Responder {
	if err != nil {
		fmt.Println(err.Error())
		msg := err.Error()
		return users.NewLoginDefault(500).WithPayload(&apiModels.Error{Message: &msg})
	}
	return nil
}
