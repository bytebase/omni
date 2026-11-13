package v1review

import (
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingUniqueKey judges the unique key names top-level statements
// introduce against cfg:
//
//   - CREATE UNIQUE INDEX: the index, judged with an empty name when the
//     statement gives none.
//   - CREATE TABLE: its named UNIQUE table constraints, then its columns'
//     named UNIQUE constraints.
//   - ALTER, command by command, whatever relation it alters: an ADD
//     CONSTRAINT's named UNIQUE constraint and an ADD COLUMN's named UNIQUE
//     column constraints.
//   - ALTER INDEX ... RENAME TO: the new name, when schema's "public"
//     schema holds an index of the old name that is unique and not primary,
//     whatever schema the statement names.
//   - ALTER TABLE ... RENAME CONSTRAINT: the new name, when the table holds
//     an index of the old name in schema that is unique and not primary.
//     The table is looked up in the statement's schema, "public" when it
//     names none.
//
// An unnamed constraint, CREATE FOREIGN TABLE, and a statement in a CREATE
// SCHEMA or in a routine body are not judged.
//
// cfg.Format is expanded per key with its table, without its schema, and
// its column list joined with "_": a created index's plain column elements,
// a table constraint's key columns, a column constraint's own column, and
// for a rename the snapshot index's expressions. A UNIQUE USING INDEX
// constraint takes the expressions of that index on the statement's table
// in schema's "public" schema, whatever schema the statement names, and an
// empty list when it is not found. For a renamed index the table is the
// snapshot's table or materialized view; everywhere else it is the
// statement's spelling. A format that does not compile once expanded gives
// one PatternInvalid finding and that key is judged no further.
//
// The name is judged as the parser stores it: folded to lower case unless
// quoted and truncated to 63 bytes. Its length is in bytes. Every finding
// anchors the whole statement; Bytebase's V1 advisor reported the line of
// the name's first occurrence in the statement text instead, and no line
// for a PatternInvalid finding. Ported from Bytebase's V1 rule
// NAMING_INDEX_UK.
func NamingUniqueKey(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema) []v1review.Finding {
	naming := &templateNaming{cfg: cfg}
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.IndexStmt:
			if n.Unique {
				naming.check(n.Loc, n.Idxname, relname(n.Relation), indexColumns(n))
			}
		case *ast.CreateStmt:
			if n.TableElts == nil {
				continue
			}
			table := relname(n.Relation)
			for _, item := range n.TableElts.Items {
				if c, ok := item.(*ast.Constraint); ok {
					checkUniqueConstraint(naming, schema, n.Loc, table, c)
				}
			}
			for _, item := range n.TableElts.Items {
				if col, ok := item.(*ast.ColumnDef); ok {
					checkUniqueColumn(naming, n.Loc, table, col)
				}
			}
		case *ast.AlterTableStmt:
			if n.Cmds == nil {
				continue
			}
			table := relname(n.Relation)
			for _, item := range n.Cmds.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if !ok {
					continue
				}
				switch ast.AlterTableType(cmd.Subtype) {
				case ast.AT_AddConstraint:
					if c, ok := cmd.Def.(*ast.Constraint); ok {
						checkUniqueConstraint(naming, schema, n.Loc, table, c)
					}
				case ast.AT_AddColumn:
					if col, ok := cmd.Def.(*ast.ColumnDef); ok {
						checkUniqueColumn(naming, n.Loc, table, col)
					}
				}
			}
		case *ast.RenameStmt:
			var table string
			var index *metadata.IndexMetadata
			switch n.RenameType {
			case ast.OBJECT_INDEX:
				if old := relname(n.Relation); old != "" {
					table, index = schema.SchemaIndex(defaultSchema, old, indexCaseSensitive)
				}
			case ast.OBJECT_TABCONSTRAINT:
				lookupSchema := defaultSchema
				if n.Relation != nil && n.Relation.Schemaname != "" {
					lookupSchema = n.Relation.Schemaname
				}
				if table = relname(n.Relation); table != "" {
					index = schema.TableIndex(lookupSchema, table, n.Subname, indexCaseSensitive)
				} else {
					table, index = schema.SchemaIndex(lookupSchema, n.Subname, indexCaseSensitive)
				}
			}
			if index != nil && index.GetUnique() && !index.GetPrimary() {
				naming.check(n.Loc, n.Newname, table, index.GetExpressions())
			}
		}
	}
	return naming.findings
}

// checkUniqueConstraint judges c when it is a named UNIQUE table constraint
// on table.
func checkUniqueConstraint(naming *templateNaming, schema v1review.Schema, loc ast.Loc, table string, c *ast.Constraint) {
	if c.Contype != ast.CONSTR_UNIQUE || c.Conname == "" {
		return
	}
	columns := names(c.Keys)
	if c.Indexname != "" {
		columns = schema.TableIndex(defaultSchema, table, c.Indexname, indexCaseSensitive).GetExpressions()
	}
	naming.check(loc, c.Conname, table, columns)
}

// checkUniqueColumn judges the named UNIQUE constraints col carries.
func checkUniqueColumn(naming *templateNaming, loc ast.Loc, table string, col *ast.ColumnDef) {
	for _, c := range columnConstraints(col) {
		if c.Contype == ast.CONSTR_UNIQUE && c.Conname != "" {
			naming.check(loc, c.Conname, table, []string{col.Colname})
		}
	}
}
