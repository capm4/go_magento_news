package controller

import "magento/bot/pkg/config"

type AppController struct {
	Website interface {
		WebsiteControllerInterface
	}
	User interface {
		UserControllerInterface
	}
	Config *config.Сonfig
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

func CreateErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Message: msg}
}

func CreateResponseMsg(msg string) *ResponseMessage {
	return &ResponseMessage{Message: msg}
}
