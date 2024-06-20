package v1

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xboshy/go-template/daemon/rest-server/api/v1/generated/model"
	"github.com/xboshy/go-template/daemon/rest-server/lib/context"
	"github.com/xboshy/go-template/data"
	"github.com/xboshy/go-template/data/dummy"
	"github.com/xboshy/go-template/protocol"
)

type ServiceInterface interface {
	Process(data.BacklogMsg[dummy.Transaction, dummy.TransactionResult])
}

// (POST /)
func (v1 *Handlers) Request(ctx echo.Context) error {
	var err error
	log := v1.Log

	req := ctx.Request()
	if req == nil {
		log.Error("request can't be nil")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	dec := protocol.NewDecoder(req.Body)
	replyTo := make(chan dummy.TransactionResult)
	defer close(replyTo)

	var response []byte
	for {
		var mReq model.Request
		err = dec.Decode(&mReq)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Errorf("could not decode body: %v", err)
			return ctx.NoContent(http.StatusBadRequest)
		}

		txn := dummy.MakeTransaction(mReq)
		work := dummy.MakeWork(txn).WithReplyTo(replyTo)
		v1.Node.Process(work)

		txnResp := <-replyTo
		err = txnResp.GetError()
		if err != nil {
			log.Errorf("could process request: %v", err)
			return ctx.NoContent(http.StatusBadRequest)
		}

		var mResp model.Response = txnResp.GetMsg()
		var resp bytes.Buffer

		writer := io.Writer(&resp)
		enc := protocol.NewEncoder(writer)
		err = enc.Encode(mResp)
		if err != nil {
			log.Errorf("could not encode response: %v", err)
			return ctx.NoContent(http.StatusInternalServerError)
		}

		response = append(response, resp.Bytes()...)
	}

	log.With(
		context.ResponseBody,
		fmt.Sprintf("%v", response),
	).Info("Debugging bytes")

	return ctx.Blob(http.StatusOK, "application/x-binary", response)
}
