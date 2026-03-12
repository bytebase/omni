package parser

import (
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
		p.match(kwNOT)
		p.match(kwEXISTS_KW)
		stmt.IfNotExists = true
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
		// Check for CREATE TABLE ... (SELECT ...) — subquery without AS
		if next := p.peekNext(); next.Type == kwSELECT {
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
		opt, ok := p.parseTableOption()
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
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	col := &nodes.ColumnDef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
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
		if !p.parseColumnOption(col) {
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
	return false
}

// parseColumnOption parses a single column option/constraint.
// Returns false if no option was found.
func (p *Parser) parseColumnOption(col *nodes.ColumnDef) bool {
	start := p.pos()

	switch p.cur.Type {
	case kwNOT:
		p.advance()
		p.match(kwNULL)
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrNotNull,
		})
		return true

	case kwNULL:
		p.advance()
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrNull,
		})
		return true

	case kwDEFAULT:
		p.advance()
		expr, err := p.parseDefaultValue()
		if err != nil {
			return false
		}
		col.DefaultValue = expr
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrDefault,
			Expr: expr,
		})
		return true

	case kwAUTO_INCREMENT:
		p.advance()
		col.AutoIncrement = true
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrAutoIncrement,
		})
		return true

	case kwUNIQUE:
		p.advance()
		p.match(kwKEY) // optional KEY
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrUnique,
		})
		return true

	case kwPRIMARY:
		p.advance()
		p.match(kwKEY)
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrPrimaryKey,
		})
		return true

	case kwKEY:
		// Bare KEY means PRIMARY KEY
		p.advance()
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrPrimaryKey,
		})
		return true

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
		return true

	case kwCOLLATE:
		p.advance()
		if p.isIdentToken() {
			collName, _, _ := p.parseIdentifier()
			col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
				Loc:  nodes.Loc{Start: start, End: p.pos()},
				Type: nodes.ColConstrCollate,
				Name: collName,
			})
		}
		return true

	case kwREFERENCES:
		constr, err := p.parseReferenceDefinition(start)
		if err != nil {
			return false
		}
		col.Constraints = append(col.Constraints, constr)
		return true

	case kwVISIBLE:
		p.advance()
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrVisible,
		})
		return true

	case kwINVISIBLE:
		p.advance()
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrInvisible,
		})
		return true

	case kwCHECK:
		p.advance()
		if _, err := p.expect('('); err != nil {
			return false
		}
		expr, err := p.parseExpr()
		if err != nil {
			return false
		}
		if _, err := p.expect(')'); err != nil {
			return false
		}
		cc := &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrCheck,
			Expr: expr,
		}
		// [[NOT] ENFORCED]
		if p.cur.Type == kwNOT {
			p.advance()
			if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "enforced") {
				p.advance()
				cc.NotEnforced = true
			}
		} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "enforced") {
			p.advance()
		}
		cc.Loc.End = p.pos()
		col.Constraints = append(col.Constraints, cc)
		return true

	case kwGENERATED:
		gen, err := p.parseGeneratedColumn()
		if err != nil {
			return false
		}
		col.Generated = gen
		return true

	case kwON:
		// ON UPDATE CURRENT_TIMESTAMP
		if next := p.peekNext(); next.Type == kwUPDATE {
			p.advance() // consume ON
			p.advance() // consume UPDATE
			expr, err := p.parseExpr()
			if err != nil {
				return false
			}
			col.OnUpdate = expr
			return true
		}
		return false

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
			return false
		}
		col.ColumnFormat = val
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrColumnFormat,
			Name: val,
		})
		return true

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
			return false
		}
		col.Storage = val
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Type: nodes.ColConstrStorage,
			Name: val,
		})
		return true
	}

	return false
}

// parseDefaultValue parses a DEFAULT value expression.
// Handles parenthesized expressions, literals, CURRENT_TIMESTAMP, NULL, etc.
func (p *Parser) parseDefaultValue() (nodes.ExprNode, error) {
	// Handle parenthesized expression: DEFAULT (expr)
	if p.cur.Type == '(' {
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return expr, nil
	}
	return p.parseExpr()
}

// parseReferenceDefinition parses a REFERENCES clause.
//
//	REFERENCES tbl_name (col,...) [MATCH FULL|PARTIAL|SIMPLE]
//	    [ON DELETE action] [ON UPDATE action]
//
//	action: RESTRICT | CASCADE | SET NULL | SET DEFAULT | NO ACTION
func (p *Parser) parseReferenceDefinition(start int) (*nodes.ColumnConstraint, error) {
	p.advance() // consume REFERENCES

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
		cols, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		constr.RefColumns = cols
	}

	// MATCH
	if _, ok := p.match(kwMATCH); ok {
		if p.isIdentToken() {
			// FULL, PARTIAL, SIMPLE
			constr.Match, _, _ = p.parseIdentifier()
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
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "action") {
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
	if p.cur.Type == kwALWAYS || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "always")) {
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
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
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
			constr.Name, _, _ = p.parseIdentifier()
		}
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance()
		p.match(kwKEY)
		constr.Type = nodes.ConstrPrimaryKey
		// Optional index type
		p.parseIndexTypeClause(constr)
		idxCols, err := p.parseParenIndexKeyParts()
		if err != nil {
			return nil, err
		}
		constr.IndexColumns = idxCols
		constr.Columns = indexColumnsToNames(idxCols)
		p.parseConstraintIndexOptions(constr)

	case kwUNIQUE:
		p.advance()
		p.match(kwINDEX, kwKEY) // optional INDEX or KEY
		constr.Type = nodes.ConstrUnique
		// Optional index name
		if p.isIdentToken() && p.cur.Type != '(' {
			constr.Name, _, _ = p.parseIdentifier()
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
			constr.Name, _, _ = p.parseIdentifier()
		}
		cols, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		constr.Columns = cols
		// REFERENCES
		if p.cur.Type == kwREFERENCES {
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			constr.RefTable = ref
			if p.cur.Type == '(' {
				refCols, err := p.parseParenIdentList()
				if err != nil {
					return nil, err
				}
				constr.RefColumns = refCols
			}
			// MATCH
			if _, ok := p.match(kwMATCH); ok {
				if p.isIdentToken() {
					constr.Match, _, _ = p.parseIdentifier()
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
		constr.Expr = expr
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		// [[NOT] ENFORCED]
		if p.cur.Type == kwNOT {
			p.advance()
			if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "enforced") {
				p.advance()
				constr.NotEnforced = true
			}
		} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "enforced") {
			p.advance()
		}

	case kwINDEX, kwKEY:
		p.advance()
		constr.Type = nodes.ConstrIndex
		// Optional index name
		if p.isIdentToken() && p.cur.Type != '(' {
			constr.Name, _, _ = p.parseIdentifier()
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
			constr.Name, _, _ = p.parseIdentifier()
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
			constr.Name, _, _ = p.parseIdentifier()
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
		if p.isIdentToken() {
			constr.IndexType, _, _ = p.parseIdentifier()
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
func (p *Parser) parseTableOption() (*nodes.TableOption, bool) {
	start := p.pos()

	switch p.cur.Type {
	case kwENGINE:
		p.advance()
		p.match('=')
		val := p.consumeOptionValue()
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "ENGINE", Value: val}, true

	case kwAUTO_INCREMENT:
		p.advance()
		p.match('=')
		val := p.consumeOptionValue()
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "AUTO_INCREMENT", Value: val}, true

	case kwDEFAULT:
		// DEFAULT CHARSET or DEFAULT CHARACTER SET or DEFAULT COLLATE
		next := p.peekNext()
		if next.Type == kwCHARSET || next.Type == kwCHARACTER || next.Type == kwCOLLATE {
			p.advance() // consume DEFAULT
			return p.parseTableOption()
		}
		return nil, false

	case kwCHARSET:
		p.advance()
		p.match('=')
		val := p.consumeOptionValue()
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "CHARSET", Value: val}, true

	case kwCHARACTER:
		p.advance()
		if _, ok := p.match(kwSET); ok {
			p.match('=')
			val := p.consumeOptionValue()
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "CHARACTER SET", Value: val}, true
		}
		return nil, false

	case kwCOLLATE:
		p.advance()
		p.match('=')
		val := p.consumeOptionValue()
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "COLLATE", Value: val}, true

	case kwCOMMENT:
		p.advance()
		p.match('=')
		val := ""
		if p.cur.Type == tokSCONST {
			val = p.cur.Str
			p.advance()
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "COMMENT", Value: val}, true

	case kwROW_FORMAT:
		p.advance()
		p.match('=')
		val := p.consumeOptionValue()
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "ROW_FORMAT", Value: val}, true

	case kwCONNECTION:
		p.advance()
		p.match('=')
		val := ""
		if p.cur.Type == tokSCONST {
			val = p.cur.Str
			p.advance()
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "CONNECTION", Value: val}, true

	case kwPASSWORD:
		p.advance()
		p.match('=')
		val := ""
		if p.cur.Type == tokSCONST {
			val = p.cur.Str
			p.advance()
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "PASSWORD", Value: val}, true

	case kwSTORAGE:
		p.advance()
		p.match('=')
		val := p.consumeOptionValue()
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "STORAGE", Value: val}, true
	}

	// Handle identifier-based options: KEY_BLOCK_SIZE, STATS_AUTO_RECALC, etc.
	if p.cur.Type == tokIDENT {
		optName := p.cur.Str
		switch {
		case eqFold(optName, "key_block_size"),
			eqFold(optName, "stats_auto_recalc"),
			eqFold(optName, "stats_persistent"),
			eqFold(optName, "stats_sample_pages"),
			eqFold(optName, "max_rows"),
			eqFold(optName, "min_rows"),
			eqFold(optName, "avg_row_length"),
			eqFold(optName, "pack_keys"),
			eqFold(optName, "checksum"),
			eqFold(optName, "delay_key_write"),
			eqFold(optName, "compression"),
			eqFold(optName, "encryption"),
			eqFold(optName, "insert_method"),
			eqFold(optName, "data_directory"),
			eqFold(optName, "index_directory"),
			eqFold(optName, "tablespace"),
			eqFold(optName, "union"),
			eqFold(optName, "secondary_engine"),
			eqFold(optName, "secondary_engine_attribute"),
			eqFold(optName, "autoextend_size"),
			eqFold(optName, "engine_attribute"):
			p.advance()
			p.match('=')
			val := p.consumeOptionValue()
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: optName, Value: val}, true
		}
	}

	return nil, false
}

// consumeOptionValue reads the next token as an option value string.
func (p *Parser) consumeOptionValue() string {
	switch p.cur.Type {
	case tokSCONST:
		val := p.cur.Str
		p.advance()
		return val
	case tokICONST:
		val := p.cur.Str
		if val == "" {
			// Ival is set but Str may be empty for integer tokens
			val = string(rune('0' + p.cur.Ival))
		}
		p.advance()
		return val
	default:
		if p.isIdentToken() {
			name, _, _ := p.parseIdentifier()
			return name
		}
		return ""
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
			cols, err := p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
			part.Columns = cols
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
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "partitions") {
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
			cols, err := p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
			part.SubPartColumns = cols
		case p.cur.Type == tokIDENT && eqFold(p.cur.Str, "linear"):
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
				part.SubPartExpr = expr
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
			} else if p.cur.Type == kwKEY {
				p.advance()
				part.SubPartType = nodes.PartitionKey
				if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "algorithm") {
					p.advance()
					p.match('=')
					if p.cur.Type == tokICONST {
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
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "subpartitions") {
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

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	pd := &nodes.PartitionDef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// VALUES LESS THAN or VALUES IN
	if _, ok := p.match(kwVALUES); ok {
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "less") {
			p.advance() // LESS
			if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "than") {
				p.advance() // THAN
			}
			if p.cur.Type == '(' {
				p.advance()
				if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "maxvalue") {
					pd.Values = &nodes.String{Str: "MAXVALUE"}
					p.advance()
				} else {
					expr, err := p.parseExpr()
					if err != nil {
						return nil, err
					}
					pd.Values = expr
				}
				p.match(')')
			} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "maxvalue") {
				pd.Values = &nodes.String{Str: "MAXVALUE"}
				p.advance()
			}
		} else if _, ok := p.match(kwIN); ok {
			if p.cur.Type == '(' {
				p.advance()
				var vals []nodes.ExprNode
				for {
					expr, err := p.parseExpr()
					if err != nil {
						return nil, err
					}
					vals = append(vals, expr)
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
				p.match(')')
				// Store as List node
				items := make([]nodes.Node, len(vals))
				for i, v := range vals {
					items[i] = v
				}
				pd.Values = &nodes.List{Items: items}
			}
		}
	}

	// Optional table options for partition
	for {
		opt, ok := p.parseTableOption()
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
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	spd := &nodes.SubPartitionDef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}
	for {
		opt, ok := p.parseTableOption()
		if !ok {
			break
		}
		spd.Options = append(spd.Options, opt)
	}
	spd.Loc.End = p.pos()
	return spd, nil
}
