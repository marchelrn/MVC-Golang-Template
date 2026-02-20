package server

import (
	"fmt"
	"net/http"

	"github.com/your-org/your-app/config"
	"github.com/your-org/your-app/handler"
	"github.com/your-org/your-app/repository"
	"github.com/your-org/your-app/routes"
	"github.com/your-org/your-app/service"
)

func Run(cfg config.Config) error {
	repo := repository.NewHealthRepository()
	svc := service.NewHealthService(repo)
	h := handler.NewHealthHandler(svc)

	r := routes.NewRouter(h)
	addr := fmt.Sprintf(":%s", cfg.Port)

	return http.ListenAndServe(addr, r)
}
