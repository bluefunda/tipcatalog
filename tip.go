// Package tipcatalog defines the shared tip/suggestion content schema
// consumed by the bluefunda CLI, iOS app, and editor plugins, plus the
// loading, validation, and signing helpers around it.
package tipcatalog

// Known surface identifiers. A Tip's Surfaces field gates which clients may
// show it; Render must carry a matching entry for each declared surface.
const (
	SurfaceCLI    = "cli"
	SurfaceIOS    = "ios"
	SurfaceVSCode = "vscode"
	SurfaceADT    = "adt"
)

// knownSurfaces is the set of Surfaces values validate.go accepts.
var knownSurfaces = map[string]bool{
	SurfaceCLI:    true,
	SurfaceIOS:    true,
	SurfaceVSCode: true,
	SurfaceADT:    true,
}

// Tip is one entry in the shared catalog.
type Tip struct {
	ID                string    `json:"id"`
	Family            string    `json:"family"`
	Surfaces          []string  `json:"surfaces"`
	DomainScope       []string  `json:"domain_scope,omitempty"`
	PersonaGate       string    `json:"persona_gate,omitempty"`
	TriggerConditions []string  `json:"trigger_conditions,omitempty"`
	MinTier           string    `json:"min_tier,omitempty"`
	Cooldown          string    `json:"cooldown,omitempty"`
	Render            Render    `json:"render"`
	DeepLink          string    `json:"deep_link,omitempty"`
	Embedding         []float64 `json:"embedding"`
	CatalogVersion    string    `json:"catalog_version"`
}

// Render holds per-surface copy for a Tip. A field is set only when the
// corresponding surface appears in the Tip's Surfaces list.
type Render struct {
	CLI    *RenderContent `json:"cli,omitempty"`
	IOS    *RenderContent `json:"ios,omitempty"`
	VSCode *RenderContent `json:"vscode,omitempty"`
	ADT    *RenderContent `json:"adt,omitempty"`
}

// RenderContent is the copy shown for one surface.
type RenderContent struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// forSurface returns the RenderContent registered for surface, or nil if
// none is set or surface is not one of the known constants.
func (r Render) forSurface(surface string) *RenderContent {
	switch surface {
	case SurfaceCLI:
		return r.CLI
	case SurfaceIOS:
		return r.IOS
	case SurfaceVSCode:
		return r.VSCode
	case SurfaceADT:
		return r.ADT
	default:
		return nil
	}
}
