package catalog

import "strings"

type Database struct {
	Name      string
	Charset   string
	Collation string
	Tables    map[string]*Table // lowered name -> Table
	Views     map[string]*View
	// ExternalTables holds the lowered names of the external tables a change
	// created. Their columns, location and format live in
	// metadata.ExternalTableMetadata, which the catalog holds from a sync, so a
	// name is all it keeps here; the name shares the relations' namespace (1050).
	ExternalTables map[string]bool
	// MaterializedViewLogs holds the base tables a materialized view log of
	// the database is created on. The log itself is not modeled.
	MaterializedViewLogs map[string]bool
	// MaterializedViewLogsKnown marks the tables whose log state this catalog
	// knows. A table loaded from a snapshot has no entry, so a log of it is
	// unknown rather than known to be missing.
	MaterializedViewLogsKnown map[string]bool
	Functions                 map[string]*Routine // lowered name -> stored function
	Procedures                map[string]*Routine // lowered name -> stored procedure
	Triggers                  map[string]*Trigger // lowered name -> trigger
	Events                    map[string]*Event   // lowered name -> event
}

func newDatabase(name, charset, collation string) *Database {
	return &Database{
		Name:                      name,
		Charset:                   charset,
		Collation:                 collation,
		Tables:                    make(map[string]*Table),
		Views:                     make(map[string]*View),
		ExternalTables:            make(map[string]bool),
		MaterializedViewLogs:      make(map[string]bool),
		MaterializedViewLogsKnown: make(map[string]bool),
		Functions:                 make(map[string]*Routine),
		Procedures:                make(map[string]*Routine),
		Triggers:                  make(map[string]*Trigger),
		Events:                    make(map[string]*Event),
	}
}

func (db *Database) GetTable(name string) *Table {
	return db.Tables[toLower(name)]
}

func toLower(s string) string {
	return strings.ToLower(s)
}
