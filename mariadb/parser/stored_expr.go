package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// What MariaDB refuses of an expression stored with a table, while it parses. The action that
// closes a column's CHECK, DEFAULT or generated expression walks it and answers 1901,
// ER_GENERATED_COLUMN_FUNCTION_IS_NOT_ALLOWED, for what it may not hold, and CHECK answers
// 1970, ER_SUBQUERIES_NOT_SUPPORTED, for a subquery or NEXT VALUE FOR. A table's CHECK is held
// only to the second of those here; the rest of it waits for the table to open. Read from
// MariaDB 11.8.8 by trying every name in information_schema.SQL_FUNCTIONS and every keyword of
// this lexer in each place.
//
// Ref: Virtual_column_info::fix_and_check_expr and check_vcol_func_processor in MariaDB 11.8
// sql/field.cc and sql/item.cc.

// storedExprKind says which stored expression the parser is inside.
type storedExprKind int

const (
	storedExprNone storedExprKind = iota
	storedExprColumnCheck
	storedExprTableCheck
	storedExprDefault
	storedExprGenerated
)

// unlistedNativeFunctionNames are native too, and information_schema.SQL_FUNCTIONS, which the list
// below follows, does not name them: the spatial functions and the dynamic-column ones.
// mariadb:11.8.8 answers each with no 1305, and takes CREATE TABLE t (a POINT, x DOUBLE DEFAULT
// x(a)) and CREATE TABLE t (d BLOB DEFAULT (COLUMN_CREATE('a', 1))).
var unlistedNativeFunctionNames = func() map[string]bool {
	set := make(map[string]bool)
	for _, name := range strings.Fields(`
AREA ASBINARY ASTEXT ASWKB ASWKT BOUNDARY BUFFER CENTROID CONVEXHULL CROSSES DIMENSION
DISJOINT ENDPOINT ENVELOPE EQUALS EXTERIORRING GEOMCOLLFROMTEXT GEOMCOLLFROMWKB
GEOMETRYCOLLECTIONFROMTEXT GEOMETRYCOLLECTIONFROMWKB GEOMETRYFROMTEXT GEOMETRYFROMWKB
GEOMETRYN GEOMETRYTYPE GEOMFROMTEXT GEOMFROMWKB GLENGTH INTERIORRINGN INTERSECTS ISCLOSED
ISEMPTY ISRING ISSIMPLE LINEFROMTEXT LINEFROMWKB LINESTRINGFROMTEXT LINESTRINGFROMWKB
MBRCONTAINS MBRDISJOINT MBREQUAL MBREQUALS MBRINTERSECTS MBROVERLAPS MBRTOUCHES MBRWITHIN
MLINEFROMTEXT MLINEFROMWKB MPOINTFROMTEXT MPOINTFROMWKB MPOLYFROMTEXT MPOLYFROMWKB
MULTILINESTRINGFROMTEXT MULTILINESTRINGFROMWKB MULTIPOINTFROMTEXT MULTIPOINTFROMWKB
MULTIPOLYGONFROMTEXT MULTIPOLYGONFROMWKB NUMGEOMETRIES NUMINTERIORRINGS NUMPOINTS
POINTFROMTEXT POINTFROMWKB POINTN POINTONSURFACE POLYFROMTEXT POLYFROMWKB POLYGONFROMTEXT
POLYGONFROMWKB STARTPOINT ST_AREA ST_ASBINARY ST_ASGEOJSON ST_ASTEXT ST_ASWKB ST_ASWKT
ST_BOUNDARY ST_BUFFER ST_CENTROID ST_CONTAINS ST_CONVEXHULL ST_CROSSES ST_DIFFERENCE
ST_DIMENSION ST_DISJOINT ST_DISTANCE ST_DISTANCE_SPHERE ST_ENDPOINT ST_ENVELOPE ST_EQUALS
ST_EXTERIORRING ST_GEOMCOLLFROMTEXT ST_GEOMCOLLFROMWKB ST_GEOMETRYCOLLECTIONFROMTEXT
ST_GEOMETRYCOLLECTIONFROMWKB ST_GEOMETRYFROMTEXT ST_GEOMETRYFROMWKB ST_GEOMETRYN
ST_GEOMETRYTYPE ST_GEOMFROMGEOJSON ST_GEOMFROMTEXT ST_GEOMFROMWKB ST_INTERIORRINGN
ST_INTERSECTION ST_INTERSECTS ST_ISCLOSED ST_ISEMPTY ST_ISRING ST_ISSIMPLE ST_LENGTH
ST_LINEFROMTEXT ST_LINEFROMWKB ST_LINESTRINGFROMTEXT ST_LINESTRINGFROMWKB ST_MLINEFROMTEXT
ST_MLINEFROMWKB ST_MPOINTFROMTEXT ST_MPOINTFROMWKB ST_MPOLYFROMTEXT ST_MPOLYFROMWKB
ST_MULTILINESTRINGFROMTEXT ST_MULTILINESTRINGFROMWKB ST_MULTIPOINTFROMTEXT
ST_MULTIPOINTFROMWKB ST_MULTIPOLYGONFROMTEXT ST_MULTIPOLYGONFROMWKB ST_NUMGEOMETRIES
ST_NUMINTERIORRINGS ST_NUMPOINTS ST_OVERLAPS ST_POINTFROMTEXT ST_POINTFROMWKB ST_POINTN
ST_POINTONSURFACE ST_POLYFROMTEXT ST_POLYFROMWKB ST_POLYGONFROMTEXT ST_POLYGONFROMWKB
ST_RELATE ST_SRID ST_STARTPOINT ST_SYMDIFFERENCE ST_TOUCHES ST_UNION ST_WITHIN ST_X ST_Y
TOUCHES WITHIN X Y
COLUMN_ADD COLUMN_CHECK COLUMN_CREATE COLUMN_DELETE COLUMN_EXISTS COLUMN_GET COLUMN_JSON COLUMN_LIST
`) {
		set[name] = true
	}
	return set
}()

// builtinFunctionNames holds the names that call no stored function.
var builtinFunctionNames = map[string]bool{
	"ABS": true, "ACOS": true, "ADD_MONTHS": true, "ADDDATE": true, "ADDTIME": true,
	"AES_DECRYPT": true, "AES_ENCRYPT": true, "ASCII": true, "ASIN": true, "ATAN": true,
	"ATAN2": true, "AVG": true, "BENCHMARK": true, "BIN": true, "BINLOG_GTID_POS": true,
	"BIT_AND": true, "BIT_COUNT": true, "BIT_LENGTH": true, "BIT_OR": true, "BIT_XOR": true,
	"CAST": true, "CEIL": true, "CEILING": true, "CHAR": true, "CONVERT": true, "GET_FORMAT": true, "INSERT": true, "INTERVAL": true, "TIMESTAMPADD": true, "TIMESTAMPDIFF": true, "CHAR_LENGTH": true,
	"CHARACTER": true, "CHARACTER_LENGTH": true, "CHARSET": true, "CHR": true, "COALESCE": true,
	"COERCIBILITY": true, "COLLATION": true, "COLUMN_CHECK": true, "COLUMN_EXISTS": true,
	"COLUMN_JSON": true, "COLUMN_LIST": true, "COMPRESS": true, "CONCAT": true,
	"CONCAT_OPERATOR_ORACLE": true, "CONCAT_WS": true, "CONNECTION_ID": true, "CONTAINS": true,
	"CONV": true, "CONVERT_TZ": true, "COS": true, "COT": true, "COUNT": true, "CRC32": true,
	"CRC32C": true, "CUME_DIST": true, "CURDATE": true, "CURRENT_DATE": true,
	"CURRENT_TIME": true, "CURRENT_TIMESTAMP": true, "CURRENT_USER": true, "CURTIME": true,
	"DATABASE": true, "DATE": true, "DATE_ADD": true, "DATE_FORMAT": true, "DATE_SUB": true,
	"DATEDIFF": true, "DAY": true, "DAYNAME": true, "DAYOFMONTH": true, "DAYOFWEEK": true,
	"DAYOFYEAR": true, "DECODE": true, "DECODE_HISTOGRAM": true, "DECODE_ORACLE": true,
	"DEGREES": true, "DENSE_RANK": true, "DES_DECRYPT": true, "DES_ENCRYPT": true, "ELT": true,
	"ENCODE": true, "ENCRYPT": true, "EXP": true, "EXPORT_SET": true, "EXTRACT": true,
	"EXTRACTVALUE": true, "FIELD": true, "FIND_IN_SET": true, "FIRST_VALUE": true, "FLOOR": true,
	"FORMAT": true, "FORMAT_BYTES": true, "FORMAT_PICO_TIME": true, "FOUND_ROWS": true,
	"FROM_BASE64": true, "FROM_DAYS": true, "FROM_UNIXTIME": true, "GEOMETRYCOLLECTION": true,
	"GET_LOCK": true, "GREATEST": true, "GROUP_CONCAT": true, "HEX": true, "HOUR": true,
	"IF": true, "IFNULL": true, "INET6_ATON": true, "INET6_NTOA": true, "INET_ATON": true,
	"INET_NTOA": true, "INSTR": true, "IS_FREE_LOCK": true, "IS_IPV4": true,
	"IS_IPV4_COMPAT": true, "IS_IPV4_MAPPED": true, "IS_IPV6": true, "IS_USED_LOCK": true,
	"ISNULL": true, "JSON_ARRAY": true, "JSON_ARRAY_APPEND": true, "JSON_ARRAY_INSERT": true,
	"JSON_ARRAY_INTERSECT": true, "JSON_ARRAYAGG": true, "JSON_COMPACT": true,
	"JSON_CONTAINS": true, "JSON_CONTAINS_PATH": true, "JSON_DEPTH": true, "JSON_DETAILED": true,
	"JSON_EQUALS": true, "JSON_EXISTS": true, "JSON_EXTRACT": true, "JSON_INSERT": true,
	"JSON_KEY_VALUE": true, "JSON_KEYS": true, "JSON_LENGTH": true, "JSON_LOOSE": true,
	"JSON_MERGE": true, "JSON_MERGE_PATCH": true, "JSON_MERGE_PRESERVE": true,
	"JSON_NORMALIZE": true, "JSON_OBJECT": true, "JSON_OBJECT_FILTER_KEYS": true,
	"JSON_OBJECT_TO_ARRAY": true, "JSON_OBJECTAGG": true, "JSON_OVERLAPS": true,
	"JSON_PRETTY": true, "JSON_QUERY": true, "JSON_QUOTE": true, "JSON_REMOVE": true,
	"JSON_REPLACE": true, "JSON_SCHEMA_VALID": true, "JSON_SEARCH": true, "JSON_SET": true,
	"JSON_TYPE": true, "JSON_UNQUOTE": true, "JSON_VALID": true, "JSON_VALUE": true, "KDF": true,
	"LAG": true, "LAST_DAY": true, "LAST_INSERT_ID": true, "LAST_VALUE": true, "LCASE": true,
	"LEAD": true, "LEAST": true, "LEFT": true, "LENGTH": true, "LENGTHB": true,
	"LINESTRING": true, "LN": true, "LOAD_FILE": true, "LOCALTIME": true, "LOCALTIMESTAMP": true,
	"LOCATE": true, "LOG": true, "LOG10": true, "LOG2": true, "LOWER": true, "LPAD": true,
	"LPAD_ORACLE": true, "LTRIM": true, "LTRIM_ORACLE": true, "MAKE_SET": true, "MAKEDATE": true,
	"MAKETIME": true, "MASTER_GTID_WAIT": true, "MASTER_POS_WAIT": true, "MAX": true, "MD5": true,
	"MEDIAN": true, "MICROSECOND": true, "MID": true, "MIN": true, "MINUTE": true, "MOD": true,
	"MONTH": true, "MONTHNAME": true, "MULTILINESTRING": true, "MULTIPOINT": true,
	"MULTIPOLYGON": true, "NAME_CONST": true, "NATURAL_SORT_KEY": true, "NOW": true,
	"NTH_VALUE": true, "NTILE": true, "NULLIF": true, "NVL": true, "NVL2": true, "OCT": true,
	"OCTET_LENGTH": true, "OLD_PASSWORD": true, "ORD": true, "OVERLAPS": true, "PASSWORD": true,
	"PERCENT_RANK": true, "PERCENTILE_CONT": true, "PERCENTILE_DISC": true, "PERIOD_ADD": true,
	"PERIOD_DIFF": true, "PI": true, "POINT": true, "POLYGON": true, "POSITION": true,
	"POW": true, "POWER": true, "QUARTER": true, "QUOTE": true, "RADIANS": true, "RAND": true,
	"RANDOM_BYTES": true, "RANK": true, "REGEXP_INSTR": true, "REGEXP_REPLACE": true,
	"REGEXP_SUBSTR": true, "RELEASE_ALL_LOCKS": true, "RELEASE_LOCK": true, "REPEAT": true,
	"REPLACE": true, "REPLACE_ORACLE": true, "REVERSE": true, "RIGHT": true, "ROUND": true,
	"ROW_COUNT": true, "RPAD": true, "RPAD_ORACLE": true, "RTRIM": true, "RTRIM_ORACLE": true,
	"SCHEMA": true, "SCHEMAS": true, "SEC_TO_TIME": true, "SECOND": true, "SESSION_USER": true,
	"SFORMAT": true, "SHA": true, "SHA1": true, "SHA2": true, "SIGN": true, "SIN": true,
	"SLEEP": true, "SOUNDEX": true, "SPACE": true, "SQL_TSI_DAY": true, "SQL_TSI_HOUR": true,
	"SQL_TSI_MINUTE": true, "SQL_TSI_MONTH": true, "SQL_TSI_SECOND": true, "SQL_TSI_YEAR": true,
	"SQRT": true, "SRID": true, "STD": true, "STDDEV": true, "STDDEV_POP": true,
	"STDDEV_SAMP": true, "STR_TO_DATE": true, "STRCMP": true, "SUBDATE": true, "SUBSTR": true,
	"SUBSTR_ORACLE": true, "SUBSTRING": true, "SUBSTRING_INDEX": true, "SUBTIME": true,
	"SUM": true, "SYS_GUID": true, "SYSDATE": true, "SYSTEM_USER": true, "TAN": true,
	"TIME": true, "TIME_FORMAT": true, "TIME_TO_SEC": true, "TIMEDIFF": true, "TIMESTAMP": true,
	"TO_BASE64": true, "TO_CHAR": true, "TO_DAYS": true, "TO_SECONDS": true, "TRIM": true,
	"TRIM_ORACLE": true, "TRUNCATE": true, "UCASE": true, "UNCOMPRESS": true,
	"UNCOMPRESSED_LENGTH": true, "UNHEX": true, "UNIX_TIMESTAMP": true, "UPDATEXML": true,
	"UPPER": true, "USER": true, "UTC_DATE": true, "UTC_TIME": true, "UTC_TIMESTAMP": true,
	"UUID": true, "UUID_SHORT": true, "UUID_V4": true, "UUID_V7": true, "VAR_POP": true,
	"VAR_SAMP": true, "VARIANCE": true, "VEC_DISTANCE": true, "VEC_DISTANCE_COSINE": true,
	"VEC_DISTANCE_EUCLIDEAN": true, "VEC_FROMTEXT": true, "VEC_TOTEXT": true, "VERSION": true,
	"WEEK": true, "WEEKDAY": true, "WEEKOFYEAR": true, "WEIGHT_STRING": true,
	"WSREP_LAST_SEEN_GTID": true, "WSREP_LAST_WRITTEN_GTID": true,
	"WSREP_SYNC_WAIT_UPTO_GTID": true, "YEAR": true, "YEARWEEK": true,
}

// storedExprRefusedFunctions may stand in no column's stored expression.
var storedExprRefusedFunctions = map[string]bool{
	"AVG": true, "BENCHMARK": true, "BIT_AND": true, "BIT_OR": true, "BIT_XOR": true,
	"COUNT": true, "CUME_DIST": true, "DENSE_RANK": true, "FIRST_VALUE": true, "FOUND_ROWS": true,
	"GET_LOCK": true, "GROUP_CONCAT": true, "IS_FREE_LOCK": true, "IS_USED_LOCK": true,
	"JSON_ARRAYAGG": true, "JSON_OBJECTAGG": true, "LAG": true, "LAST_INSERT_ID": true,
	"LAST_VALUE": true, "LEAD": true, "LOAD_FILE": true, "MASTER_GTID_WAIT": true,
	"MASTER_POS_WAIT": true, "MAX": true, "MEDIAN": true, "MIN": true, "NAME_CONST": true,
	"NTH_VALUE": true, "NTILE": true, "PERCENT_RANK": true, "PERCENTILE_CONT": true,
	"PERCENTILE_DISC": true, "RANK": true, "RELEASE_ALL_LOCKS": true, "RELEASE_LOCK": true,
	"ROW_COUNT": true, "ROW_NUMBER": true, "SLEEP": true, "STD": true, "STDDEV": true,
	"STDDEV_POP": true, "STDDEV_SAMP": true, "SUM": true, "VAR_POP": true, "VAR_SAMP": true,
	"VARIANCE": true,
}

// columnCheckRefusedFunctions may stand in a DEFAULT or a generated column and not in a
// column's CHECK, which has to give the same answer every time.
var columnCheckRefusedFunctions = map[string]bool{
	"AES_DECRYPT": true, "AES_ENCRYPT": true, "CONNECTION_ID": true, "CURDATE": true,
	"CURRENT_DATE": true, "CURRENT_TIME": true, "CURRENT_TIMESTAMP": true, "CURRENT_USER": true,
	"CURTIME": true, "DATABASE": true, "DAYNAME": true,
	"FROM_UNIXTIME": true, "KDF": true, "LOCALTIME": true, "LOCALTIMESTAMP": true,
	"MONTHNAME": true, "NATURAL_SORT_KEY": true, "NOW": true, "RAND": true, "RANDOM_BYTES": true,
	"SCHEMA": true, "SESSION_USER": true, "SYS_GUID": true, "SYSDATE": true, "SYSTEM_USER": true,
	"TO_CHAR": true, "USER": true, "UTC_DATE": true, "UTC_TIME": true,
	"UTC_TIMESTAMP": true, "UUID": true, "UUID_SHORT": true, "UUID_V4": true, "UUID_V7": true,
	"VERSION": true,
}

// columnCheckRefusedArgCounts are the functions a column's CHECK refuses with one count of
// arguments only, the one that reads a session default: WEEK without a mode, DATE_FORMAT without
// a locale, UNIX_TIMESTAMP without a time. mariadb:11.8.8 takes WEEK(d,3), DATE_FORMAT(d,'%a',
// 'en_US') and UNIX_TIMESTAMP(d) there, and answers 1901 for the others while parsing.
var columnCheckRefusedArgCounts = map[string]int{"WEEK": 1, "DATE_FORMAT": 2, "UNIX_TIMESTAMP": 0}

// refuseStoredExprArgCount is the answer to one of those once its arguments are read.
func (p *Parser) refuseStoredExprArgCount(fc *nodes.FuncCallExpr, qualified bool) error {
	if p.storedExpr != storedExprColumnCheck || qualified {
		return nil
	}
	if n, ok := columnCheckRefusedArgCounts[strings.ToUpper(fc.Name)]; ok && len(fc.Args) == n {
		return &ParseError{Message: "Function or expression '" + fc.Name + "' cannot be used here", Position: fc.Loc.Start}
	}
	return nil
}

// enterStoredExpr marks the parser inside a stored expression and returns what restores it.
func (p *Parser) enterStoredExpr(kind storedExprKind) func() {
	outer := p.storedExpr
	p.storedExpr = kind
	return func() { p.storedExpr = outer }
}

// refuseStoredExprSubquery is the answer to a subquery or a sequence in a stored expression. A
// DEFAULT may hold a sequence.
func (p *Parser) refuseStoredExprSubquery(pos int, sequence bool) error {
	switch p.storedExpr {
	case storedExprNone:
		return nil
	case storedExprColumnCheck, storedExprTableCheck:
		return &ParseError{Message: "CHECK does not support subqueries or stored functions", Position: pos}
	case storedExprDefault:
		if sequence {
			return nil
		}
	}
	return &ParseError{Message: "Function or expression is not allowed in a stored expression", Position: pos}
}

// refuseStoredExprFunction is the answer to a call in a stored expression. A qualified name or
// one no built-in has is a stored function, which only a table's CHECK reads past.
func (p *Parser) refuseStoredExprFunction(name string, qualified bool, pos int) error {
	if !qualified && (name == "NEXTVAL" || name == "LASTVAL" || name == "SETVAL") {
		return p.refuseStoredExprSubquery(pos, true)
	}
	if p.storedExpr == storedExprNone || p.storedExpr == storedExprTableCheck {
		return nil
	}
	refused := qualified || (!builtinFunctionNames[name] && !unlistedNativeFunctionNames[name]) || storedExprRefusedFunctions[name]
	if p.storedExpr == storedExprColumnCheck {
		refused = refused || columnCheckRefusedFunctions[name]
	}
	if refused {
		return &ParseError{Message: "Function or expression '" + name + "' cannot be used here", Position: pos}
	}
	return nil
}

// refuseStoredExprClock is the answer to a word that reads the clock or the session with no
// parentheses, which only a column's CHECK turns away.
func (p *Parser) refuseStoredExprClock(pos int) error {
	if p.storedExpr == storedExprColumnCheck {
		return &ParseError{Message: "Function or expression cannot be used in the CHECK clause", Position: pos}
	}
	return nil
}

// refuseStoredExprVariable is the answer to a user or system variable, which only a column's
// CHECK turns away.
func (p *Parser) refuseStoredExprVariable(pos int) error {
	if p.storedExpr == storedExprColumnCheck {
		return &ParseError{Message: "Function or expression '@' cannot be used in the CHECK clause", Position: pos}
	}
	return nil
}
