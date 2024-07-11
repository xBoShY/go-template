package server

import (
	"time"
)

type Transaction = *transaction
type StringFn string

const (
	ReverseStringFn StringFn = "reverse"
	UpperStringFn   StringFn = "upper"
	LowerStringFn   StringFn = "lower"
)

type transaction struct {
	uuid      string
	timestamp time.Time
	ttl       time.Duration
	fns       []StringFn
	msg       string
}

func MakeTransaction(fns []StringFn, msg string) Transaction {
	return &transaction{
		timestamp: time.Now(),
		ttl:       ttlMs(int64(^uint64(0) >> 1)),
		fns:       fns,
		msg:       msg,
	}
}

func (txn Transaction) WithUuid(uuid string) Transaction {
	txn.uuid = uuid
	return txn
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

func (txn Transaction) GetUUID() string {
	return txn.uuid
}

func (txn Transaction) GetMsg() string {
	return txn.msg
}

func (txn Transaction) GetTimestamp() time.Time {
	return txn.timestamp
}

func (txn Transaction) GetTTL() time.Duration {
	return txn.ttl
}

func ttlMs(ttl int64) time.Duration {
	return time.Duration(ttl) * time.Millisecond
}
