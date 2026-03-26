package completion

import (
	"github.com/bytebase/omni/mysql/catalog"
	"github.com/bytebase/omni/mysql/parser"
)

// Built-in MySQL function names for func_name / function_ref resolution.
var builtinFunctions = []string{
	// Aggregate
	"COUNT", "SUM", "AVG", "MAX", "MIN", "GROUP_CONCAT", "JSON_ARRAYAGG", "JSON_OBJECTAGG",
	"BIT_AND", "BIT_OR", "BIT_XOR", "STD", "STDDEV", "STDDEV_POP", "STDDEV_SAMP",
	"VAR_POP", "VAR_SAMP", "VARIANCE",
	// Window
	"ROW_NUMBER", "RANK", "DENSE_RANK", "CUME_DIST", "PERCENT_RANK",
	"NTILE", "LAG", "LEAD", "FIRST_VALUE", "LAST_VALUE", "NTH_VALUE",
	// String
	"CONCAT", "CONCAT_WS", "SUBSTRING", "SUBSTR", "LEFT", "RIGHT", "LENGTH", "CHAR_LENGTH",
	"CHARACTER_LENGTH", "UPPER", "LOWER", "LCASE", "UCASE", "TRIM", "LTRIM", "RTRIM",
	"REPLACE", "REVERSE", "INSERT", "LPAD", "RPAD", "REPEAT", "SPACE", "FORMAT",
	"LOCATE", "INSTR", "POSITION", "FIELD", "FIND_IN_SET", "ELT", "MAKE_SET",
	"QUOTE", "SOUNDEX", "HEX", "UNHEX", "ORD", "ASCII", "BIN", "OCT",
	// Numeric
	"ABS", "CEIL", "CEILING", "FLOOR", "ROUND", "TRUNCATE", "MOD", "POW", "POWER",
	"SQRT", "EXP", "LOG", "LOG2", "LOG10", "LN", "SIGN", "PI", "RAND",
	"CRC32", "CONV", "RADIANS", "DEGREES", "SIN", "COS", "TAN", "ASIN", "ACOS", "ATAN", "COT",
	// Date/Time
	"NOW", "CURDATE", "CURTIME", "CURRENT_DATE", "CURRENT_TIME", "CURRENT_TIMESTAMP",
	"SYSDATE", "UTC_DATE", "UTC_TIME", "UTC_TIMESTAMP", "LOCALTIME", "LOCALTIMESTAMP",
	"DATE", "TIME", "YEAR", "MONTH", "DAY", "HOUR", "MINUTE", "SECOND", "MICROSECOND",
	"DAYNAME", "DAYOFMONTH", "DAYOFWEEK", "DAYOFYEAR", "WEEK", "WEEKDAY", "WEEKOFYEAR",
	"QUARTER", "YEARWEEK", "LAST_DAY", "MAKEDATE", "MAKETIME",
	"DATE_ADD", "DATE_SUB", "ADDDATE", "SUBDATE", "ADDTIME", "SUBTIME",
	"DATEDIFF", "TIMEDIFF", "TIMESTAMPDIFF", "TIMESTAMPADD",
	"DATE_FORMAT", "TIME_FORMAT", "STR_TO_DATE", "FROM_UNIXTIME", "UNIX_TIMESTAMP",
	"EXTRACT", "GET_FORMAT", "PERIOD_ADD", "PERIOD_DIFF", "SEC_TO_TIME", "TIME_TO_SEC",
	"FROM_DAYS", "TO_DAYS", "TO_SECONDS",
	// Control flow
	"IF", "IFNULL", "NULLIF", "COALESCE", "GREATEST", "LEAST", "INTERVAL",
	// Cast
	"CAST", "CONVERT", "BINARY",
	// JSON
	"JSON_ARRAY", "JSON_OBJECT", "JSON_QUOTE", "JSON_EXTRACT", "JSON_UNQUOTE",
	"JSON_CONTAINS", "JSON_CONTAINS_PATH", "JSON_KEYS", "JSON_SEARCH",
	"JSON_SET", "JSON_INSERT", "JSON_REPLACE", "JSON_REMOVE", "JSON_MERGE_PRESERVE",
	"JSON_MERGE_PATCH", "JSON_DEPTH", "JSON_LENGTH", "JSON_TYPE", "JSON_VALID",
	"JSON_ARRAYAGG", "JSON_OBJECTAGG", "JSON_PRETTY", "JSON_STORAGE_FREE", "JSON_STORAGE_SIZE",
	"JSON_TABLE", "JSON_VALUE",
	// Info
	"DATABASE", "SCHEMA", "USER", "CURRENT_USER", "SESSION_USER", "SYSTEM_USER",
	"VERSION", "CONNECTION_ID", "LAST_INSERT_ID", "ROW_COUNT", "FOUND_ROWS",
	"BENCHMARK", "CHARSET", "COLLATION", "COERCIBILITY",
	// Encryption
	"MD5", "SHA1", "SHA2", "AES_ENCRYPT", "AES_DECRYPT",
	"RANDOM_BYTES",
	// Misc
	"UUID", "UUID_SHORT", "UUID_TO_BIN", "BIN_TO_UUID",
	"SLEEP", "VALUES", "DEFAULT", "INET_ATON", "INET_NTOA", "INET6_ATON", "INET6_NTOA",
	"IS_IPV4", "IS_IPV6", "IS_UUID",
	"ANY_VALUE", "GROUPING",
}

// Known charsets for "charset" rule.
var knownCharsets = []string{
	"utf8mb4", "utf8mb3", "utf8", "latin1", "ascii", "binary",
	"big5", "cp1250", "cp1251", "cp1256", "cp1257", "cp850", "cp852", "cp866",
	"dec8", "eucjpms", "euckr", "gb2312", "gb18030", "gbk", "geostd8",
	"greek", "hebrew", "hp8", "keybcs2", "koi8r", "koi8u",
	"latin2", "latin5", "latin7", "macce", "macroman",
	"sjis", "swe7", "tis620", "ucs2", "ujis", "utf16", "utf16le", "utf32",
}

// Known engines for "engine" rule.
var knownEngines = []string{
	"InnoDB", "MyISAM", "MEMORY", "CSV", "ARCHIVE",
	"BLACKHOLE", "MERGE", "FEDERATED", "NDB", "NDBCLUSTER",
}

// MySQL type keywords for "type_name" rule.
var typeKeywords = []string{
	// Integer types
	"TINYINT", "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "BIGINT",
	// Fixed-point
	"DECIMAL", "NUMERIC", "DEC", "FIXED",
	// Floating-point
	"FLOAT", "DOUBLE", "REAL",
	// Bit
	"BIT", "BOOL", "BOOLEAN",
	// String types
	"CHAR", "VARCHAR", "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT",
	"BINARY", "VARBINARY", "TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB",
	// Date/time
	"DATE", "DATETIME", "TIMESTAMP", "TIME", "YEAR",
	// JSON
	"JSON",
	// Spatial
	"GEOMETRY", "POINT", "LINESTRING", "POLYGON",
	"MULTIPOINT", "MULTILINESTRING", "MULTIPOLYGON", "GEOMETRYCOLLECTION",
	// Enum/Set
	"ENUM", "SET",
	// Serial
	"SERIAL",
}

// resolveRules converts parser rule candidates into typed Candidate values
// using the catalog for name resolution.
func resolveRules(cs *parser.CandidateSet, cat *catalog.Catalog, sql string, cursorOffset int) []Candidate {
	if cs == nil {
		return nil
	}
	var result []Candidate
	for _, rc := range cs.Rules {
		result = append(result, resolveRule(rc.Rule, cat, sql, cursorOffset)...)
	}
	return result
}

// resolveRule resolves a single grammar rule name into completion candidates.
func resolveRule(rule string, cat *catalog.Catalog, sql string, cursorOffset int) []Candidate {
	switch rule {
	case "table_ref":
		return resolveTableRef(cat)
	case "columnref":
		return resolveColumnRefScoped(cat, sql, cursorOffset)
	case "database_ref":
		return resolveDatabaseRef(cat)
	case "function_ref", "func_name":
		return resolveFunctionRef(cat)
	case "procedure_ref":
		return resolveProcedureRef(cat)
	case "index_ref":
		return resolveIndexRef(cat)
	case "trigger_ref":
		return resolveTriggerRef(cat)
	case "event_ref":
		return resolveEventRef(cat)
	case "view_ref":
		return resolveViewRef(cat)
	case "charset":
		return resolveCharset()
	case "engine":
		return resolveEngine()
	case "type_name":
		return resolveTypeName()
	}
	return nil
}

// resolveTableRef returns tables and views from the current database.
func resolveTableRef(cat *catalog.Catalog) []Candidate {
	if cat == nil {
		return nil
	}
	db := currentDB(cat)
	if db == nil {
		return nil
	}
	var result []Candidate
	for _, t := range db.Tables {
		result = append(result, Candidate{Text: t.Name, Type: CandidateTable})
	}
	for _, v := range db.Views {
		result = append(result, Candidate{Text: v.Name, Type: CandidateView})
	}
	return result
}

// resolveColumnRefScoped returns columns scoped to the tables referenced in
// the SQL statement. If no table refs are found, falls back to all columns
// in the current database.
func resolveColumnRefScoped(cat *catalog.Catalog, sql string, cursorOffset int) []Candidate {
	if cat == nil {
		return nil
	}
	db := currentDB(cat)
	if db == nil {
		return nil
	}

	refs := extractTableRefs(sql, cursorOffset)
	if len(refs) == 0 {
		return resolveColumnRef(cat)
	}

	seen := make(map[string]bool)
	var result []Candidate
	for _, ref := range refs {
		// Resolve table in the appropriate database.
		targetDB := db
		if ref.Database != "" {
			targetDB = cat.GetDatabase(ref.Database)
			if targetDB == nil {
				continue
			}
		}
		// Look up the table.
		for _, t := range targetDB.Tables {
			if t.Name == ref.Table {
				for _, col := range t.Columns {
					if !seen[col.Name] {
						seen[col.Name] = true
						result = append(result, Candidate{Text: col.Name, Type: CandidateColumn})
					}
				}
				break
			}
		}
		// Also check views.
		for _, v := range targetDB.Views {
			if v.Name == ref.Table {
				for _, colName := range v.Columns {
					if !seen[colName] {
						seen[colName] = true
						result = append(result, Candidate{Text: colName, Type: CandidateColumn})
					}
				}
				break
			}
		}
	}
	// If we found refs but couldn't resolve any columns (e.g. CTE name not
	// matching any catalog table), fall back to all columns.
	if len(result) == 0 {
		return resolveColumnRef(cat)
	}
	return result
}

// resolveColumnRef returns columns from all tables in the current database.
func resolveColumnRef(cat *catalog.Catalog) []Candidate {
	if cat == nil {
		return nil
	}
	db := currentDB(cat)
	if db == nil {
		return nil
	}
	seen := make(map[string]bool)
	var result []Candidate
	for _, t := range db.Tables {
		for _, col := range t.Columns {
			if !seen[col.Name] {
				seen[col.Name] = true
				result = append(result, Candidate{Text: col.Name, Type: CandidateColumn})
			}
		}
	}
	for _, v := range db.Views {
		for _, colName := range v.Columns {
			if !seen[colName] {
				seen[colName] = true
				result = append(result, Candidate{Text: colName, Type: CandidateColumn})
			}
		}
	}
	return result
}

// resolveDatabaseRef returns all databases from the catalog.
func resolveDatabaseRef(cat *catalog.Catalog) []Candidate {
	if cat == nil {
		return nil
	}
	var result []Candidate
	for _, db := range cat.Databases() {
		result = append(result, Candidate{Text: db.Name, Type: CandidateDatabase})
	}
	return result
}

// resolveFunctionRef returns catalog functions + built-in function names.
func resolveFunctionRef(cat *catalog.Catalog) []Candidate {
	var result []Candidate
	// Built-in functions always available.
	for _, name := range builtinFunctions {
		result = append(result, Candidate{Text: name, Type: CandidateFunction})
	}
	// Catalog functions from current database.
	if cat != nil {
		if db := currentDB(cat); db != nil {
			for _, fn := range db.Functions {
				result = append(result, Candidate{Text: fn.Name, Type: CandidateFunction})
			}
		}
	}
	return result
}

// resolveProcedureRef returns catalog procedures from the current database.
func resolveProcedureRef(cat *catalog.Catalog) []Candidate {
	if cat == nil {
		return nil
	}
	db := currentDB(cat)
	if db == nil {
		return nil
	}
	var result []Candidate
	for _, p := range db.Procedures {
		result = append(result, Candidate{Text: p.Name, Type: CandidateProcedure})
	}
	return result
}

// resolveIndexRef returns indexes from all tables in the current database.
// Table-scoped resolution will be refined in later phases.
func resolveIndexRef(cat *catalog.Catalog) []Candidate {
	if cat == nil {
		return nil
	}
	db := currentDB(cat)
	if db == nil {
		return nil
	}
	seen := make(map[string]bool)
	var result []Candidate
	for _, t := range db.Tables {
		for _, idx := range t.Indexes {
			if idx.Name != "" && !seen[idx.Name] {
				seen[idx.Name] = true
				result = append(result, Candidate{Text: idx.Name, Type: CandidateIndex})
			}
		}
	}
	return result
}

// resolveTriggerRef returns triggers from the current database.
func resolveTriggerRef(cat *catalog.Catalog) []Candidate {
	if cat == nil {
		return nil
	}
	db := currentDB(cat)
	if db == nil {
		return nil
	}
	var result []Candidate
	for _, tr := range db.Triggers {
		result = append(result, Candidate{Text: tr.Name, Type: CandidateTrigger})
	}
	return result
}

// resolveEventRef returns events from the current database.
func resolveEventRef(cat *catalog.Catalog) []Candidate {
	if cat == nil {
		return nil
	}
	db := currentDB(cat)
	if db == nil {
		return nil
	}
	var result []Candidate
	for _, ev := range db.Events {
		result = append(result, Candidate{Text: ev.Name, Type: CandidateEvent})
	}
	return result
}

// resolveViewRef returns views from the current database.
func resolveViewRef(cat *catalog.Catalog) []Candidate {
	if cat == nil {
		return nil
	}
	db := currentDB(cat)
	if db == nil {
		return nil
	}
	var result []Candidate
	for _, v := range db.Views {
		result = append(result, Candidate{Text: v.Name, Type: CandidateView})
	}
	return result
}

// resolveCharset returns known MySQL charset names.
func resolveCharset() []Candidate {
	result := make([]Candidate, len(knownCharsets))
	for i, name := range knownCharsets {
		result[i] = Candidate{Text: name, Type: CandidateCharset}
	}
	return result
}

// resolveEngine returns known MySQL engine names.
func resolveEngine() []Candidate {
	result := make([]Candidate, len(knownEngines))
	for i, name := range knownEngines {
		result[i] = Candidate{Text: name, Type: CandidateEngine}
	}
	return result
}

// resolveTypeName returns MySQL type keywords.
func resolveTypeName() []Candidate {
	result := make([]Candidate, len(typeKeywords))
	for i, name := range typeKeywords {
		result[i] = Candidate{Text: name, Type: CandidateType_}
	}
	return result
}

// currentDB returns the current database from the catalog, or nil.
func currentDB(cat *catalog.Catalog) *catalog.Database {
	name := cat.CurrentDatabase()
	if name == "" {
		return nil
	}
	return cat.GetDatabase(name)
}
