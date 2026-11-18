package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseMethodOptStmt parses method_opt, the statistics method DBMS_STATS takes
// as an argument and OceanBase's grammar also reaches as a statement of its own
// (stmt: method_opt). The server parses such a statement and runs nothing for
// it, so the tree is kept and no visitor reads it.
//
//	method_opt:  method_list
//	method_list: method | method_list method
//	method:      for_all | for_columns
//
// Ref: method_opt and the rules under it in sql_parser_mysql_mode.y
// (oceanbase v4.3.5_CE_BP6, lines 23802-23941).
func (p *Parser) parseMethodOptStmt() (*nodes.MethodOptStmt, error) {
	stmt := &nodes.MethodOptStmt{Loc: nodes.Loc{Start: p.pos()}}
	for p.cur.Type == kwFOR {
		method, err := p.parseStatsMethod()
		if err != nil {
			return nil, err
		}
		stmt.Methods = append(stmt.Methods, method)
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseStatsMethod parses one method of method_opt, at its FOR.
//
//	for_all:     FOR ALL [INDEXED | HIDDEN] COLUMNS [size_clause]
//	for_columns: FOR COLUMNS [for_columns_list]
func (p *Parser) parseStatsMethod() (*nodes.StatsMethod, error) {
	method := &nodes.StatsMethod{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // FOR

	if p.cur.Type == kwALL {
		method.ForAll = true
		p.advance()
		if p.isBareWord("INDEXED") || p.isBareWord("HIDDEN") {
			method.Indexed = strings.ToUpper(p.advance().Str)
		}
		if !p.isBareWord("COLUMNS") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if p.isBareWord("SIZE") {
			size, err := p.parseSizeClause()
			if err != nil {
				return nil, err
			}
			method.Size = size
		}
		method.Loc.End = p.prev.End
		return method, nil
	}

	if !p.isBareWord("COLUMNS") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// for_columns_list: one item, then an item beside it or after a comma. An
	// item that ends the list is one the next token does not continue, and what
	// stands there is the next method or the end of the statement. The rule also
	// has a FOR COLUMNS with no item at all (for_columns' second alternative).
	for p.atStatsForColumnsItem() {
		item, err := p.parseStatsForColumnsItem()
		if err != nil {
			return nil, err
		}
		method.Items = append(method.Items, item)
		if p.cur.Type == ',' {
			p.advance()
			if !p.atStatsForColumnsItem() {
				return nil, p.syntaxErrorAtCur()
			}
		}
	}
	method.Loc.End = p.prev.End
	return method, nil
}

// atStatsForColumnsItem reports whether the current token opens a
// for_columns_item: a size_clause, or a column_clause, which is a column_name or
// an extension's parenthesized name list.
func (p *Parser) atStatsForColumnsItem() bool {
	if p.isBareWord("SIZE") || p.cur.Type == '(' {
		return true
	}
	if p.cur.Type == tokIDENT {
		return !p.curIsVariableOrIntroducer()
	}
	return p.cur.Type >= 700 && !isReserved(p.cur.Type)
}

// parseStatsForColumnsItem parses one for_columns_item:
//
//	column_clause [size_clause] | size_clause
//
// where column_clause is a column_name or '(' column_name_list ')'.
func (p *Parser) parseStatsForColumnsItem() (*nodes.StatsForColumnsItem, error) {
	item := &nodes.StatsForColumnsItem{Loc: nodes.Loc{Start: p.pos()}}
	switch {
	case p.isBareWord("SIZE") && p.sizeValueAhead():
		// A size_clause alone, which is the item SIZE 10 and not a column
		// called size: SIZE's precedence is above LOWER_PARENS. SIZE with no
		// value after it is the column: `FOR COLUMNS SIZE` parses.
	case p.cur.Type == '(':
		names, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		item.Columns = names
	default:
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		item.Columns = []string{name}
	}
	if p.isBareWord("SIZE") {
		size, err := p.parseSizeClause()
		if err != nil {
			return nil, err
		}
		item.Size = size
	}
	item.Loc.End = p.prev.End
	return item, nil
}

// parseSizeClause parses size_clause, at its SIZE:
//
//	size_clause: SIZE {AUTO | REPEAT | SKEWONLY | number_literal}
//
// The number is kept as written.
func (p *Parser) parseSizeClause() (string, error) {
	if !p.isBareWord("SIZE") {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	switch {
	case p.isBareWord("AUTO"), p.isBareWord("REPEAT"), p.isBareWord("SKEWONLY"):
		return strings.ToUpper(p.advance().Str), nil
	case p.cur.Type == tokICONST, p.cur.Type == tokFCONST:
		return p.advance().Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// sizeValueAhead reports whether the token after SIZE makes it a size_clause:
// AUTO, REPEAT, SKEWONLY or a number. Otherwise SIZE names a column, as in
// `FOR COLUMNS SIZE`, which the pinned server accepts
// (number_literal/INTNUM!drop-1).
func (p *Parser) sizeValueAhead() bool {
	next := p.peekNext()
	if next.Type == tokICONST || next.Type == tokFCONST {
		return true
	}
	if next.Type != tokIDENT && next.Type < 700 {
		return false
	}
	if p.backtickQuotedAt(next.Loc) {
		return false
	}
	return strings.EqualFold(next.Str, "AUTO") || strings.EqualFold(next.Str, "REPEAT") ||
		strings.EqualFold(next.Str, "SKEWONLY")
}
