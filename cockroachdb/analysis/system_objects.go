package analysis

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

var (
	// systemDatabases is the list of system or internal databases.
	systemDatabases = map[string]bool{
		"postgres": true,
		"system":   true,
	}

	// systemSchemas is the list of system schemas that a sync excludes.
	// https://www.cockroachlabs.com/docs/stable/system-catalogs
	systemSchemas = map[string]bool{
		"crdb_internal":      true,
		"information_schema": true,
		"pg_catalog":         true,
		"pg_extension":       true,
	}

	// SystemSchemaWhereClause lists systemSchemas as SQL string literals, in
	// name order so that a query built from it is the same text every run.
	SystemSchemaWhereClause = func() string {
		var schemas []string
		for _, schema := range slices.Sorted(maps.Keys(systemSchemas)) {
			schemas = append(schemas, fmt.Sprintf("'%s'", schema))
		}
		return strings.Join(schemas, ",")
	}()
)

// IsSystemDatabase reports whether database is one CockroachDB manages.
func IsSystemDatabase(database string) bool {
	_, ok := systemDatabases[database]
	return ok
}

// IsSystemSchema reports whether schema is one CockroachDB manages.
func IsSystemSchema(schema string) bool {
	_, ok := systemSchemas[schema]
	return ok
}
