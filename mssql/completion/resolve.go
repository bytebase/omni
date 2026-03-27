package completion

import (
	"strings"

	"github.com/bytebase/omni/mssql/parser"
)

// resolve converts parser CandidateSet into completion Candidates.
// cat may be nil; in that case only keyword and hardcoded type candidates are returned.
func resolve(cs *parser.CandidateSet, cat interface{}, sql string, cursorOffset int) []Candidate {
	if cs == nil {
		return nil
	}
	var result []Candidate

	// Token candidates -> keywords
	for _, tok := range cs.Tokens {
		name := parser.TokenName(tok)
		if name == "" {
			continue
		}
		result = append(result, Candidate{Text: name, Type: CandidateKeyword})
	}

	// Rule candidates -> catalog objects or hardcoded lists
	for _, rc := range cs.Rules {
		result = append(result, resolveRule(rc.Rule, cat, sql, cursorOffset)...)
	}

	return dedup(result)
}

// tsqlTypes are the hardcoded T-SQL type keywords returned for the "type_name" rule.
var tsqlTypes = []string{
	"INT", "BIGINT", "SMALLINT", "TINYINT",
	"VARCHAR", "NVARCHAR", "CHAR", "NCHAR",
	"TEXT", "NTEXT",
	"DATETIME", "DATETIME2", "DATE", "TIME",
	"DECIMAL", "NUMERIC", "FLOAT", "REAL",
	"BIT",
	"MONEY", "SMALLMONEY",
	"UNIQUEIDENTIFIER",
	"XML",
	"VARBINARY", "IMAGE",
	"SQL_VARIANT",
}

// resolveRule resolves a single grammar rule into completion candidates.
func resolveRule(rule string, cat interface{}, sql string, cursorOffset int) []Candidate {
	// type_name always returns hardcoded T-SQL types regardless of catalog.
	if rule == "type_name" {
		return resolveTypeNames()
	}

	// All other rules require a catalog.
	if cat == nil {
		return nil
	}

	switch rule {
	case "table_ref":
		return resolveCatalogRule(cat, CandidateTable)
	case "columnref":
		// Extract table refs visible at cursor for context-aware column completion.
		_ = extractTableRefs(sql, cursorOffset)
		return resolveCatalogRule(cat, CandidateColumn)
	case "schema_ref":
		return resolveCatalogRule(cat, CandidateSchema)
	case "func_name":
		return resolveCatalogRule(cat, CandidateFunction)
	case "proc_ref":
		return resolveCatalogRule(cat, CandidateProcedure)
	case "index_ref":
		return resolveCatalogRule(cat, CandidateIndex)
	case "trigger_ref":
		return resolveCatalogRule(cat, CandidateTrigger)
	case "database_ref":
		// Database names — requires catalog support.
		return nil
	case "sequence_ref":
		return resolveCatalogRule(cat, CandidateSequence)
	case "login_ref":
		return resolveCatalogRule(cat, CandidateLogin)
	case "user_ref":
		return resolveCatalogRule(cat, CandidateUser)
	case "role_ref":
		return resolveCatalogRule(cat, CandidateRole)
	}
	return nil
}

// resolveTypeNames returns hardcoded T-SQL type keyword candidates.
func resolveTypeNames() []Candidate {
	result := make([]Candidate, 0, len(tsqlTypes))
	for _, t := range tsqlTypes {
		result = append(result, Candidate{Text: t, Type: CandidateType_})
	}
	return result
}

// resolveCatalogRule is a placeholder that will resolve catalog-backed rules
// once the mssql/catalog package exists. Currently returns nil because no
// catalog implementation is available yet.
func resolveCatalogRule(_ interface{}, _ CandidateType) []Candidate {
	// TODO: resolve against mssql/catalog once it exists.
	return nil
}

// dedup removes duplicate candidates (same text+type, case-insensitive).
func dedup(cs []Candidate) []Candidate {
	type key struct {
		text string
		typ  CandidateType
	}
	seen := make(map[key]bool, len(cs))
	result := make([]Candidate, 0, len(cs))
	for _, c := range cs {
		k := key{strings.ToLower(c.Text), c.Type}
		if seen[k] {
			continue
		}
		seen[k] = true
		result = append(result, c)
	}
	return result
}
