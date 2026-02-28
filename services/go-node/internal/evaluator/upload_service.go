package evaluator

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"digital-eval-system/services/go-node/internal/pybridge"
)

// UploadService handles uploading evaluated scripts to IPFS via pybridge.
type UploadService struct {
	pyClient *pybridge.Client
}

// NewUploadService creates a new upload service.
func NewUploadService(py *pybridge.Client) *UploadService {
	if py == nil {
		// fallback default
		py = pybridge.NewClient("http://127.0.0.1:8081", 300*time.Second)
	}
	return &UploadService{
		pyClient: py,
	}
}

// UploadEvaluatedScript sends the local file to Python service for IPFS upload.
// Returns the IPFS CID and the PDF path (if returned by extractor).
func (s *UploadService) UploadEvaluatedScript(ctx context.Context, filePath string) (string, string, error) {
	// We reuse pybridge.Extract because it handles the IPFS upload logic.
	// Even though we might not need the metadata extraction, this is the established path for "File -> IPFS".
	resp, err := s.pyClient.Extract(ctx, filePath)
	if err != nil {
		logrus.Warnf("evaluator upload: pybridge extract failed: %v", err)
		return "", "", fmt.Errorf("upload failed: %w", err)
	}

	if resp.Status != "success" {
		return "", "", fmt.Errorf("upload returned non-success status: %s (error: %s)", resp.Status, resp.Error)
	}

	// Return CID and Path
	return resp.PDFCid, resp.PDFPath, nil
}
