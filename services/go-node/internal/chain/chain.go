package chain

import (
	"errors"
	"sync"

	"digital-eval-system/services/go-node/internal/block"
	"digital-eval-system/services/go-node/internal/storage"
)

var (
	ErrChainEmpty = errors.New("chain empty")
)

// Chain provides in-memory view + storage backend
type Chain struct {
	store storage.Storage
	lock  sync.RWMutex
}

// NewChain creates chain wrapper
func NewChain(store storage.Storage) *Chain {
	return &Chain{
		store: store,
	}
}

// AppendBlock persists block and updates indices. Caller must ensure block.Header.Signature set.
func (c *Chain) AppendBlock(b *block.Block) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	// compute hash
	h, err := block.BlockHash(b)
	if err != nil {
		return "", err
	}
	// store block bytes keyed by hash
	if err := c.store.PutBlock(h, b); err != nil {
		return "", err
	}
	// update head
	if err := c.store.SetHead(h); err != nil {
		return "", err
	}
	return h, nil
}

// GetBlock retrieves block by hash
func (c *Chain) GetBlock(hash string) (*block.Block, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.store.GetBlock(hash)
}

// Head returns current head hash
func (c *Chain) Head() (string, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.store.Head()
}
