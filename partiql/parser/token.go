// Package parser implements the PartiQL parser. This file declares the
// token type used by the lexer (lexer.go) and the full set of token
// type constants. Token positions use ast.Loc directly to eliminate the
// need for the parser to convert at every AST node construction site.
package parser

import "github.com/bytebase/omni/partiql/ast"

// Token is a single PartiQL lexer token.
//
// Type is one of the tok* constants below.
//
// Str holds the raw source text for most tokens. For tokSCONST (single-quoted
// string literals) and tokIDENT_QUOTED (double-quoted identifiers), Str holds
// the *decoded* value with the doubled-quote escape collapsed: two consecutive
// quote characters inside the literal represent a single quote in the decoded
// value (e.g., the SCONST source spelled i, t, quote, quote, s with surrounding
// single quotes decodes to the Go string "it's"). For tokION_LITERAL, Str is
// the verbatim inner content between the backticks (no decoding).
//
// Loc is the half-open byte range covering the token in the source string.
type Token struct {
	Type int
	Str  string
	Loc  ast.Loc
}

// ===========================================================================
// Special tokens.
// ===========================================================================

const (
	tokEOF     = 0 // end of input or after lex error
	tokInvalid = 1 // sentinel for unknown token type (never returned by Next)
)

// ===========================================================================
// Literal tokens — group 1000.
// ===========================================================================

const (
	tokSCONST       = iota + 1000 // single-quoted string literal: 'hello'
	tokICONST                     // integer literal: 42
	tokFCONST                     // decimal/float literal: 3.14, 1e10, .5
	tokIDENT                      // unquoted identifier (case-insensitive lookup)
	tokIDENT_QUOTED               // double-quoted identifier (case-sensitive): "Foo"
	tokION_LITERAL                // backtick-delimited Ion blob (body deferred to DAG node 17)
)

// ===========================================================================
// Operator and punctuation tokens — group 2000.
//
// Names generally follow PartiQLLexer.g4 rule names verbatim for traceability
// against the grammar; tokLT and tokGT substitute for the grammar's
// ANGLE_LEFT and ANGLE_RIGHT for readability and SQL-parser idiom.
//
// Ordering is semantic (arithmetic, then comparison, then concat, then the
// PartiQL bag-bracket pair, then paired delimiters, then punctuation) rather
// than the grammar's lexer-rule declaration order.
// ===========================================================================

const (
	tokPLUS               = iota + 2000 // +
	tokMINUS                            // -
	tokASTERISK                         // *
	tokSLASH_FORWARD                    // /
	tokPERCENT                          // %
	tokCARET                            // ^
	tokTILDE                            // ~
	tokAT_SIGN                          // @
	tokEQ                               // =
	tokNEQ                              // <> or !=
	tokLT                               // < (ANGLE_LEFT in grammar)
	tokGT                               // > (ANGLE_RIGHT in grammar)
	tokLT_EQ                            // <=
	tokGT_EQ                            // >=
	tokCONCAT                           // ||
	tokANGLE_DOUBLE_LEFT                // <<  (PartiQL bag-literal start)
	tokANGLE_DOUBLE_RIGHT               // >>  (PartiQL bag-literal end)
	tokPAREN_LEFT                       // (
	tokPAREN_RIGHT                      // )
	tokBRACKET_LEFT                     // [
	tokBRACKET_RIGHT                    // ]
	tokBRACE_LEFT                       // {
	tokBRACE_RIGHT                      // }
	tokCOLON                            // :
	tokCOLON_SEMI                       // ;
	tokCOMMA                            // ,
	tokPERIOD                           // .
	tokQUESTION_MARK                    // ?
)

// tokenName returns the canonical printable name for a token type constant.
// Used by error messages, test failure output, and future debugging.
//
// Covers all 302 tok* constants in the package.
func tokenName(t int) string {
	switch t {
	// Specials.
	case tokEOF:
		return "EOF"
	case tokInvalid:
		return "INVALID"

	// Literals.
	case tokSCONST:
		return "SCONST"
	case tokICONST:
		return "ICONST"
	case tokFCONST:
		return "FCONST"
	case tokIDENT:
		return "IDENT"
	case tokIDENT_QUOTED:
		return "IDENT_QUOTED"
	case tokION_LITERAL:
		return "ION_LITERAL"

	// Operators / punctuation.
	case tokPLUS:
		return "PLUS"
	case tokMINUS:
		return "MINUS"
	case tokASTERISK:
		return "ASTERISK"
	case tokSLASH_FORWARD:
		return "SLASH_FORWARD"
	case tokPERCENT:
		return "PERCENT"
	case tokCARET:
		return "CARET"
	case tokTILDE:
		return "TILDE"
	case tokAT_SIGN:
		return "AT_SIGN"
	case tokEQ:
		return "EQ"
	case tokNEQ:
		return "NEQ"
	case tokLT:
		return "LT"
	case tokGT:
		return "GT"
	case tokLT_EQ:
		return "LT_EQ"
	case tokGT_EQ:
		return "GT_EQ"
	case tokCONCAT:
		return "CONCAT"
	case tokANGLE_DOUBLE_LEFT:
		return "ANGLE_DOUBLE_LEFT"
	case tokANGLE_DOUBLE_RIGHT:
		return "ANGLE_DOUBLE_RIGHT"
	case tokPAREN_LEFT:
		return "PAREN_LEFT"
	case tokPAREN_RIGHT:
		return "PAREN_RIGHT"
	case tokBRACKET_LEFT:
		return "BRACKET_LEFT"
	case tokBRACKET_RIGHT:
		return "BRACKET_RIGHT"
	case tokBRACE_LEFT:
		return "BRACE_LEFT"
	case tokBRACE_RIGHT:
		return "BRACE_RIGHT"
	case tokCOLON:
		return "COLON"
	case tokCOLON_SEMI:
		return "COLON_SEMI"
	case tokCOMMA:
		return "COMMA"
	case tokPERIOD:
		return "PERIOD"
	case tokQUESTION_MARK:
		return "QUESTION_MARK"

	// Keywords (alphabetical).
	case tokABSOLUTE:
		return "ABSOLUTE"
	case tokACTION:
		return "ACTION"
	case tokADD:
		return "ADD"
	case tokALL:
		return "ALL"
	case tokALLOCATE:
		return "ALLOCATE"
	case tokALTER:
		return "ALTER"
	case tokAND:
		return "AND"
	case tokANY:
		return "ANY"
	case tokARE:
		return "ARE"
	case tokAS:
		return "AS"
	case tokASC:
		return "ASC"
	case tokASSERTION:
		return "ASSERTION"
	case tokAT:
		return "AT"
	case tokAUTHORIZATION:
		return "AUTHORIZATION"
	case tokAVG:
		return "AVG"
	case tokBAG:
		return "BAG"
	case tokBEGIN:
		return "BEGIN"
	case tokBETWEEN:
		return "BETWEEN"
	case tokBIGINT:
		return "BIGINT"
	case tokBIT:
		return "BIT"
	case tokBIT_LENGTH:
		return "BIT_LENGTH"
	case tokBLOB:
		return "BLOB"
	case tokBOOL:
		return "BOOL"
	case tokBOOLEAN:
		return "BOOLEAN"
	case tokBY:
		return "BY"
	case tokCAN_CAST:
		return "CAN_CAST"
	case tokCAN_LOSSLESS_CAST:
		return "CAN_LOSSLESS_CAST"
	case tokCASCADE:
		return "CASCADE"
	case tokCASCADED:
		return "CASCADED"
	case tokCASE:
		return "CASE"
	case tokCAST:
		return "CAST"
	case tokCATALOG:
		return "CATALOG"
	case tokCHAR:
		return "CHAR"
	case tokCHAR_LENGTH:
		return "CHAR_LENGTH"
	case tokCHARACTER:
		return "CHARACTER"
	case tokCHARACTER_LENGTH:
		return "CHARACTER_LENGTH"
	case tokCHECK:
		return "CHECK"
	case tokCLOB:
		return "CLOB"
	case tokCLOSE:
		return "CLOSE"
	case tokCOALESCE:
		return "COALESCE"
	case tokCOLLATE:
		return "COLLATE"
	case tokCOLLATION:
		return "COLLATION"
	case tokCOLUMN:
		return "COLUMN"
	case tokCOMMIT:
		return "COMMIT"
	case tokCONFLICT:
		return "CONFLICT"
	case tokCONNECT:
		return "CONNECT"
	case tokCONNECTION:
		return "CONNECTION"
	case tokCONSTRAINT:
		return "CONSTRAINT"
	case tokCONSTRAINTS:
		return "CONSTRAINTS"
	case tokCONTINUE:
		return "CONTINUE"
	case tokCONVERT:
		return "CONVERT"
	case tokCORRESPONDING:
		return "CORRESPONDING"
	case tokCOUNT:
		return "COUNT"
	case tokCREATE:
		return "CREATE"
	case tokCROSS:
		return "CROSS"
	case tokCURRENT:
		return "CURRENT"
	case tokCURRENT_DATE:
		return "CURRENT_DATE"
	case tokCURRENT_TIME:
		return "CURRENT_TIME"
	case tokCURRENT_TIMESTAMP:
		return "CURRENT_TIMESTAMP"
	case tokCURRENT_USER:
		return "CURRENT_USER"
	case tokCURSOR:
		return "CURSOR"
	case tokDATE:
		return "DATE"
	case tokDATE_ADD:
		return "DATE_ADD"
	case tokDATE_DIFF:
		return "DATE_DIFF"
	case tokDEALLOCATE:
		return "DEALLOCATE"
	case tokDEC:
		return "DEC"
	case tokDECIMAL:
		return "DECIMAL"
	case tokDECLARE:
		return "DECLARE"
	case tokDEFAULT:
		return "DEFAULT"
	case tokDEFERRABLE:
		return "DEFERRABLE"
	case tokDEFERRED:
		return "DEFERRED"
	case tokDELETE:
		return "DELETE"
	case tokDESC:
		return "DESC"
	case tokDESCRIBE:
		return "DESCRIBE"
	case tokDESCRIPTOR:
		return "DESCRIPTOR"
	case tokDIAGNOSTICS:
		return "DIAGNOSTICS"
	case tokDISCONNECT:
		return "DISCONNECT"
	case tokDISTINCT:
		return "DISTINCT"
	case tokDO:
		return "DO"
	case tokDOMAIN:
		return "DOMAIN"
	case tokDOUBLE:
		return "DOUBLE"
	case tokDROP:
		return "DROP"
	case tokELSE:
		return "ELSE"
	case tokEND:
		return "END"
	case tokEND_EXEC:
		return "END_EXEC"
	case tokESCAPE:
		return "ESCAPE"
	case tokEXCEPT:
		return "EXCEPT"
	case tokEXCEPTION:
		return "EXCEPTION"
	case tokEXCLUDED:
		return "EXCLUDED"
	case tokEXEC:
		return "EXEC"
	case tokEXECUTE:
		return "EXECUTE"
	case tokEXISTS:
		return "EXISTS"
	case tokEXPLAIN:
		return "EXPLAIN"
	case tokEXTERNAL:
		return "EXTERNAL"
	case tokEXTRACT:
		return "EXTRACT"
	case tokFALSE:
		return "FALSE"
	case tokFETCH:
		return "FETCH"
	case tokFIRST:
		return "FIRST"
	case tokFLOAT:
		return "FLOAT"
	case tokFOR:
		return "FOR"
	case tokFOREIGN:
		return "FOREIGN"
	case tokFOUND:
		return "FOUND"
	case tokFROM:
		return "FROM"
	case tokFULL:
		return "FULL"
	case tokGET:
		return "GET"
	case tokGLOBAL:
		return "GLOBAL"
	case tokGO:
		return "GO"
	case tokGOTO:
		return "GOTO"
	case tokGRANT:
		return "GRANT"
	case tokGROUP:
		return "GROUP"
	case tokHAVING:
		return "HAVING"
	case tokIDENTITY:
		return "IDENTITY"
	case tokIMMEDIATE:
		return "IMMEDIATE"
	case tokIN:
		return "IN"
	case tokINDEX:
		return "INDEX"
	case tokINDICATOR:
		return "INDICATOR"
	case tokINITIALLY:
		return "INITIALLY"
	case tokINNER:
		return "INNER"
	case tokINPUT:
		return "INPUT"
	case tokINSENSITIVE:
		return "INSENSITIVE"
	case tokINSERT:
		return "INSERT"
	case tokINT:
		return "INT"
	case tokINT2:
		return "INT2"
	case tokINT4:
		return "INT4"
	case tokINT8:
		return "INT8"
	case tokINTEGER:
		return "INTEGER"
	case tokINTEGER2:
		return "INTEGER2"
	case tokINTEGER4:
		return "INTEGER4"
	case tokINTEGER8:
		return "INTEGER8"
	case tokINTERSECT:
		return "INTERSECT"
	case tokINTERVAL:
		return "INTERVAL"
	case tokINTO:
		return "INTO"
	case tokIS:
		return "IS"
	case tokISOLATION:
		return "ISOLATION"
	case tokJOIN:
		return "JOIN"
	case tokKEY:
		return "KEY"
	case tokLAG:
		return "LAG"
	case tokLANGUAGE:
		return "LANGUAGE"
	case tokLAST:
		return "LAST"
	case tokLATERAL:
		return "LATERAL"
	case tokLEAD:
		return "LEAD"
	case tokLEFT:
		return "LEFT"
	case tokLET:
		return "LET"
	case tokLEVEL:
		return "LEVEL"
	case tokLIKE:
		return "LIKE"
	case tokLIMIT:
		return "LIMIT"
	case tokLIST:
		return "LIST"
	case tokLOCAL:
		return "LOCAL"
	case tokLOWER:
		return "LOWER"
	case tokMATCH:
		return "MATCH"
	case tokMAX:
		return "MAX"
	case tokMIN:
		return "MIN"
	case tokMISSING:
		return "MISSING"
	case tokMODIFIED:
		return "MODIFIED"
	case tokMODULE:
		return "MODULE"
	case tokNAMES:
		return "NAMES"
	case tokNATIONAL:
		return "NATIONAL"
	case tokNATURAL:
		return "NATURAL"
	case tokNCHAR:
		return "NCHAR"
	case tokNEW:
		return "NEW"
	case tokNEXT:
		return "NEXT"
	case tokNO:
		return "NO"
	case tokNOT:
		return "NOT"
	case tokNOTHING:
		return "NOTHING"
	case tokNULL:
		return "NULL"
	case tokNULLIF:
		return "NULLIF"
	case tokNULLS:
		return "NULLS"
	case tokNUMERIC:
		return "NUMERIC"
	case tokOCTET_LENGTH:
		return "OCTET_LENGTH"
	case tokOF:
		return "OF"
	case tokOFFSET:
		return "OFFSET"
	case tokOLD:
		return "OLD"
	case tokON:
		return "ON"
	case tokONLY:
		return "ONLY"
	case tokOPEN:
		return "OPEN"
	case tokOPTION:
		return "OPTION"
	case tokOR:
		return "OR"
	case tokORDER:
		return "ORDER"
	case tokOUTER:
		return "OUTER"
	case tokOUTPUT:
		return "OUTPUT"
	case tokOVER:
		return "OVER"
	case tokOVERLAPS:
		return "OVERLAPS"
	case tokOVERLAY:
		return "OVERLAY"
	case tokPAD:
		return "PAD"
	case tokPARTIAL:
		return "PARTIAL"
	case tokPARTITION:
		return "PARTITION"
	case tokPIVOT:
		return "PIVOT"
	case tokPLACING:
		return "PLACING"
	case tokPOSITION:
		return "POSITION"
	case tokPRECISION:
		return "PRECISION"
	case tokPREPARE:
		return "PREPARE"
	case tokPRESERVE:
		return "PRESERVE"
	case tokPRIMARY:
		return "PRIMARY"
	case tokPRIOR:
		return "PRIOR"
	case tokPRIVILEGES:
		return "PRIVILEGES"
	case tokPROCEDURE:
		return "PROCEDURE"
	case tokPUBLIC:
		return "PUBLIC"
	case tokREAD:
		return "READ"
	case tokREAL:
		return "REAL"
	case tokREFERENCES:
		return "REFERENCES"
	case tokRELATIVE:
		return "RELATIVE"
	case tokREMOVE:
		return "REMOVE"
	case tokREPLACE:
		return "REPLACE"
	case tokRESTRICT:
		return "RESTRICT"
	case tokRETURNING:
		return "RETURNING"
	case tokREVOKE:
		return "REVOKE"
	case tokRIGHT:
		return "RIGHT"
	case tokROLLBACK:
		return "ROLLBACK"
	case tokROWS:
		return "ROWS"
	case tokSCHEMA:
		return "SCHEMA"
	case tokSCROLL:
		return "SCROLL"
	case tokSECTION:
		return "SECTION"
	case tokSELECT:
		return "SELECT"
	case tokSESSION:
		return "SESSION"
	case tokSESSION_USER:
		return "SESSION_USER"
	case tokSET:
		return "SET"
	case tokSEXP:
		return "SEXP"
	case tokSHORTEST:
		return "SHORTEST"
	case tokSIZE:
		return "SIZE"
	case tokSMALLINT:
		return "SMALLINT"
	case tokSOME:
		return "SOME"
	case tokSPACE:
		return "SPACE"
	case tokSQL:
		return "SQL"
	case tokSQLCODE:
		return "SQLCODE"
	case tokSQLERROR:
		return "SQLERROR"
	case tokSQLSTATE:
		return "SQLSTATE"
	case tokSTRING:
		return "STRING"
	case tokSTRUCT:
		return "STRUCT"
	case tokSUBSTRING:
		return "SUBSTRING"
	case tokSUM:
		return "SUM"
	case tokSYMBOL:
		return "SYMBOL"
	case tokSYSTEM_USER:
		return "SYSTEM_USER"
	case tokTABLE:
		return "TABLE"
	case tokTEMPORARY:
		return "TEMPORARY"
	case tokTHEN:
		return "THEN"
	case tokTIME:
		return "TIME"
	case tokTIMESTAMP:
		return "TIMESTAMP"
	case tokTO:
		return "TO"
	case tokTRANSACTION:
		return "TRANSACTION"
	case tokTRANSLATE:
		return "TRANSLATE"
	case tokTRANSLATION:
		return "TRANSLATION"
	case tokTRIM:
		return "TRIM"
	case tokTRUE:
		return "TRUE"
	case tokTUPLE:
		return "TUPLE"
	case tokUNION:
		return "UNION"
	case tokUNIQUE:
		return "UNIQUE"
	case tokUNKNOWN:
		return "UNKNOWN"
	case tokUNPIVOT:
		return "UNPIVOT"
	case tokUPDATE:
		return "UPDATE"
	case tokUPPER:
		return "UPPER"
	case tokUPSERT:
		return "UPSERT"
	case tokUSAGE:
		return "USAGE"
	case tokUSER:
		return "USER"
	case tokUSING:
		return "USING"
	case tokVALUE:
		return "VALUE"
	case tokVALUES:
		return "VALUES"
	case tokVARCHAR:
		return "VARCHAR"
	case tokVARYING:
		return "VARYING"
	case tokVIEW:
		return "VIEW"
	case tokWHEN:
		return "WHEN"
	case tokWHENEVER:
		return "WHENEVER"
	case tokWHERE:
		return "WHERE"
	case tokWITH:
		return "WITH"
	case tokWORK:
		return "WORK"
	case tokWRITE:
		return "WRITE"
	case tokZONE:
		return "ZONE"
	}
	return ""
}

// ===========================================================================
// Keyword tokens — group 3000.
//
// 266 keywords from PartiQLLexer.g4 lines 13–295, alphabetical.
//
// Includes standard SQL keywords (ABSOLUTE..ZONE), window keywords
// (LAG, LEAD, OVER, PARTITION), PartiQL extension keywords (CAN_CAST,
// CAN_LOSSLESS_CAST, MISSING, PIVOT, UNPIVOT, LIMIT, OFFSET, REMOVE,
// INDEX, LET, CONFLICT, DO, RETURNING, MODIFIED, NEW, OLD, NOTHING,
// EXCLUDED, SHORTEST, MATCH), and data type keywords (TUPLE, INT2/4/8,
// INTEGER2/4/8, BIGINT, BOOL, BOOLEAN, STRING, SYMBOL, CLOB, BLOB,
// STRUCT, LIST, SEXP, BAG).
// ===========================================================================

const (
	tokABSOLUTE = iota + 3000
	tokACTION
	tokADD
	tokALL
	tokALLOCATE
	tokALTER
	tokAND
	tokANY
	tokARE
	tokAS
	tokASC
	tokASSERTION
	tokAT
	tokAUTHORIZATION
	tokAVG
	tokBAG
	tokBEGIN
	tokBETWEEN
	tokBIGINT
	tokBIT
	tokBIT_LENGTH
	tokBLOB
	tokBOOL
	tokBOOLEAN
	tokBY
	tokCAN_CAST
	tokCAN_LOSSLESS_CAST
	tokCASCADE
	tokCASCADED
	tokCASE
	tokCAST
	tokCATALOG
	tokCHAR
	tokCHAR_LENGTH
	tokCHARACTER
	tokCHARACTER_LENGTH
	tokCHECK
	tokCLOB
	tokCLOSE
	tokCOALESCE
	tokCOLLATE
	tokCOLLATION
	tokCOLUMN
	tokCOMMIT
	tokCONFLICT
	tokCONNECT
	tokCONNECTION
	tokCONSTRAINT
	tokCONSTRAINTS
	tokCONTINUE
	tokCONVERT
	tokCORRESPONDING
	tokCOUNT
	tokCREATE
	tokCROSS
	tokCURRENT
	tokCURRENT_DATE
	tokCURRENT_TIME
	tokCURRENT_TIMESTAMP
	tokCURRENT_USER
	tokCURSOR
	tokDATE
	tokDATE_ADD
	tokDATE_DIFF
	tokDEALLOCATE
	tokDEC
	tokDECIMAL
	tokDECLARE
	tokDEFAULT
	tokDEFERRABLE
	tokDEFERRED
	tokDELETE
	tokDESC
	tokDESCRIBE
	tokDESCRIPTOR
	tokDIAGNOSTICS
	tokDISCONNECT
	tokDISTINCT
	tokDO
	tokDOMAIN
	tokDOUBLE
	tokDROP
	tokELSE
	tokEND
	tokEND_EXEC
	tokESCAPE
	tokEXCEPT
	tokEXCEPTION
	tokEXCLUDED
	tokEXEC
	tokEXECUTE
	tokEXISTS
	tokEXPLAIN
	tokEXTERNAL
	tokEXTRACT
	tokFALSE
	tokFETCH
	tokFIRST
	tokFLOAT
	tokFOR
	tokFOREIGN
	tokFOUND
	tokFROM
	tokFULL
	tokGET
	tokGLOBAL
	tokGO
	tokGOTO
	tokGRANT
	tokGROUP
	tokHAVING
	tokIDENTITY
	tokIMMEDIATE
	tokIN
	tokINDEX
	tokINDICATOR
	tokINITIALLY
	tokINNER
	tokINPUT
	tokINSENSITIVE
	tokINSERT
	tokINT
	tokINT2
	tokINT4
	tokINT8
	tokINTEGER
	tokINTEGER2
	tokINTEGER4
	tokINTEGER8
	tokINTERSECT
	tokINTERVAL
	tokINTO
	tokIS
	tokISOLATION
	tokJOIN
	tokKEY
	tokLAG
	tokLANGUAGE
	tokLAST
	tokLATERAL
	tokLEAD
	tokLEFT
	tokLET
	tokLEVEL
	tokLIKE
	tokLIMIT
	tokLIST
	tokLOCAL
	tokLOWER
	tokMATCH
	tokMAX
	tokMIN
	tokMISSING
	tokMODIFIED
	tokMODULE
	tokNAMES
	tokNATIONAL
	tokNATURAL
	tokNCHAR
	tokNEW
	tokNEXT
	tokNO
	tokNOT
	tokNOTHING
	tokNULL
	tokNULLIF
	tokNULLS
	tokNUMERIC
	tokOCTET_LENGTH
	tokOF
	tokOFFSET
	tokOLD
	tokON
	tokONLY
	tokOPEN
	tokOPTION
	tokOR
	tokORDER
	tokOUTER
	tokOUTPUT
	tokOVER
	tokOVERLAPS
	tokOVERLAY
	tokPAD
	tokPARTIAL
	tokPARTITION
	tokPIVOT
	tokPLACING
	tokPOSITION
	tokPRECISION
	tokPREPARE
	tokPRESERVE
	tokPRIMARY
	tokPRIOR
	tokPRIVILEGES
	tokPROCEDURE
	tokPUBLIC
	tokREAD
	tokREAL
	tokREFERENCES
	tokRELATIVE
	tokREMOVE
	tokREPLACE
	tokRESTRICT
	tokRETURNING
	tokREVOKE
	tokRIGHT
	tokROLLBACK
	tokROWS
	tokSCHEMA
	tokSCROLL
	tokSECTION
	tokSELECT
	tokSESSION
	tokSESSION_USER
	tokSET
	tokSEXP
	tokSHORTEST
	tokSIZE
	tokSMALLINT
	tokSOME
	tokSPACE
	tokSQL
	tokSQLCODE
	tokSQLERROR
	tokSQLSTATE
	tokSTRING
	tokSTRUCT
	tokSUBSTRING
	tokSUM
	tokSYMBOL
	tokSYSTEM_USER
	tokTABLE
	tokTEMPORARY
	tokTHEN
	tokTIME
	tokTIMESTAMP
	tokTO
	tokTRANSACTION
	tokTRANSLATE
	tokTRANSLATION
	tokTRIM
	tokTRUE
	tokTUPLE
	tokUNION
	tokUNIQUE
	tokUNKNOWN
	tokUNPIVOT
	tokUPDATE
	tokUPPER
	tokUPSERT
	tokUSAGE
	tokUSER
	tokUSING
	tokVALUE
	tokVALUES
	tokVARCHAR
	tokVARYING
	tokVIEW
	tokWHEN
	tokWHENEVER
	tokWHERE
	tokWITH
	tokWORK
	tokWRITE
	tokZONE
)
