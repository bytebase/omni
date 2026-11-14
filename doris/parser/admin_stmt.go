package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// ADMIN statements
// ---------------------------------------------------------------------------

// parseAdminStmt parses an ADMIN statement. ADMIN has been consumed by the
// caller; cur is the word after it.
//
// Doris reads a statement with its ANTLR parser and refuses what that grammar
// refuses. The forms of unsupportedAdminStatement (DorisParser.g4:485) are
// then handed to the legacy CUP parser, whose admin_stmt
// (DorisLegacyParser.g4:3019) has to take them as well, so each of those
// forms reads here as what both grammars accept: ADMIN REBALANCE DISK ON
// ('a') is legacy-valid and refused by the ANTLR grammar, which wants exactly
// two backends; ADMIN SET REPLICA STATUS PROPERTIES(a = a) is ANTLR-valid and
// refused by the legacy key_value_map, which wants strings. The
// supportedAdminStatement (DorisParser.g4:512) and showConfig
// (DorisParser.g4:48) forms are read by the ANTLR parser alone.
func (p *Parser) parseAdminStmt(adminLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.AdminStmt{}
	var err error
	switch p.cur.Kind {
	case kwSHOW:
		err = p.parseAdminShow(stmt)
	case kwSET:
		err = p.parseAdminSet(stmt)
	case kwREPAIR:
		// ADMIN REPAIR TABLE baseTableRef
		err = p.parseAdminTableRefForm(stmt, "TABLE", kwTABLE)
	case kwCANCEL:
		err = p.parseAdminCancel(stmt)
	case kwCOMPACT:
		// ADMIN COMPACT TABLE baseTableRef wildWhere?
		if err = p.parseAdminTableRefForm(stmt, "TABLE", kwTABLE); err == nil {
			err = p.parseAdminWildWhere(stmt)
		}
	case kwCHECK:
		// ADMIN CHECK tabletList propertyClause? (DorisParser.g4:495); the
		// legacy form spells the list TABLET '(' integer_list ')'.
		p.adminVerb(stmt, "TABLET")
		if _, err = p.expect(kwTABLET); err == nil {
			start := p.cur.Loc.Start
			if err = p.parseAdminIntegerList(); err == nil {
				p.adminArgs(stmt, start)
				err = p.parseAdminStringProperties(stmt)
			}
		}
	case kwREBALANCE:
		// ADMIN REBALANCE DISK (ON ('be', 'be'))?
		p.adminVerb(stmt, "DISK")
		if _, err = p.expect(kwDISK); err == nil {
			err = p.parseAdminBackendsOn(stmt)
		}
	case kwCLEAN:
		// ADMIN CLEAN TRASH (ON ('be', 'be'))?
		p.adminVerb(stmt, "TRASH")
		if _, err = p.expect(kwTRASH); err == nil {
			err = p.parseAdminBackendsOn(stmt)
		}
	case kwDIAGNOSE:
		// ADMIN DIAGNOSE TABLET INTEGER_VALUE
		p.adminVerb(stmt, "TABLET")
		if _, err = p.expect(kwTABLET); err == nil {
			err = p.parseAdminTabletID(stmt)
		}
	case kwCOPY:
		// ADMIN COPY TABLET INTEGER_VALUE propertyClause?
		p.adminVerb(stmt, "TABLET")
		if _, err = p.expect(kwTABLET); err == nil {
			if err = p.parseAdminTabletID(stmt); err == nil {
				err = p.parseAdminStringProperties(stmt)
			}
		}
	case kwROTATE:
		// ADMIN ROTATE TDE ROOT KEY propertyClause? (DorisParser.g4:513), a
		// supported statement: the ANTLR propertyItem, identifiers included.
		p.adminVerb(stmt, "TDE ROOT KEY")
		if err = p.expectAdminWords(kwTDE, kwROOT, kwKEY); err == nil && p.cur.Kind == kwPROPERTIES {
			p.advance() // consume PROPERTIES
			if _, err = p.expect(int('(')); err == nil {
				stmt.Properties, err = p.parsePropertyItems(p.parsePropertyItem)
			}
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}
	stmt.Loc = ast.Loc{Start: adminLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// parseAdminShow reads the ADMIN SHOW forms. cur is SHOW.
func (p *Parser) parseAdminShow(stmt *ast.AdminStmt) error {
	p.adminVerb(stmt, "")
	switch p.cur.Kind {
	case kwFRONTEND, kwBACKEND:
		// ADMIN? SHOW (FRONTEND | BACKEND) CONFIG (LIKE valueExpression)?
		// (FROM INTEGER_VALUE)? (DorisParser.g4:48), read by the ANTLR parser
		// alone.
		stmt.Object = strings.ToUpper(p.advance().Str) + " CONFIG"
		if _, err := p.expect(kwCONFIG); err != nil {
			return err
		}
		start := p.cur.Loc.Start
		if p.cur.Kind == kwLIKE {
			p.advance() // consume LIKE
			if _, err := p.parseExprPrec(bpCmp); err != nil {
				return err
			}
		}
		if p.cur.Kind == kwFROM {
			p.advance() // consume FROM
			if _, err := p.expect(tokInt); err != nil {
				return err
			}
		}
		p.adminArgs(stmt, start)
		return nil
	case kwREPLICA:
		p.advance() // consume REPLICA
		switch p.cur.Kind {
		case kwSTATUS:
			// ADMIN SHOW REPLICA STATUS FROM baseTableRef wildWhere?
			stmt.Object = "REPLICA STATUS"
			p.advance() // consume STATUS
			if _, err := p.expect(kwFROM); err != nil {
				return err
			}
			if err := p.parseAdminBaseTableRef(stmt); err != nil {
				return err
			}
			return p.parseAdminWildWhere(stmt)
		case kwDISTRIBUTION:
			// ADMIN SHOW REPLICA DISTRIBUTION FROM baseTableRef
			stmt.Object = "REPLICA DISTRIBUTION"
			p.advance() // consume DISTRIBUTION
			if _, err := p.expect(kwFROM); err != nil {
				return err
			}
			return p.parseAdminBaseTableRef(stmt)
		}
		return p.syntaxErrorAtCur()
	case kwTABLET:
		// ADMIN SHOW TABLET STORAGE FORMAT VERBOSE?
		stmt.Object = "TABLET STORAGE FORMAT"
		p.advance() // consume TABLET
		if err := p.expectAdminWords(kwSTORAGE, kwFORMAT); err != nil {
			return err
		}
		if p.cur.Kind == kwVERBOSE {
			p.adminArgs(stmt, p.advance().Loc.Start)
		}
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseAdminSet reads the ADMIN SET forms. cur is SET.
func (p *Parser) parseAdminSet(stmt *ast.AdminStmt) error {
	p.adminVerb(stmt, "")
	switch p.cur.Kind {
	case kwREPLICA:
		// ADMIN SET REPLICA (STATUS | VERSION) PROPERTIES '(' ... ')': the
		// legacy key_value_map takes string pairs only.
		p.advance() // consume REPLICA
		if p.cur.Kind != kwSTATUS && p.cur.Kind != kwVERSION {
			return p.syntaxErrorAtCur()
		}
		stmt.Object = "REPLICA " + strings.ToUpper(p.advance().Str)
		if p.cur.Kind != kwPROPERTIES {
			return p.syntaxErrorAtCur()
		}
		return p.parseAdminStringProperties(stmt)
	case kwFRONTEND, kwALL:
		// ADMIN SET (FRONTEND | ALL FRONTENDS) CONFIG ('(' ... ')')? ALL?
		// (DorisParser.g4:493). The legacy grammar takes a trailing ALL
		// after FRONTEND only, and string pairs only.
		all := p.cur.Kind == kwALL
		p.advance() // consume FRONTEND / ALL
		stmt.Object = "FRONTEND CONFIG"
		if all {
			stmt.Object = "ALL FRONTENDS CONFIG"
			if _, err := p.expect(kwFRONTENDS); err != nil {
				return err
			}
		}
		if _, err := p.expect(kwCONFIG); err != nil {
			return err
		}
		if p.cur.Kind == int('(') {
			p.advance() // consume '('
			props, err := p.parsePropertyItems(p.parseStringProperty)
			if err != nil {
				return err
			}
			stmt.Properties = props
		}
		if !all && p.cur.Kind == kwALL {
			p.adminArgs(stmt, p.advance().Loc.Start)
		}
		return nil
	case kwTABLE:
		// ADMIN SET TABLE name PARTITION VERSION propertyClause?
		// ADMIN SET TABLE name STATUS propertyClause?
		p.advance() // consume TABLE
		if err := p.parseAdminTableName(stmt); err != nil {
			return err
		}
		switch p.cur.Kind {
		case kwPARTITION:
			stmt.Object = "TABLE PARTITION VERSION"
			p.advance() // consume PARTITION
			if _, err := p.expect(kwVERSION); err != nil {
				return err
			}
		case kwSTATUS:
			stmt.Object = "TABLE STATUS"
			p.advance() // consume STATUS
		default:
			return p.syntaxErrorAtCur()
		}
		return p.parseAdminStringProperties(stmt)
	}
	// ADMIN SET ENCRYPTION ROOT KEY (DorisParser.g4:507) has no legacy
	// admin_stmt form, so Doris refuses it whatever follows.
	return p.syntaxErrorAtCur()
}

// parseAdminCancel reads ADMIN CANCEL REPAIR TABLE baseTableRef and ADMIN
// CANCEL REBALANCE DISK (ON (...))?. cur is CANCEL.
func (p *Parser) parseAdminCancel(stmt *ast.AdminStmt) error {
	p.adminVerb(stmt, "")
	switch p.cur.Kind {
	case kwREPAIR:
		p.advance() // consume REPAIR
		stmt.Object = "REPAIR TABLE"
		if _, err := p.expect(kwTABLE); err != nil {
			return err
		}
		return p.parseAdminBaseTableRef(stmt)
	case kwREBALANCE:
		p.advance() // consume REBALANCE
		stmt.Object = "REBALANCE DISK"
		if _, err := p.expect(kwDISK); err != nil {
			return err
		}
		return p.parseAdminBackendsOn(stmt)
	}
	return p.syntaxErrorAtCur()
}

// adminVerb records the verb at cur (and the object, when the form fixes one)
// and consumes the verb.
func (p *Parser) adminVerb(stmt *ast.AdminStmt, object string) {
	stmt.Verb = strings.ToUpper(p.advance().Str)
	stmt.Object = object
}

// expectAdminWords consumes kinds in order.
func (p *Parser) expectAdminWords(kinds ...TokenKind) error {
	for _, k := range kinds {
		if _, err := p.expect(k); err != nil {
			return err
		}
	}
	return nil
}

// adminArgs appends the source text from start to the last consumed token to
// stmt.Args, which holds what a form takes beyond its verb, object, target
// and properties.
func (p *Parser) adminArgs(stmt *ast.AdminStmt, start int) {
	if p.prev.Loc.End <= start {
		return
	}
	text := p.sourceText(ast.Loc{Start: start, End: p.prev.Loc.End})
	if stmt.Args != "" {
		stmt.Args += " "
	}
	stmt.Args += text
}

// parseAdminTableRefForm reads verb, object keyword and a baseTableRef: ADMIN
// REPAIR TABLE t, ADMIN COMPACT TABLE t.
func (p *Parser) parseAdminTableRefForm(stmt *ast.AdminStmt, object string, objectKind TokenKind) error {
	p.adminVerb(stmt, object)
	if _, err := p.expect(objectKind); err != nil {
		return err
	}
	return p.parseAdminBaseTableRef(stmt)
}

// parseAdminTableName reads a table name into stmt.Target: the ANTLR
// multipartIdentifier held to the legacy table_name (DorisLegacyParser.g4:2175),
// which has at most three parts.
func (p *Parser) parseAdminTableName(stmt *ast.AdminStmt) error {
	start := p.cur
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return err
	}
	if len(name.Parts) > 3 {
		return syntaxErrorAt(start)
	}
	stmt.Target = name
	return nil
}

// parseAdminBaseTableRef reads the baseTableRef of DorisParser.g4:516, held to
// the legacy base_table_ref (DorisLegacyParser.g4:2256) that reads the same
// clauses in the same order:
//
//	name ('@' ident '(' ('k' = 'v', ...)? ')')?
//	(FOR VERSION AS OF INTEGER_VALUE | FOR TIME AS OF STRING_LITERAL)?
//	(TEMPORARY? PARTITION (ident | '(' ident, ... ')') | TEMPORARY? PARTITIONS '(' ident, ... ')')?
//	(TABLET '(' INTEGER_VALUE, ... ')')?
//	(AS? ident)?
//	(TABLESAMPLE '(' INTEGER_VALUE (PERCENT | ROWS) ')' (REPEATABLE INTEGER_VALUE)?)?
//	('[' ident, ... ']' | '/*+' ident, ... '*/')?
//
// The name goes to stmt.Target and the text of the clauses after it to
// stmt.Args.
func (p *Parser) parseAdminBaseTableRef(stmt *ast.AdminStmt) error {
	if err := p.parseAdminTableName(stmt); err != nil {
		return err
	}
	start := p.cur.Loc.Start

	// optScanParams (DorisParser.g4:1359) and opt_scan_params_ref
	// (DorisLegacyParser.g4:2194): the legacy list is string pairs, or empty.
	if p.cur.Kind == int('@') {
		p.advance() // consume '@'
		if _, _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if _, err := p.expect(int('(')); err != nil {
			return err
		}
		if p.cur.Kind == int(')') {
			p.advance() // consume ')'
		} else if _, err := p.parsePropertyItems(p.parseStringProperty); err != nil {
			return err
		}
	}

	// tableSnapshot (DorisParser.g4:1788) and table_snapshot
	// (DorisLegacyParser.g4:2267): a version is an integer, a time a string.
	if p.cur.Kind == kwFOR {
		p.advance() // consume FOR
		value := tokInt
		switch p.cur.Kind {
		case kwVERSION:
		case kwTIME:
			value = tokString
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume VERSION / TIME
		if err := p.expectAdminWords(kwAS, kwOF, value); err != nil {
			return err
		}
	}

	// specifiedPartition (DorisParser.g4:1663) and partition_names
	// (DorisLegacyParser.g4:2311). PARTITIONS is non-reserved, so it is the
	// clause only when its list follows; otherwise it is left to the alias.
	temporary := false
	if p.cur.Kind == kwTEMPORARY {
		if next := p.peekNext().Kind; next == kwPARTITION || next == kwPARTITIONS {
			p.advance() // consume TEMPORARY
			temporary = true
		}
	}
	switch {
	case p.cur.Kind == kwPARTITION:
		p.advance() // consume PARTITION
		if p.cur.Kind == int('(') {
			if err := p.parseAdminIdentList(int('('), int(')')); err != nil {
				return err
			}
		} else if _, _, err := p.parseIdentifier(); err != nil {
			return err
		}
	case p.cur.Kind == kwPARTITIONS && (temporary || p.peekNext().Kind == int('(')):
		p.advance() // consume PARTITIONS
		if err := p.parseAdminIdentList(int('('), int(')')); err != nil {
			return err
		}
	}

	// tabletList (DorisParser.g4:1480).
	if p.cur.Kind == kwTABLET {
		p.advance() // consume TABLET
		if err := p.parseAdminIntegerList(); err != nil {
			return err
		}
	}

	// tableAlias (DorisParser.g4:1393) and opt_table_alias
	// (DorisLegacyParser.g4:2286): AS? ident; neither the ANTLR column list
	// nor the legacy '=' ident is in both.
	if p.cur.Kind == kwAS {
		p.advance() // consume AS
		if _, _, err := p.parseIdentifier(); err != nil {
			return err
		}
	} else if isIdentifierToken(p.cur.Kind) {
		p.advance() // consume the alias
	}

	// sample (DorisParser.g4:1779) and table_sample
	// (DorisLegacyParser.g4:2330): the legacy form wants the method.
	if p.cur.Kind == kwTABLESAMPLE {
		p.advance() // consume TABLESAMPLE
		if err := p.expectAdminWords(int('('), tokInt); err != nil {
			return err
		}
		if p.cur.Kind != kwPERCENT && p.cur.Kind != kwROWS {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume PERCENT / ROWS
		if _, err := p.expect(int(')')); err != nil {
			return err
		}
		if p.cur.Kind == kwREPEATABLE {
			p.advance() // consume REPEATABLE
			if _, err := p.expect(tokInt); err != nil {
				return err
			}
		}
	}

	// relationHint (DorisParser.g4:1260) and opt_common_hints
	// (DorisLegacyParser.g4:2273).
	switch p.cur.Kind {
	case int('['):
		if err := p.parseAdminIdentList(int('['), int(']')); err != nil {
			return err
		}
	case tokHintStart:
		if err := p.parseAdminIdentList(tokHintStart, tokHintEnd); err != nil {
			return err
		}
	}

	p.adminArgs(stmt, start)
	return nil
}

// parseAdminIdentList reads open ident (',' ident)* close.
func (p *Parser) parseAdminIdentList(open, close TokenKind) error {
	if _, err := p.expect(open); err != nil {
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
	_, err := p.expect(close)
	return err
}

// parseAdminIntegerList reads '(' INTEGER_VALUE (',' INTEGER_VALUE)* ')', the
// list of tabletList (DorisParser.g4:1480).
func (p *Parser) parseAdminIntegerList() error {
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
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

// parseAdminTabletID reads the INTEGER_VALUE tablet id of ADMIN DIAGNOSE /
// COPY TABLET into stmt.Args.
func (p *Parser) parseAdminTabletID(stmt *ast.AdminStmt) error {
	tok, err := p.expect(tokInt)
	if err != nil {
		return err
	}
	p.adminArgs(stmt, tok.Loc.Start)
	return nil
}

// parseAdminWildWhere reads an optional wildWhere (DorisParser.g4:521), the
// legacy opt_wild_where (DorisLegacyParser.g4:1589): LIKE STRING_LITERAL or
// WHERE expression.
func (p *Parser) parseAdminWildWhere(stmt *ast.AdminStmt) error {
	start := p.cur.Loc.Start
	switch p.cur.Kind {
	case kwLIKE:
		p.advance() // consume LIKE
		if _, err := p.expect(tokString); err != nil {
			return err
		}
	case kwWHERE:
		p.advance() // consume WHERE
		if _, err := p.parseExpr(); err != nil {
			return err
		}
	default:
		return nil
	}
	p.adminArgs(stmt, start)
	return nil
}

// parseAdminBackendsOn reads the optional backend list of ADMIN REBALANCE
// DISK, ADMIN CANCEL REBALANCE DISK and ADMIN CLEAN TRASH:
//
//	(ON LEFT_PAREN STRING_LITERAL (COMMA STRING_LITERAL) RIGHT_PAREN)?
//
// (DorisParser.g4:496-501). The ANTLR grammar wants exactly two backends, no
// more and no fewer, and Doris holds a statement to it before the legacy
// string_list could take another count.
func (p *Parser) parseAdminBackendsOn(stmt *ast.AdminStmt) error {
	if p.cur.Kind != kwON {
		return nil
	}
	start := p.advance().Loc.Start // consume ON
	if err := p.expectAdminWords(int('('), tokString, int(','), tokString, int(')')); err != nil {
		return err
	}
	p.adminArgs(stmt, start)
	return nil
}

// parseAdminStringProperties reads an optional PROPERTIES '(' 'k' = 'v', ...
// ')' into stmt.Properties: the ANTLR propertyClause (DorisParser.g4:1377)
// held to the legacy properties (DorisLegacyParser.g4:1242), whose
// key_value_map takes string pairs only.
func (p *Parser) parseAdminStringProperties(stmt *ast.AdminStmt) error {
	if p.cur.Kind != kwPROPERTIES {
		return nil
	}
	props, err := p.parseProperties()
	if err != nil {
		return err
	}
	stmt.Properties = props
	return nil
}
