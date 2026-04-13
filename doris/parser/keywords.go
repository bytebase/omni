package parser

import "strings"

// Keyword token constants. Values start at 700 and grow with iota.
// Extracted from /Users/h3n4l/OpenSource/parser/doris/DorisLexer.g4 (531 keywords).
// Numeric values are NOT stable across edits — do not persist them.
const (
	kwACCOUNT_LOCK   TokenKind = 700 + iota
	kwACCOUNT_UNLOCK            // ACCOUNT_UNLOCK
	kwACTIONS
	kwADD
	kwADMIN
	kwAFTER
	kwAGG_STATE
	kwAGGREGATE
	kwALIAS
	kwALL
	kwALTER
	kwALWAYS
	kwANALYZE
	kwANALYZED
	kwANALYZER
	kwAND
	kwANN
	kwANTI
	kwAPPEND
	kwARRAY
	kwAS
	kwASC
	kwAT
	kwAUTHORS
	kwAUTO
	kwAUTO_INCREMENT
	kwBACKEND
	kwBACKENDS
	kwBACKUP
	kwBEGIN
	kwBELONG
	kwBETWEEN
	kwBIGINT
	kwBIN
	kwBINARY
	kwBINLOG
	kwBITAND
	kwBITMAP
	kwBITMAP_EMPTY
	kwBITMAP_UNION
	kwBITOR
	kwBITXOR
	kwBLOB
	kwBOOLEAN
	kwBOTH
	kwBRANCH
	kwBRIEF
	kwBROKER
	kwBUCKETS
	kwBUILD
	kwBUILTIN
	kwBULK
	kwBY
	kwCACHE
	kwCACHED
	kwCALL
	kwCANCEL
	kwCASE
	kwCAST
	kwCATALOG
	kwCATALOGS
	kwCHAIN
	kwCHAR // also matches 'CHARACTER'
	kwCHAR_FILTER
	kwCHARSET
	kwCHECK
	kwCLEAN
	kwCLUSTER
	kwCLUSTERS
	kwCOLLATE
	kwCOLLATION
	kwCOLLECT
	kwCOLOCATE
	kwCOLUMN
	kwCOLUMNS
	kwCOMMENT
	kwCOMMIT
	kwCOMMITTED
	kwCOMPACT
	kwCOMPLETE
	kwCOMPRESS_TYPE
	kwCOMPUTE
	kwCONDITIONS
	kwCONFIG
	kwCONNECTION
	kwCONNECTION_ID
	kwCONSISTENT
	kwCONSTRAINT
	kwCONSTRAINTS
	kwCONVERT
	kwCONVERT_LSC // matches 'CONVERT_LIGHT_SCHEMA_CHANGE_PROCESS'
	kwCOPY
	kwCOUNT
	kwCREATE
	kwCREATION
	kwCRON
	kwCROSS
	kwCUBE
	kwCURRENT
	kwCURRENT_CATALOG
	kwCURRENT_DATE
	kwCURRENT_TIME
	kwCURRENT_TIMESTAMP
	kwCURRENT_USER
	kwDATA
	kwDATABASE
	kwDATABASES
	kwDATE
	kwDATETIME
	kwDATETIMEV1
	kwDATETIMEV2
	kwDATEV1
	kwDATEV2
	kwDAY
	kwDAY_HOUR
	kwDAY_SECOND
	kwDAYS
	kwDECIMAL
	kwDECIMALV2
	kwDECIMALV3
	kwDECOMMISSION
	kwDEFAULT
	kwDEFERRED
	kwDELETE
	kwDEMAND
	kwDESC
	kwDESCRIBE
	kwDIAGNOSE
	kwDIAGNOSIS
	kwDICTIONARIES
	kwDICTIONARY
	kwDISK
	kwDISTINCT
	kwDISTINCTPC
	kwDISTINCTPCSA
	kwDISTRIBUTED
	kwDISTRIBUTION
	kwDIV
	kwDO
	kwDORIS_INTERNAL_TABLE_ID
	kwDOUBLE
	kwDROP
	kwDROPP
	kwDUAL
	kwDUMP
	kwDUPLICATE
	kwDYNAMIC
	kwE
	kwELSE
	kwENABLE
	kwENCRYPTION
	kwENCRYPTKEY
	kwENCRYPTKEYS
	kwEND
	kwENDS
	kwENGINE
	kwENGINES
	kwENTER
	kwERRORS
	kwESCAPE
	kwEVENTS
	kwEVERY
	kwEXCEPT
	kwEXCLUDE
	kwEXECUTE
	kwEXISTS
	kwEXPIRED
	kwEXPLAIN
	kwEXPORT
	kwEXTENDED
	kwEXTERNAL
	kwEXTRACT
	kwFAILED_LOGIN_ATTEMPTS
	kwFALSE
	kwFAST
	kwFEATURE
	kwFIELDS
	kwFILE
	kwFILTER
	kwFIRST
	kwFLOAT
	kwFOLLOWER
	kwFOLLOWING
	kwFOR
	kwFOREIGN
	kwFORCE
	kwFORMAT
	kwFREE
	kwFROM
	kwFRONTEND
	kwFRONTENDS
	kwFULL
	kwFUNCTION
	kwFUNCTIONS
	kwGENERATED
	kwGENERIC
	kwGLOBAL
	kwGRANT
	kwGRANTS
	kwGRAPH
	kwGROUP
	kwGROUP_CONCAT
	kwGROUPING
	kwGROUPS
	kwHASH
	kwHASH_MAP
	kwHAVING
	kwHDFS
	kwHELP
	kwHISTOGRAM
	kwHLL
	kwHLL_UNION
	kwHOSTNAME
	kwHOTSPOT
	kwHOUR
	kwHOURS
	kwHUB
	kwIDENTIFIED
	kwIF
	kwIGNORE
	kwIMMEDIATE
	kwIN
	kwINCREMENTAL
	kwINDEX
	kwINDEXES
	kwINFILE
	kwINNER
	kwINSERT
	kwINSTALL
	kwINT
	kwINTEGER
	kwINTERMEDIATE
	kwINTERSECT
	kwINTERVAL
	kwINTO
	kwINVERTED
	kwIP_TRIE
	kwIPV4
	kwIPV6
	kwIS
	kwIS_NOT_NULL_PRED
	kwIS_NULL_PRED
	kwISNULL
	kwISOLATION
	kwJOB
	kwJOBS
	kwJOIN
	kwJSON
	kwJSONB
	kwKEY
	kwKEYS
	kwKILL
	kwLABEL
	kwLARGEINT
	kwLAST
	kwLATERAL
	kwLAYOUT
	kwLDAP
	kwLDAP_ADMIN_PASSWORD
	kwLEADING
	kwLEFT
	kwLESS
	kwLEVEL
	kwLIKE
	kwLIMIT
	kwLINES
	kwLINK
	kwLIST
	kwLOAD
	kwLOCAL
	kwLOCALTIME
	kwLOCALTIMESTAMP
	kwLOCATION
	kwLOCK
	kwLOGICAL
	kwLOW_PRIORITY
	kwMANUAL
	kwMAP
	kwMATCH
	kwMATCH_ALL
	kwMATCH_ANY
	kwMATCH_NAME
	kwMATCH_NAME_GLOB
	kwMATCH_PHRASE
	kwMATCH_PHRASE_EDGE
	kwMATCH_PHRASE_PREFIX
	kwMATCH_REGEXP
	kwMATCHED
	kwMATERIALIZED
	kwMAX
	kwMAXVALUE
	kwMEMO
	kwMERGE
	kwMID
	kwMIGRATE
	kwMIGRATIONS
	kwMIN
	kwMINUS
	kwMINUTE
	kwMINUTE_SECOND
	kwMINUTES
	kwMODIFY
	kwMONTH
	kwMTMV
	kwNAME
	kwNAMES
	kwNATURAL
	kwNEGATIVE
	kwNEVER
	kwNEXT
	kwNGRAM_BF
	kwNO
	kwNO_USE_MV
	kwNON_NULLABLE
	kwNOT
	kwNULL
	kwNULLS
	kwOBSERVER
	kwOF
	kwOFF
	kwOFFSET
	kwON
	kwONLY
	kwOPEN
	kwOPTIMIZE
	kwOPTIMIZED
	kwOR
	kwORDER
	kwOUTER
	kwOUTFILE
	kwOVER
	kwOVERWRITE
	kwPARAMETER
	kwPARSED
	kwPARTITION
	kwPARTITIONS
	kwPASSWORD
	kwPASSWORD_EXPIRE
	kwPASSWORD_HISTORY
	kwPASSWORD_LOCK_TIME
	kwPASSWORD_REUSE
	kwPATH
	kwPAUSE
	kwPERCENT
	kwPERIOD
	kwPERMISSIVE
	kwPHYSICAL
	kwPI
	kwPLAN
	kwPLAY
	kwPLUGIN
	kwPLUGINS
	kwPOLICY
	kwPOSITION
	kwPRECEDING
	kwPREPARE
	kwPRIMARY
	kwPRIVILEGES
	kwPROC
	kwPROCEDURE
	kwPROCESS
	kwPROCESSLIST
	kwPROFILE
	kwPROPERTIES
	kwPROPERTY
	kwQUALIFY
	kwQUANTILE_STATE
	kwQUANTILE_UNION
	kwQUARTER
	kwQUERY
	kwQUEUED
	kwQUOTA
	kwRANDOM
	kwRANGE
	kwREAD
	kwREAL
	kwREBALANCE
	kwRECENT
	kwRECOVER
	kwRECYCLE
	kwREFERENCES
	kwREFRESH
	kwREGEXP
	kwRELEASE
	kwRENAME
	kwREPAIR
	kwREPEATABLE
	kwREPLACE
	kwREPLACE_IF_NOT_NULL
	kwREPLAYER
	kwREPLICA
	kwREPOSITORIES
	kwREPOSITORY
	kwRESOURCE
	kwRESOURCES
	kwRESTORE
	kwRESTRICTIVE
	kwRESUME
	kwRETAIN
	kwRETENTION
	kwRETURNS
	kwREVOKE
	kwREWRITTEN
	kwRIGHT
	kwRLIKE
	kwROLE
	kwROLES
	kwROLLBACK
	kwROLLUP
	kwROOT
	kwROTATE
	kwROUTINE
	kwROW
	kwROWS
	kwS3
	kwSAMPLE
	kwSCHEDULE
	kwSCHEDULER
	kwSCHEMA
	kwSCHEMAS
	kwSECOND
	kwSELECT
	kwSEMI
	kwSEPARATOR
	kwSERIALIZABLE
	kwSESSION
	kwSESSION_USER
	kwSET
	kwSET_SESSION_VARIABLE
	kwSETS
	kwSHAPE
	kwSHOW
	kwSIGNED
	kwSKEW
	kwSMALLINT
	kwSNAPSHOT
	kwSNAPSHOTS
	kwSONAME
	kwSPLIT
	kwSQL
	kwSQL_BLOCK_RULE
	kwSTAGE
	kwSTAGES
	kwSTART
	kwSTARTS
	kwSTATS
	kwSTATUS
	kwSTOP
	kwSTORAGE
	kwSTREAM
	kwSTREAMING
	kwSTRING
	kwSTRUCT
	kwSUBSTR
	kwSUBSTRING
	kwSUM
	kwSUPERUSER
	kwSWITCH
	kwSYNC
	kwSYSTEM
	kwTABLE
	kwTABLES
	kwTABLESAMPLE
	kwTABLET
	kwTABLETS
	kwTAG
	kwTASK
	kwTASKS
	kwTDE
	kwTEMPORARY
	kwTERMINATED
	kwTEXT
	kwTHAN
	kwTHEN
	kwTIME
	kwTIMESTAMP
	kwTINYINT
	kwTO
	kwTOKEN_FILTER
	kwTOKENIZER
	kwTRAILING
	kwTRANSACTION
	kwTRASH
	kwTREE
	kwTRIGGERS
	kwTRIM
	kwTRUE
	kwTRUNCATE
	kwTRY_CAST
	kwTYPE
	kwTYPECAST // matches 'TYPE_CAST'
	kwTYPES
	kwUNBOUNDED
	kwUNCOMMITTED
	kwUNINSTALL
	kwUNION
	kwUNIQUE
	kwUNLOCK
	kwUNSET
	kwUNSIGNED
	kwUP
	kwUPDATE
	kwUSE
	kwUSE_MV
	kwUSER
	kwUSING
	kwVALUE
	kwVALUES
	kwVARBINARY
	kwVARCHAR
	kwVARIABLE
	kwVARIABLES
	kwVARIANT
	kwVAULT
	kwVAULTS
	kwVERBOSE
	kwVERSION
	kwVIEW
	kwVIEWS
	kwWARM
	kwWARNINGS
	kwWEEK
	kwWHEN
	kwWHERE
	kwWHITELIST
	kwWITH
	kwWORK
	kwWORKLOAD
	kwWRITE
	kwXOR
	kwYEAR
)

// keywordMap maps lowercase keyword strings to their token kinds.
// Case-insensitive lookup is performed by lowercasing the input in KeywordToken.
var keywordMap = map[string]TokenKind{
	"account_lock":                         kwACCOUNT_LOCK,
	"account_unlock":                       kwACCOUNT_UNLOCK,
	"actions":                              kwACTIONS,
	"add":                                  kwADD,
	"admin":                                kwADMIN,
	"after":                                kwAFTER,
	"agg_state":                            kwAGG_STATE,
	"aggregate":                            kwAGGREGATE,
	"alias":                                kwALIAS,
	"all":                                  kwALL,
	"alter":                                kwALTER,
	"always":                               kwALWAYS,
	"analyze":                              kwANALYZE,
	"analyzed":                             kwANALYZED,
	"analyzer":                             kwANALYZER,
	"and":                                  kwAND,
	"ann":                                  kwANN,
	"anti":                                 kwANTI,
	"append":                               kwAPPEND,
	"array":                                kwARRAY,
	"as":                                   kwAS,
	"asc":                                  kwASC,
	"at":                                   kwAT,
	"authors":                              kwAUTHORS,
	"auto":                                 kwAUTO,
	"auto_increment":                       kwAUTO_INCREMENT,
	"backend":                              kwBACKEND,
	"backends":                             kwBACKENDS,
	"backup":                               kwBACKUP,
	"begin":                                kwBEGIN,
	"belong":                               kwBELONG,
	"between":                              kwBETWEEN,
	"bigint":                               kwBIGINT,
	"bin":                                  kwBIN,
	"binary":                               kwBINARY,
	"binlog":                               kwBINLOG,
	"bitand":                               kwBITAND,
	"bitmap":                               kwBITMAP,
	"bitmap_empty":                         kwBITMAP_EMPTY,
	"bitmap_union":                         kwBITMAP_UNION,
	"bitor":                                kwBITOR,
	"bitxor":                               kwBITXOR,
	"blob":                                 kwBLOB,
	"boolean":                              kwBOOLEAN,
	"both":                                 kwBOTH,
	"branch":                               kwBRANCH,
	"brief":                                kwBRIEF,
	"broker":                               kwBROKER,
	"buckets":                              kwBUCKETS,
	"build":                                kwBUILD,
	"builtin":                              kwBUILTIN,
	"bulk":                                 kwBULK,
	"by":                                   kwBY,
	"cache":                                kwCACHE,
	"cached":                               kwCACHED,
	"call":                                 kwCALL,
	"cancel":                               kwCANCEL,
	"case":                                 kwCASE,
	"cast":                                 kwCAST,
	"catalog":                              kwCATALOG,
	"catalogs":                             kwCATALOGS,
	"chain":                                kwCHAIN,
	"char":                                 kwCHAR,
	"character":                            kwCHAR, // alias for CHAR
	"char_filter":                          kwCHAR_FILTER,
	"charset":                              kwCHARSET,
	"check":                                kwCHECK,
	"clean":                                kwCLEAN,
	"cluster":                              kwCLUSTER,
	"clusters":                             kwCLUSTERS,
	"collate":                              kwCOLLATE,
	"collation":                            kwCOLLATION,
	"collect":                              kwCOLLECT,
	"colocate":                             kwCOLOCATE,
	"column":                               kwCOLUMN,
	"columns":                              kwCOLUMNS,
	"comment":                              kwCOMMENT,
	"commit":                               kwCOMMIT,
	"committed":                            kwCOMMITTED,
	"compact":                              kwCOMPACT,
	"complete":                             kwCOMPLETE,
	"compress_type":                        kwCOMPRESS_TYPE,
	"compute":                              kwCOMPUTE,
	"conditions":                           kwCONDITIONS,
	"config":                               kwCONFIG,
	"connection":                           kwCONNECTION,
	"connection_id":                        kwCONNECTION_ID,
	"consistent":                           kwCONSISTENT,
	"constraint":                           kwCONSTRAINT,
	"constraints":                          kwCONSTRAINTS,
	"convert":                              kwCONVERT,
	"convert_light_schema_change_process":  kwCONVERT_LSC,
	"copy":                                 kwCOPY,
	"count":                                kwCOUNT,
	"create":                               kwCREATE,
	"creation":                             kwCREATION,
	"cron":                                 kwCRON,
	"cross":                                kwCROSS,
	"cube":                                 kwCUBE,
	"current":                              kwCURRENT,
	"current_catalog":                      kwCURRENT_CATALOG,
	"current_date":                         kwCURRENT_DATE,
	"current_time":                         kwCURRENT_TIME,
	"current_timestamp":                    kwCURRENT_TIMESTAMP,
	"current_user":                         kwCURRENT_USER,
	"data":                                 kwDATA,
	"database":                             kwDATABASE,
	"databases":                            kwDATABASES,
	"date":                                 kwDATE,
	"datetime":                             kwDATETIME,
	"datetimev1":                           kwDATETIMEV1,
	"datetimev2":                           kwDATETIMEV2,
	"datev1":                               kwDATEV1,
	"datev2":                               kwDATEV2,
	"day":                                  kwDAY,
	"day_hour":                             kwDAY_HOUR,
	"day_second":                           kwDAY_SECOND,
	"days":                                 kwDAYS,
	"decimal":                              kwDECIMAL,
	"decimalv2":                            kwDECIMALV2,
	"decimalv3":                            kwDECIMALV3,
	"decommission":                         kwDECOMMISSION,
	"default":                              kwDEFAULT,
	"deferred":                             kwDEFERRED,
	"delete":                               kwDELETE,
	"demand":                               kwDEMAND,
	"desc":                                 kwDESC,
	"describe":                             kwDESCRIBE,
	"diagnose":                             kwDIAGNOSE,
	"diagnosis":                            kwDIAGNOSIS,
	"dictionaries":                         kwDICTIONARIES,
	"dictionary":                           kwDICTIONARY,
	"disk":                                 kwDISK,
	"distinct":                             kwDISTINCT,
	"distinctpc":                           kwDISTINCTPC,
	"distinctpcsa":                         kwDISTINCTPCSA,
	"distributed":                          kwDISTRIBUTED,
	"distribution":                         kwDISTRIBUTION,
	"div":                                  kwDIV,
	"do":                                   kwDO,
	"doris_internal_table_id":              kwDORIS_INTERNAL_TABLE_ID,
	"double":                               kwDOUBLE,
	"drop":                                 kwDROP,
	"dropp":                                kwDROPP,
	"dual":                                 kwDUAL,
	"dump":                                 kwDUMP,
	"duplicate":                            kwDUPLICATE,
	"dynamic":                              kwDYNAMIC,
	"e":                                    kwE,
	"else":                                 kwELSE,
	"enable":                               kwENABLE,
	"encryption":                           kwENCRYPTION,
	"encryptkey":                           kwENCRYPTKEY,
	"encryptkeys":                          kwENCRYPTKEYS,
	"end":                                  kwEND,
	"ends":                                 kwENDS,
	"engine":                               kwENGINE,
	"engines":                              kwENGINES,
	"enter":                                kwENTER,
	"errors":                               kwERRORS,
	"escape":                               kwESCAPE,
	"events":                               kwEVENTS,
	"every":                                kwEVERY,
	"except":                               kwEXCEPT,
	"exclude":                              kwEXCLUDE,
	"execute":                              kwEXECUTE,
	"exists":                               kwEXISTS,
	"expired":                              kwEXPIRED,
	"explain":                              kwEXPLAIN,
	"export":                               kwEXPORT,
	"extended":                             kwEXTENDED,
	"external":                             kwEXTERNAL,
	"extract":                              kwEXTRACT,
	"failed_login_attempts":                kwFAILED_LOGIN_ATTEMPTS,
	"false":                                kwFALSE,
	"fast":                                 kwFAST,
	"feature":                              kwFEATURE,
	"fields":                               kwFIELDS,
	"file":                                 kwFILE,
	"filter":                               kwFILTER,
	"first":                                kwFIRST,
	"float":                                kwFLOAT,
	"follower":                             kwFOLLOWER,
	"following":                            kwFOLLOWING,
	"for":                                  kwFOR,
	"foreign":                              kwFOREIGN,
	"force":                                kwFORCE,
	"format":                               kwFORMAT,
	"free":                                 kwFREE,
	"from":                                 kwFROM,
	"frontend":                             kwFRONTEND,
	"frontends":                            kwFRONTENDS,
	"full":                                 kwFULL,
	"function":                             kwFUNCTION,
	"functions":                            kwFUNCTIONS,
	"generated":                            kwGENERATED,
	"generic":                              kwGENERIC,
	"global":                               kwGLOBAL,
	"grant":                                kwGRANT,
	"grants":                               kwGRANTS,
	"graph":                                kwGRAPH,
	"group":                                kwGROUP,
	"group_concat":                         kwGROUP_CONCAT,
	"grouping":                             kwGROUPING,
	"groups":                               kwGROUPS,
	"hash":                                 kwHASH,
	"hash_map":                             kwHASH_MAP,
	"having":                               kwHAVING,
	"hdfs":                                 kwHDFS,
	"help":                                 kwHELP,
	"histogram":                            kwHISTOGRAM,
	"hll":                                  kwHLL,
	"hll_union":                            kwHLL_UNION,
	"hostname":                             kwHOSTNAME,
	"hotspot":                              kwHOTSPOT,
	"hour":                                 kwHOUR,
	"hours":                                kwHOURS,
	"hub":                                  kwHUB,
	"identified":                           kwIDENTIFIED,
	"if":                                   kwIF,
	"ignore":                               kwIGNORE,
	"immediate":                            kwIMMEDIATE,
	"in":                                   kwIN,
	"incremental":                          kwINCREMENTAL,
	"index":                                kwINDEX,
	"indexes":                              kwINDEXES,
	"infile":                               kwINFILE,
	"inner":                                kwINNER,
	"insert":                               kwINSERT,
	"install":                              kwINSTALL,
	"int":                                  kwINT,
	"integer":                              kwINTEGER,
	"intermediate":                         kwINTERMEDIATE,
	"intersect":                            kwINTERSECT,
	"interval":                             kwINTERVAL,
	"into":                                 kwINTO,
	"inverted":                             kwINVERTED,
	"ip_trie":                              kwIP_TRIE,
	"ipv4":                                 kwIPV4,
	"ipv6":                                 kwIPV6,
	"is":                                   kwIS,
	"is_not_null_pred":                     kwIS_NOT_NULL_PRED,
	"is_null_pred":                         kwIS_NULL_PRED,
	"isnull":                               kwISNULL,
	"isolation":                            kwISOLATION,
	"job":                                  kwJOB,
	"jobs":                                 kwJOBS,
	"join":                                 kwJOIN,
	"json":                                 kwJSON,
	"jsonb":                                kwJSONB,
	"key":                                  kwKEY,
	"keys":                                 kwKEYS,
	"kill":                                 kwKILL,
	"label":                                kwLABEL,
	"largeint":                             kwLARGEINT,
	"last":                                 kwLAST,
	"lateral":                              kwLATERAL,
	"layout":                               kwLAYOUT,
	"ldap":                                 kwLDAP,
	"ldap_admin_password":                  kwLDAP_ADMIN_PASSWORD,
	"leading":                              kwLEADING,
	"left":                                 kwLEFT,
	"less":                                 kwLESS,
	"level":                                kwLEVEL,
	"like":                                 kwLIKE,
	"limit":                                kwLIMIT,
	"lines":                                kwLINES,
	"link":                                 kwLINK,
	"list":                                 kwLIST,
	"load":                                 kwLOAD,
	"local":                                kwLOCAL,
	"localtime":                            kwLOCALTIME,
	"localtimestamp":                       kwLOCALTIMESTAMP,
	"location":                             kwLOCATION,
	"lock":                                 kwLOCK,
	"logical":                              kwLOGICAL,
	"low_priority":                         kwLOW_PRIORITY,
	"manual":                               kwMANUAL,
	"map":                                  kwMAP,
	"match":                                kwMATCH,
	"match_all":                            kwMATCH_ALL,
	"match_any":                            kwMATCH_ANY,
	"match_name":                           kwMATCH_NAME,
	"match_name_glob":                      kwMATCH_NAME_GLOB,
	"match_phrase":                         kwMATCH_PHRASE,
	"match_phrase_edge":                    kwMATCH_PHRASE_EDGE,
	"match_phrase_prefix":                  kwMATCH_PHRASE_PREFIX,
	"match_regexp":                         kwMATCH_REGEXP,
	"matched":                              kwMATCHED,
	"materialized":                         kwMATERIALIZED,
	"max":                                  kwMAX,
	"maxvalue":                             kwMAXVALUE,
	"memo":                                 kwMEMO,
	"merge":                                kwMERGE,
	"mid":                                  kwMID,
	"migrate":                              kwMIGRATE,
	"migrations":                           kwMIGRATIONS,
	"min":                                  kwMIN,
	"minus":                                kwMINUS,
	"minute":                               kwMINUTE,
	"minute_second":                        kwMINUTE_SECOND,
	"minutes":                              kwMINUTES,
	"modify":                               kwMODIFY,
	"month":                                kwMONTH,
	"mtmv":                                 kwMTMV,
	"name":                                 kwNAME,
	"names":                                kwNAMES,
	"natural":                              kwNATURAL,
	"negative":                             kwNEGATIVE,
	"never":                                kwNEVER,
	"next":                                 kwNEXT,
	"ngram_bf":                             kwNGRAM_BF,
	"no":                                   kwNO,
	"no_use_mv":                            kwNO_USE_MV,
	"non_nullable":                         kwNON_NULLABLE,
	"not":                                  kwNOT,
	"null":                                 kwNULL,
	"nulls":                                kwNULLS,
	"observer":                             kwOBSERVER,
	"of":                                   kwOF,
	"off":                                  kwOFF,
	"offset":                               kwOFFSET,
	"on":                                   kwON,
	"only":                                 kwONLY,
	"open":                                 kwOPEN,
	"optimize":                             kwOPTIMIZE,
	"optimized":                            kwOPTIMIZED,
	"or":                                   kwOR,
	"order":                                kwORDER,
	"outer":                                kwOUTER,
	"outfile":                              kwOUTFILE,
	"over":                                 kwOVER,
	"overwrite":                            kwOVERWRITE,
	"parameter":                            kwPARAMETER,
	"parsed":                               kwPARSED,
	"partition":                            kwPARTITION,
	"partitions":                           kwPARTITIONS,
	"password":                             kwPASSWORD,
	"password_expire":                      kwPASSWORD_EXPIRE,
	"password_history":                     kwPASSWORD_HISTORY,
	"password_lock_time":                   kwPASSWORD_LOCK_TIME,
	"password_reuse":                       kwPASSWORD_REUSE,
	"path":                                 kwPATH,
	"pause":                                kwPAUSE,
	"percent":                              kwPERCENT,
	"period":                               kwPERIOD,
	"permissive":                           kwPERMISSIVE,
	"physical":                             kwPHYSICAL,
	"pi":                                   kwPI,
	"plan":                                 kwPLAN,
	"play":                                 kwPLAY,
	"plugin":                               kwPLUGIN,
	"plugins":                              kwPLUGINS,
	"policy":                               kwPOLICY,
	"position":                             kwPOSITION,
	"preceding":                            kwPRECEDING,
	"prepare":                              kwPREPARE,
	"primary":                              kwPRIMARY,
	"privileges":                           kwPRIVILEGES,
	"proc":                                 kwPROC,
	"procedure":                            kwPROCEDURE,
	"process":                              kwPROCESS,
	"processlist":                          kwPROCESSLIST,
	"profile":                              kwPROFILE,
	"properties":                           kwPROPERTIES,
	"property":                             kwPROPERTY,
	"qualify":                              kwQUALIFY,
	"quantile_state":                       kwQUANTILE_STATE,
	"quantile_union":                       kwQUANTILE_UNION,
	"quarter":                              kwQUARTER,
	"query":                                kwQUERY,
	"queued":                               kwQUEUED,
	"quota":                                kwQUOTA,
	"random":                               kwRANDOM,
	"range":                                kwRANGE,
	"read":                                 kwREAD,
	"real":                                 kwREAL,
	"rebalance":                            kwREBALANCE,
	"recent":                               kwRECENT,
	"recover":                              kwRECOVER,
	"recycle":                              kwRECYCLE,
	"references":                           kwREFERENCES,
	"refresh":                              kwREFRESH,
	"regexp":                               kwREGEXP,
	"release":                              kwRELEASE,
	"rename":                               kwRENAME,
	"repair":                               kwREPAIR,
	"repeatable":                           kwREPEATABLE,
	"replace":                              kwREPLACE,
	"replace_if_not_null":                  kwREPLACE_IF_NOT_NULL,
	"replayer":                             kwREPLAYER,
	"replica":                              kwREPLICA,
	"repositories":                         kwREPOSITORIES,
	"repository":                           kwREPOSITORY,
	"resource":                             kwRESOURCE,
	"resources":                            kwRESOURCES,
	"restore":                              kwRESTORE,
	"restrictive":                          kwRESTRICTIVE,
	"resume":                               kwRESUME,
	"retain":                               kwRETAIN,
	"retention":                            kwRETENTION,
	"returns":                              kwRETURNS,
	"revoke":                               kwREVOKE,
	"rewritten":                            kwREWRITTEN,
	"right":                                kwRIGHT,
	"rlike":                                kwRLIKE,
	"role":                                 kwROLE,
	"roles":                                kwROLES,
	"rollback":                             kwROLLBACK,
	"rollup":                               kwROLLUP,
	"root":                                 kwROOT,
	"rotate":                               kwROTATE,
	"routine":                              kwROUTINE,
	"row":                                  kwROW,
	"rows":                                 kwROWS,
	"s3":                                   kwS3,
	"sample":                               kwSAMPLE,
	"schedule":                             kwSCHEDULE,
	"scheduler":                            kwSCHEDULER,
	"schema":                               kwSCHEMA,
	"schemas":                              kwSCHEMAS,
	"second":                               kwSECOND,
	"select":                               kwSELECT,
	"semi":                                 kwSEMI,
	"separator":                            kwSEPARATOR,
	"serializable":                         kwSERIALIZABLE,
	"session":                              kwSESSION,
	"session_user":                         kwSESSION_USER,
	"set":                                  kwSET,
	"set_session_variable":                 kwSET_SESSION_VARIABLE,
	"sets":                                 kwSETS,
	"shape":                                kwSHAPE,
	"show":                                 kwSHOW,
	"signed":                               kwSIGNED,
	"skew":                                 kwSKEW,
	"smallint":                             kwSMALLINT,
	"snapshot":                             kwSNAPSHOT,
	"snapshots":                            kwSNAPSHOTS,
	"soname":                               kwSONAME,
	"split":                                kwSPLIT,
	"sql":                                  kwSQL,
	"sql_block_rule":                       kwSQL_BLOCK_RULE,
	"stage":                                kwSTAGE,
	"stages":                               kwSTAGES,
	"start":                                kwSTART,
	"starts":                               kwSTARTS,
	"stats":                                kwSTATS,
	"status":                               kwSTATUS,
	"stop":                                 kwSTOP,
	"storage":                              kwSTORAGE,
	"stream":                               kwSTREAM,
	"streaming":                            kwSTREAMING,
	"string":                               kwSTRING,
	"struct":                               kwSTRUCT,
	"substr":                               kwSUBSTR,
	"substring":                            kwSUBSTRING,
	"sum":                                  kwSUM,
	"superuser":                            kwSUPERUSER,
	"switch":                               kwSWITCH,
	"sync":                                 kwSYNC,
	"system":                               kwSYSTEM,
	"table":                                kwTABLE,
	"tables":                               kwTABLES,
	"tablesample":                          kwTABLESAMPLE,
	"tablet":                               kwTABLET,
	"tablets":                              kwTABLETS,
	"tag":                                  kwTAG,
	"task":                                 kwTASK,
	"tasks":                                kwTASKS,
	"tde":                                  kwTDE,
	"temporary":                            kwTEMPORARY,
	"terminated":                           kwTERMINATED,
	"text":                                 kwTEXT,
	"than":                                 kwTHAN,
	"then":                                 kwTHEN,
	"time":                                 kwTIME,
	"timestamp":                            kwTIMESTAMP,
	"tinyint":                              kwTINYINT,
	"to":                                   kwTO,
	"token_filter":                         kwTOKEN_FILTER,
	"tokenizer":                            kwTOKENIZER,
	"trailing":                             kwTRAILING,
	"transaction":                          kwTRANSACTION,
	"trash":                                kwTRASH,
	"tree":                                 kwTREE,
	"triggers":                             kwTRIGGERS,
	"trim":                                 kwTRIM,
	"true":                                 kwTRUE,
	"truncate":                             kwTRUNCATE,
	"try_cast":                             kwTRY_CAST,
	"type":                                 kwTYPE,
	"type_cast":                            kwTYPECAST, // alias for TYPECAST
	"types":                                kwTYPES,
	"unbounded":                            kwUNBOUNDED,
	"uncommitted":                          kwUNCOMMITTED,
	"uninstall":                            kwUNINSTALL,
	"union":                                kwUNION,
	"unique":                               kwUNIQUE,
	"unlock":                               kwUNLOCK,
	"unset":                                kwUNSET,
	"unsigned":                             kwUNSIGNED,
	"up":                                   kwUP,
	"update":                               kwUPDATE,
	"use":                                  kwUSE,
	"use_mv":                               kwUSE_MV,
	"user":                                 kwUSER,
	"using":                                kwUSING,
	"value":                                kwVALUE,
	"values":                               kwVALUES,
	"varbinary":                            kwVARBINARY,
	"varchar":                              kwVARCHAR,
	"variable":                             kwVARIABLE,
	"variables":                            kwVARIABLES,
	"variant":                              kwVARIANT,
	"vault":                                kwVAULT,
	"vaults":                               kwVAULTS,
	"verbose":                              kwVERBOSE,
	"version":                              kwVERSION,
	"view":                                 kwVIEW,
	"views":                                kwVIEWS,
	"warm":                                 kwWARM,
	"warnings":                             kwWARNINGS,
	"week":                                 kwWEEK,
	"when":                                 kwWHEN,
	"where":                                kwWHERE,
	"whitelist":                            kwWHITELIST,
	"with":                                 kwWITH,
	"work":                                 kwWORK,
	"workload":                             kwWORKLOAD,
	"write":                                kwWRITE,
	"xor":                                  kwXOR,
	"year":                                 kwYEAR,
}

// nonReservedKeywords is the set of keywords that can be used as unquoted
// identifiers. Seeded from the nonReserved rule in DorisParser.g4 (lines 1906-2257).
// Note: COMMENT_START (/*), LEFT_BRACE ({), RIGHT_BRACE (}) and HINT_START/HINT_END
// appear in the nonReserved grammar rule but are operators/punctuation, not keywords
// handled here.
var nonReservedKeywords = map[TokenKind]bool{
	kwACTIONS:                true,
	kwAFTER:                  true,
	kwAGG_STATE:              true,
	kwAGGREGATE:              true,
	kwALIAS:                  true,
	kwALWAYS:                 true,
	kwANALYZED:               true,
	kwANN:                    true,
	kwARRAY:                  true,
	kwAT:                     true,
	kwAUTHORS:                true,
	kwAUTO_INCREMENT:         true,
	kwBACKENDS:               true,
	kwBACKUP:                 true,
	kwBEGIN:                  true,
	kwBELONG:                 true,
	kwBIN:                    true,
	kwBITAND:                 true,
	kwBITMAP:                 true,
	kwBITMAP_EMPTY:           true,
	kwBITMAP_UNION:           true,
	kwBITOR:                  true,
	kwBITXOR:                 true,
	kwBLOB:                   true,
	kwBOOLEAN:                true,
	kwBRANCH:                 true,
	kwBRIEF:                  true,
	kwBROKER:                 true,
	kwBUCKETS:                true,
	kwBUILD:                  true,
	kwBUILTIN:                true,
	kwBULK:                   true,
	kwCACHE:                  true,
	kwCACHED:                 true,
	kwCALL:                   true,
	kwCATALOG:                true,
	kwCATALOGS:               true,
	kwCHAIN:                  true,
	kwCHAR:                   true,
	kwCHARSET:                true,
	kwCHECK:                  true,
	kwCLUSTER:                true,
	kwCLUSTERS:               true,
	kwCOLLATION:              true,
	kwCOLLECT:                true,
	kwCOLOCATE:               true,
	kwCOLUMNS:                true,
	kwCOMMENT:                true,
	kwCOMMIT:                 true,
	kwCOMMITTED:              true,
	kwCOMPACT:                true,
	kwCOMPLETE:               true,
	kwCOMPRESS_TYPE:          true,
	kwCOMPUTE:                true,
	kwCONDITIONS:             true,
	kwCONFIG:                 true,
	kwCONNECTION:             true,
	kwCONNECTION_ID:          true,
	kwCONSISTENT:             true,
	kwCONSTRAINTS:            true,
	kwCONVERT:                true,
	kwCONVERT_LSC:            true,
	kwCOPY:                   true,
	kwCOUNT:                  true,
	kwCREATION:               true,
	kwCRON:                   true,
	kwCURRENT_CATALOG:        true,
	kwCURRENT_DATE:           true,
	kwCURRENT_TIME:           true,
	kwCURRENT_TIMESTAMP:      true,
	kwCURRENT_USER:           true,
	kwDATA:                   true,
	kwDATE:                   true,
	kwDATETIME:               true,
	kwDATETIMEV1:             true,
	kwDATETIMEV2:             true,
	kwDATEV1:                 true,
	kwDATEV2:                 true,
	kwDAY:                    true,
	kwDAYS:                   true,
	kwDECIMAL:                true,
	kwDECIMALV2:              true,
	kwDECIMALV3:              true,
	kwDEFERRED:               true,
	kwDEMAND:                 true,
	kwDIAGNOSE:               true,
	kwDIAGNOSIS:              true,
	kwDICTIONARIES:           true,
	kwDICTIONARY:             true,
	kwDISTINCTPC:             true,
	kwDISTINCTPCSA:           true,
	kwDO:                     true,
	kwDORIS_INTERNAL_TABLE_ID: true,
	kwDUAL:                   true,
	kwDYNAMIC:                true,
	kwE:                      true,
	kwENABLE:                 true,
	kwENCRYPTION:             true,
	kwENCRYPTKEY:             true,
	kwENCRYPTKEYS:            true,
	kwEND:                    true,
	kwENDS:                   true,
	kwENGINE:                 true,
	kwENGINES:                true,
	kwERRORS:                 true,
	kwESCAPE:                 true,
	kwEVENTS:                 true,
	kwEVERY:                  true,
	kwEXCLUDE:                true,
	kwEXPIRED:                true,
	kwEXTERNAL:               true,
	kwFAILED_LOGIN_ATTEMPTS:  true,
	kwFAST:                   true,
	kwFEATURE:                true,
	kwFIELDS:                 true,
	kwFILE:                   true,
	kwFILTER:                 true,
	kwFIRST:                  true,
	kwFORMAT:                 true,
	kwFREE:                   true,
	kwFRONTENDS:              true,
	kwFUNCTION:               true,
	kwGENERATED:              true,
	kwGENERIC:                true,
	kwGLOBAL:                 true,
	kwGRAPH:                  true,
	kwGROUPING:               true,
	kwGROUPS:                 true,
	kwGROUP_CONCAT:           true,
	kwHASH:                   true,
	kwHASH_MAP:               true,
	kwHDFS:                   true,
	kwHELP:                   true,
	kwHISTOGRAM:              true,
	kwHLL_UNION:              true,
	kwHOSTNAME:               true,
	kwHOTSPOT:                true,
	kwHOUR:                   true,
	kwHOURS:                  true,
	kwHUB:                    true,
	kwIDENTIFIED:             true,
	kwIGNORE:                 true,
	kwIMMEDIATE:              true,
	kwINCREMENTAL:            true,
	kwINDEXES:                true,
	kwINSERT:                 true,
	kwINVERTED:               true,
	kwIP_TRIE:                true,
	kwIPV4:                   true,
	kwIPV6:                   true,
	kwIS_NOT_NULL_PRED:       true,
	kwIS_NULL_PRED:           true,
	kwISNULL:                 true,
	kwISOLATION:              true,
	kwJOB:                    true,
	kwJOBS:                   true,
	kwJSON:                   true,
	kwJSONB:                  true,
	kwLABEL:                  true,
	kwLAST:                   true,
	kwLDAP:                   true,
	kwLDAP_ADMIN_PASSWORD:    true,
	kwLESS:                   true,
	kwLEVEL:                  true,
	kwLINES:                  true,
	kwLINK:                   true,
	kwLOCAL:                  true,
	kwLOCALTIME:              true,
	kwLOCALTIMESTAMP:         true,
	kwLOCATION:               true,
	kwLOCK:                   true,
	kwLOGICAL:                true,
	kwMANUAL:                 true,
	kwMAP:                    true,
	kwMATCHED:                true,
	kwMATCH_ALL:              true,
	kwMATCH_ANY:              true,
	kwMATCH_PHRASE:           true,
	kwMATCH_PHRASE_EDGE:      true,
	kwMATCH_PHRASE_PREFIX:    true,
	kwMATCH_REGEXP:           true,
	kwMATCH_NAME:             true,
	kwMATCH_NAME_GLOB:        true,
	kwMATERIALIZED:           true,
	kwMAX:                    true,
	kwMEMO:                   true,
	kwMERGE:                  true,
	kwMID:                    true,
	kwMIGRATE:                true,
	kwMIGRATIONS:             true,
	kwMIN:                    true,
	kwMINUTE:                 true,
	kwMINUTES:                true,
	kwMODIFY:                 true,
	kwMONTH:                  true,
	kwMTMV:                   true,
	kwNAME:                   true,
	kwNAMES:                  true,
	kwNEGATIVE:               true,
	kwNEVER:                  true,
	kwNEXT:                   true,
	kwNGRAM_BF:               true,
	kwNO:                     true,
	kwNON_NULLABLE:           true,
	kwNULLS:                  true,
	kwOF:                     true,
	kwOFF:                    true,
	kwOFFSET:                 true,
	kwONLY:                   true,
	kwOPEN:                   true,
	kwOPTIMIZE:               true,
	kwOPTIMIZED:              true,
	kwPARAMETER:              true,
	kwPARSED:                 true,
	kwPASSWORD:               true,
	kwPASSWORD_EXPIRE:        true,
	kwPASSWORD_HISTORY:       true,
	kwPASSWORD_LOCK_TIME:     true,
	kwPASSWORD_REUSE:         true,
	kwPARTITIONS:             true,
	kwPATH:                   true,
	kwPAUSE:                  true,
	kwPERCENT:                true,
	kwPERIOD:                 true,
	kwPERMISSIVE:             true,
	kwPHYSICAL:               true,
	kwPI:                     true,
	kwPLAN:                   true,
	kwPLUGIN:                 true,
	kwPLUGINS:                true,
	kwPOLICY:                 true,
	kwPOSITION:               true,
	kwPRIVILEGES:             true,
	kwPROC:                   true,
	kwPROCESS:                true,
	kwPROCESSLIST:            true,
	kwPROFILE:                true,
	kwPROPERTIES:             true,
	kwPROPERTY:               true,
	kwQUANTILE_STATE:         true,
	kwQUANTILE_UNION:         true,
	kwQUARTER:                true,
	kwQUERY:                  true,
	kwQUOTA:                  true,
	kwQUALIFY:                true,
	kwQUEUED:                 true,
	kwRANDOM:                 true,
	kwRECENT:                 true,
	kwRECOVER:                true,
	kwRECYCLE:                true,
	kwREFRESH:                true,
	kwREPEATABLE:             true,
	kwREPLACE:                true,
	kwREPLACE_IF_NOT_NULL:    true,
	kwREPLAYER:               true,
	kwREPOSITORIES:           true,
	kwREPOSITORY:             true,
	kwRESOURCE:               true,
	kwRESOURCES:              true,
	kwRESTORE:                true,
	kwRESTRICTIVE:            true,
	kwRESUME:                 true,
	kwRETAIN:                 true,
	kwRETENTION:              true,
	kwRETURNS:                true,
	kwREWRITTEN:              true,
	kwRLIKE:                  true,
	kwROLLBACK:               true,
	kwROLLUP:                 true,
	kwROOT:                   true,
	kwROTATE:                 true,
	kwROUTINE:                true,
	kwS3:                     true,
	kwSAMPLE:                 true,
	kwSCHEDULE:               true,
	kwSCHEDULER:              true,
	kwSCHEMA:                 true,
	kwSECOND:                 true,
	kwSEPARATOR:              true,
	kwSERIALIZABLE:           true,
	kwSET_SESSION_VARIABLE:   true,
	kwSESSION:                true,
	kwSESSION_USER:           true,
	kwSHAPE:                  true,
	kwSKEW:                   true,
	kwSNAPSHOT:               true,
	kwSNAPSHOTS:              true,
	kwSONAME:                 true,
	kwSPLIT:                  true,
	kwSQL:                    true,
	kwSTAGE:                  true,
	kwSTAGES:                 true,
	kwSTART:                  true,
	kwSTARTS:                 true,
	kwSTATS:                  true,
	kwSTATUS:                 true,
	kwSTOP:                   true,
	kwSTORAGE:                true,
	kwSTREAM:                 true,
	kwSTREAMING:              true,
	kwSTRING:                 true,
	kwSTRUCT:                 true,
	kwSUBSTR:                 true,
	kwSUBSTRING:              true,
	kwSUM:                    true,
	kwTABLES:                 true,
	kwTAG:                    true,
	kwTASK:                   true,
	kwTASKS:                  true,
	kwTDE:                    true,
	kwTEMPORARY:              true,
	kwTEXT:                   true,
	kwTHAN:                   true,
	kwTIME:                   true,
	kwTIMESTAMP:              true,
	kwTRANSACTION:            true,
	kwTREE:                   true,
	kwTRIGGERS:               true,
	kwTRUNCATE:               true,
	kwTYPE:                   true,
	kwTYPES:                  true,
	kwUNCOMMITTED:            true,
	kwUNLOCK:                 true,
	kwUNSET:                  true,
	kwUP:                     true,
	kwUSER:                   true,
	kwVALUE:                  true,
	kwVARBINARY:              true,
	kwVARCHAR:                true,
	kwVARIABLE:               true,
	kwVARIABLES:              true,
	kwVARIANT:                true,
	kwVAULT:                  true,
	kwVAULTS:                 true,
	kwVERBOSE:                true,
	kwVERSION:                true,
	kwVIEW:                   true,
	kwVIEWS:                  true,
	kwWARM:                   true,
	kwWARNINGS:               true,
	kwWEEK:                   true,
	kwWORK:                   true,
	kwYEAR:                   true,
}

// KeywordToken returns the keyword TokenKind for a name, or (0, false)
// if name is not a keyword. Lookup is case-insensitive.
func KeywordToken(name string) (TokenKind, bool) {
	kind, ok := keywordMap[strings.ToLower(name)]
	return kind, ok
}

// IsReserved reports whether kind is a reserved keyword that cannot be
// used as an unquoted identifier.
func IsReserved(kind TokenKind) bool {
	return kind >= 700 && !nonReservedKeywords[kind]
}

// tokenNames is lazily initialized by TokenName.
var tokenNames map[TokenKind]string

// TokenName returns a human-readable name for a TokenKind.
func TokenName(kind TokenKind) string {
	if tokenNames == nil {
		tokenNames = make(map[TokenKind]string, len(keywordMap)+20)
		for name, k := range keywordMap {
			tokenNames[k] = strings.ToUpper(name)
		}
		tokenNames[tokEOF] = "EOF"
		tokenNames[tokInvalid] = "INVALID"
		tokenNames[tokInt] = "INT"
		tokenNames[tokFloat] = "FLOAT"
		tokenNames[tokString] = "STRING"
		tokenNames[tokIdent] = "IDENT"
		tokenNames[tokQuotedIdent] = "QUOTED_IDENT"
		tokenNames[tokHexLiteral] = "HEX_LITERAL"
		tokenNames[tokBitLiteral] = "BIT_LITERAL"
		tokenNames[tokPlaceholder] = "PLACEHOLDER"
		tokenNames[tokLessEq] = "<="
		tokenNames[tokGreaterEq] = ">="
		tokenNames[tokNotEq] = "<>"
		tokenNames[tokNullSafeEq] = "<=>"
		tokenNames[tokLogicalAnd] = "&&"
		tokenNames[tokDoublePipes] = "||"
		tokenNames[tokShiftLeft] = "<<"
		tokenNames[tokShiftRight] = ">>"
		tokenNames[tokArrow] = "->"
		tokenNames[tokDoubleAt] = "@@"
		tokenNames[tokAssign] = ":="
		tokenNames[tokDotDotDot] = "..."
		tokenNames[tokHintStart] = "HINT_START"
		tokenNames[tokHintEnd] = "HINT_END"
	}
	if name, ok := tokenNames[kind]; ok {
		return name
	}
	if kind > 0 && kind < 128 {
		return string(rune(kind))
	}
	return "UNKNOWN"
}
