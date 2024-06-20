package v1

import (
	"encoding/json"
	"net/http"

	"github.com/xboshy/go-template/config"
	"github.com/xboshy/go-template/node"

	"github.com/labstack/echo/v4"
)

type CommonInterface interface {
	Status() (s node.StatusReport, err error)
	Config() config.Template
}

func (h *Handlers) HealthCheck(ctx echo.Context) error {
	w := ctx.Response().Writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})

	return nil
}
