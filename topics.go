package tipcatalog

// Topics is the shared taxonomy tip content is tagged against (via
// DomainScope) and the client's interest vector is scored against, in
// place of a real embedding model. Order is significant: it fixes each
// topic's position in every derived Embedding vector.
//
// Appending a new topic is safe — existing embeddings just gain a new
// always-zero dimension until tips are re-tagged to use it. Reordering or
// removing a topic invalidates every previously-computed embedding and
// client-side interest vector; don't do either without a coordinated
// re-embed of the whole catalog (and bumping catalog_version).
var Topics = []string{
	"auth",
	"sessions",
	"mcp",
	"memory",
	"plugins",
	"worktree",
	"cost-budget",
	"model-selection",
	"output-format",
	"config",
	"diagnostics",
	"updates",
	"automation",
	"onboarding",
	"errors",
}

// EmbeddingDim is the length of every Tip's derived Embedding vector — one
// dimension per entry in Topics.
var EmbeddingDim = len(Topics)

// topicIndex returns topic's position in Topics, or -1 if unrecognized.
func topicIndex(topic string) int {
	for i, t := range Topics {
		if t == topic {
			return i
		}
	}
	return -1
}

// EmbeddingFromDomainScope derives a multi-hot vector over Topics from a
// tip's DomainScope: 1.0 at each recognized topic's position, 0 elsewhere.
// Unrecognized entries are ignored here — Validate separately rejects an
// unknown domain_scope topic so authoring typos surface immediately
// instead of silently vanishing from the embedding.
func EmbeddingFromDomainScope(domainScope []string) []float64 {
	v := make([]float64, EmbeddingDim)
	for _, d := range domainScope {
		if i := topicIndex(d); i >= 0 {
			v[i] = 1.0
		}
	}
	return v
}
