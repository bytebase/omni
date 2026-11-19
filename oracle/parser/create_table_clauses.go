package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads the clauses of CREATE TABLE's physical_properties and
// table_properties by the grammar of the CREATE TABLE page.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-TABLE.html
// BNF: oracle/parser/bnf/CREATE-TABLE.bnf

// parseCreateTableSharing reads the sharing clause of an application
// container's table.
//
//	SHARING = { METADATA | DATA | EXTENDED DATA | NONE }
func (p *Parser) parseCreateTableSharing(stmt *nodes.CreateTableStmt) error {
	p.advance() // consume SHARING
	if err := p.expectCreateTableToken('='); err != nil {
		return err
	}
	switch {
	case p.isKeywordStr("METADATA"):
		stmt.Sharing = "METADATA"
	case p.isKeywordStr("DATA"):
		stmt.Sharing = "DATA"
	case p.isKeywordStr("NONE"):
		stmt.Sharing = "NONE"
	case p.isKeywordStr("EXTENDED"):
		p.advance() // consume EXTENDED
		stmt.Sharing = "EXTENDED DATA"
		if !p.isKeywordStr("DATA") {
			return p.syntaxErrorAtCur()
		}
	default:
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// createTableOptionValue returns the text of the tokens from byte start to
// the end of the last consumed token, joined by single spaces, the form
// DDLOption values of CREATE TABLE have always taken.
func (p *Parser) createTableOptionValue(start int) string {
	end := p.prev.End
	if end <= start {
		return ""
	}
	lexer := NewLexerRange(p.source, start, end)
	var parts []string
	for {
		tok := lexer.NextToken()
		if tok.Type == tokEOF || tok.Loc >= end {
			break
		}
		parts = append(parts, p.ddlOptionTokenText(tok))
	}
	return strings.Join(parts, " ")
}

// expectCreateTableWord consumes the unquoted word s or fails at the
// current token.
func (p *Parser) expectCreateTableWord(s string) error {
	if !p.isKeywordStr(s) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// expectCreateTableToken consumes a token of type t or fails at the
// current token.
func (p *Parser) expectCreateTableToken(t int) error {
	if p.cur.Type != t {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseCreateTableInteger consumes an integer literal.
func (p *Parser) parseCreateTableInteger() error {
	return p.expectCreateTableToken(tokICONST)
}

// isCreateTableNameStart reports whether a name starts here: a quoted
// identifier, or a word that is not reserved.
func (p *Parser) isCreateTableNameStart() bool {
	return p.cur.Type == tokQIDENT || p.isIdentLike() && !isOracleSQLReservedKeyword(p.cur)
}

// parseCreateTableName consumes a name: a quoted identifier, or a word
// that is not reserved.
func (p *Parser) parseCreateTableName() error {
	if !p.isCreateTableNameStart() {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseCreateTableColumnNames reads ( column [, column ] ... ).
func (p *Parser) parseCreateTableColumnNames() error {
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	for {
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return p.expectCreateTableToken(')')
}

// isCreateTableLoggingStart reports whether a logging_clause starts here.
//
//	logging_clause: { LOGGING | NOLOGGING | FILESYSTEM_LIKE_LOGGING }
func (p *Parser) isCreateTableLoggingStart() bool {
	return p.cur.Type == kwLOGGING || p.cur.Type == kwNOLOGGING || p.isKeywordStr("FILESYSTEM_LIKE_LOGGING")
}

// ---------------------------------------------------------------------------
// inmemory_table_clause
// ---------------------------------------------------------------------------

// isCreateTableMemcompressStart reports whether an inmemory_memcompress
// starts here: MEMCOMPRESS, or NO MEMCOMPRESS.
func (p *Parser) isCreateTableMemcompressStart() bool {
	if p.isKeywordStr("MEMCOMPRESS") {
		return true
	}
	return p.isKeywordStr("NO") && isCreateTableWordTok(p.peekNext(), "MEMCOMPRESS")
}

// isCreateTableWordTok reports whether tok is the unquoted word s.
func isCreateTableWordTok(tok Token, s string) bool {
	return tok.Type != tokQIDENT && (tok.Type == tokIDENT || tok.Type >= 2000) && tok.Str == s
}

// parseCreateTableMemcompress reads inmemory_memcompress.
//
//	{ MEMCOMPRESS FOR { DML | QUERY [ LOW | HIGH ] | CAPACITY [ LOW | HIGH ] }
//	| NO MEMCOMPRESS
//	}
func (p *Parser) parseCreateTableMemcompress() error {
	if p.isKeywordStr("NO") {
		p.advance() // consume NO
		p.advance() // consume MEMCOMPRESS
		return nil
	}
	p.advance() // consume MEMCOMPRESS
	if err := p.expectCreateTableToken(kwFOR); err != nil {
		return err
	}
	switch {
	case p.isKeywordStr("DML"):
		p.advance()
	case p.isKeywordStr("QUERY"), p.isKeywordStr("CAPACITY"):
		p.advance()
		if p.isKeywordStr("LOW") || p.isKeywordStr("HIGH") {
			p.advance()
		}
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseCreateTableInmemoryAttributes reads inmemory_attributes, whose
// parts Oracle takes in any order, each once.
//
//	[ inmemory_memcompress ] [ inmemory_priority ]
//	[ inmemory_distribute ] [ inmemory_duplicate ]
//
//	inmemory_priority: PRIORITY { NONE | LOW | MEDIUM | HIGH | CRITICAL }
//	inmemory_distribute:
//	    DISTRIBUTE [ AUTO | BY { ROWID RANGE | PARTITION | SUBPARTITION } ]
//	    [ FOR SERVICE { DEFAULT | ALL | service_name | NONE } ]
//	inmemory_duplicate: { DUPLICATE | DUPLICATE ALL | NO DUPLICATE }
//
// A service name is a quoted identifier: Oracle 23ai refuses FOR SERVICE a
// and FOR SERVICE 'a' (ORA-00922) and takes FOR SERVICE "a".
func (p *Parser) parseCreateTableInmemoryAttributes(memcompressSeen bool) error {
	var priority, distribute, duplicate bool
	for {
		switch {
		case !memcompressSeen && p.isCreateTableMemcompressStart():
			if err := p.parseCreateTableMemcompress(); err != nil {
				return err
			}
			memcompressSeen = true
		case !priority && p.isKeywordStr("PRIORITY"):
			p.advance() // consume PRIORITY
			switch {
			case p.isKeywordStr("NONE"), p.isKeywordStr("LOW"), p.isKeywordStr("MEDIUM"),
				p.isKeywordStr("HIGH"), p.isKeywordStr("CRITICAL"):
				p.advance()
			default:
				return p.syntaxErrorAtCur()
			}
			priority = true
		case !distribute && p.isKeywordStr("DISTRIBUTE"):
			p.advance() // consume DISTRIBUTE
			switch {
			case p.isKeywordStr("AUTO"):
				p.advance()
			case p.cur.Type == kwBY:
				p.advance() // consume BY
				switch {
				case p.cur.Type == kwROWID:
					p.advance() // consume ROWID
					if err := p.expectCreateTableWord("RANGE"); err != nil {
						return err
					}
				case p.cur.Type == kwPARTITION, p.isKeywordStr("SUBPARTITION"):
					p.advance()
				default:
					return p.syntaxErrorAtCur()
				}
			}
			if p.cur.Type == kwFOR {
				p.advance() // consume FOR
				if err := p.expectCreateTableWord("SERVICE"); err != nil {
					return err
				}
				switch {
				case p.cur.Type == kwDEFAULT, p.cur.Type == kwALL, p.isKeywordStr("NONE"), p.cur.Type == tokQIDENT:
					p.advance()
				default:
					return p.syntaxErrorAtCur()
				}
			}
			distribute = true
		case !duplicate && p.isKeywordStr("DUPLICATE"):
			p.advance() // consume DUPLICATE
			if p.cur.Type == kwALL {
				p.advance()
			}
			duplicate = true
		case !duplicate && p.isKeywordStr("NO") && isCreateTableWordTok(p.peekNext(), "DUPLICATE"):
			p.advance() // consume NO
			p.advance() // consume DUPLICATE
			duplicate = true
		default:
			return nil
		}
	}
}

// isCreateTableInmemoryStart reports whether an inmemory_table_clause or
// an inmemory_column_clause starts here: INMEMORY, or NO INMEMORY.
func (p *Parser) isCreateTableInmemoryStart() bool {
	if p.isKeywordStr("INMEMORY") {
		return true
	}
	return p.isKeywordStr("NO") && isCreateTableWordTok(p.peekNext(), "INMEMORY")
}

// parseCreateTableInmemory reads an inmemory_table_clause, or one
// inmemory_column_clause item, which Oracle takes on its own among the
// table's properties (CREATE TABLE a(a INT) INMEMORY(a) is accepted).
//
//	inmemory_table_clause:
//	    { INMEMORY [ inmemory_attributes ] [ inmemory_column_clause ]
//	    | NO INMEMORY
//	    }
//	inmemory_column_clause:
//	    { INMEMORY [ inmemory_memcompress ] | NO INMEMORY } ( column [, column ] ... )
//	    [ ... ] ...
//
// An inmemory_column_clause names columns of the column list. A table
// without one (CREATE TABLE a INMEMORY INMEMORY(a)) is refused by Oracle
// at the token after the clause (ORA-00922), even AS subquery.
//
// A table takes one inmemory_table_clause: Oracle 23ai refuses a second
// with ORA-64350, at the word after its INMEMORY or at its NO INMEMORY.
//
// Each clause is kept as a DDLOption: INMEMORY with the text after the
// keyword (its attributes, or its compression and columns), or NO INMEMORY
// with its columns. tableClause tells whether the table has its
// inmemory_table_clause yet.
//
// A materialized view's page writes its inmemory_column_clause the other
// way round, "( column [, column ] ... ) { INMEMORY [ inmemory_memcompress ]
// | NO INMEMORY }", so after its columns the keyword is required: without
// one (CREATE MATERIALIZED VIEW a INMEMORY(a) AS SELECT 1 FROM dual) the
// server answers ORA-00922 where the query starts, and the page derives no
// such clause either. mview says which page is read.
func (p *Parser) parseCreateTableInmemory(stmt *nodes.CreateTableStmt, hasColumns bool, tableClause *bool, mview bool) error {
	optStart := p.pos()
	key := "INMEMORY"
	columnClause := false
	var second Token
	if p.isKeywordStr("NO") {
		key = "NO INMEMORY"
		second = p.cur
		p.advance() // consume NO
		p.advance() // consume INMEMORY
		if p.cur.Type == '(' {
			if mview {
				return p.syntaxErrorAtCur()
			}
			if err := p.parseCreateTableColumnNames(); err != nil {
				return err
			}
			columnClause = true
		}
	} else {
		p.advance() // consume INMEMORY
		second = p.cur
		memcompress := false
		if p.isCreateTableMemcompressStart() {
			if err := p.parseCreateTableMemcompress(); err != nil {
				return err
			}
			memcompress = true
		}
		if p.cur.Type == '(' {
			if err := p.parseCreateTableColumnNames(); err != nil {
				return err
			}
			columnClause = true
			if mview {
				if err := p.parseMviewInmemoryColumnKeyword(); err != nil {
					return err
				}
			}
		} else if !*tableClause {
			if err := p.parseCreateTableInmemoryAttributes(memcompress); err != nil {
				return err
			}
		}
	}
	if !columnClause {
		if *tableClause {
			return p.syntaxErrorAtTok(second)
		}
		*tableClause = true
	}
	value := strings.TrimSpace(strings.TrimPrefix(p.createTableOptionValue(optStart), key))
	appendCreateTableOption(stmt, key, value, nodes.Loc{Start: optStart, End: p.prev.End})
	if columnClause && !hasColumns {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseMviewInmemoryColumnKeyword reads the { INMEMORY [ inmemory_memcompress ]
// | NO INMEMORY } that follows a materialized view's inmemory_column_clause
// columns, and fails where neither stands.
func (p *Parser) parseMviewInmemoryColumnKeyword() error {
	switch {
	case p.isKeywordStr("INMEMORY"):
		p.advance() // consume INMEMORY
		if p.isCreateTableMemcompressStart() {
			return p.parseCreateTableMemcompress()
		}
		return nil
	case p.isKeywordStr("NO") && isCreateTableWordTok(p.peekNext(), "INMEMORY"):
		p.advance() // consume NO
		p.advance() // consume INMEMORY
		return nil
	}
	return p.syntaxErrorAtCur()
}

// ---------------------------------------------------------------------------
// LOB_storage_clause, LOB_storage_parameters
// ---------------------------------------------------------------------------

// parseCreateTableLobStorageParameters reads ( LOB_storage_parameters ).
//
//	{ { TABLESPACE tablespace | storage_clause | logging_clause } [ ... ]
//	| { ENABLE | DISABLE } STORAGE IN ROW
//	| CHUNK integer
//	| PCTVERSION integer
//	| FREEPOOLS integer
//	| LOB_retention_clause
//	| LOB_deduplicate_clause
//	| LOB_compression_clause
//	| { ENCRYPT encryption_spec | DECRYPT }
//	| { CACHE | NOCACHE | CACHE READS } [ logging_clause ]
//	} [ ... ]
//
//	LOB_retention_clause: RETENTION [ { MAX | MIN integer | AUTO | NONE } ]
//	LOB_deduplicate_clause: { DEDUPLICATE | KEEP_DUPLICATES }
//	LOB_compression_clause: { COMPRESS [ HIGH | MEDIUM | LOW ] | NOCOMPRESS }
//
// The list holds at least one parameter. A logging_clause stands only
// after CACHE, NOCACHE or CACHE READS: Oracle 23ai refuses STORE
// AS(LOGGING) and STORE AS(TABLESPACE t LOGGING) with ORA-22853 at the
// LOGGING.
func (p *Parser) parseCreateTableLobStorageParameters() error {
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	for first := true; first || p.cur.Type != ')'; first = false {
		switch {
		case p.cur.Type == kwTABLESPACE:
			p.advance() // consume TABLESPACE
			if err := p.parseCreateTableName(); err != nil {
				return err
			}
		case p.isKeywordStr("STORAGE"):
			if err := p.parseCreateTableStorageClause(); err != nil {
				return err
			}
		case p.cur.Type == kwENABLE, p.cur.Type == kwDISABLE:
			p.advance() // consume ENABLE / DISABLE
			if err := p.expectCreateTableWord("STORAGE"); err != nil {
				return err
			}
			if err := p.expectCreateTableToken(kwIN); err != nil {
				return err
			}
			if err := p.expectCreateTableToken(kwROW); err != nil {
				return err
			}
		case p.isKeywordStr("CHUNK"), p.isKeywordStr("PCTVERSION"), p.isKeywordStr("FREEPOOLS"):
			p.advance()
			if err := p.parseCreateTableInteger(); err != nil {
				return err
			}
		case p.isKeywordStr("RETENTION"):
			p.advance() // consume RETENTION
			switch {
			case p.isKeywordStr("MAX"), p.isKeywordStr("AUTO"), p.isKeywordStr("NONE"):
				p.advance()
			case p.isKeywordStr("MIN"):
				p.advance()
				if err := p.parseCreateTableInteger(); err != nil {
					return err
				}
			}
		case p.isKeywordStr("DEDUPLICATE"), p.isKeywordStr("KEEP_DUPLICATES"), p.cur.Type == kwNOCOMPRESS,
			p.isKeywordStr("DECRYPT"):
			p.advance()
		case p.cur.Type == kwCOMPRESS:
			p.advance() // consume COMPRESS
			if p.isKeywordStr("HIGH") || p.isKeywordStr("MEDIUM") || p.isKeywordStr("LOW") {
				p.advance()
			}
		case p.isKeywordStr("ENCRYPT"):
			p.advance() // consume ENCRYPT
			if err := p.parseEncryptionSpec(&nodes.ColumnDef{}); err != nil {
				return err
			}
		case p.cur.Type == kwCACHE, p.isKeywordStr("NOCACHE"):
			cache := p.cur.Type == kwCACHE
			p.advance()
			if cache && p.isKeywordStr("READS") {
				p.advance()
			}
			if p.isCreateTableLoggingStart() {
				p.advance()
			}
		default:
			return p.syntaxErrorAtCur()
		}
	}
	p.advance() // consume ')'
	return nil
}

// isCreateTableLobStorageStart reports whether a LOB_storage_clause starts
// here: LOB followed by its item list.
func (p *Parser) isCreateTableLobStorageStart() bool {
	return p.isKeywordStr("LOB") && p.peekNext().Type == '('
}

// parseCreateTableLobStorageClause reads a LOB_storage_clause.
//
//	LOB { ( LOB_item [, LOB_item ] ... ) STORE AS [ SECUREFILE | BASICFILE ]
//	      { ( LOB_storage_parameters ) | LOB_segname [ ( LOB_storage_parameters ) ] }
//	    | ( LOB_item ) STORE AS [ SECUREFILE | BASICFILE ]
//	      { LOB_segname | ( LOB_storage_parameters ) | LOB_segname ( LOB_storage_parameters ) }
//	    }
//
// Oracle 23ai takes STORE AS SECUREFILE or BASICFILE alone, and names a
// segment only for one LOB item (ORA-22855 for LOB(a, b) STORE AS s).
func (p *Parser) parseCreateTableLobStorageClause() error {
	if err := p.expectCreateTableWord("LOB"); err != nil {
		return err
	}
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	items := 0
	for {
		if err := p.parseCreateTableItem(); err != nil {
			return err
		}
		items++
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return err
	}
	if err := p.expectCreateTableWord("STORE"); err != nil {
		return err
	}
	if err := p.expectCreateTableToken(kwAS); err != nil {
		return err
	}
	return p.parseCreateTableLobStorageTarget(items == 1)
}

// parseCreateTableLobStorageTarget reads what follows STORE AS in a
// LOB_storage_clause:
//
//	[ SECUREFILE | BASICFILE ] { LOB_segname [ ( LOB_storage_parameters ) ] | ( LOB_storage_parameters ) }
func (p *Parser) parseCreateTableLobStorageTarget(segname bool) error {
	file := false
	if p.isKeywordStr("SECUREFILE") || p.isKeywordStr("BASICFILE") {
		p.advance()
		file = true
	}
	if p.cur.Type == '(' {
		return p.parseCreateTableLobStorageParameters()
	}
	if p.isCreateTableNameStart() && !p.isCreateTableLobSegmentEnd() {
		if !segname || file && p.isCreateTableRefusedFileSegname() {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume LOB_segname
		if p.cur.Type == '(' {
			return p.parseCreateTableLobStorageParameters()
		}
		return nil
	}
	if file {
		return nil
	}
	return p.syntaxErrorAtCur()
}

// isCreateTableLobSegmentEnd reports whether the word after STORE AS
// [ SECUREFILE | BASICFILE ] begins the next table property rather than
// naming a LOB segment.
func (p *Parser) isCreateTableLobSegmentEnd() bool {
	return p.isCreateTableOptionBoundary()
}

// ---------------------------------------------------------------------------
// varray_col_properties, nested_table_col_properties
// ---------------------------------------------------------------------------

// isCreateTableSubstitutableStart reports whether a
// substitutable_column_clause starts here.
func (p *Parser) isCreateTableSubstitutableStart() bool {
	switch {
	case p.isKeywordStr("ELEMENT"), p.isKeywordStr("SUBSTITUTABLE"):
		return true
	case p.cur.Type == kwNOT:
		return isCreateTableWordTok(p.peekNext(), "SUBSTITUTABLE")
	}
	return false
}

// parseCreateTableSubstitutable reads a substitutable_column_clause.
//
//	{ ELEMENT IS OF [ TYPE ] ( [ ONLY ] type )
//	| [ NOT ] SUBSTITUTABLE AT ALL LEVELS
//	}
func (p *Parser) parseCreateTableSubstitutable() error {
	if p.isKeywordStr("ELEMENT") {
		p.advance() // consume ELEMENT
		if err := p.expectCreateTableToken(kwIS); err != nil {
			return err
		}
		if err := p.expectCreateTableToken(kwOF); err != nil {
			return err
		}
		if p.cur.Type == kwTYPE {
			p.advance()
		}
		if err := p.expectCreateTableToken('('); err != nil {
			return err
		}
		if p.isKeywordStr("ONLY") {
			p.advance()
		}
		typeName, err := p.parseRequiredObjectName()
		if err != nil {
			return err
		}
		_ = typeName
		return p.expectCreateTableToken(')')
	}
	if p.cur.Type == kwNOT {
		p.advance() // consume NOT
	}
	p.advance() // consume SUBSTITUTABLE
	if err := p.expectCreateTableWord("AT"); err != nil {
		return err
	}
	if err := p.expectCreateTableToken(kwALL); err != nil {
		return err
	}
	return p.expectCreateTableWord("LEVELS")
}

// parseCreateTableVarrayColProperties reads varray_col_properties.
//
//	VARRAY varray_item
//	{ [ substitutable_column_clause ] STORE AS [ SECUREFILE | BASICFILE ] LOB
//	  { [ LOB_segname ] ( LOB_storage_parameters ) | LOB_segname }
//	| substitutable_column_clause
//	}
func (p *Parser) parseCreateTableVarrayColProperties(hasColumns bool) error {
	p.advance() // consume VARRAY
	if err := p.parseCreateTableItem(); err != nil {
		return err
	}
	substitutable := false
	if p.isCreateTableSubstitutableStart() {
		if err := p.parseCreateTableSubstitutable(); err != nil {
			return err
		}
		substitutable = true
	}
	if !p.isKeywordStr("STORE") {
		if substitutable {
			return nil
		}
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume STORE
	if err := p.expectCreateTableToken(kwAS); err != nil {
		return err
	}
	return p.parseCreateTableVarrayLobTarget()
}

// parseCreateTableVarrayLobTarget reads what follows STORE AS in
// varray_col_properties:
//
//	[ SECUREFILE | BASICFILE ] LOB { [ LOB_segname ] ( LOB_storage_parameters ) | LOB_segname }
func (p *Parser) parseCreateTableVarrayLobTarget() error {
	file := false
	if p.isKeywordStr("SECUREFILE") || p.isKeywordStr("BASICFILE") {
		p.advance()
		file = true
	}
	if err := p.expectCreateTableWord("LOB"); err != nil {
		return err
	}
	if p.cur.Type == '(' {
		return p.parseCreateTableLobStorageParameters()
	}
	if file && p.isCreateTableRefusedFileSegname() {
		return p.syntaxErrorAtCur()
	}
	if err := p.parseCreateTableName(); err != nil {
		return err
	}
	if p.cur.Type == '(' {
		return p.parseCreateTableLobStorageParameters()
	}
	return nil
}

// parseCreateTableNestedTableColProperties reads nested_table_col_properties.
//
//	NESTED TABLE { nested_item | COLUMN_VALUE }
//	[ substitutable_column_clause ]
//	[ LOCAL | GLOBAL ] STORE AS storage_table
//	[ ( ( object_properties ) [ physical_properties ] [ column_properties ] ) ]
//	[ RETURN [ AS ] { LOCATOR | VALUE } ]
//
// The storage table's parenthesized properties are kept whole: Oracle
// answers what is in them after reading the statement (ORA-22912 for
// STORE AS a((a) zzz)), not as a syntax error.
func (p *Parser) parseCreateTableNestedTableColProperties(hasColumns bool) error {
	p.advance() // consume NESTED
	if err := p.expectCreateTableToken(kwTABLE); err != nil {
		return err
	}
	if p.cur.Type == kwCOLUMN_VALUE {
		p.advance()
	} else if err := p.parseCreateTableNestedItem(); err != nil {
		return err
	}
	if p.isCreateTableSubstitutableStart() {
		if err := p.parseCreateTableSubstitutable(); err != nil {
			return err
		}
	}
	if p.cur.Type == kwGLOBAL || p.isKeywordStr("LOCAL") {
		p.advance()
	}
	if err := p.expectCreateTableWord("STORE"); err != nil {
		return err
	}
	if err := p.expectCreateTableToken(kwAS); err != nil {
		return err
	}
	if err := p.parseCreateTableName(); err != nil {
		return err
	}
	if p.cur.Type == '(' {
		p.skipParenthesized()
		if p.prev.Type != ')' {
			return p.syntaxErrorAtCur()
		}
	}
	if p.cur.Type == kwRETURN {
		p.advance() // consume RETURN
		if p.cur.Type == kwAS {
			p.advance()
		}
		if !p.isKeywordStr("LOCATOR") && !p.isKeywordStr("VALUE") {
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
	return nil
}

// parseCreateTableNestedItem reads the nested_item of a nested table: a
// column, or an attribute path of an object column.
func (p *Parser) parseCreateTableNestedItem() error {
	if !p.isCreateTableNameStart() {
		return p.syntaxErrorAtCur()
	}
	return p.parseCreateTableItem()
}

// parseCreateTableItem reads a LOB, VARRAY or nested table item: a column,
// or an attribute path of an object column.
func (p *Parser) parseCreateTableItem() error {
	item, err := p.parseColumnPath()
	if err != nil {
		return err
	}
	_ = item
	return nil
}

// ---------------------------------------------------------------------------
// storage_clause
// ---------------------------------------------------------------------------

// parseCreateTableStorageClause reads a storage_clause.
//
//	STORAGE ( { INITIAL size_clause | NEXT size_clause | MINEXTENTS integer
//	          | MAXEXTENTS { integer | UNLIMITED } | maxsize_clause
//	          | PCTINCREASE integer | FREELISTS integer | FREELIST GROUPS integer
//	          | OPTIMAL [ size_clause | NULL ]
//	          | BUFFER_POOL { KEEP | RECYCLE | DEFAULT }
//	          | FLASH_CACHE { KEEP | NONE | DEFAULT }
//	          | CELL_FLASH_CACHE { KEEP | NONE | DEFAULT }
//	          | ENCRYPT
//	          } ... )
//
//	maxsize_clause: MAXSIZE { UNLIMITED | size_clause }
//	size_clause: integer [ K | M | G | T | P | E ]
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/storage_clause.html
func (p *Parser) parseCreateTableStorageClause() error {
	p.advance() // consume STORAGE
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	for first := true; first || p.cur.Type != ')'; first = false {
		switch {
		case p.cur.Type == kwINITIAL, p.isKeywordStr("NEXT"):
			p.advance()
			if err := p.parseCreateTableSize(); err != nil {
				return err
			}
		case p.isKeywordStr("MINEXTENTS"), p.isKeywordStr("PCTINCREASE"), p.isKeywordStr("FREELISTS"):
			p.advance()
			if err := p.parseCreateTableInteger(); err != nil {
				return err
			}
		case p.cur.Type == kwMAXEXTENTS:
			p.advance()
			if p.isKeywordStr("UNLIMITED") {
				p.advance()
			} else if err := p.parseCreateTableInteger(); err != nil {
				return err
			}
		case p.isKeywordStr("MAXSIZE"):
			p.advance()
			if p.isKeywordStr("UNLIMITED") {
				p.advance()
			} else if err := p.parseCreateTableSize(); err != nil {
				return err
			}
		case p.isKeywordStr("FREELIST"):
			p.advance()
			if err := p.expectCreateTableWord("GROUPS"); err != nil {
				return err
			}
			if err := p.parseCreateTableInteger(); err != nil {
				return err
			}
		case p.isKeywordStr("OPTIMAL"):
			p.advance()
			if p.cur.Type == kwNULL {
				p.advance()
			} else if p.cur.Type == tokICONST {
				if err := p.parseCreateTableSize(); err != nil {
					return err
				}
			}
		case p.isKeywordStr("BUFFER_POOL"):
			p.advance()
			if !p.isKeywordStr("KEEP") && !p.isKeywordStr("RECYCLE") && p.cur.Type != kwDEFAULT {
				return p.syntaxErrorAtCur()
			}
			p.advance()
		case p.isKeywordStr("FLASH_CACHE"), p.isKeywordStr("CELL_FLASH_CACHE"):
			p.advance()
			if !p.isKeywordStr("KEEP") && !p.isKeywordStr("NONE") && p.cur.Type != kwDEFAULT {
				return p.syntaxErrorAtCur()
			}
			p.advance()
		case p.isKeywordStr("ENCRYPT"):
			p.advance()
		default:
			return p.syntaxErrorAtCur()
		}
	}
	p.advance() // consume ')'
	return nil
}

// parseCreateTableSize reads a size_clause: integer [ K | M | G | T | P | E ].
func (p *Parser) parseCreateTableSize() error {
	if err := p.parseCreateTableInteger(); err != nil {
		return err
	}
	switch {
	case p.isKeywordStr("K"), p.isKeywordStr("M"), p.isKeywordStr("G"),
		p.isKeywordStr("T"), p.isKeywordStr("P"), p.isKeywordStr("E"):
		p.advance()
	}
	return nil
}

// ---------------------------------------------------------------------------
// ilm_clause
// ---------------------------------------------------------------------------

// parseCreateTableIlmClause reads an ilm_clause.
//
//	ILM { ADD POLICY ilm_policy_clause
//	    | { DELETE | ENABLE | DISABLE } POLICY ilm_policy_name
//	    | DELETE_ALL | ENABLE_ALL | DISABLE_ALL
//	    }
func (p *Parser) parseCreateTableIlmClause() error {
	p.advance() // consume ILM
	switch {
	case p.cur.Type == kwADD:
		p.advance() // consume ADD
		if err := p.expectCreateTableWord("POLICY"); err != nil {
			return err
		}
		return p.parseCreateTableIlmPolicy()
	case p.cur.Type == kwDELETE, p.cur.Type == kwENABLE, p.cur.Type == kwDISABLE:
		p.advance()
		if err := p.expectCreateTableWord("POLICY"); err != nil {
			return err
		}
		err := p.parseCreateTableName()
		return err
	case p.isKeywordStr("DELETE_ALL"), p.isKeywordStr("ENABLE_ALL"), p.isKeywordStr("DISABLE_ALL"):
		p.advance()
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseCreateTableIlmPolicy reads an ilm_policy_clause.
//
//	ilm_compression_policy:
//	    table_compression { SEGMENT | GROUP | ROW }
//	    { AFTER ilm_time_period OF { NO ACCESS | NO MODIFICATION | CREATION } | ON function_name }
//	ilm_tiering_policy:
//	    TIER TO tablespace [ READ ONLY ] [ SEGMENT | GROUP ] [ ON function_name ]
//	ilm_inmemory_policy:
//	    { SET INMEMORY [ inmemory_attributes ] | MODIFY INMEMORY inmemory_memcompress | NO INMEMORY }
//	    [ SEGMENT ]
//	    { AFTER ilm_time_period OF { NO ACCESS | NO MODIFICATION | CREATION } | ON function_name }
//
// Oracle 23ai takes ROW STORE COMPRESS, COLUMN STORE COMPRESS and COMPRESS
// as a compression policy's action and refuses NOCOMPRESS (ORA-38339), a
// REBUILD or MERGE action, and COMPRESS ADVANCED (ORA-38322). A tiering
// policy ends at its scope or its ON: AFTER there is ORA-00922, as is a
// ROW scope.
func (p *Parser) parseCreateTableIlmPolicy() error {
	switch {
	case p.isKeywordStr("TIER"):
		p.advance() // consume TIER
		if err := p.expectCreateTableToken(kwTO); err != nil {
			return err
		}
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
		if p.cur.Type == kwREAD {
			p.advance() // consume READ
			if err := p.expectCreateTableWord("ONLY"); err != nil {
				return err
			}
		}
		if p.isKeywordStr("SEGMENT") || p.cur.Type == kwGROUP {
			p.advance()
		}
		if p.cur.Type == kwON {
			return p.parseCreateTableIlmFunction()
		}
		return nil
	case p.cur.Type == kwSET:
		p.advance() // consume SET
		if err := p.expectCreateTableWord("INMEMORY"); err != nil {
			return err
		}
		if err := p.parseCreateTableInmemoryAttributes(false); err != nil {
			return err
		}
		return p.parseCreateTableIlmInmemoryCondition()
	case p.cur.Type == kwMODIFY:
		p.advance() // consume MODIFY
		if err := p.expectCreateTableWord("INMEMORY"); err != nil {
			return err
		}
		if !p.isCreateTableMemcompressStart() {
			return p.syntaxErrorAtCur()
		}
		if err := p.parseCreateTableMemcompress(); err != nil {
			return err
		}
		return p.parseCreateTableIlmInmemoryCondition()
	case p.isKeywordStr("NO") && isCreateTableWordTok(p.peekNext(), "INMEMORY"):
		p.advance() // consume NO
		p.advance() // consume INMEMORY
		return p.parseCreateTableIlmInmemoryCondition()
	case p.cur.Type == kwROW, p.isKeywordStr("COLUMN"), p.cur.Type == kwCOMPRESS:
		if err := p.parseCreateTableIlmCompression(); err != nil {
			return err
		}
		if !p.isKeywordStr("SEGMENT") && p.cur.Type != kwGROUP && p.cur.Type != kwROW {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return p.parseCreateTableIlmCondition()
	}
	return p.syntaxErrorAtCur()
}

// parseCreateTableIlmCompression reads the table_compression of an ILM
// compression policy.
//
//	{ ROW STORE COMPRESS [ BASIC | ADVANCED ]
//	| COLUMN STORE COMPRESS FOR { QUERY | ARCHIVE } [ LOW | HIGH ]
//	| COMPRESS
//	}
func (p *Parser) parseCreateTableIlmCompression() error {
	switch {
	case p.cur.Type == kwCOMPRESS:
		p.advance()
		return nil
	case p.cur.Type == kwROW:
		p.advance() // consume ROW
		if err := p.expectCreateTableWord("STORE"); err != nil {
			return err
		}
		if err := p.expectCreateTableToken(kwCOMPRESS); err != nil {
			return err
		}
		if p.isKeywordStr("BASIC") || p.isKeywordStr("ADVANCED") {
			p.advance()
		}
		return nil
	}
	p.advance() // consume COLUMN
	if err := p.expectCreateTableWord("STORE"); err != nil {
		return err
	}
	if err := p.expectCreateTableToken(kwCOMPRESS); err != nil {
		return err
	}
	if err := p.expectCreateTableToken(kwFOR); err != nil {
		return err
	}
	if !p.isKeywordStr("QUERY") && !p.isKeywordStr("ARCHIVE") {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.isKeywordStr("LOW") || p.isKeywordStr("HIGH") {
		p.advance()
	}
	return nil
}

// parseCreateTableIlmInmemoryCondition reads the end of an ILM in-memory
// policy: [ SEGMENT ] and its condition. Oracle 23ai takes SET INMEMORY ON
// f without SEGMENT and refuses GROUP there (ORA-03048).
func (p *Parser) parseCreateTableIlmInmemoryCondition() error {
	if p.isKeywordStr("SEGMENT") {
		p.advance()
	}
	return p.parseCreateTableIlmCondition()
}

// parseCreateTableIlmCondition reads an ILM policy's condition.
//
//	{ AFTER ilm_time_period OF { NO ACCESS | NO MODIFICATION | CREATION }
//	| ON function_name
//	}
//	ilm_time_period: integer [ DAY | DAYS | MONTH | MONTHS | YEAR | YEARS ]
//
// Oracle 23ai takes a period without its unit (AFTER 1 OF CREATION).
func (p *Parser) parseCreateTableIlmCondition() error {
	if p.cur.Type == kwON {
		return p.parseCreateTableIlmFunction()
	}
	if p.cur.Type != kwAFTER {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume AFTER
	if err := p.parseCreateTableInteger(); err != nil {
		return err
	}
	switch {
	case p.isKeywordStr("DAY"), p.isKeywordStr("DAYS"), p.isKeywordStr("MONTH"),
		p.isKeywordStr("MONTHS"), p.isKeywordStr("YEAR"), p.isKeywordStr("YEARS"):
		p.advance()
	}
	if err := p.expectCreateTableToken(kwOF); err != nil {
		return err
	}
	switch {
	case p.isKeywordStr("CREATION"):
		p.advance()
		return nil
	case p.isKeywordStr("NO"):
		p.advance() // consume NO
		if !p.isKeywordStr("ACCESS") && !p.isKeywordStr("MODIFICATION") {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseCreateTableIlmFunction reads ON function_name.
func (p *Parser) parseCreateTableIlmFunction() error {
	p.advance() // consume ON
	err := p.parseCreateTableName()
	return err
}

// isCreateTableRefusedFileSegname reports whether the LOB segment name
// after SECUREFILE or BASICFILE is one Oracle 23ai refuses there with
// ORA-00922: the one-letter names A, D, E, G, K, M, P and T, quoted or
// not (B, C, F, R, X, Y, Z and longer names are taken).
func (p *Parser) isCreateTableRefusedFileSegname() bool {
	switch p.cur.Str {
	case "A", "D", "E", "G", "K", "M", "P", "T":
		return p.isIdentLike()
	}
	return false
}

// isCreateTableConstraintStateStart reports whether a constraint_state
// starts here (NOT DEFERRABLE aside, which starts with NOT).
func (p *Parser) isCreateTableConstraintStateStart() bool {
	switch {
	case p.cur.Type == kwDEFERRABLE, p.cur.Type == kwINITIALLY, p.cur.Type == kwRELY,
		p.cur.Type == kwENABLE, p.cur.Type == kwDISABLE, p.cur.Type == kwVALIDATE:
		return true
	case p.cur.Type == kwUSING:
		return p.peekNext().Type == kwINDEX
	}
	return p.isKeywordStr("NORELY") || p.isKeywordStr("NOVALIDATE")
}

// ---------------------------------------------------------------------------
// external_table_clause
// ---------------------------------------------------------------------------

// parseCreateTableExternalData reads the parenthesized part of an
// external_table_clause, after ORGANIZATION EXTERNAL.
//
//	( [ TYPE access_driver_type ]
//	  [ DEFAULT DIRECTORY directory ]
//	  [ ACCESS PARAMETERS { ( opaque_format_spec ) | USING CLOB subquery } ]
//	  [ LOCATION ( [ directory: ] 'location_specifier'
//	               [, [ directory: ] 'location_specifier' ] ... ) ]
//	)
//
// The opaque_format_spec is the access driver's, kept whole.
func (p *Parser) parseCreateTableExternalData() error {
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	if p.cur.Type == kwTYPE {
		p.advance() // consume TYPE
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
	}
	if p.cur.Type == kwDEFAULT {
		p.advance() // consume DEFAULT
		if err := p.expectCreateTableToken(kwDIRECTORY); err != nil {
			return err
		}
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
	}
	if p.cur.Type == kwACCESS {
		p.advance() // consume ACCESS
		if err := p.expectCreateTableWord("PARAMETERS"); err != nil {
			return err
		}
		if p.cur.Type == '(' {
			p.skipParenthesized()
			if p.prev.Type != ')' {
				return p.syntaxErrorAtCur()
			}
		} else {
			if err := p.expectCreateTableToken(kwUSING); err != nil {
				return err
			}
			if err := p.expectCreateTableToken(kwCLOB); err != nil {
				return err
			}
			query, err := p.parseSelectStmt()
			if err != nil {
				return err
			}
			if query == nil {
				return p.syntaxErrorAtCur()
			}
		}
	}
	if p.isKeywordStr("LOCATION") {
		p.advance() // consume LOCATION
		if err := p.expectCreateTableToken('('); err != nil {
			return err
		}
		for {
			if p.cur.Type != tokSCONST {
				if err := p.parseCreateTableName(); err != nil {
					return err
				}
				if err := p.expectCreateTableToken(':'); err != nil {
					return err
				}
			}
			if err := p.expectCreateTableToken(tokSCONST); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
		if err := p.expectCreateTableToken(')'); err != nil {
			return err
		}
	}
	return p.expectCreateTableToken(')')
}

// ---------------------------------------------------------------------------
// column_definition, virtual_column_definition, period_definition
// ---------------------------------------------------------------------------

// isCreateTableColumnClauseStart reports whether the word after a column's
// name starts one of its clauses rather than naming its datatype: COLLATE,
// ENCRYPT, ANNOTATIONS, or SCOPE IS.
func (p *Parser) isCreateTableColumnClauseStart() bool {
	switch {
	case p.isKeywordStr("COLLATE"), p.isKeywordStr("ENCRYPT"), p.isKeywordStr("ANNOTATIONS"):
		return true
	case p.isKeywordStr("SCOPE"):
		return p.peekNext().Type == kwIS
	}
	return false
}

// parseCreateTableAnnotations reads a column's annotations_clause.
//
//	ANNOTATIONS ( [ ADD ] annotation_name [ 'annotation_value' ] [, ...] )
//
// CREATE TABLE adds annotations: Oracle 23ai refuses DROP there, and a name
// given twice, at the second.
func (p *Parser) parseCreateTableAnnotations() (*nodes.List, error) {
	anns := &nodes.List{}
	p.advance() // consume ANNOTATIONS
	if err := p.expectCreateTableToken('('); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for {
		start := p.pos()
		op := ""
		if p.cur.Type == kwADD {
			op = "ADD"
			p.advance()
		}
		if seen[p.cur.Str] && p.isCreateTableNameStart() {
			return nil, p.syntaxErrorAtCur()
		}
		seen[p.cur.Str] = true
		name := p.cur.Str
		if err := p.parseCreateTableName(); err != nil {
			return nil, err
		}
		opt := &nodes.DDLOption{Key: op, Value: name, Loc: nodes.Loc{Start: start}}
		if p.cur.Type == tokSCONST {
			opt.Items = &nodes.List{Items: []nodes.Node{&nodes.StringLiteral{Val: p.cur.Str, Loc: nodes.Loc{Start: p.cur.Loc, End: p.cur.End}}}}
			p.advance()
		}
		opt.Loc.End = p.prev.End
		anns.Items = append(anns.Items, opt)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return nil, err
	}
	return anns, nil
}

// parseCreateTableEditionClauses reads a virtual column's
// evaluation_edition_clause and unusable_editions_clause.
//
//	evaluation_edition_clause:
//	    EVALUATE USING { CURRENT EDITION | EDITION edition | NULL EDITION }
//	unusable_editions_clause:
//	    [ UNUSABLE BEFORE { CURRENT EDITION | EDITION edition } ]
//	    [ UNUSABLE BEGINNING WITH { CURRENT EDITION | EDITION edition | NULL EDITION } ]
func (p *Parser) parseCreateTableEditionClauses() error {
	if p.isKeywordStr("EVALUATE") {
		p.advance() // consume EVALUATE
		if err := p.expectCreateTableToken(kwUSING); err != nil {
			return err
		}
		if err := p.parseCreateTableEdition(true); err != nil {
			return err
		}
	}
	before := false
	for p.isKeywordStr("UNUSABLE") {
		p.advance() // consume UNUSABLE
		switch {
		case !before && p.isKeywordStr("BEFORE"):
			p.advance()
			before = true
			if err := p.parseCreateTableEdition(false); err != nil {
				return err
			}
		case p.isKeywordStr("BEGINNING"):
			p.advance()
			if err := p.expectCreateTableToken(kwWITH); err != nil {
				return err
			}
			return p.parseCreateTableEdition(true)
		default:
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}

// parseCreateTableEdition reads { CURRENT EDITION | EDITION edition |
// NULL EDITION }, NULL EDITION where null allows it.
func (p *Parser) parseCreateTableEdition(null bool) error {
	switch {
	case p.cur.Type == kwCURRENT, null && p.cur.Type == kwNULL:
		p.advance()
		return p.expectCreateTableWord("EDITION")
	case p.isKeywordStr("EDITION"):
		p.advance()
		return p.parseCreateTableName()
	}
	return p.syntaxErrorAtCur()
}

// isCreateTablePeriodStart reports whether a period_definition starts here.
func (p *Parser) isCreateTablePeriodStart() bool {
	return p.isKeywordStr("PERIOD") && p.peekNext().Type == kwFOR
}

// parseCreateTablePeriod reads a period_definition, which adds the hidden
// start and end columns of its valid time, and keeps it as a PERIOD FOR
// DDLOption with its text.
//
//	PERIOD FOR valid_time_column [ ( start_time_column , end_time_column ) ]
func (p *Parser) parseCreateTablePeriod(stmt *nodes.CreateTableStmt) error {
	start := p.pos()
	if err := p.parseCreateTablePeriodDefinition(); err != nil {
		return err
	}
	appendCreateTableOption(stmt, "PERIOD FOR", p.createTableOptionValue(start+len("PERIOD FOR")), nodes.Loc{Start: start, End: p.prev.End})
	return nil
}

// parseCreateTablePeriodDefinition reads the words of a period_definition.
func (p *Parser) parseCreateTablePeriodDefinition() error {
	p.advance() // consume PERIOD
	p.advance() // consume FOR
	if err := p.parseCreateTableName(); err != nil {
		return err
	}
	if p.cur.Type != '(' {
		return nil
	}
	p.advance() // consume '('
	if err := p.parseCreateTableName(); err != nil {
		return err
	}
	if err := p.expectCreateTableToken(','); err != nil {
		return err
	}
	if err := p.parseCreateTableName(); err != nil {
		return err
	}
	return p.expectCreateTableToken(')')
}

// ---------------------------------------------------------------------------
// object_table, XMLType_table
// ---------------------------------------------------------------------------

// parseCreateTableObjectTable reads the part of an object_table or an
// XMLType_table before its physical and table properties, after OF.
//
//	object_table:
//	    OF [ schema. ] object_type
//	    [ object_table_substitution ]
//	    [ ( object_properties ) ]
//	    [ ON COMMIT { DELETE | PRESERVE } ROWS ]
//	    [ oid_clause ] [ oid_index_clause ]
//	XMLType_table:
//	    OF XMLTYPE
//	    [ ( object_properties ) ]
//	    [ XMLSchema_spec ]
//	    [ ON COMMIT { DELETE | PRESERVE } ROWS ]
//	    [ oid_clause ] [ oid_index_clause ]
//
//	object_table_substitution: [ NOT ] SUBSTITUTABLE AT ALL LEVELS
//	oid_clause: OBJECT IDENTIFIER IS { SYSTEM GENERATED | PRIMARY KEY }
//	oid_index_clause:
//	    OIDINDEX [ index ] ( [ physical_attributes_clause | TABLESPACE tablespace ] [ ... ] )
//
// Oracle 23ai takes OBJECT IDENTIFIER without IS, and refuses STORE AS
// (XMLType_storage) after OF XMLTYPE at STORE.
func (p *Parser) parseCreateTableObjectTable(stmt *nodes.CreateTableStmt) error {
	p.advance() // consume OF
	if !p.isCreateTableNameStart() {
		return p.syntaxErrorAtCur()
	}
	typeName, err := p.parseObjectName()
	if err != nil {
		return err
	}
	stmt.OfType = typeName
	xmlType := typeName.Schema == "" && typeName.Name == "XMLTYPE"
	if !xmlType && (p.isKeywordStr("SUBSTITUTABLE") || p.cur.Type == kwNOT && isCreateTableWordTok(p.peekNext(), "SUBSTITUTABLE")) {
		if err := p.parseCreateTableSubstitutable(); err != nil {
			return err
		}
	}
	if p.cur.Type == '(' {
		if err := p.parseCreateTableObjectProperties(stmt); err != nil {
			return err
		}
	}
	if xmlType && (p.isKeywordStr("XMLSCHEMA") || p.isKeywordStr("ELEMENT")) {
		if err := p.parseCreateTableXMLSchemaSpec(); err != nil {
			return err
		}
	}
	if p.isKeywordStr("OBJECT") {
		p.advance() // consume OBJECT
		if err := p.expectCreateTableWord("IDENTIFIER"); err != nil {
			return err
		}
		if p.cur.Type == kwIS {
			p.advance()
		}
		switch {
		case p.isKeywordStr("SYSTEM"):
			p.advance()
			if err := p.expectCreateTableWord("GENERATED"); err != nil {
				return err
			}
		case p.cur.Type == kwPRIMARY:
			p.advance()
			if err := p.expectCreateTableToken(kwKEY); err != nil {
				return err
			}
		default:
			return p.syntaxErrorAtCur()
		}
	}
	if p.isKeywordStr("OIDINDEX") {
		p.advance() // consume OIDINDEX
		if p.cur.Type != '(' {
			if err := p.parseCreateTableName(); err != nil {
				return err
			}
		}
		if err := p.expectCreateTableToken('('); err != nil {
			return err
		}
		for p.cur.Type != ')' {
			switch {
			case p.cur.Type == kwTABLESPACE:
				p.advance()
				if err := p.parseCreateTableName(); err != nil {
					return err
				}
			case p.cur.Type == kwPCTFREE, p.isKeywordStr("PCTUSED"), p.isKeywordStr("INITRANS"), p.isKeywordStr("MAXTRANS"):
				p.advance()
				if err := p.parseCreateTableInteger(); err != nil {
					return err
				}
			case p.isKeywordStr("STORAGE"):
				if err := p.parseCreateTableStorageClause(); err != nil {
					return err
				}
			default:
				return p.syntaxErrorAtCur()
			}
		}
		p.advance() // consume ')'
	}
	return nil
}

// parseCreateTableXMLSchemaSpec reads an XMLSchema_spec.
//
//	[ XMLSCHEMA XMLSchema_URL ]
//	ELEMENT { element | XMLSchema_URL # element }
//	[ STORE ALL VARRAYS AS { LOBS | TABLES } ]
//	[ { ALLOW | DISALLOW } NONSCHEMA ]
//	[ { ALLOW | DISALLOW } ANYSCHEMA ]
func (p *Parser) parseCreateTableXMLSchemaSpec() error {
	if p.isKeywordStr("XMLSCHEMA") {
		p.advance()
		if err := p.expectCreateTableToken(tokSCONST); err != nil {
			return err
		}
	}
	if err := p.expectCreateTableWord("ELEMENT"); err != nil {
		return err
	}
	if p.cur.Type == tokSCONST {
		p.advance()
		if err := p.expectCreateTableToken('#'); err != nil {
			return err
		}
	}
	if err := p.parseCreateTableName(); err != nil {
		return err
	}
	if p.isKeywordStr("STORE") {
		p.advance()
		if err := p.expectCreateTableToken(kwALL); err != nil {
			return err
		}
		if err := p.expectCreateTableWord("VARRAYS"); err != nil {
			return err
		}
		if err := p.expectCreateTableToken(kwAS); err != nil {
			return err
		}
		if !p.isKeywordStr("LOBS") && !p.isKeywordStr("TABLES") {
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
	for _, schema := range []string{"NONSCHEMA", "ANYSCHEMA"} {
		if (p.isKeywordStr("ALLOW") || p.isKeywordStr("DISALLOW")) && isCreateTableWordTok(p.peekNext(), schema) {
			p.advance()
			p.advance()
		}
	}
	return nil
}

// parseCreateTableObjectProperties reads ( object_properties ).
//
//	{ { column | attribute }
//	  [ DEFAULT expr ]
//	  [ { inline_constraint | inline_ref_constraint } [ ... ] ]
//	| out_of_line_constraint
//	| out_of_line_ref_constraint
//	| supplemental_logging_props
//	} [, ... ]
//
// An attribute takes a default or a constraint: Oracle 23ai refuses one
// alone (ORA-00922 at the word after it).
func (p *Parser) parseCreateTableObjectProperties(stmt *nodes.CreateTableStmt) error {
	p.advance() // consume '('
	for {
		switch {
		case p.isKeywordStr("SUPPLEMENTAL"):
			if err := p.parseCreateTableSupplementalLogging(); err != nil {
				return err
			}
		case p.isOutOfLineRefConstraintStart():
			tc, err := p.parseOutOfLineRefConstraint()
			if err != nil {
				return err
			}
			stmt.Constraints.Items = append(stmt.Constraints.Items, tc)
		case p.isTableConstraintStart():
			tc, err := p.parseTableConstraint()
			if err != nil {
				return err
			}
			if tc == nil {
				return p.syntaxErrorAtCur()
			}
			stmt.Constraints.Items = append(stmt.Constraints.Items, tc)
		default:
			if err := p.parseCreateTableObjectAttribute(stmt); err != nil {
				return err
			}
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return p.expectCreateTableToken(')')
}

// parseCreateTableObjectAttribute reads an attribute of object_properties:
// its name, then [ DEFAULT expr ] and its constraints, at least one of
// them.
func (p *Parser) parseCreateTableObjectAttribute(stmt *nodes.CreateTableStmt) error {
	start := p.pos()
	if !p.isCreateTableNameStart() {
		return p.syntaxErrorAtCur()
	}
	col := &nodes.ColumnDef{Constraints: &nodes.List{}, Loc: nodes.Loc{Start: start}}
	var err error
	if col.Name, err = p.parseIdentifier(); err != nil {
		return err
	}
	if p.cur.Type == kwDEFAULT {
		p.advance()
		if col.Default, err = p.parseExpr(); err != nil {
			return err
		}
		if col.Default == nil {
			return p.syntaxErrorAtCur()
		}
	}
	for {
		switch {
		case p.cur.Type == kwNOT && p.peekNext().Type == kwNULL, p.cur.Type == kwNULL:
			if col.Null || col.NotNull {
				return p.syntaxErrorAtCur()
			}
			if p.cur.Type == kwNOT {
				p.advance()
				col.NotNull = true
			} else {
				col.Null = true
			}
			p.advance() // consume NULL
			var cs constraintState
			if err := p.parseConstraintState(&cs); err != nil {
				return err
			}
		case p.cur.Type == kwCONSTRAINT:
			cc, err := p.parseColumnConstraint()
			if err != nil {
				return err
			}
			if cc != nil {
				col.Constraints.Items = append(col.Constraints.Items, cc)
			}
		case p.cur.Type == kwPRIMARY, p.cur.Type == kwUNIQUE, p.cur.Type == kwCHECK, p.cur.Type == kwREFERENCES:
			cc, err := p.parseColumnConstraintInline()
			if err != nil {
				return err
			}
			if cc != nil {
				col.Constraints.Items = append(col.Constraints.Items, cc)
			}
		case p.cur.Type == kwWITH && p.peekNext().Type == kwROWID:
			start := p.pos()
			p.advance() // consume WITH
			p.advance() // consume ROWID
			col.Constraints.Items = append(col.Constraints.Items, &nodes.ColumnConstraint{Type: nodes.CONSTRAINT_WITH_ROWID, Loc: nodes.Loc{Start: start, End: p.prev.End}})
		case p.isKeywordStr("SCOPE") && p.peekNext().Type == kwIS:
			start := p.pos()
			p.advance() // consume SCOPE
			p.advance() // consume IS
			table, err := p.parseRefScopeTable()
			if err != nil {
				return err
			}
			col.Constraints.Items = append(col.Constraints.Items, &nodes.ColumnConstraint{Type: nodes.CONSTRAINT_SCOPE, RefTable: table, Loc: nodes.Loc{Start: start, End: p.prev.End}})
		default:
			if col.Default == nil && !col.Null && !col.NotNull && col.Constraints.Len() == 0 {
				return p.syntaxErrorAtCur()
			}
			col.Loc.End = p.prev.End
			stmt.Columns.Items = append(stmt.Columns.Items, col)
			return nil
		}
	}
}

// parseCreateTableSupplementalLogging reads supplemental_logging_props.
//
//	SUPPLEMENTAL LOG { supplemental_log_grp_clause | supplemental_id_key_clause }
//	supplemental_log_grp_clause:
//	    GROUP log_group ( column [ NO LOG ] [, column [ NO LOG ] ] ... ) [ ALWAYS ]
//	supplemental_id_key_clause:
//	    DATA ( { ALL | PRIMARY KEY | UNIQUE | FOREIGN KEY } [, ...] ) COLUMNS
func (p *Parser) parseCreateTableSupplementalLogging() error {
	p.advance() // consume SUPPLEMENTAL
	if err := p.expectCreateTableWord("LOG"); err != nil {
		return err
	}
	if p.isKeywordStr("DATA") {
		p.advance()
		if err := p.expectCreateTableToken('('); err != nil {
			return err
		}
		for {
			switch {
			case p.cur.Type == kwALL, p.cur.Type == kwUNIQUE:
				p.advance()
			case p.cur.Type == kwPRIMARY, p.cur.Type == kwFOREIGN:
				p.advance()
				if err := p.expectCreateTableToken(kwKEY); err != nil {
					return err
				}
			default:
				return p.syntaxErrorAtCur()
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if err := p.expectCreateTableToken(')'); err != nil {
			return err
		}
		return p.expectCreateTableToken(kwCOLUMNS)
	}
	if err := p.expectCreateTableToken(kwGROUP); err != nil {
		return err
	}
	if err := p.parseCreateTableName(); err != nil {
		return err
	}
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	for {
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
		if p.isKeywordStr("NO") && isCreateTableWordTok(p.peekNext(), "LOG") {
			p.advance()
			p.advance()
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return err
	}
	if p.cur.Type == kwALWAYS {
		p.advance()
	}
	return nil
}
