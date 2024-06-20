package dummy

type TransactionResult = *transactionResult

type transactionResult struct {
	err error
	msg string
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

func (txnRes TransactionResult) GetMsg() string {
	return txnRes.msg
}

func (txnRes TransactionResult) GetError() error {
	return txnRes.err
}
