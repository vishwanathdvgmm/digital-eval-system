package evaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"digital-eval-system/services/go-node/internal/core"
	"digital-eval-system/services/go-node/internal/rootdir"
)

// UploadHandler handles file uploads for evaluators.
type UploadHandler struct {
	service   *UploadService
	uploadDir string
}

func NewUploadHandler(svc *UploadService, uploadDir string) *UploadHandler {
	if uploadDir == "" {
		uploadDir = rootdir.Resolve("data/uploads")
	}
	// ensure dir exists
	_ = os.MkdirAll(uploadDir, 0755)
	return &UploadHandler{
		service:   svc,
		uploadDir: uploadDir,
	}
}

// ServeHTTP handles POST /api/v1/evaluator/upload
// Expects multipart/form-data with "file" field.
func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	trace := core.TraceIDFromContext(r.Context())

	// 1. Parse Multipart
	const maxUploadSize = 20 << 20 // 20 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 2. Save to temp local file
	filename := filepath.Base(header.Filename)
	savedPath := filepath.Join(h.uploadDir, fmt.Sprintf("eval_%d_%s", time.Now().UnixNano(), filename))

	dst, err := os.Create(savedPath)
	if err != nil {
		logrus.Errorf("failed to create temp file: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(dst, file); err != nil {
		dst.Close()
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}
	dst.Close()

	// cleanup on exit (best effort)
	defer os.Remove(savedPath)

	// 3. Call Service
	ctx, cancel := context.WithTimeout(r.Context(), 300*time.Second)
	defer cancel()

	cid, pdfPath, err := h.service.UploadEvaluatedScript(ctx, savedPath)
	if err != nil {
		logrus.WithField("trace", trace).Warnf("evaluator upload failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// 4. Response
	resp := map[string]string{
		"cid":      cid,
		"pdf_path": pdfPath,
		"status":   "uploaded",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
