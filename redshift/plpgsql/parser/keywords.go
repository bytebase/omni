package parser

import "strings"

// PLKeywordCategory represents PL/pgSQL keyword categories.
// These are distinct from the SQL parser keyword categories:
// PL/pgSQL has its own reserved/unreserved classification
// (see PostgreSQL src/pl/plpgsql/src/pl_reserved_kwlist.h and pl_unreserved_kwlist.h).
type PLKeywordCategory int

const (
	// PLReserved keywords cannot be used as variable names in PL/pgSQL.
	PLReserved PLKeywordCategory = iota
	// PLUnreserved keywords are recognized by the PL/pgSQL parser but
	// can still be used as variable names in declarations.
	PLUnreserved
)

// PLKeyword represents a PL/pgSQL keyword with its name and category.
type PLKeyword struct {
	Name     string
	Category PLKeywordCategory
}

// PLReservedKeywords is the complete list of PL/pgSQL reserved keywords.
// Source: PostgreSQL src/pl/plpgsql/src/pl_reserved_kwlist.h
var PLReservedKeywords = []PLKeyword{
	{"all", PLReserved},
	{"begin", PLReserved},
	{"by", PLReserved},
	{"case", PLReserved},
	{"declare", PLReserved},
	{"else", PLReserved},
	{"end", PLReserved},
	{"execute", PLReserved},
	{"for", PLReserved},
	{"foreach", PLReserved},
	{"from", PLReserved},
	{"if", PLReserved},
	{"in", PLReserved},
	{"into", PLReserved},
	{"loop", PLReserved},
	{"not", PLReserved},
	{"null", PLReserved},
	{"or", PLReserved},
	{"strict", PLReserved},
	{"then", PLReserved},
	{"to", PLReserved},
	{"using", PLReserved},
	{"when", PLReserved},
	{"while", PLReserved},
}

// PLUnreservedKeywords is the complete list of PL/pgSQL unreserved keywords.
// Source: PostgreSQL src/pl/plpgsql/src/pl_unreserved_kwlist.h
var PLUnreservedKeywords = []PLKeyword{
	{"absolute", PLUnreserved},
	{"alias", PLUnreserved},
	{"and", PLUnreserved},
	{"array", PLUnreserved},
	{"assert", PLUnreserved},
	{"backward", PLUnreserved},
	{"call", PLUnreserved},
	{"chain", PLUnreserved},
	{"close", PLUnreserved},
	{"collate", PLUnreserved},
	{"column", PLUnreserved},
	{"column_name", PLUnreserved},
	{"commit", PLUnreserved},
	{"constant", PLUnreserved},
	{"constraint", PLUnreserved},
	{"constraint_name", PLUnreserved},
	{"continue", PLUnreserved},
	{"current", PLUnreserved},
	{"cursor", PLUnreserved},
	{"datatype", PLUnreserved},
	{"debug", PLUnreserved},
	{"default", PLUnreserved},
	{"detail", PLUnreserved},
	{"diagnostics", PLUnreserved},
	{"do", PLUnreserved},
	{"dump", PLUnreserved},
	{"elseif", PLUnreserved},
	{"elsif", PLUnreserved},
	{"errcode", PLUnreserved},
	{"error", PLUnreserved},
	{"exception", PLUnreserved},
	{"exit", PLUnreserved},
	{"fetch", PLUnreserved},
	{"first", PLUnreserved},
	{"forward", PLUnreserved},
	{"get", PLUnreserved},
	{"hint", PLUnreserved},
	{"import", PLUnreserved},
	{"info", PLUnreserved},
	{"insert", PLUnreserved},
	{"is", PLUnreserved},
	{"last", PLUnreserved},
	{"log", PLUnreserved},
	{"merge", PLUnreserved},
	{"message", PLUnreserved},
	{"message_text", PLUnreserved},
	{"move", PLUnreserved},
	{"next", PLUnreserved},
	{"no", PLUnreserved},
	{"notice", PLUnreserved},
	{"open", PLUnreserved},
	{"option", PLUnreserved},
	{"perform", PLUnreserved},
	{"pg_context", PLUnreserved},
	{"pg_datatype_name", PLUnreserved},
	{"pg_exception_context", PLUnreserved},
	{"pg_exception_detail", PLUnreserved},
	{"pg_exception_hint", PLUnreserved},
	{"pg_routine_oid", PLUnreserved},
	{"print_strict_params", PLUnreserved},
	{"prior", PLUnreserved},
	{"query", PLUnreserved},
	{"raise", PLUnreserved},
	{"relative", PLUnreserved},
	{"return", PLUnreserved},
	{"returned_sqlstate", PLUnreserved},
	{"reverse", PLUnreserved},
	{"rollback", PLUnreserved},
	{"row_count", PLUnreserved},
	{"rowtype", PLUnreserved},
	{"schema", PLUnreserved},
	{"schema_name", PLUnreserved},
	{"scroll", PLUnreserved},
	{"slice", PLUnreserved},
	{"sqlstate", PLUnreserved},
	{"stacked", PLUnreserved},
	{"table", PLUnreserved},
	{"table_name", PLUnreserved},
	{"type", PLUnreserved},
	{"use_column", PLUnreserved},
	{"use_variable", PLUnreserved},
	{"variable_conflict", PLUnreserved},
	{"warning", PLUnreserved},
}

// plKeywordMap maps lowercase keyword names to their definitions.
var plKeywordMap = func() map[string]*PLKeyword {
	m := make(map[string]*PLKeyword, len(PLReservedKeywords)+len(PLUnreservedKeywords))
	for i := range PLReservedKeywords {
		m[PLReservedKeywords[i].Name] = &PLReservedKeywords[i]
	}
	for i := range PLUnreservedKeywords {
		m[PLUnreservedKeywords[i].Name] = &PLUnreservedKeywords[i]
	}
	return m
}()

// LookupPLKeyword looks up a PL/pgSQL keyword by name (case-insensitive).
// Returns the keyword category and true if found, or (0, false) if not a PL/pgSQL keyword.
func LookupPLKeyword(name string) (PLKeywordCategory, bool) {
	kw, ok := plKeywordMap[strings.ToLower(name)]
	if !ok {
		return 0, false
	}
	return kw.Category, true
}
