package review

import (
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// Oracle built-in datatype classifiers used for column-level semantic checks.
var (
	oracleScalarTypes = map[string]bool{
		"NUMBER": true, "FLOAT": true, "BINARY_FLOAT": true, "BINARY_DOUBLE": true,
		"INTEGER": true, "INT": true, "SMALLINT": true, "DECIMAL": true, "DEC": true,
		"NUMERIC": true, "REAL": true, "DOUBLE": true, "DOUBLE PRECISION": true,
		"VARCHAR2": true, "VARCHAR": true, "NVARCHAR2": true, "CHAR": true,
		"CHARACTER": true, "NCHAR": true, "DATE": true, "TIMESTAMP": true,
		"INTERVAL": true, "RAW": true, "LONG": true, "LONG RAW": true,
		"ROWID": true, "UROWID": true, "BOOLEAN": true,
	}

	oracleNonCharTypes = map[string]bool{
		"NUMBER": true, "FLOAT": true, "BINARY_FLOAT": true, "BINARY_DOUBLE": true,
		"INTEGER": true, "INT": true, "SMALLINT": true, "DECIMAL": true, "DEC": true,
		"NUMERIC": true, "REAL": true, "DOUBLE": true, "DOUBLE PRECISION": true,
		"DATE": true, "TIMESTAMP": true, "INTERVAL": true, "RAW": true, "LONG RAW": true,
		"ROWID": true, "UROWID": true, "BOOLEAN": true, "BLOB": true, "BFILE": true,
		"VECTOR": true, "JSON": true,
	}

	oracleBuiltinTypeNames = map[string]bool{
		"NUMBER": true, "NUMERIC": true, "DECIMAL": true, "DEC": true,
		"INTEGER": true, "INT": true, "SMALLINT": true, "FLOAT": true,
		"DOUBLE": true, "REAL": true, "BINARY_FLOAT": true, "BINARY_DOUBLE": true,
		"CHAR": true, "CHARACTER": true, "VARCHAR": true, "VARCHAR2": true,
		"NCHAR": true, "NVARCHAR2": true, "LONG": true, "RAW": true,
		"DATE": true, "TIMESTAMP": true, "INTERVAL": true,
		"BLOB": true, "CLOB": true, "NCLOB": true, "BFILE": true,
		"ROWID": true, "UROWID": true, "BOOLEAN": true, "JSON": true,
		"VECTOR": true, "XMLTYPE": true, "URITYPE": true, "HTTPURITYPE": true,
		"DBURITYPE": true, "XDBURITYPE": true, "SDO_GEOMETRY": true,
		"ANYDATA": true, "ANYTYPE": true, "ANYDATASET": true, "MLSLABEL": true,
		"PLS_INTEGER": true, "BINARY_INTEGER": true, "NATURAL": true, "POSITIVE": true,
		"SIGNTYPE": true, "SIMPLE_INTEGER": true, "SIMPLE_FLOAT": true, "SIMPLE_DOUBLE": true,
		"STRING": true, "NATIONAL": true,
	}

	// SQL aggregate functions valid inside PIVOT (... FOR ... IN (...)).
	oracleAggregateBuiltins = map[string]bool{
		"ANY_VALUE": true, "APPROX_COUNT": true, "APPROX_COUNT_DISTINCT": true,
		"APPROX_COUNT_DISTINCT_AGG": true, "APPROX_MEDIAN": true, "APPROX_PERCENTILE": true,
		"APPROX_PERCENTILE_AGG": true, "APPROX_SUM": true, "AVG": true,
		"BIT_AND_AGG": true, "BIT_OR_AGG": true, "BIT_XOR_AGG": true,
		"BOOLEAN_AND_AGG": true, "BOOLEAN_OR_AGG": true, "CHECKSUM": true,
		"COLLECT": true, "CORR": true, "CORR_K": true, "CORR_S": true,
		"COUNT": true, "COVAR_POP": true, "COVAR_SAMP": true,
		"CUME_DIST": true, "DENSE_RANK": true, "FIRST": true, "GROUP_ID": true,
		"GROUPING": true, "GROUPING_ID": true, "JSON_ARRAYAGG": true,
		"JSON_DATAGUIDE": true, "JSON_OBJECTAGG": true, "KURTOSIS_POP": true,
		"KURTOSIS_SAMP": true, "LAST": true, "LISTAGG": true, "MAX": true,
		"MEDIAN": true, "MIN": true, "PERCENT_RANK": true, "PERCENTILE_CONT": true,
		"PERCENTILE_DISC": true, "RANK": true, "REGR_AVGX": true, "REGR_AVGY": true,
		"REGR_COUNT": true, "REGR_INTERCEPT": true, "REGR_R2": true, "REGR_SLOPE": true,
		"REGR_SXX": true, "REGR_SXY": true, "REGR_SYY": true, "SKEWNESS_POP": true,
		"SKEWNESS_SAMP": true, "STATS_BINOMIAL_TEST": true, "STATS_CROSSTAB": true,
		"STATS_F_TEST": true, "STATS_KS_TEST": true, "STATS_MODE": true,
		"STATS_MW_TEST": true, "STATS_ONE_WAY_ANOVA": true, "STATS_T_TEST_INDEP": true,
		"STATS_T_TEST_INDEPU": true, "STATS_T_TEST_ONE": true, "STATS_T_TEST_PAIRED": true,
		"STATS_WSR_TEST": true, "STDDEV": true, "STDDEV_POP": true, "STDDEV_SAMP": true,
		"SUM": true, "SYS_OP_ZONE_ID": true, "SYS_XMLAGG": true,
		"TO_APPROX_COUNT_DISTINCT": true, "TO_APPROX_PERCENTILE": true,
		"VAR_POP": true, "VAR_SAMP": true, "VARIANCE": true, "XMLAGG": true,
	}

	oraclePseudoColumns = map[string]bool{
		"ROWID": true, "ROWNUM": true, "LEVEL": true, "SYSDATE": true,
		"SYSTIMESTAMP": true, "USER": true, "UID": true, "CURRENT_DATE": true,
		"CURRENT_TIMESTAMP": true, "LOCALTIMESTAMP": true, "DBTIMEZONE": true,
		"SESSIONTIMEZONE": true, "ORA_ROWSCN": true, "COLUMN_VALUE": true,
		"OBJECT_ID": true, "OBJECT_VALUE": true, "XMLDATA": true,
		"CONNECT_BY_ISLEAF": true, "CONNECT_BY_ISCYCLE": true,
		"VERSIONS_STARTTIME": true, "VERSIONS_STARTSCN": true,
		"VERSIONS_ENDTIME": true, "VERSIONS_ENDSCN": true,
		"VERSIONS_XID": true, "VERSIONS_OPERATION": true,
	}

	// Keywords that can immediately follow PARTITION / SUBPARTITION in an
	// unnamed partition specification rather than being a partition name.
	partitionOptionKeywords = map[string]bool{
		"VALUES": true, "TABLESPACE": true, "LOGGING": true, "NOLOGGING": true,
		"COMPRESS": true, "NOCOMPRESS": true, "OVERFLOW": true, "SEGMENT": true,
		"READ": true, "INDEXING": true, "INMEMORY": true, "NO": true,
		"PCTFREE": true, "PCTUSED": true, "INITRANS": true, "MAXTRANS": true,
		"STORAGE": true, "LOB": true, "VARRAY": true, "NESTED": true,
		"SUBPARTITIONS": true, "SUBPARTITION": true, "PARTITION": true,
		"USABLE": true, "UNUSABLE": true, "PARAMETERS": true, "PARALLEL": true,
		"NOPARALLEL": true, "CACHE": true, "NOCACHE": true, "ILM": true,
		"MEMOPTIMIZE": true, "ROW": true, "COLUMN": true, "UPDATE": true,
		"ONLINE": true, "INCLUDING": true, "ALLOW": true, "DISALLOW": true,
		"COMPUTATION": true, "USAGE": true, "UNUSED": true, "EXTERNAL": true,
		"DEFAULT": true,
	}
)

func isOracleScalarType(typ string) bool {
	return oracleScalarTypes[baseOracleType(typ)]
}

func isOracleNonCharType(typ string) bool {
	return oracleNonCharTypes[baseOracleType(typ)]
}

func typeNameBase(tn *ast.TypeName) string {
	if tn == nil || tn.Names == nil || len(tn.Names.Items) == 0 {
		return ""
	}
	if len(tn.Names.Items) == 1 {
		if s, ok := tn.Names.Items[0].(*ast.String); ok {
			return strings.ToUpper(strings.TrimSpace(s.Str))
		}
	}
	// Multi-part name like schema.type_name or DOUBLE PRECISION / LONG RAW.
	var parts []string
	for _, item := range tn.Names.Items {
		if s, ok := item.(*ast.String); ok {
			parts = append(parts, strings.ToUpper(strings.TrimSpace(s.Str)))
		}
	}
	joined := strings.Join(parts, " ")
	if oracleBuiltinTypeNames[joined] || oracleScalarTypes[joined] {
		return joined
	}
	if len(parts) == 2 && (parts[0] == "SYS" || parts[0] == "MDSYS" || parts[0] == "PUBLIC") && oracleBuiltinTypeNames[parts[1]] {
		return parts[1]
	}
	if len(parts) > 0 && (parts[0] == "TIMESTAMP" || parts[0] == "INTERVAL" || parts[0] == "NATIONAL" || parts[0] == "DOUBLE" || parts[0] == "CHARACTER" || parts[0] == "LONG") {
		return parts[0]
	}
	return joined
}

func tokWord(toks []parser.Token, i int, word string) bool {
	return i >= 0 && i < len(toks) &&
		toks[i].Type != tokenQuoted && toks[i].Type != tokenString && toks[i].Type != tokenNString &&
		strings.EqualFold(toks[i].Str, word)
}

// checkAlterUpgrade rejects column_properties subclauses that Oracle 23ai
// always refuses with ORA-01735 inside ALTER TABLE ... UPGRADE [NOT] INCLUDING
// DATA (verified on Oracle 23ai Free against plain, nested-table/varray/JSON,
// and object tables): NESTED TABLE, VARRAY, JSON(...), and XMLTYPE ... STORE
// AS {BASICFILE | SECUREFILE | BINARY XML | CLOB <segname>}.
func (b *builder) checkAlterUpgrade(loc ast.Loc, rng review.Range) bool {
	toks := tokens(b.w.sql, loc)
	start := 0
	for i := 0; i < len(toks); i++ {
		if tokWord(toks, i, "UPGRADE") {
			start = i + 1
			if tokWord(toks, start, "NOT") {
				start++
			}
			if tokWord(toks, start, "INCLUDING") && tokWord(toks, start+1, "DATA") {
				start += 2
			}
			break
		}
	}
	depth := 0
	for i := start; i < len(toks); i++ {
		switch toks[i].Type {
		case tokenOpenParen:
			depth++
			continue
		case tokenCloseParen:
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth != 0 {
			continue
		}
		switch {
		case tokWord(toks, i, "NESTED") && tokWord(toks, i+1, "TABLE"):
			b.fail(rng, review.UnsupportedForm, "ALTER TABLE UPGRADE does not accept NESTED TABLE storage clauses.")
			return true
		case tokWord(toks, i, "VARRAY"):
			b.fail(rng, review.UnsupportedForm, "ALTER TABLE UPGRADE does not accept VARRAY storage clauses.")
			return true
		case tokWord(toks, i, "JSON") && i+1 < len(toks) && toks[i+1].Type == tokenOpenParen:
			b.fail(rng, review.UnsupportedForm, "ALTER TABLE UPGRADE does not accept JSON storage clauses.")
			return true
		case tokWord(toks, i, "XMLTYPE"):
			j := i + 1
			if tokWord(toks, j, "COLUMN") {
				j++
			}
			if j < len(toks) {
				j++ // column name
			}
			if tokWord(toks, j, "STORE") && tokWord(toks, j+1, "AS") {
				k := j + 2
				switch {
				case tokWord(toks, k, "BASICFILE"), tokWord(toks, k, "SECUREFILE"), tokWord(toks, k, "BINARY"):
					b.fail(rng, review.UnsupportedForm, "ALTER TABLE UPGRADE does not accept XMLTYPE STORE AS "+strings.ToUpper(toks[k].Str)+".")
					return true
				case tokWord(toks, k, "CLOB") && k+1 < len(toks):
					next := toks[k+1]
					if _, isName := tokenName(next); isName && (next.Type == tokenQuoted || (!tokWord(toks, k+1, "XMLTYPE") && !tokWord(toks, k+1, "ELEMENT") && !tokWord(toks, k+1, "ALLOW") && !tokWord(toks, k+1, "DISALLOW") && !tokWord(toks, k+1, "XMLSCHEMA") && !tokWord(toks, k+1, "LOGGING") && !tokWord(toks, k+1, "NOLOGGING") && !tokWord(toks, k+1, "CACHE") && !tokWord(toks, k+1, "NOCACHE") && !tokWord(toks, k+1, "TABLESPACE") && !tokWord(toks, k+1, "PCTFREE") && !tokWord(toks, k+1, "PARALLEL") && !tokWord(toks, k+1, "NOPARALLEL") && !tokWord(toks, k+1, "COMPRESS") && !tokWord(toks, k+1, "NOCOMPRESS"))) {
						b.fail(rng, review.UnsupportedForm, "ALTER TABLE UPGRADE does not accept a LOB segment name after XMLTYPE STORE AS CLOB.")
						return true
					}
				}
			}
		}
	}
	return false
}

// checkPartitionTextContradictions catches self-contradictory or malformed
// partition maintenance clauses that fail on every Oracle 23ai table or
// materialized view regardless of existing partition metadata.
func (b *builder) checkPartitionTextContradictions(loc ast.Loc, tableRef objRef, isAlterTable bool, rng review.Range) bool {
	toks := tokens(b.w.sql, loc)
	if len(toks) == 0 {
		return false
	}
	for i := 0; i < len(toks); i++ {
		switch {
		// MERGE PARTITIONS / SUBPARTITIONS p1, p2  or  p1 TO p2
		case tokWord(toks, i, "MERGE") && (tokWord(toks, i+1, "PARTITIONS") || tokWord(toks, i+1, "SUBPARTITIONS")):
			j := i + 2
			if j >= len(toks) {
				continue
			}
			// MERGE PARTITIONS (1), a -> ORA-14108
			if toks[j].Type == tokenOpenParen {
				b.fail(rng, review.UnsupportedForm, "MERGE PARTITIONS requires partition names or FOR (...) expressions.")
				return true
			}
			// MERGE PARTITIONS FOR(...), a -> ORA-14764 (mixing FOR and partition name, or same FOR)
			if tokWord(toks, j, "FOR") && j+1 < len(toks) && toks[j+1].Type == tokenOpenParen {
				k := skipBalancedParens(toks, j+1)
				if k < len(toks) && toks[k].Type == tokenComma && k+1 < len(toks) && !tokWord(toks, k+1, "FOR") {
					b.fail(rng, review.UnsupportedForm, "MERGE PARTITIONS cannot mix FOR (...) with a partition name.")
					return true
				}
				continue
			}
			if n1, ok := tokenName(toks[j]); ok && j+2 < len(toks) && (toks[j+1].Type == tokenComma || tokWord(toks, j+1, "TO")) {
				if n2, ok2 := tokenName(toks[j+2]); ok2 && n1 == n2 {
					b.fail(rng, review.DuplicateName, "Cannot merge partition "+review.Code(ident(n1))+" with itself.")
					return true
				}
			}

		// RENAME [SUB]PARTITION p1 TO p2 with p1 == p2
		case tokWord(toks, i, "RENAME") && (tokWord(toks, i+1, "PARTITION") || tokWord(toks, i+1, "SUBPARTITION")):
			if i+4 < len(toks) && tokWord(toks, i+3, "TO") {
				n1, ok1 := tokenName(toks[i+2])
				n2, ok2 := tokenName(toks[i+4])
				if ok1 && ok2 && n1 == n2 {
					b.fail(rng, review.DuplicateName, "Cannot rename partition "+review.Code(ident(n1))+" to the same name.")
					return true
				}
			}

		// ALTER TABLE t EXCHANGE [SUB]PARTITION ... WITH TABLE t (same table)
		case isAlterTable && tokWord(toks, i, "EXCHANGE") && (tokWord(toks, i+1, "PARTITION") || tokWord(toks, i+1, "SUBPARTITION")):
			for k := i + 2; k+2 < len(toks); k++ {
				if tokWord(toks, k, "WITH") && tokWord(toks, k+1, "TABLE") {
					schema, name, ok := readQualifiedTokenName(toks, k+2, b.w.current)
					if ok && schema == tableRef.schema && name == tableRef.name {
						b.fail(rng, review.UnsupportedForm, "Cannot exchange a partition of "+review.Code(tableRef.written)+" with itself.")
						return true
					}
					break
				}
			}

		// SPLIT ... INTO ( ... )
		case tokWord(toks, i, "INTO") && i+1 < len(toks) && toks[i+1].Type == tokenOpenParen:
			end := skipBalancedParens(toks, i+1)
			inner := toks[i+2 : end]
			if end > i+2 && toks[end-1].Type == tokenCloseParen {
				inner = toks[i+2 : end-1]
			}
			if b.checkSplitIntoParens(inner, rng) {
				return true
			}

		// ADD PARTITION ... ( SUBPARTITION SUBPARTITION )
		case tokWord(toks, i, "SUBPARTITION") && tokWord(toks, i+1, "SUBPARTITION"):
			b.fail(rng, review.UnsupportedForm, "Missing comma between subpartition specifications.")
			return true

		// DROP / TRUNCATE [SUB]PARTITION (expr) without FOR (ORA-14108),
		// or DROP [SUB]PARTITION p1, p1 duplicate name (ORA-14332).
		case (tokWord(toks, i, "DROP") || tokWord(toks, i, "TRUNCATE")) &&
			(tokWord(toks, i+1, "PARTITION") || tokWord(toks, i+1, "PARTITIONS") || tokWord(toks, i+1, "SUBPARTITION") || tokWord(toks, i+1, "SUBPARTITIONS")):
			j := i + 2
			if j < len(toks) && toks[j].Type == tokenOpenParen {
				b.fail(rng, review.UnsupportedForm, "Partition specification requires a partition name or FOR (...).")
				return true
			}
			if tokWord(toks, i, "DROP") {
				seen := make(map[string]bool)
				for k := j; k < len(toks); k++ {
					if toks[k].Type == tokenComma {
						continue
					}
					if tokWord(toks, k, "FOR") && k+1 < len(toks) && toks[k+1].Type == tokenOpenParen {
						k = skipBalancedParens(toks, k+1) - 1
						continue
					}
					if tokWord(toks, k, "UPDATE") || tokWord(toks, k, "INVALIDATE") || tokWord(toks, k, "PARALLEL") || tokWord(toks, k, "NOPARALLEL") {
						break
					}
					if name, ok := tokenName(toks[k]); ok {
						if seen[name] {
							b.fail(rng, review.DuplicateName, twice("names", "partition", name))
							return true
						}
						seen[name] = true
					} else {
						break
					}
				}
			}

		}
	}
	if b.checkUpdateIndexesText(loc, rng) {
		return true
	}
	return false
}

// checkUpdateIndexesText checks UPDATE INDEXES (...) clauses for duplicate
// index names (ORA-14293) and missing commas between partition specs.
func (b *builder) checkUpdateIndexesText(loc ast.Loc, rng review.Range) bool {
	toks := tokens(b.w.sql, loc)
	for i := 0; i+2 < len(toks); i++ {
		if !tokWord(toks, i, "UPDATE") || !tokWord(toks, i+1, "INDEXES") || toks[i+2].Type != tokenOpenParen {
			continue
		}
		end := skipBalancedParens(toks, i+2)
		inner := toks[i+3 : end]
		if end > i+3 && toks[end-1].Type == tokenCloseParen {
			inner = toks[i+3 : end-1]
		}
		seenIdx := make(map[string]bool)
		depth := 0
		for j := 0; j < len(inner); j++ {
			if inner[j].Type == tokenOpenParen {
				depth++
				continue
			}
			if inner[j].Type == tokenCloseParen {
				if depth > 0 {
					depth--
				}
				continue
			}
			if depth == 0 {
				if inner[j].Type == tokenComma {
					continue
				}
				if name, ok := tokenName(inner[j]); ok && j+1 < len(inner) && inner[j+1].Type == tokenOpenParen {
					if seenIdx[name] {
						b.fail(rng, review.DuplicateName, twice("names", "index", name))
						return true
					}
					seenIdx[name] = true
				}
			} else if (tokWord(inner, j, "PARTITION") && tokWord(inner, j+1, "PARTITION")) ||
				(tokWord(inner, j, "SUBPARTITION") && tokWord(inner, j+1, "SUBPARTITION")) {
				b.fail(rng, review.UnsupportedForm, "Missing comma between index partition specifications.")
				return true
			}
		}
	}
	return false
}

func isIOTPropertySubtype(subtype string) bool {
	switch subtype {
	case "OVERFLOW", "COALESCE", "MAPPING", "ADD OVERFLOW":
		return true
	}
	return false
}

// isExternalPropertyClause reports whether an ALTER TABLE clause alters an
// external table's properties (alter_external_table): DEFAULT DIRECTORY,
// ACCESS PARAMETERS, LOCATION, REJECT LIMIT, or PROJECT COLUMN. Any other
// table refuses them (01735).
func isExternalPropertyClause(toks []parser.Token) bool {
	return tokensSay(toks, "DEFAULT", "DIRECTORY") || tokensSay(toks, "ACCESS", "PARAMETERS") ||
		tokensSay(toks, "LOCATION") || tokensSay(toks, "REJECT", "LIMIT") || tokensSay(toks, "PROJECT", "COLUMN")
}

// tokensSay reports whether a clause's tokens start with words.
func tokensSay(toks []parser.Token, words ...string) bool {
	for i, w := range words {
		if !tokWord(toks, i, w) {
			return false
		}
	}
	return true
}

// notExternal reports whether the walk knows r is a table of the schema
// that is neither external nor hybrid partitioned: the sync lists it with
// no partitions, which leaves out only system and reference partitioning,
// neither of which an external table takes, and with an index that backs
// no constraint, which no external table has (30657).
func (b *builder) notExternal(r objRef) bool {
	if b.heapTable(r) {
		return true
	}
	if b.syncedPartMeta(r) == nil || b.syncedIndexMeta(r) == nil {
		return false
	}
	if st := b.w.snap.tables[r.name]; st == nil || st.partCount > 0 {
		return false
	}
	for name, table := range b.w.snap.indexes {
		if table == r.name && !b.w.snap.constraintIndexes[name] {
			return true
		}
	}
	return false
}

// heapTable reports whether the walk knows r is a heap table of no
// partitions: the sync gave it a tablespace (syncedTable.heap) and the
// change has not touched its partitions.
func (b *builder) heapTable(r objRef) bool {
	st := b.syncedPartMeta(r)
	return st != nil && st.heap
}

// checkIOTPartitionClause rejects IOT-only partition clauses (OVERFLOW,
// MAPPING TABLE, or bare COALESCE without SUBPARTITION) when the table is
// known to have no primary key (an IOT must have a primary key, ORA-25175).
func (b *builder) checkIOTPartitionClause(r objRef, pkAbsent bool, loc ast.Loc, rng review.Range) bool {
	if !pkAbsent {
		return false
	}
	toks := tokens(b.w.sql, loc)
	for i := 0; i < len(toks); i++ {
		if tokWord(toks, i, "OVERFLOW") {
			b.fail(rng, review.UnsupportedForm, "OVERFLOW is only allowed on an index-organized table, not "+review.Code(r.written)+".")
			return true
		}
		if tokWord(toks, i, "MAPPING") && tokWord(toks, i+1, "TABLE") {
			b.fail(rng, review.UnsupportedForm, "MAPPING TABLE is only allowed on an index-organized table, not "+review.Code(r.written)+".")
			return true
		}
		if tokWord(toks, i, "MODIFY") && tokWord(toks, i+1, "PARTITION") {
			for k := i + 2; k < len(toks); k++ {
				if tokWord(toks, k, "COALESCE") && !tokWord(toks, k+1, "SUBPARTITION") && !tokWord(toks, k+1, "PARTITION") {
					b.fail(rng, review.UnsupportedForm, "MODIFY PARTITION COALESCE is only allowed on an index-organized table, not "+review.Code(r.written)+".")
					return true
				}
			}
		}
	}
	return false
}

// checkIOTMoveClause rejects MOVE MAPPING TABLE / MOVE NOMAPPING on a table
// known to have no primary key.
func (b *builder) checkIOTMoveClause(r objRef, pkAbsent bool, loc ast.Loc, rng review.Range) bool {
	if !pkAbsent {
		return false
	}
	toks := tokens(b.w.sql, loc)
	for i := 0; i < len(toks); i++ {
		if tokWord(toks, i, "NOMAPPING") || (tokWord(toks, i, "MAPPING") && tokWord(toks, i+1, "TABLE")) {
			b.fail(rng, review.UnsupportedForm, "MAPPING TABLE / NOMAPPING is only allowed on an index-organized table, not "+review.Code(r.written)+".")
			return true
		}
	}
	return false
}

// checkSplitIntoParens checks SPLIT ... INTO (...) for missing commas
// (PARTITION PARTITION / SUBPARTITION SUBPARTITION) and duplicate output
// partition names (PARTITION a, PARTITION a).
func (b *builder) checkSplitIntoParens(inner []parser.Token, rng review.Range) bool {
	seen := make(map[string]bool)
	depth := 0
	for i := 0; i < len(inner); i++ {
		switch inner[i].Type {
		case tokenOpenParen:
			depth++
			continue
		case tokenCloseParen:
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth != 0 {
			continue
		}
		if (tokWord(inner, i, "PARTITION") && tokWord(inner, i+1, "PARTITION")) ||
			(tokWord(inner, i, "SUBPARTITION") && tokWord(inner, i+1, "SUBPARTITION")) {
			b.fail(rng, review.UnsupportedForm, "Missing comma between split partition specifications.")
			return true
		}
		if tokWord(inner, i, "PARTITION") || tokWord(inner, i, "SUBPARTITION") {
			if i+1 < len(inner) && inner[i+1].Type != tokenComma && inner[i+1].Type != tokenOpenParen && inner[i+1].Type != tokenCloseParen {
				if inner[i+1].Type == tokenQuoted || !partitionOptionKeywords[strings.ToUpper(inner[i+1].Str)] {
					if name, ok := tokenName(inner[i+1]); ok {
						if seen[name] {
							b.fail(rng, review.DuplicateName, twice("names", "partition", name))
							return true
						}
						seen[name] = true
					}
				}
			}
		}
	}
	return false
}

func skipBalancedParens(toks []parser.Token, openIdx int) int {
	depth := 0
	for i := openIdx; i < len(toks); i++ {
		if toks[i].Type == tokenOpenParen {
			depth++
		} else if toks[i].Type == tokenCloseParen {
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return len(toks)
}

func readQualifiedTokenName(toks []parser.Token, start int, defaultSchema string) (schema, name string, ok bool) {
	if start >= len(toks) {
		return "", "", false
	}
	n1, ok1 := tokenName(toks[start])
	if !ok1 {
		return "", "", false
	}
	if start+2 < len(toks) && toks[start+1].Str == "." {
		n2, ok2 := tokenName(toks[start+2])
		if !ok2 {
			return "", "", false
		}
		return n1, n2, true
	}
	return defaultSchema, n1, true
}

// checkPartitionMeta validates a partition maintenance clause (from ALTER TABLE
// or table-delegated ALTER MATERIALIZED VIEW) against a synced table's
// partition metadata (oracle/driver/sync.go).
func (b *builder) checkPartitionMeta(r objRef, st *syncedTable, loc ast.Loc, rng review.Range) bool {
	if st == nil {
		return false
	}
	toks := tokens(b.w.sql, loc)
	if len(toks) == 0 {
		return false
	}

	// Tables with no synced RANGE/LIST/HASH partitions may be non-partitioned,
	// REFERENCE-partitioned (requires a foreign key), or SYSTEM-partitioned
	// (omitted by oracle/driver/sync.go). When no foreign key exists, reject
	// only partition clauses that are invalid on BOTH non-partitioned and
	// SYSTEM-partitioned tables: subpartitions, SPLIT/MERGE/COALESCE,
	// SET INTERVAL/PARTITIONING/SUBPARTITION TEMPLATE, PARTITION FOR (...),
	// ADD/MODIFY PARTITION with VALUES, and MODIFY PARTITIONSET.
	if st.heap {
		b.fail(rng, review.UnsupportedForm, "Table "+review.Code(r.written)+" is not partitioned.")
		return true
	}
	if st.partType == metadata.TablePartitionMetadata_TYPE_UNSPECIFIED {
		if st.hasForeignKey {
			return false
		}
		for i := 0; i < len(toks); i++ {
			if tokWord(toks, i, "SUBPARTITION") || tokWord(toks, i, "SUBPARTITIONS") {
				b.fail(rng, review.UnsupportedForm, "Table "+review.Code(r.written)+" is not composite-partitioned.")
				return true
			}
			if tokWord(toks, i, "PARTITION") && tokWord(toks, i+1, "FOR") && i+2 < len(toks) && toks[i+2].Type == tokenOpenParen {
				b.fail(rng, review.UnsupportedForm, "PARTITION FOR (...) requires a range-, list-, or hash-partitioned table, not "+review.Code(r.written)+".")
				return true
			}
			if tokWord(toks, i, "PARTITIONSET") {
				b.fail(rng, review.UnsupportedForm, "Table "+review.Code(r.written)+" is not partitioned.")
				return true
			}
		}
		verbIdx := 0
		for verbIdx < len(toks) && !(tokWord(toks, verbIdx, "ADD") || tokWord(toks, verbIdx, "DROP") ||
			tokWord(toks, verbIdx, "TRUNCATE") || tokWord(toks, verbIdx, "SPLIT") ||
			tokWord(toks, verbIdx, "MERGE") || tokWord(toks, verbIdx, "COALESCE") ||
			tokWord(toks, verbIdx, "MODIFY") || tokWord(toks, verbIdx, "MOVE") ||
			tokWord(toks, verbIdx, "RENAME") || tokWord(toks, verbIdx, "EXCHANGE") ||
			tokWord(toks, verbIdx, "SET")) {
			verbIdx++
		}
		if verbIdx >= len(toks) {
			return false
		}
		switch {
		case tokWord(toks, verbIdx, "SPLIT"), tokWord(toks, verbIdx, "MERGE"), tokWord(toks, verbIdx, "COALESCE"):
			b.fail(rng, review.UnsupportedForm, "Table "+review.Code(r.written)+" is not range-, list-, or hash-partitioned.")
			return true
		case tokWord(toks, verbIdx, "SET") && (tokWord(toks, verbIdx+1, "INTERVAL") || tokWord(toks, verbIdx+1, "PARTITIONING") || tokWord(toks, verbIdx+1, "AUTOMATIC") || tokWord(toks, verbIdx+1, "SUBPARTITION")):
			b.fail(rng, review.UnsupportedForm, "Table "+review.Code(r.written)+" is not range- or list-partitioned.")
			return true
		case (tokWord(toks, verbIdx, "ADD") || tokWord(toks, verbIdx, "MODIFY")) && tokWord(toks, verbIdx+1, "PARTITION"):
			for k := verbIdx + 2; k < len(toks); k++ {
				if tokWord(toks, k, "VALUES") {
					b.fail(rng, review.UnsupportedForm, "Partition VALUES clause requires a range- or list-partitioned table, not "+review.Code(r.written)+".")
					return true
				}
			}
			if tokWord(toks, verbIdx, "ADD") && b.checkAddPartitionNonHashOpts(toks[verbIdx+2:], rng) {
				return true
			}
		}
		return false
	}

	// Subpartition operations on a non-composite partitioned table.
	if !st.hasSubparts {
		for i := 0; i < len(toks); i++ {
			if tokWord(toks, i, "SUBPARTITION") || tokWord(toks, i, "SUBPARTITIONS") {
				b.fail(rng, review.UnsupportedForm, "Table "+review.Code(r.written)+" is not composite-partitioned.")
				return true
			}
		}
	}

	// Identify the leading partition verb inside loc.
	verbIdx := 0
	for verbIdx < len(toks) && !(tokWord(toks, verbIdx, "ADD") || tokWord(toks, verbIdx, "DROP") ||
		tokWord(toks, verbIdx, "TRUNCATE") || tokWord(toks, verbIdx, "SPLIT") ||
		tokWord(toks, verbIdx, "MERGE") || tokWord(toks, verbIdx, "COALESCE") ||
		tokWord(toks, verbIdx, "MODIFY") || tokWord(toks, verbIdx, "MOVE") ||
		tokWord(toks, verbIdx, "RENAME") || tokWord(toks, verbIdx, "EXCHANGE") ||
		tokWord(toks, verbIdx, "SET")) {
		verbIdx++
	}
	if verbIdx >= len(toks) {
		return false
	}

	// On known RANGE/LIST/HASH tables, check PARTITION FOR (...) arity (ORA-14092)
	// and RANGE SPLIT PARTITION bound order (ORA-14080).
	if st.partKeyCount > 0 {
		for i := 0; i+2 < len(toks); i++ {
			if tokWord(toks, i, "PARTITION") && tokWord(toks, i+1, "FOR") && toks[i+2].Type == tokenOpenParen {
				if n := countTopLevelExprs(toks, i+2); n > 0 && n != st.partKeyCount {
					b.fail(rng, review.UnsupportedForm, "PARTITION FOR value count does not match partitioning key count of "+review.Code(r.written)+".")
					return true
				}
			}
		}
	}
	if st.partType == metadata.TablePartitionMetadata_RANGE && tokWord(toks, verbIdx, "SPLIT") && tokWord(toks, verbIdx+1, "PARTITION") && verbIdx+2 < len(toks) {
		if pname, ok := tokenName(toks[verbIdx+2]); ok {
			if b.checkRangeSplitBound(r, st, pname, toks[verbIdx+3:], rng) {
				return true
			}
		}
	}

	switch st.partType {
	case metadata.TablePartitionMetadata_RANGE:
		switch {
		case tokWord(toks, verbIdx, "COALESCE") && tokWord(toks, verbIdx+1, "PARTITION"):
			b.fail(rng, review.UnsupportedForm, "COALESCE PARTITION requires a hash-partitioned table.")
			return true
		case tokWord(toks, verbIdx, "SET") && (tokWord(toks, verbIdx+1, "PARTITIONING") || tokWord(toks, verbIdx+1, "AUTOMATIC")):
			b.fail(rng, review.UnsupportedForm, "SET PARTITIONING AUTOMATIC/MANUAL requires a list-partitioned table.")
			return true
		case tokWord(toks, verbIdx, "MODIFY") && tokWord(toks, verbIdx+1, "PARTITION"):
			for k := verbIdx + 2; k+1 < len(toks); k++ {
				if (tokWord(toks, k, "ADD") || tokWord(toks, k, "DROP")) && tokWord(toks, k+1, "VALUES") {
					b.fail(rng, review.UnsupportedForm, "MODIFY PARTITION ADD/DROP VALUES requires a list-partitioned table.")
					return true
				}
			}
		case tokWord(toks, verbIdx, "SPLIT") && tokWord(toks, verbIdx+1, "PARTITION"):
			for k := verbIdx + 2; k+1 < len(toks); k++ {
				if tokWord(toks, k, "VALUES") && toks[k+1].Type == tokenOpenParen {
					b.fail(rng, review.UnsupportedForm, "SPLIT PARTITION VALUES (...) requires a list-partitioned table.")
					return true
				}
			}
		case tokWord(toks, verbIdx, "ADD") && tokWord(toks, verbIdx+1, "PARTITION"):
			if b.checkAddPartitionNonHashOpts(toks[verbIdx+2:], rng) {
				return true
			}
			return b.checkRangeAddPartition(r, st, toks[verbIdx+2:], rng)
		}

	case metadata.TablePartitionMetadata_LIST:
		switch {
		case tokWord(toks, verbIdx, "COALESCE") && tokWord(toks, verbIdx+1, "PARTITION"):
			b.fail(rng, review.UnsupportedForm, "COALESCE PARTITION requires a hash-partitioned table.")
			return true
		case tokWord(toks, verbIdx, "SET") && tokWord(toks, verbIdx+1, "INTERVAL"):
			b.fail(rng, review.UnsupportedForm, "SET INTERVAL requires a range-partitioned table.")
			return true
		case tokWord(toks, verbIdx, "ADD") && tokWord(toks, verbIdx+1, "PARTITION"):
			if b.checkAddPartitionNonHashOpts(toks[verbIdx+2:], rng) {
				return true
			}
			for k := verbIdx + 2; k+2 < len(toks); k++ {
				if tokWord(toks, k, "VALUES") && tokWord(toks, k+1, "LESS") && tokWord(toks, k+2, "THAN") {
					b.fail(rng, review.UnsupportedForm, "ADD PARTITION VALUES LESS THAN requires a range-partitioned table.")
					return true
				}
			}
		}

	case metadata.TablePartitionMetadata_HASH:
		switch {
		case tokWord(toks, verbIdx, "SET") && (tokWord(toks, verbIdx+1, "INTERVAL") || tokWord(toks, verbIdx+1, "PARTITIONING") || tokWord(toks, verbIdx+1, "AUTOMATIC")):
			b.fail(rng, review.UnsupportedForm, "Cannot use SET INTERVAL or SET PARTITIONING on a hash-partitioned table.")
			return true
		case (tokWord(toks, verbIdx, "SPLIT") || tokWord(toks, verbIdx, "MERGE") || tokWord(toks, verbIdx, "DROP")) &&
			(tokWord(toks, verbIdx+1, "PARTITION") || tokWord(toks, verbIdx+1, "PARTITIONS")):
			b.fail(rng, review.UnsupportedForm, "Cannot SPLIT, MERGE, or DROP partitions on a hash-partitioned table.")
			return true
		case tokWord(toks, verbIdx, "ADD") && tokWord(toks, verbIdx+1, "PARTITION"):
			for k := verbIdx + 2; k < len(toks); k++ {
				if tokWord(toks, k, "VALUES") {
					b.fail(rng, review.UnsupportedForm, "ADD PARTITION on a hash-partitioned table cannot specify VALUES.")
					return true
				}
			}
		}
	}
	return false
}

// checkRangeAddPartition validates ADD PARTITION on a RANGE-partitioned table:
// requires VALUES LESS THAN (...), checks bound tuple arity (ORA-14017), and
// checks monotonically increasing integer high values (ORA-14074).
func (b *builder) checkRangeAddPartition(r objRef, st *syncedTable, rest []parser.Token, rng review.Range) bool {
	depth := 0
	vltIdx := -1
	for i := 0; i < len(rest); i++ {
		if rest[i].Type == tokenOpenParen {
			depth++
			continue
		}
		if rest[i].Type == tokenCloseParen {
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth != 0 {
			continue
		}
		if tokWord(rest, i, "VALUES") {
			if tokWord(rest, i+1, "LESS") && tokWord(rest, i+2, "THAN") {
				vltIdx = i + 3
				break
			}
			b.fail(rng, review.UnsupportedForm, "ADD PARTITION VALUES (...) requires a list-partitioned table, not range-partitioned "+review.Code(r.written)+".")
			return true
		}
	}
	if vltIdx < 0 {
		b.fail(rng, review.UnsupportedForm, "ADD PARTITION on range-partitioned table "+review.Code(r.written)+" requires VALUES LESS THAN.")
		return true
	}
	if vltIdx >= len(rest) || rest[vltIdx].Type != tokenOpenParen {
		return false
	}
	end := skipBalancedParens(rest, vltIdx)
	if end <= vltIdx+1 {
		return false
	}
	inner := rest[vltIdx+1 : end-1]
	// Count top-level comma-separated bound expressions.
	items := 1
	d := 0
	for _, t := range inner {
		if t.Type == tokenOpenParen {
			d++
		} else if t.Type == tokenCloseParen && d > 0 {
			d--
		} else if t.Type == tokenComma && d == 0 {
			items++
		}
	}
	if st.partKeyCount > 0 && items != st.partKeyCount {
		b.fail(rng, review.UnsupportedForm, "Partition bound value count does not match partitioning key count of "+review.Code(r.written)+".")
		return true
	}
	if st.allRangeInt && items == 1 && len(inner) == 1 && inner[0].Type != tokenQuoted {
		if val, err := strconv.ParseInt(inner[0].Str, 10, 64); err == nil && val <= st.maxRangeHighVal {
			b.fail(rng, review.UnsupportedForm, "Partition bound for "+review.Code(r.written)+" must be higher than the last partition's bound.")
			return true
		}
	}
	return false
}

// alterMviewRejectsSynonym reports whether an ALTER MATERIALIZED VIEW
// partition action does NOT follow synonyms in Oracle 23ai (returns ORA-00942
// on a synonym even when the synonym targets a partitioned table, verified on
// Oracle 23ai Free): MODIFY PARTITION (except ADD/COALESCE SUBPARTITION),
// MODIFY SUBPARTITION, RENAME PARTITION, RENAME SUBPARTITION, DROP PARTITION,
// TRUNCATE PARTITION.
func (b *builder) alterMviewRejectsSynonym(v *ast.AlterMaterializedViewStmt) bool {
	if v.Name == nil || v.Name.Loc.End <= 0 || v.Name.Loc.End >= v.Loc.End {
		return false
	}
	toks := tokens(b.w.sql, ast.Loc{Start: v.Name.Loc.End, End: v.Loc.End})
	if len(toks) < 2 {
		return false
	}
	switch {
	case tokWord(toks, 0, "RENAME") && (tokWord(toks, 1, "PARTITION") || tokWord(toks, 1, "SUBPARTITION")):
		return true
	case tokWord(toks, 0, "MODIFY") && tokWord(toks, 1, "SUBPARTITION"):
		return true
	case tokWord(toks, 0, "MODIFY") && tokWord(toks, 1, "PARTITION"):
		// MODIFY PARTITION p ADD SUBPARTITION / COALESCE SUBPARTITION follows synonyms;
		// other MODIFY PARTITION actions reject synonyms with ORA-00942.
		for i := 2; i+1 < len(toks); i++ {
			if (tokWord(toks, i, "ADD") || tokWord(toks, i, "COALESCE")) && tokWord(toks, i+1, "SUBPARTITION") {
				return false
			}
		}
		return true
	case (tokWord(toks, 0, "DROP") || tokWord(toks, 0, "TRUNCATE")) && (tokWord(toks, 1, "PARTITION") || tokWord(toks, 1, "SUBPARTITION")):
		return true
	}
	return false
}

// checkCreateIndexType validates INDEXTYPE IS <name> on CREATE INDEX:
//   - rejects an indextype name already occupied by a non-indextype object
//     (TABLE, VIEW, MATERIALIZED VIEW, SEQUENCE, routine, package, or synonym
//     to a table) (ORA-29833)
//   - rejects INDEXTYPE IS XDB.XMLINDEX on a known built-in scalar column
//     (ORA-29958/29855)
func (b *builder) checkCreateIndexType(v *ast.CreateIndexStmt, targetRef objRef, st *syncedTable, rng review.Range) bool {
	if v.IndexType == nil || v.IndexType.DBLink != "" {
		return false
	}
	if ir, ok := b.ref(v.IndexType); ok {
		if f, kind := b.look(ir.key(spaceObject)); f == walk.Present && kind != "" && kind != "INDEXTYPE" {
			if kind != kindSynonym {
				b.need(ir.key(spaceObject), walk.Present, []string{"INDEXTYPE"}, ir.rng, "", wrongKind(ir.written, kind, "indextype"))
				return true
			}
			if _, sst, sok := b.resolveSynonymTable(v.IndexType, ir); sok && sst != nil {
				b.fail(ir.rng, review.WrongObjectKind, wrongKind(ir.written, kindTable, "indextype"))
				return true
			}
		}
	}
	if v.IndexType.Schema == "XDB" && v.IndexType.Name == "XMLINDEX" && st != nil && v.Columns != nil && len(v.Columns.Items) == 1 {
		if ic, ok := v.Columns.Items[0].(*ast.IndexColumn); ok {
			if cr, ok := ic.Expr.(*ast.ColumnRef); ok && cr.Column != "" && cr.Schema == "" {
				if sc, ok := st.colMeta[cr.Column]; ok && isOracleScalarType(sc.typ) {
					b.fail(rng, review.UnsupportedForm, "XMLINDEX cannot be created on scalar column "+review.Code(ident(cr.Column))+" of "+review.Code(targetRef.written)+".")
					return true
				}
			}
		}
	}
	return false
}

// checkCreateIndexLocal checks CREATE INDEX ... LOCAL:
//   - pure-text errors: PARTITION PARTITION / SUBPARTITION SUBPARTITION (missing
//     comma), STORE IN (...,) trailing comma
//   - SYSTEM-partition-safe checks on tables with no synced RANGE/LIST/HASH
//     partitions (and no FK): UNIQUE LOCAL (ORA-14039), LOCAL SUBPARTITION,
//     LOCAL STORE IN (ORA-14177)
//   - non-composite RANGE/LIST tables: LOCAL SUBPARTITION and LOCAL STORE IN
//   - explicit LOCAL (PARTITION ...) count mismatch on non-interval partitioned
//     tables (ORA-14024)
func (b *builder) checkCreateIndexLocal(v *ast.CreateIndexStmt, targetRef objRef, st *syncedTable, rng review.Range) bool {
	if !v.Local {
		return false
	}
	toks := tokens(b.w.sql, v.Loc)
	hasSubpart := false
	hasStoreIn := false
	for i := 0; i < len(toks); i++ {
		if (tokWord(toks, i, "PARTITION") && tokWord(toks, i+1, "PARTITION")) ||
			(tokWord(toks, i, "SUBPARTITION") && tokWord(toks, i+1, "SUBPARTITION")) {
			b.fail(rng, review.UnsupportedForm, "Missing comma between LOCAL index partition specifications.")
			return true
		}
		if tokWord(toks, i, "SUBPARTITION") || tokWord(toks, i, "SUBPARTITIONS") {
			hasSubpart = true
		}
		if tokWord(toks, i, "STORE") && tokWord(toks, i+1, "IN") && i+2 < len(toks) && toks[i+2].Type == tokenOpenParen {
			hasStoreIn = true
			end := skipBalancedParens(toks, i+2)
			if end > i+3 && toks[end-1].Type == tokenCloseParen && toks[end-2].Type == tokenComma {
				b.fail(rng, review.UnsupportedForm, "Trailing comma in LOCAL index STORE IN clause.")
				return true
			}
		}
	}
	if st == nil {
		return false
	}
	if st.heap {
		b.fail(rng, review.UnsupportedForm, "A LOCAL index needs a partitioned table, not "+review.Code(targetRef.written)+".")
		return true
	}
	if st.partType == metadata.TablePartitionMetadata_TYPE_UNSPECIFIED {
		// SYSTEM- and REFERENCE-partitioned tables are omitted by oracle/driver/sync.go
		// and do support non-unique LOCAL indexes without subpartitions or STORE IN.
		if st.hasForeignKey {
			return false
		}
		if v.Unique {
			b.fail(rng, review.UnsupportedForm, "Cannot create a UNIQUE LOCAL index without partitioning key columns on "+review.Code(targetRef.written)+".")
			return true
		}
		if hasSubpart {
			b.fail(rng, review.UnsupportedForm, "Table "+review.Code(targetRef.written)+" is not composite-partitioned.")
			return true
		}
		if hasStoreIn {
			b.fail(rng, review.UnsupportedForm, "LOCAL index STORE IN requires a hash- or composite-hash-partitioned table, not "+review.Code(targetRef.written)+".")
			return true
		}
		return false
	}
	if !st.hasSubparts && hasSubpart {
		b.fail(rng, review.UnsupportedForm, "Table "+review.Code(targetRef.written)+" is not composite-partitioned.")
		return true
	}
	if !st.hasSubparts && st.partType != metadata.TablePartitionMetadata_HASH && hasStoreIn {
		b.fail(rng, review.UnsupportedForm, "LOCAL index STORE IN requires a hash- or composite-hash-partitioned table, not "+review.Code(targetRef.written)+".")
		return true
	}
	if st.hasInterval || st.partCount == 0 || st.partType == metadata.TablePartitionMetadata_LIST {
		return false
	}
	for i := 0; i+1 < len(toks); i++ {
		if tokWord(toks, i, "LOCAL") && toks[i+1].Type == tokenOpenParen {
			end := skipBalancedParens(toks, i+1)
			inner := toks[i+2 : end]
			if end > i+2 && toks[end-1].Type == tokenCloseParen {
				inner = toks[i+2 : end-1]
			}
			// Count top-level comma-separated partition specs.
			count := 0
			depth := 0
			hasPart := false
			for j := 0; j < len(inner); j++ {
				if inner[j].Type == tokenOpenParen {
					depth++
					continue
				}
				if inner[j].Type == tokenCloseParen {
					if depth > 0 {
						depth--
					}
					continue
				}
				if depth != 0 {
					continue
				}
				if tokWord(inner, j, "PARTITION") {
					hasPart = true
					count++
				}
			}
			if hasPart && count != st.partCount {
				b.fail(rng, review.UnsupportedForm, "LOCAL index partition count does not match table "+review.Code(targetRef.written)+".")
				return true
			}
			break
		}
	}
	return false
}

// checkDuplicateIndexCols fails CREATE INDEX when a visible plain-column index
// is created on a column list that already has a visible index on the synced
// table (ORA-01408). Oracle permits multiple indexes on the same columns only
// when at most one is visible and their type, partitioning, or uniqueness differ.
func (b *builder) checkDuplicateIndexCols(v *ast.CreateIndexStmt, targetRef objRef, st *syncedTable, rng review.Range) bool {
	if st == nil || len(st.indexedCols) == 0 || v.Columns == nil || len(v.Columns.Items) == 0 ||
		v.Invisible || v.Reverse || v.FromTables != nil || v.IndexType != nil || v.Multivalue {
		return false
	}
	cols := make([]string, 0, len(v.Columns.Items))
	for _, item := range v.Columns.Items {
		ic, ok := item.(*ast.IndexColumn)
		if !ok || ic.Dir == ast.SORTBY_DESC {
			return false
		}
		cr, ok := ic.Expr.(*ast.ColumnRef)
		if !ok || cr.Column == "" || cr.Schema != "" {
			return false
		}
		cols = append(cols, cr.Column)
	}
	for _, existing := range st.indexedCols {
		if len(existing) != len(cols) {
			continue
		}
		match := true
		for i := range cols {
			if existing[i] != cols[i] {
				match = false
				break
			}
		}
		if match {
			b.fail(rng, review.UnsupportedForm, "Column list is already indexed on "+review.Code(targetRef.written)+".")
			return true
		}
	}
	return false
}

// checkColTypeKnownNonType fails when a column definition (or CREATE TABLE OF)
// writes an unqualified type name that is already occupied in the object
// namespace by a non-type (TABLE, VIEW, MATERIALIZED VIEW, SEQUENCE, routine,
// package) (ORA-00902).
func (b *builder) checkColTypeKnownNonType(tn *ast.TypeName, rng review.Range) bool {
	if tn == nil || tn.Names == nil || len(tn.Names.Items) != 1 || tn.IsPercType || tn.IsPercRowtype {
		return false
	}
	s, ok := tn.Names.Items[0].(*ast.String)
	if !ok || s.Str == "" {
		return false
	}
	upper := strings.ToUpper(s.Str)
	if oracleBuiltinTypeNames[upper] || b.w.current == "" {
		return false
	}
	k := walk.Key{Space: spaceObject, Parent: walk.Path(b.w.current), Name: s.Str}
	if f, kind := b.look(k); f == walk.Present && kind != "" && kind != kindType && kind != kindSynonym {
		b.need(k, walk.Present, []string{kindType}, rng, "", wrongKind(ident(s.Str), kind, kindWord(kindType)))
		return true
	}
	return false
}

// checkColCollationAndRef checks column-level COLLATE on non-character types
// (ORA-43929) and WITH ROWID / SCOPE IS on built-in scalar types (ORA-22893).
func (b *builder) checkColCollationAndRef(c *ast.ColumnDef, existing *syncedCol, rng review.Range) bool {
	baseType := ""
	if c.TypeName != nil && !c.TypeName.IsRef {
		baseType = typeNameBase(c.TypeName)
	} else if existing != nil {
		baseType = existing.typ
	}
	if c.Collation != "" && baseType != "" && isOracleNonCharType(baseType) {
		b.fail(rng, review.UnsupportedForm, "COLLATE is not allowed on non-character column "+review.Code(ident(c.Name))+".")
		return true
	}
	if baseType != "" && isOracleScalarType(baseType) {
		toks := tokens(b.w.sql, c.Loc)
		for i := 0; i+1 < len(toks); i++ {
			if (tokWord(toks, i, "WITH") && tokWord(toks, i+1, "ROWID")) ||
				(tokWord(toks, i, "SCOPE") && tokWord(toks, i+1, "IS")) {
				b.fail(rng, review.UnsupportedForm, "REF attribute clause is not allowed on non-REF column "+review.Code(ident(c.Name))+".")
				return true
			}
		}
	}
	return false
}

// checkModifyColMeta validates ALTER TABLE MODIFY (col ...) against synced
// column metadata (ORA-30673 identity, ORA-54026 virtual, ORA-01451 already
// nullable, ORA-22858 invalid alteration to non-builtin type, ORA-43929).
func (b *builder) checkModifyColMeta(r objRef, st *syncedTable, col *ast.ColumnDef, rng review.Range) bool {
	if st == nil {
		return b.checkColCollationAndRef(col, nil, rng)
	}
	sc, ok := st.colMeta[col.Name]
	if !ok {
		return b.checkColCollationAndRef(col, nil, rng)
	}
	if (col.Identity != nil || col.DropIdentity) && !sc.identity {
		b.fail(rng, review.UnsupportedForm, "Column "+review.Code(ident(col.Name))+" in "+review.Code(r.written)+" is not an identity column.")
		return true
	}
	if col.Virtual != nil && !sc.virtual {
		b.fail(rng, review.UnsupportedForm, "Cannot modify non-virtual column "+review.Code(ident(col.Name))+" in "+review.Code(r.written)+" to a virtual column.")
		return true
	}
	if col.Null && sc.nullable && !col.NotNull {
		b.fail(rng, review.UnsupportedForm, "Column "+review.Code(ident(col.Name))+" in "+review.Code(r.written)+" is already nullable.")
		return true
	}
	if col.TypeName != nil && !col.TypeName.IsRef && !col.TypeName.IsPercType && !col.TypeName.IsPercRowtype && isOracleScalarType(sc.typ) {
		base := typeNameBase(col.TypeName)
		if base != "" && !oracleBuiltinTypeNames[base] {
			b.fail(rng, review.UnsupportedForm, "Cannot alter datatype of column "+review.Code(ident(col.Name))+" to "+review.Code(ident(base))+".")
			return true
		}
	}
	return b.checkColCollationAndRef(col, &sc, rng)
}

// checkDuplicateAnnotations checks ANNOTATIONS (...) for duplicate annotation
// keys within a single clause (ORA-11602, ORA-11562).
func (b *builder) checkDuplicateAnnotations(loc ast.Loc, rng review.Range) bool {
	toks := tokens(b.w.sql, loc)
	for i := 0; i+1 < len(toks); i++ {
		if !tokWord(toks, i, "ANNOTATIONS") || toks[i+1].Type != tokenOpenParen {
			continue
		}
		end := skipBalancedParens(toks, i+1)
		inner := toks[i+2 : end]
		if end > i+2 && toks[end-1].Type == tokenCloseParen {
			inner = toks[i+2 : end-1]
		}
		seen := make(map[string]bool)
		// Split by top-level commas.
		for j := 0; j < len(inner); {
			k := j
			for k < len(inner) && inner[k].Type != tokenComma {
				k++
			}
			seg := inner[j:k]
			j = k + 1
			if len(seg) == 0 {
				continue
			}
			idx := 0
			if tokWord(seg, idx, "ADD") || tokWord(seg, idx, "REPLACE") || tokWord(seg, idx, "DROP") {
				idx++
				if tokWord(seg, idx, "IF") {
					idx++
					if tokWord(seg, idx, "NOT") {
						idx++
					}
					if tokWord(seg, idx, "EXISTS") {
						idx++
					}
				}
			}
			if idx < len(seg) {
				if name, ok := tokenName(seg[idx]); ok {
					if seen[name] {
						b.fail(rng, review.DuplicateName, twice("names", "annotation", name))
						return true
					}
					seen[name] = true
				}
			}
		}
	}
	return false
}

// checkAlterInmemory checks column lists in ALTER TABLE ... INMEMORY (...)
// / NO INMEMORY (...) for duplicate column names (ORA-00957, ORA-64360) and
// missing columns on r (ORA-00904).
func (b *builder) checkAlterInmemory(r objRef, loc ast.Loc, rng review.Range) bool {
	toks := tokens(b.w.sql, loc)
	seenAcross := make(map[string]bool)
	for i := 0; i < len(toks); i++ {
		if toks[i].Type != tokenOpenParen {
			continue
		}
		// Only inspect paren lists directly attached to INMEMORY / MEMCOMPRESS
		// (skip PRIORITY/DISTRIBUTE/DUPLICATE/other non-column parens if any).
		if i == 0 {
			continue
		}
		prev := toks[i-1]
		if prev.Type == tokenQuoted {
			continue
		}
		pw := strings.ToUpper(prev.Str)
		if pw != "INMEMORY" && pw != "MEMCOMPRESS" && pw != "DML" && pw != "LOW" && pw != "HIGH" && pw != "QUERY" && pw != "CAPACITY" {
			continue
		}
		end := skipBalancedParens(toks, i)
		inner := toks[i+1 : end]
		if end > i+1 && toks[end-1].Type == tokenCloseParen {
			inner = toks[i+1 : end-1]
		}
		seenLocal := make(map[string]bool)
		for _, tok := range inner {
			if tok.Type == tokenComma {
				continue
			}
			name, ok := tokenName(tok)
			if !ok {
				break
			}
			if seenLocal[name] || seenAcross[name] {
				b.fail(rng, review.DuplicateName, twice("lists", "column", name))
				return true
			}
			seenLocal[name] = true
			seenAcross[name] = true
			b.columnNeeds(r, name, rng)
		}
		i = end - 1
	}
	return false
}

// checkSupplementalLogGroup checks column lists in ALTER TABLE ... ADD
// SUPPLEMENTAL LOG GROUP g (c1, c2, ...) for duplicate columns (ORA-00957)
// and missing columns on r (ORA-00904).
func (b *builder) checkSupplementalLogGroup(r objRef, loc ast.Loc, rng review.Range) bool {
	toks := tokens(b.w.sql, loc)
	for i := 0; i+3 < len(toks); i++ {
		if tokWord(toks, i, "SUPPLEMENTAL") && tokWord(toks, i+1, "LOG") && tokWord(toks, i+2, "GROUP") {
			// Skip group name, find '('
			j := i + 3
			for j < len(toks) && toks[j].Type != tokenOpenParen {
				j++
			}
			if j >= len(toks) {
				return false
			}
			end := skipBalancedParens(toks, j)
			inner := toks[j+1 : end]
			if end > j+1 && toks[end-1].Type == tokenCloseParen {
				inner = toks[j+1 : end-1]
			}
			seen := make(map[string]bool)
			for k := 0; k < len(inner); k++ {
				if inner[k].Type == tokenComma {
					continue
				}
				if tokWord(inner, k, "NO") && tokWord(inner, k+1, "LOG") {
					k++
					continue
				}
				name, ok := tokenName(inner[k])
				if !ok {
					continue
				}
				if seen[name] {
					b.fail(rng, review.DuplicateName, twice("lists", "column", name))
					return true
				}
				seen[name] = true
				b.columnNeeds(r, name, rng)
			}
			return false
		}
	}
	return false
}

// checkCreateTableExtra checks CREATE TABLE clauses beyond names/duplicates:
// - all columns virtual (ORA-54037)
// - ENABLE/DISABLE PRIMARY KEY without a declared PK (ORA-02432/02433)
// - ENABLE/DISABLE UNIQUE (cols) with duplicate cols (ORA-00957) or without matching UNIQUE (ORA-02434)
// - ILM DELETE_ALL / DISABLE_ALL / ENABLE_ALL / DELETE POLICY (ORA-38346)
// - PARTITION BY REFERENCE (c) without a declared FK named c (ORA-14655)
// - self-referencing REFERENCES without PK/UNIQUE (ORA-02268/02270) or missing ref table (ORA-00942)
// - duplicate annotations (ORA-11562/11602)
func (b *builder) checkCreateTableExtra(v *ast.CreateTableStmt, r objRef, columns []*ast.ColumnDef, d *tableDecl) bool {
	if v.IfNotExists {
		return false
	}
	// All columns virtual (ORA-54037).
	if len(columns) > 0 && v.AsQuery == nil && v.OfType == nil {
		allVirtual := true
		for _, c := range columns {
			if c.Virtual == nil {
				allVirtual = false
				break
			}
		}
		if allVirtual {
			b.fail(r.rng, review.UnsupportedForm, "Table "+review.Code(r.written)+" must have at least one non-virtual column.")
			return true
		}
	}
	// Pass 1: collect all declared PK/UNIQUE column sets and FK names, and check column types/collations.
	declaredFKs := make(map[string]bool)
	uniqueColSets := make(map[string]bool)
	type fkCheckItem struct {
		refTable *ast.ObjectName
		refCols  *ast.List
		rng      review.Range
	}
	var fkChecks []fkCheckItem
	for _, c := range columns {
		crng := rangeOf(c.Loc)
		if crng == (review.Range{}) {
			crng = r.rng
		}
		if b.checkColTypeKnownNonType(c.TypeName, crng) {
			return true
		}
		if b.checkColTypeSelfTable(c.TypeName, r, crng) {
			return true
		}
		if b.checkColCollationAndRef(c, nil, crng) {
			return true
		}
		if c.Constraints != nil {
			for _, item := range c.Constraints.Items {
				cc, ok := item.(*ast.ColumnConstraint)
				if !ok {
					continue
				}
				switch cc.Type {
				case ast.CONSTRAINT_PRIMARY, ast.CONSTRAINT_UNIQUE:
					uniqueColSets[c.Name] = true
				case ast.CONSTRAINT_FOREIGN:
					if cc.Name != "" {
						declaredFKs[cc.Name] = true
					}
					fkChecks = append(fkChecks, fkCheckItem{cc.RefTable, cc.RefColumns, crng})
				}
			}
		}
	}
	for _, c := range outOfLine(v.Constraints) {
		crng := rangeOf(c.Loc)
		if crng == (review.Range{}) {
			crng = r.rng
		}
		switch c.Type {
		case ast.CONSTRAINT_PRIMARY, ast.CONSTRAINT_UNIQUE:
			if cols := stringList(c.Columns); len(cols) > 0 {
				if dup := firstDupString(cols); dup != "" {
					b.fail(crng, review.DuplicateName, twice("lists", "column", dup))
					return true
				}
				uniqueColSets[sortedColSetKey(cols)] = true
			}
		case ast.CONSTRAINT_FOREIGN:
			if c.Name != "" {
				declaredFKs[c.Name] = true
			}
			fkChecks = append(fkChecks, fkCheckItem{c.RefTable, c.RefColumns, crng})
		}
	}
	// Pass 2: validate all foreign key REFERENCES after all PK/UNIQUE keys are collected.
	for _, fk := range fkChecks {
		if b.checkForeignKeyRef(r, fk.refTable, fk.refCols, d.primaries > 0, uniqueColSets, fk.rng) {
			return true
		}
	}
	// Check trailing ENABLE/DISABLE PRIMARY KEY / UNIQUE (...), ILM, PARTITION BY REFERENCE, ANNOTATIONS.
	toks := tokens(b.w.sql, v.Loc)
	depth := 0
	for i := 0; i < len(toks); i++ {
		if toks[i].Type == tokenOpenParen {
			depth++
			continue
		}
		if toks[i].Type == tokenCloseParen {
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth != 0 {
			continue
		}
		switch {
		case tokWord(toks, i, "ILM") && (tokWord(toks, i+1, "DELETE_ALL") || tokWord(toks, i+1, "DISABLE_ALL") || tokWord(toks, i+1, "ENABLE_ALL") || (tokWord(toks, i+1, "DELETE") && tokWord(toks, i+2, "POLICY"))):
			b.fail(r.rng, review.UnsupportedForm, "ILM "+strings.ToUpper(toks[i+1].Str)+" is not allowed in CREATE TABLE.")
			return true
		case tokWord(toks, i, "PARTITION") && tokWord(toks, i+1, "BY") && tokWord(toks, i+2, "REFERENCE") && i+4 < len(toks) && toks[i+3].Type == tokenOpenParen:
			if cname, ok := tokenName(toks[i+4]); ok && !declaredFKs[cname] {
				b.fail(r.rng, review.ObjectNotFound, "Reference partitioning constraint "+review.Code(ident(cname))+" is not defined as a foreign key on "+review.Code(r.written)+".")
				return true
			}
		case tokWord(toks, i, "ENABLE") || tokWord(toks, i, "DISABLE"):
			j := i + 1
			if tokWord(toks, j, "VALIDATE") || tokWord(toks, j, "NOVALIDATE") {
				j++
			}
			if tokWord(toks, j, "PRIMARY") && tokWord(toks, j+1, "KEY") && d.primaries == 0 {
				b.fail(r.rng, review.ObjectNotFound, review.Code(r.written)+" has no primary key.")
				return true
			}
			if tokWord(toks, j, "UNIQUE") && j+1 < len(toks) && toks[j+1].Type == tokenOpenParen {
				end := skipBalancedParens(toks, j+1)
				inner := toks[j+2 : end]
				if end > j+2 && toks[end-1].Type == tokenCloseParen {
					inner = toks[j+2 : end-1]
				}
				var ucols []string
				for _, tok := range inner {
					if tok.Type == tokenComma {
						continue
					}
					if name, ok := tokenName(tok); ok {
						ucols = append(ucols, name)
					}
				}
				if dup := firstDupString(ucols); dup != "" {
					b.fail(r.rng, review.DuplicateName, twice("lists", "column", dup))
					return true
				}
				if len(ucols) > 0 && !uniqueColSets[sortedColSetKey(ucols)] {
					b.fail(r.rng, review.ObjectNotFound, "Unique key on "+review.Code(r.written)+" is not defined.")
					return true
				}
			}
		}
	}
	if b.checkDuplicateAnnotations(v.Loc, r.rng) {
		return true
	}
	return false
}

// checkForeignKeyRef validates a REFERENCES clause in CREATE TABLE or ALTER TABLE:
// - missing referenced table in current schema (no private/public synonym) -> ORA-00942
// - self-reference without column list when table has no primary key -> ORA-02268
// - self-reference to columns with no PK/UNIQUE in CREATE TABLE -> ORA-02270
func (b *builder) checkForeignKeyRef(tableRef objRef, refTable *ast.ObjectName, refCols *ast.List, selfHasPK bool, selfUniqueSets map[string]bool, rng review.Range) bool {
	if refTable == nil || refTable.Name == "" || refTable.DBLink != "" {
		return false
	}
	rr, ok := b.ref(refTable)
	if !ok {
		return false
	}
	isSelf := rr.schema == tableRef.schema && rr.name == tableRef.name
	if isSelf {
		cols := stringList(refCols)
		if len(cols) == 0 {
			if !selfHasPK {
				b.fail(rng, review.ObjectNotFound, "Referenced table "+review.Code(rr.written)+" has no primary key.")
				return true
			}
		} else if selfUniqueSets != nil && !selfUniqueSets[sortedColSetKey(cols)] {
			b.fail(rng, review.UnsupportedForm, "Referenced columns do not match a primary or unique key on "+review.Code(rr.written)+".")
			return true
		}
		return false
	}
	// Non-self reference: check if referenced table is absent (and no synonym).
	if f, _ := b.look(rr.key(spaceTable)); f == walk.Absent {
		if objF, objK := b.look(rr.key(spaceObject)); objF == walk.Present && objK != "" && objK != kindSynonym {
			b.wrongTarget(rr, false)
			return true
		}
		if !generated(spaceTable, rr.name) {
			keys := []walk.Key{rr.key(spaceTable), rr.key(spaceSynonym)}
			if refTable.Schema == "" {
				keys = append(keys, publicKey(rr.name))
			}
			allAbsent := true
			for _, sk := range keys {
				if sf, _ := b.look(sk); sf != walk.Absent {
					allAbsent = false
					break
				}
			}
			if allAbsent {
				b.Fail(keys, walk.Absent, noteOf(rng, review.ObjectNotFound, missingTable(rr.written)))
				return true
			}
		}
	}
	return false
}

func stringList(l *ast.List) []string {
	if l == nil {
		return nil
	}
	var out []string
	for _, item := range l.Items {
		switch x := item.(type) {
		case *ast.String:
			if x.Str != "" {
				out = append(out, x.Str)
			}
		case *ast.ColumnRef:
			if x.Column != "" {
				out = append(out, x.Column)
			}
		}
	}
	return out
}

func firstDupString(xs []string) string {
	seen := make(map[string]bool, len(xs))
	for _, x := range xs {
		if seen[x] {
			return x
		}
		seen[x] = true
	}
	return ""
}

func sortedColSetKey(cols []string) string {
	cp := append([]string(nil), cols...)
	slices.Sort(cp)
	return strings.Join(cp, "\x00")
}

// checkDMLPartitionExt validates PARTITION(...) / SUBPARTITION(...) extensions
// on DML targets (ORA-14109 on views, ORA-14501 on non-partitioned tables,
// ORA-14253/14173 on non-composite tables).
func (b *builder) checkDMLPartitionExt(o *ast.ObjectName, r objRef, pe *ast.PartitionExtClause) bool {
	if pe == nil {
		return false
	}
	rng := rangeOf(pe.Loc)
	if rng == (review.Range{}) {
		rng = r.rng
	}
	if f, kind := b.look(r.key(spaceObject)); f == walk.Present && kind == kindView {
		b.need(r.key(spaceObject), walk.Present, []string{kindTable, kindMView}, rng, "", "Partition-extended name is not allowed on view "+review.Code(r.written)+".")
		return true
	}
	tr := r
	st := b.syncedPartMeta(r)
	if st == nil {
		if sr, sst, ok := b.resolveSynonymTable(o, r); ok {
			tr = sr
			st = b.syncedPartMeta(sr)
			_ = sst
		}
	}
	if st == nil {
		return false
	}
	if st.heap {
		b.fail(rng, review.UnsupportedForm, "Table "+review.Code(tr.written)+" is not partitioned.")
		return true
	}
	if st.partType == metadata.TablePartitionMetadata_TYPE_UNSPECIFIED {
		if st.hasForeignKey {
			return false
		}
		// SYSTEM-partitioned tables (omitted by oracle/driver/sync.go) allow
		// PARTITION (name), but never allow SUBPARTITION or PARTITION FOR (...).
		if pe.IsSubpartition {
			b.fail(rng, review.UnsupportedForm, "Table "+review.Code(tr.written)+" is not composite-partitioned.")
			return true
		}
		if pe.IsFor {
			b.fail(rng, review.UnsupportedForm, "PARTITION FOR (...) requires a range-, list-, or hash-partitioned table, not "+review.Code(tr.written)+".")
			return true
		}
		return false
	}
	if pe.IsSubpartition && !st.hasSubparts {
		b.fail(rng, review.UnsupportedForm, "Table "+review.Code(tr.written)+" is not composite-partitioned.")
		return true
	}
	if pe.IsFor && !pe.IsSubpartition && st.partKeyCount > 0 {
		toks := tokens(b.w.sql, pe.Loc)
		for i := 0; i+1 < len(toks); i++ {
			if tokWord(toks, i, "FOR") && toks[i+1].Type == tokenOpenParen {
				if n := countTopLevelExprs(toks, i+1); n > 0 && n != st.partKeyCount {
					b.fail(rng, review.UnsupportedForm, "PARTITION FOR value count does not match partitioning key count of "+review.Code(tr.written)+".")
					return true
				}
				break
			}
		}
	}
	return false
}

// checkDMLReturning fails standalone SQL DML with RETURNING ... INTO <ident>
// when a target variable after INTO is not a bind variable (ORA-00905).
func (b *builder) checkDMLReturning(loc ast.Loc, returning *ast.List) bool {
	if returning == nil {
		return false
	}
	toks := tokens(b.w.sql, loc)
	for i := 0; i < len(toks); i++ {
		if !tokWord(toks, i, "RETURNING") && !tokWord(toks, i, "RETURN") {
			continue
		}
		for j := i + 1; j+1 < len(toks); j++ {
			if tokWord(toks, j, "INTO") {
				// Next token after INTO must be a bind variable (starts with ':' in source text).
				next := toks[j+1]
				isBind := strings.HasPrefix(next.Str, ":") ||
					(next.Loc >= 0 && next.Loc < len(b.w.sql) && b.w.sql[next.Loc] == ':')
				if !isBind {
					b.fail(b.anchor(), review.UnsupportedForm, "RETURNING INTO requires bind variables in SQL.")
					return true
				}
				return false
			}
		}
	}
	return false
}

// checkDMLInlineSubquery validates UPDATE/DELETE targeting an inline view
// subquery: fails on constant-only / DUAL subqueries (ORA-01732, ORA-41900)
// and checks missing columns when the subquery reads a single known table
// (ORA-00904).
func (b *builder) checkDMLInlineSubquery(target ast.TableExpr) bool {
	sub, ok := target.(*ast.SubqueryRef)
	if !ok || sub == nil {
		return false
	}
	sel, ok := sub.Subquery.(*ast.SelectStmt)
	if !ok || sel == nil {
		return false
	}
	for sel.Larg != nil {
		sel = sel.Larg
	}
	rng := rangeOf(sub.Loc)
	if rng == (review.Range{}) {
		rng = b.anchor()
	}
	// Check FROM SYS.DUAL or single base table in subquery.
	if sel.WithClause == nil && sel.FromClause != nil && len(sel.FromClause.Items) == 1 {
		if tr, ok := sel.FromClause.Items[0].(*ast.TableRef); ok && tr.Name != nil && tr.Dblink == "" {
			if b.isSysDual(tr.Name) {
				b.fail(rng, review.UnsupportedForm, "Data manipulation is not legal on an inline view over DUAL.")
				return true
			}
			// Single base table in subquery: verify bare SELECT columns exist.
			if baseRef, ok := b.ref(tr.Name); ok {
				if st := b.syncedTable(baseRef); st != nil && st.complete {
					if sel.TargetList != nil {
						for _, item := range sel.TargetList.Items {
							rt, ok := item.(*ast.ResTarget)
							if !ok || rt.Expr == nil {
								continue
							}
							if cr, ok := rt.Expr.(*ast.ColumnRef); ok && cr.Column != "" && cr.Column != "*" && cr.Schema == "" {
								if cr.Table == "" || cr.Table == baseRef.name || cr.Table == aliasName(tr.Alias) {
									if !b.isValidBareIdentifier(baseRef, cr.Column) {
										b.columnNeeds(baseRef, cr.Column, rng)
										return true
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

// isSysDual reports whether o resolves to Oracle's SYS.DUAL table (either
// explicitly qualified as SYS.DUAL, or unqualified DUAL with no shadowing
// object in the current schema and PUBLIC.DUAL synced to SYS.DUAL). Exact case
// is required so a quoted table named "dual" does not match.
func (b *builder) isSysDual(o *ast.ObjectName) bool {
	if o == nil || o.DBLink != "" || o.Name != "DUAL" {
		return false
	}
	if o.Schema == "SYS" {
		return true
	}
	if o.Schema != "" {
		return false
	}
	if b.w.otherSpecs || b.w.current == "" || b.w.mutatedSynonyms[publicSchema+".DUAL"] {
		return false
	}
	target, ok := b.w.snap.publicSynonymTargets["DUAL"]
	if !ok || target != [2]string{"SYS", "DUAL"} {
		return false
	}
	dr := objRef{schema: b.w.current, name: "DUAL"}
	if f, _ := b.look(dr.key(spaceTable)); f != walk.Absent {
		return false
	}
	if f, _ := b.look(dr.key(spaceView)); f != walk.Absent {
		return false
	}
	if f, _ := b.look(dr.key(spaceSynonym)); f != walk.Absent {
		return false
	}
	if f, _ := b.look(dr.key(spaceObject)); f == walk.Present {
		return false
	}
	return true
}

// isValidBareIdentifier reports whether a bare identifier in a SQL expression
// over a single table r could resolve to a column in the current ledger state,
// a pseudocolumn, builtin, or synced function/synonym.
func (b *builder) isValidBareIdentifier(r objRef, name string) bool {
	if f, _ := b.look(r.childKey(spaceColumn, name)); f != walk.Absent {
		return true
	}
	upper := strings.ToUpper(name)
	if oraclePseudoColumns[upper] || noParenBuiltins[upper] || sqlBuiltins[upper] || generated(spaceColumn, name) {
		return true
	}
	if b.w.current != "" {
		if f, kind := b.look(walk.Key{Space: spaceObject, Parent: walk.Path(b.w.current), Name: name}); f == walk.Present && (kind == kindFunction || kind == kindSynonym || kind == kindPackage) {
			return true
		}
		if f, _ := b.look(publicKey(name)); f == walk.Present {
			return true
		}
	}
	return false
}

// checkDMLWhereColumns checks bare ColumnRef identifiers in a simple WHERE
// clause of a single-table UPDATE or DELETE against a known complete table
// (ORA-00904).
func (b *builder) checkDMLWhereColumns(r objRef, st *syncedTable, alias string, where ast.ExprNode) bool {
	if where == nil || st == nil || !st.complete {
		return false
	}
	hasSubquery := false
	var refs []*ast.ColumnRef
	ast.Inspect(where, func(m ast.Node) bool {
		switch v := m.(type) {
		case *ast.SelectStmt:
			hasSubquery = true
			return false
		case *ast.ColumnRef:
			refs = append(refs, v)
		}
		return true
	})
	if hasSubquery {
		return false
	}
	for _, cr := range refs {
		if cr.Column == "" || cr.Column == "*" || cr.Schema != "" {
			continue
		}
		if cr.Table != "" && cr.Table != r.name && (alias == "" || cr.Table != alias) {
			continue
		}
		if cr.Table != "" {
			// Table-qualified reference MUST be a column of r.
			if f, _ := b.look(r.childKey(spaceColumn, cr.Column)); f == walk.Absent && !generated(spaceColumn, cr.Column) && !oraclePseudoColumns[strings.ToUpper(cr.Column)] {
				b.columnNeeds(r, cr.Column, rangeOf(cr.Loc))
				return true
			}
			continue
		}
		if !b.isValidBareIdentifier(r, cr.Column) {
			rng := rangeOf(cr.Loc)
			if rng == (review.Range{}) {
				rng = r.rng
			}
			b.columnNeeds(r, cr.Column, rng)
			return true
		}
	}
	return false
}

// checkDMLFromTree inspects FROM / USING / JOIN / PIVOT / UNPIVOT / ANALYTIC
// VIEW / TABLE(...) / FLASHBACK clauses inside a DML statement.
func (b *builder) checkDMLFromTree(te ast.TableExpr) bool {
	if te == nil {
		return false
	}
	switch v := te.(type) {
	case *ast.PivotClause:
		rng := rangeOf(v.Loc)
		if rng == (review.Range{}) {
			rng = b.anchor()
		}
		// Empty IN () list.
		if (v.InItems == nil || len(v.InItems.Items) == 0) && (v.InList == nil || len(v.InList.Items) == 0) {
			b.fail(rng, review.UnsupportedForm, "PIVOT requires a non-empty IN list.")
			return true
		}
		// Inspect aggregate expressions.
		var aggExprs []ast.ExprNode
		if v.Aggregates != nil && len(v.Aggregates.Items) > 0 {
			for _, item := range v.Aggregates.Items {
				if pa, ok := item.(*ast.PivotAggregate); ok {
					aggExprs = append(aggExprs, pa.Expr)
				}
			}
		} else if v.AggFuncs != nil {
			for _, item := range v.AggFuncs.Items {
				if rt, ok := item.(*ast.ResTarget); ok {
					aggExprs = append(aggExprs, rt.Expr)
				} else if e, ok := item.(ast.ExprNode); ok {
					aggExprs = append(aggExprs, e)
				}
			}
		}
		for _, expr := range aggExprs {
			fc, ok := expr.(*ast.FuncCallExpr)
			if !ok || fc == nil || fc.FuncName == nil || fc.FuncName.Name == "" {
				b.fail(rng, review.UnsupportedForm, "PIVOT requires an aggregate function.")
				return true
			}
			if !fc.Star && (fc.Args == nil || len(fc.Args.Items) == 0) {
				b.fail(rng, review.UnsupportedForm, "Aggregate function in PIVOT requires an argument.")
				return true
			}
			if !oracleAggregateBuiltins[strings.ToUpper(fc.FuncName.Name)] {
				if fr, fok := b.ref(fc.FuncName); fok {
					if f, kind := b.look(fr.key(spaceObject)); f == walk.Present && kind != "" && kind != kindFunction && kind != kindPackage {
						if kind != kindSynonym {
							b.need(fr.key(spaceObject), walk.Present, []string{kindFunction}, rng, "", wrongKind(fr.written, kind, "aggregate function"))
							return true
						}
						if _, sst, sok := b.resolveSynonymTable(fc.FuncName, fr); sok && sst != nil {
							b.fail(rng, review.WrongObjectKind, wrongKind(fr.written, kindTable, "aggregate function"))
							return true
						}
					}
				}
			}
		}
		return b.checkDMLFromTree(v.Source)

	case *ast.UnpivotClause:
		rng := rangeOf(v.Loc)
		if rng == (review.Range{}) {
			rng = b.anchor()
		}
		// Collect value column names and pivot (FOR) column names.
		valNames := make(map[string]bool)
		addCols := func(e ast.ExprNode, target map[string]bool) bool {
			if e == nil {
				return false
			}
			switch x := e.(type) {
			case *ast.ColumnRef:
				if x.Column != "" {
					if target[x.Column] {
						return true
					}
					target[x.Column] = true
				}
			case *ast.ParenExpr:
				if cr, ok := x.Expr.(*ast.ColumnRef); ok && cr.Column != "" {
					if target[cr.Column] {
						return true
					}
					target[cr.Column] = true
				}
			}
			return false
		}
		if v.ValueColumns != nil && len(v.ValueColumns.Items) > 0 {
			for _, item := range v.ValueColumns.Items {
				if e, ok := item.(ast.ExprNode); ok {
					if addCols(e, valNames) {
						b.fail(rng, review.DuplicateName, "UNPIVOT declares duplicate output column.")
						return true
					}
				}
			}
		} else {
			addCols(v.ValueCol, valNames)
		}
		pivNames := make(map[string]bool)
		if v.PivotColumn != nil {
			addCols(v.PivotColumn, pivNames)
		} else {
			addCols(v.PivotCol, pivNames)
		}
		// Check FOR (col1, col2) in tokens for duplicate FOR columns (ORA-00918).
		{
			toks := tokens(b.w.sql, v.Loc)
			for i := 0; i+1 < len(toks); i++ {
				if tokWord(toks, i, "FOR") && toks[i+1].Type == tokenOpenParen {
					end := skipBalancedParens(toks, i+1)
					inner := toks[i+2 : end]
					if end > i+2 && toks[end-1].Type == tokenCloseParen {
						inner = toks[i+2 : end-1]
					}
					seenFor := make(map[string]bool)
					for _, tok := range inner {
						if tok.Type == tokenComma {
							continue
						}
						if name, ok := tokenName(tok); ok {
							if seenFor[name] {
								b.fail(rng, review.DuplicateName, twice("lists", "column", name))
								return true
							}
							seenFor[name] = true
						}
					}
					break
				}
			}
		}
		// Oracle 23ai accepts UNPIVOT INCLUDE NULLS(a FOR a IN(a)) in DELETE FROM,
		// so only flag value/FOR collision when INCLUDE NULLS is not specified.
		if !v.IncludeNulls {
			for name := range pivNames {
				if valNames[name] {
					b.fail(rng, review.DuplicateName, "UNPIVOT value column and FOR column cannot share name "+review.Code(ident(name))+".")
					return true
				}
			}
		}
		// Check UNPIVOT IN (...) source columns against source table if known, and for duplicates.
		var srcRef objRef
		var srcTable *syncedTable
		if tr, ok := v.Source.(*ast.TableRef); ok && tr.Name != nil && tr.Dblink == "" {
			if r, ok := b.ref(tr.Name); ok {
				srcRef = r
				srcTable = b.syncedTable(r)
				if srcTable == nil {
					if sr, sst, sok := b.resolveSynonymTable(tr.Name, r); sok {
						srcRef = sr
						srcTable = sst
					}
				}
			}
		}
		if v.InputMappings != nil {
			seenIn := make(map[string]bool)
			for _, item := range v.InputMappings.Items {
				ui, ok := item.(*ast.UnpivotInItem)
				if !ok || ui.InputColumns == nil {
					continue
				}
				for _, ci := range ui.InputColumns.Items {
					if cr, ok := ci.(*ast.ColumnRef); ok && cr.Column != "" {
						if seenIn[cr.Column] {
							b.fail(rng, review.DuplicateName, twice("lists", "column", cr.Column))
							return true
						}
						seenIn[cr.Column] = true
						if srcTable != nil && srcTable.complete {
							b.columnNeeds(srcRef, cr.Column, rng)
						}
					}
				}
			}
		}
		return b.checkDMLFromTree(v.Source)

	case *ast.InlineAnalyticView:
		if v.Clause != nil && v.Clause.Using != nil && v.Clause.Using.DBLink == "" {
			if ur, ok := b.ref(v.Clause.Using); ok {
				if f, kind := b.look(ur.key(spaceObject)); f == walk.Present && kind != "" {
					if kind != kindSynonym {
						b.need(ur.key(spaceObject), walk.Present, []string{"ANALYTIC VIEW"}, ur.rng, "", wrongKind(ur.written, kind, "analytic view"))
						return true
					}
					if _, sst, sok := b.resolveSynonymTable(v.Clause.Using, ur); sok && sst != nil {
						b.fail(ur.rng, review.WrongObjectKind, wrongKind(ur.written, kindTable, "analytic view"))
						return true
					}
				}
			}
		}

	case *ast.TableCollectionExpr:
		rng := rangeOf(v.Loc)
		if rng == (review.Range{}) {
			rng = b.anchor()
		}
		expr := v.Expr
		if expr == nil {
			expr = v.CollectionExpr
		}
		switch e := expr.(type) {
		case *ast.NumberLiteral, *ast.StringLiteral:
			b.fail(rng, review.UnsupportedForm, "TABLE() requires a collection expression, not a scalar literal.")
			return true
		case *ast.FuncCallExpr:
			if e.FuncName != nil && e.FuncName.Name != "" && e.FuncName.DBLink == "" {
				if fr, ok := b.ref(e.FuncName); ok {
					if f, kind := b.look(fr.key(spaceObject)); f == walk.Present && kind != "" && kind != kindFunction && kind != kindPackage && kind != kindType {
						if kind != kindSynonym {
							b.need(fr.key(spaceObject), walk.Present, []string{kindFunction}, fr.rng, "", wrongKind(fr.written, kind, "table function"))
							return true
						}
						if _, sst, sok := b.resolveSynonymTable(e.FuncName, fr); sok && sst != nil {
							b.fail(fr.rng, review.WrongObjectKind, wrongKind(fr.written, kindTable, "table function"))
							return true
						}
					}
				}
			}
		}

	case *ast.JoinClause:
		rng := rangeOf(v.Loc)
		if rng == (review.Range{}) {
			rng = b.anchor()
		}
		// NATURAL FULL JOIN self-join -> ORA-00918.
		if v.Type == ast.JOIN_NATURAL_FULL {
			if lref, lok := v.Left.(*ast.TableRef); lok && lref.Name != nil {
				if rref, rok := v.Right.(*ast.TableRef); rok && rref.Name != nil {
					lr, ok1 := b.ref(lref.Name)
					rr, ok2 := b.ref(rref.Name)
					if ok1 && ok2 && lr.schema == rr.schema && lr.name == rr.name {
						b.fail(rng, review.UnsupportedForm, "Ambiguous NATURAL FULL JOIN on "+review.Code(lr.written)+".")
						return true
					}
				}
			}
		}
		if v.Using != nil && len(v.Using.Items) > 0 {
			for _, item := range v.Using.Items {
				s, ok := item.(*ast.String)
				if !ok || s.Str == "" {
					continue
				}
				for _, side := range []ast.TableExpr{v.Left, v.Right} {
					if tr, ok := side.(*ast.TableRef); ok && tr.Name != nil && tr.Dblink == "" {
						if sr, sok := b.ref(tr.Name); sok {
							if st := b.syncedTable(sr); st != nil && st.complete {
								b.columnNeeds(sr, s.Str, rng)
							}
						}
					}
				}
			}
		}
		return b.checkDMLFromTree(v.Left) || b.checkDMLFromTree(v.Right)

	case *ast.TableRef:
		if v.Name != nil && v.Dblink == "" {
			if r, ok := b.ref(v.Name); ok {
				if b.checkDMLPartitionExt(v.Name, r, v.PartitionExt) {
					return true
				}
			}
		}
		if v.Flashback != nil {
			if strings.EqualFold(v.Flashback.Type, "TIMESTAMP") {
				if _, isNum := v.Flashback.Expr.(*ast.NumberLiteral); isNum {
					rng := rangeOf(v.Flashback.Loc)
					if rng == (review.Range{}) {
						rng = b.anchor()
					}
					b.fail(rng, review.UnsupportedForm, "AS OF TIMESTAMP requires a timestamp expression, not a number.")
					return true
				}
			}
			if b.checkFlashbackVersionsPeriod(v.Flashback, v.Name, v.Dblink) {
				return true
			}
		}
	}
	return false
}

// checkDMLTableFuncRef checks FROM a(1) / FROM a(a) / FROM TABLES(1) where
// the parser produces a TableRef or FuncCall-like FROM item naming an existing
// non-function object (ORA-00904).
func (b *builder) checkDMLFromList(from *ast.List) bool {
	if from == nil {
		return false
	}
	for _, item := range from.Items {
		if te, ok := item.(ast.TableExpr); ok {
			if b.checkDMLFromTree(te) {
				return true
			}
		}
		if fc, ok := item.(*ast.FuncCallExpr); ok && fc.FuncName != nil {
			if fr, ok := b.ref(fc.FuncName); ok {
				if f, kind := b.look(fr.key(spaceObject)); f == walk.Present && kind != "" && kind != kindFunction && kind != kindPackage {
					if kind != kindSynonym {
						b.need(fr.key(spaceObject), walk.Present, []string{kindFunction}, fr.rng, "", wrongKind(fr.written, kind, "table function"))
						return true
					}
					if _, sst, sok := b.resolveSynonymTable(fc.FuncName, fr); sok && sst != nil {
						b.fail(fr.rng, review.WrongObjectKind, wrongKind(fr.written, kindTable, "table function"))
						return true
					}
				}
			}
		}
	}
	return false
}

// checkDMLDuplicateSetAndValue checks SET clauses on UPDATE/INSERT/MERGE for:
// - duplicate target columns (ORA-00957)
// - SET VALUE(alias) = ... on a relational (non-object) synced table (ORA-00904)
func (b *builder) checkDMLDuplicateSetAndValue(r objRef, st *syncedTable, setClauses *ast.List, alias string) bool {
	if setClauses == nil || len(setClauses.Items) == 0 {
		return false
	}
	seen := make(map[string]bool)
	checkCol := func(cr *ast.ColumnRef) bool {
		if cr == nil || cr.Column == "" || cr.Schema != "" {
			return false
		}
		if cr.Table != "" && cr.Table != r.name && (alias == "" || cr.Table != alias) {
			return false
		}
		if seen[cr.Column] {
			rng := rangeOf(cr.Loc)
			if rng == (review.Range{}) {
				rng = r.rng
			}
			b.fail(rng, review.DuplicateName, twice("assigns", "column", cr.Column))
			return true
		}
		seen[cr.Column] = true
		return false
	}
	for _, item := range setClauses.Items {
		sc, ok := item.(*ast.SetClause)
		if !ok {
			continue
		}
		if sc.ValueOf != "" && st != nil && !st.skipDump {
			rng := rangeOf(sc.Loc)
			if rng == (review.Range{}) {
				rng = r.rng
			}
			b.fail(rng, review.UnsupportedForm, "SET VALUE(...) is only valid on an object table, not "+review.Code(r.written)+".")
			return true
		}
		if checkCol(sc.Column) {
			return true
		}
		if sc.Columns != nil {
			for _, ci := range sc.Columns.Items {
				if cr, ok := ci.(*ast.ColumnRef); ok && checkCol(cr) {
					return true
				}
			}
		}
	}
	return false
}

// checkCreateMviewExtra checks CREATE MATERIALIZED VIEW options and query:
// - ORGANIZATION INDEX without PRIMARY KEY (ORA-25175)
// - INMEMORY DISTRIBUTE BY [SUB]PARTITION without PARTITION BY (ORA-64370/64372)
// - single-table SELECT query referencing a missing column on a synced complete table (ORA-00904)
func (b *builder) checkCreateMviewExtra(v *ast.CreateViewStmt, r objRef) bool {
	if v.Name == nil || v.Name.Loc.End <= 0 || v.Query == nil {
		return false
	}
	qStart := v.Loc.End
	if sel, ok := v.Query.(*ast.SelectStmt); ok && sel.Loc.Start > v.Name.Loc.End {
		qStart = sel.Loc.Start
	}
	if qStart > v.Name.Loc.End {
		toks := tokens(b.w.sql, ast.Loc{Start: v.Name.Loc.End, End: qStart})
		hasPartBy := false
		hasOrgIndex := false
		hasWithRowid := v.WithRowID
		hasDistPart := false
		for i := 0; i < len(toks); i++ {
			if tokWord(toks, i, "PARTITION") && tokWord(toks, i+1, "BY") {
				hasPartBy = true
			}
			if tokWord(toks, i, "ORGANIZATION") && tokWord(toks, i+1, "INDEX") {
				hasOrgIndex = true
			}
			if tokWord(toks, i, "WITH") && tokWord(toks, i+1, "ROWID") {
				hasWithRowid = true
			}
			if tokWord(toks, i, "DISTRIBUTE") && tokWord(toks, i+1, "BY") && (tokWord(toks, i+2, "PARTITION") || tokWord(toks, i+2, "SUBPARTITION")) {
				hasDistPart = true
			}
		}
		if hasOrgIndex {
			missingPK := hasWithRowid
			if !missingPK && !v.WithPK {
				// Oracle defaults materialized views to WITH PRIMARY KEY, so
				// ORGANIZATION INDEX fails only when the query source is known
				// to have no primary key (e.g. SYS.DUAL or a single base table
				// whose primary key is absent).
				if sel, ok := v.Query.(*ast.SelectStmt); ok && sel != nil && sel.WithClause == nil && sel.Larg == nil && sel.FromClause != nil && len(sel.FromClause.Items) == 1 {
					if tr, ok := sel.FromClause.Items[0].(*ast.TableRef); ok && tr.Name != nil && tr.Dblink == "" {
						if b.isSysDual(tr.Name) {
							missingPK = true
						} else if br, ok := b.ref(tr.Name); ok {
							if f, _ := b.look(br.childKey(spacePrimary, "")); f == walk.Absent {
								missingPK = true
							}
						}
					}
				}
			}
			if missingPK {
				b.fail(r.rng, review.UnsupportedForm, "Index-organized materialized view "+review.Code(r.written)+" requires a PRIMARY KEY.")
				return true
			}
		}
		if hasDistPart && !hasPartBy {
			b.fail(r.rng, review.UnsupportedForm, "INMEMORY DISTRIBUTE BY PARTITION requires a partitioned materialized view.")
			return true
		}
	}
	// Check single-table SELECT column references against a synced complete table.
	if sel, ok := v.Query.(*ast.SelectStmt); ok && sel != nil && sel.WithClause == nil && sel.Larg == nil && sel.FromClause != nil && len(sel.FromClause.Items) == 1 {
		if tr, ok := sel.FromClause.Items[0].(*ast.TableRef); ok && tr.Name != nil && tr.Dblink == "" {
			if br, ok := b.ref(tr.Name); ok {
				if st := b.syncedTable(br); st != nil && st.complete && sel.TargetList != nil {
					for _, item := range sel.TargetList.Items {
						rt, ok := item.(*ast.ResTarget)
						if !ok || rt.Expr == nil {
							continue
						}
						if cr, ok := rt.Expr.(*ast.ColumnRef); ok && cr.Column != "" && cr.Column != "*" && cr.Schema == "" {
							if cr.Table == "" || cr.Table == br.name || cr.Table == aliasName(tr.Alias) {
								if !b.isValidBareIdentifier(br, cr.Column) {
									b.columnNeeds(br, cr.Column, r.rng)
									return true
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

// checkTriggerUpdateOf verifies that columns listed in BEFORE/AFTER UPDATE OF
// col1, col2 ON table exist on the target table (ORA-00904).
func (b *builder) checkTriggerUpdateOf(v *ast.CreateTriggerStmt, t objRef) {
	if v.Table == nil || v.Table.Loc.Start <= v.Loc.Start {
		return
	}
	toks := tokens(b.w.sql, ast.Loc{Start: v.Loc.Start, End: v.Table.Loc.Start})
	for i := 0; i+2 < len(toks); i++ {
		if !tokWord(toks, i, "UPDATE") || !tokWord(toks, i+1, "OF") {
			continue
		}
		for j := i + 2; j < len(toks); j++ {
			if tokWord(toks, j, "ON") || tokWord(toks, j, "OR") {
				break
			}
			if toks[j].Type == tokenComma {
				continue
			}
			if name, ok := tokenName(toks[j]); ok {
				b.columnNeeds(t, name, t.rng)
			}
		}
		break
	}
}

// dmlCheckInsertTableMeta checks INSERT column-count mismatches when no
// explicit column list is given (using visibleInsCols from synced metadata),
// omitted NOT NULL columns without defaults when a row is guaranteed to be
// inserted without error logging or triggers (ORA-01400), and INSERT BY NAME
// SELECT column names (ORA-63875).
func (b *builder) dmlCheckInsertTableMeta(o *ast.ObjectName, r objRef, cols *ast.List, vals *ast.List, sel *ast.SelectStmt, byName, guaranteedRow bool, rng review.Range) bool {
	if rng == (review.Range{}) {
		rng = r.rng
	}
	tr := r
	st := b.syncedColMeta(r)
	if st == nil {
		if sr, _, sok := b.resolveSynonymTable(o, r); sok {
			tr = sr
			st = b.syncedColMeta(sr)
		}
	}
	if st == nil || !st.complete {
		return false
	}
	// INSERT BY NAME SELECT ...: verify every SELECT output name exists in st.
	if byName && sel != nil {
		s := sel
		for s.Larg != nil {
			s = s.Larg
		}
		if s.TargetList != nil {
			for _, item := range s.TargetList.Items {
				rt, ok := item.(*ast.ResTarget)
				if !ok {
					continue
				}
				if _, isStar := rt.Expr.(*ast.Star); isStar {
					continue
				}
				if cr, isCR := rt.Expr.(*ast.ColumnRef); isCR && cr.Column == "*" {
					continue
				}
				colName := rt.Name
				if colName == "" {
					switch e := rt.Expr.(type) {
					case *ast.ColumnRef:
						colName = e.Column
					case *ast.NumberLiteral:
						colName = e.Val
					}
				}
				if colName != "" && !st.columns[colName] && !generated(spaceColumn, colName) {
					b.columnNeeds(tr, colName, rng)
					return true
				}
			}
		}
		return false
	}
	// Explicit column list: check for omitted NOT NULL columns with no default/identity/virtual
	// only when a row is guaranteed to be inserted without error logging and no trigger may fire
	// (ORA-01400). INSERT ... SELECT and conditional multi-table INSERT may insert zero rows.
	if cols != nil && len(cols.Items) > 0 {
		if !guaranteedRow || b.w.triggers {
			return false
		}
		listed := make(map[string]bool, len(cols.Items))
		for _, item := range cols.Items {
			if cr, ok := item.(*ast.ColumnRef); ok && cr.Column != "" {
				listed[cr.Column] = true
			}
		}
		for colName, sc := range st.colMeta {
			if !listed[colName] && !sc.nullable && !sc.hasDefault && !sc.identity && !sc.virtual && !sc.invisible {
				b.fail(rng, review.UnsupportedForm, "INSERT into "+review.Code(tr.written)+" omits NOT NULL column "+review.Code(ident(colName))+".")
				return true
			}
		}
		return false
	}
	// Implicit column list: compare VALUES or non-star SELECT count against visibleInsCols.
	if vals != nil && len(vals.Items) > 0 {
		if len(vals.Items) != st.visibleInsCols {
			b.fail(rng, review.UnsupportedForm, "INSERT value count does not match column count of "+review.Code(tr.written)+".")
			return true
		}
		return false
	}
	if sel != nil {
		if nSel, ok := selectColumnCount(sel); ok && nSel != st.visibleInsCols {
			b.fail(rng, review.UnsupportedForm, "INSERT SELECT list count does not match column count of "+review.Code(tr.written)+".")
			return true
		}
	}
	return false
}

// checkModifyPartitionBy validates ALTER TABLE t MODIFY PARTITION BY ...:
//   - checks PARTITION BY {RANGE|LIST|HASH} (cols) and SUBPARTITION BY {RANGE|LIST|HASH} (cols)
//     for duplicate column names (ORA-00957) and missing columns on r (ORA-00904)
//   - checks PARTITION BY REFERENCE (c) when synced table r has no foreign keys (ORA-25129)
func (b *builder) checkModifyPartitionBy(r objRef, loc ast.Loc, rng review.Range) bool {
	toks := tokens(b.w.sql, loc)
	for i := 0; i+3 < len(toks); i++ {
		if !(tokWord(toks, i, "PARTITION") || tokWord(toks, i, "SUBPARTITION")) || !tokWord(toks, i+1, "BY") {
			continue
		}
		if tokWord(toks, i+2, "REFERENCE") && i+4 < len(toks) && toks[i+3].Type == tokenOpenParen {
			if st := b.syncedColMeta(r); st != nil && !st.hasForeignKey {
				if cname, ok := tokenName(toks[i+4]); ok {
					b.fail(rng, review.ObjectNotFound, "Reference partitioning constraint "+review.Code(ident(cname))+" is not a foreign key on "+review.Code(r.written)+".")
					return true
				}
			}
			continue
		}
		if !tokWord(toks, i+2, "RANGE") && !tokWord(toks, i+2, "LIST") && !tokWord(toks, i+2, "HASH") {
			continue
		}
		j := i + 3
		if toks[j].Type != tokenOpenParen {
			continue
		}
		seen := make(map[string]bool)
		for j++; j < len(toks) && toks[j].Type != tokenCloseParen; j++ {
			if toks[j].Type == tokenComma {
				continue
			}
			name, ok := tokenName(toks[j])
			if !ok {
				break
			}
			if seen[name] {
				b.fail(rng, review.DuplicateName, twice("lists", "column", name))
				return true
			}
			seen[name] = true
			b.columnNeeds(r, name, rng)
		}
	}
	return false
}

// checkPartitionColClauses checks column-specific storage subclauses inside
// partition maintenance actions: SPLIT PARTITION ... NESTED TABLE col ...
// (ORA-22912) and MODIFY PARTITION ... LOB (col) (...) (ORA-00904).
func (b *builder) checkPartitionColClauses(r objRef, st *syncedTable, loc ast.Loc, rng review.Range) bool {
	toks := tokens(b.w.sql, loc)
	for i := 0; i < len(toks); i++ {
		if tokWord(toks, i, "NESTED") && tokWord(toks, i+1, "TABLE") && i+2 < len(toks) {
			if col, ok := tokenName(toks[i+2]); ok && !strings.EqualFold(col, "COLUMN_VALUE") {
				b.columnNeeds(r, col, rng)
				if st != nil {
					if sc, ok := st.colMeta[col]; ok && isOracleScalarType(sc.typ) {
						b.fail(rng, review.UnsupportedForm, "Column "+review.Code(ident(col))+" in "+review.Code(r.written)+" is not a nested table column.")
						return true
					}
				}
			}
		}
		if tokWord(toks, i, "LOB") && i+2 < len(toks) && toks[i+1].Type == tokenOpenParen {
			if col, ok := tokenName(toks[i+2]); ok && !strings.EqualFold(col, "COLUMN_VALUE") {
				b.columnNeeds(r, col, rng)
				if st != nil {
					if sc, ok := st.colMeta[col]; ok && isOracleScalarType(sc.typ) {
						b.fail(rng, review.UnsupportedForm, "Column "+review.Code(ident(col))+" in "+review.Code(r.written)+" is not a LOB column.")
						return true
					}
				}
			}
		}
	}
	return false
}

// moveHasPartition reports whether an ALTER TABLE ... MOVE action is a
// MOVE PARTITION or MOVE SUBPARTITION clause (not a whole-table MOVE).
func (b *builder) moveHasPartition(loc ast.Loc) bool {
	toks := tokens(b.w.sql, loc)
	for i := 0; i+1 < len(toks); i++ {
		if tokWord(toks, i, "MOVE") && (tokWord(toks, i+1, "PARTITION") || tokWord(toks, i+1, "SUBPARTITION")) {
			return true
		}
	}
	return false
}

func hasTokenWord(toks []parser.Token, word string) bool {
	for i := range toks {
		if tokWord(toks, i, word) {
			return true
		}
	}
	return false
}

// checkAddPartitionNonHashOpts rejects ADD PARTITION options that Oracle 23ai
// permits only on HASH-partitioned tables (verified on Oracle 23ai Free against
// RANGE, LIST, HASH, and SYSTEM tables): top-level PARALLEL/NOPARALLEL
// (ORA-14267) and INVALIDATE GLOBAL INDEXES (ORA-30564).
func (b *builder) checkAddPartitionNonHashOpts(rest []parser.Token, rng review.Range) bool {
	depth := 0
	for i := 0; i < len(rest); i++ {
		if rest[i].Type == tokenOpenParen {
			depth++
			continue
		}
		if rest[i].Type == tokenCloseParen {
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth != 0 {
			continue
		}
		if tokWord(rest, i, "PARALLEL") || tokWord(rest, i, "NOPARALLEL") {
			b.fail(rng, review.UnsupportedForm, strings.ToUpper(rest[i].Str)+" is not allowed in ADD PARTITION on a non-hash-partitioned table.")
			return true
		}
		if tokWord(rest, i, "INVALIDATE") && tokWord(rest, i+1, "GLOBAL") && tokWord(rest, i+2, "INDEXES") {
			b.fail(rng, review.UnsupportedForm, "INVALIDATE GLOBAL INDEXES is not allowed in ADD PARTITION on a non-hash-partitioned table.")
			return true
		}
	}
	return false
}

// checkObjectViewSelfRef rejects CREATE VIEW v OF v (same schema and name when
// not OR REPLACE) and CREATE VIEW v OF ... UNDER v (a view cannot inherit
// under itself).
func (b *builder) checkObjectViewSelfRef(v *ast.CreateViewStmt, r objRef) bool {
	if v.Name == nil || v.Name.Loc.End <= 0 || v.Name.Loc.End >= v.Loc.End {
		return false
	}
	toks := tokens(b.w.sql, ast.Loc{Start: v.Name.Loc.End, End: v.Loc.End})
	if len(toks) < 2 || !tokWord(toks, 0, "OF") {
		return false
	}
	if !v.OrReplace {
		if schema, name, ok := readQualifiedTokenName(toks, 1, b.w.current); ok && schema == r.schema && name == r.name {
			b.fail(r.rng, review.UnsupportedForm, "View "+review.Code(r.written)+" cannot have the same name as its object type.")
			return true
		}
	}
	for i := 1; i+1 < len(toks); i++ {
		if tokWord(toks, i, "AS") {
			break
		}
		if tokWord(toks, i, "UNDER") {
			if schema, name, ok := readQualifiedTokenName(toks, i+1, b.w.current); ok && schema == r.schema && name == r.name {
				b.fail(r.rng, review.UnsupportedForm, "View "+review.Code(r.written)+" cannot be created UNDER itself.")
				return true
			}
		}
	}
	return false
}

// dmlCheckInsertSetNotNull checks guaranteed single-row INSERT INTO t SET ...
// (without LOG ERRORS or triggers) for omitted NOT NULL columns without
// defaults/identity/virtual (ORA-01400).
func (b *builder) dmlCheckInsertSetNotNull(o *ast.ObjectName, r objRef, setClauses *ast.List, alias string, rng review.Range) bool {
	if b.w.triggers || setClauses == nil || len(setClauses.Items) == 0 {
		return false
	}
	tr := r
	st := b.syncedColMeta(r)
	if st == nil {
		if sr, _, sok := b.resolveSynonymTable(o, r); sok {
			tr = sr
			st = b.syncedColMeta(sr)
		}
	}
	if st == nil || !st.complete {
		return false
	}
	listed := make(map[string]bool)
	record := func(cr *ast.ColumnRef) {
		if cr == nil || cr.Column == "" || cr.Schema != "" {
			return
		}
		if cr.Table != "" && cr.Table != r.name && cr.Table != tr.name && (alias == "" || cr.Table != alias) {
			return
		}
		listed[cr.Column] = true
	}
	for _, item := range setClauses.Items {
		sc, ok := item.(*ast.SetClause)
		if !ok {
			continue
		}
		record(sc.Column)
		if sc.Columns != nil {
			for _, ci := range sc.Columns.Items {
				if cr, ok := ci.(*ast.ColumnRef); ok {
					record(cr)
				}
			}
		}
	}
	for colName, sc := range st.colMeta {
		if !listed[colName] && !sc.nullable && !sc.hasDefault && !sc.identity && !sc.virtual && !sc.invisible {
			b.fail(rng, review.UnsupportedForm, "INSERT into "+review.Code(tr.written)+" omits NOT NULL column "+review.Code(ident(colName))+".")
			return true
		}
	}
	return false
}

// checkFlashbackVersionsPeriod rejects AS OF / VERSIONS PERIOD FOR <p> with
// numeric literal bounds (ORA-55603 / ORA-00932: valid-time period bounds must
// be date or timestamp expressions).
func (b *builder) checkFlashbackVersionsPeriod(fb *ast.FlashbackClause, name *ast.ObjectName, dblink string) bool {
	if fb == nil || !fb.IsPeriodFor {
		return false
	}
	rng := rangeOf(fb.Loc)
	if rng == (review.Range{}) {
		rng = b.anchor()
	}
	if _, isNum := fb.Expr.(*ast.NumberLiteral); isNum {
		b.fail(rng, review.UnsupportedForm, "PERIOD FOR requires a date or timestamp expression, not a number.")
		return true
	}
	if _, isNum := fb.VersionsLow.(*ast.NumberLiteral); isNum {
		b.fail(rng, review.UnsupportedForm, "VERSIONS PERIOD FOR requires date or timestamp bounds, not a number.")
		return true
	}
	if _, isNum := fb.VersionsHigh.(*ast.NumberLiteral); isNum {
		b.fail(rng, review.UnsupportedForm, "VERSIONS PERIOD FOR requires date or timestamp bounds, not a number.")
		return true
	}
	return false
}

// checkAlterTableAddRefOnScalar rejects table-level ADD REF(col) WITH ROWID /
// ADD SCOPE FOR (col) IS ... when col is a known built-in scalar column
// (ORA-22893).
func (b *builder) checkAlterTableAddRefOnScalar(r objRef, loc ast.Loc, rng review.Range) bool {
	st := b.syncedColMeta(r)
	if st == nil {
		return false
	}
	toks := tokens(b.w.sql, loc)
	for i := 0; i+3 < len(toks); i++ {
		if tokWord(toks, i, "REF") && toks[i+1].Type == tokenOpenParen {
			if col, ok := tokenName(toks[i+2]); ok && toks[i+3].Type == tokenCloseParen {
				if sc, ok := st.colMeta[col]; ok && isOracleScalarType(sc.typ) {
					b.fail(rng, review.UnsupportedForm, "REF attribute clause is not allowed on non-REF column "+review.Code(ident(col))+".")
					return true
				}
			}
		}
		if tokWord(toks, i, "SCOPE") && tokWord(toks, i+1, "FOR") && i+4 < len(toks) && toks[i+2].Type == tokenOpenParen {
			if col, ok := tokenName(toks[i+3]); ok && toks[i+4].Type == tokenCloseParen {
				if sc, ok := st.colMeta[col]; ok && isOracleScalarType(sc.typ) {
					b.fail(rng, review.UnsupportedForm, "SCOPE FOR clause is not allowed on non-REF column "+review.Code(ident(col))+".")
					return true
				}
			}
		}
	}
	return false
}

// countTopLevelExprs counts comma-separated expressions inside balanced parens
// starting at toks[openIdx].
func countTopLevelExprs(toks []parser.Token, openIdx int) int {
	if openIdx < 0 || openIdx >= len(toks) || toks[openIdx].Type != tokenOpenParen {
		return 0
	}
	end := skipBalancedParens(toks, openIdx)
	if end <= openIdx+2 || toks[end-1].Type != tokenCloseParen {
		return 0
	}
	items := 1
	d := 0
	for _, t := range toks[openIdx+1 : end-1] {
		if t.Type == tokenOpenParen {
			d++
		} else if t.Type == tokenCloseParen && d > 0 {
			d--
		} else if t.Type == tokenComma && d == 0 {
			items++
		}
	}
	return items
}

// checkRangeSplitBound rejects SPLIT PARTITION p AT (v) when v is an integer
// literal outside the open interval (prevHigh, high) of p (ORA-14080).
func (b *builder) checkRangeSplitBound(r objRef, st *syncedTable, pname string, rest []parser.Token, rng review.Range) bool {
	hi, hasHi := st.partRangeHigh[pname]
	lo := st.partRangePrev[pname]
	hasPrev := st.partRangeHasPrev[pname]
	if !hasHi && !hasPrev {
		return false
	}
	for i := 0; i+1 < len(rest); i++ {
		if tokWord(rest, i, "AT") && rest[i+1].Type == tokenOpenParen {
			end := skipBalancedParens(rest, i+1)
			if end <= i+2 || rest[end-1].Type != tokenCloseParen {
				return false
			}
			inner := rest[i+2 : end-1]
			if len(inner) == 1 && inner[0].Type != tokenQuoted && inner[0].Type != tokenString {
				if val, err := strconv.ParseInt(inner[0].Str, 10, 64); err == nil {
					if hasHi && val >= hi {
						b.fail(rng, review.UnsupportedForm, "Split bound for partition "+review.Code(ident(pname))+" of "+review.Code(r.written)+" must be less than its high bound.")
						return true
					}
					if hasPrev && val <= lo {
						b.fail(rng, review.UnsupportedForm, "Split bound for partition "+review.Code(ident(pname))+" of "+review.Code(r.written)+" must be greater than the preceding partition's high bound.")
						return true
					}
				}
			}
			break
		}
	}
	return false
}

// checkColTypeSelfTable rejects a CREATE TABLE column whose non-builtin
// datatype names the table itself in the same schema (tables and types share
// Oracle's object namespace).
func (b *builder) checkColTypeSelfTable(tn *ast.TypeName, r objRef, rng review.Range) bool {
	if tn == nil || tn.Names == nil || tn.IsRef || tn.IsPercType || tn.IsPercRowtype {
		return false
	}
	switch len(tn.Names.Items) {
	case 1:
		s, ok := tn.Names.Items[0].(*ast.String)
		if !ok || s.Str == "" || oracleBuiltinTypeNames[strings.ToUpper(s.Str)] {
			return false
		}
		if r.schema == b.w.current && s.Str == r.name {
			b.fail(rng, review.UnsupportedForm, "Column datatype cannot have the same name as table "+review.Code(r.written)+".")
			return true
		}
	case 2:
		s0, ok0 := tn.Names.Items[0].(*ast.String)
		s1, ok1 := tn.Names.Items[1].(*ast.String)
		if !ok0 || !ok1 || s1.Str == "" || oracleBuiltinTypeNames[strings.ToUpper(s1.Str)] {
			return false
		}
		if s0.Str == r.schema && s1.Str == r.name {
			b.fail(rng, review.UnsupportedForm, "Column datatype cannot have the same name as table "+review.Code(r.written)+".")
			return true
		}
	}
	return false
}

// checkSelectFromTables checks that direct FROM table references in a SELECT
// (and EXISTS/IN subqueries in its WHERE clause) exist, respecting CTE names
// declared in WITH clauses (including recursive CTEs).
func (b *builder) checkSelectFromTables(stmt ast.StmtNode, outerCTEs map[string]bool) {
	sel, ok := stmt.(*ast.SelectStmt)
	if !ok || sel == nil {
		return
	}
	ctes := outerCTEs
	if sel.WithClause != nil && sel.WithClause.CTEs != nil && len(sel.WithClause.CTEs.Items) > 0 {
		merged := make(map[string]bool, len(outerCTEs)+len(sel.WithClause.CTEs.Items))
		for k, v := range outerCTEs {
			merged[k] = v
		}
		for _, item := range sel.WithClause.CTEs.Items {
			if cte, ok := item.(*ast.CTE); ok && cte.Name != "" {
				merged[cte.Name] = true
			}
		}
		for _, item := range sel.WithClause.CTEs.Items {
			if cte, ok := item.(*ast.CTE); ok && cte.Query != nil {
				b.checkSelectFromTables(cte.Query, merged)
			}
		}
		ctes = merged
	}
	if sel.Larg != nil {
		b.checkSelectFromTables(sel.Larg, ctes)
	}
	if sel.Rarg != nil {
		b.checkSelectFromTables(sel.Rarg, ctes)
	}
	if sel.FromClause != nil {
		b.checkFromListTables(sel.FromClause, ctes)
	}
	if sel.WhereClause != nil {
		b.checkWhereSubqueryTables(sel.WhereClause, ctes)
	}
}

func (b *builder) checkFromListTables(from *ast.List, ctes map[string]bool) {
	if from == nil {
		return
	}
	for _, item := range from.Items {
		if te, ok := item.(ast.TableExpr); ok {
			b.checkFromExprTables(te, ctes)
		}
	}
}

func (b *builder) checkFromExprTables(te ast.TableExpr, ctes map[string]bool) {
	if te == nil {
		return
	}
	switch v := te.(type) {
	case *ast.TableRef:
		if v.Name == nil || v.Name.Name == "" || v.Dblink != "" || v.Name.DBLink != "" {
			return
		}
		if v.Name.Schema == "" && ctes[v.Name.Name] {
			return
		}
		if b.isSysDual(v.Name) {
			return
		}
		if r, ok := b.ref(v.Name); ok {
			b.missingTarget(v.Name, r)
		}
	case *ast.JoinClause:
		b.checkFromExprTables(v.Left, ctes)
		b.checkFromExprTables(v.Right, ctes)
	case *ast.SubqueryRef:
		b.checkSelectFromTables(v.Subquery, ctes)
	case *ast.PivotClause:
		b.checkFromExprTables(v.Source, ctes)
	case *ast.UnpivotClause:
		b.checkFromExprTables(v.Source, ctes)
	}
}

func (b *builder) checkWhereSubqueryTables(where ast.ExprNode, ctes map[string]bool) {
	if where == nil {
		return
	}
	ast.Inspect(where, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectStmt); ok {
			b.checkSelectFromTables(sel, ctes)
			return false
		}
		return true
	})
}

// checkDMLSetSubqueryCounts rejects multi-column SET (c1, c2) = (SELECT ...)
// when the subquery column count does not match len(Columns) (ORA-00947/00913).
func (b *builder) checkDMLSetSubqueryCounts(setClauses *ast.List) bool {
	if setClauses == nil {
		return false
	}
	for _, item := range setClauses.Items {
		sc, ok := item.(*ast.SetClause)
		if !ok || sc.Columns == nil || len(sc.Columns.Items) == 0 || sc.Value == nil {
			continue
		}
		var sel *ast.SelectStmt
		switch e := sc.Value.(type) {
		case *ast.SubqueryExpr:
			sel, _ = e.Subquery.(*ast.SelectStmt)
		case *ast.ParenExpr:
			if sub, ok := e.Expr.(*ast.SubqueryExpr); ok {
				sel, _ = sub.Subquery.(*ast.SelectStmt)
			}
		}
		if sel != nil {
			if nSel, ok := selectColumnCount(sel); ok && nSel != len(sc.Columns.Items) {
				rng := rangeOf(sc.Loc)
				if rng == (review.Range{}) {
					rng = b.anchor()
				}
				b.fail(rng, review.UnsupportedForm, "SET column count does not match subquery SELECT list count.")
				return true
			}
		}
	}
	return false
}

// checkDMLSetRhsAndSubquery checks multi-column SET subquery column counts
// (ORA-00947) and bare-identifier RHS values on single-table UPDATE SET /
// INSERT SET (ORA-00904 / ORA-00984).
func (b *builder) checkDMLSetRhsAndSubquery(r objRef, st *syncedTable, setClauses *ast.List, alias string, isInsert bool) bool {
	if b.checkDMLSetSubqueryCounts(setClauses) {
		return true
	}
	if setClauses == nil || st == nil || !st.complete {
		return false
	}
	for _, item := range setClauses.Items {
		sc, ok := item.(*ast.SetClause)
		if !ok || sc.Value == nil {
			continue
		}
		cr, ok := sc.Value.(*ast.ColumnRef)
		if !ok || cr.Column == "" || cr.Column == "*" || cr.Schema != "" {
			continue
		}
		if cr.Table != "" && cr.Table != r.name && (alias == "" || cr.Table != alias) {
			continue
		}
		rng := rangeOf(cr.Loc)
		if rng == (review.Range{}) {
			rng = r.rng
		}
		if isInsert {
			// Column references are not allowed in INSERT ... SET values (ORA-00984).
			upper := strings.ToUpper(cr.Column)
			if oraclePseudoColumns[upper] || noParenBuiltins[upper] || sqlBuiltins[upper] || upper == "DEFAULT" {
				continue
			}
			if b.w.current != "" {
				if f, kind := b.look(walk.Key{Space: spaceObject, Parent: walk.Path(b.w.current), Name: cr.Column}); f == walk.Present && (kind == kindFunction || kind == kindSynonym || kind == kindPackage) {
					continue
				}
				if f, _ := b.look(publicKey(cr.Column)); f == walk.Present {
					continue
				}
			}
			b.fail(rng, review.UnsupportedForm, "Column reference "+review.Code(ident(cr.Column))+" is not allowed in INSERT SET.")
			return true
		}
		if cr.Table != "" {
			if f, _ := b.look(r.childKey(spaceColumn, cr.Column)); f == walk.Absent && !generated(spaceColumn, cr.Column) && !oraclePseudoColumns[strings.ToUpper(cr.Column)] {
				b.columnNeeds(r, cr.Column, rng)
				return true
			}
			continue
		}
		if !strings.EqualFold(cr.Column, "DEFAULT") && !b.isValidBareIdentifier(r, cr.Column) {
			b.columnNeeds(r, cr.Column, rng)
			return true
		}
	}
	return false
}
