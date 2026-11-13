package parser

import (
	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// AUDIT and NOAUDIT, written as alter_system.go writes ALTER SYSTEM: forms
// over the same matchers.
//
//	{AUDIT | NOAUDIT} {shortcut [, shortcut] ... | ALL | ALL STATEMENTS}
//	    [ON {table | DEFAULT} | BY user [, user] ...]
//	    [BY ACCESS] [WHENEVER [NOT] SUCCESSFUL]
//
// Ref: audit_stmt, audit_clause and audit_all_shortcut in
// sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6). The rule's rows in
// testdata/compat.tsv hold every alternative, and its near-misses, to the
// pinned server.
var (
	// audit_all_shortcut.
	sysAuditShortcut = sysAlt(
		"ALTER SYSTEM", "ALTER TABLE", "ALTER", "AUDIT", "CLUSTER", "COMMENT TABLE", "COMMENT", "CONTEXT",
		"DELETE TABLE", "DELETE", "EXECUTE PROCEDURE", "EXECUTE", "FLASHBACK", "GRANT PROCEDURE", "GRANT TABLE",
		"GRANT TYPE", "GRANT", "INDEX", "INSERT TABLE", "INSERT", "MATERIALIZED VIEW", "NOT EXISTS", "OUTLINE",
		"PROCEDURE", "PROFILE", "RENAME", "SELECT TABLE", "SELECT", "SESSION", "SYSTEM AUDIT", "SYSTEM GRANT",
		"TABLESPACE", "TABLE", "TRIGGER", "TYPE", "UPDATE TABLE", "UPDATE", "USER", "VIEW",
	)
	// normal_relation_factor: a table's name with no leading dot.
	sysNormalFactor = sysLeaf(func(p *Parser) bool {
		if p.cur.Type == '.' {
			return false
		}
		_, err := p.parseRelationFactor()
		return err == nil
	})
	// audit_user: STRING_VALUE | NAME_OB | unreserved_keyword_normal, which
	// leaves out the words a role may not be called and four more.
	sysAuditUser = sysAlt(sysStr, sysLeaf(func(p *Parser) bool {
		if isAmbiguousRoleWord(p.cur.Type) {
			return false
		}
		for _, w := range []string{"ACCESS", "PASSWORD", "MAXVALUE", "RANGE"} {
			if p.isBareWord(w) {
				return false
			}
		}
		_, _, err := p.parseIdent()
		return err == nil
	}))
	// opt_host_name: the @host of a user, as parseUserHostOpt reads it.
	sysUserHost = sysLeaf(func(p *Parser) bool {
		_, hasHost, err := p.parseUserHostOpt()
		return err == nil && hasHost
	})
	// audit_clause.
	auditForms = sysSeq(
		sysAlt(sysList(sysAuditShortcut, sysComma), "ALL STATEMENTS", "ALL"),
		sysOpt(sysAlt(
			sysSeq("ON", sysAlt("DEFAULT", sysNormalFactor)),
			sysSeq("BY", sysList(sysSeq(sysAuditUser, sysOpt(sysUserHost)), sysComma)),
		)),
		sysOpt("BY ACCESS"),
		sysOpt(sysAlt("WHENEVER NOT SUCCESSFUL", "WHENEVER SUCCESSFUL")),
	)
)

// parseAuditStmt parses AUDIT or NOAUDIT at that word.
func (p *Parser) parseAuditStmt() (*nodes.AuditStmt, error) {
	stmt := &nodes.AuditStmt{Loc: nodes.Loc{Start: p.pos()}, NoAudit: p.isBareWord("NOAUDIT")}
	p.advance()
	if _, err := p.matchSystemForms(auditForms); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}
