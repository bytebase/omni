// Package parser - create_index.go implements T-SQL CREATE INDEX statement parsing.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateIndexStmt parses a CREATE INDEX statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-index-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-columnstore-index-transact-sql
//
//	createIndexStatement (TSql160.g:16831)
//	    : Create ( Unique )? ( ( Clustered | NonClustered ) ( createRelationalIndexStatement
//	      | createColumnStoreIndexStatement ) | createRelationalIndexStatement )
//	createRelationalIndexStatement (TSql160.g:16865)
//	    : Index identifier On schemaObjectThreePartName ( indexKeyColumnList )
//	      ( INCLUDE identifierColumnList )? ( filterClause )?
//	      ( With ( indexLegacyOptionList | indexOptionList ) )?
//	      ( On filegroupOrPartitionScheme )? fileStreamOnOpt
//	createColumnStoreIndexStatement (TSql160.g:16918)
//	    : COLUMNSTORE Index identifier On schemaObjectFourPartName
//	      ( ( columnstoreIndexKeyColumnList ) )? ( Order identifierColumnList )?
//	      ( filterClause )? ( With indexOptionList )? ( On filegroupOrPartitionScheme )?
//
// SQL Server refuses, while it parses, a columnstore index that is UNIQUE
// (Msg 35301), a clustered one with a key list (35335), a WHERE (102) or an
// ORDER on a nonclustered one (35342), a nonclustered or rowstore index
// without a key list (35336), and a sort order on a columnstore key (35302).
func (p *Parser) parseCreateIndexStmt(unique bool) (*nodes.CreateIndexStmt, error) {
	loc := p.pos()

	stmt := &nodes.CreateIndexStmt{
		Unique: unique,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// CLUSTERED / NONCLUSTERED
	if p.cur.Type == kwCLUSTERED {
		p.advance()
		v := true
		stmt.Clustered = &v
	} else if p.cur.Type == kwNONCLUSTERED {
		p.advance()
		v := false
		stmt.Clustered = &v
	}

	// COLUMNSTORE (optional)
	if p.cur.Type == kwCOLUMNSTORE {
		if unique {
			return nil, p.newParseError(p.cur.Loc, "a columnstore index cannot be unique")
		}
		p.advance()
		stmt.Columnstore = true
	}

	if _, err := p.expect(kwINDEX); err != nil {
		return nil, err
	}

	// Index name
	name, nameOK := p.parseIdentifier()
	if !nameOK {
		return nil, p.unexpectedToken()
	}
	stmt.Name = name

	// ON table
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	// Completion: after CREATE INDEX ... ON → table_ref
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, errCollecting
	}
	table, err := p.indexTargetTable(3)
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	if stmt.Columnstore {
		if err := p.indexColumnstoreBody(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	}

	// ( indexKeyColumnList )
	if p.cur.Type != '(' {
		return nil, p.unexpectedToken()
	}
	stmt.Columns, err = p.parseIndexColumnList()
	if err != nil {
		return nil, err
	}

	// INCLUDE ( column [ ,...n ] )
	if p.cur.Type == kwINCLUDE {
		p.advance()
		if p.cur.Type != '(' {
			return nil, p.unexpectedToken()
		}
		stmt.IncludeCols, err = p.indexIncludeColumnList()
		if err != nil {
			return nil, err
		}
	}

	// WHERE (filtered index)
	if p.cur.Type == kwWHERE {
		p.advance()
		stmt.WhereClause, err = p.indexFilterExpression()
		if err != nil {
			return nil, err
		}
	}

	// WITH ( <relational_index_option> [ ,...n ] ) | WITH <legacy option> [ ,...n ]
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			stmt.Options, err = p.indexOptionList(indexOptCreate)
		} else {
			stmt.Options, err = p.indexLegacyOptionList()
		}
		if err != nil {
			return nil, err
		}
	}

	// ON { partition_scheme_name ( column_name ) | filegroup_name | "default" }
	if p.cur.Type == kwON {
		p.advance()
		stmt.OnFileGroup, err = p.indexFilegroup(true)
		if err != nil {
			return nil, err
		}
	}

	// FILESTREAM_ON { filestream_filegroup_name | partition_scheme_name | "NULL" }
	if p.cur.Type == kwFILESTREAM_ON {
		p.advance()
		stmt.FilestreamOn, err = p.indexFilegroup(false)
		if err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// indexColumnstoreBody parses what follows ON <table> in CREATE COLUMNSTORE
// INDEX (TSql160.g:16918 and its actions at 16941 and 16964).
func (p *Parser) indexColumnstoreBody(stmt *nodes.CreateIndexStmt) error {
	clustered := stmt.Clustered != nil && *stmt.Clustered
	var err error

	// ( column [ ,...n ] ): required on a nonclustered columnstore index,
	// refused on a clustered one.
	if p.cur.Type == '(' {
		if clustered {
			return p.newParseError(p.cur.Loc, "a clustered columnstore index takes no key list")
		}
		p.advance()
		items, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
			if p.collectMode() {
				p.addRuleCandidate("columnref")
				return nil, errCollecting
			}
			loc := p.pos()
			name, err := p.parseIndexColumnName()
			if err != nil {
				return nil, err
			}
			return &nodes.IndexColumn{Name: name, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		})
		if err != nil {
			return err
		}
		if _, err := p.expect(')'); err != nil {
			return err
		}
		stmt.Columns = &nodes.List{Items: items}
	} else if !clustered {
		return p.newParseError(p.cur.Loc, "a nonclustered columnstore index needs a key list")
	}

	// ORDER ( column [ ,...n ] ): clustered columnstore only in SQL Server 2022.
	if p.cur.Type == kwORDER {
		if !clustered {
			return p.newParseError(p.cur.Loc, "ORDER is not valid on a nonclustered columnstore index")
		}
		p.advance()
		if p.cur.Type != '(' {
			return p.unexpectedToken()
		}
		stmt.OrderCols, err = p.parseParenIdentList()
		if err != nil {
			return err
		}
	}

	// WHERE <filter_expression>: nonclustered columnstore only.
	if p.cur.Type == kwWHERE {
		if clustered {
			return p.unexpectedToken()
		}
		p.advance()
		stmt.WhereClause, err = p.indexFilterExpression()
		if err != nil {
			return err
		}
	}

	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type != '(' {
			return p.unexpectedToken()
		}
		stmt.Options, err = p.indexOptionList(indexOptColumnstore)
		if err != nil {
			return err
		}
	}

	if p.cur.Type == kwON {
		p.advance()
		stmt.OnFileGroup, err = p.indexFilegroup(true)
		if err != nil {
			return err
		}
	}
	return nil
}

// indexTargetTable parses the object an index or statistics statement names
// after ON. SQL Server refuses a name with more than maxParts parts (Msg 117).
func (p *Parser) indexTargetTable(maxParts int) (*nodes.TableRef, error) {
	start := p.cur
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, p.unexpectedToken()
	}
	if indexTableRefParts(ref) > maxParts {
		return nil, p.newParseError(start.Loc, "the object name contains more than the maximum number of prefixes")
	}
	return ref, nil
}

// indexTableRefParts counts the parts of a parsed object name.
func indexTableRefParts(ref *nodes.TableRef) int {
	switch {
	case ref.Server != "":
		return 4
	case ref.Database != "":
		return 3
	case ref.Schema != "":
		return 2
	}
	return 1
}

// indexIncludeColumnList parses INCLUDE ( column [ ,...n ] ): one-part names
// or graph pseudo-columns (SQL Server refuses a dotted name, Msg 102).
func (p *Parser) indexIncludeColumnList() (*nodes.List, error) {
	p.advance() // consume (
	items, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, errCollecting
		}
		loc := p.pos()
		name, err := p.parseIndexColumnName()
		if err != nil {
			return nil, err
		}
		return &nodes.IndexColumn{Name: name, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// indexFilegroup parses filegroupOrPartitionScheme (TSql160.g:17224):
// stringOrIdentifier ( columnNameList )?. A bare DEFAULT or PRIMARY is a
// reserved word there; SQL Server takes [DEFAULT], "default" and [PRIMARY].
// FILESTREAM_ON takes no column list.
func (p *Parser) indexFilegroup(allowColumns bool) (string, error) {
	var name string
	switch {
	case p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST:
		name = p.cur.Str
		p.advance()
	default:
		id, ok := p.parseIdentifier()
		if !ok {
			return "", p.unexpectedToken()
		}
		name = id
	}
	if allowColumns && p.cur.Type == '(' {
		if _, err := p.parseParenIdentList(); err != nil {
			return "", err
		}
	}
	return name, nil
}

// parseIndexColumnList parses (col [ASC|DESC], ...).
func (p *Parser) parseIndexColumnList() (*nodes.List, error) {
	p.advance() // consume (
	items, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		// Completion: at start of each index column slot → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, errCollecting
		}
		loc := p.pos()
		name, err := p.parseIndexColumnName()
		if err != nil {
			return nil, err
		}
		dir := nodes.SortDefault
		if _, ok := p.match(kwASC); ok {
			dir = nodes.SortAsc
		} else if _, ok := p.match(kwDESC); ok {
			dir = nodes.SortDesc
		}
		return &nodes.IndexColumn{
			Name:    name,
			SortDir: dir,
			Loc:     nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// indexFilterExpression parses the WHERE of a filtered index or statistics.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-index-transact-sql#where-filter_predicate
//
//	filterExpression (TSql160.g:17113)
//	    : filterExpressionPrimary ( And filterExpression )?
//	filterExpressionPrimary (TSql160.g:17080)
//	    : filterColumn ( Is ( Not )? Null | filterComparisonPredicate | filterInPredicate )
//	    | LeftParenthesis filterExpression RightParenthesis
//	filterComparisonPredicate (TSql160.g:17176) : comparisonOperator expression
//	filterInPredicate (TSql160.g:17190) : In ( expression ( , expression )* )
//
// OR, NOT, NOT IN, BETWEEN and LIKE are refused (Msg 156). IS [NOT]
// DISTINCT FROM (TSql160.g:17080) is not read: the AST has no node for it.
func (p *Parser) indexFilterExpression() (nodes.ExprNode, error) {
	left, err := p.indexFilterPrimary()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == kwAND {
		loc := p.pos()
		p.advance()
		right, err := p.indexFilterPrimary()
		if err != nil {
			return nil, err
		}
		left = &nodes.BinaryExpr{
			Op:    nodes.BinOpAnd,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
		}
	}
	return left, nil
}

func (p *Parser) indexFilterPrimary() (nodes.ExprNode, error) {
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, errCollecting
	}
	loc := p.pos()
	if p.cur.Type == '(' {
		p.advance()
		inner, err := p.indexFilterExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.ParenExpr{Expr: inner, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	}

	// filterColumn: a column name of one to four parts, IDENTITYCOL,
	// ROWGUIDCOL, or a pseudo-column.
	if p.cur.Type != tokIDENT && p.cur.Type != tokPSEUDOCOL &&
		p.cur.Type != kwIDENTITYCOL && p.cur.Type != kwROWGUIDCOL && !isContextKeyword(p.cur.Type) {
		return nil, p.unexpectedToken()
	}
	col, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	if _, ok := col.(*nodes.ColumnRef); !ok {
		return nil, p.newParseError(loc, "a filter predicate must start with a column")
	}

	switch {
	case p.cur.Type == kwIS:
		p.advance()
		test := nodes.IsNull
		if _, ok := p.match(kwNOT); ok {
			test = nodes.IsNotNull
		}
		if _, err := p.expect(kwNULL); err != nil {
			return nil, err
		}
		return &nodes.IsExpr{Expr: col, TestType: test, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	case p.cur.Type == kwIN:
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		items, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
			e, err := p.parseAddition()
			if err != nil {
				return nil, err
			}
			if e == nil {
				return nil, p.unexpectedToken()
			}
			return e, nil
		})
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.InExpr{Expr: col, List: &nodes.List{Items: items}, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	}
	op, ok := comparisonOp(p.cur.Type)
	if !ok {
		return nil, p.unexpectedToken()
	}
	p.advance()
	right, err := p.parseAddition()
	if err != nil {
		return nil, err
	}
	if right == nil {
		return nil, p.unexpectedToken()
	}
	return &nodes.BinaryExpr{Op: op, Left: col, Right: right, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
}

// indexWordIs reports whether the current token is the option word `word`
// (upper case): its keyword token when the lexer has one, else an
// undelimited identifier with that text.
func (p *Parser) indexWordIs(word string) bool {
	if kw := lookupKeyword(word); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == word
}

// indexOptionWord returns the upper-case text of the current token when it
// can name an option (an undelimited identifier or a keyword), else "".
// SQL Server refuses a bracketed or quoted option name (Msg 102).
func (p *Parser) indexOptionWord() string {
	if p.cur.Type == tokIDENT {
		if p.isDelimited(p.cur) {
			return ""
		}
		return strings.ToUpper(p.cur.Str)
	}
	if p.cur.Type >= kwABSENT && p.cur.Str != "" {
		return strings.ToUpper(p.cur.Str)
	}
	return ""
}

// indexTokIsBinary reports whether t is a 0x... binary constant.
func indexTokIsBinary(t Token) bool {
	return (t.Type == tokICONST || t.Type == tokFCONST) && len(t.Str) >= 2 && t.Str[0] == '0' && (t.Str[1] == 'x' || t.Str[1] == 'X')
}

// indexExpectInteger consumes an unsigned integer literal and returns its text.
func (p *Parser) indexExpectInteger() (string, error) {
	if p.cur.Type != tokICONST || indexTokIsBinary(p.cur) {
		return "", p.unexpectedToken()
	}
	s := p.cur.Str
	p.advance()
	return s, nil
}

// indexExpectOnOff consumes ON or OFF and returns it.
func (p *Parser) indexExpectOnOff() (string, error) {
	switch p.cur.Type {
	case kwON:
		p.advance()
		return "ON", nil
	case kwOFF:
		p.advance()
		return "OFF", nil
	}
	return "", p.unexpectedToken()
}

// indexOptionKind names the statement an index option list belongs to. SQL
// Server checks the option names each one takes while it parses: an option
// it does not know there is Msg 155, one it knows but refuses there Msg 153,
// 35317 or 35318.
type indexOptionKind int

const (
	indexOptCreate           indexOptionKind = iota // CREATE [UNIQUE] [NON]CLUSTERED INDEX
	indexOptColumnstore                             // CREATE [NON]CLUSTERED COLUMNSTORE INDEX
	indexOptXML                                     // CREATE [PRIMARY | SELECTIVE] XML INDEX
	indexOptRebuild                                 // ALTER INDEX REBUILD [PARTITION = ALL]
	indexOptRebuildPartition                        // ALTER INDEX REBUILD PARTITION = n
	indexOptReorganize                              // ALTER INDEX REORGANIZE
	indexOptSet                                     // ALTER INDEX SET
	indexOptResume                                  // ALTER INDEX RESUME
	indexOptDrop                                    // DROP INDEX ... ON ... WITH
)

// indexOptionRule is what one statement kind takes in its WITH ( ... ).
type indexOptionRule struct {
	names             map[string]bool
	dataCompression   map[string]bool // DATA_COMPRESSION values
	onlineOn          bool            // ONLINE = ON is taken
	lowPriority       bool            // ONLINE = ON ( WAIT_AT_LOW_PRIORITY (...) ) is taken
	onPartitions      bool            // DATA_COMPRESSION / XML_COMPRESSION ... ON PARTITIONS (...)
	partitionsNeedAll bool            // ON PARTITIONS needs PARTITION = ALL (Msg 10737, 16209)
	once              bool            // an option given twice is refused (Msg 153)
}

func indexNameSet(names ...string) map[string]bool {
	m := make(map[string]bool, len(names))
	for _, n := range names {
		m[n] = true
	}
	return m
}

var (
	indexRowstoreCompression    = indexNameSet("NONE", "ROW", "PAGE")
	indexColumnstoreCompression = indexNameSet("COLUMNSTORE", "COLUMNSTORE_ARCHIVE")
	indexAnyCompression         = indexNameSet("NONE", "ROW", "PAGE", "COLUMNSTORE", "COLUMNSTORE_ARCHIVE")
)

// indexOptionRules holds, per statement kind, the options SQL Server 2022
// takes. indexOption (TSql160.g:17601) lists every option for every index
// statement and leaves the choice to VerifyAllowedIndexOption160; these sets
// are the pinned server's answers to each option in each statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-index-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-index-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-index-transact-sql
var indexOptionRules = map[indexOptionKind]indexOptionRule{
	indexOptCreate: {
		names: indexNameSet("PAD_INDEX", "FILLFACTOR", "SORT_IN_TEMPDB", "IGNORE_DUP_KEY",
			"STATISTICS_NORECOMPUTE", "STATISTICS_INCREMENTAL", "DROP_EXISTING", "ONLINE",
			"RESUMABLE", "MAX_DURATION", "ALLOW_ROW_LOCKS", "ALLOW_PAGE_LOCKS",
			"OPTIMIZE_FOR_SEQUENTIAL_KEY", "MAXDOP", "DATA_COMPRESSION", "XML_COMPRESSION"),
		dataCompression: indexRowstoreCompression,
		onlineOn:        true, lowPriority: true, onPartitions: true,
	},
	indexOptColumnstore: {
		names:           indexNameSet("DROP_EXISTING", "ONLINE", "RESUMABLE", "MAXDOP", "DATA_COMPRESSION", "COMPRESSION_DELAY"),
		dataCompression: indexColumnstoreCompression,
		onlineOn:        true, lowPriority: true, onPartitions: true,
	},
	indexOptXML: {
		names: indexNameSet("PAD_INDEX", "FILLFACTOR", "SORT_IN_TEMPDB", "STATISTICS_NORECOMPUTE",
			"DROP_EXISTING", "ONLINE", "RESUMABLE", "ALLOW_ROW_LOCKS", "ALLOW_PAGE_LOCKS",
			"MAXDOP", "XML_COMPRESSION", "COMPRESSION_DELAY", "COMPRESS_ALL_ROW_GROUPS"),
		onPartitions: true,
	},
	indexOptRebuild: {
		names: indexNameSet("PAD_INDEX", "FILLFACTOR", "SORT_IN_TEMPDB", "IGNORE_DUP_KEY",
			"STATISTICS_NORECOMPUTE", "STATISTICS_INCREMENTAL", "ONLINE", "RESUMABLE",
			"MAX_DURATION", "ALLOW_ROW_LOCKS", "ALLOW_PAGE_LOCKS", "MAXDOP",
			"DATA_COMPRESSION", "XML_COMPRESSION", "BUCKET_COUNT"),
		dataCompression: indexAnyCompression,
		onlineOn:        true, lowPriority: true, onPartitions: true,
	},
	indexOptRebuildPartition: {
		names: indexNameSet("SORT_IN_TEMPDB", "ONLINE", "RESUMABLE", "MAX_DURATION", "MAXDOP",
			"DATA_COMPRESSION", "XML_COMPRESSION"),
		dataCompression: indexAnyCompression,
		onlineOn:        true, lowPriority: true, onPartitions: true,
	},
	indexOptReorganize: {
		names: indexNameSet("LOB_COMPACTION", "COMPRESS_ALL_ROW_GROUPS"),
	},
	indexOptSet: {
		names: indexNameSet("IGNORE_DUP_KEY", "STATISTICS_NORECOMPUTE", "ALLOW_ROW_LOCKS",
			"ALLOW_PAGE_LOCKS", "OPTIMIZE_FOR_SEQUENTIAL_KEY", "COMPRESSION_DELAY"),
	},
	indexOptResume: {
		names: indexNameSet("MAX_DURATION", "MAXDOP", "WAIT_AT_LOW_PRIORITY"),
	},
	indexOptDrop: {
		names: indexNameSet("ONLINE", "MAXDOP", "DATA_COMPRESSION", "XML_COMPRESSION",
			"WAIT_AT_LOW_PRIORITY", "MOVE", "FILESTREAM_ON"),
		dataCompression: indexAnyCompression,
		onlineOn:        true, onPartitions: true, once: true,
	},
}

// indexOnOffOptions take { ON | OFF } (indexStateOption, TSql160.g:17537).
var indexOnOffOptions = indexNameSet("PAD_INDEX", "SORT_IN_TEMPDB", "STATISTICS_NORECOMPUTE",
	"STATISTICS_INCREMENTAL", "DROP_EXISTING", "RESUMABLE", "ALLOW_ROW_LOCKS", "ALLOW_PAGE_LOCKS",
	"OPTIMIZE_FOR_SEQUENTIAL_KEY", "LOB_COMPACTION", "COMPRESS_ALL_ROW_GROUPS")

// indexOptionList parses ( indexOption [ , ...n ] ) (indexOptionList,
// TSql160.g:17282) for the statement kind given. Each option becomes a
// String "NAME=VALUE", the shape mssql/review reads DROP_EXISTING from.
func (p *Parser) indexOptionList(kind indexOptionKind) (*nodes.List, error) {
	return p.indexOptionListRule(indexOptionRules[kind])
}

func (p *Parser) indexOptionListRule(rule indexOptionRule) (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	items, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		word := p.indexOptionWord()
		if rule.once && seen[word] {
			return nil, p.newParseError(p.cur.Loc, fmt.Sprintf("option %s is given more than once", word))
		}
		seen[word] = true
		return p.indexOption(rule)
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// indexOption parses one option of an index option list.
//
//	indexOption (TSql160.g:17601)
//	    : FILLFACTOR = integer | MAXDOP = integer | BUCKET_COUNT = integer
//	    | DATA_COMPRESSION = Identifier ( On PARTITIONS ( compressionPartitionRange , ... ) )?
//	    | XML_COMPRESSION = ( On | Off ) ( On PARTITIONS ( ... ) )?
//	    | ONLINE = ( On ( ( lowPriorityLockWaitOption ) )? | Off )
//	    | COMPRESSION_DELAY = integer ( MINUTES )?
//	    | IGNORE_DUP_KEY = ( On | Off ) ( ( SUPPRESS_MESSAGES = ( On | Off ) ) )?
//	    | MAX_DURATION = integer ( MINUTES )?
//	    | WAIT_AT_LOW_PRIORITY ( MAX_DURATION = integer ( MINUTES )? , ABORT_AFTER_WAIT = Identifier )
//	    | indexStateOption
//
// DROP INDEX adds MOVE TO and FILESTREAM_ON (dropClusteredIndexOption,
// TSql160.g:14977). SQL Server also takes SUPPRESS_MESSAGES after
// IGNORE_DUP_KEY = OFF, which ScriptDOM refuses.
func (p *Parser) indexOption(rule indexOptionRule) (nodes.Node, error) {
	if p.collectMode() {
		return nil, errCollecting
	}
	word := p.indexOptionWord()
	if word == "" || !rule.names[word] {
		return nil, p.unexpectedToken()
	}
	p.advance()

	var value string
	var err error
	switch word {
	case "MOVE":
		// MOVE TO { partition_scheme_name ( column_name ) | filegroup_name | "default" }
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		start := p.pos()
		if _, err := p.indexFilegroup(true); err != nil {
			return nil, err
		}
		return &nodes.String{Str: "MOVE TO=" + p.source[start:p.prevEnd()]}, nil
	case "FILESTREAM_ON":
		value, err = p.indexFilegroup(false)
		if err != nil {
			return nil, err
		}
		return &nodes.String{Str: word + "=" + value}, nil
	case "WAIT_AT_LOW_PRIORITY":
		start := p.pos()
		if err := p.indexLowPriorityBody(); err != nil {
			return nil, err
		}
		return &nodes.String{Str: word + "=" + p.source[start:p.prevEnd()]}, nil
	}

	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	switch word {
	case "FILLFACTOR":
		value, err = p.indexExpectInteger()
		if err == nil {
			if n, convErr := strconv.Atoi(value); convErr != nil || n < 1 || n > 100 {
				// Msg 129: fillfactor must be between 1 and 100.
				return nil, p.newParseError(p.prev.Loc, fmt.Sprintf("fillfactor %s is not between 1 and 100", value))
			}
		}
	case "MAXDOP", "BUCKET_COUNT":
		value, err = p.indexExpectInteger()
	case "MAX_DURATION", "COMPRESSION_DELAY":
		value, err = p.indexExpectInteger()
		if err == nil && p.cur.Type == kwMINUTES {
			p.advance()
		}
	case "DATA_COMPRESSION":
		v := p.indexOptionWord()
		if v == "" || !rule.dataCompression[v] {
			return nil, p.unexpectedToken()
		}
		p.advance()
		value = v
		err = p.indexOnPartitions(rule)
	case "XML_COMPRESSION":
		value, err = p.indexExpectOnOff()
		if err == nil {
			err = p.indexOnPartitions(rule)
		}
	case "ONLINE":
		onTok := p.cur
		value, err = p.indexExpectOnOff()
		if err == nil && value == "ON" {
			if !rule.onlineOn {
				// Msg 153: ONLINE = ON is not valid for this index.
				return nil, p.newParseError(onTok.Loc, "ONLINE = ON is not valid here")
			}
			if p.cur.Type == '(' && rule.lowPriority {
				p.advance()
				if !p.indexWordIs("WAIT_AT_LOW_PRIORITY") {
					return nil, p.unexpectedToken()
				}
				p.advance()
				if err := p.indexLowPriorityBody(); err != nil {
					return nil, err
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
			}
		}
	case "IGNORE_DUP_KEY":
		value, err = p.indexExpectOnOff()
		if err == nil && p.cur.Type == '(' {
			p.advance()
			if !p.indexWordIs("SUPPRESS_MESSAGES") {
				return nil, p.unexpectedToken()
			}
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			if _, err := p.indexExpectOnOff(); err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
	default:
		if !indexOnOffOptions[word] {
			return nil, p.unexpectedToken()
		}
		value, err = p.indexExpectOnOff()
	}
	if err != nil {
		return nil, err
	}
	return &nodes.String{Str: word + "=" + value}, nil
}

// indexOnPartitions parses the optional ON PARTITIONS ( range [ , ...n ] )
// after DATA_COMPRESSION or XML_COMPRESSION, compressionPartitionRange
// (TSql160.g:27623): expression ( TO expression )?.
func (p *Parser) indexOnPartitions(rule indexOptionRule) error {
	if p.cur.Type != kwON || p.peekNext().Type != kwPARTITIONS {
		return nil
	}
	if !rule.onPartitions || rule.partitionsNeedAll {
		// Msg 10737 / 16209: ON PARTITIONS needs PARTITION = ALL.
		return p.unexpectedToken()
	}
	p.advance() // ON
	p.advance() // PARTITIONS
	if _, err := p.expect('('); err != nil {
		return err
	}
	_, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		from, err := p.parseAddition()
		if err != nil {
			return nil, err
		}
		if from == nil {
			return nil, p.unexpectedToken()
		}
		if _, ok := p.match(kwTO); ok {
			to, err := p.parseAddition()
			if err != nil {
				return nil, err
			}
			if to == nil {
				return nil, p.unexpectedToken()
			}
		}
		return from, nil
	})
	if err != nil {
		return err
	}
	_, err = p.expect(')')
	return err
}

// indexLowPriorityBody parses the parenthesized part of
// lowPriorityLockWaitOption (TSql160.g:17417), after WAIT_AT_LOW_PRIORITY:
// ( MAX_DURATION = integer [ MINUTES ] , ABORT_AFTER_WAIT = { NONE | SELF | BLOCKERS } ).
func (p *Parser) indexLowPriorityBody() error {
	if _, err := p.expect('('); err != nil {
		return err
	}
	if !p.indexWordIs("MAX_DURATION") {
		return p.unexpectedToken()
	}
	p.advance()
	if _, err := p.expect('='); err != nil {
		return err
	}
	if _, err := p.indexExpectInteger(); err != nil {
		return err
	}
	if p.cur.Type == kwMINUTES {
		p.advance()
	}
	if _, err := p.expect(','); err != nil {
		return err
	}
	if !p.indexWordIs("ABORT_AFTER_WAIT") {
		return p.unexpectedToken()
	}
	p.advance()
	if _, err := p.expect('='); err != nil {
		return err
	}
	if !p.indexWordIs("NONE") && !p.indexWordIs("SELF") && !p.indexWordIs("BLOCKERS") {
		return p.unexpectedToken()
	}
	p.advance()
	_, err := p.expect(')')
	return err
}

// indexLegacyOptionList parses the unparenthesized WITH list of CREATE INDEX
// (indexLegacyOptionList, TSql160.g:17238): FILLFACTOR = integer, or one of
// the words SQL Server 2022 still takes alone, each of which turns its
// option on (Msg 153 for any other word).
func (p *Parser) indexLegacyOptionList() (*nodes.List, error) {
	items, err := p.parseCommaList(tokEOF, commaListStrict, func() (nodes.Node, error) {
		if p.collectMode() {
			return nil, errCollecting
		}
		word := p.indexOptionWord()
		switch word {
		case "FILLFACTOR":
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			v, err := p.indexExpectInteger()
			if err != nil {
				return nil, err
			}
			return &nodes.String{Str: "FILLFACTOR=" + v}, nil
		case "PAD_INDEX", "SORT_IN_TEMPDB", "IGNORE_DUP_KEY", "DROP_EXISTING", "STATISTICS_NORECOMPUTE":
			p.advance()
			return &nodes.String{Str: word + "=ON"}, nil
		}
		return nil, p.unexpectedToken()
	})
	if err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// indexXmlColumn parses ( xml_column_name ): exactly one identifier.
func (p *Parser) indexXmlColumn() (string, error) {
	if _, err := p.expect('('); err != nil {
		return "", err
	}
	col, ok := p.parseIdentifier()
	if !ok {
		return "", p.unexpectedToken()
	}
	if _, err := p.expect(')'); err != nil {
		return "", err
	}
	return col, nil
}

// parseCreateXmlIndexStmt parses CREATE [PRIMARY] XML INDEX.
// Caller has consumed CREATE [PRIMARY] XML INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-xml-index-transact-sql
//
//	createXmlIndexStatement (TSql160.g:9053)
//	    : Index identifier On schemaObjectThreePartName ( identifier )
//	      ( | USING XML Index identifier ( For Identifier | For ( identifier ) ) )
//	      ( With indexOptionList )?
//
// A primary XML index takes no USING (the vPrimary predicate); SQL Server
// refuses a secondary one without it (Msg 6343) and a FOR word other than
// VALUE, PATH or PROPERTY (Msg 487).
func (p *Parser) parseCreateXmlIndexStmt(primary bool) (*nodes.CreateXmlIndexStmt, error) {
	// The caller consumes INDEX only when present.
	if p.prev.Type != kwINDEX {
		return nil, p.unexpectedToken()
	}
	loc := p.pos()
	stmt := &nodes.CreateXmlIndexStmt{
		Primary: primary,
		Loc:     nodes.Loc{Start: loc, End: -1},
	}

	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.unexpectedToken()
	}
	stmt.Name = name

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	table, err := p.indexTargetTable(4)
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	stmt.XmlColumn, err = p.indexXmlColumn()
	if err != nil {
		return nil, err
	}

	if !primary {
		// USING XML INDEX parent_index_name FOR { VALUE | PATH | PROPERTY | ( path_name ) }
		if _, err := p.expect(kwUSING); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwXML); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwINDEX); err != nil {
			return nil, err
		}
		idx, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		stmt.UsingIndex = idx
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		if p.cur.Type == '(' {
			// secondarySelectiveXmlIndex (TSql160.g:9112): FOR ( identifier )
			p.advance()
			path, ok := p.parseIdentifier()
			if !ok {
				return nil, p.unexpectedToken()
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			stmt.SecondaryFor = "(" + path + ")"
		} else {
			word := p.indexOptionWord()
			if word != "VALUE" && word != "PATH" && word != "PROPERTY" {
				return nil, p.unexpectedToken()
			}
			p.advance()
			stmt.SecondaryFor = word
		}
	}

	if p.cur.Type == kwWITH {
		p.advance()
		stmt.Options, err = p.indexOptionList(indexOptXML)
		if err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateSelectiveXmlIndexStmt parses CREATE SELECTIVE XML INDEX.
// Caller has consumed CREATE SELECTIVE XML INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-selective-xml-index-transact-sql
//
//	createSelectiveXmlIndexStatement (TSql160.g:8962)
//	    : SELECTIVE XML Index identifier On schemaObjectThreePartName ( identifier )
//	      ( With xmlNamespaces )? For ( promotedSelectiveXmlIndexPath ( , ... )* )
//	      ( With indexOptionList )?
func (p *Parser) parseCreateSelectiveXmlIndexStmt() (*nodes.CreateSelectiveXmlIndexStmt, error) {
	// The caller consumes INDEX only when present.
	if p.prev.Type != kwINDEX {
		return nil, p.unexpectedToken()
	}
	loc := p.pos()
	stmt := &nodes.CreateSelectiveXmlIndexStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.unexpectedToken()
	}
	stmt.Name = name

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	table, err := p.indexTargetTable(4)
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	stmt.XmlColumn, err = p.indexXmlColumn()
	if err != nil {
		return nil, err
	}

	// WITH XMLNAMESPACES (...)
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type != kwXMLNAMESPACES {
			return nil, p.unexpectedToken()
		}
		stmt.Namespaces, err = p.indexXmlNamespaces()
		if err != nil {
			return nil, err
		}
	}

	// FOR ( promotedSelectiveXmlIndexPath [ , ...n ] )
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	paths, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		return p.indexSelectivePath()
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	stmt.Paths = &nodes.List{Items: paths}

	if p.cur.Type == kwWITH {
		p.advance()
		stmt.Options, err = p.indexOptionList(indexOptXML)
		if err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// indexXmlNamespaces parses xmlNamespaces (TSql160.g:17703), at XMLNAMESPACES:
//
//	XMLNAMESPACES ( { DEFAULT stringLiteral | stringLiteral AS identifier } [ , ...n ] )
func (p *Parser) indexXmlNamespaces() (*nodes.List, error) {
	p.advance() // XMLNAMESPACES
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Msg 6869: a prefix, or the default namespace, is declared once.
	prefixes := map[string]bool{}
	decls, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		decl := &nodes.XmlNamespaceDecl{Loc: nodes.Loc{Start: p.pos(), End: -1}}
		if _, ok := p.match(kwDEFAULT); ok {
			decl.IsDefault = true
		}
		if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
			return nil, p.unexpectedToken()
		}
		decl.URI = p.cur.Str
		p.advance()
		if !decl.IsDefault {
			if _, err := p.expect(kwAS); err != nil {
				return nil, err
			}
			prefix, ok := p.parseIdentifier()
			if !ok {
				return nil, p.unexpectedToken()
			}
			decl.Prefix = prefix
		}
		key := "\x00default"
		if !decl.IsDefault {
			key = decl.Prefix
		}
		if prefixes[key] {
			return nil, p.newParseError(decl.Loc.Start, "namespace prefix declared twice")
		}
		prefixes[key] = true
		decl.Loc.End = p.prevEnd()
		return decl, nil
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: decls}, nil
}

// indexSelectivePath parses one promoted path of a selective XML index and
// returns its text.
//
//	promotedSelectiveXmlIndexPath (TSql160.g:9133)
//	    : identifier = stringLiteral
//	      ( As ( XQUERY ( stringLiteral )? ( MAXLENGTH ( integer ) )? | SQL scalarDataType )
//	        ( SINGLETON )? )?
//
// AS XQUERY with neither a type, MAXLENGTH nor SINGLETON is refused (the
// action at TSql160.g:9179).
func (p *Parser) indexSelectivePath() (nodes.Node, error) {
	start := p.pos()
	if _, ok := p.parseIdentifier(); !ok {
		return nil, p.unexpectedToken()
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
		return nil, p.unexpectedToken()
	}
	p.advance()
	if _, ok := p.match(kwAS); ok {
		switch {
		case p.indexWordIs("XQUERY"):
			p.advance()
			typed := false
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				p.advance()
				typed = true
			}
			if p.indexWordIs("MAXLENGTH") {
				p.advance()
				if _, err := p.expect('('); err != nil {
					return nil, err
				}
				if _, err := p.indexExpectInteger(); err != nil {
					return nil, err
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				typed = true
			}
			if p.indexWordIs("SINGLETON") {
				p.advance()
				typed = true
			}
			if !typed {
				return nil, p.unexpectedToken()
			}
		case p.indexWordIs("SQL"):
			p.advance()
			dt, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			if dt == nil {
				return nil, p.unexpectedToken()
			}
			if p.indexWordIs("SINGLETON") {
				p.advance()
			}
		default:
			return nil, p.unexpectedToken()
		}
	}
	return &nodes.String{Str: p.source[start:p.prevEnd()]}, nil
}

// parseCreateSpatialIndexStmt parses CREATE SPATIAL INDEX.
// Caller has consumed CREATE SPATIAL INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-spatial-index-transact-sql
//
//	createSpatialIndexStatement (TSql160.g:24267)
//	    : SPATIAL Index identifier On schemaObjectFourPartName ( identifier )
//	      ( USING Identifier )? ( With spatialIndexOptionsList )? ( On stringOrIdentifier )?
//
// SQL Server refuses a four-part table name here (Msg 117).
func (p *Parser) parseCreateSpatialIndexStmt() (*nodes.CreateSpatialIndexStmt, error) {
	// The caller consumes INDEX only when present.
	if p.prev.Type != kwINDEX {
		return nil, p.unexpectedToken()
	}
	loc := p.pos()
	stmt := &nodes.CreateSpatialIndexStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.unexpectedToken()
	}
	stmt.Name = name

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	table, err := p.indexTargetTable(3)
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	// (spatial_column)
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	colLoc := p.pos()
	col, ok := p.parseIdentifier()
	if !ok {
		return nil, p.unexpectedToken()
	}
	stmt.SpatialColumn = col
	stmt.SpatialColumnRef = &nodes.ColumnRef{
		Column: col,
		Loc:    nodes.Loc{Start: colLoc, End: p.prevEnd()},
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// USING tessellation_scheme
	if p.cur.Type == kwUSING {
		p.advance()
		scheme, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		stmt.Using = scheme
	}

	if p.cur.Type == kwWITH {
		p.advance()
		stmt.Options, err = p.indexSpatialOptionList()
		if err != nil {
			return nil, err
		}
	}

	// ON { filegroup_name | "default" }
	if p.cur.Type == kwON {
		p.advance()
		stmt.OnFileGroup, err = p.indexFilegroup(false)
		if err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// indexSpatialOnOff are the options a spatial index takes as { ON | OFF };
// ONLINE, IGNORE_DUP_KEY and OPTIMIZE_FOR_SEQUENTIAL_KEY take only OFF
// (Msg 153 for ON). Any other name with ON or OFF is Msg 155.
var indexSpatialOnOff = indexNameSet("PAD_INDEX", "SORT_IN_TEMPDB", "STATISTICS_NORECOMPUTE",
	"DROP_EXISTING", "ONLINE", "RESUMABLE", "ALLOW_ROW_LOCKS", "ALLOW_PAGE_LOCKS",
	"COMPRESS_ALL_ROW_GROUPS", "IGNORE_DUP_KEY", "OPTIMIZE_FOR_SEQUENTIAL_KEY")

var indexSpatialOffOnly = indexNameSet("ONLINE", "IGNORE_DUP_KEY", "OPTIMIZE_FOR_SEQUENTIAL_KEY")

// indexSpatialOptionList parses spatialIndexOptionsList (TSql160.g:24319):
// ( spatialIndexOption [ , ...n ] ).
//
// SQL Server reads a spatial index option more loosely than ScriptDOM's
// spatialIndexOption (TSql160.g:24341): any name = a constant, a name or a
// parenthesized list of them parses (BOUNDING_BOX, GRIDS and CELLS_PER_OBJECT
// among them); FILLFACTOR, MAXDOP and DATA_COMPRESSION keep their own
// values, and a name given ON or OFF must be one it knows.
func (p *Parser) indexSpatialOptionList() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	items, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		return p.indexSpatialOption()
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) indexSpatialOption() (nodes.Node, error) {
	if p.collectMode() {
		return nil, errCollecting
	}
	if !p.indexSpatialNameToken() && p.cur.Type != kwFILLFACTOR {
		return nil, p.unexpectedToken()
	}
	word := strings.ToUpper(p.cur.Str)
	p.advance()
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	start := p.pos()
	var value string
	switch {
	case word == "FILLFACTOR":
		v, err := p.indexExpectInteger()
		if err != nil {
			return nil, err
		}
		if n, convErr := strconv.Atoi(v); convErr != nil || n < 1 || n > 100 {
			return nil, p.newParseError(p.prev.Loc, fmt.Sprintf("fillfactor %s is not between 1 and 100", v))
		}
		value = v
	case word == "MAXDOP":
		v, err := p.indexExpectInteger()
		if err != nil {
			return nil, err
		}
		value = v
	case word == "DATA_COMPRESSION":
		v := p.indexOptionWord()
		if v == "" || !indexAnyCompression[v] {
			return nil, p.unexpectedToken()
		}
		p.advance()
		value = v
	case word == "XML_COMPRESSION":
		// Msg 153 / 155: a spatial index takes no XML compression.
		return nil, p.newParseError(start, "XML_COMPRESSION is not valid on a spatial index")
	case p.cur.Type == kwON || p.cur.Type == kwOFF:
		if !indexSpatialOnOff[word] || (p.cur.Type == kwON && indexSpatialOffOnly[word]) {
			return nil, p.newParseError(start, fmt.Sprintf("option %s cannot be %s on a spatial index", word, strings.ToUpper(p.cur.Str)))
		}
		value = strings.ToUpper(p.cur.Str)
		p.advance()
	case p.cur.Type == '(':
		if err := p.indexSpatialValueList(); err != nil {
			return nil, err
		}
		value = p.source[start:p.prevEnd()]
	default:
		if err := p.indexSpatialScalar(); err != nil {
			return nil, err
		}
		value = p.source[start:p.prevEnd()]
		if p.prev.Type == tokIDENT || isContextKeyword(p.prev.Type) {
			value = strings.ToUpper(p.prev.Str)
		}
	}
	return &nodes.String{Str: word + "=" + value}, nil
}

// indexSpatialNameToken reports whether the current token can name a spatial
// index option or a name inside a list value: an undelimited identifier or a
// non-reserved keyword.
func (p *Parser) indexSpatialNameToken() bool {
	if p.cur.Type == tokIDENT {
		return !p.isDelimited(p.cur)
	}
	return isContextKeyword(p.cur.Type)
}

// indexSpatialScalar consumes one constant or name: [-] number, a string, an
// identifier, or NULL.
func (p *Parser) indexSpatialScalar() error {
	if p.cur.Type == '-' {
		p.advance()
		switch p.cur.Type {
		case tokICONST, tokFCONST, tokMONEY:
			p.advance()
			return nil
		}
		return p.unexpectedToken()
	}
	switch {
	case p.cur.Type == tokICONST, p.cur.Type == tokFCONST, p.cur.Type == tokMONEY,
		p.cur.Type == tokSCONST, p.cur.Type == tokNSCONST, p.cur.Type == kwNULL,
		p.cur.Type == tokIDENT, isContextKeyword(p.cur.Type):
		p.advance()
		return nil
	}
	return p.unexpectedToken()
}

// indexSpatialValueList consumes ( value [ , ...n ] ), whose items are all
// constants (boundingBoxParameterByOrd, gridParameterByOrd) or all
// name = constant (boundingBoxParameterByName, gridParameterByName).
func (p *Parser) indexSpatialValueList() error {
	p.advance() // (
	named := p.indexSpatialNameToken() && p.peekNext().Type == '='
	_, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		start := p.pos()
		if named {
			if !p.indexSpatialNameToken() {
				return nil, p.unexpectedToken()
			}
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
		}
		if err := p.indexSpatialScalar(); err != nil {
			return nil, err
		}
		return &nodes.String{Str: p.source[start:p.prevEnd()]}, nil
	})
	if err != nil {
		return err
	}
	_, err = p.expect(')')
	return err
}

// parseAlterIndexStmt parses an ALTER INDEX statement. Caller has consumed
// ALTER INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-index-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-index-selective-xml-indexes
//
//	alterIndexStatement (TSql160.g:11812)
//	    : Index ( identifier | All ) On schemaObjectFourPartName
//	      ( Set indexOptionList | DISABLE | PAUSE | ABORT
//	      | Identifier ( partitionSpecifier )? ( With indexOptionList )?
//	      | ( With xmlNamespaces )? For ( promotedSelectiveXmlIndexPathInAlter ( , ... )* ) )
//	partitionSpecifier (TSql160.g:11916) : PARTITION = ( expression | All )
//
// The Identifier action is REBUILD, REORGANIZE or RESUME; RESUME takes no
// PARTITION (Msg 102). A selective XML index path change is kept in Options
// as "ADD <path>" or "REMOVE <name>" with Action "FOR".
func (p *Parser) parseAlterIndexStmt() (*nodes.AlterIndexStmt, error) {
	loc := p.pos()
	stmt := &nodes.AlterIndexStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if p.cur.Type == kwALL {
		stmt.IndexName = p.cur.Str
		p.advance()
	} else {
		name, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		stmt.IndexName = name
	}

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	table, err := p.indexTargetTable(4)
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	switch {
	case p.cur.Type == kwSET:
		p.advance()
		stmt.Action = "SET"
		stmt.Options, err = p.indexOptionList(indexOptSet)
		if err != nil {
			return nil, err
		}
	case p.indexWordIs("DISABLE"), p.indexWordIs("PAUSE"), p.indexWordIs("ABORT"):
		stmt.Action = strings.ToUpper(p.cur.Str)
		p.advance()
	case p.indexWordIs("REBUILD"), p.indexWordIs("REORGANIZE"), p.indexWordIs("RESUME"):
		stmt.Action = strings.ToUpper(p.cur.Str)
		p.advance()
		if stmt.Action != "RESUME" && p.cur.Type == kwPARTITION {
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			if p.cur.Type == kwALL {
				stmt.Partition = "ALL"
				p.advance()
			} else {
				start := p.pos()
				e, err := p.parseAddition()
				if err != nil {
					return nil, err
				}
				if e == nil {
					return nil, p.unexpectedToken()
				}
				stmt.Partition = p.source[start:p.prevEnd()]
			}
		}
		if p.cur.Type == kwWITH {
			p.advance()
			var rule indexOptionRule
			switch {
			case stmt.Action == "REORGANIZE":
				rule = indexOptionRules[indexOptReorganize]
			case stmt.Action == "RESUME":
				rule = indexOptionRules[indexOptResume]
			case stmt.Partition == "":
				rule = indexOptionRules[indexOptRebuild]
				rule.partitionsNeedAll = true
			case stmt.Partition == "ALL":
				rule = indexOptionRules[indexOptRebuild]
			default:
				rule = indexOptionRules[indexOptRebuildPartition]
			}
			stmt.Options, err = p.indexOptionListRule(rule)
			if err != nil {
				return nil, err
			}
		}
	case p.cur.Type == kwWITH || p.cur.Type == kwFOR:
		// [ WITH XMLNAMESPACES ( ... ) ] FOR ( { ADD path | REMOVE name } [ , ...n ] )
		if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type != kwXMLNAMESPACES {
				return nil, p.unexpectedToken()
			}
			if _, err := p.indexXmlNamespaces(); err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		stmt.Action = "FOR"
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		// promotedSelectiveXmlIndexPathInAlter (TSql160.g:9220)
		items, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
			switch {
			case p.cur.Type == kwADD:
				p.advance()
				path, err := p.indexSelectivePath()
				if err != nil {
					return nil, err
				}
				return &nodes.String{Str: "ADD " + path.(*nodes.String).Str}, nil
			case p.indexWordIs("REMOVE"):
				p.advance()
				name, ok := p.parseIdentifier()
				if !ok {
					return nil, p.unexpectedToken()
				}
				return &nodes.String{Str: "REMOVE " + name}, nil
			}
			return nil, p.unexpectedToken()
		})
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.Options = &nodes.List{Items: items}
	default:
		return nil, p.unexpectedToken()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropIndexClauses parses the clauses of DROP INDEX [ IF EXISTS ] into
// stmt: Names, OnTables (nil for the table.index form) and the options of
// every clause.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-index-transact-sql
//
//	dropIndexClause (TSql160.g:14948) : indexDropObject | indexDropObjectNewNameFormat
//	indexDropObjectNewNameFormat (TSql160.g:14953)
//	    : identifier On schemaObjectThreePartName ( With ( dropClusteredIndexOption , ... ) )?
//	indexDropObject (TSql160.g:15176) : [ schema . ] table . index
//
// The old form needs the table (SQL46027) and takes no database (Msg 166).
func (p *Parser) parseDropIndexClauses(stmt *nodes.DropStmt) error {
	var names, onTables, opts []nodes.Node
	for {
		start := p.cur
		name, err := p.parseObjectRef(dropObjectCompletionRule(nodes.DropIndex))
		if err != nil {
			return err
		}
		if name == nil {
			return p.unexpectedToken()
		}
		var onTable nodes.Node
		if p.cur.Type == kwON {
			if indexTableRefParts(name) != 1 {
				return p.unexpectedToken()
			}
			p.advance()
			ref, err := p.indexTargetTable(3)
			if err != nil {
				return err
			}
			onTable = ref
			if p.cur.Type == kwWITH {
				p.advance()
				list, err := p.indexOptionList(indexOptDrop)
				if err != nil {
					return err
				}
				opts = append(opts, list.Items...)
			}
		} else if parts := indexTableRefParts(name); parts < 2 || parts > 3 {
			return p.newParseError(start.Loc, "DROP INDEX needs table.index or index ON table")
		}
		names = append(names, name)
		onTables = append(onTables, onTable)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	stmt.Names = &nodes.List{Items: names}
	stmt.OnTables = &nodes.List{Items: onTables}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	return nil
}

// parseCreateAggregateStmt parses CREATE AGGREGATE.
// Caller has consumed CREATE AGGREGATE.
//
// BNF: mssql/parser/bnf/create-aggregate-transact-sql.bnf
//
//	CREATE AGGREGATE [ schema_name . ] aggregate_name
//	        (@param_name <input_sqltype>
//	        [ ,...n ] )
//	RETURNS <return_sqltype>
//	EXTERNAL NAME assembly_name [ .class_name ]
//
//	<input_sqltype> ::=
//	        system_scalar_type | { [ udt_schema_name. ] udt_type_name }
//
//	<return_sqltype> ::=
//	        system_scalar_type | { [ udt_schema_name. ] udt_type_name }
func (p *Parser) parseCreateAggregateStmt() (*nodes.CreateAggregateStmt, error) {
	loc := p.pos()
	stmt := &nodes.CreateAggregateStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// [schema.]aggregate_name
	var err error
	stmt.Name, err = p.parseTableRef()
	if err != nil {
		return nil, err
	}

	// (@param_name type, ...)
	if p.cur.Type == '(' {
		p.advance()
		var params []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			param := &nodes.ParamDef{Loc: nodes.Loc{Start: p.pos(), End: -1}}
			// @param_name
			if p.cur.Type == tokVARIABLE {
				param.Name = p.cur.Str
				p.advance()
			}
			// data type
			param.DataType, err = p.parseDataType()
			if err != nil {
				return nil, err
			}
			params = append(params, param)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
		stmt.Params = &nodes.List{Items: params}
	}

	// RETURNS return_type
	if p.cur.Type == kwRETURNS {
		p.advance()
		stmt.ReturnType, err = p.parseDataType()
		if err != nil {
			return nil, err
		}
		if stmt.ReturnType == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}

	// EXTERNAL NAME assembly_qualified_name
	if p.cur.Type == kwEXTERNAL {
		p.advance()
		if p.cur.Type == kwNAME {
			p.advance()
		}
		// Read dotted name: assembly.class[.method]
		var parts []string
		if p.isIdentLike() {
			parts = append(parts, p.cur.Str)
			p.advance()
		}
		for {
			if _, ok := p.match('.'); ok {
				if p.isIdentLike() {
					parts = append(parts, p.cur.Str)
					p.advance()
				}
			} else {
				break
			}
		}
		stmt.ExternalName = joinDots(parts)
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropAggregateStmt parses DROP AGGREGATE [IF EXISTS] name.
// Caller has consumed DROP AGGREGATE.
//
// BNF: mssql/parser/bnf/drop-aggregate-transact-sql.bnf
//
//	DROP AGGREGATE [ IF EXISTS ] [ schema_name . ] aggregate_name
func (p *Parser) parseDropAggregateStmt() (*nodes.DropAggregateStmt, error) {
	loc := p.pos()
	stmt := &nodes.DropAggregateStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// IF EXISTS
	if _, ok := p.match(kwIF); ok {
		p.match(kwEXISTS)
		stmt.IfExists = true
	}

	var err error
	stmt.Name, err = p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// joinDots joins string parts with dots.
func joinDots(parts []string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += "."
		}
		result += p
	}
	return result
}

// parseCreateJsonIndexStmt parses a CREATE JSON INDEX statement.
// Caller has consumed CREATE JSON INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-json-index-transact-sql
//
//	CREATE JSON INDEX name ON table_name (json_column_name)
//	  [ FOR ( sql_json_path [ , ...n ] ) ]
//	  [ WITH ( <json_index_option> [ , ...n ] ) ]
//	  [ ON { filegroup_name | "default" } ]
//
//	<json_index_option> ::=
//	{
//	    OPTIMIZE_FOR_ARRAY_SEARCH = { ON | OFF }
//	  | FILLFACTOR = fillfactor
//	  | DROP_EXISTING = { ON | OFF }
//	  | ONLINE = OFF
//	  | ALLOW_ROW_LOCKS = { ON | OFF }
//	  | ALLOW_PAGE_LOCKS = { ON | OFF }
//	  | MAXDOP = max_degree_of_parallelism
//	  | DATA_COMPRESSION = { NONE | ROW | PAGE }
//	}
func (p *Parser) parseCreateJsonIndexStmt() (*nodes.CreateJsonIndexStmt, error) {
	loc := p.pos()
	stmt := &nodes.CreateJsonIndexStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Index name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	// ON table
	if _, ok := p.match(kwON); ok {
		var err error
		stmt.Table, err = p.parseTableRef()
		if err != nil {
			return nil, err
		}
	}

	// (json_column_name)
	if p.cur.Type == '(' {
		p.advance()
		col, _ := p.parseIdentifier()
		stmt.JsonColumn = col
		p.match(')')
	}

	// FOR (sql_json_path, ...)
	if _, ok := p.match(kwFOR); ok {
		if p.cur.Type == '(' {
			p.advance()
			var paths []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				// JSON paths are string literals like '$.name'
				expr, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if expr != nil {
					paths = append(paths, expr)
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.match(')')
			if len(paths) > 0 {
				stmt.ForPaths = &nodes.List{Items: paths}
			}
		}
	}

	// WITH (options)
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			var err error
			stmt.Options, err = p.parseOptionList()
			if err != nil {
				return nil, err
			}
		}
	}

	// ON filegroup
	if _, ok := p.match(kwON); ok {
		if p.isIdentLike() || p.cur.Type == kwPRIMARY {
			stmt.OnFileGroup = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateVectorIndexStmt parses a CREATE VECTOR INDEX statement.
// Caller has consumed CREATE VECTOR INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-vector-index-transact-sql
//
//	CREATE VECTOR INDEX index_name
//	ON object ( vector_column )
//	[ WITH (
//	    [ , ] METRIC = { 'cosine' | 'dot' | 'euclidean' }
//	    [ [ , ] TYPE = 'DiskANN' ]
//	    [ [ , ] MAXDOP = max_degree_of_parallelism ]
//	) ]
//	[ ON { filegroup_name | "default" } ]
func (p *Parser) parseCreateVectorIndexStmt() (*nodes.CreateVectorIndexStmt, error) {
	loc := p.pos()
	stmt := &nodes.CreateVectorIndexStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Index name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	// ON table
	if _, ok := p.match(kwON); ok {
		var err error
		stmt.Table, err = p.parseTableRef()
		if err != nil {
			return nil, err
		}
	}

	// (vector_column)
	if p.cur.Type == '(' {
		p.advance()
		col, _ := p.parseIdentifier()
		stmt.VectorCol = col
		p.match(')')
	}

	// WITH (options)
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			var err error
			stmt.Options, err = p.parseOptionList()
			if err != nil {
				return nil, err
			}
		}
	}

	// ON filegroup
	if _, ok := p.match(kwON); ok {
		if p.isIdentLike() || p.cur.Type == kwPRIMARY {
			stmt.OnFileGroup = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseOptionList parses (option = value, ...) used in WITH clauses.
//
// Option NAMES in T-SQL WITH clauses are frequently reserved keywords
// (FILLFACTOR, MAXDOP, ONLINE, DATA_COMPRESSION, ...) that are not
// permitted as unquoted identifiers in expression position. We therefore
// do not route the name through parseExpr; instead we consume any
// keyword-or-ident token as the name and parse the right-hand side as a
// normal expression.
func (p *Parser) parseOptionList() (*nodes.List, error) {
	p.advance() // consume (
	items, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		return p.parseNameValueOption()
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// parseNameValueOption consumes a single `name [= value]` option entry, where
// the name may be any keyword (Core or Context) or plain identifier.
func (p *Parser) parseNameValueOption() (nodes.Node, error) {
	loc := p.pos()
	if !p.isIdentLike() && !(p.cur.Type >= kwABSENT && p.cur.Str != "") {
		return nil, p.unexpectedToken()
	}
	name := p.cur.Str
	p.advance()
	nameRef := &nodes.ColumnRef{Column: name, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}
	if p.cur.Type != '=' {
		return nameRef, nil
	}
	p.advance() // consume =
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if rhs == nil {
		return nil, p.unexpectedToken()
	}
	return &nodes.BinaryExpr{
		Op:    nodes.BinOpEq,
		Left:  nameRef,
		Right: rhs,
		Loc:   nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}
