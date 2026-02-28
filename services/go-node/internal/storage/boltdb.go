package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"

	"digital-eval-system/services/go-node/internal/block"
)

// Storage defines required storage operations used by chain
type Storage interface {
	PutBlock(hash string, b *block.Block) error
	GetBlock(hash string) (*block.Block, error)
	ForEachBlock(fn func(*block.Block)) error
	SetHead(hash string) error
	Head() (string, error)
	Iterator(startHash string) Iterator
	Close() error
}

// Iterator for scanning blocks starting from a hash backwards (following PrevHash)
type Iterator interface {
	Next() bool
	Block() *block.Block
	Err() error
}

type boltDB struct {
	db *bolt.DB
}

// NewBoltDB opens/creates BoltDB at path
func NewBoltDB(path string, timeout time.Duration) (Storage, error) {
	if path == "" {
		return nil, errors.New("boltdb path empty")
	}
	dir := filepath.Dir(path)
	// ensure dir exists
	_ = ensureDir(dir)
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: timeout})
	if err != nil {
		return nil, err
	}
	// create buckets if not exist
	err = db.Update(func(tx *bolt.Tx) error {
		if _, e := tx.CreateBucketIfNotExists([]byte(bucketBlocks)); e != nil {
			return e
		}
		if _, e := tx.CreateBucketIfNotExists([]byte(bucketChainMeta)); e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return &boltDB{db: db}, nil
}

func (b *boltDB) Close() error {
	return b.db.Close()
}

func (b *boltDB) PutBlock(hash string, bl *block.Block) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bb := tx.Bucket([]byte(bucketBlocks))
		if bb == nil {
			return errors.New("blocks bucket missing")
		}
		// do not overwrite existing block
		if bb.Get([]byte(hash)) != nil {
			return nil
		}
		buf, err := json.Marshal(bl)
		if err != nil {
			return err
		}
		return bb.Put([]byte(hash), buf)
	})
}

func (b *boltDB) GetBlock(hash string) (*block.Block, error) {
	var bl block.Block
	err := b.db.View(func(tx *bolt.Tx) error {
		bb := tx.Bucket([]byte(bucketBlocks))
		if bb == nil {
			return errors.New("blocks bucket missing")
		}
		v := bb.Get([]byte(hash))
		if v == nil {
			return errors.New("block not found")
		}
		return json.Unmarshal(v, &bl)
	})
	if err != nil {
		return nil, err
	}
	return &bl, nil
}

// ForEachBlock iterates over every block stored in the "blocks" bucket.
func (b *boltDB) ForEachBlock(fn func(*block.Block)) error {
	return b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlocks))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(_, v []byte) error {
			if v == nil {
				return nil
			}

			var bl block.Block
			if err := json.Unmarshal(v, &bl); err != nil {
				return nil
			}

			fn(&bl)
			return nil
		})
	})
}

func (b *boltDB) SetHead(hash string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		cb := tx.Bucket([]byte(bucketChainMeta))
		if cb == nil {
			return errors.New("chain_meta bucket missing")
		}
		return cb.Put([]byte("head"), []byte(hash))
	})
}

func (b *boltDB) Head() (string, error) {
	var h string
	err := b.db.View(func(tx *bolt.Tx) error {
		cb := tx.Bucket([]byte(bucketChainMeta))
		if cb == nil {
			return errors.New("chain_meta bucket missing")
		}
		v := cb.Get([]byte("head"))
		if v == nil {
			h = ""
			return nil
		}
		h = string(v)
		return nil
	})
	if err != nil {
		return "", err
	}
	return h, nil
}

// iterator implementation
type boltIterator struct {
	db       *bolt.DB
	current  *block.Block
	currentH string
	err      error
}

func (b *boltDB) Iterator(startHash string) Iterator {
	it := &boltIterator{
		db:       b.db,
		currentH: startHash,
	}
	return it
}

func (it *boltIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if it.currentH == "" {
		return false
	}
	err := it.db.View(func(tx *bolt.Tx) error {
		bb := tx.Bucket([]byte(bucketBlocks))
		if bb == nil {
			return errors.New("blocks bucket missing")
		}
		v := bb.Get([]byte(it.currentH))
		if v == nil {
			it.current = nil
			it.currentH = ""
			return nil
		}
		var bl block.Block
		if err := json.Unmarshal(bytes.TrimSpace(v), &bl); err != nil {
			return err
		}
		it.current = &bl
		// move to previous hash
		it.currentH = bl.Header.PrevHash
		return nil
	})
	if err != nil {
		it.err = err
		return false
	}
	return it.current != nil
}

func (it *boltIterator) Block() *block.Block {
	return it.current
}

func (it *boltIterator) Err() error {
	return it.err
}

// helper: ensureDir
func ensureDir(dir string) error {
	if dir == "" {
		return nil
	}
	return nil
}
