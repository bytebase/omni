package v1review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingFullyQualified reports the objects top-level statements name
// without a schema, one NameUnqualified finding per object:
//
//   - the relation of CREATE TABLE, CREATE INDEX, CREATE TRIGGER, ALTER
//     TABLE and the other ALTERs that share its node, a RENAME that has a
//     relation, INSERT, and UPDATE, and the sequence of CREATE SEQUENCE and
//     ALTER SEQUENCE;
//   - each object of a DROP whose name is a list of name parts, whatever it
//     drops: a schema's own name is reported, and a trigger, rule, or
//     policy counts as qualified by its table;
//   - in a SELECT, each relation of the FROM clause, through joins and
//     subqueries in FROM, whose name is a table or external table in any
//     schema of schema.Metadata. A snapshot with no schemas reports every
//     such relation, and a nil snapshot none.
//
// Nothing else is judged: not DELETE, MERGE, or CREATE VIEW; not the other
// relations an INSERT or UPDATE reads; not a SELECT's relations in a WITH
// query, in a branch of a set operation, or outside the FROM clause; and
// not a statement under EXPLAIN, in a CREATE SCHEMA, or in a routine body.
//
// An object counts as qualified when its name, joined with ".", holds a
// ".", so a quoted name that holds one passes unqualified. Name is the
// relation's name as the parser stores it, folded to lower case unless
// quoted and truncated to 63 bytes; for a DROP it is the name's parts
// joined with ".". A SELECT's relation is matched against the snapshot
// case sensitively. Every finding anchors the whole statement; Bytebase's
// V1 advisor reported the statement text's last line. Ported from
// Bytebase's V1 rule NAMING_FULLY_QUALIFIED.
func NamingFullyQualified(stmts []ast.Node, schema v1review.Schema) []v1review.Finding {
	var findings []v1review.Finding
	report := func(loc ast.Loc, name string) {
		start, end := span(loc)
		findings = append(findings, v1review.Finding{Code: v1review.NameUnqualified, Name: name, Start: start, End: end})
	}
	relation := func(loc ast.Loc, rv *ast.RangeVar) {
		if rv == nil {
			return
		}
		name := rv.Relname
		if rv.Schemaname != "" {
			name = rv.Schemaname + "." + rv.Relname
		}
		if !qualified(name) {
			report(loc, name)
		}
	}
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			relation(n.Loc, n.Relation)
		case *ast.CreateSeqStmt:
			relation(n.Loc, n.Sequence)
		case *ast.CreateTrigStmt:
			relation(n.Loc, n.Relation)
		case *ast.IndexStmt:
			relation(n.Loc, n.Relation)
		case *ast.AlterTableStmt:
			relation(n.Loc, n.Relation)
		case *ast.AlterSeqStmt:
			relation(n.Loc, n.Sequence)
		case *ast.RenameStmt:
			relation(n.Loc, n.Relation)
		case *ast.InsertStmt:
			relation(n.Loc, n.Relation)
		case *ast.UpdateStmt:
			relation(n.Loc, n.Relation)
		case *ast.DropStmt:
			if n.Objects == nil {
				continue
			}
			for _, item := range n.Objects.Items {
				list, ok := item.(*ast.List)
				if !ok {
					continue
				}
				if name := strings.Join(names(list), "."); !qualified(name) {
					report(n.Loc, name)
				}
			}
		case *ast.SelectStmt:
			if schema.Metadata == nil {
				continue
			}
			tables := snapshotTables(schema)
			for _, rv := range fromRelations(n.FromClause) {
				if (tables == nil || tables[rv.Relname]) && rv.Schemaname == "" {
					report(n.Loc, rv.Relname)
				}
			}
		}
	}
	return findings
}

// qualified reports whether a name counts as naming its schema.
func qualified(name string) bool {
	return name == "" || strings.Contains(name, ".")
}

// snapshotTables returns the names of the tables and external tables in
// every schema of the snapshot, nil when the snapshot has no schemas.
func snapshotTables(schema v1review.Schema) map[string]bool {
	if schema.Metadata.Schemas == nil {
		return nil
	}
	tables := map[string]bool{}
	for _, s := range schema.Metadata.Schemas {
		for _, table := range s.Tables {
			tables[table.Name] = true
		}
		for _, table := range s.ExternalTables {
			tables[table.Name] = true
		}
	}
	return tables
}

// fromRelations returns the relations a FROM clause reads, through joins
// and the FROM clauses of its subqueries.
func fromRelations(from *ast.List) []*ast.RangeVar {
	if from == nil {
		return nil
	}
	var relations []*ast.RangeVar
	for _, item := range from.Items {
		relations = append(relations, relationsIn(item)...)
	}
	return relations
}

func relationsIn(node ast.Node) []*ast.RangeVar {
	switch n := node.(type) {
	case *ast.RangeVar:
		return []*ast.RangeVar{n}
	case *ast.JoinExpr:
		return append(relationsIn(n.Larg), relationsIn(n.Rarg)...)
	case *ast.RangeSubselect:
		if sel, ok := n.Subquery.(*ast.SelectStmt); ok {
			return fromRelations(sel.FromClause)
		}
	}
	return nil
}
