// Package parser implements a hand-written recursive-descent HiveQL parser
// targeting Apache Hive 4.0.1 syntax.
package parser

import (
	"fmt"
	"strings"
)

// ParseError describes a lexical or syntactic error encountered while parsing
// HiveQL.
type ParseError struct {
	Message    string
	Position   int
	depthLimit bool
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at byte %d: %s", e.Position, e.Message)
}

// StatementKind classifies a top-level HiveQL statement.
type StatementKind string

const (
	StmtExplain StatementKind = "EXPLAIN"
	StmtQuery   StatementKind = "QUERY"
	StmtInsert  StatementKind = "INSERT"
	StmtUpdate  StatementKind = "UPDATE"
	StmtDelete  StatementKind = "DELETE"
	StmtMerge   StatementKind = "MERGE"
	StmtLoad    StatementKind = "LOAD"
	StmtExport  StatementKind = "EXPORT"
	StmtImport  StatementKind = "IMPORT"
	StmtRepl    StatementKind = "REPL"
	StmtTx      StatementKind = "TRANSACTION"
	StmtPrepare StatementKind = "PREPARE"
	StmtExecute StatementKind = "EXECUTE"
	StmtDDL     StatementKind = "DDL"
)

// ExplainInfo holds parsed metadata for an EXPLAIN statement.
type ExplainInfo struct {
	Formatted bool
	Analyze   bool
	Rewrite   bool
	InnerKind StatementKind
}

// Statement holds metadata for one parsed HiveQL statement.
type Statement struct {
	Kind         StatementKind
	Explain      *ExplainInfo // non-nil when Kind == StmtExplain
	FirstKeyword TokenType
	HasTransform bool
	Text         string
	Loc          Loc
}

// ParseResult is the result of parsing a HiveQL statement.
type ParseResult struct {
	Statement *Statement
}

// Parse parses a single HiveQL statement according to Apache Hive 4.0.1 grammar
// rules (`statement : explainStatement EOF | execStatement EOF`).
func Parse(sql string) (*ParseResult, error) {
	toks, err := lex(sql)
	if err != nil {
		return nil, err
	}
	p := &parser{
		sql:  sql,
		toks: toks,
	}
	kind, explainInfo, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	if p.peek().Type != TokEOF {
		return nil, p.errAt(p.peek(), fmt.Sprintf("unexpected token %q", p.peek().Str))
	}
	start := 0
	end := len(sql)
	var firstKW TokenType
	if len(toks) > 1 {
		start = toks[0].Loc.Start
		end = toks[len(toks)-2].Loc.End
		firstKW = toks[0].Type
	}
	return &ParseResult{
		Statement: &Statement{
			Kind:         kind,
			Explain:      explainInfo,
			FirstKeyword: firstKW,
			HasTransform: p.hasTransform,
			Text:         sql,
			Loc:          Loc{Start: start, End: end},
		},
	}, nil
}

const maxDepth = 256

type memoRule uint8

const (
	memoExpression memoRule = iota + 1
	memoSelectStatement
	memoQueryStatement
)

type memoKey struct {
	rule memoRule
	pos  int
}

type memoEntry struct {
	endPos       int
	relDepth     int
	hasInsert    bool
	hasTransform bool
	err          error
}

type parser struct {
	sql          string
	toks         []Token
	pos          int
	depth        int
	peakDepth    int
	inTry        int
	hasInsert    bool
	hasTransform bool
	memo         map[memoKey]memoEntry
}

func (p *parser) enter() error {
	p.depth++
	if p.depth > p.peakDepth {
		p.peakDepth = p.depth
	}
	if p.depth > maxDepth {
		return &ParseError{
			Message:    "maximum nesting depth exceeded",
			Position:   p.peek().Loc.Start,
			depthLimit: true,
		}
	}
	return nil
}

func (p *parser) leave() {
	p.depth--
}

// memoize wraps a recursive entry point with depth checking and packrat
// memoization during speculative parses (p.try), preventing exponential
// backtracking on nested expressions and subqueries while adding zero
// allocations when no backtracking occurs.
func (p *parser) memoize(rule memoRule, fn func() error) error {
	if err := p.enter(); err != nil {
		return err
	}
	defer p.leave()

	startPos := p.pos
	if p.memo != nil {
		if entry, ok := p.memo[memoKey{rule: rule, pos: startPos}]; ok {
			if entry.err != nil {
				return entry.err
			}
			if p.depth+entry.relDepth <= maxDepth {
				if p.depth+entry.relDepth > p.peakDepth {
					p.peakDepth = p.depth + entry.relDepth
				}
				p.pos = entry.endPos
				p.hasInsert = p.hasInsert || entry.hasInsert
				p.hasTransform = p.hasTransform || entry.hasTransform
				return nil
			}
		}
	}

	prevPeak := p.peakDepth
	p.peakDepth = p.depth
	prevInsert, prevTransform := p.hasInsert, p.hasTransform
	p.hasInsert, p.hasTransform = false, false
	err := fn()
	subInsert, subTransform := p.hasInsert, p.hasTransform
	p.hasInsert = prevInsert || subInsert
	p.hasTransform = prevTransform || subTransform
	relDepth := p.peakDepth - p.depth
	if prevPeak > p.peakDepth {
		p.peakDepth = prevPeak
	}

	if p.inTry > 0 {
		if pe, ok := err.(*ParseError); !ok || !pe.depthLimit {
			if p.memo == nil {
				p.memo = make(map[memoKey]memoEntry)
			}
			p.memo[memoKey{rule: rule, pos: startPos}] = memoEntry{
				endPos:       p.pos,
				relDepth:     relDepth,
				hasInsert:    subInsert,
				hasTransform: subTransform,
				err:          err,
			}
		}
	}
	return err
}

func (p *parser) save() int {
	return p.pos
}

func (p *parser) restore(pos int) {
	p.pos = pos
}

func (p *parser) peek() Token {
	if p.pos < len(p.toks) {
		return p.toks[p.pos]
	}
	return p.toks[len(p.toks)-1]
}

func (p *parser) peekAt(offset int) Token {
	idx := p.pos + offset
	if idx >= 0 && idx < len(p.toks) {
		return p.toks[idx]
	}
	return p.toks[len(p.toks)-1]
}

func (p *parser) next() Token {
	t := p.peek()
	if p.pos < len(p.toks)-1 {
		p.pos++
	}
	return t
}

func (p *parser) match(tt TokenType) bool {
	if p.peek().Type == tt {
		p.next()
		return true
	}
	return false
}

func (p *parser) expect(tt TokenType) (Token, error) {
	t := p.peek()
	if t.Type != tt {
		return Token{}, p.errAt(t, fmt.Sprintf("unexpected token %q", t.Str))
	}
	p.next()
	return t, nil
}

func (p *parser) errAt(t Token, msg string) error {
	return &ParseError{Message: msg, Position: t.Loc.Start}
}

func (p *parser) try(fn func() error) bool {
	saved := p.save()
	savedInsert, savedTransform := p.hasInsert, p.hasTransform
	p.inTry++
	err := fn()
	p.inTry--
	if err == nil {
		return true
	}
	p.restore(saved)
	p.hasInsert, p.hasTransform = savedInsert, savedTransform
	return false
}

func (p *parser) parseIdentifier() (Token, error) {
	t := p.peek()
	if isIdentifier(t.Type) {
		p.next()
		return t, nil
	}
	return Token{}, p.errAt(t, fmt.Sprintf("expected identifier, got %q", t.Str))
}

func (p *parser) parseFunctionIdentifier() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.match(TokDot) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseTableName() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.peek().Type == TokDot && isIdentifier(p.peekAt(1).Type) {
		p.next()
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if p.peek().Type == TokDot && isIdentifier(p.peekAt(1).Type) {
			p.next()
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *parser) parseViewName() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.match(TokDot) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	}
	return nil
}

// checkCreateTableColName mirrors HiveParser.g's
// containExcludedCharForCreateTableColumnName (rejects '.' and ':' in column names).
func (p *parser) checkCreateTableColName(t Token) error {
	if strings.ContainsAny(t.Str, ".:") {
		return p.errAt(t, "'.' and ':' cannot be used in column name in create table statement")
	}
	return nil
}

// statement : explainStatement EOF | execStatement EOF ;
func (p *parser) parseStatement() (StatementKind, *ExplainInfo, error) {
	if err := p.enter(); err != nil {
		return "", nil, err
	}
	defer p.leave()

	if p.peek().Type == KW_EXPLAIN {
		info, err := p.parseExplainStatement()
		if err != nil {
			return "", nil, err
		}
		return StmtExplain, info, nil
	}
	kind, err := p.parseExecStatement()
	return kind, nil, err
}

// explainStatement : KW_EXPLAIN ( explainOption* execStatement | KW_REWRITE queryStatementExpression ) ;
func (p *parser) parseExplainStatement() (*ExplainInfo, error) {
	if _, err := p.expect(KW_EXPLAIN); err != nil {
		return nil, err
	}
	info := &ExplainInfo{}
	if p.match(KW_REWRITE) {
		info.Rewrite = true
		if err := p.parseQueryStatementExpression(); err != nil {
			return nil, err
		}
		if p.hasInsert {
			info.InnerKind = StmtInsert
		} else {
			info.InnerKind = StmtQuery
		}
		return info, nil
	}
	for p.parseExplainOption(info) {
	}
	innerKind, err := p.parseExecStatement()
	if err != nil {
		return nil, err
	}
	info.InnerKind = innerKind
	return info, nil
}

func (p *parser) parseExplainOption(info *ExplainInfo) bool {
	switch p.peek().Type {
	case KW_FORMATTED:
		info.Formatted = true
		p.next()
		return true
	case KW_EXTENDED, KW_DEPENDENCY, KW_LOGICAL,
		KW_AUTHORIZATION, KW_REOPTIMIZATION, KW_LOCKS, KW_AST, KW_DEBUG, KW_DDL:
		p.next()
		return true
	case KW_ANALYZE:
		// ANALYZE TABLE ... starts analyzeStatement rather than explainOption.
		if p.peekAt(1).Type == KW_TABLE {
			return false
		}
		info.Analyze = true
		p.next()
		return true
	case KW_CBO:
		p.next()
		if p.peek().Type == KW_COST || p.peek().Type == KW_JOINCOST {
			p.next()
		}
		return true
	case KW_VECTORIZATION:
		p.next()
		p.match(KW_ONLY)
		switch p.peek().Type {
		case KW_SUMMARY, KW_OPERATOR, KW_EXPRESSION, KW_DETAIL:
			p.next()
		}
		return true
	default:
		return false
	}
}

// execStatement dispatches all non-EXPLAIN HiveQL statements.
func (p *parser) parseExecStatement() (StatementKind, error) {
	switch p.peek().Type {
	case KW_WITH, KW_FROM, KW_SELECT, KW_MAP, KW_REDUCE, KW_VALUES, TokLParen:
		if err := p.parseQueryStatementExpression(); err != nil {
			return "", err
		}
		if p.hasInsert {
			return StmtInsert, nil
		}
		return StmtQuery, nil
	case KW_INSERT:
		return StmtInsert, p.parseQueryStatementExpression()
	case KW_LOAD:
		return StmtLoad, p.parseLoadStatement()
	case KW_EXPORT:
		return StmtExport, p.parseExportStatement()
	case KW_IMPORT:
		return StmtImport, p.parseImportStatement()
	case KW_REPL:
		return StmtRepl, p.parseReplStatement()
	case KW_DELETE:
		return StmtDelete, p.parseDeleteStatement()
	case KW_UPDATE:
		return StmtUpdate, p.parseUpdateStatement()
	case KW_START, KW_COMMIT, KW_ROLLBACK:
		return StmtTx, p.parseSQLTransactionStatement()
	case KW_SET:
		if p.peekAt(1).Type == KW_AUTOCOMMIT {
			return StmtTx, p.parseSetAutoCommitStatement()
		}
		return StmtDDL, p.parseSetRole()
	case KW_MERGE:
		return StmtMerge, p.parseMergeStatement()
	case KW_PREPARE:
		return StmtPrepare, p.parsePrepareStatement()
	case KW_EXECUTE:
		return StmtExecute, p.parseExecuteStatement()
	case KW_CREATE, KW_USE, KW_DROP, KW_TRUNCATE, KW_ALTER, KW_OPTIMIZE,
		KW_DESCRIBE, KW_DESC, KW_SHOW, KW_MSCK, KW_RELOAD, KW_ANALYZE,
		KW_LOCK, KW_UNLOCK, KW_GRANT, KW_REVOKE, KW_ABORT, KW_KILL,
		KW_ENABLE, KW_DISABLE, KW_REPLACE:
		return StmtDDL, p.parseDDLStatement()
	default:
		return "", p.errAt(p.peek(), fmt.Sprintf("unexpected token %q at start of statement", p.peek().Str))
	}
}
