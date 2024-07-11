package server

import "github.com/xboshy/go-template/util/uuid"

type TransactionResult = *transactionResult

type transactionResult struct {
	uuid uuid.UUID
	err  error
	msg  string
}

func MakeTransactionResult(msg string) TransactionResult {
	return &transactionResult{
		msg: msg,
	}
}

func MakeTransactionResultError(err error) TransactionResult {
	return &transactionResult{
		err: err,
	}
}

func (txnRes TransactionResult) WithUuid(uuid uuid.UUID) TransactionResult {
	txnRes.uuid = uuid
	return txnRes
}

func (txnRes TransactionResult) GetUUID() uuid.UUID {
	return txnRes.uuid
}

func (txnRes TransactionResult) GetMsg() string {
	return txnRes.msg
}

func (txnRes TransactionResult) GetError() error {
	return txnRes.err
}
