package parser

import (
	"strings"
)

// Outlines, tenant snapshots, concurrent limiting rules and sensitive rules,
// written as cluster_ddl.go writes tenants: forms over the matchers of
// alter_system.go, each parsed into an ObjectDDLStmt that names the statement
// and its object.
//
//	CREATE [OR REPLACE] [FORMAT] OUTLINE name ON explainable_stmt [TO explainable_stmt]
//	CREATE [OR REPLACE] [FORMAT] OUTLINE name ON 'sql' USING HINT /*+ ... */
//	ALTER [FORMAT] OUTLINE name ADD explainable_stmt [TO explainable_stmt]
//	DROP [FORMAT] OUTLINE relation_factor
//	CREATE SNAPSHOT [name] [FOR TENANT name]
//	DROP SNAPSHOT name [FOR TENANT name]
//	CREATE CONCURRENT_LIMITING_RULE [IF NOT EXISTS] name ON {*.* | d.* | d.t}
//	    TO user FOR {ALL | SELECT | UPDATE | INSERT | DELETE}
//	    [FILTER BY KEYWORD ('s' [, 's'] ...)] WITH MAX_CONCURRENCY [=] n [PER SQL]
//	DROP CONCURRENT_LIMITING_RULE [IF EXISTS] relation_factor
//	CREATE SENSITIVE RULE name ON t(c [, c] ...) [, ...] USING ENCRYPTION [= {DEFAULT | 's'}]
//	ALTER SENSITIVE RULE name {ADD COLUMN fields | DROP COLUMN fields | ENABLE | DISABLE | USING ENCRYPTION ...}
//	DROP SENSITIVE RULE name
//
// Ref: create_outline_stmt, alter_outline_stmt, drop_outline_stmt,
// create_tenant_snapshot_stmt, drop_tenant_snapshot_stmt, create_ccl_rule_stmt,
// drop_ccl_rule_stmt, create_sensitive_rule_stmt, alter_sensitive_rule_stmt and
// drop_sensitive_rule_stmt in sql_parser_mysql_mode.y (oceanbase
// v4.3.5_CE_BP6, lines 4817, 15648 and 25582 on). Their rows in
// testdata/compat.tsv hold every alternative, and its near-misses, to the
// pinned server.
var (
	// explainable_stmt, as EXPLAIN reads it. Once a statement has begun, its own
	// error is the statement's: the match keeps it, so a nesting error stays one
	// and a syntax error points where the inner statement failed.
	sysExplainable = sysLeaf(func(p *Parser) bool {
		_, ok, err := p.parseExplainableStmt()
		if err != nil && p.sys.err == nil {
			p.sys.err = err
		}
		return ok && err == nil
	})
	// HINT_HINT_BEGIN hint_list_with_end: the lexer reads `HINT /*+` as one
	// token, and this one drops the hint as a comment, so what is left to check
	// is that one stood after HINT. Being a comment, the hint is past the
	// statement's last token, where its Loc ends.
	sysHintComment = sysLeaf(func(p *Parser) bool {
		gap := strings.TrimLeft(p.inputText(p.prev.End, p.cur.Loc), " \t\r\n")
		return strings.HasPrefix(gap, "/*+")
	})
	sysOutlineTo  = sysOpt("TO", sysExplainable) // opt_outline_target
	createOutline = sysOutline(sysObjName, "ON", sysAlt(
		sysSeq(sysExplainable, sysOutlineTo),
		sysSeq(sysStr, "USING HINT", sysHintComment)))

	// id_dot_id: a name, a dot and a name with nothing between them, which
	// OceanBase's lexer reads as one token. Spaced, `d . t` is no table here.
	sysIDDotID = sysLeaf(func(p *Parser) bool {
		// A user variable is no name on either side: d.@t is no target.
		if !p.isIdentToken() || p.backtickQuotedAt(p.cur.Loc) || p.curIsVariableOrIntroducer() {
			return false
		}
		p.advance()
		if p.cur.Type != '.' || p.cur.Loc != p.prev.End {
			return false
		}
		p.advance()
		if (p.cur.Type != tokIDENT && p.cur.Type < 700) || p.cur.Loc != p.prev.End || p.backtickQuotedAt(p.cur.Loc) ||
			strings.HasPrefix(p.cur.Str, "@") {
			return false
		}
		p.advance()
		return true
	})
	// ccl_database_table_optition.
	sysCCLTarget = sysAlt(
		sysSeq(sysTok('*'), sysTok('.'), sysTok('*')),
		sysSeq(sysName, sysTok('.'), sysTok('*')),
		sysIDDotID,
	)
	// user_with_host_name: user opt_host_name | CURRENT_USER ['(' ')'].
	sysUserWithHost = sysAlt(
		sysSeq(sysAlt(sysStr, sysName), sysOpt(sysUserHost)),
		sysSeq("CURRENT_USER", sysOpt(sysTok('('), sysTok(')'))),
	)

	// complex_string_literal: a string, with a charset introducer or without, or
	// a national string, N'text'.
	sysComplexStr = sysAlt(sysStr, sysLeaf(func(p *Parser) bool {
		if !p.atNationalString() {
			return false
		}
		p.advance() // N
		for p.cur.Type == tokSCONST {
			p.advance()
		}
		return true
	}), sysLeaf(func(p *Parser) bool {
		if !p.curIsVariableOrIntroducer() || !isCharsetIntroducer(p.cur.Str) {
			return false
		}
		p.advance()
		if p.cur.Type != tokSCONST && p.cur.Type != tokXCONST && p.cur.Type != tokBCONST {
			return false
		}
		p.advance()
		return true
	}))
	// sensitive_field_list: relation_factor '(' column_name_list ')', comma-separated.
	sysSensitiveFields = sysList(sysSeq(sysRelFactor, sysTok('('), sysList(sysName, sysComma), sysTok(')')), sysComma)
	// sensitivity_protection_spec.
	sysSensitiveSpec = sysSeq("ENCRYPTION", sysOpt(sysTok('='), sysAlt("DEFAULT", sysComplexStr)))
)

// sysOutline is the forms of an outline statement, outline_type OUTLINE
// clauses, where outline_type is FORMAT or nothing.
func sysOutline(clauses ...any) sysPat {
	return sysAlt(sysForm("OUTLINE", clauses...), sysForm("FORMAT OUTLINE", clauses...))
}

// ruleForms holds, for each word these statements open with, the forms that
// follow it. parseClusterDDL reads them beside clusterForms.
var ruleForms = map[string]sysPat{
	"CREATE": sysAlt(
		createOutline,
		sysForm("SNAPSHOT", sysOpt(sysObjName), sysOpt("FOR TENANT", sysName)),
		sysForm("CONCURRENT_LIMITING_RULE", sysIfNotExists, sysObjName, "ON", sysCCLTarget, "TO", sysUserWithHost,
			"FOR", sysAlt("ALL", "SELECT", "UPDATE", "INSERT", "DELETE"),
			sysOpt("FILTER BY KEYWORD", sysTok('('), sysList(sysStr, sysComma), sysTok(')')),
			"WITH MAX_CONCURRENCY", sysEq, sysInt, sysOpt("PER SQL")),
		sysForm("SENSITIVE RULE", sysObjName, "ON", sysSensitiveFields, "USING", sysSensitiveSpec),
	),
	"CREATE OR REPLACE": createOutline,
	"ALTER": sysAlt(
		sysOutline(sysObjName, "ADD", sysExplainable, sysOutlineTo),
		sysForm("SENSITIVE RULE", sysObjName, sysAlt(
			sysSeq(sysAlt("ADD COLUMN", "DROP COLUMN"), sysSensitiveFields),
			"ENABLE", "DISABLE",
			sysSeq("USING", sysSensitiveSpec))),
	),
	"DROP": sysAlt(
		sysOutline(sysObjFactor),
		sysForm("SNAPSHOT", sysObjName, sysOpt("FOR TENANT", sysName)),
		sysForm("CONCURRENT_LIMITING_RULE", sysIfExists, sysObjFactor),
		sysForm("SENSITIVE RULE", sysObjName),
	),
}

// ruleObjectAhead reports whether the token after CREATE, ALTER or DROP opens
// one of these statements.
func (p *Parser) ruleObjectAhead() bool {
	switch {
	case p.isBareWord("OUTLINE"), p.isBareWord("SNAPSHOT"), p.isBareWord("CONCURRENT_LIMITING_RULE"):
		return true
	case p.isBareWord("FORMAT"):
		return p.nextIsBareWord("OUTLINE")
	case p.isBareWord("SENSITIVE"):
		return p.nextIsBareWord("RULE")
	}
	return false
}
