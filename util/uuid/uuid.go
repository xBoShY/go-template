package uuid

import (
	"fmt"

	"github.com/google/uuid"
)

type UUID = uuid.UUID

func NewUUID() UUID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(fmt.Errorf("unable to uuid.New() : %w", err))
	}

	return id
}
