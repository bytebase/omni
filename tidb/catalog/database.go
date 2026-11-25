package catalog

import "strings"

type Database struct {
	Name       string
	Charset    string
	Collation  string
	Tables     map[string]*Table // lowered name -> Table
	Views      map[string]*View
	Functions  map[string]*Routine // lowered name -> stored function
	Procedures map[string]*Routine // lowered name -> stored procedure
	Triggers   map[string]*Trigger // lowered name -> trigger
	Events     map[string]*Event   // lowered name -> event
	Sequences  map[string]*Sequence

	// TiDB-specific: database-level placement policy. Empty when unset.
	PlacementPolicy string

	// TiDB: database-level TiFlash replica count (0 = not set).
	// The replica count applies to all tables in the database that
	// don't override it at the table level (upstream resolution is
	// executed asynchronously by the placement scheduler, not the
	// DDL path). Omni stores it as declared metadata only.
	TiFlashReplica        int
	TiFlashLocationLabels []string

	// TempTables are the session's local temporary tables (CREATE TEMPORARY
	// TABLE), by lowered name. TiDB keeps them apart from the schema's
	// objects: one may share a name with a table, view or sequence, and then
	// it shadows the table wherever a statement names it (SHOW CREATE TABLE,
	// DROP TABLE, a query), while SHOW TABLES and information_schema list none
	// (read on v8.5.5). A global temporary table is a schema object and lives
	// in Tables.
	TempTables map[string]*Table
}

func newDatabase(name, charset, collation string) *Database {
	return &Database{
		Name:       name,
		Charset:    charset,
		Collation:  collation,
		Tables:     make(map[string]*Table),
		TempTables: make(map[string]*Table),
		Views:      make(map[string]*View),
		Functions:  make(map[string]*Routine),
		Procedures: make(map[string]*Routine),
		Triggers:   make(map[string]*Trigger),
		Events:     make(map[string]*Event),
		Sequences:  make(map[string]*Sequence),
	}
}

// GetTable returns the table a statement naming name reads: a local temporary
// table of that name before the schema's table it shadows.
func (db *Database) GetTable(name string) *Table {
	key := toLower(name)
	if tbl := db.TempTables[key]; tbl != nil {
		return tbl
	}
	return db.Tables[key]
}

func toLower(s string) string {
	return strings.ToLower(s)
}
