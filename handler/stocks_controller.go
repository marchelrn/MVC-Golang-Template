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

	rawTickers := strings.Split(strings.ToUpper(stockTickerParam), ",")
	tickers := make([]string, 0, len(rawTickers))
	seen := make(map[string]struct{}, len(rawTickers))
	for _, ticker := range rawTickers {
		ticker = strings.TrimSpace(ticker)
		if ticker == "" {
			continue
		}
		if _, exists := seen[ticker]; exists {
			continue
		}
		seen[ticker] = struct{}{}
		tickers = append(tickers, ticker)
	}

	if len(tickers) == 0 {
		HandleError(c, errs.BadRequest("Stock ticker is required"))
		return
	}

	response, err := s.service.GetStocks(tickers)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": response.Message,
		"data":    response.StocksData,
	})
}
