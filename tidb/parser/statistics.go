package parser

import (
	"strconv"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB's standalone statistics statements: ANALYZE TABLE, DROP STATS, LOAD
// STATS, CREATE STATISTICS and DROP STATISTICS. Only ANALYZE TABLE exists in
// MySQL, and there it is a much smaller grammar -- a table list and the
// UPDATE/DROP HISTOGRAM arms -- so this file holds TiDB's own shape of it
// rather than mysql/parser/utility.go's.
//
// TiDB spells most of the words below as keywords of its own (STATS,
// STATISTICS, INCREMENTAL, PREDICATE, TOPN, CMSKETCH, WIDTH, DEPTH, SAMPLES,
// SAMPLERATE, CARDINALITY, DEPENDENCY, CORRELATION). This parser registers
// none of them: registering a keyword changes every site that reads a name,
// and all of these are ordinary names everywhere else -- `CREATE TABLE stats
// (samples int, width int)` and `SELECT cardinality FROM t topn` both run on
// v8.5.5. Each is read as a bare word instead, through isBareWord, which
// keeps the distinction TiDB's lexer makes: a backquoted name is an
// identifier and never one of these words, so v8.5.5 answers 1064 for
// `DROP `STATS` t`, `ANALYZE TABLE t WITH 4 `buckets``, and `CREATE
// STATISTICS s (`cardinality`) ON t(a)`.
//
// Every verdict this file encodes was read from pingcap/tidb:v8.5.5 through
// PREPARE on 2026-11-04 and is held to that server by
// TestStatisticsFormsMatchTiDB.
//
// Ref: AnalyzeTableStmt, DropStatsStmt, LoadStatsStmt, CreateStatisticsStmt
// and DropStatisticsStmt in pingcap/tidb pkg/parser/parser.y.
// https://docs.pingcap.com/tidb/stable/sql-statement-analyze-table/

// parseAnalyzeTableStmt parses an ANALYZE TABLE statement. ANALYZE is the
// current token.
//
//	ANALYZE [NO_WRITE_TO_BINLOG | LOCAL] TABLE tbl_name [, tbl_name] ...
//	    [ALL COLUMNS | PREDICATE COLUMNS] [WITH analyze_option, ...]
//	ANALYZE [NO_WRITE_TO_BINLOG | LOCAL] TABLE tbl_name
//	    [PARTITION p_name [, p_name] ...]
//	    [ALL COLUMNS | PREDICATE COLUMNS | COLUMNS col, ... | INDEX [idx, ...]]
//	    [WITH analyze_option, ...]
//	ANALYZE [NO_WRITE_TO_BINLOG | LOCAL] INCREMENTAL TABLE tbl_name
//	    [PARTITION p_name, ...] INDEX [idx, ...] [WITH analyze_option, ...]
//	ANALYZE [NO_WRITE_TO_BINLOG | LOCAL] TABLE tbl_name
//	    UPDATE HISTOGRAM ON col_name [, col_name] ... [WITH analyze_option, ...]
//	ANALYZE [NO_WRITE_TO_BINLOG | LOCAL] TABLE tbl_name
//	    DROP HISTOGRAM ON col_name [, col_name] ...
//
// Which arm may name more than one table is the server's answer, not a
// guess: v8.5.5 takes `ANALYZE TABLE t1, t2 ALL COLUMNS` but answers 1064 for
// `t1, t2` with PARTITION, INDEX, a COLUMNS list, or either HISTOGRAM arm.
// The arms are mutually exclusive too -- `INDEX a ALL COLUMNS`, `ALL COLUMNS
// INDEX a` and `COLUMNS c1 INDEX a` are each 1064 -- and INCREMENTAL reaches
// only the INDEX arm, so `ANALYZE INCREMENTAL TABLE t` on its own is 1064.
//
// MySQL's `USING DATA 'json'` tail is not here: TiDB has no such form and
// v8.5.5 answers 1064 for it, where omni used to take it.
func (p *Parser) parseAnalyzeTableStmt() (*nodes.AnalyzeTableStmt, error) {
	start := p.pos()
	p.advance() // consume ANALYZE

	stmt := &nodes.AnalyzeTableStmt{Loc: nodes.Loc{Start: start}}

	// NO_WRITE_TO_BINLOG and LOCAL are the same modifier spelled two ways,
	// and both come before INCREMENTAL: v8.5.5 answers 1064 for `ANALYZE
	// INCREMENTAL LOCAL TABLE t INDEX a` and takes `ANALYZE LOCAL
	// INCREMENTAL TABLE t INDEX a`.
	if p.cur.Type == kwLOCAL || p.cur.Type == kwNO_WRITE_TO_BINLOG {
		stmt.NoWriteToBinlog = true
		p.advance()
	}
	if p.isBareWord("INCREMENTAL") {
		stmt.Incremental = true
		p.advance()
	}

	// TABLE is required. It used to be matched leniently, which left
	// `ANALYZE t` -- a 1064 on every one of these servers -- parsing as an
	// ANALYZE of t.
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	// Completion: after ANALYZE [LOCAL] TABLE, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	// Table list
	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Tables = append(stmt.Tables, ref)
		if _, ok := p.match(','); !ok {
			break
		}
	}

	if err := p.parseAnalyzeTail(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAnalyzeTail reads everything after ANALYZE's table list: the partition
// list, the column or index choice, the histogram arms, and the WITH options.
func (p *Parser) parseAnalyzeTail(stmt *nodes.AnalyzeTableStmt) error {
	// PARTITION p, ... -- one table only, and the names are bare, with no
	// parentheses: v8.5.5 answers 1064 for `PARTITION (p0)` here, where LOCK
	// STATS takes them.
	if p.cur.Type == kwPARTITION {
		if len(stmt.Tables) != 1 {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume PARTITION
		for {
			name, _, err := p.parseIdent()
			if err != nil {
				return err
			}
			stmt.Partitions = append(stmt.Partitions, name)
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	switch {
	case p.cur.Type == kwINDEX:
		if len(stmt.Tables) != 1 {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume INDEX
		stmt.HasIndex = true
		names, err := p.parseAnalyzeIndexNames()
		if err != nil {
			return err
		}
		stmt.Indexes = names

	case p.cur.Type == kwALL && p.peekNext().Type == kwCOLUMNS:
		p.advance() // consume ALL
		p.advance() // consume COLUMNS
		stmt.ColumnChoice = nodes.AnalyzeAllColumns

	case p.isBareWord("PREDICATE") && p.peekNext().Type == kwCOLUMNS:
		p.advance() // consume PREDICATE
		p.advance() // consume COLUMNS
		stmt.ColumnChoice = nodes.AnalyzePredicateColumns

	case p.cur.Type == kwCOLUMNS:
		if len(stmt.Tables) != 1 {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume COLUMNS
		stmt.ColumnChoice = nodes.AnalyzeColumnList
		names, err := p.parseAnalyzeColumnNames()
		if err != nil {
			return err
		}
		stmt.Columns = names

	case p.cur.Type == kwUPDATE || p.cur.Type == kwDROP:
		// The histogram arms take one table and reach neither INCREMENTAL
		// nor a partition list: v8.5.5 answers 1064 for `ANALYZE TABLE t1,
		// t2 UPDATE HISTOGRAM ON b`, for `ANALYZE INCREMENTAL TABLE t UPDATE
		// HISTOGRAM ON b`, and for `ANALYZE TABLE t PARTITION p0 UPDATE
		// HISTOGRAM ON b`.
		if len(stmt.Tables) != 1 || stmt.Incremental || len(stmt.Partitions) > 0 {
			return p.syntaxErrorAtCur()
		}
		if err := p.parseAnalyzeHistogram(stmt); err != nil {
			return err
		}
		// DROP HISTOGRAM takes no WITH list -- `DROP HISTOGRAM ON b WITH 4
		// BUCKETS` is 1064 -- while UPDATE HISTOGRAM takes the same list
		// every other arm does.
		if stmt.HistogramOp == "DROP" {
			if p.cur.Type == kwWITH {
				return p.syntaxErrorAtCur()
			}
			return nil
		}
	}

	// INCREMENTAL exists only on the INDEX arms: `ANALYZE INCREMENTAL TABLE
	// t`, `... ALL COLUMNS`, `... COLUMNS c1`, `... PARTITION p0` and
	// `... WITH 4 BUCKETS` are each 1064.
	if stmt.Incremental && !stmt.HasIndex {
		return p.syntaxErrorAtCur()
	}

	if p.cur.Type == kwWITH {
		opts, err := p.parseAnalyzeOptions()
		if err != nil {
			return err
		}
		stmt.Options = opts
	}

	// The arms above take one choice each, so a second choice word here is
	// a syntax error of ANALYZE's own: v8.5.5 answers 1064 for `INDEX a ALL
	// COLUMNS`, `ALL COLUMNS INDEX a`, `COLUMNS c1 INDEX a` and `PARTITION
	// p0 ALL COLUMNS INDEX a`. It is reported here rather than left to the
	// statement list, where a tail that happens to start a statement of its
	// own -- UPDATE and DROP both do -- would be read as one instead of as
	// the rejection it is.
	switch {
	case p.cur.Type == kwINDEX, p.cur.Type == kwCOLUMNS, p.cur.Type == kwALL,
		p.cur.Type == kwPARTITION, p.cur.Type == kwUPDATE, p.cur.Type == kwDROP,
		p.isBareWord("PREDICATE"):
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseAnalyzeHistogram reads the UPDATE HISTOGRAM or DROP HISTOGRAM arm.
// UPDATE or DROP is the current token.
//
// The column names are bare: TiDB's own rule takes a qualified name, but
// v8.5.5 answers 1064 for `UPDATE HISTOGRAM ON t.c1` and for `DROP HISTOGRAM
// ON t.c1, t.c2`, so a qualifier is refused here too.
func (p *Parser) parseAnalyzeHistogram(stmt *nodes.AnalyzeTableStmt) error {
	if p.cur.Type == kwUPDATE {
		stmt.HistogramOp = "UPDATE"
	} else {
		stmt.HistogramOp = "DROP"
	}
	p.advance() // consume UPDATE or DROP
	if _, err := p.expect(kwHISTOGRAM); err != nil {
		return err
	}
	if _, err := p.expect(kwON); err != nil {
		return err
	}
	for {
		col, _, err := p.parseIdent()
		if err != nil {
			return err
		}
		stmt.HistogramColumns = append(stmt.HistogramColumns, col)
		if _, ok := p.match(','); !ok {
			return nil
		}
	}
}

// parseAnalyzeIndexNames reads the index list of ANALYZE's INDEX arm, which
// may be empty.
//
// TiDB writes it left-recursive with an empty base case:
//
//	IndexNameList: /* empty */ | Identifier | IndexNameList ',' Identifier
//
// so the empty list may be followed by a comma and v8.5.5 takes `INDEX ,a`
// and `INDEX ,a,b` while `INDEX a,`, `INDEX ,,a` and `INDEX a,,b` are each
// 1064. That is reproduced rather than tidied: the corpus holds the bare
// `INDEX` form, and a reader that tidied the leading comma away would differ
// from the server on text the server takes.
//
// IndexNameList also takes PRIMARY as a name, in and beyond the first slot, so
// v8.5.5 runs `INDEX PRIMARY` and `INDEX ,PRIMARY`; PRIMARY is a reserved word
// and no Identifier, so it is read on its own.
func (p *Parser) parseAnalyzeIndexNames() ([]string, error) {
	var names []string
	if p.cur.Type != ',' {
		// An empty list is what follows a bare INDEX, so anything that
		// cannot start a name ends the arm here and WITH or the end of the
		// statement reads it.
		if !p.isIdentToken() && p.cur.Type != kwPRIMARY {
			return nil, nil
		}
		name, err := p.parseIndexNameOrPrimary()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	for {
		if _, ok := p.match(','); !ok {
			return names, nil
		}
		name, err := p.parseIndexNameOrPrimary()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
}

// parseIndexNameOrPrimary reads one entry of an IndexNameList: an identifier,
// or the PRIMARY keyword naming the primary-key index. Its string is what the
// token spells.
func (p *Parser) parseIndexNameOrPrimary() (string, error) {
	if p.cur.Type == kwPRIMARY {
		return p.advance().Str, nil
	}
	name, _, err := p.parseIdent()
	return name, err
}

// parseStatsType reads a StatsType word: CARDINALITY, DEPENDENCY or
// CORRELATION, each read bare. These are the three extended-statistics kinds
// CREATE STATISTICS and ADD STATS_EXTENDED share.
//
// Ref: StatsType in TiDB v8.5.5 pkg/parser/parser.y.
func (p *Parser) parseStatsType() (nodes.StatisticsType, error) {
	switch {
	case p.isBareWord("CARDINALITY"):
		p.advance()
		return nodes.StatisticsCardinality, nil
	case p.isBareWord("DEPENDENCY"):
		p.advance()
		return nodes.StatisticsDependency, nil
	case p.isBareWord("CORRELATION"):
		p.advance()
		return nodes.StatisticsCorrelation, nil
	}
	return 0, p.syntaxErrorAtCur()
}

// parseStatsColumnList reads the parenthesized column list of an extended
// statistic: `(col [, col] ...)`. The names may be qualified, as in CREATE
// STATISTICS, but `t.*` is not a name and the list is never empty.
func (p *Parser) parseStatsColumnList() ([]*nodes.ColumnRef, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var cols []*nodes.ColumnRef
	for {
		col, err := p.parseColumnRef()
		if err != nil {
			return nil, err
		}
		if col.Star {
			return nil, &ParseError{Message: "expected a column name", Position: col.Loc.Start}
		}
		cols = append(cols, col)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseAnalyzeColumnNames reads the column list of ANALYZE's COLUMNS arm.
// Unlike the index list it is never empty and its names carry no qualifier:
// v8.5.5 answers 1064 for a bare `COLUMNS`, for `COLUMNS c1,` and for
// `COLUMNS t.c1, t.c2`.
func (p *Parser) parseAnalyzeColumnNames() ([]string, error) {
	var names []string
	for {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if _, ok := p.match(','); !ok {
			return names, nil
		}
	}
}

// parseAnalyzeOptions reads ANALYZE's WITH list. WITH is the current token.
//
//	WITH analyze_option [, analyze_option] ...
//
// One WITH carries the whole list: v8.5.5 takes `WITH 4 BUCKETS, 4 TOPN` and
// answers 1064 for `WITH 4 BUCKETS WITH 4 TOPN` and for `WITH 4 BUCKETS 4
// TOPN`. The list is not a set -- `WITH 4 BUCKETS, 4 BUCKETS` parses -- so a
// repeat is kept rather than refused here.
func (p *Parser) parseAnalyzeOptions() ([]*nodes.AnalyzeOption, error) {
	p.advance() // consume WITH
	var opts []*nodes.AnalyzeOption
	for {
		opt, err := p.parseAnalyzeOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
		if _, ok := p.match(','); !ok {
			return opts, nil
		}
	}
}

// parseAnalyzeOption reads one `n UNIT` entry of the WITH list.
//
//	n BUCKETS | n TOPN | n CMSKETCH WIDTH | n CMSKETCH DEPTH | n SAMPLES
//	n SAMPLERATE
//
// The unit list is closed, which is the point of reading it here rather than
// taking any word: v8.5.5 answers 1064 for `WITH 4 PENGUINS`, for a bare
// `WITH 4`, for `WITH BUCKETS`, and for `WITH 4 CMSKETCH` with neither WIDTH
// nor DEPTH after it.
//
// Five of the six units take an unsigned decimal integer and SAMPLERATE takes
// any unsigned number, so `WITH 1.5 BUCKETS`, `WITH 0.5 TOPN`, `WITH -1
// BUCKETS` and `WITH '4' BUCKETS` are 1064 while `WITH .5 SAMPLERATE`, `WITH
// 1e3 SAMPLERATE` and `WITH 4 SAMPLERATE` all parse. A hex or bit literal is
// neither -- `WITH 0x10 BUCKETS` and `WITH 0x10 SAMPLERATE` are both 1064 --
// and each lexes as its own token here, so neither reaches a unit. An integer
// past uint64 is a decimal to TiDB's lexer and not an integer literal, so
// `WITH 18446744073709551616 BUCKETS` is 1064 where
// `WITH 18446744073709551615 BUCKETS` parses; the check is on the written
// text because Token.Ival saturates.
func (p *Parser) parseAnalyzeOption() (*nodes.AnalyzeOption, error) {
	start := p.pos()
	num := p.cur
	switch num.Type {
	case tokICONST, tokFCONST:
		p.advance()
	default:
		return nil, p.syntaxErrorAtCur()
	}

	opt := &nodes.AnalyzeOption{Loc: nodes.Loc{Start: start}, Value: num.Str}

	switch {
	case p.isBareWord("BUCKETS"):
		opt.Type = nodes.AnalyzeOptBuckets
	case p.isBareWord("TOPN"):
		opt.Type = nodes.AnalyzeOptTopN
	case p.isBareWord("SAMPLES"):
		opt.Type = nodes.AnalyzeOptSamples
	case p.isBareWord("SAMPLERATE"):
		opt.Type = nodes.AnalyzeOptSampleRate
	case p.isBareWord("CMSKETCH"):
		p.advance() // consume CMSKETCH
		switch {
		case p.isBareWord("WIDTH"):
			opt.Type = nodes.AnalyzeOptCMSketchWidth
		case p.isBareWord("DEPTH"):
			opt.Type = nodes.AnalyzeOptCMSketchDepth
		default:
			return nil, p.syntaxErrorAtCur()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume the unit word

	if opt.Type != nodes.AnalyzeOptSampleRate {
		if num.Type != tokICONST {
			return nil, &ParseError{Message: "expected an integer", Position: num.Loc}
		}
		if _, err := strconv.ParseUint(num.Str, 10, 64); err != nil {
			return nil, &ParseError{Message: "integer out of range", Position: num.Loc}
		}
	}

	opt.Loc.End = p.prev.End
	return opt, nil
}

// parseDropStatsStmt parses TiDB's DROP STATS. STATS is the current token and
// start is where DROP began.
//
//	DROP STATS tbl_name [, tbl_name] ...
//	DROP STATS tbl_name PARTITION p_name [, p_name] ...
//	DROP STATS tbl_name GLOBAL
//
// GLOBAL and the partition list each take one table -- `DROP STATS t1, t2
// GLOBAL` and `DROP STATS t1, t2 PARTITION p0` are both 1064 -- and they do
// not combine: `DROP STATS t GLOBAL PARTITION p0` is 1064 as well. The
// partition names take no parentheses here either, and a trailing comma is
// 1064.
func (p *Parser) parseDropStatsStmt(start int) (*nodes.DropStatsStmt, error) {
	p.advance() // consume STATS

	stmt := &nodes.DropStatsStmt{Loc: nodes.Loc{Start: start}}
	for {
		// Completion: each position in the list offers a table.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
		}
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Tables = append(stmt.Tables, ref)
		if _, ok := p.match(','); !ok {
			break
		}
	}

	switch {
	case p.cur.Type == kwPARTITION:
		if len(stmt.Tables) != 1 {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume PARTITION
		for {
			name, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			stmt.Partitions = append(stmt.Partitions, name)
			if _, ok := p.match(','); !ok {
				break
			}
		}
	case p.cur.Type == kwGLOBAL:
		if len(stmt.Tables) != 1 {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume GLOBAL
		stmt.Global = true
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseLoadStatsStmt parses TiDB's LOAD STATS. STATS is the current token and
// start is where LOAD began.
//
//	LOAD STATS 'file_path'
//
// The path is a string literal and nothing else: `LOAD STATS`, `LOAD STATS 4`
// and `LOAD STATS 'a' 'b'` are each 1064.
func (p *Parser) parseLoadStatsStmt(start int) (*nodes.LoadStatsStmt, error) {
	p.advance() // consume STATS
	tok, err := p.expect(tokSCONST)
	if err != nil {
		return nil, err
	}
	return &nodes.LoadStatsStmt{
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
		Path: tok.Str,
	}, nil
}

// parseCreateStatisticsStmt parses TiDB's CREATE STATISTICS. STATISTICS is
// the current token and start is where CREATE began.
//
//	CREATE STATISTICS [IF NOT EXISTS] stats_name (stats_type)
//	    ON tbl_name (col_name [, col_name] ...)
//
// Exactly one type in the parentheses, from a closed list: `CREATE STATISTICS
// s ON t(a)` with no type, `(frequency)` and `(cardinality, dependency)` are
// each 1064. The column list is required and not empty, and `(a,)` is 1064.
// Unlike ANALYZE's COLUMNS list these names may be qualified -- v8.5.5 takes
// `ON t(t.a)` and `ON t(d.t.a)` -- but `t.*` is not a name: `ON t(a.*)` is
// 1064.
func (p *Parser) parseCreateStatisticsStmt(start int) (*nodes.CreateStatisticsStmt, error) {
	p.advance() // consume STATISTICS

	stmt := &nodes.CreateStatisticsStmt{Loc: nodes.Loc{Start: start}}
	if p.cur.Type == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	statsType, err := p.parseStatsType()
	if err != nil {
		return nil, err
	}
	stmt.Type = statsType
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	table, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	cols, err := p.parseStatsColumnList()
	if err != nil {
		return nil, err
	}
	stmt.Columns = cols

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropStatisticsStmt parses TiDB's DROP STATISTICS. STATISTICS is the
// current token and start is where DROP began.
//
//	DROP STATISTICS stats_name
//
// One name, and no IF EXISTS: `DROP STATISTICS`, `DROP STATISTICS s1, s2` and
// `DROP STATISTICS IF EXISTS s1` are each 1064.
func (p *Parser) parseDropStatisticsStmt(start int) (*nodes.DropStatisticsStmt, error) {
	p.advance() // consume STATISTICS
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	return &nodes.DropStatisticsStmt{
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
		Name: name,
	}, nil
}
