package parser

import "github.com/bytebase/omni/googlesql/ast"

// Special tokens.
const (
	tokEOF     = 0
	tokInvalid = 1 // synthetic error-recovery token; always accompanied by a LexError
)

// Multi-character operators and literal token kinds. Single-character tokens
// use the ASCII byte value directly as their Type — no constant needed. The
// single-char GoogleSQL tokens that have no constant here are:
//
//   - - * / % ~ ! ^ & = < > ( ) [ ] { } , ; . : | @ ?
//
// (See GoogleSQLLexer.g4 §operator/§punctuation. '|' is STROKE_SYMBOL, '@' is
// AT_SYMBOL, '?' is QUESTION_SYMBOL, etc.)
const (
	tokInteger    = 600 + iota // INTEGER_LITERAL: decimal or 0x-hex
	tokFloat                   // FLOATING_POINT_LITERAL
	tokString                  // STRING_LITERAL (R? + 4 quote forms); quotes stripped
	tokBytes                   // BYTES_LITERAL ((B|RB|BR) + 4 quote forms); prefix+quotes stripped
	tokIdentifier              // IDENTIFIER: UNQUOTED_IDENTIFIER or `backtick`; backticks stripped

	// Multi-character operator / punctuation symbols (GoogleSQLLexer.g4). The
	// single-character forms (=, <, >, +, -, etc.) use their ASCII byte value.
	tokNotEqual     // != (NOT_EQUAL_OPERATOR)
	tokNotEqual2    // <> (NOT_EQUAL2_OPERATOR)
	tokLessEqual    // <= (LE_OPERATOR)
	tokGreaterEqual // >= (GE_OPERATOR)
	tokShiftLeft    // << (KL_OPERATOR)
	tokShiftRight   // >> (KR_OPERATOR)
	tokArrow        // -> (SUB_GT_BRACKET_SYMBOL) — lambda / function-type
	tokFatArrow     // => (EQUAL_GT_BRACKET_SYMBOL) — named argument
	tokPlusEqual    // += (PLUS_EQUAL_SYMBOL) — OPTIONS assignment
	tokMinusEqual   // -= (SUB_EQUAL_SYMBOL) — OPTIONS assignment
	tokPipe         // |> (PIPE_SYMBOL) — pipe-query operator (lexed; parser-gap)
	tokBoolOr       // || (BOOL_OR_SYMBOL) — logical/concat OR
	tokAtAt         // @@ (ATAT_SYMBOL) — system-variable prefix
)

// keywordBase is the first numeric value of the keyword token block. Keyword
// constants are NOT stable across edits (iota-assigned) — never persist them.
const keywordBase = 700

// Token represents a single lexical token.
type Token struct {
	Type  int     // tok* / kw* / ASCII byte value
	Str   string  // identifier text, string/bytes content (unquoted), raw operator text
	Ival  int64   // integer value for tokInteger (decimal and 0x-hex)
	Loc   ast.Loc // {Start, End} byte offsets in source text
	IsRaw bool    // true if string/bytes had an R prefix (raw); only set on tokString/tokBytes
}

// Keyword token constants, in GoogleSQLLexer.g4 declaration order. Generated
// from /Users/h3n4l/OpenSource/parser/googlesql/GoogleSQLLexer.g4 — every
// *_SYMBOL word-keyword token (308 total: 211 non-reserved + 97 reserved).
// The reserved/non-reserved split lives in keywords.go, driven by the parser
// grammar's common_keyword_as_identifier rule. Numeric values are NOT stable.
//
// Note: kwKW_MATCH_RECOGNIZE_NONRESERVED is a ZetaSQL placeholder token whose
// literal text is the string "KW_MATCH_RECOGNIZE_NONRESERVED"; it is a real
// token rule (used in an error alternative) but effectively never appears in
// real SQL. It is included for full grammar parity.
const (
	kwARRAY = keywordBase + iota
	kwALL
	kwAS
	kwASC
	kwBY
	kwCROSS
	kwJOIN
	kwDELTA
	kwDESC
	kwDIFFERENTIAL_PRIVACY
	kwDISTINCT
	kwEPSILON
	kwEXCEPT
	kwEXCLUDE
	kwFOR
	kwFROM
	kwFULL
	kwIN
	kwINCLUDE
	kwINNER
	kwINTERSECT
	kwLEFT
	kwLIMIT
	kwMAX_GROUPS_CONTRIBUTED
	kwNULL
	kwNULLS
	kwOF
	kwOFFSET
	kwON
	kwOPTIONS
	kwORDER
	kwOUTER
	kwPERCENT
	kwPIVOT
	kwPRIVACY_UNIT_COLUMN
	kwRIGHT
	kwRECURSIVE
	kwREPLACE
	kwUNPIVOT
	kwSELECT
	kwSTRUCT
	kwSYSTEM
	kwSYSTEM_TIME
	kwTABLESAMPLE
	kwUNION
	kwUNNEST
	kwUSING
	kwVALUE
	kwWITH
	kwTRUE
	kwFALSE
	kwNUMERIC
	kwDECIMAL
	kwBIGNUMERIC
	kwBIGDECIMAL
	kwNOT
	kwAND
	kwOR
	kwJSON
	kwDATE
	kwTIME
	kwDATETIME
	kwTIMESTAMP
	kwRANGE
	kwINTERVAL
	kwSIMPLE
	kwABORT
	kwACCESS
	kwACTION
	kwAGGREGATE
	kwADD
	kwALTER
	kwALWAYS
	kwANALYZE
	kwAPPROX
	kwARE
	kwASSERT
	kwBATCH
	kwBEGIN
	kwBREAK
	kwCALL
	kwCASCADE
	kwCHECK
	kwCLAMPED
	kwCLONE
	kwCOPY
	kwCLUSTER
	kwCOLUMN
	kwCOLUMNS
	kwCOMMIT
	kwCONNECTION
	kwCONSTANT
	kwCONSTRAINT
	kwCONTINUE
	kwCORRESPONDING
	kwCYCLE
	kwDATA
	kwDATABASE
	kwDECLARE
	kwDEFINER
	kwDELETE
	kwDELETION
	kwDEPTH
	kwDESCRIBE
	kwDETERMINISTIC
	kwDO
	kwDROP
	kwELSEIF
	kwENFORCED
	kwERROR
	kwEXCEPTION
	kwEXECUTE
	kwEXPLAIN
	kwEXPORT
	kwEXTEND
	kwEXTERNAL
	kwFILES
	kwFILTER
	kwFILL
	kwFIRST
	kwFOREIGN
	kwFORMAT
	kwFUNCTION
	kwGENERATED
	kwGRANT
	kwGROUP_ROWS
	kwHIDDEN
	kwIDENTITY
	kwIMMEDIATE
	kwIMMUTABLE
	kwIMPORT
	kwINCREMENT
	kwINDEX
	kwINOUT
	kwINPUT
	kwINSERT
	kwINVOKER
	kwISOLATION
	kwITERATE
	kwKEY
	kwLANGUAGE
	kwLAST
	kwLEAVE
	kwLEVEL
	kwLOAD
	kwLOOP
	kwMACRO
	kwMAP
	kwMATCH
	kwKW_MATCH_RECOGNIZE_NONRESERVED
	kwMATCHED
	kwMATERIALIZED
	kwMAX
	kwMAXVALUE
	kwMEASURES
	kwMESSAGE
	kwMETADATA
	kwMIN
	kwMINVALUE
	kwMODEL
	kwMODULE
	kwONLY
	kwOUT
	kwOUTPUT
	kwOVERWRITE
	kwPARTITIONS
	kwPATTERN
	kwPOLICIES
	kwPOLICY
	kwPRIMARY
	kwPRIVATE
	kwPRIVILEGE
	kwPRIVILEGES
	kwPROCEDURE
	kwPROJECT
	kwPUBLIC
	kwRAISE
	kwREAD
	kwREFERENCES
	kwREMOTE
	kwREMOVE
	kwRENAME
	kwREPEAT
	kwREPEATABLE
	kwREPLACE_FIELDS
	kwREPLICA
	kwREPORT
	kwRESTRICT
	kwRESTRICTION
	kwRETURNS
	kwRETURN
	kwREVOKE
	kwROLLBACK
	kwROW
	kwRUN
	kwSAFE_CAST
	kwSCHEMA
	kwSEARCH
	kwSECURITY
	kwSEQUENCE
	kwSETS
	kwSET
	kwSHOW
	kwSNAPSHOT
	kwSOURCE
	kwSQL
	kwSTABLE
	kwSTART
	kwSTATIC_DESCRIBE
	kwSTORED
	kwSTORING
	kwSTRICT
	kwTABLE
	kwTABLES
	kwTARGET
	kwTEMP
	kwTEMPORARY
	kwTRANSACTION
	kwTRANSFORM
	kwTRUNCATE
	kwTYPE
	kwUNDROP
	kwUNIQUE
	kwUNKNOWN
	kwUNTIL
	kwUPDATE
	kwVALUES
	kwVECTOR
	kwVIEW
	kwVIEWS
	kwVOLATILE
	kwWEIGHT
	kwWHILE
	kwWRITE
	kwZONE
	kwDESCRIPTOR
	kwINTERLEAVE
	kwNULL_FILTERED
	kwPARENT
	kwNEW
	kwEND
	kwCASE
	kwWHEN
	kwTHEN
	kwELSE
	kwCAST
	kwEXTRACT
	kwCOLLATE
	kwIF
	kwGROUPING
	kwHAVING
	kwGROUP
	kwROLLUP
	kwCUBE
	kwHASH
	kwPROTO
	kwPARTITION
	kwIGNORE
	kwRESPECT
	kwROWS
	kwOVER
	kwBETWEEN
	kwUNBOUNDED
	kwCURRENT
	kwPRECEDING
	kwFOLLOWING
	kwNATURAL
	kwQUALIFY
	kwDEFAULT
	kwSLASH
	kwMATCH_RECOGNIZE
	kwDEFINE
	kwLOOKUP
	kwWHERE
	kwWINDOW
	kwTO
	kwEXISTS
	kwANY
	kwSOME
	kwLIKE
	kwIS
	kwNO
	kwINTO
	kwASSERT_ROWS_MODIFIED
	kwCONFLICT
	kwNOTHING
	kwMERGE
	kwCREATE
	kwENUM
	kwDESTINATION
	kwPROPERTY
	kwGRAPH
	kwNODE
	kwPROPERTIES
	kwLABEL
	kwEDGE
	kwNEXT
	kwASCENDING
	kwDESCENDING
	kwSKIP
	kwSHORTEST
	kwPATH
	kwPATHS
	kwWALK
	kwTRAIL
	kwACYCLIC
	kwOPTIONAL
	kwLET
)

// TokenName returns a human-readable name for a token type. Used by tests and
// debug output.
//
// Returns:
//   - "EOF" for tokEOF, "INVALID" for tokInvalid
//   - the upper-case keyword name (e.g. "SELECT") for kw* constants
//   - a symbolic name for tok* literal/operator constants
//   - the literal character in quotes (e.g. "'+'") for ASCII single-char tokens
func TokenName(t int) string {
	switch t {
	case tokEOF:
		return "EOF"
	case tokInvalid:
		return "INVALID"
	case tokInteger:
		return "INTEGER"
	case tokFloat:
		return "FLOAT"
	case tokString:
		return "STRING"
	case tokBytes:
		return "BYTES"
	case tokIdentifier:
		return "IDENTIFIER"
	case tokNotEqual:
		return "!="
	case tokNotEqual2:
		return "<>"
	case tokLessEqual:
		return "<="
	case tokGreaterEqual:
		return ">="
	case tokShiftLeft:
		return "<<"
	case tokShiftRight:
		return ">>"
	case tokArrow:
		return "->"
	case tokFatArrow:
		return "=>"
	case tokPlusEqual:
		return "+="
	case tokMinusEqual:
		return "-="
	case tokPipe:
		return "|>"
	case tokBoolOr:
		return "||"
	case tokAtAt:
		return "@@"
	}
	if t >= keywordBase {
		if name, ok := keywordName(t); ok {
			return name
		}
		return "UNKNOWN_KEYWORD"
	}
	if t > 0 && t < 128 {
		return "'" + string(rune(t)) + "'"
	}
	return "UNKNOWN"
}
