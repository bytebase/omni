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
		kwOPTION, kwOR, kwORDER, kwPCTFREE, kwPRIOR, kwPRIVILEGES,
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
