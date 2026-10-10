package parser

type oracleKeywordCategory int

const (
	oracleKeywordIdentifier oracleKeywordCategory = iota
	oracleKeywordReserved
	oracleKeywordNonReserved
	oracleKeywordContext
	oracleKeywordType
	oracleKeywordFunction
	oracleKeywordPseudoColumn
	oracleKeywordClauseStarter
)

func (c oracleKeywordCategory) String() string {
	switch c {
	case oracleKeywordIdentifier:
		return "identifier"
	case oracleKeywordReserved:
		return "reserved"
	case oracleKeywordNonReserved:
		return "nonreserved"
	case oracleKeywordContext:
		return "context"
	case oracleKeywordType:
		return "type"
	case oracleKeywordFunction:
		return "function"
	case oracleKeywordPseudoColumn:
		return "pseudo-column"
	case oracleKeywordClauseStarter:
		return "clause-starter"
	default:
		return "unknown"
	}
}

var oracleReservedIdentifierTexts = map[string]struct{}{
	"MLSLABEL": {},
	"RESOURCE": {},
	"UID":      {},
}

// oracleKeywordCategoryOf classifies lexer output for keyword audit coverage.
func oracleKeywordCategoryOf(tok Token) oracleKeywordCategory {
	switch {
	case tok.Type == tokQIDENT:
		return oracleKeywordIdentifier
	case isOraclePseudoColumnKeyword(tok.Type):
		return oracleKeywordPseudoColumn
	case isOracleTypeKeyword(tok.Type):
		return oracleKeywordType
	case isOracleFunctionKeyword(tok.Type):
		return oracleKeywordFunction
	case isOracleClauseStarterKeyword(tok.Type):
		return oracleKeywordClauseStarter
	case isOracleContextKeyword(tok.Type):
		return oracleKeywordContext
	case isOracleSQLReservedKeyword(tok):
		return oracleKeywordReserved
	case tok.Type == tokIDENT:
		return oracleKeywordIdentifier
	case tok.Type >= 2000:
		return oracleKeywordNonReserved
	default:
		return oracleKeywordIdentifier
	}
}

func isOraclePseudoColumnKeyword(tokenType int) bool {
	switch tokenType {
	case kwROWID, kwROWNUM, kwLEVEL, kwSYSDATE, kwSYSTIMESTAMP, kwUSER:
		return true
	default:
		return false
	}
}

func isOracleTypeKeyword(tokenType int) bool {
	switch tokenType {
	case kwBLOB, kwCHAR, kwCLOB, kwDATE, kwDECIMAL, kwFLOAT, kwINTEGER,
		kwINTERVAL, kwLONG, kwNCHAR, kwNCLOB, kwNUMBER, kwNVARCHAR2,
		kwRAW, kwSMALLINT, kwTIMESTAMP, kwVARCHAR, kwVARCHAR2, kwVARRAY:
		return true
	default:
		return false
	}
}

func isOracleFunctionKeyword(tokenType int) bool {
	switch tokenType {
	case kwCAST, kwCOLLECT, kwDECODE, kwDENSE_RANK, kwJSON_ARRAY,
		kwJSON_EXISTS, kwJSON_MERGEPATCH, kwJSON_OBJECT, kwJSON_QUERY,
		kwJSON_VALUE, kwSYS_CONNECT_BY_PATH, kwTREAT, kwXMLAGG,
		kwXMLELEMENT, kwXMLFOREST, kwXMLPARSE, kwXMLROOT, kwXMLSERIALIZE:
		return true
	default:
		return false
	}
}

func isOracleClauseStarterKeyword(tokenType int) bool {
	switch tokenType {
	case kwCONNECT, kwFETCH, kwFROM, kwGROUP, kwHAVING, kwINTERSECT, kwJOIN,
		kwMINUS, kwMODEL, kwOFFSET, kwON, kwORDER, kwSTART, kwUNION,
		kwUSING, kwWHERE, kwWITH:
		return true
	default:
		return false
	}
}

func isOracleContextKeyword(tokenType int) bool {
	switch tokenType {
	case kwALWAYS, kwAUTOMATIC, kwBLOCK, kwCOLUMNS, kwCONTENT, kwCUBE,
		kwDECREMENT, kwDIMENSION, kwFOLLOWING, kwFORMAT, kwGENERATED,
		kwGROUPING, kwGROUPS, kwIDENTITY, kwJSON, kwJSON_TABLE, kwLATERAL,
		kwMAIN, kwMEASURES, kwNAV, kwNESTED, kwORDINALITY, kwPASSING,
		kwPATH, kwPRECEDING, kwREFERENCE, kwROLLUP, kwRULES, kwSEED,
		kwSEQUENTIAL, kwSETS, kwUNBOUNDED, kwUNTIL, kwUPDATED, kwUPSERT,
		kwVERSIONS, kwWITHIN, kwXMLTABLE:
		return true
	default:
		return false
	}
}

// plsqlReservedWords are the words Oracle 23ai rejects, unquoted, where
// PL/SQL expects a name, such as the label after END (PLS-00103). They were
// found by compiling BEGIN BEGIN NULL; END word; END; for each word of the
// documented PL/SQL reserved words, each SQL reserved word, and each PL/SQL
// statement keyword; the engine, not the list, decides.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/plsql-reserved-words-keywords.html
var plsqlReservedWords = map[string]bool{
	"ALL": true, "ALTER": true, "AND": true, "ANY": true, "AS": true, "ASC": true, "AT": true,
	"BEGIN": true, "BETWEEN": true, "BY": true, "CASE": true, "CHECK": true, "CLUSTER": true,
	"CLUSTERS": true, "COLAUTH": true, "COLUMNS": true, "COMPRESS": true, "CONNECT": true,
	"CRASH": true, "CREATE": true, "CURRENT": true, "DECLARE": true, "DEFAULT": true,
	"DELETE": true, "DESC": true, "DISTINCT": true, "DROP": true, "ELSE": true, "END": true,
	"EXCEPTION": true, "EXCLUSIVE": true, "EXISTS": true, "FETCH": true, "FOR": true,
	"FROM": true, "GOTO": true, "GRANT": true, "GROUP": true, "HAVING": true,
	"IDENTIFIED": true, "IF": true, "IN": true, "INDEX": true, "INDEXES": true, "INSERT": true,
	"INTERSECT": true, "INTO": true, "IS": true, "LIKE": true, "LOCK": true, "MINUS": true,
	"MODE": true, "NOCOMPRESS": true, "NOT": true, "NOWAIT": true, "NULL": true, "OF": true,
	"ON": true, "OPTION": true, "OR": true, "ORDER": true, "OVERLAPS": true, "PRIOR": true,
	"PROCEDURE": true, "PUBLIC": true, "RESOURCE": true, "REVOKE": true, "SELECT": true,
	"SHARE": true, "SIZE": true, "SQL": true, "START": true, "TABAUTH": true, "TABLE": true,
	"THEN": true, "TO": true, "UNION": true, "UNIQUE": true, "UPDATE": true, "VALUES": true,
	"VIEW": true, "VIEWS": true, "WHEN": true, "WHERE": true, "WITH": true,
}

// isPLSQLIdentifier reports whether the current token can name a PL/SQL
// item: a quoted identifier, or a word that is not PL/SQL reserved. SUBTYPE
// loop IS NUMBER compiles on Oracle 23ai; SUBTYPE prior IS NUMBER does not.
func (p *Parser) isPLSQLIdentifier() bool {
	if p.cur.Type == tokQIDENT {
		return true
	}
	return p.isIdentLike() && !plsqlReservedWords[p.cur.Str]
}

// atEndName reports whether the current token, just after END or END LOOP,
// is the optional label or unit name. A PL/SQL reserved word is not: in
// END LOOP END; the second END closes the enclosing block, and Oracle reports
// the missing ';' at it.
func (p *Parser) atEndName() bool {
	return p.isPLSQLIdentifier()
}

// isOracleSQLReservedKeyword returns true when tok is an Oracle SQL reserved word
// that cannot be used as a nonquoted object or column identifier.
func isOracleSQLReservedKeyword(tok Token) bool {
	switch tok.Type {
	case kwACCESS, kwADD, kwALL, kwALTER, kwAND, kwANY, kwAS, kwASC,
		kwAUDIT, kwBETWEEN, kwBY, kwCHAR, kwCHECK, kwCLUSTER, kwCOLUMN, kwCOLUMN_VALUE,
		kwCOMMENT, kwCOMPRESS, kwCONNECT, kwCREATE, kwCURRENT, kwDATE,
		kwDECIMAL, kwDEFAULT, kwDELETE, kwDESC, kwDISTINCT, kwDROP,
		kwELSE, kwEXCLUSIVE, kwEXISTS, kwFILE, kwFLOAT, kwFOR, kwFROM,
		kwGRANT, kwGROUP, kwHAVING, kwIDENTIFIED, kwIMMEDIATE, kwIN,
		kwINCREMENT, kwINDEX, kwINITIAL, kwINSERT, kwINTEGER, kwINTERSECT,
		kwINTO, kwIS, kwLEVEL, kwLIKE, kwLOCK, kwLONG, kwMAXEXTENTS,
		kwMINUS, kwMLSLABEL, kwMODE, kwMODIFY, kwNESTED_TABLE_ID, kwNOAUDIT, kwNOCOMPRESS, kwNOT,
		kwNOWAIT, kwNULL, kwNUMBER, kwOF, kwOFFLINE, kwON, kwONLINE,
		kwOPTION, kwOR, kwORDER, kwPCTFREE, kwPRIOR,
		kwPUBLIC, kwRAW, kwRENAME, kwRESOURCE, kwREVOKE, kwROW, kwROWID, kwROWNUM,
		kwROWS, kwSELECT, kwSESSION, kwSET, kwSHARE, kwSIZE, kwSMALLINT,
		kwSTART, kwSUCCESSFUL, kwSYNONYM, kwSYSDATE, kwTABLE, kwTHEN,
		kwTO, kwTRIGGER, kwUID, kwUNION, kwUNIQUE, kwUPDATE, kwUSER, kwVALIDATE,
		kwVALUES, kwVARCHAR, kwVARCHAR2, kwVIEW, kwWHENEVER, kwWHERE, kwWITH:
		return true
	default:
		_, ok := oracleReservedIdentifierTexts[tok.Str]
		return ok
	}
}

func (p *Parser) syntaxErrorIfReservedIdentifier() error {
	if isOracleSQLReservedKeyword(p.cur) {
		return p.syntaxErrorAtCur()
	}
	return nil
}

func (p *Parser) syntaxErrorIfReservedColumnIdentifier() error {
	if p.cur.Type == tokQIDENT && p.cur.Str == "ROWID" {
		return p.syntaxErrorAtCur()
	}
	return p.syntaxErrorIfReservedIdentifier()
}
