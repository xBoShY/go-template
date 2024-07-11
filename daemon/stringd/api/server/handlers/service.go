package handlers

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xboshy/go-template/daemon/stringd/api/generated/model"
	"github.com/xboshy/go-template/daemon/stringd/api/server/lib/ctxkeys"
	"github.com/xboshy/go-template/data"
	strings "github.com/xboshy/go-template/data/strings"
	"github.com/xboshy/go-template/protocol"
)

type ServiceInterface interface {
	Process(data.BacklogMsg[strings.Transaction, strings.TransactionResult])
}

var functionsMap = map[model.Function]strings.StringFn{
	model.FunctionReverse: strings.ReverseStringFn,
	model.FunctionUpper:   strings.UpperStringFn,
	model.FunctionLower:   strings.LowerStringFn,
}

func modelFunctionsToStringFn(modelFunctions []model.Function) ([]strings.StringFn, error) {
	res := make([]strings.StringFn, 0)

	for _, mf := range modelFunctions {
		sf, exists := functionsMap[mf]
		if !exists {
			return nil, fmt.Errorf("function %s does not exist", mf)
		}
		res = append(res, sf)
	}

	return res, nil
}

// (POST /)
func (v1 *Handlers) Request(ctx *fiber.Ctx) error {
	var err error
	log := v1.Log

	req := ctx.Request()
	resp := ctx.Response()
	defer func() {
		if resp.Header.StatusCode() != http.StatusOK {
			resp.ResetBody()
		}
	}()

	bodyReader := req.BodyStream()
	dec := protocol.NewDecoder(bodyReader)
	replyTo := make(chan strings.TransactionResult)
	defer close(replyTo)

	var mReq model.Request
	err = dec.Decode(&mReq)
	if err != nil {
		log.Errorf("could not decode body: %v", err)
		resp.Header.SetStatusCode(http.StatusBadRequest)
		return nil
	}

	stringFns, err := modelFunctionsToStringFn(mReq.Functions)
	if err != nil {
		log.Error(err)
		resp.Header.SetStatusCode(http.StatusBadRequest)
		return nil
	}

	txn := strings.MakeTransaction(stringFns, mReq.Message)
	if mReq.Uuid != nil {
		txn.WithUuid(mReq.Uuid.String())
	} else {
		reqId := ctx.Get(ctxkeys.RequestID)
		txn.WithUuid(reqId)
	}

	if mReq.Timestamp != nil {
		ts := time.UnixMilli(*mReq.Timestamp)
		txn.WithTimestamp(ts)
	}
	if mReq.Ttl != nil {
		txn.WithTTL(*mReq.Ttl)
	}

	work := strings.MakeWork(txn).WithReplyTo(replyTo)
	v1.Node.Process(work)

	txnResp := <-replyTo
	err = txnResp.GetError()
	if err != nil {
		log.Errorf("could process request: %v", err)
		resp.Header.SetStatusCode(http.StatusBadRequest)
		return nil
	}

	var mResp model.Response = model.Response{
		ReversedMessage: txnResp.GetMsg(),
		Uuid:            txnResp.GetUUID(),
	}

	respBodyReader, respBodyWriter := io.Pipe()
	bufWriter := bufio.NewWriter(respBodyWriter)
	bufReader := bufio.NewReader(respBodyReader)

	enc := protocol.NewEncoder(bufWriter)
	err = enc.Encode(mResp)
	if err != nil {
		log.Errorf("could not encode response: %v", err)
		resp.Header.SetStatusCode(http.StatusInternalServerError)
		return nil
	}
	err = bufWriter.Flush()
	if err != nil {
		log.Errorf("could not flush response body: %v", err)
		resp.Header.SetStatusCode(http.StatusInternalServerError)
		return nil
	}

	resp.SetStatusCode(http.StatusOK)
	resp.Header.SetContentType("application/x-binary")
	resp.SetBodyStream(bufReader, bufReader.Buffered())

	return nil
}
