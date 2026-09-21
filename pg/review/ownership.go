package review

import (
	"fmt"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
)

// pgDatabaseOwner is the role PostgreSQL 14+ gives the public schema: a
// stand-in for whoever owns the database.
const pgDatabaseOwner = "pg_database_owner"

// owners is what the synced schema says about who owns what. The catalog
// does not track ownership, so the walk reads it from here.
type owners struct {
	database string
	schemas  map[string]string
	tables   map[[2]string]string
}

func newOwners(meta *metadata.DatabaseSchemaMetadata) *owners {
	o := &owners{
		database: meta.GetOwner(),
		schemas:  make(map[string]string),
		tables:   make(map[[2]string]string),
	}
	for _, s := range meta.GetSchemas() {
		o.schemas[s.GetName()] = s.GetOwner()
		for _, t := range s.GetTables() {
			o.tables[[2]string{s.GetName(), t.GetName()}] = t.GetOwner()
		}
	}
	return o
}

// schema returns the schema's owner, or "" when the snapshot does not say.
func (o *owners) schema(name string) string {
	return o.resolve(o.schemas[name])
}

// table returns the table's owner, or "" when the snapshot does not say:
// the table is not in it, or is a view, whose owner it does not record.
func (o *owners) table(schema, name string) string {
	return o.resolve(o.tables[[2]string{schema, name}])
}

func (o *owners) resolve(owner string) string {
	if owner == pgDatabaseOwner {
		return o.database
	}
	return owner
}

// ownershipProblems lists the objects the statement acts on that the
// current role does not own, as PostgreSQL would refuse them. The check
// mirrors the server's rule for each statement kind: altering, dropping,
// renaming, or indexing a relation needs its owner; creating in a schema
// is approximated by owning the schema; creating a schema by owning the
// database. An object the snapshot has no owner for, including one the
// change itself created, is not checked.
func (s *session) ownershipProblems(n ast.Node) []string {
	c := &ownerCheck{s: s, role: s.cat.CurrentUser()}
	switch v := n.(type) {
	case *ast.CreateSchemaStmt:
		c.database()
	case *ast.AlterTableStmt:
		c.relation(v.Relation)
	case *ast.IndexStmt:
		c.relation(v.Relation)
	case *ast.CreateTrigStmt:
		c.relation(v.Relation)
	case *ast.CreatePolicyStmt:
		c.relation(v.Table)
	case *ast.RefreshMatViewStmt:
		c.relation(v.Relation)
	case *ast.AlterObjectSchemaStmt:
		if v.Relation != nil {
			c.relation(v.Relation)
		}
	case *ast.RenameStmt:
		switch v.RenameType {
		case ast.OBJECT_SCHEMA:
			c.schema(v.Subname)
		default:
			if v.Relation != nil {
				c.relation(v.Relation)
			}
		}
	case *ast.DropStmt:
		c.drop(v)
	case *ast.CreateStmt:
		c.schemaOf(v.Relation)
	case *ast.CreateForeignTableStmt:
		c.schemaOf(v.Base.Relation)
	case *ast.CreateSeqStmt:
		c.schemaOf(v.Sequence)
	case *ast.AlterSeqStmt:
		c.schemaOf(v.Sequence)
	case *ast.ViewStmt:
		c.schemaOf(v.View)
	case *ast.CreateTableAsStmt:
		if v.Into != nil {
			c.schemaOf(v.Into.Rel)
		}
	case *ast.CompositeTypeStmt:
		c.schemaOf(v.Typevar)
	case *ast.CreateEnumStmt:
		c.schemaOfName(nameParts(v.TypeName))
	case *ast.AlterEnumStmt:
		c.schemaOfName(nameParts(v.Typname))
	case *ast.CreateDomainStmt:
		c.schemaOfName(nameParts(v.Domainname))
	case *ast.AlterDomainStmt:
		c.schemaOfName(nameParts(v.Typname))
	case *ast.AlterTypeStmt:
		c.schemaOfName(nameParts(v.TypeName))
	case *ast.DefineStmt:
		c.schemaOfName(nameParts(v.Defnames))
	case *ast.CreateFunctionStmt:
		c.schemaOfName(nameParts(v.Funcname))
	case *ast.AlterFunctionStmt:
		if v.Func != nil {
			c.schemaOfName(nameParts(v.Func.Objname))
		}
	}
	return c.problems
}

type ownerCheck struct {
	s        *session
	role     string
	problems []string
}

func (c *ownerCheck) fail(kind, name, owner string) {
	c.problems = append(c.problems, fmt.Sprintf("%s %s is owned by %s, but the change runs as %s", kind, name, ident(owner), ident(c.role)))
}

func (c *ownerCheck) database() {
	if owner := c.s.owners.database; owner != "" && owner != c.role {
		c.fail("database", ident(c.s.target.Schema.GetName()), owner)
	}
}

func (c *ownerCheck) schema(name string) {
	if name == "" {
		return
	}
	if owner := c.s.owners.schema(name); owner != "" && owner != c.role {
		c.fail("schema", ident(name), owner)
	}
}

// relation checks the owner of an existing relation. A name the catalog
// cannot resolve is left to the walk, which reports it.
func (c *ownerCheck) relation(rv *ast.RangeVar) {
	if rv == nil {
		return
	}
	rel := c.s.cat.GetRelation(rv.Schemaname, rv.Relname)
	if rel == nil || rel.Schema == nil {
		return
	}
	if owner := c.s.owners.table(rel.Schema.Name, rel.Name); owner != "" && owner != c.role {
		c.fail("table", qualified([]string{rel.Schema.Name, rel.Name}), owner)
	}
}

// schemaOf checks the owner of the schema a new relation goes in.
func (c *ownerCheck) schemaOf(rv *ast.RangeVar) {
	if rv == nil {
		return
	}
	c.schemaNamed(rv.Schemaname)
}

// schemaOfName is schemaOf for a qualified name list.
func (c *ownerCheck) schemaOfName(parts []string) {
	if len(parts) == 0 {
		return
	}
	schema := ""
	if len(parts) > 1 {
		schema = parts[len(parts)-2]
	}
	c.schemaNamed(schema)
}

// schemaNamed checks the named schema, or the current schema when the
// name is empty.
func (c *ownerCheck) schemaNamed(name string) {
	if name == "" {
		cur := c.s.cat.CurrentSchema()
		if cur == nil {
			return
		}
		name = cur.Name
	}
	c.schema(name)
}

// drop checks each dropped object: relations by their owner, schemas by
// theirs, relation-bound objects by the relation's, and the rest by the
// schema they are in.
func (c *ownerCheck) drop(v *ast.DropStmt) {
	if v.Objects == nil {
		return
	}
	kind := ast.ObjectType(v.RemoveType)
	for _, obj := range v.Objects.Items {
		switch kind {
		case ast.OBJECT_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_FOREIGN_TABLE:
			c.relation(rangeVarOf(nameParts(listOf(obj))))
		case ast.OBJECT_SEQUENCE, ast.OBJECT_INDEX:
			c.schemaOfName(nameParts(listOf(obj)))
		case ast.OBJECT_TRIGGER, ast.OBJECT_POLICY, ast.OBJECT_RULE:
			parts := nameParts(listOf(obj))
			if len(parts) > 1 {
				c.relation(rangeVarOf(parts[:len(parts)-1]))
			}
		case ast.OBJECT_SCHEMA:
			if s, ok := obj.(*ast.String); ok {
				c.schema(s.Str)
			} else {
				c.schemaOfName(append(nameParts(listOf(obj)), ""))
			}
		default:
			switch o := obj.(type) {
			case *ast.List:
				c.schemaOfName(nameParts(o))
			case *ast.TypeName:
				c.schemaOfName(nameParts(o.Names))
			case *ast.ObjectWithArgs:
				c.schemaOfName(nameParts(o.Objname))
			}
		}
	}
}

func listOf(n ast.Node) *ast.List {
	l, _ := n.(*ast.List)
	return l
}

// rangeVarOf builds the RangeVar of a one- or two-part name.
func rangeVarOf(parts []string) *ast.RangeVar {
	switch len(parts) {
	case 1:
		return &ast.RangeVar{Relname: parts[0]}
	case 2:
		return &ast.RangeVar{Schemaname: parts[0], Relname: parts[1]}
	default:
		return nil
	}
}
