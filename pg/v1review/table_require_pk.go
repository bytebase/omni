package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// TableRequirePK reports the tables final, the schema after the change,
// holds without a primary key that the change is answerable for, as
// TablePrimaryKeyMissing: Name is the table and Value its schema, both as
// final names them. Tables are visited by schema name, then table name. A
// table is answerable when original, the schema before the change, lacks it
// (under its own name, or under the name a RENAME in the change gave it) and
// a CREATE TABLE in the change named it, or when original holds it without a
// primary key and the change both dropped and created it. A table original
// holds with a primary key is always answerable. A nil final reports
// nothing; a nil original holds nothing. Primary keys and tables are looked
// up as final and original say, index names case sensitive; a partition's
// name finds no primary key.
//
// The change is read only to tell what it did to each table, keyed by the
// schema and table names as written: CREATE TABLE marks it created, a
// top-level DROP TABLE dropped, and RENAME (ALTER TABLE ... RENAME TO) moves
// what is known under the old name to the new one, remembering the oldest
// name. CREATE TABLE, every top-level ALTER TABLE statement (ALTER INDEX and
// ALTER VIEW included, as the V1 advisor parsed them) and RENAME attribute
// the table to themselves; DROP TABLE does not. A table final holds in
// schema s merges what the change did under s and unqualified, taking the
// attribution of the unqualified mention only when startLines puts it on a
// later line than the qualified one. Statements nested in a function body
// are not read. A statement may be passed bare or wrapped in an
// *ast.RawStmt.
//
// Each finding's Statement is the index in stmts of the statement the table
// is attributed to, -1 when none is. startLines holds, for each statement,
// the line its text starts on as the caller numbers lines; the function only
// compares them. Ported from Bytebase's V1 rule TABLE_REQUIRE_PK.
func TableRequirePK(stmts []ast.Node, startLines []int, original, final *v1review.Schema) []v1review.Finding {
	type mention struct {
		statement, line  int
		created, dropped bool
		renamedFrom      string
	}
	mentions := map[string]*mention{}
	mentionOf := func(schema, table string) *mention {
		key := schema + "." + table
		m := mentions[key]
		if m == nil {
			m = &mention{statement: -1}
			mentions[key] = m
		}
		return m
	}
	attribute := func(m *mention, i int) {
		m.statement, m.line = i, startLines[i]
	}
	record := func(rv *ast.RangeVar, i int) *mention {
		if rv == nil || rv.Relname == "" {
			return nil
		}
		m := mentionOf(rv.Schemaname, rv.Relname)
		attribute(m, i)
		return m
	}

	for i, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			if m := record(n.Relation, i); m != nil {
				m.created = true
			}
		case *ast.AlterTableStmt:
			record(n.Relation, i)
		case *ast.DropStmt:
			if n.RemoveType != int(ast.OBJECT_TABLE) || n.Objects == nil {
				continue
			}
			for _, item := range n.Objects.Items {
				if schema, table, ok := requirePKDropName(item); ok {
					mentionOf(schema, table).dropped = true
				}
			}
		case *ast.RenameStmt:
			if n.RenameType != ast.OBJECT_TABLE || n.Relation == nil {
				continue
			}
			oldName := n.Relation.Relname
			if oldName == "" || n.Newname == "" {
				continue
			}
			oldKey := n.Relation.Schemaname + "." + oldName
			old := mentions[oldKey]
			delete(mentions, oldKey)
			m := mentionOf(n.Relation.Schemaname, n.Newname)
			attribute(m, i)
			m.renamedFrom = oldName
			if old != nil {
				m.created = m.created || old.created
				m.dropped = m.dropped || old.dropped
				if old.renamedFrom != "" {
					m.renamedFrom = old.renamedFrom
				}
			}
		default:
		}
	}

	mentionFor := func(schema, table string) *mention {
		q, u := mentions[schema+"."+table], mentions["."+table]
		switch {
		case q == nil:
			return u
		case u == nil:
			return q
		default:
		}
		m := *q
		if u.line > q.line {
			m.statement, m.line = u.statement, u.line
		}
		m.created = q.created || u.created
		m.dropped = q.dropped || u.dropped
		if m.renamedFrom == "" {
			m.renamedFrom = u.renamedFrom
		}
		return &m
	}

	if final == nil {
		return nil
	}
	var findings []v1review.Finding
	for _, schema := range final.SchemaNames() {
		for _, table := range final.TableNames(schema) {
			if final.TablePrimaryKey(schema, table, true) != nil {
				continue
			}
			m := mentionFor(schema, table)
			before := table
			exists := original != nil && original.HasTable(schema, table)
			if !exists && m != nil && m.renamedFrom != "" {
				before = m.renamedFrom
				exists = original != nil && original.HasTable(schema, before)
			}
			switch {
			case !exists:
				if m == nil || !m.created {
					continue
				}
			case original.TablePrimaryKey(schema, before, true) == nil:
				if m == nil || !m.created || !m.dropped {
					continue
				}
			default:
			}
			statement := -1
			if m != nil {
				statement = m.statement
			}
			findings = append(findings, v1review.Finding{Code: v1review.TablePrimaryKeyMissing, Name: table, Value: schema, Statement: statement})
		}
	}
	return findings
}

// requirePKDropName splits a DROP TABLE object into its schema, empty when
// unqualified, and table name; it fails for any other shape.
func requirePKDropName(item ast.Node) (string, string, bool) {
	list, ok := item.(*ast.List)
	if !ok {
		return "", "", false
	}
	var parts []string
	for _, nameItem := range list.Items {
		if s, ok := nameItem.(*ast.String); ok {
			parts = append(parts, s.Str)
		}
	}
	switch len(parts) {
	case 1:
		return "", parts[0], true
	case 2:
		return parts[0], parts[1], true
	default:
		return "", "", false
	}
}
