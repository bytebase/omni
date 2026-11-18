package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseCreateMaterializedStmt parses what follows CREATE [OR REPLACE] when
// MATERIALIZED stands there. p.cur is MATERIALIZED. MATERIALIZED VIEW LOG is
// the only object OR REPLACE takes; a materialized view itself has no OR
// REPLACE alternative in the grammar.
//
// Ref: create_mview_stmt and create_mlog_stmt in OceanBase's
// src/sql/parser/sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines 9766 and 10724).
func (p *Parser) parseCreateMaterializedStmt(start int, orReplace bool) (nodes.Node, error) {
	p.advance() // consume MATERIALIZED
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	if p.isBareWord("LOG") {
		return p.parseCreateMaterializedViewLogStmt(start, orReplace)
	}
	if orReplace {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseCreateMaterializedViewStmt(start)
}

// parseCreateMaterializedViewStmt parses CREATE MATERIALIZED VIEW. CREATE,
// MATERIALIZED and VIEW are consumed, so p.cur is the view's name.
//
//	CREATE MATERIALIZED VIEW view_name
//	    [ (column_name_list [, PRIMARY KEY [index_using_algorithm] (column_name_list)
//	       [index_using_algorithm] [COMMENT 'string']]) ]
//	    [table_option_list] [partition_option] [with_column_group]
//	    [create_mview_opts] AS view_select_stmt [opt_check_option]
func (p *Parser) parseCreateMaterializedViewStmt(start int) (*nodes.CreateMaterializedViewStmt, error) {
	stmt := &nodes.CreateMaterializedViewStmt{Loc: nodes.Loc{Start: start}}

	ref, err := p.parseRelationFactor()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref

	if p.cur.Type == '(' {
		cols, pk, err := p.parseMviewColumnList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
		stmt.PrimaryKey = pk
	}

	// A comma may stand between two table options and never ends the list; a
	// comma with no option after it is a syntax error. This is the loop
	// CREATE TABLE uses for the same rule.
	for {
		opt, ok, err := p.parseTableOption()
		if err != nil {
			return nil, err
		}
		if !ok {
			if len(stmt.Options) > 0 && p.prev.Type == ',' {
				return nil, p.syntaxErrorAtCur()
			}
			break
		}
		stmt.Options = append(stmt.Options, opt)
		p.match(',')
	}

	if p.cur.Type == kwPARTITION {
		part, err := p.parsePartitionClause(true)
		if err != nil {
			return nil, err
		}
		stmt.Partitions = part
	}

	if group, ok, err := p.parseWithColumnGroup(); err != nil {
		return nil, err
	} else if ok {
		stmt.ColumnGroup = group
	}

	opts, err := p.parseCreateMviewOptions()
	if err != nil {
		return nil, err
	}
	stmt.MviewOptions = opts

	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	selectStart := p.pos()
	sel, err := p.parseQueryStatement()
	if err != nil {
		return nil, err
	}
	stmt.Select = sel
	stmt.SelectText = strings.TrimSpace(p.inputText(selectStart, p.pos()))

	checkOption, err := p.parseViewCheckOption()
	if err != nil {
		return nil, err
	}
	stmt.CheckOption = checkOption

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseViewCheckOption parses opt_check_option: WITH [CASCADED | LOCAL] CHECK
// OPTION, and returns CASCADED, LOCAL, or "" when the clause is absent. CHECK
// and OPTION are both required once WITH stands there.
func (p *Parser) parseViewCheckOption() (string, error) {
	if p.cur.Type != kwWITH {
		return "", nil
	}
	p.advance()
	checkOption := "CASCADED"
	if p.cur.Type == kwCASCADED {
		p.advance()
	} else if p.cur.Type == kwLOCAL {
		checkOption = "LOCAL"
		p.advance()
	}
	if _, err := p.expect(kwCHECK); err != nil {
		return "", err
	}
	if _, err := p.expect(kwOPTION); err != nil {
		return "", err
	}
	return checkOption, nil
}

// parseMviewColumnList parses opt_mv_column_list's parentheses:
//
//	'(' mv_column_list ')'
//	mv_column_list:
//	    column_name_list
//	  | column_name_list ',' PRIMARY KEY opt_index_using_algorithm '(' column_name_list ')'
//	        opt_index_using_algorithm opt_comment
//	  | PRIMARY KEY opt_index_using_algorithm '(' column_name_list ')'
//	        opt_index_using_algorithm opt_comment
//
// It returns the plain columns, then the primary key clause if one is written.
func (p *Parser) parseMviewColumnList() ([]string, *nodes.MaterializedViewPrimaryKey, error) {
	if _, err := p.expect('('); err != nil {
		return nil, nil, err
	}

	if p.cur.Type == kwPRIMARY {
		pk, err := p.parseMviewPrimaryKey()
		if err != nil {
			return nil, nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, nil, err
		}
		return nil, pk, nil
	}

	var cols []string
	for {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, nil, err
		}
		cols = append(cols, name)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
		if p.cur.Type == kwPRIMARY {
			pk, err := p.parseMviewPrimaryKey()
			if err != nil {
				return nil, nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, nil, err
			}
			return cols, pk, nil
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, nil, err
	}
	return cols, nil, nil
}

// parseMviewPrimaryKey parses PRIMARY KEY in a materialized view's column
// list. p.cur is PRIMARY.
func (p *Parser) parseMviewPrimaryKey() (*nodes.MaterializedViewPrimaryKey, error) {
	pk := &nodes.MaterializedViewPrimaryKey{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // consume PRIMARY
	if _, err := p.expect(kwKEY); err != nil {
		return nil, err
	}

	if p.cur.Type == kwUSING {
		p.advance()
		indexType, err := p.parseIndexTypeName()
		if err != nil {
			return nil, err
		}
		pk.IndexType = indexType
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	for {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		pk.Columns = append(pk.Columns, name)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if p.cur.Type == kwUSING {
		p.advance()
		indexType, err := p.parseIndexTypeName()
		if err != nil {
			return nil, err
		}
		pk.IndexType = indexType
	}

	if p.cur.Type == kwCOMMENT {
		p.advance()
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		pk.Comment = p.cur.Str
		p.advance()
	}

	pk.Loc.End = p.prev.End
	return pk, nil
}

// parseCreateMviewOptions parses create_mview_opts:
//
//	mview_refresh_opt [mview_enable_disable ON QUERY COMPUTATION]
//	                   [mview_enable_disable QUERY REWRITE]
//
// The two ENABLE/DISABLE clauses may stand in either order. It returns nil when
// neither a refresh clause nor a switch is written.
func (p *Parser) parseCreateMviewOptions() (*nodes.MviewOptions, error) {
	start := p.pos()

	refresh, err := p.parseCreateMviewRefresh()
	if err != nil {
		return nil, err
	}

	opts := &nodes.MviewOptions{Loc: nodes.Loc{Start: start}, Refresh: refresh}
	for p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE {
		kind, enable, err := p.parseMviewEnableDisable()
		if err != nil {
			return nil, err
		}
		if kind == "ON QUERY COMPUTATION" {
			if opts.OnQueryComputation != "" {
				return nil, p.syntaxErrorAtCur()
			}
			opts.OnQueryComputation = mviewEnableWord(enable)
		} else {
			if opts.QueryRewrite != "" {
				return nil, p.syntaxErrorAtCur()
			}
			opts.QueryRewrite = mviewEnableWord(enable)
		}
	}

	if opts.Refresh == nil && opts.OnQueryComputation == "" && opts.QueryRewrite == "" {
		return nil, nil
	}
	opts.Loc.End = p.prev.End
	return opts, nil
}

func mviewEnableWord(enable bool) string {
	if enable {
		return "ENABLE"
	}
	return "DISABLE"
}

// parseMviewEnableDisable parses one mview_enable_disable clause: ENABLE or
// DISABLE, then ON QUERY COMPUTATION or QUERY REWRITE.
func (p *Parser) parseMviewEnableDisable() (kind string, enable bool, err error) {
	enable = p.cur.Type == kwENABLE
	p.advance() // consume ENABLE or DISABLE

	if p.cur.Type == kwON {
		p.advance()
		if _, err := p.expect(kwQUERY); err != nil {
			return "", false, err
		}
		if !p.isBareWord("COMPUTATION") {
			return "", false, p.syntaxErrorAtCur()
		}
		p.advance()
		return "ON QUERY COMPUTATION", enable, nil
	}
	if p.cur.Type == kwQUERY {
		p.advance()
		if !p.isBareWord("REWRITE") {
			return "", false, p.syntaxErrorAtCur()
		}
		p.advance()
		return "QUERY REWRITE", enable, nil
	}
	return "", false, p.syntaxErrorAtCur()
}

// parseCreateMviewRefresh parses create_mview_opts' mview_refresh_opt:
//
//	REFRESH mv_refresh_method mv_nested_refresh_opt mv_refresh_on_clause mv_refresh_interval
//	| NEVER REFRESH
//	| /* empty */
//
// Unlike ALTER MATERIALIZED VIEW's REFRESH action, REFRESH here requires a
// method.
func (p *Parser) parseCreateMviewRefresh() (*nodes.MviewRefreshInfo, error) {
	start := p.pos()

	if p.isBareWord("NEVER") {
		p.advance()
		if !p.isBareWord("REFRESH") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		return &nodes.MviewRefreshInfo{Loc: nodes.Loc{Start: start, End: p.prev.End}, Never: true}, nil
	}
	if !p.isBareWord("REFRESH") {
		return nil, nil
	}
	p.advance() // consume REFRESH

	method, ok, err := p.parseMvRefreshMethod()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	info := &nodes.MviewRefreshInfo{Loc: nodes.Loc{Start: start}, Method: method}

	nested, ok, err := p.parseMvNestedRefreshMode()
	if err != nil {
		return nil, err
	}
	if ok {
		info.NestedRefresh = nested
	}

	if p.cur.Type == kwON {
		p.advance()
		mode, ok, err := p.parseMvRefreshMode()
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		info.On = mode
	}

	if err := p.parseMvRefreshInterval(info); err != nil {
		return nil, err
	}

	info.Loc.End = p.prev.End
	return info, nil
}

// parseMvRefreshMethod reads FAST, COMPLETE or FORCE.
func (p *Parser) parseMvRefreshMethod() (string, bool, error) {
	for _, word := range []string{"FAST", "COMPLETE", "FORCE"} {
		if p.isBareWord(word) {
			p.advance()
			return word, true, nil
		}
	}
	return "", false, nil
}

// parseMvRefreshMode reads DEMAND, COMMIT or STATEMENT.
func (p *Parser) parseMvRefreshMode() (string, bool, error) {
	for _, word := range []string{"DEMAND", "COMMIT", "STATEMENT"} {
		if p.isBareWord(word) {
			p.advance()
			return word, true, nil
		}
	}
	return "", false, nil
}

// parseMvNestedRefreshMode reads INDIVIDUAL, INCONSISTENT or CONSISTENT.
func (p *Parser) parseMvNestedRefreshMode() (string, bool, error) {
	for _, word := range []string{"INDIVIDUAL", "INCONSISTENT", "CONSISTENT"} {
		if p.isBareWord(word) {
			p.advance()
			return word, true, nil
		}
	}
	return "", false, nil
}

// parseMvRefreshInterval reads mv_refresh_interval: an optional START WITH
// bit_expr followed by an optional NEXT bit_expr. Both are optional, so an
// empty interval is accepted.
func (p *Parser) parseMvRefreshInterval(info *nodes.MviewRefreshInfo) error {
	if p.cur.Type == kwSTART {
		p.advance()
		if !p.isBareWord("WITH") {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		expr, err := p.parseExprPrec(precBitOr)
		if err != nil {
			return err
		}
		info.StartWith = expr
	}
	if p.cur.Type == kwNEXT {
		p.advance()
		expr, err := p.parseExprPrec(precBitOr)
		if err != nil {
			return err
		}
		info.Next = expr
	}
	return nil
}

// parseCreateMaterializedViewLogStmt parses CREATE [OR REPLACE] MATERIALIZED
// VIEW LOG. p.cur is LOG; CREATE, MATERIALIZED and VIEW are consumed.
//
//	CREATE [OR REPLACE] MATERIALIZED VIEW LOG ON relation_factor
//	    opt_mlog_option_list opt_mlog_with opt_mlog_new_values opt_mlog_purge
func (p *Parser) parseCreateMaterializedViewLogStmt(start int, orReplace bool) (*nodes.CreateMaterializedViewLogStmt, error) {
	p.advance() // consume LOG
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}

	stmt := &nodes.CreateMaterializedViewLogStmt{
		Loc:       nodes.Loc{Start: start},
		OrReplace: orReplace,
	}

	ref, err := p.parseRelationFactor()
	if err != nil {
		return nil, err
	}
	stmt.Table = ref

	for {
		opt, ok, err := p.parseParallelOption()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		stmt.Options = append(stmt.Options, opt)
	}

	if p.cur.Type == kwWITH {
		withStart := p.pos()
		p.advance() // consume WITH
		with, err := p.parseMlogWith(withStart)
		if err != nil {
			return nil, err
		}
		stmt.With = with
	}

	if p.isBareWord("INCLUDING") || p.isBareWord("EXCLUDING") {
		stmt.NewValues = strings.ToUpper(p.advance().Str)
		if !p.isBareWord("NEW") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if !p.isBareWord("VALUES") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	if p.cur.Type == kwPURGE {
		purge, err := p.parseMlogPurge()
		if err != nil {
			return nil, err
		}
		stmt.Purge = purge
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseMlogWith parses what follows WITH in a materialized view log: any
// special columns, then an optional parenthesized reference column list. WITH
// alone is accepted, and so is WITH ().
func (p *Parser) parseMlogWith(start int) (*nodes.MlogWith, error) {
	with := &nodes.MlogWith{Loc: nodes.Loc{Start: start}}

	if p.isBareWord("PRIMARY") || p.isBareWord("ROWID") || p.isBareWord("SEQUENCE") {
		for {
			word, ok, err := p.parseMlogSpecialColumn()
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, p.syntaxErrorAtCur()
			}
			with.SpecialColumns = append(with.SpecialColumns, word)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
	}

	if p.cur.Type == '(' {
		p.advance() // consume '('
		with.HasReferenceList = true
		if p.cur.Type != ')' {
			for {
				name, _, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				with.ReferenceColumns = append(with.ReferenceColumns, name)
				if _, ok := p.match(','); !ok {
					break
				}
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	with.Loc.End = p.prev.End
	return with, nil
}

// parseMlogSpecialColumn reads ROWID, PRIMARY KEY or SEQUENCE.
func (p *Parser) parseMlogSpecialColumn() (string, bool, error) {
	if p.isBareWord("PRIMARY") {
		p.advance()
		if _, err := p.expect(kwKEY); err != nil {
			return "", false, err
		}
		return "PRIMARY KEY", true, nil
	}
	for _, word := range []string{"ROWID", "SEQUENCE"} {
		if p.isBareWord(word) {
			p.advance()
			return word, true, nil
		}
	}
	return "", false, nil
}

// parseMlogPurge parses opt_mlog_purge: PURGE IMMEDIATE [SYNCHRONOUS |
// ASYNCHRONOUS], or PURGE [START WITH bit_expr] [NEXT bit_expr]. p.cur is
// PURGE.
func (p *Parser) parseMlogPurge() (*nodes.MlogPurge, error) {
	purge := &nodes.MlogPurge{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // consume PURGE

	if p.isBareWord("IMMEDIATE") {
		p.advance()
		purge.Immediate = true
		if p.isBareWord("SYNCHRONOUS") || p.isBareWord("ASYNCHRONOUS") {
			purge.Sync = strings.ToUpper(p.advance().Str)
		}
	} else {
		info := &nodes.MviewRefreshInfo{}
		if err := p.parseMvRefreshInterval(info); err != nil {
			return nil, err
		}
		purge.StartWith = info.StartWith
		purge.Next = info.Next
	}

	purge.Loc.End = p.prev.End
	return purge, nil
}

// parseParallelOption reads PARALLEL [=] n or NOPARALLEL.
func (p *Parser) parseParallelOption() (*nodes.ParallelOption, bool, error) {
	start := p.pos()

	if p.cur.Type == kwPARALLEL {
		p.advance()
		p.match('=')
		if p.cur.Type != tokICONST || p.cur.Ival < 1 {
			return nil, false, p.syntaxErrorAtCur()
		}
		degree := p.cur.Ival
		p.advance()
		return &nodes.ParallelOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Degree: degree}, true, nil
	}
	if p.isBareWord("NOPARALLEL") {
		p.advance()
		return &nodes.ParallelOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, NoParallel: true}, true, nil
	}
	return nil, false, nil
}

// parseAlterMaterializedStmt parses ALTER MATERIALIZED VIEW and ALTER
// MATERIALIZED VIEW LOG. p.cur is MATERIALIZED.
func (p *Parser) parseAlterMaterializedStmt(start int) (nodes.Node, error) {
	p.advance() // consume MATERIALIZED
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	if p.isBareWord("LOG") {
		return p.parseAlterMaterializedViewLogStmt(start)
	}
	return p.parseAlterMaterializedViewStmt(start)
}

// parseAlterMaterializedViewStmt parses ALTER MATERIALIZED VIEW. ALTER,
// MATERIALIZED and VIEW are consumed, so p.cur is the view's name.
func (p *Parser) parseAlterMaterializedViewStmt(start int) (*nodes.AlterMaterializedViewStmt, error) {
	ref, err := p.parseRelationFactor()
	if err != nil {
		return nil, err
	}
	stmt := &nodes.AlterMaterializedViewStmt{Loc: nodes.Loc{Start: start}, Name: ref}

	for {
		action, err := p.parseAlterMviewAction()
		if err != nil {
			return nil, err
		}
		stmt.Actions = append(stmt.Actions, action)
		if _, ok := p.match(','); !ok {
			break
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterMviewAction parses one alter_mview_action.
func (p *Parser) parseAlterMviewAction() (*nodes.AlterMviewAction, error) {
	start := p.pos()

	if p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE {
		kind, enable, err := p.parseMviewEnableDisable()
		if err != nil {
			return nil, err
		}
		action := &nodes.AlterMviewAction{Loc: nodes.Loc{Start: start}, Enable: enable}
		if kind == "ON QUERY COMPUTATION" {
			action.Type = nodes.AMViewEnableOnQueryComputation
		} else {
			action.Type = nodes.AMViewEnableQueryRewrite
		}
		action.Loc.End = p.prev.End
		return action, nil
	}

	if p.cur.Type == kwSET {
		p.advance()
		par, ok, err := p.parseParallelOption()
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.AlterMviewAction{Loc: nodes.Loc{Start: start, End: p.prev.End}, Type: nodes.AMViewParallel, Parallel: par}, nil
	}
	if par, ok, err := p.parseParallelOption(); err != nil {
		return nil, err
	} else if ok {
		return &nodes.AlterMviewAction{Loc: nodes.Loc{Start: start, End: p.prev.End}, Type: nodes.AMViewParallel, Parallel: par}, nil
	}

	if p.isBareWord("NEVER") {
		p.advance()
		if !p.isBareWord("REFRESH") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		info := &nodes.MviewRefreshInfo{Loc: nodes.Loc{Start: start, End: p.prev.End}, Never: true}
		return &nodes.AlterMviewAction{Loc: nodes.Loc{Start: start, End: p.prev.End}, Type: nodes.AMViewNeverRefresh, Refresh: info}, nil
	}

	if p.isBareWord("REFRESH") {
		p.advance() // consume REFRESH

		if method, ok, err := p.parseMvRefreshMethod(); err != nil {
			return nil, err
		} else if ok {
			info := &nodes.MviewRefreshInfo{Loc: nodes.Loc{Start: start, End: p.prev.End}, Method: method}
			return &nodes.AlterMviewAction{Loc: nodes.Loc{Start: start, End: p.prev.End}, Type: nodes.AMViewRefreshMethod, Refresh: info}, nil
		}

		if nested, ok, err := p.parseMvNestedRefreshMode(); err != nil {
			return nil, err
		} else if ok {
			info := &nodes.MviewRefreshInfo{Loc: nodes.Loc{Start: start, End: p.prev.End}, NestedRefresh: nested}
			return &nodes.AlterMviewAction{Loc: nodes.Loc{Start: start, End: p.prev.End}, Type: nodes.AMViewNestedRefresh, Refresh: info}, nil
		}

		info := &nodes.MviewRefreshInfo{Loc: nodes.Loc{Start: start}}
		if err := p.parseMvRefreshInterval(info); err != nil {
			return nil, err
		}
		info.Loc.End = p.prev.End
		return &nodes.AlterMviewAction{Loc: nodes.Loc{Start: start, End: p.prev.End}, Type: nodes.AMViewRefreshInterval, Refresh: info}, nil
	}

	return nil, p.syntaxErrorAtCur()
}

// parseAlterMaterializedViewLogStmt parses ALTER MATERIALIZED VIEW LOG ON
// <table>. p.cur is LOG.
func (p *Parser) parseAlterMaterializedViewLogStmt(start int) (*nodes.AlterMaterializedViewLogStmt, error) {
	p.advance() // consume LOG
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	ref, err := p.parseRelationFactor()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.AlterMaterializedViewLogStmt{Loc: nodes.Loc{Start: start}, Table: ref}
	for {
		action, err := p.parseAlterMlogAction()
		if err != nil {
			return nil, err
		}
		stmt.Actions = append(stmt.Actions, action)
		if _, ok := p.match(','); !ok {
			break
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterMlogAction parses one alter_mlog_action.
func (p *Parser) parseAlterMlogAction() (*nodes.AlterMlogAction, error) {
	start := p.pos()

	if p.cur.Type == kwPURGE {
		p.advance()
		info := &nodes.MviewRefreshInfo{}
		if err := p.parseMvRefreshInterval(info); err != nil {
			return nil, err
		}
		return &nodes.AlterMlogAction{
			Loc:       nodes.Loc{Start: start, End: p.prev.End},
			Type:      nodes.AMLogPurge,
			StartWith: info.StartWith,
			Next:      info.Next,
		}, nil
	}

	if p.isBareWord("LOB_INROW_THRESHOLD") {
		p.advance()
		p.match('=')
		if p.cur.Type != tokICONST {
			return nil, p.syntaxErrorAtCur()
		}
		threshold := p.cur.Ival
		p.advance()
		return &nodes.AlterMlogAction{
			Loc:       nodes.Loc{Start: start, End: p.prev.End},
			Type:      nodes.AMLogLobInRowThreshold,
			Threshold: threshold,
		}, nil
	}

	if par, ok, err := p.parseParallelOption(); err != nil {
		return nil, err
	} else if ok {
		return &nodes.AlterMlogAction{
			Loc:      nodes.Loc{Start: start, End: p.prev.End},
			Type:     nodes.AMLogParallel,
			Parallel: par,
		}, nil
	}

	return nil, p.syntaxErrorAtCur()
}

// parseDropMaterializedStmt parses DROP MATERIALIZED VIEW and DROP
// MATERIALIZED VIEW LOG. p.cur is MATERIALIZED.
func (p *Parser) parseDropMaterializedStmt(start int) (nodes.Node, error) {
	p.advance() // consume MATERIALIZED
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	if p.isBareWord("LOG") {
		p.advance() // consume LOG
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		ref, err := p.parseRelationFactor()
		if err != nil {
			return nil, err
		}
		return &nodes.DropMaterializedViewLogStmt{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Table: ref,
		}, nil
	}

	// ALTER, MATERIALIZED and VIEW are consumed, and parseDropViewBody reads the
	// rest. The materialized flag keeps a materialized view apart from a plain
	// view.
	stmt, err := p.parseDropViewBody(start)
	if err != nil {
		return nil, err
	}
	stmt.Materialized = true
	return stmt, nil
}
