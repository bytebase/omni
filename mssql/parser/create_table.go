// Package parser - create_table.go implements T-SQL CREATE TABLE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateTableStmt parses a CREATE TABLE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-table-transact-sql
//
//	CREATE TABLE
//	    { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	    [ AS FileTable ]
//	    ( { <column_definition>
//	        | <computed_column_definition>
//	        | <column_set_definition>
//	        | [ <table_constraint> ] [ ,...n ]
//	        | [ <table_index> ] }
//	          [ ,...n ]
//	          [ PERIOD FOR SYSTEM_TIME ( system_start_time_column_name
//	             , system_end_time_column_name ) ]
//	      )
//	    [ ON { partition_scheme_name ( partition_column_name )
//	           | filegroup
//	           | "default" } ]
//	    [ TEXTIMAGE_ON { filegroup | "default" } ]
//	    [ FILESTREAM_ON { partition_scheme_name | filegroup | "default" } ]
//	    [ WITH ( <table_option> [ ,...n ] ) ]
func (p *Parser) parseCreateTableStmt() *nodes.CreateTableStmt {
	loc := p.pos()

	stmt := &nodes.CreateTableStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Table name
	stmt.Name = p.parseTableRef()

	// AS NODE | AS EDGE (graph tables) | AS FILETABLE
	// Note: AS CLONE OF is handled at the dispatch level (parser.go)
	if p.cur.Type == kwAS {
		next := p.peekNext()
		if next.Type == tokIDENT && strings.EqualFold(next.Str, "node") {
			p.advance() // AS
			p.advance() // NODE
			stmt.IsNode = true
		} else if next.Type == tokIDENT && strings.EqualFold(next.Str, "edge") {
			p.advance() // AS
			p.advance() // EDGE
			stmt.IsEdge = true
		}
	}

	// Column and constraint definitions
	if _, err := p.expect('('); err != nil {
		stmt.Loc.End = p.pos()
		return stmt
	}

	var cols []nodes.Node
	var constraints []nodes.Node

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// Check for PERIOD FOR SYSTEM_TIME
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "period") {
			if p.peekNext().Type == kwFOR {
				p.advance() // PERIOD
				p.advance() // FOR
				// SYSTEM_TIME
				if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "system_time") {
					p.advance()
				}
				// ( start_col , end_col )
				if p.cur.Type == '(' {
					p.advance()
					startCol, _ := p.parseIdentifier()
					stmt.PeriodStartCol = startCol
					p.match(',')
					endCol, _ := p.parseIdentifier()
					stmt.PeriodEndCol = endCol
					p.expect(')')
				}
				if _, ok := p.match(','); !ok {
					break
				}
				continue
			}
		}

		// Check for inline INDEX definition
		if p.cur.Type == kwINDEX {
			idx := p.parseInlineTableIndex()
			if idx != nil {
				if stmt.Indexes == nil {
					stmt.Indexes = &nodes.List{}
				}
				stmt.Indexes.Items = append(stmt.Indexes.Items, idx)
			}
		} else if p.cur.Type == kwCONSTRAINT || p.cur.Type == kwPRIMARY ||
			p.cur.Type == kwUNIQUE || p.cur.Type == kwCHECK ||
			p.cur.Type == kwFOREIGN {
			constraint := p.parseTableConstraint()
			if constraint != nil {
				constraints = append(constraints, constraint)
			}
		} else {
			col := p.parseColumnDef()
			if col != nil {
				cols = append(cols, col)
			}
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')

	if len(cols) > 0 {
		stmt.Columns = &nodes.List{Items: cols}
	}
	if len(constraints) > 0 {
		stmt.Constraints = &nodes.List{Items: constraints}
	}

	// ON filegroup
	if p.cur.Type == kwON {
		p.advance()
		if p.isIdentLike() {
			stmt.OnFilegroup = p.cur.Str
			p.advance()
		}
	}

	// TEXTIMAGE_ON filegroup
	if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "textimage_on") {
		p.advance()
		if p.isIdentLike() {
			stmt.TextImageOn = p.cur.Str
			p.advance()
		}
	}

	// FILESTREAM_ON filegroup
	if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "filestream_on") {
		p.advance()
		if p.isIdentLike() {
			stmt.FilestreamOn = p.cur.Str
			p.advance()
		}
	}

	// WITH ( table_options )
	if p.cur.Type == kwWITH {
		p.advance()
		stmt.TableOptions = p.parseTableOptions()
	}

	// AS NODE | AS EDGE (graph tables) -- can also appear after closing paren
	if !stmt.IsNode && !stmt.IsEdge && p.cur.Type == kwAS {
		next := p.peekNext()
		if next.Type == tokIDENT && strings.EqualFold(next.Str, "node") {
			p.advance() // AS
			p.advance() // NODE
			stmt.IsNode = true
		} else if next.Type == tokIDENT && strings.EqualFold(next.Str, "edge") {
			p.advance() // AS
			p.advance() // EDGE
			stmt.IsEdge = true
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseTableOptions parses ( option = value [, ...] ) for CREATE TABLE WITH clause.
//
//	<table_option> ::=
//	    DATA_COMPRESSION = { NONE | ROW | PAGE }
//	  | XML_COMPRESSION = { ON | OFF }
//	  | MEMORY_OPTIMIZED = ON
//	  | DURABILITY = { SCHEMA_ONLY | SCHEMA_AND_DATA }
//	  | SYSTEM_VERSIONING = ON [ ( HISTORY_TABLE = schema.table [, DATA_CONSISTENCY_CHECK = { ON | OFF } ] ) ]
//	  | LEDGER = ON | OFF
func (p *Parser) parseTableOptions() *nodes.List {
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()

	var items []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		opt := p.parseOneTableOption()
		if opt != nil {
			items = append(items, opt)
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	p.expect(')')

	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

// parseOneTableOption parses a single NAME = VALUE table option.
func (p *Parser) parseOneTableOption() *nodes.TableOption {
	loc := p.pos()
	if !p.isIdentLike() {
		return nil
	}
	name := strings.ToUpper(p.cur.Str)
	p.advance()

	opt := &nodes.TableOption{
		Name: name,
		Loc:  nodes.Loc{Start: loc},
	}

	// Expect '='
	if p.cur.Type != '=' {
		opt.Loc.End = p.pos()
		return opt
	}
	p.advance()

	// SYSTEM_VERSIONING = ON [ ( ... ) ]
	if name == "SYSTEM_VERSIONING" {
		if p.isIdentLike() {
			opt.Value = strings.ToUpper(p.cur.Str)
			p.advance()
		}
		// Optional sub-options
		if opt.Value == "ON" && p.cur.Type == '(' {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.isIdentLike() {
					subName := strings.ToUpper(p.cur.Str)
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
					}
					switch subName {
					case "HISTORY_TABLE":
						// schema.table - collect dotted name
						var parts []string
						if p.isIdentLike() {
							parts = append(parts, p.cur.Str)
							p.advance()
						}
						for p.cur.Type == '.' {
							p.advance()
							if p.isIdentLike() {
								parts = append(parts, p.cur.Str)
								p.advance()
							}
						}
						opt.HistoryTable = strings.Join(parts, ".")
					case "DATA_CONSISTENCY_CHECK":
						if p.isIdentLike() {
							opt.DataConsistencyCheck = strings.ToUpper(p.cur.Str)
							p.advance()
						}
					case "HISTORY_RETENTION_PERIOD":
						// Collect value tokens until comma or paren
						var valParts []string
						for p.cur.Type != ',' && p.cur.Type != ')' && p.cur.Type != tokEOF {
							valParts = append(valParts, p.cur.Str)
							p.advance()
						}
						opt.HistoryRetentionPeriod = strings.Join(valParts, " ")
					default:
						// Consume unknown sub-option value structurally (ON/OFF/identifier/number)
						if p.cur.Type == kwON {
							p.advance()
						} else if p.cur.Type == kwOFF {
							p.advance()
						} else if p.isIdentLike() {
							p.advance()
						} else if p.cur.Type == tokICONST {
							p.advance()
						}
					}
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.expect(')')
		}
	} else {
		// Simple NAME = VALUE
		if p.isIdentLike() {
			opt.Value = strings.ToUpper(p.cur.Str)
			p.advance()
		} else if p.cur.Type == kwON {
			opt.Value = "ON"
			p.advance()
		} else if p.cur.Type == kwOFF {
			opt.Value = "OFF"
			p.advance()
		}
	}

	opt.Loc.End = p.pos()
	return opt
}

// parseColumnDef parses a column definition.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-table-transact-sql
//
//	<column_definition> ::=
//	column_name <data_type>
//	    [ FILESTREAM ]
//	    [ COLLATE collation_name ]
//	    [ SPARSE ]
//	    [ MASKED WITH ( FUNCTION = 'mask_function' ) ]
//	    [ [ CONSTRAINT constraint_name ] DEFAULT constant_expression ]
//	    [ IDENTITY [ ( seed , increment ) ] ]
//	    [ NOT FOR REPLICATION ]
//	    [ GENERATED ALWAYS AS { ROW | TRANSACTION_ID | SEQUENCE_NUMBER }
//	      { START | END } [ HIDDEN ] ]
//	    [ [ CONSTRAINT constraint_name ] { NULL | NOT NULL } ]
//	    [ ROWGUIDCOL ]
//	    [ ENCRYPTED WITH
//	        ( COLUMN_ENCRYPTION_KEY = key_name ,
//	          ENCRYPTION_TYPE = { DETERMINISTIC | RANDOMIZED } ,
//	          ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'
//	        ) ]
//	    [ <column_constraint> [ ,...n ] ]
func (p *Parser) parseColumnDef() *nodes.ColumnDef {
	loc := p.pos()

	name, ok := p.parseIdentifier()
	if !ok {
		return nil
	}

	col := &nodes.ColumnDef{
		Name: name,
		Loc:  nodes.Loc{Start: loc},
	}

	// Check for computed column: name AS expr
	if p.cur.Type == kwAS {
		p.advance()
		compLoc := p.pos()
		expr := p.parseExpr()
		persisted := false
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "persisted") {
			persisted = true
			p.advance()
		}
		col.Computed = &nodes.ComputedColumnDef{
			Expr:      expr,
			Persisted: persisted,
			Loc:       nodes.Loc{Start: compLoc},
		}
		col.Loc.End = p.pos()
		return col
	}

	// Data type
	col.DataType = p.parseDataType()

	// Column-level options in any order
	for {
		consumed := false

		// FILESTREAM
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "filestream") {
			col.Filestream = true
			p.advance()
			consumed = true
		}

		// SPARSE
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "sparse") {
			col.Sparse = true
			p.advance()
			consumed = true
		}

		// ROWGUIDCOL
		if p.cur.Type == kwROWGUIDCOL {
			col.Rowguidcol = true
			p.advance()
			consumed = true
		}

		// HIDDEN (standalone, not after GENERATED ALWAYS)
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "hidden") {
			col.Hidden = true
			p.advance()
			consumed = true
		}

		// MASKED WITH (FUNCTION = 'mask_function')
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "masked") {
			p.advance() // MASKED
			if p.cur.Type == kwWITH {
				p.advance() // WITH
				if p.cur.Type == '(' {
					p.advance()
					// FUNCTION = 'value'
					if p.isIdentLike() && strings.EqualFold(p.cur.Str, "function") {
						p.advance()
						if p.cur.Type == '=' {
							p.advance()
						}
						if p.cur.Type == tokSCONST {
							col.MaskFunction = p.cur.Str
							p.advance()
						}
					}
					p.expect(')')
				}
			}
			consumed = true
		}

		// ENCRYPTED WITH (...)
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "encrypted") {
			col.EncryptedWith = p.parseEncryptedWith()
			consumed = true
		}

		// GENERATED ALWAYS AS { ROW | TRANSACTION_ID | SEQUENCE_NUMBER } { START | END } [ HIDDEN ]
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "generated") {
			col.GeneratedAlways = p.parseGeneratedAlways()
			// GENERATED ALWAYS ... HIDDEN
			if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "hidden") {
				col.Hidden = true
				p.advance()
			}
			consumed = true
		}

		// IDENTITY
		if p.cur.Type == kwIDENTITY {
			col.Identity = p.parseIdentitySpec()
			consumed = true
		}

		// NOT FOR REPLICATION / NOT NULL
		if p.cur.Type == kwNOT {
			next := p.peekNext()
			if next.Type == kwNULL {
				p.advance() // NOT
				p.advance() // NULL
				col.Nullable = &nodes.NullableSpec{NotNull: true, Loc: nodes.Loc{Start: p.pos()}}
				consumed = true
			} else if next.Type == kwFOR {
				p.advance() // NOT
				p.advance() // FOR
				if p.cur.Type == kwREPLICATION {
					p.advance()
				}
				col.NotForReplication = true
				consumed = true
			}
		}

		// NULL
		if p.cur.Type == kwNULL {
			p.advance()
			col.Nullable = &nodes.NullableSpec{NotNull: false, Loc: nodes.Loc{Start: p.pos()}}
			consumed = true
		}

		// DEFAULT
		if p.cur.Type == kwDEFAULT {
			p.advance()
			col.DefaultExpr = p.parseExpr()
			consumed = true
		}

		// COLLATE
		if p.cur.Type == kwCOLLATE {
			p.advance()
			if p.isIdentLike() {
				col.Collation = p.cur.Str
				p.advance()
			}
			consumed = true
		}

		// CONSTRAINT (inline column constraint)
		if p.cur.Type == kwCONSTRAINT {
			constraint := p.parseColumnConstraint()
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			consumed = true
		}

		// PRIMARY KEY / UNIQUE (without CONSTRAINT keyword)
		if p.cur.Type == kwPRIMARY || p.cur.Type == kwUNIQUE {
			constraint := p.parseInlineConstraint("")
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			consumed = true
		}

		// CHECK (without CONSTRAINT keyword)
		if p.cur.Type == kwCHECK {
			constraint := p.parseInlineConstraint("")
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			consumed = true
		}

		// REFERENCES (inline FK without CONSTRAINT keyword)
		if p.cur.Type == kwREFERENCES {
			constraint := p.parseInlineConstraint("")
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			consumed = true
		}

		if !consumed {
			break
		}
	}

	col.Loc.End = p.pos()
	return col
}

// parseEncryptedWith parses ENCRYPTED WITH (...).
//
//	ENCRYPTED WITH
//	    ( COLUMN_ENCRYPTION_KEY = key_name ,
//	      ENCRYPTION_TYPE = { DETERMINISTIC | RANDOMIZED } ,
//	      ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'
//	    )
func (p *Parser) parseEncryptedWith() *nodes.EncryptedWithSpec {
	loc := p.pos()
	p.advance() // ENCRYPTED

	spec := &nodes.EncryptedWithSpec{
		Loc: nodes.Loc{Start: loc},
	}

	if p.cur.Type != kwWITH {
		spec.Loc.End = p.pos()
		return spec
	}
	p.advance() // WITH

	if p.cur.Type != '(' {
		spec.Loc.End = p.pos()
		return spec
	}
	p.advance()

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.isIdentLike() {
			optName := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			switch optName {
			case "COLUMN_ENCRYPTION_KEY":
				if p.isIdentLike() {
					spec.ColumnEncryptionKey = p.cur.Str
					p.advance()
				}
			case "ENCRYPTION_TYPE":
				if p.isIdentLike() {
					spec.EncryptionType = strings.ToUpper(p.cur.Str)
					p.advance()
				}
			case "ALGORITHM":
				if p.cur.Type == tokSCONST {
					spec.Algorithm = p.cur.Str
					p.advance()
				}
			default:
				if p.isIdentLike() || p.cur.Type == tokSCONST {
					p.advance()
				}
			}
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	p.expect(')')

	spec.Loc.End = p.pos()
	return spec
}

// parseGeneratedAlways parses GENERATED ALWAYS AS { ROW | TRANSACTION_ID | SEQUENCE_NUMBER } { START | END }.
//
//	GENERATED ALWAYS AS { ROW | TRANSACTION_ID | SEQUENCE_NUMBER } { START | END } [ HIDDEN ]
func (p *Parser) parseGeneratedAlways() *nodes.GeneratedAlwaysSpec {
	loc := p.pos()
	p.advance() // GENERATED

	spec := &nodes.GeneratedAlwaysSpec{
		Loc: nodes.Loc{Start: loc},
	}

	// ALWAYS
	if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "always") {
		p.advance()
	}

	// AS
	if p.cur.Type == kwAS {
		p.advance()
	}

	// ROW | TRANSACTION_ID | SEQUENCE_NUMBER
	if p.isIdentLike() {
		kind := strings.ToUpper(p.cur.Str)
		p.advance()
		spec.Kind = kind
	}

	// START | END
	if p.isIdentLike() {
		startEnd := strings.ToUpper(p.cur.Str)
		if startEnd == "START" || startEnd == "END" {
			spec.StartEnd = startEnd
			p.advance()
		}
	}
	// Handle END keyword (it's a reserved keyword kwEND)
	if p.cur.Type == kwEND {
		spec.StartEnd = "END"
		p.advance()
	}

	spec.Loc.End = p.pos()
	return spec
}

// parseIdentitySpec parses IDENTITY(seed, increment).
func (p *Parser) parseIdentitySpec() *nodes.IdentitySpec {
	loc := p.pos()
	p.advance() // consume IDENTITY

	spec := &nodes.IdentitySpec{
		Seed:      1,
		Increment: 1,
		Loc:       nodes.Loc{Start: loc},
	}

	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == tokICONST {
			spec.Seed = p.cur.Ival
			p.advance()
		}
		if _, ok := p.match(','); ok {
			if p.cur.Type == tokICONST {
				spec.Increment = p.cur.Ival
				p.advance()
			}
		}
		_, _ = p.expect(')')
	}

	spec.Loc.End = p.pos()
	return spec
}

// parseColumnConstraint parses CONSTRAINT name followed by constraint type.
func (p *Parser) parseColumnConstraint() *nodes.ConstraintDef {
	p.advance() // consume CONSTRAINT
	name, _ := p.parseIdentifier()
	return p.parseInlineConstraint(name)
}

// parseInlineConstraint parses a constraint type (PRIMARY KEY, UNIQUE, CHECK, DEFAULT, REFERENCES).
func (p *Parser) parseInlineConstraint(name string) *nodes.ConstraintDef {
	loc := p.pos()
	cd := &nodes.ConstraintDef{
		Name: name,
		Loc:  nodes.Loc{Start: loc},
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // PRIMARY
		p.match(kwKEY)
		cd.Type = nodes.ConstraintPrimaryKey
		p.parseClusteredOption(cd)
	case kwUNIQUE:
		p.advance()
		cd.Type = nodes.ConstraintUnique
		p.parseClusteredOption(cd)
	case kwCHECK:
		p.advance()
		cd.Type = nodes.ConstraintCheck
		if _, err := p.expect('('); err == nil {
			cd.Expr = p.parseExpr()
			_, _ = p.expect(')')
		}
	case kwDEFAULT:
		p.advance()
		cd.Type = nodes.ConstraintDefault
		cd.Expr = p.parseExpr()
	case kwREFERENCES:
		p.advance()
		cd.Type = nodes.ConstraintForeignKey
		cd.RefTable = p.parseTableRef()
		if p.cur.Type == '(' {
			cd.RefColumns = p.parseParenIdentList()
		}
		p.parseReferentialActions(cd)
	default:
		return nil
	}

	cd.Loc.End = p.pos()
	return cd
}

// parseTableConstraint parses a table-level constraint.
func (p *Parser) parseTableConstraint() *nodes.ConstraintDef {
	loc := p.pos()
	var name string

	if p.cur.Type == kwCONSTRAINT {
		p.advance()
		name, _ = p.parseIdentifier()
	}

	cd := &nodes.ConstraintDef{
		Name: name,
		Loc:  nodes.Loc{Start: loc},
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // PRIMARY
		p.match(kwKEY)
		cd.Type = nodes.ConstraintPrimaryKey
		p.parseClusteredOption(cd)
		if p.cur.Type == '(' {
			cd.Columns = p.parseParenIdentList()
		}
	case kwUNIQUE:
		p.advance()
		cd.Type = nodes.ConstraintUnique
		p.parseClusteredOption(cd)
		if p.cur.Type == '(' {
			cd.Columns = p.parseParenIdentList()
		}
	case kwCHECK:
		p.advance()
		cd.Type = nodes.ConstraintCheck
		if _, err := p.expect('('); err == nil {
			cd.Expr = p.parseExpr()
			_, _ = p.expect(')')
		}
	case kwFOREIGN:
		p.advance() // FOREIGN
		p.match(kwKEY)
		cd.Type = nodes.ConstraintForeignKey
		if p.cur.Type == '(' {
			cd.Columns = p.parseParenIdentList()
		}
		if _, ok := p.match(kwREFERENCES); ok {
			cd.RefTable = p.parseTableRef()
			if p.cur.Type == '(' {
				cd.RefColumns = p.parseParenIdentList()
			}
			p.parseReferentialActions(cd)
		}
	case kwDEFAULT:
		p.advance()
		cd.Type = nodes.ConstraintDefault
		cd.Expr = p.parseExpr()
		// FOR column
		if _, ok := p.match(kwFOR); ok {
			p.parseIdentifier() // column name (not stored separately but consumed)
		}
	default:
		// Check for EDGE CONSTRAINT: CONSTRAINT name CONNECTION (...)
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "CONNECTION") {
			cd.Type = nodes.ConstraintEdge
			cd.EdgeConnections = p.parseEdgeConstraintConnections()
			p.parseReferentialActions(cd)
		} else {
			return nil
		}
	}

	cd.Loc.End = p.pos()
	return cd
}

// parseEdgeConstraintConnections parses CONNECTION ( from_table TO to_table [, ...] ).
//
// Ref: https://learn.microsoft.com/en-us/sql/relational-databases/tables/graph-edge-constraints
//
//	CONNECTION ( from_table TO to_table [ , from_table TO to_table ... ] )
func (p *Parser) parseEdgeConstraintConnections() *nodes.List {
	p.advance() // consume CONNECTION

	var conns []nodes.Node
	if _, err := p.expect('('); err != nil {
		return nil
	}

	for {
		loc := p.pos()
		fromTable := p.parseTableRef()
		p.match(kwTO)
		toTable := p.parseTableRef()
		conn := &nodes.EdgeConnectionDef{
			FromTable: fromTable,
			ToTable:   toTable,
			Loc:       nodes.Loc{Start: loc, End: p.pos()},
		}
		conns = append(conns, conn)
		if _, ok := p.match(','); !ok {
			break
		}
	}

	p.expect(')')
	return &nodes.List{Items: conns}
}

// parseClusteredOption parses optional CLUSTERED/NONCLUSTERED.
func (p *Parser) parseClusteredOption(cd *nodes.ConstraintDef) {
	if p.cur.Type == kwCLUSTERED {
		p.advance()
		v := true
		cd.Clustered = &v
	} else if p.cur.Type == kwNONCLUSTERED {
		p.advance()
		v := false
		cd.Clustered = &v
	}
}

// parseReferentialActions parses ON DELETE/UPDATE actions.
func (p *Parser) parseReferentialActions(cd *nodes.ConstraintDef) {
	for p.cur.Type == kwON {
		p.advance()
		if p.cur.Type == kwDELETE {
			p.advance()
			cd.OnDelete = p.parseRefAction()
		} else if p.cur.Type == kwUPDATE {
			p.advance()
			cd.OnUpdate = p.parseRefAction()
		} else {
			break
		}
	}
}

// parseRefAction parses a referential action (CASCADE, SET NULL, SET DEFAULT, NO ACTION).
func (p *Parser) parseRefAction() nodes.ReferentialAction {
	if _, ok := p.match(kwCASCADE); ok {
		return nodes.RefActCascade
	}
	if p.cur.Type == kwSET {
		p.advance()
		if _, ok := p.match(kwNULL); ok {
			return nodes.RefActSetNull
		}
		if _, ok := p.match(kwDEFAULT); ok {
			return nodes.RefActSetDefault
		}
	}
	// NO ACTION
	if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "no") {
		p.advance()
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "action") {
			p.advance()
		}
		return nodes.RefActNoAction
	}
	return nodes.RefActNone
}

// parseInlineTableIndex parses an inline INDEX definition inside a CREATE TABLE body.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-table-transact-sql
//
//	INDEX index_name [ UNIQUE ] [ CLUSTERED | NONCLUSTERED ]
//	    ( column_name [ ASC | DESC ] [ ,...n ] )
//	    [ INCLUDE ( column_name [ ,...n ] ) ]
//	    [ WHERE <filter_predicate> ]
//	    [ WITH ( <index_option> [ ,...n ] ) ]
func (p *Parser) parseInlineTableIndex() *nodes.InlineIndexDef {
	loc := p.pos()
	p.advance() // consume INDEX

	idx := &nodes.InlineIndexDef{
		Loc: nodes.Loc{Start: loc},
	}

	// index_name
	name, _ := p.parseIdentifier()
	idx.Name = name

	// [ UNIQUE ]
	if p.cur.Type == kwUNIQUE {
		idx.Unique = true
		p.advance()
	}

	// [ CLUSTERED | NONCLUSTERED ]
	if p.cur.Type == kwCLUSTERED {
		v := true
		idx.Clustered = &v
		p.advance()
	} else if p.cur.Type == kwNONCLUSTERED {
		v := false
		idx.Clustered = &v
		p.advance()
	}

	// ( column_name [ ASC | DESC ] [ ,...n ] )
	if p.cur.Type == '(' {
		idx.Columns = p.parseIndexColumnList()
	}

	// [ INCLUDE ( column_name [ ,...n ] ) ]
	if p.cur.Type == kwINCLUDE {
		p.advance()
		if p.cur.Type == '(' {
			idx.IncludeCols = p.parseParenIdentList()
		}
	}

	// [ WHERE <filter_predicate> ]
	if _, ok := p.match(kwWHERE); ok {
		idx.WhereClause = p.parseExpr()
	}

	// [ WITH ( <index_option> [ ,...n ] ) ]
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			idx.Options = p.parseOptionList()
		}
	}

	idx.Loc.End = p.pos()
	return idx
}

// parseParenIdentList parses (ident, ident, ...).
func (p *Parser) parseParenIdentList() *nodes.List {
	p.advance() // consume (
	var items []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		name, ok := p.parseIdentifier()
		if !ok {
			break
		}
		items = append(items, &nodes.String{Str: name})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')
	return &nodes.List{Items: items}
}
