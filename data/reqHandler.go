package data

import (
	"context"
	"sync"

	"github.com/xboshy/go-template/logging"
)

type ReqHandler[REQ any, RESP any] struct {
	ctx          context.Context
	ctxCancel    context.CancelFunc
	backlogWg    sync.WaitGroup
	backlogQueue chan BacklogMsg[REQ, RESP]
	log          logging.Logger
	resource     Resource[REQ, RESP]
}

type Resource[REQ any, RESP any] interface {
	Process(REQ) RESP
}

type BacklogMsg[REQ any, RESP any] interface {
	Request() REQ
	Reply(RESP)
}

func MakeReqHandler[REQ any, RESP any](log logging.Logger, maxBacklogSize int) (*ReqHandler[REQ, RESP], error) {
	backlogSize := maxBacklogSize

	handler := &ReqHandler[REQ, RESP]{
		backlogQueue: make(chan BacklogMsg[REQ, RESP], backlogSize),
		log:          log,
	}

	return handler, nil
}

func (handler *ReqHandler[REQ, RESP]) Start() {
	handler.ctx, handler.ctxCancel = context.WithCancel(context.Background())
	handler.backlogWg.Add(1)
	go handler.handler()
}

func (handler *ReqHandler[REQ, RESP]) Stop() {
	handler.ctxCancel()
	handler.backlogWg.Wait()
}

func (handler *ReqHandler[REQ, RESP]) Process(msg BacklogMsg[REQ, RESP]) {
	handler.backlogQueue <- msg
}

func (handler *ReqHandler[REQ, RESP]) handler() {
	defer handler.backlogWg.Done()
	for {
		select {
		case msg := <-handler.backlogQueue:
			req := msg.Request()
			resp := handler.resource.Process(req)
			msg.Reply(resp)
		case <-handler.ctx.Done():
			return
		}
	}
}
