package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marchelrn/stock_api/contract"
)

type BrokersController struct {
	service *contract.Service
}

func ImplBrokersController(service *contract.Service) *BrokersController {
	return &BrokersController{service: service}
}

func (c *BrokersController) InitService(svc *contract.Service) {
	c.service = svc
}

func (c *BrokersController) GetBrokersDetails(ctx *gin.Context) {
	brokerName := ctx.Param("name")
	brokerNames := strings.Split(brokerName, ",")

	brokerDetails, err := c.service.Brokers.GetBrokersDetails(brokerNames)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": brokerDetails.Message,
		"data":    brokerDetails.BrokersData,
	})
}
