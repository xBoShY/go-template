package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/xboshy/go-template/config"
	"github.com/xboshy/go-template/node"
)

type CommonInterface interface {
	Status() (s node.StringsNodeStatus, err error)
	Config() config.Template
}

func (h *Handlers) HealthCheck(ctx *fiber.Ctx) error {
	resp := ctx.Response()
	resp.Header.SetContentType("application/json")
	resp.Header.SetStatusCode(http.StatusOK)
	bw := resp.BodyWriter()
	json.NewEncoder(bw).Encode(struct{}{})

	return nil
}
