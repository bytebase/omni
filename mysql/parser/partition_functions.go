package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// What MySQL 8.0 lets a partition clause call. The action that closes a partition expression or
// a partition's values walks it and answers 1064 for a function it can't partition by, and a
// name no built-in has is a stored function, which it never can. Read from the server by trying
// each built-in's name in both places, and a name that is not listed counts as refused, which
// is also what a built-in the sweep did not reach would answer after parsing.
//
// Ref: check_partition_func_processor in MySQL 8.0 sql/item_func.h and sql/partition_info.cc.

// partitionExprKind says which part of a partition clause the parser is inside.
type partitionExprKind int

const (
	partitionExprNone partitionExprKind = iota
	partitionExprMethod
	partitionExprValue
)

// partitionMethodFunctions may stand in PARTITION BY HASH, RANGE or LIST.
var partitionMethodFunctions = map[string]bool{
	"ABS": true, "ACOS": true, "ADDDATE": true, "ADDTIME": true, "ASIN": true, "ATAN": true,
	"ATAN2": true, "BIN": true, "BIT_AND": true, "BIT_COUNT": true, "BIT_LENGTH": true,
	"BIT_OR": true, "BIT_XOR": true, "CAST": true, "CEIL": true, "CEILING": true, "CHAR": true,
	"CHAR_LENGTH": true, "CHARACTER_LENGTH": true, "COALESCE": true, "COERCIBILITY": true,
	"COLLATION": true, "COMPRESS": true, "CONCAT": true, "CONCAT_WS": true, "CONV": true,
	"CONVERT": true, "CONVERT_TZ": true, "COS": true, "COT": true, "COUNT": true, "CRC32": true,
	"DATE": true, "DATE_ADD": true, "DATE_FORMAT": true, "DATEDIFF": true, "DAY": true,
	"DAYNAME": true, "DAYOFMONTH": true, "DAYOFWEEK": true, "DAYOFYEAR": true, "DEGREES": true,
	"ELT": true, "EXP": true, "EXPORT_SET": true, "EXTRACT": true, "EXTRACTVALUE": true,
	"FIELD": true, "FIND_IN_SET": true, "FLOOR": true, "FORMAT": true, "FORMAT_BYTES": true,
	"FORMAT_PICO_TIME": true, "FROM_BASE64": true, "FROM_DAYS": true, "FROM_UNIXTIME": true,
	"GREATEST": true, "GROUP_CONCAT": true, "HEX": true, "HOUR": true, "IF": true, "IFNULL": true,
	"INET6_ATON": true, "INET6_NTOA": true, "INET_ATON": true, "INET_NTOA": true, "INSERT": true,
	"INSTR": true, "IS_IPV4": true, "IS_IPV4_COMPAT": true, "IS_IPV4_MAPPED": true,
	"IS_IPV6": true, "ISNULL": true, "JSON_ARRAY": true, "JSON_ARRAY_APPEND": true,
	"JSON_ARRAY_INSERT": true, "JSON_ARRAYAGG": true, "JSON_CONTAINS": true,
	"JSON_CONTAINS_PATH": true, "JSON_DEPTH": true, "JSON_EXTRACT": true, "JSON_INSERT": true,
	"JSON_KEYS": true, "JSON_LENGTH": true, "JSON_MERGE": true, "JSON_MERGE_PATCH": true,
	"JSON_MERGE_PRESERVE": true, "JSON_OBJECT": true, "JSON_OBJECTAGG": true,
	"JSON_OVERLAPS": true, "JSON_PRETTY": true, "JSON_QUOTE": true, "JSON_REMOVE": true,
	"JSON_REPLACE": true, "JSON_SCHEMA_VALID": true, "JSON_SEARCH": true, "JSON_SET": true,
	"JSON_TYPE": true, "JSON_UNQUOTE": true, "JSON_VALID": true, "LAST_DAY": true, "LCASE": true,
	"LEAST": true, "LEFT": true, "LENGTH": true, "LN": true, "LOCATE": true, "LOG": true,
	"LOG10": true, "LOG2": true, "LOWER": true, "LPAD": true, "LTRIM": true, "MAKE_SET": true,
	"MAKEDATE": true, "MAKETIME": true, "MAX": true, "MD5": true, "MICROSECOND": true,
	"MID": true, "MIN": true, "MINUTE": true, "MOD": true, "MONTH": true, "MONTHNAME": true,
	"NAME_CONST": true, "NULLIF": true, "OCT": true, "OCTET_LENGTH": true, "ORD": true,
	"PERIOD_ADD": true, "PERIOD_DIFF": true, "PI": true, "POW": true, "POWER": true,
	"QUARTER": true, "QUOTE": true, "RADIANS": true, "REGEXP_INSTR": true, "REGEXP_REPLACE": true,
	"REGEXP_SUBSTR": true, "REPEAT": true, "REPLACE": true, "REVERSE": true, "RIGHT": true,
	"ROUND": true, "RPAD": true, "RTRIM": true, "SEC_TO_TIME": true, "SECOND": true, "SHA": true,
	"SHA1": true, "SHA2": true, "SIGN": true, "SIN": true, "SOUNDEX": true, "SPACE": true,
	"SQRT": true, "STD": true, "STDDEV": true, "STDDEV_POP": true, "STDDEV_SAMP": true,
	"STR_TO_DATE": true, "STRCMP": true, "SUBDATE": true, "SUBSTR": true, "SUBSTRING": true,
	"SUBSTRING_INDEX": true, "SUBTIME": true, "SUM": true, "TAN": true, "TIME": true,
	"TIME_FORMAT": true, "TIME_TO_SEC": true, "TIMEDIFF": true, "TIMESTAMP": true,
	"TIMESTAMPDIFF": true, "TO_BASE64": true, "TO_DAYS": true, "TO_SECONDS": true, "TRIM": true,
	"TRUNCATE": true, "UCASE": true, "UNCOMPRESS": true, "UNCOMPRESSED_LENGTH": true,
	"UNHEX": true, "UNIX_TIMESTAMP": true, "UPDATEXML": true, "UPPER": true, "VAR_POP": true,
	"VAR_SAMP": true, "VARIANCE": true, "VERSION": true, "WEEK": true, "WEEKDAY": true,
	"WEEKOFYEAR": true, "YEAR": true, "YEARWEEK": true,
}

// partitionValueFunctions may stand in VALUES LESS THAN or VALUES IN.
var partitionValueFunctions = map[string]bool{
	"ABS": true, "ACOS": true, "ADDDATE": true, "ADDTIME": true, "ASIN": true, "ATAN": true,
	"ATAN2": true, "BIN": true, "BIT_AND": true, "BIT_COUNT": true, "BIT_LENGTH": true,
	"BIT_OR": true, "BIT_XOR": true, "CAST": true, "CEIL": true, "CEILING": true, "CHAR": true,
	"CHAR_LENGTH": true, "CHARACTER_LENGTH": true, "COALESCE": true, "COERCIBILITY": true,
	"COLLATION": true, "COMPRESS": true, "CONCAT": true, "CONCAT_WS": true, "CONV": true,
	"CONVERT": true, "CONVERT_TZ": true, "COS": true, "COT": true, "COUNT": true, "CRC32": true,
	"DATE": true, "DATE_ADD": true, "DATE_FORMAT": true, "DATEDIFF": true, "DAY": true,
	"DAYNAME": true, "DAYOFMONTH": true, "DAYOFWEEK": true, "DAYOFYEAR": true, "DEGREES": true,
	"ELT": true, "EXP": true, "EXPORT_SET": true, "EXTRACT": true, "EXTRACTVALUE": true,
	"FIELD": true, "FIND_IN_SET": true, "FLOOR": true, "FORMAT": true, "FORMAT_BYTES": true,
	"FORMAT_PICO_TIME": true, "FROM_BASE64": true, "FROM_DAYS": true, "FROM_UNIXTIME": true,
	"GREATEST": true, "GROUP_CONCAT": true, "HEX": true, "HOUR": true, "IF": true, "IFNULL": true,
	"INET6_ATON": true, "INET6_NTOA": true, "INET_ATON": true, "INET_NTOA": true, "INSERT": true,
	"INSTR": true, "IS_IPV4": true, "IS_IPV4_COMPAT": true, "IS_IPV4_MAPPED": true,
	"IS_IPV6": true, "ISNULL": true, "JSON_ARRAY": true, "JSON_ARRAY_APPEND": true,
	"JSON_ARRAY_INSERT": true, "JSON_ARRAYAGG": true, "JSON_CONTAINS": true,
	"JSON_CONTAINS_PATH": true, "JSON_DEPTH": true, "JSON_EXTRACT": true, "JSON_INSERT": true,
	"JSON_KEYS": true, "JSON_LENGTH": true, "JSON_MERGE": true, "JSON_MERGE_PATCH": true,
	"JSON_MERGE_PRESERVE": true, "JSON_OBJECT": true, "JSON_OBJECTAGG": true,
	"JSON_OVERLAPS": true, "JSON_PRETTY": true, "JSON_QUOTE": true, "JSON_REMOVE": true,
	"JSON_REPLACE": true, "JSON_SCHEMA_VALID": true, "JSON_SEARCH": true, "JSON_SET": true,
	"JSON_TYPE": true, "JSON_UNQUOTE": true, "JSON_VALID": true, "LAST_DAY": true, "LCASE": true,
	"LEAST": true, "LEFT": true, "LENGTH": true, "LN": true, "LOCATE": true, "LOG": true,
	"LOG10": true, "LOG2": true, "LOWER": true, "LPAD": true, "LTRIM": true, "MAKE_SET": true,
	"MAKEDATE": true, "MAKETIME": true, "MAX": true, "MD5": true, "MICROSECOND": true,
	"MID": true, "MIN": true, "MINUTE": true, "MOD": true, "MONTH": true, "MONTHNAME": true,
	"NAME_CONST": true, "NULLIF": true, "OCT": true, "OCTET_LENGTH": true, "ORD": true,
	"PERIOD_ADD": true, "PERIOD_DIFF": true, "PI": true, "POW": true, "POWER": true,
	"QUARTER": true, "QUOTE": true, "RADIANS": true, "REGEXP_INSTR": true, "REGEXP_REPLACE": true,
	"REGEXP_SUBSTR": true, "REPEAT": true, "REPLACE": true, "REVERSE": true, "RIGHT": true,
	"ROUND": true, "RPAD": true, "RTRIM": true, "SEC_TO_TIME": true, "SECOND": true, "SHA": true,
	"SHA1": true, "SHA2": true, "SIGN": true, "SIN": true, "SOUNDEX": true, "SPACE": true,
	"SQRT": true, "STD": true, "STDDEV": true, "STDDEV_POP": true, "STDDEV_SAMP": true,
	"STR_TO_DATE": true, "STRCMP": true, "SUBDATE": true, "SUBSTR": true, "SUBSTRING": true,
	"SUBSTRING_INDEX": true, "SUBTIME": true, "SUM": true, "TAN": true, "TIME": true,
	"TIME_FORMAT": true, "TIME_TO_SEC": true, "TIMEDIFF": true, "TIMESTAMP": true,
	"TIMESTAMPDIFF": true, "TO_BASE64": true, "TO_DAYS": true, "TO_SECONDS": true, "TRIM": true,
	"TRUNCATE": true, "UCASE": true, "UNCOMPRESS": true, "UNCOMPRESSED_LENGTH": true,
	"UNHEX": true, "UNIX_TIMESTAMP": true, "UPDATEXML": true, "UPPER": true, "VAR_POP": true,
	"VAR_SAMP": true, "VARIANCE": true, "VERSION": true, "WEEK": true, "WEEKDAY": true,
	"WEEKOFYEAR": true, "YEAR": true, "YEARWEEK": true,
}

// parsePartitionExpr reads an expression of a partition clause.
func (p *Parser) parsePartitionExpr(kind partitionExprKind) (nodes.ExprNode, error) {
	outer := p.partitionExpr
	p.partitionExpr = kind
	defer func() { p.partitionExpr = outer }()
	return p.parseBitExpr()
}

// refusePartitionVariable is the answer to a variable in a partition clause: a 1064 from the
// action, like a function it can't partition by.
func (p *Parser) refusePartitionVariable(pos int) error {
	if p.partitionExpr != partitionExprNone {
		return &ParseError{Message: "This partition function is not allowed", Position: pos}
	}
	return nil
}

// refusePartitionFunction holds a call to the list for where it is.
func (p *Parser) refusePartitionFunction(name string, qualified bool, pos int) error {
	allowed := true
	switch p.partitionExpr {
	case partitionExprNone:
		return nil
	case partitionExprMethod:
		allowed = !qualified && partitionMethodFunctions[name]
	case partitionExprValue:
		allowed = !qualified && partitionValueFunctions[name]
	}
	// UNIX_TIMESTAMP reads its argument and may be used, and reads the clock with none.
	if name == "UNIX_TIMESTAMP" && p.peekNext().Type == ')' {
		allowed = false
	}
	if !allowed {
		return &ParseError{Message: "Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed", Position: pos}
	}
	return nil
}
