package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// parseDynColFunc reads one of the dynamic column functions, which each have
// their own production in function_call_nonkeyword (MariaDB 11.8
// sql/sql_yacc.yy):
//
//	COLUMN_ADD_SYM '(' expr ',' dyncall_create_list ')'
//	COLUMN_DELETE_SYM '(' expr ',' expr_list ')'
//	COLUMN_CREATE_SYM '(' dyncall_create_list ')'
//	COLUMN_GET_SYM '(' expr ',' expr AS cast_type ')'
//
//	dyncall_create_list: dyncall_create_element | dyncall_create_list ',' dyncall_create_element
//	dyncall_create_element: expr ',' expr opt_dyncol_type
//	opt_dyncol_type: %empty | AS dyncol_type
//
// Each name is a keyword, not a token the lexer folds out of `name(`, so a space
// before '(' changes nothing; a quoted or qualified name is an ordinary call.
// mariadb:11.8.8 answers 1064 for COLUMN_GET(a, a), COLUMN_CREATE(a),
// COLUMN_ADD(a, a) and COLUMN_DELETE(a), and takes COLUMN_CREATE(a, a AS INT).
//
// p.cur is just past the '('. fc.ArgTypes runs beside fc.Args and holds each
// AS type.
func (p *Parser) parseDynColFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	arg := func() error {
		e, err := p.parseExpr()
		if err != nil {
			return err
		}
		if e == nil || isBareInterval(e) {
			return p.syntaxErrorAtCur()
		}
		fc.Args = append(fc.Args, e)
		fc.ArgTypes = append(fc.ArgTypes, nil)
		return nil
	}
	comma := func() error {
		_, err := p.expect(',')
		return err
	}

	if fc.Name != "COLUMN_CREATE" {
		// The blob every function but COLUMN_CREATE reads first.
		if err := arg(); err != nil {
			return nil, err
		}
		if err := comma(); err != nil {
			return nil, err
		}
	}
	switch fc.Name {
	case "COLUMN_DELETE":
		for {
			if err := arg(); err != nil {
				return nil, err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	case "COLUMN_GET":
		if err := arg(); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwAS); err != nil {
			return nil, err
		}
		dt, err := p.parseCastDataType()
		if err != nil {
			return nil, err
		}
		fc.ArgTypes[len(fc.ArgTypes)-1] = dt
	default:
		// COLUMN_CREATE and COLUMN_ADD: name, value [AS dyncol_type], ...
		for {
			if err := arg(); err != nil {
				return nil, err
			}
			if err := comma(); err != nil {
				return nil, err
			}
			if err := arg(); err != nil {
				return nil, err
			}
			if _, ok := p.match(kwAS); ok {
				dt, err := p.parseDynColType()
				if err != nil {
					return nil, err
				}
				fc.ArgTypes[len(fc.ArgTypes)-1] = dt
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	// No production of the four has a window.
	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseDynColType reads dyncol_type, a shorter list than CAST's:
//
//	numeric_dyncol_type: INT_SYM | UNSIGNED INT_SYM | DOUBLE_SYM | REAL | FLOAT_SYM
//	                   | DECIMAL_SYM float_options
//	temporal_dyncol_type: DATE_SYM | TIME_SYM opt_field_scale | DATETIME opt_field_scale
//	string_dyncol_type: char opt_binary | nchar
//
// The lexer reads INTEGER and INT4 as INT, FLOAT4 as FLOAT, FLOAT8 as DOUBLE and
// DEC as DECIMAL (sql/lex.h). mariadb:11.8.8 answers 1064 for SIGNED, UNSIGNED
// alone, NUMERIC, DOUBLE(3,2), DATE(1), CHAR(3), NCHAR(2), VARCHAR(2), BINARY
// and TIMESTAMP here, and checks no size: DECIMAL(100) and TIME(7) parse.
func (p *Parser) parseDynColType() (*nodes.DataType, error) {
	dt := &nodes.DataType{Loc: nodes.Loc{Start: p.pos()}}
	switch p.cur.Type {
	case kwINT, kwINTEGER, kwINT4:
		dt.Name = "INT"
		p.advance()
	case kwUNSIGNED:
		p.advance()
		if _, ok := p.match(kwINT, kwINTEGER, kwINT4); !ok {
			return nil, p.syntaxErrorAtCur()
		}
		dt.Name, dt.Unsigned = "INT", true
	case kwDOUBLE, kwFLOAT8:
		dt.Name = "DOUBLE"
		p.advance()
	case kwREAL:
		dt.Name = "REAL"
		p.advance()
	case kwFLOAT, kwFLOAT4:
		dt.Name = "FLOAT"
		p.advance()
	case kwDECIMAL, kwDEC:
		dt.Name = "DECIMAL"
		p.advance()
		if err := p.parseOptionalPrecision(dt); err != nil {
			return nil, err
		}
	case kwDATE:
		dt.Name = "DATE"
		p.advance()
	case kwTIME, kwDATETIME:
		dt.Name = "TIME"
		if p.cur.Type == kwDATETIME {
			dt.Name = "DATETIME"
		}
		p.advance()
		if err := p.parseCastSize(dt, true); err != nil {
			return nil, err
		}
	case kwCHAR, kwCHARACTER:
		dt.Name = "CHAR"
		p.advance()
		if err := p.parseCastCharOptions(dt); err != nil {
			return nil, err
		}
	case kwNCHAR:
		dt.Name = "NCHAR"
		p.advance()
	case kwNATIONAL:
		p.advance()
		if _, ok := p.match(kwCHAR, kwCHARACTER); !ok {
			return nil, p.syntaxErrorAtCur()
		}
		dt.Name = "NCHAR"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	dt.Loc.End = p.prev.End
	return dt, nil
}

// parsePercentileFunc reads an inverse distribution function after its '(':
//
//	inverse_distribution_function:
//	    percentile_function OVER_SYM '(' opt_window_partition_clause ')'
//	percentile_function:
//	    inverse_distribution_function_def WITHIN GROUP_SYM '(' order_by_single_element_list ')'
//	  | MEDIAN_SYM '(' expr ')'
//	inverse_distribution_function_def:
//	    PERCENTILE_CONT_SYM '(' expr ')' | PERCENTILE_DISC_SYM '(' expr ')'
//	order_by_single_element_list: ORDER_SYM BY order_ident order_dir
//
// Each is a window function with a window of its own: a PARTITION BY list or
// nothing. mariadb:11.8.8 answers 1064 for PERCENTILE_CONT(a) without WITHIN
// GROUP or without OVER, for two ORDER BY items, and for OVER w, OVER (w) and
// OVER (ORDER BY b), and for MEDIAN as for the others.
func (p *Parser) parsePercentileFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	arg, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if arg == nil || isBareInterval(arg) {
		return nil, p.syntaxErrorAtCur()
	}
	fc.Args = []nodes.ExprNode{arg}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	if fc.Name != "MEDIAN" {
		// WITHIN is no keyword of this lexer, so the word is matched as written.
		if p.cur.Type != tokIDENT || !strings.EqualFold(p.cur.Str, "WITHIN") || p.backtickQuotedAt(p.cur.Loc) {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		for _, t := range []int{kwGROUP, '(', kwORDER, kwBY} {
			if _, err := p.expect(t); err != nil {
				return nil, err
			}
		}
		start := p.pos()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		item := &nodes.OrderByItem{Loc: nodes.Loc{Start: start}, Expr: expr}
		if _, ok := p.match(kwDESC); ok {
			item.Direction, item.Desc = nodes.OrderDirectionDesc, true
		} else if _, ok := p.match(kwASC); ok {
			item.Direction = nodes.OrderDirectionAsc
		}
		item.Loc.End = p.prev.End
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		fc.OrderBy, fc.WithinGroup = []*nodes.OrderByItem{item}, true
	}
	over := &nodes.WindowDef{Loc: nodes.Loc{Start: p.pos()}}
	if _, err := p.expect(kwOVER); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, ok := p.match(kwPARTITION); ok {
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		exprs, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		over.PartitionBy = exprs
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	over.Loc.End = p.prev.End
	fc.Over = over
	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseValueFunc reads VALUE's column after its '(':
//
//	VALUE_SYM '(' simple_ident_nospvar ')'
//	simple_ident_nospvar: ident | ident '.' ident | '.' ident '.' ident
//	                    | ident '.' ident '.' ident
//
// A leading dot names a table of no database. mariadb:11.8.8 answers 1064 for
// VALUE(), VALUE(1), VALUE(@a) and VALUE(a, b).
func (p *Parser) parseValueFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	var col *nodes.ColumnRef
	if p.cur.Type == '.' {
		start := p.pos()
		p.advance()
		table, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('.'); err != nil {
			return nil, err
		}
		column, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		col = &nodes.ColumnRef{Loc: nodes.Loc{Start: start, End: p.prev.End}, Table: table, Column: column}
	} else {
		if p.isVariableRef() {
			return nil, p.syntaxErrorAtCur()
		}
		var err error
		if col, err = p.parseColumnRef(); err != nil {
			return nil, err
		}
		if col.Star {
			return nil, p.syntaxErrorAtTok(p.prev)
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Args = []nodes.ExprNode{col}
	fc.Loc.End = p.prev.End
	return fc, nil
}

// parsePackageFuncCall reads a call of a stored package's routine, from the
// '(' after its three-part name:
//
//	function_call_generic: ident_cli '.' ident_cli '.' ident_cli '(' opt_expr_list ')'
//
// The arguments are ordinary ones and the call has no window. mariadb:11.8.8
// parses DO a.a.a() and answers 1305 for the routine it does not have.
func (p *Parser) parsePackageFuncCall(start int, db, pkg, routine string) (nodes.ExprNode, error) {
	if err := p.refuseStoredExprFunction(strings.ToUpper(routine), true, start); err != nil {
		return nil, err
	}
	if err := refuseDatabaseName(db, start); err != nil {
		return nil, err
	}
	if err := p.refusePartitionFunction(strings.ToUpper(routine), true, start); err != nil {
		return nil, err
	}
	p.advance() // consume '('
	fc := &nodes.FuncCallExpr{
		Loc:       nodes.Loc{Start: start},
		Schema:    db,
		Package:   pkg,
		Name:      strings.ToUpper(routine),
		RawName:   routine,
		HasParens: true,
	}
	if p.cur.Type != ')' {
		for {
			arg, err := p.parseOrdinaryCallArg()
			if err != nil {
				return nil, err
			}
			fc.Args = append(fc.Args, arg)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prev.End
	return fc, nil
}
