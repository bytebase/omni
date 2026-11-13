package analysis

import "github.com/bytebase/omni/snowflake/ast"

// normalizeIdentifier returns the folded form of an identifier: unquoted is
// upper-cased, quoted is kept verbatim.
func normalizeIdentifier(id ast.Ident) string {
	if id.IsEmpty() {
		return ""
	}
	return id.Normalize()
}

// normalizeObjectName splits an object name into its folded (database,
// schema, table) parts, using the fallbacks for the parts the name omits.
func normalizeObjectName(objectName *ast.ObjectName, fallbackDatabaseName, fallbackSchemaName string) (database, schema, table string) {
	if objectName == nil {
		return "", "", ""
	}

	database = fallbackDatabaseName
	if d := normalizeIdentifier(objectName.Database); d != "" {
		database = d
	}

	schema = fallbackSchemaName
	if s := normalizeIdentifier(objectName.Schema); s != "" {
		schema = s
	}

	table = normalizeIdentifier(objectName.Name)
	return database, schema, table
}
