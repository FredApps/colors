package codegen

import "github.com/thmalt/colors/gen/codegen/writer"

func GenerateInternalPkgs(ctx *Context) {
	generateHexLUTPkg(ctx)
	generateConvertBootstrapLUT(ctx)
}

func generateConvertBootstrapLUT(ctx *Context) {
	if !ctx.Opts.genLUT() {
		return
	}

	pkg := Pkg{Name: "convert", Path: "gen/codegen/internal/convert"}
	w := newWriter(ctx)
	emitGoFile(ctx, pkg, w, "srgb_lut_init", func(w *writer.GoWriter) {
		// The generator needs its own LUTs to build, including when the public
		// convert package uses static tables. Always bootstrap these at runtime.
		w.SetBuildTags()
		genConvertPkgRuntimeInitLUT(w, 8)
		genConvertPkgRuntimeInitLUT(w, 16)
	})
}
