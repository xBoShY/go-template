package server

type Work = *work
type work struct {
	request Transaction
	replyTo chan TransactionResult
}

func MakeWork(request Transaction) Work {
	return &work{
		request: request,
		replyTo: nil,
	}
}

func (w Work) WithReplyTo(replyTo chan TransactionResult) Work {
	w.replyTo = replyTo
	return w
}

func (w Work) Request() Transaction {
	return w.request
}

func (w Work) Reply(r TransactionResult) {
	if w.replyTo == nil {
		return
	}
	w.replyTo <- r
}
