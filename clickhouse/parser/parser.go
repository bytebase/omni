// Package parser implements a hand-written recursive-descent SQL parser for
// ClickHouse.
package parser

import (
	"fmt"
	"strings"
)

// Parse parses a single ClickHouse SQL statement (with optional trailing
// semicolon) and returns its AST. On lexical or syntax error, it returns a
// *ParseError carrying the byte offset.
func Parse(sql string) ([]Statement, error) {
	p := &parser{
		src: sql,
		lex: &lexer{src: sql},
	}
	if p.cur().Type == TokIllegal && p.lexErr != nil {
		return nil, p.lexErr
	}
	if p.cur().Type == TokEOF {
		return nil, p.errAt(p.cur(), "empty query")
	}
	stmt, err := p.parseTopLevel()
	if err != nil {
		return nil, err
	}
	if p.cur().Type == TokSemicolon {
		p.advance()
	}
	if p.cur().Type == TokIllegal && p.lexErr != nil {
		return nil, p.lexErr
	}
	if p.cur().Type != TokEOF {
		return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q", p.cur().Raw))
	}
	if treeTooHigh(stmt, len(sql)) {
		return nil, &ParseError{Message: nestingErrorMessage, Position: stmt.StartPos()}
	}
	return []Statement{stmt}, nil
}

type parser struct {
	src            string
	lex            *lexer
	tokens         []Token
	pos            int
	depth          int
	eof            bool
	lexErr         error
	noLambdaBefore int          // token index bound caching failed unparenthesized lambda lookahead
	noIntervalAt   map[int]bool // token positions where INTERVAL expr unit lookahead failed
	noTypeArgAt    map[int]bool // token positions where nested DataType lookahead failed
	noCaseAt       map[int]bool // token positions where CASE operand WHEN lookahead failed
}

func (p *parser) ensure(idx int) {
	for len(p.tokens) <= idx && !p.eof {
		tok, err := p.lex.nextToken()
		if err != nil {
			p.lexErr = err
			pos := p.lex.pos
			if pe, ok := err.(*ParseError); ok {
				pos = pe.Position
			}
			p.tokens = append(p.tokens, Token{Type: TokIllegal, Pos: pos, End: pos})
			p.eof = true
			return
		}
		p.tokens = append(p.tokens, tok)
		if tok.Type == TokEOF {
			p.eof = true
		}
	}
}

// consumeRawTail stops SQL tokenization at the current position and treats all
// remaining bytes as an opaque raw payload (used for INSERT FORMAT/VALUES data).
func (p *parser) consumeRawTail(rawStart int) string {
	data := ""
	if rawStart < len(p.src) {
		data = p.src[rawStart:]
	}
	p.tokens = append(p.tokens[:p.pos], Token{Type: TokEOF, Pos: len(p.src), End: len(p.src)})
	p.eof = true
	p.lexErr = nil
	return data
}

func (p *parser) cur() Token {
	p.ensure(p.pos)
	if p.pos >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.pos]
}

func (p *parser) peek(offset int) Token {
	idx := p.pos + offset
	if idx < 0 {
		p.ensure(0)
		return p.tokens[0]
	}
	p.ensure(idx)
	if idx >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[idx]
}

func (p *parser) advance() Token {
	tok := p.cur()
	if tok.Type != TokEOF && tok.Type != TokIllegal {
		p.pos++
	}
	return tok
}

func (p *parser) expect(tt TokenType) (Token, error) {
	tok := p.cur()
	if tok.Type != tt {
		return Token{}, p.errAt(tok, fmt.Sprintf("unexpected token %q", tok.Raw))
	}
	return p.advance(), nil
}

func (p *parser) match(tt TokenType) bool {
	if p.cur().Type == tt {
		p.advance()
		return true
	}
	return false
}

func (p *parser) errAt(tok Token, msg string) error {
	if tok.Type == TokIllegal && p.lexErr != nil {
		return p.lexErr
	}
	return &ParseError{Message: msg, Position: tok.Pos}
}

func (p *parser) parseTopLevel() (Statement, error) {
	switch p.cur().Type {
	case KwInsert:
		return p.parseInsertStmt()
	case KwDelete:
		return p.parseDeleteStmt()
	case KwUpdate:
		return p.parseUpdateStmt()
	}
	stmt, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	out, err := p.parseOptionalOutputClause()
	if err != nil {
		return nil, err
	}
	attachOutputClause(stmt, out)
	return stmt, nil
}

func (p *parser) parseQuery() (Statement, error) {
	if err := p.enter(); err != nil {
		return nil, err
	}
	defer p.leave()

	switch p.cur().Type {
	case KwSelect, KwWith, KwFrom, TokLParen:
		return p.parseSelectUnionStmt()
	case KwCreate, KwReplace:
		return p.parseCreateStmt()
	case KwAttach:
		return p.parseAttachOrCreateStmt()
	case KwAlter:
		return p.parseAlterStmt()
	case KwDrop, KwDetach:
		return p.parseDropStmt()
	case KwTruncate:
		return p.parseTruncateStmt()
	case KwRename:
		return p.parseRenameStmt()
	case KwExchange:
		return p.parseExchangeStmt()
	case KwUndrop:
		return p.parseUndropStmt()
	case KwShow:
		return p.parseShowStmt()
	case KwDescribe, KwDesc:
		return p.parseDescribeStmt()
	case KwExplain:
		return p.parseExplainStmt()
	case KwExists:
		return p.parseExistsStmt()
	case KwCheck:
		return p.parseCheckStmt()
	case KwOptimize:
		return p.parseOptimizeStmt()
	case KwKill:
		return p.parseKillStmt()
	case KwSystem:
		return p.parseSystemStmt()
	case KwSet:
		return p.parseSetStmt()
	case KwUse:
		return p.parseUseStmt()
	case KwWatch:
		return p.parseWatchStmt()
	case KwGrant, KwRevoke:
		return p.parseGrantRevokeStmt()
	case KwBackup, KwRestore:
		return p.parseBackupRestoreStmt()
	default:
		return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q at start of query", p.cur().Raw))
	}
}

func (p *parser) parseOptionalOutputClause() (*OutputClause, error) {
	var out *OutputClause
	seenOutfile := false
	seenFormat := false
	seenSettings := false
	for {
		switch {
		case !seenOutfile && p.cur().Type == KwInto && p.peek(1).Type == KwOutfile:
			seenOutfile = true
			start := p.advance().Pos
			p.advance()
			str, err := p.parseStringLiteral()
			if err != nil {
				return nil, err
			}
			end := str.End
			comp := ""
			if p.match(KwCompression) {
				c, err := p.parseStringLiteral()
				if err != nil {
					return nil, err
				}
				comp = c.Val
				end = c.End
			}
			if out == nil {
				out = &OutputClause{span: span{Pos: start, End: end}}
			} else {
				out.End = end
			}
			out.Outfile = str.Val
			out.Compression = comp
		case !seenFormat && p.match(KwFormat):
			seenFormat = true
			fmtStart := p.peek(-1).Pos
			name, err := p.parseIdentifierOrNull()
			if err != nil {
				return nil, err
			}
			end := p.peek(-1).End
			if out == nil {
				out = &OutputClause{span: span{Pos: fmtStart, End: end}}
			} else {
				out.End = end
			}
			out.Format = name
		case !seenSettings && p.cur().Type == KwSettings:
			seenSettings = true
			setStart := p.cur().Pos
			settings, err := p.parseSettingsClause()
			if err != nil {
				return nil, err
			}
			end := p.peek(-1).End
			if out == nil {
				out = &OutputClause{span: span{Pos: setStart, End: end}}
			} else {
				out.End = end
			}
			out.Settings = settings
		default:
			return out, nil
		}
	}
}

func attachOutputClause(stmt Statement, out *OutputClause) {
	if out == nil {
		return
	}
	switch s := stmt.(type) {
	case *SelectUnionStmt:
		s.Output, s.End = out, out.End
	case *CreateDatabaseStmt:
		s.Output, s.End = out, out.End
	case *CreateTableStmt:
		s.Output, s.End = out, out.End
	case *CreateViewStmt:
		s.Output, s.End = out, out.End
	case *CreateDictionaryStmt:
		s.Output, s.End = out, out.End
	case *CreateFunctionStmt:
		s.Output, s.End = out, out.End
	case *AlterTableStmt:
		s.Output, s.End = out, out.End
	case *DropStmt:
		s.Output, s.End = out, out.End
	case *AttachStmt:
		s.Output, s.End = out, out.End
	case *RenameStmt:
		s.Output, s.End = out, out.End
	case *ExchangeStmt:
		s.Output, s.End = out, out.End
	case *TruncateStmt:
		s.Output, s.End = out, out.End
	case *UndropStmt:
		s.Output, s.End = out, out.End
	case *ShowStmt:
		s.Output, s.End = out, out.End
	case *DescribeStmt:
		s.Output, s.End = out, out.End
	case *ExplainStmt:
		s.Output, s.End = out, out.End
	case *ExistsStmt:
		s.Output, s.End = out, out.End
	case *CheckStmt:
		s.Output, s.End = out, out.End
	case *OptimizeStmt:
		s.Output, s.End = out, out.End
	case *KillStmt:
		s.Output, s.End = out, out.End
	case *SystemStmt:
		s.Output, s.End = out, out.End
	case *SetStmt:
		s.Output, s.End = out, out.End
	case *UseStmt:
		s.Output, s.End = out, out.End
	case *WatchStmt:
		s.Output, s.End = out, out.End
	case *AccessControlStmt:
		s.Output, s.End = out, out.End
	case *BackupRestoreStmt:
		s.Output, s.End = out, out.End
	}
}

// --- Identifier helpers ---

// isIdentifier reports whether tok can be used as a ClickHouse identifier.
// In ClickHouse, almost all keywords (except NULL, INF/INFINITY, NAN) are valid
// unquoted identifiers in identifier positions.
func isIdentifier(tok Token) bool {
	switch tok.Type {
	case TokIdent, TokQuotedIdent:
		return true
	case KwNull, KwInf, KwNan:
		return false
	default:
		return tok.IsKeyword()
	}
}

func isIdentifierOrNull(tok Token) bool {
	return isIdentifier(tok) || tok.Type == KwNull
}

func (p *parser) parseIdentifier() (string, error) {
	tok := p.cur()
	if !isIdentifier(tok) {
		return "", p.errAt(tok, fmt.Sprintf("expected identifier, got %q", tok.Raw))
	}
	p.advance()
	return tok.Val, nil
}

func (p *parser) parseIdentifierOrNull() (string, error) {
	tok := p.cur()
	if !isIdentifierOrNull(tok) {
		return "", p.errAt(tok, fmt.Sprintf("expected identifier, got %q", tok.Raw))
	}
	p.advance()
	return tok.Val, nil
}

func (p *parser) parseNestedIdentifier() (string, error) {
	first, err := p.parseIdentifier()
	if err != nil {
		return "", err
	}
	if p.cur().Type != TokDot {
		return first, nil
	}
	var b strings.Builder
	b.WriteString(first)
	for p.match(TokDot) {
		part, err := p.parseIdentifier()
		if err != nil {
			return "", err
		}
		b.WriteByte('.')
		b.WriteString(part)
	}
	return b.String(), nil
}

func (p *parser) parseTableIdentifier() (db, table string, err error) {
	first, err := p.parseIdentifier()
	if err != nil {
		return "", "", err
	}
	if p.match(TokDot) {
		second, err := p.parseIdentifier()
		if err != nil {
			return "", "", err
		}
		return first, second, nil
	}
	return "", first, nil
}

// isBareAlias reports whether tok can be used as an implicit alias (without AS).
// Matches ClickHouse's alias rules: clause-introducing keywords and interval
// units cannot act as implicit aliases.
func isBareAlias(tok Token) bool {
	if tok.Type == TokIdent || tok.Type == TokQuotedIdent {
		return true
	}
	if !isIdentifier(tok) {
		return false
	}
	switch tok.Type {
	// SQL clause / operator / expression keywords that cannot be bare aliases
	case KwAll, KwAnd, KwAnti, KwAny, KwArray, KwAs, KwAsc, KwAsof,
		KwBetween, KwBy, KwCross, KwDesc, KwDescending, KwDistinct,
		KwElse, KwEnd, KwExcept, KwFinal, KwFor, KwFormat,
		KwFrom, KwFull, KwGlobal, KwGroup, KwGrouping, KwHaving,
		KwIlike, KwIn, KwInner, KwInterpolate, KwIntersect, KwInterval,
		KwInto, KwIs, KwJoin, KwLast, KwLeft, KwLike, KwLimit, KwLocal,
		KwNot, KwNulls, KwOffset, KwOn, KwOr, KwOrder, KwOuter,
		KwPrewhere, KwQualify, KwRight, KwSample, KwSemi, KwSettings,
		KwStep, KwThen, KwTo, KwTop, KwUnion, KwUsing, KwWhen, KwWhere,
		KwWindow, KwWith,
		// Interval units (ambiguous with INTERVAL 1 DAY etc.)
		KwNanosecond, KwMicrosecond, KwMillisecond, KwSecond,
		KwMinute, KwHour, KwDay, KwWeek, KwMonth, KwQuarter, KwYear,
		// Literals
		KwTrue, KwFalse, KwNull, KwInf, KwNan:
		return false
	default:
		return true
	}
}

func isIntervalUnit(tt TokenType) bool {
	switch tt {
	case KwNanosecond, KwMicrosecond, KwMillisecond, KwSecond,
		KwMinute, KwHour, KwDay, KwWeek, KwMonth, KwQuarter, KwYear:
		return true
	default:
		return false
	}
}

func (p *parser) parseClusterClause() (string, error) {
	if !(p.cur().Type == KwOn && p.peek(1).Type == KwCluster) {
		return "", nil
	}
	p.advance()
	p.advance()
	if isStringTok(p.cur().Type) {
		lit, err := p.parseStringLiteral()
		if err != nil {
			return "", err
		}
		return lit.Val, nil
	}
	return p.parseIdentifier()
}

func (p *parser) parseSettingExprList() ([]*SettingExpr, error) {
	return p.parseSettingExprListMode(false)
}

func (p *parser) parseAlterSettingExprList() ([]*SettingExpr, error) {
	return p.parseSettingExprListMode(true)
}

func (p *parser) parseSettingExprListMode(alterAware bool) ([]*SettingExpr, error) {
	var list []*SettingExpr
	for {
		start := p.cur().Pos
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(TokEq); err != nil {
			return nil, err
		}
		val, err := p.parseSettingValue()
		if err != nil {
			return nil, err
		}
		list = append(list, &SettingExpr{
			span:  span{Pos: start, End: val.EndPos()},
			Name:  name,
			Value: val,
		})
		if alterAware {
			if p.cur().Type == TokComma && isIdentifier(p.peek(1)) && p.peek(2).Type == TokEq {
				p.advance()
				continue
			}
			break
		}
		if !p.match(TokComma) {
			break
		}
	}
	return list, nil
}

func (p *parser) parseSettingValue() (Expr, error) {
	if isLiteralStart(p.cur(), p.peek(1)) {
		return p.parseLiteral()
	}
	// ClickHouse also allows identifiers/expressions in some settings (e.g. engine settings).
	return p.parsePrimaryExpr()
}

func (p *parser) parseSettingsClause() ([]*SettingExpr, error) {
	if !p.match(KwSettings) {
		return nil, nil
	}
	return p.parseSettingExprList()
}

func (p *parser) parseStringLiteral() (*LiteralExpr, error) {
	tok := p.cur()
	if !isStringTok(tok.Type) {
		return nil, p.errAt(tok, fmt.Sprintf("expected string literal, got %q", tok.Raw))
	}
	p.advance()
	kind := "STRING"
	if tok.Type == TokHexString {
		kind = "HEX_STRING"
	} else if tok.Type == TokBinString {
		kind = "BIN_STRING"
	}
	return &LiteralExpr{
		span: span{Pos: tok.Pos, End: tok.End},
		Kind: kind,
		Raw:  tok.Raw,
		Val:  tok.Val,
	}, nil
}

func isStringTok(tt TokenType) bool {
	return tt == TokString || tt == TokHexString || tt == TokBinString
}

func isNumberStart(cur, next Token) bool {
	switch cur.Type {
	case TokInt, TokFloat, KwInf, KwNan:
		return true
	case TokPlus, TokMinus:
		return next.Type == TokInt || next.Type == TokFloat || next.Type == KwInf || next.Type == KwNan
	default:
		return false
	}
}

func isLiteralStart(cur, next Token) bool {
	if isNumberStart(cur, next) || isStringTok(cur.Type) {
		return true
	}
	return cur.Type == KwTrue || cur.Type == KwFalse || cur.Type == KwNull
}

func (p *parser) parseNumberLiteral() (*LiteralExpr, error) {
	start := p.cur().Pos
	sign := ""
	if p.cur().Type == TokPlus || p.cur().Type == TokMinus {
		sign = p.advance().Raw
	}
	tok := p.cur()
	switch tok.Type {
	case TokInt:
		p.advance()
		return &LiteralExpr{span: span{Pos: start, End: tok.End}, Kind: "INT", Raw: sign + tok.Raw, Val: sign + tok.Val}, nil
	case TokFloat:
		p.advance()
		return &LiteralExpr{span: span{Pos: start, End: tok.End}, Kind: "FLOAT", Raw: sign + tok.Raw, Val: sign + tok.Val}, nil
	case KwInf:
		p.advance()
		return &LiteralExpr{span: span{Pos: start, End: tok.End}, Kind: "INF", Raw: sign + tok.Raw, Val: sign + tok.Val}, nil
	case KwNan:
		p.advance()
		return &LiteralExpr{span: span{Pos: start, End: tok.End}, Kind: "NAN", Raw: sign + tok.Raw, Val: sign + tok.Val}, nil
	default:
		return nil, p.errAt(tok, fmt.Sprintf("expected number literal, got %q", tok.Raw))
	}
}

func (p *parser) parseLiteral() (*LiteralExpr, error) {
	tok := p.cur()
	switch tok.Type {
	case TokString, TokHexString, TokBinString:
		return p.parseStringLiteral()
	case KwTrue, KwFalse:
		p.advance()
		return &LiteralExpr{span: span{Pos: tok.Pos, End: tok.End}, Kind: "BOOL", Raw: tok.Raw, Val: strings.ToLower(tok.Val)}, nil
	case KwNull:
		p.advance()
		return &LiteralExpr{span: span{Pos: tok.Pos, End: tok.End}, Kind: "NULL", Raw: tok.Raw, Val: "NULL"}, nil
	default:
		return p.parseNumberLiteral()
	}
}
