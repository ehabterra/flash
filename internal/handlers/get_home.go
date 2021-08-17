package handlers

import (
	"github.com/ehabterra/flash_api/api/restapi/operations/home"
	"github.com/go-openapi/runtime/middleware"
)

type homeHandler struct{}

func NewHomeHandler() home.GetHandler {
	return &homeHandler{}
}

func (c *homeHandler) Handle(params home.GetParams) middleware.Responder {
	return home.NewGetOK().WithPayload("Welcome!")
}
