package parser

import "strings"

// Keyword token constants. Values start at 700 and grow with iota.
// Extracted from /Users/h3n4l/OpenSource/parser/trino/TrinoLexer.g4 (295 keyword
// tokens). Numeric values are NOT stable across edits and must never be persisted;
// always compare against the named kw* constant.
//
// kwSTRING maps the literal text "STRING" (Trino's STRING type alias). Its token
// name in the legacy grammar is TEXT_STRING_, renamed here so it does not read as
// the string-literal token (which is tokString in tokens.go).
const (
	kwABSENT TokenKind = 700 + iota
	kwADD
	kwADMIN
	kwAFTER
	kwALL
	kwALTER
	kwANALYZE
	kwAND
	kwANY
	kwARRAY
	kwAS
	kwASC
	kwAT
	kwAUTHORIZATION
	kwBEGIN
	kwBERNOULLI
	kwBETWEEN
	kwBOTH
	kwBY
	kwCALL
	kwCALLED
	kwCASCADE
	kwCASE
	kwCAST
	kwCATALOG
	kwCATALOGS
	kwCOLUMN
	kwCOLUMNS
	kwCOMMENT
	kwCOMMIT
	kwCOMMITTED
	kwCONDITIONAL
	kwCONSTRAINT
	kwCOUNT
	kwCOPARTITION
	kwCREATE
	kwCROSS
	kwCUBE
	kwCURRENT
	kwCURRENT_CATALOG
	kwCURRENT_DATE
	kwCURRENT_PATH
	kwCURRENT_ROLE
	kwCURRENT_SCHEMA
	kwCURRENT_TIME
	kwCURRENT_TIMESTAMP
	kwCURRENT_USER
	kwDATA
	kwDATE
	kwDAY
	kwDEALLOCATE
	kwDECLARE
	kwDEFAULT
	kwDEFINE
	kwDEFINER
	kwDELETE
	kwDENY
	kwDESC
	kwDESCRIBE
	kwDESCRIPTOR
	kwDETERMINISTIC
	kwDISTINCT
	kwDISTRIBUTED
	kwDO
	kwDOUBLE
	kwDROP
	kwELSE
	kwEMPTY
	kwELSEIF
	kwENCODING
	kwEND
	kwERROR
	kwESCAPE
	kwEXCEPT
	kwEXCLUDING
	kwEXECUTE
	kwEXISTS
	kwEXPLAIN
	kwEXTRACT
	kwFALSE
	kwFETCH
	kwFILTER
	kwFINAL
	kwFIRST
	kwFOLLOWING
	kwFOR
	kwFORMAT
	kwFROM
	kwFULL
	kwFUNCTION
	kwFUNCTIONS
	kwGRACE
	kwGRANT
	kwGRANTED
	kwGRANTS
	kwGRAPHVIZ
	kwGROUP
	kwGROUPING
	kwGROUPS
	kwHAVING
	kwHOUR
	kwIF
	kwIGNORE
	kwIMMEDIATE
	kwIN
	kwINCLUDING
	kwINITIAL
	kwINNER
	kwINPUT
	kwINSERT
	kwINTERSECT
	kwINTERVAL
	kwINTO
	kwINVOKER
	kwIO
	kwIS
	kwISOLATION
	kwITERATE
	kwJOIN
	kwJSON
	kwJSON_ARRAY
	kwJSON_EXISTS
	kwJSON_OBJECT
	kwJSON_QUERY
	kwJSON_TABLE
	kwJSON_VALUE
	kwKEEP
	kwKEY
	kwKEYS
	kwLANGUAGE
	kwLAST
	kwLATERAL
	kwLEADING
	kwLEAVE
	kwLEFT
	kwLEVEL
	kwLIKE
	kwLIMIT
	kwLISTAGG
	kwLOCAL
	kwLOCALTIME
	kwLOCALTIMESTAMP
	kwLOGICAL
	kwLOOP
	kwMAP
	kwMATCH
	kwMATCHED
	kwMATCHES
	kwMATCH_RECOGNIZE
	kwMATERIALIZED
	kwMEASURES
	kwMERGE
	kwMINUTE
	kwMONTH
	kwNATURAL
	kwNESTED
	kwNEXT
	kwNFC
	kwNFD
	kwNFKC
	kwNFKD
	kwNO
	kwNONE
	kwNORMALIZE
	kwNOT
	kwNULL
	kwNULLIF
	kwNULLS
	kwOBJECT
	kwOF
	kwOFFSET
	kwOMIT
	kwON
	kwONE
	kwONLY
	kwOPTION
	kwOR
	kwORDER
	kwORDINALITY
	kwOUTER
	kwOUTPUT
	kwOVER
	kwOVERFLOW
	kwPARTITION
	kwPARTITIONS
	kwPASSING
	kwPAST
	kwPATH
	kwPATTERN
	kwPER
	kwPERIOD
	kwPERMUTE
	kwPLAN
	kwPOSITION
	kwPRECEDING
	kwPRECISION
	kwPREPARE
	kwPRIVILEGES
	kwPROPERTIES
	kwPRUNE
	kwQUOTES
	kwRANGE
	kwREAD
	kwRECURSIVE
	kwREFRESH
	kwRENAME
	kwREPEAT
	kwREPEATABLE
	kwREPLACE
	kwRESET
	kwRESPECT
	kwRESTRICT
	kwRETURN
	kwRETURNING
	kwRETURNS
	kwREVOKE
	kwRIGHT
	kwROLE
	kwROLES
	kwROLLBACK
	kwROLLUP
	kwROW
	kwROWS
	kwRUNNING
	kwSCALAR
	kwSCHEMA
	kwSCHEMAS
	kwSECOND
	kwSECURITY
	kwSEEK
	kwSELECT
	kwSERIALIZABLE
	kwSESSION
	kwSET
	kwSETS
	kwSHOW
	kwSKIP
	kwSOME
	kwSTART
	kwSTATS
	kwSUBSET
	kwSUBSTRING
	kwSYSTEM
	kwTABLE
	kwTABLES
	kwTABLESAMPLE
	kwTEXT
	kwSTRING
	kwTHEN
	kwTIES
	kwTIME
	kwTIMESTAMP
	kwTO
	kwTRAILING
	kwTRANSACTION
	kwTRIM
	kwTRUE
	kwTRUNCATE
	kwTRY_CAST
	kwTYPE
	kwUESCAPE
	kwUNBOUNDED
	kwUNCOMMITTED
	kwUNCONDITIONAL
	kwUNION
	kwUNIQUE
	kwUNKNOWN
	kwUNMATCHED
	kwUNNEST
	kwUNTIL
	kwUPDATE
	kwUSE
	kwUSER
	kwUSING
	kwUTF16
	kwUTF32
	kwUTF8
	kwVALIDATE
	kwVALUE
	kwVALUES
	kwVERBOSE
	kwVERSION
	kwVIEW
	kwWHEN
	kwWHERE
	kwWHILE
	kwWINDOW
	kwWITH
	kwWITHIN
	kwWITHOUT
	kwWORK
	kwWRAPPER
	kwWRITE
	kwYEAR
	kwZONE
)

// keywordMap maps a lowercase keyword string to its token kind. Trino keywords
// are case-insensitive (the lexer grammar sets caseInsensitive = true), so
// KeywordToken lowercases its input before lookup.
var keywordMap = map[string]TokenKind{
	"absent":            kwABSENT,
	"add":               kwADD,
	"admin":             kwADMIN,
	"after":             kwAFTER,
	"all":               kwALL,
	"alter":             kwALTER,
	"analyze":           kwANALYZE,
	"and":               kwAND,
	"any":               kwANY,
	"array":             kwARRAY,
	"as":                kwAS,
	"asc":               kwASC,
	"at":                kwAT,
	"authorization":     kwAUTHORIZATION,
	"begin":             kwBEGIN,
	"bernoulli":         kwBERNOULLI,
	"between":           kwBETWEEN,
	"both":              kwBOTH,
	"by":                kwBY,
	"call":              kwCALL,
	"called":            kwCALLED,
	"cascade":           kwCASCADE,
	"case":              kwCASE,
	"cast":              kwCAST,
	"catalog":           kwCATALOG,
	"catalogs":          kwCATALOGS,
	"column":            kwCOLUMN,
	"columns":           kwCOLUMNS,
	"comment":           kwCOMMENT,
	"commit":            kwCOMMIT,
	"committed":         kwCOMMITTED,
	"conditional":       kwCONDITIONAL,
	"constraint":        kwCONSTRAINT,
	"count":             kwCOUNT,
	"copartition":       kwCOPARTITION,
	"create":            kwCREATE,
	"cross":             kwCROSS,
	"cube":              kwCUBE,
	"current":           kwCURRENT,
	"current_catalog":   kwCURRENT_CATALOG,
	"current_date":      kwCURRENT_DATE,
	"current_path":      kwCURRENT_PATH,
	"current_role":      kwCURRENT_ROLE,
	"current_schema":    kwCURRENT_SCHEMA,
	"current_time":      kwCURRENT_TIME,
	"current_timestamp": kwCURRENT_TIMESTAMP,
	"current_user":      kwCURRENT_USER,
	"data":              kwDATA,
	"date":              kwDATE,
	"day":               kwDAY,
	"deallocate":        kwDEALLOCATE,
	"declare":           kwDECLARE,
	"default":           kwDEFAULT,
	"define":            kwDEFINE,
	"definer":           kwDEFINER,
	"delete":            kwDELETE,
	"deny":              kwDENY,
	"desc":              kwDESC,
	"describe":          kwDESCRIBE,
	"descriptor":        kwDESCRIPTOR,
	"deterministic":     kwDETERMINISTIC,
	"distinct":          kwDISTINCT,
	"distributed":       kwDISTRIBUTED,
	"do":                kwDO,
	"double":            kwDOUBLE,
	"drop":              kwDROP,
	"else":              kwELSE,
	"empty":             kwEMPTY,
	"elseif":            kwELSEIF,
	"encoding":          kwENCODING,
	"end":               kwEND,
	"error":             kwERROR,
	"escape":            kwESCAPE,
	"except":            kwEXCEPT,
	"excluding":         kwEXCLUDING,
	"execute":           kwEXECUTE,
	"exists":            kwEXISTS,
	"explain":           kwEXPLAIN,
	"extract":           kwEXTRACT,
	"false":             kwFALSE,
	"fetch":             kwFETCH,
	"filter":            kwFILTER,
	"final":             kwFINAL,
	"first":             kwFIRST,
	"following":         kwFOLLOWING,
	"for":               kwFOR,
	"format":            kwFORMAT,
	"from":              kwFROM,
	"full":              kwFULL,
	"function":          kwFUNCTION,
	"functions":         kwFUNCTIONS,
	"grace":             kwGRACE,
	"grant":             kwGRANT,
	"granted":           kwGRANTED,
	"grants":            kwGRANTS,
	"graphviz":          kwGRAPHVIZ,
	"group":             kwGROUP,
	"grouping":          kwGROUPING,
	"groups":            kwGROUPS,
	"having":            kwHAVING,
	"hour":              kwHOUR,
	"if":                kwIF,
	"ignore":            kwIGNORE,
	"immediate":         kwIMMEDIATE,
	"in":                kwIN,
	"including":         kwINCLUDING,
	"initial":           kwINITIAL,
	"inner":             kwINNER,
	"input":             kwINPUT,
	"insert":            kwINSERT,
	"intersect":         kwINTERSECT,
	"interval":          kwINTERVAL,
	"into":              kwINTO,
	"invoker":           kwINVOKER,
	"io":                kwIO,
	"is":                kwIS,
	"isolation":         kwISOLATION,
	"iterate":           kwITERATE,
	"join":              kwJOIN,
	"json":              kwJSON,
	"json_array":        kwJSON_ARRAY,
	"json_exists":       kwJSON_EXISTS,
	"json_object":       kwJSON_OBJECT,
	"json_query":        kwJSON_QUERY,
	"json_table":        kwJSON_TABLE,
	"json_value":        kwJSON_VALUE,
	"keep":              kwKEEP,
	"key":               kwKEY,
	"keys":              kwKEYS,
	"language":          kwLANGUAGE,
	"last":              kwLAST,
	"lateral":           kwLATERAL,
	"leading":           kwLEADING,
	"leave":             kwLEAVE,
	"left":              kwLEFT,
	"level":             kwLEVEL,
	"like":              kwLIKE,
	"limit":             kwLIMIT,
	"listagg":           kwLISTAGG,
	"local":             kwLOCAL,
	"localtime":         kwLOCALTIME,
	"localtimestamp":    kwLOCALTIMESTAMP,
	"logical":           kwLOGICAL,
	"loop":              kwLOOP,
	"map":               kwMAP,
	"match":             kwMATCH,
	"matched":           kwMATCHED,
	"matches":           kwMATCHES,
	"match_recognize":   kwMATCH_RECOGNIZE,
	"materialized":      kwMATERIALIZED,
	"measures":          kwMEASURES,
	"merge":             kwMERGE,
	"minute":            kwMINUTE,
	"month":             kwMONTH,
	"natural":           kwNATURAL,
	"nested":            kwNESTED,
	"next":              kwNEXT,
	"nfc":               kwNFC,
	"nfd":               kwNFD,
	"nfkc":              kwNFKC,
	"nfkd":              kwNFKD,
	"no":                kwNO,
	"none":              kwNONE,
	"normalize":         kwNORMALIZE,
	"not":               kwNOT,
	"null":              kwNULL,
	"nullif":            kwNULLIF,
	"nulls":             kwNULLS,
	"object":            kwOBJECT,
	"of":                kwOF,
	"offset":            kwOFFSET,
	"omit":              kwOMIT,
	"on":                kwON,
	"one":               kwONE,
	"only":              kwONLY,
	"option":            kwOPTION,
	"or":                kwOR,
	"order":             kwORDER,
	"ordinality":        kwORDINALITY,
	"outer":             kwOUTER,
	"output":            kwOUTPUT,
	"over":              kwOVER,
	"overflow":          kwOVERFLOW,
	"partition":         kwPARTITION,
	"partitions":        kwPARTITIONS,
	"passing":           kwPASSING,
	"past":              kwPAST,
	"path":              kwPATH,
	"pattern":           kwPATTERN,
	"per":               kwPER,
	"period":            kwPERIOD,
	"permute":           kwPERMUTE,
	"plan":              kwPLAN,
	"position":          kwPOSITION,
	"preceding":         kwPRECEDING,
	"precision":         kwPRECISION,
	"prepare":           kwPREPARE,
	"privileges":        kwPRIVILEGES,
	"properties":        kwPROPERTIES,
	"prune":             kwPRUNE,
	"quotes":            kwQUOTES,
	"range":             kwRANGE,
	"read":              kwREAD,
	"recursive":         kwRECURSIVE,
	"refresh":           kwREFRESH,
	"rename":            kwRENAME,
	"repeat":            kwREPEAT,
	"repeatable":        kwREPEATABLE,
	"replace":           kwREPLACE,
	"reset":             kwRESET,
	"respect":           kwRESPECT,
	"restrict":          kwRESTRICT,
	"return":            kwRETURN,
	"returning":         kwRETURNING,
	"returns":           kwRETURNS,
	"revoke":            kwREVOKE,
	"right":             kwRIGHT,
	"role":              kwROLE,
	"roles":             kwROLES,
	"rollback":          kwROLLBACK,
	"rollup":            kwROLLUP,
	"row":               kwROW,
	"rows":              kwROWS,
	"running":           kwRUNNING,
	"scalar":            kwSCALAR,
	"schema":            kwSCHEMA,
	"schemas":           kwSCHEMAS,
	"second":            kwSECOND,
	"security":          kwSECURITY,
	"seek":              kwSEEK,
	"select":            kwSELECT,
	"serializable":      kwSERIALIZABLE,
	"session":           kwSESSION,
	"set":               kwSET,
	"sets":              kwSETS,
	"show":              kwSHOW,
	"skip":              kwSKIP,
	"some":              kwSOME,
	"start":             kwSTART,
	"stats":             kwSTATS,
	"subset":            kwSUBSET,
	"substring":         kwSUBSTRING,
	"system":            kwSYSTEM,
	"table":             kwTABLE,
	"tables":            kwTABLES,
	"tablesample":       kwTABLESAMPLE,
	"text":              kwTEXT,
	"string":            kwSTRING,
	"then":              kwTHEN,
	"ties":              kwTIES,
	"time":              kwTIME,
	"timestamp":         kwTIMESTAMP,
	"to":                kwTO,
	"trailing":          kwTRAILING,
	"transaction":       kwTRANSACTION,
	"trim":              kwTRIM,
	"true":              kwTRUE,
	"truncate":          kwTRUNCATE,
	"try_cast":          kwTRY_CAST,
	"type":              kwTYPE,
	"uescape":           kwUESCAPE,
	"unbounded":         kwUNBOUNDED,
	"uncommitted":       kwUNCOMMITTED,
	"unconditional":     kwUNCONDITIONAL,
	"union":             kwUNION,
	"unique":            kwUNIQUE,
	"unknown":           kwUNKNOWN,
	"unmatched":         kwUNMATCHED,
	"unnest":            kwUNNEST,
	"until":             kwUNTIL,
	"update":            kwUPDATE,
	"use":               kwUSE,
	"user":              kwUSER,
	"using":             kwUSING,
	"utf16":             kwUTF16,
	"utf32":             kwUTF32,
	"utf8":              kwUTF8,
	"validate":          kwVALIDATE,
	"value":             kwVALUE,
	"values":            kwVALUES,
	"verbose":           kwVERBOSE,
	"version":           kwVERSION,
	"view":              kwVIEW,
	"when":              kwWHEN,
	"where":             kwWHERE,
	"while":             kwWHILE,
	"window":            kwWINDOW,
	"with":              kwWITH,
	"within":            kwWITHIN,
	"without":           kwWITHOUT,
	"work":              kwWORK,
	"wrapper":           kwWRAPPER,
	"write":             kwWRITE,
	"year":              kwYEAR,
	"zone":              kwZONE,
}

// nonReservedKeywords is the set of keyword kinds that may be used as an
// unquoted identifier (the nonReserved rule in TrinoParser.g4, 213 tokens).
// Every keyword NOT in this set is reserved and must be quoted to be used as an
// identifier. See IsReserved.
var nonReservedKeywords = map[TokenKind]bool{
	kwABSENT:          true,
	kwADD:             true,
	kwADMIN:           true,
	kwAFTER:           true,
	kwALL:             true,
	kwANALYZE:         true,
	kwANY:             true,
	kwARRAY:           true,
	kwASC:             true,
	kwAT:              true,
	kwAUTHORIZATION:   true,
	kwBEGIN:           true,
	kwBERNOULLI:       true,
	kwBOTH:            true,
	kwCALL:            true,
	kwCALLED:          true,
	kwCASCADE:         true,
	kwCATALOG:         true,
	kwCATALOGS:        true,
	kwCOLUMN:          true,
	kwCOLUMNS:         true,
	kwCOMMENT:         true,
	kwCOMMIT:          true,
	kwCOMMITTED:       true,
	kwCONDITIONAL:     true,
	kwCOUNT:           true,
	kwCOPARTITION:     true,
	kwCURRENT:         true,
	kwDATA:            true,
	kwDATE:            true,
	kwDAY:             true,
	kwDECLARE:         true,
	kwDEFAULT:         true,
	kwDEFINE:          true,
	kwDEFINER:         true,
	kwDENY:            true,
	kwDESC:            true,
	kwDESCRIPTOR:      true,
	kwDETERMINISTIC:   true,
	kwDISTRIBUTED:     true,
	kwDO:              true,
	kwDOUBLE:          true,
	kwEMPTY:           true,
	kwELSEIF:          true,
	kwENCODING:        true,
	kwERROR:           true,
	kwEXCLUDING:       true,
	kwEXPLAIN:         true,
	kwFETCH:           true,
	kwFILTER:          true,
	kwFINAL:           true,
	kwFIRST:           true,
	kwFOLLOWING:       true,
	kwFORMAT:          true,
	kwFUNCTION:        true,
	kwFUNCTIONS:       true,
	kwGRACE:           true,
	kwGRANT:           true,
	kwGRANTED:         true,
	kwGRANTS:          true,
	kwGRAPHVIZ:        true,
	kwGROUPS:          true,
	kwHOUR:            true,
	kwIF:              true,
	kwIGNORE:          true,
	kwIMMEDIATE:       true,
	kwINCLUDING:       true,
	kwINITIAL:         true,
	kwINPUT:           true,
	kwINTERVAL:        true,
	kwINVOKER:         true,
	kwIO:              true,
	kwISOLATION:       true,
	kwITERATE:         true,
	kwJSON:            true,
	kwKEEP:            true,
	kwKEY:             true,
	kwKEYS:            true,
	kwLANGUAGE:        true,
	kwLAST:            true,
	kwLATERAL:         true,
	kwLEADING:         true,
	kwLEAVE:           true,
	kwLEVEL:           true,
	kwLIMIT:           true,
	kwLOCAL:           true,
	kwLOGICAL:         true,
	kwLOOP:            true,
	kwMAP:             true,
	kwMATCH:           true,
	kwMATCHED:         true,
	kwMATCHES:         true,
	kwMATCH_RECOGNIZE: true,
	kwMATERIALIZED:    true,
	kwMEASURES:        true,
	kwMERGE:           true,
	kwMINUTE:          true,
	kwMONTH:           true,
	kwNESTED:          true,
	kwNEXT:            true,
	kwNFC:             true,
	kwNFD:             true,
	kwNFKC:            true,
	kwNFKD:            true,
	kwNO:              true,
	kwNONE:            true,
	kwNULLIF:          true,
	kwNULLS:           true,
	kwOBJECT:          true,
	kwOF:              true,
	kwOFFSET:          true,
	kwOMIT:            true,
	kwONE:             true,
	kwONLY:            true,
	kwOPTION:          true,
	kwORDINALITY:      true,
	kwOUTPUT:          true,
	kwOVER:            true,
	kwOVERFLOW:        true,
	kwPARTITION:       true,
	kwPARTITIONS:      true,
	kwPASSING:         true,
	kwPAST:            true,
	kwPATH:            true,
	kwPATTERN:         true,
	kwPER:             true,
	kwPERIOD:          true,
	kwPERMUTE:         true,
	kwPLAN:            true,
	kwPOSITION:        true,
	kwPRECEDING:       true,
	kwPRECISION:       true,
	kwPRIVILEGES:      true,
	kwPROPERTIES:      true,
	kwPRUNE:           true,
	kwQUOTES:          true,
	kwRANGE:           true,
	kwREAD:            true,
	kwREFRESH:         true,
	kwRENAME:          true,
	kwREPEAT:          true,
	kwREPEATABLE:      true,
	kwREPLACE:         true,
	kwRESET:           true,
	kwRESPECT:         true,
	kwRESTRICT:        true,
	kwRETURN:          true,
	kwRETURNING:       true,
	kwRETURNS:         true,
	kwREVOKE:          true,
	kwROLE:            true,
	kwROLES:           true,
	kwROLLBACK:        true,
	kwROW:             true,
	kwROWS:            true,
	kwRUNNING:         true,
	kwSCALAR:          true,
	kwSCHEMA:          true,
	kwSCHEMAS:         true,
	kwSECOND:          true,
	kwSECURITY:        true,
	kwSEEK:            true,
	kwSERIALIZABLE:    true,
	kwSESSION:         true,
	kwSET:             true,
	kwSETS:            true,
	kwSHOW:            true,
	kwSOME:            true,
	kwSTART:           true,
	kwSTATS:           true,
	kwSUBSET:          true,
	kwSUBSTRING:       true,
	kwSYSTEM:          true,
	kwTABLES:          true,
	kwTABLESAMPLE:     true,
	kwTEXT:            true,
	kwSTRING:          true,
	kwTIES:            true,
	kwTIME:            true,
	kwTIMESTAMP:       true,
	kwTO:              true,
	kwTRAILING:        true,
	kwTRANSACTION:     true,
	kwTRUNCATE:        true,
	kwTRY_CAST:        true,
	kwTYPE:            true,
	kwUNBOUNDED:       true,
	kwUNCOMMITTED:     true,
	kwUNCONDITIONAL:   true,
	kwUNIQUE:          true,
	kwUNKNOWN:         true,
	kwUNMATCHED:       true,
	kwUNTIL:           true,
	kwUPDATE:          true,
	kwUSE:             true,
	kwUSER:            true,
	kwUTF16:           true,
	kwUTF32:           true,
	kwUTF8:            true,
	kwVALIDATE:        true,
	kwVALUE:           true,
	kwVERBOSE:         true,
	kwVERSION:         true,
	kwVIEW:            true,
	kwWHILE:           true,
	kwWINDOW:          true,
	kwWITHIN:          true,
	kwWITHOUT:         true,
	kwWORK:            true,
	kwWRAPPER:         true,
	kwWRITE:           true,
	kwYEAR:            true,
	kwZONE:            true,
}

// KeywordToken returns the keyword TokenKind for name, or (0, false) if name is
// not a keyword. Lookup is case-insensitive.
func KeywordToken(name string) (TokenKind, bool) {
	kind, ok := keywordMap[strings.ToLower(name)]
	return kind, ok
}

// IsReserved reports whether kind is a reserved keyword — one that cannot be used
// as an unquoted identifier. It is meaningful only for keyword kinds (>= 700);
// for any non-keyword kind (literals, operators, EOF) it returns false, since
// those are never keywords. Callers that only ever pass kinds emitted by this
// lexer are safe: every kind >= 700 this lexer produces is a defined keyword.
func IsReserved(kind TokenKind) bool {
	return kind >= 700 && kind <= maxKeywordKind && !nonReservedKeywords[kind]
}

// maxKeywordKind is the largest keyword TokenKind, used by IsReserved to reject
// out-of-range values rather than treating any kind >= 700 as a reserved keyword.
const maxKeywordKind = kwZONE

// tokenNames is lazily initialized by TokenName.
var tokenNames map[TokenKind]string

// TokenName returns a human-readable name for a TokenKind, suitable for error
// messages. Keyword kinds yield their uppercase spelling; single-byte ASCII
// kinds yield the character itself.
func TokenName(kind TokenKind) string {
	if tokenNames == nil {
		tokenNames = make(map[TokenKind]string, len(keywordMap)+32)
		for name, k := range keywordMap {
			tokenNames[k] = strings.ToUpper(name)
		}
		tokenNames[tokEOF] = "EOF"
		tokenNames[tokInvalid] = "INVALID"
		tokenNames[tokInteger] = "INTEGER_VALUE"
		tokenNames[tokDecimal] = "DECIMAL_VALUE"
		tokenNames[tokDouble] = "DOUBLE_VALUE"
		tokenNames[tokString] = "STRING"
		tokenNames[tokUnicodeString] = "UNICODE_STRING"
		tokenNames[tokBinaryLiteral] = "BINARY_LITERAL"
		tokenNames[tokIdent] = "IDENTIFIER"
		tokenNames[tokQuotedIdent] = "QUOTED_IDENTIFIER"
		tokenNames[tokBackquotedIdent] = "BACKQUOTED_IDENTIFIER"
		tokenNames[tokDigitIdent] = "DIGIT_IDENTIFIER"
		tokenNames[tokQuestion] = "?"
		tokenNames[tokNotEq] = "<>"
		tokenNames[tokLessEq] = "<="
		tokenNames[tokGreaterEq] = ">="
		tokenNames[tokConcat] = "||"
		tokenNames[tokArrow] = "->"
		tokenNames[tokLeftArrow] = "<-"
		tokenNames[tokDoubleArrow] = "=>"
		tokenNames[tokLCurlyHyphen] = "{-"
		tokenNames[tokRCurlyHyphen] = "-}"
		// kwSTRING aliases the literal text STRING; pin its display name so the
		// map iteration above cannot leave it keyed on a different spelling.
		tokenNames[kwSTRING] = "STRING"
	}
	if name, ok := tokenNames[kind]; ok {
		return name
	}
	if kind > 0 && kind < 128 {
		return string(rune(kind))
	}
	return "UNKNOWN"
}
