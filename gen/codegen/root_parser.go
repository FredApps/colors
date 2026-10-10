package codegen

import (
	"math"
	"strings"

	"github.com/thmalt/colors/gen/codegen/data"
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgParser(ctx *Context, w *writer.GoWriter) {
	w.Separate()
	w.Method("p *parser", "parseColorSpace")
	w.FuncParams("originColor Color")
	w.FuncResults("Color, error")
	w.FuncBody()
	w.LineWriteln("var buf [", ctx.MaxCssNameLength+1, "]byte")
	w.LineWriteln("n := copy(buf[:], p.t.Text())")
	w.LineWriteln("ascii.ToLowerBytes(buf[:n])")

	w.Separate()
	w.LineWriteln("start, end := p.t.TokenStart(), p.t.TokenEnd()")
	w.LineWriteln("p.next()")

	w.Separate()
	w.Switch("string(buf[:n])")

	sub := w.SubWriter()
	for _, s := range ctx.BuiltSpaces {
		if s.UseGenericColorFunction {
			sub.Reset()
			sub.Write('"', strings.ToLower(s.CssName), '"')
			for _, alias := range s.CssAliases {
				sub.Write(", ", '"', strings.ToLower(alias), '"')
			}
			w.Case(sub.Bytes())

			switch s.ParseKind {
			case model.Parse01:
				w.Return("p.parseColorSpace01(originColor, ", ctx.SpacePkg.Join(s.Name), ", ", ctx.SpacePkg.Join(chIdentSetType+s.ChannelIdentFlagsName()), ")")

			case model.ParseLabLike:
				scale := s.Channels[1].Max
				if ps := s.Channels[1].PercentScale; ps != 0 {
					scale = ps
				}
				w.Return("p.parseLabLike(originColor, ", ctx.SpacePkg.Join(s.Name), ", ", s.Channels[0].Max, ", ", scale, ")")

			case model.ParseLchLike:
				scale := s.Channels[1].Max
				if ps := s.Channels[1].PercentScale; ps != 0 {
					scale = ps
				}
				w.Return("p.parseLchLike(originColor, ", ctx.SpacePkg.Join(s.Name), ", ", s.Channels[0].Max, ", ", scale, ")")
			}
		}
	}
	w.Default()
	w.LineWriteln("p.setErr(errInvalidColor)")
	w.LineWriteln(`p.setErrMsg("unknown color space")`)
	w.LineWriteln("p.setErrTextAt(start, end)")
	w.Return("Color{}, p.fail()")
	w.End()
	w.End()
}

func genRootPkgNamedColor(ctx *Context, w *writer.GoWriter) {
	w.Separate()
	w.Func("lookupNamedColor")
	w.FuncParams("name string")
	w.FuncResults("Color, bool")
	w.FuncBody()
	w.LineWriteln("var buf [", data.MaxNamedColorLength+1, "]byte")
	w.LineWriteln("n := copy(buf[:], name)")
	w.LineWriteln("ascii.ToLowerBytes(buf[:n])")
	w.Separate()
	w.Switch("string(buf[:n])")

	for _, nc := range ctx.NamedColors {
		w.Case('"', nc.Lower(), '"')
		w.ReturnInline("Color{")
		w.In()
		w.LineWriteln("space: ", ctx.SpacePkg.Join("Srgb"), ',')
		w.LineWriteln("c1: ", float64(nc.RGBA[0])/math.MaxUint8, ',')
		w.LineWriteln("c2: ", float64(nc.RGBA[1])/math.MaxUint8, ',')
		w.LineWriteln("c3: ", float64(nc.RGBA[2])/math.MaxUint8, ',')
		w.LineWriteln("alpha: ", float64(nc.RGBA[3])/math.MaxUint8, ',')
		w.Out()
		w.LineWriteln("}, true")
	}

	w.Default()
	w.Return("Color{}, false")
	w.End()
	w.End()
}
