package old

import (
	"regexp/syntax"
	"strings"
)

func Dump(re *syntax.Regexp) string {
	buf := new(strings.Builder)
	dump(buf, "", re)
	return buf.String()
}

func dump(buf *strings.Builder, prefix string, re *syntax.Regexp) {

	const inc = "   " // increment for prefix

	if re == nil {
		buf.WriteString(prefix)
		buf.WriteString("<nil>")
		return
	}
	switch re.Op {
	case syntax.OpNoMatch:
		buf.WriteString(prefix)
		buf.WriteString("<NoMatch>")
	case syntax.OpEmptyMatch:
		buf.WriteString(prefix)
		buf.WriteString("<EmptyMatch>")
	case syntax.OpLiteral:
		buf.WriteString(prefix)
		buf.WriteString("<Literal>")
		buf.WriteString(string(re.Rune))
	case syntax.OpCharClass:
		buf.WriteString(prefix)
		buf.WriteString("<CharClass>")
		buf.WriteString(re.String())
	case syntax.OpAlternate:
		buf.WriteString(prefix)
		buf.WriteString("<Alternate>(\n")
		for _, sub := range re.Sub {
			dump(buf, prefix+inc, sub)
			buf.WriteString("\n")
		}
		buf.WriteString(prefix + ")")
	case syntax.OpAnyChar, syntax.OpAnyCharNotNL:
		buf.WriteString(prefix)
		buf.WriteString("<AnyChar not implemented>")
	case syntax.OpBeginLine, syntax.OpBeginText, syntax.OpEndLine, syntax.OpEndText:
		buf.WriteString(prefix)
		buf.WriteString("<start/end not implemented>")
	case syntax.OpCapture:
		buf.WriteString(prefix)
		buf.WriteString("(\n")
		dump(buf, prefix+inc, re.Sub[0])
		buf.WriteString("\n" + prefix + ")")
	case syntax.OpConcat:
		buf.WriteString(prefix)
		buf.WriteString("<concat>(\n")
		for _, sub := range re.Sub {
			dump(buf, prefix+inc, sub)
			buf.WriteString("\n")
		}
		buf.WriteString(prefix + ")")
	case syntax.OpQuest:
		buf.WriteString(prefix)
		buf.WriteString("<ZeroOrOne>(\n")
		dump(buf, prefix+inc, re.Sub0[0])
		buf.WriteString("\n" + prefix + ")")
	case syntax.OpPlus:
		buf.WriteString(prefix)
		buf.WriteString("<Plus>(\n")
		dump(buf, prefix+inc, re.Sub0[0])
		buf.WriteString("\n" + prefix + ")")
	case syntax.OpStar:
		buf.WriteString(prefix)
		buf.WriteString("<Star>(\n")
		dump(buf, prefix+inc, re.Sub0[0])
		buf.WriteString("\n" + prefix + ")")

	default:
		buf.WriteString(prefix)
		buf.WriteString("<not implemented : " + re.String() + ">")
	}

}
