package dummy

import (
	"fmt"

	"github.com/xboshy/go-template/data"
)

type resource struct{}

func MakeResource() data.Resource[Transaction, TransactionResult] {
	return &resource{}
}

func (r *resource) Process(req Transaction) TransactionResult {
	if req.HasExpired() {
		return MakeTransactionResultError(
			fmt.Errorf(""),
		)
	}

	return MakeTransactionResult(reverseString(req.msg))
}

func reverseString(s string) string {
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
