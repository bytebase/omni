package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// SHOW
// ---------------------------------------------------------------------------

// parseShow parses a SHOW statement.
//
// Strategy: consume SHOW, then dispatch on the next 1-3 keywords to identify
// the variant. For common variants (TABLES, DATABASES, COLUMNS, CREATE TABLE,
// VARIABLES, PARTITIONS, GRANTS) populate the ShowStmt fields properly.
// For all other variants, capture the remaining tokens as raw text in Args.
//
// On entry, SHOW has NOT yet been consumed (cur == kwSHOW).
func (p *Parser) parseShow() (ast.Node, error) {
	startLoc := p.cur.Loc
	p.advance() // consume SHOW

	stmt := &ast.ShowStmt{Loc: startLoc}

	// Optional FULL / EXTENDED / BRIEF / ALL modifiers before the variant keyword.
	if p.cur.Kind == kwFULL {
		stmt.Full = true
		p.advance()
	} else if p.cur.Kind == kwEXTENDED {
		stmt.Extended = true
		p.advance()
	} else if p.cur.Kind == kwBRIEF {
		// SHOW BRIEF CREATE TABLE ... — record and consume
		stmt.Args = "BRIEF"
		p.advance()
	} else if p.cur.Kind == kwALL {
		// SHOW ALL GRANTS — record ALL in args, let dispatch handle GRANTS
		stmt.Args = "ALL"
		p.advance()
	} else if p.cur.Kind == kwTEMPORARY {
		// SHOW TEMPORARY PARTITIONS — record TEMPORARY
		stmt.Args = "TEMPORARY"
		p.advance()
	}

	// SHOW JOB (T8.1) — special-case routing
	if p.cur.Kind == kwJOB {
		return p.parseShowJob(startLoc)
	}

	// SHOW ANALYZE / SHOW STATS / SHOW CONSTRAINTS (T8.3) — special-case routing
	if p.cur.Kind == kwANALYZE {
		all := stmt.Args == "ALL"
		return p.parseShowAnalyze(startLoc, all, false)
	}
	if p.cur.Kind == kwQUEUED && p.peekNext().Kind == kwANALYZE {
		p.advance() // consume QUEUED
		return p.parseShowAnalyze(startLoc, false, true)
	}
	if p.cur.Kind == kwCONSTRAINTS {
		return p.parseShowConstraints(startLoc)
	}
	if p.cur.Kind == kwSTATS {
		return p.parseShowStats(startLoc, "")
	}
	// SHOW COLUMN/INDEX/PARTITION/TABLE STATS
	if p.peekNext().Kind == kwSTATS {
		switch p.cur.Kind {
		case kwCOLUMN:
			p.advance()
			return p.parseShowStats(startLoc, "COLUMN")
		case kwINDEX:
			p.advance()
			return p.parseShowStats(startLoc, "INDEX")
		case kwPARTITION:
			p.advance()
			return p.parseShowStats(startLoc, "PARTITION")
		case kwTABLE:
			p.advance()
			return p.parseShowStats(startLoc, "TABLE")
		}
	}

	// SHOW [ALL] ROUTINE LOAD (T6.2) — special-case routing
	if p.cur.Kind == kwROUTINE {
		// SHOW ALL ROUTINE LOAD case: the ALL was already consumed and stored in stmt.Args.
		// parseShowRoutineLoad needs to know about it via a passthrough.
		all := stmt.Args == "ALL"
		node, err := p.parseShowRoutineLoad(startLoc)
		if err != nil {
			return nil, err
		}
		if all {
			if rl, ok := node.(*ast.ShowRoutineLoadStmt); ok {
				rl.All = true
			}
		}
		return node, nil
	}

	switch p.cur.Kind {
	case kwTABLES:
		return p.parseShowTables(stmt)
	case kwDATABASES, kwSCHEMAS:
		return p.parseShowDatabases(stmt)
	case kwCOLUMNS, kwFIELDS:
		return p.parseShowColumns(stmt)
	case kwCREATE:
		// SHOW CREATE ROUTINE LOAD FOR job_name (T6.2)
		if p.peekNext().Kind == kwROUTINE {
			p.advance() // consume CREATE
			return p.parseShowRoutineLoad(startLoc)
		}
		return p.parseShowCreate(stmt)
	case kwVARIABLES:
		return p.parseShowVariables(stmt)
	case kwPARTITIONS:
		return p.parseShowPartitions(stmt)
	case kwGRANTS:
		return p.parseShowGrants(stmt)
	case kwROLES:
		return p.parseShowGenericNamed(stmt, "ROLES")
	case kwCATALOGS:
		return p.parseShowCatalogs(stmt)
	case kwALTER:
		return p.parseShowAlterTable(stmt)
	case kwTABLE:
		// SHOW TABLE STATUS [FROM db] [LIKE 'pat']
		return p.parseShowTableStatus(stmt)
	default:
		// Generic fallthrough: capture the variant type word(s) plus remaining
		// tokens as raw text.
		return p.parseShowGeneric(stmt)
	}
}

// parseShowTables parses: SHOW [FULL] TABLES [FROM db] [LIKE 'pat'] [WHERE ...]
// On entry, cur == kwTABLES.
func (p *Parser) parseShowTables(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume TABLES
	stmt.Type = "TABLES"

	p.parseShowFromLikeWhere(stmt)

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowDatabases parses: SHOW DATABASES [LIKE 'pat'] [WHERE ...]
// On entry, cur == kwDATABASES or kwSCHEMAS.
func (p *Parser) parseShowDatabases(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume DATABASES/SCHEMAS
	stmt.Type = "DATABASES"

	p.parseShowLikeWhere(stmt)

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowColumns parses: SHOW [FULL] COLUMNS FROM table [FROM db] [LIKE 'pat']
// On entry, cur == kwCOLUMNS or kwFIELDS.
func (p *Parser) parseShowColumns(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume COLUMNS/FIELDS
	stmt.Type = "COLUMNS"

	// FROM table
	if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
		p.advance() // consume FROM/IN
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = target
	}

	// Optional second FROM db
	if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
		p.advance() // consume FROM/IN
		dbName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.From = dbName
	}

	p.parseShowLikeWhere(stmt)

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowCreate parses: SHOW [BRIEF] CREATE {TABLE|VIEW|DATABASE|CATALOG|...} name
// On entry, cur == kwCREATE.
func (p *Parser) parseShowCreate(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume CREATE

	// Build type: "CREATE TABLE", "CREATE VIEW", "CREATE DATABASE", etc.
	typeWord := "CREATE"
	switch p.cur.Kind {
	case kwTABLE:
		typeWord = "CREATE TABLE"
		p.advance()
	case kwVIEW:
		typeWord = "CREATE VIEW"
		p.advance()
	case kwDATABASE, kwSCHEMA:
		typeWord = "CREATE DATABASE"
		p.advance()
	case kwCATALOG:
		typeWord = "CREATE CATALOG"
		p.advance()
	case kwMATERIALIZED:
		p.advance() // consume MATERIALIZED
		if p.cur.Kind == kwVIEW {
			p.advance()
		}
		typeWord = "CREATE MATERIALIZED VIEW"
	case kwFUNCTION:
		typeWord = "CREATE FUNCTION"
		p.advance()
	default:
		// Unrecognized: let it fall through, name may be the next token
	}

	// Prepend BRIEF if it was the first modifier
	if stmt.Args == "BRIEF" {
		stmt.Type = "BRIEF " + typeWord
		stmt.Args = ""
	} else {
		stmt.Type = typeWord
	}

	// Optional target name
	if isIdentifierToken(p.cur.Kind) || p.cur.Kind == tokString {
		target, err := p.parseMultipartIdentifier()
		if err == nil {
			stmt.Target = target
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowVariables parses: SHOW [GLOBAL|SESSION] VARIABLES [LIKE 'pat'] [WHERE ...]
// On entry, cur == kwVARIABLES.
func (p *Parser) parseShowVariables(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume VARIABLES
	stmt.Type = "VARIABLES"

	p.parseShowLikeWhere(stmt)

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowPartitions parses:
//
//	SHOW [TEMPORARY] PARTITIONS FROM table [PARTITION(p1,...)]
//	    [WHERE ...] [ORDER BY ...] [LIMIT n]
//
// On entry, cur == kwPARTITIONS. stmt.Args may already contain "TEMPORARY".
func (p *Parser) parseShowPartitions(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume PARTITIONS
	stmt.Type = "PARTITIONS"
	// Note: stmt.Args may be "TEMPORARY" (set by parseShow before dispatch)

	// FROM table
	if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
		p.advance() // consume FROM
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = target
	}

	// Remaining tokens captured as raw args (WHERE, ORDER BY, LIMIT, etc.)
	// Preserve pre-set Args (e.g. "TEMPORARY").
	rest := p.collectRemainingRaw()
	if rest != "" {
		if stmt.Args != "" {
			stmt.Args = stmt.Args + " " + rest
		} else {
			stmt.Args = rest
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowGrants parses: SHOW [ALL] GRANTS [FOR user]
// On entry, cur == kwGRANTS. stmt.Args may contain "ALL".
func (p *Parser) parseShowGrants(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume GRANTS
	stmt.Type = "GRANTS"

	// Optional FOR user_ident; preserve pre-set Args (e.g. "ALL").
	rest := p.collectRemainingRaw()
	if rest != "" {
		if stmt.Args != "" {
			stmt.Args = stmt.Args + " " + rest
		} else {
			stmt.Args = rest
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowTableStatus parses: SHOW TABLE STATUS [FROM db] [LIKE 'pat']
// On entry, cur == kwTABLE.
func (p *Parser) parseShowTableStatus(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume TABLE
	stmt.Type = "TABLE"

	// Optional STATUS
	if p.cur.Kind == kwSTATUS {
		p.advance()
		stmt.Type = "TABLE STATUS"
	}

	p.parseShowFromLikeWhere(stmt)

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowCatalogs parses: SHOW CATALOGS [LIKE 'pat']
// On entry, cur == kwCATALOGS.
func (p *Parser) parseShowCatalogs(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume CATALOGS
	stmt.Type = "CATALOGS"

	p.parseShowLikeWhere(stmt)

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowAlterTable parses: SHOW ALTER TABLE {COLUMN|ROLLUP} [FROM db] [WHERE ...] ...
// On entry, cur == kwALTER.
func (p *Parser) parseShowAlterTable(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume ALTER

	if p.cur.Kind != kwTABLE {
		// Fallback: type is just "ALTER"
		stmt.Type = "ALTER"
		stmt.Args = p.collectRemainingRaw()
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}
	p.advance() // consume TABLE

	subType := ""
	switch p.cur.Kind {
	case kwCOLUMN:
		subType = "COLUMN"
		p.advance()
	case kwROLLUP:
		subType = "ROLLUP"
		p.advance()
	default:
		subType = ""
	}

	if subType == "" {
		stmt.Type = "ALTER TABLE"
	} else {
		stmt.Type = "ALTER TABLE " + subType
	}

	// Optional FROM db
	if p.cur.Kind == kwFROM {
		p.advance()
		dbName, _, err := p.parseIdentifier()
		if err == nil {
			stmt.From = dbName
		}
	}

	// Remaining tokens (WHERE, ORDER BY, LIMIT)
	stmt.Args = p.collectRemainingRaw()
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowGenericNamed handles SHOW variants where the type is already known
// (e.g. ROLES, CATALOGS). It consumes the type keyword and captures the rest
// as Args.
func (p *Parser) parseShowGenericNamed(stmt *ast.ShowStmt, typeName string) (ast.Node, error) {
	p.advance() // consume the type keyword
	stmt.Type = typeName
	stmt.Args = p.collectRemainingRaw()
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowGeneric handles all SHOW variants not specifically recognized above.
// It reads the next 1-2 tokens as the type word(s), then captures the rest
// as raw args.
func (p *Parser) parseShowGeneric(stmt *ast.ShowStmt) (ast.Node, error) {
	// Build the type from the next keyword(s).
	var typeParts []string

	// The BRIEF modifier may have been prepended in Args already.
	prefix := ""
	if stmt.Args == "BRIEF" {
		prefix = "BRIEF "
		stmt.Args = ""
	}

	// First type word.
	if p.cur.Kind != tokEOF {
		typeParts = append(typeParts, strings.ToUpper(p.cur.Str))
		p.advance()
	}

	// Second type word (for two-keyword variants like ALTER TABLE, TABLE STATUS, etc.)
	// Heuristic: consume a second keyword only when both are keyword tokens and
	// the combination commonly appears as a SHOW variant.
	if p.cur.Kind >= 700 && p.cur.Kind != tokEOF {
		second := strings.ToUpper(p.cur.Str)
		// Known two-word type prefixes.
		switch strings.Join(typeParts, "") + second {
		case "ALTERTABLE", "ALTERTABLECOLUMN",
			"TABLETSTATUS", "TABLESTATUS",
			"ALTERCOLUMN", "ALTERROLLUP":
			typeParts = append(typeParts, second)
			p.advance()
		default:
			// Keep only the first word.
		}
	}

	stmt.Type = prefix + strings.Join(typeParts, " ")
	stmt.Args = p.collectRemainingRaw()

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Shared SHOW helpers
// ---------------------------------------------------------------------------

// parseShowFromLikeWhere optionally consumes [FROM db] [LIKE 'pat'] [WHERE expr]
// for SHOW TABLES and similar forms.
func (p *Parser) parseShowFromLikeWhere(stmt *ast.ShowStmt) {
	// Optional FROM db
	if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
		p.advance()
		if isIdentifierToken(p.cur.Kind) {
			dbName, _, err := p.parseIdentifier()
			if err == nil {
				stmt.From = dbName
			}
		}
	}
	p.parseShowLikeWhere(stmt)
}

// parseShowLikeWhere optionally consumes [LIKE 'pat'] [WHERE expr].
func (p *Parser) parseShowLikeWhere(stmt *ast.ShowStmt) {
	if p.cur.Kind == kwLIKE {
		p.advance() // consume LIKE
		if p.cur.Kind == tokString {
			stmt.Like = p.cur.Str
			p.advance()
		}
		return
	}
	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		expr, err := p.parseExpr()
		if err == nil {
			stmt.Where = expr
		}
	}
}

// collectRemainingRaw consumes all remaining tokens (up to EOF) and returns
// them joined by a single space. Used for variant-specific args we don't parse
// in detail.
func (p *Parser) collectRemainingRaw() string {
	var parts []string
	for p.cur.Kind != tokEOF {
		if p.cur.Str != "" {
			parts = append(parts, p.cur.Str)
		} else {
			// Single-char punctuation: use the actual character.
			if p.cur.Kind > 0 && p.cur.Kind < 128 {
				parts = append(parts, string(rune(p.cur.Kind)))
			}
		}
		p.advance()
	}
	return strings.Join(parts, " ")
}

// ---------------------------------------------------------------------------
// DESCRIBE / DESC
// ---------------------------------------------------------------------------

// parseDescribe parses DESCRIBE/DESC statements.
//
//	DESCRIBE [FULL] table_name [ALL VERBOSE]
//	DESC [FULL] table_name
//
// On entry, DESCRIBE or DESC has NOT yet been consumed.
func (p *Parser) parseDescribe() (ast.Node, error) {
	startLoc := p.cur.Loc
	p.advance() // consume DESCRIBE/DESC

	stmt := &ast.DescribeStmt{Loc: startLoc}

	// Optional FULL modifier.
	if p.cur.Kind == kwFULL {
		stmt.Full = true
		p.advance()
	}

	// Target table/function name.
	if isIdentifierToken(p.cur.Kind) || p.cur.Kind == tokString {
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = target
	}

	// Optional ALL VERBOSE suffix.
	if p.cur.Kind == kwALL {
		p.advance() // consume ALL
		if p.cur.Kind == kwVERBOSE {
			p.advance() // consume VERBOSE
			stmt.AllVerbose = true
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// EXPLAIN
// ---------------------------------------------------------------------------

// parseExplain parses EXPLAIN statements.
//
//	EXPLAIN [VERBOSE|GRAPH|PARSED|ANALYZED|REWRITTEN|PLAN|PHYSICAL|
//	         MEMO|SHAPE|DUMP|OPTIMIZED|PLAN PROCESS] query
//
// On entry, EXPLAIN has NOT yet been consumed.
func (p *Parser) parseExplain() (ast.Node, error) {
	startLoc := p.cur.Loc
	p.advance() // consume EXPLAIN

	stmt := &ast.ExplainStmt{Loc: startLoc}

	// Optional explain type modifier.
	switch p.cur.Kind {
	case kwVERBOSE:
		stmt.Type = "VERBOSE"
		p.advance()
	case kwGRAPH:
		stmt.Type = "GRAPH"
		p.advance()
	case kwPARSED:
		stmt.Type = "PARSED"
		p.advance()
	case kwANALYZED:
		stmt.Type = "ANALYZED"
		p.advance()
	case kwREWRITTEN:
		stmt.Type = "REWRITTEN"
		p.advance()
	case kwPLAN:
		p.advance() // consume PLAN
		if p.cur.Kind == kwPROCESS {
			stmt.Type = "PLAN PROCESS"
			p.advance()
		} else {
			stmt.Type = "PLAN"
		}
	case kwPHYSICAL:
		stmt.Type = "PHYSICAL"
		p.advance()
	case kwMEMO:
		stmt.Type = "MEMO"
		p.advance()
	case kwSHAPE:
		stmt.Type = "SHAPE"
		p.advance()
	case kwDUMP:
		stmt.Type = "DUMP"
		p.advance()
	case kwOPTIMIZED:
		stmt.Type = "OPTIMIZED"
		p.advance()
	default:
		// No modifier — empty Type.
	}

	// Parse the explained query (best-effort; use RawQuery as fallback).
	if p.cur.Kind != tokEOF {
		query, err := p.parseStmt()
		if err == nil && query != nil {
			stmt.Query = query
		} else {
			// Fallback: wrap remaining tokens as a RawQuery.
			stmt.Query = &ast.RawQuery{RawText: p.collectRemainingRaw(), Loc: p.cur.Loc}
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// USE
// ---------------------------------------------------------------------------

// parseUse parses a USE statement.
//
//	USE db_name
//	USE catalog_name.db_name
//	USE db_name@cluster_name
//
// On entry, USE has NOT yet been consumed.
func (p *Parser) parseUse() (ast.Node, error) {
	startLoc := p.cur.Loc
	p.advance() // consume USE

	stmt := &ast.UseStmt{Loc: startLoc}

	// Parse the name: may be plain, catalog.db, or db@cluster.
	first, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Check for catalog.db form.
	if p.cur.Kind == int('.') {
		p.advance() // consume '.'
		db, _, err := p.parseIdentifierQualified()
		if err != nil {
			return nil, err
		}
		stmt.Catalog = first
		stmt.Database = db
	} else {
		stmt.Database = first
	}

	// Check for db@cluster form.
	if p.cur.Kind == int('@') {
		p.advance() // consume '@'
		cluster, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Cluster = cluster
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// SET (generic variable assignment)
// ---------------------------------------------------------------------------

// parseGenericSet parses generic SET statements that are not handled by the
// dedicated SET DEFAULT STORAGE VAULT or SET PASSWORD paths:
//
//	SET [GLOBAL|SESSION|LOCAL] var = expr [, ...]
//	SET @@[GLOBAL.|SESSION.]var = expr
//	SET NAMES 'charset' [COLLATE 'collation']
//	SET CHARSET 'charset'
//	SET TRANSACTION { READ ONLY | READ WRITE | ISOLATION LEVEL ... }
//
// On entry, SET has already been consumed; startLoc is the SET token's Loc.
func (p *Parser) parseGenericSet(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.SetStmt{Loc: startLoc, Type: "VARIABLE"}

	// NAMES form
	if p.cur.Kind == kwNAMES {
		return p.parseSetNames(startLoc)
	}
	// CHARSET form
	if p.cur.Kind == kwCHARSET {
		return p.parseSetCharset(startLoc)
	}
	// TRANSACTION form
	if p.cur.Kind == kwTRANSACTION {
		return p.parseSetTransaction(startLoc)
	}

	// One or more variable assignments.
	for {
		item, err := p.parseSetItem()
		if err != nil {
			return nil, err
		}
		stmt.Items = append(stmt.Items, item)
		if _, ok := p.match(int(',')); !ok {
			break
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseSetNames parses: SET NAMES 'charset' [COLLATE 'collation']
// On entry cur == kwNAMES; startLoc is the SET Loc.
func (p *Parser) parseSetNames(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume NAMES
	stmt := &ast.SetStmt{Loc: startLoc, Type: "NAMES"}

	item := &ast.SetItem{Name: "names"}
	if p.cur.Kind == tokString || isIdentifierToken(p.cur.Kind) {
		val, loc, err := p.parseIdentifierOrString()
		if err == nil {
			item.Raw = val
			item.Loc = loc
		}
	}

	// Optional COLLATE
	if p.cur.Kind == kwCOLLATE {
		p.advance()
		collation, _, _ := p.parseIdentifierOrString()
		if collation != "" {
			item.Raw += " COLLATE " + collation
		}
	}

	stmt.Items = []*ast.SetItem{item}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseSetCharset parses: SET CHARSET 'charset'
// On entry cur == kwCHARSET; startLoc is the SET Loc.
func (p *Parser) parseSetCharset(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume CHARSET
	stmt := &ast.SetStmt{Loc: startLoc, Type: "CHARSET"}

	item := &ast.SetItem{Name: "charset"}
	if p.cur.Kind == tokString || isIdentifierToken(p.cur.Kind) {
		val, loc, err := p.parseIdentifierOrString()
		if err == nil {
			item.Raw = val
			item.Loc = loc
		}
	}
	stmt.Items = []*ast.SetItem{item}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseSetTransaction parses: SET TRANSACTION { READ ONLY | READ WRITE | ISOLATION LEVEL ... }
// On entry cur == kwTRANSACTION; startLoc is the SET Loc.
func (p *Parser) parseSetTransaction(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume TRANSACTION
	stmt := &ast.SetStmt{Loc: startLoc, Type: "TRANSACTION"}

	item := &ast.SetItem{Name: "transaction", Raw: p.collectRemainingRaw()}
	stmt.Items = []*ast.SetItem{item}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseSetItem parses a single SET assignment:
//
//	[GLOBAL|SESSION|LOCAL] var_name = expr
//	@@[GLOBAL.|SESSION.]var_name = expr
func (p *Parser) parseSetItem() (*ast.SetItem, error) {
	item := &ast.SetItem{Loc: p.cur.Loc}

	// Optional @@ prefix.
	if p.cur.Kind == tokDoubleAt {
		p.advance() // consume @@
		// Optional scope prefix: GLOBAL. or SESSION.
		if (p.cur.Kind == kwGLOBAL || p.cur.Kind == kwSESSION) && p.peekNext().Kind == int('.') {
			item.Scope = strings.ToUpper(p.cur.Str)
			p.advance() // consume GLOBAL/SESSION
			p.advance() // consume '.'
		}
		name, _, err := p.parseIdentifierQualified()
		if err != nil {
			return nil, err
		}
		item.Name = name
	} else {
		// Optional scope keyword.
		switch p.cur.Kind {
		case kwGLOBAL:
			item.Scope = "GLOBAL"
			p.advance()
		case kwSESSION, kwLOCAL:
			item.Scope = "SESSION"
			p.advance()
		}
		// Variable name.
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		item.Name = name
	}

	// Consume = or :=
	if p.cur.Kind == int('=') || p.cur.Kind == tokAssign {
		p.advance()
	}

	// Parse the value expression.
	if p.cur.Kind != tokEOF && p.cur.Kind != int(',') {
		expr, err := p.parseExpr()
		if err == nil {
			item.Value = expr
		} else {
			// Fallback to raw.
			item.Raw = p.collectUntilCommaOrEOF()
		}
	}

	item.Loc.End = p.prev.Loc.End
	return item, nil
}

// collectUntilCommaOrEOF reads raw tokens until ',' or EOF.
func (p *Parser) collectUntilCommaOrEOF() string {
	var parts []string
	for p.cur.Kind != tokEOF && p.cur.Kind != int(',') {
		if p.cur.Str != "" {
			parts = append(parts, p.cur.Str)
		} else if p.cur.Kind > 0 && p.cur.Kind < 128 {
			parts = append(parts, string(rune(p.cur.Kind)))
		}
		p.advance()
	}
	return strings.Join(parts, " ")
}

// ---------------------------------------------------------------------------
// UNSET
// ---------------------------------------------------------------------------

// parseGenericUnset parses generic UNSET statements that are not handled by
// the dedicated UNSET DEFAULT STORAGE VAULT path:
//
//	UNSET [GLOBAL|SESSION] VARIABLE name [, name ...]
//	UNSET [GLOBAL|SESSION] VARIABLE ALL
//	UNSET [GLOBAL|SESSION] VARIABLE *
//
// On entry, UNSET has already been consumed; startLoc is the UNSET token Loc.
func (p *Parser) parseGenericUnset(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.UnsetStmt{Loc: startLoc, Type: "VARIABLE"}

	// Optional scope.
	switch p.cur.Kind {
	case kwGLOBAL:
		stmt.Scope = "GLOBAL"
		p.advance()
	case kwSESSION, kwLOCAL:
		stmt.Scope = "SESSION"
		p.advance()
	}

	// Expect VARIABLE keyword.
	if p.cur.Kind == kwVARIABLE {
		p.advance()
	}

	// ALL or * — unset all variables.
	if p.cur.Kind == kwALL || p.cur.Kind == int('*') {
		stmt.All = true
		p.advance()
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// One or more variable names.
	for {
		name, _, err := p.parseIdentifier()
		if err != nil {
			break
		}
		stmt.Names = append(stmt.Names, name)
		if _, ok := p.match(int(',')); !ok {
			break
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// HELP
// ---------------------------------------------------------------------------

// parseHelp parses: HELP 'mask'
// On entry, HELP has NOT yet been consumed.
func (p *Parser) parseHelp() (ast.Node, error) {
	startLoc := p.cur.Loc
	p.advance() // consume HELP

	stmt := &ast.HelpStmt{Loc: startLoc}

	if p.cur.Kind == tokString || isIdentifierToken(p.cur.Kind) {
		mask, _, err := p.parseIdentifierOrString()
		if err == nil {
			stmt.Mask = mask
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
