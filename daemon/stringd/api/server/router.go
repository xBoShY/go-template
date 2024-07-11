package server

import (
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/xboshy/go-template/daemon/stringd/api/generated/routes"
	"github.com/xboshy/go-template/daemon/stringd/api/server/handlers"
	"github.com/xboshy/go-template/daemon/stringd/api/server/lib/ctxkeys"
	"github.com/xboshy/go-template/logging"
	"github.com/xboshy/go-template/node"
	"github.com/xboshy/go-template/util/uuid"
)

type APINode struct {
	*node.StringsNode
}

type APINodeInterface interface {
	handlers.NodeInterface
}

func NewRouter(logger logging.Logger, node APINodeInterface, shutdown <-chan struct{}, listener net.Listener) (*fiber.App, error) {
	f := fiber.New()
	f.Server().Logger = logger

	v1Handler := handlers.Handlers{
		Node:     node,
		Log:      logger,
		Shutdown: shutdown,
	}

	options := routes.FiberServerOptions{
		BaseURL: "",
		Middlewares: []routes.MiddlewareFunc{
			// ADD REQUEST LOGGING
			requestid.New(
				requestid.Config{
					Header:     fiber.HeaderXRequestID,
					Generator:  uuid.NewUUID().String,
					ContextKey: ctxkeys.RequestID,
				},
			),
		},
	}

	routes.RegisterHandlersWithOptions(f, &v1Handler, options)

	return f, nil
}
