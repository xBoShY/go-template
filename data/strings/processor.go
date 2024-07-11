package server

import (
	"fmt"
	"strings"

	"github.com/xboshy/go-template/data"
)

type processor struct{}

func MakeProcessor() data.Processor[Transaction, TransactionResult] {
	return &processor{}
}

func (r *processor) Process(req Transaction) TransactionResult {
	if req.HasExpired() {
		return MakeTransactionResultError(
			fmt.Errorf("Transaction has expired"),
		)
	}

	resp := req.msg
	for _, fn := range req.fns {
		switch fn {
		case ReverseStringFn:
			resp = reverseStr(resp)
		case UpperStringFn:
			resp = upperStr(resp)
		case LowerStringFn:
			resp = lowerStr(resp)
		}
	}

	return MakeTransactionResult(resp)
}

func reverseStr(s string) string {
	rune := []rune(s)
	n := len(rune)
	rune = rune[0:n]

	// Reverse
	for i := 0; i < n/2; i++ {
		rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
	}

	// Convert back to UTF-8.
	return string(rune)
}

func upperStr(s string) string {
	return strings.ToLower(s)
}

func lowerStr(s string) string {
	return strings.ToLower(s)
}
