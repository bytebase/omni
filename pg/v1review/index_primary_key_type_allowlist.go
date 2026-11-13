package v1review

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexPrimaryKeyTypeAllowlist reports the primary key columns top-level
// CREATE TABLE statements define with a type that is not in allow: each
// column with a PRIMARY KEY constraint, once per such constraint, and then
// each column a table-level PRIMARY KEY lists that the statement defines.
// ALTER TABLE, CREATE TABLE AS, CREATE FOREIGN TABLE, and a statement in a
// CREATE SCHEMA or in a routine body are not judged, nor is a listed column
// the statement does not define.
//
// A column's type is the last part of its type's name as the parser stores
// it, so INT is int4 and pg_catalog.varchar is varchar, followed by its
// integer modifiers in parentheses and [] for an array, then normalized by
// normalizePostgreSQLType: VARCHAR(20) is `character varying(20)` and
// SERIAL is `integer`. A column defined without a type has the empty type.
// Where a statement defines a column twice, a listed column takes the type
// of its last definition. A type is allowed when isTypeInList finds it
// equivalent to an entry of allow.
//
// Name is the column as the parser stores it, Table the table's name
// without its schema, and Value the type as judged. Every finding anchors the
// whole statement; Bytebase's V1 advisor reported the line of the column
// name's first occurrence in the statement text, matched case sensitively
// and so not necessarily the column's own, which is Search. Findings are in
// statement order, within one the column constraints first. Ported from
// Bytebase's V1 rule INDEX_PRIMARY_KEY_TYPE_ALLOWLIST.
func IndexPrimaryKeyTypeAllowlist(stmts []ast.Node, allow []string) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		n, ok := stmt.(*ast.CreateStmt)
		if !ok || n.TableElts == nil {
			continue
		}
		table := relname(n.Relation)
		start, end := span(n.Loc)
		report := func(column, columnType string) {
			findings = append(findings, v1review.Finding{Code: v1review.PrimaryKeyTypeNotAllowed, Name: column, Table: table, Value: columnType, Search: column, Start: start, End: end})
		}
		columnTypes := make(map[string]string)
		for _, item := range n.TableElts.Items {
			col, ok := item.(*ast.ColumnDef)
			if !ok {
				continue
			}
			columnType := normalizePostgreSQLType(typeNameFull(col.TypeName))
			columnTypes[col.Colname] = columnType
			for _, c := range columnConstraints(col) {
				if c.Contype == ast.CONSTR_PRIMARY && !isTypeInList(columnType, allow) {
					report(col.Colname, columnType)
				}
			}
		}
		for _, item := range n.TableElts.Items {
			c, ok := item.(*ast.Constraint)
			if !ok || c.Contype != ast.CONSTR_PRIMARY {
				continue
			}
			for _, column := range names(c.Keys) {
				if columnType, ok := columnTypes[column]; ok && !isTypeInList(columnType, allow) {
					report(column, columnType)
				}
			}
		}
	}
	return findings
}

// typeNameFull returns a type's name as the V1 advisors spelled it: the last
// part of its qualified name, its integer modifiers in parentheses, and []
// for an array, so character varying(20) and int4[]. A type without a name
// is the empty string.
func typeNameFull(tn *ast.TypeName) string {
	if tn == nil || tn.Names == nil {
		return ""
	}
	var name string
	for _, item := range tn.Names.Items {
		if s, ok := item.(*ast.String); ok {
			name = s.Str
		}
	}
	if name == "" {
		return ""
	}
	if tn.Typmods != nil && len(tn.Typmods.Items) > 0 {
		var mods []string
		for _, item := range tn.Typmods.Items {
			switch v := item.(type) {
			case *ast.Integer:
				mods = append(mods, fmt.Sprintf("%d", v.Ival))
			case *ast.A_Const:
				if iv, ok := v.Val.(*ast.Integer); ok {
					mods = append(mods, fmt.Sprintf("%d", iv.Ival))
				}
			}
		}
		if len(mods) > 0 {
			name += "(" + strings.Join(mods, ",") + ")"
		}
	}
	if tn.ArrayBounds != nil && len(tn.ArrayBounds.Items) > 0 {
		name += "[]"
	}
	return name
}
