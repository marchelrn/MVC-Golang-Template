package app

import (
	"github.com/marchelrn/stock_api/config"
	"github.com/marchelrn/stock_api/internal/server"
)

func Run() {
	config.Load()
	server.Run()
}	
