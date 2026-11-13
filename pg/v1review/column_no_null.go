package v1review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnNoNull judges the statements together for columns left nullable,
// keyed by schema (public when unqualified), table and column name as parsed.
// texts holds each statement's source text, index for index with stmts.
//
// A CREATE TABLE records each column it defines, then forgets one that says
// NOT NULL or carries a NOT NULL or PRIMARY KEY constraint, and the columns
// a PRIMARY KEY table constraint names. An ALTER TABLE walks its actions:
// ADD COLUMN records and forgets the same way, SET NOT NULL forgets the
// column, DROP NOT NULL records it, and ADD PRIMARY KEY forgets the columns
// it names and, for PRIMARY KEY USING INDEX, every expression of that index
// in schema, looked up by the table's schema, table and index name. A
// statement naming no table is skipped, and a later recording of a column
// replaces its anchor.
//
// The anchor is where V1 found the column by searching the statement text
// for its name: in a CREATE TABLE each search starts after the previous
// column's match, in an ADD COLUMN at the start, and an unfound name, like
// DROP NOT NULL, has no range, so the caller takes the statement's first
// content line. The search can match inside an earlier word, as V1's did.
// Statement says which statement holds the anchor. Each finding has code
// ColumnNullable, Name the column, Table the table and Schema its schema,
// ordered by schema, table and column, compared bytewise. Ported from
// Bytebase's V1 rule COLUMN_NO_NULL.
func ColumnNoNull(stmts []ast.Node, texts []string, schema v1review.Schema) []v1review.Finding {
	type key struct{ schema, table, column string }
	type anchor struct{ statement, start, end int }
	nullable := make(map[key]anchor)
	for i, stmt := range stmts {
		text := texts[i]
		search := func(name string, from int) (anchor, int) {
			if from < len(text) {
				if idx := strings.Index(text[from:], name); idx >= 0 {
					pos := from + idx
					return anchor{i, pos, pos + len(name)}, pos + len(name)
				}
			}
			return anchor{statement: i}, from
		}
		switch n := statementNode(stmt).(type) {
		case *ast.CreateStmt:
			if n.Relation == nil || n.Relation.Relname == "" {
				continue
			}
			s, table := schemaOrDefault(n.Relation), n.Relation.Relname
			from := 0
			var constraints []*ast.Constraint
			if n.TableElts != nil {
				for _, item := range n.TableElts.Items {
					switch elt := item.(type) {
					case *ast.ColumnDef:
						var at anchor
						at, from = search(elt.Colname, from)
						nullable[key{s, table, elt.Colname}] = at
						if columnNotNull(elt) {
							delete(nullable, key{s, table, elt.Colname})
						}
					case *ast.Constraint:
						constraints = append(constraints, elt)
					default:
					}
				}
			}
			for _, c := range constraints {
				if c.Contype == ast.CONSTR_PRIMARY {
					for _, column := range names(c.Keys) {
						delete(nullable, key{s, table, column})
					}
				}
			}
		case *ast.AlterTableStmt:
			if n.Relation == nil || n.Relation.Relname == "" {
				continue
			}
			s, table := schemaOrDefault(n.Relation), n.Relation.Relname
			for _, cmd := range alterCommands(n) {
				switch ast.AlterTableType(cmd.Subtype) {
				case ast.AT_AddColumn:
					col, ok := cmd.Def.(*ast.ColumnDef)
					if !ok || col == nil {
						continue
					}
					at, _ := search(col.Colname, 0)
					nullable[key{s, table, col.Colname}] = at
					if columnNotNull(col) {
						delete(nullable, key{s, table, col.Colname})
					}
				case ast.AT_SetNotNull:
					delete(nullable, key{s, table, cmd.Name})
				case ast.AT_DropNotNull:
					nullable[key{s, table, cmd.Name}] = anchor{statement: i}
				case ast.AT_AddConstraint:
					c, ok := cmd.Def.(*ast.Constraint)
					if !ok || c == nil || c.Contype != ast.CONSTR_PRIMARY {
						continue
					}
					for _, column := range names(c.Keys) {
						delete(nullable, key{s, table, column})
					}
					if c.Indexname != "" {
						if index := schema.TableIndex(s, table, c.Indexname, true); index != nil {
							for _, expression := range index.GetExpressions() {
								delete(nullable, key{s, table, expression})
							}
						}
					}
				default:
				}
			}
		default:
		}
	}

	keys := make([]key, 0, len(nullable))
	for k := range nullable {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b key) int {
		if c := strings.Compare(a.schema, b.schema); c != 0 {
			return c
		}
		if c := strings.Compare(a.table, b.table); c != 0 {
			return c
		}
		return strings.Compare(a.column, b.column)
	})
	findings := make([]v1review.Finding, 0, len(keys))
	for _, k := range keys {
		at := nullable[k]
		findings = append(findings, v1review.Finding{
			Code:      v1review.ColumnNullable,
			Name:      k.column,
			Table:     k.table,
			Schema:    k.schema,
			Start:     at.start,
			End:       at.end,
			Statement: at.statement,
		})
	}
	return findings
}

// columnNotNull reports whether a column says NOT NULL or carries a NOT NULL
// or PRIMARY KEY constraint.
func columnNotNull(col *ast.ColumnDef) bool {
	if col.IsNotNull {
		return true
	}
	for _, c := range columnConstraints(col) {
		if c.Contype == ast.CONSTR_NOTNULL || c.Contype == ast.CONSTR_PRIMARY {
			return true
		}
	}
	return false
}
