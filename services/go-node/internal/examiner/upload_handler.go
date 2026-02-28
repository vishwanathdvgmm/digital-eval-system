package examiner

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"digital-eval-system/services/go-node/internal/core"
	"digital-eval-system/services/go-node/internal/rootdir"
)

// UploadHandler expects multipart/form-data with field "file".
//
// Example curl:
// curl -F "file=@/path/to/script.pdf" http://127.0.0.1:8443/api/v1/examiner/upload
type UploadHandler struct {
	service *Service
	// directory where uploaded files are saved
	uploadDir string
	// temp retention in seconds for manual cleanup
	retention time.Duration
}

func NewUploadHandler(svc *Service, uploadDir string, retentionSeconds int) *UploadHandler {
	if uploadDir == "" {
		uploadDir = rootdir.Resolve("data/uploads")
	}
	if retentionSeconds <= 0 {
		retentionSeconds = 24 * 3600
	}
	return &UploadHandler{
		service:   svc,
		uploadDir: uploadDir,
		retention: time.Duration(retentionSeconds) * time.Second,
	}
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// trace id from context
	trace := core.TraceIDFromContext(r.Context())

	// limit upload size (e.g., 20MB)
	const maxUploadSize = 20 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		logrus.WithField("trace", trace).Warnf("parse multipart error: %v", err)
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		logrus.WithField("trace", trace).Warnf("missing file field: %v", err)
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// ensure filename safe
	filename := filepath.Base(header.Filename)
	savedPath, err := SaveUploadedFile(h.uploadDir, filename, file)
	if err != nil {
		logrus.WithField("trace", trace).Errorf("failed save uploaded file: %v", err)
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}
	logrus.WithFields(logrus.Fields{"trace": trace, "saved": savedPath}).Info("file saved")

	// we need to re-open file for passing to service (because original reader is exhausted)
	fh, err := os.Open(savedPath)
	if err != nil {
		logrus.WithField("trace", trace).Errorf("failed reopen saved file: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		RemoveFile(savedPath)
		return
	}
	_ = fh.Close() // service uses path

	ctx, cancel := context.WithTimeout(r.Context(), 300*time.Second)
	defer cancel()

	rec, blockHash, err := h.service.ProcessUpload(ctx, savedPath)
	if err != nil {
		logrus.WithField("trace", trace).Warnf("process upload failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		// keep file for debugging — do not delete automatically on error
		return
	}

	// success response
	resp := map[string]interface{}{
		"script_id":  rec.ScriptID,
		"block_hash": blockHash,
		"pdf_cid":    rec.PDFCid,
		"pdf_path":   rec.PDFPath,
		"metadata":   rec.Metadata,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)

	// schedule background cleanup (best effort)
	RemoveFile(savedPath)
}
