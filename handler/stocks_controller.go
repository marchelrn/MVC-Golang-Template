package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marchelrn/stock_api/contract"
	errs "github.com/marchelrn/stock_api/pkg/error"
)

type StocksController struct {
	service contract.StocksService
}

func (s *StocksController) InitService(svc *contract.Service) {
	fmt.Println("DEBUG: Initializing StocksController with StocksService")
	if svc == nil {
		fmt.Println("ERROR: Provided service is nil")
		return
	}
	s.service = svc.Stocks
}

func (s *StocksController) GetStocks(c *gin.Context) {
	stockTickerParam := c.Param("ticker")
	if stockTickerParam == "" {
		HandleError(c, errs.BadRequest("Stock ticker is required"))
		return
	}

	tickers := strings.Split(strings.ToUpper(stockTickerParam), ",")

	reponse, err := s.service.GetStocks(tickers)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": reponse.Message,
		"data": reponse.StocksData,
	})
}
