package parser

import (
	"math"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
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

	// CREATE TABLE ... AS SELECT
	if p.cur.Type == kwAS {
		p.advance()
		sel, err := p.parseSelectStmt()
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
		// Check for CREATE TABLE ... (query) — query source without AS: a
		// parenthesized query expression, whichever query primary it opens
		// with (SELECT, WITH, TABLE or VALUES). Let parseSelectStmt consume the
		// parentheses and any following set-op / ORDER BY / LIMIT so
		// "(SELECT 1) UNION (SELECT 2)" stays one source.
		if p.isQueryExpressionStart() {
			sel, err := p.parseSelectStmt()
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
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		stmt.Select = sel
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Table options. A comma may stand between two of them and never ends the
	// list (create_table_options in MySQL 8.0 sql/sql_yacc.yy), so a comma
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

	// IGNORE / REPLACE before SELECT in CTAS
	if p.cur.Type == kwIGNORE {
		p.advance()
		stmt.Ignore = true
	} else if p.cur.Type == kwREPLACE {
		p.advance()
		stmt.Replace = true
	}

	// [AS] query_expression after the columns, options and partitioning:
	// as_create_query_expression in MySQL 8.0 sql/sql_yacc.yy, which takes any
	// query expression with or without AS -- SELECT, WITH, TABLE, VALUES or a
	// parenthesized query -- and which IGNORE or REPLACE must be followed by
	// (opt_duplicate_as_qe). Reading only SELECT here used to end the statement
	// in front of "TABLE u" or "(SELECT 1)", which then parsed as a second
	// statement.
	switch {
	case p.cur.Type == kwAS:
		p.advance()
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		stmt.Select = sel
	case p.isQueryExpressionStart():
		sel, err := p.parseSelectStmt()
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

	// Column constraints and options. The generated form is an arm of field_def of its
	// own: only COLLATE comes between the type and AS, and the attributes come after it,
	// so `b INT NOT NULL AS (a)` and `b INT AS (a) AS (a)` are both a 1064.
	//
	// Ref: field_def in MySQL 8.0 sql/sql_yacc.yy.
	sawAttribute := false
	for {
		if sawAttribute && (p.cur.Type == kwGENERATED || (p.cur.Type == kwAS && p.peekNext().Type == '(')) {
			return nil, p.syntaxErrorAtCur()
		}
		isCollate := p.cur.Type == kwCOLLATE
		ok, err := p.parseColumnOption(col)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		if !isCollate {
			sawAttribute = true
		}
	}

	col.Loc.End = p.prev.End
	return col, nil
}

// columnHasCollation reports a column that wrote COLLATE already, beside its type or as an
// attribute. MySQL takes one: `CHAR(5) COLLATE a COLLATE a` is a 1064 even with one name.
func columnHasCollation(col *nodes.ColumnDef) bool {
	if col.TypeName != nil && col.TypeName.Collate != "" {
		return true
	}
	for _, c := range col.Constraints {
		if c.Type == nodes.ColConstrCollate {
			return true
		}
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

	// SERIAL DEFAULT VALUE, the three words together, which
	// https://dev.mysql.com/doc/refman/8.0/en/numeric-type-syntax.html gives as
	// another way to write NOT NULL AUTO_INCREMENT UNIQUE. It is recorded as
	// those three and not as an attribute of its own, so everything downstream —
	// the catalog's implicit primary key among it — sees the column the engine
	// builds. All three words are required: MySQL 8.0 answers 1064 for
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
		if columnHasCollation(col) {
			return false, p.syntaxErrorAtCur()
		}
		p.advance()
		// collate_clause is COLLATE collation_name, and collation_name is
		// ident_or_text | BINARY_SYM (sql_yacc.yy): `c JSON COLLATE BINARY`
		// parses on MySQL 8.0.
		collName, err := p.parseCollationName()
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
		// CHECK is the only continuation: MySQL 8.0.46 runs
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
			expr, err := p.parseNowFunction()
			if err != nil {
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

	// ENGINE_ATTRIBUTE and SECONDARY_ENGINE_ATTRIBUTE column options. The string
	// is required, as it is on the table option: MySQL 8.0 answers 1064 for
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
	// Unparenthesized, it is now_or_signed_literal: a signed literal or a NOW function, and no
	// other expression. `DEFAULT rand()`, `DEFAULT b`, `DEFAULT @a` and `DEFAULT 1 + 1` are each a
	// 1064, and the value goes in parentheses to be one.
	//
	// Ref: column_attribute's DEFAULT_SYM arms and now_or_signed_literal in MySQL 8.0
	// sql/sql_yacc.yy.
	if startsNowFunction(p.cur) {
		return p.parseNowFunction()
	}
	switch p.cur.Type {
	case '+', '-':
		if next := p.peekNext().Type; next != tokICONST && next != tokFCONST {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseExpr()
	case tokICONST, tokFCONST, tokSCONST, tokXCONST, tokBCONST, kwNULL, kwTRUE, kwFALSE:
	default:
		// A literal with an introducer or a temporal keyword reads as one, and nothing else.
		if !p.startsPrefixedLiteral() {
			return nil, p.syntaxErrorAtCur()
		}
	}
	expr, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return expr, nil
}

// startsNowFunction reports NOW, CURRENT_TIMESTAMP, LOCALTIME or LOCALTIMESTAMP, unquoted.
func startsNowFunction(tok Token) bool {
	switch strings.ToUpper(tok.Str) {
	case "NOW", "CURRENT_TIMESTAMP", "LOCALTIME", "LOCALTIMESTAMP":
		return tok.Type != tokSCONST
	}
	return false
}

// parseNowFunction reads what DEFAULT and ON UPDATE take of the clock: one of the names, with
// an integer precision or none. `ON UPDATE 1`, `ON UPDATE NOW() + 1` and `ON UPDATE NOW(1+1)`
// are each a 1064, and a backquoted `NOW` is a stored function.
//
// Ref: now in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseNowFunction() (nodes.ExprNode, error) {
	if !startsNowFunction(p.cur) || p.backtickQuotedAt(p.cur.Loc) {
		return nil, p.syntaxErrorAtCur()
	}
	// NOW is the one of them that is no word on its own: ON UPDATE NOW is a 1064.
	if strings.EqualFold(p.cur.Str, "NOW") && p.peekNext().Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	expr, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if call, ok := expr.(*nodes.FuncCallExpr); ok && len(call.Args) > 0 {
		if _, isInt := call.Args[0].(*nodes.IntLit); !isInt || len(call.Args) > 1 {
			return nil, p.syntaxErrorAtTok(p.prev)
		}
	}
	return expr, nil
}

// startsPrefixedLiteral reports a literal written with an introducer, a national marker or
// the word of a temporal type, each of which the expression reader takes as a primary. An
// introducer takes a string, a hex or a bit literal: `_binary X'FF'` and `_binary 0b1`.
func (p *Parser) startsPrefixedLiteral() bool {
	next := p.peekNext().Type
	switch p.cur.Type {
	case kwDATE, kwTIME, kwTIMESTAMP:
		return next == tokSCONST
	}
	return strings.HasPrefix(p.cur.Str, "_") && (next == tokSCONST || next == tokXCONST || next == tokBCONST)
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
		idxCols, err := p.parseParenIndexKeyParts()
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
		// The column list is a key-part list, not a plain name list: a prefix
		// length and an ordering direction are both allowed. MySQL 8.0.46 runs
		// `ADD FOREIGN KEY (ident, ident(1) ASC) REFERENCES t (a)`.
		//
		// A functional key part is not allowed, though, because a foreign
		// key's list is key_list and not the key_list_with_expression an
		// index takes: MySQL 8.0.46 answers 1064 for
		// `ADD FOREIGN KEY ((a + 1)) REFERENCES u (b)` and for the same list
		// on CREATE TABLE, where it runs `ADD INDEX ((a + 1))`. (TiDB v8.5.5
		// does parse it, through its own IndexPartSpecificationList, and then
		// crashes building the key -- 1105, a nil dereference -- so the fork
		// there keeps the wider list and this one does not.)
		//
		// Ref: `FOREIGN KEY opt_index_name '(' key_list ')' references` in
		// MySQL 8.0.46 sql/sql_yacc.yy.
		fkCols, err := p.parseParenIndexKeyParts()
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
		val, _, err := p.parseIdentOrText()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "ENGINE", Value: val}, true, nil

	case kwAUTO_INCREMENT:
		p.advance()
		p.matchEqual()
		val, err := p.consumeNumericOptionValue()
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
		val, err := p.parseCharsetName()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "CHARSET", Value: val}, true, nil

	case kwCHARACTER:
		p.advance()
		if _, ok := p.match(kwSET); ok {
			p.matchEqual()
			val, err := p.parseCharsetName()
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
		val, err := p.parseCollationName()
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
		val, err := p.consumeWordOptionValue("DEFAULT", "FIXED", "DYNAMIC", "COMPRESSED", "REDUNDANT", "COMPACT")
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
		val, err := p.consumeWordOptionValue("DISK", "MEMORY")
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
	// CHECKSUM_SYM. The value is a number and only a number: MySQL answers 1064
	// for `CHECKSUM = 'abc'`, for `= DEFAULT` and for `= -1`. A decimal is part
	// of `ulong_num`, so `= 1.5` is accepted.
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
		p.matchEqual()
		val, _, err := p.parseIdent()
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
		// ENCRYPTION [=] 'string', which a table holds to Y or N only after parsing.
		p.advance()
		p.matchEqual()
		val, err := p.consumeStringOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "ENCRYPTION", Value: val}, true, nil

	case kwUNION:
		// UNION [=] '(' [tbl_name [, tbl_name]...] ')'
		//
		// The parentheses are required and the list inside them may be empty,
		// and each name may be qualified: MySQL 8.0.46 runs `UNION = ()`,
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
			val, err := p.consumeStringOptionValue()
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
			val, err := p.consumeStringOptionValue()
			if err != nil {
				return nil, false, err
			}
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "DATA DIRECTORY", Value: val}, true, nil
		}
		return nil, false, nil
	}

	// Handle keyword-token-based table options: COMPRESSION, INSERT_METHOD.
	if p.cur.Type == kwCOMPRESSION {
		optName := p.cur.Str
		p.advance()
		p.matchEqual()
		val, err := p.consumeStringOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil
	}
	if p.cur.Type == kwINSERT_METHOD {
		optName := p.cur.Str
		p.advance()
		p.matchEqual()
		val, err := p.consumeWordOptionValue("NO", "FIRST", "LAST")
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil
	}

	// Handle identifier-based options: KEY_BLOCK_SIZE, MAX_ROWS, etc.
	switch p.cur.Type {
	// STATS_AUTO_RECALC [=] {DEFAULT | 0 | 1}, and nothing else: MySQL
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
	// required: MySQL answers 1064 for the bare name and for an `=` with
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
		optType := p.cur.Type
		optName := p.cur.Str
		p.advance()
		p.matchEqual()
		var val string
		var err error
		switch optType {
		case kwKEY_BLOCK_SIZE:
			val, err = p.consumeBoundedOptionValue(0, 65535)
		case kwAVG_ROW_LENGTH:
			val, err = p.consumeBoundedOptionValue(0, math.MaxUint32)
		case kwMAX_ROWS, kwMIN_ROWS, kwDELAY_KEY_WRITE:
			val, err = p.consumeNumericOptionValue()
		case kwSTATS_PERSISTENT, kwPACK_KEYS:
			val, err = p.consumeZeroOneOrDefaultOptionValue()
		case kwSTATS_SAMPLE_PAGES:
			if tok, ok := p.match(kwDEFAULT); ok {
				val = tok.Str
			} else {
				val, err = p.consumeBoundedOptionValue(1, 65535)
			}
		case kwSECONDARY_ENGINE:
			if tok, ok := p.match(kwNULL); ok {
				val = tok.Str
			} else {
				val, _, err = p.parseIdentOrText()
			}
		case kwAUTOEXTEND_SIZE:
			val, err = p.consumeSizeOptionValue()
		}
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil
	}

	return nil, false, nil
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
		// A reserved word is no value, except DEFAULT, which several options take:
		// `MAX_ROWS = NULL` and `PACK_KEYS = SELECT` are each a 1064 (measured).
		if p.cur.Type == tokIDENT || (p.cur.Type >= 700 && (!isReserved(p.cur.Type) || p.cur.Type == kwDEFAULT)) {
			name, _, err := p.parseKeywordOrIdent()
			if err != nil {
				return "", err
			}
			return name, nil
		}
		// Every caller reads an option that takes a value, and MySQL 8.0 answers
		// 1064 when it is missing (`ENGINE =`, `KEY_BLOCK_SIZE` at end of input).
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
	return p.consumeBoundedOptionValue(0, 1)
}

// consumeBoundedOptionValue reads an option value that must be a number the action holds
// to a range. ulong_num reads a decimal with my_strtoll10, which stops at the point, so
// 1.9 and 1e2 are both 1 and it is that the range sees.
//
// Ref: ulong_num and ternary_option in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) consumeBoundedOptionValue(minimum, maximum uint64) (string, error) {
	if p.cur.Type != tokICONST && p.cur.Type != tokFCONST {
		return "", p.syntaxErrorAtCur()
	}
	digits := p.cur.Str
	if i := strings.IndexFunc(digits, func(r rune) bool { return r < '0' || r > '9' }); i >= 0 {
		digits = digits[:i]
	}
	n, err := strconv.ParseUint("0"+digits, 10, 64)
	if err != nil || n < minimum || n > maximum {
		return "", p.syntaxErrorAtCur()
	}
	val := p.cur.Str
	p.advance()
	return val, nil
}

// consumeWordOptionValue reads an option value that is one of a closed list of words,
// which the grammar spells as tokens: the words in any case, and neither a string nor a
// backquoted name.
func (p *Parser) consumeWordOptionValue(words ...string) (string, error) {
	if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
		if !strings.HasPrefix(p.inputText(p.cur.Loc, p.cur.End), "`") {
			for _, word := range words {
				if strings.EqualFold(p.cur.Str, word) {
					return p.advance().Str, nil
				}
			}
		}
	}
	return "", p.syntaxErrorAtCur()
}

// consumeSizeOptionValue reads size_number: an integer, or an integer with a K, M or G
// after it, which lexes as a name. The action answers 1531 for any other name.
//
// Ref: size_number in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) consumeSizeOptionValue() (string, error) {
	if p.cur.Type == tokICONST {
		return p.advance().Str, nil
	}
	if p.cur.Type == tokIDENT {
		digits := strings.TrimRight(p.cur.Str, "kKmMgG")
		if len(p.cur.Str)-len(digits) == 1 && strings.Trim(digits, "0123456789") == "" {
			return p.advance().Str, nil
		}
	}
	return "", p.syntaxErrorAtCur()
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
		expr, err := p.parsePartitionExpr(partitionExprMethod)
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
		algorithm, err := p.parsePartitionKeyAlgorithm()
		if err != nil {
			return nil, err
		}
		part.Algorithm = algorithm
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
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

	case kwRANGE:
		if part.Linear {
			return nil, p.syntaxErrorAtCur()
		}
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
			expr, err := p.parsePartitionExpr(partitionExprMethod)
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
		if part.Linear {
			return nil, p.syntaxErrorAtCur()
		}
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
			expr, err := p.parsePartitionExpr(partitionExprMethod)
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
		n, err := p.parsePartitionCount("partitions")
		if err != nil {
			return nil, err
		}
		part.NumParts = n
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
			expr, err := p.parsePartitionExpr(partitionExprMethod)
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
			algorithm, err := p.parsePartitionKeyAlgorithm()
			if err != nil {
				return nil, err
			}
			part.SubPartAlgo = algorithm
			// Unlike PARTITION BY KEY, the column list is required: `KEY ()` is a 1064.
			cols, err := p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
			part.SubPartColumns = cols
		case p.cur.Type == kwLINEAR:
			p.advance() // LINEAR
			if p.cur.Type == kwHASH {
				p.advance()
				part.SubPartType = nodes.PartitionHash
				if _, err := p.expect('('); err != nil {
					return nil, err
				}
				expr, err := p.parsePartitionExpr(partitionExprMethod)
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
				algorithm, err := p.parsePartitionKeyAlgorithm()
				if err != nil {
					return nil, err
				}
				part.SubPartAlgo = algorithm
				cols, err := p.parseParenIdentList()
				if err != nil {
					return nil, err
				}
				part.SubPartColumns = cols
			} else {
				return nil, p.syntaxErrorAtCur()
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}

		// SUBPARTITIONS num
		if p.cur.Type == kwSUBPARTITIONS {
			p.advance()
			n, err := p.parsePartitionCount("subpartitions")
			if err != nil {
				return nil, err
			}
			part.NumSubParts = n
		}
	}

	// Optional partition definitions
	if p.cur.Type == '(' {
		p.advance()
		var bounds []partitionBound
		for {
			pd, bound, err := p.parsePartitionDefAndBound()
			if err != nil {
				return nil, err
			}
			part.Partitions = append(part.Partitions, pd)
			bounds = append(bounds, bound)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if err := validatePartitionDefList(part.Partitions); err != nil {
			return nil, err
		}
		if err := validatePartitionClause(part, bounds); err != nil {
			return nil, err
		}
	}

	part.Loc.End = p.prev.End
	return part, nil
}

// parsePartitionDef parses a single partition definition.
//
//	PARTITION partition_name
//	    [VALUES {LESS THAN (expr|MAXVALUE) | IN (value_list)}]
//	    [table_options]
func (p *Parser) parsePartitionDef() (*nodes.PartitionDef, error) {
	pd, _, err := p.parsePartitionDefAndBound()
	return pd, err
}

// parsePartitionDefAndBound parses a single partition definition and says which bound it
// wrote.
func (p *Parser) parsePartitionDefAndBound() (*nodes.PartitionDef, partitionBound, error) {
	start := p.pos()
	if _, err := p.expect(kwPARTITION); err != nil {
		return nil, 0, err
	}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, 0, err
	}

	pd := &nodes.PartitionDef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}
	bound := boundNone

	// VALUES LESS THAN or VALUES IN
	if _, ok := p.match(kwVALUES); ok {
		if p.cur.Type == kwLESS {
			bound = boundLessThan
			p.advance() // LESS
			if _, err := p.expect(kwTHAN); err != nil {
				return nil, 0, err
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
						expr, err := p.parsePartitionExpr(partitionExprValue)
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
				return nil, 0, p.syntaxErrorAtCur()
			}
		} else if _, ok := p.match(kwIN); ok {
			bound = boundIn
			if p.cur.Type == '(' {
				p.advance()
				var items []nodes.Node
				for {
					if p.cur.Type == '(' {
						// Tuple: (val1, val2, ...) for multi-column LIST COLUMNS
						p.advance()
						var tupleVals []nodes.Node
						for {
							if p.startsDefaultKeyword() {
								return nil, 0, p.syntaxErrorAtCur()
							}
							expr, err := p.parsePartitionExpr(partitionExprValue)
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
						if p.startsDefaultKeyword() {
							return nil, 0, p.syntaxErrorAtCur()
						}
						expr, err := p.parsePartitionExpr(partitionExprValue)
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
				if err := validatePartitionInRows(items, start); err != nil {
					return nil, 0, err
				}
				pd.Values = &nodes.List{Items: items}
			} else {
				return nil, 0, p.syntaxErrorAtCur()
			}
		} else {
			// VALUES with neither LESS THAN nor IN after it is no bound.
			return nil, 0, p.syntaxErrorAtCur()
		}
	}

	// Optional partition-definition options
	for {
		opt, ok, err := p.parsePartDefOption()
		if err != nil {
			return nil, 0, err
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
				return nil, 0, err
			}
			pd.SubPartitions = append(pd.SubPartitions, spd)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, 0, err
		}
	}

	pd.Loc.End = p.prev.End
	return pd, bound, nil
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
// three are a keyword list, not an identifier: MySQL 8.0.46 answers 1064 for
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
// MySQL 8.0.46 runs `ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE
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
		isEngine := p.cur.Type == kwENGINE
		p.advance()
		p.matchEqual()
		var val string
		var err error
		if isEngine {
			val, _, err = p.parseIdentOrText()
		} else {
			val, _, err = p.parseIdent()
		}
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
