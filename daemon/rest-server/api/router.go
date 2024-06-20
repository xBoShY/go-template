package api

import (
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v1 "github.com/xboshy/go-template/daemon/rest-server/api/v1"
	"github.com/xboshy/go-template/daemon/rest-server/api/v1/generated/common"
	"github.com/xboshy/go-template/daemon/rest-server/api/v1/generated/service"
	"github.com/xboshy/go-template/daemon/rest-server/lib/middlewares"
	"github.com/xboshy/go-template/logging"
	"github.com/xboshy/go-template/node"
)

const (
	BaseURL = "v1"
)

type APINode struct {
	*node.Node
}

type APINodeInterface interface {
	v1.NodeInterface
}

func NewRouter(logger logging.Logger, node APINodeInterface, shutdown <-chan struct{}, listener net.Listener) *echo.Echo {
	e := echo.New()
	e.Logger = logger.MakeEchoLogger()

	e.Listener = listener
	e.HideBanner = true

	e.Pre(
		middleware.RemoveTrailingSlash(),
	)
	e.Use(
		middlewares.MakeLogger(logger),
	)

	v1Handler := v1.Handlers{
		Node:     node,
		Log:      logger,
		Shutdown: shutdown,
	}

	common.RegisterHandlersWithBaseURL(e, &v1Handler, BaseURL)
	service.RegisterHandlersWithBaseURL(e, &v1Handler, BaseURL)

	return e
}
