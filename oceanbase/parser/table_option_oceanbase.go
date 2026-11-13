package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// OceanBase's own table options: the alternatives of table_option that MySQL
// does not have, each a word and a value in the shapes alter_system.go's
// matchers already read. The option's value is kept as it was written.
//
// Ref: table_option and the rules under it in sql_parser_mysql_mode.y
// (oceanbase v4.3.5_CE_BP6, line 7836 on). The rule's rows in
// testdata/compat.tsv hold each alternative to the pinned server.
var (
	sysBool = sysAlt("TRUE", "FALSE") // BOOL_VALUE
	sysExpr = sysLeaf(func(p *Parser) bool {
		expr, err := p.parseExpr()
		return err == nil && expr != nil
	})
	sysOptComma = sysOpt(sysComma)
	// compression_name and binary_format: NAME_OB | STRING_VALUE | unreserved_keyword.
	sysNameOrStr = sysAlt(sysName, sysStr)
	// file_size_const.
	sysFileSize = sysAlt(sysInt, sysStr)

	// external_file_format, the entries of FORMAT = (...).
	sysExternalFormat = sysAlt(
		sysSeq(sysAlt("TYPE", "ENCODING", "COLUMN_INDEX_TYPE", "FILE_EXTENSION"), sysTok('='), sysStr),
		sysSeq(sysAlt("FIELD_DELIMITER", "LINE_DELIMITER", "ESCAPE", "FIELD_OPTIONALLY_ENCLOSED_BY", "FIELD_ENCLOSED_BY"), sysTok('='), sysExpr),
		sysSeq(sysAlt("SKIP_HEADER", "ROW_INDEX_STRIDE"), sysTok('='), sysInt),
		sysSeq(sysAlt("SKIP_BLANK_LINES", "TRIM_SPACE", "EMPTY_FIELD_AS_NULL", "PARSE_HEADER", "IGNORE_LAST_EMPTY_COLUMN"), sysTok('='), sysBool),
		sysSeq("NULL_IF", sysTok('='), sysTok('('), sysList(sysExpr, sysComma), sysTok(')')),
		sysSeq(sysAlt("ROW_GROUP_SIZE", "STRIPE_SIZE", "COMPRESSION_BLOCK_SIZE"), sysTok('='), sysFileSize),
		sysSeq(sysAlt("COMPRESSION", "BINARY_FORMAT"), sysTok('='), sysNameOrStr),
		sysSeq("COLUMN_BLOOM_FILTER", sysTok('='), sysTok('('), sysList(sysInt, sysOptComma), sysTok(')')),
	)
	sysFormatExpr = sysSeq("FORMAT", sysEq, sysTok('('), sysList(sysExternalFormat, sysOptComma), sysTok(')')) // format_expr

	// external_properties, the entries of PROPERTIES = (...).
	sysExternalProperty = sysAlt(
		sysSeq(sysAlt("TYPE", "ACCESSTYPE", "ACCESSID", "ACCESSKEY", "STSTOKEN", "ENDPOINT", "TUNNEL_ENDPOINT", "PROJECT_NAME",
			"SCHEMA_NAME", "QUOTA_NAME", "TABLE_NAME", "REGION", "API_MODE", "SPLIT"), sysTok('='), sysStr),
		sysSeq(sysAlt("COMPRESSION_CODE", "COMPRESSION"), sysTok('='), sysNameOrStr),
		sysSeq("COLLECT_STATISTICS_ON_CREATE", sysTok('='), sysBool),
	)
	sysExternalProperties = sysList(sysExternalProperty, sysOptComma)

	// ttl_expr: column_definition_ref '+' INTERVAL INTNUM ttl_unit. The column is
	// kept as it is read, so the catalog needs no second reading of the option.
	sysTTLExpr = sysSeq(sysTTLColumn, sysTok('+'), "INTERVAL", sysInt,
		sysAlt("SECOND", "MINUTE", "HOUR", "DAY", "MONTH", "YEAR"))
	// storage_cache_time_policy_attribute. GRANULARITY is in the rule and is no
	// keyword of the lexer's.
	sysCacheTimePolicy = sysAlt(
		sysSeq("BOUNDARY_COLUMN", sysTok('='), sysName),
		sysSeq("BOUNDARY_COLUMN_UNIT", sysTok('='), sysStr),
		sysSeq("HOT_RETENTION", sysTok('='), sysInt, sysAlt("YEAR", "MONTH", "WEEK", "DAY", "HOUR", "MINUTE", "SECOND")),
	)
)

// oceanBaseTableOptions is the alternatives of table_option that
// parseTableOption does not read itself.
var oceanBaseTableOptions = sysAlt(
	sysForm("TABLE_MODE", sysEq, sysStr),
	sysForm("DUPLICATE_SCOPE", sysEq, sysStr),
	sysForm("EXPIRE_INFO", sysEq, sysTok('('), sysExpr, sysTok(')')),
	sysForm("PROGRESSIVE_MERGE_NUM", sysEq, sysInt),
	sysForm("TABLE_ID", sysEq, sysInt),
	sysForm("STORAGE_FORMAT_VERSION", sysEq, sysInt),
	sysForm("TABLEGROUP", sysEq, sysNS),
	sysForm("READ ONLY"),
	sysForm("READ WRITE"),
	sysForm("MAX_USED_PART_ID", sysEq, sysInt),
	sysForm("PARALLEL", sysEq, sysInt),
	sysForm("NOPARALLEL"),
	sysForm("AUTO_INCREMENT_CACHE_SIZE", sysEq, sysInt),
	sysForm("ENABLE_EXTENDED_ROWID", sysEq, sysBool),
	sysForm("PARTITION_TYPE", sysEq, "USER_SPECIFIED"),
	sysForm("LOCATION", sysEq, sysStr),
	sysForm("FORMAT", sysEq, sysTok('('), sysList(sysExternalFormat, sysOptComma), sysTok(')')),
	sysForm("PROPERTIES", sysEq, sysTok('('), sysExternalProperties, sysTok(')')),
	sysForm("PATTERN", sysEq, sysStr),
	sysForm("TTL", sysEq, sysTok('('), sysList(sysTTLExpr, sysComma), sysTok(')')),
	sysForm("KV_ATTRIBUTES", sysEq, sysStr),
	sysForm("DEFAULT_LOB_INROW_THRESHOLD", sysEq, sysInt),
	sysForm("LOB_INROW_THRESHOLD", sysEq, sysInt),
	// lob_storage_clause.
	sysForm("JSON", sysTok('('), sysName, sysTok(')'), "STORE AS", sysTok('('), sysList(sysSeq("CHUNK", sysAlt(sysInt, sysStr)), sysSeq()), sysTok(')')),
	sysForm("MICRO_INDEX_CLUSTERED", sysEq, sysBool),
	sysForm("AUTO_REFRESH", sysEq, sysAlt("OFF", "IMMEDIATE", "INTERVAL")),
	sysForm("DUPLICATE_READ_CONSISTENCY", sysEq, sysStr),
	sysForm("STORAGE_CACHE_POLICY", sysEq, sysTok('('), sysAlt(sysSeq("GLOBAL", sysTok('='), sysStr), sysList(sysCacheTimePolicy, sysComma)), sysTok(')')),
	sysForm("SEMISTRUCT_ENCODING_TYPE", sysEq, sysStr),
	sysForm("MERGE_ENGINE", sysEq, sysAlt("PARTIAL_UPDATE", "DELETE_INSERT")),
)

// parseOceanBaseTableOption reads one of oceanBaseTableOptions at the current
// token, and reports whether one stands there. An option whose word is there
// and whose value is not what the grammar takes is a syntax error at the
// furthest token the option refused.
func (p *Parser) parseOceanBaseTableOption(start int) (*nodes.TableOption, bool, error) {
	p.sys = sysMatch{}
	if !oceanBaseTableOptions(p, func() bool { return true }) {
		if p.sys.action != "" {
			return nil, false, p.syntaxErrorAtTok(p.sys.far)
		}
		return nil, false, nil
	}
	value := strings.TrimSpace(p.inputText(p.sys.nameEnd, p.prev.End))
	value = strings.TrimSpace(strings.TrimPrefix(value, "="))
	return &nodes.TableOption{
		Loc:        nodes.Loc{Start: start, End: p.prev.End},
		Name:       p.sys.action,
		Value:      value,
		TTLColumns: p.sys.ttlColumns,
	}, true, nil
}

// sysTTLColumn reads a TTL entry's column_definition_ref, which may be written
// qualified, and keeps the column itself: the last name of the reference. It gives
// the name back when what follows does not match, as sysLeaf gives back the tokens.
var sysTTLColumn sysPat = func(p *Parser, k func() bool) bool {
	s := p.saveScan()
	kept := len(p.sys.ttlColumns)
	// column_definition_ref takes at most three parts, db.table.column.
	name, _, err := p.parseIdent()
	for parts := 1; err == nil && parts < 3 && p.cur.Type == '.'; parts++ {
		p.advance()
		name, _, err = p.parseIdent()
	}
	if err != nil {
		p.restoreScan(s)
		p.sysMiss()
		return false
	}
	p.sys.ttlColumns = append(p.sys.ttlColumns, name)
	if k() {
		return true
	}
	p.sys.ttlColumns = p.sys.ttlColumns[:kept]
	p.restoreScan(s)
	return false
}
