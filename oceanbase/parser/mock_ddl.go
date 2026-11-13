package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// MySQL's tablespace, server, logfile group and instance statements, which
// OceanBase parses and does not run: mock_stmt and its neighbours, with fewer
// options than MySQL's. Each is a table of forms over alter_system.go's
// matchers, built per statement so its leaves can fill the node.
//
// Ref: create_tablespace_stmt, alter_tablespace_stmt, drop_tablespace_stmt,
// create_server_stmt, alter_server_stmt, drop_server_stmt,
// create_logfile_group_stmt, alter_logfile_group_stmt, drop_logfile_group_stmt
// and rotate_master_key_stmt in sql_parser_mysql_mode.y (oceanbase
// v4.3.5_CE_BP6, line 23065 on). Their rows in testdata/compat.tsv hold every
// alternative, and its near-misses, to the pinned server.

// mockStr is STRING_VALUE stored in dst. Strings side by side are one value,
// as OceanBase's lexer reads them.
func mockStr(dst *string) sysPat {
	return sysLeaf(func(p *Parser) bool {
		if p.cur.Type != tokSCONST {
			return false
		}
		var b strings.Builder
		for p.cur.Type == tokSCONST {
			b.WriteString(p.advance().Str)
		}
		*dst = b.String()
		return true
	})
}

// mockInt is INTNUM stored in dst as written.
func mockInt(dst *string) sysPat {
	return sysLeaf(func(p *Parser) bool {
		if p.cur.Type != tokICONST {
			return false
		}
		*dst = p.advance().Str
		return true
	})
}

// mockNameOB is NAME_OB, an identifier that is no keyword, stored in dst.
func mockNameOB(dst *string) sysPat {
	return sysLeaf(func(p *Parser) bool {
		name, _, err := p.parseNameOB()
		*dst = name
		return err == nil
	})
}

// mockSet sets *dst where it matches.
func mockSet(dst *bool, parts ...any) sysPat {
	match := sysSeq(parts...)
	return func(p *Parser, k func() bool) bool {
		return match(p, func() bool {
			*dst = true
			return k()
		})
	}
}

// mockSpaceOrCommaList is a list whose items stand side by side or with a
// comma between: logfile_group_options and its kin.
func mockSpaceOrCommaList(item sysPat) sysPat {
	return sysList(item, sysOptComma)
}

// parseCreateTablespaceStmt parses CREATE TABLESPACE at TABLESPACE:
//
//	CREATE TABLESPACE name
//	    [ENCRYPTION [=] 'e' [, ENCRYPTION [=] 'e'] ...
//	    | option [option] ...]
//
//	option: ADD DATAFILE 'f' | FILE_BLOCK_SIZE = n | USE LOGFILE GROUP 'g'
//	      | {EXTENT_SIZE | INITIAL_SIZE | AUTOEXTEND_SIZE | MAX_SIZE | NODEGROUP} [=] n
//	      | WAIT | COMMENT [=] 'c' | ENGINE [=] {name | 'name' | ALL}
func (p *Parser) parseCreateTablespaceStmt(start int) (*nodes.CreateTablespaceStmt, error) {
	p.advance() // consume TABLESPACE
	stmt := &nodes.CreateTablespaceStmt{Loc: nodes.Loc{Start: start}}
	option := sysAlt(
		sysSeq("ADD DATAFILE", mockStr(&stmt.DataFile)),
		sysSeq("FILE_BLOCK_SIZE", sysTok('='), mockInt(&stmt.FileBlockSize)),
		sysSeq("USE LOGFILE GROUP", mockStr(&stmt.UseLogfileGroup)),
		sysSeq("EXTENT_SIZE", sysEq, mockInt(&stmt.ExtentSize)),
		sysSeq("INITIAL_SIZE", sysEq, mockInt(&stmt.InitialSize)),
		sysSeq("AUTOEXTEND_SIZE", sysEq, mockInt(&stmt.AutoextendSize)),
		sysSeq("MAX_SIZE", sysEq, mockInt(&stmt.MaxSize)),
		sysSeq("NODEGROUP", sysEq, mockInt(&stmt.NodeGroup)),
		mockSet(&stmt.Wait, "WAIT"),
		sysSeq("COMMENT", sysEq, mockStr(&stmt.Comment)),
		sysSeq("ENGINE", sysEq, p.mockEngineName(&stmt.Engine)),
	)
	forms := sysSeq(mockNameOB(&stmt.Name), sysOpt(sysAlt(
		sysList(sysSeq("ENCRYPTION", sysEq, mockStr(&stmt.Encryption)), sysComma), // permanent_tablespace_options
		sysList(option, sysSeq()), // opt_tablespace_options
	)))
	if _, err := p.matchSystemForms(forms); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterTablespaceStmt parses ALTER TABLESPACE at TABLESPACE:
//
//	ALTER TABLESPACE name {[SET] ENCRYPTION [=] 'e' [, [SET] ENCRYPTION [=] 'e'] | option [option] ...}
//
//	option: {ADD | DROP} DATAFILE 'f' | INITIAL_SIZE [=] n | WAIT | ENGINE [=] {name | 'name' | ALL}
//
// alter_tablespace_actions holds one action or two, and no more.
func (p *Parser) parseAlterTablespaceStmt(start int) (*nodes.AlterTablespaceStmt, error) {
	p.advance() // consume TABLESPACE
	stmt := &nodes.AlterTablespaceStmt{Loc: nodes.Loc{Start: start}}
	action := sysSeq(sysOpt("SET"), "ENCRYPTION", sysEq, mockStr(&stmt.Encryption))
	option := sysAlt(
		sysSeq("ADD DATAFILE", mockStr(&stmt.AddDataFile)),
		sysSeq("DROP DATAFILE", mockStr(&stmt.DropDataFile)),
		sysSeq("INITIAL_SIZE", sysEq, mockInt(&stmt.InitialSize)),
		mockSet(&stmt.Wait, "WAIT"),
		sysSeq("ENGINE", sysEq, p.mockEngineName(&stmt.Engine)),
	)
	forms := sysSeq(mockNameOB(&stmt.Name), sysAlt(
		sysSeq(action, sysOpt(sysComma, action)),
		sysList(option, sysSeq()),
	))
	if _, err := p.matchSystemForms(forms); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropTablespaceStmt parses DROP TABLESPACE at TABLESPACE:
//
//	DROP TABLESPACE name [ENGINE [=] {name | 'name' | ALL}]
func (p *Parser) parseDropTablespaceStmt(start int) (*nodes.DropTablespaceStmt, error) {
	p.advance() // consume TABLESPACE
	stmt := &nodes.DropTablespaceStmt{Loc: nodes.Loc{Start: start}}
	forms := sysSeq(mockNameOB(&stmt.Name), sysOpt("ENGINE", sysEq, p.mockEngineName(&stmt.Engine)))
	if _, err := p.matchSystemForms(forms); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// mockEngineName is relation_name_or_string, the engine of
// opt_tablespace_engine: a name, a string or ALL.
func (p *Parser) mockEngineName(dst *string) sysPat {
	return sysAlt(
		sysLeaf(func(p *Parser) bool {
			name, _, err := p.parseIdent()
			*dst = name
			return err == nil
		}),
		mockStr(dst),
		sysLeaf(func(p *Parser) bool {
			if p.cur.Type != kwALL {
				return false
			}
			*dst = p.advance().Str
			return true
		}),
	)
}

// mockServerOptions is OPTIONS '(' server_options_list ')': each option is
// one of seven words and its value, kept as the word and the value.
func mockServerOptions(dst *[]string) sysPat {
	option := func(word string, value func(*string) sysPat) sysPat {
		var v string
		match := sysSeq(word, value(&v))
		return func(p *Parser, k func() bool) bool {
			return match(p, func() bool {
				*dst = append(*dst, word+" "+v)
				if k() {
					return true
				}
				*dst = (*dst)[:len(*dst)-1]
				return false
			})
		}
	}
	item := sysAlt(
		option("USER", mockStr), option("HOST", mockStr), option("DATABASE", mockStr),
		option("OWNER", mockStr), option("PASSWORD", mockStr), option("SOCKET", mockStr),
		option("PORT", mockInt),
	)
	return sysSeq("OPTIONS", sysTok('('), sysList(item, sysComma), sysTok(')'))
}

// parseCreateServerStmt parses CREATE SERVER at SERVER:
//
//	CREATE SERVER name FOREIGN DATA WRAPPER name OPTIONS (option [, option] ...)
//
//	option: {USER | HOST | DATABASE | OWNER | PASSWORD | SOCKET} 'value' | PORT n
func (p *Parser) parseCreateServerStmt(start int) (*nodes.CreateServerStmt, error) {
	p.advance() // consume SERVER
	stmt := &nodes.CreateServerStmt{Loc: nodes.Loc{Start: start}}
	forms := sysSeq(mockNameOB(&stmt.Name), "FOREIGN DATA WRAPPER", mockNameOB(&stmt.WrapperName), mockServerOptions(&stmt.Options))
	if _, err := p.matchSystemForms(forms); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterServerStmt parses ALTER SERVER at SERVER:
//
//	ALTER SERVER name OPTIONS (option [, option] ...)
func (p *Parser) parseAlterServerStmt(start int) (*nodes.AlterServerStmt, error) {
	p.advance() // consume SERVER
	stmt := &nodes.AlterServerStmt{Loc: nodes.Loc{Start: start}}
	if _, err := p.matchSystemForms(sysSeq(mockNameOB(&stmt.Name), mockServerOptions(&stmt.Options))); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropServerStmt parses DROP SERVER at SERVER:
//
//	DROP SERVER [IF EXISTS] name
func (p *Parser) parseDropServerStmt(start int) (*nodes.DropServerStmt, error) {
	p.advance() // consume SERVER
	stmt := &nodes.DropServerStmt{Loc: nodes.Loc{Start: start}}
	if _, err := p.matchSystemForms(sysSeq(sysOpt(mockSet(&stmt.IfExists, "IF EXISTS")), mockNameOB(&stmt.Name))); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// mockLogfileEngine is opt_ts_engine: [STORAGE] ENGINE [=] {NAME_OB | 'name'}.
func mockLogfileEngine(dst *string) sysPat {
	return sysSeq(sysOpt("STORAGE"), "ENGINE", sysEq, sysAlt(mockNameOB(dst), mockStr(dst)))
}

// mockAddLogFile is add_log_file: ADD {UNDOFILE | REDOFILE} 'file'. The node
// keeps an undo file only.
func mockAddLogFile(undo *string) sysPat {
	var redo string
	return sysSeq("ADD", sysAlt(sysSeq("UNDOFILE", mockStr(undo)), sysSeq("REDOFILE", mockStr(&redo))))
}

// parseCreateLogfileGroupStmt parses CREATE LOGFILE GROUP at LOGFILE:
//
//	CREATE LOGFILE GROUP name ADD {UNDOFILE | REDOFILE} 'file' [option [[,] option] ...]
//
//	option: {INITIAL_SIZE | UNDO_BUFFER_SIZE | REDO_BUFFER_SIZE | NODEGROUP} [=] n
//	      | [STORAGE] ENGINE [=] {name | 'name'} | WAIT | NO_WAIT | COMMENT [=] 'c'
func (p *Parser) parseCreateLogfileGroupStmt(start int) (*nodes.CreateLogfileGroupStmt, error) {
	p.advance() // consume LOGFILE
	stmt := &nodes.CreateLogfileGroupStmt{Loc: nodes.Loc{Start: start}}
	option := sysAlt(
		sysSeq("INITIAL_SIZE", sysEq, mockInt(&stmt.InitialSize)),
		sysSeq("UNDO_BUFFER_SIZE", sysEq, mockInt(&stmt.UndoBufferSize)),
		sysSeq("REDO_BUFFER_SIZE", sysEq, mockInt(&stmt.RedoBufferSize)),
		sysSeq("NODEGROUP", sysEq, mockInt(&stmt.NodeGroup)),
		mockLogfileEngine(&stmt.Engine),
		mockSet(&stmt.Wait, "WAIT"),
		"NO_WAIT",
		sysSeq("COMMENT", sysEq, mockStr(&stmt.Comment)),
	)
	forms := sysSeq("GROUP", mockNameOB(&stmt.Name), mockAddLogFile(&stmt.UndoFile), sysOpt(mockSpaceOrCommaList(option)))
	if _, err := p.matchSystemForms(forms); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterLogfileGroupStmt parses ALTER LOGFILE GROUP at LOGFILE:
//
//	ALTER LOGFILE GROUP name ADD {UNDOFILE | REDOFILE} 'file' [option [[,] option] ...]
//
//	option: INITIAL_SIZE [=] n | [STORAGE] ENGINE [=] {name | 'name'} | WAIT | NO_WAIT
func (p *Parser) parseAlterLogfileGroupStmt(start int) (*nodes.AlterLogfileGroupStmt, error) {
	p.advance() // consume LOGFILE
	stmt := &nodes.AlterLogfileGroupStmt{Loc: nodes.Loc{Start: start}}
	option := sysAlt(
		sysSeq("INITIAL_SIZE", sysEq, mockInt(&stmt.InitialSize)),
		mockLogfileEngine(&stmt.Engine),
		mockSet(&stmt.Wait, "WAIT"),
		"NO_WAIT",
	)
	forms := sysSeq("GROUP", mockNameOB(&stmt.Name), mockAddLogFile(&stmt.UndoFile), sysOpt(mockSpaceOrCommaList(option)))
	if _, err := p.matchSystemForms(forms); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropLogfileGroupStmt parses DROP LOGFILE GROUP at LOGFILE:
//
//	DROP LOGFILE GROUP name [[,] option [[,] option] ...]
//
//	option: [STORAGE] ENGINE [=] {name | 'name'} | WAIT | NO_WAIT
//
// drop_ts_options_list may be empty before its first comma, so `DROP LOGFILE
// GROUP g, WAIT` is taken.
func (p *Parser) parseDropLogfileGroupStmt(start int) (*nodes.DropLogfileGroupStmt, error) {
	p.advance() // consume LOGFILE
	stmt := &nodes.DropLogfileGroupStmt{Loc: nodes.Loc{Start: start}}
	option := sysAlt(mockLogfileEngine(&stmt.Engine), "WAIT", "NO_WAIT")
	forms := sysSeq("GROUP", mockNameOB(&stmt.Name), sysOpt(sysOptComma, mockSpaceOrCommaList(option)))
	if _, err := p.matchSystemForms(forms); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterInstanceStmt parses ALTER INSTANCE at INSTANCE. OceanBase has one
// form, rotate_master_key_stmt:
//
//	ALTER INSTANCE ROTATE INNODB MASTER KEY
func (p *Parser) parseAlterInstanceStmt(start int) (*nodes.AlterInstanceStmt, error) {
	p.advance() // consume INSTANCE
	stmt := &nodes.AlterInstanceStmt{Loc: nodes.Loc{Start: start}}
	if _, err := p.matchSystemForms(sysSeq("ROTATE INNODB MASTER KEY")); err != nil {
		return nil, err
	}
	stmt.Action = "ROTATE INNODB MASTER KEY"
	stmt.Loc.End = p.prev.End
	return stmt, nil
}
