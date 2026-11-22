package parser

import "fmt"

func (p *parser) parseInsertStmt() (*InsertStmt, error) {
	start := p.advance().Pos // INSERT
	if _, err := p.expect(KwInto); err != nil {
		return nil, err
	}
	p.match(KwTable)

	stmt := &InsertStmt{span: span{Pos: start}}
	if p.match(KwFunction) {
		fn, err := p.parseTableFunctionExpr()
		if err != nil {
			return nil, err
		}
		stmt.Func = fn
	} else {
		db, tbl, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Database = db
		stmt.Table = tbl
	}

	// Optional column list: (col1, col2, ...) with optional trailing comma
	if p.cur().Type == TokLParen && !p.looksLikeSelect() {
		p.advance()
		for {
			col, err := p.parseNestedIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Columns = append(stmt.Columns, col)
			if !p.match(TokComma) {
				break
			}
			if p.cur().Type == TokRParen {
				break
			}
		}
		if _, err := p.expect(TokRParen); err != nil {
			return nil, err
		}
	}

	// Optional SETTINGS before FORMAT/VALUES/SELECT
	if p.cur().Type == KwSettings {
		settings, err := p.parseSettingsClause()
		if err != nil {
			return nil, err
		}
		stmt.Settings = settings
	}

	// dataClause: FORMAT ident | VALUES ... | selectUnionStmt
	switch p.cur().Type {
	case KwFormat:
		p.advance()
		fmtName, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Format = fmtName
		// Stop SQL tokenization immediately after `FORMAT <name>`: any remaining
		// bytes are raw inline payload in the named format (CSV, JSONEachRow, etc.).
		stmt.Data = p.consumeRawTail(p.peek(-1).End)
		stmt.End = len(p.src)
		return stmt, nil

	case KwValues:
		valTok := p.advance()
		// ClickHouse's INSERT ... VALUES treats everything after VALUES as inline
		// Values-format data. Try parsing structured tuples when the remaining
		// input is valid SQL tuples; otherwise capture the raw payload bytes without
		// imposing SQL lexing errors on non-SQL inline data.
		savePos := p.pos
		saveLexPos := p.lex.pos
		saveLexErrObj := p.lex.err
		saveTokLen := len(p.tokens)
		saveEOF := p.eof
		saveLexErr := p.lexErr

		parsedTuples := false
		lastTupleEnd := 0
		if p.cur().Type == TokLParen {
			var rows [][]Expr
			ok := true
			for p.cur().Type == TokLParen {
				p.advance()
				var row []Expr
				if p.cur().Type != TokRParen {
					for {
						expr, err := p.parseExpr()
						if err != nil {
							ok = false
							break
						}
						row = append(row, expr)
						if !p.match(TokComma) {
							break
						}
					}
				}
				if !ok || !p.match(TokRParen) {
					ok = false
					break
				}
				lastTupleEnd = p.peek(-1).End
				rows = append(rows, row)
				if !p.match(TokComma) {
					break
				}
				if p.cur().Type != TokLParen {
					ok = false
					break
				}
			}
			if ok && (p.cur().Type == TokEOF || (p.cur().Type == TokSemicolon && p.peek(1).Type == TokEOF)) {
				stmt.Values = rows
				parsedTuples = true
			}
		}
		if parsedTuples {
			stmt.End = lastTupleEnd
		} else {
			p.pos = savePos
			p.lex.pos = saveLexPos
			p.lex.err = saveLexErrObj
			p.tokens = p.tokens[:saveTokLen]
			p.eof = saveEOF
			p.lexErr = saveLexErr
			stmt.Data = p.consumeRawTail(valTok.End)
			stmt.End = len(p.src)
		}
		return stmt, nil

	case KwSelect, KwWith, TokLParen:
		sub, err := p.parseSelectUnionStmt()
		if err != nil {
			return nil, err
		}
		stmt.Select = sub
		stmt.End = sub.End
		return stmt, nil

	default:
		return nil, p.errAt(p.cur(), fmt.Sprintf("expected FORMAT, VALUES, or SELECT in INSERT, got %q", p.cur().Raw))
	}
}

func (p *parser) parseDeleteStmt() (*DeleteStmt, error) {
	start := p.advance().Pos // DELETE
	if _, err := p.expect(KwFrom); err != nil {
		return nil, err
	}
	table, err := p.parseNestedIdentifier()
	if err != nil {
		return nil, err
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	var part *PartitionClause
	if p.cur().Type == KwIn && p.peek(1).Type == KwPartition {
		p.advance() // IN
		part, err = p.parsePartitionClause()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(KwWhere); err != nil {
		return nil, err
	}
	where, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	var settings []*SettingExpr
	if p.cur().Type == KwSettings {
		settings, err = p.parseSettingsClause()
		if err != nil {
			return nil, err
		}
	}
	return &DeleteStmt{
		span:      span{Pos: start, End: p.peek(-1).End},
		Table:     table,
		Cluster:   cluster,
		Partition: part,
		Where:     where,
		Settings:  settings,
	}, nil
}

func (p *parser) parseUpdateStmt() (*UpdateStmt, error) {
	start := p.advance().Pos // UPDATE
	table, err := p.parseNestedIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(KwSet); err != nil {
		return nil, err
	}
	assigns, err := p.parseAssignmentExprList()
	if err != nil {
		return nil, err
	}
	var part *PartitionClause
	if p.cur().Type == KwIn && p.peek(1).Type == KwPartition {
		p.advance() // IN
		if p.peek(1).Type == KwWhere {
			return nil, p.errAt(p.peek(1), "expected partition expression after IN PARTITION")
		}
		part, err = p.parsePartitionClause()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(KwWhere); err != nil {
		return nil, err
	}
	where, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	var settings []*SettingExpr
	if p.cur().Type == KwSettings {
		settings, err = p.parseSettingsClause()
		if err != nil {
			return nil, err
		}
	}
	return &UpdateStmt{
		span:        span{Pos: start, End: p.peek(-1).End},
		Table:       table,
		Assignments: assigns,
		Partition:   part,
		Where:       where,
		Settings:    settings,
	}, nil
}

func (p *parser) parseAssignmentExprList() ([]*AssignmentExpr, error) {
	var list []*AssignmentExpr
	for {
		start := p.cur().Pos
		col, err := p.parseNestedIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(TokEq); err != nil {
			return nil, err
		}
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		list = append(list, &AssignmentExpr{
			span:   span{Pos: start, End: val.EndPos()},
			Column: col,
			Value:  val,
		})
		if !p.match(TokComma) {
			break
		}
	}
	return list, nil
}
