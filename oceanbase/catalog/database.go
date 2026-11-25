package catalog

import "strings"

type Database struct {
	Name      string
	Charset   string
	Collation string
	Tables    map[string]*Table // lowered name -> Table
	Views     map[string]*View
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
	// Sequences are the database's sequences, whose names are apart from its
	// tables' and views'.
	Sequences map[string]*Sequence // lowered name -> sequence
}

func newDatabase(name, charset, collation string) *Database {
	return &Database{
		Name:                      name,
		Charset:                   charset,
		Collation:                 collation,
		Tables:                    make(map[string]*Table),
		Views:                     make(map[string]*View),
		MaterializedViewLogs:      make(map[string]bool),
		MaterializedViewLogsKnown: make(map[string]bool),
		Functions:                 make(map[string]*Routine),
		Procedures:                make(map[string]*Routine),
		Triggers:                  make(map[string]*Trigger),
		Events:                    make(map[string]*Event),
		Sequences:                 make(map[string]*Sequence),
	}
}

func (db *Database) GetTable(name string) *Table {
	return db.Tables[toLower(name)]
}

func toLower(s string) string {
	return strings.ToLower(s)
}
