# Adding tip content

This is the practical, step-by-step guide for authoring real tips — as opposed to
[CONTRIBUTING.md](CONTRIBUTING.md), which covers code changes to the Go package itself.

## How it fits together

Each tip is one JSON file under [`tips/`](tips/). You author `id`, `family`, `surfaces`,
`domain_scope`, `render` copy, and the other fields described below — **you never write the
`embedding` field**. It's derived automatically from `domain_scope` when the catalog loads (see
`EmbeddingFromDomainScope` in [`topics.go`](topics.go)): a multi-hot vector, one dimension per
entry in [`Topics`](topics.go). Anything you put in `embedding` in a source file is silently
overwritten.

This means the entire "does this tip get shown to the right person" mechanism boils down to
picking accurate `domain_scope` tags — that's the one piece of authoring that actually drives
targeting, everything else is copy.

## The topic taxonomy

The current list, in order (order matters — see the comment in `topics.go` before changing it):

```
auth, sessions, mcp, memory, plugins, worktree, cost-budget, model-selection,
output-format, config, diagnostics, updates, automation, onboarding, errors
```

Tag a tip with every topic it's genuinely relevant to. Most tips need just one or two — don't
over-tag; a tip tagged with everything ranks no better than one tagged with nothing; `Validate`
rejects any topic not in this list, so a typo fails fast.

## Adding a tip

1. **Create the file.** One JSON file per tip in `tips/`, filename matching `id`:

   ```bash
   cd ~/src/tipcatalog
   git checkout -b content/my-new-tip
   ```

   ```json
   // tips/my-new-tip.json
   {
     "id": "my-new-tip",
     "family": "some-family",
     "surfaces": ["cli"],
     "domain_scope": ["mcp"],
     "persona_gate": "",
     "trigger_conditions": [],
     "min_tier": "free",
     "cooldown": "24h",
     "render": {
       "cli": {
         "title": "Tip",
         "body": "Terse, monospace-safe, one clear action."
       }
     },
     "deep_link": "",
     "catalog_version": "1"
   }
   ```

   Field notes:
   - `family` groups related tips for anti-annoyance suppression — dismissing one suppresses the
     whole family, so give genuinely-related tips the same family name.
   - `surfaces` gates which clients may show it (`cli`, `ios`, `vscode`, `adt`); `render` must
     have a matching key for every surface listed.
   - `persona_gate` is optional (e.g. `"new_user"`); leave `""` for no gate.
   - `min_tier` is currently informational only — no client has real subscription-tier data on
     the hot path yet (see `AGENTS.md`).
   - `cooldown` is a Go duration string (`"24h"`, `"72h"`); a client falls back to a 1h floor if
     you leave it empty.
   - CLI copy should be terse and monospace-safe — this repo's spec explicitly rules out
     LLM-rewritten copy at render time, so write it exactly as it should appear.

2. **Validate.**

   ```bash
   go test ./...
   ```

   This is the same suite CI runs: required fields, per-surface render coverage, a recognized
   `domain_scope` (typos in topic names fail here), duplicate IDs, and the schema/`Topics` drift
   check.

3. **Sanity-check the derived embedding** (optional, but useful while you're learning the
   mechanism):

   ```bash
   cat <<'EOF' > /tmp/check_embedding.go
   package main

   import (
       "fmt"
       tipcatalog "github.com/bluefunda/tipcatalog"
   )

   func main() {
       tips, err := tipcatalog.LoadDir("tips")
       if err != nil {
           panic(err)
       }
       for _, t := range tips {
           fmt.Println(t.ID, t.DomainScope, t.Embedding)
       }
   }
   EOF
   go run /tmp/check_embedding.go
   ```

4. **Open a PR.**

   ```bash
   git add tips/my-new-tip.json
   git commit -m "feat: add my-new-tip"
   git push -u origin content/my-new-tip
   gh pr create --title "feat: add my-new-tip" --body "..."
   ```

5. **What happens on merge.** Nothing publishes immediately — release-please batches merged
   commits into a version-bump PR. Once *that's* merged and tagged, `publish-manifest.yml`
   compiles every tip in `tips/` into `catalog.json` (deriving embeddings fresh at that point),
   signs it, and attaches `catalog.json` + `catalog.json.sig` to the GitHub Release. Clients pick
   it up on their next opportunistic refresh (at most once per 24h).

## Adding a new topic

Rarer, and has consequences outside this repo, so don't do it casually:

1. Add the entry to `Topics` in `topics.go` — **append, don't insert or reorder** (see the
   comment there for why: reordering invalidates every previously-computed embedding and every
   client's in-progress interest vector).
2. Add the same string to `schema/tip.schema.json`'s `domain_scope` enum.
3. Run `go test ./...` — `TestSchemaDomainScopeEnumMatchesTopics` fails loudly if the two drift.
4. Tell whoever maintains the client-side command→topic mapping (currently
   `commandTopics`/`flagTopics` in `bluefunda-ai`'s `internal/tips/signal.go`) that a new topic
   exists — a topic nothing ever maps to on the client side will never accumulate interest, so
   tips tagged with it will only ever score by whatever *other* topics they also carry.

## Removing or renaming a topic

Don't, without a coordinated re-embed: every tip's embedding and every client's stored interest
vector encodes topic *position*, not name. If you must, bump `catalog_version` and communicate it
as a breaking change.
