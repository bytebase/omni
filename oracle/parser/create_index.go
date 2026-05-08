package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateIndexStmt parses a CREATE INDEX statement.
// The CREATE keyword has already been consumed. The current token is
// UNIQUE, BITMAP, MULTIVALUE, or INDEX.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-INDEX.html
//
//	CREATE [ UNIQUE ] [ BITMAP ] [ MULTIVALUE ] INDEX [ IF NOT EXISTS ]
//	    [ schema. ] index_name
//	    [ index_ilm_clause ]
//	    { cluster_index_clause | table_index_clause | bitmap_join_index_clause }
//	    [ { DEFERRED | IMMEDIATE } INVALIDATION ]
func (p *Parser) parseCreateIndexStmt(start int) (*nodes.CreateIndexStmt, error) {
	stmt := &nodes.CreateIndexStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional UNIQUE
	if p.cur.Type == kwUNIQUE {
		stmt.Unique = true
		p.advance()
	}

	// Optional BITMAP
	if p.cur.Type == kwBITMAP {
		stmt.Bitmap = true
		p.advance()
	}

	// Optional MULTIVALUE
	if p.isIdentLikeStr("MULTIVALUE") {
		stmt.Multivalue = true
		p.advance()
	}

	// INDEX keyword
	if p.cur.Type == kwINDEX {
		p.advance()
	}

	// Optional IF NOT EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			stmt.IfNotExists = true
			p.advance() // consume IF
			p.advance() // consume NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // consume EXISTS
			}
		}
	}
	var parseErr410 error

	// Index name
	stmt.Name, parseErr410 = p.parseReservedCheckedObjectName()
	if parseErr410 !=

		// ON clause
		nil {
		return nil, parseErr410
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	if p.cur.Type != kwON {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume ON

	// cluster_index_clause: ON CLUSTER [schema.]cluster
	if p.cur.Type == kwCLUSTER {
		p.advance()
		var // consume CLUSTER
		parseErr411 error
		stmt.Cluster, parseErr411 = p.parseObjectName()
		if parseErr411 !=
			// Parse index_attributes
			nil {
			return nil, parseErr411
		}
		if stmt.Cluster == nil || stmt.Cluster.Name == "" {
			return nil, p.syntaxErrorAtCur()
		}
		parseErr412 := p.parseCreateIndexAttributes(stmt)
		if parseErr412 !=
			// Parse optional INVALIDATION
			nil {
			return nil, parseErr412
		}
		parseErr413 := p.parseCreateIndexInvalidation(stmt)
		if parseErr413 != nil {
			return nil, parseErr413
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	var parseErr414 error

	// table_index_clause or bitmap_join_index_clause
	// Both start with [schema.]table
	stmt.Table, parseErr414 = p.parseObjectName()
	if parseErr414 !=

		// Optional table alias (single identifier before '(')
		nil {
		return nil, parseErr414
	}
	if stmt.Table == nil || stmt.Table.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	if p.isIdentLike() && p.cur.Type != '(' && !p.isKeywordLike() {
		var parseErr415 error
		stmt.Alias, parseErr415 = p.parseIdentifier()
		if parseErr415 !=

			// ( column_list )
			nil {
			return nil, parseErr415
		}
	}

	if p.cur.Type == '(' {
		p.advance()
		var parseErr416 error
		stmt.Columns, parseErr416 = p.parseIndexColumnList()
		if parseErr416 != nil {
			return nil, parseErr416
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	// Check for bitmap_join_index_clause: FROM ... WHERE ...
	if p.cur.Type == kwFROM {
		p.advance() // consume FROM
		stmt.FromTables = &nodes.List{}
		// Parse table list
		tbl, parseErr417 := p.parseObjectName()
		if parseErr417 != nil {
			return nil, parseErr417
		}
		stmt.FromTables.Items = append(stmt.FromTables.Items, tbl)
		for p.cur.Type == ',' {
			p.advance()
			var // consume ','
			parseErr418 error
			tbl, parseErr418 = p.parseObjectName()
			if parseErr418 != nil {
				return nil, parseErr418
			}
			stmt.FromTables.Items = append(stmt.FromTables.Items, tbl)
		}
		// WHERE join_condition
		if p.cur.Type == kwWHERE {
			p.advance()
			var // consume WHERE
			parseErr419 error
			stmt.Where, parseErr419 = p.parseExpr()
			if parseErr419 !=

				// Parse trailing options (index_properties + index_attributes)
				nil {
				return nil, parseErr419
			}
		}
	}
	parseErr420 := p.parseCreateIndexAttributes(stmt)
	if parseErr420 !=

		// Parse optional INVALIDATION
		nil {
		return nil, parseErr420
	}
	parseErr421 := p.parseCreateIndexInvalidation(stmt)
	if parseErr421 != nil {
		return nil, parseErr421
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateIndexInvalidation parses optional { DEFERRED | IMMEDIATE } INVALIDATION.
func (p *Parser) parseCreateIndexInvalidation(stmt *nodes.CreateIndexStmt) error {
	if p.cur.Type == kwDEFERRED {
		stmt.Invalidation = "DEFERRED"
		p.advance()
		if p.isIdentLikeStr("INVALIDATION") {
			p.advance()
		}
	} else if p.cur.Type == kwIMMEDIATE {
		stmt.Invalidation = "IMMEDIATE"
		p.advance()
		if p.isIdentLikeStr("INVALIDATION") {
			p.advance()
		}
	}
	return nil
}

// isKeywordLike checks if current token looks like a known clause-starting keyword
// used to distinguish table aliases from trailing clauses.
func (p *Parser) isKeywordLike() bool {
	switch p.cur.Type {
	case kwREVERSE, kwTABLESPACE, kwLOCAL, kwGLOBAL, kwONLINE,
		kwPARALLEL, kwNOPARALLEL, kwCOMPRESS, kwNOCOMPRESS,
		kwLOGGING, kwNOLOGGING, kwINVISIBLE, kwPCTFREE,
		kwFROM, kwWHERE, kwPARTITION, kwDEFERRED, kwIMMEDIATE:
		return true
	}
	if p.isIdentLikeStr("SORT") || p.isIdentLikeStr("NOSORT") ||
		p.isIdentLikeStr("VISIBLE") || p.isIdentLikeStr("INDEXTYPE") ||
		p.isIdentLikeStr("INDEXING") || p.isIdentLikeStr("PARAMETERS") {
		return true
	}
	return false
}

// parseCreateIndexAttributes parses index_attributes and index_properties for CREATE INDEX.
//
//	index_attributes:
//	    [ physical_attributes_clause ]
//	    [ TABLESPACE { tablespace | DEFAULT } ]
//	    [ index_compression ]
//	    [ partial_index_clause ]
//	    [ parallel_clause ]
//	    [ { SORT | NOSORT } ]
//	    [ REVERSE ]
//	    [ { VISIBLE | INVISIBLE } ]
//	    [ logging_clause ]
//	    [ ONLINE ]
//
//	index_properties:
//	    [ global_partitioned_index | local_partitioned_index ]
//	    [ index_attributes ]
//	    [ domain_index_clause | XMLIndex_clause ]
func (p *Parser) parseCreateIndexAttributes(stmt *nodes.CreateIndexStmt) error {
	for {
		switch p.cur.Type {
		case kwREVERSE:
			stmt.Reverse = true
			p.advance()
		case kwTABLESPACE:
			p.advance()
			var parseErr422 error
			stmt.Tablespace, parseErr422 = p.parseIdentifier()
			if parseErr422 != nil {
				return parseErr422
			}
		case kwLOCAL:
			optStart := p.pos()
			stmt.Local = true
			p.advance()
			value := p.collectCreateIndexOptionValue()
			appendCreateIndexOption(stmt, "LOCAL", value, nodes.Loc{Start: optStart, End: p.prev.End})
		case kwGLOBAL:
			optStart := p.pos()
			stmt.Global = true
			p.advance()
			value := p.collectCreateIndexOptionValue()
			appendCreateIndexOption(stmt, "GLOBAL", value, nodes.Loc{Start: optStart, End: p.prev.End})
		case kwONLINE:
			stmt.Online = true
			p.advance()
		case kwPARALLEL:
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.Parallel = p.cur.Str
				p.advance()
			} else {
				stmt.Parallel = "PARALLEL"
			}
		case kwNOPARALLEL:
			stmt.Parallel = "NOPARALLEL"
			p.advance()
		case kwCOMPRESS:
			p.advance()
			// COMPRESS ADVANCED [LOW | HIGH] or COMPRESS [integer]
			if p.isIdentLikeStr("ADVANCED") {
				p.advance() // consume ADVANCED
				if p.isIdentLikeStr("LOW") {
					stmt.Compress = "ADVANCED LOW"
					p.advance()
				} else if p.isIdentLikeStr("HIGH") {
					stmt.Compress = "ADVANCED HIGH"
					p.advance()
				} else {
					stmt.Compress = "ADVANCED"
				}
			} else if p.cur.Type == tokICONST {
				stmt.Compress = p.cur.Str
				p.advance()
			} else {
				stmt.Compress = "COMPRESS"
			}
		case kwNOCOMPRESS:
			stmt.Compress = "NOCOMPRESS"
			p.advance()
		case kwLOGGING:
			stmt.Logging = true
			p.advance()
		case kwNOLOGGING:
			stmt.NoLogging = true
			p.advance()
		case kwINVISIBLE:
			stmt.Invisible = true
			p.advance()
		case kwPCTFREE:
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.PctFree = p.cur.Str
				p.advance()
			}
		default:
			// Check identifier-based keywords
			if p.isIdentLikeStr("VISIBLE") {
				stmt.Visible = true
				p.advance()
			} else if p.isIdentLikeStr("SORT") {
				stmt.Sort = true
				p.advance()
			} else if p.isIdentLikeStr("NOSORT") {
				stmt.NoSort = true
				p.advance()
			} else if p.isIdentLikeStr("INITRANS") {
				p.advance()
				if p.cur.Type == tokICONST {
					stmt.InitTrans = p.cur.Str
					p.advance()
				}
			} else if p.isIdentLikeStr("MAXTRANS") {
				p.advance()
				if p.cur.Type == tokICONST {
					stmt.MaxTrans = p.cur.Str
					p.advance()
				}
			} else if p.isIdentLikeStr("INDEXTYPE") {
				// domain_index_clause: INDEXTYPE IS [schema.]indextype
				p.advance() // consume INDEXTYPE
				if p.cur.Type == kwIS {
					p.advance() // consume IS
				}
				var parseErr423 error
				stmt.IndexType, parseErr423 = p.parseObjectName()
				if parseErr423 !=
					// Optional LOCAL and PARAMETERS
					nil {
					return parseErr423
				}

				if p.cur.Type == kwLOCAL {
					optStart := p.pos()
					stmt.Local = true
					p.advance()
					value := p.collectCreateIndexOptionValue()
					appendCreateIndexOption(stmt, "LOCAL", value, nodes.Loc{Start: optStart, End: p.prev.End})
				}
				// Optional parallel_clause
				if p.cur.Type == kwPARALLEL {
					p.advance()
					if p.cur.Type == tokICONST {
						stmt.Parallel = p.cur.Str
						p.advance()
					} else {
						stmt.Parallel = "PARALLEL"
					}
				} else if p.cur.Type == kwNOPARALLEL {
					stmt.Parallel = "NOPARALLEL"
					p.advance()
				}
				// Optional PARAMETERS
				if p.isIdentLikeStr("PARAMETERS") {
					p.advance() // consume PARAMETERS
					if p.cur.Type == '(' {
						p.advance()
						if p.cur.Type == tokSCONST {
							stmt.Parameters = p.cur.Str
							p.advance()
						}
						if p.cur.Type == ')' {
							p.advance()
						}
					}
				}
			} else if p.isIdentLikeStr("INDEXING") {
				p.advance() // consume INDEXING
				if p.isIdentLikeStr("FULL") {
					stmt.IndexingFull = true
					p.advance()
				} else if p.isIdentLikeStr("PARTIAL") {
					stmt.IndexingPartial = true
					p.advance()
				}
			} else if p.isIdentLikeStr("PARAMETERS") {
				p.advance() // consume PARAMETERS
				if p.cur.Type == '(' {
					p.advance()
					if p.cur.Type == tokSCONST {
						stmt.Parameters = p.cur.Str
						p.advance()
					}
					if p.cur.Type == ')' {
						p.advance()
					}
				}
			} else if p.isIdentLikeStr("STORAGE") {
				optStart := p.pos()
				p.advance()
				if p.cur.Type != '(' {
					return p.syntaxErrorAtCur()
				}
				value := p.collectCreateIndexOptionValue()
				appendCreateIndexOption(stmt, "STORAGE", value, nodes.Loc{Start: optStart, End: p.prev.End})
			} else if p.isIdentLikeStr("ANNOTATIONS") {
				optStart := p.pos()
				p.advance()
				value := p.collectCreateIndexOptionValue()
				appendCreateIndexOption(stmt, "ANNOTATIONS", value, nodes.Loc{Start: optStart, End: p.prev.End})
			} else {
				return nil
			}
		}
	}
	return nil
}

func appendCreateIndexOption(stmt *nodes.CreateIndexStmt, key, value string, loc nodes.Loc) {
	if loc.End <= loc.Start {
		return
	}
	if stmt.Options == nil {
		stmt.Options = &nodes.List{}
	}
	stmt.Options.Items = append(stmt.Options.Items, &nodes.DDLOption{Key: key, Value: value, Loc: loc})
}

func (p *Parser) collectCreateIndexOptionValue() string {
	tokens := p.collectDDLTokensUntil(p.isCreateIndexOptionBoundary)
	return strings.Join(tokens, " ")
}

func (p *Parser) isCreateIndexOptionBoundary() bool {
	switch p.cur.Type {
	case ';', tokEOF, kwCOMPRESS, kwDEFERRED, kwGLOBAL, kwIMMEDIATE, kwINVISIBLE,
		kwLOCAL, kwLOGGING, kwNOCOMPRESS, kwNOLOGGING, kwONLINE, kwPARALLEL,
		kwNOPARALLEL, kwPCTFREE, kwREVERSE, kwTABLESPACE:
		return true
	}
	if p.isIdentLike() {
		switch p.cur.Str {
		case "ANNOTATIONS", "INDEXING", "INDEXTYPE", "INITRANS", "MAXTRANS",
			"NOSORT", "PARAMETERS", "SORT", "STORAGE", "VISIBLE":
			return true
		}
	}
	return false
}

// collectParenthesizedBlock consumes and returns a balanced parenthesized block if present.
func (p *Parser) collectParenthesizedBlock() string {
	if p.cur.Type != '(' {
		return ""
	}
	var tokens []string
	depth := 1
	tokens = append(tokens, p.ddlOptionTokenText(p.cur))
	p.advance()
	for depth > 0 && p.cur.Type != tokEOF {
		tokens = append(tokens, p.ddlOptionTokenText(p.cur))
		if p.cur.Type == '(' {
			depth++
		} else if p.cur.Type == ')' {
			depth--
		}
		p.advance()
	}
	return strings.Join(tokens, " ")
}

// parseIndexColumnList parses a comma-separated list of index columns.
func (p *Parser) parseIndexColumnList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		col, parseErr424 := p.parseIndexColumn()
		if parseErr424 != nil {
			return nil, parseErr424
		}
		if col == nil {
			break
		}
		list.Items = append(list.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return list, nil
}

// parseIndexColumn parses a single index column expression with optional ASC/DESC.
//
//	index_expr [ ASC | DESC ]
func (p *Parser) parseIndexColumn() (*nodes.IndexColumn, error) {
	start := p.pos()
	expr, parseErr425 := p.parseExpr()
	if parseErr425 != nil {
		return nil, parseErr425
	}
	if expr == nil {
		return nil, nil
	}

	col := &nodes.IndexColumn{
		Expr: expr,
		Loc:  nodes.Loc{Start: start},
	}

	// ASC | DESC
	switch p.cur.Type {
	case kwASC:
		col.Dir = nodes.SORTBY_ASC
		p.advance()
	case kwDESC:
		col.Dir = nodes.SORTBY_DESC
		p.advance()
	}

	// NULLS FIRST | NULLS LAST
	if p.cur.Type == kwNULLS {
		p.advance()
		switch p.cur.Type {
		case kwFIRST:
			col.NullOrder = nodes.SORTBY_NULLS_FIRST
			p.advance()
		case kwLAST:
			col.NullOrder = nodes.SORTBY_NULLS_LAST
			p.advance()
		}
	}

	col.Loc.End = p.prev.End
	return col, nil
}

// parseIndexOptions is a legacy wrapper kept for backward compatibility.
// It delegates to parseCreateIndexAttributes.
func (p *Parser) parseIndexOptions(stmt *nodes.CreateIndexStmt) error {
	parseErr426 := p.parseCreateIndexAttributes(stmt)
	if parseErr426 !=

		// ---------------------------------------------------------------------------
		// CREATE INDEXTYPE
		// ---------------------------------------------------------------------------
		nil {
		return parseErr426
	}
	return nil
}

// parseCreateIndextypeStmt parses a CREATE INDEXTYPE statement.
// Called after CREATE [OR REPLACE] has been consumed and INDEXTYPE consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-INDEXTYPE.html
//
//	CREATE [ OR REPLACE ] INDEXTYPE [ IF NOT EXISTS ] [ schema. ] indextype
//	    [ SHARING = { METADATA | NONE } ]
//	    FOR [ schema. ] operator ( parameter_type [, parameter_type ]... )
//	        [, [ schema. ] operator ( parameter_type [, parameter_type ]... ) ]...
//	    using_type_clause
//	    [ WITH LOCAL [ RANGE ] PARTITION ]
//	    [ storage_table_clause ]
//	    [ array_DML_clause ]
func (p *Parser) parseCreateIndextypeStmt(start int, orReplace bool) (*nodes.CreateIndextypeStmt, error) {
	stmt := &nodes.CreateIndextypeStmt{
		OrReplace: orReplace,
		Loc:       nodes.Loc{Start: start},
	}

	// Optional IF NOT EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			stmt.IfNotExists = true
			p.advance() // consume IF
			p.advance() // consume NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // consume EXISTS
			}
		}
	}
	var parseErr427 error

	// Indextype name
	stmt.Name, parseErr427 = p.parseObjectName()
	if parseErr427 !=

		// Optional SHARING = { METADATA | NONE }
		nil {
		return nil, parseErr427
	}

	if p.isIdentLikeStr("SHARING") {
		p.advance() // consume SHARING
		if p.cur.Type == '=' {
			p.advance() // consume =
		}
		if p.isIdentLikeStr("METADATA") {
			stmt.Sharing = "METADATA"
			p.advance()
		} else if p.isIdentLikeStr("NONE") {
			stmt.Sharing = "NONE"
			p.advance()
		} else if p.isIdentLike() {
			stmt.Sharing = p.cur.Str
			p.advance()
		}
	}

	// FOR operator_list
	if p.cur.Type == kwFOR {
		p.advance()
		var // consume FOR
		parseErr428 error
		stmt.Operators, parseErr428 = p.parseIndextypeOperatorList()
		if parseErr428 !=

			// USING implementation_type
			nil {
			return nil, parseErr428
		}
	}

	if p.cur.Type == kwUSING {
		p.advance()
		var // consume USING
		parseErr429 error
		stmt.UsingType, parseErr429 = p.parseObjectName()
		if parseErr429 !=

			// Optional WITH LOCAL [RANGE] PARTITION
			nil {
			return nil, parseErr429
		}
	}

	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == kwLOCAL {
			stmt.WithLocal = true
			p.advance() // consume LOCAL
			if p.isIdentLikeStr("RANGE") {
				stmt.WithRange = true
				p.advance() // consume RANGE
			}
			if p.cur.Type == kwPARTITION {
				p.advance() // consume PARTITION
			}
		} else if p.isIdentLikeStr("SYSTEM") || p.cur.Type == kwUSER {
			var parseErr430 error
			// storage_table_clause
			stmt.StorageTable, parseErr430 = p.parseIdentifier()
			if parseErr430 !=
				// MANAGED STORAGE TABLES
				nil {
				return nil, parseErr430
			}

			if p.isIdentLikeStr("MANAGED") {
				p.advance()
			}
			if p.isIdentLikeStr("STORAGE") {
				p.advance()
			}
			if p.isIdentLikeStr("TABLES") {
				p.advance()
			}
		} else if p.isIdentLikeStr("ARRAY") {
			// array_DML_clause
			stmt.ArrayDML = true
			p.advance() // consume ARRAY
			if p.isIdentLikeStr("DML") {
				p.advance() // consume DML
			}
			// Optional type list
			p.collectParenthesizedBlock()
		}
	}

	// More optional WITH clauses
	for p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == kwLOCAL {
			stmt.WithLocal = true
			p.advance()
			if p.isIdentLikeStr("RANGE") {
				stmt.WithRange = true
				p.advance()
			}
			if p.cur.Type == kwPARTITION {
				p.advance()
			}
		} else if p.isIdentLikeStr("SYSTEM") || p.cur.Type == kwUSER {
			var parseErr431 error
			stmt.StorageTable, parseErr431 = p.parseIdentifier()
			if parseErr431 != nil {
				return nil, parseErr431
			}
			if p.isIdentLikeStr("MANAGED") {
				p.advance()
			}
			if p.isIdentLikeStr("STORAGE") {
				p.advance()
			}
			if p.isIdentLikeStr("TABLES") {
				p.advance()
			}
		} else if p.isIdentLikeStr("ARRAY") {
			stmt.ArrayDML = true
			p.advance()
			if p.isIdentLikeStr("DML") {
				p.advance()
			}
			p.collectParenthesizedBlock()
		} else {
			break
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseIndextypeOperatorList parses operator list for CREATE INDEXTYPE.
//
//	[ schema. ] operator ( parameter_type [, parameter_type ]... )
//	[, [ schema. ] operator ( parameter_type [, parameter_type ]... ) ]...
func (p *Parser) parseIndextypeOperatorList() ([]*nodes.IndextypeOp, error) {
	var ops []*nodes.IndextypeOp
	for {
		op, parseErr432 := p.parseIndextypeOp()
		if parseErr432 != nil {
			return nil, parseErr432
		}
		if op == nil {
			break
		}
		ops = append(ops, op)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return ops, nil
}

// parseIndextypeOp parses a single operator with parameter types.
func (p *Parser) parseIndextypeOp() (*nodes.IndextypeOp, error) {
	if !p.isIdentLike() {
		return nil, nil
	}
	start := p.pos()
	name, parseErr433 := p.parseObjectName()
	if parseErr433 != nil {
		return nil, parseErr433
	}
	op := &nodes.IndextypeOp{
		Name: name,
		Loc:  nodes.Loc{Start: start},
	}
	// ( parameter_type [, parameter_type ]... )
	if p.cur.Type == '(' {
		p.advance()
		var parseErr434 error
		op.ParamTypes, parseErr434 = p.parseTypeNameList()
		if parseErr434 != nil {
			return nil, parseErr434
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}
	op.Loc.End = p.prev.End
	return op, nil
}

// parseTypeNameList parses a comma-separated list of type names.
func (p *Parser) parseTypeNameList() ([]string, error) {
	var types []string
	for {
		tn, parseErr435 := p.parseTypeNameStr()
		if parseErr435 != nil {
			return nil, parseErr435
		}
		if tn == "" {
			break
		}
		types = append(types, tn)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return types, nil
}

// parseTypeName parses a single type name (e.g. NUMBER, VARCHAR2, schema.type_name).
func (p *Parser) parseTypeNameStr() (string, error) {
	if !p.isIdentLike() {
		return "", nil
	}
	name, parseErr436 := p.parseIdentifier()
	if parseErr436 !=
		// Handle schema.type_name
		nil {
		return "", parseErr436
	}

	if p.cur.Type == '.' {
		p.advance()
		parseValue63, parseErr64 := p.parseIdentifier()
		if parseErr64 !=

			// Handle parameterized types like VARCHAR2(100)
			nil {
			return "", parseErr64
		}
		name += "." + parseValue63
	}

	if p.cur.Type == '(' {
		name += "("
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.cur.Type == tokICONST {
				name += p.cur.Str
			} else if p.isIdentLike() {
				name += p.cur.Str
			} else if p.cur.Type == ',' {
				name += ","
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			name += ")"
			p.advance()
		}
	}
	return name, nil
}

// ---------------------------------------------------------------------------
// ALTER INDEXTYPE
// ---------------------------------------------------------------------------

// parseAlterIndextypeStmt parses an ALTER INDEXTYPE statement.
// Called after ALTER has been consumed and INDEXTYPE consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-INDEXTYPE.html
//
//	ALTER INDEXTYPE [ IF EXISTS ] [ schema. ] indextype
//	    { { ADD | DROP } [ schema. ] operator ( parameter_type [, parameter_type ]... )
//	        [ , { ADD | DROP } [ schema. ] operator ( parameter_type [, parameter_type ]... ) ]...
//	        [ using_type_clause ]
//	    | COMPILE
//	    }
//	    [ WITH LOCAL PARTITION ]
//	    [ storage_table_clause ]
func (p *Parser) parseAlterIndextypeStmt(start int) (*nodes.AlterIndextypeStmt, error) {
	stmt := &nodes.AlterIndextypeStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}
	var parseErr437 error

	// Indextype name
	stmt.Name, parseErr437 = p.parseObjectName()
	if parseErr437 !=

		// Action
		nil {
		return nil, parseErr437
	}

	if p.isIdentLikeStr("COMPILE") {
		stmt.Action = "COMPILE"
		p.advance()
	} else if p.cur.Type == kwADD || p.cur.Type == kwDROP {
		stmt.Action = "ADD_DROP"
		// Parse modifications list
		for p.cur.Type == kwADD || p.cur.Type == kwDROP {
			modStart := p.pos()
			isAdd := p.cur.Type == kwADD
			p.advance() // consume ADD/DROP
			name, parseErr438 := p.parseObjectName()
			if parseErr438 != nil {
				return nil, parseErr438
			}
			mod := &nodes.IndextypeModOp{
				Add:  isAdd,
				Name: name,
				Loc:  nodes.Loc{Start: modStart},
			}
			if p.cur.Type == '(' {
				p.advance()
				var parseErr439 error
				mod.ParamTypes, parseErr439 = p.parseTypeNameList()
				if parseErr439 != nil {
					return nil, parseErr439
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
			mod.Loc.End = p.prev.End
			stmt.Modifications = append(stmt.Modifications, mod)
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		// Optional using_type_clause
		if p.cur.Type == kwUSING {
			p.advance()
			var // consume USING
			parseErr440 error
			stmt.UsingType, parseErr440 = p.parseObjectName()
			if parseErr440 !=
				// Optional array_DML_clause
				nil {
				return nil, parseErr440
			}

			if p.cur.Type == kwWITH {
				p.advance()
				if p.isIdentLikeStr("ARRAY") {
					stmt.ArrayDML = true
					p.advance() // consume ARRAY
					if p.isIdentLikeStr("DML") {
						p.advance() // consume DML
					}
					p.collectParenthesizedBlock()
				}
			}
		}
	}

	// Optional WITH LOCAL PARTITION / storage_table_clause
	for p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == kwLOCAL {
			stmt.WithLocal = true
			p.advance()
			if p.cur.Type == kwPARTITION {
				p.advance()
			}
		} else if p.isIdentLikeStr("SYSTEM") || p.cur.Type == kwUSER {
			var parseErr441 error
			stmt.StorageTable, parseErr441 = p.parseIdentifier()
			if parseErr441 != nil {
				return nil, parseErr441
			}
			if p.isIdentLikeStr("MANAGED") {
				p.advance()
			}
			if p.isIdentLikeStr("STORAGE") {
				p.advance()
			}
			if p.isIdentLikeStr("TABLES") {
				p.advance()
			}
		} else {
			break
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE OPERATOR
// ---------------------------------------------------------------------------

// parseCreateOperatorStmt parses a CREATE OPERATOR statement.
// Called after CREATE [OR REPLACE | IF NOT EXISTS] has been consumed and OPERATOR consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-OPERATOR.html
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ] OPERATOR
//	    [ schema. ] operator
//	    binding_clause
//	    [ SHARING = { METADATA | NONE } ]
func (p *Parser) parseCreateOperatorStmt(start int, orReplace bool, ifNotExists bool) (*nodes.CreateOperatorStmt, error) {
	stmt := &nodes.CreateOperatorStmt{
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Loc:         nodes.Loc{Start: start},
	}
	var parseErr442 error

	// Operator name
	stmt.Name, parseErr442 = p.parseObjectName()
	if parseErr442 !=

		// binding_clause: BINDING (types) RETURN type USING func
		// There can be multiple bindings (comma-separated in the BNF for CREATE, but
		// typically one binding per CREATE OPERATOR)
		nil {
		return nil, parseErr442
	}

	if p.isIdentLikeStr("BINDING") {
		binding, parseErr443 := p.parseOperatorBinding()
		if parseErr443 != nil {
			return nil, parseErr443
		}
		if binding != nil {
			stmt.Bindings = append(stmt.Bindings, binding)
		}
	}

	// Optional SHARING = { METADATA | NONE }
	if p.isIdentLikeStr("SHARING") {
		p.advance() // consume SHARING
		if p.cur.Type == '=' {
			p.advance() // consume =
		}
		if p.isIdentLikeStr("METADATA") {
			stmt.Sharing = "METADATA"
			p.advance()
		} else if p.isIdentLikeStr("NONE") {
			stmt.Sharing = "NONE"
			p.advance()
		} else if p.isIdentLike() {
			stmt.Sharing = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseOperatorBinding parses a BINDING clause.
//
//	BINDING ( [ parameter_type [, parameter_type ]... ] )
//	    RETURN return_type
//	    [ implementation_clause ]
//	    using_function_clause
func (p *Parser) parseOperatorBinding() (*nodes.OperatorBinding, error) {
	if !p.isIdentLikeStr("BINDING") {
		return nil, nil
	}
	start := p.pos()
	p.advance() // consume BINDING

	binding := &nodes.OperatorBinding{
		Loc: nodes.Loc{Start: start},
	}

	// ( parameter_types )
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type != ')' {
			var parseErr444 error
			binding.ParamTypes, parseErr444 = p.parseTypeNameList()
			if parseErr444 != nil {
				return nil, parseErr444
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// RETURN return_type
	if p.cur.Type == kwRETURN {
		p.advance() // consume RETURN
		// Return type may be in parens: RETURN ( datatype )
		if p.cur.Type == '(' {
			p.advance()
			var parseErr445 error
			binding.ReturnType, parseErr445 = p.parseTypeNameStr()
			if parseErr445 != nil {
				return nil, parseErr445
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		} else {
			var parseErr446 error
			binding.ReturnType, parseErr446 = p.parseTypeNameStr()
			if parseErr446 !=

				// Optional implementation_clause
				// { ANCILLARY TO primary_operator (types) | context_clause }
				nil {
				return nil, parseErr446
			}
		}
	}

	if p.isIdentLikeStr("ANCILLARY") {
		p.advance() // consume ANCILLARY
		if p.cur.Type == kwTO {
			p.advance() // consume TO
		}
		var parseErr447 error
		binding.AncillaryTo, parseErr447 = p.parseObjectName()
		if parseErr447 != nil {
			return nil, parseErr447
		}
		if p.cur.Type == '(' {
			p.advance()
			var parseErr448 error
			binding.AncillaryParams, parseErr448 = p.parseTypeNameList()
			if parseErr448 != nil {
				return nil, parseErr448
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	} else if p.cur.Type == kwWITH {
		// context_clause: WITH INDEX CONTEXT, SCAN CONTEXT implementation_type
		p.advance() // consume WITH
		if p.cur.Type == kwINDEX {
			p.advance() // consume INDEX
			if p.isIdentLikeStr("CONTEXT") {
				binding.WithIndexCtx = true
				p.advance() // consume CONTEXT
			}
			if p.cur.Type == ',' {
				p.advance() // consume ','
			}
			if p.isIdentLikeStr("SCAN") {
				p.advance() // consume SCAN
				if p.isIdentLikeStr("CONTEXT") {
					p.advance() // consume CONTEXT
				}
				var parseErr449 error
				binding.ScanCtxType, parseErr449 = p.parseIdentifier()
				if parseErr449 != nil {
					return nil, parseErr449
				}
			}
		} else if p.isIdentLikeStr("COLUMN") {
			binding.WithColumnCtx = true
			p.advance() // consume COLUMN
			if p.isIdentLikeStr("CONTEXT") {
				p.advance() // consume CONTEXT
			}
		}
	}

	// COMPUTE ANCILLARY DATA
	if p.isIdentLikeStr("COMPUTE") {
		binding.ComputeAnc = true
		p.advance() // consume COMPUTE
		if p.isIdentLikeStr("ANCILLARY") {
			p.advance() // consume ANCILLARY
		}
		if p.isIdentLikeStr("DATA") {
			p.advance() // consume DATA
		}
	}

	// WITH COLUMN CONTEXT (can appear after COMPUTE ANCILLARY DATA)
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.isIdentLikeStr("COLUMN") {
			binding.WithColumnCtx = true
			p.advance() // consume COLUMN
			if p.isIdentLikeStr("CONTEXT") {
				p.advance() // consume CONTEXT
			}
		}
	}

	// USING function
	if p.cur.Type == kwUSING {
		p.advance()
		var // consume USING
		parseErr450 error
		binding.UsingFunc, parseErr450 = p.parseObjectName()
		if parseErr450 != nil {
			return nil, parseErr450
		}
	}

	binding.Loc.End = p.prev.End
	return binding, nil
}

// ---------------------------------------------------------------------------
// ALTER OPERATOR
// ---------------------------------------------------------------------------

// parseAlterOperatorStmt parses an ALTER OPERATOR statement.
// Called after ALTER has been consumed and OPERATOR consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-OPERATOR.html
//
//	ALTER OPERATOR [ IF EXISTS ] [ schema. ] operator
//	    { add_binding_clause
//	    | drop_binding_clause
//	    | COMPILE
//	    }
func (p *Parser) parseAlterOperatorStmt(start int) (*nodes.AlterOperatorStmt, error) {
	stmt := &nodes.AlterOperatorStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}
	var parseErr451 error

	// Operator name
	stmt.Name, parseErr451 = p.parseObjectName()
	if parseErr451 !=

		// Action
		nil {
		return nil, parseErr451
	}

	switch {
	case p.isIdentLikeStr("COMPILE"):
		stmt.Action = "COMPILE"
		p.advance()
	case p.cur.Type == kwADD:
		stmt.Action = "ADD_BINDING"
		p.advance() // consume ADD
		binding, parseErr452 := p.parseOperatorBinding()
		if parseErr452 != nil {
			return nil, parseErr452
		}
		stmt.Binding = binding
	case p.cur.Type == kwDROP:
		stmt.Action = "DROP_BINDING"
		p.advance() // consume DROP
		if p.isIdentLikeStr("BINDING") {
			p.advance() // consume BINDING
		}
		// ( datatype [, datatype ]... )
		if p.cur.Type == '(' {
			p.advance()
			var parseErr453 error
			stmt.DropTypes, parseErr453 = p.parseTypeNameList()
			if parseErr453 != nil {
				return nil, parseErr453
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		// Optional FORCE
		if p.cur.Type == kwFORCE {
			stmt.DropForce = true
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
