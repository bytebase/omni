package v1review

import (
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingPrimaryKey judges the primary key names top-level statements
// introduce against cfg:
//
//   - CREATE TABLE: its named PRIMARY KEY table constraints. One written
//     USING INDEX is judged a second time, after all of them, with that
//     index's columns.
//   - ALTER, command by command, whatever relation it alters: an ADD
//     CONSTRAINT's named PRIMARY KEY constraint.
//   - ALTER TABLE ... RENAME CONSTRAINT: the new name, when the table holds
//     a primary index of the old name in schema.
//   - ALTER INDEX ... RENAME TO: the new name, when schema holds a primary
//     index of the old name.
//
// A column's PRIMARY KEY constraint, in a CREATE TABLE or an ADD COLUMN, is
// not judged, nor are an unnamed constraint, CREATE FOREIGN TABLE, and a
// statement in a CREATE SCHEMA or in a routine body. Every lookup is in the
// statement's schema, "public" when it names none.
//
// cfg.Format is expanded per key with its table, without its schema, and
// its column list joined with "_": a constraint's key columns, and for a
// rename the snapshot index's expressions. A PRIMARY KEY USING INDEX
// constraint takes the expressions of that index on the statement's table
// in schema, and an empty list when it is not found. For a renamed index
// the table is the snapshot's; everywhere else it is the statement's
// spelling. A format that does not compile once expanded gives one
// PatternInvalid finding and that key is judged no further.
//
// The name is judged as the parser stores it: folded to lower case unless
// quoted and truncated to 63 bytes. Its length is in bytes. Every finding
// anchors the whole statement; Bytebase's V1 advisor reported the line of
// the name's first occurrence in the statement text instead. Ported from
// Bytebase's V1 rule NAMING_INDEX_PK.
func NamingPrimaryKey(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema) []v1review.Finding {
	naming := &templateNaming{cfg: cfg}
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			if n.TableElts == nil {
				continue
			}
			table := relname(n.Relation)
			for _, item := range n.TableElts.Items {
				if c, ok := item.(*ast.Constraint); ok && c.Contype == ast.CONSTR_PRIMARY && c.Conname != "" {
					naming.check(n.Loc, c.Conname, table, names(c.Keys))
				}
			}
			for _, item := range n.TableElts.Items {
				if c, ok := item.(*ast.Constraint); ok && c.Contype == ast.CONSTR_PRIMARY && c.Conname != "" && c.Indexname != "" {
					index := schema.TableIndex(schemaOrDefault(n.Relation), table, c.Indexname, indexCaseSensitive)
					naming.check(n.Loc, c.Conname, table, index.GetExpressions())
				}
			}
		case *ast.AlterTableStmt:
			if n.Cmds == nil {
				continue
			}
			table := relname(n.Relation)
			for _, item := range n.Cmds.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if !ok || ast.AlterTableType(cmd.Subtype) != ast.AT_AddConstraint {
					continue
				}
				c, ok := cmd.Def.(*ast.Constraint)
				if !ok || c.Contype != ast.CONSTR_PRIMARY || c.Conname == "" {
					continue
				}
				columns := names(c.Keys)
				if c.Indexname != "" {
					columns = schema.TableIndex(schemaOrDefault(n.Relation), table, c.Indexname, indexCaseSensitive).GetExpressions()
				}
				naming.check(n.Loc, c.Conname, table, columns)
			}
		case *ast.RenameStmt:
			var table string
			var index *metadata.IndexMetadata
			switch n.RenameType {
			case ast.OBJECT_TABCONSTRAINT:
				if n.Subname != "" {
					table = relname(n.Relation)
					index = schema.TableIndex(schemaOrDefault(n.Relation), table, n.Subname, indexCaseSensitive)
				}
			case ast.OBJECT_INDEX:
				if old := relname(n.Relation); old != "" {
					table, index = schema.SchemaIndex(schemaOrDefault(n.Relation), old, indexCaseSensitive)
				}
			}
			if index != nil && index.GetPrimary() {
				naming.check(n.Loc, n.Newname, table, index.GetExpressions())
			}
		}
	}
	return naming.findings
}
