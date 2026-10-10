package codegen

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateConvertBootstrapLUT(t *testing.T) {
	const relativePath = "internal/convert/srgb_lut_init_gen.go"
	want, err := os.ReadFile(relativePath)
	if err != nil {
		t.Fatal(err)
	}

	ctx := NewContext(Options{
		LUT:          true,
		FormatSource: true,
		BuildTags:    "custom",
	})
	ctx.Module = "github.com/thmalt/colors"
	ctx.Path = "./gen/codegen"
	ctx.Directory = t.TempDir()
	GenerateInternalPkgs(ctx)

	got, err := os.ReadFile(filepath.Join(ctx.Directory, "gen/codegen", relativePath))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(got, []byte("//go:build")) {
		t.Fatal("generator bootstrap LUT must be available under every build tag")
	}
	if !bytes.Equal(got, want) {
		t.Fatal("generator bootstrap LUT is stale; run GEN_LUT=1 go generate ./gen")
	}
}
