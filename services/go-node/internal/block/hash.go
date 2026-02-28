package block

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// BlockHash returns SHA256 of header bytes + merkle root (canonical)
func BlockHash(b *Block) (string, error) {
	hb, err := b.headerBytes()
	if err != nil {
		return "", err
	}
	// include serialized transactions to avoid collisions
	txb, err := json.Marshal(b.Transactions)
	if err != nil {
		return "", err
	}
	combined := append(hb, txb...)
	sum := sha256.Sum256(combined)
	return hex.EncodeToString(sum[:]), nil
}

// computeMerkleRoot computes a simple merkle root: hash of concatenated tx hashes.
// For Phase1 deterministic and simple; replace with real merkle tree in Phase2/3.
func computeMerkleRoot(txs []Transaction) string {
	if len(txs) == 0 {
		zero := sha256.Sum256([]byte{})
		return hex.EncodeToString(zero[:])
	}
	var agg []byte
	for _, tx := range txs {
		b, _ := json.Marshal(tx)
		h := sha256.Sum256(b)
		agg = append(agg, h[:]...)
	}
	root := sha256.Sum256(agg)
	return hex.EncodeToString(root[:])
}
