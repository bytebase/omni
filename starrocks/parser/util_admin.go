package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// The readers in this file follow StarRocks.g4 and, where it refuses a tree
// the grammar admits, StarRocks's AST builder (AstBuilder.java in StarRocks
// 3.4, the pinned version), whose ParsingException is a parse error too.

// ---------------------------------------------------------------------------
// BACKUP SNAPSHOT
// ---------------------------------------------------------------------------

// parseBackup parses (StarRocks.g4:1765 backupStatement):
//
//	BACKUP SNAPSHOT [db.]label TO repo_name
//	    [ON (backupRestoreObjectDesc [, ...])]
//	    [PROPERTIES("key"="value", ...)]
//
// On entry, cur is SNAPSHOT (BACKUP was consumed by the dispatch).
func (p *Parser) parseBackup(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwSNAPSHOT); err != nil {
		return nil, err
	}

	stmt := &ast.BackupStmt{}

	// [db.]label — parse as multipart identifier
	labelName, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Label = labelName.String()
	stmt.LabelParts = labelName.Parts
	endLoc := ast.NodeLoc(labelName)

	// TO repo_name
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	repo, repoLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Repo = repo
	endLoc = repoLoc

	// Optional ON (object_list)
	if p.cur.Kind == kwON {
		p.advance() // consume ON
		tables, loc, err := p.parseBackupObjectList(false)
		if err != nil {
			return nil, err
		}
		stmt.Tables = tables
		endLoc = loc
	}

	// Optional PROPERTIES(...)
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseBackupObjectList parses the ON list of BACKUP and RESTORE:
// '(' backupRestoreObjectDesc (',' backupRestoreObjectDesc)* ')'
// (StarRocks.g4:1768, 1785), where each entry is (StarRocks.g4:2564,
// backupRestoreObjectDesc and backupRestoreTableDesc):
//
//	name [partitionNames] [AS alias]
//	{TABLE | TABLES} name [partitionNames] [AS alias]
//	{FUNCTION | FUNCTIONS | VIEW | VIEWS | MATERIALIZED {VIEW | VIEWS}} name [AS alias]
//	ALL {FUNCTION | FUNCTIONS | TABLE | TABLES | VIEW | VIEWS | MATERIALIZED {VIEW | VIEWS}}
//
// A BACKUP entry takes no alias: the AST builder refuses one
// (unsupportedSepcifyAliasInBackupStmt). Every named object is returned,
// whatever its kind; ALL entries return nothing.
func (p *Parser) parseBackupObjectList(restore bool) ([]*ast.ObjectName, ast.Loc, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, ast.NoLoc(), err
	}

	var names []*ast.ObjectName
	for {
		name, err := p.parseBackupObjectDesc(restore)
		if err != nil {
			return nil, ast.NoLoc(), err
		}
		if name != nil {
			names = append(names, name)
		}
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, ast.NoLoc(), err
	}
	return names, closeTok.Loc, nil
}

// parseBackupObjectDesc parses one entry of a BACKUP or RESTORE ON list (see
// parseBackupObjectList). A kind keyword that can also be a table name (VIEW,
// TABLES, ...) is the kind only when an object name follows it, as in
// ANTLR's prediction: `ON (VIEW)` backs up a table named view.
func (p *Parser) parseBackupObjectDesc(restore bool) (*ast.ObjectName, error) {
	isKind := func(kind int) bool {
		switch kind {
		case kwFUNCTION, kwFUNCTIONS, kwTABLE, kwTABLES, kwVIEW, kwVIEWS, kwMATERIALIZED:
			return true
		}
		return false
	}

	if p.cur.Kind == kwALL {
		p.advance() // consume ALL
		if !isKind(p.cur.Kind) {
			return nil, p.syntaxErrorAtCur()
		}
		mv := p.cur.Kind == kwMATERIALIZED
		p.advance()
		if mv {
			if _, ok := p.match(kwVIEW, kwVIEWS); !ok {
				return nil, p.syntaxErrorAtCur()
			}
		}
		return nil, nil
	}

	tableDesc := true
	if isKind(p.cur.Kind) {
		next := p.peekNext().Kind
		switch {
		case p.cur.Kind == kwMATERIALIZED:
			if next == kwVIEW || next == kwVIEWS {
				p.advance() // consume MATERIALIZED
				p.advance() // consume VIEW | VIEWS
				tableDesc = false
			}
		case isIdentifierToken(next) || p.cur.Kind == kwFUNCTION:
			// FUNCTION is not in StarRocks.g4's nonReserved, so it is never
			// a table name here.
			tableDesc = p.cur.Kind == kwTABLE || p.cur.Kind == kwTABLES
			p.advance() // consume the kind keyword
		}
	}

	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	if tableDesc {
		if err := p.skipBackupPartitionNames(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwAS {
		if !restore {
			return nil, &ParseError{
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
				Message: "BACKUP does not take an alias for the objects it backs up",
			}
		}
		p.advance() // consume AS
		if _, _, err := p.parseIdentifier(); err != nil {
			return nil, err
		}
	}
	return name, nil
}

// skipBackupPartitionNames consumes an optional partitionNames
// (StarRocks.g4:2268):
//
//	[TEMPORARY] {PARTITION | PARTITIONS} {name_or_string | (name_or_string, ...)}
//	PARTITION (column = literal, ...)                       -- keyPartitions
func (p *Parser) skipBackupPartitionNames() error {
	temporary := false
	if p.cur.Kind == kwTEMPORARY {
		p.advance() // consume TEMPORARY
		temporary = true
		if p.cur.Kind != kwPARTITION && p.cur.Kind != kwPARTITIONS {
			return p.syntaxErrorAtCur()
		}
	}
	if p.cur.Kind != kwPARTITION && p.cur.Kind != kwPARTITIONS {
		return nil
	}
	single := p.cur.Kind == kwPARTITION
	p.advance() // consume PARTITION | PARTITIONS
	if p.cur.Kind != int('(') {
		_, _, err := p.parseIdentifierOrString()
		return err
	}
	p.advance() // consume '('
	keyed := single && !temporary && isIdentifierToken(p.cur.Kind) && p.peekNext().Kind == int('=')
	for {
		if keyed {
			if _, _, err := p.parseIdentifier(); err != nil {
				return err
			}
			if _, err := p.expect(int('=')); err != nil {
				return err
			}
			if _, err := p.parseExpr(); err != nil {
				return err
			}
		} else if _, _, err := p.parseIdentifierOrString(); err != nil {
			return err
		}
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	_, err := p.expect(int(')'))
	return err
}

// ---------------------------------------------------------------------------
// RESTORE SNAPSHOT
// ---------------------------------------------------------------------------

// parseRestore parses (StarRocks.g4:1780 restoreStatement):
//
//	RESTORE SNAPSHOT [db.]label FROM repo_name
//	    [ON (backupRestoreObjectDesc [, ...])]
//	    [PROPERTIES("key"="value", ...)]
//
// An ON list needs the database in the label: the AST builder refuses one
// otherwise (unsupportedOnClauseWithoutAnyDbNameInRestoreStmt).
//
// On entry, cur is SNAPSHOT (RESTORE was consumed by the dispatch).
func (p *Parser) parseRestore(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwSNAPSHOT); err != nil {
		return nil, err
	}

	stmt := &ast.RestoreStmt{}

	labelName, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Label = labelName.String()
	stmt.LabelParts = labelName.Parts
	endLoc := ast.NodeLoc(labelName)

	// FROM repo_name
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	repo, repoLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Repo = repo
	endLoc = repoLoc

	// Optional ON (object_list)
	if p.cur.Kind == kwON {
		if len(labelName.Parts) < 2 {
			return nil, &ParseError{
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
				Message: "RESTORE with ON needs the database in the snapshot label",
			}
		}
		p.advance()
		tables, loc, err := p.parseBackupObjectList(true)
		if err != nil {
			return nil, err
		}
		stmt.Tables = tables
		endLoc = loc
	}

	// Optional PROPERTIES(...)
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// KILL
// ---------------------------------------------------------------------------

// parseKill parses (StarRocks.g4:728 killStatement):
//
//	KILL [CONNECTION] integer
//	KILL QUERY {integer | 'query_id'}
//
// The grammar admits a string after KILL [CONNECTION], but the AST builder
// refuses it ("connection id ... should be a positive integer").
//
// On entry, cur is CONNECTION | QUERY | the id token (KILL was consumed by dispatch).
func (p *Parser) parseKill(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.KillStmt{}

	switch p.cur.Kind {
	case kwCONNECTION:
		stmt.Kind = "CONNECTION"
		p.advance()
	case kwQUERY:
		stmt.Kind = "QUERY"
		p.advance()
	}

	// Target: an integer id, or a string query id after QUERY only.
	if p.cur.Kind != tokInt && (p.cur.Kind != tokString || stmt.Kind != "QUERY") {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Target = p.cur.Str
	stmt.Loc = startLoc.Merge(p.cur.Loc)
	p.advance()
	return stmt, nil
}

// ---------------------------------------------------------------------------
// LOCK TABLES
// ---------------------------------------------------------------------------

// parseLockTables parses (StarRocks.g4:2010):
//
//	LOCK TABLES lock_item [, ...]
//
// On entry, cur is TABLES (LOCK was consumed by dispatch).
func (p *Parser) parseLockTables(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwTABLES); err != nil {
		return nil, err
	}

	stmt := &ast.LockTablesStmt{}
	for {
		item, err := p.parseLockItem()
		if err != nil {
			return nil, err
		}
		stmt.Items = append(stmt.Items, item)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}

	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// parseLockItem parses one lock entry (StarRocks.g4:2014 lock_item,
// 2018 lock_type):
//
//	tbl [[AS] alias] {READ [LOCAL] | [LOW_PRIORITY] WRITE}
//
// The table is a single identifier, and the lock mode is required.
func (p *Parser) parseLockItem() (*ast.LockItem, error) {
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	tbl := &ast.ObjectName{Parts: []string{name}, Loc: nameLoc}
	item := &ast.LockItem{Table: tbl}

	// Optional [AS] alias. A lock-mode keyword is the alias only when a lock
	// mode still follows it.
	lockStart := func(kind int) bool {
		return kind == kwREAD || kind == kwWRITE || kind == kwLOW_PRIORITY
	}
	if p.cur.Kind == kwAS {
		p.advance()
		alias, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		item.Alias = alias
	} else if isIdentifierToken(p.cur.Kind) && (!lockStart(p.cur.Kind) || lockStart(p.peekNext().Kind)) {
		alias, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		item.Alias = alias
	}

	// Lock mode
	switch p.cur.Kind {
	case kwREAD:
		p.advance()
		item.Mode = "READ"
		if p.cur.Kind == kwLOCAL {
			p.advance()
			item.Mode = "READ LOCAL"
		}
	case kwLOW_PRIORITY:
		p.advance()
		if _, err := p.expect(kwWRITE); err != nil {
			return nil, err
		}
		item.Mode = "LOW_PRIORITY WRITE"
	case kwWRITE:
		p.advance()
		item.Mode = "WRITE"
	default:
		return nil, p.syntaxErrorAtCur()
	}

	item.Loc = nameLoc.Merge(p.prev.Loc)
	return item, nil
}

// parseUnlockTables parses UNLOCK TABLES (StarRocks.g4:2011).
//
// On entry, cur is TABLES (UNLOCK was consumed by dispatch).
func (p *Parser) parseUnlockTables(startLoc ast.Loc) (ast.Node, error) {
	tablesTok, err := p.expect(kwTABLES)
	if err != nil {
		return nil, err
	}
	return &ast.UnlockTablesStmt{Loc: startLoc.Merge(tablesTok.Loc)}, nil
}

// ---------------------------------------------------------------------------
// INSTALL / UNINSTALL PLUGIN
// ---------------------------------------------------------------------------

// parseInstallPlugin parses:
//
//	INSTALL PLUGIN FROM 'source'
//	INSTALL PLUGIN FROM SONAME 'library'
//	(with optional PROPERTIES(...))
//
// On entry, cur is PLUGIN (INSTALL was consumed by dispatch).
func (p *Parser) parseInstallPlugin(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwPLUGIN); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	stmt := &ast.InstallPluginStmt{}

	// Optional SONAME keyword
	if p.cur.Kind == kwSONAME {
		stmt.IsSoname = true
		p.advance()
	}

	// Source path / URL / library name — string literal
	src, srcLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Source = src
	endLoc := srcLoc

	// Optional PROPERTIES(...)
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseUninstallPlugin parses:
//
//	UNINSTALL PLUGIN name
//
// On entry, cur is PLUGIN (UNINSTALL was consumed by dispatch).
func (p *Parser) parseUninstallPlugin(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwPLUGIN); err != nil {
		return nil, err
	}

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &ast.UninstallPluginStmt{
		Name: name,
		Loc:  startLoc.Merge(nameLoc),
	}, nil
}

// ---------------------------------------------------------------------------
// CLEAN
// ---------------------------------------------------------------------------

// parseClean parses the one CLEAN statement StarRocks has
// (StarRocks.g4:505 cleanTemporaryTableStatement):
//
//	CLEAN TEMPORARY TABLE ON SESSION 'session_id'
//
// The AST builder reads the session id with Java's UUID.fromString and
// refuses one it cannot read ("invalid session id format"). Verb is
// "TEMPORARY TABLE"; Target is the session id.
//
// On entry, cur is the keyword after CLEAN (CLEAN was consumed by dispatch).
func (p *Parser) parseClean(startLoc ast.Loc) (ast.Node, error) {
	for _, kind := range []int{kwTEMPORARY, kwTABLE, kwON, kwSESSION} {
		if _, err := p.expect(kind); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	if !isJavaUUIDString(p.cur.Str) {
		return nil, &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "invalid session id format",
		}
	}
	stmt := &ast.CleanStmt{Verb: "TEMPORARY TABLE", Target: p.cur.Str, Loc: startLoc.Merge(p.cur.Loc)}
	p.advance()
	return stmt, nil
}

// isJavaUUIDString reports whether java.util.UUID.fromString accepts s: at
// most 36 characters in five dash-separated groups of hexadecimal digits,
// each group fitting a long.
func isJavaUUIDString(s string) bool {
	if len(s) > 36 {
		return false
	}
	groups := strings.Split(s, "-")
	if len(groups) != 5 {
		return false
	}
	for _, g := range groups {
		if g == "" || len(g) > 16 {
			return false
		}
		for _, c := range g {
			if !('0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F') {
				return false
			}
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// CANCEL (generic)
// ---------------------------------------------------------------------------

// parseCancelGeneric parses the CANCEL statements of StarRocks.g4 other than
// DECOMMISSION (parseCancelDecommissionBackend):
//
//	CANCEL LOAD [FROM db] [WHERE expr]                                (1495)
//	CANCEL EXPORT [{FROM | IN} [catalog.]db] [LIKE 'pattern' | WHERE expr]   (1877)
//	CANCEL ALTER TABLE [COLUMN | ROLLUP | OPTIMIZE] FROM table [(id, ...)]  (585)
//	CANCEL ALTER MATERIALIZED VIEW FROM name                          (587)
//	CANCEL BACKUP {{FROM | IN} db | FOR EXTERNAL CATALOG}             (1772)
//	CANCEL RESTORE {{FROM | IN} db | FOR EXTERNAL CATALOG}            (1789)
//	CANCEL REFRESH MATERIALIZED VIEW name [FORCE]                     (688)
//	CANCEL REFRESH DICTIONARY name                                    (860)
//	CANCEL COMPACTION WHERE expr                                      (1506)
//
// The grammar makes the database of CANCEL BACKUP and RESTORE optional, but
// the AST builder refuses a statement that names neither a database nor
// the external catalogs (nullIdentifierCancelBackupRestore).
// Verb identifies the operation; Target holds the text after it.
//
// On entry, cur is the keyword immediately after CANCEL.
func (p *Parser) parseCancelGeneric(cancelLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.CancelStmt{}
	var targetStart int

	switch p.cur.Kind {
	case kwLOAD:
		stmt.Verb = "LOAD"
		p.advance()
		targetStart = p.cur.Loc.Start
		if p.cur.Kind == kwFROM {
			p.advance()
			if _, _, err := p.parseIdentifier(); err != nil {
				return nil, err
			}
		}
		if p.cur.Kind == kwWHERE {
			p.advance()
			if _, err := p.parseExpr(); err != nil {
				return nil, err
			}
		}
	case kwEXPORT:
		stmt.Verb = "EXPORT"
		p.advance()
		targetStart = p.cur.Loc.Start
		if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
			p.advance()
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
		}
		switch p.cur.Kind {
		case kwLIKE:
			p.advance()
			if _, err := p.expect(tokString); err != nil {
				return nil, err
			}
		case kwWHERE:
			p.advance()
			if _, err := p.parseExpr(); err != nil {
				return nil, err
			}
		}
	case kwBACKUP, kwRESTORE:
		stmt.Verb = strings.ToUpper(p.cur.Str)
		p.advance()
		targetStart = p.cur.Loc.Start
		switch p.cur.Kind {
		case kwFROM, kwIN:
			p.advance()
			if _, _, err := p.parseIdentifier(); err != nil {
				return nil, err
			}
		case kwFOR:
			p.advance()
			if _, err := p.expect(kwEXTERNAL); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwCATALOG); err != nil {
				return nil, err
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case kwALTER:
		p.advance() // consume ALTER
		switch p.cur.Kind {
		case kwTABLE:
			p.advance()
			stmt.Verb = "ALTER TABLE"
			targetStart = p.cur.Loc.Start
			p.match(kwCOLUMN, kwROLLUP, kwOPTIMIZE) // the job kind is optional
			if _, err := p.expect(kwFROM); err != nil {
				return nil, err
			}
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
			if err := p.skipUtilAdminJobIDList(); err != nil {
				return nil, err
			}
		case kwMATERIALIZED:
			p.advance()
			if _, err := p.expect(kwVIEW); err != nil {
				return nil, err
			}
			stmt.Verb = "ALTER MATERIALIZED VIEW"
			targetStart = p.cur.Loc.Start
			if _, err := p.expect(kwFROM); err != nil {
				return nil, err
			}
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case kwREFRESH:
		p.advance() // consume REFRESH
		switch p.cur.Kind {
		case kwMATERIALIZED:
			p.advance()
			if _, err := p.expect(kwVIEW); err != nil {
				return nil, err
			}
			stmt.Verb = "REFRESH MATERIALIZED VIEW"
			targetStart = p.cur.Loc.Start
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
			if p.cur.Kind == kwFORCE {
				p.advance()
			}
		case kwDICTIONARY:
			p.advance()
			stmt.Verb = "REFRESH DICTIONARY"
			targetStart = p.cur.Loc.Start
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case tokIdent:
		// COMPACTION is a plain identifier to this lexer.
		if !strings.EqualFold(p.cur.Str, "compaction") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume COMPACTION
		stmt.Verb = "COMPACTION"
		targetStart = p.cur.Loc.Start
		if _, err := p.expect(kwWHERE); err != nil {
			return nil, err
		}
		if _, err := p.parseExpr(); err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Target = p.utilAdminRawText(targetStart)
	stmt.Loc = cancelLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// parseCancelDecommissionBackend parses (StarRocks.g4:750
// cancelAlterSystemStatement):
//
//	CANCEL DECOMMISSION BACKEND 'host:port' [, 'host:port' ...]
//
// CANCEL has already been consumed; cur is DECOMMISSION on entry.
func (p *Parser) parseCancelDecommissionBackend(cancelLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwDECOMMISSION); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwBACKEND); err != nil {
		return nil, err
	}
	stmt := &ast.CancelDecommissionStmt{Object: "BACKEND"}
	for {
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Hosts = append(stmt.Hosts, p.cur.Str)
		p.advance()
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	stmt.Loc = cancelLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// RECOVER
// ---------------------------------------------------------------------------

// parseRecover parses:
//
//	RECOVER {DATABASE | SCHEMA} name              (StarRocks.g4:369 recoverDbStmt)
//	RECOVER TABLE [db.]name                       (StarRocks.g4:577 recoverTableStatement)
//	RECOVER PARTITION name FROM [db.]table        (StarRocks.g4:596 recoverPartitionStatement)
//
// On entry, cur is DATABASE | SCHEMA | TABLE | PARTITION (RECOVER was consumed by dispatch).
func (p *Parser) parseRecover(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.RecoverStmt{}

	switch p.cur.Kind {
	case kwDATABASE, kwSCHEMA:
		stmt.Verb = "DATABASE"
	case kwTABLE:
		stmt.Verb = "TABLE"
	case kwPARTITION:
		stmt.Verb = "PARTITION"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// The table name may be qualified; a database or partition name may not.
	if stmt.Verb == "TABLE" {
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	} else {
		name, loc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Name = &ast.ObjectName{Parts: []string{name}, Loc: loc}
	}

	if stmt.Verb == "PARTITION" {
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		fromTable, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.FromTable = fromTable
	}

	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// skipUtilAdminJobIDList consumes an optional job id list:
// '(' INTEGER_VALUE (',' INTEGER_VALUE)* ')' (StarRocks.g4:586).
func (p *Parser) skipUtilAdminJobIDList() error {
	if p.cur.Kind != int('(') {
		return nil
	}
	p.advance() // consume '('
	for {
		if _, err := p.expect(tokInt); err != nil {
			return err
		}
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	_, err := p.expect(int(')'))
	return err
}

// utilAdminRawText returns the source text from the absolute offset start to
// the end of the last consumed token, or "" when nothing was consumed.
func (p *Parser) utilAdminRawText(start int) string {
	end := p.prev.Loc.End
	if start < p.baseOffset || end <= start || end-p.baseOffset > len(p.input) {
		return ""
	}
	return p.input[start-p.baseOffset : end-p.baseOffset]
}
