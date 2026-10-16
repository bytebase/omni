package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
)

// backupDatabase is the database where Bytebase writes a SQL Server prior
// backup.
const backupDatabase = "bbdataarchive"

// schemaIndex looks objects up in a target's synced schema, the database
// the change runs in. SQL Server compares names under the database's
// collation, usually without case; the index compares without case.
type schemaIndex struct {
	schema *metadata.DatabaseSchemaMetadata
}

// inDatabase reports whether a reference may mean an object of the
// target's database: it names no server, and no other database.
func (x schemaIndex) inDatabase(ref *ast.TableRef) bool {
	return ref != nil && ref.Server == "" && (ref.Database == "" || sameName(ref.Database, x.schema.GetName()))
}

// table returns the synced table a reference means: in the schema it
// names, or, unqualified, the one table of the name in any schema when
// that schema is dbo. It returns nil when the reference may mean another
// database, or no table, or several, or an unqualified name only another
// schema holds: SQL Server resolves an unqualified name in the session's
// default schema, which the review does not know, and then in dbo.
func (x schemaIndex) table(ref *ast.TableRef) *metadata.TableMetadata {
	if x.schema == nil || !x.inDatabase(ref) {
		return nil
	}
	var found *metadata.TableMetadata
	for _, s := range x.schema.GetSchemas() {
		if ref.Schema != "" && !sameName(s.GetName(), ref.Schema) {
			continue
		}
		for _, t := range s.GetTables() {
			if sameName(t.GetName(), ref.Object) {
				if found != nil || !resolves(ref, s) {
					return nil
				}
				found = t
			}
		}
	}
	return found
}

// resolves reports whether a reference to an object of the synced schema
// is known to mean it: the reference names the schema, or the schema is
// dbo, where an unqualified name resolves when the default schema holds
// none of the name.
func resolves(ref *ast.TableRef, s *metadata.SchemaMetadata) bool {
	return ref.Schema != "" || sameName(s.GetName(), "dbo")
}

// relationKind returns "table" or "view" for the synced relation a
// reference means, as table resolves it, "" for none or more than one.
func (x schemaIndex) relationKind(ref *ast.TableRef) string {
	if x.schema == nil || !x.inDatabase(ref) {
		return ""
	}
	kind := ""
	for _, s := range x.schema.GetSchemas() {
		if ref.Schema != "" && !sameName(s.GetName(), ref.Schema) {
			continue
		}
		for _, t := range s.GetTables() {
			if sameName(t.GetName(), ref.Object) {
				if kind != "" || !resolves(ref, s) {
					return ""
				}
				kind = "table"
			}
		}
		for _, v := range s.GetViews() {
			if sameName(v.GetName(), ref.Object) {
				if kind != "" || !resolves(ref, s) {
					return ""
				}
				kind = "view"
			}
		}
	}
	return kind
}

// holdsRelation reports whether a synced schema the reference may mean
// holds a table or view of its name: the schema named, or any.
func (x schemaIndex) holdsRelation(ref *ast.TableRef) bool {
	if x.schema == nil {
		return false
	}
	for _, s := range x.schema.GetSchemas() {
		if ref.Schema != "" && !sameName(s.GetName(), ref.Schema) {
			continue
		}
		for _, t := range s.GetTables() {
			if sameName(t.GetName(), ref.Object) {
				return true
			}
		}
		for _, v := range s.GetViews() {
			if sameName(v.GetName(), ref.Object) {
				return true
			}
		}
	}
	return false
}

// constraintKind returns the kind of the synced table's constraint of the
// name: the primary key or a unique constraint, by its index, a foreign
// key, or a check; ok is false when the table has none of the name, which
// may be a default constraint. A unique index that is not a constraint is
// none: DROP CONSTRAINT refuses it.
func constraintKind(t *metadata.TableMetadata, name string) (ast.ConstraintType, bool) {
	for _, idx := range t.GetIndexes() {
		if strings.EqualFold(idx.GetName(), name) {
			switch {
			case idx.GetPrimary():
				return ast.ConstraintPrimaryKey, true
			case idx.GetUnique() && idx.GetIsConstraint():
				return ast.ConstraintUnique, true
			}
		}
	}
	for _, fk := range t.GetForeignKeys() {
		if strings.EqualFold(fk.GetName(), name) {
			return ast.ConstraintForeignKey, true
		}
	}
	for _, ck := range t.GetCheckConstraints() {
		if strings.EqualFold(ck.GetName(), name) {
			return ast.ConstraintCheck, true
		}
	}
	return 0, false
}

// beforeUse returns the position of the change's first USE, the number of
// statements when it has none: what runs from there may mean another
// database by a name without one.
func beforeUse(stmts []statement) int {
	for i := range stmts {
		if hasUse(stmts[i].node) {
			return i
		}
	}
	return len(stmts)
}

// usesDatabase reports whether a statement before the given position is
// USE, which moves the database names mean.
func usesDatabase(stmts []statement, before int) bool {
	for i := 0; i < before && i < len(stmts); i++ {
		if hasUse(stmts[i].node) {
			return true
		}
	}
	return false
}

// hasUse reports whether a statement is or holds a USE, in a block or
// anywhere it may run.
func hasUse(n ast.Node) bool {
	found := false
	ast.Inspect(n, func(m ast.Node) bool {
		if _, ok := m.(*ast.UseStmt); ok {
			found = true
		}
		return !found
	})
	return found
}
