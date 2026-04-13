package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// parseCreateTable parses a CREATE TABLE statement. The CREATE keyword has
// already been consumed by the caller. cur may be kwEXTERNAL, kwTEMPORARY,
// or kwTABLE on entry.
//
// Syntax:
//
//	CREATE [EXTERNAL] [TEMPORARY] TABLE [IF NOT EXISTS] name
//	    (column_defs [, index_defs]) | LIKE other_table
//	    [ENGINE = name]
//	    [AGGREGATE|UNIQUE|DUPLICATE KEY (cols) [CLUSTER BY (cols)]]
//	    [COMMENT 'text']
//	    [PARTITION BY ...]
//	    [DISTRIBUTED BY HASH(cols)|RANDOM [BUCKETS n|AUTO]]
//	    [ROLLUP (...)]
//	    [PROPERTIES (...)]
//	    [BROKER ext_properties]
//	    [AS query]
func (p *Parser) parseCreateTable() (ast.Node, error) {
	startLoc := p.prev.Loc // loc of CREATE token

	stmt := &ast.CreateTableStmt{}

	// Optional EXTERNAL
	if p.cur.Kind == kwEXTERNAL {
		stmt.External = true
		p.advance()
	}

	// Optional TEMPORARY
	if p.cur.Kind == kwTEMPORARY {
		stmt.Temporary = true
		p.advance()
	}

	// TABLE keyword
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	// Optional IF NOT EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// Table name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Check for LIKE (CREATE TABLE ... LIKE other_table)
	if p.cur.Kind == kwLIKE {
		p.advance() // consume LIKE
		likeName, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Like = likeName
		stmt.Loc = startLoc.Merge(ast.NodeLoc(likeName))
		return stmt, nil
	}

	// Check for CTAS without column defs: CREATE TABLE t [PROPERTIES(...)] AS SELECT ...
	// or CREATE TABLE t (col_list) ...
	if p.cur.Kind == int('(') {
		// Column definitions
		if err := p.parseColumnDefsBlock(stmt); err != nil {
			return nil, err
		}
	}

	// Parse optional trailing clauses in any order.
	if err := p.parseCreateTableClauses(stmt); err != nil {
		return nil, err
	}

	// Compute location.
	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// parseColumnDefsBlock parses the parenthesized column definitions block,
// including inline index definitions and trailing commas.
//
//	'(' column_def [, column_def]* [, index_def]* ')'
func (p *Parser) parseColumnDefsBlock(stmt *ast.CreateTableStmt) error {
	if _, err := p.expect(int('(')); err != nil {
		return err
	}

	// Parse column defs and index defs.
	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		// Check if this is an INDEX definition
		if p.cur.Kind == kwINDEX {
			idx, err := p.parseInlineIndexDef()
			if err != nil {
				return err
			}
			stmt.Indexes = append(stmt.Indexes, idx)
		} else if p.cur.Kind == kwCONSTRAINT || p.cur.Kind == kwPRIMARY || p.cur.Kind == kwUNIQUE {
			// Table-level constraint
			constraint, err := p.parseTableConstraint()
			if err != nil {
				return err
			}
			stmt.Constraints = append(stmt.Constraints, constraint)
		} else {
			col, err := p.parseColumnDef()
			if err != nil {
				return err
			}
			stmt.Columns = append(stmt.Columns, col)
		}

		// Optional comma separator
		if p.cur.Kind == int(',') {
			p.advance()
		}
	}

	if _, err := p.expect(int(')')); err != nil {
		return err
	}
	return nil
}

// parseColumnDef parses a single column definition.
//
//	col_name data_type [KEY] [agg_type]
//	    [GENERATED ALWAYS AS (expr)]
//	    [NOT NULL | NULL] [AUTO_INCREMENT]
//	    [DEFAULT value] [ON UPDATE CURRENT_TIMESTAMP]
//	    [COMMENT 'str']
func (p *Parser) parseColumnDef() (*ast.ColumnDef, error) {
	startLoc := p.cur.Loc

	// Column name
	colName, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Data type
	colType, err := p.parseDataType()
	if err != nil {
		return nil, err
	}

	col := &ast.ColumnDef{
		Name: colName,
		Type: colType,
		Loc:  ast.Loc{Start: startLoc.Start, End: colType.Loc.End},
	}

	// Parse optional column attributes in flexible order.
	// Each attribute can appear at most once, and they can appear in various orders
	// based on the grammar.
	for {
		switch p.cur.Kind {
		case kwKEY:
			col.IsKey = true
			col.Loc.End = p.cur.Loc.End
			p.advance()
			continue

		case kwMAX, kwMIN, kwSUM, kwREPLACE, kwREPLACE_IF_NOT_NULL,
			kwHLL_UNION, kwBITMAP_UNION, kwQUANTILE_UNION, kwGENERIC:
			col.AggType = strings.ToUpper(p.cur.Str)
			col.Loc.End = p.cur.Loc.End
			p.advance()
			continue

		case kwGENERATED:
			p.advance() // consume GENERATED
			// Optional ALWAYS
			p.match(kwALWAYS)
			if _, err := p.expect(kwAS); err != nil {
				return nil, err
			}
			if _, err := p.expect(int('(')); err != nil {
				return nil, err
			}
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			closeTok, err := p.expect(int(')'))
			if err != nil {
				return nil, err
			}
			col.Generated = expr
			col.Loc.End = closeTok.Loc.End
			continue

		case kwNOT:
			// NOT NULL
			p.advance()
			if _, err := p.expect(kwNULL); err != nil {
				return nil, err
			}
			b := false
			col.Nullable = &b
			col.Loc.End = p.prev.Loc.End
			continue

		case kwNULL:
			b := true
			col.Nullable = &b
			col.Loc.End = p.cur.Loc.End
			p.advance()
			continue

		case kwAUTO_INCREMENT:
			col.AutoInc = true
			col.Loc.End = p.cur.Loc.End
			p.advance()
			// Optional (initial_value)
			if p.cur.Kind == int('(') {
				p.advance()
				if p.cur.Kind == tokInt {
					p.advance()
				}
				if _, err := p.expect(int(')')); err != nil {
					return nil, err
				}
				col.Loc.End = p.prev.Loc.End
			}
			continue

		case kwDEFAULT:
			p.advance() // consume DEFAULT
			def, err := p.parseDefaultValue()
			if err != nil {
				return nil, err
			}
			col.Default = def
			col.Loc.End = ast.NodeLoc(def).End
			continue

		case kwON:
			// ON UPDATE CURRENT_TIMESTAMP
			p.advance() // consume ON
			if _, err := p.expect(kwUPDATE); err != nil {
				return nil, err
			}
			if p.cur.Kind != kwCURRENT_TIMESTAMP {
				return nil, p.syntaxErrorAtCur()
			}
			col.OnUpdate = "CURRENT_TIMESTAMP"
			col.Loc.End = p.cur.Loc.End
			p.advance()
			// Optional precision: (n)
			if p.cur.Kind == int('(') {
				p.advance()
				if p.cur.Kind == tokInt {
					p.advance()
				}
				closeTok, err := p.expect(int(')'))
				if err != nil {
					return nil, err
				}
				col.Loc.End = closeTok.Loc.End
			}
			continue

		case kwCOMMENT:
			p.advance()
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			col.Comment = p.cur.Str
			col.Loc.End = p.cur.Loc.End
			p.advance()
			continue
		}
		break
	}

	return col, nil
}

// parseDefaultValue parses the value after DEFAULT keyword.
// Supports: NULL, integers, decimals, PI, E, BITMAP_EMPTY, string literals,
// CURRENT_DATE, CURRENT_TIMESTAMP with optional precision.
func (p *Parser) parseDefaultValue() (ast.Node, error) {
	switch p.cur.Kind {
	case kwNULL:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitNull, Value: "NULL", Loc: tok.Loc}, nil

	case kwCURRENT_DATE:
		tok := p.advance()
		return &ast.FuncCallExpr{
			Name: &ast.ObjectName{Parts: []string{"CURRENT_DATE"}, Loc: tok.Loc},
			Loc:  tok.Loc,
		}, nil

	case kwCURRENT_TIMESTAMP:
		tok := p.advance()
		fc := &ast.FuncCallExpr{
			Name: &ast.ObjectName{Parts: []string{"CURRENT_TIMESTAMP"}, Loc: tok.Loc},
			Loc:  tok.Loc,
		}
		// Optional precision: (n)
		if p.cur.Kind == int('(') {
			p.advance()
			if p.cur.Kind == tokInt {
				fc.Args = append(fc.Args, &ast.Literal{
					Kind: ast.LitInt, Value: p.cur.Str, Loc: p.cur.Loc,
				})
				p.advance()
			}
			closeTok, err := p.expect(int(')'))
			if err != nil {
				return nil, err
			}
			fc.Loc.End = closeTok.Loc.End
		}
		return fc, nil

	case kwPI:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitFloat, Value: "PI", Loc: tok.Loc}, nil

	case kwE:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitFloat, Value: "E", Loc: tok.Loc}, nil

	case kwBITMAP_EMPTY:
		tok := p.advance()
		return &ast.FuncCallExpr{
			Name: &ast.ObjectName{Parts: []string{"BITMAP_EMPTY"}, Loc: tok.Loc},
			Loc:  tok.Loc,
		}, nil

	case int('-'):
		// Negative number: -INTEGER_VALUE or -DECIMAL_VALUE
		minusTok := p.advance()
		if p.cur.Kind == tokInt {
			valTok := p.advance()
			return &ast.Literal{
				Kind:  ast.LitInt,
				Value: "-" + valTok.Str,
				Loc:   ast.Loc{Start: minusTok.Loc.Start, End: valTok.Loc.End},
			}, nil
		}
		if p.cur.Kind == tokFloat {
			valTok := p.advance()
			return &ast.Literal{
				Kind:  ast.LitFloat,
				Value: "-" + valTok.Str,
				Loc:   ast.Loc{Start: minusTok.Loc.Start, End: valTok.Loc.End},
			}, nil
		}
		return nil, p.syntaxErrorAtCur()

	case tokInt:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Loc: tok.Loc}, nil

	case tokFloat:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitFloat, Value: tok.Str, Loc: tok.Loc}, nil

	case tokString:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc}, nil

	default:
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: fmt.Sprintf("expected DEFAULT value, got %q", p.cur.Str),
		}
	}
}

// parseInlineIndexDef parses an inline INDEX definition in CREATE TABLE.
//
//	INDEX [IF NOT EXISTS] name (col1, col2, ...) [USING type] [PROPERTIES(...)] [COMMENT 'str']
func (p *Parser) parseInlineIndexDef() (*ast.IndexDef, error) {
	startLoc := p.cur.Loc
	p.advance() // consume INDEX

	idx := &ast.IndexDef{
		Loc: startLoc,
	}

	// Optional IF NOT EXISTS
	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		idx.IfNotExists = true
	}

	// Index name
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	idx.Name = name

	// Column list: (col1, col2, ...)
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	cols, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	idx.Columns = cols
	idx.Loc.End = closeTok.Loc.End

	// Optional USING index_type
	if p.cur.Kind == kwUSING {
		p.advance()
		if !isIdentifierToken(p.cur.Kind) {
			return nil, p.syntaxErrorAtCur()
		}
		idx.IndexType = strings.ToLower(p.cur.Str)
		idx.Loc.End = p.cur.Loc.End
		p.advance()
	}

	// Optional PROPERTIES(...)
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		idx.Properties = props
		idx.Loc.End = p.prev.Loc.End
	}

	// Optional COMMENT
	if p.cur.Kind == kwCOMMENT {
		p.advance()
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		idx.Comment = p.cur.Str
		idx.Loc.End = p.cur.Loc.End
		p.advance()
	}

	return idx, nil
}

// parseTableConstraint parses a table-level constraint.
//
//	[CONSTRAINT name] PRIMARY KEY (col1, col2, ...)
//	[CONSTRAINT name] UNIQUE [KEY] (col1, col2, ...)
func (p *Parser) parseTableConstraint() (*ast.TableConstraint, error) {
	startLoc := p.cur.Loc
	tc := &ast.TableConstraint{Loc: startLoc}

	// Optional CONSTRAINT name
	if p.cur.Kind == kwCONSTRAINT {
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		tc.Name = name
	}

	switch p.cur.Kind {
	case kwPRIMARY:
		p.advance() // consume PRIMARY
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		tc.Type = ast.ConstraintPrimaryKey
	case kwUNIQUE:
		p.advance() // consume UNIQUE
		p.match(kwKEY) // optional KEY
		tc.Type = ast.ConstraintUnique
	default:
		return nil, p.syntaxErrorAtCur()
	}

	// Column list
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	cols, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	tc.Columns = cols
	tc.Loc.End = closeTok.Loc.End

	return tc, nil
}

// parseCreateTableClauses parses the optional trailing clauses of CREATE TABLE
// in flexible order.
func (p *Parser) parseCreateTableClauses(stmt *ast.CreateTableStmt) error {
	for {
		switch p.cur.Kind {
		case kwENGINE:
			p.advance() // consume ENGINE
			if _, err := p.expect(int('=')); err != nil {
				return err
			}
			engineName, _, err := p.parseIdentifier()
			if err != nil {
				return err
			}
			stmt.Engine = engineName
			continue

		case kwAGGREGATE, kwUNIQUE, kwDUPLICATE:
			keyDesc, err := p.parseKeyDesc()
			if err != nil {
				return err
			}
			stmt.KeyDesc = keyDesc
			// Optional CLUSTER BY (cols) -- skip for now
			if p.cur.Kind == kwCLUSTER {
				p.advance()
				if _, err := p.expect(kwBY); err != nil {
					return err
				}
				if _, err := p.expect(int('(')); err != nil {
					return err
				}
				if _, err := p.parseIdentifierList(); err != nil {
					return err
				}
				if _, err := p.expect(int(')')); err != nil {
					return err
				}
			}
			continue

		case kwCOMMENT:
			p.advance()
			if p.cur.Kind != tokString {
				return p.syntaxErrorAtCur()
			}
			stmt.Comment = p.cur.Str
			p.advance()
			continue

		case kwAUTO:
			// AUTO PARTITION BY ...
			partDesc, err := p.parsePartitionBy()
			if err != nil {
				return err
			}
			stmt.PartitionBy = partDesc
			continue

		case kwPARTITION:
			partDesc, err := p.parsePartitionBy()
			if err != nil {
				return err
			}
			stmt.PartitionBy = partDesc
			continue

		case kwDISTRIBUTED:
			distDesc, err := p.parseDistributedBy()
			if err != nil {
				return err
			}
			stmt.DistributedBy = distDesc
			continue

		case kwROLLUP:
			rollups, err := p.parseRollup()
			if err != nil {
				return err
			}
			stmt.Rollup = rollups
			continue

		case kwPROPERTIES:
			props, err := p.parseProperties()
			if err != nil {
				return err
			}
			stmt.Properties = props
			continue

		case kwBROKER:
			// BROKER ext_properties -- skip the broker properties clause
			p.advance() // consume BROKER
			if p.cur.Kind == kwPROPERTIES {
				if _, err := p.parseProperties(); err != nil {
					return err
				}
			}
			continue

		case kwAS:
			// AS query (CTAS)
			p.advance() // consume AS
			rawQuery, err := p.parseRawQuery()
			if err != nil {
				return err
			}
			stmt.AsSelect = rawQuery
			continue
		}

		// Check for bare SELECT/WITH without AS (also valid CTAS per grammar: AS? query)
		if p.cur.Kind == kwSELECT || p.cur.Kind == kwWITH {
			rawQuery, err := p.parseRawQuery()
			if err != nil {
				return err
			}
			stmt.AsSelect = rawQuery
			continue
		}

		break
	}
	return nil
}

// parseKeyDesc parses: AGGREGATE|UNIQUE|DUPLICATE KEY (col1, col2, ...)
func (p *Parser) parseKeyDesc() (*ast.KeyDesc, error) {
	startLoc := p.cur.Loc
	keyType := strings.ToUpper(p.cur.Str)
	p.advance() // consume AGGREGATE/UNIQUE/DUPLICATE

	if _, err := p.expect(kwKEY); err != nil {
		return nil, err
	}

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	cols, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}

	return &ast.KeyDesc{
		Type:    keyType,
		Columns: cols,
		Loc:     ast.Loc{Start: startLoc.Start, End: closeTok.Loc.End},
	}, nil
}

// parsePartitionBy parses:
//
//	[AUTO] PARTITION BY [RANGE|LIST] (col_or_func, ...) (partition_defs...)
func (p *Parser) parsePartitionBy() (*ast.PartitionDesc, error) {
	startLoc := p.cur.Loc
	pd := &ast.PartitionDesc{
		Loc: startLoc,
	}

	// Optional AUTO
	if p.cur.Kind == kwAUTO {
		pd.Auto = true
		p.advance()
	}

	// PARTITION BY
	if _, err := p.expect(kwPARTITION); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}

	// Optional RANGE or LIST
	switch p.cur.Kind {
	case kwRANGE:
		pd.Type = "RANGE"
		p.advance()
	case kwLIST:
		pd.Type = "LIST"
		p.advance()
	default:
		// Type may be inferred from context (AUTO partition can omit it)
		pd.Type = "RANGE" // default for AUTO
	}

	// Parse partition column/function list: (col1, func(col2), ...)
	if p.cur.Kind == int('(') {
		cols, funcs, err := p.parseIdentityOrFunctionList()
		if err != nil {
			return nil, err
		}
		pd.Columns = cols
		pd.FuncExprs = funcs
	}

	// Parse partition definitions: (...)
	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
			item, err := p.parsePartitionItem()
			if err != nil {
				return nil, err
			}
			pd.Partitions = append(pd.Partitions, item)

			if p.cur.Kind == int(',') {
				p.advance()
			}
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		pd.Loc.End = closeTok.Loc.End
	}

	return pd, nil
}

// parseIdentityOrFunctionList parses a parenthesized list of identifiers or
// function calls used in PARTITION BY clauses.
// Returns (columns, funcExprs) where funcExprs contains raw function text.
func (p *Parser) parseIdentityOrFunctionList() ([]string, []string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, nil, err
	}

	var cols []string
	var funcs []string

	for {
		// Check if it's a function call: identifier followed by '('
		if isIdentifierToken(p.cur.Kind) {
			next := p.peekNext()
			if next.Kind == int('(') {
				// Function call -- capture raw text
				funcStart := p.cur.Loc.Start
				p.advance() // consume function name
				p.advance() // consume '('
				depth := 1
				for depth > 0 && p.cur.Kind != tokEOF {
					if p.cur.Kind == int('(') {
						depth++
					} else if p.cur.Kind == int(')') {
						depth--
						if depth == 0 {
							break
						}
					}
					p.advance()
				}
				funcEnd := p.cur.Loc.End
				p.advance() // consume closing ')'
				funcText := p.input[funcStart:funcEnd]
				funcs = append(funcs, strings.TrimSpace(funcText))
			} else {
				// Simple identifier
				ident, _, err := p.parseIdentifier()
				if err != nil {
					return nil, nil, err
				}
				cols = append(cols, ident)
			}
		} else {
			return nil, nil, p.syntaxErrorAtCur()
		}

		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, nil, err
	}

	return cols, funcs, nil
}

// parsePartitionItem parses a single partition definition.
func (p *Parser) parsePartitionItem() (*ast.PartitionItem, error) {
	startLoc := p.cur.Loc

	// Step partition: FROM (...) TO (...) INTERVAL n [unit]
	if p.cur.Kind == kwFROM {
		return p.parseStepPartition(startLoc)
	}

	// Named partition: PARTITION [IF NOT EXISTS] name ...
	if p.cur.Kind != kwPARTITION {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume PARTITION

	// Optional IF NOT EXISTS
	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
	}

	// Partition name
	partName, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	item := &ast.PartitionItem{
		Name: partName,
		Loc:  startLoc,
	}

	// VALUES clause
	if p.cur.Kind == kwVALUES {
		p.advance() // consume VALUES

		switch p.cur.Kind {
		case kwLESS:
			// LESS THAN (values) | MAXVALUE
			p.advance() // consume LESS
			if _, err := p.expect(kwTHAN); err != nil {
				return nil, err
			}

			if p.cur.Kind == kwMAXVALUE {
				item.IsMaxValue = true
				item.Loc.End = p.cur.Loc.End
				p.advance()
			} else {
				vals, endLoc, err := p.parsePartitionValueList()
				if err != nil {
					return nil, err
				}
				item.Values = vals
				item.Loc.End = endLoc
			}

		case kwIN:
			// IN ((val1, val2), (val3, val4)) or IN (val1, val2)
			p.advance() // consume IN
			inVals, endLoc, err := p.parseInPartitionValues()
			if err != nil {
				return nil, err
			}
			item.InValues = inVals
			item.Loc.End = endLoc

		case int('['):
			// Fixed range: [lower, upper)
			vals, endLoc, err := p.parseFixedRangePartition()
			if err != nil {
				return nil, err
			}
			item.Values = vals
			item.Loc.End = endLoc

		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	// Optional partition properties
	if p.cur.Kind == int('(') {
		// Skip inline partition properties
		p.advance()
		for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
			p.advance()
		}
		if p.cur.Kind == int(')') {
			item.Loc.End = p.cur.Loc.End
			p.advance()
		}
	}

	return item, nil
}

// parseStepPartition parses: FROM (values) TO (values) INTERVAL n [unit]
func (p *Parser) parseStepPartition(startLoc ast.Loc) (*ast.PartitionItem, error) {
	p.advance() // consume FROM

	fromVals, _, err := p.parsePartitionValueList()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}

	toVals, _, err := p.parsePartitionValueList()
	if err != nil {
		return nil, err
	}

	item := &ast.PartitionItem{
		IsStep:     true,
		FromValues: fromVals,
		ToValues:   toVals,
		Loc:        startLoc,
	}

	if _, err := p.expect(kwINTERVAL); err != nil {
		return nil, err
	}

	// Interval amount
	if p.cur.Kind != tokInt {
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: fmt.Sprintf("expected integer interval amount, got %q", p.cur.Str),
		}
	}
	item.Interval = p.cur.Str
	item.Loc.End = p.cur.Loc.End
	p.advance()

	// Optional interval unit
	if p.cur.Kind == kwDAY || p.cur.Kind == kwWEEK || p.cur.Kind == kwMONTH ||
		p.cur.Kind == kwYEAR || p.cur.Kind == kwHOUR || p.cur.Kind == kwMINUTE ||
		p.cur.Kind == kwSECOND || p.cur.Kind == tokIdent {
		item.IntervalUnit = strings.ToUpper(p.cur.Str)
		item.Loc.End = p.cur.Loc.End
		p.advance()
	}

	return item, nil
}

// parsePartitionValueList parses: '(' value [, value]* ')'
// Returns the values as strings and the end byte offset.
func (p *Parser) parsePartitionValueList() ([]string, int, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, 0, err
	}

	var vals []string
	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		val, err := p.parsePartitionValueDef()
		if err != nil {
			return nil, 0, err
		}
		vals = append(vals, val)

		if p.cur.Kind == int(',') {
			p.advance()
		}
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, 0, err
	}

	return vals, closeTok.Loc.End, nil
}

// parsePartitionValueDef parses a single partition value:
// [-] INTEGER | STRING | MAXVALUE | NULL
func (p *Parser) parsePartitionValueDef() (string, error) {
	switch p.cur.Kind {
	case kwMAXVALUE:
		tok := p.advance()
		return strings.ToUpper(tok.Str), nil
	case kwNULL:
		tok := p.advance()
		return strings.ToUpper(tok.Str), nil
	case tokString:
		tok := p.advance()
		return tok.Str, nil
	case tokInt:
		tok := p.advance()
		return tok.Str, nil
	case int('-'):
		p.advance() // consume '-'
		if p.cur.Kind == tokInt {
			tok := p.advance()
			return "-" + tok.Str, nil
		}
		if p.cur.Kind == tokFloat {
			tok := p.advance()
			return "-" + tok.Str, nil
		}
		return "", p.syntaxErrorAtCur()
	default:
		return "", &ParseError{
			Loc: p.cur.Loc,
			Msg: fmt.Sprintf("expected partition value, got %q", p.cur.Str),
		}
	}
}

// parseInPartitionValues parses IN partition values:
// '(' (value_list) [, (value_list)]* ')' or '(' value [, value]* ')'
func (p *Parser) parseInPartitionValues() ([][]string, int, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, 0, err
	}

	var result [][]string

	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		if p.cur.Kind == int('(') {
			// Parenthesized value list
			vals, _, err := p.parsePartitionValueList()
			if err != nil {
				return nil, 0, err
			}
			result = append(result, vals)
		} else {
			// Single value
			val, err := p.parsePartitionValueDef()
			if err != nil {
				return nil, 0, err
			}
			result = append(result, []string{val})
		}

		if p.cur.Kind == int(',') {
			p.advance()
		}
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, 0, err
	}

	return result, closeTok.Loc.End, nil
}

// parseFixedRangePartition parses: '[' lower_values ',' upper_values ')'
// This is for fixed range partition syntax like: VALUES [("2020-01-01"), ("2020-02-01"))
func (p *Parser) parseFixedRangePartition() ([]string, int, error) {
	// Consume '['
	p.advance()

	// Parse lower bound
	lowerVals, _, err := p.parsePartitionValueList()
	if err != nil {
		return nil, 0, err
	}

	if _, err := p.expect(int(',')); err != nil {
		return nil, 0, err
	}

	// Parse upper bound
	upperVals, _, err := p.parsePartitionValueList()
	if err != nil {
		return nil, 0, err
	}

	// Expect ')' (half-open interval)
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, 0, err
	}

	// Combine: store as lower1, lower2, ..., upper1, upper2, ...
	all := append(lowerVals, upperVals...)
	return all, closeTok.Loc.End, nil
}

// parseDistributedBy parses:
//
//	DISTRIBUTED BY HASH(col1, col2, ...) [BUCKETS n | BUCKETS AUTO]
//	DISTRIBUTED BY RANDOM [BUCKETS n | BUCKETS AUTO]
func (p *Parser) parseDistributedBy() (*ast.DistributionDesc, error) {
	startLoc := p.cur.Loc
	p.advance() // consume DISTRIBUTED

	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}

	dd := &ast.DistributionDesc{
		Loc: startLoc,
	}

	switch p.cur.Kind {
	case kwHASH:
		dd.Type = "HASH"
		p.advance() // consume HASH

		// Column list: (col1, col2, ...)
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		cols, err := p.parseIdentifierList()
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		dd.Columns = cols
		dd.Loc.End = closeTok.Loc.End

	case kwRANDOM:
		dd.Type = "RANDOM"
		dd.Loc.End = p.cur.Loc.End
		p.advance()

	default:
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: fmt.Sprintf("expected HASH or RANDOM after DISTRIBUTED BY, got %q", p.cur.Str),
		}
	}

	// Optional BUCKETS n | BUCKETS AUTO
	if p.cur.Kind == kwBUCKETS {
		p.advance() // consume BUCKETS
		if p.cur.Kind == kwAUTO {
			dd.Auto = true
			dd.Loc.End = p.cur.Loc.End
			p.advance()
		} else if p.cur.Kind == tokInt {
			dd.Buckets = int(p.cur.Ival)
			dd.Loc.End = p.cur.Loc.End
			p.advance()
		} else {
			return nil, &ParseError{
				Loc: p.cur.Loc,
				Msg: fmt.Sprintf("expected integer or AUTO after BUCKETS, got %q", p.cur.Str),
			}
		}
	}

	return dd, nil
}

// parseRollup parses:
//
//	ROLLUP '(' rollup_def [, rollup_def]* ')'
//
// where rollup_def is: name (col1, col2, ...) [DUPLICATE KEY (key_cols)] [PROPERTIES(...)]
func (p *Parser) parseRollup() ([]*ast.RollupDef, error) {
	p.advance() // consume ROLLUP

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	var rollups []*ast.RollupDef
	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		rd, err := p.parseRollupDef()
		if err != nil {
			return nil, err
		}
		rollups = append(rollups, rd)

		if p.cur.Kind == int(',') {
			p.advance()
		}
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	return rollups, nil
}

// parseRollupDef parses a single rollup definition:
//
//	name (col1, col2, ...) [DUPLICATE KEY (key_cols)] [PROPERTIES(...)]
func (p *Parser) parseRollupDef() (*ast.RollupDef, error) {
	startLoc := p.cur.Loc

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Column list
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	cols, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}

	rd := &ast.RollupDef{
		Name:    name,
		Columns: cols,
		Loc:     ast.Loc{Start: startLoc.Start, End: closeTok.Loc.End},
	}

	// Optional DUPLICATE KEY (key_cols)
	if p.cur.Kind == kwDUPLICATE {
		p.advance() // consume DUPLICATE
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		keys, err := p.parseIdentifierList()
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		rd.DupKeys = keys
		rd.Loc.End = closeTok.Loc.End
	}

	// Optional PROPERTIES
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		rd.Properties = props
		rd.Loc.End = p.prev.Loc.End
	}

	return rd, nil
}

// parseRawQuery consumes remaining tokens as a raw SQL query (for CTAS).
// The AS keyword has already been consumed if present.
func (p *Parser) parseRawQuery() (*ast.RawQuery, error) {
	startLoc := p.cur.Loc
	start := p.cur.Loc.Start

	// Consume all remaining tokens until EOF or semicolon at depth 0.
	depth := 0
	for p.cur.Kind != tokEOF {
		switch p.cur.Kind {
		case int('('):
			depth++
		case int(')'):
			if depth > 0 {
				depth--
			}
		case int(';'):
			if depth == 0 {
				goto done
			}
		}
		p.advance()
	}
done:
	end := p.prev.Loc.End
	rawText := strings.TrimSpace(p.input[start:end])

	return &ast.RawQuery{
		RawText: rawText,
		Loc:     ast.Loc{Start: startLoc.Start, End: end},
	}, nil
}
