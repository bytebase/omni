// Package parser implements a recursive descent SQL parser for MySQL.
package parser

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Token type constants for literals and operators.
const (
	tokEOF = 0

	// Literal tokens (600+)
	tokICONST = 600 + iota
	tokFCONST
	tokSCONST
	tokBCONST // bit string b'...'
	tokXCONST // hex string X'...'
	tokIDENT

	// Operators
	tokLessEq      // <=
	tokGreaterEq   // >=
	tokNotEq       // != or <>
	tokNullSafeEq  // <=>
	tokShiftLeft   // <<
	tokShiftRight  // >>
	tokAssign      // :=
	tokColonColon  // :: (not MySQL, but for compat)
	tokJsonExtract // ->
	tokJsonUnquote // ->>
)

// Keyword token constants. Values start at 700.
const (
	kwSELECT = 700 + iota
	kwINSERT
	kwUPDATE
	kwDELETE
	kwFROM
	kwWHERE
	kwSET
	kwINTO
	kwVALUES
	kwCREATE
	kwALTER
	kwDROP
	kwTABLE
	kwINDEX
	kwVIEW
	kwDATABASE
	kwSCHEMA
	kwIF
	kwNOT
	_ // was kwEXISTS — use kwEXISTS_KW instead (mapped to "exists")
	kwNULL
	kwTRUE
	kwFALSE
	kwAND
	kwOR
	kwIS
	kwIN
	kwBETWEEN
	kwLIKE
	kwREGEXP
	kwRLIKE
	kwCASE
	kwWHEN
	kwTHEN
	kwELSE
	kwEND
	kwAS
	kwON
	kwUSING
	kwJOIN
	kwINNER
	kwLEFT
	kwRIGHT
	kwCROSS
	kwNATURAL
	kwOUTER
	kwFULL
	kwORDER
	kwBY
	kwGROUP
	kwHAVING
	kwLIMIT
	kwOFFSET
	kwUNION
	kwINTERSECT
	kwEXCEPT
	kwALL
	kwDISTINCT
	kwDISTINCTROW
	kwASC
	kwDESC
	kwNULLS
	kwFIRST
	kwLAST
	kwFOR
	kwSHARE
	kwLOCK
	kwNOWAIT
	kwSKIP
	kwLOCKED
	kwPRIMARY
	kwKEY
	kwUNIQUE
	kwCHECK
	kwCONSTRAINT
	kwREFERENCES
	kwFOREIGN
	kwDEFAULT
	kwAUTO_INCREMENT
	kwCOMMENT
	kwCOLUMN
	kwADD
	kwMODIFY
	kwCHANGE
	kwRENAME
	kwTO
	kwTRUNCATE
	kwTEMPORARY
	kwCASCADE
	kwRESTRICT
	kwENGINE
	kwCHARSET
	kwCHARACTER
	kwCOLLATE
	kwINT
	kwINTEGER
	kwSMALLINT
	kwTINYINT
	kwMEDIUMINT
	kwBIGINT
	kwFLOAT
	kwDOUBLE
	kwDECIMAL
	kwNUMERIC
	kwVARCHAR
	kwCHAR
	kwTEXT
	kwTINYTEXT
	kwMEDIUMTEXT
	kwLONGTEXT
	kwBLOB
	kwTINYBLOB
	kwMEDIUMBLOB
	kwLONGBLOB
	kwDATE
	kwDATETIME
	kwTIMESTAMP
	kwTIME
	kwYEAR
	kwBOOL
	kwBOOLEAN
	kwENUM
	kwJSON
	kwUNSIGNED
	kwZEROFILL
	kwBINARY
	kwVARBINARY
	kwBIT
	kwFULLTEXT
	kwSPATIAL
	kwBTREE
	kwHASH
	kwCAST
	kwCONVERT
	kwINTERVAL
	kwCOALESCE
	kwNULLIF
	kwGREATEST
	kwLEAST
	_ // was kwIF_FUNC — unused, IF is handled via kwIF contextually
	kwCONCAT
	kwSUBSTRING
	kwTRIM
	kwPOSITION
	kwEXTRACT
	kwCOUNT
	kwMAX
	kwMIN
	kwSUM
	kwAVG
	kwGROUP_CONCAT
	kwOVER
	kwPARTITION
	kwROW
	kwROWS
	kwRANGE
	kwGROUPS
	kwUNBOUNDED
	kwPRECEDING
	kwFOLLOWING
	kwCURRENT
	kwWINDOW
	kwSOUNDS
	kwREPLACE
	kwIGNORE
	kwDUPLICATE
	kwLOW_PRIORITY
	kwDELAYED
	kwHIGH_PRIORITY
	kwSTRAIGHT_JOIN
	kwSQL_CALC_FOUND_ROWS
	kwDIV
	kwMOD
	kwXOR
	kwCURRENT_DATE
	kwCURRENT_TIME
	kwCURRENT_TIMESTAMP
	kwCURRENT_USER
	kwLOCALTIME
	kwLOCALTIMESTAMP
	kwUSE
	kwSHOW
	kwDESCRIBE
	kwEXPLAIN
	kwBEGIN
	kwCOMMIT
	kwROLLBACK
	kwSAVEPOINT
	kwSTART
	kwTRANSACTION
	kwGRANT
	kwREVOKE
	kwFUNCTION
	kwPROCEDURE
	kwTRIGGER
	kwEVENT
	kwLOAD
	kwDATA
	kwINFILE
	kwPREPARE
	kwEXECUTE
	kwDEALLOCATE
	kwANALYZE
	kwOPTIMIZE
	kwFLUSH
	kwRESET
	kwKILL
	kwDO
	kwESCAPE
	kwMATCH
	kwAGAINST
	kwEXISTS_KW
	kwOUTFILE
	kwDUMPFILE
	kwLINES
	kwFIELDS
	kwTERMINATED
	kwENCLOSED
	kwESCAPED
	kwSTARTING
	kwOPTIONALLY
	kwLOCAL
	kwGLOBAL
	kwSESSION
	kwREAD
	kwWRITE
	kwONLY
	kwREPEATABLE
	kwCOMMITTED
	kwUNCOMMITTED
	kwSERIALIZABLE
	kwISOLATION
	kwLEVEL
	kwUNLOCK
	kwTABLES
	kwREPAIR
	kwQUICK
	kwEXTENDED
	kwWITH
	kwROLLUP
	kwDATABASES
	kwCOLUMNS
	kwSTATUS
	kwVARIABLES
	kwWARNINGS
	kwERRORS
	kwPROCESSLIST
	_ // was kwCREATES — unused, not mapped in keywords
	kwAFTER
	kwBEFORE
	kwEACH
	kwFOLLOWS
	kwPRECEDES
	kwDETERMINISTIC
	kwCONTAINS
	kwSQL
	kwNO
	kwMODIFIES
	kwREADS
	kwRETURNS
	kwLANGUAGE
	kwOUT
	kwINOUT
	kwAT
	kwDEFINER
	kwINVOKER
	kwSECURITY
	kwAGGREGATE
	kwALGORITHM
	kwUNDEFINED
	kwMERGE
	kwTEMPTABLE
	kwCASCADED
	_ // was kwCHECKED — unused, not mapped in keywords
	kwUSER
	kwIDENTIFIED
	kwPASSWORD
	kwNONE
	kwROLE
	_ // was kwSEP — unused duplicate, use kwSEPARATOR
	kwSEPARATOR
	kwBOTH
	kwLEADING
	kwTRAILING
	kwOVERLAY
	kwPLACING
	kwSTORED
	kwVIRTUAL
	kwGENERATED
	kwALWAYS
	_ // was kwPARTITIONING — unused, not mapped in keywords
	kwLINEAR
	kwLIST
	kwSUBPARTITION
	kwFIXED
	kwDYNAMIC
	kwCOMPRESSED
	kwREDUNDANT
	kwCOMPACT
	kwROW_FORMAT
	kwHANDLER
	kwOPEN
	kwCLOSE
	kwNEXT
	kwPREV
	kwQUERY
	kwCONNECTION
	kwCOLUMN_FORMAT
	kwSTORAGE
	kwDISK
	kwMEMORY
	kwBINLOG
	kwMASTER
	kwSLAVE
	kwCHAIN
	kwRELEASE
	kwCONSISTENT
	kwSNAPSHOT
	_ // was kwON_KW — duplicate of kwON
	kwSIGNAL
	kwRESIGNAL
	kwGET
	kwDIAGNOSTICS
	kwCONDITION
	kwFORCE
	_ // was kwBY_KW — duplicate of kwBY
	kwTYPE
	kwRECURSIVE
	kwMEMBER
	kwJSON_TABLE
	kwORDINALITY
	kwNESTED
	kwPATH
	kwEMPTY
	kwERROR_KW
	kwXA
	kwSUSPEND
	kwMIGRATE
	kwPHASE
	kwRECOVER
	kwRESUME
	kwONE
	kwCALL
	kwDECLARE
	kwCURSOR
	kwCONTINUE
	kwEXIT
	kwUNDO
	kwFOUND
	kwENGINES
	kwPLUGINS
	kwREPLICA
	kwPRIVILEGES
	kwPROFILES
	kwRELAYLOG
	kwCOLLATION
	kwLOGS
	kwELSEIF
	kwWHILE
	kwREPEAT
	kwUNTIL
	kwLOOP
	kwLEAVE
	kwITERATE
	kwRETURN
	kwFETCH
	kwCHECKSUM
	kwSHUTDOWN
	kwRESTART
	kwCLONE
	kwINSTANCE
	kwDIRECTORY
	kwREQUIRE
	kwSSL
	kwINSTALL
	kwUNINSTALL
	kwPLUGIN
	kwCOMPONENT
	kwSONAME
	kwTABLESPACE
	kwSERVER
	kwDATAFILE
	kwWRAPPER
	kwOPTIONS
	kwENCRYPTION
	kwLOGFILE
	kwRESOURCE
	kwENABLE
	kwDISABLE
	kwLATERAL
	kwREPLICATION
	kwSOURCE
	kwFILTER
	kwCHANNEL
	kwPURGE
	kwSTOP
	kwIMPORT
	kwPERSIST
	kwBACKUP
	kwHELP
	kwCACHE
	kwREORGANIZE
	kwEXCHANGE
	kwREBUILD
	kwREMOVE
	kwDISCARD
	kwVALIDATION
	kwWITHOUT
	kwPARTITIONING
	kwVISIBLE
	kwINVISIBLE
	kwKEYS
	kwSQL_SMALL_RESULT
	kwSQL_BIG_RESULT
	kwSQL_BUFFER_RESULT
	kwSQL_NO_CACHE
	kwMODE
	kwEXPANSION
	kwRANDOM
	kwRETAIN
	kwOLD
	kwREAL
	kwDEC
	kwACCESSIBLE
	kwASENSITIVE
	kwCUBE
	kwCUME_DIST
	kwDENSE_RANK
	kwDUAL
	kwFIRST_VALUE
	kwGROUPING
	kwINSENSITIVE
	kwLAG
	kwLAST_VALUE
	kwLEAD
	kwNTH_VALUE
	kwNTILE
	kwOF
	kwOPTIMIZER_COSTS
	kwPERCENT_RANK
	kwRANK
	kwROW_NUMBER
	kwSENSITIVE
	kwSPECIFIC
	kwUSAGE
	kwVARYING
	kwDAY_HOUR
	kwDAY_MICROSECOND
	kwDAY_MINUTE
	kwDAY_SECOND
	kwHOUR_MICROSECOND
	kwHOUR_MINUTE
	kwHOUR_SECOND
	kwMINUTE_MICROSECOND
	kwMINUTE_SECOND
	kwSECOND_MICROSECOND
	kwYEAR_MONTH
	kwUTC_DATE
	kwUTC_TIME
	kwUTC_TIMESTAMP
	kwMAXVALUE
	kwNO_WRITE_TO_BINLOG
	kwIO_AFTER_GTIDS
	kwIO_BEFORE_GTIDS
	kwSQLEXCEPTION
	kwSQLSTATE
	kwSQLWARNING
	kwGEOMETRY
	kwPOINT
	kwLINESTRING
	kwPOLYGON
	kwMULTIPOINT
	kwMULTILINESTRING
	kwMULTIPOLYGON
	kwGEOMETRYCOLLECTION
	kwSERIAL
	kwNATIONAL
	kwNCHAR
	kwNVARCHAR
	kwSIGNED
	kwPRECISION
	kwSRID
	kwENFORCED
	kwLESS
	kwTHAN
	kwSUBPARTITIONS
	kwLEAVES
	kwPARSER
	kwCOMPRESSION
	kwINSERT_METHOD
	kwACTION
	kwPARTIAL
	kwFORMAT
	kwXML
	kwCONCURRENT
	kwWORK
	kwXID
	kwEXPORT
	kwUPGRADE
	kwFAST
	kwMEDIUM
	kwCHANGED
	kwCODE
	kwEVENTS
	kwINDEXES
	kwGRANTS
	kwTRIGGERS
	kwSCHEMAS
	kwPARTITIONS
	kwHOSTS
	kwMUTEX
	kwPROFILE
	kwREPLICAS
	kwNAMES
	kwACCOUNT
	kwOPTION
	kwPROXY
	kwROUTINE
	kwEXPIRE
	kwNEVER
	kwDAY
	kwHISTORY
	kwREUSE
	kwOPTIONAL
	kwX509
	kwISSUER
	kwSUBJECT
	kwCIPHER
	kwSCHEDULE
	kwCOMPLETION
	kwPRESERVE
	kwEVERY
	kwSTARTS
	kwENDS
	kwVALUE
	kwSTACKED
	kwUNKNOWN
	kwWAIT
	kwACTIVE
	kwINACTIVE
	kwATTRIBUTE
	kwADMIN
	kwDESCRIPTION
	kwORGANIZATION
	kwREFERENCE
	kwDEFINITION
	kwNAME
	kwSYSTEM
	kwROTATE
	kwKEYRING
	kwTLS
	kwSTREAM
	kwGENERATE
	kwPROCESS
	kwRELOAD
)

// keywords maps lowercase keyword strings to their token types.
var keywords = map[string]int{
	"select":              kwSELECT,
	"insert":              kwINSERT,
	"update":              kwUPDATE,
	"delete":              kwDELETE,
	"from":                kwFROM,
	"where":               kwWHERE,
	"set":                 kwSET,
	"into":                kwINTO,
	"values":              kwVALUES,
	"create":              kwCREATE,
	"alter":               kwALTER,
	"drop":                kwDROP,
	"table":               kwTABLE,
	"index":               kwINDEX,
	"view":                kwVIEW,
	"database":            kwDATABASE,
	"schema":              kwSCHEMA,
	"if":                  kwIF,
	"not":                 kwNOT,
	"exists":              kwEXISTS_KW,
	"null":                kwNULL,
	"true":                kwTRUE,
	"false":               kwFALSE,
	"and":                 kwAND,
	"or":                  kwOR,
	"is":                  kwIS,
	"in":                  kwIN,
	"between":             kwBETWEEN,
	"like":                kwLIKE,
	"regexp":              kwREGEXP,
	"rlike":               kwRLIKE,
	"case":                kwCASE,
	"when":                kwWHEN,
	"then":                kwTHEN,
	"else":                kwELSE,
	"end":                 kwEND,
	"as":                  kwAS,
	"on":                  kwON,
	"using":               kwUSING,
	"join":                kwJOIN,
	"inner":               kwINNER,
	"left":                kwLEFT,
	"right":               kwRIGHT,
	"cross":               kwCROSS,
	"natural":             kwNATURAL,
	"outer":               kwOUTER,
	"full":                kwFULL,
	"order":               kwORDER,
	"by":                  kwBY,
	"group":               kwGROUP,
	"having":              kwHAVING,
	"limit":               kwLIMIT,
	"offset":              kwOFFSET,
	"union":               kwUNION,
	"intersect":           kwINTERSECT,
	"except":              kwEXCEPT,
	"all":                 kwALL,
	"distinct":            kwDISTINCT,
	"distinctrow":         kwDISTINCTROW,
	"asc":                 kwASC,
	"desc":                kwDESC,
	"nulls":               kwNULLS,
	"first":               kwFIRST,
	"last":                kwLAST,
	"for":                 kwFOR,
	"share":               kwSHARE,
	"lock":                kwLOCK,
	"nowait":              kwNOWAIT,
	"skip":                kwSKIP,
	"locked":              kwLOCKED,
	"primary":             kwPRIMARY,
	"key":                 kwKEY,
	"unique":              kwUNIQUE,
	"check":               kwCHECK,
	"constraint":          kwCONSTRAINT,
	"references":          kwREFERENCES,
	"foreign":             kwFOREIGN,
	"default":             kwDEFAULT,
	"auto_increment":      kwAUTO_INCREMENT,
	"comment":             kwCOMMENT,
	"column":              kwCOLUMN,
	"add":                 kwADD,
	"modify":              kwMODIFY,
	"change":              kwCHANGE,
	"rename":              kwRENAME,
	"to":                  kwTO,
	"truncate":            kwTRUNCATE,
	"temporary":           kwTEMPORARY,
	"cascade":             kwCASCADE,
	"restrict":            kwRESTRICT,
	"engine":              kwENGINE,
	"charset":             kwCHARSET,
	"character":           kwCHARACTER,
	"collate":             kwCOLLATE,
	"int":                 kwINT,
	"integer":             kwINTEGER,
	"smallint":            kwSMALLINT,
	"tinyint":             kwTINYINT,
	"mediumint":           kwMEDIUMINT,
	"bigint":              kwBIGINT,
	"float":               kwFLOAT,
	"double":              kwDOUBLE,
	"real":                kwREAL,
	"decimal":             kwDECIMAL,
	"numeric":             kwNUMERIC,
	"dec":                 kwDEC,
	"varchar":             kwVARCHAR,
	"char":                kwCHAR,
	"text":                kwTEXT,
	"tinytext":            kwTINYTEXT,
	"mediumtext":          kwMEDIUMTEXT,
	"longtext":            kwLONGTEXT,
	"blob":                kwBLOB,
	"tinyblob":            kwTINYBLOB,
	"mediumblob":          kwMEDIUMBLOB,
	"longblob":            kwLONGBLOB,
	"date":                kwDATE,
	"datetime":            kwDATETIME,
	"timestamp":           kwTIMESTAMP,
	"time":                kwTIME,
	"year":                kwYEAR,
	"bool":                kwBOOL,
	"boolean":             kwBOOLEAN,
	"enum":                kwENUM,
	"json":                kwJSON,
	"unsigned":            kwUNSIGNED,
	"zerofill":            kwZEROFILL,
	"binary":              kwBINARY,
	"varbinary":           kwVARBINARY,
	"bit":                 kwBIT,
	"fulltext":            kwFULLTEXT,
	"spatial":             kwSPATIAL,
	"btree":               kwBTREE,
	"hash":                kwHASH,
	"cast":                kwCAST,
	"convert":             kwCONVERT,
	"interval":            kwINTERVAL,
	"coalesce":            kwCOALESCE,
	"nullif":              kwNULLIF,
	"greatest":            kwGREATEST,
	"least":               kwLEAST,
	"concat":              kwCONCAT,
	"substring":           kwSUBSTRING,
	"trim":                kwTRIM,
	"position":            kwPOSITION,
	"extract":             kwEXTRACT,
	"count":               kwCOUNT,
	"max":                 kwMAX,
	"min":                 kwMIN,
	"sum":                 kwSUM,
	"avg":                 kwAVG,
	"group_concat":        kwGROUP_CONCAT,
	"over":                kwOVER,
	"partition":           kwPARTITION,
	"row":                 kwROW,
	"rows":                kwROWS,
	"range":               kwRANGE,
	"groups":              kwGROUPS,
	"unbounded":           kwUNBOUNDED,
	"preceding":           kwPRECEDING,
	"following":           kwFOLLOWING,
	"current":             kwCURRENT,
	"window":              kwWINDOW,
	"sounds":              kwSOUNDS,
	"replace":             kwREPLACE,
	"ignore":              kwIGNORE,
	"duplicate":           kwDUPLICATE,
	"low_priority":        kwLOW_PRIORITY,
	"delayed":             kwDELAYED,
	"high_priority":       kwHIGH_PRIORITY,
	"straight_join":       kwSTRAIGHT_JOIN,
	"sql_calc_found_rows": kwSQL_CALC_FOUND_ROWS,
	"div":                 kwDIV,
	"mod":                 kwMOD,
	"xor":                 kwXOR,
	"current_date":        kwCURRENT_DATE,
	"current_time":        kwCURRENT_TIME,
	"current_timestamp":   kwCURRENT_TIMESTAMP,
	"current_user":        kwCURRENT_USER,
	"localtime":           kwLOCALTIME,
	"localtimestamp":      kwLOCALTIMESTAMP,
	"use":                 kwUSE,
	"show":                kwSHOW,
	"describe":            kwDESCRIBE,
	"explain":             kwEXPLAIN,
	"begin":               kwBEGIN,
	"commit":              kwCOMMIT,
	"rollback":            kwROLLBACK,
	"savepoint":           kwSAVEPOINT,
	"start":               kwSTART,
	"transaction":         kwTRANSACTION,
	"grant":               kwGRANT,
	"revoke":              kwREVOKE,
	"function":            kwFUNCTION,
	"procedure":           kwPROCEDURE,
	"trigger":             kwTRIGGER,
	"event":               kwEVENT,
	"load":                kwLOAD,
	"data":                kwDATA,
	"infile":              kwINFILE,
	"prepare":             kwPREPARE,
	"execute":             kwEXECUTE,
	"deallocate":          kwDEALLOCATE,
	"analyze":             kwANALYZE,
	"optimize":            kwOPTIMIZE,
	"flush":               kwFLUSH,
	"reset":               kwRESET,
	"kill":                kwKILL,
	"do":                  kwDO,
	"escape":              kwESCAPE,
	"match":               kwMATCH,
	"against":             kwAGAINST,
	"outfile":             kwOUTFILE,
	"dumpfile":            kwDUMPFILE,
	"lines":               kwLINES,
	"fields":              kwFIELDS,
	"terminated":          kwTERMINATED,
	"enclosed":            kwENCLOSED,
	"escaped":             kwESCAPED,
	"starting":            kwSTARTING,
	"optionally":          kwOPTIONALLY,
	"local":               kwLOCAL,
	"global":              kwGLOBAL,
	"session":             kwSESSION,
	"read":                kwREAD,
	"write":               kwWRITE,
	"only":                kwONLY,
	"repeatable":          kwREPEATABLE,
	"committed":           kwCOMMITTED,
	"uncommitted":         kwUNCOMMITTED,
	"serializable":        kwSERIALIZABLE,
	"isolation":           kwISOLATION,
	"level":               kwLEVEL,
	"unlock":              kwUNLOCK,
	"tables":              kwTABLES,
	"repair":              kwREPAIR,
	"quick":               kwQUICK,
	"extended":            kwEXTENDED,
	"with":                kwWITH,
	"rollup":              kwROLLUP,
	"databases":           kwDATABASES,
	"columns":             kwCOLUMNS,
	"status":              kwSTATUS,
	"variables":           kwVARIABLES,
	"warnings":            kwWARNINGS,
	"errors":              kwERRORS,
	"processlist":         kwPROCESSLIST,
	"after":               kwAFTER,
	"before":              kwBEFORE,
	"each":                kwEACH,
	"follows":             kwFOLLOWS,
	"precedes":            kwPRECEDES,
	"deterministic":       kwDETERMINISTIC,
	"contains":            kwCONTAINS,
	"sql":                 kwSQL,
	"no":                  kwNO,
	"modifies":            kwMODIFIES,
	"reads":               kwREADS,
	"returns":             kwRETURNS,
	"language":            kwLANGUAGE,
	"lateral":             kwLATERAL,
	"replication":         kwREPLICATION,
	"source":              kwSOURCE,
	"filter":              kwFILTER,
	"channel":             kwCHANNEL,
	"purge":               kwPURGE,
	"stop":                kwSTOP,
	"import":              kwIMPORT,
	"persist":             kwPERSIST,
	"backup":              kwBACKUP,
	"help":                kwHELP,
	"cache":               kwCACHE,
	"out":                 kwOUT,
	"inout":               kwINOUT,
	"at":                  kwAT,
	"definer":             kwDEFINER,
	"invoker":             kwINVOKER,
	"security":            kwSECURITY,
	"aggregate":           kwAGGREGATE,
	"algorithm":           kwALGORITHM,
	"undefined":           kwUNDEFINED,
	"merge":               kwMERGE,
	"temptable":           kwTEMPTABLE,
	"cascaded":            kwCASCADED,
	"user":                kwUSER,
	"identified":          kwIDENTIFIED,
	"password":            kwPASSWORD,
	"none":                kwNONE,
	"role":                kwROLE,
	"separator":           kwSEPARATOR,
	"both":                kwBOTH,
	"leading":             kwLEADING,
	"trailing":            kwTRAILING,
	"overlay":             kwOVERLAY,
	"placing":             kwPLACING,
	"stored":              kwSTORED,
	"virtual":             kwVIRTUAL,
	"generated":           kwGENERATED,
	"always":              kwALWAYS,
	"linear":              kwLINEAR,
	"list":                kwLIST,
	"subpartition":        kwSUBPARTITION,
	"fixed":               kwFIXED,
	"dynamic":             kwDYNAMIC,
	"compressed":          kwCOMPRESSED,
	"redundant":           kwREDUNDANT,
	"compact":             kwCOMPACT,
	"row_format":          kwROW_FORMAT,
	"column_format":       kwCOLUMN_FORMAT,
	"storage":             kwSTORAGE,
	"disk":                kwDISK,
	"memory":              kwMEMORY,
	"handler":             kwHANDLER,
	"open":                kwOPEN,
	"close":               kwCLOSE,
	"next":                kwNEXT,
	"prev":                kwPREV,
	"query":               kwQUERY,
	"connection":          kwCONNECTION,
	"binlog":              kwBINLOG,
	"master":              kwMASTER,
	"slave":               kwSLAVE,
	"chain":               kwCHAIN,
	"release":             kwRELEASE,
	"consistent":          kwCONSISTENT,
	"snapshot":            kwSNAPSHOT,
	"signal":              kwSIGNAL,
	"resignal":            kwRESIGNAL,
	"get":                 kwGET,
	"diagnostics":         kwDIAGNOSTICS,
	"condition":           kwCONDITION,
	"force":               kwFORCE,
	"type":                kwTYPE,
	"recursive":           kwRECURSIVE,
	"member":              kwMEMBER,
	"json_table":          kwJSON_TABLE,
	"ordinality":          kwORDINALITY,
	"nested":              kwNESTED,
	"path":                kwPATH,
	"empty":               kwEMPTY,
	"error":               kwERROR_KW,
	"xa":                  kwXA,
	"suspend":             kwSUSPEND,
	"migrate":             kwMIGRATE,
	"phase":               kwPHASE,
	"recover":             kwRECOVER,
	"resume":              kwRESUME,
	"one":                 kwONE,
	"call":                kwCALL,
	"declare":             kwDECLARE,
	"cursor":              kwCURSOR,
	"continue":            kwCONTINUE,
	"exit":                kwEXIT,
	"undo":                kwUNDO,
	"found":               kwFOUND,
	"engines":             kwENGINES,
	"plugins":             kwPLUGINS,
	"replica":             kwREPLICA,
	"privileges":          kwPRIVILEGES,
	"profiles":            kwPROFILES,
	"relaylog":            kwRELAYLOG,
	"collation":           kwCOLLATION,
	"logs":                kwLOGS,
	"elseif":              kwELSEIF,
	"while":               kwWHILE,
	"repeat":              kwREPEAT,
	"until":               kwUNTIL,
	"loop":                kwLOOP,
	"leave":               kwLEAVE,
	"iterate":             kwITERATE,
	"return":              kwRETURN,
	"fetch":               kwFETCH,
	"clone":               kwCLONE,
	"instance":            kwINSTANCE,
	"directory":           kwDIRECTORY,
	"require":             kwREQUIRE,
	"ssl":                 kwSSL,
	"install":             kwINSTALL,
	"uninstall":           kwUNINSTALL,
	"plugin":              kwPLUGIN,
	"component":           kwCOMPONENT,
	"soname":              kwSONAME,
	"checksum":            kwCHECKSUM,
	"shutdown":            kwSHUTDOWN,
	"restart":             kwRESTART,
	"tablespace":          kwTABLESPACE,
	"server":              kwSERVER,
	"datafile":            kwDATAFILE,
	"wrapper":             kwWRAPPER,
	"options":             kwOPTIONS,
	"encryption":          kwENCRYPTION,
	"logfile":             kwLOGFILE,
	"resource":            kwRESOURCE,
	"enable":              kwENABLE,
	"disable":             kwDISABLE,
	"reorganize":          kwREORGANIZE,
	"exchange":            kwEXCHANGE,
	"rebuild":             kwREBUILD,
	"remove":              kwREMOVE,
	"discard":             kwDISCARD,
	"validation":          kwVALIDATION,
	"without":             kwWITHOUT,
	"partitioning":        kwPARTITIONING,
	"visible":             kwVISIBLE,
	"invisible":           kwINVISIBLE,
	"keys":                kwKEYS,
	"sql_small_result":    kwSQL_SMALL_RESULT,
	"sql_big_result":      kwSQL_BIG_RESULT,
	"sql_buffer_result":   kwSQL_BUFFER_RESULT,
	"sql_no_cache":        kwSQL_NO_CACHE,
	"mode":                kwMODE,
	"expansion":           kwEXPANSION,
	"random":              kwRANDOM,
	"retain":              kwRETAIN,
	"old":                 kwOLD,
	"accessible":          kwACCESSIBLE,
	"asensitive":          kwASENSITIVE,
	"cube":                kwCUBE,
	"cume_dist":           kwCUME_DIST,
	"dense_rank":          kwDENSE_RANK,
	"dual":                kwDUAL,
	"first_value":         kwFIRST_VALUE,
	"grouping":            kwGROUPING,
	"insensitive":         kwINSENSITIVE,
	"lag":                 kwLAG,
	"last_value":          kwLAST_VALUE,
	"lead":                kwLEAD,
	"nth_value":           kwNTH_VALUE,
	"ntile":               kwNTILE,
	"of":                  kwOF,
	"optimizer_costs":     kwOPTIMIZER_COSTS,
	"percent_rank":        kwPERCENT_RANK,
	"rank":                kwRANK,
	"row_number":          kwROW_NUMBER,
	"sensitive":           kwSENSITIVE,
	"specific":            kwSPECIFIC,
	"usage":               kwUSAGE,
	"varying":             kwVARYING,
	"day_hour":            kwDAY_HOUR,
	"day_microsecond":     kwDAY_MICROSECOND,
	"day_minute":          kwDAY_MINUTE,
	"day_second":          kwDAY_SECOND,
	"hour_microsecond":    kwHOUR_MICROSECOND,
	"hour_minute":         kwHOUR_MINUTE,
	"hour_second":         kwHOUR_SECOND,
	"minute_microsecond":  kwMINUTE_MICROSECOND,
	"minute_second":       kwMINUTE_SECOND,
	"second_microsecond":  kwSECOND_MICROSECOND,
	"year_month":          kwYEAR_MONTH,
	"utc_date":            kwUTC_DATE,
	"utc_time":            kwUTC_TIME,
	"utc_timestamp":       kwUTC_TIMESTAMP,
	"maxvalue":            kwMAXVALUE,
	"no_write_to_binlog":  kwNO_WRITE_TO_BINLOG,
	"io_after_gtids":      kwIO_AFTER_GTIDS,
	"io_before_gtids":     kwIO_BEFORE_GTIDS,
	"sqlexception":        kwSQLEXCEPTION,
	"sqlstate":            kwSQLSTATE,
	"sqlwarning":          kwSQLWARNING,
	"geometry":            kwGEOMETRY,
	"point":               kwPOINT,
	"linestring":          kwLINESTRING,
	"polygon":             kwPOLYGON,
	"multipoint":          kwMULTIPOINT,
	"multilinestring":     kwMULTILINESTRING,
	"multipolygon":        kwMULTIPOLYGON,
	"geometrycollection":  kwGEOMETRYCOLLECTION,
	"serial":              kwSERIAL,
	"national":            kwNATIONAL,
	"nchar":               kwNCHAR,
	"nvarchar":            kwNVARCHAR,
	"signed":              kwSIGNED,
	"precision":           kwPRECISION,
	"srid":                kwSRID,
	"enforced":            kwENFORCED,
	"less":                kwLESS,
	"than":                kwTHAN,
	"subpartitions":       kwSUBPARTITIONS,
	"leaves":              kwLEAVES,
	"parser":              kwPARSER,
	"compression":         kwCOMPRESSION,
	"insert_method":       kwINSERT_METHOD,
	"action":              kwACTION,
	"partial":             kwPARTIAL,
	"format":              kwFORMAT,
	"xml":                 kwXML,
	"concurrent":          kwCONCURRENT,
	"work":                kwWORK,
	"xid":                 kwXID,
	"export":              kwEXPORT,
	"upgrade":             kwUPGRADE,
	"fast":                kwFAST,
	"medium":              kwMEDIUM,
	"changed":             kwCHANGED,
	"code":                kwCODE,
	"events":              kwEVENTS,
	"indexes":             kwINDEXES,
	"grants":              kwGRANTS,
	"triggers":            kwTRIGGERS,
	"schemas":             kwSCHEMAS,
	"partitions":          kwPARTITIONS,
	"hosts":               kwHOSTS,
	"mutex":               kwMUTEX,
	"profile":             kwPROFILE,
	"replicas":            kwREPLICAS,
	"names":               kwNAMES,
	"account":             kwACCOUNT,
	"option":              kwOPTION,
	"proxy":               kwPROXY,
	"routine":             kwROUTINE,
	"expire":              kwEXPIRE,
	"never":               kwNEVER,
	"day":                 kwDAY,
	"history":             kwHISTORY,
	"reuse":               kwREUSE,
	"optional":            kwOPTIONAL,
	"x509":                kwX509,
	"issuer":              kwISSUER,
	"subject":             kwSUBJECT,
	"cipher":              kwCIPHER,
	"schedule":            kwSCHEDULE,
	"completion":          kwCOMPLETION,
	"preserve":            kwPRESERVE,
	"every":               kwEVERY,
	"starts":              kwSTARTS,
	"ends":                kwENDS,
	"value":               kwVALUE,
	"stacked":             kwSTACKED,
	"unknown":             kwUNKNOWN,
	"wait":                kwWAIT,
	"active":              kwACTIVE,
	"inactive":            kwINACTIVE,
	"attribute":           kwATTRIBUTE,
	"admin":               kwADMIN,
	"description":         kwDESCRIPTION,
	"organization":        kwORGANIZATION,
	"reference":           kwREFERENCE,
	"definition":          kwDEFINITION,
	"name":                kwNAME,
	"system":              kwSYSTEM,
	"rotate":              kwROTATE,
	"keyring":             kwKEYRING,
	"tls":                 kwTLS,
	"stream":              kwSTREAM,
	"generate":            kwGENERATE,
	"process":             kwPROCESS,
	"reload":              kwRELOAD,
}

// Token represents a lexical token.
type Token struct {
	Type int    // token type
	Str  string // string value for identifiers, operators, string literals
	Ival int64  // integer value for ICONST
	Loc  int    // byte offset in source text
}

// Lexer implements a MySQL SQL lexer.
type Lexer struct {
	input     string
	pos       int
	start     int
	prevToken int // type of the previously emitted token
}

// NewLexer creates a new MySQL lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	tok := l.nextToken()
	l.prevToken = tok.Type
	return tok
}

func (l *Lexer) nextToken() Token {
	l.skipWhitespaceAndComments()

	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: l.pos}
	}

	l.start = l.pos
	ch := l.input[l.pos]

	// User variable @name or system variable @@name
	if ch == '@' {
		return l.scanVariable()
	}

	// Backtick-quoted identifier
	if ch == '`' {
		return l.scanBacktickIdent()
	}

	// String literals (single or double quoted)
	if ch == '\'' {
		return l.scanString('\'')
	}
	if ch == '"' {
		return l.scanString('"')
	}

	// Hex literal: 0x or X'...'
	if ch == 'X' || ch == 'x' {
		if l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
			l.pos++ // skip X
			tok := l.scanString('\'')
			tok.Type = tokXCONST
			tok.Loc = l.start
			return tok
		}
	}

	// Bit literal: b'...' or B'...'
	if (ch == 'b' || ch == 'B') && l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
		l.pos++ // skip b/B
		tok := l.scanString('\'')
		tok.Type = tokBCONST
		tok.Loc = l.start
		return tok
	}

	// Number: 0x hex, 0b binary, or decimal/float
	if ch >= '0' && ch <= '9' {
		return l.scanNumber()
	}
	// Also handle .N float (e.g., .5)
	if ch == '.' && l.pos+1 < len(l.input) && l.input[l.pos+1] >= '0' && l.input[l.pos+1] <= '9' {
		return l.scanNumber()
	}

	// Identifiers and keywords (including _charset prefix)
	if isIdentStart(ch) {
		return l.scanIdentOrKeyword()
	}

	// Multi-character operators
	if l.pos+1 < len(l.input) {
		next := l.input[l.pos+1]
		switch {
		case ch == '<' && next == '=':
			if l.pos+2 < len(l.input) && l.input[l.pos+2] == '>' {
				l.pos += 3
				return Token{Type: tokNullSafeEq, Str: "<=>", Loc: l.start}
			}
			l.pos += 2
			return Token{Type: tokLessEq, Str: "<=", Loc: l.start}
		case ch == '>' && next == '=':
			l.pos += 2
			return Token{Type: tokGreaterEq, Str: ">=", Loc: l.start}
		case ch == '<' && next == '>':
			l.pos += 2
			return Token{Type: tokNotEq, Str: "<>", Loc: l.start}
		case ch == '!' && next == '=':
			l.pos += 2
			return Token{Type: tokNotEq, Str: "!=", Loc: l.start}
		case ch == '<' && next == '<':
			l.pos += 2
			return Token{Type: tokShiftLeft, Str: "<<", Loc: l.start}
		case ch == '>' && next == '>':
			l.pos += 2
			return Token{Type: tokShiftRight, Str: ">>", Loc: l.start}
		case ch == ':' && next == '=':
			l.pos += 2
			return Token{Type: tokAssign, Str: ":=", Loc: l.start}
		case ch == '-' && next == '>':
			if l.pos+2 < len(l.input) && l.input[l.pos+2] == '>' {
				l.pos += 3
				return Token{Type: tokJsonUnquote, Str: "->>", Loc: l.start}
			}
			l.pos += 2
			return Token{Type: tokJsonExtract, Str: "->", Loc: l.start}
		}
	}

	// Single-character tokens
	l.pos++
	return Token{Type: int(ch), Str: string(ch), Loc: l.start}
}

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		// Whitespace
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.pos++
			continue
		}

		// Line comment: -- must be followed by a space, tab, newline, or end-of-input (per MySQL spec).
		if ch == '-' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '-' {
			// Check third character: must be space, tab, newline, or end of input.
			if l.pos+2 >= len(l.input) || l.input[l.pos+2] == ' ' || l.input[l.pos+2] == '\t' || l.input[l.pos+2] == '\n' || l.input[l.pos+2] == '\r' {
				l.pos += 2
				for l.pos < len(l.input) && l.input[l.pos] != '\n' {
					l.pos++
				}
				continue
			}
			// Not a comment; break and let the main scanner handle '-' as a token.
			break
		}

		// Line comment: #
		if ch == '#' {
			l.pos++
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			continue
		}

		// Block comment: /* ... */
		if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
			// MySQL conditional comments: /*!NNNNN ... */ or /*! ... */
			// These should be parsed as SQL, not skipped.
			if l.pos+2 < len(l.input) && l.input[l.pos+2] == '!' {
				// Skip /*!
				innerStart := l.pos + 3
				// Skip optional version number (digits)
				vpos := innerStart
				for vpos < len(l.input) && l.input[vpos] >= '0' && l.input[vpos] <= '9' {
					vpos++
				}
				// Find the matching */
				end := vpos
				depth := 1
				for end < len(l.input) && depth > 0 {
					if l.input[end] == '*' && end+1 < len(l.input) && l.input[end+1] == '/' {
						depth--
						if depth == 0 {
							break
						}
						end += 2
					} else if l.input[end] == '/' && end+1 < len(l.input) && l.input[end+1] == '*' {
						depth++
						end += 2
					} else {
						end++
					}
				}
				// Extract the inner content and splice it into the input,
				// replacing the conditional comment with its content.
				inner := l.input[vpos:end]
				l.input = l.input[:l.pos] + inner + l.input[end+2:]
				// Don't advance l.pos — re-scan from the start of the inner content.
				continue
			}

			l.pos += 2
			// Regular block comment: skip everything.
			depth := 1
			for l.pos < len(l.input) && depth > 0 {
				if l.input[l.pos] == '*' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
					depth--
					l.pos += 2
				} else if l.input[l.pos] == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
					depth++
					l.pos += 2
				} else {
					l.pos++
				}
			}
			continue
		}

		break
	}
}

func (l *Lexer) scanVariable() Token {
	start := l.pos
	l.pos++ // skip first @
	system := false
	if l.pos < len(l.input) && l.input[l.pos] == '@' {
		system = true
		l.pos++ // skip second @
	}

	// Scan variable name: can be backtick-quoted or identifier chars
	var name string
	if l.pos < len(l.input) && l.input[l.pos] == '`' {
		l.pos++ // skip opening backtick
		nameStart := l.pos
		for l.pos < len(l.input) && l.input[l.pos] != '`' {
			l.pos++
		}
		name = l.input[nameStart:l.pos]
		if l.pos < len(l.input) {
			l.pos++ // skip closing backtick
		}
	} else {
		nameStart := l.pos
		for l.pos < len(l.input) {
			ch := l.input[l.pos]
			if isIdentChar(ch) || ch == '.' {
				l.pos++
			} else {
				break
			}
		}
		name = l.input[nameStart:l.pos]
	}

	prefix := "@"
	if system {
		prefix = "@@"
	}
	return Token{Type: tokIDENT, Str: prefix + name, Loc: start}
}

func (l *Lexer) scanBacktickIdent() Token {
	start := l.pos
	l.pos++ // skip opening backtick
	var sb strings.Builder
	for l.pos < len(l.input) {
		if l.input[l.pos] == '`' {
			// Double backtick is an escape
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '`' {
				sb.WriteByte('`')
				l.pos += 2
			} else {
				l.pos++ // skip closing backtick
				break
			}
		} else {
			sb.WriteByte(l.input[l.pos])
			l.pos++
		}
	}
	return Token{Type: tokIDENT, Str: sb.String(), Loc: start}
}

func (l *Lexer) scanString(quote byte) Token {
	start := l.pos
	l.pos++ // skip opening quote
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == quote {
			// Double quote is an escape
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == quote {
				sb.WriteByte(quote)
				l.pos += 2
			} else {
				l.pos++ // skip closing quote
				break
			}
		} else if ch == '\\' {
			l.pos++
			if l.pos < len(l.input) {
				esc := l.input[l.pos]
				switch esc {
				case 'n':
					sb.WriteByte('\n')
				case 't':
					sb.WriteByte('\t')
				case 'r':
					sb.WriteByte('\r')
				case '0':
					sb.WriteByte(0)
				case '\\':
					sb.WriteByte('\\')
				case '\'':
					sb.WriteByte('\'')
				case '"':
					sb.WriteByte('"')
				case 'b':
					sb.WriteByte(0x08) // backspace
				case 'Z':
					sb.WriteByte(0x1A) // Ctrl-Z
				case '_':
					sb.WriteByte('\\')
					sb.WriteByte('_') // preserve backslash for LIKE patterns
				case '%':
					sb.WriteByte('\\')
					sb.WriteByte('%') // preserve backslash for LIKE patterns
				default:
					sb.WriteByte(esc)
				}
				l.pos++
			}
		} else {
			sb.WriteByte(ch)
			l.pos++
		}
	}
	return Token{Type: tokSCONST, Str: sb.String(), Loc: start}
}

func (l *Lexer) scanNumber() Token {
	start := l.pos

	// Handle 0x hex literals
	if l.input[l.pos] == '0' && l.pos+1 < len(l.input) && (l.input[l.pos+1] == 'x' || l.input[l.pos+1] == 'X') {
		l.pos += 2
		for l.pos < len(l.input) && isHexDigit(l.input[l.pos]) {
			l.pos++
		}
		return Token{Type: tokXCONST, Str: l.input[start:l.pos], Loc: start}
	}

	// Handle 0b binary literals
	if l.input[l.pos] == '0' && l.pos+1 < len(l.input) && (l.input[l.pos+1] == 'b' || l.input[l.pos+1] == 'B') {
		l.pos += 2
		for l.pos < len(l.input) && (l.input[l.pos] == '0' || l.input[l.pos] == '1') {
			l.pos++
		}
		return Token{Type: tokBCONST, Str: l.input[start:l.pos], Loc: start}
	}

	// Integer part (or start of float)
	isFloat := false
	if l.input[l.pos] == '.' {
		isFloat = true
	}
	for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
		l.pos++
	}

	// Decimal point
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	}

	// Exponent
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		isFloat = true
		l.pos++
		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	}

	s := l.input[start:l.pos]
	if isFloat {
		return Token{Type: tokFCONST, Str: s, Loc: start}
	}

	ival, _ := strconv.ParseInt(s, 10, 64)
	return Token{Type: tokICONST, Str: s, Ival: ival, Loc: start}
}

func (l *Lexer) scanIdentOrKeyword() Token {
	start := l.pos
	for l.pos < len(l.input) && isIdentChar(l.input[l.pos]) {
		l.pos++
	}
	word := l.input[start:l.pos]

	// After '.', suppress keyword lookup — MySQL treats any word after dot as
	// a plain identifier (sql_lex.cc MY_LEX_IDENT_START state).
	if l.prevToken == '.' {
		return Token{Type: tokIDENT, Str: word, Loc: start}
	}

	lower := strings.ToLower(word)
	if kwType, ok := keywords[lower]; ok {
		return Token{Type: kwType, Str: word, Loc: start}
	}

	return Token{Type: tokIDENT, Str: word, Loc: start}
}

func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' || ch == '$' || ch > 127
}

func isIdentChar(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9')
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// isIdentRune checks if a rune is a valid identifier character (for multi-byte chars).
func isIdentRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '$'
}

// These functions exist to ensure the imports are used.
var _ = utf8.RuneLen
var _ = isIdentRune
