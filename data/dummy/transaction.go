package dummy

import "time"

type Transaction = *transaction

type transaction struct {
	timestamp time.Time
	ttl       time.Duration
	msg       string
}

func MakeTransaction(msg string) Transaction {
	return &transaction{
		timestamp: time.Now(),
		ttl:       ttlMs(int64(^uint64(0) >> 1)),
		msg:       msg,
	}
}

func (txn Transaction) WithTimestamp(timestamp time.Time) Transaction {
	txn.timestamp = timestamp
	return txn
}

func (txn Transaction) WithTTL(ttl int64) Transaction {
	txn.ttl = ttlMs(ttl)
	return txn
}

func (txn Transaction) HasExpired() bool {
	return txn.timestamp.Add(txn.ttl).Before(time.Now())
}

func (txn Transaction) GetMsg() string {
	return txn.msg
}

func ttlMs(ttl int64) time.Duration {
	return time.Duration(ttl) * time.Millisecond
}
