package tipcatalog

import (
	"encoding/json"
	"os"
	"sort"
	"testing"
)

// TestSchemaRequiredMatchesValidator guards against schema/tip.schema.json
// (the Swift codegen source of truth) drifting from what validate.go
// actually enforces as required. If you add/remove a required field in one,
// update the other.
func TestSchemaRequiredMatchesValidator(t *testing.T) {
	b, err := os.ReadFile("schema/tip.schema.json")
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var doc struct {
		Required []string `json:"required"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	want := []string{"catalog_version", "embedding", "family", "id", "render", "surfaces"}

	got := append([]string(nil), doc.Required...)
	sort.Strings(got)
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("schema required = %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("schema required = %v, want %v", got, want)
		}
	}
}

// TestSchemaDomainScopeEnumMatchesTopics guards against
// schema/tip.schema.json's domain_scope enum drifting from Topics
// (topics.go). If you add/remove/rename a topic, update both.
func TestSchemaDomainScopeEnumMatchesTopics(t *testing.T) {
	b, err := os.ReadFile("schema/tip.schema.json")
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var doc struct {
		Properties struct {
			DomainScope struct {
				Items struct {
					Enum []string `json:"enum"`
				} `json:"items"`
			} `json:"domain_scope"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	got := append([]string(nil), doc.Properties.DomainScope.Items.Enum...)
	want := append([]string(nil), Topics...)
	sort.Strings(got)
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("schema domain_scope enum = %v, want %v (Topics)", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("schema domain_scope enum = %v, want %v (Topics)", got, want)
		}
	}
}
