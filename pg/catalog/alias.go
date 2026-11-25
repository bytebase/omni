package catalog

import "strings"

// typeAliases maps SQL-standard type names to their PostgreSQL internal names.
var typeAliases = map[string]string{
	"integer":                     "int4",
	"int":                         "int4",
	"smallint":                    "int2",
	"bigint":                      "int8",
	"real":                        "float4",
	"double precision":            "float8",
	"decimal":                     "numeric",
	"character varying":           "varchar",
	"character":                   "bpchar",
	"boolean":                     "bool",
	"timestamp without time zone": "timestamp",
	"timestamp with time zone":    "timestamptz",
	"time without time zone":      "time",
	"time with time zone":         "timetz",
}

// resolveAlias returns the internal PG type name for the given SQL name.
// If no alias exists, the name is returned as given: the parser has already
// folded an unquoted name and kept a quoted one's case, so folding it again
// would turn "Mood" into mood.
func resolveAlias(name string) string {
	if alias, ok := typeAliases[strings.ToLower(name)]; ok {
		return alias
	}
	return name
}
