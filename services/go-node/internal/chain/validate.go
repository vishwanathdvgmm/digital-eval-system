package chain

import (
	"crypto/rsa"
	"errors"

	"digital-eval-system/services/go-node/internal/block"
)

// ValidateChain performs a forward/backward scan to verify each block links properly and signatures (RSA) verify using provided public key loader.
// pubKeyLoader returns a *rsa.PublicKey for a given signerID.
func (c *Chain) ValidateChain(pubKeyLoader func(signerID string) (*rsa.PublicKey, error)) error {
	headHash, err := c.Head()
	if err != nil {
		return err
	}
	if headHash == "" {
		return errors.New("empty head")
	}
	it := c.store.Iterator(headHash)
	for it.Next() {
		b := it.Block()
		pub, err := pubKeyLoader(b.Header.SignerID)
		if err != nil {
			return err
		}
		if pub == nil {
			return errors.New("public key loader returned nil")
		}
		if err := b.VerifyHeaderRSA(pub); err != nil {
			return err
		}
		// validate transactions
		for _, tx := range b.Transactions {
			if !block.ValidateTransaction(&tx) {
				return errors.New("invalid transaction in block")
			}
		}
	}
	if err := it.Err(); err != nil {
		return err
	}
	return nil
}
