package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// keywordCategory classifies MySQL 8.0 keywords into 6 categories matching
// the sql_yacc.yy grammar's identifier context rules.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/keywords.html
// Ref: mysql-server sql/sql_yacc.yy (ident, label_ident, role_ident, lvalue_ident rules)
type keywordCategory int

const (
	kwCatReserved    keywordCategory = iota // Cannot be used as identifiers without quoting
	kwCatUnambiguous                        // ident_keywords_unambiguous — allowed in all identifier contexts
	kwCatAmbiguous1                         // ident_keywords_ambiguous_1_roles_and_labels — NOT allowed as label or role
	kwCatAmbiguous2                         // ident_keywords_ambiguous_2_labels — NOT allowed as label
	kwCatAmbiguous3                         // ident_keywords_ambiguous_3_roles — NOT allowed as role
	kwCatAmbiguous4                         // ident_keywords_ambiguous_4_system_variables — NOT allowed as lvalue
)

// keywordCategories maps keyword token types to their category. Keywords not
// present in this map are not registered keywords (they lex as tokIDENT).
//
// All 65 original reserved keywords are migrated here with kwCatReserved.
// All other registered keywords default to kwCatUnambiguous and will be
// refined to their correct ambiguous category in later sections.
var keywordCategories = map[int]keywordCategory{
	kwSELECT:     kwCatReserved,
	kwINSERT:     kwCatReserved,
	kwUPDATE:     kwCatReserved,
	kwDELETE:     kwCatReserved,
	kwFROM:       kwCatReserved,
	kwWHERE:      kwCatReserved,
	kwCREATE:     kwCatReserved,
	kwDROP:       kwCatReserved,
	kwALTER:      kwCatReserved,
	kwTABLE:      kwCatReserved,
	kwINTO:       kwCatReserved,
	kwVALUES:     kwCatReserved,
	kwSET:        kwCatReserved,
	kwJOIN:       kwCatReserved,
	kwLEFT:       kwCatReserved,
	kwRIGHT:      kwCatReserved,
	kwINNER:      kwCatReserved,
	kwOUTER:      kwCatReserved,
	kwON:         kwCatReserved,
	kwAND:        kwCatReserved,
	kwOR:         kwCatReserved,
	kwNOT:        kwCatReserved,
	kwNULL:       kwCatReserved,
	kwTRUE:       kwCatReserved,
	kwFALSE:      kwCatReserved,
	kwIN:         kwCatReserved,
	kwBETWEEN:    kwCatReserved,
	kwLIKE:       kwCatReserved,
	kwORDER:      kwCatReserved,
	kwGROUP:      kwCatReserved,
	kwBY:         kwCatReserved,
	kwHAVING:     kwCatReserved,
	kwLIMIT:      kwCatReserved,
	kwAS:         kwCatReserved,
	kwIS:         kwCatReserved,
	kwEXISTS_KW:  kwCatReserved,
	kwCASE:       kwCatReserved,
	kwWHEN:       kwCatReserved,
	kwTHEN:       kwCatReserved,
	kwELSE:       kwCatReserved,
	kwEND:        kwCatAmbiguous2, // demoted from reserved — MySQL 8.0 ambiguous_2 (not label)
	kwIF:         kwCatReserved,
	kwFOR:        kwCatReserved,
	kwWHILE:      kwCatReserved,
	kwINDEX:      kwCatReserved,
	kwKEY:        kwCatReserved,
	kwPRIMARY:    kwCatReserved,
	kwFOREIGN:    kwCatReserved,
	kwREFERENCES: kwCatReserved,
	kwCONSTRAINT: kwCatReserved,
	kwUNIQUE:     kwCatReserved,
	kwCHECK:      kwCatReserved,
	kwDEFAULT:    kwCatReserved,
	kwCOLUMN:     kwCatReserved,
	kwADD:        kwCatReserved,
	kwCHANGE:     kwCatReserved,
	kwMODIFY:     kwCatReserved,
	kwRENAME:     kwCatReserved,
	kwGRANT:      kwCatReserved,
	kwREVOKE:     kwCatReserved,
	kwALL:        kwCatReserved,
	kwDISTINCT:   kwCatReserved,
	kwUNION:      kwCatReserved,
	kwINTERSECT:      kwCatReserved,
	kwEXCEPT:         kwCatReserved,
	kwACCESSIBLE:     kwCatReserved,
	kwASENSITIVE:     kwCatReserved,
	kwCUBE:           kwCatReserved,
	kwCUME_DIST:      kwCatReserved,
	kwDENSE_RANK:     kwCatReserved,
	kwDUAL:           kwCatReserved,
	kwFIRST_VALUE:    kwCatReserved,
	kwGROUPING:       kwCatReserved,
	kwINSENSITIVE:    kwCatReserved,
	kwLAG:            kwCatReserved,
	kwLAST_VALUE:     kwCatReserved,
	kwLEAD:           kwCatReserved,
	kwNTH_VALUE:      kwCatReserved,
	kwNTILE:          kwCatReserved,
	kwOF:             kwCatReserved,
	kwOPTIMIZER_COSTS: kwCatReserved,
	kwPERCENT_RANK:   kwCatReserved,
	kwRANK:           kwCatReserved,
	kwROW_NUMBER:     kwCatReserved,
	kwSENSITIVE:      kwCatReserved,
	kwSPECIFIC:       kwCatReserved,
	kwUSAGE:          kwCatReserved,
	kwVARYING:            kwCatReserved,
	kwDAY_HOUR:           kwCatReserved,
	kwDAY_MICROSECOND:    kwCatReserved,
	kwDAY_MINUTE:         kwCatReserved,
	kwDAY_SECOND:         kwCatReserved,
	kwHOUR_MICROSECOND:   kwCatReserved,
	kwHOUR_MINUTE:        kwCatReserved,
	kwHOUR_SECOND:        kwCatReserved,
	kwMINUTE_MICROSECOND: kwCatReserved,
	kwMINUTE_SECOND:      kwCatReserved,
	kwSECOND_MICROSECOND: kwCatReserved,
	kwYEAR_MONTH:         kwCatReserved,
	kwUTC_DATE:           kwCatReserved,
	kwUTC_TIME:           kwCatReserved,
	kwUTC_TIMESTAMP:      kwCatReserved,
	kwMAXVALUE:           kwCatReserved,
	kwNO_WRITE_TO_BINLOG: kwCatReserved,
	kwIO_AFTER_GTIDS:     kwCatReserved,
	kwIO_BEFORE_GTIDS:    kwCatReserved,
	kwSQLEXCEPTION:       kwCatReserved,
	kwSQLSTATE:           kwCatReserved,
	kwSQLWARNING:         kwCatReserved,
	kwCROSS:              kwCatReserved,
	kwNATURAL:            kwCatReserved,
	kwUSING:              kwCatReserved,
	kwASC:                kwCatReserved,
	kwDESC:               kwCatReserved,
	kwTO:                 kwCatReserved,
	kwDIV:                kwCatReserved,
	kwMOD:                kwCatReserved,
	kwXOR:                kwCatReserved,
	kwREGEXP:             kwCatReserved,
	kwBINARY:             kwCatReserved,
	kwINTERVAL:           kwCatReserved,
	kwMATCH:              kwCatReserved,
	kwCURRENT_DATE:       kwCatReserved,
	kwCURRENT_TIME:       kwCatReserved,
	kwCURRENT_TIMESTAMP:  kwCatReserved,
	kwCURRENT_USER:       kwCatReserved,
	kwDATABASE:           kwCatReserved,
	kwFUNCTION:           kwCatReserved,
	kwPROCEDURE:          kwCatReserved,
	kwTRIGGER:            kwCatReserved,
	kwPARTITION:          kwCatReserved,
	kwRANGE:              kwCatReserved,
	kwROW:                kwCatReserved,
	kwROWS:               kwCatReserved,
	kwOVER:               kwCatReserved,
	kwWINDOW:             kwCatReserved,
	kwFORCE:              kwCatReserved,
	kwCONVERT:            kwCatReserved,
	kwCAST:               kwCatReserved,
	kwWITH:               kwCatReserved,
	kwREPLACE:            kwCatReserved,
	kwIGNORE:             kwCatReserved,
	kwLOAD:               kwCatReserved,
	kwUSE:                kwCatReserved,
	kwKILL:               kwCatReserved,
	kwEXPLAIN:            kwCatReserved,
	kwSPATIAL:            kwCatReserved,
	kwFULLTEXT:           kwCatReserved,
	kwOUTFILE:            kwCatReserved,
	kwGEOMETRY:           kwCatUnambiguous,
	kwPOINT:              kwCatUnambiguous,
	kwLINESTRING:         kwCatUnambiguous,
	kwPOLYGON:            kwCatUnambiguous,
	kwMULTIPOINT:         kwCatUnambiguous,
	kwMULTILINESTRING:    kwCatUnambiguous,
	kwMULTIPOLYGON:       kwCatUnambiguous,
	kwGEOMETRYCOLLECTION: kwCatUnambiguous,
	kwSERIAL:             kwCatUnambiguous,
	kwNATIONAL:           kwCatUnambiguous,
	kwNCHAR:              kwCatUnambiguous,
	kwNVARCHAR:           kwCatUnambiguous,
	kwSIGNED:             kwCatAmbiguous2,
	kwPRECISION:          kwCatUnambiguous,
	kwBOOL:               kwCatUnambiguous,
	kwBOOLEAN:            kwCatUnambiguous,
	kwSRID:               kwCatUnambiguous,
	kwENFORCED:           kwCatUnambiguous,
	kwLESS:               kwCatUnambiguous,
	kwTHAN:               kwCatUnambiguous,
	kwSUBPARTITIONS:      kwCatUnambiguous,
	kwLEAVES:             kwCatUnambiguous,
	kwPARSER:             kwCatUnambiguous,
	kwCOMPRESSION:        kwCatUnambiguous,
	kwINSERT_METHOD:      kwCatUnambiguous,
	kwACTION:             kwCatUnambiguous,
	kwPARTIAL:            kwCatUnambiguous,
	kwFORMAT:             kwCatUnambiguous,
	kwXML:                kwCatUnambiguous,
	kwCONCURRENT:         kwCatUnambiguous,
	kwWORK:               kwCatUnambiguous,
	kwXID:                kwCatUnambiguous,
	kwEXPORT:             kwCatUnambiguous,
	kwUPGRADE:            kwCatUnambiguous,
	kwFAST:               kwCatUnambiguous,
	kwMEDIUM:             kwCatUnambiguous,
	kwCHANGED:            kwCatUnambiguous,
	kwCODE:               kwCatUnambiguous,
	kwEVENTS:             kwCatUnambiguous,
	kwINDEXES:            kwCatUnambiguous,
	kwGRANTS:             kwCatUnambiguous,
	kwTRIGGERS:           kwCatUnambiguous,
	kwSCHEMAS:            kwCatUnambiguous,
	kwPARTITIONS:         kwCatUnambiguous,
	kwHOSTS:              kwCatUnambiguous,
	kwMUTEX:              kwCatUnambiguous,
	kwPROFILE:            kwCatUnambiguous,
	kwREPLICAS:           kwCatUnambiguous,
	kwNAMES:              kwCatUnambiguous,
	kwACCOUNT:            kwCatUnambiguous,
	kwOPTION:             kwCatUnambiguous,
	kwPROXY:              kwCatAmbiguous3,
	kwROUTINE:            kwCatUnambiguous,
	kwEXPIRE:             kwCatUnambiguous,
	kwNEVER:              kwCatUnambiguous,
	kwDAY:                kwCatUnambiguous,
	kwHISTORY:            kwCatUnambiguous,
	kwREUSE:              kwCatUnambiguous,
	kwOPTIONAL:           kwCatUnambiguous,
	kwX509:               kwCatUnambiguous,
	kwISSUER:             kwCatUnambiguous,
	kwSUBJECT:            kwCatUnambiguous,
	kwCIPHER:             kwCatUnambiguous,
	kwSCHEDULE:           kwCatUnambiguous,
	kwCOMPLETION:         kwCatUnambiguous,
	kwPRESERVE:           kwCatUnambiguous,
	kwEVERY:              kwCatUnambiguous,
	kwSTARTS:             kwCatUnambiguous,
	kwENDS:               kwCatUnambiguous,
	kwVALUE:              kwCatUnambiguous,
	kwSTACKED:            kwCatUnambiguous,
	kwUNKNOWN:            kwCatUnambiguous,
	kwWAIT:               kwCatUnambiguous,
	kwACTIVE:             kwCatUnambiguous,
	kwINACTIVE:           kwCatUnambiguous,
	kwATTRIBUTE:          kwCatUnambiguous,
	kwADMIN:              kwCatUnambiguous,
	kwDESCRIPTION:        kwCatUnambiguous,
	kwORGANIZATION:       kwCatUnambiguous,
	kwREFERENCE:          kwCatUnambiguous,
	kwDEFINITION:         kwCatUnambiguous,
	kwNAME:               kwCatUnambiguous,
	kwSYSTEM:             kwCatUnambiguous,
	kwROTATE:             kwCatUnambiguous,
	kwKEYRING:            kwCatUnambiguous,
	kwTLS:                kwCatUnambiguous,
	kwSTREAM:             kwCatUnambiguous,
	kwGENERATE:           kwCatUnambiguous,
	// Section 1.5: Ambiguous category classifications
	// Ambiguous 1 (not label, not role)
	kwEXECUTE:     kwCatAmbiguous1,
	// Ambiguous 2 (not label)
	kwBEGIN:       kwCatAmbiguous2,
	kwCOMMIT:      kwCatAmbiguous2,
	kwCONTAINS:    kwCatAmbiguous2,
	kwDO:          kwCatAmbiguous2,
	kwFLUSH:       kwCatAmbiguous2,
	kwFOLLOWS:     kwCatAmbiguous2,
	kwPRECEDES:    kwCatAmbiguous2,
	kwPREPARE:     kwCatAmbiguous2,
	kwREPAIR:      kwCatAmbiguous2,
	kwRESET:       kwCatAmbiguous2,
	kwROLLBACK:    kwCatAmbiguous2,
	kwSAVEPOINT:   kwCatAmbiguous2,
	kwSLAVE:       kwCatAmbiguous2,
	kwSTART:       kwCatAmbiguous2,
	kwSTOP:        kwCatAmbiguous2,
	kwTRUNCATE:    kwCatAmbiguous2,
	kwXA:          kwCatAmbiguous2,
	// Ambiguous 3 (not role)
	kwEVENT:       kwCatAmbiguous3,
	kwPROCESS:     kwCatAmbiguous3,
	kwRELOAD:      kwCatAmbiguous3,
	kwREPLICATION: kwCatAmbiguous3,
	// Ambiguous 4 (not lvalue)
	kwGLOBAL:      kwCatAmbiguous4,
	kwSESSION:     kwCatAmbiguous4,
	kwLOCAL:       kwCatAmbiguous4,
}

// isReserved returns true if the token type is a reserved keyword that cannot
// be used as an unquoted identifier.
func isReserved(t int) bool {
	cat, ok := keywordCategories[t]
	return ok && cat == kwCatReserved
}

// isIdentKeyword returns true if the token type is a non-reserved keyword that
// can be used as an identifier. This covers all 5 non-reserved categories:
// unambiguous, ambiguous_1, ambiguous_2, ambiguous_3, and ambiguous_4.
func isIdentKeyword(t int) bool {
	cat, ok := keywordCategories[t]
	return ok && cat != kwCatReserved
}

// isLabelKeyword returns true if the token type is a non-reserved keyword that
// can be used as a statement label. Includes: unambiguous, ambiguous_3, ambiguous_4.
// Excludes: ambiguous_1 (not label, not role), ambiguous_2 (not label).
func isLabelKeyword(t int) bool {
	cat, ok := keywordCategories[t]
	if !ok {
		return false
	}
	return cat == kwCatUnambiguous || cat == kwCatAmbiguous3 || cat == kwCatAmbiguous4
}

// isRoleKeyword returns true if the token type is a non-reserved keyword that
// can be used as a role name. Includes: unambiguous, ambiguous_2, ambiguous_4.
// Excludes: ambiguous_1 (not label, not role), ambiguous_3 (not role).
func isRoleKeyword(t int) bool {
	cat, ok := keywordCategories[t]
	if !ok {
		return false
	}
	return cat == kwCatUnambiguous || cat == kwCatAmbiguous2 || cat == kwCatAmbiguous4
}

// isLvalueKeyword returns true if the token type is a non-reserved keyword that
// can be used as an lvalue (SET target). Includes: unambiguous, ambiguous_1, ambiguous_2, ambiguous_3.
// Excludes: ambiguous_4 (system variables like GLOBAL, SESSION, LOCAL).
func isLvalueKeyword(t int) bool {
	cat, ok := keywordCategories[t]
	if !ok {
		return false
	}
	return cat == kwCatUnambiguous || cat == kwCatAmbiguous1 || cat == kwCatAmbiguous2 || cat == kwCatAmbiguous3
}

// parseIdent parses an identifier matching MySQL's `ident` grammar rule.
// Accepts tokIDENT plus all 5 non-reserved keyword categories (unambiguous,
// ambiguous_1, ambiguous_2, ambiguous_3, ambiguous_4).
//
// Ref: mysql-server sql/sql_yacc.yy — ident rule
// Ref: https://dev.mysql.com/doc/refman/8.0/en/identifiers.html
func (p *Parser) parseIdent() (string, int, error) {
	if p.cur.Type == tokIDENT {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	// Accept non-reserved keyword tokens as identifiers.
	if p.cur.Type >= 700 && !isReserved(p.cur.Type) {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type == tokEOF {
		return "", 0, p.syntaxErrorAtCur()
	}
	return "", 0, &ParseError{
		Message:  "expected identifier",
		Position: p.cur.Loc,
	}
}

// parseIdentifier is a thin alias for parseIdent, preserved for gradual migration
// of existing call sites. New code should use parseIdent directly.
func (p *Parser) parseIdentifier() (string, int, error) {
	return p.parseIdent()
}

// parseLabelIdent parses an identifier matching MySQL's `label_ident` grammar rule.
// Accepts tokIDENT plus unambiguous, ambiguous_3, and ambiguous_4 keywords.
// Excludes ambiguous_1 (not label, not role) and ambiguous_2 (not label).
//
// Ref: mysql-server sql/sql_yacc.yy — label_ident rule
func (p *Parser) parseLabelIdent() (string, int, error) {
	if p.cur.Type == tokIDENT {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type >= 700 && isLabelKeyword(p.cur.Type) {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type == tokEOF {
		return "", 0, p.syntaxErrorAtCur()
	}
	return "", 0, &ParseError{
		Message:  "expected identifier",
		Position: p.cur.Loc,
	}
}

// parseRoleIdent parses an identifier matching MySQL's `role_ident` grammar rule.
// Accepts tokIDENT plus unambiguous, ambiguous_2, and ambiguous_4 keywords.
// Excludes ambiguous_1 (not label, not role) and ambiguous_3 (not role).
//
// Ref: mysql-server sql/sql_yacc.yy — role_ident rule
func (p *Parser) parseRoleIdent() (string, int, error) {
	if p.cur.Type == tokIDENT {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type >= 700 && isRoleKeyword(p.cur.Type) {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type == tokEOF {
		return "", 0, p.syntaxErrorAtCur()
	}
	return "", 0, &ParseError{
		Message:  "expected identifier",
		Position: p.cur.Loc,
	}
}

// parseLvalueIdent parses an identifier matching MySQL's `lvalue_ident` grammar rule.
// Accepts tokIDENT plus unambiguous, ambiguous_1, ambiguous_2, and ambiguous_3 keywords.
// Excludes ambiguous_4 (system variables like GLOBAL, SESSION, LOCAL).
//
// Ref: mysql-server sql/sql_yacc.yy — lvalue_ident rule
func (p *Parser) parseLvalueIdent() (string, int, error) {
	if p.cur.Type == tokIDENT {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type >= 700 && isLvalueKeyword(p.cur.Type) {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type == tokEOF {
		return "", 0, p.syntaxErrorAtCur()
	}
	return "", 0, &ParseError{
		Message:  "expected identifier",
		Position: p.cur.Loc,
	}
}

// parseKeywordOrIdent parses an identifier or ANY keyword token (including reserved words).
// Use this in contexts where the grammar expects a fixed enum value, action word, or
// option name that may collide with reserved keywords. Examples:
//   - ALGORITHM = DEFAULT/INSTANT/INPLACE/COPY
//   - LOCK = NONE/SHARED/EXCLUSIVE
//   - REQUIRE_TABLE_PRIMARY_KEY_CHECK = ON/OFF/STREAM/GENERATE
//   - ALTER INSTANCE action words: ROTATE/RELOAD/ENABLE/DISABLE
//   - INTERVAL units: DAY/HOUR/MINUTE/SECOND (and compound forms)
//   - EXTRACT units, EXPLAIN FORMAT values, etc.
//
// This matches MySQL's grammar behavior where specific productions explicitly
// list keyword tokens as valid alternatives (e.g., ON_SYM in option values).
func (p *Parser) parseKeywordOrIdent() (string, int, error) {
	if p.cur.Type == tokIDENT {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	// Accept ANY keyword token, including reserved words.
	if p.cur.Type >= 700 {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type == tokEOF {
		return "", 0, p.syntaxErrorAtCur()
	}
	return "", 0, &ParseError{
		Message:  "expected identifier or keyword",
		Position: p.cur.Loc,
	}
}

// parseColumnRef parses a column reference, which may be qualified:
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/identifier-qualifiers.html
//
//	column_ref:
//	    identifier
//	    | identifier '.' identifier
//	    | identifier '.' identifier '.' identifier
//	    | identifier '.' '*'
//	    | identifier '.' identifier '.' '*'
func (p *Parser) parseColumnRef() (*nodes.ColumnRef, error) {
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	ref := &nodes.ColumnRef{
		Loc:    nodes.Loc{Start: start},
		Column: name,
	}

	// Check for dot-qualification
	if p.cur.Type == '.' {
		p.advance() // consume '.'

		// Check for table.* or schema.table.*
		if p.cur.Type == '*' {
			p.advance()
			ref.Table = name
			ref.Column = ""
			ref.Star = true
			ref.Loc.End = p.pos()
			return ref, nil
		}

		name2, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}

		// Check for second dot: schema.table.col or schema.table.*
		if p.cur.Type == '.' {
			p.advance() // consume second '.'

			if p.cur.Type == '*' {
				p.advance()
				ref.Schema = name
				ref.Table = name2
				ref.Column = ""
				ref.Star = true
				ref.Loc.End = p.pos()
				return ref, nil
			}

			name3, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			ref.Schema = name
			ref.Table = name2
			ref.Column = name3
			ref.Loc.End = p.pos()
			return ref, nil
		}

		// table.col
		ref.Table = name
		ref.Column = name2
	}

	ref.Loc.End = p.pos()
	return ref, nil
}

// parseTableRef parses a table reference (possibly qualified with schema).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/identifier-qualifiers.html
//
//	table_ref:
//	    identifier
//	    | identifier '.' identifier
func (p *Parser) parseTableRef() (*nodes.TableRef, error) {
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	ref := &nodes.TableRef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// Check for schema.table
	if p.cur.Type == '.' {
		p.advance() // consume '.'

		// Completion: after "db.", offer table_ref qualified with database.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
		}

		name2, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		ref.Schema = name
		ref.Name = name2
	}

	ref.Loc.End = p.pos()
	return ref, nil
}

// parseTableRefWithAlias parses a table reference with an optional alias.
//
//	table_ref_alias:
//	    table_ref [AS identifier | identifier]
func (p *Parser) parseTableRefWithAlias() (*nodes.TableRef, error) {
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}

	// Optional PARTITION (p0, p1, ...)
	if p.cur.Type == kwPARTITION {
		p.advance()
		parts, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		ref.Partitions = parts
		ref.Loc.End = p.pos()
	}

	// Optional AS alias
	if _, ok := p.match(kwAS); ok {
		alias, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		ref.Alias = alias
		ref.Loc.End = p.pos()
	} else if p.cur.Type == tokIDENT {
		// Alias without AS keyword
		alias, _, _ := p.parseIdentifier()
		ref.Alias = alias
		ref.Loc.End = p.pos()
	}

	// Optional index hints: USE/FORCE/IGNORE {INDEX|KEY} ...
	if p.cur.Type == kwUSE || p.cur.Type == kwFORCE || p.cur.Type == kwIGNORE {
		hints, err := p.parseIndexHints()
		if err != nil {
			return nil, err
		}
		ref.IndexHints = hints
		ref.Loc.End = p.pos()
	}

	return ref, nil
}

// parseVariableRef parses a user variable (@var) or system variable (@@var).
// The lexer emits these as tokIDENT with "@" or "@@" prefix in the Str.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/user-variables.html
// Ref: https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html
//
//	variable_ref:
//	    '@' identifier
//	    | '@@' [GLOBAL | SESSION | LOCAL] '.' identifier
func (p *Parser) parseVariableRef() (*nodes.VariableRef, error) {
	if p.cur.Type != tokIDENT {
		return nil, &ParseError{
			Message:  "expected variable reference",
			Position: p.cur.Loc,
		}
	}

	tok := p.cur
	str := tok.Str

	// Check for @@ prefix (system variable)
	if len(str) > 2 && str[0] == '@' && str[1] == '@' {
		p.advance()
		name := str[2:]
		ref := &nodes.VariableRef{
			Loc:    nodes.Loc{Start: tok.Loc},
			System: true,
		}

		// Check for scope prefix: @@global.var, @@session.var, @@local.var
		if dotIdx := indexOf(name, '.'); dotIdx >= 0 {
			scope := name[:dotIdx]
			varName := name[dotIdx+1:]
			switch {
			case eqFold(scope, "global"):
				ref.Scope = "GLOBAL"
			case eqFold(scope, "session"):
				ref.Scope = "SESSION"
			case eqFold(scope, "local"):
				ref.Scope = "LOCAL"
			default:
				// Not a scope, treat as qualified name
				ref.Name = name
				ref.Loc.End = p.pos()
				return ref, nil
			}
			ref.Name = varName
		} else {
			ref.Name = name
		}

		ref.Loc.End = p.pos()
		return ref, nil
	}

	// Check for @ prefix (user variable)
	if len(str) > 1 && str[0] == '@' {
		p.advance()
		ref := &nodes.VariableRef{
			Loc:  nodes.Loc{Start: tok.Loc},
			Name: str[1:],
		}
		ref.Loc.End = p.pos()
		return ref, nil
	}

	return nil, &ParseError{
		Message:  "expected variable reference",
		Position: p.cur.Loc,
	}
}

// isIdentToken returns true if the current token can be used as an identifier.
func (p *Parser) isIdentToken() bool {
	return p.cur.Type == tokIDENT || (p.cur.Type >= 700 && !isReserved(p.cur.Type))
}

// isLabelIdentToken returns true if the current token can be used as a label identifier.
// This matches MySQL's label_ident rule: tokIDENT + unambiguous + ambiguous_3 + ambiguous_4.
// Excludes ambiguous_1 and ambiguous_2 keywords (e.g., BEGIN, COMMIT, END).
func (p *Parser) isLabelIdentToken() bool {
	return p.cur.Type == tokIDENT || (p.cur.Type >= 700 && isLabelKeyword(p.cur.Type))
}

// isLvalueIdentToken returns true if the current token can be used as an lvalue identifier.
// This matches MySQL's lvalue_ident rule: tokIDENT + unambiguous + ambiguous_1 + ambiguous_2 + ambiguous_3.
// Excludes ambiguous_4 keywords (GLOBAL, SESSION, LOCAL).
func (p *Parser) isLvalueIdentToken() bool {
	return p.cur.Type == tokIDENT || (p.cur.Type >= 700 && isLvalueKeyword(p.cur.Type))
}

// isVariableRef returns true if the current token is a variable reference.
func (p *Parser) isVariableRef() bool {
	if p.cur.Type != tokIDENT {
		return false
	}
	return len(p.cur.Str) > 0 && p.cur.Str[0] == '@'
}

// indexOf returns the index of the first occurrence of ch in s, or -1.
func indexOf(s string, ch byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ch {
			return i
		}
	}
	return -1
}

// eqFold reports whether s and t are equal under Unicode case-folding (ASCII only).
func eqFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		a, b := s[i], t[i]
		if a >= 'A' && a <= 'Z' {
			a += 'a' - 'A'
		}
		if b >= 'A' && b <= 'Z' {
			b += 'a' - 'A'
		}
		if a != b {
			return false
		}
	}
	return true
}
