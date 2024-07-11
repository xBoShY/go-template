package node

import (
	"context"
	"fmt"

	"github.com/xboshy/go-deadlock"
	"github.com/xboshy/go-template/config"
	"github.com/xboshy/go-template/data"
	strings "github.com/xboshy/go-template/data/strings"
	"github.com/xboshy/go-template/logging"
)

type StringsNode struct {
	mu         deadlock.Mutex
	ctx        context.Context
	config     config.Template
	cancelCtx  context.CancelFunc
	txnHandler *data.ReqHandler[strings.Transaction, strings.TransactionResult]
	log        logging.Logger
}

type StringsNodeStatus struct {
}

func MakeStringsNode(log logging.Logger, rootDir string, cfg config.Template) (*StringsNode, error) {
	var err error

	node := new(StringsNode)
	node.log = log.With("name", cfg.EndpointAddress)
	node.config = cfg

	processor := strings.MakeProcessor()

	node.txnHandler, err = data.MakeReqHandler(log, 10, processor)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialize the transaction handler: %s", err)
	}

	return node, nil
}

func (node *StringsNode) Config() config.Template {
	return node.config
}

func (node *StringsNode) Status() (StringsNodeStatus, error) {
	var s StringsNodeStatus
	var err error

	return s, err
}

func (node *StringsNode) Start() {
	node.mu.Lock()
	defer node.mu.Unlock()

	node.ctx, node.cancelCtx = context.WithCancel(context.Background())
	node.txnHandler.Start()
}

func (node *StringsNode) Stop() {
	node.cancelCtx()
}

func (node *StringsNode) Process(msg data.BacklogMsg[strings.Transaction, strings.TransactionResult]) {
	node.txnHandler.Process(msg)
}
