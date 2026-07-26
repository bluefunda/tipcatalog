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

func TestLoadDir_DerivesEmbeddingFromDomainScope(t *testing.T) {
	tips, err := LoadDir("tips")
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	for _, tp := range tips {
		want := EmbeddingFromDomainScope(tp.DomainScope)
		if len(tp.Embedding) != len(want) {
			t.Fatalf("tip %q: embedding len = %d, want %d", tp.ID, len(tp.Embedding), len(want))
		}
		for i := range want {
			if tp.Embedding[i] != want[i] {
				t.Fatalf("tip %q: embedding does not match its own domain_scope-derived value at index %d", tp.ID, i)
			}
		}
	}
}

func TestLoadDir_IgnoresHandAuthoredEmbedding(t *testing.T) {
	dir := t.TempDir()
	tp := validTip()
	tp.DomainScope = []string{"auth"}
	tp.Embedding = make([]float64, EmbeddingDim) // wrong: doesn't match domain_scope
	for i := range tp.Embedding {
		tp.Embedding[i] = 99 // deliberately bogus, to prove it gets overwritten
	}
	b, err := json.Marshal(tp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "t.json"), b, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	loaded, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	want := EmbeddingFromDomainScope([]string{"auth"})
	for i := range want {
		if loaded[0].Embedding[i] != want[i] {
			t.Fatalf("expected hand-authored embedding to be overwritten by the domain_scope-derived value, got %v", loaded[0].Embedding)
		}
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
