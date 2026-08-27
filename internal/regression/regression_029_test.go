package regression

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/store"
)

func TestBug29_AuditWriteHonorsCanceledContext(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db")); if err != nil { t.Fatal(err) }; defer repo.Close(); canceled, cancel := context.WithCancel(context.Background()); cancel()
	err = repo.RecordAudit(canceled, model.Audit{SubjectType: "kiln", SubjectID: "k1", Action: "state", Detail: "ready"})
	if !errors.Is(err, context.Canceled) { t.Fatalf("error = %v, want context.Canceled", err) }
}
