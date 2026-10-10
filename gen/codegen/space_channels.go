package codegen

import (
	"slices"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

const (
	chIdentType    = "ChIdent"
	chIdentSetType = "ChIdentSet"
)

func genSpacePkgChannelIdent(ctx *Context, w *writer.GoWriter) {
	identMap := make(map[string]string)
	sub := w.SubWriter()
	groups := newGroupSpace()

	none := "none"
	identMap[none] = chIdentName(none)

	alpha := "alpha"
	identMap[alpha] = chIdentName(alpha)

	for i, s := range ctx.BuiltSpaces {
		sub.Reset()
		for j, ident := range s.ChannelIdent() {
			lid := strings.ToLower(ident)
			chId := chIdentName(lid)
			identMap[lid] = chId

			if j > 0 {
				sub.Write(" | ")
			}
			sub.Write("1<<", chId)
		}
		sub.Write(" | ", "1<<", identMap[alpha])
		groups.AppendExtra(string(sub.Bytes()), s.ChannelCount(), i, s, chIdentSetType+s.ChannelIdentFlagsName())
	}

	flagsGroups := groups.Slice()

	var identList = make([]string, 0, len(identMap))
	identList = append(identList, none, alpha)

	for k := range identMap {
		if k != alpha && k != none {
			identList = append(identList, k)
		}
	}

	slices.SortFunc(identList[2:], func(a, b string) int {
		if n := len(a) - len(b); n != 0 {
			return n
		}
		return strings.Compare(a, b)
	})

	bitLen := uintBits(max(32, len(identList)))

	w.Separate()
	w.LineWriteln("const MaxChannelCount = ", ctx.MaxChannelCount)
	w.Separate()
	w.BeginGroup("type")
	w.LineWriteln(chIdentType, " uint8")
	w.LineWriteln(chIdentSetType, " uint", bitLen)
	w.End()

	w.Separate()
	w.BeginGroup("const")
	for i, ident := range identList {
		w.LineWrite(identMap[ident])
		if i == 0 {
			w.Write(" ", chIdentType, " = iota")
		}
		if ident == none {
			ident = ""
		}
		w.InlineComment(ident)
	}
	w.End()

	w.Separate()
	w.BeginGroup("const")
	w.LineWrite(chIdentSetType, toUpperCaseFirstChar(alpha), " ", chIdentSetType, " = 1 << ", identMap[alpha])
	w.InlineComment(alpha)

	flagsGroupsSorted := slices.Clone(flagsGroups)
	slices.SortStableFunc(flagsGroupsSorted, func(a, b groupSpaceValue) int {
		if n := len(a.Extra) - len(b.Extra); n != 0 {
			return n
		}
		return strings.Compare(a.Extra, b.Extra)
	})

	for _, group := range flagsGroupsSorted {
		sub.Reset()
		for i, ident := range group.Spaces[0].ChannelIdent() {
			if i > 0 {
				sub.Write(", ")
			}
			sub.Write(strings.ToLower(ident))
		}
		sub.Write(", ", alpha)

		w.LineWrite(group.Extra, " ", chIdentSetType, " = ", group.Key)
		w.InlineComment(sub.Bytes())
	}
	w.End()

	var identLenMap = make(map[int][]string)
	for ident := range identMap {
		if ident == none {
			continue
		}
		identLenMap[len(ident)] = append(identLenMap[len(ident)], ident)
	}
	var identLenList = make([]int, 0, len(identLenMap))
	for k := range identLenMap {
		identLenList = append(identLenList, k)
	}
	slices.Sort(identLenList)

	w.Separate()
	// func (ch ChIdent) Set() ChIdentSet
	w.Method("ch "+chIdentType, "Set")
	w.FuncResults(chIdentSetType)
	w.FuncBody()
	w.Return("1 << (ch&", bitLen-1, ')')
	w.End()

	w.Separate()
	// func (s ChIdentSet) Has(ch ChIdent) bool
	w.Method("s "+chIdentSetType, "Has")
	w.FuncParams("ch ", chIdentType)
	w.FuncResults("bool")
	w.FuncBody()
	w.Return("s&(1<<(ch&", bitLen-1, ")) != 0")
	w.End()

	w.Separate()
	// func (s ChIdentSet) HasAny(o ChIdentSet) bool
	w.Method("s "+chIdentSetType, "HasAny")
	w.FuncParams("o ", chIdentSetType)
	w.FuncResults("bool")
	w.FuncBody()
	w.Return("s&o != 0")
	w.End()

	w.Separate()
	// func (s ChIdentSet) HasAll(o ChIdentSet) bool
	w.Method("s "+chIdentSetType, "HasAll")
	w.FuncParams("o ", chIdentSetType)
	w.FuncResults("bool")
	w.FuncBody()
	w.Return("s&o == o")
	w.End()

	w.Separate()
	// func (ch ChIdent) String() string
	w.Method("ch "+chIdentType, "String")
	w.FuncResults("string")
	w.FuncBody()
	w.Switch("ch")
	for _, ident := range identList {
		w.Case(identMap[ident])
		if ident == none {
			ident = ""
		}
		w.Return('"', ident, '"')
	}
	w.Separate()
	w.Default()
	w.Return(`"`, chIdentType, `(" + strconv.Itoa(int(ch)) + ")"`)
	w.End()
	w.End()

	w.Separate()
	// func LookupChIdentFlags(ident string) ChIdent
	w.Func("LookupChIdentFlags")
	w.FuncParams("ident string")
	w.FuncResults(chIdentType)
	w.FuncBody()
	w.Switch("len(ident)")
	for _, l := range identLenList {
		slices.Sort(identLenMap[l])
		if l == 1 {
			w.Case(l)
			w.Switch("ident[0] | ('a' - 'A')")
			for _, v := range identLenMap[l] {
				w.Case("'", v, "'")
				w.Return(identMap[v])
			}
			w.End()
		} else {
			w.Case(l)
			for i, v := range identLenMap[l] {
				if i == 0 {
					w.If("ascii.EqualFold(ident, ", '"', v, '"', ')')
				} else {
					w.ElseIf("ascii.EqualFold(ident, ", '"', v, '"', ')')
				}
				w.Return(identMap[v])
			}
			w.End()
		}
	}
	w.End()
	w.Return("0")
	w.End()

	w.Separate()
	// func (s Space) ChIdentSet() ChIdentSet
	w.Method("s Space", chIdentSetType)
	w.FuncResults(chIdentSetType)
	w.FuncBody()
	w.Switch("s")
	for _, group := range flagsGroups {
		sub.Reset()
		for i, s := range group.Spaces {
			if i > 0 {
				if i%4 == 0 {
					sub.Write(",")
					sub.Indent()
				} else {
					sub.Write(", ")
				}
			}
			sub.Write(s.Name)
		}
		w.Case(sub.Bytes())
		w.Return(group.Extra)
	}
	w.Default()
	w.Return("0")
	w.End()
	w.End()

	groups.Reset()
	for i, s := range ctx.BuiltSpaces {
		sub.Reset()
		sub.Switch("ch")
		for j, ident := range s.ChannelIdent() {
			sub.Case(identMap[strings.ToLower(ident)])
			sub.Return(j)
		}
		sub.End()
		groups.Append(string(sub.Bytes()), s.ChannelCount(), i, s)
	}

	w.Separate()
	// func (s Space) ChIndex(ch ChIdent) int
	w.Method("s Space", "ChIndex")
	w.FuncParams("ch ", chIdentType)
	w.FuncResults("int")
	w.FuncBody()
	w.Switch("s")
	for _, group := range groups.Slice() {
		sub.Reset()
		for i, s := range group.Spaces {
			if i > 0 {
				if i%4 == 0 {
					sub.Write(",")
					sub.Indent()
				} else {
					sub.Write(", ")
				}
			}
			sub.Write(s.Name)
		}
		w.Case(sub.Bytes())
		w.Write(group.Key)
	}
	w.End()
	w.Separate()
	w.If("ch == ", identMap[alpha])
	w.If("n := s.ChannelCount(); n > 0")
	w.Return("n")
	w.End()
	w.End()
	w.Separate()
	w.Return("-1")
	w.End()
}

func chIdentName(ident string) string {
	return "ChIdent" + toUpperCaseFirstChar(ident)
}
