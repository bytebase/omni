package driver

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

var (
	systemDatabases = map[string]bool{
		"postgres": true,
		"system":   true,
	}
	systemSchemas = map[string]bool{
		// https://www.cockroachlabs.com/docs/stable/system-catalogs
		"crdb_internal":      true,
		"information_schema": true,
		"pg_catalog":         true,
		"pg_extension":       true,
	}

	// systemSchemaWhereClause is an optimization for getting fewer schema
	// objects. It is in name order so the query text is the same every run.
	systemSchemaWhereClause = func() string {
		schemas := make([]string, 0, len(systemSchemas))
		for _, schema := range slices.Sorted(maps.Keys(systemSchemas)) {
			schemas = append(schemas, fmt.Sprintf("'%s'", schema))
		}
		return strings.Join(schemas, ",")
	}()
)

func isSystemDatabase(database string) bool {
	_, ok := systemDatabases[database]
	return ok
}

func isSystemSchema(schema string) bool {
	_, ok := systemSchemas[schema]
	return ok
}
