package parser

import "strings"

// Keyword token constants, one per keyword of StarRocksLex.g4, named for the
// word it matches, and one per word of the Doris residue below. Values start
// at 700 and grow with iota. Numeric values are NOT stable across edits — do
// not persist them.
const (
	kwACCESS TokenKind = 700 + iota
	kwACTIVE
	kwADD
	kwADMIN
	kwADVISOR
	kwAFTER
	kwAGGREGATE
	kwALL
	kwALTER
	kwANALYZE
	kwAND
	kwANTI
	kwAPPLY
	kwARRAY
	kwARRAY_AGG
	kwARRAY_AGG_DISTINCT
	kwAS
	kwASC
	kwASYNC
	kwAUTHENTICATION
	kwAUTHORS
	kwAUTOMATED
	kwAUTO_INCREMENT
	kwAVG
	kwBACKEND
	kwBACKENDS
	kwBACKUP
	kwBASE
	kwBEFORE
	kwBEGIN
	kwBETWEEN
	kwBIGINT
	kwBINARY
	kwBITMAP
	kwBITMAP_UNION
	kwBITSHIFTLEFT
	kwBITSHIFTRIGHT
	kwBITSHIFTRIGHTLOGICAL
	kwBLACKHOLE
	kwBLACKLIST
	kwBODY
	kwBOOLEAN
	kwBOTH
	kwBRANCH
	kwBROKER
	kwBUCKETS
	kwBUILTIN
	kwBY
	kwCACHE
	kwCANCEL
	kwCASE
	kwCAST
	kwCATALOG
	kwCATALOGS
	kwCEIL
	kwCHAIN
	kwCHAR
	kwCHARACTER
	kwCHARSET
	kwCHECK
	kwCLEAN
	kwCLEAR
	kwCLUSTER
	kwCLUSTERS
	kwCOLLATE
	kwCOLLATION
	kwCOLUMN
	kwCOLUMNS
	kwCOMMENT
	kwCOMMIT
	kwCOMMITTED
	kwCOMPACT
	kwCOMPACTION
	kwCOMPUTE
	kwCONFIG
	kwCONNECTION
	kwCONSISTENT
	kwCONVERT
	kwCOSTS
	kwCOUNT
	kwCREATE
	kwCROSS
	kwCUBE
	kwCUME_DIST
	kwCUMULATIVE
	kwCURRENT
	kwCURRENT_DATE
	kwCURRENT_ROLE
	kwCURRENT_TIME
	kwCURRENT_TIMESTAMP
	kwCURRENT_USER
	kwDATA
	kwDATABASE
	kwDATABASES
	kwDATACACHE
	kwDATE
	kwDATETIME
	kwDAY
	kwDAYS
	kwDEALLOCATE
	kwDECIMAL
	kwDECIMAL128
	kwDECIMAL32
	kwDECIMAL64
	kwDECIMALV2
	kwDECOMMISSION
	kwDEFAULT
	kwDEFERRED
	kwDELETE
	kwDENSE_RANK
	kwDESC
	kwDESCRIBE
	kwDICTIONARY
	kwDICTIONARY_GET
	kwDISABLE
	kwDISK
	kwDISTINCT
	kwDISTRIBUTED
	kwDISTRIBUTION
	kwDIV
	kwDOUBLE
	kwDROP
	kwDUAL
	kwDUPLICATE
	kwDYNAMIC
	kwELSE
	kwENABLE
	kwENCLOSE
	kwEND
	kwENGINE
	kwENGINES
	kwERRORS
	kwESCAPE
	kwEVENTS
	kwEVERY
	kwEXCEPT
	kwEXECUTE
	kwEXISTS
	kwEXPLAIN
	kwEXPORT
	kwEXTERNAL
	kwEXTRACT
	kwFAILPOINT
	kwFAILPOINTS
	kwFALSE
	kwFIELD
	kwFIELDS
	kwFILE
	kwFILES
	kwFILTER
	kwFIRST
	kwFIRST_VALUE
	kwFLOAT
	kwFLOOR
	kwFN
	kwFOLLOWER
	kwFOLLOWING
	kwFOR
	kwFORCE
	kwFORMAT
	kwFREE
	kwFROM
	kwFRONTEND
	kwFRONTENDS
	kwFULL
	kwFUNCTION
	kwFUNCTIONS
	kwGIN
	kwGLOBAL
	kwGRANT
	kwGRANTS
	kwGROUP
	kwGROUPING
	kwGROUPING_ID
	kwGROUPS
	kwGROUP_CONCAT
	kwHASH
	kwHAVING
	kwHELP
	kwHISTOGRAM
	kwHLL
	kwHLL_UNION
	kwHOST
	kwHOUR
	kwHOURS
	kwHUB
	kwIDENTIFIED
	kwIF
	kwIGNORE
	kwIMAGE
	kwIMMEDIATE
	kwIMPERSONATE
	kwIN
	kwINACTIVE
	kwINCREMENTAL
	kwINDEX
	kwINDEXES
	kwINFILE
	kwINNER
	kwINSERT
	kwINSTALL
	kwINT
	kwINTEGER
	kwINTEGRATION
	kwINTEGRATIONS
	kwINTERMEDIATE
	kwINTERSECT
	kwINTERVAL
	kwINTO
	kwINVOKER
	kwIS
	kwISOLATION
	kwJOB
	kwJOIN
	kwJSON
	kwKEY
	kwKEYS
	kwKILL
	kwLABEL
	kwLAG
	kwLARGEINT
	kwLAST
	kwLAST_VALUE
	kwLATERAL
	kwLEAD
	kwLEFT
	kwLESS
	kwLEVEL
	kwLIKE
	kwLIMIT
	kwLIST
	kwLOAD
	kwLOCAL
	kwLOCALTIME
	kwLOCALTIMESTAMP
	kwLOCATION
	kwLOCATIONS
	kwLOCK
	kwLOGICAL
	kwLOGS
	kwLOW_PRIORITY
	kwMANUAL
	kwMAP
	kwMAPPING
	kwMAPPINGS
	kwMASKING
	kwMATCH
	kwMATERIALIZED
	kwMAX
	kwMAXVALUE
	kwMERGE
	kwMETA
	kwMICROSECOND
	kwMILLISECOND
	kwMIN
	kwMINUS
	kwMINUTE
	kwMINUTES
	kwMOD
	kwMODE
	kwMODIFY
	kwMONTH
	kwNAME
	kwNAMES
	kwNEGATIVE
	kwNGRAMBF
	kwNO
	kwNODE
	kwNODES
	kwNONE
	kwNOT
	kwNTILE
	kwNULL
	kwNULLS
	kwNUMBER
	kwNUMERIC
	kwOBSERVER
	kwOF
	kwOFF
	kwOFFSET
	kwON
	kwONLY
	kwOPEN
	kwOPERATE
	kwOPTIMIZE
	kwOPTIMIZER
	kwOPTION
	kwOR
	kwORDER
	kwOUTER
	kwOUTFILE
	kwOVER
	kwOVERWRITE
	kwPARTITION
	kwPARTITIONS
	kwPASSWORD
	kwPATH
	kwPAUSE
	kwPENDING
	kwPERCENTILE
	kwPERCENTILE_UNION
	kwPERCENT_RANK
	kwPERSISTENT
	kwPIPE
	kwPIPES
	kwPIVOT
	kwPLAN
	kwPLUGIN
	kwPLUGINS
	kwPOLICIES
	kwPOLICY
	kwPRECEDING
	kwPREPARE
	kwPRIMARY
	kwPRIORITY
	kwPRIVILEGES
	kwPROBABILITY
	kwPROC
	kwPROCEDURE
	kwPROCESSLIST
	kwPROFILE
	kwPROFILELIST
	kwPROPERTIES
	kwPROPERTY
	kwQUALIFY
	kwQUARTER
	kwQUERIES
	kwQUERY
	kwQUEUE
	kwQUOTA
	kwRANDOM
	kwRANGE
	kwRANK
	kwREAD
	kwREASON
	kwRECOVER
	kwREFRESH
	kwREGEXP
	kwRELEASE
	kwREMOVE
	kwRENAME
	kwREPAIR
	kwREPEATABLE
	kwREPLACE
	kwREPLACE_IF_NOT_NULL
	kwREPLICA
	kwREPOSITORIES
	kwREPOSITORY
	kwRESOURCE
	kwRESOURCES
	kwRESTORE
	kwRESUME
	kwRETAIN
	kwRETENTION
	kwRETRY
	kwRETURNS
	kwREVERT
	kwREVOKE
	kwREWRITE
	kwRIGHT
	kwRLIKE
	kwROLE
	kwROLES
	kwROLLBACK
	kwROLLUP
	kwROUTINE
	kwROW
	kwROWS
	kwROW_NUMBER
	kwRULE
	kwRULES
	kwRUNNING
	kwSAMPLE
	kwSCHEDULE
	kwSCHEDULER
	kwSCHEMA
	kwSCHEMAS
	kwSECOND
	kwSECURITY
	kwSELECT
	kwSEMI
	kwSEPARATOR
	kwSERIALIZABLE
	kwSESSION
	kwSET
	kwSETS
	kwSET_VAR
	kwSHOW
	kwSIGNED
	kwSKIP_HEADER
	kwSMALLINT
	kwSNAPSHOT
	kwSNAPSHOTS
	kwSQLBLACKLIST
	kwSTART
	kwSTATS
	kwSTATUS
	kwSTOP
	kwSTORAGE
	kwSTREAM
	kwSTRING
	kwSTRUCT
	kwSUBMIT
	kwSUM
	kwSUSPEND
	kwSWAP
	kwSYNC
	kwSYSTEM
	kwSYSTEM_TIME
	kwTABLE
	kwTABLES
	kwTABLET
	kwTABLETS
	kwTAG
	kwTASK
	kwTEMPORARY
	kwTERMINATED
	kwTEXT
	kwTHAN
	kwTHEN
	kwTIME
	kwTIMES
	kwTIMESTAMP
	kwTIMESTAMPADD
	kwTIMESTAMPDIFF
	kwTINYINT
	kwTO
	kwTRACE
	kwTRANSACTION
	kwTRANSLATE
	kwTRIGGERS
	kwTRIM_SPACE
	kwTRUE
	kwTRUNCATE
	kwTYPE
	kwTYPES
	kwUNBOUNDED
	kwUNCOMMITTED
	kwUNINSTALL
	kwUNION
	kwUNIQUE
	kwUNLOCK
	kwUNSET
	kwUNSIGNED
	kwUPDATE
	kwUSAGE
	kwUSE
	kwUSER
	kwUSERS
	kwUSING
	kwVALUE
	kwVALUES
	kwVARBINARY
	kwVARCHAR
	kwVARIABLES
	kwVECTOR
	kwVERBOSE
	kwVERSION
	kwVIEW
	kwVIEWS
	kwVOLUME
	kwVOLUMES
	kwWAREHOUSE
	kwWAREHOUSES
	kwWARNINGS
	kwWEEK
	kwWHEN
	kwWHERE
	kwWHITELIST
	kwWITH
	kwWORK
	kwWRITE
	kwYEAR

	// The Doris residue: Doris words this parser still reads, which
	// StarRocksLex.g4 has no token for. See dorisOnlyKeywords.
	kwACCOUNT_LOCK
	kwACCOUNT_UNLOCK
	kwACTIONS
	kwAGG_STATE
	kwALWAYS
	kwAUTO
	kwBRIEF
	kwCALL
	kwCOMPLETE
	kwCONDITIONS
	kwCONSTRAINT
	kwCONSTRAINTS
	kwDATETIMEV1
	kwDATETIMEV2
	kwDATEV1
	kwDATEV2
	kwDECIMALV3
	kwDO
	kwENCRYPTKEY
	kwEXTENDED
	kwFAILED_LOGIN_ATTEMPTS
	kwFEATURE
	kwFOREIGN
	kwGENERATED
	kwGENERIC
	kwHDFS
	kwIPV4
	kwIPV6
	kwJOBS
	kwJSONB
	kwLAYOUT
	kwMATCHED
	kwNEVER
	kwPASSWORD_EXPIRE
	kwPASSWORD_HISTORY
	kwPASSWORD_LOCK_TIME
	kwPERMISSIVE
	kwQUANTILE_STATE
	kwQUANTILE_UNION
	kwQUEUED
	kwREFERENCES
	kwRESTRICTIVE
	kwS3
	kwSONAME
	kwSQL_BLOCK_RULE
	kwSTAGE
	kwSTARTS
	kwVARIABLE
	kwVARIANT
	kwVAULT
	kwWORKLOAD
)

// keywordMap maps a keyword text, lowercased, to its token. It is the pinned
// release's own grammar, not the Doris fork's: StarRocksLex.g4's keyword
// rules, with the reservedness StarRocks.g4's nonReserved rule gives them,
// both read from ../testdata. TestKeywordTableMatchesGrammar fails a table
// that stops matching them, so a keyword the engine adds, drops or un-reserves
// is a test failure rather than a silent over-acceptance.
var keywordMap = map[string]TokenKind{
	"access":               kwACCESS,
	"active":               kwACTIVE,
	"add":                  kwADD,
	"admin":                kwADMIN,
	"advisor":              kwADVISOR,
	"after":                kwAFTER,
	"aggregate":            kwAGGREGATE,
	"all":                  kwALL,
	"alter":                kwALTER,
	"analyze":              kwANALYZE,
	"and":                  kwAND,
	"anti":                 kwANTI,
	"apply":                kwAPPLY,
	"array":                kwARRAY,
	"array_agg":            kwARRAY_AGG,
	"array_agg_distinct":   kwARRAY_AGG_DISTINCT,
	"as":                   kwAS,
	"asc":                  kwASC,
	"async":                kwASYNC,
	"authentication":       kwAUTHENTICATION,
	"authors":              kwAUTHORS,
	"automated":            kwAUTOMATED,
	"auto_increment":       kwAUTO_INCREMENT,
	"avg":                  kwAVG,
	"backend":              kwBACKEND,
	"backends":             kwBACKENDS,
	"backup":               kwBACKUP,
	"base":                 kwBASE,
	"before":               kwBEFORE,
	"begin":                kwBEGIN,
	"between":              kwBETWEEN,
	"bigint":               kwBIGINT,
	"binary":               kwBINARY,
	"bitmap":               kwBITMAP,
	"bitmap_union":         kwBITMAP_UNION,
	"bitshiftleft":         kwBITSHIFTLEFT,
	"bitshiftright":        kwBITSHIFTRIGHT,
	"bitshiftrightlogical": kwBITSHIFTRIGHTLOGICAL,
	"blackhole":            kwBLACKHOLE,
	"blacklist":            kwBLACKLIST,
	"body":                 kwBODY,
	"boolean":              kwBOOLEAN,
	"both":                 kwBOTH,
	"branch":               kwBRANCH,
	"broker":               kwBROKER,
	"buckets":              kwBUCKETS,
	"builtin":              kwBUILTIN,
	"by":                   kwBY,
	"cache":                kwCACHE,
	"cancel":               kwCANCEL,
	"case":                 kwCASE,
	"cast":                 kwCAST,
	"catalog":              kwCATALOG,
	"catalogs":             kwCATALOGS,
	"ceil":                 kwCEIL,
	"chain":                kwCHAIN,
	"char":                 kwCHAR,
	"character":            kwCHARACTER,
	"charset":              kwCHARSET,
	"check":                kwCHECK,
	"clean":                kwCLEAN,
	"clear":                kwCLEAR,
	"cluster":              kwCLUSTER,
	"clusters":             kwCLUSTERS,
	"collate":              kwCOLLATE,
	"collation":            kwCOLLATION,
	"column":               kwCOLUMN,
	"columns":              kwCOLUMNS,
	"comment":              kwCOMMENT,
	"commit":               kwCOMMIT,
	"committed":            kwCOMMITTED,
	"compact":              kwCOMPACT,
	"compaction":           kwCOMPACTION,
	"compute":              kwCOMPUTE,
	"config":               kwCONFIG,
	"connection":           kwCONNECTION,
	"consistent":           kwCONSISTENT,
	"convert":              kwCONVERT,
	"costs":                kwCOSTS,
	"count":                kwCOUNT,
	"create":               kwCREATE,
	"cross":                kwCROSS,
	"cube":                 kwCUBE,
	"cume_dist":            kwCUME_DIST,
	"cumulative":           kwCUMULATIVE,
	"current":              kwCURRENT,
	"current_date":         kwCURRENT_DATE,
	"current_role":         kwCURRENT_ROLE,
	"current_time":         kwCURRENT_TIME,
	"current_timestamp":    kwCURRENT_TIMESTAMP,
	"current_user":         kwCURRENT_USER,
	"data":                 kwDATA,
	"database":             kwDATABASE,
	"databases":            kwDATABASES,
	"datacache":            kwDATACACHE,
	"date":                 kwDATE,
	"datetime":             kwDATETIME,
	"day":                  kwDAY,
	"days":                 kwDAYS,
	"deallocate":           kwDEALLOCATE,
	"decimal":              kwDECIMAL,
	"decimal128":           kwDECIMAL128,
	"decimal32":            kwDECIMAL32,
	"decimal64":            kwDECIMAL64,
	"decimalv2":            kwDECIMALV2,
	"decommission":         kwDECOMMISSION,
	"default":              kwDEFAULT,
	"deferred":             kwDEFERRED,
	"delete":               kwDELETE,
	"dense_rank":           kwDENSE_RANK,
	"desc":                 kwDESC,
	"describe":             kwDESCRIBE,
	"dictionary":           kwDICTIONARY,
	"dictionary_get":       kwDICTIONARY_GET,
	"disable":              kwDISABLE,
	"disk":                 kwDISK,
	"distinct":             kwDISTINCT,
	"distributed":          kwDISTRIBUTED,
	"distribution":         kwDISTRIBUTION,
	"div":                  kwDIV,
	"double":               kwDOUBLE,
	"drop":                 kwDROP,
	"dual":                 kwDUAL,
	"duplicate":            kwDUPLICATE,
	"dynamic":              kwDYNAMIC,
	"else":                 kwELSE,
	"enable":               kwENABLE,
	"enclose":              kwENCLOSE,
	"end":                  kwEND,
	"engine":               kwENGINE,
	"engines":              kwENGINES,
	"errors":               kwERRORS,
	"escape":               kwESCAPE,
	"events":               kwEVENTS,
	"every":                kwEVERY,
	"except":               kwEXCEPT,
	"execute":              kwEXECUTE,
	"exists":               kwEXISTS,
	"explain":              kwEXPLAIN,
	"export":               kwEXPORT,
	"external":             kwEXTERNAL,
	"extract":              kwEXTRACT,
	"failpoint":            kwFAILPOINT,
	"failpoints":           kwFAILPOINTS,
	"false":                kwFALSE,
	"field":                kwFIELD,
	"fields":               kwFIELDS,
	"file":                 kwFILE,
	"files":                kwFILES,
	"filter":               kwFILTER,
	"first":                kwFIRST,
	"first_value":          kwFIRST_VALUE,
	"float":                kwFLOAT,
	"floor":                kwFLOOR,
	"fn":                   kwFN,
	"follower":             kwFOLLOWER,
	"following":            kwFOLLOWING,
	"for":                  kwFOR,
	"force":                kwFORCE,
	"format":               kwFORMAT,
	"free":                 kwFREE,
	"from":                 kwFROM,
	"frontend":             kwFRONTEND,
	"frontends":            kwFRONTENDS,
	"full":                 kwFULL,
	"function":             kwFUNCTION,
	"functions":            kwFUNCTIONS,
	"gin":                  kwGIN,
	"global":               kwGLOBAL,
	"grant":                kwGRANT,
	"grants":               kwGRANTS,
	"group":                kwGROUP,
	"grouping":             kwGROUPING,
	"grouping_id":          kwGROUPING_ID,
	"groups":               kwGROUPS,
	"group_concat":         kwGROUP_CONCAT,
	"hash":                 kwHASH,
	"having":               kwHAVING,
	"help":                 kwHELP,
	"histogram":            kwHISTOGRAM,
	"hll":                  kwHLL,
	"hll_union":            kwHLL_UNION,
	"host":                 kwHOST,
	"hour":                 kwHOUR,
	"hours":                kwHOURS,
	"hub":                  kwHUB,
	"identified":           kwIDENTIFIED,
	"if":                   kwIF,
	"ignore":               kwIGNORE,
	"image":                kwIMAGE,
	"immediate":            kwIMMEDIATE,
	"impersonate":          kwIMPERSONATE,
	"in":                   kwIN,
	"inactive":             kwINACTIVE,
	"incremental":          kwINCREMENTAL,
	"index":                kwINDEX,
	"indexes":              kwINDEXES,
	"infile":               kwINFILE,
	"inner":                kwINNER,
	"insert":               kwINSERT,
	"install":              kwINSTALL,
	"int":                  kwINT,
	"integer":              kwINTEGER,
	"integration":          kwINTEGRATION,
	"integrations":         kwINTEGRATIONS,
	"intermediate":         kwINTERMEDIATE,
	"intersect":            kwINTERSECT,
	"interval":             kwINTERVAL,
	"into":                 kwINTO,
	"invoker":              kwINVOKER,
	"is":                   kwIS,
	"isolation":            kwISOLATION,
	"job":                  kwJOB,
	"join":                 kwJOIN,
	"json":                 kwJSON,
	"key":                  kwKEY,
	"keys":                 kwKEYS,
	"kill":                 kwKILL,
	"label":                kwLABEL,
	"lag":                  kwLAG,
	"largeint":             kwLARGEINT,
	"last":                 kwLAST,
	"last_value":           kwLAST_VALUE,
	"lateral":              kwLATERAL,
	"lead":                 kwLEAD,
	"left":                 kwLEFT,
	"less":                 kwLESS,
	"level":                kwLEVEL,
	"like":                 kwLIKE,
	"limit":                kwLIMIT,
	"list":                 kwLIST,
	"load":                 kwLOAD,
	"local":                kwLOCAL,
	"localtime":            kwLOCALTIME,
	"localtimestamp":       kwLOCALTIMESTAMP,
	"location":             kwLOCATION,
	"locations":            kwLOCATIONS,
	"lock":                 kwLOCK,
	"logical":              kwLOGICAL,
	"logs":                 kwLOGS,
	"low_priority":         kwLOW_PRIORITY,
	"manual":               kwMANUAL,
	"map":                  kwMAP,
	"mapping":              kwMAPPING,
	"mappings":             kwMAPPINGS,
	"masking":              kwMASKING,
	"match":                kwMATCH,
	"materialized":         kwMATERIALIZED,
	"max":                  kwMAX,
	"maxvalue":             kwMAXVALUE,
	"merge":                kwMERGE,
	"meta":                 kwMETA,
	"microsecond":          kwMICROSECOND,
	"millisecond":          kwMILLISECOND,
	"min":                  kwMIN,
	"minus":                kwMINUS,
	"minute":               kwMINUTE,
	"minutes":              kwMINUTES,
	"mod":                  kwMOD,
	"mode":                 kwMODE,
	"modify":               kwMODIFY,
	"month":                kwMONTH,
	"name":                 kwNAME,
	"names":                kwNAMES,
	"negative":             kwNEGATIVE,
	"ngrambf":              kwNGRAMBF,
	"no":                   kwNO,
	"node":                 kwNODE,
	"nodes":                kwNODES,
	"none":                 kwNONE,
	"not":                  kwNOT,
	"ntile":                kwNTILE,
	"null":                 kwNULL,
	"nulls":                kwNULLS,
	"number":               kwNUMBER,
	"numeric":              kwNUMERIC,
	"observer":             kwOBSERVER,
	"of":                   kwOF,
	"off":                  kwOFF,
	"offset":               kwOFFSET,
	"on":                   kwON,
	"only":                 kwONLY,
	"open":                 kwOPEN,
	"operate":              kwOPERATE,
	"optimize":             kwOPTIMIZE,
	"optimizer":            kwOPTIMIZER,
	"option":               kwOPTION,
	"or":                   kwOR,
	"order":                kwORDER,
	"outer":                kwOUTER,
	"outfile":              kwOUTFILE,
	"over":                 kwOVER,
	"overwrite":            kwOVERWRITE,
	"partition":            kwPARTITION,
	"partitions":           kwPARTITIONS,
	"password":             kwPASSWORD,
	"path":                 kwPATH,
	"pause":                kwPAUSE,
	"pending":              kwPENDING,
	"percentile":           kwPERCENTILE,
	"percentile_union":     kwPERCENTILE_UNION,
	"percent_rank":         kwPERCENT_RANK,
	"persistent":           kwPERSISTENT,
	"pipe":                 kwPIPE,
	"pipes":                kwPIPES,
	"pivot":                kwPIVOT,
	"plan":                 kwPLAN,
	"plugin":               kwPLUGIN,
	"plugins":              kwPLUGINS,
	"policies":             kwPOLICIES,
	"policy":               kwPOLICY,
	"preceding":            kwPRECEDING,
	"prepare":              kwPREPARE,
	"primary":              kwPRIMARY,
	"priority":             kwPRIORITY,
	"privileges":           kwPRIVILEGES,
	"probability":          kwPROBABILITY,
	"proc":                 kwPROC,
	"procedure":            kwPROCEDURE,
	"processlist":          kwPROCESSLIST,
	"profile":              kwPROFILE,
	"profilelist":          kwPROFILELIST,
	"properties":           kwPROPERTIES,
	"property":             kwPROPERTY,
	"qualify":              kwQUALIFY,
	"quarter":              kwQUARTER,
	"queries":              kwQUERIES,
	"query":                kwQUERY,
	"queue":                kwQUEUE,
	"quota":                kwQUOTA,
	"random":               kwRANDOM,
	"range":                kwRANGE,
	"rank":                 kwRANK,
	"read":                 kwREAD,
	"reason":               kwREASON,
	"recover":              kwRECOVER,
	"refresh":              kwREFRESH,
	"regexp":               kwREGEXP,
	"release":              kwRELEASE,
	"remove":               kwREMOVE,
	"rename":               kwRENAME,
	"repair":               kwREPAIR,
	"repeatable":           kwREPEATABLE,
	"replace":              kwREPLACE,
	"replace_if_not_null":  kwREPLACE_IF_NOT_NULL,
	"replica":              kwREPLICA,
	"repositories":         kwREPOSITORIES,
	"repository":           kwREPOSITORY,
	"resource":             kwRESOURCE,
	"resources":            kwRESOURCES,
	"restore":              kwRESTORE,
	"resume":               kwRESUME,
	"retain":               kwRETAIN,
	"retention":            kwRETENTION,
	"retry":                kwRETRY,
	"returns":              kwRETURNS,
	"revert":               kwREVERT,
	"revoke":               kwREVOKE,
	"rewrite":              kwREWRITE,
	"right":                kwRIGHT,
	"rlike":                kwRLIKE,
	"role":                 kwROLE,
	"roles":                kwROLES,
	"rollback":             kwROLLBACK,
	"rollup":               kwROLLUP,
	"routine":              kwROUTINE,
	"row":                  kwROW,
	"rows":                 kwROWS,
	"row_number":           kwROW_NUMBER,
	"rule":                 kwRULE,
	"rules":                kwRULES,
	"running":              kwRUNNING,
	"sample":               kwSAMPLE,
	"schedule":             kwSCHEDULE,
	"scheduler":            kwSCHEDULER,
	"schema":               kwSCHEMA,
	"schemas":              kwSCHEMAS,
	"second":               kwSECOND,
	"security":             kwSECURITY,
	"select":               kwSELECT,
	"semi":                 kwSEMI,
	"separator":            kwSEPARATOR,
	"serializable":         kwSERIALIZABLE,
	"session":              kwSESSION,
	"set":                  kwSET,
	"sets":                 kwSETS,
	"set_var":              kwSET_VAR,
	"show":                 kwSHOW,
	"signed":               kwSIGNED,
	"skip_header":          kwSKIP_HEADER,
	"smallint":             kwSMALLINT,
	"snapshot":             kwSNAPSHOT,
	"snapshots":            kwSNAPSHOTS,
	"sqlblacklist":         kwSQLBLACKLIST,
	"start":                kwSTART,
	"stats":                kwSTATS,
	"status":               kwSTATUS,
	"stop":                 kwSTOP,
	"storage":              kwSTORAGE,
	"stream":               kwSTREAM,
	"string":               kwSTRING,
	"struct":               kwSTRUCT,
	"submit":               kwSUBMIT,
	"sum":                  kwSUM,
	"suspend":              kwSUSPEND,
	"swap":                 kwSWAP,
	"sync":                 kwSYNC,
	"system":               kwSYSTEM,
	"system_time":          kwSYSTEM_TIME,
	"table":                kwTABLE,
	"tables":               kwTABLES,
	"tablet":               kwTABLET,
	"tablets":              kwTABLETS,
	"tag":                  kwTAG,
	"task":                 kwTASK,
	"temporary":            kwTEMPORARY,
	"terminated":           kwTERMINATED,
	"text":                 kwTEXT,
	"than":                 kwTHAN,
	"then":                 kwTHEN,
	"time":                 kwTIME,
	"times":                kwTIMES,
	"timestamp":            kwTIMESTAMP,
	"timestampadd":         kwTIMESTAMPADD,
	"timestampdiff":        kwTIMESTAMPDIFF,
	"tinyint":              kwTINYINT,
	"to":                   kwTO,
	"trace":                kwTRACE,
	"transaction":          kwTRANSACTION,
	"translate":            kwTRANSLATE,
	"triggers":             kwTRIGGERS,
	"trim_space":           kwTRIM_SPACE,
	"true":                 kwTRUE,
	"truncate":             kwTRUNCATE,
	"type":                 kwTYPE,
	"types":                kwTYPES,
	"unbounded":            kwUNBOUNDED,
	"uncommitted":          kwUNCOMMITTED,
	"uninstall":            kwUNINSTALL,
	"union":                kwUNION,
	"unique":               kwUNIQUE,
	"unlock":               kwUNLOCK,
	"unset":                kwUNSET,
	"unsigned":             kwUNSIGNED,
	"update":               kwUPDATE,
	"usage":                kwUSAGE,
	"use":                  kwUSE,
	"user":                 kwUSER,
	"users":                kwUSERS,
	"using":                kwUSING,
	"value":                kwVALUE,
	"values":               kwVALUES,
	"varbinary":            kwVARBINARY,
	"varchar":              kwVARCHAR,
	"variables":            kwVARIABLES,
	"vector":               kwVECTOR,
	"verbose":              kwVERBOSE,
	"version":              kwVERSION,
	"view":                 kwVIEW,
	"views":                kwVIEWS,
	"volume":               kwVOLUME,
	"volumes":              kwVOLUMES,
	"warehouse":            kwWAREHOUSE,
	"warehouses":           kwWAREHOUSES,
	"warnings":             kwWARNINGS,
	"week":                 kwWEEK,
	"when":                 kwWHEN,
	"where":                kwWHERE,
	"whitelist":            kwWHITELIST,
	"with":                 kwWITH,
	"work":                 kwWORK,
	"write":                kwWRITE,
	"year":                 kwYEAR,
}

// dorisOnlyKeywords maps the Doris words this hand-written parser still reads
// to their tokens. StarRocksLex.g4 has no rule for any of them, so the engine
// reads each as an identifier and has no statement or clause that uses it;
// what they reach here is the fork's Doris-only grammar, which follows the
// words the parser matched: CALL, MERGE, BUILD INDEX, the WORKLOAD GROUP,
// WORKLOAD POLICY and SQL_BLOCK_RULE statements, STAGE, STORAGE VAULT and
// ENCRYPT KEY objects, Doris's column types, SHOW modifiers, MERGE and MTMV
// clauses and INSERT-time constraints. They stay non-reserved, so an
// identifier use of each reads as the engine reads it; the statements and
// clauses behind them are pinned in testdata/compat.tsv as rows only omni
// accepts, so removing one makes the pin fail rather than go unnoticed.
var dorisOnlyKeywords = map[string]TokenKind{
	"account_lock":          kwACCOUNT_LOCK,
	"account_unlock":        kwACCOUNT_UNLOCK,
	"actions":               kwACTIONS,
	"agg_state":             kwAGG_STATE,
	"always":                kwALWAYS,
	"auto":                  kwAUTO,
	"brief":                 kwBRIEF,
	"call":                  kwCALL,
	"complete":              kwCOMPLETE,
	"conditions":            kwCONDITIONS,
	"constraint":            kwCONSTRAINT,
	"constraints":           kwCONSTRAINTS,
	"datetimev1":            kwDATETIMEV1,
	"datetimev2":            kwDATETIMEV2,
	"datev1":                kwDATEV1,
	"datev2":                kwDATEV2,
	"decimalv3":             kwDECIMALV3,
	"do":                    kwDO,
	"encryptkey":            kwENCRYPTKEY,
	"extended":              kwEXTENDED,
	"failed_login_attempts": kwFAILED_LOGIN_ATTEMPTS,
	"feature":               kwFEATURE,
	"foreign":               kwFOREIGN,
	"generated":             kwGENERATED,
	"generic":               kwGENERIC,
	"hdfs":                  kwHDFS,
	"ipv4":                  kwIPV4,
	"ipv6":                  kwIPV6,
	"jobs":                  kwJOBS,
	"jsonb":                 kwJSONB,
	"layout":                kwLAYOUT,
	"matched":               kwMATCHED,
	"never":                 kwNEVER,
	"password_expire":       kwPASSWORD_EXPIRE,
	"password_history":      kwPASSWORD_HISTORY,
	"password_lock_time":    kwPASSWORD_LOCK_TIME,
	"permissive":            kwPERMISSIVE,
	"quantile_state":        kwQUANTILE_STATE,
	"quantile_union":        kwQUANTILE_UNION,
	"queued":                kwQUEUED,
	"references":            kwREFERENCES,
	"restrictive":           kwRESTRICTIVE,
	"s3":                    kwS3,
	"soname":                kwSONAME,
	"sql_block_rule":        kwSQL_BLOCK_RULE,
	"stage":                 kwSTAGE,
	"starts":                kwSTARTS,
	"variable":              kwVARIABLE,
	"variant":               kwVARIANT,
	"vault":                 kwVAULT,
	"workload":              kwWORKLOAD,
}

// dorisOnlyKeywordKinds holds the kinds of dorisOnlyKeywords.
var dorisOnlyKeywordKinds = func() map[TokenKind]string {
	m := make(map[TokenKind]string, len(dorisOnlyKeywords))
	for word, kind := range dorisOnlyKeywords {
		m[kind] = word
	}
	return m
}()

// nonReservedKeywords holds the keywords StarRocks.g4's nonReserved rule
// allows as bare identifiers. A keyword the rule does not name is reserved.
var nonReservedKeywords = map[TokenKind]bool{
	kwACCESS:              true,
	kwACTIVE:              true,
	kwADMIN:               true,
	kwADVISOR:             true,
	kwAFTER:               true,
	kwAGGREGATE:           true,
	kwANTI:                true,
	kwAPPLY:               true,
	kwARRAY_AGG:           true,
	kwARRAY_AGG_DISTINCT:  true,
	kwASYNC:               true,
	kwAUTHENTICATION:      true,
	kwAUTHORS:             true,
	kwAUTOMATED:           true,
	kwAUTO_INCREMENT:      true,
	kwAVG:                 true,
	kwBACKEND:             true,
	kwBACKENDS:            true,
	kwBACKUP:              true,
	kwBASE:                true,
	kwBEFORE:              true,
	kwBEGIN:               true,
	kwBINARY:              true,
	kwBITMAP_UNION:        true,
	kwBLACKHOLE:           true,
	kwBLACKLIST:           true,
	kwBODY:                true,
	kwBOOLEAN:             true,
	kwBRANCH:              true,
	kwBROKER:              true,
	kwBUCKETS:             true,
	kwBUILTIN:             true,
	kwCACHE:               true,
	kwCANCEL:              true,
	kwCAST:                true,
	kwCATALOG:             true,
	kwCATALOGS:            true,
	kwCEIL:                true,
	kwCHAIN:               true,
	kwCHARSET:             true,
	kwCLEAN:               true,
	kwCLEAR:               true,
	kwCLUSTER:             true,
	kwCLUSTERS:            true,
	kwCOLLATION:           true,
	kwCOLUMNS:             true,
	kwCOMMENT:             true,
	kwCOMMIT:              true,
	kwCOMMITTED:           true,
	kwCOMPACT:             true,
	kwCOMPUTE:             true,
	kwCONFIG:              true,
	kwCONNECTION:          true,
	kwCONSISTENT:          true,
	kwCOSTS:               true,
	kwCOUNT:               true,
	kwCUME_DIST:           true,
	kwCUMULATIVE:          true,
	kwCURRENT:             true,
	kwDATA:                true,
	kwDATACACHE:           true,
	kwDATE:                true,
	kwDATETIME:            true,
	kwDAY:                 true,
	kwDAYS:                true,
	kwDEALLOCATE:          true,
	kwDECOMMISSION:        true,
	kwDICTIONARY:          true,
	kwDICTIONARY_GET:      true,
	kwDISABLE:             true,
	kwDISK:                true,
	kwDISTRIBUTED:         true,
	kwDISTRIBUTION:        true,
	kwDUPLICATE:           true,
	kwDYNAMIC:             true,
	kwENABLE:              true,
	kwENCLOSE:             true,
	kwEND:                 true,
	kwENGINE:              true,
	kwENGINES:             true,
	kwERRORS:              true,
	kwESCAPE:              true,
	kwEVENTS:              true,
	kwEVERY:               true,
	kwEXECUTE:             true,
	kwEXPORT:              true,
	kwEXTERNAL:            true,
	kwEXTRACT:             true,
	kwFAILPOINT:           true,
	kwFAILPOINTS:          true,
	kwFIELD:               true,
	kwFIELDS:              true,
	kwFILE:                true,
	kwFILTER:              true,
	kwFIRST:               true,
	kwFLOOR:               true,
	kwFN:                  true,
	kwFOLLOWER:            true,
	kwFOLLOWING:           true,
	kwFORMAT:              true,
	kwFREE:                true,
	kwFRONTEND:            true,
	kwFRONTENDS:           true,
	kwFUNCTIONS:           true,
	kwGLOBAL:              true,
	kwGRANTS:              true,
	kwGROUP_CONCAT:        true,
	kwHASH:                true,
	kwHELP:                true,
	kwHISTOGRAM:           true,
	kwHLL_UNION:           true,
	kwHOST:                true,
	kwHOUR:                true,
	kwHOURS:               true,
	kwHUB:                 true,
	kwIDENTIFIED:          true,
	kwIMAGE:               true,
	kwIMPERSONATE:         true,
	kwINACTIVE:            true,
	kwINCREMENTAL:         true,
	kwINDEXES:             true,
	kwINSTALL:             true,
	kwINTEGRATION:         true,
	kwINTEGRATIONS:        true,
	kwINTERMEDIATE:        true,
	kwINTERVAL:            true,
	kwISOLATION:           true,
	kwJOB:                 true,
	kwLABEL:               true,
	kwLAST:                true,
	kwLESS:                true,
	kwLEVEL:               true,
	kwLIST:                true,
	kwLOCAL:               true,
	kwLOCATION:            true,
	kwLOCATIONS:           true,
	kwLOCK:                true,
	kwLOGICAL:             true,
	kwLOGS:                true,
	kwLOW_PRIORITY:        true,
	kwMANUAL:              true,
	kwMAP:                 true,
	kwMAPPING:             true,
	kwMAPPINGS:            true,
	kwMASKING:             true,
	kwMATCH:               true,
	kwMATERIALIZED:        true,
	kwMAX:                 true,
	kwMERGE:               true,
	kwMETA:                true,
	kwMIN:                 true,
	kwMINUS:               true,
	kwMINUTE:              true,
	kwMINUTES:             true,
	kwMODE:                true,
	kwMODIFY:              true,
	kwMONTH:               true,
	kwNAME:                true,
	kwNAMES:               true,
	kwNEGATIVE:            true,
	kwNGRAMBF:             true,
	kwNO:                  true,
	kwNODE:                true,
	kwNODES:               true,
	kwNONE:                true,
	kwNULLS:               true,
	kwNUMBER:              true,
	kwNUMERIC:             true,
	kwOBSERVER:            true,
	kwOF:                  true,
	kwOFF:                 true,
	kwOFFSET:              true,
	kwONLY:                true,
	kwOPEN:                true,
	kwOPERATE:             true,
	kwOPTIMIZER:           true,
	kwOPTION:              true,
	kwOVERWRITE:           true,
	kwPARTITIONS:          true,
	kwPASSWORD:            true,
	kwPATH:                true,
	kwPAUSE:               true,
	kwPENDING:             true,
	kwPERCENTILE_UNION:    true,
	kwPERCENT_RANK:        true,
	kwPERSISTENT:          true,
	kwPIPE:                true,
	kwPIPES:               true,
	kwPIVOT:               true,
	kwPLAN:                true,
	kwPLUGIN:              true,
	kwPLUGINS:             true,
	kwPOLICIES:            true,
	kwPOLICY:              true,
	kwPRECEDING:           true,
	kwPRIORITY:            true,
	kwPRIVILEGES:          true,
	kwPROBABILITY:         true,
	kwPROC:                true,
	kwPROCESSLIST:         true,
	kwPROFILE:             true,
	kwPROFILELIST:         true,
	kwPROPERTIES:          true,
	kwPROPERTY:            true,
	kwQUALIFY:             true,
	kwQUARTER:             true,
	kwQUERIES:             true,
	kwQUERY:               true,
	kwQUEUE:               true,
	kwQUOTA:               true,
	kwRANDOM:              true,
	kwRANK:                true,
	kwREASON:              true,
	kwRECOVER:             true,
	kwREFRESH:             true,
	kwREMOVE:              true,
	kwREPAIR:              true,
	kwREPEATABLE:          true,
	kwREPLACE_IF_NOT_NULL: true,
	kwREPLICA:             true,
	kwREPOSITORIES:        true,
	kwREPOSITORY:          true,
	kwRESOURCE:            true,
	kwRESOURCES:           true,
	kwRESTORE:             true,
	kwRESUME:              true,
	kwRETAIN:              true,
	kwRETENTION:           true,
	kwRETRY:               true,
	kwRETURNS:             true,
	kwREVERT:              true,
	kwREWRITE:             true,
	kwROLE:                true,
	kwROLES:               true,
	kwROLLBACK:            true,
	kwROLLUP:              true,
	kwROUTINE:             true,
	kwROW:                 true,
	kwRULE:                true,
	kwRULES:               true,
	kwRUNNING:             true,
	kwSAMPLE:              true,
	kwSCHEDULE:            true,
	kwSCHEDULER:           true,
	kwSECOND:              true,
	kwSECURITY:            true,
	kwSEMI:                true,
	kwSEPARATOR:           true,
	kwSERIALIZABLE:        true,
	kwSESSION:             true,
	kwSETS:                true,
	kwSIGNED:              true,
	kwSKIP_HEADER:         true,
	kwSNAPSHOT:            true,
	kwSNAPSHOTS:           true,
	kwSQLBLACKLIST:        true,
	kwSTART:               true,
	kwSTATS:               true,
	kwSTATUS:              true,
	kwSTOP:                true,
	kwSTORAGE:             true,
	kwSTREAM:              true,
	kwSTRING:              true,
	kwSTRUCT:              true,
	kwSUBMIT:              true,
	kwSUM:                 true,
	kwSUSPEND:             true,
	kwSWAP:                true,
	kwSYNC:                true,
	kwSYSTEM_TIME:         true,
	kwTABLES:              true,
	kwTABLET:              true,
	kwTABLETS:             true,
	kwTAG:                 true,
	kwTASK:                true,
	kwTEMPORARY:           true,
	kwTHAN:                true,
	kwTIME:                true,
	kwTIMES:               true,
	kwTIMESTAMP:           true,
	kwTIMESTAMPADD:        true,
	kwTIMESTAMPDIFF:       true,
	kwTRACE:               true,
	kwTRANSACTION:         true,
	kwTRANSLATE:           true,
	kwTRIGGERS:            true,
	kwTRIM_SPACE:          true,
	kwTRUNCATE:            true,
	kwTYPE:                true,
	kwTYPES:               true,
	kwUNBOUNDED:           true,
	kwUNCOMMITTED:         true,
	kwUNINSTALL:           true,
	kwUNLOCK:              true,
	kwUNSET:               true,
	kwUSAGE:               true,
	kwUSER:                true,
	kwUSERS:               true,
	kwVALUE:               true,
	kwVARBINARY:           true,
	kwVARIABLES:           true,
	kwVECTOR:              true,
	kwVERBOSE:             true,
	kwVERSION:             true,
	kwVIEW:                true,
	kwVIEWS:               true,
	kwVOLUME:              true,
	kwVOLUMES:             true,
	kwWAREHOUSE:           true,
	kwWAREHOUSES:          true,
	kwWARNINGS:            true,
	kwWEEK:                true,
	kwWHITELIST:           true,
	kwWORK:                true,
	kwWRITE:               true,
	kwYEAR:                true,
}

// IsReserved reports whether kind is a reserved keyword that cannot be used as
// an unquoted identifier.
func IsReserved(kind TokenKind) bool {
	if kind >= 700 && dorisOnlyKeywordKinds[kind] != "" {
		return false
	}
	return kind >= 700 && !nonReservedKeywords[kind]
}

// KeywordToken returns the keyword TokenKind for a name, or (0, false)
// if name is not a keyword. Lookup is case-insensitive.
func KeywordToken(name string) (TokenKind, bool) {
	lower := strings.ToLower(name)
	if kind, ok := keywordMap[lower]; ok {
		return kind, true
	}
	kind, ok := dorisOnlyKeywords[lower]
	return kind, ok
}

// tokenNames maps a TokenKind to its display name. It is built at package
// init and only read afterwards, so concurrent callers of TokenName need no
// locking.
var tokenNames = buildTokenNames()

func buildTokenNames() map[TokenKind]string {
	m := make(map[TokenKind]string, len(keywordMap)+len(dorisOnlyKeywords)+20)
	for name, k := range keywordMap {
		m[k] = strings.ToUpper(name)
	}
	for kind, word := range dorisOnlyKeywordKinds {
		m[kind] = strings.ToUpper(word)
	}
	m[tokEOF] = "EOF"
	m[tokInvalid] = "INVALID"
	m[tokInt] = "INT"
	m[tokFloat] = "FLOAT"
	m[tokString] = "STRING"
	m[tokIdent] = "IDENT"
	m[tokQuotedIdent] = "QUOTED_IDENT"
	m[tokDotIdent] = "DOT_IDENTIFIER"
	m[tokHexLiteral] = "HEX_LITERAL"
	m[tokBitLiteral] = "BIT_LITERAL"
	m[tokPlaceholder] = "PLACEHOLDER"
	m[tokLessEq] = "<="
	m[tokGreaterEq] = ">="
	m[tokNotEq] = "<>"
	m[tokNullSafeEq] = "<=>"
	m[tokLogicalAnd] = "&&"
	m[tokDoublePipes] = "||"
	m[tokShiftLeft] = "<<"
	m[tokShiftRight] = ">>"
	m[tokArrow] = "->"
	m[tokDoubleAt] = "@@"
	m[tokAssign] = ":="
	m[tokDotDotDot] = "..."
	m[tokHintStart] = "HINT_START"
	m[tokHintEnd] = "HINT_END"
	return m
}

// TokenName returns a human-readable name for a TokenKind.
func TokenName(kind TokenKind) string {
	if name, ok := tokenNames[kind]; ok {
		return name
	}
	if kind > 0 && kind < 128 {
		return string(rune(kind))
	}
	return "UNKNOWN"
}
