package regression

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/service"
	"github.com/jb843051627/kilnward/internal/store"
)

func TestBug27_CanceledCommandDoesNotRun(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db")); if err != nil { t.Fatal(err) }; defer repo.Close(); app := service.NewLab(repo); defer app.Close(); ctx := context.Background(); kiln, err := app.RegisterKiln(ctx, model.Kiln{Name: "K1", Location: "east", MaxTempC: 1500, ProbeCount: 2}); if err != nil { t.Fatal(err) }
 command, err := app.IssueCommand(ctx, model.OperatorCommand{KilnID: kiln.ID, Action: model.CommandCalibrate, RequestedBy: "operator", Reason: "daily check"}); if err != nil { t.Fatal(err) }; canceled, cancel := context.WithCancel(ctx); cancel(); _, err = app.RunCommand(canceled, command.ID); if !errors.Is(err, context.Canceled) { t.Fatalf("error = %v, want context.Canceled", err) }
}
