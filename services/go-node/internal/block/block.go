package block

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"time"
)

// BlockHeader contains the immutable header fields
type BlockHeader struct {
	PrevHash   string `json:"prev_hash"`
	Timestamp  int64  `json:"timestamp"`
	MerkleRoot string `json:"merkle_root"`
	SignerID   string `json:"signer_id"`
	Signature  []byte `json:"signature"` // signature of header bytes
}

// Transaction represents the canonical transaction/metadata stored in block body.
type Transaction struct {
	ScriptID      string            `json:"script_id"`
	USN           string            `json:"usn"`
	CourseID      string            `json:"course_id"`
	Semester      string            `json:"semester"`
	AcademicYear  string            `json:"academic_year"`
	CourseCredits int               `json:"credits"`
	CID           string            `json:"cid"` // IPFS CID
	Meta          map[string]string `json:"meta,omitempty"`
	CreatedAt     int64             `json:"created_at"`
	SignerID      string            `json:"signer_id"`
	ExtraSig      []byte            `json:"extra_sig,omitempty"`
}

// Block contains header + body
type Block struct {
	Header       BlockHeader   `json:"header"`
	Transactions []Transaction `json:"transactions"`
}

// NewBlock creates a block with provided prevHash and transactions. Signature left empty until Sign() call.
func NewBlock(prevHash string, txs []Transaction, signerID string) *Block {
	return &Block{
		Header: BlockHeader{
			PrevHash:   prevHash,
			Timestamp:  time.Now().Unix(),
			MerkleRoot: computeMerkleRoot(txs),
			SignerID:   signerID,
			Signature:  nil,
		},
		Transactions: txs,
	}
}

func (b *Block) headerBytes() ([]byte, error) {
	// produce canonical header bytes for signing/hashing (excluding Signature)
	h := struct {
		PrevHash   string `json:"prev_hash"`
		Timestamp  int64  `json:"timestamp"`
		MerkleRoot string `json:"merkle_root"`
		SignerID   string `json:"signer_id"`
	}{
		PrevHash:   b.Header.PrevHash,
		Timestamp:  b.Header.Timestamp,
		MerkleRoot: b.Header.MerkleRoot,
		SignerID:   b.Header.SignerID,
	}
	return json.Marshal(h)
}

// SignHeaderRSA computes RSA-SHA256 signature over header bytes and sets Header.Signature
func (b *Block) SignHeaderRSA(priv *rsa.PrivateKey) error {
	if priv == nil {
		return errors.New("private key nil")
	}
	hb, err := b.headerBytes()
	if err != nil {
		return err
	}
	hash := sha256.Sum256(hb)
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash[:])
	if err != nil {
		return err
	}
	b.Header.Signature = sig
	return nil
}

func (b *Block) VerifyHeaderRSA(pub *rsa.PublicKey) error {
	if pub == nil {
		return errors.New("public key nil")
	}
	hb, err := b.headerBytes()
	if err != nil {
		return err
	}
	hash := sha256.Sum256(hb)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], b.Header.Signature)
}
