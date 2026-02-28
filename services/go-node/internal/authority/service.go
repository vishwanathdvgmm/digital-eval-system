package authority

import (
	"context"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"

	"digital-eval-system/services/go-node/internal/block"
	"digital-eval-system/services/go-node/internal/db"
	"digital-eval-system/services/go-node/internal/storage"
)

// Service handles authority operations
type Service struct {
	db    *db.PostgresDB
	store storage.Storage
	rand  *rand.Rand
}

// NewService constructs authority service
func NewService(pg *db.PostgresDB, store storage.Storage) *Service {
	return &Service{
		db:    pg,
		store: store,
		rand:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ListPendingRequests returns pending requests
func (s *Service) ListPending(ctx context.Context) ([]db.EvalRequestRow, error) {
	return s.db.ListPendingRequests(ctx)
}

// ListHistory returns approved/rejected requests
func (s *Service) ListHistory(ctx context.Context) ([]db.EvalRequestRow, error) {
	return s.db.ListRequestHistory(ctx)
}

// ApproveRequest approves and assigns random scripts
func (s *Service) ApproveRequest(ctx context.Context, requestID int64, evaluatorID, courseID, semester, academicYear string, assignNum int) ([]string, error) {
	// find candidate scripts from BoltDB store where CourseID & Semester match and not yet assigned (we'll use store.Iterator)

	var eligible []string

	err := s.store.ForEachBlock(func(blk *block.Block) {
		for _, tx := range blk.Transactions {
			if tx.CourseID == courseID && tx.Semester == semester {
				eligible = append(eligible, tx.ScriptID)
			}
		}
	})
	if err != nil {
		return nil, err
	}

	// shuffle and pick assignNum
	s.rand.Shuffle(len(eligible), func(i, j int) { eligible[i], eligible[j] = eligible[j], eligible[i] })
	if assignNum <= 0 {
		assignNum = 5
	}
	n := assignNum
	if n > len(eligible) {
		n = len(eligible)
	}
	selected := eligible[:n]

	assigned := []string{}
	for _, sid := range selected {
		_, err := s.db.CreateAssignment(ctx, sid, evaluatorID, courseID, semester, academicYear, 0)
		if err != nil {
			logrus.Warnf("failed to create assignment for script %s: %v", sid, err)
			continue
		}
		assigned = append(assigned, sid)
	}
	// mark request approved
	if err := s.db.UpdateRequestStatus(ctx, requestID, "approved"); err != nil {
		logrus.Warnf("failed to update request status: %v", err)
	}
	return assigned, nil
}

// RejectRequest marks request rejected
func (s *Service) RejectRequest(ctx context.Context, requestID int64) error {
	return s.db.UpdateRequestStatus(ctx, requestID, "rejected")
}
