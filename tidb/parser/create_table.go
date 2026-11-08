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
//
// TiDB adds GLOBAL TEMPORARY, which ends in ON COMMIT and is the one spelling
// that does: see parseOnCommit. It also lets AS stand with no query after it:
// see parseCreateTableQueryOpt.
func (p *Parser) parseCreateTableStmt(start int, temporary, global bool) (*nodes.CreateTableStmt, error) {
	p.advance() // consume TABLE

	stmt := &nodes.CreateTableStmt{
		Loc:       nodes.Loc{Start: start},
		Temporary: temporary,
		Global:    global,
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
		if err := p.parseOnCommit(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// CREATE TABLE ... AS [query]
	if p.cur.Type == kwAS {
		p.advance()
		if err := p.parseCreateTableQueryOpt(stmt); err != nil {
			return nil, err
		}
		if err := p.parseOnCommit(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Parenthesized create definitions (columns and constraints)
	if p.cur.Type == '(' {
		next := p.peekNext()
		// A '(' right after the table name opens the column list and never a
		// query: TiDB v8.5.5 parser.y:4525 (CreateTableStmt) reads
		// TableElementListOpt there, so "CREATE TABLE t (SELECT 1)" is 1064 while
		// "CREATE TABLE t (a INT) (SELECT 1)" and "CREATE TABLE t AS (SELECT 1)"
		// parse.
		//
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
			if err := p.parseOnCommit(stmt); err != nil {
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

	// CREATE TABLE t (cols) AS [query]
	if p.cur.Type == kwAS {
		p.advance()
		if err := p.parseCreateTableQueryOpt(stmt); err != nil {
			return nil, err
		}
		if err := p.parseOnCommit(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Table options. A comma may stand between two of them and never ends the
	// list (TableOptionList in TiDB v8.5.5 parser.y:12811), so a comma with no
	// option after it is a syntax error rather than the end of the options:
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

	// IGNORE / REPLACE before SELECT in CTAS
	if p.cur.Type == kwIGNORE {
		p.advance()
		stmt.Ignore = true
	} else if p.cur.Type == kwREPLACE {
		p.advance()
		stmt.Replace = true
	}

	// [AS] [query] after the columns, options and partitioning: AsOpt
	// CreateTableSelectOpt in TiDB v8.5.5 parser.y:4525 (CreateTableStmt), see
	// parseCreateTableQueryOpt. Reading only SELECT here used to end the
	// statement in front of "TABLE u" or "(SELECT 1)", which then parsed as a
	// second statement.
	p.match(kwAS)
	if err := p.parseCreateTableQueryOpt(stmt); err != nil {
		return nil, err
	}

	if err := p.parseOnCommit(stmt); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateTableQueryOpt reads the query of CREATE TABLE ... [AS] query, if
// one starts here, into stmt.Select: CreateTableSelectOpt in TiDB v8.5.5
// parser.y:5031, whose first arm is empty.
//
// AS does not call for a query either (AsOpt, parser.y:5026), so AS may end
// the statement or stand in front of the ON COMMIT clause that ends a GLOBAL
// TEMPORARY table. pingcap/tidb:v8.5.5 parses "CREATE TABLE t (a INT) AS",
// "CREATE TABLE t (a INT) IGNORE AS", "CREATE TABLE t AS" and "CREATE GLOBAL
// TEMPORARY TABLE t (a INT) AS ON COMMIT DELETE ROWS" (measured through
// PREPARE on 2026-11-05). Run, the first creates t with its one column, and
// "CREATE TABLE t AS" fails with 1113, a table needs at least one column,
// which is the server's to say rather than a syntax error. A word after AS
// that opens no query is a syntax error: "CREATE TABLE t AS IGNORE" is 1064.
func (p *Parser) parseCreateTableQueryOpt(stmt *nodes.CreateTableStmt) error {
	if !p.isCreateQueryStart() {
		return nil
	}
	sel, err := p.parseSelectStmt()
	if err != nil {
		return err
	}
	stmt.Select = sel
	return nil
}

// isCreateQueryStart reports whether the current token opens the query of
// CREATE TABLE ... [AS] query after the column list, the table options or the
// partitioning: CreateTableSelectOpt in TiDB v8.5.5 parser.y:5031, which is a
// SELECT, a WITH query, TABLE, VALUES, or a parenthesized query (SubSelect).
func (p *Parser) isCreateQueryStart() bool {
	switch p.cur.Type {
	case kwSELECT, kwWITH, kwTABLE, kwVALUES, '(':
		return true
	}
	return false
}

// parseOnCommit reads the ON COMMIT clause a GLOBAL TEMPORARY table ends with,
// and holds every other table to not having one.
//
//	ON COMMIT {DELETE | PRESERVE} ROWS
//
// TiDB's grammar hangs this clause on GLOBAL TEMPORARY alone and makes it
// mandatory there (parser.y, rule CreateTableStmt): v8.5.5 answers 1064 for
// `create global temporary table t (a int)` with no clause, and for
// `create temporary table t (a int) on commit delete rows` and a plain
// `create table` with one. It comes last, after the options and after the
// partition clause -- `... comment 'x' on commit delete rows` parses and
// `... on commit delete rows comment 'x'` does not -- which is why this runs at
// the end of every path through parseCreateTableStmt, the LIKE one included:
// TiDB takes `create global temporary table t like t1 on commit delete rows`.
//
// PRESERVE ROWS parses and the server then refuses it with 8200, "TiDB doesn't
// support ON COMMIT PRESERVE ROWS for now". That is a refusal and not a syntax
// error, so the clause is taken here and the engine answers for it.
func (p *Parser) parseOnCommit(stmt *nodes.CreateTableStmt) error {
	if p.cur.Type != kwON || p.peekNext().Type != kwCOMMIT {
		if stmt.Global {
			return p.syntaxErrorAtCur()
		}
		return nil
	}
	if !stmt.Global {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume ON
	p.advance() // consume COMMIT
	var action string
	switch p.cur.Type {
	case kwDELETE:
		action = "DELETE"
	case kwPRESERVE:
		action = "PRESERVE"
	default:
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if _, err := p.expect(kwROWS); err != nil {
		return err
	}
	stmt.OnCommit = action + " ROWS"
	return nil
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
	case kwPRIMARY, kwUNIQUE, kwINDEX, kwKEY, kwFULLTEXT, kwFOREIGN, kwCHECK:
		// No SPATIAL: TiDB's ConstraintElem has no spatial key, and v8.5.5
		// answers 1064 for one in a table definition and in ALTER TABLE ADD.
		return true
	case kwCONSTRAINT:
		return true
	case kwVECTOR:
		// TiDB's VECTOR INDEX. VECTOR is also a column type, so the word only
		// starts a constraint when INDEX follows it: `a VECTOR(3)` is a column
		// and `VECTOR INDEX (...)` a constraint. v8.5.5 answers 1064 for
		// `VECTOR KEY (...)`, so INDEX is the only continuation.
		return p.peekNext().Type == kwINDEX
	}
	return false
}

// parseReferencedColumns reads the column list of REFERENCES, whose parts are ColumnNames. The
// names it keeps are the columns', as the qualifier names nothing the foreign key uses.
func (p *Parser) parseReferencedColumns() ([]string, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var names []string
	for {
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, &ParseError{Message: "collecting"}
		}
		_, _, name, err := p.parseColumnName()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
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

// parseColumnName reads ColumnName, parser.y:3511: an Identifier, with up to two more before
// it. So `a`, `t.a` and `d.t.a` are column names and `c.d.t.a` is a 1064. What a qualifier must
// match is the statement's to say: CREATE TABLE, MODIFY and CHANGE check it, and ADD, ALTER and
// DROP ignore it, which the catalog holds.
func (p *Parser) parseColumnName() (schema, table, name string, err error) {
	name, _, err = p.parseIdent()
	if err != nil {
		return "", "", "", err
	}
	if p.cur.Type != '.' {
		return "", "", name, nil
	}
	p.advance()
	table = name
	name, _, err = p.parseIdent()
	if err != nil {
		return "", "", "", err
	}
	if p.cur.Type != '.' {
		return "", table, name, nil
	}
	p.advance()
	schema = table
	table = name
	name, _, err = p.parseIdent()
	if err != nil {
		return "", "", "", err
	}
	return schema, table, name, nil
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
	schema, table, name, err := p.parseColumnName()
	if err != nil {
		return nil, err
	}

	col := &nodes.ColumnDef{
		Loc:    nodes.Loc{Start: start},
		Schema: schema,
		Table:  table,
		Name:   name,
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

	// ColumnDef's action calls Validate, which answers 1221 for these three beside a
	// generated column in either order.
	//
	// Ref: invalidOptionForGeneratedColumn and ColumnDef.Validate in TiDB v8.5.5
	// pkg/parser/ast/ddl.go.
	if col.Generated != nil && (col.DefaultValue != nil || col.AutoIncrement || col.OnUpdate != nil) {
		return nil, &ParseError{
			Message:  "a generated column takes no DEFAULT, AUTO_INCREMENT or ON UPDATE",
			Position: start,
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
		for _, t := range []int{
			kwNOT, kwNULL, kwDEFAULT, kwAUTO_INCREMENT, kwPRIMARY, kwUNIQUE,
			kwCOMMENT, kwCOLLATE, kwREFERENCES, kwCHECK, kwGENERATED,
			// TiDB: AUTO_RANDOM belongs on the column constraint list
			// alongside AUTO_INCREMENT (same grammar position).
			kwAUTO_RANDOM,
		} {
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
	// builds. All three words are required: TiDB answers 1064 for
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
		cc := &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrUnique,
		}
		p.parseColumnKeyGlobalOrLocal()
		cc.Loc.End = p.prev.End
		col.Constraints = append(col.Constraints, cc)
		return true, nil

	case kwPRIMARY:
		p.advance()
		if _, err := p.expect(kwKEY); err != nil {
			return false, err
		}
		cc := &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrPrimaryKey,
		}
		if cl := p.parseClusteredAttr(); cl != nil {
			cc.Clustered = cl
		}
		p.parseColumnKeyGlobalOrLocal()
		cc.Loc.End = p.prev.End
		col.Constraints = append(col.Constraints, cc)
		return true, nil

	case kwKEY:
		// Bare KEY means PRIMARY KEY
		p.advance()
		cc := &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrPrimaryKey,
		}
		if cl := p.parseClusteredAttr(); cl != nil {
			cc.Clustered = cl
		}
		p.parseColumnKeyGlobalOrLocal()
		cc.Loc.End = p.prev.End
		col.Constraints = append(col.Constraints, cc)
		return true, nil

	case kwCONSTRAINT:
		// CONSTRAINT [symbol] CHECK (expr) [[NOT] ENFORCED], the column-level
		// check the CONSTRAINT keyword may name. The symbol is optional, and
		// CHECK is the only continuation: TiDB v8.5.5 and MySQL 8.0.46 both
		// run `ADD a INT CONSTRAINT c CHECK (a > 1) NOT ENFORCED` and
		// `ADD a INT CONSTRAINT CHECK (a > 1)`, and both answer 1064 for
		// `ADD a INT CONSTRAINT c NOT NULL` and for a bare
		// `ADD a INT CONSTRAINT c`. MariaDB 11.8.8 answers 1064 for the whole
		// form, so this is not ported to that fork.
		//
		// Ref: `ColumnOption: ConstraintKeywordOpt "CHECK" '(' Expression ')'
		// EnforcedOrNotOrNotNullOpt` in TiDB v8.5.5 pkg/parser/parser.y.
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
		// ColumnOption "COLLATE" CollationName (parser.y:3760), where
		// CollationName is StringName | binaryType: an identifier, a string
		// literal, or the reserved word BINARY. `c JSON COLLATE BINARY` parses
		// on TiDB v8.5.5.
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

	// No VISIBLE or INVISIBLE: TiDB v8.5.5 has no column visibility and answers
	// 1064 for either word after a column type, generated or not. They are index
	// options only.
	//
	// Ref: ColumnOption in TiDB v8.5.5 pkg/parser/parser.y, which has no
	// "VISIBLE" or "INVISIBLE" arm; IndexInvisible is used by index options and
	// `ALTER INDEX` alone.

	case kwCHECK:
		cc, err := p.parseColumnCheckConstraint(start)
		if err != nil {
			return false, err
		}
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
			expr, err := p.parseOnUpdateValue()
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
		// "STORAGE" StorageMedia (parser.y:3768, v8.5.5): DEFAULT, DISK or MEMORY,
		// parsed and ignored. A bare STORAGE is a 1064.
		p.advance()
		var val string
		switch p.cur.Type {
		case kwDEFAULT:
			val = "DEFAULT"
			p.advance()
		case kwDISK:
			val = "DISK"
			p.advance()
		case kwMEMORY:
			val = "MEMORY"
			p.advance()
		default:
			return false, p.syntaxErrorAtCur()
		}
		col.Storage = val
		col.Constraints = append(col.Constraints, &nodes.ColumnConstraint{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Type: nodes.ColConstrStorage,
			Name: val,
		})
		return true, nil
	}

	// ENGINE_ATTRIBUTE and SECONDARY_ENGINE_ATTRIBUTE are not read here: TiDB
	// answers 1064 for both as a column option, as it does for both as a table
	// option, so taking them only passed statements the server refuses.

	// No SRID: TiDB has no spatial types and no SRID column attribute, and
	// v8.5.5 answers 1064 for one.

	return false, nil
}

// parseDefaultValue parses a DEFAULT value expression.
// Handles parenthesized expressions, literals, CURRENT_TIMESTAMP, NULL, etc.
func (p *Parser) parseDefaultValue() (nodes.ExprNode, error) {
	start := p.pos()
	// Handle parenthesized expression: DEFAULT (expr)
	if p.cur.Type == '(' {
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
		paren := &nodes.ParenExpr{Loc: nodes.Loc{Start: start, End: p.prev.End}, Expr: expr}
		if !validColumnDefault(paren) {
			return nil, &ParseError{Message: "Invalid default value", Position: start}
		}
		return paren, nil
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if !validColumnDefault(expr) {
		return nil, &ParseError{Message: "Invalid default value", Position: start}
	}
	return expr, nil
}

// validColumnDefault reports whether expr is a column DEFAULT the way TiDB's
// DefaultValueExpr spells it (parser.y:4052, v8.5.5): a SignedLiteral, a NowSym
// form or NEXT VALUE FOR (both FuncCallExpr here), a BuiltinFunction, or a
// function call wrapped in parentheses. A parenthesized literal, column or
// arithmetic is not one: `DEFAULT (0)`, `(b)`, `(1+1)` and `DEFAULT 1+1` are
// 1064s on v8.5.5, while `DEFAULT 0`, `-1`, `now()` and `(now())` parse.
func validColumnDefault(expr nodes.ExprNode) bool {
	switch n := expr.(type) {
	case *nodes.IntLit, *nodes.FloatLit, *nodes.StringLit, *nodes.HexLit,
		*nodes.BitLit, *nodes.NullLit, *nodes.BoolLit, *nodes.TemporalLit,
		*nodes.FuncCallExpr:
		return true
	case *nodes.UnaryExpr:
		// SignedLiteral: a numeric literal with a leading + or -.
		if n.Op != nodes.UnaryPlus && n.Op != nodes.UnaryMinus {
			return false
		}
		switch n.Operand.(type) {
		case *nodes.IntLit, *nodes.FloatLit:
			return true
		}
		return false
	case *nodes.ParenExpr:
		// ( BuiltinFunction ), with any number of enclosing parens.
		inner := n.Expr
		for {
			pe, ok := inner.(*nodes.ParenExpr)
			if !ok {
				break
			}
			inner = pe.Expr
		}
		_, ok := inner.(*nodes.FuncCallExpr)
		return ok
	}
	return false
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
		cols, err := p.parseReferencedColumns()
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
		// Optional index name. TiDB names the key after the CONSTRAINT
		// symbol when one is given (pkg/parser/parser.y, rule Constraint;
		// verified on v8.5.5).
		if p.isIdentToken() && p.cur.Type != '(' {
			name, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			if constr.Name == "" {
				constr.Name = name
			}
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
		ine, err := p.parseIfNotExists()
		if err != nil {
			return nil, err
		}
		constr.IfNotExists = ine
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
		// length and an ordering direction are both allowed. TiDB v8.5.5,
		// MySQL 8.0.46 and MariaDB 11.8.8 all run
		// `ADD FOREIGN KEY (ident, ident(1) ASC) REFERENCES t (a)`.
		//
		// Ref: `"FOREIGN" "KEY" IfNotExists IndexName '(' IndexPartSpecificationList
		// ')' ReferDef` in TiDB v8.5.5 pkg/parser/parser.y; key_list in MySQL
		// 8.0.46 sql/sql_yacc.yy.
		fkCols, err := p.parseParenIndexKeyParts()
		if err != nil {
			return nil, err
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
				refCols, err := p.parseReferencedColumns()
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
		ine, err := p.parseIfNotExists()
		if err != nil {
			return nil, err
		}
		constr.IfNotExists = ine
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

	case kwVECTOR:
		// VECTOR INDEX [IF NOT EXISTS] [index_name] [USING type]
		//     '(' key_part, ... ')' index_option...
		//
		// It is its own arm of ADD rather than a Constraint, so no CONSTRAINT
		// symbol may precede it: v8.5.5 answers 1064 for
		// `ADD CONSTRAINT c VECTOR INDEX ((lower(a)))` and for `VECTOR KEY`,
		// and runs `ADD VECTOR INDEX idx USING HNSW ((lower(a)))`,
		// `ADD VECTOR INDEX ((lower(a))) GLOBAL` and
		// `CREATE TABLE t (a vector(3), VECTOR INDEX ((VEC_COSINE_DISTANCE(a))))`.
		// The key list is read as an ordinary key-part list and then held to the
		// action's shape: v8.5.5 answers 1149 for `VECTOR INDEX (a)` and for
		// `VECTOR INDEX (a, (VEC_COSINE_DISTANCE(a)))` while it parses.
		//
		// Ref: ConstraintVectorIndex in TiDB v8.5.5 pkg/parser/parser.y.
		if constr.Name != "" {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume VECTOR
		if _, err := p.expect(kwINDEX); err != nil {
			return nil, err
		}
		constr.Type = nodes.ConstrVector
		ine, err := p.parseIfNotExists()
		if err != nil {
			return nil, err
		}
		constr.IfNotExists = ine
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
		if vectorConstraintRefused(constr) {
			return nil, refuseVectorIndexShape(constr.Loc.Start)
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
		// CLUSTERED and NONCLUSTERED are index options in the grammar, so they
		// arrive here wherever they were written, but every consumer reads the
		// primary key's clustering from Clustered. Record it there as well as
		// in the option list, and let a later option win, as the grammar's own
		// action does.
		switch opt.Name {
		case "CLUSTERED":
			v := true
			constr.Clustered = &v
		case "NONCLUSTERED":
			v := false
			constr.Clustered = &v
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
//
// parseClusteredAttr parses optional CLUSTERED/NONCLUSTERED after PRIMARY KEY.
// Returns nil if neither keyword is present.
func (p *Parser) parseClusteredAttr() *bool {
	// Completion: after a PRIMARY KEY (col, ...) grammar position,
	// TiDB accepts an optional CLUSTERED/NONCLUSTERED modifier. Offer
	// both candidates whenever we're collecting here.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwCLUSTERED)
		p.addTokenCandidate(kwNONCLUSTERED)
		// Note: don't swallow completion state with a ParseError —
		// parseClusteredAttr is called in the middle of a larger
		// constraint/option chain and the caller expects it to return
		// nil when nothing matches. The completion candidates are
		// already recorded.
		return nil
	}
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

	// ForceOpt, the FORCE that may precede AUTO_INCREMENT and
	// AUTO_RANDOM_BASE and no other table option. It is TiDB's own: v8.5.5
	// runs `ALTER TABLE t FORCE AUTO_INCREMENT 3`,
	// `FORCE AUTO_RANDOM_BASE = 3`, `AUTO_INCREMENT 30, FORCE
	// AUTO_RANDOM_BASE 40` and `CREATE TABLE t (a int) FORCE AUTO_INCREMENT
	// 3`, and answers 1064 for `FORCE ENGINE = InnoDB`, `FORCE AUTO_ID_CACHE
	// = 3` and a `FORCE AUTO_INCREMENT` with no value. MySQL 8.0.46 and
	// MariaDB 11.8.8 answer 1064 for every FORCE-prefixed option, keeping
	// FORCE as the standalone ALTER TABLE action all three also have.
	//
	// Ref: `ForceOpt "AUTO_INCREMENT" EqOpt LengthNum` and
	// `ForceOpt "AUTO_RANDOM_BASE" EqOpt LengthNum` in TiDB v8.5.5
	// pkg/parser/parser.y.
	if p.cur.Type == kwFORCE {
		switch p.peekNext().Type {
		case kwAUTO_INCREMENT, kwAUTO_RANDOM_BASE:
			p.advance() // consume FORCE
			name := strings.ToUpper(p.cur.Str)
			p.advance()
			p.match('=')
			val, err := p.consumeIntegerOptionValue()
			if err != nil {
				return nil, false, err
			}
			return &nodes.TableOption{
				Loc:   nodes.Loc{Start: start, End: p.prev.End},
				Name:  name,
				Value: val,
				Force: true,
			}, true, nil
		}
	}

	// AFFINITY [=] name is a TiDB table option. omni lexes AFFINITY as an
	// identifier, not a keyword token (it is an unreserved keyword upstream),
	// so it is matched as a bare word rather than in the type switch below.
	// Its value is a StringName — a quoted string or an identifier, not a
	// number.
	//
	// Ref: `"AFFINITY" EqOpt StringName` (TableOption) and StringName in TiDB
	// v8.5.5 pkg/parser/parser.y.
	if p.isBareWord("AFFINITY") {
		p.advance()
		p.match('=')
		val, err := p.consumeStringNameOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "AFFINITY", Value: val}, true, nil
	}

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
		// CHARACTER without SET is no option. Returning "not an option" after
		// consuming it dropped the word, so `CREATE TABLE t (a INT) CHARACTER`
		// and `ALTER TABLE t COMMENT 'x' CHARACTER` parsed; the engine answers
		// 1064 for both.
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
		val, err := p.consumeOptionValue()
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

	case kwSTORAGE:
		// STORAGE ENGINE [=] name is a table option of its own — the same one
		// ENGINE names, taking a StringName. STORAGE {MEMORY|DISK|...} on its
		// own is the storage media and takes a keyword.
		//
		// Ref: `"STORAGE" "ENGINE" EqOpt StringName` (PartDefOption) and
		// `"STORAGE" "MEMORY"` / `"STORAGE" "DISK"` (TableOption) in TiDB
		// v8.5.5 pkg/parser/parser.y.
		if p.peekNext().Type == kwENGINE {
			p.advance() // consume STORAGE
			p.advance() // consume ENGINE
			p.match('=')
			val, err := p.consumeStringNameOptionValue()
			if err != nil {
				return nil, false, err
			}
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "STORAGE ENGINE", Value: val}, true, nil
		}
		p.advance()
		p.match('=')
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
	// CHECKSUM_SYM. The value is a number and only a number: TiDB answers 1064
	// for `CHECKSUM = 'abc'`, for `= DEFAULT` and for `= -1`, and for `= 1.5`
	// as well, since its own `LengthNum` is narrower than MySQL's `ulong_num`.
	case kwCHECKSUM, kwTABLE_CHECKSUM:
		optName := "CHECKSUM"
		if p.cur.Type == kwTABLE_CHECKSUM {
			optName = "TABLE_CHECKSUM"
		}
		p.advance()
		p.match('=')
		val, err := p.consumeIntegerOptionValue()
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
		p.match('=')
		val, err := p.parseEncryptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "ENCRYPTION", Value: val}, true, nil

	case kwUNION:
		// UNION [=] '(' [tbl_name [, tbl_name]...] ')'
		//
		// The parentheses are required and the list inside them may be empty,
		// and each name may be qualified: TiDB v8.5.5, MySQL 8.0.46 and
		// MariaDB 11.8.8 all run `UNION = ()`, `UNION ()` and
		// `UNION = (d.t1, d.t2)`, and all three answer 1064 for a bare
		// `UNION`, for `UNION = (t1,)` and for `UNION = (,)`.
		//
		// Ref: `"UNION" EqOpt '(' TableNameListOpt ')'` in TiDB v8.5.5
		// pkg/parser/parser.y.
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

	case kwSHARD_ROW_ID_BITS:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "SHARD_ROW_ID_BITS", Value: val}, true, nil

	case kwPRE_SPLIT_REGIONS:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "PRE_SPLIT_REGIONS", Value: val}, true, nil

	case kwAUTO_ID_CACHE:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "AUTO_ID_CACHE", Value: val}, true, nil

	case kwAUTO_RANDOM_BASE:
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "AUTO_RANDOM_BASE", Value: val}, true, nil

	case kwTTL_ENABLE:
		// TTL_ENABLE [=] {'ON' | 'OFF'}. The value is a string literal, and
		// the grammar's own action rejects any other word with a parse error
		// rather than a later one, so this is the parser's business: TiDB
		// v8.5.5 answers 1064 for `TTL_ENABLE = 'test_case'` and for
		// `/*T![ttl] TTL_ENABLE = 'test_case' */`.
		//
		// Ref: `"TTL_ENABLE" EqOpt stringLit` in TiDB v8.5.5
		// pkg/parser/parser.y, whose action calls yylex.AppendError and
		// returns 1 unless the value is "on" or "off".
		p.advance()
		p.match('=')
		val, err := p.consumeStringOptionValue()
		if err != nil {
			return nil, false, err
		}
		if !strings.EqualFold(val, "ON") && !strings.EqualFold(val, "OFF") {
			return nil, false, &ParseError{
				Message:  "the TTL_ENABLE option has to be set 'ON' or 'OFF'",
				Position: p.prev.Loc,
			}
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "TTL_ENABLE", Value: val}, true, nil

	case kwTTL_JOB_INTERVAL:
		p.advance()
		p.match('=')
		// A string read by TiDB's own duration parser (TableOption in v8.5.5
		// pkg/parser/parser.y): `TTL_JOB_INTERVAL = 'foo'` is 1064.
		if p.cur.Type != tokSCONST {
			return nil, false, p.syntaxErrorAtCur()
		}
		if !tidbDuration(p.cur.Str) {
			return nil, false, &ParseError{Message: "The TTL_JOB_INTERVAL option is not a valid duration", Position: p.cur.Loc}
		}
		val := p.advance().Str
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "TTL_JOB_INTERVAL", Value: val}, true, nil

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
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "PLACEMENT POLICY", Value: val}, true, nil

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
		// Use inputText (not raw slicing) so that multi-statement inputs
		// — where parseSingle runs this segment with a non-zero
		// baseOffset — don't index lexer.input out of bounds. Every
		// other body-capture site in the parser uses inputText.
		val := strings.TrimSpace(p.inputText(exprStart, exprEnd))
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "TTL", Value: val}, true, nil
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
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil
	}

	// Handle identifier-based options: KEY_BLOCK_SIZE, MAX_ROWS, etc.
	//
	// ENGINE_ATTRIBUTE and SECONDARY_ENGINE_ATTRIBUTE are absent: they are MySQL
	// 8.0 options TiDB does not have, and it answers 1064 for both as a table
	// option, as a column option and on ALTER TABLE. Reading them here only
	// passed statements the server refuses.
	switch p.cur.Type {
	// TiDB's own statistics options, which MySQL does not have
	// (parser.y, TableOptionStatsBuckets and its neighbours). Each takes the
	// value its own kind: a count is an integer, the sample rate a number, and
	// the column choice and list strings. TiDB answers 1064 for the others --
	// `STATS_BUCKETS="1"`, `STATS_BUCKETS=1.5`, `STATS_COL_CHOICE=1` and the
	// bare `STATS_COL_CHOICE=all` among them.
	//
	// STATS_OPTIONS is deliberately absent: TiDB takes it on ALTER TABLE and
	// answers 1064 for it on CREATE TABLE, and this function reads the options
	// of both, so it belongs with the ALTER TABLE family and the split that
	// family needs.
	case kwSTATS_BUCKETS, kwSTATS_TOPN:
		optName := p.cur.Str
		p.advance()
		p.match('=')
		val, err := p.consumeIntegerOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil

	case kwSTATS_SAMPLE_RATE:
		optName := p.cur.Str
		p.advance()
		p.match('=')
		val, err := p.consumeNumericOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil

	case kwSTATS_COL_CHOICE, kwSTATS_COL_LIST:
		optName := p.cur.Str
		p.advance()
		p.match('=')
		val, err := p.consumeStringOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil

	// STATS_AUTO_RECALC [=] {DEFAULT | 0 | 1}, and nothing else: TiDB
	// answers 1064 for `= 2`, `= 10`, `= -1`, `= 'abc'` and `= '1'`
	// (https://dev.mysql.com/doc/refman/8.0/en/create-table.html).
	case kwSTATS_AUTO_RECALC:
		optName := p.cur.Str
		p.advance()
		p.match('=')
		val, err := p.consumeZeroOneOrDefaultOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: optName, Value: val}, true, nil

	case kwNODEGROUP:
		// NODEGROUP [=] number, a table/partition option TiDB accepts on
		// CREATE TABLE and CREATE SEQUENCE alike. The value is a LengthNum —
		// an integer, not a string.
		//
		// Ref: `"NODEGROUP" EqOpt LengthNum` (PartDefOption) in TiDB v8.5.5
		// pkg/parser/parser.y.
		p.advance()
		p.match('=')
		val, err := p.consumeIntegerOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "NODEGROUP", Value: val}, true, nil

	case kwKEY_BLOCK_SIZE, kwSTATS_PERSISTENT,
		kwSTATS_SAMPLE_PAGES, kwMAX_ROWS, kwMIN_ROWS, kwAVG_ROW_LENGTH,
		kwPACK_KEYS, kwDELAY_KEY_WRITE, kwSECONDARY_ENGINE,
		kwAUTOEXTEND_SIZE:
		optName := p.cur.Str
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
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
		if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
			name, _, err := p.parseKeywordOrIdent()
			if err != nil {
				return "", err
			}
			return name, nil
		}
		// Every caller reads an option that takes a value, and TiDB answers 1064
		// when it is missing (`ENGINE =`, `PLACEMENT POLICY` at end of input).
		return "", p.syntaxErrorAtCur()
	}
}

// consumeStringNameOptionValue reads a StringName option value, the way STORAGE
// ENGINE and AFFINITY take it: a quoted string or an Identifier. Identifier is a
// plain or backquoted name or an unreserved keyword, so a bare number and a
// reserved word (`SELECT`) are each a 1064.
//
// Ref: `StringName: stringLit | Identifier` in TiDB v8.5.5 pkg/parser/parser.y.
func (p *Parser) consumeStringNameOptionValue() (string, error) {
	if p.cur.Type == tokSCONST {
		val := p.cur.Str
		p.advance()
		return val, nil
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return "", err
	}
	return name, nil
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
		algorithm, err := p.parsePartitionKeyAlgorithm()
		if err != nil {
			return nil, err
		}
		part.Algorithm = algorithm
		// KEY () with empty column list means "use PK columns". The parentheses are
		// required: `PARTITION BY KEY` alone is 1064.
		if p.cur.Type == '(' && p.peekNext().Type == ')' {
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
		if p.cur.Type == kwINTERVAL {
			interval, err := p.parsePartitionInterval()
			if err != nil {
				return nil, err
			}
			part.Interval = interval
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

	default:
		if !part.Linear && p.isBareWord("SYSTEM_TIME") {
			if err := p.parseSystemTimeMethod(part); err != nil {
				return nil, err
			}
			break
		}
		return nil, &ParseError{
			Message:  "expected HASH, KEY, RANGE, LIST, or SYSTEM_TIME after PARTITION BY",
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
			algorithm, err := p.parsePartitionKeyAlgorithm()
			if err != nil {
				return nil, err
			}
			part.SubPartAlgo = algorithm
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
	var bounds []partitionBound
	if p.cur.Type == '(' {
		p.advance()
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
	}

	// PartitionOpt's action validates what the productions read.
	if err := validatePartitionClause(part, bounds); err != nil {
		return nil, err
	}

	if err := p.parseUpdateIndexes(part); err != nil {
		return nil, err
	}

	part.Loc.End = p.prev.End
	return part, nil
}

// parseUpdateIndexes parses the UPDATE INDEXES clause TiDB allows at the end
// of a partitioning clause, which says for each named index whether it stays
// global or becomes local once the table is partitioned.
//
//	UPDATE INDEXES '(' index_name {GLOBAL | LOCAL} [, ...] ')'
//
// GLOBAL or LOCAL is required on every entry and the list may not be empty:
// TiDB v8.5.5 runs `PARTITION BY HASH(a) UPDATE INDEXES (i GLOBAL)`,
// `UPDATE INDEXES (i GLOBAL, j LOCAL)`,
// `PARTITION BY HASH(a) PARTITIONS 2 UPDATE INDEXES (i GLOBAL)`,
// `PARTITION BY RANGE(a) (PARTITION p VALUES LESS THAN (1)) UPDATE INDEXES
// (i GLOBAL)` and the same clause on CREATE TABLE, and answers 1064 for
// `UPDATE INDEXES ()`, `UPDATE INDEXES (i)` and `UPDATE INDEXES (i GLOBAL,)`.
// MySQL 8.0.46 and MariaDB 11.8.8 answer 1064 for the clause itself.
//
// Ref: UpdateIndexesOpt and UpdateIndexesList in TiDB v8.5.5
// pkg/parser/parser.y.
func (p *Parser) parseUpdateIndexes(part *nodes.PartitionClause) error {
	if p.cur.Type != kwUPDATE || p.peekNext().Type != kwINDEXES {
		return nil
	}
	p.advance() // consume UPDATE
	p.advance() // consume INDEXES
	if _, err := p.expect('('); err != nil {
		return err
	}
	for {
		start := p.pos()
		name, _, err := p.parseIdent()
		if err != nil {
			return err
		}
		tok, ok := p.match(kwGLOBAL, kwLOCAL)
		if !ok {
			return p.syntaxErrorAtCur()
		}
		part.UpdateIndexes = append(part.UpdateIndexes, &nodes.UpdateIndexSpec{
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
			Name:   name,
			Global: tok.Type == kwGLOBAL,
		})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	return nil
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

// parsePartitionDefAndBound is parsePartitionDef, and says which kind of bound the
// definition wrote, which PartitionDef.Values alone leaves open for a list of values.
func (p *Parser) parsePartitionDefAndBound() (*nodes.PartitionDef, partitionBound, error) {
	start := p.pos()
	if _, err := p.expect(kwPARTITION); err != nil {
		return nil, 0, err
	}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, 0, err
	}
	bound := boundNone

	pd := &nodes.PartitionDef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// PartDefValuesOpt's bare DEFAULT, which names the catch-all partition
	// without the VALUES IN (DEFAULT) spelling: TiDB v8.5.5 and MariaDB
	// 11.8.8 run `ADD PARTITION (PARTITION p DEFAULT)`, MySQL 8.0.46 answers
	// 1064.
	if _, ok := p.match(kwDEFAULT); ok {
		pd.Values = &nodes.String{Str: "DEFAULT"}
		bound = boundIn
	}

	// PartDefValuesOpt's HISTORY and CURRENT, which only SYSTEM_TIME takes.
	if _, ok := p.match(kwHISTORY); ok {
		pd.Values = &nodes.String{Str: "HISTORY"}
		bound = boundHistory
	} else if _, ok := p.match(kwCURRENT); ok {
		pd.Values = &nodes.String{Str: "CURRENT"}
		bound = boundHistory
	}

	// VALUES LESS THAN or VALUES IN
	if pd.Values == nil {
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
							// DEFAULT is a VALUES IN bound, not a LESS THAN one:
							// `VALUES LESS THAN (DEFAULT)` is a 1064 on v8.5.5.
							if p.cur.Type == kwDEFAULT {
								return nil, 0, p.syntaxErrorAtCur()
							}
							expr, err := p.parseBitExpr()
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
						} else if p.cur.Type == kwDEFAULT {
							// VALUES IN (DEFAULT) names the catch-all list
							// partition, and DEFAULT may sit anywhere in the list:
							// TiDB v8.5.5 runs `VALUES IN (DEFAULT)`,
							// `VALUES IN (DEFAULT, 1)` and `VALUES IN (1, DEFAULT)`,
							// and on CREATE TABLE too. MySQL 8.0.46 and MariaDB
							// 11.8.8 answer 1064 for all three.
							//
							// Ref: DefaultOrExpressionList in TiDB v8.5.5
							// pkg/parser/parser.y.
							items = append(items, &nodes.String{Str: "DEFAULT"})
							p.advance()
						} else {
							expr, err := p.parseBitExpr()
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
					pd.Values = &nodes.List{Items: items}
				} else {
					return nil, 0, p.syntaxErrorAtCur()
				}
			} else {
				// VALUES with neither LESS THAN nor IN after it is no bound.
				return nil, 0, p.syntaxErrorAtCur()
			}
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

// parseIfNotExists consumes the IF NOT EXISTS a grammar position allows and
// reports whether one was there. A bare IF, or IF NOT without EXISTS, is a
// syntax error rather than nothing: the words are the clause.
//
// Ref: IfNotExists in TiDB v8.5.5 pkg/parser/parser.y, which reads
// `/* empty */ | "IF" NotSym "EXISTS"`.
func (p *Parser) parseIfNotExists() (bool, error) {
	if p.cur.Type != kwIF {
		return false, nil
	}
	p.advance() // consume IF
	if _, err := p.expect(kwNOT); err != nil {
		return false, err
	}
	if _, err := p.expect(kwEXISTS_KW); err != nil {
		return false, err
	}
	return true, nil
}

// parseIfExists consumes the IF EXISTS a grammar position allows and reports
// whether one was there.
//
// Ref: IfExists in TiDB v8.5.5 pkg/parser/parser.y, which reads
// `/* empty */ | "IF" "EXISTS"`.
func (p *Parser) parseIfExists() (bool, error) {
	if p.cur.Type != kwIF {
		return false, nil
	}
	p.advance() // consume IF
	if _, err := p.expect(kwEXISTS_KW); err != nil {
		return false, err
	}
	return true, nil
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

// parseColumnKeyGlobalOrLocal consumes the GLOBAL or LOCAL that TiDB allows
// after a column's PRIMARY KEY, KEY or UNIQUE [KEY] option. One of the two at
// most, after the CLUSTERED attribute rather than before it: v8.5.5 runs
// `a INT PRIMARY KEY CLUSTERED GLOBAL`, `a INT UNIQUE LOCAL` and
// `a INT KEY GLOBAL`, and answers 1064 for `PRIMARY KEY GLOBAL CLUSTERED` and
// for `UNIQUE GLOBAL LOCAL`. MySQL 8.0.46 and MariaDB 11.8.8 answer 1064 for
// every one of them.
//
// The word is consumed and not recorded: the grammar's own action keeps only
// whether the index is global, which the column form never reaches a
// consumer by -- TiDB's parser.y folds LOCAL and the empty case to the same
// value.
//
// Ref: GlobalOrLocalOpt in TiDB v8.5.5 pkg/parser/parser.y.
func (p *Parser) parseColumnKeyGlobalOrLocal() {
	p.match(kwGLOBAL, kwLOCAL)
}

// parsePartDefOption parses one of the options a partition or subpartition
// definition takes, and reports whether one was there. The list is its own,
// narrower than a table's: a partition definition takes STORAGE ENGINE, which
// no table option spells, and takes none of ENGINE's neighbours such as
// CHARSET, AUTO_INCREMENT or the TTL options.
//
//	part_def_option:
//	    COMMENT [=] 'string'
//	  | ENGINE [=] name
//	  | STORAGE ENGINE [=] name
//	  | INSERT_METHOD [=] name
//	  | DATA DIRECTORY [=] 'string'
//	  | INDEX DIRECTORY [=] 'string'
//	  | MAX_ROWS [=] number
//	  | MIN_ROWS [=] number
//	  | TABLESPACE [=] name
//	  | NODEGROUP [=] number
//	  | PLACEMENT POLICY [=] {name | 'string' | DEFAULT}
//	  | PLACEMENT POLICY SET DEFAULT
//
// TiDB v8.5.5, MySQL 8.0.46 and MariaDB 11.8.8 all run
// `ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE ENGINE text)`,
// `... STORAGE ENGINE = text`, `... MAX_ROWS 12` and
// `... (SUBPARTITION s MIN_ROWS 11)`. PLACEMENT POLICY is TiDB's own: v8.5.5
// runs `... PLACEMENT POLICY x` where the other two answer 1064.
//
// Ref: PartDefOption and PlacementPolicyOption in TiDB v8.5.5
// pkg/parser/parser.y; opt_part_option in MySQL 8.0.46 sql/sql_yacc.yy.
func (p *Parser) parsePartDefOption() (*nodes.TableOption, bool, error) {
	start := p.pos()

	switch p.cur.Type {
	case kwCOMMENT:
		p.advance()
		p.match('=')
		val, err := p.consumeStringOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "COMMENT", Value: val}, true, nil

	case kwENGINE, kwINSERT_METHOD, kwTABLESPACE:
		name := strings.ToUpper(p.cur.Str)
		p.advance()
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: name, Value: val}, true, nil

	case kwSTORAGE:
		// STORAGE ENGINE [=] name. STORAGE on its own is no partition option:
		// the table option STORAGE {DISK | MEMORY} has no partition form, and
		// TiDB v8.5.5 answers 1064 for `PARTITION p STORAGE DISK`.
		if p.peekNext().Type != kwENGINE {
			return nil, false, nil
		}
		p.advance() // consume STORAGE
		p.advance() // consume ENGINE
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "STORAGE ENGINE", Value: val}, true, nil

	case kwMAX_ROWS, kwMIN_ROWS, kwNODEGROUP:
		name := strings.ToUpper(p.cur.Str)
		p.advance()
		p.match('=')
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
		p.match('=')
		val, err := p.consumeStringOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: name, Value: val}, true, nil

	case kwPLACEMENT:
		if p.peekNext().Type != kwPOLICY {
			return nil, false, nil
		}
		p.advance() // consume PLACEMENT
		p.advance() // consume POLICY
		// PLACEMENT POLICY SET DEFAULT is its own spelling of the default,
		// with no = before SET.
		if p.cur.Type == kwSET {
			p.advance()
			if _, err := p.expect(kwDEFAULT); err != nil {
				return nil, false, err
			}
			return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "PLACEMENT POLICY", Value: "DEFAULT"}, true, nil
		}
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "PLACEMENT POLICY", Value: val}, true, nil
	}

	return nil, false, nil
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
// three are a keyword list, not an identifier: MySQL 8.0.46, MariaDB 11.8.8
// and TiDB v8.5.5 all answer 1064 for
// `ADD FOREIGN KEY (a) REFERENCES u(b) MATCH foo`.
//
// Ref: `Match: "MATCH" "FULL" | "MATCH" "PARTIAL" | "MATCH" "SIMPLE"` in TiDB
// v8.5.5 pkg/parser/parser.y; opt_match_clause in MySQL 8.0.46
// sql/sql_yacc.yy.
func (p *Parser) parseMatchType() (string, error) {
	if tok, ok := p.match(kwFULL, kwPARTIAL, kwSIMPLE); ok {
		return strings.ToUpper(tok.Str), nil
	}
	return "", p.syntaxErrorAtCur()
}

// refuseColumnValue refuses a bare name as a column's DEFAULT or ON UPDATE value.
// DefaultValueExpr, parser.y:4052, takes a literal, a function call, or those in
// parentheses, and ON UPDATE only a function call, so `a INT DEFAULT b` and `a TIMESTAMP
// DEFAULT now` are 1064s.
func refuseColumnValue(expr nodes.ExprNode, start int) error {
	if _, ok := expr.(*nodes.ColumnRef); ok {
		return &ParseError{Message: "syntax error", Position: start}
	}
	return nil
}

// parseOnUpdateValue reads NowSymOptionFraction, parser.y:4091, which is what ON UPDATE takes and
// all it takes: CURRENT_TIMESTAMP, LOCALTIME or LOCALTIMESTAMP alone, one of them or NOW with `()` or
// `(NUM)`, or CURRENT_DATE or CURDATE. So parentheses around one, as in `ON UPDATE (NOW())`, and any
// other expression are 1064s, where DEFAULT reads a wider rule.
func (p *Parser) parseOnUpdateValue() (nodes.ExprNode, error) {
	first := p.cur.Type
	switch first {
	case kwCURRENT_TIMESTAMP, kwLOCALTIME, kwLOCALTIMESTAMP, kwNOW, kwCURRENT_DATE, kwCURDATE:
	default:
		return nil, p.syntaxErrorAtCur()
	}
	start := p.pos()
	expr, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}
	fn, ok := expr.(*nodes.FuncCallExpr)
	if !ok {
		return nil, &ParseError{Message: "syntax error", Position: start}
	}
	switch len(fn.Args) {
	case 0:
	case 1:
		// A fraction is a number, and a date has none.
		_, isNum := fn.Args[0].(*nodes.IntLit)
		if !isNum || first == kwCURRENT_DATE || first == kwCURDATE {
			return nil, &ParseError{Message: "syntax error", Position: start}
		}
	default:
		return nil, &ParseError{Message: "syntax error", Position: start}
	}
	return expr, nil
}
