package routes

import (
	"net/http"

	"github.com/your-org/your-app/handler"
)

func NewRouter(healthHandler *handler.HealthHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler.GetHealth)
	return mux
}
