package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
)

// backupDatabase is the database where Bytebase writes a MySQL prior
// backup and stages gh-ost's tables.
const backupDatabase = "bbdataarchive"

// schemaIndex looks tables up in a target's synced schema, the database
// the change runs in, as MySQL resolves their names.
type schemaIndex struct {
	schema *metadata.DatabaseSchemaMetadata
	// caseSensitive is lower_case_table_names = 0: table and database
	// names compare with case; otherwise they compare without.
	caseSensitive bool
}

func newSchemaIndex(t schemaTarget) schemaIndex {
	return schemaIndex{schema: t.schema, caseSensitive: t.lowerCaseTableNames == 0}
}

func (x schemaIndex) same(a, b string) bool {
	if x.caseSensitive {
		return a == b
	}
	return strings.EqualFold(a, b)
}

// table returns the synced table a reference means, nil when the schema
// has none of the name, or the reference may mean a table of another
// database: one it names, or, when moved is set, the current database a
// USE earlier in the change may have changed.
func (x schemaIndex) table(ref *ast.TableRef, moved bool) *metadata.TableMetadata {
	if x.schema == nil || ref == nil {
		return nil
	}
	if ref.Schema == "" && moved || ref.Schema != "" && !x.same(ref.Schema, x.schema.GetName()) {
		return nil
	}
	for _, s := range x.schema.GetSchemas() {
		for _, t := range s.GetTables() {
			if x.same(t.GetName(), ref.Name) {
				return t
			}
		}
	}
	return nil
}

// hasRelation reports whether the synced schema holds a table or view of
// the name, which a query of the database may read.
func (x schemaIndex) hasRelation(name string) bool {
	if x.schema == nil {
		return false
	}
	for _, s := range x.schema.GetSchemas() {
		for _, t := range s.GetTables() {
			if x.same(t.GetName(), name) {
				return true
			}
		}
		for _, v := range s.GetViews() {
			if x.same(v.GetName(), name) {
				return true
			}
		}
	}
	return false
}

// referenced reports whether a foreign key of a table of the synced
// schema references the table of the name in the same database.
func (x schemaIndex) referenced(name string) bool {
	if x.schema == nil {
		return false
	}
	for _, s := range x.schema.GetSchemas() {
		for _, t := range s.GetTables() {
			for _, fk := range t.GetForeignKeys() {
				db := fk.GetReferencedSchema()
				if (db == "" || x.same(db, x.schema.GetName())) && x.same(fk.GetReferencedTable(), name) {
					return true
				}
			}
		}
	}
	return false
}

// isView reports whether the synced schema holds a view of the name, and
// no table of it, in the target's database.
func (x schemaIndex) isView(name string) bool {
	if x.schema == nil {
		return false
	}
	view := false
	for _, s := range x.schema.GetSchemas() {
		for _, t := range s.GetTables() {
			if x.same(t.GetName(), name) {
				return false
			}
		}
		for _, v := range s.GetViews() {
			view = view || x.same(v.GetName(), name)
		}
	}
	return view
}

// uniqueIndex reports whether the synced table has a unique or primary
// index of the name; index names compare without case.
func uniqueIndex(t *metadata.TableMetadata, name string) (unique, found bool) {
	for _, idx := range t.GetIndexes() {
		if strings.EqualFold(idx.GetName(), name) {
			return idx.GetUnique() || idx.GetPrimary(), true
		}
	}
	return false, false
}

// beforeUse returns the position of the change's first USE, the number of
// statements when it has none: what runs from there may mean another
// database by a name without one.
func beforeUse(stmts []statement) int {
	for i := range stmts {
		if _, ok := stmts[i].node.(*ast.UseStmt); ok {
			return i
		}
	}
	return len(stmts)
}

// usesDatabase reports whether a statement before the given one is USE,
// which moves the database unqualified names mean.
func usesDatabase(stmts []statement, before int) bool {
	for i := 0; i < before && i < len(stmts); i++ {
		if _, ok := stmts[i].node.(*ast.UseStmt); ok {
			return true
		}
	}
	return false
}
