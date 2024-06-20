package v1

import "github.com/xboshy/go-template/logging"

// NodeInterface represents node fns used by the handlers.
type NodeInterface interface {
	CommonInterface
	ServiceInterface
}

type Handlers struct {
	Node     NodeInterface
	Log      logging.Logger
	Shutdown <-chan struct{}
}
