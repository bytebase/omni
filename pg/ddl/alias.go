package ddl

import "strings"

// typeAliases maps SQL-standard type names to their PostgreSQL internal names.
var typeAliases = map[string]string{
	"integer":                      "int4",
	"int":                          "int4",
	"smallint":                     "int2",
	"bigint":                       "int8",
	"real":                         "float4",
	"double precision":             "float8",
	"decimal":                      "numeric",
	"character varying":            "varchar",
	"character":                    "bpchar",
	"boolean":                      "bool",
	"timestamp without time zone":  "timestamp",
	"timestamp with time zone":     "timestamptz",
	"time without time zone":       "time",
	"time with time zone":          "timetz",
}

// resolveAlias returns the internal PG type name for the given SQL name.
// If no alias exists, the original name is returned (lowercased).
func resolveAlias(name string) string {
	lower := strings.ToLower(name)
	if alias, ok := typeAliases[lower]; ok {
		return alias
	}
	return lower
}
