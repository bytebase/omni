package review

import (
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/snowflake/ast"
)

// schemaIndex looks objects up in a target's synced schema, the database
// the change runs in. Bytebase syncs a Snowflake database's names from its
// INFORMATION_SCHEMA, as Snowflake stores them, so a name of the SQL is
// compared by its normalized form, exactly. The sync reads tables, columns,
// views, streams, and tasks, and no keys: until it fills indexes and
// foreign keys, a DROP CONSTRAINT is classified only from what the change
// itself defines.
//
// bytebase: backend/plugin/db/snowflake/sync.go — SyncDBSchema, getTableSchema
type schemaIndex struct {
	schema *metadata.DatabaseSchemaMetadata
}

// table returns the synced table a name means: in the schema it names, or,
// unqualified, the one table of the name in any schema. DDL augments an
// unqualified name with the session's current schema, which the review does
// not know; when one schema holds a table of the name, a statement naming
// it either means that table or fails. It returns nil when the name may
// mean another database, or no table, or several. A database..name, the
// PUBLIC schema's, is not judged.
//
// sql: https://docs.snowflake.com/en/sql-reference/name-resolution
func (x schemaIndex) table(ref *ast.ObjectName) *metadata.TableMetadata {
	if x.schema == nil || ref == nil {
		return nil
	}
	if !ref.Database.IsEmpty() && (ref.Database.Normalize() != x.schema.GetName() || ref.Schema.IsEmpty()) {
		return nil
	}
	name := ref.Name.Normalize()
	var found *metadata.TableMetadata
	for _, s := range x.schema.GetSchemas() {
		if !ref.Schema.IsEmpty() && ref.Schema.Normalize() != s.GetName() {
			continue
		}
		for _, t := range s.GetTables() {
			if t.GetName() == name {
				if found != nil {
					return nil
				}
				found = t
			}
		}
	}
	return found
}

// syncedConstraint returns the kind of the synced table's constraint of the
// normalized name: the primary key or a unique constraint, by its index, a
// foreign key, or a check; ok is false when the table has none of the name.
// CREATE INDEX, which makes a hybrid table's secondary index, has no UNIQUE,
// so a unique index backs a unique constraint.
//
// sql: https://docs.snowflake.com/en/sql-reference/sql/create-index
func syncedConstraint(t *metadata.TableMetadata, name string) (constraintKind, bool) {
	for _, idx := range t.GetIndexes() {
		if idx.GetName() != name {
			continue
		}
		switch {
		case idx.GetPrimary():
			return primaryKey, true
		case idx.GetUnique():
			return uniqueKey, true
		}
	}
	for _, fk := range t.GetForeignKeys() {
		if fk.GetName() == name {
			return foreignKey, true
		}
	}
	for _, ck := range t.GetCheckConstraints() {
		if ck.GetName() == name {
			return checkKey, true
		}
	}
	return unknownKind, false
}
