package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseMleIfNotExists parses the optional IF NOT EXISTS that follows MLE ENV
// or MLE MODULE. Oracle 23ai takes it there, not before MLE (ORA-00901,
// refused by the CREATE dispatcher); it rejects any other spelling at the IF
// (ORA-11543) and OR REPLACE beside it at the name that follows (ORA-11541).
func (p *Parser) parseMleIfNotExists(stmt *nodes.AdminDDLStmt) error {
	if p.cur.Type != kwIF {
		return nil
	}
	ifTok := p.cur
	p.advance() // consume IF
	if p.cur.Type != kwNOT || p.peekNext().Type != kwEXISTS {
		return p.syntaxErrorAtTok(ifTok)
	}
	p.advance() // consume NOT
	p.advance() // consume EXISTS
	if stmt.OrReplace {
		return p.syntaxErrorAtCur()
	}
	stmt.IfNotExists = true
	return nil
}

// isMleStringStart reports whether the current token is a string, national
// (N'...') or not, as an import name, options string, or version takes.
func (p *Parser) isMleStringStart() bool {
	return p.cur.Type == tokSCONST || p.cur.Type == tokNCHARLIT
}

// parseCreateMLEEnvStmt parses a CREATE MLE ENV statement. It is called after
// CREATE [ OR REPLACE ] MLE ENV.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/create-mle-env.html
//
//	CREATE [ OR REPLACE ] MLE ENV [ IF NOT EXISTS ]
//	    [ schema. ] environment_name
//	    { CLONE [ schema. ] existing_environment
//	    | [ imports_clause ] [ language_options_clause ] }
//
//	imports_clause:
//	    IMPORTS ( import_item [, import_item ]... )
//
//	import_item:
//	    import_name MODULE [ schema. ] module_name
//
//	language_options_clause:
//	    LANGUAGE OPTIONS language_options_string
//
// The page puts IF NOT EXISTS and PURE before MLE and lets CLONE precede the
// other clauses. Oracle 23ai rejects both words there (ORA-00901), takes IF
// NOT EXISTS after ENV, and takes nothing after CLONE's name (ORA-03049). An
// import name is a string (ORA-04114), as is the options string (ORA-04113).
func (p *Parser) parseCreateMLEEnvStmt(start int, orReplace bool) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_MLE_ENV,
		OrReplace:  orReplace,
		Loc:        nodes.Loc{Start: start},
	}
	if err := p.parseMleIfNotExists(stmt); err != nil {
		return nil, err
	}
	name, err := p.parsePlainObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	opts := &nodes.List{}
	if p.isKeywordStr("CLONE") {
		optStart := p.pos()
		p.advance() // consume CLONE
		clone, err := p.parsePlainObjectName()
		if err != nil {
			return nil, err
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key:   "CLONE",
			Value: clone.Name,
			Items: &nodes.List{Items: []nodes.Node{clone}},
			Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
		})
	} else {
		if p.isKeywordStr("IMPORTS") {
			opt, err := p.parseMleImportsClause()
			if err != nil {
				return nil, err
			}
			opts.Items = append(opts.Items, opt)
		}
		if p.isKeywordStr("LANGUAGE") {
			optStart := p.pos()
			p.advance() // consume LANGUAGE
			if err := p.expectCreateTableWord("OPTIONS"); err != nil {
				return nil, err
			}
			if !p.isMleStringStart() {
				return nil, p.syntaxErrorAtCur()
			}
			value := p.cur.Str
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key:   "LANGUAGE OPTIONS",
				Value: value,
				Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
			})
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseMleImportsClause parses the imports_clause of CREATE MLE ENV, from its
// IMPORTS, as an IMPORTS option whose items hold one option per import, keyed
// by the import name, with the module as its item.
func (p *Parser) parseMleImportsClause() (*nodes.DDLOption, error) {
	optStart := p.pos()
	p.advance() // consume IMPORTS
	if err := p.expectCreateTableToken('('); err != nil {
		return nil, err
	}
	items := &nodes.List{}
	for {
		itemStart := p.pos()
		if !p.isMleStringStart() {
			return nil, p.syntaxErrorAtCur()
		}
		importName := p.cur.Str
		p.advance()
		if err := p.expectCreateTableWord("MODULE"); err != nil {
			return nil, err
		}
		module, err := p.parsePlainObjectName()
		if err != nil {
			return nil, err
		}
		items.Items = append(items.Items, &nodes.DDLOption{
			Key:   importName,
			Value: module.Name,
			Items: &nodes.List{Items: []nodes.Node{module}},
			Loc:   nodes.Loc{Start: itemStart, End: p.prev.End},
		})
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ,
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return nil, err
	}
	return &nodes.DDLOption{
		Key:   "IMPORTS",
		Items: items,
		Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
	}, nil
}

// parseCreateMLEModuleStmt parses a CREATE MLE MODULE statement. It is called
// after CREATE [ OR REPLACE ] MLE MODULE.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/create-mle-module.html
//
//	CREATE [ OR REPLACE ] MLE MODULE [ IF NOT EXISTS ]
//	    [ schema. ] module_name
//	    LANGUAGE language_name
//	    [ VERSION version_string ]
//	    { USING { CLOB ( subquery )
//	            | BLOB ( subquery )
//	            | BFILE ( directory_object_name , server_file_name ) }
//	    | AS source_code }
//
// The page puts IF NOT EXISTS before MLE, which Oracle 23ai rejects
// (ORA-00901), and names JAVASCRIPT as the language; Oracle takes any name
// there and checks it later. The version is a string (ORA-04112). BFILE takes
// a directory and a file name only (ORA-02000 for a subquery). The subquery
// stands alone in its parentheses: Oracle rejects a WITH clause, ORDER BY, or
// FOR UPDATE on it, a query without SELECT, and anything after the closing
// parenthesis (ORA-29506, at CLOB or BLOB).
func (p *Parser) parseCreateMLEModuleStmt(start int, orReplace bool) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_MLE_MODULE,
		OrReplace:  orReplace,
		Loc:        nodes.Loc{Start: start},
	}
	if err := p.parseMleIfNotExists(stmt); err != nil {
		return nil, err
	}
	name, err := p.parsePlainObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	opts := &nodes.List{}

	// LANGUAGE language_name
	optStart := p.pos()
	if !p.isKeywordStr("LANGUAGE") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume LANGUAGE
	if !p.isCreateTableNameStart() {
		return nil, p.syntaxErrorAtCur()
	}
	language := p.cur.Str
	p.advance()
	opts.Items = append(opts.Items, &nodes.DDLOption{Key: "LANGUAGE", Value: language, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})

	// [ VERSION version_string ]
	if p.isKeywordStr("VERSION") {
		optStart := p.pos()
		p.advance() // consume VERSION
		if !p.isMleStringStart() {
			return nil, p.syntaxErrorAtCur()
		}
		value := p.cur.Str
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "VERSION", Value: value, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
	}

	switch {
	case p.cur.Type == kwUSING:
		opt, err := p.parseMleUsingClause()
		if err != nil {
			return nil, err
		}
		opts.Items = append(opts.Items, opt)
	case p.cur.Type == kwAS:
		// module_text is JavaScript source, not SQL; see
		// embeddedSourceHead. Take it verbatim.
		if err := p.parseSourceText(opts); err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Options = opts
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseMleUsingClause parses the USING clause of CREATE MLE MODULE, from its
// USING, as a USING option whose value is CLOB, BLOB, or BFILE and whose
// items hold the subquery or the directory and file names.
func (p *Parser) parseMleUsingClause() (*nodes.DDLOption, error) {
	optStart := p.pos()
	p.advance() // consume USING
	kindTok := p.cur
	switch {
	case p.isKeywordStr("BFILE"):
		p.advance() // consume BFILE
		if err := p.expectCreateTableToken('('); err != nil {
			return nil, err
		}
		if !p.isCreateTableNameStart() {
			return nil, p.syntaxErrorAtCur()
		}
		dir := p.cur.Str
		p.advance()
		if err := p.expectCreateTableToken(','); err != nil {
			return nil, err
		}
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		file := p.cur.Str
		p.advance()
		if err := p.expectCreateTableToken(')'); err != nil {
			return nil, err
		}
		return &nodes.DDLOption{
			Key:   "USING",
			Value: "BFILE",
			Items: &nodes.List{Items: []nodes.Node{&nodes.String{Str: dir}, &nodes.String{Str: file}}},
			Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
		}, nil
	case p.isKeywordStr("CLOB"), p.isKeywordStr("BLOB"):
		p.advance() // consume CLOB or BLOB
		if err := p.expectCreateTableToken('('); err != nil {
			return nil, err
		}
		if p.cur.Type != kwSELECT && p.cur.Type != '(' {
			return nil, p.syntaxErrorAtTok(kindTok)
		}
		sub, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		if sub.WithClause != nil || sub.OrderBy != nil || sub.ForUpdate != nil || p.cur.Type != ')' {
			return nil, p.syntaxErrorAtTok(kindTok)
		}
		p.advance() // consume )
		if p.cur.Type != ';' && p.cur.Type != tokEOF {
			return nil, p.syntaxErrorAtTok(kindTok)
		}
		return &nodes.DDLOption{
			Key:   "USING",
			Value: kindTok.Str,
			Items: &nodes.List{Items: []nodes.Node{sub}},
			Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
		}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}
