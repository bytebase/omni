package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// TableComment judges, after reading every statement, the comment each
// table a top-level CREATE TABLE creates is left with by top-level COMMENT
// ON TABLE statements. A table is keyed by its schema, "public" when
// unqualified, and its name as parsed, so COMMENT ON TABLE must spell the
// schema as the CREATE TABLE did; a COMMENT naming three parts is ignored.
// The last COMMENT for a key wins, wherever it stands relative to the
// CREATE. With no comment, or an empty or NULL one, the table is
// TableCommentMissing if cfg.Required; otherwise a comment longer than a
// positive cfg.MaxLength, in bytes, is TableCommentTooLong. Name is the
// table as V1 printed it, prefixed with its schema unless that is "public".
//
// A missing comment anchors the table's last CREATE TABLE statement, a
// too-long one its last COMMENT statement; Statement is that statement's
// index in stmts. One finding per key, in the order the keys were first
// created; V1 reported them in Go map order. CREATE TABLE AS, a CREATE in a
// CREATE SCHEMA, and statements under EXPLAIN or in a routine body are not
// judged. A statement may be passed bare or wrapped in an *ast.RawStmt.
// Ported from Bytebase's V1 rule TABLE_COMMENT.
//
// A nil cfg panics at the first created table, as the V1 advisor did when
// the rule had no payload.
func TableComment(stmts []ast.Node, cfg *v1review.TableCommentConfig) []v1review.Finding {
	type anchor struct {
		statement int
		loc       ast.Loc
	}
	type created struct {
		name string
		anchor
	}
	type comment struct {
		text string
		anchor
	}
	var keys []string
	tables := make(map[string]created)
	comments := make(map[string]comment)
	for i, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			if n.Relation == nil {
				continue
			}
			schema, name := schemaOrDefault(n.Relation), n.Relation.Relname
			key := schema + "." + name
			if _, ok := tables[key]; !ok {
				keys = append(keys, key)
			}
			if schema != defaultSchema {
				name = key
			}
			tables[key] = created{name: name, anchor: anchor{i, n.Loc}}
		case *ast.CommentStmt:
			if n.Objtype != ast.OBJECT_TABLE {
				continue
			}
			key, ok := commentedTable(n.Object)
			if !ok {
				continue
			}
			comments[key] = comment{text: n.Comment, anchor: anchor{i, n.Loc}}
		}
	}
	var findings []v1review.Finding
	finding := func(code v1review.Code, name string, a anchor) {
		start, end := span(a.loc)
		findings = append(findings, v1review.Finding{Code: code, Name: name, Start: start, End: end, Statement: a.statement})
	}
	for _, key := range keys {
		table := tables[key]
		c, ok := comments[key]
		if !ok || c.text == "" {
			if cfg.Required {
				finding(v1review.TableCommentMissing, table.name, table.anchor)
			}
		} else if cfg.MaxLength > 0 && len(c.text) > cfg.MaxLength {
			finding(v1review.TableCommentTooLong, table.name, c.anchor)
		}
	}
	return findings
}

// commentedTable returns the schema.table key a COMMENT ON TABLE names,
// false when it names three parts or none.
func commentedTable(object ast.Node) (string, bool) {
	switch obj := object.(type) {
	case *ast.List:
		switch len(obj.Items) {
		case 1:
			return defaultSchema + "." + stringValue(obj.Items[0]), true
		case 2:
			return stringValue(obj.Items[0]) + "." + stringValue(obj.Items[1]), true
		}
	case *ast.String:
		return defaultSchema + "." + obj.Str, true
	}
	return "", false
}

// stringValue returns an *ast.String's text, "" for any other node.
func stringValue(n ast.Node) string {
	if s, ok := n.(*ast.String); ok {
		return s.Str
	}
	return ""
}
