package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// parseCreateTableStmt parses a CREATE TABLE statement. start is the offset of
// CREATE; the caller has consumed CREATE [OR REPLACE] [TEMPORARY].
//
// Ref: https://mariadb.com/kb/en/create-table/
//
//	CREATE [OR REPLACE] [TEMPORARY] TABLE [IF NOT EXISTS] tbl_name
//	    (create_definition,...) [table_options] [partition_options]
//	CREATE [OR REPLACE] [TEMPORARY] TABLE [IF NOT EXISTS] tbl_name
//	    LIKE old_tbl_name
//	CREATE [OR REPLACE] [TEMPORARY] TABLE [IF NOT EXISTS] tbl_name
//	    [AS] query_expression
func (p *Parser) parseCreateTableStmt(start int, orReplace, temporary bool) (*nodes.CreateTableStmt, error) {
	p.advance() // consume TABLE

	stmt := &nodes.CreateTableStmt{
		Loc:       nodes.Loc{Start: start},
		OrReplace: orReplace,
		Temporary: temporary,
	}

	ifNotExists, err := p.parseCreateIfNotExists(orReplace)
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	// Completion: after CREATE TABLE, identifier context (no specific candidates).
	p.checkCursor()
	if p.collectMode() {
		// No specific candidates for table name — user defines a new name.
		return nil, &ParseError{Message: "collecting"}
	}

	// Table name
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = ref

	// CREATE TABLE ... LIKE
	if p.cur.Type == kwLIKE {
		p.advance()
		// Completion: after LIKE, offer table_ref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
		}
		likeRef, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Like = likeRef
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// create sets parsing_place to BEFORE_OPT_LIST before the table name
	// (MariaDB 11.8.8 sql/sql_yacc.yy line 2531), and until something sets it
	// back, Lex_input_stream::lex_token (sql/sql_lex.cc line 1888) reads two
	// tokens differently:
	//   - VALUE is VALUES, starting a table value constructor: CREATE TABLE
	//     t2 VALUE (1) and CREATE TABLE t2 AS VALUE (1) are queries, CREATE
	//     TABLE t2 (a INT) VALUE (1) is not;
	//   - '(' followed by WITH is LEFT_PAREN_WITH, which only
	//     create_select_query_expression takes, the query of CREATE TABLE and
	//     of INSERT (see parseCreateSelect).
	// The '(' of the column list, PARTITION, and the UNION of a UNION table
	// option set it back. Table options before the query, IGNORE, REPLACE and
	// AS leave it.
	beforeOptList := true

	// CREATE TABLE ... AS SELECT
	if p.cur.Type == kwAS {
		p.advance()
		sel, err := p.parseCreateSelect(beforeOptList)
		if err != nil {
			return nil, err
		}
		stmt.Select = sel
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Parenthesized create definitions (columns and constraints)
	if p.cur.Type == '(' {
		next := p.peekNext()
		// Check for CREATE TABLE ... (query) — query source without AS. Let
		// parseCreateSelect consume the parentheses and any following set-op /
		// ORDER BY / LIMIT so "(SELECT 1) UNION (SELECT 2)" stays one source.
		if next.Type == kwSELECT || next.Type == kwVALUES || next.Type == kwWITH || next.Type == '(' {
			sel, err := p.parseCreateSelect(beforeOptList)
			if err != nil {
				return nil, err
			}
			stmt.Select = sel
			stmt.Loc.End = p.prev.End
			return stmt, nil
		}
		// Check for CREATE TABLE ... (LIKE old_tbl_name) — parenthesized LIKE
		if next.Type == kwLIKE {
			p.advance() // consume '('
			p.advance() // consume LIKE
			likeRef, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.Like = likeRef
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			stmt.Loc.End = p.prev.End
			return stmt, nil
		}

		beforeOptList = false
		p.advance() // consume '('
		if err := p.parseCreateDefinitions(stmt); err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// CREATE TABLE t (cols) [IGNORE|REPLACE] [AS] SELECT ...
	if p.cur.Type == kwAS {
		p.advance()
		sel, err := p.parseCreateSelect(beforeOptList)
		if err != nil {
			return nil, err
		}
		stmt.Select = sel
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Table options. A comma may stand between two of them and never ends the
	// list (create_table_options in MariaDB 11.8 sql/sql_yacc.yy), so a comma
	// with no option after it is a syntax error rather than the end of the
	// options: "ENGINE=InnoDB, (SELECT 1)" is 1064.
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
		if opt.Name == "UNION" {
			beforeOptList = false
		}
		p.match(',')
	}

	// Partition clause
	if p.cur.Type == kwPARTITION {
		beforeOptList = false
		part, err := p.parsePartitionClause()
		if err != nil {
			return nil, err
		}
		stmt.Partitions = part
	}

	// IGNORE / REPLACE before SELECT in CTAS
	if p.cur.Type == kwIGNORE {
		p.advance()
		stmt.Ignore = true
	} else if p.cur.Type == kwREPLACE {
		p.advance()
		stmt.Replace = true
	}

	// [AS] query — the query of opt_create_select (MariaDB 11.8.8
	// sql/sql_yacc.yy line 4821), with or without AS, after the options and
	// partitions: SELECT, a table value constructor, WITH, or a
	// parenthesized query. IGNORE or REPLACE is opt_duplicate of that rule,
	// so a query must follow it.
	switch {
	case p.cur.Type == kwAS:
		p.advance()
		sel, err := p.parseCreateSelect(beforeOptList)
		if err != nil {
			return nil, err
		}
		stmt.Select = sel
	case p.cur.Type == kwSELECT || p.cur.Type == kwVALUES || p.cur.Type == kwWITH || p.cur.Type == '(' ||
		p.cur.Type == kwVALUE && beforeOptList:
		sel, err := p.parseCreateSelect(beforeOptList)
		if err != nil {
			return nil, err
		}
		stmt.Select = sel
	case stmt.Ignore || stmt.Replace:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateSelect parses the query of CREATE TABLE ... [AS] query,
// create_select_query_expression in MariaDB 11.8.8 sql/sql_yacc.yy (line
// 4832):
//
//	create_select_query_expression:
//	      query_expression
//	    | LEFT_PAREN_WITH with_clause query_expression_no_with_clause ')'
//
// beforeOptList reports whether the query starts where MariaDB's lexer still
// reads VALUE as VALUES and '(' WITH as LEFT_PAREN_WITH (see
// parseCreateTableStmt).
//
// query_expression is no statement, so it takes no INTO, and neither does any
// query nested in it (see intoRefused).
func (p *Parser) parseCreateSelect(beforeOptList bool) (*nodes.SelectStmt, error) {
	saved := p.intoRefused
	p.intoRefused = true
	defer func() { p.intoRefused = saved }()

	switch {
	case p.cur.Type == kwVALUE && beforeOptList:
		values, err := p.parseValuesConstructor()
		if err != nil {
			return nil, err
		}
		return p.parseQueryAfterValues(values)
	case p.cur.Type == '(' && p.peekNext().Type == kwWITH:
		return p.parseCreateSelectParenWith(beforeOptList)
	}
	return p.parseSelectStmt()
}

// parseCreateSelectParenWith parses a query of CREATE TABLE that starts with
// '(' WITH. A parenthesized query in query_expression holds no WITH
// (query_expression_body_ext_parens, sql/sql_yacc.yy line 9140), so the
// query is the LEFT_PAREN_WITH alternative or nothing, and the lexer makes
// that token only while parsing_place is BEFORE_OPT_LIST:
//
//	CREATE TABLE t (WITH c AS (SELECT 1) SELECT * FROM c)        -- a query
//	CREATE TABLE t ENGINE=InnoDB (WITH c AS (SELECT 1) SELECT * FROM c)
//	CREATE TABLE t (a INT) (WITH c AS (SELECT 1) SELECT * FROM c) -- 1064
//
// The alternative ends at its ')': no set operation, ORDER BY, LIMIT or
// locking clause follows it, so "(WITH ...) UNION SELECT 2" is 1064 where
// "(WITH ... UNION SELECT 2)" is a query. The tree is the one a
// parenthesized query has: a SelectStmt whose ParenSource carries the WITH.
func (p *Parser) parseCreateSelectParenWith(beforeOptList bool) (*nodes.SelectStmt, error) {
	start := p.pos()
	p.advance() // consume '('
	if !beforeOptList {
		return nil, p.syntaxErrorAtCur()
	}
	inner, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	if inner == nil { // collecting completion candidates
		return nil, nil
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.SelectStmt{
		Loc:         nodes.Loc{Start: start, End: p.prev.End},
		ParenSource: inner,
	}, nil
}

// parseCreateDefinitions parses column definitions and table constraints
// inside the parentheses of a CREATE TABLE statement.
func (p *Parser) parseCreateDefinitions(stmt *nodes.CreateTableStmt) error {
	for {

		// Completion: at start of a column/constraint definition, offer constraint keywords.
		p.checkCursor()
		if p.collectMode() {
			for _, t := range []int{kwPRIMARY, kwUNIQUE, kwINDEX, kwKEY, kwFOREIGN, kwCHECK, kwCONSTRAINT} {
				p.addTokenCandidate(t)
			}
			return &ParseError{Message: "collecting"}
		}

		// PERIOD FOR SYSTEM_TIME (start_col, end_col) — a table element, not a
		// column. Intercepted before parseColumnDef would read PERIOD as a name.
		if p.atPeriodForSystemTime() {
			if err := p.parsePeriodFor(stmt); err != nil {
				return err
			}
		} else if p.isTableConstraintStart() {
			// Table-level constraints: PRIMARY KEY, UNIQUE, INDEX, KEY, FULLTEXT, SPATIAL, FOREIGN KEY, CONSTRAINT, CHECK
			constr, err := p.parseTableConstraint()
			if err != nil {
				return err
			}
			stmt.Constraints = append(stmt.Constraints, constr)
		} else {
			// Column definition
			col, err := p.parseColumnDef()
			if err != nil {
				return err
			}
			stmt.Columns = append(stmt.Columns, col)
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return nil
}

// isTableConstraintStart returns true if current token starts a table-level constraint.
func (p *Parser) isTableConstraintStart() bool {
	switch p.cur.Type {
	case kwPRIMARY, kwUNIQUE, kwINDEX, kwKEY, kwFULLTEXT, kwSPATIAL, kwFOREIGN, kwCHECK:
		return true
	case kwCONSTRAINT:
		return true
	}
	return false
}

// atPeriodForSystemTime reports whether the cursor begins a
// PERIOD FOR SYSTEM_TIME table element. PERIOD is non-reserved (matched by
// text); FOR follows. Application-time PERIOD FOR <name> is not yet supported.
func (p *Parser) atPeriodForSystemTime() bool {
	return p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "PERIOD") &&
		p.peekNext().Type == kwFOR
}

// parsePeriodFor parses a CREATE TABLE PERIOD FOR <name> (start, end) element,
// routing SYSTEM_TIME to the system-versioning period and any other name to the
// application-time period.
func (p *Parser) parsePeriodFor(stmt *nodes.CreateTableStmt) error {
	name, startCol, endCol, err := p.parsePeriodForCols()
	if err != nil {
		return err
	}
	if strings.EqualFold(name, "SYSTEM_TIME") {
		stmt.PeriodStartCol = startCol
		stmt.PeriodEndCol = endCol
	} else {
		if stmt.AppPeriodName != "" {
			stmt.AppPeriodDuplicate = true
		}
		stmt.AppPeriodName = name
		stmt.AppPeriodStartCol = startCol
		stmt.AppPeriodEndCol = endCol
	}
	return nil
}

// parsePeriodForCols parses PERIOD FOR <name> (start_col, end_col) and returns
// the period name and the two column names. Shared by CREATE TABLE and
// ALTER TABLE ... ADD PERIOD. PERIOD is the current token.
func (p *Parser) parsePeriodForCols() (name, startCol, endCol string, err error) {
	p.advance() // PERIOD
	p.advance() // FOR
	// The period name follows the identifier rule (non-reserved keywords allowed),
	// like the start/end column names below.
	if name, _, err = p.parseIdent(); err != nil {
		return "", "", "", err
	}
	if _, err = p.expect('('); err != nil {
		return "", "", "", err
	}
	if startCol, _, err = p.parseIdent(); err != nil {
		return "", "", "", err
	}
	if _, err = p.expect(','); err != nil {
		return "", "", "", err
	}
	if endCol, _, err = p.parseIdent(); err != nil {
		return "", "", "", err
	}
	if _, err = p.expect(')'); err != nil {
		return "", "", "", err
	}
	return name, startCol, endCol, nil
}

// parseAlterDropPeriod parses the DROP PERIOD arms of alter_list_item in
// MariaDB 11.8 sql/sql_yacc.yy into an ATDropPeriod. PERIOD is the current
// token.
//
//	DROP PERIOD_SYM FOR_SYSTEM_TIME_SYM
//	DROP PERIOD_SYM opt_if_exists_table_element FOR_SYM ident
//
// The first drops the system-versioning period, the second an
// application-time period. MariaDB's lexer merges FOR and an unquoted
// SYSTEM_TIME into one FOR_SYSTEM_TIME_SYM token, which the IF EXISTS arm does
// not take, so IF EXISTS FOR SYSTEM_TIME is 1064; a quoted `SYSTEM_TIME` stays
// an ident there. The period name follows the identifier rule (non-reserved
// keywords allowed).
func (p *Parser) parseAlterDropPeriod(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // PERIOD
	ifExists, err := p.parseOptIfExists()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	if ifExists && p.atUnquotedKeyword("SYSTEM_TIME") {
		return nil, p.syntaxErrorAtCur()
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	cmd.Type = nodes.ATDropPeriod
	cmd.IfExists = ifExists
	cmd.PeriodName = name
	cmd.Loc.End = p.prev.End
	return cmd, nil
}

// parseRowBound parses ROW START | ROW END (the system-versioning period bound a
// generated column carries). ROW is the current token.
func (p *Parser) parseRowBound() (nodes.RowBoundKind, error) {
	p.advance() // ROW
	switch p.cur.Type {
	case kwSTART:
		p.advance()
		return nodes.RowBoundStart, nil
	case kwEND:
		p.advance()
		return nodes.RowBoundEnd, nil
	default:
		return nodes.RowBoundNone, p.syntaxErrorAtCur()
	}
}

// parseColumnDef parses a column definition.
//
//	col_name data_type [NOT NULL | NULL] [DEFAULT expr] [AUTO_INCREMENT]
//	    [UNIQUE [KEY] | [PRIMARY] KEY] [COMMENT 'string']
//	    [COLLATE collation] [REFERENCES ...] [CHECK (...)]
//	    [GENERATED ALWAYS AS (expr) [VIRTUAL | STORED]]
//	    [ON UPDATE expr]
func (p *Parser) parseColumnDef() (*nodes.ColumnDef, error) {
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	col := &nodes.ColumnDef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// Completion: after column name, offer type candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		return nil, &ParseError{Message: "collecting"}
	}

	// Data type (optional for generated columns with GENERATED ALWAYS AS)
	if p.cur.Type != kwGENERATED && !p.isColumnConstraintStart() {
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		col.TypeName = dt
	}

	// Column constraints and options
	for {
		ok, err := p.parseColumnOption(col)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
	}

	// MariaDB rejects a column-level WITH/WITHOUT SYSTEM VERSIONING on any
	// generated column — both GENERATED ALWAYS AS (expr) and a ROW START/END
	// period column.
	if col.SystemVersioning != nodes.ColVersioningNone && col.Generated != nil {
		return nil, &ParseError{
			Message:  "SYSTEM VERSIONING is not allowed on a generated column",
			Position: p.cur.Loc,
		}
	}

	col.Loc.End = p.prev.End
	return col, nil
}

// isColumnConstraintStart returns true if current token starts a column constraint/option.
func (p *Parser) isColumnConstraintStart() bool {
	switch p.cur.Type {
	case kwNOT, kwNULL, kwDEFAULT, kwAUTO_INCREMENT, kwUNIQUE, kwPRIMARY,
		kwKEY, kwCOMMENT, kwCOLLATE, kwREFERENCES, kwCHECK, kwGENERATED, kwON,
		kwCOLUMN_FORMAT, kwSTORAGE:
		return true
	}
	if p.cur.Type == kwENGINE_ATTRIBUTE || p.cur.Type == kwSECONDARY_ENGINE_ATTRIBUTE {
		return true
	}
	if p.cur.Type == kwSRID {
		return true
	}
	return false
}

// parseColumnOption parses a single column option/constraint.
// Returns false if no option was found.
func (p *Parser) parseColumnOption(col *nodes.ColumnDef) (bool, error) {
	// Completion: offer column option keywords.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwNOT, kwNULL, kwDEFAULT, kwAUTO_INCREMENT, kwPRIMARY, kwUNIQUE, kwCOMMENT, kwCOLLATE, kwREFERENCES, kwCHECK, kwGENERATED} {
			p.addTokenCandidate(t)
		}
		return false, &ParseError{Message: "collecting"}
	}

	start := p.pos()

	switch p.cur.Type {
	case kwNOT:
		p.advance()
		if _, err := p.expect(kwNULL); err != nil {
			return false, err
		}
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrNotNull,
		})
		return true, nil

	case kwNULL:
		p.advance()
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrNull,
		})
		return true, nil

	case kwDEFAULT:
		p.advance()
		expr, err := p.parseDefaultValue()
		if err != nil {
			return false, err
		}
		col.DefaultValue = expr
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrDefault,
			Expr: expr,
		})
		return true, nil

	case kwAUTO_INCREMENT:
		p.advance()
		col.AutoIncrement = true
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrAutoIncrement,
		})
		return true, nil

	// SERIAL DEFAULT VALUE, the three words together, which
	// https://dev.mysql.com/doc/refman/8.0/en/numeric-type-syntax.html gives as
	// another way to write NOT NULL AUTO_INCREMENT UNIQUE. It is recorded as
	// those three and not as an attribute of its own, so everything downstream —
	// the catalog's implicit primary key among it — sees the column the engine
	// builds. All three words are required: MariaDB answers 1064 for
	// `a int serial` and for `a int serial default`.
	case kwSERIAL:
		p.advance()
		if _, err := p.expect(kwDEFAULT); err != nil {
			return false, err
		}
		if _, err := p.expect(kwVALUE); err != nil {
			return false, err
		}
		col.AutoIncrement = true
		loc := nodes.Loc{Start: start, End: p.prev.End}
		col.Constraints = append(col.Constraints,
			&nodes.ColumnConstraint{Loc: loc, Type: nodes.ColConstrNotNull},
			&nodes.ColumnConstraint{Loc: loc, Type: nodes.ColConstrAutoIncrement},
			&nodes.ColumnConstraint{Loc: loc, Type: nodes.ColConstrUnique},
		)
		return true, nil

	case kwUNIQUE:
		p.advance()
		p.match(kwKEY) // optional KEY
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrUnique,
		})
		return true, nil

	case kwPRIMARY:
		p.advance()
		if _, err := p.expect(kwKEY); err != nil {
			return false, err
		}
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrPrimaryKey,
		})
		return true, nil

	case kwKEY:
		// Bare KEY means PRIMARY KEY
		p.advance()
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrPrimaryKey,
		})
		return true, nil

	case kwCOMMENT:
		p.advance()
		// The comment is a quoted string and is required: the server answers 1064
		// for a bare COMMENT, as it does for the table option.
		if p.cur.Type != tokSCONST {
			return false, p.syntaxErrorAtCur()
		}
		col.Comment = p.cur.Str
		p.advance()
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrComment,
		})
		return true, nil

	case kwCOLLATE:
		p.advance()
		collName, _, err := p.parseIdent()
		if err != nil {
			return false, err
		}
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrCollate,
			Name: collName,
		})
		return true, nil

	case kwREFERENCES:
		constr, err := p.parseReferenceDefinition(start)
		if err != nil {
			return false, err
		}
		col.Constraints = append(col.Constraints, constr)
		return true, nil

	case kwVISIBLE:
		p.advance()
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrVisible,
		})
		return true, nil

	case kwINVISIBLE:
		p.advance()
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrInvisible,
		})
		return true, nil

	case kwCHECK:
		p.advance()
		if _, err := p.expect('('); err != nil {
			return false, err
		}
		expr, err := p.parseExpr()
		if err != nil {
			return false, err
		}
		if expr == nil {
			return false, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return false, err
		}
		cc := &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrCheck,
			Expr: expr,
		}
		// [[NOT] ENFORCED]
		if p.cur.Type == kwNOT {
			p.advance()
			if p.cur.Type == kwENFORCED {
				p.advance()
				cc.NotEnforced = true
			}
		} else if p.cur.Type == kwENFORCED {
			p.advance()
		}
		cc.Loc.End = p.prev.End
		col.Constraints = append(col.Constraints, cc)
		return true, nil

	case kwGENERATED:
		gen, err := p.parseGeneratedColumn()
		if err != nil {
			return false, err
		}
		col.Generated = gen
		return true, nil

	case kwAS:
		// AS (expr) [VIRTUAL|STORED] — shorthand for GENERATED ALWAYS AS (expr)
		if p.peekNext().Type == '(' {
			gen, err := p.parseGeneratedColumnShorthand()
			if err != nil {
				return false, err
			}
			col.Generated = gen
			return true, nil
		}
		// AS ROW START | ROW END — shorthand for GENERATED ALWAYS AS ROW START/END
		if p.peekNext().Type == kwROW {
			start := p.pos()
			p.advance() // AS
			rb, err := p.parseRowBound()
			if err != nil {
				return false, err
			}
			col.Generated = &nodes.GeneratedColumn{Loc: nodes.Loc{Start: start, End: p.prev.End}, RowBound: rb}
			return true, nil
		}
		return false, nil

	case kwWITH, kwWITHOUT:
		// WITH / WITHOUT SYSTEM VERSIONING — column-level system versioning.
		if p.peekNext().Type == kwSYSTEM {
			without := p.cur.Type == kwWITHOUT
			p.advance() // WITH / WITHOUT
			p.advance() // SYSTEM
			if p.cur.Type != tokIDENT || !strings.EqualFold(p.cur.Str, "VERSIONING") {
				return false, p.syntaxErrorAtCur()
			}
			p.advance() // VERSIONING
			if without {
				col.SystemVersioning = nodes.ColVersioningWithout
			} else {
				col.SystemVersioning = nodes.ColVersioningWith
				col.SystemVersioningWithSeen = true
			}
			return true, nil
		}
		return false, nil

	case kwON:
		// ON UPDATE CURRENT_TIMESTAMP
		if next := p.peekNext(); next.Type == kwUPDATE {
			p.advance() // consume ON
			p.advance() // consume UPDATE
			expr, err := p.parseExpr()
			if err != nil {
				return false, err
			}
			if expr == nil {
				return false, p.syntaxErrorAtCur()
			}
			col.OnUpdate = expr
			return true, nil
		}
		return false, nil

	case kwCOLUMN_FORMAT:
		// COLUMN_FORMAT {FIXED | DYNAMIC | DEFAULT}
		p.advance()
		var val string
		switch p.cur.Type {
		case kwFIXED:
			val = "FIXED"
			p.advance()
		case kwDYNAMIC:
			val = "DYNAMIC"
			p.advance()
		case kwDEFAULT:
			val = "DEFAULT"
			p.advance()
		default:
			return false, nil
		}
		col.ColumnFormat = val
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrColumnFormat,
			Name: val,
		})
		return true, nil

	case kwSTORAGE:
		// STORAGE {DISK | MEMORY}
		p.advance()
		var val string
		switch p.cur.Type {
		case kwDISK:
			val = "DISK"
			p.advance()
		case kwMEMORY:
			val = "MEMORY"
			p.advance()
		default:
			return false, nil
		}
		col.Storage = val
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrStorage,
			Name: val,
		})
		return true, nil
	}

	// ENGINE_ATTRIBUTE and SECONDARY_ENGINE_ATTRIBUTE column options. The string
	// is required, as it is on the table option: MariaDB answers 1064 for
	// `id INT SECONDARY_ENGINE_ATTRIBUTE`, for `id INT SECONDARY_ENGINE_ATTRIBUTE=`
	// and for both with ENGINE_ATTRIBUTE. Taking the value as optional parsed
	// each of them.
	if p.cur.Type == kwENGINE_ATTRIBUTE || p.cur.Type == kwSECONDARY_ENGINE_ATTRIBUTE {
		secondary := p.cur.Type == kwSECONDARY_ENGINE_ATTRIBUTE
		p.advance()
		p.matchEqual()
		val, err := p.consumeStringOptionValue()
		if err != nil {
			return false, err
		}
		constrType := nodes.ColConstrEngineAttribute
		if secondary {
			col.SecondaryEngineAttribute = val
			constrType = nodes.ColConstrSecondaryEngineAttribute
		} else {
			col.EngineAttribute = val
		}
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: constrType,
			Name: val,
		})
		return true, nil
	}

	// SRID column option (now lexes as kwSRID keyword token)
	if p.cur.Type == kwSRID {
		p.advance()
		if p.cur.Type != tokICONST {
			return false, nil
		}
		col.TypeName.SRID = int(p.cur.Ival)
		p.advance()
		return true, nil
	}

	return false, nil
}

// parseDefaultValue parses a DEFAULT value expression.
// Handles parenthesized expressions, literals, CURRENT_TIMESTAMP, NULL, etc.
func (p *Parser) parseDefaultValue() (nodes.ExprNode, error) {
	// Handle parenthesized expression: DEFAULT (expr)
	if p.cur.Type == '(' {
		start := p.pos()
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.ParenExpr{Loc: nodes.Loc{Start: start, End: p.prev.End}, Expr: expr}, nil
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return expr, nil
}

// parseReferenceDefinition parses a REFERENCES clause.
//
//	REFERENCES tbl_name (col,...) [MATCH FULL|PARTIAL|SIMPLE]
//	    [ON DELETE action] [ON UPDATE action]
//
//	action: RESTRICT | CASCADE | SET NULL | SET DEFAULT | NO ACTION
func (p *Parser) parseReferenceDefinition(start int) (*nodes.ColumnConstraint, error) {
	p.advance() // consume REFERENCES

	// Completion: after REFERENCES, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	constr := &nodes.ColumnConstraint{
		Loc:  nodes.Loc{Start: start},
		Type: nodes.ColConstrReferences,
	}

	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	constr.RefTable = ref

	// Optional column list
	if p.cur.Type == '(' {
		cols, err := p.parseParenIdentListWithCompletion("columnref")
		if err != nil {
			return nil, err
		}
		constr.RefColumns = cols
	}

	// MATCH
	if _, ok := p.match(kwMATCH); ok {
		m, err := p.parseMatchType()
		if err != nil {
			return nil, err
		}
		constr.Match = m
	}

	// ON DELETE / ON UPDATE
	for p.cur.Type == kwON {
		p.advance()
		if _, ok := p.match(kwDELETE); ok {
			constr.OnDelete = p.parseReferenceAction()
		} else if _, ok := p.match(kwUPDATE); ok {
			constr.OnUpdate = p.parseReferenceAction()
		} else {
			break
		}
	}

	constr.Loc.End = p.prev.End
	return constr, nil
}

// parseReferenceAction parses a foreign key action.
//
//	RESTRICT | CASCADE | SET NULL | SET DEFAULT | NO ACTION
func (p *Parser) parseReferenceAction() nodes.ReferenceAction {
	switch p.cur.Type {
	case kwRESTRICT:
		p.advance()
		return nodes.RefActRestrict
	case kwCASCADE:
		p.advance()
		return nodes.RefActCascade
	case kwSET:
		p.advance()
		if _, ok := p.match(kwNULL); ok {
			return nodes.RefActSetNull
		}
		if _, ok := p.match(kwDEFAULT); ok {
			return nodes.RefActSetDefault
		}
		return nodes.RefActNone
	case kwNO:
		p.advance()
		if p.cur.Type == kwACTION {
			p.advance()
		}
		return nodes.RefActNoAction
	}
	return nodes.RefActNone
}

// matchUnquotedKeyword consumes the current token when it is an UNQUOTED
// identifier equal (case-insensitively) to kw. A backtick-quoted identifier
// (e.g. `PERSISTENT`) is a quoted identifier, not the keyword, and must not match
// — MariaDB rejects a quoted keyword in a keyword position (1064).
func (p *Parser) matchUnquotedKeyword(kw string) bool {
	if !p.atUnquotedKeyword(kw) {
		return false
	}
	p.advance()
	return true
}

// atUnquotedKeyword reports whether the current token is an UNQUOTED
// identifier equal (case-insensitively) to kw, without consuming it. The lexer
// gives a backtick-quoted identifier the same tokIDENT as an unquoted one, so
// the source text tells them apart.
func (p *Parser) atUnquotedKeyword(kw string) bool {
	if p.cur.Type != tokIDENT || !strings.EqualFold(p.cur.Str, kw) {
		return false
	}
	// p.cur.Loc is absolute (includes baseOffset); p.lexer.input is the current
	// split segment, so subtract baseOffset before indexing (cf. the offset
	// helpers in parser.go). A leading backtick means the keyword was quoted.
	off := p.cur.Loc - p.lexer.baseOffset
	return !(off >= 0 && off < len(p.lexer.input) && p.lexer.input[off] == '`')
}

// parseGeneratedColumn parses a GENERATED ALWAYS AS (expr) [VIRTUAL|STORED].
func (p *Parser) parseGeneratedColumn() (*nodes.GeneratedColumn, error) {
	start := p.pos()
	p.advance() // consume GENERATED

	// Optional ALWAYS
	if p.cur.Type == kwALWAYS {
		p.advance()
	}

	// AS
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// GENERATED ALWAYS AS ROW START | ROW END (system-versioning period column)
	if p.cur.Type == kwROW {
		rb, err := p.parseRowBound()
		if err != nil {
			return nil, err
		}
		return &nodes.GeneratedColumn{Loc: nodes.Loc{Start: start, End: p.prev.End}, RowBound: rb}, nil
	}

	// (expr)
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Completion: inside GENERATED ALWAYS AS (...), offer expression context.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, &ParseError{Message: "collecting"}
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	gen := &nodes.GeneratedColumn{
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
		Expr: expr,
	}

	// VIRTUAL or STORED
	if _, ok := p.match(kwVIRTUAL); ok {
		gen.Stored = false
	} else if _, ok := p.match(kwSTORED); ok {
		gen.Stored = true
	} else if p.matchUnquotedKeyword("PERSISTENT") {
		// PERSISTENT is MariaDB's non-reserved synonym for STORED.
		gen.Stored = true
	}

	gen.Loc.End = p.prev.End
	return gen, nil
}

// parseGeneratedColumnShorthand parses AS (expr) [VIRTUAL|STORED] — the
// shorthand form without GENERATED ALWAYS prefix.
func (p *Parser) parseGeneratedColumnShorthand() (*nodes.GeneratedColumn, error) {
	start := p.pos()
	p.advance() // consume AS

	// (expr)
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Completion: inside AS (...), offer expression context.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, &ParseError{Message: "collecting"}
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	gen := &nodes.GeneratedColumn{
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
		Expr: expr,
	}

	// VIRTUAL or STORED
	if _, ok := p.match(kwVIRTUAL); ok {
		gen.Stored = false
	} else if _, ok := p.match(kwSTORED); ok {
		gen.Stored = true
	} else if p.matchUnquotedKeyword("PERSISTENT") {
		// PERSISTENT is MariaDB's non-reserved synonym for STORED.
		gen.Stored = true
	}

	gen.Loc.End = p.prev.End
	return gen, nil
}

// parseTableConstraint parses a table-level constraint.
//
//	[CONSTRAINT [symbol]] PRIMARY KEY (col,...) [index_option]
//	[CONSTRAINT [symbol]] UNIQUE [INDEX|KEY] [name] (col,...) [index_option]
//	[CONSTRAINT [symbol]] FOREIGN KEY [name] (col,...) reference_definition
//	[CONSTRAINT [symbol]] CHECK (expr)
//	{INDEX | KEY} [name] (col,...) [index_option]
//	{FULLTEXT | SPATIAL} [INDEX|KEY] [name] (col,...) [index_option]
func (p *Parser) parseTableConstraint() (*nodes.Constraint, error) {
	start := p.pos()
	constr := &nodes.Constraint{Loc: nodes.Loc{Start: start}}

	// Optional CONSTRAINT [name]
	if _, ok := p.match(kwCONSTRAINT); ok {
		if p.isIdentToken() && p.cur.Type != kwPRIMARY && p.cur.Type != kwUNIQUE &&
			p.cur.Type != kwFOREIGN && p.cur.Type != kwCHECK {
			constr.Name, _, _ = p.parseIdent()
		}
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance()
		p.match(kwKEY)
		constr.Type = nodes.ConstrPrimaryKey
		// Optional index name (undocumented but accepted by MySQL 8.0, silently ignored)
		if p.isIdentToken() && p.cur.Type != '(' {
			p.parseIdent() // consume and discard — MySQL ignores PK index names
		}
		// Optional index type
		if err := p.parseIndexTypeClause(constr); err != nil {
			return nil, err
		}
		idxCols, err := p.parseParenIndexKeyParts(constrAllowsOverlaps(constr.Type))
		if err != nil {
			return nil, err
		}
		constr.IndexColumns = idxCols
		constr.Columns = indexColumnsToNames(idxCols)
		if err := p.parseConstraintIndexOptions(constr); err != nil {
			return nil, err
		}

	case kwUNIQUE:
		p.advance()
		p.match(kwINDEX, kwKEY) // optional INDEX or KEY
		constr.Type = nodes.ConstrUnique
		// Optional index name
		if p.isIdentToken() && p.cur.Type != '(' {
			constr.Name, _, _ = p.parseIdent()
		}
		// Optional index type
		if err := p.parseIndexTypeClause(constr); err != nil {
			return nil, err
		}
		idxCols, err := p.parseParenIndexKeyParts(constrAllowsOverlaps(constr.Type))
		if err != nil {
			return nil, err
		}
		constr.IndexColumns = idxCols
		constr.Columns = indexColumnsToNames(idxCols)
		if err := p.parseConstraintIndexOptions(constr); err != nil {
			return nil, err
		}

	case kwFOREIGN:
		p.advance()
		p.match(kwKEY)
		constr.Type = nodes.ConstrForeignKey
		// Optional index name, kept apart from the CONSTRAINT symbol, which
		// names the foreign key when given (see nodes.Constraint).
		if p.isIdentToken() && p.cur.Type != '(' {
			name, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			constr.IndexName = name
		}
		// The column list is a key-part list, not a plain name list: a prefix
		// length and an ordering direction are both allowed. MariaDB 11.8.8 runs
		// `ADD FOREIGN KEY (ident, ident(1) ASC) REFERENCES t (a)`.
		//
		// A functional key part is not allowed, though, because a foreign
		// key's list is key_list and not the key_list_with_expression an
		// index takes: MariaDB 11.8.8 answers 1064 for
		// `ADD FOREIGN KEY ((a + 1)) REFERENCES u (b)` and for the same list
		// on CREATE TABLE, where it runs `ADD INDEX ((a + 1))`. (TiDB v8.5.5
		// does parse it, through its own IndexPartSpecificationList, and then
		// crashes building the key -- 1105, a nil dereference -- so the fork
		// there keeps the wider list and this one does not.)
		//
		// Ref: `FOREIGN KEY opt_index_name '(' key_list ')' references` in
		// MySQL 8.0.46 sql/sql_yacc.yy.
		fkCols, err := p.parseParenIndexKeyParts(constrAllowsOverlaps(constr.Type))
		if err != nil {
			return nil, err
		}
		for _, kp := range fkCols {
			if kp.Functional {
				return nil, &ParseError{
					Message:  "a foreign key takes no functional key part",
					Position: kp.Loc.Start,
				}
			}
		}
		constr.IndexColumns = fkCols
		constr.Columns = indexColumnsToNames(fkCols)
		// REFERENCES
		if p.cur.Type == kwREFERENCES {
			p.advance()
			// Completion: after FOREIGN KEY (...) REFERENCES, offer table_ref.
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("table_ref")
				return nil, &ParseError{Message: "collecting"}
			}
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			constr.RefTable = ref
			if p.cur.Type == '(' {
				refCols, err := p.parseParenIdentListWithCompletion("columnref")
				if err != nil {
					return nil, err
				}
				constr.RefColumns = refCols
			}
			// MATCH
			if _, ok := p.match(kwMATCH); ok {
				m, err := p.parseMatchType()
				if err != nil {
					return nil, err
				}
				constr.Match = m
			}
			// ON DELETE / ON UPDATE
			for p.cur.Type == kwON {
				p.advance()
				if _, ok := p.match(kwDELETE); ok {
					constr.OnDelete = p.parseReferenceAction()
				} else if _, ok := p.match(kwUPDATE); ok {
					constr.OnUpdate = p.parseReferenceAction()
				} else {
					break
				}
			}
		}

	case kwCHECK:
		p.advance()
		constr.Type = nodes.ConstrCheck
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		constr.Expr = expr
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		// [[NOT] ENFORCED]
		if p.cur.Type == kwNOT {
			p.advance()
			if p.cur.Type == kwENFORCED {
				p.advance()
				constr.NotEnforced = true
			}
		} else if p.cur.Type == kwENFORCED {
			p.advance()
		}

	case kwINDEX, kwKEY:
		p.advance()
		constr.Type = nodes.ConstrIndex
		// Optional index name
		if p.isIdentToken() && p.cur.Type != '(' {
			constr.Name, _, _ = p.parseIdent()
		}
		if err := p.parseIndexTypeClause(constr); err != nil {
			return nil, err
		}
		idxCols, err := p.parseParenIndexKeyParts(constrAllowsOverlaps(constr.Type))
		if err != nil {
			return nil, err
		}
		constr.IndexColumns = idxCols
		constr.Columns = indexColumnsToNames(idxCols)
		if err := p.parseConstraintIndexOptions(constr); err != nil {
			return nil, err
		}

	case kwFULLTEXT:
		p.advance()
		p.match(kwINDEX, kwKEY)
		constr.Type = nodes.ConstrFulltextIndex
		if p.isIdentToken() && p.cur.Type != '(' {
			constr.Name, _, _ = p.parseIdent()
		}
		idxCols, err := p.parseParenIndexKeyParts(constrAllowsOverlaps(constr.Type))
		if err != nil {
			return nil, err
		}
		constr.IndexColumns = idxCols
		constr.Columns = indexColumnsToNames(idxCols)
		if err := p.parseConstraintIndexOptions(constr); err != nil {
			return nil, err
		}

	case kwSPATIAL:
		p.advance()
		p.match(kwINDEX, kwKEY)
		constr.Type = nodes.ConstrSpatialIndex
		if p.isIdentToken() && p.cur.Type != '(' {
			constr.Name, _, _ = p.parseIdent()
		}
		idxCols, err := p.parseParenIndexKeyParts(constrAllowsOverlaps(constr.Type))
		if err != nil {
			return nil, err
		}
		constr.IndexColumns = idxCols
		constr.Columns = indexColumnsToNames(idxCols)
		if err := p.parseConstraintIndexOptions(constr); err != nil {
			return nil, err
		}

	default:
		return nil, &ParseError{
			Message:  "expected constraint definition",
			Position: p.cur.Loc,
		}
	}

	constr.Loc.End = p.prev.End
	return constr, nil
}

// parseConstraintIndexOptions parses index_option list after a table constraint's key parts.
// An option that starts and then fails, such as a COMMENT without its string,
// is a syntax error, as it is after CREATE INDEX.
func (p *Parser) parseConstraintIndexOptions(constr *nodes.Constraint) error {
	for {
		opt, ok, err := p.parseIndexOption()
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		constr.IndexOptions = append(constr.IndexOptions, opt)
	}
}

// parseIndexTypeClause parses the optional index_type clause a key definition
// takes in front of its column list.
//
// MySQL reads that position as opt_index_name_and_type, where the USING
// spelling may stand alone and the TYPE spelling may not: TYPE is an ordinary
// identifier, so `KEY TYPE BTREE (a)` names the index `type` and then fails on
// BTREE, which is the 1064 mysql:8.0 answers. The caller takes the name before
// this runs and takes a bare TYPE with it, so cur is at TYPE here only once a
// name has been read, and nothing has to be passed in to tell the two apart.
//
// Ref: opt_index_name_and_type in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseIndexTypeClause(constr *nodes.Constraint) error {
	if !p.isIndexTypeIntro() {
		return nil
	}
	p.advance()
	typeName, err := p.parseIndexTypeName()
	if err != nil {
		return err
	}
	constr.IndexType = typeName
	return nil
}

// parseTableOption parses a single table option.
//
//	ENGINE [=] engine_name
//	AUTO_INCREMENT [=] value
//	DEFAULT? CHARSET [=] charset | CHARACTER SET [=] charset
//	COLLATE [=] collation
//	COMMENT [=] 'string'
//	ROW_FORMAT [=] format
//	KEY_BLOCK_SIZE [=] value
//	And many more...
func (p *Parser) parseTableOption() (*nodes.TableOption, bool, error) {
	// Completion: offer table option keywords.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwENGINE, kwDEFAULT, kwCHARSET, kwCOLLATE, kwCOMMENT, kwAUTO_INCREMENT, kwROW_FORMAT, kwPARTITION} {
			p.addTokenCandidate(t)
		}
		return nil, false, &ParseError{Message: "collecting"}
	}

	start := p.pos()

	switch p.cur.Type {
	case kwENGINE:
		p.advance()
		p.matchEqual()
		// Completion: after ENGINE=, offer engine candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("engine")
			return nil, false, &ParseError{Message: "collecting"}
		}
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "ENGINE", Value: val}, true, nil

	case kwWITH:
		// WITH SYSTEM VERSIONING — system-versioned table option. VERSIONING is
		// non-reserved, so it is matched by identifier text.
		if p.peekNext().Type == kwSYSTEM {
			p.advance() // WITH
			p.advance() // SYSTEM
			if p.cur.Type != tokIDENT || !strings.EqualFold(p.cur.Str, "VERSIONING") {
				return nil, false, p.syntaxErrorAtCur()
			}
			p.advance() // VERSIONING
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "SYSTEM VERSIONING"}, true, nil
		}
		return nil, false, nil

	case kwAUTO_INCREMENT:
		p.advance()
		p.matchEqual()
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "AUTO_INCREMENT", Value: val}, true, nil

	case kwDEFAULT:
		// DEFAULT CHARSET or DEFAULT CHARACTER SET or DEFAULT COLLATE
		next := p.peekNext()
		if next.Type == kwCHARSET || next.Type == kwCHARACTER || next.Type == kwCOLLATE {
			p.advance() // consume DEFAULT
			return p.parseTableOption()
		}
		return nil, false, nil

	case kwCHARSET:
		p.advance()
		p.matchEqual()
		val, err := p.parseTableCharsetName()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "CHARSET", Value: val}, true, nil

	case kwCHARACTER:
		p.advance()
		if _, ok := p.match(kwSET); ok {
			p.matchEqual()
			val, err := p.parseTableCharsetName()
			if err != nil {
				return nil, false, err
			}
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "CHARACTER SET", Value: val}, true, nil
		}
		// CHARACTER without SET is no option. Returning "not an option" after
		// consuming it dropped the word, so `CREATE TABLE t (a INT) CHARACTER`
		// and `ALTER TABLE t COMMENT 'x' CHARACTER` parsed; the engine answers
		// 1064 for both.
		return nil, false, p.syntaxErrorAtCur()

	case kwCOLLATE:
		p.advance()
		p.matchEqual()
		val, err := p.parseTableCollationName()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "COLLATE", Value: val}, true, nil

	case kwCOMMENT:
		p.advance()
		p.matchEqual()
		// The value is a quoted string and is required: the server answers 1064
		// for a bare COMMENT, CONNECTION or PASSWORD.
		if p.cur.Type != tokSCONST {
			return nil, false, p.syntaxErrorAtCur()
		}
		val := p.cur.Str
		p.advance()
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "COMMENT", Value: val}, true, nil

	case kwROW_FORMAT:
		p.advance()
		p.matchEqual()
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "ROW_FORMAT", Value: val}, true, nil

	case kwCONNECTION:
		p.advance()
		p.matchEqual()
		if p.cur.Type != tokSCONST {
			return nil, false, p.syntaxErrorAtCur()
		}
		val := p.cur.Str
		p.advance()
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "CONNECTION", Value: val}, true, nil

	case kwPASSWORD:
		p.advance()
		p.matchEqual()
		if p.cur.Type != tokSCONST {
			return nil, false, p.syntaxErrorAtCur()
		}
		val := p.cur.Str
		p.advance()
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "PASSWORD", Value: val}, true, nil

	case kwSTORAGE:
		// STORAGE {DISK | MEMORY} has no equals sign: STORAGE = DISK is a 1064
		// (measured).
		p.advance()
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "STORAGE", Value: val}, true, nil

	case kwSTART:
		// START TRANSACTION (used for versioned tables)
		if next := p.peekNext(); next.Type == kwTRANSACTION {
			p.advance() // consume START
			p.advance() // consume TRANSACTION
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "START TRANSACTION", Value: ""}, true, nil
		}
		return nil, false, nil

	// CHECKSUM [=] n, which TABLE_CHECKSUM is the other name for: MySQL's
	// sql_yacc.yy gives TABLE_CHECKSUM_SYM the same `opt_equal ulong_num` as
	// CHECKSUM_SYM. The value is a number and only a number: MariaDB answers
	// 1064 for `CHECKSUM = 'abc'`, for `= DEFAULT` and for `= -1`. A decimal is
	// part of `ulong_num`, so `= 1.5` is accepted.
	case kwCHECKSUM, kwTABLE_CHECKSUM:
		optName := "CHECKSUM"
		if p.cur.Type == kwTABLE_CHECKSUM {
			optName = "TABLE_CHECKSUM"
		}
		p.advance()
		p.matchEqual()
		val, err := p.consumeNumericOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil

	case kwTABLESPACE:
		// TABLESPACE tablespace_name [STORAGE {DISK | MEMORY}]
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		opt := &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "TABLESPACE", Value: val}
		if p.cur.Type == kwSTORAGE {
			p.advance()
			storageVal, err := p.consumeOptionValue()
			if err != nil {
				return nil, false, err
			}
			opt.Storage = storageVal
		}
		opt.Loc.End = p.prev.End
		return opt, true, nil

	case kwENCRYPTION:
		// ENCRYPTION [=] {'Y' | 'N'}
		p.advance()
		p.matchEqual()
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "ENCRYPTION", Value: val}, true, nil

	case kwUNION:
		// UNION [=] '(' [tbl_name [, tbl_name]...] ')'
		//
		// The parentheses are required and the list inside them may be empty,
		// and each name may be qualified: MariaDB 11.8.8 runs `UNION = ()`,
		// `UNION ()` and `UNION = (d.t1, d.t2)`, and answers 1064 for a bare
		// `UNION`, for `UNION = (t1,)` and for `UNION = (,)`.
		//
		// Ref: `UNION opt_equal '(' opt_table_list ')'` in MySQL 8.0.46
		// sql/sql_yacc.yy.
		p.advance()
		p.matchEqual()
		names, err := p.parseParenTableNameListOpt()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "UNION", Value: strings.Join(names, ",")}, true, nil

	case kwINDEX:
		// INDEX DIRECTORY [=] 'path' (as table option)
		if next := p.peekNext(); next.Type == kwDIRECTORY {
			p.advance() // consume INDEX
			p.advance() // consume DIRECTORY
			p.matchEqual()
			val, err := p.consumeOptionValue()
			if err != nil {
				return nil, false, err
			}
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "INDEX DIRECTORY", Value: val}, true, nil
		}
		return nil, false, nil

	case kwDATA:
		// DATA DIRECTORY [=] 'path'
		if next := p.peekNext(); next.Type == kwDIRECTORY {
			p.advance() // consume DATA
			p.advance() // consume DIRECTORY
			p.matchEqual()
			val, err := p.consumeOptionValue()
			if err != nil {
				return nil, false, err
			}
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "DATA DIRECTORY", Value: val}, true, nil
		}
		return nil, false, nil
	}

	// Handle keyword-token-based table options: COMPRESSION, INSERT_METHOD.
	if p.cur.Type == kwCOMPRESSION || p.cur.Type == kwINSERT_METHOD {
		optName := p.cur.Str
		p.advance()
		p.matchEqual()
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil
	}

	// Handle identifier-based options: KEY_BLOCK_SIZE, MAX_ROWS, etc.
	switch p.cur.Type {
	// STATS_AUTO_RECALC [=] {DEFAULT | 0 | 1}, and nothing else: MariaDB
	// answers 1064 for `= 2`, `= 10`, `= -1`, `= 'abc'` and `= '1'`
	// (https://dev.mysql.com/doc/refman/8.0/en/create-table.html).
	case kwSTATS_AUTO_RECALC:
		optName := p.cur.Str
		p.advance()
		p.matchEqual()
		val, err := p.consumeZeroOneOrDefaultOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil

	// ENGINE_ATTRIBUTE and SECONDARY_ENGINE_ATTRIBUTE [=] 'json'. The string is
	// required: MariaDB answers 1064 for the bare name and for an `=` with
	// nothing after it, the way it does for a bare COMMENT
	// (https://dev.mysql.com/doc/refman/8.0/en/create-table.html).
	case kwENGINE_ATTRIBUTE, kwSECONDARY_ENGINE_ATTRIBUTE:
		optName := p.cur.Str
		p.advance()
		p.matchEqual()
		val, err := p.consumeStringOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil

	case kwKEY_BLOCK_SIZE, kwSTATS_PERSISTENT,
		kwSTATS_SAMPLE_PAGES, kwMAX_ROWS, kwMIN_ROWS, kwAVG_ROW_LENGTH,
		kwPACK_KEYS, kwDELAY_KEY_WRITE, kwSECONDARY_ENGINE,
		kwAUTOEXTEND_SIZE:
		optName := p.cur.Str
		p.advance()
		p.matchEqual()
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil
	}

	// Generic MariaDB table option: `IDENT = value`. MariaDB's parser accepts an
	// unknown option name and rejects it semantically (ER_UNKNOWN_OPTION 1911),
	// not as a syntax error — so a parse-only check accepts the shape. This covers
	// PAGE_COMPRESSED, PAGE_COMPRESSION_LEVEL, TRANSACTIONAL, PAGE_CHECKSUM, and
	// any future plain option. The '=' is required (a bare option name is 1064),
	// which also guards against consuming a trailing clause.
	if next := p.peekNext().Type; p.cur.Type == tokIDENT && (next == '=' || next == tokAssign) {
		optName := p.cur.Str
		p.advance()                        // option name
		p.advance()                        // '='
		val, err := p.consumeOptionValue() // errors if no value token follows
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil
	}

	return nil, false, nil
}

// parseTableCharsetName parses the name of a table's CHARACTER SET option. It is what
// a database option takes, or DEFAULT: charset_name_or_default, which MySQL's table
// option does not have (measured).
func (p *Parser) parseTableCharsetName() (string, error) {
	if p.cur.Type == kwDEFAULT {
		return p.advance().Str, nil
	}
	return p.parseCharsetName()
}

// parseTableCollationName parses the name of a table's COLLATE option, which is
// collation_name_or_default.
func (p *Parser) parseTableCollationName() (string, error) {
	return p.parseTableCharsetName()
}

// consumeOptionValue reads the next token as an option value string.
func (p *Parser) consumeOptionValue() (string, error) {
	switch p.cur.Type {
	case tokSCONST:
		val := p.cur.Str
		p.advance()
		return val, nil
	case tokICONST:
		val := p.cur.Str
		if val == "" {
			// Ival is set but Str may be empty for integer tokens
			val = string(rune('0' + p.cur.Ival))
		}
		p.advance()
		return val, nil
	default:
		// A reserved word is no value, except DEFAULT: engine_defined_option takes a
		// string, an ident, a number or DEFAULT, so `foo = null` and `foo = select`
		// are each a 1064 (measured).
		if p.cur.Type == tokIDENT || (p.cur.Type >= 700 && (!isReserved(p.cur.Type) || p.cur.Type == kwDEFAULT)) {
			name, _, err := p.parseKeywordOrIdent()
			if err != nil {
				return "", err
			}
			return name, nil
		}
		// No value token (e.g. `ENGINE=` at end of statement). MariaDB rejects an
		// option with no value (1064); do not silently return an empty value.
		return "", p.syntaxErrorAtCur()
	}
}

// consumeIntegerOptionValue reads an option value that must be an integer.
// A string, a decimal, a keyword and a leading sign are each a 1064.
func (p *Parser) consumeIntegerOptionValue() (string, error) {
	if p.cur.Type != tokICONST {
		return "", p.syntaxErrorAtCur()
	}
	val := p.cur.Str
	p.advance()
	return val, nil
}

// consumeNumericOptionValue reads an option value that must be a number,
// integer or decimal. A string, a keyword and a leading sign are each a 1064.
func (p *Parser) consumeNumericOptionValue() (string, error) {
	if p.cur.Type != tokICONST && p.cur.Type != tokFCONST {
		return "", p.syntaxErrorAtCur()
	}
	val := p.cur.Str
	p.advance()
	return val, nil
}

// consumeStringOptionValue reads an option value that must be a quoted string,
// the way COMMENT's is. An unquoted word, a number and a missing value are each
// a 1064.
func (p *Parser) consumeStringOptionValue() (string, error) {
	if p.cur.Type != tokSCONST {
		return "", p.syntaxErrorAtCur()
	}
	val := p.cur.Str
	p.advance()
	return val, nil
}

// consumeZeroOneOrDefaultOptionValue reads an option value that must be 0, 1 or
// DEFAULT, which is what STATS_AUTO_RECALC takes. Any other number, and any
// string, is a 1064.
func (p *Parser) consumeZeroOneOrDefaultOptionValue() (string, error) {
	if p.cur.Type == kwDEFAULT {
		val := p.cur.Str
		p.advance()
		return val, nil
	}
	if p.cur.Type != tokICONST || (p.cur.Ival != 0 && p.cur.Ival != 1) {
		return "", p.syntaxErrorAtCur()
	}
	val := p.cur.Str
	p.advance()
	return val, nil
}

// parsePartitionClause parses a PARTITION BY clause.
//
//	PARTITION BY
//	    { [LINEAR] HASH(expr) | [LINEAR] KEY [ALGORITHM={1|2}] (column_list)
//	    | RANGE(expr) | RANGE COLUMNS(column_list)
//	    | LIST(expr) | LIST COLUMNS(column_list) }
//	[PARTITIONS num]
//	[(partition_definition [, partition_definition] ...)]
func (p *Parser) parsePartitionClause() (*nodes.PartitionClause, error) {
	start := p.pos()
	p.advance() // consume PARTITION

	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}

	part := &nodes.PartitionClause{Loc: nodes.Loc{Start: start}}

	// Optional LINEAR
	if p.cur.Type == kwLINEAR {
		p.advance()
		part.Linear = true
	}

	switch p.cur.Type {
	case kwHASH:
		p.advance()
		part.Type = nodes.PartitionHash
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		expr, err := p.parseBitExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		part.Expr = expr
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

	case kwKEY:
		p.advance()
		part.Type = nodes.PartitionKey
		// Optional ALGORITHM={1|2}
		if p.cur.Type == kwALGORITHM {
			p.advance()
			p.match('=')
			if p.cur.Type == tokICONST {
				part.Algorithm = int(p.cur.Ival)
				p.advance()
			}
		}
		if p.cur.Type == '(' {
			// KEY () with empty column list means "use PK columns".
			if p.peekNext().Type == ')' {
				p.advance() // consume '('
				p.advance() // consume ')'
			} else {
				cols, err := p.parseParenIdentList()
				if err != nil {
					return nil, err
				}
				part.Columns = cols
			}
		}

	case kwRANGE:
		p.advance()
		part.Type = nodes.PartitionRange
		if _, ok := p.match(kwCOLUMNS, kwFIELDS); ok {
			cols, err := p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
			part.Columns = cols
		} else {
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			expr, err := p.parseBitExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.syntaxErrorAtCur()
			}
			part.Expr = expr
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}

	case kwLIST:
		p.advance()
		part.Type = nodes.PartitionList
		if _, ok := p.match(kwCOLUMNS, kwFIELDS); ok {
			cols, err := p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
			part.Columns = cols
		} else {
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			// MariaDB reads VALUES here as the VALUES(col) function
			// (IN_PART_FUNC, MariaDB 11.8.8 sql/sql_yacc.yy line 4964),
			// which then fails as a partition function (error 1564).
			p.valuesIsFunction = true
			expr, err := p.parseBitExpr()
			p.valuesIsFunction = false
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.syntaxErrorAtCur()
			}
			part.Expr = expr
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}

	default:
		// SYSTEM_TIME is non-reserved (matched by identifier text).
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "SYSTEM_TIME") {
			p.advance() // SYSTEM_TIME
			part.Type = nodes.PartitionSystemTime
			if err := p.parseSystemTimePartitionSpec(part); err != nil {
				return nil, err
			}
		} else {
			return nil, &ParseError{
				Message:  "expected HASH, KEY, RANGE, LIST, or SYSTEM_TIME after PARTITION BY",
				Position: p.cur.Loc,
			}
		}
	}

	// PARTITIONS num (the count is required)
	if p.cur.Type == kwPARTITIONS {
		p.advance()
		if p.cur.Type != tokICONST {
			return nil, p.syntaxErrorAtCur()
		}
		part.NumParts = int(p.cur.Ival)
		p.advance()
	}

	// SUBPARTITION BY {HASH(expr) | KEY [ALGORITHM=N] (column_list)}
	if p.cur.Type == kwSUBPARTITION {
		p.advance() // SUBPARTITION
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		switch {
		case p.cur.Type == kwHASH:
			p.advance()
			part.SubPartType = nodes.PartitionHash
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			expr, err := p.parseBitExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.syntaxErrorAtCur()
			}
			part.SubPartExpr = expr
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		case p.cur.Type == kwKEY:
			p.advance()
			part.SubPartType = nodes.PartitionKey
			// Optional ALGORITHM = N
			if p.cur.Type == kwALGORITHM {
				p.advance()
				p.match('=')
				if p.cur.Type == tokICONST {
					part.SubPartAlgo = int(p.cur.Ival)
					p.advance()
				}
			}
			// KEY () with empty column list means "use PK columns".
			if p.cur.Type == '(' && p.peekNext().Type == ')' {
				p.advance() // consume '('
				p.advance() // consume ')'
			} else {
				cols, err := p.parseParenIdentList()
				if err != nil {
					return nil, err
				}
				part.SubPartColumns = cols
			}
		case p.cur.Type == kwLINEAR:
			p.advance() // LINEAR
			if p.cur.Type == kwHASH {
				p.advance()
				part.SubPartType = nodes.PartitionHash
				if _, err := p.expect('('); err != nil {
					return nil, err
				}
				expr, err := p.parseBitExpr()
				if err != nil {
					return nil, err
				}
				if expr == nil {
					return nil, p.syntaxErrorAtCur()
				}
				part.SubPartExpr = expr
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
			} else if p.cur.Type == kwKEY {
				p.advance()
				part.SubPartType = nodes.PartitionKey
				if p.cur.Type == kwALGORITHM {
					p.advance()
					p.match('=')
					if p.cur.Type == tokICONST {
						part.SubPartAlgo = int(p.cur.Ival)
						p.advance()
					}
				}
				cols, err := p.parseParenIdentList()
				if err != nil {
					return nil, err
				}
				part.SubPartColumns = cols
			}
		}

		// SUBPARTITIONS num
		if p.cur.Type == kwSUBPARTITIONS {
			p.advance()
			if p.cur.Type == tokICONST {
				part.NumSubParts = int(p.cur.Ival)
				p.advance()
			}
		}
	}

	// Optional partition definitions
	if p.cur.Type == '(' {
		p.advance()
		for {
			pd, err := p.parsePartitionDef()
			if err != nil {
				return nil, err
			}
			// A SYSTEM_TIME HISTORY/CURRENT partition cannot also carry VALUES
			// (MariaDB 1064). VALUES without HISTORY/CURRENT is left to the catalog
			// (1480); HISTORY/CURRENT under another type is left to the catalog (4113).
			if part.Type == nodes.PartitionSystemTime && pd.SystemTime != "" && pd.Values != nil {
				return nil, &ParseError{
					Position: pd.Loc.Start,
					Message:  "VALUES is not allowed for a SYSTEM_TIME HISTORY/CURRENT partition",
				}
			}
			part.Partitions = append(part.Partitions, pd)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		// An empty partition list "()" is a syntax error.
		if len(part.Partitions) == 0 {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// A PARTITIONS count and an explicit partition list must agree.
	if part.NumParts > 0 && len(part.Partitions) > 0 && len(part.Partitions) != part.NumParts {
		return nil, &ParseError{
			Position: part.Loc.Start,
			Message:  "Wrong number of partitions defined, mismatch with previous setting",
		}
	}

	part.Loc.End = p.prev.End
	return part, nil
}

// parseSystemTimePartitionSpec parses the SYSTEM_TIME specifics that follow
// PARTITION BY SYSTEM_TIME, before the PARTITIONS count and partition list:
//
// Each arm encodes the exact MariaDB grammar (STARTS only after INTERVAL; AUTO
// not valid bare); other shapes fall through unconsumed and fail as 1064:
//
//	INTERVAL value unit [STARTS expr] [AUTO]
//	LIMIT n [AUTO]
//	(no modifier)
func (p *Parser) parseSystemTimePartitionSpec(part *nodes.PartitionClause) error {
	switch p.cur.Type {
	case kwINTERVAL:
		p.advance() // INTERVAL
		val, err := p.parseExpr()
		if err != nil {
			return err
		}
		part.IntervalValue = val
		// The unit is required and must be a valid interval unit (mirrors
		// parseIntervalExpr; compound units like DAY_HOUR are reserved keywords).
		unit, _, err := p.parseKeywordOrIdent()
		if err != nil {
			return err
		}
		upper := strings.ToUpper(unit)
		if !isValidIntervalUnit(upper) {
			return &ParseError{Position: p.pos(), Message: "invalid INTERVAL unit: " + unit}
		}
		part.IntervalUnit = upper
		// STARTS is only valid after INTERVAL.
		if p.cur.Type == kwSTARTS {
			p.advance() // STARTS
			expr, err := p.parseExpr()
			if err != nil {
				return err
			}
			part.Starts = expr
		}
		if p.cur.Type == kwAUTO {
			p.advance() // AUTO
			part.Auto = true
		}
	case kwLIMIT:
		p.advance() // LIMIT
		if p.cur.Type != tokICONST {
			return p.syntaxErrorAtCur()
		}
		part.Limit = int(p.cur.Ival)
		p.advance()
		if p.cur.Type == kwAUTO {
			p.advance() // AUTO
			part.Auto = true
		}
	}
	return nil
}

// parsePartitionDef parses a single partition definition.
//
//	PARTITION partition_name
//	    [VALUES {LESS THAN (expr|MAXVALUE) | IN (value_list)} | {HISTORY | CURRENT}]
//	    [table_options]
func (p *Parser) parsePartitionDef() (*nodes.PartitionDef, error) {
	start := p.pos()
	if _, err := p.expect(kwPARTITION); err != nil {
		return nil, err
	}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	pd := &nodes.PartitionDef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// SYSTEM_TIME partitions are tagged HISTORY or CURRENT instead of VALUES.
	if p.cur.Type == kwHISTORY {
		p.advance()
		pd.SystemTime = "HISTORY"
	} else if p.cur.Type == kwCURRENT {
		p.advance()
		pd.SystemTime = "CURRENT"
	}

	// VALUES LESS THAN or VALUES IN
	if _, ok := p.match(kwVALUES); ok {
		if p.cur.Type == kwLESS {
			p.advance() // LESS
			if _, err := p.expect(kwTHAN); err != nil {
				return nil, err
			}
			if p.cur.Type == '(' {
				p.advance()
				// Parse comma-separated list of values (may include MAXVALUE tokens and expressions).
				var items []nodes.Node
				for {
					if p.cur.Type == kwMAXVALUE {
						items = append(items, &nodes.String{Str: "MAXVALUE"})
						p.advance()
					} else {
						expr, err := p.parseBitExpr()
						if err != nil {
							return nil, err
						}
						if expr == nil {
							return nil, p.syntaxErrorAtCur()
						}
						items = append(items, expr)
					}
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				// A single bound is stored as itself, MAXVALUE as a String; a
				// tuple keeps every column, MAXVALUE ones included, so
				// "(MAXVALUE, MAXVALUE)" reads back with both.
				if len(items) == 1 {
					pd.Values = items[0]
				} else {
					pd.Values = &nodes.List{Items: items}
				}
			} else if p.cur.Type == kwMAXVALUE {
				pd.Values = &nodes.String{Str: "MAXVALUE"}
				p.advance()
			} else {
				return nil, p.syntaxErrorAtCur()
			}
		} else if _, ok := p.match(kwIN); ok {
			if p.cur.Type == '(' {
				p.advance()
				var items []nodes.Node
				for {
					if p.cur.Type == '(' {
						// Tuple: (val1, val2, ...) for multi-column LIST COLUMNS
						p.advance()
						var tupleVals []nodes.Node
						for {
							expr, err := p.parseBitExpr()
							if err != nil {
								return nil, err
							}
							if expr == nil {
								return nil, p.syntaxErrorAtCur()
							}
							tupleVals = append(tupleVals, expr)
							if p.cur.Type != ',' {
								break
							}
							p.advance()
						}
						if _, err := p.expect(')'); err != nil {
							return nil, err
						}
						items = append(items, &nodes.List{Items: tupleVals})
					} else {
						expr, err := p.parseBitExpr()
						if err != nil {
							return nil, err
						}
						if expr == nil {
							return nil, p.syntaxErrorAtCur()
						}
						items = append(items, expr)
					}
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				pd.Values = &nodes.List{Items: items}
			} else {
				return nil, p.syntaxErrorAtCur()
			}
		}
	}

	// Optional partition-definition options
	for {
		opt, ok, err := p.parsePartDefOption()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		pd.Options = append(pd.Options, opt)
	}

	// Optional subpartition definitions: (SUBPARTITION name [table_options], ...)
	if p.cur.Type == '(' {
		p.advance()
		for {
			spd, err := p.parseSubPartitionDef()
			if err != nil {
				return nil, err
			}
			pd.SubPartitions = append(pd.SubPartitions, spd)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	pd.Loc.End = p.prev.End
	return pd, nil
}

// parseSubPartitionDef parses a single subpartition definition.
//
//	SUBPARTITION logical_name [table_options]
func (p *Parser) parseSubPartitionDef() (*nodes.SubPartitionDef, error) {
	start := p.pos()
	if _, err := p.expect(kwSUBPARTITION); err != nil {
		return nil, err
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	spd := &nodes.SubPartitionDef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}
	for {
		opt, ok, err := p.parsePartDefOption()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		spd.Options = append(spd.Options, opt)
	}
	spd.Loc.End = p.prev.End
	return spd, nil
}

// parseParenTableNameListOpt parses a parenthesized, possibly empty list of
// table names, each optionally qualified by a database name. It is the value
// of the UNION table option.
func (p *Parser) parseParenTableNameListOpt() ([]string, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var names []string
	if p.cur.Type == ')' {
		p.advance()
		return names, nil
	}
	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		if ref.Schema != "" {
			names = append(names, ref.Schema+"."+ref.Name)
		} else {
			names = append(names, ref.Name)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return names, nil
}

// parseMatchType reads the FULL, PARTIAL or SIMPLE a MATCH clause takes. The
// three are a keyword list, not an identifier: MariaDB 11.8.8 answers 1064 for
// `ADD FOREIGN KEY (a) REFERENCES u(b) MATCH foo`.
//
// Ref: opt_match_clause in MySQL 8.0.46 sql/sql_yacc.yy.
func (p *Parser) parseMatchType() (string, error) {
	if tok, ok := p.match(kwFULL, kwPARTIAL, kwSIMPLE); ok {
		return strings.ToUpper(tok.Str), nil
	}
	return "", p.syntaxErrorAtCur()
}

// parsePartDefOption parses one of the options a partition or subpartition
// definition takes, and reports whether one was there. The list is its own,
// narrower than a table's: a partition definition takes STORAGE ENGINE, which
// no table option spells, and takes none of ENGINE's neighbours such as
// CHARSET, AUTO_INCREMENT or ROW_FORMAT.
//
//	part_def_option:
//	    COMMENT [=] 'string'
//	  | ENGINE [=] name
//	  | STORAGE ENGINE [=] name
//	  | DATA DIRECTORY [=] 'string'
//	  | INDEX DIRECTORY [=] 'string'
//	  | MAX_ROWS [=] number
//	  | MIN_ROWS [=] number
//	  | TABLESPACE [=] name
//	  | NODEGROUP [=] number
//
// MariaDB 11.8.8 runs `ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE
// ENGINE text)`, `... STORAGE ENGINE = text`, `... MAX_ROWS 12`,
// `... NODEGROUP 2`, `... TABLESPACE ts`, `... DATA DIRECTORY '/x'` and
// `... (SUBPARTITION s MIN_ROWS 11)`, and answers 1064 for
// `... STORAGE DISK`, `... CHARSET utf8`, `... AUTO_INCREMENT = 1`,
// `... ROW_FORMAT = DYNAMIC` and `... INSERT_METHOD = FIRST`, all of which
// the table-option reader used here before took.
//
// Ref: opt_part_option in MySQL 8.0.46 sql/sql_yacc.yy.
func (p *Parser) parsePartDefOption() (*nodes.TableOption, bool, error) {
	start := p.pos()

	switch p.cur.Type {
	case kwCOMMENT:
		p.advance()
		p.matchEqual()
		val, err := p.consumeStringOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "COMMENT", Value: val}, true, nil

	case kwENGINE, kwTABLESPACE:
		name := strings.ToUpper(p.cur.Str)
		p.advance()
		p.matchEqual()
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: name, Value: val}, true, nil

	case kwSTORAGE:
		// STORAGE ENGINE [=] name. STORAGE on its own is no partition option:
		// the table option STORAGE {DISK | MEMORY} has no partition form.
		if p.peekNext().Type != kwENGINE {
			return nil, false, nil
		}
		p.advance() // consume STORAGE
		p.advance() // consume ENGINE
		p.matchEqual()
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "STORAGE ENGINE", Value: val}, true, nil

	case kwMAX_ROWS, kwMIN_ROWS, kwNODEGROUP:
		name := strings.ToUpper(p.cur.Str)
		p.advance()
		p.matchEqual()
		val, err := p.consumeIntegerOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: name, Value: val}, true, nil

	case kwDATA, kwINDEX:
		if p.peekNext().Type != kwDIRECTORY {
			return nil, false, nil
		}
		name := strings.ToUpper(p.cur.Str) + " DIRECTORY"
		p.advance() // consume DATA or INDEX
		p.advance() // consume DIRECTORY
		p.matchEqual()
		val, err := p.consumeStringOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: name, Value: val}, true, nil
	}

	return nil, false, nil
}
