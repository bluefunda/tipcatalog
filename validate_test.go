package tipcatalog

import "testing"

func validEmbedding() []float64 {
	v := make([]float64, EmbeddingDim)
	for i := range v {
		v[i] = float64(i) / float64(EmbeddingDim)
	}
	return v
}

func validTip() Tip {
	return Tip{
		ID:             "t1",
		Family:         "fam",
		Surfaces:       []string{SurfaceCLI},
		Render:         Render{CLI: &RenderContent{Title: "Tip", Body: "Body"}},
		Embedding:      validEmbedding(),
		CatalogVersion: "1",
	}
}

func TestValidate_OK(t *testing.T) {
	if err := Validate([]Tip{validTip()}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_MissingRequiredField(t *testing.T) {
	cases := map[string]func(*Tip){
		"id":              func(tp *Tip) { tp.ID = "" },
		"family":          func(tp *Tip) { tp.Family = "" },
		"catalog_version": func(tp *Tip) { tp.CatalogVersion = "" },
		"surfaces":        func(tp *Tip) { tp.Surfaces = nil },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			tp := validTip()
			mutate(&tp)
			if err := Validate([]Tip{tp}); err == nil {
				t.Fatalf("expected error when %s is missing", name)
			}
		})
	}
}

func TestValidate_DuplicateID(t *testing.T) {
	a, b := validTip(), validTip()
	if err := Validate([]Tip{a, b}); err == nil {
		t.Fatal("expected error for duplicate id")
	}
}

func TestValidate_WrongEmbeddingDimensionality(t *testing.T) {
	tp := validTip()
	tp.Embedding = tp.Embedding[:len(tp.Embedding)-1]
	if err := Validate([]Tip{tp}); err == nil {
		t.Fatal("expected error for wrong embedding length")
	}
}

func TestValidate_SurfaceWithoutMatchingRender(t *testing.T) {
	tp := validTip()
	tp.Surfaces = []string{SurfaceCLI, SurfaceIOS} // no render.ios set
	if err := Validate([]Tip{tp}); err == nil {
		t.Fatal("expected error for surface declared without matching render key")
	}
}

func TestValidate_UnknownSurface(t *testing.T) {
	tp := validTip()
	tp.Surfaces = []string{"carrier-pigeon"}
	if err := Validate([]Tip{tp}); err == nil {
		t.Fatal("expected error for unknown surface")
	}
}

func TestValidate_EmptyRenderContent(t *testing.T) {
	tp := validTip()
	tp.Render.CLI = &RenderContent{Title: "", Body: ""}
	if err := Validate([]Tip{tp}); err == nil {
		t.Fatal("expected error for empty render title/body")
	}
}
