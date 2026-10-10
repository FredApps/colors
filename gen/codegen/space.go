package codegen

import (
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateSpacePkg(ctx *Context) {
	pkg := ctx.SpacePkg
	if pkg.Name == "" {
		return
	}

	w := newWriter(ctx)

	emitGoFile(ctx, pkg, w, pkg.Name, func(w *writer.GoWriter) {
		w.Import("strconv")

		genSpacePkgSpace(ctx, w)
	})

	emitGoFile(ctx, pkg, w, pkg.Name+"_info", func(w *writer.GoWriter) {
		genSpacePkgSpaceInfo(ctx, w)
	})

	emitGoFile(ctx, pkg, w, pkg.Name+"_table", func(w *writer.GoWriter) {
		genSpacePkgTables(ctx, w)
	})

	emitGoFile(ctx, pkg, w, "channel", func(w *writer.GoWriter) {
		w.Import(
			"strconv",
			"/internal/ascii",
		)
		genSpacePkgChannelIdent(ctx, w)
	})

	emitGoFile(ctx, pkg, w, "whitepoint", func(w *writer.GoWriter) {
		genSpacePkgWhitePoint(ctx, w)
	})
}

func genSpacePkgSpace(ctx *Context, w *writer.GoWriter) {
	w.LineWriteln("type Space uint", uintBitsForValue(len(ctx.BuiltSpaces)))

	w.Separate()
	w.BeginGroup("const ")
	w.LineWriteln("InvalidSpace Space = iota")

	w.Separate()
	var firstSpace, lastSpace string

	var aliasSpaces []*model.Space

	for i, space := range ctx.BuiltSpaces {
		if i == 0 {
			firstSpace = space.Name
		}

		if space.Description != "" {
			if ctx.Opts.SeparateAfterComment {
				w.Separate()
			}

			w.Comment(space.Description)
		}
		w.LineWriteln(space.Name)

		if len(space.Aliases) > 0 {
			aliasSpaces = append(aliasSpaces, space)
		}

		lastSpace = space.Name
	}

	w.Separate()
	w.Comment("SpaceCount is the number of defined spaces, including InvalidSpace.")
	w.LineWriteln("SpaceCount")

	w.Separate()
	w.Comment("FirstSpace is the first valid color space.")
	w.LineWriteln("FirstSpace = ", firstSpace)

	w.Separate()
	w.Comment("LastSpace is the last valid color space.")
	w.LineWriteln("LastSpace = ", lastSpace)

	w.End()

	w.Separate()
	w.BeginGroup("const ")
	for _, space := range aliasSpaces {
		for _, alias := range space.Aliases {
			w.Comment(alias, " is an alias for [", space.Name, "].")
			w.LineWrite(alias, " = ", space.Name)
		}
	}
	w.End()

	w.Separate()
	w.Method("s Space", "String")
	w.FuncResults("string")
	w.FuncBody()

	w.Switch("s")
	w.Case("InvalidSpace")
	w.Return(`"Invalid"`)

	for _, space := range ctx.BuiltSpaces {
		w.Case(space.Name)
		w.Return('"', space.Name, '"')
	}

	w.Default()

	w.Return(`"Space(" + strconv.FormatUint(uint64(s), 10) + ")"`)

	w.End()

	w.End()
}
