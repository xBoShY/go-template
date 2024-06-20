package node

import (
	"context"
	"fmt"

	"github.com/xboshy/go-deadlock"
	"github.com/xboshy/go-template/config"
	"github.com/xboshy/go-template/data"
	"github.com/xboshy/go-template/data/dummy"
	"github.com/xboshy/go-template/logging"
)

type Node struct {
	mu         deadlock.Mutex
	ctx        context.Context
	config     config.Template
	cancelCtx  context.CancelFunc
	txnHandler *data.ReqHandler[dummy.Transaction, dummy.TransactionResult]
	log        logging.Logger
}

type StatusReport struct {
}

func MakeNode(log logging.Logger, rootDir string, cfg config.Template) (*Node, error) {
	var err error

	node := new(Node)
	node.log = log.With("name", cfg.NetAddress)
	node.config = cfg

	node.txnHandler, err = data.MakeReqHandler[dummy.Transaction, dummy.TransactionResult](log, 10)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialize the transaction handler: %s", err)
	}

	return node, nil
}

func (node *Node) Config() config.Template {
	return node.config
}

func (node *Node) Status() (StatusReport, error) {
	var s StatusReport
	var err error

	return s, err
}

func (node *Node) Start() {
	node.mu.Lock()
	defer node.mu.Unlock()

	node.ctx, node.cancelCtx = context.WithCancel(context.Background())
	node.txnHandler.Start()
}

func (node *Node) Stop() {
	node.cancelCtx()
}

func (node *Node) Process(msg data.BacklogMsg[dummy.Transaction, dummy.TransactionResult]) {
	node.txnHandler.Process(msg)
}
