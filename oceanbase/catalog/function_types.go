package catalog

import "strings"

// functionReturnType returns the inferred return type for a known function.
// Returns nil for unknown functions (Phase 3 covers ~50 common functions).
func functionReturnType(name string, args []AnalyzedExpr) *ResolvedType {
	switch strings.ToLower(name) {
	// String functions
	case "concat", "concat_ws", "lower", "upper", "trim", "ltrim", "rtrim",
		"substring", "substr", "left", "right", "replace", "reverse",
		"lpad", "rpad", "repeat", "space", "format":
		return &ResolvedType{BaseType: BaseTypeVarchar}

	// Numeric functions
	case "abs", "ceil", "ceiling", "floor", "round", "truncate", "mod":
		return &ResolvedType{BaseType: BaseTypeDecimal}
	case "rand":
		return &ResolvedType{BaseType: BaseTypeDouble}

	// Aggregate functions
	case "count":
		return &ResolvedType{BaseType: BaseTypeBigInt}
	case "sum", "avg":
		return &ResolvedType{BaseType: BaseTypeDecimal}
	case "min", "max":
		return nil // type depends on argument
	case "group_concat":
		return &ResolvedType{BaseType: BaseTypeText}

	// Date/time functions
	case "now", "current_timestamp", "sysdate", "localtime", "localtimestamp":
		return &ResolvedType{BaseType: BaseTypeDateTime}
	case "curdate", "current_date":
		return &ResolvedType{BaseType: BaseTypeDate}
	case "curtime", "current_time":
		return &ResolvedType{BaseType: BaseTypeTime}
	case "year", "month", "day", "hour", "minute", "second",
		"dayofweek", "dayofmonth", "dayofyear", "weekday",
		"quarter", "week", "yearweek":
		return &ResolvedType{BaseType: BaseTypeInt}
	case "date":
		return &ResolvedType{BaseType: BaseTypeDate}
	case "time":
		return &ResolvedType{BaseType: BaseTypeTime}
	case "timestamp":
		return &ResolvedType{BaseType: BaseTypeTimestamp}

	// Type conversion
	case "cast", "convert":
		return nil // handled by CastExprQ

	// Control flow
	case "if":
		return nil // type depends on arguments
	case "nullif":
		return nil // type depends on first argument

	// JSON functions
	case "json_extract", "json_unquote":
		return &ResolvedType{BaseType: BaseTypeJSON}
	case "json_length", "json_depth", "json_valid":
		return &ResolvedType{BaseType: BaseTypeInt}
	case "json_type":
		return &ResolvedType{BaseType: BaseTypeVarchar}
	case "json_array", "json_object", "json_merge_preserve", "json_merge_patch":
		return &ResolvedType{BaseType: BaseTypeJSON}

	// Misc
	case "coalesce", "ifnull":
		return nil // type depends on arguments
	case "uuid":
		return &ResolvedType{BaseType: BaseTypeVarchar}
	case "version":
		return &ResolvedType{BaseType: BaseTypeVarchar}
	case "database", "schema", "user", "current_user", "session_user", "system_user":
		return &ResolvedType{BaseType: BaseTypeVarchar}
	case "last_insert_id", "row_count", "found_rows":
		return &ResolvedType{BaseType: BaseTypeBigInt}
	case "connection_id":
		return &ResolvedType{BaseType: BaseTypeBigInt}
	case "charset", "collation":
		return &ResolvedType{BaseType: BaseTypeVarchar}
	case "length", "char_length", "character_length", "bit_length", "octet_length":
		return &ResolvedType{BaseType: BaseTypeInt}
	}
	return nil
}
