package catalog

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/oracle/ast"
)

// The bytes one character takes in the database character set (AL32UTF8)
// and in the national character set (AL16UTF16), which a CHAR-semantics
// length and an N-type length are stored in.
const (
	dbCharBytes       = 4
	nationalCharBytes = 2
	maxVarchar2Bytes  = 4000
	maxCharBytes      = 2000
	// longestVarchar2 is the longest VARCHAR2 a function returns: 32767
	// bytes under MAX_STRING_SIZE=EXTENDED, which a view's columns are cut
	// to maxVarchar2Bytes under STANDARD (clampStringTypes).
	longestVarchar2 = 32767
)

// A qtype is the data type Oracle gives a query column, as USER_TAB_COLS
// reads it back: the base type and its length, precision, scale and
// fractional-second or leading precision.
type qtype struct {
	// kind is the base type: NUMBER, FLOAT, VARCHAR2, CHAR, NVARCHAR2,
	// NCHAR, CLOB, NCLOB, BLOB, RAW, LONG, DATE, TIMESTAMP, TIMESTAMP WITH
	// TIME ZONE, TIMESTAMP WITH LOCAL TIME ZONE, INTERVAL YEAR TO MONTH,
	// INTERVAL DAY TO SECOND, BINARY_FLOAT, BINARY_DOUBLE, ROWID, UROWID,
	// BOOLEAN, XMLTYPE, JSON, NULL (an untyped NULL), or a user-defined
	// type's name.
	kind string
	// bytes is the length in bytes of a character or RAW type.
	bytes int
	// chars is the length in characters of a CHAR-semantics or N type.
	chars int
	// charSem reports CHAR length semantics.
	charSem bool
	// prec is a NUMBER's or FLOAT's precision, -1 when absent; scale is a
	// NUMBER's scale, noScale when absent (a scale may be negative).
	prec, scale int
	// frac is a TIMESTAMP's or INTERVAL DAY TO SECOND's fractional-second
	// precision; lead is an INTERVAL's leading-field precision.
	frac, lead int
}

// noScale marks a NUMBER without a scale.
const noScale = -1 << 20

func numberType() qtype { return qtype{kind: "NUMBER", prec: -1, scale: noScale} }

func nullType() qtype { return qtype{kind: "NULL", prec: -1, scale: noScale} }

func simpleType(kind string) qtype { return qtype{kind: kind, prec: -1, scale: noScale} }

func varchar2Bytes(n int) qtype {
	return qtype{kind: "VARCHAR2", bytes: min(n, longestVarchar2), chars: n, prec: -1, scale: noScale}
}

func charLiteralType(val string) qtype {
	n := len(val)
	return qtype{kind: "CHAR", bytes: n, chars: utf8.RuneCountInString(val), prec: -1, scale: noScale}
}

func timestampType(kind string, frac int) qtype {
	return qtype{kind: kind, frac: frac, prec: -1, scale: noScale}
}

func (t qtype) isNull() bool { return t.kind == "NULL" }

func (t qtype) isCharacter() bool {
	switch t.kind {
	case "VARCHAR2", "CHAR", "NVARCHAR2", "NCHAR":
		return true
	}
	return false
}

func (t qtype) isNational() bool { return t.kind == "NVARCHAR2" || t.kind == "NCHAR" }

func (t qtype) isNumeric() bool {
	switch t.kind {
	case "NUMBER", "FLOAT", "BINARY_FLOAT", "BINARY_DOUBLE":
		return true
	}
	return false
}

func (t qtype) isDatetime() bool {
	switch t.kind {
	case "DATE", "TIMESTAMP", "TIMESTAMP WITH TIME ZONE", "TIMESTAMP WITH LOCAL TIME ZONE":
		return true
	}
	return false
}

func (t qtype) isInterval() bool {
	return t.kind == "INTERVAL YEAR TO MONTH" || t.kind == "INTERVAL DAY TO SECOND"
}

func (t qtype) isLOB() bool {
	switch t.kind {
	case "CLOB", "NCLOB", "BLOB", "BFILE":
		return true
	}
	return false
}

// byteLength is how many bytes a character value of the type takes once
// converted to the database character set.
func (t qtype) byteLength() int {
	switch t.kind {
	case "VARCHAR2", "CHAR":
		return t.bytes
	case "NVARCHAR2", "NCHAR":
		return t.chars * dbCharBytes
	case "NUMBER", "FLOAT", "BINARY_FLOAT", "BINARY_DOUBLE":
		return 40
	case "RAW":
		return t.bytes * 2
	case "NULL":
		return 0
	}
	return 0
}

// String renders the type as the catalog writes a view column's type: the
// dictionary's form, as dictionaryType in the simulation renders it.
func (t qtype) String() string {
	switch t.kind {
	case "VARCHAR2", "CHAR":
		if t.charSem {
			return fmt.Sprintf("%s(%d CHAR)", t.kind, t.chars)
		}
		return fmt.Sprintf("%s(%d BYTE)", t.kind, t.bytes)
	case "NVARCHAR2", "NCHAR":
		return fmt.Sprintf("%s(%d)", t.kind, t.chars)
	case "RAW", "UROWID":
		return fmt.Sprintf("%s(%d)", t.kind, t.bytes)
	case "NULL":
		return "VARCHAR2(0 BYTE)"
	case "NUMBER":
		switch {
		case t.prec < 0 && t.scale != noScale:
			return fmt.Sprintf("NUMBER(*,%d)", t.scale)
		case t.prec < 0:
			return "NUMBER"
		case t.scale != 0 && t.scale != noScale:
			return fmt.Sprintf("NUMBER(%d,%d)", t.prec, t.scale)
		default:
			return fmt.Sprintf("NUMBER(%d)", t.prec)
		}
	case "FLOAT":
		if t.prec > 0 {
			return fmt.Sprintf("FLOAT(%d)", t.prec)
		}
		return "FLOAT"
	case "TIMESTAMP":
		return fmt.Sprintf("TIMESTAMP(%d)", t.frac)
	case "TIMESTAMP WITH TIME ZONE":
		return fmt.Sprintf("TIMESTAMP(%d) WITH TIME ZONE", t.frac)
	case "TIMESTAMP WITH LOCAL TIME ZONE":
		return fmt.Sprintf("TIMESTAMP(%d) WITH LOCAL TIME ZONE", t.frac)
	case "INTERVAL YEAR TO MONTH":
		return fmt.Sprintf("INTERVAL YEAR(%d) TO MONTH", t.lead)
	case "INTERVAL DAY TO SECOND":
		return fmt.Sprintf("INTERVAL DAY(%d) TO SECOND(%d)", t.lead, t.frac)
	}
	return t.kind
}

// sameType reports whether two types are the same type with the same
// length, precision and scale.
func sameType(a, b qtype) bool {
	return a.String() == b.String()
}

// parseColumnType reads a catalog column's type text (as normalizeDataTypeText
// writes it, or as a synced schema holds it) into a qtype.
func parseColumnType(text string) qtype {
	s := strings.ToUpper(strings.TrimSpace(text))
	if s == "" {
		return simpleType("UNKNOWN")
	}
	base, args, rest := splitTypeText(s)
	switch base {
	case "VARCHAR2", "VARCHAR", "CHAR", "CHARACTER", "NVARCHAR2", "NCHAR", "CHARACTER VARYING", "CHAR VARYING", "NATIONAL CHARACTER VARYING", "NATIONAL CHAR VARYING", "NATIONAL CHARACTER", "NATIONAL CHAR", "NCHAR VARYING":
		kind := map[string]string{
			"VARCHAR2": "VARCHAR2", "VARCHAR": "VARCHAR2", "CHARACTER VARYING": "VARCHAR2", "CHAR VARYING": "VARCHAR2",
			"CHAR": "CHAR", "CHARACTER": "CHAR",
			"NVARCHAR2": "NVARCHAR2", "NATIONAL CHARACTER VARYING": "NVARCHAR2", "NATIONAL CHAR VARYING": "NVARCHAR2", "NCHAR VARYING": "NVARCHAR2",
			"NCHAR": "NCHAR", "NATIONAL CHARACTER": "NCHAR", "NATIONAL CHAR": "NCHAR",
		}[base]
		n := 1
		sem := ""
		if len(args) > 0 {
			f := strings.Fields(args[0])
			if len(f) > 0 {
				if v, err := strconv.Atoi(f[0]); err == nil {
					n = v
				}
			}
			if len(f) > 1 {
				sem = f[1]
			}
		}
		return characterType(kind, n, sem == "CHAR")
	case "NUMBER", "NUMERIC", "DECIMAL", "DEC":
		t := numberType()
		if len(args) > 0 {
			if args[0] != "*" {
				t.prec = atoiOr(args[0], -1)
			}
			if len(args) > 1 {
				t.scale = atoiOr(args[1], noScale)
			} else if t.prec >= 0 {
				t.scale = 0
			}
		}
		if base != "NUMBER" && len(args) == 0 {
			t.scale = 0
		}
		return t
	case "INTEGER", "INT", "SMALLINT":
		t := numberType()
		t.scale = 0
		return t
	case "FLOAT":
		t := simpleType("FLOAT")
		t.prec = 126
		if len(args) > 0 {
			t.prec = atoiOr(args[0], 126)
		}
		return t
	case "REAL":
		t := simpleType("FLOAT")
		t.prec = 63
		return t
	case "DOUBLE PRECISION":
		t := simpleType("FLOAT")
		t.prec = 126
		return t
	case "RAW":
		t := simpleType("RAW")
		if len(args) > 0 {
			t.bytes = atoiOr(args[0], 0)
		}
		return t
	case "UROWID":
		t := simpleType("UROWID")
		t.bytes = 4000
		if len(args) > 0 {
			t.bytes = atoiOr(args[0], 4000)
		}
		return t
	case "TIMESTAMP":
		frac := 6
		if len(args) > 0 {
			frac = atoiOr(args[0], 6)
		}
		switch strings.Join(strings.Fields(rest), " ") {
		case "WITH TIME ZONE":
			return timestampType("TIMESTAMP WITH TIME ZONE", frac)
		case "WITH LOCAL TIME ZONE":
			return timestampType("TIMESTAMP WITH LOCAL TIME ZONE", frac)
		}
		return timestampType("TIMESTAMP", frac)
	case "INTERVAL YEAR":
		t := simpleType("INTERVAL YEAR TO MONTH")
		t.lead = 2
		if len(args) > 0 {
			t.lead = atoiOr(args[0], 2)
		}
		return t
	case "INTERVAL DAY":
		t := simpleType("INTERVAL DAY TO SECOND")
		t.lead, t.frac = 2, 6
		if len(args) > 0 {
			t.lead = atoiOr(args[0], 2)
		}
		if _, sargs, _ := splitTypeText(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(rest), "TO"))); len(sargs) > 0 {
			t.frac = atoiOr(sargs[0], 6)
		}
		return t
	case "LONG RAW", "LONGRAW":
		return simpleType("LONG RAW")
	}
	return simpleType(strings.Join(strings.Fields(s), " "))
}

// characterType is a character type of n units: bytes, or characters when
// charSem is set or the type is an N type.
func characterType(kind string, n int, charSem bool) qtype {
	t := qtype{kind: kind, prec: -1, scale: noScale}
	limit := longestVarchar2
	if kind == "CHAR" || kind == "NCHAR" {
		limit = maxCharBytes
	}
	switch {
	case kind == "NVARCHAR2" || kind == "NCHAR":
		t.chars = n
		t.bytes = min(n*nationalCharBytes, limit)
	case charSem:
		t.charSem = true
		t.chars = n
		t.bytes = min(n*dbCharBytes, limit)
	default:
		t.bytes = n
		t.chars = n
	}
	return t
}

// splitTypeText splits "BASE(args) rest" into the base name, the
// comma-separated arguments of its first parenthesis, and the text after it.
func splitTypeText(s string) (base string, args []string, rest string) {
	open := strings.IndexByte(s, '(')
	if open < 0 {
		// INTERVAL DAY TO SECOND and the like without parentheses.
		for _, prefix := range []string{"INTERVAL YEAR", "INTERVAL DAY", "TIMESTAMP"} {
			if strings.HasPrefix(s, prefix) {
				return prefix, nil, strings.TrimSpace(s[len(prefix):])
			}
		}
		return strings.Join(strings.Fields(s), " "), nil, ""
	}
	closeIdx := strings.IndexByte(s[open:], ')')
	if closeIdx < 0 {
		return strings.TrimSpace(s[:open]), nil, ""
	}
	base = strings.Join(strings.Fields(s[:open]), " ")
	for _, a := range strings.Split(s[open+1:open+closeIdx], ",") {
		args = append(args, strings.TrimSpace(a))
	}
	return base, args, strings.TrimSpace(s[open+closeIdx+1:])
}

func atoiOr(s string, def int) int {
	v, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return def
	}
	return v
}

// typeNameType reads a CAST target or a declared type into a qtype.
func typeNameType(tn *ast.TypeName) qtype {
	if tn == nil {
		return simpleType("UNKNOWN")
	}
	var names []string
	for _, item := range listItems(tn.Names) {
		if s, ok := item.(*ast.String); ok {
			names = append(names, strings.ToUpper(s.Str))
		}
	}
	var mods []string
	for _, item := range listItems(tn.TypeMods) {
		switch m := item.(type) {
		case *ast.Integer:
			mods = append(mods, strconv.FormatInt(m.Ival, 10))
		case *ast.String:
			mods = append(mods, strings.ToUpper(m.Str))
		case *ast.NumberLiteral:
			mods = append(mods, m.Val)
		}
	}
	name := strings.Join(names, " ")
	switch {
	case strings.HasPrefix(name, "TIMESTAMP"):
		text := "TIMESTAMP"
		if len(mods) > 0 {
			text += "(" + mods[0] + ")"
		}
		return parseColumnType(text + " " + strings.TrimSpace(strings.TrimPrefix(name, "TIMESTAMP")))
	case strings.HasPrefix(name, "INTERVAL YEAR"):
		if len(mods) > 0 {
			return parseColumnType("INTERVAL YEAR(" + mods[0] + ") TO MONTH")
		}
		return parseColumnType("INTERVAL YEAR TO MONTH")
	case strings.HasPrefix(name, "INTERVAL DAY"):
		lead, frac := "2", "6"
		if len(mods) > 0 {
			lead = mods[0]
		}
		if len(mods) > 1 {
			frac = mods[1]
		}
		return parseColumnType("INTERVAL DAY(" + lead + ") TO SECOND(" + frac + ")")
	}
	if len(mods) == 0 {
		switch name {
		case "VARCHAR2", "VARCHAR", "NVARCHAR2", "RAW":
			// Oracle requires a length; the parser accepted none.
			return parseColumnType(name)
		}
		return parseColumnType(name)
	}
	// "20 CHAR" arrives as two mods.
	if len(mods) == 2 && (mods[1] == "CHAR" || mods[1] == "BYTE") {
		return parseColumnType(name + "(" + mods[0] + " " + mods[1] + ")")
	}
	return parseColumnType(name + "(" + strings.Join(mods, ",") + ")")
}

// clampStringTypes cuts a view's character columns to the longest VARCHAR2
// and NVARCHAR2 the server's MAX_STRING_SIZE allows.
func clampStringTypes(cols []qcol, limit int) []qcol {
	for i := range cols {
		t := &cols[i].typ
		switch t.kind {
		case "VARCHAR2":
			if t.bytes > limit {
				t.bytes = limit
				if !t.charSem {
					t.chars = limit
				}
			}
		case "NVARCHAR2":
			if t.chars*nationalCharBytes > limit {
				t.chars = limit / nationalCharBytes
				t.bytes = limit
			}
		}
	}
	return cols
}
