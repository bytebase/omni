package parser

import "strings"

// TokenKind identifies the lexical class of a Databricks/Spark SQL token.
type TokenKind int

const (
	tokEOF TokenKind = iota
	tokInvalid

	// Identifiers and literals.
	tokIdent        // unquoted ASCII identifier ([A-Za-z0-9_]+)
	tokRawWord      // non-ASCII word or scheme:// URI (valid only in raw .*? tails, rejected as SQL identifier)
	tokBackquoted   // `...`
	tokString       // '...' or R'...' or R"..."
	tokDoubleQuoted // "..."
	tokInt          // 123
	tokDecimal      // 1.5, .5, 1.
	tokExponent     // 1E1, 1.5E-2
	tokBigInt       // 1L
	tokSmallInt     // 1S
	tokTinyInt      // 1Y
	tokFloat        // 1F, 1.5F
	tokDouble       // 1D, 1.5D
	tokBigDecimal   // 1BD, 1.5BD

	// Punctuation and operators.
	tokSemicolon    // ;
	tokLParen       // (
	tokRParen       // )
	tokComma        // ,
	tokDot          // .
	tokLBracket     // [
	tokRBracket     // ]
	tokBang         // !
	tokEq           // = or ==
	tokNullSafeEq   // <=>
	tokNeq          // <>
	tokNeqBang      // !=
	tokLT           // <
	tokLTE          // <= or !>
	tokGT           // >
	tokGTE          // >= or !<
	tokShiftLeft    // <<
	tokShiftRight   // >>
	tokShiftRightU  // >>>
	tokPlus         // +
	tokMinus        // -
	tokAsterisk     // *
	tokSlash        // /
	tokPercent      // %
	tokTilde        // ~
	tokAmpersand    // &
	tokPipe         // |
	tokConcatPipe   // ||
	tokOperatorPipe // |>
	tokHat          // ^
	tokColon        // :
	tokDoubleColon  // ::
	tokArrow        // ->
	tokFatArrow     // =>
	tokHintStart    // /*+
	tokHintEnd      // */
	tokQuestion     // ?

	// Keywords (matching SqlBaseLexer.g4).
	kwStart
	kwADD
	kwAFTER
	kwAGGREGATE
	kwALL
	kwALTER
	kwALWAYS
	kwANALYZE
	kwAND
	kwANTI
	kwANY
	kwANY_VALUE
	kwARCHIVE
	kwARRAY
	kwAS
	kwASC
	kwAT
	kwATOMIC
	kwAUTHORIZATION
	kwBEGIN
	kwBETWEEN
	kwBIGINT
	kwBINARY
	kwBINARY_HEX // X
	kwBINDING
	kwBOOLEAN
	kwBOTH
	kwBUCKET
	kwBUCKETS
	kwBY
	kwBYTE
	kwCACHE
	kwCALL
	kwCALLED
	kwCASCADE
	kwCASE
	kwCAST
	kwCATALOG
	kwCATALOGS
	kwCHANGE
	kwCHAR
	kwCHARACTER
	kwCHECK
	kwCLEAR
	kwCLUSTER
	kwCLUSTERED
	kwCODEGEN
	kwCOLLATE
	kwCOLLATION
	kwCOLLECTION
	kwCOLUMN
	kwCOLUMNS
	kwCOMMENT
	kwCOMMIT
	kwCOMPACT
	kwCOMPACTIONS
	kwCOMPENSATION
	kwCOMPUTE
	kwCONCATENATE
	kwCONDITION
	kwCONSTRAINT
	kwCONTAINS
	kwCONTINUE
	kwCOST
	kwCREATE
	kwCROSS
	kwCUBE
	kwCURRENT
	kwCURRENT_DATE
	kwCURRENT_TIME
	kwCURRENT_TIMESTAMP
	kwCURRENT_USER
	kwDAY
	kwDAYS
	kwDAYOFYEAR
	kwDATA
	kwDATE
	kwDATABASE
	kwDATABASES
	kwDATEADD
	kwDATE_ADD
	kwDATEDIFF
	kwDATE_DIFF
	kwDBPROPERTIES
	kwDEC
	kwDECIMAL
	kwDECLARE
	kwDEFAULT
	kwDEFINED
	kwDEFINER
	kwDELETE
	kwDELIMITED
	kwDESC
	kwDESCRIBE
	kwDETERMINISTIC
	kwDFS
	kwDIRECTORIES
	kwDIRECTORY
	kwDISTINCT
	kwDISTRIBUTE
	kwDIV
	kwDO
	kwDOUBLE
	kwDROP
	kwELSE
	kwELSEIF
	kwEND
	kwESCAPE
	kwESCAPED
	kwEVOLUTION
	kwEXCEPT
	kwEXCHANGE
	kwEXCLUDE
	kwEXECUTE
	kwEXISTS
	kwEXIT
	kwEXPLAIN
	kwEXPORT
	kwEXTEND
	kwEXTENDED
	kwEXTERNAL
	kwEXTRACT
	kwFALSE
	kwFETCH
	kwFIELDS
	kwFILTER
	kwFILEFORMAT
	kwFIRST
	kwFLOAT
	kwFOLLOWING
	kwFOR
	kwFOREIGN
	kwFORMAT
	kwFORMATTED
	kwFOUND
	kwFROM
	kwFULL
	kwFUNCTION
	kwFUNCTIONS
	kwGENERATED
	kwGLOBAL
	kwGRANT
	kwGROUP
	kwGROUPING
	kwHANDLER
	kwHAVING
	kwHOUR
	kwHOURS
	kwIDENTIFIER
	kwIDENTITY
	kwIF
	kwIGNORE
	kwILIKE
	kwIMMEDIATE
	kwIMPORT
	kwIN
	kwINCLUDE
	kwINCREMENT
	kwINDEX
	kwINDEXES
	kwINNER
	kwINPATH
	kwINPUT
	kwINPUTFORMAT
	kwINSERT
	kwINTERSECT
	kwINTERVAL
	kwINT
	kwINTEGER
	kwINTO
	kwINVOKER
	kwIS
	kwITEMS
	kwITERATE
	kwJOIN
	kwJSON
	kwKEYS
	kwLANGUAGE
	kwLAST
	kwLATERAL
	kwLAZY
	kwLEADING
	kwLEAVE
	kwLEFT
	kwLIKE
	kwLIMIT
	kwLINES
	kwLIST
	kwLOAD
	kwLOCAL
	kwLOCATION
	kwLOCK
	kwLOCKS
	kwLOGICAL
	kwLONG
	kwLOOP
	kwMACRO
	kwMAP
	kwMATCHED
	kwMERGE
	kwMICROSECOND
	kwMICROSECONDS
	kwMILLISECOND
	kwMILLISECONDS
	kwMINUS
	kwMINUTE
	kwMINUTES
	kwMODIFIES
	kwMONTH
	kwMONTHS
	kwMSCK
	kwNAME
	kwNAMESPACE
	kwNAMESPACES
	kwNANOSECOND
	kwNANOSECONDS
	kwNATURAL
	kwNO
	kwNONE
	kwNOT
	kwNULL
	kwNULLS
	kwNUMERIC
	kwOF
	kwOFFSET
	kwON
	kwONLY
	kwOPTION
	kwOPTIONS
	kwOR
	kwORDER
	kwOUT
	kwOUTER
	kwOUTPUTFORMAT
	kwOVER
	kwOVERLAPS
	kwOVERLAY
	kwOVERWRITE
	kwPARTITION
	kwPARTITIONED
	kwPARTITIONS
	kwPERCENT
	kwPIVOT
	kwPLACING
	kwPOSITION
	kwPRECEDING
	kwPRIMARY
	kwPRINCIPALS
	kwPROPERTIES
	kwPURGE
	kwQUARTER
	kwQUERY
	kwRANGE
	kwREADS
	kwREAL
	kwRECORDREADER
	kwRECORDWRITER
	kwRECOVER
	kwRECURSIVE
	kwREDUCE
	kwREFERENCES
	kwREFRESH
	kwRENAME
	kwREPAIR
	kwREPEAT
	kwREPEATABLE
	kwREPLACE
	kwRESET
	kwRESPECT
	kwRESTRICT
	kwRETURN
	kwRETURNS
	kwREVOKE
	kwRIGHT
	kwRLIKE
	kwROLE
	kwROLES
	kwROLLBACK
	kwROLLUP
	kwROW
	kwROWS
	kwSECOND
	kwSECONDS
	kwSCHEMA
	kwSCHEMAS
	kwSECURITY
	kwSELECT
	kwSEMI
	kwSEPARATED
	kwSERDE
	kwSERDEPROPERTIES
	kwSESSION_USER
	kwSET
	kwSETS
	kwSHORT
	kwSHOW
	kwSINGLE
	kwSKEWED
	kwSMALLINT
	kwSOME
	kwSORT
	kwSORTED
	kwSOURCE
	kwSPECIFIC
	kwSQL
	kwSQLEXCEPTION
	kwSQLSTATE
	kwSTART
	kwSTATISTICS
	kwSTORED
	kwSTRATIFY
	kwSTRING
	kwSTRUCT
	kwSUBSTR
	kwSUBSTRING
	kwSYNC
	kwSYSTEM_TIME
	kwSYSTEM_VERSION
	kwTABLE
	kwTABLES
	kwTABLESAMPLE
	kwTARGET
	kwTBLPROPERTIES
	kwTEMPORARY
	kwTERMINATED
	kwTHEN
	kwTIME
	kwTIMEDIFF
	kwTIMESTAMP
	kwTIMESTAMP_LTZ
	kwTIMESTAMP_NTZ
	kwTIMESTAMPADD
	kwTIMESTAMPDIFF
	kwTINYINT
	kwTO
	kwTOUCH
	kwTRAILING
	kwTRANSACTION
	kwTRANSACTIONS
	kwTRANSFORM
	kwTRIM
	kwTRUE
	kwTRUNCATE
	kwTRY_CAST
	kwTYPE
	kwUNARCHIVE
	kwUNBOUNDED
	kwUNCACHE
	kwUNION
	kwUNIQUE
	kwUNKNOWN
	kwUNLOCK
	kwUNPIVOT
	kwUNSET
	kwUNTIL
	kwUPDATE
	kwUSE
	kwUSER
	kwUSING
	kwVALUE
	kwVALUES
	kwVARCHAR
	kwVAR
	kwVARIABLE
	kwVARIANT
	kwVERSION
	kwVIEW
	kwVIEWS
	kwVOID
	kwWEEK
	kwWEEKS
	kwWHEN
	kwWHERE
	kwWHILE
	kwWINDOW
	kwWITH
	kwWITHIN
	kwYEAR
	kwYEARS
	kwZONE
	kwEnd
)

// Token represents a single lexed token with byte offsets.
type Token struct {
	Kind  TokenKind
	Text  string
	Start int
	End   int
}

func (k TokenKind) isKeyword() bool {
	return k > kwStart && k < kwEnd
}

// isStrictNonReservedOnly reports whether k is in Spark's strictNonReserved
// list, meaning it cannot be used as a strictIdentifier (e.g. a table alias)
// in Spark's default mode (SQL_standard_keyword_behavior=false).
func (k TokenKind) isStrictNonReservedOnly() bool {
	switch k {
	case kwANTI, kwCROSS, kwEXCEPT, kwFULL, kwINNER, kwINTERSECT, kwJOIN,
		kwLATERAL, kwLEFT, kwMINUS, kwNATURAL, kwON, kwRIGHT, kwSEMI, kwUNION, kwUSING:
		return true
	}
	return false
}

// isStrictIdentifier reports whether k can be parsed as a Spark strictIdentifier
// under Spark's default quoting configuration (double_quoted_identifiers=false).
func (k TokenKind) isStrictIdentifier() bool {
	if k == tokIdent || k == tokBackquoted {
		return true
	}
	return k.isKeyword() && !k.isStrictNonReservedOnly()
}

// isIdentifier reports whether k can be parsed as a Spark identifier under
// Spark's default quoting configuration (double_quoted_identifiers=false).
func (k TokenKind) isIdentifier() bool {
	return k == tokIdent || k == tokBackquoted || k.isKeyword()
}

var keywords = map[string]TokenKind{
	"ADD":               kwADD,
	"AFTER":             kwAFTER,
	"AGGREGATE":         kwAGGREGATE,
	"ALL":               kwALL,
	"ALTER":             kwALTER,
	"ALWAYS":            kwALWAYS,
	"ANALYZE":           kwANALYZE,
	"AND":               kwAND,
	"ANTI":              kwANTI,
	"ANY":               kwANY,
	"ANY_VALUE":         kwANY_VALUE,
	"ARCHIVE":           kwARCHIVE,
	"ARRAY":             kwARRAY,
	"AS":                kwAS,
	"ASC":               kwASC,
	"AT":                kwAT,
	"ATOMIC":            kwATOMIC,
	"AUTHORIZATION":     kwAUTHORIZATION,
	"BEGIN":             kwBEGIN,
	"BETWEEN":           kwBETWEEN,
	"BIGINT":            kwBIGINT,
	"BINARY":            kwBINARY,
	"X":                 kwBINARY_HEX,
	"BINDING":           kwBINDING,
	"BOOLEAN":           kwBOOLEAN,
	"BOTH":              kwBOTH,
	"BUCKET":            kwBUCKET,
	"BUCKETS":           kwBUCKETS,
	"BY":                kwBY,
	"BYTE":              kwBYTE,
	"CACHE":             kwCACHE,
	"CALL":              kwCALL,
	"CALLED":            kwCALLED,
	"CASCADE":           kwCASCADE,
	"CASE":              kwCASE,
	"CAST":              kwCAST,
	"CATALOG":           kwCATALOG,
	"CATALOGS":          kwCATALOGS,
	"CHANGE":            kwCHANGE,
	"CHAR":              kwCHAR,
	"CHARACTER":         kwCHARACTER,
	"CHECK":             kwCHECK,
	"CLEAR":             kwCLEAR,
	"CLUSTER":           kwCLUSTER,
	"CLUSTERED":         kwCLUSTERED,
	"CODEGEN":           kwCODEGEN,
	"COLLATE":           kwCOLLATE,
	"COLLATION":         kwCOLLATION,
	"COLLECTION":        kwCOLLECTION,
	"COLUMN":            kwCOLUMN,
	"COLUMNS":           kwCOLUMNS,
	"COMMENT":           kwCOMMENT,
	"COMMIT":            kwCOMMIT,
	"COMPACT":           kwCOMPACT,
	"COMPACTIONS":       kwCOMPACTIONS,
	"COMPENSATION":      kwCOMPENSATION,
	"COMPUTE":           kwCOMPUTE,
	"CONCATENATE":       kwCONCATENATE,
	"CONDITION":         kwCONDITION,
	"CONSTRAINT":        kwCONSTRAINT,
	"CONTAINS":          kwCONTAINS,
	"CONTINUE":          kwCONTINUE,
	"COST":              kwCOST,
	"CREATE":            kwCREATE,
	"CROSS":             kwCROSS,
	"CUBE":              kwCUBE,
	"CURRENT":           kwCURRENT,
	"CURRENT_DATE":      kwCURRENT_DATE,
	"CURRENT_TIME":      kwCURRENT_TIME,
	"CURRENT_TIMESTAMP": kwCURRENT_TIMESTAMP,
	"CURRENT_USER":      kwCURRENT_USER,
	"DAY":               kwDAY,
	"DAYS":              kwDAYS,
	"DAYOFYEAR":         kwDAYOFYEAR,
	"DATA":              kwDATA,
	"DATE":              kwDATE,
	"DATABASE":          kwDATABASE,
	"DATABASES":         kwDATABASES,
	"DATEADD":           kwDATEADD,
	"DATE_ADD":          kwDATE_ADD,
	"DATEDIFF":          kwDATEDIFF,
	"DATE_DIFF":         kwDATE_DIFF,
	"DBPROPERTIES":      kwDBPROPERTIES,
	"DEC":               kwDEC,
	"DECIMAL":           kwDECIMAL,
	"DECLARE":           kwDECLARE,
	"DEFAULT":           kwDEFAULT,
	"DEFINED":           kwDEFINED,
	"DEFINER":           kwDEFINER,
	"DELETE":            kwDELETE,
	"DELIMITED":         kwDELIMITED,
	"DESC":              kwDESC,
	"DESCRIBE":          kwDESCRIBE,
	"DETERMINISTIC":     kwDETERMINISTIC,
	"DFS":               kwDFS,
	"DIRECTORIES":       kwDIRECTORIES,
	"DIRECTORY":         kwDIRECTORY,
	"DISTINCT":          kwDISTINCT,
	"DISTRIBUTE":        kwDISTRIBUTE,
	"DIV":               kwDIV,
	"DO":                kwDO,
	"DOUBLE":            kwDOUBLE,
	"DROP":              kwDROP,
	"ELSE":              kwELSE,
	"ELSEIF":            kwELSEIF,
	"END":               kwEND,
	"ESCAPE":            kwESCAPE,
	"ESCAPED":           kwESCAPED,
	"EVOLUTION":         kwEVOLUTION,
	"EXCEPT":            kwEXCEPT,
	"EXCHANGE":          kwEXCHANGE,
	"EXCLUDE":           kwEXCLUDE,
	"EXECUTE":           kwEXECUTE,
	"EXISTS":            kwEXISTS,
	"EXIT":              kwEXIT,
	"EXPLAIN":           kwEXPLAIN,
	"EXPORT":            kwEXPORT,
	"EXTEND":            kwEXTEND,
	"EXTENDED":          kwEXTENDED,
	"EXTERNAL":          kwEXTERNAL,
	"EXTRACT":           kwEXTRACT,
	"FALSE":             kwFALSE,
	"FETCH":             kwFETCH,
	"FIELDS":            kwFIELDS,
	"FILTER":            kwFILTER,
	"FILEFORMAT":        kwFILEFORMAT,
	"FIRST":             kwFIRST,
	"FLOAT":             kwFLOAT,
	"FOLLOWING":         kwFOLLOWING,
	"FOR":               kwFOR,
	"FOREIGN":           kwFOREIGN,
	"FORMAT":            kwFORMAT,
	"FORMATTED":         kwFORMATTED,
	"FOUND":             kwFOUND,
	"FROM":              kwFROM,
	"FULL":              kwFULL,
	"FUNCTION":          kwFUNCTION,
	"FUNCTIONS":         kwFUNCTIONS,
	"GENERATED":         kwGENERATED,
	"GLOBAL":            kwGLOBAL,
	"GRANT":             kwGRANT,
	"GROUP":             kwGROUP,
	"GROUPING":          kwGROUPING,
	"HANDLER":           kwHANDLER,
	"HAVING":            kwHAVING,
	"HOUR":              kwHOUR,
	"HOURS":             kwHOURS,
	"IDENTIFIER":        kwIDENTIFIER,
	"IDENTITY":          kwIDENTITY,
	"IF":                kwIF,
	"IGNORE":            kwIGNORE,
	"ILIKE":             kwILIKE,
	"IMMEDIATE":         kwIMMEDIATE,
	"IMPORT":            kwIMPORT,
	"IN":                kwIN,
	"INCLUDE":           kwINCLUDE,
	"INCREMENT":         kwINCREMENT,
	"INDEX":             kwINDEX,
	"INDEXES":           kwINDEXES,
	"INNER":             kwINNER,
	"INPATH":            kwINPATH,
	"INPUT":             kwINPUT,
	"INPUTFORMAT":       kwINPUTFORMAT,
	"INSERT":            kwINSERT,
	"INTERSECT":         kwINTERSECT,
	"INTERVAL":          kwINTERVAL,
	"INT":               kwINT,
	"INTEGER":           kwINTEGER,
	"INTO":              kwINTO,
	"INVOKER":           kwINVOKER,
	"IS":                kwIS,
	"ITEMS":             kwITEMS,
	"ITERATE":           kwITERATE,
	"JOIN":              kwJOIN,
	"JSON":              kwJSON,
	"KEYS":              kwKEYS,
	"LANGUAGE":          kwLANGUAGE,
	"LAST":              kwLAST,
	"LATERAL":           kwLATERAL,
	"LAZY":              kwLAZY,
	"LEADING":           kwLEADING,
	"LEAVE":             kwLEAVE,
	"LEFT":              kwLEFT,
	"LIKE":              kwLIKE,
	"LIMIT":             kwLIMIT,
	"LINES":             kwLINES,
	"LIST":              kwLIST,
	"LOAD":              kwLOAD,
	"LOCAL":             kwLOCAL,
	"LOCATION":          kwLOCATION,
	"LOCK":              kwLOCK,
	"LOCKS":             kwLOCKS,
	"LOGICAL":           kwLOGICAL,
	"LONG":              kwLONG,
	"LOOP":              kwLOOP,
	"MACRO":             kwMACRO,
	"MAP":               kwMAP,
	"MATCHED":           kwMATCHED,
	"MERGE":             kwMERGE,
	"MICROSECOND":       kwMICROSECOND,
	"MICROSECONDS":      kwMICROSECONDS,
	"MILLISECOND":       kwMILLISECOND,
	"MILLISECONDS":      kwMILLISECONDS,
	"MINUS":             kwMINUS,
	"MINUTE":            kwMINUTE,
	"MINUTES":           kwMINUTES,
	"MODIFIES":          kwMODIFIES,
	"MONTH":             kwMONTH,
	"MONTHS":            kwMONTHS,
	"MSCK":              kwMSCK,
	"NAME":              kwNAME,
	"NAMESPACE":         kwNAMESPACE,
	"NAMESPACES":        kwNAMESPACES,
	"NANOSECOND":        kwNANOSECOND,
	"NANOSECONDS":       kwNANOSECONDS,
	"NATURAL":           kwNATURAL,
	"NO":                kwNO,
	"NONE":              kwNONE,
	"NOT":               kwNOT,
	"NULL":              kwNULL,
	"NULLS":             kwNULLS,
	"NUMERIC":           kwNUMERIC,
	"OF":                kwOF,
	"OFFSET":            kwOFFSET,
	"ON":                kwON,
	"ONLY":              kwONLY,
	"OPTION":            kwOPTION,
	"OPTIONS":           kwOPTIONS,
	"OR":                kwOR,
	"ORDER":             kwORDER,
	"OUT":               kwOUT,
	"OUTER":             kwOUTER,
	"OUTPUTFORMAT":      kwOUTPUTFORMAT,
	"OVER":              kwOVER,
	"OVERLAPS":          kwOVERLAPS,
	"OVERLAY":           kwOVERLAY,
	"OVERWRITE":         kwOVERWRITE,
	"PARTITION":         kwPARTITION,
	"PARTITIONED":       kwPARTITIONED,
	"PARTITIONS":        kwPARTITIONS,
	"PERCENT":           kwPERCENT,
	"PIVOT":             kwPIVOT,
	"PLACING":           kwPLACING,
	"POSITION":          kwPOSITION,
	"PRECEDING":         kwPRECEDING,
	"PRIMARY":           kwPRIMARY,
	"PRINCIPALS":        kwPRINCIPALS,
	"PROPERTIES":        kwPROPERTIES,
	"PURGE":             kwPURGE,
	"QUARTER":           kwQUARTER,
	"QUERY":             kwQUERY,
	"RANGE":             kwRANGE,
	"READS":             kwREADS,
	"REAL":              kwREAL,
	"RECORDREADER":      kwRECORDREADER,
	"RECORDWRITER":      kwRECORDWRITER,
	"RECOVER":           kwRECOVER,
	"RECURSIVE":         kwRECURSIVE,
	"REDUCE":            kwREDUCE,
	"REFERENCES":        kwREFERENCES,
	"REFRESH":           kwREFRESH,
	"REGEXP":            kwRLIKE,
	"RENAME":            kwRENAME,
	"REPAIR":            kwREPAIR,
	"REPEAT":            kwREPEAT,
	"REPEATABLE":        kwREPEATABLE,
	"REPLACE":           kwREPLACE,
	"RESET":             kwRESET,
	"RESPECT":           kwRESPECT,
	"RESTRICT":          kwRESTRICT,
	"RETURN":            kwRETURN,
	"RETURNS":           kwRETURNS,
	"REVOKE":            kwREVOKE,
	"RIGHT":             kwRIGHT,
	"RLIKE":             kwRLIKE,
	"ROLE":              kwROLE,
	"ROLES":             kwROLES,
	"ROLLBACK":          kwROLLBACK,
	"ROLLUP":            kwROLLUP,
	"ROW":               kwROW,
	"ROWS":              kwROWS,
	"SECOND":            kwSECOND,
	"SECONDS":           kwSECONDS,
	"SCHEMA":            kwSCHEMA,
	"SCHEMAS":           kwSCHEMAS,
	"SECURITY":          kwSECURITY,
	"SELECT":            kwSELECT,
	"SEMI":              kwSEMI,
	"SEPARATED":         kwSEPARATED,
	"SERDE":             kwSERDE,
	"SERDEPROPERTIES":   kwSERDEPROPERTIES,
	"SESSION_USER":      kwSESSION_USER,
	"SET":               kwSET,
	"SETS":              kwSETS,
	"SHORT":             kwSHORT,
	"SHOW":              kwSHOW,
	"SINGLE":            kwSINGLE,
	"SKEWED":            kwSKEWED,
	"SMALLINT":          kwSMALLINT,
	"SOME":              kwSOME,
	"SORT":              kwSORT,
	"SORTED":            kwSORTED,
	"SOURCE":            kwSOURCE,
	"SPECIFIC":          kwSPECIFIC,
	"SQL":               kwSQL,
	"SQLEXCEPTION":      kwSQLEXCEPTION,
	"SQLSTATE":          kwSQLSTATE,
	"START":             kwSTART,
	"STATISTICS":        kwSTATISTICS,
	"STORED":            kwSTORED,
	"STRATIFY":          kwSTRATIFY,
	"STRING":            kwSTRING,
	"STRUCT":            kwSTRUCT,
	"SUBSTR":            kwSUBSTR,
	"SUBSTRING":         kwSUBSTRING,
	"SYNC":              kwSYNC,
	"SYSTEM_TIME":       kwSYSTEM_TIME,
	"SYSTEM_VERSION":    kwSYSTEM_VERSION,
	"TABLE":             kwTABLE,
	"TABLES":            kwTABLES,
	"TABLESAMPLE":       kwTABLESAMPLE,
	"TARGET":            kwTARGET,
	"TBLPROPERTIES":     kwTBLPROPERTIES,
	"TEMP":              kwTEMPORARY,
	"TEMPORARY":         kwTEMPORARY,
	"TERMINATED":        kwTERMINATED,
	"THEN":              kwTHEN,
	"TIME":              kwTIME,
	"TIMEDIFF":          kwTIMEDIFF,
	"TIMESTAMP":         kwTIMESTAMP,
	"TIMESTAMP_LTZ":     kwTIMESTAMP_LTZ,
	"TIMESTAMP_NTZ":     kwTIMESTAMP_NTZ,
	"TIMESTAMPADD":      kwTIMESTAMPADD,
	"TIMESTAMPDIFF":     kwTIMESTAMPDIFF,
	"TINYINT":           kwTINYINT,
	"TO":                kwTO,
	"TOUCH":             kwTOUCH,
	"TRAILING":          kwTRAILING,
	"TRANSACTION":       kwTRANSACTION,
	"TRANSACTIONS":      kwTRANSACTIONS,
	"TRANSFORM":         kwTRANSFORM,
	"TRIM":              kwTRIM,
	"TRUE":              kwTRUE,
	"TRUNCATE":          kwTRUNCATE,
	"TRY_CAST":          kwTRY_CAST,
	"TYPE":              kwTYPE,
	"UNARCHIVE":         kwUNARCHIVE,
	"UNBOUNDED":         kwUNBOUNDED,
	"UNCACHE":           kwUNCACHE,
	"UNION":             kwUNION,
	"UNIQUE":            kwUNIQUE,
	"UNKNOWN":           kwUNKNOWN,
	"UNLOCK":            kwUNLOCK,
	"UNPIVOT":           kwUNPIVOT,
	"UNSET":             kwUNSET,
	"UNTIL":             kwUNTIL,
	"UPDATE":            kwUPDATE,
	"USE":               kwUSE,
	"USER":              kwUSER,
	"USING":             kwUSING,
	"VALUE":             kwVALUE,
	"VALUES":            kwVALUES,
	"VARCHAR":           kwVARCHAR,
	"VAR":               kwVAR,
	"VARIABLE":          kwVARIABLE,
	"VARIANT":           kwVARIANT,
	"VERSION":           kwVERSION,
	"VIEW":              kwVIEW,
	"VIEWS":             kwVIEWS,
	"VOID":              kwVOID,
	"WEEK":              kwWEEK,
	"WEEKS":             kwWEEKS,
	"WHEN":              kwWHEN,
	"WHERE":             kwWHERE,
	"WHILE":             kwWHILE,
	"WINDOW":            kwWINDOW,
	"WITH":              kwWITH,
	"WITHIN":            kwWITHIN,
	"YEAR":              kwYEAR,
	"YEARS":             kwYEARS,
	"ZONE":              kwZONE,
}

func lookupKeyword(word string) TokenKind {
	if kw, ok := keywords[strings.ToUpper(word)]; ok {
		return kw
	}
	return tokIdent
}
