package catalog

// The names PL/SQL resolves without a declaration: package STANDARD's
// values, exceptions and procedures, and the built-in types. A name the
// checker cannot find among a unit's declarations, the schema's objects and
// these is one Oracle reports as PLS-00201. Each list errs on the side of a
// name: one missing here makes the catalog call a valid unit invalid.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/predefined-exceptions.html
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/plsql-data-types.html

// plsqlValues are the names STANDARD declares that an expression reads
// without arguments.
var plsqlValues = map[string]bool{
	"TRUE": true, "FALSE": true, "NULL": true,
	"SYSDATE": true, "SYSTIMESTAMP": true, "CURRENT_DATE": true, "CURRENT_TIMESTAMP": true,
	"LOCALTIMESTAMP": true, "DBTIMEZONE": true, "SESSIONTIMEZONE": true,
	"USER": true, "UID": true, "SQLCODE": true, "SQLERRM": true,
	"SYS_GUID": true, "ORA_INVOKING_USER": true, "ORA_INVOKING_USERID": true,
	// The floating-point constants.
	"BINARY_FLOAT_NAN": true, "BINARY_FLOAT_INFINITY": true,
	"BINARY_FLOAT_MAX_NORMAL": true, "BINARY_FLOAT_MIN_NORMAL": true,
	"BINARY_FLOAT_MAX_SUBNORMAL": true, "BINARY_FLOAT_MIN_SUBNORMAL": true,
	"BINARY_DOUBLE_NAN": true, "BINARY_DOUBLE_INFINITY": true,
	"BINARY_DOUBLE_MAX_NORMAL": true, "BINARY_DOUBLE_MIN_NORMAL": true,
	"BINARY_DOUBLE_MAX_SUBNORMAL": true, "BINARY_DOUBLE_MIN_SUBNORMAL": true,
}

// plsqlExceptions are STANDARD's predefined exceptions.
var plsqlExceptions = map[string]bool{
	"ACCESS_INTO_NULL": true, "CASE_NOT_FOUND": true, "COLLECTION_IS_NULL": true,
	"CURSOR_ALREADY_OPEN": true, "DUP_VAL_ON_INDEX": true, "INVALID_CURSOR": true,
	"INVALID_NUMBER": true, "LOGIN_DENIED": true, "NO_DATA_FOUND": true,
	"NO_DATA_NEEDED": true, "NOT_LOGGED_ON": true, "PROGRAM_ERROR": true,
	"ROWTYPE_MISMATCH": true, "SELF_IS_NULL": true, "STORAGE_ERROR": true,
	"SUBSCRIPT_BEYOND_COUNT": true, "SUBSCRIPT_OUTSIDE_LIMIT": true,
	"SYS_INVALID_ROWID": true, "TIMEOUT_ON_RESOURCE": true, "TOO_MANY_ROWS": true,
	"VALUE_ERROR": true, "ZERO_DIVIDE": true, "OTHERS": true,
}

// plsqlTypes are the built-in PL/SQL and SQL type names, by their first
// word.
var plsqlTypes = map[string]bool{
	"NUMBER": true, "NUMERIC": true, "DECIMAL": true, "DEC": true, "INTEGER": true, "INT": true,
	"SMALLINT": true, "FLOAT": true, "REAL": true, "DOUBLE": true, "BINARY_INTEGER": true,
	"PLS_INTEGER": true, "NATURAL": true, "NATURALN": true, "POSITIVE": true, "POSITIVEN": true,
	"SIGNTYPE": true, "SIMPLE_INTEGER": true, "BINARY_FLOAT": true, "BINARY_DOUBLE": true,
	"SIMPLE_FLOAT": true, "SIMPLE_DOUBLE": true,
	"VARCHAR2": true, "VARCHAR": true, "CHAR": true, "CHARACTER": true, "NCHAR": true,
	"NVARCHAR2": true, "STRING": true, "LONG": true, "RAW": true, "ROWID": true, "UROWID": true,
	"NATIONAL": true, "MLSLABEL": true,
	"DATE": true, "TIMESTAMP": true, "INTERVAL": true, "TIME": true,
	"BOOLEAN": true, "BOOL": true, "CLOB": true, "NCLOB": true, "BLOB": true, "BFILE": true,
	"XMLTYPE": true, "SYS_REFCURSOR": true, "JSON": true, "VECTOR": true, "ANYDATA": true,
	"ANYTYPE": true, "ANYDATASET": true, "URITYPE": true, "SDO_GEOMETRY": true,
	"EXCEPTION": true, "RECORD": true,
	"DSINTERVAL_UNCONSTRAINED": true, "YMINTERVAL_UNCONSTRAINED": true,
	"TIMESTAMP_UNCONSTRAINED": true, "TIMESTAMP_TZ_UNCONSTRAINED": true,
	"TIMESTAMP_LTZ_UNCONSTRAINED": true, "TIME_UNCONSTRAINED": true, "TIME_TZ_UNCONSTRAINED": true,
}

// plsqlFunctions are the built-in functions a PL/SQL expression may call
// without a declaration: SQL's single-row, aggregate and analytic
// functions, and STANDARD's own. builtinArity adds those the query analyzer
// checks the arguments of.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Functions.html
var plsqlFunctions = func() map[string]bool {
	m := map[string]bool{}
	for _, name := range []string{
		// Numeric.
		"ABS", "ACOS", "ASIN", "ATAN", "ATAN2", "BITAND", "CEIL", "COS", "COSH", "EXP", "FLOOR",
		"LN", "LOG", "MOD", "NANVL", "POWER", "REMAINDER", "ROUND", "SIGN", "SIN", "SINH", "SQRT",
		"TAN", "TANH", "TRUNC", "WIDTH_BUCKET", "CEILING",
		// Character.
		"CHR", "CONCAT", "INITCAP", "LOWER", "LPAD", "LTRIM", "NCHR", "NLS_INITCAP", "NLS_LOWER",
		"NLS_UPPER", "NLSSORT", "REGEXP_REPLACE", "REGEXP_SUBSTR", "REGEXP_INSTR", "REGEXP_COUNT",
		"REGEXP_LIKE", "REPLACE", "RPAD", "RTRIM", "SOUNDEX", "SUBSTR", "SUBSTRB", "SUBSTRC",
		"SUBSTR2", "SUBSTR4", "TRANSLATE", "TRIM", "UPPER", "ASCII", "INSTR", "INSTRB", "INSTRC",
		"INSTR2", "INSTR4", "LENGTH", "LENGTHB", "LENGTHC", "LENGTH2", "LENGTH4", "NLS_CHARSET_DECL_LEN",
		"NLS_CHARSET_ID", "NLS_CHARSET_NAME", "ASCIISTR", "UNISTR", "COLLATION",
		// Datetime.
		"ADD_MONTHS", "CURRENT_DATE", "CURRENT_TIMESTAMP", "DBTIMEZONE", "EXTRACT", "FROM_TZ",
		"LAST_DAY", "LOCALTIMESTAMP", "MONTHS_BETWEEN", "NEW_TIME", "NEXT_DAY", "NUMTODSINTERVAL",
		"NUMTOYMINTERVAL", "ORA_DST_AFFECTED", "ORA_DST_CONVERT", "ORA_DST_ERROR", "ROUND",
		"SESSIONTIMEZONE", "SYS_EXTRACT_UTC", "SYSDATE", "SYSTIMESTAMP", "TO_CHAR", "TO_DSINTERVAL",
		"TO_TIMESTAMP", "TO_TIMESTAMP_TZ", "TO_YMINTERVAL", "TZ_OFFSET",
		// Comparison and NULL.
		"GREATEST", "LEAST", "COALESCE", "LNNVL", "NULLIF", "NVL", "NVL2", "DECODE",
		// Conversion.
		"ASCIISTR", "BIN_TO_NUM", "CAST", "CHARTOROWID", "COMPOSE", "CONVERT", "DECOMPOSE",
		"HEXTORAW", "NUMTODSINTERVAL", "RAWTOHEX", "RAWTONHEX", "ROWIDTOCHAR", "ROWIDTONCHAR",
		"SCN_TO_TIMESTAMP", "TIMESTAMP_TO_SCN", "TO_BINARY_DOUBLE", "TO_BINARY_FLOAT", "TO_BLOB",
		"TO_CLOB", "TO_DATE", "TO_LOB", "TO_MULTI_BYTE", "TO_NCHAR", "TO_NCLOB", "TO_NUMBER",
		"TO_SINGLE_BYTE", "TREAT", "VALIDATE_CONVERSION", "TO_BOOLEAN",
		// Large objects, collections, hierarchy, XML and JSON.
		"BFILENAME", "EMPTY_BLOB", "EMPTY_CLOB", "CARDINALITY", "COLLECT", "POWERMULTISET",
		"POWERMULTISET_BY_CARDINALITY", "SET", "SYS_CONNECT_BY_PATH", "XMLAGG", "XMLCAST",
		"XMLCDATA", "XMLCOLATTVAL", "XMLCOMMENT", "XMLCONCAT", "XMLDIFF", "XMLELEMENT", "XMLEXISTS",
		"XMLFOREST", "XMLISVALID", "XMLPARSE", "XMLPATCH", "XMLPI", "XMLQUERY", "XMLROOT",
		"XMLSEQUENCE", "XMLSERIALIZE", "XMLTABLE", "XMLTRANSFORM", "EXTRACTVALUE", "EXISTSNODE",
		"JSON_ARRAY", "JSON_ARRAYAGG", "JSON_OBJECT", "JSON_OBJECTAGG", "JSON_QUERY", "JSON_SCALAR",
		"JSON_SERIALIZE", "JSON_TABLE", "JSON_TRANSFORM", "JSON_VALUE", "JSON_EXISTS", "JSON_MERGEPATCH",
		// Environment and identifiers.
		"SYS_CONTEXT", "SYS_GUID", "SYS_TYPEID", "UID", "USER", "USERENV", "ORA_HASH",
		"STANDARD_HASH", "DUMP", "VSIZE", "NLS_CHARSET_ID", "ORA_INVOKING_USER", "SQLCODE", "SQLERRM",
		"CON_ID_TO_CON_NAME", "CON_NAME_TO_ID", "CON_DBID_TO_ID", "CON_GUID_TO_ID", "COST",
		// Aggregate and analytic.
		"APPROX_COUNT_DISTINCT", "AVG", "CORR", "COUNT", "COVAR_POP", "COVAR_SAMP", "CUME_DIST",
		"DENSE_RANK", "FIRST", "FIRST_VALUE", "GROUP_ID", "GROUPING", "GROUPING_ID", "LAG", "LAST",
		"LAST_VALUE", "LEAD", "LISTAGG", "MAX", "MEDIAN", "MIN", "NTH_VALUE", "NTILE", "PERCENT_RANK",
		"PERCENTILE_CONT", "PERCENTILE_DISC", "RANK", "RATIO_TO_REPORT", "ROW_NUMBER", "STATS_MODE",
		"STDDEV", "STDDEV_POP", "STDDEV_SAMP", "SUM", "VAR_POP", "VAR_SAMP", "VARIANCE", "ANY_VALUE",
		// PL/SQL STANDARD.
		"RAISE_APPLICATION_ERROR",
	} {
		m[name] = true
	}
	for name := range builtinArity {
		m[name] = true
	}
	return m
}()
