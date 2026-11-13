package v1review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnTypeDisallowList judges the column types top-level statements write
// against disallowed: each column a CREATE TABLE defines, each column an
// ALTER ... ADD COLUMN adds, and the new type an ALTER COLUMN ... TYPE gives,
// the ALTER forms whatever relation they name. A column without a type is not
// judged, a statement whose relation has no name is skipped, and a statement
// in a routine body is not judged.
//
// Only the last part of the type's name as the parser stores it is compared,
// without modifiers or array bounds, so INT is compared as int4, VARCHAR(20)
// as varchar and INT[] as int4. It is compared, case-insensitively, against
// each disallowed type with the equivalence of Bytebase's
// areTypesEquivalent: equal after normalizing (serial is integer, varchar is
// character varying) or in one family (int, integer, int4; bigint, int8;
// smallint, int2; real, float4; double precision, float8; boolean, bool).
// A disallowed type with a modifier, such as varchar(20), never matches, and
// an empty one matching a type with no name reports nothing for that column.
//
// Each finding reports ColumnTypeDisallowed with Name the column, Table the
// relation's name without its schema, and Value the matched disallowed type
// lower-cased, in the order the columns are written. Where several
// disallowed types match, Value is the first in disallowed's order; the V1
// advisor took whichever a Go map gave first. The finding anchors the whole
// statement and sets Search to the column name: Bytebase's V1 advisor
// reported the line of the name's first occurrence in the statement text.
// Ported from Bytebase's V1 rule COLUMN_TYPE_DISALLOW_LIST.
func ColumnTypeDisallowList(stmts []ast.Node, disallowed []string) []v1review.Finding {
	types := make([]string, 0, len(disallowed))
	for _, t := range disallowed {
		types = append(types, strings.ToLower(t))
	}
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var table string
		var loc ast.Loc
		var columns []typedColumn
		switch n := statementNode(stmt).(type) {
		case *ast.CreateStmt:
			table, loc = relname(n.Relation), n.Loc
			for _, col := range tableColumns(n) {
				if col.TypeName != nil {
					columns = append(columns, typedColumn{col.Colname, col.TypeName})
				}
			}
		case *ast.AlterTableStmt:
			table, loc = relname(n.Relation), n.Loc
			columns = alteredColumnTypes(n)
		default:
			continue
		}
		if table == "" {
			continue
		}
		for _, col := range columns {
			typ := lastTypeName(col.typ)
			for _, t := range types {
				if !typesEquivalent(typ, t) {
					continue
				}
				if t == "" {
					// V1 took an empty match for no match.
					break
				}
				start, end := span(loc)
				findings = append(findings, v1review.Finding{
					Code:   v1review.ColumnTypeDisallowed,
					Name:   col.name,
					Table:  table,
					Value:  t,
					Search: col.name,
					Start:  start,
					End:    end,
				})
				break
			}
		}
	}
	return findings
}

// typesEquivalent reports whether two type names are the same type, copied
// from Bytebase's areTypesEquivalent.
func typesEquivalent(a, b string) bool {
	a, b = normalizeType(a), normalizeType(b)
	if a == b {
		return true
	}
	for _, family := range [][]string{
		{"int", "integer", "int4"},
		{"bigint", "int8"},
		{"smallint", "int2"},
		{"real", "float4"},
		{"double precision", "float8"},
		{"boolean", "bool"},
	} {
		if slices.Contains(family, a) && slices.Contains(family, b) {
			return true
		}
	}
	if strings.Contains(a, "(") && strings.Contains(b, "(") {
		baseA, paramsA := splitTypeParams(a)
		baseB, paramsB := splitTypeParams(b)
		if baseTypesEquivalent(baseA, baseB) {
			return paramsA == paramsB
		}
	}
	return false
}

// normalizeType lower-cases a type name, collapses its spaces, drops the
// spaces in its modifiers, and spells serial types and varchar as the types
// they stand for.
func normalizeType(t string) string {
	t = strings.ToLower(strings.TrimSpace(t))
	if i := strings.Index(t, "("); i == -1 {
		t = collapseSpaces(t)
	} else {
		t = collapseSpaces(t[:i]) + strings.ReplaceAll(t[i:], " ", "")
	}
	switch t {
	case "serial":
		return "integer"
	case "bigserial":
		return "bigint"
	case "smallserial":
		return "smallint"
	default:
	}
	if t == "varchar" {
		return "character varying"
	}
	if strings.HasPrefix(t, "varchar(") {
		return "character varying" + strings.TrimPrefix(t, "varchar")
	}
	return t
}

func collapseSpaces(s string) string {
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return strings.TrimSpace(s)
}

func splitTypeParams(t string) (string, string) {
	i := strings.Index(t, "(")
	if i == -1 {
		return t, ""
	}
	return strings.TrimSpace(t[:i]), t[i:]
}

func baseTypesEquivalent(a, b string) bool {
	a, b = strings.ToLower(strings.TrimSpace(a)), strings.ToLower(strings.TrimSpace(b))
	if a == b {
		return true
	}
	for _, family := range [][]string{
		{"character varying", "varchar"},
		{"character", "char"},
	} {
		if slices.Contains(family, a) && slices.Contains(family, b) {
			return true
		}
	}
	return false
}
