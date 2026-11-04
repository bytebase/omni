package parser

import (
	"math"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// ADMIN is TiDB's own family of maintenance statements: no MySQL, MariaDB, or
// OceanBase server has one, and all three answer 1064 for every form below.
//
// The family is broad rather than deep -- some thirty subcommands, each with
// its own tail -- and TiDB spells every word of it as a keyword of its own
// (DDL, JOBS, QUERIES, NEXT_ROW_ID, BDR, ...). omni registers none of them:
// registering a keyword changes every site that reads a name, and these words
// are ordinary names everywhere outside ADMIN (`CREATE TABLE jobs (queries
// int)`, `DELETE FROM t ddl WHERE ddl.id = 1`, and `SELECT next_row_id FROM t`
// all run on v8.5.5). Each word is read here as an unquoted word instead,
// through adminWord, which keeps the keyword/identifier distinction TiDB has:
// a quoted name is an identifier in TiDB and never one of these words, so
// v8.5.5 answers 1064 for ADMIN SHOW `ddl`, ADMIN `checksum` TABLE t, and
// ADMIN SET BDR ROLE `primary`, and adminWord reports no word for any of them.
//
// Ref: TiDB v8.5.5 parser.y:11204-11501 (AdminStmt and its tails).
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-admin/

// adminWord returns the current token's text, uppercased, when the token is an
// unquoted word: a plain identifier or any keyword, since a word TiDB spells
// as a keyword may be registered in omni, unregistered, or reserved, and the
// spelling is what decides an ADMIN subcommand either way. A quoted name is
// not a word: a backtick-quoted identifier lexes as tokIDENT whose Loc is the
// backtick, so the check is on the byte the token starts at. Everything else
// -- a literal, an operator, end of input -- returns "".
func (p *Parser) adminWord() string {
	return p.unquotedWord(p.cur)
}

// unquotedWord is adminWord for any token, so a lookahead reads a word the
// same way the current token does.
func (p *Parser) unquotedWord(tok Token) string {
	if tok.Type != tokIDENT && tok.Type < 700 {
		return ""
	}
	i := tok.Loc - p.lexer.baseOffset
	if i < 0 || i >= len(p.lexer.input) || !isIdentStart(p.lexer.input[i]) {
		return ""
	}
	return strings.ToUpper(tok.Str)
}

// matchAdminWord consumes the current token when it is the unquoted word w.
func (p *Parser) matchAdminWord(w string) bool {
	if p.adminWord() != w {
		return false
	}
	p.advance()
	return true
}

// expectAdminWord consumes the unquoted word w or fails at the current token.
func (p *Parser) expectAdminWord(w string) error {
	if !p.matchAdminWord(w) {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseAdminStmt parses an ADMIN statement. ADMIN is the current token.
func (p *Parser) parseAdminStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume ADMIN

	stmt := &nodes.AdminStmt{Loc: nodes.Loc{Start: start}}

	var err error
	switch p.adminWord() {
	case "SHOW":
		p.advance()
		err = p.parseAdminShow(stmt)
	case "CHECK":
		p.advance()
		err = p.parseAdminCheck(stmt)
	case "CHECKSUM":
		p.advance()
		stmt.Kind = nodes.AdminChecksumTable
		if err = p.expectAdminWord("TABLE"); err == nil {
			stmt.Tables, err = p.parseAdminTableList()
		}
	case "RECOVER":
		p.advance()
		stmt.Kind = nodes.AdminRecoverIndex
		err = p.parseAdminIndexTarget(stmt)
	case "CLEANUP":
		p.advance()
		err = p.parseAdminCleanup(stmt)
	case "CANCEL":
		p.advance()
		stmt.Kind = nodes.AdminCancelDDLJobs
		err = p.parseAdminDDLJobIDs(stmt)
	case "PAUSE":
		p.advance()
		stmt.Kind = nodes.AdminPauseDDLJobs
		err = p.parseAdminDDLJobIDs(stmt)
	case "RESUME":
		p.advance()
		stmt.Kind = nodes.AdminResumeDDLJobs
		err = p.parseAdminDDLJobIDs(stmt)
	case "ALTER":
		p.advance()
		err = p.parseAdminAlterDDLJobs(stmt)
	case "RELOAD":
		p.advance()
		err = p.parseAdminReload(stmt)
	case "PLUGINS":
		p.advance()
		err = p.parseAdminPlugins(stmt)
	case "FLUSH":
		p.advance()
		err = p.parseAdminFlush(stmt)
	case "CAPTURE":
		p.advance()
		stmt.Kind = nodes.AdminCaptureBindings
		err = p.expectAdminWord("BINDINGS")
	case "EVOLVE":
		p.advance()
		stmt.Kind = nodes.AdminEvolveBindings
		err = p.expectAdminWord("BINDINGS")
	case "SET":
		p.advance()
		stmt.Kind = nodes.AdminSetBDRRole
		if err = p.parseAdminBDRRole(); err == nil {
			stmt.BDRRole, err = p.parseAdminBDRRoleName()
		}
	case "UNSET":
		p.advance()
		stmt.Kind = nodes.AdminUnsetBDRRole
		err = p.parseAdminBDRRole()
	case "REPAIR":
		p.advance()
		err = p.parseAdminRepairTable(stmt)
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAdminShow parses the tail of ADMIN SHOW. SHOW is consumed.
//
// DDL, SLOW, and BDR are words of the grammar here, but each is also an
// unreserved keyword in TiDB and so a table name too, and ADMIN SHOW's last
// form takes a table: v8.5.5 runs `admin show ddl next_row_id` and `admin show
// slow.t next_row_id` as that form, and answers 1064 for `admin show ddl
// extra` and `admin show slow top all next_row_id`. One token of lookahead
// settles it, as it does for TiDB's LALR(1) tables: the word opens its own
// subcommand unless NEXT_ROW_ID or a qualifying '.' follows, which only the
// table form takes.
func (p *Parser) parseAdminShow(stmt *nodes.AdminStmt) error {
	switch w := p.adminWord(); {
	case w == "DDL" && !p.adminShowTableFollows():
		p.advance()
		return p.parseAdminShowDDL(stmt)
	case w == "SLOW" && !p.adminShowTableFollows():
		p.advance()
		stmt.Kind = nodes.AdminShowSlow
		slow, err := p.parseAdminSlow()
		if err != nil {
			return err
		}
		stmt.Slow = slow
		return nil
	case w == "BDR" && !p.adminShowTableFollows():
		p.advance()
		stmt.Kind = nodes.AdminShowBDRRole
		return p.expectAdminWord("ROLE")
	}

	// ADMIN SHOW TableName NEXT_ROW_ID
	stmt.Kind = nodes.AdminShowNextRowID
	ref, err := p.parseTableRef()
	if err != nil {
		return err
	}
	stmt.Tables = []*nodes.TableRef{ref}
	return p.expectAdminWord("NEXT_ROW_ID")
}

// adminShowTableFollows reports whether the token after the current word can
// only continue a table name: NEXT_ROW_ID, which ends ADMIN SHOW's table form,
// or the '.' of a qualified name.
func (p *Parser) adminShowTableFollows() bool {
	next := p.peekNext()
	return next.Type == '.' || p.unquotedWord(next) == "NEXT_ROW_ID"
}

// parseAdminShowDDL parses ADMIN SHOW DDL and its JOB forms. DDL is consumed.
//
//	ADMIN SHOW DDL
//	ADMIN SHOW DDL JOBS [n] [WHERE expr]
//	ADMIN SHOW DDL JOB QUERIES id [, id] ...
//	ADMIN SHOW DDL JOB QUERIES LIMIT n [, n | OFFSET n]
func (p *Parser) parseAdminShowDDL(stmt *nodes.AdminStmt) error {
	switch p.adminWord() {
	case "JOBS":
		p.advance()
		stmt.Kind = nodes.AdminShowDDLJobs
		// Int64Num, so the count is unsigned in the grammar: v8.5.5
		// answers 1064 for `admin show ddl jobs -1`.
		if p.cur.Type == tokICONST {
			n, err := p.parseAdminInt64()
			if err != nil {
				return err
			}
			stmt.JobNumber, stmt.HasJobNumber = n, true
		}
		if p.cur.Type == kwWHERE {
			p.advance()
			where, err := p.parseExpr()
			if err != nil {
				return err
			}
			if where == nil {
				return p.syntaxErrorAtCur()
			}
			stmt.Where = where
		}
		return nil

	case "JOB":
		p.advance()
		if err := p.expectAdminWord("QUERIES"); err != nil {
			return err
		}
		// The id list and the LIMIT tail are separate productions, so
		// neither takes the other: v8.5.5 answers 1064 for `admin show ddl
		// job queries 1 limit 2` and for `admin show ddl job queries limit
		// 2 where id = 1`.
		if p.cur.Type == kwLIMIT {
			stmt.Kind = nodes.AdminShowDDLJobQueriesWithRange
			limit, err := p.parseAdminLimit()
			if err != nil {
				return err
			}
			stmt.Limit = limit
			return nil
		}
		stmt.Kind = nodes.AdminShowDDLJobQueries
		ids, err := p.parseAdminNumList()
		if err != nil {
			return err
		}
		stmt.JobIDs = ids
		return nil
	}

	stmt.Kind = nodes.AdminShowDDL
	return nil
}

// parseAdminLimit parses ADMIN SHOW DDL JOB QUERIES' LIMIT tail. LIMIT is the
// current token and the counts are unsigned, as LengthNum is.
//
//	LIMIT count | LIMIT offset ',' count | LIMIT count OFFSET offset
func (p *Parser) parseAdminLimit() (*nodes.AdminLimit, error) {
	start := p.pos()
	p.advance() // consume LIMIT

	first, err := p.parseAdminUint64()
	if err != nil {
		return nil, err
	}
	limit := &nodes.AdminLimit{Loc: nodes.Loc{Start: start}, Count: first}

	switch p.cur.Type {
	case ',':
		p.advance()
		count, err := p.parseAdminUint64()
		if err != nil {
			return nil, err
		}
		limit.Offset, limit.Count = first, count
	case kwOFFSET:
		p.advance()
		offset, err := p.parseAdminUint64()
		if err != nil {
			return nil, err
		}
		limit.Offset = offset
	}

	limit.Loc.End = p.prev.End
	return limit, nil
}

// parseAdminSlow parses ADMIN SHOW SLOW's tail. Both forms end in a count, so
// neither is optional: v8.5.5 answers 1064 for `admin show slow`, `admin show
// slow top`, and `admin show slow recent`.
//
//	RECENT n | TOP n | TOP INTERNAL n | TOP ALL n
func (p *Parser) parseAdminSlow() (*nodes.AdminSlowSpec, error) {
	start := p.pos()
	slow := &nodes.AdminSlowSpec{Loc: nodes.Loc{Start: start}}

	switch {
	case p.matchAdminWord("RECENT"):
		slow.Recent = true
	case p.matchAdminWord("TOP"):
		switch {
		case p.matchAdminWord("INTERNAL"):
			slow.Kind = nodes.AdminSlowKindInternal
		case p.cur.Type == kwALL:
			p.advance()
			slow.Kind = nodes.AdminSlowKindAll
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	count, err := p.parseAdminUint64()
	if err != nil {
		return nil, err
	}
	slow.Count = count
	slow.Loc.End = p.prev.End
	return slow, nil
}

// parseAdminCheck parses ADMIN CHECK's tail. CHECK is consumed.
//
//	CHECK TABLE tbl [, tbl] ...
//	CHECK INDEX tbl index [handle_range [, handle_range] ...]
func (p *Parser) parseAdminCheck(stmt *nodes.AdminStmt) error {
	if p.matchAdminWord("TABLE") {
		stmt.Kind = nodes.AdminCheckTable
		tables, err := p.parseAdminTableList()
		if err != nil {
			return err
		}
		stmt.Tables = tables
		return nil
	}

	stmt.Kind = nodes.AdminCheckIndex
	if err := p.parseAdminIndexTarget(stmt); err != nil {
		return err
	}
	if p.cur.Type != '(' {
		return nil
	}
	// Only CHECK INDEX takes handle ranges: v8.5.5 answers 1064 for `admin
	// recover index t idx (1,2)` and `admin cleanup index t idx (1,2)`.
	stmt.Kind = nodes.AdminCheckIndexRange
	for {
		rng, err := p.parseAdminHandleRange()
		if err != nil {
			return err
		}
		stmt.Ranges = append(stmt.Ranges, rng)
		if p.cur.Type != ',' {
			return nil
		}
		p.advance()
	}
}

// parseAdminHandleRange parses one '(' begin ',' end ')'. Both bounds are
// unsigned integer literals, so v8.5.5 answers 1064 for `(-1, 2)`, and the
// pair is required: it answers 1064 for `(1)` and `()` too.
func (p *Parser) parseAdminHandleRange() (*nodes.AdminHandleRange, error) {
	start := p.pos()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	begin, err := p.parseAdminInt64()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	end, err := p.parseAdminInt64()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.AdminHandleRange{Loc: nodes.Loc{Start: start, End: p.prev.End}, Begin: begin, End: end}, nil
}

// parseAdminCleanup parses ADMIN CLEANUP's tail. CLEANUP is consumed.
//
//	CLEANUP INDEX tbl index
//	CLEANUP TABLE LOCK tbl [, tbl] ...
func (p *Parser) parseAdminCleanup(stmt *nodes.AdminStmt) error {
	if p.matchAdminWord("TABLE") {
		stmt.Kind = nodes.AdminCleanupTableLock
		if err := p.expectAdminWord("LOCK"); err != nil {
			return err
		}
		tables, err := p.parseAdminTableList()
		if err != nil {
			return err
		}
		stmt.Tables = tables
		return nil
	}

	stmt.Kind = nodes.AdminCleanupIndex
	return p.parseAdminIndexTarget(stmt)
}

// parseAdminIndexTarget parses `INDEX tbl index`, the target of ADMIN CHECK,
// RECOVER, and CLEANUP INDEX. The index name is an Identifier, so a reserved
// keyword needs quoting: v8.5.5 runs `admin check index t `select“ and
// answers 1064 for `admin check index t select`.
func (p *Parser) parseAdminIndexTarget(stmt *nodes.AdminStmt) error {
	if err := p.expectAdminWord("INDEX"); err != nil {
		return err
	}
	ref, err := p.parseTableRef()
	if err != nil {
		return err
	}
	stmt.Tables = []*nodes.TableRef{ref}
	name, _, err := p.parseIdent()
	if err != nil {
		return err
	}
	stmt.Index = name
	return nil
}

// parseAdminDDLJobIDs parses `DDL JOBS id [, id] ...`, the tail of ADMIN
// CANCEL, PAUSE, and RESUME. The list is required and holds integers only:
// v8.5.5 answers 1064 for `admin pause ddl jobs` and `admin pause ddl jobs
// str_not_num`.
func (p *Parser) parseAdminDDLJobIDs(stmt *nodes.AdminStmt) error {
	if err := p.expectAdminWord("DDL"); err != nil {
		return err
	}
	if err := p.expectAdminWord("JOBS"); err != nil {
		return err
	}
	ids, err := p.parseAdminNumList()
	if err != nil {
		return err
	}
	stmt.JobIDs = ids
	return nil
}

// parseAdminAlterDDLJobs parses ADMIN ALTER DDL JOBS. ALTER is consumed.
//
//	ADMIN ALTER DDL JOBS job_id name = value [, name = value] ...
//
// One job id, not a list, and at least one option: v8.5.5 answers 1064 for
// `admin alter ddl jobs 1, 2 thread = 2` and for `admin alter ddl jobs 1`.
func (p *Parser) parseAdminAlterDDLJobs(stmt *nodes.AdminStmt) error {
	stmt.Kind = nodes.AdminAlterDDLJobs
	if err := p.expectAdminWord("DDL"); err != nil {
		return err
	}
	if err := p.expectAdminWord("JOBS"); err != nil {
		return err
	}
	n, err := p.parseAdminInt64()
	if err != nil {
		return err
	}
	stmt.JobNumber, stmt.HasJobNumber = n, true

	for {
		opt, err := p.parseAdminJobOption()
		if err != nil {
			return err
		}
		stmt.JobOptions = append(stmt.JobOptions, opt)
		if p.cur.Type != ',' {
			return nil
		}
		p.advance()
	}
}

// parseAdminJobOption parses one `name = value`. The name is an Identifier,
// which a quoted name is too (v8.5.5 runs `admin alter ddl jobs 1 `thread` =
// 2`), and the value a signed literal. Both the '=' and the value are
// required: v8.5.5 answers 1064 for `admin alter ddl jobs 1 thread` and for
// `admin alter ddl jobs 1 thread =`.
func (p *Parser) parseAdminJobOption() (*nodes.AdminJobOption, error) {
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	value, err := p.parseAdminSignedLiteral()
	if err != nil {
		return nil, err
	}
	return &nodes.AdminJobOption{
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
		Name:  strings.ToLower(name),
		Value: value,
	}, nil
}

// parseAdminSignedLiteral parses TiDB's SignedLiteral: a literal, or a sign
// before a numeric one. The sign takes nothing else -- v8.5.5 answers 1064 for
// `-'a'` and `-0xFF` -- and a literal is the whole grammar here, so an
// expression, a variable, a parenthesized value, and a temporal literal are
// all 1064 as well.
//
// Ref: TiDB v8.5.5 parser.y:16021-16035 (SignedLiteral), 14612-14699 (Literal).
func (p *Parser) parseAdminSignedLiteral() (nodes.ExprNode, error) {
	if p.cur.Type == '+' || p.cur.Type == '-' {
		tok := p.advance()
		if p.cur.Type != tokICONST && p.cur.Type != tokFCONST {
			return nil, p.syntaxErrorAtCur()
		}
		operand, err := p.parseAdminLiteral()
		if err != nil {
			return nil, err
		}
		if tok.Type == '+' {
			return operand, nil
		}
		return &nodes.UnaryExpr{
			Loc:        nodes.Loc{Start: tok.Loc, End: p.prev.End},
			Op:         nodes.UnaryMinus,
			Operand:    operand,
			OriginalOp: "-",
		}, nil
	}
	return p.parseAdminLiteral()
}

// parseAdminLiteral parses TiDB's Literal. Adjacent string literals are one
// string, as they are in TiDB (v8.5.5 takes `'a' 'b'`), and a charset
// introducer or the national-string N prefix applies to the string that
// follows it. A charset introducer before a hex or bit literal only sets that
// literal's collation, which omni's HexLit and BitLit do not carry, so the
// introducer is read and dropped.
func (p *Parser) parseAdminLiteral() (nodes.ExprNode, error) {
	start := p.pos()

	// _charset'str' / _charset x'41', and N'str'. The lexer leaves the
	// introducer as a plain identifier.
	charset := ""
	if p.cur.Type == tokIDENT && p.peekNext().Type != tokEOF {
		switch next := p.peekNext().Type; {
		case strings.HasPrefix(p.cur.Str, "_") && (next == tokSCONST || next == tokXCONST || next == tokBCONST):
			charset = p.advance().Str
		case strings.EqualFold(p.cur.Str, "N") && next == tokSCONST:
			p.advance()
			charset = "binary"
		}
	}

	switch p.cur.Type {
	case tokSCONST:
		var sb strings.Builder
		for p.cur.Type == tokSCONST {
			sb.WriteString(p.advance().Str)
		}
		return &nodes.StringLit{Loc: nodes.Loc{Start: start, End: p.prev.End}, Value: sb.String(), Charset: charset}, nil
	case tokXCONST:
		tok := p.advance()
		return &nodes.HexLit{Loc: nodes.Loc{Start: start, End: p.prev.End}, Value: tok.Str}, nil
	case tokBCONST:
		tok := p.advance()
		return &nodes.BitLit{Loc: nodes.Loc{Start: start, End: p.prev.End}, Value: tok.Str}, nil
	}
	if charset != "" {
		// An introducer with nothing it can apply to.
		return nil, p.syntaxErrorAtCur()
	}

	switch p.cur.Type {
	case tokICONST:
		tok := p.advance()
		return &nodes.IntLit{Loc: nodes.Loc{Start: start, End: p.prev.End}, Value: tok.Ival}, nil
	case tokFCONST:
		tok := p.advance()
		return &nodes.FloatLit{Loc: nodes.Loc{Start: start, End: p.prev.End}, Value: tok.Str}, nil
	case kwTRUE:
		p.advance()
		return &nodes.BoolLit{Loc: nodes.Loc{Start: start, End: p.prev.End}, Value: true}, nil
	case kwFALSE:
		p.advance()
		return &nodes.BoolLit{Loc: nodes.Loc{Start: start, End: p.prev.End}, Value: false}, nil
	case kwNULL:
		p.advance()
		return &nodes.NullLit{Loc: nodes.Loc{Start: start, End: p.prev.End}}, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseAdminReload parses ADMIN RELOAD's tail. RELOAD is consumed. The target
// is required and is one of five words: v8.5.5 answers 1064 for `admin reload`
// and `admin reload foo`.
func (p *Parser) parseAdminReload(stmt *nodes.AdminStmt) error {
	switch p.adminWord() {
	case "EXPR_PUSHDOWN_BLACKLIST":
		stmt.Kind = nodes.AdminReloadExprPushdownBlacklist
	case "OPT_RULE_BLACKLIST":
		stmt.Kind = nodes.AdminReloadOptRuleBlacklist
	// TiDB gives STATISTICS and STATS_EXTENDED the same kind: both reload
	// the extended statistics.
	case "STATISTICS", "STATS_EXTENDED":
		stmt.Kind = nodes.AdminReloadStatistics
	case "BINDINGS":
		stmt.Kind = nodes.AdminReloadBindings
	default:
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseAdminPlugins parses ADMIN PLUGINS. PLUGINS is consumed. The names are
// Identifiers, comma-separated, and at least one is required: v8.5.5 answers
// 1064 for `admin plugins enable` and `admin plugins enable audit whitelist`.
func (p *Parser) parseAdminPlugins(stmt *nodes.AdminStmt) error {
	switch p.cur.Type {
	case kwENABLE:
		stmt.Kind = nodes.AdminPluginEnable
	case kwDISABLE:
		stmt.Kind = nodes.AdminPluginDisable
	default:
		return p.syntaxErrorAtCur()
	}
	p.advance()

	for {
		name, _, err := p.parseIdent()
		if err != nil {
			return err
		}
		stmt.Plugins = append(stmt.Plugins, name)
		if p.cur.Type != ',' {
			return nil
		}
		p.advance()
	}
}

// parseAdminFlush parses ADMIN FLUSH's tail. FLUSH is consumed.
//
//	ADMIN FLUSH BINDINGS
//	ADMIN FLUSH [SESSION | GLOBAL | INSTANCE] PLAN_CACHE
//
// The scope is optional and means SESSION when left out, as TiDB's
// StatementScope does; v8.5.5 runs `admin flush plan_cache` and answers 1064
// for `admin flush local plan_cache` and `admin flush session`.
func (p *Parser) parseAdminFlush(stmt *nodes.AdminStmt) error {
	if p.matchAdminWord("BINDINGS") {
		stmt.Kind = nodes.AdminFlushBindings
		return nil
	}

	stmt.Kind = nodes.AdminFlushPlanCache
	switch p.cur.Type {
	case kwSESSION:
		p.advance()
		stmt.Scope = nodes.AdminScopeSession
	case kwGLOBAL:
		p.advance()
		stmt.Scope = nodes.AdminScopeGlobal
	case kwINSTANCE:
		p.advance()
		stmt.Scope = nodes.AdminScopeInstance
	}
	return p.expectAdminWord("PLAN_CACHE")
}

// parseAdminBDRRole consumes the `BDR ROLE` of ADMIN SET, SHOW, and UNSET BDR
// ROLE.
func (p *Parser) parseAdminBDRRole() error {
	if err := p.expectAdminWord("BDR"); err != nil {
		return err
	}
	return p.expectAdminWord("ROLE")
}

// parseAdminBDRRoleName parses ADMIN SET BDR ROLE's role. The grammar takes
// the two keywords and nothing else, so v8.5.5 answers 1064 for `admin set bdr
// role foo` and for the quoted `admin set bdr role `primary“.
func (p *Parser) parseAdminBDRRoleName() (string, error) {
	switch p.cur.Type {
	case kwPRIMARY:
		p.advance()
		return "PRIMARY", nil
	case kwSECONDARY:
		p.advance()
		return "SECONDARY", nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseAdminRepairTable parses ADMIN REPAIR TABLE. REPAIR is consumed.
//
//	ADMIN REPAIR TABLE tbl_name CreateTableStmt
//
// The tail is a whole CREATE TABLE, in every form CREATE TABLE takes: v8.5.5
// accepts IF NOT EXISTS, TEMPORARY, LIKE, and AS SELECT after it, and answers
// 1064 for `ADMIN REPAIR TABLE t CREATE VIEW v AS SELECT 1`. Outside repair
// mode the statement itself answers 8215, whatever it carries.
func (p *Parser) parseAdminRepairTable(stmt *nodes.AdminStmt) error {
	stmt.Kind = nodes.AdminRepairTable
	if err := p.expectAdminWord("TABLE"); err != nil {
		return err
	}
	ref, err := p.parseTableRef()
	if err != nil {
		return err
	}
	stmt.Tables = []*nodes.TableRef{ref}

	if p.cur.Type != kwCREATE {
		return p.syntaxErrorAtCur()
	}
	created, err := p.parseCreateDispatch()
	if err != nil {
		return err
	}
	create, ok := created.(*nodes.CreateTableStmt)
	if !ok {
		return &ParseError{Message: "ADMIN REPAIR TABLE takes a CREATE TABLE statement", Position: ref.Loc.Start}
	}
	stmt.CreateTable = create
	return nil
}

// parseAdminTableList parses `tbl [, tbl] ...`, at least one name.
func (p *Parser) parseAdminTableList() ([]*nodes.TableRef, error) {
	var refs []*nodes.TableRef
	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
		if p.cur.Type != ',' {
			return refs, nil
		}
		p.advance()
	}
}

// parseAdminNumList parses `n [, n] ...`, at least one integer.
func (p *Parser) parseAdminNumList() ([]int64, error) {
	var nums []int64
	for {
		n, err := p.parseAdminInt64()
		if err != nil {
			return nil, err
		}
		nums = append(nums, n)
		if p.cur.Type != ',' {
			return nums, nil
		}
		p.advance()
	}
}

// parseAdminInt64 parses TiDB's Int64Num: an unsigned integer literal that
// fits in an int64. The grammar's own range check makes a larger one a parse
// error, so v8.5.5 answers 1064 for `admin show ddl jobs 9223372036854775808`.
func (p *Parser) parseAdminInt64() (int64, error) {
	if p.cur.Type != tokICONST {
		return 0, p.syntaxErrorAtCur()
	}
	tok := p.cur
	u, err := strconv.ParseUint(tok.Str, 10, 64)
	if err != nil || u > uint64(math.MaxInt64) {
		return 0, &ParseError{
			Message:  "number out of range [0, 9223372036854775807]",
			Position: tok.Loc,
		}
	}
	p.advance()
	return int64(u), nil
}

// parseAdminUint64 parses TiDB's LengthNum: an unsigned integer literal that
// fits in a uint64.
func (p *Parser) parseAdminUint64() (uint64, error) {
	if p.cur.Type != tokICONST {
		return 0, p.syntaxErrorAtCur()
	}
	tok := p.cur
	u, err := strconv.ParseUint(tok.Str, 10, 64)
	if err != nil {
		return 0, &ParseError{
			Message:  "number out of range [0, 18446744073709551615]",
			Position: tok.Loc,
		}
	}
	p.advance()
	return u, nil
}
