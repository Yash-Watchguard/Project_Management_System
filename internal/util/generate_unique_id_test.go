package util

import (
	"testing"
	"strings"

	"github.com/google/uuid"

)


func TestGenerateUniqueUUID(t *testing.T) {
	uuid1 := GenerateUniqueUUID()
	uuid2 := GenerateUniqueUUID()

	// Check if UUIDs are not empty
	if uuid1 == "" || uuid2 == "" {
		t.Fatal("expected non-empty UUIDs")
	}

	// Check if UUIDs are unique
	if uuid1 == uuid2 {
		t.Fatal("expected unique UUIDs, got duplicates")
	}

	// Check if UUIDs are valid
	if _, err := uuid.Parse(uuid1); err != nil {
		t.Fatalf("invalid UUID generated: %v", err)
	}
	if _, err := uuid.Parse(uuid2); err != nil {
		t.Fatalf("invalid UUID generated: %v", err)
	}

	// Optional: check format contains hyphens
	if !strings.Contains(uuid1, "-") || !strings.Contains(uuid2, "-") {
		t.Fatal("expected UUIDs to contain hyphens")
	}
}
