package history

import (
	"path/filepath"
	"testing"
	"time"
)

func TestDBOperations(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_history.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	now := time.Now().Truncate(time.Millisecond)

	move1 := Move{
		Source:      "/downloads/doc1.pdf",
		Destination: "/documents/doc1.pdf",
		Rule:        "Documents",
		Action:      "moved",
		MovedAt:     now,
	}

	appended1, err := db.Append(move1)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}
	if appended1.ID <= 0 {
		t.Fatalf("expected inserted ID > 0, got %d", appended1.ID)
	}

	moves, err := db.Latest(10)
	if err != nil {
		t.Fatalf("Latest failed: %v", err)
	}
	if len(moves) != 1 {
		t.Fatalf("expected 1 move in Latest, got %d", len(moves))
	}

	pending, err := db.Pending(10)
	if err != nil {
		t.Fatalf("Pending failed: %v", err)
	}
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending move, got %d", len(pending))
	}

	found, ok, err := db.FindPending(appended1.ID)
	if err != nil {
		t.Fatalf("FindPending failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected to find pending move by ID %d", appended1.ID)
	}
	if found.ID != appended1.ID {
		t.Errorf("expected ID %d, got %d", appended1.ID, found.ID)
	}

	err = db.MarkUndone(appended1.ID, time.Now())
	if err != nil {
		t.Fatalf("MarkUndone failed: %v", err)
	}

	pendingAfter, err := db.Pending(10)
	if err != nil {
		t.Fatalf("Pending after MarkUndone failed: %v", err)
	}
	if len(pendingAfter) != 0 {
		t.Errorf("expected 0 pending moves after undo, got %d", len(pendingAfter))
	}
}
