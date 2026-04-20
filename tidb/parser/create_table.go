package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
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
func (p *Parser) parseCreateTableStmt(temporary bool) (*nodes.CreateTableStmt, error) {
	start := p.pos()
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
		stmt.Loc.End = p.pos()
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
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// Parenthesized create definitions (columns and constraints)
	if p.cur.Type == '(' {
		next := p.peekNext()
		// Check for CREATE TABLE ... (SELECT ...) — subquery without AS
		if next.Type == kwSELECT {
			// CREATE TABLE t (SELECT ...)
			p.advance() // consume '('
			sel, err := p.parseSelectStmt()
			if err != nil {
				return nil, err
			}
			p.match(')')
			stmt.Select = sel
			stmt.Loc.End = p.pos()
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
			stmt.Loc.End = p.pos()
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
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// Table options
	for {
		opt, ok, err := p.parseTableOption()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		stmt.Options = append(stmt.Options, opt)
		// Optional comma between table options
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

	// [AS] SELECT ... — bare SELECT or AS SELECT after options/partitions
	if p.cur.Type == kwAS {
		p.advance()
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		stmt.Select = sel
	} else if p.cur.Type == kwSELECT {
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		stmt.Select = sel
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseCreateDefinitions parses column definitions and table constraints
// inside the parentheses of a CREATE TABLE statement.
func (p *Parser) parseCreateDefinitions(stmt *nodes.CreateTableStmt) error {
	for {
		if p.cur.Type == ')' {
			break
		}

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

	col.Loc.End = p.pos()
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
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrNotNull,
		})
		return true, nil

	case kwNULL:
		p.advance()
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
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
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrDefault,
			Expr: expr,
		})
		return true, nil

	case kwAUTO_INCREMENT:
		p.advance()
		col.AutoIncrement = true
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrAutoIncrement,
		})
		return true, nil

	case kwAUTO_RANDOM:
		p.advance()
		col.AutoRandom = true
		// Optional (shard_bits[, range_bits])
		if p.cur.Type == '(' {
			p.advance()
			if p.cur.Type != tokICONST {
				return false, p.syntaxErrorAtCur()
			}
			col.AutoRandomShardBits = int(p.cur.Ival)
			p.advance()
			if p.cur.Type == ',' {
				p.advance()
				if p.cur.Type != tokICONST {
					return false, p.syntaxErrorAtCur()
				}
				col.AutoRandomRangeBits = int(p.cur.Ival)
				p.advance()
			}
			if _, ok := p.match(')'); !ok {
				return false, p.syntaxErrorAtCur()
			}
		}
		return true, nil

	case kwUNIQUE:
		p.advance()
		p.match(kwKEY) // optional KEY
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrUnique,
		})
		return true, nil

	case kwPRIMARY:
		p.advance()
		if _, err := p.expect(kwKEY); err != nil {
			return false, err
		}
		cc := &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrPrimaryKey,
		}
		if cl := p.parseClusteredAttr(); cl != nil {
			cc.Clustered = cl
			cc.Loc.End = p.pos()
		}
		col.Constraints = append(col.Constraints, cc)
		return true, nil

	case kwKEY:
		// Bare KEY means PRIMARY KEY
		p.advance()
		cc := &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrPrimaryKey,
		}
		if cl := p.parseClusteredAttr(); cl != nil {
			cc.Clustered = cl
			cc.Loc.End = p.pos()
		}
		col.Constraints = append(col.Constraints, cc)
		return true, nil

	case kwCOMMENT:
		p.advance()
		if p.cur.Type == tokSCONST {
			col.Comment = p.cur.Str
			p.advance()
		}
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
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
			Loc:  nodes.Loc{Start: start, End: p.pos()},
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
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrVisible,
		})
		return true, nil

	case kwINVISIBLE:
		p.advance()
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
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
			Loc:  nodes.Loc{Start: start, End: p.pos()},
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
		cc.Loc.End = p.pos()
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
			Loc:  nodes.Loc{Start: start, End: p.pos()},
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
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrStorage,
			Name: val,
		})
		return true, nil
	}

	// Handle ENGINE_ATTRIBUTE, SECONDARY_ENGINE_ATTRIBUTE column options
	if p.cur.Type == kwENGINE_ATTRIBUTE {
		p.advance()
		p.match('=')
		val := ""
		if p.cur.Type == tokSCONST {
			val = p.cur.Str
			p.advance()
		}
		col.EngineAttribute = val
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrEngineAttribute,
			Name: val,
		})
		return true, nil
	}
	if p.cur.Type == kwSECONDARY_ENGINE_ATTRIBUTE {
		p.advance()
		p.match('=')
		val := ""
		if p.cur.Type == tokSCONST {
			val = p.cur.Str
			p.advance()
		}
		col.SecondaryEngineAttribute = val
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrSecondaryEngineAttribute,
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
		return &nodes.ParenExpr{Loc: nodes.Loc{Start: start, End: p.pos()}, Expr: expr}, nil
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
		if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
			// FULL, PARTIAL, SIMPLE
			constr.Match, _, _ = p.parseKeywordOrIdent()
		}
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

	constr.Loc.End = p.pos()
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
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Expr: expr,
	}

	// VIRTUAL or STORED
	if _, ok := p.match(kwVIRTUAL); ok {
		gen.Stored = false
	} else if _, ok := p.match(kwSTORED); ok {
		gen.Stored = true
	}

	gen.Loc.End = p.pos()
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
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Expr: expr,
	}

	// VIRTUAL or STORED
	if _, ok := p.match(kwVIRTUAL); ok {
		gen.Stored = false
	} else if _, ok := p.match(kwSTORED); ok {
		gen.Stored = true
	}

	gen.Loc.End = p.pos()
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
		p.parseIndexTypeClause(constr)
		idxCols, err := p.parseParenIndexKeyParts()
		if err != nil {
			return nil, err
		}
		constr.IndexColumns = idxCols
		constr.Columns = indexColumnsToNames(idxCols)
		p.parseConstraintIndexOptions(constr)
		constr.Clustered = p.parseClusteredAttr()

	case kwUNIQUE:
		p.advance()
		p.match(kwINDEX, kwKEY) // optional INDEX or KEY
		constr.Type = nodes.ConstrUnique
		// Optional index name
		if p.isIdentToken() && p.cur.Type != '(' {
			constr.Name, _, _ = p.parseIdent()
		}
		// Optional index type
		p.parseIndexTypeClause(constr)
		idxCols, err := p.parseParenIndexKeyParts()
		if err != nil {
			return nil, err
		}
		constr.IndexColumns = idxCols
		constr.Columns = indexColumnsToNames(idxCols)
		p.parseConstraintIndexOptions(constr)

	case kwFOREIGN:
		p.advance()
		p.match(kwKEY)
		constr.Type = nodes.ConstrForeignKey
		// Optional index name
		if p.isIdentToken() && p.cur.Type != '(' {
			constr.Name, _, _ = p.parseIdent()
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
				if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
					constr.Match, _, _ = p.parseKeywordOrIdent()
				}
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
		p.parseIndexTypeClause(constr)
		idxCols, err := p.parseParenIndexKeyParts()
		if err != nil {
			return nil, err
		}
		constr.IndexColumns = idxCols
		constr.Columns = indexColumnsToNames(idxCols)
		p.parseConstraintIndexOptions(constr)

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
		p.parseConstraintIndexOptions(constr)

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
		p.parseConstraintIndexOptions(constr)

	default:
		return nil, &ParseError{
			Message:  "expected constraint definition",
			Position: p.cur.Loc,
		}
	}

	constr.Loc.End = p.pos()
	return constr, nil
}

// parseConstraintIndexOptions parses index_option list after a table constraint's key parts.
func (p *Parser) parseConstraintIndexOptions(constr *nodes.Constraint) {
	for {
		opt, ok, err := p.parseIndexOption()
		if err != nil || !ok {
			break
		}
		constr.IndexOptions = append(constr.IndexOptions, opt)
	}
}

// parseIndexTypeClause parses an optional USING {BTREE|HASH} clause.
func (p *Parser) parseIndexTypeClause(constr *nodes.Constraint) {
	if p.cur.Type == kwUSING {
		p.advance()
		if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
			constr.IndexType, _, _ = p.parseKeywordOrIdent()
		}
	}
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
// parseClusteredAttr parses optional CLUSTERED/NONCLUSTERED after PRIMARY KEY.
// Returns nil if neither keyword is present.
func (p *Parser) parseClusteredAttr() *bool {
	if p.cur.Type == kwCLUSTERED {
		p.advance()
		b := true
		return &b
	} else if p.cur.Type == kwNONCLUSTERED {
		p.advance()
		b := false
		return &b
	}
	return nil
}

func (p *Parser) parseTableOption() (*nodes.TableOption, bool, error) {
	// Completion: offer table option keywords.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{
			kwENGINE, kwDEFAULT, kwCHARSET, kwCOLLATE, kwCOMMENT,
			kwAUTO_INCREMENT, kwROW_FORMAT, kwPARTITION,
			// TiDB-specific table options (PR2 keywords; exposed here
			// so completion offers them in the table_option list).
			kwSHARD_ROW_ID_BITS, kwPRE_SPLIT_REGIONS, kwAUTO_ID_CACHE,
			kwAUTO_RANDOM_BASE, kwTTL, kwTTL_ENABLE, kwTTL_JOB_INTERVAL,
			kwPLACEMENT,
		} {
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
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "ENGINE", Value: val}, true, nil

	case kwAUTO_INCREMENT:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "AUTO_INCREMENT", Value: val}, true, nil

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
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "CHARSET", Value: val}, true, nil

	case kwCHARACTER:
		p.advance()
		if _, ok := p.match(kwSET); ok {
			p.match('=')
			val, err := p.consumeOptionValue()
			if err != nil {
				return nil, false, err
			}
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "CHARACTER SET", Value: val}, true, nil
		}
		return nil, false, nil

	case kwCOLLATE:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "COLLATE", Value: val}, true, nil

	case kwCOMMENT:
		p.advance()
		p.match('=')
		val := ""
		if p.cur.Type == tokSCONST {
			val = p.cur.Str
			p.advance()
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "COMMENT", Value: val}, true, nil

	case kwROW_FORMAT:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "ROW_FORMAT", Value: val}, true, nil

	case kwCONNECTION:
		p.advance()
		p.match('=')
		val := ""
		if p.cur.Type == tokSCONST {
			val = p.cur.Str
			p.advance()
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "CONNECTION", Value: val}, true, nil

	case kwPASSWORD:
		p.advance()
		p.match('=')
		val := ""
		if p.cur.Type == tokSCONST {
			val = p.cur.Str
			p.advance()
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "PASSWORD", Value: val}, true, nil

	case kwSTORAGE:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "STORAGE", Value: val}, true, nil

	case kwSTART:
		// START TRANSACTION (used for versioned tables)
		if next := p.peekNext(); next.Type == kwTRANSACTION {
			p.advance() // consume START
			p.advance() // consume TRANSACTION
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "START TRANSACTION", Value: ""}, true, nil
		}
		return nil, false, nil

	case kwCHECKSUM:
		// CHECKSUM [=] {0 | 1}
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "CHECKSUM", Value: val}, true, nil

	case kwTABLESPACE:
		// TABLESPACE tablespace_name [STORAGE {DISK | MEMORY}]
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		opt := &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "TABLESPACE", Value: val}
		if p.cur.Type == kwSTORAGE {
			p.advance()
			storageVal, err := p.consumeOptionValue()
			if err != nil {
				return nil, false, err
			}
			opt.Storage = storageVal
		}
		opt.Loc.End = p.pos()
		return opt, true, nil

	case kwENCRYPTION:
		// ENCRYPTION [=] {'Y' | 'N'}
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "ENCRYPTION", Value: val}, true, nil

	case kwUNION:
		// UNION [=] (tbl_name[,tbl_name]...)
		p.advance()
		p.match('=')
		if p.cur.Type == '(' {
			names, err := p.parseParenIdentList()
			if err == nil {
				val := strings.Join(names, ",")
				return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "UNION", Value: val}, true, nil
			}
		}
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "UNION", Value: val}, true, nil

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
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "INDEX DIRECTORY", Value: val}, true, nil
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
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "DATA DIRECTORY", Value: val}, true, nil
		}
		return nil, false, nil

	case kwSHARD_ROW_ID_BITS:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "SHARD_ROW_ID_BITS", Value: val}, true, nil

	case kwPRE_SPLIT_REGIONS:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "PRE_SPLIT_REGIONS", Value: val}, true, nil

	case kwAUTO_ID_CACHE:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "AUTO_ID_CACHE", Value: val}, true, nil

	case kwAUTO_RANDOM_BASE:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "AUTO_RANDOM_BASE", Value: val}, true, nil

	case kwTTL_ENABLE:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "TTL_ENABLE", Value: val}, true, nil

	case kwTTL_JOB_INTERVAL:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "TTL_JOB_INTERVAL", Value: val}, true, nil

	case kwPLACEMENT:
		p.advance()
		if _, ok := p.match(kwPOLICY); !ok {
			return nil, false, p.syntaxErrorAtCur()
		}
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "PLACEMENT POLICY", Value: val}, true, nil

	case kwTTL:
		// TTL = expr (e.g., created_at + INTERVAL 1 YEAR, or DATE_ADD(col, INTERVAL 1 DAY))
		// Parse as a full expression to handle function calls, parens, string literals correctly.
		// Store raw source text in Value for round-trip fidelity.
		p.advance()
		p.match('=')
		exprStart := p.cur.Loc
		_, err := p.parseExpr()
		if err != nil {
			return nil, false, err
		}
		exprEnd := p.prev.End
		val := strings.TrimSpace(p.lexer.input[exprStart:exprEnd])
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "TTL", Value: val}, true, nil
	}

	// Handle keyword-token-based table options: COMPRESSION, INSERT_METHOD.
	if p.cur.Type == kwCOMPRESSION || p.cur.Type == kwINSERT_METHOD {
		optName := p.cur.Str
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: optName, Value: val}, true, nil
	}

	// Handle identifier-based options: KEY_BLOCK_SIZE, STATS_AUTO_RECALC, etc.
	switch p.cur.Type {
	case kwKEY_BLOCK_SIZE, kwSTATS_AUTO_RECALC, kwSTATS_PERSISTENT,
		kwSTATS_SAMPLE_PAGES, kwMAX_ROWS, kwMIN_ROWS, kwAVG_ROW_LENGTH,
		kwPACK_KEYS, kwDELAY_KEY_WRITE, kwSECONDARY_ENGINE,
		kwSECONDARY_ENGINE_ATTRIBUTE, kwAUTOEXTEND_SIZE, kwENGINE_ATTRIBUTE:
		optName := p.cur.Str
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: optName, Value: val}, true, nil
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
		if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
			name, _, err := p.parseKeywordOrIdent()
			if err != nil {
				return "", err
			}
			return name, nil
		}
		return "", nil
	}
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
			expr, err := p.parseExpr()
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
				expr, err := p.parseExpr()
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
			if p.cur.Type == ')' {
				break
			}
			pd, err := p.parsePartitionDef()
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

	part.Loc.End = p.pos()
	return part, nil
}

// parsePartitionDef parses a single partition definition.
//
//	PARTITION partition_name
//	    [VALUES {LESS THAN (expr|MAXVALUE) | IN (value_list)}]
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

	// VALUES LESS THAN or VALUES IN
	if _, ok := p.match(kwVALUES); ok {
		if p.cur.Type == kwLESS {
			p.advance() // LESS
			if p.cur.Type == kwTHAN {
				p.advance() // THAN
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
				p.match(')')
				if allMax && len(items) > 0 {
					// All MAXVALUE — store as single string marker
					pd.Values = &nodes.String{Str: "MAXVALUE"}
				} else if len(items) == 1 {
					pd.Values = items[0]
				} else {
					pd.Values = &nodes.List{Items: items}
				}
			} else if p.cur.Type == kwMAXVALUE {
				pd.Values = &nodes.String{Str: "MAXVALUE"}
				p.advance()
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
							expr, err := p.parseExpr()
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
						expr, err := p.parseExpr()
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
				p.match(')')
				pd.Values = &nodes.List{Items: items}
			}
		}
	}

	// Optional table options for partition
	for {
		opt, ok, err := p.parseTableOption()
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
			if p.cur.Type == ')' {
				break
			}
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

	pd.Loc.End = p.pos()
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
		opt, ok, err := p.parseTableOption()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		spd.Options = append(spd.Options, opt)
	}
	spd.Loc.End = p.pos()
	return spd, nil
}
