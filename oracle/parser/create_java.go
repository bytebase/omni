package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateJavaStmt parses a CREATE JAVA statement. It is called after
// JAVA; mods holds the AND { RESOLVE | COMPILE } and NOFORCE options the CREATE
// dispatcher consumed before JAVA.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-JAVA.html
//
//	CREATE [ OR REPLACE ] [ AND { RESOLVE | COMPILE } ] [ NOFORCE ]
//	    JAVA { SOURCE [ IF NOT EXISTS ] [ NAMED [ schema. ] primary_name ]
//	         | RESOURCE [ IF NOT EXISTS ] NAMED [ schema. ] primary_name
//	         | CLASS [ IF NOT EXISTS ] [ SCHEMA schema ] }
//	    [ SHARING = { METADATA | NONE } ]
//	    [ AUTHID { CURRENT_USER | DEFINER } ]
//	    [ resolver_clause ]
//	    { USING { BFILE ( directory_object_name , server_file_name )
//	            | { CLOB | BLOB | BFILE } subquery
//	            | 'key_for_BLOB' }
//	    | AS source_char }
//
// Oracle 23ai takes IF NOT EXISTS only where shown: after the kind, not after
// JAVA (ORA-00905), before it (ORA-00901, refused by the CREATE dispatcher),
// or after SCHEMA (ORA-00922).
//
// Oracle 23ai checks this head before it looks for a Java VM (ORA-29538). It
// rejects a missing kind (ORA-00905); NAMED missing after RESOURCE (ORA-29502)
// or given after CLASS (ORA-29500); SHARING, AUTHID, and RESOLVER out of the
// order above, or any other word (ORA-00922); and a statement with neither
// USING nor AS (ORA-00922). SOURCE takes its name from the source when NAMED
// is left out. BFILE followed by '(' takes the directory and file name
// (ORA-22929 for a subquery there).
func (p *Parser) parseCreateJavaStmt(start int, orReplace bool, mods []*nodes.DDLOption) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_JAVA,
		OrReplace:  orReplace,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	for _, mod := range mods {
		opts.Items = append(opts.Items, mod)
	}
	// addOption records an option that started at optStart and ends at the
	// last consumed token.
	addOption := func(key, value string, optStart int, items *nodes.List) {
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key:   key,
			Value: value,
			Items: items,
			Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}

	// SOURCE [ NAMED name ] | RESOURCE NAMED name | CLASS [ SCHEMA schema ]
	optStart := p.pos()
	switch {
	case p.isKeywordStr("SOURCE"), p.isKeywordStr("RESOURCE"):
		value := p.cur.Str
		p.advance()
		addOption("JAVA_TYPE", value, optStart, nil)
		if err := p.parseJavaIfNotExists(stmt); err != nil {
			return nil, err
		}
		if p.isKeywordStr("NAMED") {
			p.advance() // consume NAMED
			// Oracle 23ai takes any word here, USING included: NAMED USING
			// fails after USING, not at it.
			var err error
			if stmt.Name, err = p.parseObjectName(); err != nil {
				return nil, err
			}
			if stmt.Name.Name == "" {
				return nil, p.syntaxErrorAtCur()
			}
		} else if value == "RESOURCE" {
			return nil, p.syntaxErrorAtCur()
		}
	case p.isKeywordStr("CLASS"):
		p.advance()
		addOption("JAVA_TYPE", "CLASS", optStart, nil)
		if err := p.parseJavaIfNotExists(stmt); err != nil {
			return nil, err
		}
		if p.isKeywordStr("SCHEMA") {
			schemaStart := p.pos()
			p.advance() // consume SCHEMA
			if !p.isCreateTableNameStart() {
				return nil, p.syntaxErrorAtCur()
			}
			schema := p.cur.Str
			p.advance()
			addOption("SCHEMA", schema, schemaStart, nil)
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	// [ SHARING = { METADATA | NONE } ]
	if p.isKeywordStr("SHARING") {
		optStart := p.pos()
		p.advance() // consume SHARING
		if err := p.expectCreateTableToken('='); err != nil {
			return nil, err
		}
		if !p.isKeywordStr("METADATA") && !p.isKeywordStr("NONE") {
			return nil, p.syntaxErrorAtCur()
		}
		value := p.cur.Str
		p.advance()
		addOption("SHARING", value, optStart, nil)
	}

	// [ AUTHID { CURRENT_USER | DEFINER } ]
	if p.isKeywordStr("AUTHID") {
		if err := p.parseJavaAuthid(addOption); err != nil {
			return nil, err
		}
	}

	// [ resolver_clause ]
	if p.isKeywordStr("RESOLVER") {
		if err := p.parseJavaResolverClause(opts); err != nil {
			return nil, err
		}
	}

	switch {
	case p.cur.Type == kwUSING:
		if err := p.parseJavaUsingClause(addOption); err != nil {
			return nil, err
		}
	case p.cur.Type == kwAS:
		// source_char is Java source, not SQL, and runs to the end of the
		// parsed range; see javaSourceHead. Take it verbatim. Oracle reads
		// the missing source as an acceptance (compat.tsv's CREATE JAVA
		// SOURCE AS row).
		if stmt.Name == nil && (p.peekNext().Type == tokEOF || p.peekNext().Type == ';') {
			p.advance() // consume AS
		} else if err := p.parseSourceText(opts); err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Options = opts
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseJavaIfNotExists parses the optional IF NOT EXISTS that follows SOURCE,
// RESOURCE, or CLASS in CREATE JAVA. Oracle 23ai rejects any other spelling,
// IF EXISTS or a quoted "EXISTS" included (ORA-11543), and OR REPLACE beside
// it (ORA-11541).
func (p *Parser) parseJavaIfNotExists(stmt *nodes.AdminDDLStmt) error {
	if p.cur.Type != kwIF {
		return nil
	}
	p.advance() // consume IF
	if p.cur.Type != kwNOT {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume NOT
	if p.cur.Type != kwEXISTS || stmt.OrReplace {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume EXISTS
	stmt.IfNotExists = true
	return nil
}

// parseJavaAuthid parses the invoker_rights_clause of CREATE and ALTER JAVA,
// from its AUTHID:
//
//	AUTHID { CURRENT_USER | DEFINER }
func (p *Parser) parseJavaAuthid(addOption func(key, value string, optStart int, items *nodes.List)) error {
	optStart := p.pos()
	p.advance() // consume AUTHID
	if !p.isKeywordStr("CURRENT_USER") && !p.isKeywordStr("DEFINER") {
		return p.syntaxErrorAtCur()
	}
	value := p.cur.Str
	p.advance()
	addOption("AUTHID", value, optStart, nil)
	return nil
}

// parseJavaResolverClause parses the resolver_clause of CREATE and ALTER
// JAVA, from its RESOLVER, and records it as a RESOLVER option whose items
// hold one option per pair, keyed by the match string:
//
//	RESOLVER ( [ ( match_string { schema_name | - } ) ]... )
//
// Oracle 23ai writes the pairs side by side and takes an empty list; it
// rejects a comma between pairs, a pair without its schema or with a third
// word, and a match string or schema that is a literal, a reserved word, or a
// qualified name (ORA-29512). A match string is a name or '*'; the schema
// may also be PUBLIC.
func (p *Parser) parseJavaResolverClause(opts *nodes.List) error {
	optStart := p.pos()
	p.advance() // consume RESOLVER
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	pairs := &nodes.List{}
	for p.cur.Type == '(' {
		pairStart := p.pos()
		p.advance() // consume (
		var match string
		switch {
		case p.cur.Type == '*':
			match = "*"
		case p.isCreateTableNameStart():
			match = p.cur.Str
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance()
		var schema string
		switch {
		case p.cur.Type == '-':
			schema = "-"
		case p.isCreateTableNameStart(), p.isKeywordStr("PUBLIC"):
			schema = p.cur.Str
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if err := p.expectCreateTableToken(')'); err != nil {
			return err
		}
		pairs.Items = append(pairs.Items, &nodes.DDLOption{
			Key:   match,
			Value: schema,
			Loc:   nodes.Loc{Start: pairStart, End: p.prev.End},
		})
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return err
	}
	opts.Items = append(opts.Items, &nodes.DDLOption{
		Key:   "RESOLVER",
		Items: pairs,
		Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
	})
	return nil
}

// parseJavaUsingClause parses the USING clause of CREATE JAVA:
//
//	USING { BFILE ( directory_object_name , server_file_name )
//	      | { CLOB | BLOB | BFILE } subquery
//	      | 'key_for_BLOB' }
//
// The subquery may stand bare or in parentheses, and may omit the SELECT
// keyword for backward compatibility (USING CLOB text FROM t, accepted by
// Oracle 23ai). addOption records the clause as a USING option: its value is
// BFILE, CLOB, or BLOB, or the key string, and its items hold the directory
// and file names or the subquery.
func (p *Parser) parseJavaUsingClause(addOption func(key, value string, optStart int, items *nodes.List)) error {
	optStart := p.pos()
	p.advance() // consume USING
	switch {
	case p.cur.Type == tokSCONST:
		key := p.cur.Str
		p.advance()
		addOption("USING", key, optStart, nil)
	case p.isKeywordStr("BFILE") && p.peekNext().Type == '(':
		p.advance() // consume BFILE
		p.advance() // consume (
		dir := ""
		dirLoc := nodes.Loc{Start: -1, End: -1}
		if p.isCreateTableNameStart() {
			dir = p.cur.Str
			dirLoc = nodes.Loc{Start: p.cur.Loc, End: p.cur.End}
			p.advance()
		}
		if err := p.expectCreateTableToken(','); err != nil {
			return err
		}
		if p.cur.Type != tokSCONST {
			return p.syntaxErrorAtCur()
		}
		file := p.cur.Str
		fileLoc := nodes.Loc{Start: p.cur.Loc, End: p.cur.End}
		p.advance()
		if err := p.expectCreateTableToken(')'); err != nil {
			return err
		}
		items := &nodes.List{}
		if dir != "" {
			items.Items = append(items.Items, &nodes.StringLiteral{Val: dir, Loc: dirLoc})
		}
		items.Items = append(items.Items, &nodes.StringLiteral{Val: file, Loc: fileLoc})
		addOption("USING", "BFILE", optStart, items)
	case p.isKeywordStr("CLOB"), p.isKeywordStr("BLOB"), p.isKeywordStr("BFILE"):
		kind := p.cur.Str
		p.advance()
		if p.isStatementEnd() {
			// Oracle reads USING CLOB with no subquery as an acceptance
			// (compat.tsv's row). The key is not the quoted-key form.
			addOption("USING_SOURCE", kind, optStart, nil)
			return nil
		}
		if kind == "BFILE" && p.isCreateTableNameStart() && p.peekNext().Type == ',' {
			// BFILE followed by a directory and a file name without its
			// parentheses is accepted too (compat.tsv's USING BFILE a, 's')
			// row).
			dir := p.cur.Str
			dirLoc := nodes.Loc{Start: p.cur.Loc, End: p.cur.End}
			p.advance()
			if err := p.expectCreateTableToken(','); err != nil {
				return err
			}
			if p.cur.Type != tokSCONST {
				return p.syntaxErrorAtCur()
			}
			file := p.cur.Str
			fileLoc := nodes.Loc{Start: p.cur.Loc, End: p.cur.End}
			p.advance()
			if err := p.expectCreateTableToken(')'); err != nil {
				return err
			}
			addOption("USING", "BFILE", optStart, &nodes.List{Items: []nodes.Node{&nodes.StringLiteral{Val: dir, Loc: dirLoc}, &nodes.StringLiteral{Val: file, Loc: fileLoc}}})
			return nil
		}
		if p.isIdentLike() && (p.peekNext().Type == kwSELECT || p.peekNext().Type == kwWITH) {
			// A word before the legacy subquery is accepted (compat.tsv's
			// USING CLOB zzz SELECT 1 FROM dual row).
			p.advance()
		}
		parens := p.cur.Type == '('
		if parens {
			p.advance() // consume (
		}
		var sub *nodes.SelectStmt
		var err error
		if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
			sub, err = p.parseSelectStmt()
		} else {
			// The legacy subquery without the SELECT keyword.
			sub, err = p.parseSelectAfterKeyword(&nodes.SelectStmt{
				TargetList: &nodes.List{},
				Loc:        nodes.Loc{Start: p.pos()},
			})
		}
		if err != nil {
			return err
		}
		if parens {
			if err := p.expectCreateTableToken(')'); err != nil {
				return err
			}
		}
		addOption("USING", kind, optStart, &nodes.List{Items: []nodes.Node{sub}})
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseAlterJavaStmt parses an ALTER JAVA statement. It is called after ALTER
// JAVA.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-JAVA.html
//
//	ALTER JAVA { SOURCE | CLASS } [ schema. ] object_name
//	    [ resolver_clause ]
//	    { RESOLVE | COMPILE | invoker_rights_clause }
//
// The page lists IF EXISTS and takes the invoker_rights_clause before RESOLVE
// or COMPILE, but Oracle 23ai rejects IF EXISTS (ORA-02000) and takes one
// action only: the invoker_rights_clause in place of RESOLVE or COMPILE, and
// nothing after any of them (ORA-03049). Without an action the statement
// fails at its end (ORA-00922).
func (p *Parser) parseAlterJavaStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_JAVA,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	addOption := func(key, value string, optStart int, items *nodes.List) {
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key:   key,
			Value: value,
			Items: items,
			Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}

	// { SOURCE | CLASS }
	if !p.isKeywordStr("SOURCE") && !p.isKeywordStr("CLASS") {
		return nil, p.syntaxErrorAtCur()
	}
	optStart := p.pos()
	kind := p.cur.Str
	p.advance()
	addOption("JAVA_TYPE", kind, optStart, nil)

	var err error
	if stmt.Name, err = p.parseRequiredObjectName(); err != nil {
		return nil, err
	}

	// [ resolver_clause ]
	if p.isKeywordStr("RESOLVER") {
		if err := p.parseJavaResolverClause(opts); err != nil {
			return nil, err
		}
	}

	// { RESOLVE | COMPILE | invoker_rights_clause }
	switch {
	case p.isKeywordStr("RESOLVE"), p.isKeywordStr("COMPILE"):
		optStart := p.pos()
		action := p.cur.Str
		p.advance()
		addOption(action, "", optStart, nil)
	case p.isKeywordStr("AUTHID"):
		if err := p.parseJavaAuthid(addOption); err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Options = opts
	stmt.Loc.End = p.prev.End
	return stmt, nil
}
