package analysis

import (
	"strings"

	"github.com/bytebase/omni/redshift"
	"github.com/bytebase/omni/redshift/ast"
)

// ChangedResource is a resource a script writes.
type ChangedResource struct {
	// Database is the resource's database. An unqualified table is in the
	// current database.
	Database string
	// Schema is the resource's schema. An unqualified table is in the current
	// schema, which a SET search_path moves.
	Schema string
	// Table is the changed table.
	Table string
	// AffectsData marks a table change that touches the table's existing rows,
	// such as ALTER, DROP or TRUNCATE. An INSERT that only adds rows does not.
	AffectsData bool
}

// StatementChanges is what a script changes.
type StatementChanges struct {
	// Resources holds a resource each time the script names one, duplicates
	// included, so a caller that deduplicates keeps a flag any occurrence set.
	Resources []ChangedResource
	// DMLStatements holds the text of each DML statement whose rows can be
	// estimated, ending with ';' and prefixed with the search path in force.
	DMLStatements []string
	// DMLCount counts the DML statements, including those DMLStatements omits.
	DMLCount int
	// InsertCount counts the rows INSERT ... VALUES writes.
	InsertCount int
}

// ExtractChangedResources returns the resources a script writes and its DML
// counts. database and schema are the session's current database and schema;
// an empty schema means "public".
func ExtractChangedResources(statement, database, schema string) (*StatementChanges, error) {
	if schema == "" {
		schema = "public"
	}
	stmts, err := redshift.Parse(statement)
	if err != nil {
		return nil, err
	}
	rawSummary := redshift.ExtractChangedResourcesFromStatements(stmts, database, schema)
	ddlAffected := ddlAffectedTables(stmts, database, schema)

	changes := &StatementChanges{}
	for _, table := range rawSummary.Tables {
		changes.Resources = append(changes.Resources, ChangedResource{
			Database:    table.Database,
			Schema:      table.Schema,
			Table:       table.Name,
			AffectsData: ddlAffected[tableKey{database: table.Database, schema: table.Schema, name: table.Name}],
		})
	}

	searchPath := schema
	// replaySearchPath is the search path set after the starting one, which the EXPLAIN connection
	// does not have.
	var replaySearchPath []string
	sampleDML := func(text string) {
		changes.DMLCount++
		text = strings.TrimSpace(text)
		if len(replaySearchPath) > 0 {
			text = withSearchPath(text, replaySearchPath)
		}
		changes.DMLStatements = append(changes.DMLStatements, text)
	}
	addTable := func(rv *ast.RangeVar, affectedTable bool) {
		tableDatabase, tableSchema := rv.Catalogname, rv.Schemaname
		if tableDatabase == "" {
			tableDatabase = database
		}
		if tableSchema == "" {
			tableSchema = searchPath
		}
		changes.Resources = append(changes.Resources, ChangedResource{
			Database:    tableDatabase,
			Schema:      tableSchema,
			Table:       rv.Relname,
			AffectsData: affectedTable,
		})
	}
	for _, stmt := range stmts {
		switch n := stmt.AST.(type) {
		case *ast.InsertStmt:
			if rows, ok := insertValuesRowCount(n); ok {
				changes.InsertCount += rows
			} else {
				sampleDML(stmt.Text)
			}
		case *ast.UpdateStmt, *ast.DeleteStmt:
			sampleDML(stmt.Text)
		case *ast.MergeStmt:
			// Redshift documents EXPLAIN only for SELECT, CREATE TABLE AS, INSERT, UPDATE, and DELETE.
			changes.DMLCount++
			// omni records no MERGE target.
			if n.Relation != nil {
				addTable(n.Relation, false)
			}
		case *ast.TruncateStmt:
			if n.Relations == nil {
				continue
			}
			for _, item := range n.Relations.Items {
				if rv, ok := item.(*ast.RangeVar); ok {
					addTable(rv, true)
				}
			}
		case *ast.VariableSetStmt:
			if schemas, ok := searchPathFromSet(n); ok {
				searchPath, replaySearchPath = schema, schemas
				if len(schemas) > 0 {
					searchPath = schemas[0]
				}
			}
		default:
		}
	}
	return changes, nil
}

type tableKey struct {
	database string
	schema   string
	name     string
}

// ddlAffectedTables returns the tables whose existing rows the script's DDL affects.
// redshift's walk also flags DML targets and merges flags across statements, so it
// reads the script's DDL statements only.
func ddlAffectedTables(stmts []redshift.Statement, database, schema string) map[tableKey]bool {
	var ddlStatements []redshift.Statement
	for _, stmt := range stmts {
		switch stmt.AST.(type) {
		case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
		default:
			ddlStatements = append(ddlStatements, stmt)
		}
	}
	summary := redshift.ExtractChangedResourcesFromStatements(ddlStatements, database, schema)
	affected := make(map[tableKey]bool)
	for _, table := range summary.Tables {
		if table.Affected {
			affected[tableKey{database: table.Database, schema: table.Schema, name: table.Name}] = true
		}
	}
	return affected
}

// insertValuesRowCount returns how many rows INSERT ... VALUES or INSERT ... DEFAULT
// VALUES adds, and false for an INSERT whose rows come from a query.
func insertValuesRowCount(n *ast.InsertStmt) (int, bool) {
	if n.SelectStmt == nil {
		return 1, true
	}
	if s, ok := n.SelectStmt.(*ast.SelectStmt); ok && s.ValuesLists != nil {
		return len(s.ValuesLists.Items), true
	}
	return 0, false
}

// searchPathFromSet returns the schemas a SET search_path names, skipping "$user", or none for RESET
// or SET ... TO DEFAULT, which restore the starting search path. ok is false for any other statement.
// Like redshift.ExtractChangedResources, unqualified names resolve to the first schema.
func searchPathFromSet(n *ast.VariableSetStmt) ([]string, bool) {
	if !strings.EqualFold(n.Name, "search_path") {
		return nil, false
	}
	if n.Kind == ast.VAR_SET_DEFAULT || n.Kind == ast.VAR_RESET {
		return nil, true
	}
	if n.Args == nil {
		return nil, false
	}
	var schemas []string
	for _, arg := range n.Args.Items {
		c, ok := arg.(*ast.A_Const)
		if !ok {
			continue
		}
		s, ok := c.Val.(*ast.String)
		if !ok || s.Str == "" || s.Str == "$user" {
			continue
		}
		schemas = append(schemas, s.Str)
	}
	return schemas, len(schemas) > 0
}

// searchPathPrefix starts the statement withSearchPath puts before a DML statement.
const searchPathPrefix = "SET LOCAL search_path TO "

// withSearchPath returns statement prefixed by the SET LOCAL that gives the
// EXPLAIN connection the search path it was sampled under.
func withSearchPath(statement string, schemas []string) string {
	quoted := make([]string, len(schemas))
	for i, schema := range schemas {
		quoted[i] = `"` + strings.ReplaceAll(schema, `"`, `""`) + `"`
	}
	return searchPathPrefix + strings.Join(quoted, ", ") + ";\n" + statement
}
