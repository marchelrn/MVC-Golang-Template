package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marchelrn/stock_api/contract"
)

type HealthController struct {
	service contract.HealthService
}

func (c *HealthController) InitService(s *contract.Service) {
	fmt.Println("DEBUG: Initializing UserController with UserService")
	if s == nil {
		fmt.Println("ERROR: Provided service is nil")
		return
	}
}

func (h *HealthController) GetHealth(c *gin.Context) {

	if h.service == nil {
		fmt.Println("ERROR: HealthService is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service not available"})
		return
	}

	response := h.service.GetStatus()
	c.JSON(http.StatusOK, gin.H{
		"status": response.Status,
	})
}
