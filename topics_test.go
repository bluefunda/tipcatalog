package tipcatalog

import "testing"

func TestEmbeddingFromDomainScope_MultiHot(t *testing.T) {
	v := EmbeddingFromDomainScope([]string{"auth", "errors"})
	if len(v) != EmbeddingDim {
		t.Fatalf("len(v) = %d, want %d", len(v), EmbeddingDim)
	}

	authIdx := topicIndex("auth")
	errorsIdx := topicIndex("errors")
	if authIdx < 0 || errorsIdx < 0 {
		t.Fatal("expected 'auth' and 'errors' to be known topics")
	}

	for i, x := range v {
		want := 0.0
		if i == authIdx || i == errorsIdx {
			want = 1.0
		}
		if x != want {
			t.Fatalf("v[%d] = %v, want %v", i, x, want)
		}
	}
}

func TestEmbeddingFromDomainScope_IgnoresUnknownTopics(t *testing.T) {
	v := EmbeddingFromDomainScope([]string{"not-a-real-topic"})
	for i, x := range v {
		if x != 0 {
			t.Fatalf("v[%d] = %v, want 0 (unknown topics contribute nothing)", i, x)
		}
	}
}

func TestEmbeddingFromDomainScope_EmptyIsZeroVector(t *testing.T) {
	v := EmbeddingFromDomainScope(nil)
	if len(v) != EmbeddingDim {
		t.Fatalf("len(v) = %d, want %d", len(v), EmbeddingDim)
	}
	for i, x := range v {
		if x != 0 {
			t.Fatalf("v[%d] = %v, want 0", i, x)
		}
	}
}

func TestTopicIndex_Unknown(t *testing.T) {
	if topicIndex("not-a-real-topic") != -1 {
		t.Fatal("expected -1 for an unrecognized topic")
	}
}

func TestTopics_NoDuplicates(t *testing.T) {
	seen := make(map[string]bool, len(Topics))
	for _, topic := range Topics {
		if seen[topic] {
			t.Fatalf("duplicate topic %q in Topics", topic)
		}
		seen[topic] = true
	}
}
