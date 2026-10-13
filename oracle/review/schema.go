package review

import (
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
)

// backupDatabase is the schema Bytebase writes an Oracle prior backup to.
//
// bytebase: backend/common/engine.go — BackupDatabaseNameOfEngine
const backupDatabase = "BBDATAARCHIVE"

// unknownSchema is the key schema of a name in the target's database when
// the database's name is not known: it is the same schema as itself, and
// may be any schema a reference names.
const unknownSchema = "\x00"

// schemaIndex looks objects up in a target's synced schema. An Oracle
// database is a schema, and Bytebase syncs it into one SchemaMetadata
// named "", holding the objects of the schema the database is named
// after; a SchemaMetadata named otherwise holds that schema's objects. A
// name another schema qualifies is in no synced schema. Names compare
// exactly: Oracle folds an unquoted name to upper case, which the parser
// has done, and the snapshot holds the dictionary's names.
//
// bytebase: backend/plugin/db/oracle/sync.go — SyncDBSchema
type schemaIndex struct {
	schema *metadata.DatabaseSchemaMetadata
}

// db returns the schema unqualified names resolve in when the change
// starts: the target's database, which Bytebase makes the session's
// CURRENT_SCHEMA, or unknownSchema.
//
// bytebase: backend/plugin/db/oracle/oracle.go — Open (ALTER SESSION SET CURRENT_SCHEMA)
func (x schemaIndex) db() string {
	if name := x.schema.GetName(); name != "" {
		return name
	}
	return unknownSchema
}

// own reports whether a key names an object of the target's database.
func (x schemaIndex) own(k tableKey) bool {
	return k.schema != "" && k.schema != unknownSchema && k.schema == x.db()
}

// schemaOf returns the synced schema a key names an object of, nil when
// the snapshot holds none of it.
func (x schemaIndex) schemaOf(k tableKey) *metadata.SchemaMetadata {
	if x.schema == nil || k.schema == "" || k.schema == unknownSchema {
		return nil
	}
	for _, s := range x.schema.GetSchemas() {
		name := s.GetName()
		if name == "" {
			name = x.db()
		}
		if name == k.schema {
			return s
		}
	}
	return nil
}

// table returns the synced table a key names, nil for none.
func (x schemaIndex) table(k tableKey) *metadata.TableMetadata {
	for _, t := range x.schemaOf(k).GetTables() {
		if t.GetName() == k.name {
			return t
		}
	}
	return nil
}

// relationKind returns "table" or "view" for the synced table or view a
// key names, "" for anything else or nothing.
func (x schemaIndex) relationKind(k tableKey) string {
	s := x.schemaOf(k)
	for _, t := range s.GetTables() {
		if t.GetName() == k.name {
			return "table"
		}
	}
	for _, v := range s.GetViews() {
		if v.GetName() == k.name {
			return "view"
		}
	}
	return ""
}

// lacksRelation reports whether the snapshot certainly holds no table,
// view, or materialized view a key names: it holds the key's schema, and
// none of the name.
func (x schemaIndex) lacksRelation(k tableKey) bool {
	s := x.schemaOf(k)
	if s == nil {
		return false
	}
	for _, t := range s.GetTables() {
		if t.GetName() == k.name {
			return false
		}
	}
	for _, v := range s.GetViews() {
		if v.GetName() == k.name {
			return false
		}
	}
	for _, v := range s.GetMaterializedViews() {
		if v.GetName() == k.name {
			return false
		}
	}
	return true
}

// constraintKind returns the kind of the synced table's constraint of the
// name: the primary key or a unique constraint, by its index, a foreign
// key, or a check, which in Oracle a named NOT NULL constraint is too; ok
// is false when the snapshot lists none of the name. A unique index that
// backs no constraint is none: DROP CONSTRAINT refuses it.
func constraintKind(t *metadata.TableMetadata, name string) (ast.ConstraintType, bool) {
	for _, idx := range t.GetIndexes() {
		if idx.GetName() != name {
			continue
		}
		switch {
		case idx.GetPrimary():
			return ast.CONSTRAINT_PRIMARY, true
		case idx.GetUnique() && idx.GetIsConstraint():
			return ast.CONSTRAINT_UNIQUE, true
		}
	}
	for _, fk := range t.GetForeignKeys() {
		if fk.GetName() == name {
			return ast.CONSTRAINT_FOREIGN, true
		}
	}
	for _, ck := range t.GetCheckConstraints() {
		if ck.GetName() == name {
			return ast.CONSTRAINT_CHECK, true
		}
	}
	return 0, false
}

// syncedConstraint is a constraint of the target's database: the table it
// is on and its kind.
type syncedConstraint struct {
	table tableKey
	kind  ast.ConstraintType
}

// constraintNames returns the constraints of the target's database by
// name. Oracle names constraints per schema, not per table.
func (x schemaIndex) constraintNames() map[string]syncedConstraint {
	out := make(map[string]syncedConstraint)
	db := x.db()
	for _, t := range x.schemaOf(tableKey{schema: db}).GetTables() {
		k := tableKey{schema: db, name: t.GetName()}
		names := make([]string, 0, len(t.GetIndexes())+len(t.GetForeignKeys())+len(t.GetCheckConstraints()))
		for _, idx := range t.GetIndexes() {
			names = append(names, idx.GetName())
		}
		for _, fk := range t.GetForeignKeys() {
			names = append(names, fk.GetName())
		}
		for _, ck := range t.GetCheckConstraints() {
			names = append(names, ck.GetName())
		}
		for _, name := range names {
			if kind, ok := constraintKind(t, name); ok {
				out[name] = syncedConstraint{table: k, kind: kind}
			}
		}
	}
	return out
}
