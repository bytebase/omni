package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// pgDatabaseOwner is the owner PostgreSQL gives the public schema since 15:
// whoever owns the database.
const pgDatabaseOwner = "pg_database_owner"

// StatementObjectOwnerCheck reports each object a top-level statement
// changes or creates in that is owned by a role other than currentRole, the
// role Bytebase read from the database (CURRENT_USER) or, in tenant mode,
// the snapshot's database owner. Owners come from the pre-change schema; an
// owner of pg_database_owner means the database's owner.
//
// CREATE SCHEMA is judged against the database's owner (Name "Database",
// Table the database's name). ALTER TABLE is judged against the table's
// owner (Name "Table", Table the table as written, without its schema),
// when the snapshot has the table. Every other judged statement is judged
// against the owner of the schema it names, "public" when it names none
// (Name "Schema", Table the schema), when the snapshot has the schema:
// CREATE and ALTER SEQUENCE, ALTER TYPE, CREATE TYPE ... AS, CREATE
// EXTENSION (its SCHEMA option), CREATE FUNCTION, ALTER FUNCTION, CREATE
// INDEX (the indexed table's schema), CREATE TABLE, CREATE TRIGGER, CREATE
// MATERIALIZED VIEW, every name a DROP lists, ALTER SCHEMA ... RENAME (the
// renamed schema) and ALTER INDEX ... RENAME. A qualified name's schema is
// its first part only when it has two or more parts. Detail is the owner.
// Findings anchor the statement (Start = End = 0). OwnerMismatch is the only
// code. Ported from Bytebase's V1 rule STATEMENT_OBJECT_OWNER_CHECK.
func StatementObjectOwnerCheck(stmts []ast.Node, schema v1review.Schema, currentRole string) []v1review.Finding {
	c := &ownerCheck{schema: schema, role: currentRole}
	for i, stmt := range stmts {
		c.statement = i
		c.check(unwrapRaw(stmt))
	}
	return c.findings
}

type ownerCheck struct {
	schema    v1review.Schema
	role      string
	statement int
	findings  []v1review.Finding
}

func (c *ownerCheck) check(node ast.Node) {
	switch n := node.(type) {
	case *ast.AlterSeqStmt:
		if n.Sequence != nil {
			c.schemaOwner(rangeVarSchema(n.Sequence))
		}
	case *ast.AlterTableStmt:
		if n.Relation != nil {
			c.tableOwner(rangeVarSchema(n.Relation), n.Relation.Relname)
		}
	case *ast.AlterTypeStmt:
		c.schemaOwner(qualifier(stringItems(n.TypeName)))
	case *ast.CreateExtensionStmt:
		schemaName := defaultSchema
		if n.Options != nil {
			for _, item := range n.Options.Items {
				if elem, ok := item.(*ast.DefElem); ok && elem.Defname == "schema" {
					if s, ok := elem.Arg.(*ast.String); ok {
						schemaName = s.Str
					}
				}
			}
		}
		c.schemaOwner(schemaName)
	case *ast.CreateFunctionStmt:
		c.schemaOwner(qualifier(stringItems(n.Funcname)))
	case *ast.IndexStmt:
		if n.Relation != nil {
			c.schemaOwner(rangeVarSchema(n.Relation))
		}
	case *ast.CreateSchemaStmt:
		owner := c.schema.Metadata.GetOwner()
		if owner != c.role {
			c.findings = append(c.findings, v1review.Finding{Code: v1review.OwnerMismatch, Name: "Database", Table: c.schema.Metadata.GetName(), Detail: owner, Statement: c.statement})
		}
	case *ast.CreateSeqStmt:
		if n.Sequence != nil {
			c.schemaOwner(rangeVarSchema(n.Sequence))
		}
	case *ast.CreateStmt:
		if n.Relation != nil {
			c.schemaOwner(rangeVarSchema(n.Relation))
		}
	case *ast.CreateTrigStmt:
		if n.Relation != nil {
			c.schemaOwner(rangeVarSchema(n.Relation))
		}
	case *ast.DefineStmt:
		if n.Kind == ast.OBJECT_TYPE {
			c.schemaOwner(qualifier(stringItems(n.Defnames)))
		}
	case *ast.CreateTableAsStmt:
		if n.Objtype == ast.OBJECT_MATVIEW && n.Into != nil && n.Into.Rel != nil {
			c.schemaOwner(rangeVarSchema(n.Into.Rel))
		}
	case *ast.AlterFunctionStmt:
		if n.Func != nil {
			c.schemaOwner(qualifier(stringItems(n.Func.Objname)))
		}
	case *ast.DropStmt:
		if n.Objects == nil {
			return
		}
		for _, item := range n.Objects.Items {
			list, ok := item.(*ast.List)
			if !ok {
				continue
			}
			if parts := stringItems(list); len(parts) > 0 {
				c.schemaOwner(qualifier(parts))
			}
		}
	case *ast.RenameStmt:
		switch n.RenameType {
		case ast.OBJECT_SCHEMA:
			if n.Subname != "" {
				c.schemaOwner(n.Subname)
			}
		case ast.OBJECT_INDEX:
			if n.Relation != nil {
				c.schemaOwner(rangeVarSchema(n.Relation))
			}
		}
	}
}

func (c *ownerCheck) schemaOwner(schemaName string) {
	if schemaName == "" {
		schemaName = defaultSchema
	}
	schema := c.schema.SchemaByName(schemaName)
	if schema == nil {
		return
	}
	owner := schema.GetOwner()
	if owner == pgDatabaseOwner {
		owner = c.schema.Metadata.GetOwner()
	}
	if owner != c.role {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.OwnerMismatch, Name: "Schema", Table: schemaName, Detail: owner, Statement: c.statement})
	}
}

func (c *ownerCheck) tableOwner(schemaName, tableName string) {
	if schemaName == "" {
		schemaName = defaultSchema
	}
	if c.schema.SchemaByName(schemaName) == nil {
		return
	}
	table := c.schema.Table(schemaName, tableName)
	if table == nil {
		return
	}
	owner := table.GetOwner()
	if owner == pgDatabaseOwner {
		owner = c.schema.Metadata.GetOwner()
	}
	if owner != c.role {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.OwnerMismatch, Name: "Table", Table: tableName, Detail: owner, Statement: c.statement})
	}
}

// rangeVarSchema returns a relation's schema, "public" when unqualified.
func rangeVarSchema(rv *ast.RangeVar) string {
	if rv.Schemaname == "" {
		return defaultSchema
	}
	return rv.Schemaname
}

// stringItems returns the string items of a list, in order.
func stringItems(list *ast.List) []string {
	if list == nil {
		return nil
	}
	var parts []string
	for _, item := range list.Items {
		if s, ok := item.(*ast.String); ok {
			parts = append(parts, s.Str)
		}
	}
	return parts
}

// qualifier returns a qualified name's schema: its first part when it has
// two or more, else empty.
func qualifier(parts []string) string {
	if len(parts) >= 2 {
		return parts[0]
	}
	return ""
}
