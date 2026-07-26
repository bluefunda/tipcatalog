package tipcatalog

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"
)

// LoadDir reads every *.json file in dir, parses it as a Tip, and validates
// the resulting set. Tips are returned sorted by ID.
func LoadDir(dir string) ([]Tip, error) {
	return loadFS(os.DirFS(dir), ".")
}

// Compile validates tips and marshals them into the single JSON document
// published as catalog.json.
func Compile(tips []Tip) ([]byte, error) {
	if err := Validate(tips); err != nil {
		return nil, err
	}
	return json.MarshalIndent(tips, "", "  ")
}

// loadFS reads every *.json file under dir in fsys, parses each as a Tip,
// and validates the resulting set. Shared by LoadDir (os.DirFS) and
// Embedded (the go:embed FS).
func loadFS(fsys fs.FS, dir string) ([]Tip, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}

	var tips []Tip
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := e.Name()
		if dir != "." {
			path = dir + "/" + e.Name()
		}
		b, err := fs.ReadFile(fsys, path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
		var t Tip
		if err := json.Unmarshal(b, &t); err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		// Embedding is always derived from DomainScope, never hand-authored —
		// this overwrites whatever (if anything) was in the source JSON.
		t.Embedding = EmbeddingFromDomainScope(t.DomainScope)
		tips = append(tips, t)
	}

	sort.Slice(tips, func(i, j int) bool { return tips[i].ID < tips[j].ID })

	if err := Validate(tips); err != nil {
		return nil, err
	}
	return tips, nil
}
