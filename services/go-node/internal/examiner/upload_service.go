package examiner

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"digital-eval-system/services/go-node/internal/block"
	"digital-eval-system/services/go-node/internal/pybridge"
	"digital-eval-system/services/go-node/internal/rootdir"
	"digital-eval-system/services/go-node/internal/storage"
)

// Service encapsulates upload logic. It requires chain + storage + pybridge client + signer info.
type Service struct {
	chain interface {
		AppendBlock(*block.Block) (string, error)
	} // minimal interface
	store    storage.Storage
	pyClient *pybridge.Client
	signerID string
	privKey  *rsa.PrivateKey
}

// NewService constructs a new upload service.
// jwtPrivPath optional: path to RSA private key used to sign blocks (falls back to infra/certs/jwt_private.pem)
func NewService(ch interface {
	AppendBlock(*block.Block) (string, error)
}, st storage.Storage, py *pybridge.Client, signerID, jwtPrivPath string) (*Service, error) {
	if jwtPrivPath == "" {
		jwtPrivPath = rootdir.Resolve("infra/certs/jwt_private.pem")
	}
	return &Service{
		chain:    ch,
		store:    st,
		pyClient: py,
		signerID: signerID,
	}, nil
}

// ProcessUpload receives a local file path (already saved), calls python extractor, validates and writes a block.
// Returns scriptID, blockHash, any error.
func (s *Service) ProcessUpload(ctx context.Context, savedFilePath string) (*ScriptUpload, string, error) {
	// call python extract
	exOut, err := s.pyClient.Extract(ctx, savedFilePath)
	if err != nil {
		logrus.Warnf("python extract error: %v", err)
		return nil, "", fmt.Errorf("extract failed: %w", err)
	}
	if exOut == nil || exOut.Status != "success" {
		return nil, "", fmt.Errorf("extract did not return success: %+v", exOut)
	}

	// canonicalize metadata (map[string]string expected)
	meta := exOut.Metadata
	// ensure required fields exist
	usn := meta["USN"]
	course := meta["CourseID"]
	sem := meta["Semester"]

	// generate script id
	scriptID := uuid.NewString()

	// build transaction payload
	tx := block.Transaction{
		ScriptID:  scriptID,
		USN:       usn,
		CourseID:  course,
		Semester:  sem,
		CID:       exOut.PDFCid,
		Meta:      meta,
		CreatedAt: time.Now().Unix(),
		SignerID:  s.signerID,
	}

	// create block - previous hash is current chain head if available
	var prevHash string
	h, _ := s.store.Head() // ignore error if any; prevHash empty allowed
	prevHash = h

	newBlock := block.NewBlock(prevHash, []block.Transaction{tx}, s.signerID)

	// sign header if private key available
	if s.privKey != nil {
		if signErr := newBlock.SignHeaderRSA(s.privKey); signErr != nil {
			logrus.Warnf("failed to sign block header: %v - continuing to store unsigned block", signErr)
		}
	}

	// append to chain
	blockHash, err := s.chain.AppendBlock(newBlock)
	if err != nil {
		logrus.Errorf("append block failed: %v", err)
		return nil, "", fmt.Errorf("failed to persist block: %w", err)
	}

	// populate ScriptUpload record
	rec := &ScriptUpload{
		ScriptID:   scriptID,
		USN:        usn,
		CourseID:   course,
		Semester:   sem,
		CourseName: meta["CourseName"],
		PDFCid:     exOut.PDFCid,
		PDFPath:    exOut.PDFPath,
		Metadata:   meta,
		Status:     "validated",
		CreatedAt:  time.Now().UTC(),
	}

	// done
	return rec, blockHash, nil
}
