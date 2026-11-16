package analysis

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// ChangedResource is a resource a statement writes.
type ChangedResource struct {
	// Database is the resource's database. An unqualified table is in the
	// current database.
	Database string
	// Table is the changed table, empty for a database-only change.
	Table string
	// DatabaseOnly marks a non-table object (view, materialized view, routine)
	// whose own name carries a database qualifier. It has no table to scope,
	// only its target database.
	DatabaseOnly bool
	// AffectsData marks a table change that touches the table's existing rows,
	// such as DROP, TRUNCATE or ALTER.
	AffectsData bool
}

// StatementChanges is what one statement changes.
type StatementChanges struct {
	// Resources holds the resources the statement writes, in the order it
	// names them, duplicates included.
	Resources []ChangedResource
	// DMLStatements holds the text of a DML statement whose rows can be
	// estimated, ending with ';'.
	DMLStatements []string
	// DMLCount counts the DML statements, including those DMLStatements omits.
	DMLCount int
	// InsertCount counts the rows INSERT ... VALUES writes.
	InsertCount int
}

// ExtractChangedResources returns the resources one statement writes and its
// DML counts. database is the current database, and text is the statement's
// text, which a DML statement is sampled from.
func ExtractChangedResources(node ast.Node, database, text string) StatementChanges {
	var changes StatementChanges
	addTableIn := func(name *ast.ObjectName, defaultDatabase string, affectsData bool) {
		db, table := objectName(name, defaultDatabase)
		if table == "" {
			return
		}
		changes.Resources = append(changes.Resources, ChangedResource{Database: db, Table: table, AffectsData: affectsData})
	}
	addTable := func(name *ast.ObjectName, affectsData bool) {
		addTableIn(name, database, affectsData)
	}
	// addObjectDatabase records a database-only write target for a non-table
	// object DDL (view, materialized view, routine) ONLY when the object's own
	// name carries an explicit database qualifier. An unqualified name records
	// nothing, leaving the request database to stand in.
	addObjectDatabase := func(name *ast.ObjectName) {
		if db, _ := objectName(name, ""); db != "" {
			changes.Resources = append(changes.Resources, ChangedResource{Database: db, DatabaseOnly: true})
		}
	}

	switch n := node.(type) {
	case *ast.CreateTableStmt:
		addTable(n.Name, false)

	case *ast.DropTableStmt:
		addTable(n.Name, true)

	case *ast.TruncateTableStmt:
		addTable(n.Target, true)

	case *ast.AlterTableStmt:
		db, table := objectName(n.Name, database)
		if table != "" {
			changes.Resources = append(changes.Resources, ChangedResource{Database: db, Table: table, AffectsData: true})
		}
		// SWAP WITH exchanges the altered table's data with the partner's,
		// which resolves in the altered table's database (engine-verified: a
		// session in another database still finds it there), and both tables
		// change.
		for _, action := range n.Actions {
			if action != nil && action.Type == ast.AlterSwapTable && action.NewName != "" {
				changes.Resources = append(changes.Resources, ChangedResource{Database: db, Table: action.NewName, AffectsData: true})
			}
		}

	case *ast.CreateIndexStmt:
		addTable(n.Table, false)

	case *ast.DropIndexStmt:
		addTable(n.Table, false)

	case *ast.LoadDataStmt:
		// A LOAD LABEL's database qualifier is the database the job runs in,
		// and an unqualified data-description target resolves against it
		// (engine-verified: with the session in another database, LOAD LABEL
		// db_b.job (... INTO TABLE t ...) loads db_b.t). LOAD DATA carries no
		// label and keeps the request database.
		loadDatabase := database
		if label := labelDatabase(n.Label); label != "" {
			loadDatabase = label
		}
		for _, desc := range n.DataDescs {
			if desc != nil {
				addTableIn(desc.Target, loadDatabase, false)
			}
		}

	case *ast.CopyIntoStmt:
		// COPY INTO writes its target table; the stage is a source. StarRocks's
		// grammar has no COPY statement and the pinned 3.4.10 server rejects
		// it, but this parser still builds the node it inherited from Doris, so
		// the case keeps the two engines' analysis consistent.
		addTable(n.Target, false)

	case *ast.InsertStmt:
		// INSERT INTO FILES(...) writes to external storage and leaves Target
		// nil, so it records no table while still counting as a DML statement.
		addTable(n.Target, false)
		if len(n.Values) > 0 {
			changes.InsertCount += len(n.Values)
			break
		}
		changes.DMLCount++
		changes.DMLStatements = append(changes.DMLStatements, statementText(text))

	case *ast.UpdateStmt:
		// StarRocks updates a single target; FROM names read sources, not write
		// targets.
		addTable(n.Target, false)
		changes.DMLCount++
		changes.DMLStatements = append(changes.DMLStatements, statementText(text))

	case *ast.DeleteStmt:
		// USING names read sources, as FROM does for UPDATE.
		addTable(n.Target, false)
		changes.DMLCount++
		changes.DMLStatements = append(changes.DMLStatements, statementText(text))

	case *ast.MergeStmt:
		addTable(n.Target, false)
		changes.DMLCount++
		changes.DMLStatements = append(changes.DMLStatements, statementText(text))

	// Non-table object DDL: gate by the object's own explicit database
	// qualifier only. CREATE FUNCTION is absent because its AST keeps the name
	// inside RawSignature rather than as an ObjectName.
	case *ast.CreateViewStmt:
		addObjectDatabase(n.Name)
	case *ast.AlterViewStmt:
		addObjectDatabase(n.Name)
	case *ast.DropViewStmt:
		addObjectDatabase(n.Name)
	case *ast.CreateMTMVStmt:
		addObjectDatabase(n.Name)
	case *ast.AlterMTMVStmt:
		addObjectDatabase(n.Name)
	case *ast.DropMTMVStmt:
		addObjectDatabase(n.Name)
	case *ast.DropFunctionStmt:
		addObjectDatabase(n.Name)
	case *ast.CreateProcedureStmt:
		addObjectDatabase(n.Name)
	case *ast.DropProcedureStmt:
		addObjectDatabase(n.Name)
	}

	return changes
}

// objectName splits an ObjectName into its database and object name. StarRocks
// writes a name as table, db.table or catalog.db.table, so the object is the
// last part and the database the one before it; a catalog qualifier is
// dropped, since changed resources are tracked per database.
func objectName(name *ast.ObjectName, defaultDatabase string) (string, string) {
	if name == nil || len(name.Parts) == 0 {
		return defaultDatabase, ""
	}
	parts := name.Parts
	database := defaultDatabase
	if len(parts) >= 2 {
		database = parts[len(parts)-2]
	}
	return database, parts[len(parts)-1]
}

// statementText returns the text of a statement, ensuring it ends with a
// semicolon.
func statementText(text string) string {
	text = strings.TrimSpace(text)
	if !strings.HasSuffix(text, ";") {
		text += ";"
	}
	return text
}

// labelDatabase returns the database part of a LOAD LABEL, which the parser
// joins as "db.label", or "" when the label carries no qualifier.
func labelDatabase(label string) string {
	if dot := strings.LastIndexByte(label, '.'); dot >= 0 {
		return label[:dot]
	}
	return ""
}
