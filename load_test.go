package tipcatalog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDir_RoundTrip(t *testing.T) {
	tips, err := LoadDir("tips")
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	if len(tips) == 0 {
		t.Fatal("expected at least one tip in tips/")
	}
	for i := 1; i < len(tips); i++ {
		if tips[i-1].ID >= tips[i].ID {
			t.Fatalf("tips not sorted by id: %s >= %s", tips[i-1].ID, tips[i].ID)
		}
	}
}

func TestLoadDir_RejectsInvalidTip(t *testing.T) {
	dir := t.TempDir()
	bad := validTip()
	bad.ID = "" // missing required field
	b, err := json.Marshal(bad)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bad.json"), b, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if _, err := LoadDir(dir); err == nil {
		t.Fatal("expected LoadDir to reject an invalid tip")
	}
}

func TestCompile_RoundTrip(t *testing.T) {
	tips, err := LoadDir("tips")
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	compiled, err := Compile(tips)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	var roundTripped []Tip
	if err := json.Unmarshal(compiled, &roundTripped); err != nil {
		t.Fatalf("unmarshal compiled catalog: %v", err)
	}
	if len(roundTripped) != len(tips) {
		t.Fatalf("got %d tips after round-trip, want %d", len(roundTripped), len(tips))
	}
	if err := Validate(roundTripped); err != nil {
		t.Fatalf("round-tripped tips failed validation: %v", err)
	}
}

func TestCompile_RejectsInvalidSet(t *testing.T) {
	tp := validTip()
	tp.ID = ""
	if _, err := Compile([]Tip{tp}); err == nil {
		t.Fatal("expected Compile to reject an invalid tip set")
	}
}

func TestEmbedded_MatchesLoadDir(t *testing.T) {
	fromDir, err := LoadDir("tips")
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	fromEmbed, err := Embedded()
	if err != nil {
		t.Fatalf("Embedded: %v", err)
	}
	if len(fromEmbed) != len(fromDir) {
		t.Fatalf("Embedded returned %d tips, LoadDir returned %d", len(fromEmbed), len(fromDir))
	}
}
