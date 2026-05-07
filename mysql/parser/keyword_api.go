package parser

import "strings"

// IsKeyword reports whether s is a registered MySQL keyword.
//
// The lookup is case-insensitive and includes both reserved and non-reserved
// keyword categories. It follows MySQL 8.0's keyword symbols, not the parser's
// function-only lexer tokens.
func IsKeyword(s string) bool {
	lower := strings.ToLower(s)
	_, excluded := mysql80KeywordAPIExclusions[lower]
	_, ok := keywords[lower]
	return ok && !excluded
}

var mysql80KeywordAPIExclusions = map[string]struct{}{
	"absent":                 {},
	"adddate":                {},
	"allow_missing_files":    {},
	"auto":                   {},
	"auto_refresh":           {},
	"auto_refresh_source":    {},
	"bernoulli":              {},
	"bit_and":                {},
	"bit_or":                 {},
	"bit_xor":                {},
	"cast":                   {},
	"count":                  {},
	"curdate":                {},
	"curtime":                {},
	"date_add":               {},
	"date_sub":               {},
	"duality":                {},
	"external":               {},
	"external_format":        {},
	"extract":                {},
	"file_format":            {},
	"file_name":              {},
	"file_pattern":           {},
	"file_prefix":            {},
	"files":                  {},
	"group_concat":           {},
	"gtids":                  {},
	"guided":                 {},
	"header":                 {},
	"json_arrayagg":          {},
	"json_duality_object":    {},
	"json_objectagg":         {},
	"library":                {},
	"log":                    {},
	"manual":                 {},
	"materialized":           {},
	"max":                    {},
	"mid":                    {},
	"min":                    {},
	"now":                    {},
	"parallel":               {},
	"parameters":             {},
	"parse_tree":             {},
	"position":               {},
	"qualify":                {},
	"relational":             {},
	"s3":                     {},
	"session_user":           {},
	"sets":                   {},
	"st_collect":             {},
	"std":                    {},
	"stddev":                 {},
	"stddev_pop":             {},
	"stddev_samp":            {},
	"strict_load":            {},
	"subdate":                {},
	"substr":                 {},
	"substring":              {},
	"sum":                    {},
	"sysdate":                {},
	"system_user":            {},
	"tablesample":            {},
	"trim":                   {},
	"uri":                    {},
	"validate":               {},
	"var_pop":                {},
	"var_samp":               {},
	"variance":               {},
	"vector":                 {},
	"verify_key_constraints": {},
}

// IsReservedKeyword reports whether s is a registered MySQL reserved keyword.
func IsReservedKeyword(s string) bool {
	lower := strings.ToLower(s)
	if !IsKeyword(lower) {
		return false
	}
	tok, ok := keywords[lower]
	return ok && isReserved(tok)
}
