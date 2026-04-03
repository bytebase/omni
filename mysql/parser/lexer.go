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
	tokAt2         // @@ (system variable prefix)
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
	// --- 253 missing MySQL 8.0 keywords (appended in one pass) ---
	kwABSENT
	kwADDDATE
	kwALLOW_MISSING_FILES
	kwANY
	kwARRAY
	kwASCII
	kwASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS
	kwAUTHENTICATION
	kwAUTO
	kwAUTO_REFRESH
	kwAUTO_REFRESH_SOURCE
	kwAUTOEXTEND_SIZE
	kwAVG_ROW_LENGTH
	kwBERNOULLI
	kwBIT_AND
	kwBIT_OR
	kwBIT_XOR
	kwBLOCK
	kwBUCKETS
	kwBULK
	kwBYTE
	kwCATALOG_NAME
	kwCHALLENGE_RESPONSE
	kwCLASS_ORIGIN
	kwCLIENT
	kwCOLUMN_NAME
	kwCONSTRAINT_CATALOG
	kwCONSTRAINT_NAME
	kwCONSTRAINT_SCHEMA
	kwCONTEXT
	kwCPU
	kwCURDATE
	kwCURSOR_NAME
	kwCURTIME
	kwDATE_ADD
	kwDATE_SUB
	kwDEFAULT_AUTH
	kwDELAY_KEY_WRITE
	kwDUALITY
	kwENGINE_ATTRIBUTE
	kwEXCLUDE
	kwEXTENT_SIZE
	kwEXTERNAL
	kwEXTERNAL_FORMAT
	kwFACTOR
	kwFAILED_LOGIN_ATTEMPTS
	kwFAULTS
	kwFILE
	kwFILE_BLOCK_SIZE
	kwFILE_FORMAT
	kwFILE_NAME
	kwFILE_PATTERN
	kwFILE_PREFIX
	kwFILES
	kwFINISH
	kwFLOAT4
	kwFLOAT8
	kwGENERAL
	kwGEOMCOLLECTION
	kwGET_FORMAT
	kwGET_SOURCE_PUBLIC_KEY
	kwGROUP_REPLICATION
	kwGTID_ONLY
	kwGTIDS
	kwGUIDED
	kwHEADER
	kwHISTOGRAM
	kwHOST
	kwHOUR
	kwIGNORE_SERVER_IDS
	kwINITIAL
	kwINITIAL_SIZE
	kwINITIATE
	kwINT1
	kwINT2
	kwINT3
	kwINT4
	kwINT8
	kwIO
	kwIO_THREAD
	kwIPC
	kwJSON_ARRAYAGG
	kwJSON_DUALITY_OBJECT
	kwJSON_OBJECTAGG
	kwJSON_VALUE
	kwKEY_BLOCK_SIZE
	kwLIBRARY
	kwLOCKS
	kwLOG
	kwLONG
	kwMANUAL
	kwMATERIALIZED
	kwMAX_CONNECTIONS_PER_HOUR
	kwMAX_QUERIES_PER_HOUR
	kwMAX_ROWS
	kwMAX_SIZE
	kwMAX_UPDATES_PER_HOUR
	kwMAX_USER_CONNECTIONS
	kwMESSAGE_TEXT
	kwMICROSECOND
	kwMID
	kwMIDDLEINT
	kwMIN_ROWS
	kwMINUTE
	kwMONTH
	kwMYSQL_ERRNO
	kwNDB
	kwNDBCLUSTER
	kwNETWORK_NAMESPACE
	kwNEW
	kwNO_WAIT
	kwNODEGROUP
	kwNOW
	kwNUMBER
	kwOFF
	kwOJ
	kwOTHERS
	kwOWNER
	kwPACK_KEYS
	kwPAGE
	kwPARALLEL
	kwPARAMETERS
	kwPARSE_TREE
	kwPASSWORD_LOCK_TIME
	kwPERSIST_ONLY
	kwPLUGIN_DIR
	kwPORT
	kwPRIVILEGE_CHECKS_USER
	kwQUALIFY
	kwQUARTER
	kwREAD_ONLY
	kwREAD_WRITE
	kwREDO_BUFFER_SIZE
	kwREGISTRATION
	kwRELATIONAL
	kwRELAY
	kwRELAY_LOG_FILE
	kwRELAY_LOG_POS
	kwRELAY_THREAD
	kwREPLICATE_DO_DB
	kwREPLICATE_DO_TABLE
	kwREPLICATE_IGNORE_DB
	kwREPLICATE_IGNORE_TABLE
	kwREPLICATE_REWRITE_DB
	kwREPLICATE_WILD_DO_TABLE
	kwREPLICATE_WILD_IGNORE_TABLE
	kwREQUIRE_ROW_FORMAT
	kwREQUIRE_TABLE_PRIMARY_KEY_CHECK
	kwRESPECT
	kwRESTORE
	kwRETURNED_SQLSTATE
	kwRETURNING
	kwREVERSE
	kwROW_COUNT
	kwRTREE
	kwS3
	kwSCHEMA_NAME
	kwSECOND
	kwSECONDARY
	kwSECONDARY_ENGINE
	kwSECONDARY_ENGINE_ATTRIBUTE
	kwSECONDARY_LOAD
	kwSECONDARY_UNLOAD
	kwSESSION_USER
	kwSETS
	kwSIMPLE
	kwSLOW
	kwSOCKET
	kwSOME
	kwSOURCE_AUTO_POSITION
	kwSOURCE_BIND
	kwSOURCE_COMPRESSION_ALGORITHMS
	kwSOURCE_CONNECT_RETRY
	kwSOURCE_CONNECTION_AUTO_FAILOVER
	kwSOURCE_DELAY
	kwSOURCE_HEARTBEAT_PERIOD
	kwSOURCE_HOST
	kwSOURCE_LOG_FILE
	kwSOURCE_LOG_POS
	kwSOURCE_PASSWORD
	kwSOURCE_PORT
	kwSOURCE_PUBLIC_KEY_PATH
	kwSOURCE_RETRY_COUNT
	kwSOURCE_SSL
	kwSOURCE_SSL_CA
	kwSOURCE_SSL_CAPATH
	kwSOURCE_SSL_CERT
	kwSOURCE_SSL_CIPHER
	kwSOURCE_SSL_CRL
	kwSOURCE_SSL_CRLPATH
	kwSOURCE_SSL_KEY
	kwSOURCE_SSL_VERIFY_SERVER_CERT
	kwSOURCE_TLS_CIPHERSUITES
	kwSOURCE_TLS_VERSION
	kwSOURCE_USER
	kwSOURCE_ZSTD_COMPRESSION_LEVEL
	kwSQL_AFTER_GTIDS
	kwSQL_AFTER_MTS_GAPS
	kwSQL_BEFORE_GTIDS
	kwSQL_THREAD
	kwSQL_TSI_DAY
	kwSQL_TSI_HOUR
	kwSQL_TSI_MINUTE
	kwSQL_TSI_MONTH
	kwSQL_TSI_QUARTER
	kwSQL_TSI_SECOND
	kwSQL_TSI_WEEK
	kwSQL_TSI_YEAR
	kwST_COLLECT
	kwSTATS_AUTO_RECALC
	kwSTATS_PERSISTENT
	kwSTATS_SAMPLE_PAGES
	kwSTD
	kwSTDDEV
	kwSTDDEV_POP
	kwSTDDEV_SAMP
	kwSTRICT_LOAD
	kwSTRING
	kwSUBCLASS_ORIGIN
	kwSUBDATE
	kwSUBSTR
	kwSUPER
	kwSWAPS
	kwSWITCHES
	kwSYSDATE
	kwSYSTEM_USER
	kwTABLE_CHECKSUM
	kwTABLE_NAME
	kwTABLESAMPLE
	kwTHREAD_PRIORITY
	kwTIES
	kwTIMESTAMPADD
	kwTIMESTAMPDIFF
	kwTYPES
	kwUNDO_BUFFER_SIZE
	kwUNDOFILE
	kwUNICODE
	kwUNREGISTER
	kwURI
	kwURL
	kwUSE_FRM
	kwUSER_RESOURCES
	kwVALIDATE
	kwVAR_POP
	kwVAR_SAMP
	kwVARCHARACTER
	kwVARIANCE
	kwVCPU
	kwVECTOR
	kwVERIFY_KEY_CONSTRAINTS
	kwWEEK
	kwWEIGHT_STRING
	kwZONE
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
	"reload":                                kwRELOAD,
	// --- 253 missing MySQL 8.0 keywords ---
	"absent":                                kwABSENT,
	"adddate":                               kwADDDATE,
	"allow_missing_files":                   kwALLOW_MISSING_FILES,
	"any":                                   kwANY,
	"array":                                 kwARRAY,
	"ascii":                                 kwASCII,
	"assign_gtids_to_anonymous_transactions": kwASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS,
	"authentication":                        kwAUTHENTICATION,
	"auto":                                  kwAUTO,
	"auto_refresh":                          kwAUTO_REFRESH,
	"auto_refresh_source":                   kwAUTO_REFRESH_SOURCE,
	"autoextend_size":                       kwAUTOEXTEND_SIZE,
	"avg_row_length":                        kwAVG_ROW_LENGTH,
	"bernoulli":                             kwBERNOULLI,
	"bit_and":                               kwBIT_AND,
	"bit_or":                                kwBIT_OR,
	"bit_xor":                               kwBIT_XOR,
	"block":                                 kwBLOCK,
	"buckets":                               kwBUCKETS,
	"bulk":                                  kwBULK,
	"byte":                                  kwBYTE,
	"catalog_name":                          kwCATALOG_NAME,
	"challenge_response":                    kwCHALLENGE_RESPONSE,
	"class_origin":                          kwCLASS_ORIGIN,
	"client":                                kwCLIENT,
	"column_name":                           kwCOLUMN_NAME,
	"constraint_catalog":                    kwCONSTRAINT_CATALOG,
	"constraint_name":                       kwCONSTRAINT_NAME,
	"constraint_schema":                     kwCONSTRAINT_SCHEMA,
	"context":                               kwCONTEXT,
	"cpu":                                   kwCPU,
	"curdate":                               kwCURDATE,
	"cursor_name":                           kwCURSOR_NAME,
	"curtime":                               kwCURTIME,
	"date_add":                              kwDATE_ADD,
	"date_sub":                              kwDATE_SUB,
	"default_auth":                          kwDEFAULT_AUTH,
	"delay_key_write":                       kwDELAY_KEY_WRITE,
	"duality":                               kwDUALITY,
	"engine_attribute":                      kwENGINE_ATTRIBUTE,
	"exclude":                               kwEXCLUDE,
	"extent_size":                           kwEXTENT_SIZE,
	"external":                              kwEXTERNAL,
	"external_format":                       kwEXTERNAL_FORMAT,
	"factor":                                kwFACTOR,
	"failed_login_attempts":                 kwFAILED_LOGIN_ATTEMPTS,
	"faults":                                kwFAULTS,
	"file":                                  kwFILE,
	"file_block_size":                       kwFILE_BLOCK_SIZE,
	"file_format":                           kwFILE_FORMAT,
	"file_name":                             kwFILE_NAME,
	"file_pattern":                          kwFILE_PATTERN,
	"file_prefix":                           kwFILE_PREFIX,
	"files":                                 kwFILES,
	"finish":                                kwFINISH,
	"float4":                                kwFLOAT4,
	"float8":                                kwFLOAT8,
	"general":                               kwGENERAL,
	"geomcollection":                        kwGEOMCOLLECTION,
	"get_format":                            kwGET_FORMAT,
	"get_source_public_key":                 kwGET_SOURCE_PUBLIC_KEY,
	"group_replication":                     kwGROUP_REPLICATION,
	"gtid_only":                             kwGTID_ONLY,
	"gtids":                                 kwGTIDS,
	"guided":                                kwGUIDED,
	"header":                                kwHEADER,
	"histogram":                             kwHISTOGRAM,
	"host":                                  kwHOST,
	"hour":                                  kwHOUR,
	"ignore_server_ids":                     kwIGNORE_SERVER_IDS,
	"initial":                               kwINITIAL,
	"initial_size":                          kwINITIAL_SIZE,
	"initiate":                              kwINITIATE,
	"int1":                                  kwINT1,
	"int2":                                  kwINT2,
	"int3":                                  kwINT3,
	"int4":                                  kwINT4,
	"int8":                                  kwINT8,
	"io":                                    kwIO,
	"io_thread":                             kwIO_THREAD,
	"ipc":                                   kwIPC,
	"json_arrayagg":                         kwJSON_ARRAYAGG,
	"json_duality_object":                   kwJSON_DUALITY_OBJECT,
	"json_objectagg":                        kwJSON_OBJECTAGG,
	"json_value":                            kwJSON_VALUE,
	"key_block_size":                        kwKEY_BLOCK_SIZE,
	"library":                               kwLIBRARY,
	"locks":                                 kwLOCKS,
	"log":                                   kwLOG,
	"long":                                  kwLONG,
	"manual":                                kwMANUAL,
	"materialized":                          kwMATERIALIZED,
	"max_connections_per_hour":              kwMAX_CONNECTIONS_PER_HOUR,
	"max_queries_per_hour":                  kwMAX_QUERIES_PER_HOUR,
	"max_rows":                              kwMAX_ROWS,
	"max_size":                              kwMAX_SIZE,
	"max_updates_per_hour":                  kwMAX_UPDATES_PER_HOUR,
	"max_user_connections":                  kwMAX_USER_CONNECTIONS,
	"message_text":                          kwMESSAGE_TEXT,
	"microsecond":                           kwMICROSECOND,
	"mid":                                   kwMID,
	"middleint":                             kwMIDDLEINT,
	"min_rows":                              kwMIN_ROWS,
	"minute":                                kwMINUTE,
	"month":                                 kwMONTH,
	"mysql_errno":                           kwMYSQL_ERRNO,
	"ndb":                                   kwNDB,
	"ndbcluster":                            kwNDBCLUSTER,
	"network_namespace":                     kwNETWORK_NAMESPACE,
	"new":                                   kwNEW,
	"no_wait":                               kwNO_WAIT,
	"nodegroup":                             kwNODEGROUP,
	"now":                                   kwNOW,
	"number":                                kwNUMBER,
	"off":                                   kwOFF,
	"oj":                                    kwOJ,
	"others":                                kwOTHERS,
	"owner":                                 kwOWNER,
	"pack_keys":                             kwPACK_KEYS,
	"page":                                  kwPAGE,
	"parallel":                              kwPARALLEL,
	"parameters":                            kwPARAMETERS,
	"parse_tree":                            kwPARSE_TREE,
	"password_lock_time":                    kwPASSWORD_LOCK_TIME,
	"persist_only":                          kwPERSIST_ONLY,
	"plugin_dir":                            kwPLUGIN_DIR,
	"port":                                  kwPORT,
	"privilege_checks_user":                 kwPRIVILEGE_CHECKS_USER,
	"qualify":                               kwQUALIFY,
	"quarter":                               kwQUARTER,
	"read_only":                             kwREAD_ONLY,
	"read_write":                            kwREAD_WRITE,
	"redo_buffer_size":                      kwREDO_BUFFER_SIZE,
	"registration":                          kwREGISTRATION,
	"relational":                            kwRELATIONAL,
	"relay":                                 kwRELAY,
	"relay_log_file":                        kwRELAY_LOG_FILE,
	"relay_log_pos":                         kwRELAY_LOG_POS,
	"relay_thread":                          kwRELAY_THREAD,
	"replicate_do_db":                       kwREPLICATE_DO_DB,
	"replicate_do_table":                    kwREPLICATE_DO_TABLE,
	"replicate_ignore_db":                   kwREPLICATE_IGNORE_DB,
	"replicate_ignore_table":                kwREPLICATE_IGNORE_TABLE,
	"replicate_rewrite_db":                  kwREPLICATE_REWRITE_DB,
	"replicate_wild_do_table":               kwREPLICATE_WILD_DO_TABLE,
	"replicate_wild_ignore_table":           kwREPLICATE_WILD_IGNORE_TABLE,
	"require_row_format":                    kwREQUIRE_ROW_FORMAT,
	"require_table_primary_key_check":       kwREQUIRE_TABLE_PRIMARY_KEY_CHECK,
	"respect":                               kwRESPECT,
	"restore":                               kwRESTORE,
	"returned_sqlstate":                     kwRETURNED_SQLSTATE,
	"returning":                             kwRETURNING,
	"reverse":                               kwREVERSE,
	"row_count":                             kwROW_COUNT,
	"rtree":                                 kwRTREE,
	"s3":                                    kwS3,
	"schema_name":                           kwSCHEMA_NAME,
	"second":                                kwSECOND,
	"secondary":                             kwSECONDARY,
	"secondary_engine":                      kwSECONDARY_ENGINE,
	"secondary_engine_attribute":            kwSECONDARY_ENGINE_ATTRIBUTE,
	"secondary_load":                        kwSECONDARY_LOAD,
	"secondary_unload":                      kwSECONDARY_UNLOAD,
	"session_user":                          kwSESSION_USER,
	"sets":                                  kwSETS,
	"simple":                                kwSIMPLE,
	"slow":                                  kwSLOW,
	"socket":                                kwSOCKET,
	"some":                                  kwSOME,
	"source_auto_position":                  kwSOURCE_AUTO_POSITION,
	"source_bind":                           kwSOURCE_BIND,
	"source_compression_algorithms":         kwSOURCE_COMPRESSION_ALGORITHMS,
	"source_connect_retry":                  kwSOURCE_CONNECT_RETRY,
	"source_connection_auto_failover":       kwSOURCE_CONNECTION_AUTO_FAILOVER,
	"source_delay":                          kwSOURCE_DELAY,
	"source_heartbeat_period":               kwSOURCE_HEARTBEAT_PERIOD,
	"source_host":                           kwSOURCE_HOST,
	"source_log_file":                       kwSOURCE_LOG_FILE,
	"source_log_pos":                        kwSOURCE_LOG_POS,
	"source_password":                       kwSOURCE_PASSWORD,
	"source_port":                           kwSOURCE_PORT,
	"source_public_key_path":                kwSOURCE_PUBLIC_KEY_PATH,
	"source_retry_count":                    kwSOURCE_RETRY_COUNT,
	"source_ssl":                            kwSOURCE_SSL,
	"source_ssl_ca":                         kwSOURCE_SSL_CA,
	"source_ssl_capath":                     kwSOURCE_SSL_CAPATH,
	"source_ssl_cert":                       kwSOURCE_SSL_CERT,
	"source_ssl_cipher":                     kwSOURCE_SSL_CIPHER,
	"source_ssl_crl":                        kwSOURCE_SSL_CRL,
	"source_ssl_crlpath":                    kwSOURCE_SSL_CRLPATH,
	"source_ssl_key":                        kwSOURCE_SSL_KEY,
	"source_ssl_verify_server_cert":         kwSOURCE_SSL_VERIFY_SERVER_CERT,
	"source_tls_ciphersuites":               kwSOURCE_TLS_CIPHERSUITES,
	"source_tls_version":                    kwSOURCE_TLS_VERSION,
	"source_user":                           kwSOURCE_USER,
	"source_zstd_compression_level":         kwSOURCE_ZSTD_COMPRESSION_LEVEL,
	"sql_after_gtids":                       kwSQL_AFTER_GTIDS,
	"sql_after_mts_gaps":                    kwSQL_AFTER_MTS_GAPS,
	"sql_before_gtids":                      kwSQL_BEFORE_GTIDS,
	"sql_thread":                            kwSQL_THREAD,
	"sql_tsi_day":                           kwSQL_TSI_DAY,
	"sql_tsi_hour":                          kwSQL_TSI_HOUR,
	"sql_tsi_minute":                        kwSQL_TSI_MINUTE,
	"sql_tsi_month":                         kwSQL_TSI_MONTH,
	"sql_tsi_quarter":                       kwSQL_TSI_QUARTER,
	"sql_tsi_second":                        kwSQL_TSI_SECOND,
	"sql_tsi_week":                          kwSQL_TSI_WEEK,
	"sql_tsi_year":                          kwSQL_TSI_YEAR,
	"st_collect":                            kwST_COLLECT,
	"stats_auto_recalc":                     kwSTATS_AUTO_RECALC,
	"stats_persistent":                      kwSTATS_PERSISTENT,
	"stats_sample_pages":                    kwSTATS_SAMPLE_PAGES,
	"std":                                   kwSTD,
	"stddev":                                kwSTDDEV,
	"stddev_pop":                            kwSTDDEV_POP,
	"stddev_samp":                           kwSTDDEV_SAMP,
	"strict_load":                           kwSTRICT_LOAD,
	"string":                                kwSTRING,
	"subclass_origin":                       kwSUBCLASS_ORIGIN,
	"subdate":                               kwSUBDATE,
	"substr":                                kwSUBSTR,
	"super":                                 kwSUPER,
	"swaps":                                 kwSWAPS,
	"switches":                              kwSWITCHES,
	"sysdate":                               kwSYSDATE,
	"system_user":                           kwSYSTEM_USER,
	"table_checksum":                        kwTABLE_CHECKSUM,
	"table_name":                            kwTABLE_NAME,
	"tablesample":                           kwTABLESAMPLE,
	"thread_priority":                       kwTHREAD_PRIORITY,
	"ties":                                  kwTIES,
	"timestampadd":                          kwTIMESTAMPADD,
	"timestampdiff":                         kwTIMESTAMPDIFF,
	"types":                                 kwTYPES,
	"undo_buffer_size":                      kwUNDO_BUFFER_SIZE,
	"undofile":                              kwUNDOFILE,
	"unicode":                               kwUNICODE,
	"unregister":                            kwUNREGISTER,
	"uri":                                   kwURI,
	"url":                                   kwURL,
	"use_frm":                               kwUSE_FRM,
	"user_resources":                        kwUSER_RESOURCES,
	"validate":                              kwVALIDATE,
	"var_pop":                               kwVAR_POP,
	"var_samp":                              kwVAR_SAMP,
	"varcharacter":                          kwVARCHARACTER,
	"variance":                              kwVARIANCE,
	"vcpu":                                  kwVCPU,
	"vector":                                kwVECTOR,
	"verify_key_constraints":                kwVERIFY_KEY_CONSTRAINTS,
	"week":                                  kwWEEK,
	"weight_string":                         kwWEIGHT_STRING,
	"zone":                                  kwZONE,
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
	input      string
	pos        int
	start      int
	prevToken  int // type of the previously emitted token
	baseOffset int // added to all token Loc values for absolute positioning
}

// NewLexer creates a new MySQL lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NewLexerWithOffset creates a Lexer whose tokens have Loc = localOffset + baseOffset.
// This allows parsing a segment of a larger SQL string while producing absolute byte positions.
func NewLexerWithOffset(input string, baseOffset int) *Lexer {
	return &Lexer{input: input, baseOffset: baseOffset}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	tok := l.nextToken()
	tok.Loc += l.baseOffset // apply base offset for absolute positioning
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

	// System variable: @@ → emit tokAt2 and let the parser handle the rest
	if l.pos < len(l.input) && l.input[l.pos] == '@' {
		l.pos++ // skip second @
		return Token{Type: tokAt2, Str: "@@", Loc: start}
	}

	// User variable: @name — scan the variable name as one token
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

	return Token{Type: tokIDENT, Str: "@" + name, Loc: start}
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
