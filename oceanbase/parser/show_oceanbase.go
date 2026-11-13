package parser

import (
	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseOceanBaseShow parses the SHOW statements OceanBase has and MySQL does
// not, after SHOW, and reports whether one stood there:
//
//	SHOW SERVER STATUS [cond]
//	SHOW SCHEMA
//	SHOW TRACE [FORMAT = 'fmt'] [cond]
//	SHOW ALL PARAMETERS [cond]
//	SHOW PARAMETERS [cond] [TENANT [=] {name | 'name' | ALL}]
//	SHOW TABLEGROUPS [cond]
//	SHOW TENANT [STATUS]
//	SHOW QUERY_RESPONSE_TIME | RECYCLEBIN | RESTORE PREVIEW | CATALOGS
//	SHOW SEQUENCES [cond] [{FROM | IN} db]
//	SHOW JOB STATUS [WHERE JOB = 'job']
//	SHOW SENSITIVE RULES [{FROM | IN} {t [{FROM | IN} db] | DATABASE db}] [cond]
//
// Ref: show_stmt and show_sensitive_rule_option in sql_parser_mysql_mode.y
// (oceanbase v4.3.5_CE_BP6, line 16105 on). Their rows in testdata/compat.tsv
// hold each form to the pinned server.
func (p *Parser) parseOceanBaseShow(stmt *nodes.ShowStmt) (bool, error) {
	// words reads the bare words given, in order, and reports whether they all
	// stood there; it reads nothing otherwise.
	words := func(ws ...string) bool {
		if !p.isBareWord(ws[0]) || (len(ws) > 1 && !p.nextIsBareWord(ws[1])) {
			return false
		}
		for range ws {
			p.advance()
		}
		return true
	}
	var err error
	switch {
	case words("SERVER", "STATUS"):
		stmt.Type = "SERVER STATUS"
		err = p.parseShowLikeOrWhere(stmt)
	case words("SCHEMA"):
		stmt.Type = "SCHEMA"
	case words("TRACE"):
		stmt.Type = "TRACE"
		if p.isBareWord("FORMAT") {
			p.advance()
			if _, err = p.expect('='); err != nil {
				return true, err
			}
			if p.cur.Type != tokSCONST {
				return true, p.syntaxErrorAtCur()
			}
			p.advance()
		}
		err = p.parseShowLikeOrWhere(stmt)
	case words("ALL", "PARAMETERS"):
		stmt.Type = "ALL PARAMETERS"
		err = p.parseShowLikeOrWhere(stmt)
	case words("PARAMETERS"):
		stmt.Type = "PARAMETERS"
		if err = p.parseShowLikeOrWhere(stmt); err != nil {
			return true, err
		}
		// opt_tenant_name: TENANT [=] relation_name_or_string.
		if p.isBareWord("TENANT") {
			p.advance()
			p.match('=')
			switch {
			case p.cur.Type == tokSCONST:
				stmt.From = &nodes.TableRef{Loc: nodes.Loc{Start: p.cur.Loc, End: p.cur.End}, Name: p.cur.Str}
				p.advance()
			case p.cur.Type == kwALL:
				p.advance()
			default:
				ref := &nodes.TableRef{Loc: nodes.Loc{Start: p.cur.Loc}}
				if ref.Name, _, err = p.parseIdent(); err != nil {
					return true, err
				}
				ref.Loc.End = p.prev.End
				stmt.From = ref
			}
		}
	case words("TABLEGROUPS"):
		stmt.Type = "TABLEGROUPS"
		err = p.parseShowLikeOrWhere(stmt)
	case words("TENANT"):
		stmt.Type = "TENANT"
		if p.cur.Type == kwSTATUS {
			p.advance()
			stmt.Type = "TENANT STATUS"
		}
	case words("QUERY_RESPONSE_TIME"):
		stmt.Type = "QUERY_RESPONSE_TIME"
	case words("RECYCLEBIN"):
		stmt.Type = "RECYCLEBIN"
	case words("RESTORE", "PREVIEW"):
		stmt.Type = "RESTORE PREVIEW"
	case words("CATALOGS"):
		stmt.Type = "CATALOGS"
	case words("SEQUENCES"):
		// The condition comes before the database here, as the grammar writes it.
		stmt.Type = "SEQUENCES"
		if err = p.parseShowLikeOrWhere(stmt); err != nil {
			return true, err
		}
		if p.matchFromOrIn() {
			stmt.From, err = p.parseTableRef()
		}
	case words("JOB", "STATUS"):
		stmt.Type = "JOB STATUS"
		if p.cur.Type == kwWHERE {
			p.advance()
			if !p.isBareWord("JOB") {
				return true, p.syntaxErrorAtCur()
			}
			p.advance()
			if _, err = p.expect('='); err != nil {
				return true, err
			}
			if p.cur.Type != tokSCONST {
				return true, p.syntaxErrorAtCur()
			}
			p.advance()
		}
	case words("SENSITIVE", "RULES"):
		stmt.Type = "SENSITIVE RULES"
		if p.matchFromOrIn() {
			if p.cur.Type == kwDATABASE {
				p.advance()
				stmt.From, err = p.parseTableRef()
			} else if stmt.From, err = p.parseRelationFactor(); err == nil && p.matchFromOrIn() {
				var db string
				if db, _, err = p.parseIdent(); err == nil {
					stmt.From.Schema = db
				}
			}
			if err != nil {
				return true, err
			}
		}
		err = p.parseShowLikeOrWhere(stmt)
	default:
		return false, nil
	}
	return true, err
}
