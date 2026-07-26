package tipcatalog

import (
	"errors"
	"fmt"
)

// Validate checks every tip in tips for required fields, valid surfaces,
// matching per-surface render copy, correct embedding dimensionality,
// recognized domain_scope topics, and duplicate IDs across the set. It
// returns the first error found.
func Validate(tips []Tip) error {
	seen := make(map[string]bool, len(tips))
	for _, t := range tips {
		if err := validateOne(t); err != nil {
			return fmt.Errorf("tip %q: %w", t.ID, err)
		}
		if seen[t.ID] {
			return fmt.Errorf("duplicate tip id %q", t.ID)
		}
		seen[t.ID] = true
	}
	return nil
}

func validateOne(t Tip) error {
	if t.ID == "" {
		return errors.New("id is required")
	}
	if t.Family == "" {
		return errors.New("family is required")
	}
	if t.CatalogVersion == "" {
		return errors.New("catalog_version is required")
	}
	if len(t.Surfaces) == 0 {
		return errors.New("surfaces must be non-empty")
	}
	if len(t.Embedding) != EmbeddingDim {
		return fmt.Errorf("embedding must have %d dimensions, got %d", EmbeddingDim, len(t.Embedding))
	}
	for _, d := range t.DomainScope {
		if topicIndex(d) < 0 {
			return fmt.Errorf("unknown domain_scope topic %q (see Topics in topics.go)", d)
		}
	}
	for _, s := range t.Surfaces {
		if !knownSurfaces[s] {
			return fmt.Errorf("unknown surface %q", s)
		}
		content := t.Render.forSurface(s)
		if content == nil {
			return fmt.Errorf("surface %q declared but render.%s is missing", s, s)
		}
		if content.Title == "" || content.Body == "" {
			return fmt.Errorf("render.%s must have both title and body", s)
		}
	}
	return nil
}
