package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// Doris hands every statement of this file but KILL to the legacy CUP parser
// (DorisParser.g4 `unsupportedStatement`), so a statement is accepted only
// when both grammars take it: each reader follows the narrower of
// DorisParser.g4 and DorisLegacyParser.g4.

// ---------------------------------------------------------------------------
// BACKUP SNAPSHOT
// ---------------------------------------------------------------------------

// parseBackup parses (DorisParser.g4:247 #backup, DorisLegacyParser.g4:1826
// backup_stmt):
//
//	BACKUP SNAPSHOT [db.]label TO repo_name
//	    [ON (baseTableRef [, ...])]
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

	// Optional ON (table_list)
	if p.cur.Kind == kwON {
		p.advance() // consume ON
		tables, loc, err := p.parseBackupTableList()
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

// parseBackupTableList parses the parenthesized table list of a BACKUP or
// RESTORE ON clause: '(' baseTableRef (',' baseTableRef)* ')'
// (DorisParser.g4:248, DorisLegacyParser.g4:1863 opt_backup_table_ref_list).
// The list holds at least one table. Each entry is
//
//	name [[TEMPORARY] PARTITION[S] ...] [[AS] alias]
//
// the part of baseTableRef (DorisParser.g4:516) that the legacy base_table_ref
// (DorisLegacyParser.g4:2256) shares with it. Partition names are consumed but
// not stored.
func (p *Parser) parseBackupTableList() ([]*ast.ObjectName, ast.Loc, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, ast.NoLoc(), err
	}

	var tables []*ast.ObjectName
	for {
		tbl, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, ast.NoLoc(), err
		}
		tables = append(tables, tbl)
		if err := p.skipBackupTableRefTail(); err != nil {
			return nil, ast.NoLoc(), err
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
	return tables, closeTok.Loc, nil
}

// skipBackupTableRefTail consumes the optional partition restriction and
// alias after a table name in a BACKUP or RESTORE ON list:
//
//	specifiedPartition (DorisParser.g4:1663):
//	    [TEMPORARY] PARTITION (identifier | identifierList)
//	    [TEMPORARY] PARTITIONS identifierList
//	tableAlias (DorisParser.g4:1393, DorisLegacyParser.g4:2286): [AS] identifier
func (p *Parser) skipBackupTableRefTail() error {
	if p.cur.Kind == kwTEMPORARY {
		p.advance() // consume TEMPORARY
		if p.cur.Kind != kwPARTITION && p.cur.Kind != kwPARTITIONS {
			return p.syntaxErrorAtCur()
		}
	}
	switch p.cur.Kind {
	case kwPARTITION:
		p.advance() // consume PARTITION
		if p.cur.Kind == int('(') {
			if err := p.skipBackupIdentifierList(); err != nil {
				return err
			}
		} else if _, _, err := p.parseIdentifier(); err != nil {
			return err
		}
	case kwPARTITIONS:
		p.advance() // consume PARTITIONS
		if err := p.skipBackupIdentifierList(); err != nil {
			return err
		}
	}
	if p.cur.Kind == kwAS {
		p.advance() // consume AS
		_, _, err := p.parseIdentifier()
		return err
	}
	if isIdentifierToken(p.cur.Kind) {
		_, _, err := p.parseIdentifier()
		return err
	}
	return nil
}

// skipBackupIdentifierList consumes '(' identifier (',' identifier)* ')'.
func (p *Parser) skipBackupIdentifierList() error {
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	for {
		if _, _, err := p.parseIdentifier(); err != nil {
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

// parseRestore parses (DorisParser.g4:250 #restore, DorisLegacyParser.g4:1875
// restore_stmt):
//
//	RESTORE SNAPSHOT [db.]label FROM repo_name
//	    [ON (baseTableRef [, ...])]
//	    [PROPERTIES("key"="value", ...)]
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

	// Optional ON (table_list)
	if p.cur.Kind == kwON {
		p.advance()
		tables, loc, err := p.parseBackupTableList()
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

// parseKill parses (DorisParser.g4:232 supportedKillStatement,
// DorisLegacyParser.g4:1880 kill_stmt):
//
//	KILL [CONNECTION] integer
//	KILL QUERY {integer | 'query_id'}
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

// parseLockTables parses (DorisParser.g4:241 #lockTables,
// DorisLegacyParser.g4:1858 lock_tables_stmt):
//
//	LOCK TABLES [lockTable [, ...]]
//
// On entry, cur is TABLES (LOCK was consumed by dispatch).
func (p *Parser) parseLockTables(startLoc ast.Loc) (ast.Node, error) {
	tablesTok, err := p.expect(kwTABLES)
	if err != nil {
		return nil, err
	}

	stmt := &ast.LockTablesStmt{}
	endLoc := tablesTok.Loc

	// The table list is optional: a bare LOCK TABLES is accepted.
	if p.cur.Kind != tokEOF {
		for {
			item, err := p.parseLockItem()
			if err != nil {
				return nil, err
			}
			stmt.Items = append(stmt.Items, item)
			endLoc = ast.NodeLoc(item)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseLockItem parses one lock entry (DorisParser.g4:260 lockTable,
// DorisLegacyParser.g4:1843 lock_table):
//
//	tbl [AS {alias | 'alias'}] {READ [LOCAL] | [LOW_PRIORITY] WRITE}
//
// The lock mode is required.
func (p *Parser) parseLockItem() (*ast.LockItem, error) {
	tbl, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	item := &ast.LockItem{Table: tbl}

	// Optional AS alias: identifierOrText (DorisLegacyParser.g4:1836 lock_alias).
	if p.cur.Kind == kwAS {
		p.advance()
		alias, _, err := p.parseIdentifierOrString()
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

	item.Loc = ast.NodeLoc(tbl).Merge(p.prev.Loc)
	return item, nil
}

// parseUnlockTables parses UNLOCK TABLES (DorisParser.g4:242 #unlockTables,
// DorisLegacyParser.g4:1831 unlock_tables_stmt).
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
// WARM UP
// ---------------------------------------------------------------------------

// parseWarmUp parses (DorisLegacyParser.g4:206 warm_up_stmt, narrower than
// DorisParser.g4:243 #warmUpCluster: the source names the same kind of
// object as the destination, and only it takes properties):
//
//	WARM UP CLUSTER dest WITH CLUSTER src [FORCE] [PROPERTIES(...)]
//	WARM UP COMPUTE GROUP dest WITH COMPUTE GROUP src [FORCE] [PROPERTIES(...)]
//	WARM UP {CLUSTER | COMPUTE GROUP} dest WITH warmUpItem [AND ...] [FORCE]
//
// warmUpItem is TABLE name [PARTITION partition] (DorisLegacyParser.g4:214).
// Verb holds CLUSTER or COMPUTE GROUP; Target holds the text after it.
//
// On entry, cur is UP (WARM was consumed by dispatch).
func (p *Parser) parseWarmUp(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwUP); err != nil {
		return nil, err
	}

	stmt := &ast.WarmUpStmt{}
	sameKind := func() (string, error) {
		switch p.cur.Kind {
		case kwCLUSTER:
			p.advance()
			return "CLUSTER", nil
		case kwCOMPUTE:
			p.advance()
			if _, err := p.expect(kwGROUP); err != nil {
				return "", err
			}
			return "COMPUTE GROUP", nil
		}
		return "", p.syntaxErrorAtCur()
	}

	verb, err := sameKind()
	if err != nil {
		return nil, err
	}
	stmt.Verb = verb
	targetStart := p.cur.Loc.Start

	if _, _, err := p.parseIdentifier(); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}

	if p.cur.Kind == kwTABLE {
		for {
			if _, err := p.expect(kwTABLE); err != nil {
				return nil, err
			}
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
			if p.cur.Kind == kwPARTITION {
				p.advance() // consume PARTITION
				if _, _, err := p.parseIdentifier(); err != nil {
					return nil, err
				}
			}
			if p.cur.Kind != kwAND {
				break
			}
			p.advance() // consume AND
		}
		if p.cur.Kind == kwFORCE {
			p.advance()
		}
	} else {
		srcKind := p.cur
		kind, err := sameKind()
		if err != nil {
			return nil, err
		}
		if kind != verb {
			return nil, syntaxErrorAt(srcKind)
		}
		if _, _, err := p.parseIdentifier(); err != nil {
			return nil, err
		}
		if p.cur.Kind == kwFORCE {
			p.advance()
		}
		if p.cur.Kind == kwPROPERTIES {
			if _, err := p.parseProperties(); err != nil {
				return nil, err
			}
		}
	}

	stmt.Target = p.utilAdminRawText(targetStart)
	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CLEAN
// ---------------------------------------------------------------------------

// parseClean parses (DorisLegacyParser.g4:156 clean_stmt, the same forms as
// DorisParser.g4:445 unsupportedCleanStatement):
//
//	CLEAN LABEL [label] {FROM | IN} db
//	CLEAN ALL PROFILE
//	CLEAN QUERY STATS FOR db
//	CLEAN QUERY STATS {FROM | IN} [db.]table
//	CLEAN ALL QUERY STATS
//
// Verb holds the keywords after CLEAN; Target holds the text after them.
//
// On entry, cur is the keyword after CLEAN (CLEAN was consumed by dispatch).
func (p *Parser) parseClean(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.CleanStmt{}

	switch p.cur.Kind {
	case kwALL:
		p.advance() // consume ALL
		switch p.cur.Kind {
		case kwPROFILE:
			p.advance()
			stmt.Verb = "ALL PROFILE"
		case kwQUERY:
			p.advance()
			if _, err := p.expect(kwSTATS); err != nil {
				return nil, err
			}
			stmt.Verb = "ALL QUERY STATS"
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case kwLABEL:
		p.advance() // consume LABEL
		stmt.Verb = "LABEL"
		targetStart := p.cur.Loc.Start
		if p.cur.Kind != kwFROM && p.cur.Kind != kwIN {
			if _, _, err := p.parseIdentifier(); err != nil {
				return nil, err
			}
		}
		if _, err := p.parseUtilAdminFromOrIn(); err != nil {
			return nil, err
		}
		if _, _, err := p.parseIdentifier(); err != nil {
			return nil, err
		}
		stmt.Target = p.utilAdminRawText(targetStart)
	case kwQUERY:
		p.advance() // consume QUERY
		if _, err := p.expect(kwSTATS); err != nil {
			return nil, err
		}
		stmt.Verb = "QUERY STATS"
		targetStart := p.cur.Loc.Start
		if p.cur.Kind == kwFOR {
			p.advance()
			if _, _, err := p.parseIdentifier(); err != nil {
				return nil, err
			}
		} else {
			if _, err := p.parseUtilAdminFromOrIn(); err != nil {
				return nil, err
			}
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
		}
		stmt.Target = p.utilAdminRawText(targetStart)
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CANCEL (generic)
// ---------------------------------------------------------------------------

// parseCancelGeneric parses the CANCEL forms of DorisParser.g4:461
// unsupportedCancelStatement and DorisLegacyParser.g4:1655 cancel_param,
// except DECOMMISSION (parseCancelDecommissionBackend):
//
//	CANCEL LOAD [{FROM | IN} db] [LIKE 'pattern' | WHERE expr]
//	CANCEL EXPORT [{FROM | IN} db] [LIKE 'pattern' | WHERE expr]
//	CANCEL ALTER TABLE {ROLLUP | MATERIALIZED VIEW | COLUMN} FROM table [(id, ...)]
//	CANCEL BUILD INDEX ON table [(id, ...)]
//	CANCEL BACKUP [{FROM | IN} db]
//	CANCEL RESTORE [{FROM | IN} db]
//	CANCEL WARM UP JOB [LIKE 'pattern' | WHERE expr]
//
// Verb identifies the operation; Target holds the text after it.
//
// On entry, cur is the keyword immediately after CANCEL.
func (p *Parser) parseCancelGeneric(cancelLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.CancelStmt{}
	var targetStart int

	switch p.cur.Kind {
	case kwLOAD, kwEXPORT:
		stmt.Verb = strings.ToUpper(p.cur.Str)
		p.advance()
		targetStart = p.cur.Loc.Start
		if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
			p.advance()
			if _, _, err := p.parseIdentifier(); err != nil {
				return nil, err
			}
		}
		if err := p.skipUtilAdminWildWhere(); err != nil {
			return nil, err
		}
	case kwBACKUP, kwRESTORE:
		stmt.Verb = strings.ToUpper(p.cur.Str)
		p.advance()
		targetStart = p.cur.Loc.Start
		if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
			p.advance()
			if _, _, err := p.parseIdentifier(); err != nil {
				return nil, err
			}
		}
	case kwALTER:
		p.advance() // consume ALTER
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		stmt.Verb = "ALTER TABLE"
		targetStart = p.cur.Loc.Start
		switch p.cur.Kind {
		case kwROLLUP, kwCOLUMN:
			p.advance()
		case kwMATERIALIZED:
			p.advance()
			if _, err := p.expect(kwVIEW); err != nil {
				return nil, err
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		if _, err := p.parseMultipartIdentifier(); err != nil {
			return nil, err
		}
		if err := p.skipUtilAdminJobIDList(); err != nil {
			return nil, err
		}
	case kwBUILD:
		p.advance() // consume BUILD
		if _, err := p.expect(kwINDEX); err != nil {
			return nil, err
		}
		stmt.Verb = "BUILD INDEX"
		targetStart = p.cur.Loc.Start
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		if _, err := p.parseMultipartIdentifier(); err != nil {
			return nil, err
		}
		if err := p.skipUtilAdminJobIDList(); err != nil {
			return nil, err
		}
	case kwWARM:
		p.advance() // consume WARM
		if _, err := p.expect(kwUP); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwJOB); err != nil {
			return nil, err
		}
		stmt.Verb = "WARM UP JOB"
		targetStart = p.cur.Loc.Start
		if err := p.skipUtilAdminWildWhere(); err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Target = p.utilAdminRawText(targetStart)
	stmt.Loc = cancelLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// parseCancelDecommissionBackend parses (DorisParser.g4:470
// #cancelDecommisionBackend, DorisLegacyParser.g4:1660):
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

// parseRecover parses (DorisParser.g4:477 unsupportedRecoverStatement,
// DorisLegacyParser.g4:991 recover_stmt):
//
//	RECOVER DATABASE name [id] [AS new_name]
//	RECOVER TABLE [db.]name [id] [AS new_name]
//	RECOVER PARTITION name [id] [AS new_name] FROM [db.]table
//
// On entry, cur is DATABASE | TABLE | PARTITION (RECOVER was consumed by dispatch).
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

	// Optional numeric id.
	if p.cur.Kind == tokInt {
		stmt.ID = p.cur.Str
		p.advance()
	}

	// Optional AS new_name.
	if p.cur.Kind == kwAS {
		p.advance()
		newName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.NewName = newName
	}

	// PARTITION: FROM [db.]table is required.
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
// UNSET VARIABLE
// ---------------------------------------------------------------------------

// parseUnsetVariableStmt parses (DorisParser.g4:923 unsupoortedUnsetStatement,
// DorisLegacyParser.g4:1900 unset_var_stmt):
//
//	UNSET [GLOBAL | SESSION | LOCAL] VARIABLE {ALL | name}
//
// UNSET DEFAULT STORAGE VAULT is read by parseUnsetDefaultStorageVault.
// On entry, UNSET has been consumed; startLoc is the UNSET token Loc.
func (p *Parser) parseUnsetVariableStmt(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.UnsetStmt{Type: "VARIABLE"}

	switch p.cur.Kind {
	case kwGLOBAL:
		stmt.Scope = "GLOBAL"
		p.advance()
	case kwSESSION, kwLOCAL:
		stmt.Scope = "SESSION"
		p.advance()
	}

	if _, err := p.expect(kwVARIABLE); err != nil {
		return nil, err
	}

	if p.cur.Kind == kwALL {
		stmt.All = true
		p.advance()
	} else {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Names = []string{name}
	}

	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// parseUtilAdminFromOrIn consumes FROM or IN (DorisLegacyParser.g4:1576
// from_or_in).
func (p *Parser) parseUtilAdminFromOrIn() (Token, error) {
	if tok, ok := p.match(kwFROM, kwIN); ok {
		return tok, nil
	}
	return Token{}, p.syntaxErrorAtCur()
}

// skipUtilAdminWildWhere consumes an optional wildWhere (DorisParser.g4:521,
// DorisLegacyParser.g4:1589 opt_wild_where): LIKE 'pattern' | WHERE expr.
func (p *Parser) skipUtilAdminWildWhere() error {
	switch p.cur.Kind {
	case kwLIKE:
		p.advance()
		if _, err := p.expect(tokString); err != nil {
			return err
		}
	case kwWHERE:
		p.advance()
		if _, err := p.parseExpr(); err != nil {
			return err
		}
	}
	return nil
}

// skipUtilAdminJobIDList consumes an optional job id list
// (DorisLegacyParser.g4:1649 job_id_list): '(' integer (',' integer)* ')'.
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
