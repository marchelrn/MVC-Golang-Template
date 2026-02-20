package handler

import (
	"log"
	"mini_jira/contract"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	InitService(service *contract.Service)
}

func New(app *gin.Engine, service *contract.Service) {
	allControllers := []Controller{
		&UserController{},
		&EmailVerificationController{},
	}

	for _, ctrl := range allControllers {
		ctrl.InitService(service)
		log.Println("Initialized controller:", ctrl)
	}
}
