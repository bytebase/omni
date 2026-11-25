package catalog

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/trino/parser"
)

// errUnknownType is a type name Trino does not know; columnOf names the
// column in the error it becomes.
type errUnknownType struct{ name string }

func (e errUnknownType) Error() string { return "unknown type " + e.name }

// simpleTypes are the parameterless types Trino 482 knows, by the name a
// statement spells them, with the name Trino stores.
var simpleTypes = map[string]string{
	"boolean": "boolean", "tinyint": "tinyint", "smallint": "smallint",
	"integer": "integer", "int": "integer", "bigint": "bigint",
	"real": "real", "double": "double", "double precision": "double",
	"varbinary": "varbinary", "json": "json", "date": "date", "uuid": "uuid",
	"ipaddress": "ipaddress", "hyperloglog": "HyperLogLog", "p4hyperloglog": "P4HyperLogLog",
	"tdigest": "tdigest", "setdigest": "SetDigest", "color": "color",
	"geometry": "Geometry", "sphericalgeography": "SphericalGeography", "bingtile": "BingTile",
	"jsonpath": "JsonPath", "number": "number",
}

// typeName is the name Trino stores for a column type, as
// information_schema.columns.data_type reads it.
func typeName(dt *parser.DataType) (string, error) {
	switch dt.Kind {
	case parser.TypeGeneric:
		name := strings.ToLower(strings.Join(strings.Fields(dt.Name), " "))
		params := dt.Params
		ints := func() ([]string, bool) {
			var out []string
			for _, p := range params {
				if !p.IsInt {
					return nil, false
				}
				out = append(out, p.IntText)
			}
			return out, true
		}
		switch name {
		case "decimal", "dec", "numeric":
			n, ok := ints()
			if !ok || len(n) > 2 {
				return "", errUnknownType{dt.String()}
			}
			switch len(n) {
			case 0:
				return "decimal(38,0)", nil
			case 1:
				n = append(n, "0")
			}
			// Verified on Trino 482: decimal(39,0), decimal(0,0) and
			// decimal(5,6) are refused so.
			p, _ := strconv.Atoi(n[0])
			sc, _ := strconv.Atoi(n[1])
			if p < 1 || p > 38 {
				return "", errorf("INVALID_FUNCTION_ARGUMENT", "DECIMAL precision must be in range [1, 38]: %s", n[0])
			}
			if sc > p {
				return "", errorf("INVALID_FUNCTION_ARGUMENT", "DECIMAL scale must be in range [0, precision (%d)]: %s", p, n[1])
			}
			return "decimal(" + n[0] + "," + n[1] + ")", nil
		case "varchar", "char":
			n, ok := ints()
			if !ok || len(n) > 1 {
				return "", errUnknownType{dt.String()}
			}
			if len(n) == 1 {
				return name + "(" + n[0] + ")", nil
			}
			if name == "char" {
				return "char(1)", nil
			}
			return "varchar", nil
		case "array":
			if len(params) == 1 && params[0].Type != nil {
				el, err := typeName(params[0].Type)
				return "array(" + el + ")", err
			}
		case "map":
			if len(params) == 2 && params[0].Type != nil && params[1].Type != nil {
				k, err := typeName(params[0].Type)
				if err != nil {
					return "", err
				}
				v, err := typeName(params[1].Type)
				return "map(" + k + ", " + v + ")", err
			}
		default:
			if s, ok := simpleTypes[name]; ok && params == nil {
				return s, nil
			}
		}
		return "", errUnknownType{dt.Name}
	case parser.TypeDateTime:
		p := "3"
		if dt.Precision != nil {
			if !dt.Precision.IsInt {
				return "", errUnknownType{dt.String()}
			}
			p = dt.Precision.IntText
		}
		s := strings.ToLower(dt.Name) + "(" + p + ")"
		if dt.WithTimeZone {
			s += " with time zone"
		}
		return s, nil
	case parser.TypeArray:
		el, err := typeName(dt.ElementType)
		return "array(" + el + ")", err
	case parser.TypeMap:
		k, err := typeName(dt.KeyType)
		if err != nil {
			return "", err
		}
		v, err := typeName(dt.ValueType)
		return "map(" + k + ", " + v + ")", err
	case parser.TypeRow:
		var fields []string
		for _, f := range dt.Fields {
			t, err := typeName(f.Type)
			if err != nil {
				return "", err
			}
			if f.Name != nil {
				t = quoteName(fold(f.Name)) + " " + t
			}
			fields = append(fields, t)
		}
		return "row(" + strings.Join(fields, ", ") + ")", nil
	case parser.TypeInterval:
		return "interval " + strings.ToLower(dt.IntervalFrom.String()) + " to " + strings.ToLower(dt.IntervalTo.String()), nil
	}
	return "", errUnknownType{dt.String()}
}

// sqlType writes a stored type name as SHOW CREATE TABLE prints it: a
// parameter list spaced after each comma, and ROW in upper case.
func sqlType(name string) string {
	var b strings.Builder
	for i := 0; i < len(name); i++ {
		switch {
		case name[i] == ',' && i+1 < len(name) && name[i+1] != ' ':
			b.WriteString(", ")
		case strings.HasPrefix(name[i:], "row(") && (i == 0 || !isNameByte(name[i-1])):
			b.WriteString("ROW(")
			i += 3
		default:
			b.WriteByte(name[i])
		}
	}
	return b.String()
}

func isNameByte(c byte) bool {
	return c == '_' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9'
}

// typeFamily groups a stored type name by what a default literal must be.
func typeFamily(name string) string {
	base, _, _ := strings.Cut(name, "(")
	switch base {
	case "tinyint", "smallint", "integer", "bigint", "real", "double", "decimal":
		return "number"
	case "varchar", "char":
		return "string"
	case "boolean":
		return "boolean"
	case "timestamp", "time", "date":
		return base
	}
	return ""
}

// defaultText checks a column's DEFAULT literal against its type, as Trino
// does when it stores one, and returns it as SHOW CREATE TABLE prints it. A
// number, boolean or NULL may default a number; a string, one that reads as
// a number; and a date or timestamp takes only NULL or a literal of its own
// type (verified on Trino 482: TIMESTAMP DEFAULT '2020-01-01 00:00:00' is
// refused).
func defaultText(e parser.Expr, colType, sql string) (string, error) {
	text, kind := literalText(e)
	family := typeFamily(colType)
	ok := false
	switch kind {
	case "null":
		ok = true
	case "number", "boolean":
		ok = family == "number" || family == "string" || kind == "boolean" && family == "boolean"
	case "string":
		value := ""
		if lit, isLit := e.(*parser.Literal); isLit {
			value = lit.Value
		}
		switch family {
		case "string":
			n, bounded := typeLength(colType)
			ok = !bounded || utf8.RuneCountInString(value) <= n
		case "number":
			_, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
			ok = err == nil
		}
	case "date":
		ok = family == "date"
	}
	if !ok {
		return "", errorf("INVALID_DEFAULT_COLUMN_VALUE", "'%s' is not a valid %s literal", exprText(e, sql), strings.ToUpper(colType))
	}
	return text, nil
}

// typeLength reads the length of varchar(n) or char(n).
func typeLength(name string) (int, bool) {
	_, rest, ok := strings.Cut(name, "(")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(strings.TrimSuffix(rest, ")"))
	return n, err == nil
}

// literalText writes a literal as Trino's formatter prints it, with its
// kind: number, boolean, string, null, date, or "" for anything else.
func literalText(e parser.Expr) (text, kind string) {
	switch n := e.(type) {
	case *parser.Literal:
		switch n.Kind {
		case parser.LiteralNull:
			return "null", "null"
		case parser.LiteralBool:
			return strings.ToLower(n.Value), "boolean"
		case parser.LiteralString:
			return quoteString(n.Value), "string"
		case parser.LiteralInteger:
			return n.Value, "number"
		case parser.LiteralDecimal:
			return "DECIMAL '" + n.Value + "'", "number"
		case parser.LiteralDouble:
			return strings.ToUpper(n.Value), "number"
		case parser.LiteralBinary:
			return strings.ToUpper(n.Value), "binary"
		}
	case *parser.UnaryExpr:
		inner, kind := literalText(n.Operand)
		if kind != "number" {
			return "", ""
		}
		if n.Op == "+" {
			return inner, kind
		}
		if strings.HasPrefix(inner, "DECIMAL '") {
			return "DECIMAL '-" + strings.TrimPrefix(inner, "DECIMAL '"), kind
		}
		return "-" + inner, kind
	case *parser.TypeConstructor:
		if strings.EqualFold(n.Name, "DATE") {
			return "DATE " + quoteString(n.Value), "date"
		}
		return strings.ToUpper(n.Name) + " " + quoteString(n.Value), "typed"
	}
	return "", ""
}
