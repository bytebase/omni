package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseCreateTableStmt parses a CREATE TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-table.html
//
//	CREATE [TEMPORARY] TABLE [IF NOT EXISTS] tbl_name
//	    (create_definition,...) [table_options] [partition_options]
//	CREATE [TEMPORARY] TABLE [IF NOT EXISTS] tbl_name
//	    LIKE old_tbl_name
//	CREATE [TEMPORARY] TABLE [IF NOT EXISTS] tbl_name
//	    [AS] query_expression
func (p *Parser) parseCreateTableStmt(start int, temporary bool) (*nodes.CreateTableStmt, error) {
	p.advance() // consume TABLE

	stmt := &nodes.CreateTableStmt{
		Loc:       nodes.Loc{Start: start},
		Temporary: temporary,
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// Completion: after CREATE TABLE, identifier context (no specific candidates).
	p.checkCursor()
	if p.collectMode() {
		// No specific candidates for table name — user defines a new name.
		return nil, &ParseError{Message: "collecting"}
	}

	// Table name
	ref, err := p.parseRelationFactor()
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
		likeRef, err := p.parseRelationFactor()
		if err != nil {
			return nil, err
		}
		stmt.Like = likeRef
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// CREATE TABLE ... AS SELECT
	if p.cur.Type == kwAS {
		p.advance()
		sel, err := p.parseQueryStatement()
		if err != nil {
			return nil, err
		}
		stmt.Select = sel
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Parenthesized create definitions (columns and constraints)
	hasDefinitions := false
	if p.cur.Type == '(' {
		next := p.peekNext()
		// Check for CREATE TABLE ... (query) — query source without AS: a
		// parenthesized query expression. Let parseSelectStmt consume the
		// parentheses and any following set-op / ORDER BY / LIMIT so
		// "(SELECT 1) UNION (SELECT 2)" stays one source.
		if p.isQueryExpressionStart() {
			sel, err := p.parseQueryStatement()
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
			likeRef, err := p.parseRelationFactor()
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

		p.advance() // consume '('
		if err := p.parseCreateDefinitions(stmt); err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		hasDefinitions = true
	}

	// CREATE TABLE t (cols) [IGNORE|REPLACE] [AS] SELECT ...
	if p.cur.Type == kwAS {
		p.advance()
		sel, err := p.parseQueryStatement()
		if err != nil {
			return nil, err
		}
		stmt.Select = sel
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Table options. A comma may stand between two of them and never ends the
	// list (table_option_list in OceanBase CE v4.3.5_CE_BP1
	// src/sql/parser/sql_parser_mysql_mode.y:7513), so a comma with no option
	// after it is a syntax error rather than the end of the options:
	// "ENGINE=InnoDB, (SELECT 1)" is 1064.
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

	// Partition clause
	if p.cur.Type == kwPARTITION {
		part, err := p.parsePartitionClause()
		if err != nil {
			return nil, err
		}
		stmt.Partitions = part
	}

	// with_column_group, after the options and the partitioning.
	if group, ok, err := p.parseWithColumnGroup(); err != nil {
		return nil, err
	} else if ok {
		stmt.ColumnGroup = group
	}

	// IGNORE / REPLACE before SELECT in CTAS
	if p.cur.Type == kwIGNORE {
		p.advance()
		stmt.Ignore = true
	} else if p.cur.Type == kwREPLACE {
		p.advance()
		stmt.Replace = true
	}

	// [AS] query after the columns, options and partitioning, with or without
	// AS, and required after IGNORE or REPLACE: create_table_stmt in OceanBase
	// CE v4.3.5_CE_BP1 src/sql/parser/sql_parser_mysql_mode.y:5365. Reading
	// only SELECT here used to end the statement in front of "VALUES ROW(1)"
	// or "(SELECT 1)", which then parsed as a second statement.
	switch {
	case p.cur.Type == kwAS:
		p.advance()
		sel, err := p.parseQueryStatement()
		if err != nil {
			return nil, err
		}
		stmt.Select = sel
	case p.isQueryExpressionStart():
		sel, err := p.parseQueryStatement()
		if err != nil {
			return nil, err
		}
		stmt.Select = sel
	case stmt.Ignore || stmt.Replace:
		return nil, p.syntaxErrorAtCur()
	case !hasDefinitions:
		// A table with no definitions is made from a query: every alternative of
		// create_table_stmt without '(' table_element_list ')' ends in a
		// select_stmt, so CREATE TABLE t NOPARALLEL and a bare CREATE TABLE t are
		// 1064s.
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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

		// Table-level constraints: PRIMARY KEY, UNIQUE, INDEX, KEY, FULLTEXT, SPATIAL, FOREIGN KEY, CONSTRAINT, CHECK
		if p.isTableConstraintStart() {
			constr, err := p.parseTableConstraint(false)
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
	case kwVECTOR:
		// OceanBase's VECTOR {INDEX | KEY}. VECTOR is also a column's name, so the word
		// only starts a constraint when the keyword follows it: `vector VECTOR(3)` and
		// `vector INT KEY` are columns.
		next := p.peekNext().Type
		return next == kwINDEX || next == kwKEY
	}
	return false
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

	// Data type, which no column does without: every engine answers 1064 (MariaDB 4161)
	// for `a key`, `a comment 'x'` and `a generated always as (1)`, generated or not.
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	col.TypeName = dt

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

	col.Loc.End = p.prev.End
	return col, nil
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
		exprStart := p.pos()
		expr, err := p.parseDefaultValue()
		if err != nil {
			return false, err
		}
		if err := refuseColumnValue(expr, exprStart); err != nil {
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

	// MySQL also writes NOT NULL AUTO_INCREMENT UNIQUE as the column option
	// SERIAL DEFAULT VALUE. OceanBase CE 4.3.5.6 does not have it — it answers
	// 1064 for `a int serial default value` in every spelling, while taking
	// SERIAL as a type — so it is not read here.

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
		cc, err := p.parseColumnCheckConstraint(start)
		if err != nil {
			return false, err
		}
		col.Constraints = append(col.Constraints, cc)
		return true, nil

	case kwCONSTRAINT:
		// CONSTRAINT [symbol] CHECK (expr) [[NOT] ENFORCED], the column-level
		// check the CONSTRAINT keyword may name. The symbol is optional, and
		// CHECK is the only continuation: OceanBase CE 4.3.5.6 runs
		// `ADD a INT CONSTRAINT c CHECK (a > 1) NOT ENFORCED` and
		// `ADD a INT CONSTRAINT CHECK (a > 1)`, and answers 1064 for
		// `ADD a INT CONSTRAINT c NOT NULL` and for a bare
		// `ADD a INT CONSTRAINT c`.
		//
		// Ref: `opt_constraint_name CHECK '(' expr ')'
		// opt_constraint_enforcement` among the column attributes in MySQL
		// 8.0.46 sql/sql_yacc.yy.
		p.advance() // consume CONSTRAINT
		var symbol string
		if p.cur.Type != kwCHECK && p.isIdentToken() {
			var err error
			symbol, _, err = p.parseIdent()
			if err != nil {
				return false, err
			}
		}
		if p.cur.Type != kwCHECK {
			return false, p.syntaxErrorAtCur()
		}
		cc, err := p.parseColumnCheckConstraint(start)
		if err != nil {
			return false, err
		}
		cc.Name = symbol
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
		return false, nil

	case kwON:
		// ON UPDATE CURRENT_TIMESTAMP
		if next := p.peekNext(); next.Type == kwUPDATE {
			p.advance() // consume ON
			p.advance() // consume UPDATE
			exprStart := p.pos()
			expr, err := p.parseExpr()
			if err != nil {
				return false, err
			}
			if expr == nil {
				return false, p.syntaxErrorAtCur()
			}
			if err := refuseColumnValue(expr, exprStart); err != nil {
				return false, err
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

	// ENGINE_ATTRIBUTE and SECONDARY_ENGINE_ATTRIBUTE are not read here: OceanBase CE 4.3.5.6
	// answers 1064 for both as a column option, as it does for both as a table
	// option, so taking them only passed statements the server refuses.

	// SRID column option (now lexes as kwSRID keyword token). Record presence even for
	// `SRID 0`: MySQL keeps `SRID 0` as a distinct, present clause (SHOW CREATE echoes it),
	// so HasSRID — not SRID != 0 — is what downstream emit/diff must gate on.
	if p.cur.Type == kwSRID {
		p.advance()
		if p.cur.Type != tokICONST {
			return false, nil
		}
		col.TypeName.SRID = int(p.cur.Ival)
		col.TypeName.HasSRID = true
		p.advance()
		return true, nil
	}

	return false, nil
}

// refuseColumnValue refuses a bare name as a column's DEFAULT or ON UPDATE value, as the
// grammar's literal and now_or_signed_literal do, so `a TIMESTAMP DEFAULT now` is a 1064.
func refuseColumnValue(expr nodes.ExprNode, start int) error {
	if _, ok := expr.(*nodes.ColumnRef); ok {
		return &ParseError{Message: "syntax error", Position: start}
	}
	return nil
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

	ref, err := p.parseRelationFactor()
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
			act, err := p.parseReferenceAction()
			if err != nil {
				return nil, err
			}
			constr.OnDelete = act
		} else if _, ok := p.match(kwUPDATE); ok {
			act, err := p.parseReferenceAction()
			if err != nil {
				return nil, err
			}
			constr.OnUpdate = act
		} else {
			break
		}
	}

	constr.Loc.End = p.prev.End
	return constr, nil
}

// parseReferenceAction parses a foreign key reference action. The action is
// required after ON DELETE / ON UPDATE and must be complete: a bare SET, a bare
// NO, or any other token is not an action.
//
//	RESTRICT | CASCADE | SET NULL | SET DEFAULT | NO ACTION
//
// Ref: delete_option (RESTRICT | CASCADE | SET NULL | NO ACTION | SET DEFAULT)
// in MySQL 8.0.46 sql/sql_yacc.yy.
func (p *Parser) parseReferenceAction() (nodes.ReferenceAction, error) {
	switch p.cur.Type {
	case kwRESTRICT:
		p.advance()
		return nodes.RefActRestrict, nil
	case kwCASCADE:
		p.advance()
		return nodes.RefActCascade, nil
	case kwSET:
		p.advance()
		if _, ok := p.match(kwNULL); ok {
			return nodes.RefActSetNull, nil
		}
		if _, ok := p.match(kwDEFAULT); ok {
			return nodes.RefActSetDefault, nil
		}
		return nodes.RefActNone, p.syntaxErrorAtCur()
	case kwNO:
		p.advance()
		if _, ok := p.match(kwACTION); ok {
			return nodes.RefActNoAction, nil
		}
		return nodes.RefActNone, p.syntaxErrorAtCur()
	}
	return nodes.RefActNone, p.syntaxErrorAtCur()
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
//
// inAlter is true for ALTER TABLE ... ADD, whose PRIMARY KEY grammar differs from
// CREATE TABLE's (parsePrimaryKey).
func (p *Parser) parseTableConstraint(inAlter bool) (*nodes.Constraint, error) {
	start := p.pos()
	constr := &nodes.Constraint{Loc: nodes.Loc{Start: start}}

	// Optional CONSTRAINT [name]
	_, afterConstraint := p.match(kwCONSTRAINT)
	if afterConstraint {
		if p.isIdentToken() && p.cur.Type != kwPRIMARY && p.cur.Type != kwUNIQUE &&
			p.cur.Type != kwFOREIGN && p.cur.Type != kwCHECK {
			constr.Name, _, _ = p.parseIdent()
		}
		// CONSTRAINT goes with these four only: 4.3.5.6 answers 1064 for `CONSTRAINT c
		// INDEX (a)`, `CONSTRAINT KEY (a)`, `CONSTRAINT FULLTEXT (c)`, `CONSTRAINT c SPATIAL
		// (g)` and `CONSTRAINT VECTOR INDEX (v)`, here and after ALTER TABLE ADD.
		//
		// Ref: table_element, constraint_definition and add_constraint_uniq_key in
		// sql_parser_mysql_mode.y (v4.3.5_CE_BP6).
		switch p.cur.Type {
		case kwPRIMARY, kwUNIQUE, kwFOREIGN, kwCHECK:
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	switch p.cur.Type {
	case kwPRIMARY:
		if err := p.parsePrimaryKey(constr, afterConstraint, inAlter); err != nil {
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
		idxCols, err := p.parseParenIndexKeyParts()
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
		cols, err := p.parseParenIdentListWithCompletion("columnref")
		if err != nil {
			return nil, err
		}
		constr.Columns = cols
		// REFERENCES
		if p.cur.Type == kwREFERENCES {
			p.advance()
			// Completion: after FOREIGN KEY (...) REFERENCES, offer table_ref.
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("table_ref")
				return nil, &ParseError{Message: "collecting"}
			}
			ref, err := p.parseRelationFactor()
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
					act, err := p.parseReferenceAction()
					if err != nil {
						return nil, err
					}
					constr.OnDelete = act
				} else if _, ok := p.match(kwUPDATE); ok {
					act, err := p.parseReferenceAction()
					if err != nil {
						return nil, err
					}
					constr.OnUpdate = act
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
		idxCols, err := p.parseParenIndexKeyParts()
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
		// FULLTEXT opt_key_or_index opt_index_name opt_index_using_algorithm '('
		// (v4.3.5_CE_BP6): 4.3.5.6 takes FULLTEXT KEY fk USING BTREE (c). USING
		// needs its algorithm: FULLTEXT KEY fk USING (c) is a 1064.
		if p.cur.Type == kwUSING {
			if next := p.peekNext(); next.Type != tokIDENT && next.Type < 700 {
				return nil, p.syntaxErrorAtTok(next)
			}
			if err := p.parseIndexTypeClause(constr); err != nil {
				return nil, err
			}
		}
		idxCols, err := p.parseParenIndexKeyParts()
		if err != nil {
			return nil, err
		}
		constr.IndexColumns = idxCols
		constr.Columns = indexColumnsToNames(idxCols)
		if err := p.parseConstraintIndexOptions(constr); err != nil {
			return nil, err
		}

	case kwVECTOR:
		// VECTOR {INDEX | KEY} [index_name] (key_part, ...) [index_option] ...
		//
		// CREATE TABLE's only: 4.3.5.6 answers 1064 for ALTER TABLE t ADD VECTOR INDEX
		// i (v) in every spelling, adds one with CREATE VECTOR INDEX alone, and takes
		// no CONSTRAINT symbol in front of it. The key parts are an ordinary list that
		// it refuses later: VECTOR INDEX (a, b) and VECTOR INDEX ((f(a))) are 1235.
		// There is no index type in front of the key parts, where the options take
		// BTREE and HASH like any other index's and never HNSW.
		if inAlter {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume VECTOR
		p.advance() // consume INDEX or KEY; isTableConstraintStart saw it
		constr.Type = nodes.ConstrVector
		if p.isIdentToken() && p.cur.Type != '(' {
			name, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			constr.Name = name
		}
		idxCols, err := p.parseParenIndexKeyParts()
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
		idxCols, err := p.parseParenIndexKeyParts()
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

// parsePrimaryKey parses a table-level PRIMARY KEY from PRIMARY on. OceanBase's
// grammar for it is narrower than MySQL's, and differs between CREATE TABLE and
// ALTER TABLE ... ADD. sql_parser_mysql_mode.y (v4.3.5_CE_BP6) reads, in CREATE
// TABLE's table_element,
//
//	CONSTRAINT [name] PRIMARY KEY [USING type] (column_name_list) [USING type] [COMMENT 'string']
//	PRIMARY KEY [index_name] [USING type] (column_name_list) [USING type] [COMMENT 'string']
//
// and in ALTER TABLE's add_constraint_pri_key and add_primary_key,
//
//	CONSTRAINT [name] PRIMARY KEY (column_name_list) [index_option ...]
//	PRIMARY KEY [index_name] (column_name_list) [index_option ...]
//
// 4.3.5.6 answers 1064 for everything else (primaryKeyGrammarCases): KEY left
// out, an index name after CONSTRAINT, a USING before the columns in ALTER TABLE,
// a column prefix, ASC or DESC, or an expression in the column list, and in CREATE
// TABLE any option but one USING then one COMMENT, so KEY_BLOCK_SIZE, VISIBLE,
// INVISIBLE, WITH PARSER, LOCAL, GLOBAL and BLOCK_SIZE, and a COMMENT before the
// USING. Each is left unread, so the statement fails at it, as OceanBase's "near"
// says.
func (p *Parser) parsePrimaryKey(constr *nodes.Constraint, afterConstraint, inAlter bool) error {
	p.advance() // PRIMARY
	if _, err := p.expect(kwKEY); err != nil {
		return err
	}
	constr.Type = nodes.ConstrPrimaryKey
	// The index name, which the engine ignores, as MySQL does.
	if p.isIdentToken() && p.cur.Type != '(' {
		if afterConstraint {
			return p.syntaxErrorAtCur()
		}
		if _, _, err := p.parseIdent(); err != nil {
			return err
		}
	}
	if p.cur.Type == kwUSING {
		if inAlter {
			return p.syntaxErrorAtCur()
		}
		if err := p.parseIndexTypeClause(constr); err != nil {
			return err
		}
	}
	cols, err := p.parsePrimaryKeyColumns()
	if err != nil {
		return err
	}
	constr.IndexColumns = cols
	constr.Columns = indexColumnsToNames(cols)
	if inAlter {
		return p.parseConstraintIndexOptions(constr)
	}
	for _, kw := range []int{kwUSING, kwCOMMENT} {
		if p.cur.Type != kw {
			continue
		}
		opt, _, err := p.parseIndexOption(false)
		if err != nil {
			return err
		}
		constr.IndexOptions = append(constr.IndexOptions, opt)
	}
	return nil
}

// parsePrimaryKeyColumns parses a PRIMARY KEY's parenthesized column_name_list:
// plain column names, which is all OceanBase takes there.
func (p *Parser) parsePrimaryKeyColumns() ([]*nodes.IndexColumn, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var cols []*nodes.IndexColumn
	for {
		// Completion: offer columnref, as for any key part.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, &ParseError{Message: "collecting"}
		}
		start := p.pos()
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		loc := nodes.Loc{Start: start, End: p.prev.End}
		cols = append(cols, &nodes.IndexColumn{Loc: loc, Expr: &nodes.ColumnRef{Loc: loc, Column: name}})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseConstraintIndexOptions parses index_option list after a table constraint's key
// parts. An option that starts and does not finish is an error, as it is after CREATE
// INDEX: dropping it left the rest to whatever followed, so `KEY i (a) COMMENT` and
// `FULLTEXT KEY fk (c) WITH PARSER)` parsed, which OceanBase answers 1064 for.
func (p *Parser) parseConstraintIndexOptions(constr *nodes.Constraint) error {
	for {
		opt, ok, err := p.parseIndexOption(constr.Type == nodes.ConstrFulltextIndex)
		if err != nil {
			return err
		}
		if !ok {
			break
		}
		constr.IndexOptions = append(constr.IndexOptions, opt)
	}
	// A key, a unique key or a spatial key ends in with_column_group, in a table
	// element and in ALTER TABLE ... ADD (table_element and alter_index_option in
	// sql_parser_mysql_mode.y, oceanbase v4.3.5_CE_BP6). A primary, foreign or
	// fulltext key does not.
	switch constr.Type {
	case nodes.ConstrIndex, nodes.ConstrUnique, nodes.ConstrSpatialIndex:
		group, _, err := p.parseWithColumnGroup()
		if err != nil {
			return err
		}
		constr.ColumnGroup = group
	}
	return nil
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
		p.match('=')
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

	case kwAUTO_INCREMENT:
		p.advance()
		p.match('=')
		// The value is a number, where the other counts are integers: OceanBase 4.3.5
		// parses AUTO_INCREMENT = 99999999999999999999 and AUTO_INCREMENT = 1.5.
		if p.cur.Type == tokFCONST {
			val := p.advance().Str
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "AUTO_INCREMENT", Value: val}, true, nil
		}
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
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "CHARSET", Value: val}, true, nil

	case kwCHARACTER:
		p.advance()
		if _, ok := p.match(kwSET); ok {
			p.match('=')
			val, err := p.consumeOptionValue()
			if err != nil {
				return nil, false, err
			}
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "CHARACTER SET", Value: val}, true, nil
		}
		// CHARACTER without SET is no option. Returning "not an option" after consuming
		// it dropped the word, so `CREATE TABLE t (a INT) CHARACTER` and `ALTER TABLE t
		// COMMENT 'x' CHARACTER` parsed; OceanBase answers 1064 for both.
		return nil, false, p.syntaxErrorAtCur()

	case kwCOLLATE:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "COLLATE", Value: val}, true, nil

	case kwCOMMENT:
		p.advance()
		p.match('=')
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
		p.match('=')
		val, err := p.parseWordOf(rowFormatWords)
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "ROW_FORMAT", Value: val}, true, nil

	case kwCONNECTION:
		p.advance()
		p.match('=')
		if p.cur.Type != tokSCONST {
			return nil, false, p.syntaxErrorAtCur()
		}
		val := p.cur.Str
		p.advance()
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "CONNECTION", Value: val}, true, nil

	case kwPASSWORD:
		p.advance()
		p.match('=')
		if p.cur.Type != tokSCONST {
			return nil, false, p.syntaxErrorAtCur()
		}
		val := p.cur.Str
		p.advance()
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "PASSWORD", Value: val}, true, nil

	// STORAGE is deliberately not a table option here. Unlike the mysql and
	// mariadb forks, OceanBase CE 4.3.5 answers 1064 for `STORAGE DISK`,
	// `STORAGE MEMORY` and every other STORAGE value in the table-option
	// position (measured), so it is left unmatched and the statement fails on
	// the unexpected keyword. (STORAGE stays valid as a column option and
	// inside a TABLESPACE clause, which are handled elsewhere.)

	case kwSTART:
		// START TRANSACTION (used for versioned tables)
		if next := p.peekNext(); next.Type == kwTRANSACTION {
			p.advance() // consume START
			p.advance() // consume TRANSACTION
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "START TRANSACTION", Value: ""}, true, nil
		}
		return nil, false, nil

	// CHECKSUM [=] n. The value is an integer and only an integer: OceanBase CE
	// 4.3.5.6 answers 1064 for `CHECKSUM = 'abc'`, for `= DEFAULT` and for
	// `= 1.5`, where MySQL's wider `ulong_num` takes the decimal. MySQL also
	// names this option TABLE_CHECKSUM and OceanBase does not — it answers 1064
	// for every spelling of it — so that name is not taken here.
	case kwCHECKSUM:
		p.advance()
		p.match('=')
		if p.cur.Type == tokICONST && !p.cur.isInt64() {
			return nil, false, p.syntaxErrorAtCur()
		}
		val, err := p.consumeIntegerOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "CHECKSUM", Value: val}, true, nil

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
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "ENCRYPTION", Value: val}, true, nil

	case kwUNION:
		// UNION [=] '(' [tbl_name [, tbl_name]...] ')'
		//
		// The parentheses are required and the list inside them may be empty,
		// and each name may be qualified: OceanBase CE 4.3.5.6 runs `UNION = ()`,
		// `UNION ()` and `UNION = (d.t1, d.t2)`, and answers 1064 for a bare
		// `UNION`, for `UNION = (t1,)` and for `UNION = (,)`.
		//
		// Ref: `UNION opt_equal '(' opt_table_list ')'` in MySQL 8.0.46
		// sql/sql_yacc.yy.
		p.advance()
		p.match('=')
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
			p.match('=')
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
			p.match('=')
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
		insertMethod := p.cur.Type == kwINSERT_METHOD
		optName := p.cur.Str
		p.advance()
		p.match('=')
		var val string
		var err error
		if insertMethod {
			val, err = p.parseWordOf(insertMethodWords)
		} else {
			val, err = p.consumeOptionValue()
		}
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil
	}

	// DYNAMIC_PARTITION_POLICY [=] (key = value [[,] key = value] ...). Nothing
	// downstream reads inside the list, so it is kept as written rather than
	// modelled, but every assignment is parsed the way 4.3.5.6 reads it
	// (dynamicPartitionPolicyCases): the key is one of dynamicPartitionPolicyKeys,
	// written bare; the '=' is required; ENABLE takes TRUE or FALSE, and every
	// other key a string literal, adjacent literals included. The engine answers
	// 1064 for anything else, so a looser reading would pass DDL through review for
	// the server to refuse. The comma between assignments is the one exception: the
	// engine parses a list without it and refuses it afterwards with 1210, so the
	// parser takes it, as it takes any statement the grammar does.
	if p.cur.Type == kwDYNAMIC_PARTITION_POLICY {
		p.advance()
		p.match('=')
		open := p.pos()
		if _, err := p.expect('('); err != nil {
			return nil, false, err
		}
		for {
			key, ok := p.dynamicPartitionPolicyKey()
			if !ok {
				return nil, false, p.syntaxErrorAtCur()
			}
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, false, err
			}
			if key == "ENABLE" {
				if p.cur.Type != kwTRUE && p.cur.Type != kwFALSE {
					return nil, false, p.syntaxErrorAtCur()
				}
				p.advance()
			} else {
				if p.cur.Type != tokSCONST {
					return nil, false, p.syntaxErrorAtCur()
				}
				for p.cur.Type == tokSCONST {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				break
			}
			p.match(',')
		}
		if _, err := p.expect(')'); err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  "DYNAMIC_PARTITION_POLICY",
			Value: p.inputText(open, p.prev.End),
		}, true, nil
	}

	// ORGANIZATION [=] {INDEX | HEAP}: OceanBase's index-organized or heap table.
	// 4.3.5 prints ORGANIZATION INDEX on every table SHOW CREATE TABLE returns, and
	// takes either word, with or without the '=', in CREATE TABLE and ALTER TABLE
	// (organizationCases). The option keeps the name it has always had for INDEX.
	if p.cur.Type == kwORGANIZATION {
		p.advance()
		p.match('=')
		var name string
		switch {
		case p.cur.Type == kwINDEX:
			name = "ORGANIZATION INDEX"
		case p.isBareWord("HEAP"):
			name = "ORGANIZATION HEAP"
		default:
			return nil, false, p.syntaxErrorAtCur()
		}
		p.advance()
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: name}, true, nil
	}

	// OceanBase's storage options each take one kind of value. table_option in
	// sql_parser_mysql_mode.y (v4.3.5_CE) reads `REPLICA_NUM opt_equal_mark INTNUM`,
	// BLOCK_SIZE, TABLET_SIZE and PCTFREE the same way, `USE_BLOOM_FILTER
	// opt_equal_mark BOOL_VALUE`, and `AUTO_INCREMENT_MODE opt_equal_mark
	// STRING_VALUE`; ENABLE_MACRO_BLOCK_BLOOM_FILTER is newer than that tag, and
	// 4.3.5.6 takes TRUE or FALSE for it. Any other value is a 1064
	// (tableOptionValueCases), where consumeOptionValue below takes any word.
	if want, ok := oceanBaseOptionValue[p.cur.Type]; ok {
		optName := p.cur.Str
		p.advance()
		p.match('=')
		if !want(p.cur.Type) {
			return nil, false, &ParseError{Message: "expected a value for the table option", Position: p.cur.Loc}
		}
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil
	}

	// Handle identifier-based options: KEY_BLOCK_SIZE, MAX_ROWS, etc.
	//
	// ENGINE_ATTRIBUTE and SECONDARY_ENGINE_ATTRIBUTE are absent: OceanBase CE
	// 4.3.5.6 answers 1064 for both as a table option, as a column option and on
	// ALTER TABLE, so reading them here only passed statements the server
	// refuses.
	switch p.cur.Type {
	// STATS_AUTO_RECALC [=] {DEFAULT | n}. MySQL takes only 0 and 1 here;
	// OceanBase CE 4.3.5.6 takes any integer — it accepts `= 2` and `= 10` —
	// and answers 1064 for a string, for a negative and for a decimal.
	case kwSTATS_AUTO_RECALC:
		optName := p.cur.Str
		p.advance()
		p.match('=')
		val, err := p.consumeIntegerOrDefaultOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil

	case kwSTATS_PERSISTENT:
		// STATS_PERSISTENT [=] {DEFAULT | integer}: the value is DEFAULT or a
		// number, so a string and a bare identifier are each a 1064, matching
		// the mysql and mariadb forks. (Those bound the number to 0 or 1; an
		// out-of-range number is left to a separate measured change here, so
		// this only closes the string/identifier hole.)
		optName := p.cur.Str
		p.advance()
		p.match('=')
		val, err := p.consumeIntegerOrDefaultOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil

	case kwKEY_BLOCK_SIZE,
		kwSTATS_SAMPLE_PAGES, kwMAX_ROWS, kwMIN_ROWS, kwAVG_ROW_LENGTH,
		kwPACK_KEYS, kwDELAY_KEY_WRITE, kwSECONDARY_ENGINE,
		kwAUTOEXTEND_SIZE:
		optType := p.cur.Type
		optName := p.cur.Str
		p.advance()
		p.match('=')
		// Two of these take an integer of no more than 63 bits: OceanBase 4.3.5
		// answers 1064 for AVG_ROW_LENGTH = 9223372036854775808 and for
		// DELAY_KEY_WRITE = 9223372036854775808, as it does for CHECKSUM.
		if (optType == kwAVG_ROW_LENGTH || optType == kwDELAY_KEY_WRITE) && p.cur.Type == tokICONST && !p.cur.isInt64() {
			return nil, false, p.syntaxErrorAtCur()
		}
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil
	}

	return p.parseOceanBaseTableOption(start)
}

// oceanBaseOptionValue maps each OceanBase storage option to the token types its
// value may be.
var oceanBaseOptionValue = map[int]func(int) bool{
	kwBLOCK_SIZE:                      isIntegerToken,
	kwREPLICA_NUM:                     isIntegerToken,
	kwTABLET_SIZE:                     isIntegerToken,
	kwPCTFREE:                         isIntegerToken,
	kwUSE_BLOOM_FILTER:                isBoolToken,
	kwENABLE_MACRO_BLOCK_BLOOM_FILTER: isBoolToken,
	kwAUTO_INCREMENT_MODE:             isStringToken,
}

func isIntegerToken(t int) bool { return t == tokICONST }
func isBoolToken(t int) bool    { return t == kwTRUE || t == kwFALSE }
func isStringToken(t int) bool  { return t == tokSCONST }

// dynamicPartitionPolicyKeys are the keys DYNAMIC_PARTITION_POLICY takes on
// 4.3.5.6. Any other key, a quoted one included, is a 1064.
var dynamicPartitionPolicyKeys = []string{
	"ENABLE", "TIME_UNIT", "PRECREATE_TIME", "EXPIRE_TIME", "TIME_ZONE", "BIGINT_PRECISION",
}

// dynamicPartitionPolicyKey reports which of dynamicPartitionPolicyKeys the
// current token spells, without consuming it.
func (p *Parser) dynamicPartitionPolicyKey() (string, bool) {
	for _, key := range dynamicPartitionPolicyKeys {
		if p.isBareWord(key) {
			return key, true
		}
	}
	return "", false
}

// The words an option takes, each list a rule of sql_parser_mysql_mode.y
// (oceanbase v4.3.5_CE_BP6): algorithm_opt and lock_opt for ALTER TABLE's
// ALGORITHM and LOCK, row_format_option, merge_insert_types for INSERT_METHOD,
// and view_algorithm. Any other word is a 1064.
var (
	alterAlgorithmWords = []string{"DEFAULT", "INPLACE", "COPY", "INSTANT"}
	alterLockWords      = []string{"DEFAULT", "NONE", "SHARED", "EXCLUSIVE"}
	rowFormatWords      = []string{"REDUNDANT", "COMPACT", "DYNAMIC", "COMPRESSED", "CONDENSED", "DEFAULT"}
	insertMethodWords   = []string{"NO", "FIRST", "LAST"}
	viewAlgorithmWords  = []string{"UNDEFINED", "MERGE", "TEMPTABLE"}
)

// parseWordOf reads the current token where it is one of words, written bare,
// and returns it as written.
func (p *Parser) parseWordOf(words []string) (string, error) {
	if _, ok := p.curBareWordOf(words); !ok {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// isBareWord reports whether the current token is word, in any case and not
// quoted: an identifier or a keyword, but not `word`.
func (p *Parser) isBareWord(word string) bool {
	if p.cur.Type != tokIDENT && p.cur.Type < 700 {
		return false
	}
	if !strings.EqualFold(p.cur.Str, word) {
		return false
	}
	return !strings.HasPrefix(p.inputText(p.cur.Loc, p.cur.End), "`")
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
		if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
			name, _, err := p.parseKeywordOrIdent()
			if err != nil {
				return "", err
			}
			return name, nil
		}
		// Every option routed here takes a value, so a missing one is a syntax error
		// rather than an empty string: returning "" let a truncated option parse.
		return "", &ParseError{Message: "expected a value for the table option", Position: p.cur.Loc}
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

// consumeIntegerOrDefaultOptionValue reads an option value that must be an
// integer or DEFAULT. A string, a decimal and a leading sign are each a 1064.
func (p *Parser) consumeIntegerOrDefaultOptionValue() (string, error) {
	if p.cur.Type == kwDEFAULT {
		val := p.cur.Str
		p.advance()
		return val, nil
	}
	return p.consumeIntegerOptionValue()
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
		expr, err := p.parseExpr()
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
		if _, ok := p.match(kwCOLUMNS); ok {
			cols, err := p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
			part.Columns = cols
		} else {
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
			part.Expr = expr
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}

	case kwLIST:
		p.advance()
		part.Type = nodes.PartitionList
		if _, ok := p.match(kwCOLUMNS); ok {
			cols, err := p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
			part.Columns = cols
		} else {
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
			part.Expr = expr
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}

	default:
		return nil, &ParseError{
			Message:  "expected HASH, KEY, RANGE, or LIST after PARTITION BY",
			Position: p.cur.Loc,
		}
	}

	// PARTITIONS num
	if p.cur.Type == kwPARTITIONS {
		p.advance()
		if p.cur.Type == tokICONST {
			part.NumParts = int(p.cur.Ival)
			p.advance()
		}
	}

	// SUBPARTITION BY ..., then a SUBPARTITION TEMPLATE or SUBPARTITIONS n.
	if p.cur.Type == kwSUBPARTITION {
		if err := p.parseSubpartitionClause(part); err != nil {
			return nil, err
		}
		// OceanBase's partition_options also takes PARTITIONS n after the subpartition
		// clause (v4.3.5_CE_BP6): 4.3.5.6 takes PARTITION BY HASH(a) SUBPARTITION BY
		// KEY(b) SUBPARTITIONS 2 PARTITIONS 3. It takes it once, before or after.
		if part.NumParts == 0 && p.cur.Type == kwPARTITIONS {
			p.advance()
			// partition_num is PARTITIONS INTNUM: the count is not optional.
			if p.cur.Type != tokICONST {
				return nil, p.syntaxErrorAtCur()
			}
			part.NumParts = int(p.cur.Ival)
			p.advance()
		}
	}

	// Optional partition definitions
	if p.cur.Type == '(' {
		p.advance()
		for {
			pd, err := p.parsePartitionDef(false)
			if err != nil {
				return nil, err
			}
			part.Partitions = append(part.Partitions, pd)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// With a SUBPARTITION TEMPLATE no partition may list its own subpartitions:
	// 4.3.5.6 answers 1064 for `SUBPARTITION TEMPLATE (...) (PARTITION p0 VALUES LESS
	// THAN (10) (SUBPARTITION z0), ...)` (templateCases).
	if len(part.SubPartTemplate) > 0 {
		for _, pd := range part.Partitions {
			if len(pd.SubPartitions) > 0 {
				sub := pd.SubPartitions[0]
				line, col := p.lineCol(sub.Loc.Start)
				return nil, &ParseError{
					Message:  "a partition cannot list subpartitions when the table has a SUBPARTITION TEMPLATE",
					Position: sub.Loc.Start,
					Line:     line,
					Column:   col,
				}
			}
		}
	}

	part.Loc.End = p.prev.End
	return part, nil
}

// parsePartitionDef parses a single partition definition.
//
//	PARTITION partition_name
//	    [VALUES {LESS THAN (expr|MAXVALUE) | IN (value_list)}]
//	    [ID n] [ENGINE = InnoDB]
//	    [(subpartition_definition, ...)]
//
// requireValues is set where OceanBase takes only RANGE and LIST partitions:
// ALTER TABLE ... ADD PARTITION reads opt_partition_range_or_list
// (sql_parser_mysql_mode.y, v4.3.5_CE_BP6), so 4.3.5.6 answers 1064 for a partition
// there without VALUES, at the token after its name (partitionOptionCases).
func (p *Parser) parsePartitionDef(requireValues bool) (*nodes.PartitionDef, error) {
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

	if requireValues && p.cur.Type != kwVALUES {
		return nil, p.syntaxErrorAtCur()
	}

	// VALUES LESS THAN or VALUES IN
	values, _, err := p.parsePartitionValues()
	if err != nil {
		return nil, err
	}
	pd.Values = values

	// MySQL takes the table options a partition may carry (ENGINE, COMMENT, DATA
	// DIRECTORY, MAX_ROWS, TABLESPACE, ...) in any number. OceanBase's
	// range_partition_element, list_partition_element and hash_partition_element read
	// `opt_part_id opt_engine_option` (v4.3.5_CE_BP6): ID n, then ENGINE = InnoDB.
	// 4.3.5.6 answers 1064 for anything else there, COMMENT with or without '='
	// included, so it is left unread and the statement fails at it
	// (partitionOptionCases).
	if p.isBareWord("ID") {
		optStart := p.pos()
		p.advance()
		if p.cur.Type != tokICONST {
			return nil, p.syntaxErrorAtCur()
		}
		val := p.cur.Str
		p.advance()
		pd.Options = append(pd.Options, &nodes.TableOption{Loc: nodes.Loc{Start: optStart, End: p.prev.End}, Name: "ID", Value: val})
	}
	if opt, err := p.parsePartitionEngine(); err != nil {
		return nil, err
	} else if opt != nil {
		pd.Options = append(pd.Options, opt)
	}

	// Optional subpartition definitions. OceanBase's opt_subpartition_list is a list
	// of HASH, RANGE or LIST subpartitions, without VALUES, with VALUES LESS THAN or
	// with VALUES IN, and the first entry sets which (v4.3.5_CE_BP6): 4.3.5.6 answers
	// 1064 for a list that mixes them. Whether the kind matches SUBPARTITION BY it
	// checks after parsing, with 1210.
	if p.cur.Type == '(' {
		p.advance()
		kind := noPartitionValues
		for i := 0; ; i++ {
			spd, got, err := p.parseSubPartitionDef(kind, i == 0)
			if err != nil {
				return nil, err
			}
			kind = got
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

// partitionValuesKind tells a partition's VALUES LESS THAN from its VALUES IN.
type partitionValuesKind int

const (
	noPartitionValues partitionValuesKind = iota
	valuesLessThan
	valuesIn
)

// parsePartitionValues parses a partition's or a subpartition's VALUES LESS THAN or
// VALUES IN when the next token is VALUES, and reports which it read. A partition's
// VALUES has always taken THAN after LESS, and the bound's closing parenthesis, as
// optional here; strict, set for a subpartition, holds it to the grammar's LESS THAN
// and (...), as OceanBase does (subpartitionCases).
func (p *Parser) parsePartitionValues() (nodes.Node, partitionValuesKind, error) {
	var values nodes.Node
	kind := noPartitionValues
	if _, ok := p.match(kwVALUES); ok {
		if p.cur.Type == kwLESS {
			kind = valuesLessThan
			p.advance() // LESS
			if _, err := p.expect(kwTHAN); err != nil {
				return nil, 0, err
			}
			if p.cur.Type == '(' {
				p.advance()
				// Parse comma-separated list of values (may include MAXVALUE tokens and expressions).
				var items []nodes.Node
				allMax := true
				for {
					if p.cur.Type == kwMAXVALUE {
						items = append(items, &nodes.String{Str: "MAXVALUE"})
						p.advance()
					} else {
						allMax = false
						expr, err := p.parseExpr()
						if err != nil {
							return nil, 0, err
						}
						if expr == nil {
							return nil, 0, p.syntaxErrorAtCur()
						}
						items = append(items, expr)
					}
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
				if _, err := p.expect(')'); err != nil {
					return nil, 0, err
				}
				if allMax && len(items) > 0 {
					// All MAXVALUE — store as single string marker
					values = &nodes.String{Str: "MAXVALUE"}
				} else if len(items) == 1 {
					values = items[0]
				} else {
					values = &nodes.List{Items: items}
				}
			} else if p.cur.Type == kwMAXVALUE {
				values = &nodes.String{Str: "MAXVALUE"}
				p.advance()
			} else {
				return nil, 0, p.syntaxErrorAtCur()
			}
		} else if _, ok := p.match(kwIN); ok {
			kind = valuesIn
			if p.cur.Type == '(' {
				p.advance()
				var items []nodes.Node
				for {
					if p.cur.Type == '(' {
						// Tuple: (val1, val2, ...) for multi-column LIST COLUMNS
						p.advance()
						var tupleVals []nodes.Node
						for {
							expr, err := p.parseExpr()
							if err != nil {
								return nil, 0, err
							}
							if expr == nil {
								return nil, 0, p.syntaxErrorAtCur()
							}
							tupleVals = append(tupleVals, expr)
							if p.cur.Type != ',' {
								break
							}
							p.advance()
						}
						if _, err := p.expect(')'); err != nil {
							return nil, 0, err
						}
						items = append(items, &nodes.List{Items: tupleVals})
					} else {
						expr, err := p.parseExpr()
						if err != nil {
							return nil, 0, err
						}
						if expr == nil {
							return nil, 0, p.syntaxErrorAtCur()
						}
						items = append(items, expr)
					}
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
				if _, err := p.expect(')'); err != nil {
					return nil, 0, err
				}
				values = &nodes.List{Items: items}
			} else {
				return nil, 0, p.syntaxErrorAtCur()
			}
		} else {
			// VALUES with neither LESS THAN nor IN after it is no bound: 4.3.5.6
			// answers 1064 at the next token, for a partition and a subpartition
			// (bareValuesCases).
			return nil, 0, p.syntaxErrorAtCur()
		}
	}

	return values, kind, nil
}

// parseSubPartitionDef parses a single subpartition definition.
//
//	SUBPARTITION logical_name [VALUES {LESS THAN (...) | IN (...)}] [ENGINE = InnoDB]
//
// OceanBase's subpartitions take VALUES as its partitions do: a HASH or KEY
// subpartition none, a RANGE one LESS THAN and a LIST one IN. want is the kind this
// one must be, unless first is set, and the kind read is returned. A subpartition of
// the wrong kind fails at its VALUES, at LESS or IN, or where VALUES was due, as
// 4.3.5.6's "near" says (subpartitionCases).
//
// OceanBase's subpartition elements, in a SUBPARTITION TEMPLATE or under a partition,
// end in opt_engine_option alone (v4.3.5_CE_BP6): no ID, and none of the table options
// MySQL takes there, COMMENT included (partitionOptionCases).
func (p *Parser) parseSubPartitionDef(want partitionValuesKind, first bool) (*nodes.SubPartitionDef, partitionValuesKind, error) {
	start := p.pos()
	if _, err := p.expect(kwSUBPARTITION); err != nil {
		return nil, 0, err
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, 0, err
	}
	spd := &nodes.SubPartitionDef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}
	got := noPartitionValues
	if p.cur.Type == kwVALUES {
		switch next := p.peekNext(); next.Type {
		case kwLESS:
			got = valuesLessThan
		case kwIN:
			got = valuesIn
		default:
			// A bare VALUES fails where its bound was due, or at VALUES itself for a
			// HASH or KEY template entry, which takes none: 4.3.5.6 says near
			// 'VALUES)' there (subpartitionCases).
			if !first && want == noPartitionValues {
				return nil, 0, p.syntaxErrorAtCur()
			}
			return nil, 0, p.syntaxErrorAtTok(next)
		}
	}
	if !first && got != want {
		if got != noPartitionValues && want != noPartitionValues {
			return nil, 0, p.syntaxErrorAtTok(p.peekNext())
		}
		return nil, 0, p.syntaxErrorAtCur()
	}
	values, got, err := p.parsePartitionValues()
	if err != nil {
		return nil, 0, err
	}
	// VALUES LESS THAN or VALUES IN with no bound after it is no bound at all.
	if got != noPartitionValues && values == nil {
		return nil, 0, p.syntaxErrorAtCur()
	}
	if err := p.checkLoneDefaultValue(values); err != nil {
		return nil, 0, err
	}
	spd.Values = values
	spd.ValuesIn = got == valuesIn
	if opt, err := p.parsePartitionEngine(); err != nil {
		return nil, 0, err
	} else if opt != nil {
		spd.Options = append(spd.Options, opt)
	}
	spd.Loc.End = p.prev.End
	return spd, got, nil
}

// checkLoneDefaultValue refuses a VALUES IN list that holds DEFAULT beside other values.
// OceanBase's list_partition_expr is '(' list_expr ')' or '(' DEFAULT ')' (v4.3.5_CE_BP6),
// so DEFAULT stands alone: 4.3.5.6 answers 1064 for a subpartition's VALUES IN
// (DEFAULT, 1), at the comma, and at DEFAULT when it is not first.
func (p *Parser) checkLoneDefaultValue(values nodes.Node) error {
	list, ok := values.(*nodes.List)
	if !ok || len(list.Items) < 2 {
		return nil
	}
	for i, item := range list.Items {
		def, ok := item.(*nodes.DefaultExpr)
		if !ok {
			continue
		}
		at := def.Loc.Start
		if i == 0 {
			// The comma after it.
			at = def.Loc.End
			for c := p.inputText(at, at+1); c != "" && strings.ContainsAny(c, " \t\r\n"); c = p.inputText(at, at+1) {
				at++
			}
		}
		line, col := p.lineCol(at)
		return &ParseError{Message: "DEFAULT must be the only value of VALUES IN", Position: at, Line: line, Column: col}
	}
	return nil
}

// parseSubpartitionClause parses OceanBase's subpartition_option from SUBPARTITION on:
//
//	SUBPARTITION BY {HASH (expr) | KEY (column_list) | RANGE (expr)
//	                 | RANGE COLUMNS (column_list) | LIST (expr) | LIST COLUMNS (column_list)}
//	    [SUBPARTITION TEMPLATE (subpartition_definition, ...) | SUBPARTITIONS n]
//
// as subpartition_template_option and subpartition_individual_option read it
// (sql_parser_mysql_mode.y, v4.3.5_CE_BP6). MySQL subpartitions by HASH and KEY only,
// and takes LINEAR, KEY ALGORITHM = n and KEY (). OceanBase subpartitions by RANGE and
// LIST too and takes none of those three, SUBPARTITIONS n only for HASH and KEY and not
// beside a template, and template entries of the subpartitioning's kind. 4.3.5.6
// answers 1064 for the rest (subpartitionCases), and SHOW CREATE TABLE prints RANGE
// and LIST subpartitioning in this form.
func (p *Parser) parseSubpartitionClause(part *nodes.PartitionClause) error {
	p.advance() // SUBPARTITION
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
	want := noPartitionValues
	switch p.cur.Type {
	case kwHASH:
		p.advance()
		part.SubPartType = nodes.PartitionHash
		expr, err := p.parseParenPartitionExpr()
		if err != nil {
			return err
		}
		part.SubPartExpr = expr
	case kwKEY:
		p.advance()
		part.SubPartType = nodes.PartitionKey
		cols, err := p.parseParenIdentList()
		if err != nil {
			return err
		}
		part.SubPartColumns = cols
	case kwRANGE, kwLIST:
		part.SubPartType, want = nodes.PartitionRange, valuesLessThan
		if p.cur.Type == kwLIST {
			part.SubPartType, want = nodes.PartitionList, valuesIn
		}
		p.advance()
		if _, ok := p.match(kwCOLUMNS); ok {
			cols, err := p.parseParenIdentList()
			if err != nil {
				return err
			}
			part.SubPartColumns = cols
		} else {
			expr, err := p.parseParenPartitionExpr()
			if err != nil {
				return err
			}
			part.SubPartExpr = expr
		}
	default:
		return p.syntaxErrorAtCur()
	}

	switch {
	case p.cur.Type == kwSUBPARTITION && p.peekNext().Type == kwTEMPLATE:
		// OceanBase: SUBPARTITION TEMPLATE (subpartition_definition, ...) gives every
		// partition the same subpartitions, where MySQL repeats them inside each
		// partition (verified on 4.3.5).
		p.advance() // SUBPARTITION
		p.advance() // TEMPLATE
		if _, err := p.expect('('); err != nil {
			return err
		}
		for {
			sub, _, err := p.parseSubPartitionDef(want, false)
			if err != nil {
				return err
			}
			part.SubPartTemplate = append(part.SubPartTemplate, sub)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return err
		}
	case p.cur.Type == kwSUBPARTITIONS:
		if want != noPartitionValues {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type == tokICONST {
			part.NumSubParts = int(p.cur.Ival)
			p.advance()
		}
	}
	return nil
}

// parseParenPartitionExpr parses the parenthesized expression of a [SUB]PARTITION BY
// HASH, RANGE or LIST.
func (p *Parser) parseParenPartitionExpr() (nodes.ExprNode, error) {
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
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return expr, nil
}

// parsePartitionEngine parses OceanBase's opt_engine_option, `ENGINE = InnoDB`, the
// one option a partition or subpartition takes besides a partition's ID. The '=' is
// required and InnoDB is the only engine, written bare in any case: 4.3.5.6 answers
// 1064 for ENGINE InnoDB, ENGINE = MyISAM and ENGINE = 'InnoDB'. It returns nil when
// the next token is not ENGINE.
func (p *Parser) parsePartitionEngine() (*nodes.TableOption, error) {
	if p.cur.Type != kwENGINE {
		return nil, nil
	}
	start := p.pos()
	p.advance()
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	if !p.isBareWord("INNODB") {
		return nil, p.syntaxErrorAtCur()
	}
	val := p.cur.Str
	p.advance()
	return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "ENGINE", Value: val}, nil
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
// three are a keyword list, not an identifier: OceanBase CE 4.3.5.6
// answers 1064 for `ADD FOREIGN KEY (a) REFERENCES u(b) MATCH foo`.
//
// Ref: opt_match_clause in MySQL 8.0.46 sql/sql_yacc.yy.
func (p *Parser) parseMatchType() (string, error) {
	if tok, ok := p.match(kwFULL, kwPARTIAL, kwSIMPLE); ok {
		return strings.ToUpper(tok.Str), nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseColumnCheckConstraint parses the CHECK (expr) [[NOT] ENFORCED] of a
// column option, from the CHECK keyword on. start is the option's first byte,
// which is the CONSTRAINT keyword when one introduced it.
func (p *Parser) parseColumnCheckConstraint(start int) (*nodes.ColumnConstraint, error) {
	p.advance() // consume CHECK
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
	if _, err := p.expect(')'); err != nil {
		return nil, err
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
	return cc, nil
}
