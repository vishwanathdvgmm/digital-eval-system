package chain

import (
	"errors"

	"digital-eval-system/services/go-node/internal/block"
	"digital-eval-system/services/go-node/internal/storage"
)

// iterator wraps storage iterator to provide block access during validation scans.
type iterator struct {
	base storage.Iterator
	err  error
}

func (it *iterator) Next() bool {
	if it.err != nil {
		return false
	}
	return it.base.Next()
}

func (it *iterator) Block() *block.Block {
	return it.base.Block()
}

func (it *iterator) Err() error {
	return it.base.Err()
}

// expose constructor
func (c *Chain) Iterator(startHash string) *iterator {
	return &iterator{
		base: c.store.Iterator(startHash),
	}
}

// storage.Iterator is defined in storage package. Provide thin wrapper adapter here.
var _ = errors.New // silence unused import if any
