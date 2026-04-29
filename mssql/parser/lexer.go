// Package parser implements a recursive descent SQL parser for T-SQL (SQL Server).
package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// Token type constants for the T-SQL lexer.
const (
	tokEOF = 0

	// Literal tokens (non-keyword, non-single-char)
	tokICONST      = iota + 256 // integer constant
	tokFCONST                   // floating-point constant
	tokSCONST                   // string constant
	tokNSCONST                  // N'...' nvarchar string constant
	tokIDENT                    // identifier (regular or [bracketed])
	tokVARIABLE                 // @variable
	tokSYSVARIABLE              // @@system_variable

	// Multi-character operators
	tokNOTEQUAL   // <> or !=
	tokLESSEQUAL  // <=
	tokGREATEQUAL // >=
	tokNOTLESS    // !<
	tokNOTGREATER // !>
	tokCOLONCOLON // ::
	tokPLUSEQUAL  // +=
	tokMINUSEQUAL // -=
	tokMULEQUAL   // *=
	tokDIVEQUAL   // /=
	tokMODEQUAL   // %=
	tokANDEQUAL   // &=
	tokOREQUAL    // |=
	tokXOREQUAL   // ^=

	// Keywords start here. T-SQL keywords (case-insensitive).
	kwABSENT
	kwABSOLUTE
	kwACCENT_SENSITIVITY
	kwACTION
	kwACTIVATION
	kwADD
	kwAFFINITY
	kwAFTER
	kwAGGREGATE
	kwALGORITHM
	kwALL
	kwALL_SPARSE_COLUMNS
	kwALTER
	kwALWAYS
	kwASSEMBLY
	kwAND
	kwANY
	kwAPPEND
	kwAPPLICATION
	kwAPPLY
	kwAS
	kwASC
	kwASYMMETRIC
	kwAT
	kwATTACH
	kwATTACH_REBUILD_LOG
	kwAUDIT
	kwAUTHENTICATION
	kwAUTHORIZATION
	kwAUTO
	kwAVAILABILITY
	kwBACKUP
	kwBASE64
	kwBEFORE
	kwBEGIN
	kwBETWEEN
	kwBINARY
	kwBINDING
	kwBLOCK
	kwBREAK
	kwBROKER
	kwBROWSE
	kwBUCKET_COUNT
	kwBUFFER
	kwBULK
	kwBY
	kwCACHE
	kwCALLED
	kwCALLER
	kwCASCADE
	kwCASE
	kwCAST
	kwCATALOG
	kwCATCH
	kwCERTIFICATE
	kwCHANGE_TRACKING
	kwCHECK
	kwCHECKPOINT
	kwCLASSIFICATION
	kwCLASSIFIER
	kwCLEANUP
	kwCLEAR
	kwCLONE
	kwCLOSE
	kwCLUSTER
	kwCLUSTERED
	kwCOALESCE
	kwCOLLATE
	kwCOLLECTION
	kwCOLUMN
	kwCOLUMN_SET
	kwCOLUMNSTORE
	kwCOMMIT
	kwCOMMITTED
	kwCOMPUTE
	kwCONCAT
	kwCONFIGURATION
	kwCONNECTION
	kwCONSTRAINT
	kwCONTAINMENT
	kwCONTAINS
	kwCONTAINSTABLE
	kwCONTEXT
	kwCONTINUE
	kwCONTRACT
	kwCONVERSATION
	kwCONVERT
	kwCOOKIE
	kwCOPY
	kwCOUNTER
	kwCREATE
	kwCREDENTIAL
	kwCROSS
	kwCRYPTOGRAPHIC
	kwCUBE
	kwCURRENT
	kwCURRENT_DATE
	kwCYCLE
	kwCURRENT_TIME
	kwCURRENT_TIMESTAMP
	kwCURRENT_USER
	kwCURSOR
	kwDATA
	kwDATA_SOURCE
	kwDATABASE
	kwDATABASE_SNAPSHOT
	kwDAYS
	kwDBCC
	kwDEALLOCATE
	kwDECLARE
	kwDECRYPTION
	kwDEFAULT
	kwDELAY
	kwDELAYED_DURABILITY
	kwDELETE
	kwDENY
	kwDEPENDENTS
	kwDESC
	kwDESCRIPTION
	kwDIAGNOSTICS
	kwDIALOG
	kwDIRECTORY_NAME
	kwDISABLE
	kwDISTINCT
	kwDISTRIBUTED
	kwDISTRIBUTION
	kwDO
	kwDOUBLE
	kwDROP
	kwDUMP
	kwDYNAMIC
	kwEDGE
	kwELEMENTS
	kwELSE
	kwENABLE
	kwENCRYPTED
	kwENCRYPTION
	kwEND
	kwENDPOINT
	kwERRLVL
	kwERROR
	kwESCAPE
	kwEVENT
	kwEXCEPT
	kwEXEC
	kwEXECUTE
	kwEXISTS
	kwEXIT
	kwEXPAND
	kwEXTENSION
	kwEXTERNAL
	kwFAILOVER
	kwFAN_IN
	kwFAST
	kwFAST_FORWARD
	kwFEDERATION
	kwFETCH
	kwFILE
	kwFILEGROUP
	kwFILENAME
	kwFILESTREAM
	kwFILESTREAM_ON
	kwFILETABLE
	kwFILETABLE_NAMESPACE
	kwFILLFACTOR
	kwFILTER
	kwFILTERING
	kwFIRST
	kwFOLLOWING
	kwFOR
	kwFOR_APPEND
	kwFORCE
	kwFORCE_FAILOVER_ALLOW_DATA_LOSS
	kwFOREIGN
	kwFORMAT
	kwFORWARD_ONLY
	kwFREETEXT
	kwFREETEXTTABLE
	kwFROM
	kwFULL
	kwFULLTEXT
	kwFUNCTION
	kwGB
	kwGENERATED
	kwGET
	kwGLOBAL
	kwGO
	kwGOTO
	kwGOVERNOR
	kwGRANT
	kwGROUP
	kwGROUPING
	kwGROUPS
	kwHADR
	kwHARDWARE_OFFLOAD
	kwHASH
	kwHASHED
	kwHAVING
	kwHEAP
	kwHIDDEN
	kwHIGH
	kwHINT
	kwHOLDLOCK
	kwHOURS
	kwHTTP
	kwIDENTITY
	kwIDENTITY_INSERT
	kwIDENTITYCOL
	kwIF
	kwIIF
	kwIMMEDIATE
	kwIN
	kwINCLUDE
	kwINCLUDE_NULL_VALUES
	kwINCREMENT
	kwINDEX
	kwINNER
	kwINPUT
	kwINSENSITIVE
	kwINSERT
	kwINSTEAD
	kwINTERSECT
	kwINTO
	kwIS
	kwISOLATION
	kwJOB
	kwJOIN
	kwJSON
	kwKB
	kwKEEP
	kwKEEPFIXED
	kwKEY
	kwKEYS
	kwKEYSET
	kwKILL
	kwLANGUAGE
	kwLAST
	kwLEFT
	kwLEVEL
	kwLIBRARY
	kwLIFETIME
	kwLIKE
	kwLINENO
	kwLIST
	kwLISTENER
	kwLISTENER_IP
	kwLISTENER_PORT
	kwLOAD
	kwLOB_COMPACTION
	kwLOCAL
	kwLOG
	kwLOGIN
	kwLOGON
	kwLOOP
	kwLOW
	kwMANUAL
	kwMANUAL_CUTOVER
	kwMARK
	kwMASKED
	kwMASTER
	kwMATCHED
	kwMATERIALIZED
	kwMAX
	kwMAX_QUEUE_READERS
	kwMAXVALUE
	kwMAXDOP
	kwMAXRECURSION
	kwMB
	kwMEMBER
	kwMEMORY_OPTIMIZED
	kwMEMORY_OPTIMIZED_DATA
	kwMERGE
	kwMESSAGE
	kwMESSAGE_FORWARD_SIZE
	kwMESSAGE_FORWARDING
	kwMINUTES
	kwMINVALUE
	kwMIRROR
	kwMIRRORING
	kwMODE
	kwMODEL
	kwMODIFY
	kwMOVE
	kwMUST_CHANGE
	kwNAME
	kwNATIONAL
	kwNATIVE_COMPILATION
	kwNEXT
	kwNO
	kwNOCHECK
	kwNOCOUNT
	kwNODE
	kwNOLOCK
	kwNONCLUSTERED
	kwNONE
	kwNORESET
	kwNOT
	kwNOTIFICATION
	kwNOWAIT
	kwNULL
	kwNULLIF
	kwNUMANODE
	kwOBJECT
	kwOF
	kwOFF
	kwOFFLINE
	kwOFFSET
	kwOFFSETS
	kwOLD_PASSWORD
	kwON
	kwONLY
	kwOPEN
	kwOPENDATASOURCE
	kwOPENJSON
	kwOPENQUERY
	kwOPENROWSET
	kwOPENXML
	kwOPTIMIZE
	kwOPTIMISTIC
	kwOPTION
	kwOR
	kwORDER
	kwOUT
	kwOUTER
	kwOUTPUT
	kwOVER
	kwOVERRIDE
	kwOWNER
	kwPAGE
	kwPARAMETERIZATION
	kwPARTITION
	kwPARTITIONS
	kwPASSWORD
	kwPATH
	kwPAUSE
	kwPERCENT
	kwPERIOD
	kwPERMISSION_SET
	kwPERSISTED
	kwPIVOT
	kwPLAN
	kwPLATFORM
	kwPOISON_MESSAGE_HANDLING
	kwPOLICY
	kwPOOL
	kwPOPULATION
	kwPRECEDING
	kwPRECISION
	kwPREDICATE
	kwPREDICT
	kwPRIMARY
	kwPRINT
	kwPRIOR
	kwPRIORITY
	kwPRIVILEGES
	kwPROC
	kwPROCEDURE
	kwPROCEDURE_CACHE
	kwPROCEDURE_NAME
	kwPROCESS
	kwPROPERTY
	kwPROVIDER
	kwPUBLIC
	kwQUERY
	kwQUERYTRACEON
	kwQUEUE
	kwRAISERROR
	kwRANGE
	kwRAW
	kwREAD
	kwREAD_ONLY
	kwREAD_WRITE_FILEGROUPS
	kwREADONLY
	kwREADTEXT
	kwREBUILD
	kwRECEIVE
	kwRECOMPILE
	kwRECONFIGURE
	kwREFERENCES
	kwREGENERATE
	kwRELATED_CONVERSATION
	kwRELATED_CONVERSATION_GROUP
	kwRELATIVE
	kwREMOTE
	kwREMOVE
	kwRENAME
	kwREORGANIZE
	kwREPEATABLE
	kwRESET
	kwREPLICA
	kwREPLICATION
	kwRESAMPLE
	kwRESOURCE
	kwRESOURCE_POOL
	kwRESTART
	kwRESTORE
	kwRESTRICT
	kwRESULT
	kwRESUME
	kwRETENTION
	kwRETURN
	kwRETURNS
	kwREVERT
	kwREVOKE
	kwRIGHT
	kwROBUST
	kwROLE
	kwROLLBACK
	kwROLLUP
	kwROOT
	kwROUND_ROBIN
	kwROUTE
	kwROW
	kwROWCOUNT
	kwROWGUIDCOL
	kwROWS
	kwRULE
	kwSAMPLE
	kwSAVE
	kwSCHEDULER
	kwSCHEMA
	kwSCHEMABINDING
	kwSCHEME
	kwSCROLL
	kwSCROLL_LOCKS
	kwSCOPED
	kwSEARCH
	kwSECONDARY
	kwSECONDS
	kwSECURITY
	kwSECURITYAUDIT
	kwSELECT
	kwSELECTIVE
	kwSELF
	kwSEMANTICKEYPHRASETABLE
	kwSEMANTICSIMILARITYDETAILSTABLE
	kwSEMANTICSIMILARITYTABLE
	kwSEMIJOIN
	kwSEND
	kwSENSITIVITY
	kwSEQUENCE
	kwSENT
	kwSERIALIZABLE
	kwSERVER
	kwSERVICE
	kwSESSION
	kwSESSION_USER
	kwSET
	kwSETS
	kwSETUSER
	kwSHUTDOWN
	kwSIGNATURE
	kwSIZE
	kwSNAPSHOT
	kwSOFTNUMA
	kwSOME
	kwSOURCE
	kwSPARSE
	kwSPATIAL
	kwSPECIFICATION
	kwSPLIT
	kwSTART
	kwSTATE
	kwSTATIC
	kwSTATISTICAL_SEMANTICS
	kwSTATISTICS
	kwSTATS
	kwSTATUS
	kwSTATUSONLY
	kwSTOP
	kwSTOPLIST
	kwSTREAM
	kwSTREAMING
	kwSUBSCRIPTION
	kwSUSPEND_FOR_SNAPSHOT_BACKUP
	kwSWITCH
	kwSYMMETRIC
	kwSYNONYM
	kwSYSTEM
	kwSYSTEM_TIME
	kwSYSTEM_USER
	kwTABLE
	kwTABLESAMPLE
	kwTARGET
	kwTB
	kwTCP
	kwTEMPDB_METADATA
	kwTEXTIMAGE_ON
	kwTEXTSIZE
	kwTHEN
	kwTHROW
	kwTIES
	kwTIME
	kwTIMEOUT
	kwTIMER
	kwTO
	kwTOP
	kwTRAN
	kwTRANSACTION
	kwTRANSFER
	kwTRIGGER
	kwTRUNCATE
	kwTRY
	kwTRY_CAST
	kwTRY_CONVERT
	kwTSEQUAL
	kwTYPE
	kwTYPE_WARNING
	kwUNBOUNDED
	kwUNCOMMITTED
	kwUNDEFINED
	kwUNION
	kwUNIQUE
	kwUNKNOWN
	kwUNLIMITED
	kwUNLOCK
	kwUNPIVOT
	kwUPDATE
	kwUPDATETEXT
	kwURL
	kwUSE
	kwUSED
	kwUSER
	kwUSING
	kwVALIDATION
	kwVALUE
	kwVALUES
	kwVARYING
	kwVECTOR
	kwVIEW
	kwVIEWS
	kwWAITFOR
	kwWHEN
	kwWHERE
	kwWHILE
	kwWINDOW
	kwWINDOWS
	kwWITH
	kwWITHIN
	kwWITHOUT
	kwWITHOUT_ARRAY_WRAPPER
	kwWORK
	kwWORKLOAD
	kwWRITE
	kwWRITETEXT
	kwXACT_ABORT
	kwXML
	kwXMLDATA
	kwXMLNAMESPACES
	kwXMLSCHEMA
	kwXSINIL
	kwZONE
)

// keywordMap maps lowercase keyword strings to token types.
var keywordMap map[string]int

func init() {
	keywordMap = map[string]int{
		"accent_sensitivity": kwACCENT_SENSITIVITY, "action": kwACTION, "activation": kwACTIVATION, "add": kwADD,
		"affinity": kwAFFINITY, "after": kwAFTER, "aggregate": kwAGGREGATE, "algorithm": kwALGORITHM,
		"all": kwALL, "all_sparse_columns": kwALL_SPARSE_COLUMNS, "alter": kwALTER, "always": kwALWAYS,
		"and": kwAND, "any": kwANY, "append": kwAPPEND, "application": kwAPPLICATION,
		"apply": kwAPPLY, "as": kwAS, "asc": kwASC, "asymmetric": kwASYMMETRIC,
		"at": kwAT, "attach": kwATTACH, "attach_rebuild_log": kwATTACH_REBUILD_LOG, "audit": kwAUDIT,
		"authentication": kwAUTHENTICATION, "authorization": kwAUTHORIZATION, "auto": kwAUTO, "availability": kwAVAILABILITY,
		"backup": kwBACKUP, "before": kwBEFORE, "begin": kwBEGIN, "between": kwBETWEEN,
		"binding": kwBINDING, "block": kwBLOCK, "break": kwBREAK, "broker": kwBROKER,
		"browse": kwBROWSE, "buffer": kwBUFFER, "bulk": kwBULK, "by": kwBY,
		"called": kwCALLED, "caller": kwCALLER, "cascade": kwCASCADE, "case": kwCASE,
		"cast": kwCAST, "catalog": kwCATALOG, "catch": kwCATCH, "certificate": kwCERTIFICATE,
		"change_tracking": kwCHANGE_TRACKING, "check": kwCHECK, "checkpoint": kwCHECKPOINT, "classification": kwCLASSIFICATION,
		"classifier": kwCLASSIFIER, "cleanup": kwCLEANUP, "clear": kwCLEAR, "clone": kwCLONE,
		"close": kwCLOSE, "cluster": kwCLUSTER, "clustered": kwCLUSTERED, "coalesce": kwCOALESCE,
		"collate": kwCOLLATE, "collection": kwCOLLECTION, "column": kwCOLUMN, "column_set": kwCOLUMN_SET,
		"columnstore": kwCOLUMNSTORE, "commit": kwCOMMIT, "committed": kwCOMMITTED, "compute": kwCOMPUTE,
		"concat": kwCONCAT, "configuration": kwCONFIGURATION, "connection": kwCONNECTION, "constraint": kwCONSTRAINT,
		"containment": kwCONTAINMENT, "contains": kwCONTAINS, "containstable": kwCONTAINSTABLE, "context": kwCONTEXT,
		"continue": kwCONTINUE, "contract": kwCONTRACT, "conversation": kwCONVERSATION, "convert": kwCONVERT,
		"cookie": kwCOOKIE, "copy": kwCOPY, "counter": kwCOUNTER, "create": kwCREATE,
		"credential": kwCREDENTIAL, "cross": kwCROSS, "cryptographic": kwCRYPTOGRAPHIC, "cube": kwCUBE,
		"current": kwCURRENT, "current_date": kwCURRENT_DATE, "current_time": kwCURRENT_TIME, "current_timestamp": kwCURRENT_TIMESTAMP,
		"current_user": kwCURRENT_USER, "cursor": kwCURSOR,
		"data": kwDATA, "data_source": kwDATA_SOURCE, "database": kwDATABASE, "database_snapshot": kwDATABASE_SNAPSHOT,
		"days": kwDAYS, "dbcc": kwDBCC, "deallocate": kwDEALLOCATE, "declare": kwDECLARE,
		"decryption": kwDECRYPTION, "default": kwDEFAULT, "delay": kwDELAY, "delete": kwDELETE,
		"deny": kwDENY, "dependents": kwDEPENDENTS, "desc": kwDESC, "description": kwDESCRIPTION,
		"diagnostics": kwDIAGNOSTICS, "dialog": kwDIALOG, "directory_name": kwDIRECTORY_NAME, "disable": kwDISABLE,
		"distinct": kwDISTINCT, "distributed": kwDISTRIBUTED, "distribution": kwDISTRIBUTION, "do": kwDO,
		"double": kwDOUBLE, "drop": kwDROP, "dump": kwDUMP,
		"edge": kwEDGE, "else": kwELSE, "enable": kwENABLE, "encrypted": kwENCRYPTED,
		"encryption": kwENCRYPTION, "end": kwEND, "endpoint": kwENDPOINT, "errlvl": kwERRLVL,
		"error": kwERROR, "escape": kwESCAPE, "event": kwEVENT, "except": kwEXCEPT,
		"exec": kwEXEC, "execute": kwEXECUTE, "exists": kwEXISTS, "exit": kwEXIT,
		"expand": kwEXPAND, "extension": kwEXTENSION, "external": kwEXTERNAL,
		"failover": kwFAILOVER, "fan_in": kwFAN_IN, "fast": kwFAST, "federation": kwFEDERATION,
		"fetch": kwFETCH, "file": kwFILE, "filegroup": kwFILEGROUP, "filename": kwFILENAME,
		"filestream": kwFILESTREAM, "filestream_on": kwFILESTREAM_ON, "filetable": kwFILETABLE, "filetable_namespace": kwFILETABLE_NAMESPACE,
		"fillfactor": kwFILLFACTOR, "filter": kwFILTER, "first": kwFIRST, "following": kwFOLLOWING,
		"for": kwFOR, "for_append": kwFOR_APPEND, "force": kwFORCE, "force_failover_allow_data_loss": kwFORCE_FAILOVER_ALLOW_DATA_LOSS,
		"foreign": kwFOREIGN, "format": kwFORMAT, "freetext": kwFREETEXT, "freetexttable": kwFREETEXTTABLE,
		"from": kwFROM, "full": kwFULL, "fulltext": kwFULLTEXT, "function": kwFUNCTION,
		"gb": kwGB, "generated": kwGENERATED, "get": kwGET, "go": kwGO,
		"goto": kwGOTO, "governor": kwGOVERNOR, "grant": kwGRANT, "group": kwGROUP,
		"grouping": kwGROUPING, "groups": kwGROUPS,
		"hadr": kwHADR, "hardware_offload": kwHARDWARE_OFFLOAD, "hash": kwHASH, "hashed": kwHASHED,
		"having": kwHAVING, "heap": kwHEAP, "hidden": kwHIDDEN, "high": kwHIGH,
		"hint": kwHINT, "holdlock": kwHOLDLOCK, "hours": kwHOURS, "http": kwHTTP,
		"identity": kwIDENTITY, "identity_insert": kwIDENTITY_INSERT, "identitycol": kwIDENTITYCOL, "if": kwIF,
		"iif": kwIIF, "immediate": kwIMMEDIATE, "in": kwIN, "include": kwINCLUDE,
		"index": kwINDEX, "inner": kwINNER, "input": kwINPUT, "insert": kwINSERT,
		"instead": kwINSTEAD, "intersect": kwINTERSECT, "into": kwINTO, "is": kwIS,
		"isolation": kwISOLATION,
		"job": kwJOB, "join": kwJOIN, "json": kwJSON,
		"kb": kwKB, "keep": kwKEEP, "keepfixed": kwKEEPFIXED, "key": kwKEY, "keys": kwKEYS,
		"kill": kwKILL,
		"language": kwLANGUAGE, "left": kwLEFT, "level": kwLEVEL, "library": kwLIBRARY,
		"lifetime": kwLIFETIME, "like": kwLIKE, "lineno": kwLINENO, "list": kwLIST,
		"listener": kwLISTENER, "listener_ip": kwLISTENER_IP, "listener_port": kwLISTENER_PORT, "load": kwLOAD,
		"lob_compaction": kwLOB_COMPACTION, "local": kwLOCAL, "log": kwLOG, "login": kwLOGIN,
		"loop": kwLOOP, "low": kwLOW,
		"manual": kwMANUAL, "manual_cutover": kwMANUAL_CUTOVER, "masked": kwMASKED, "master": kwMASTER,
		"matched": kwMATCHED, "materialized": kwMATERIALIZED, "max": kwMAX, "max_queue_readers": kwMAX_QUEUE_READERS,
		"maxdop": kwMAXDOP, "maxrecursion": kwMAXRECURSION, "mb": kwMB, "member": kwMEMBER,
		"memory_optimized": kwMEMORY_OPTIMIZED, "memory_optimized_data": kwMEMORY_OPTIMIZED_DATA, "merge": kwMERGE, "message": kwMESSAGE,
		"message_forward_size": kwMESSAGE_FORWARD_SIZE, "message_forwarding": kwMESSAGE_FORWARDING, "minutes": kwMINUTES, "mirror": kwMIRROR,
		"mirroring": kwMIRRORING, "mode": kwMODE, "model": kwMODEL, "modify": kwMODIFY,
		"move": kwMOVE, "must_change": kwMUST_CHANGE,
		"name": kwNAME, "national": kwNATIONAL, "native_compilation": kwNATIVE_COMPILATION, "next": kwNEXT,
		"no": kwNO, "nocheck": kwNOCHECK, "nocount": kwNOCOUNT, "node": kwNODE,
		"nolock": kwNOLOCK, "nonclustered": kwNONCLUSTERED, "none": kwNONE, "not": kwNOT,
		"notification": kwNOTIFICATION, "nowait": kwNOWAIT, "null": kwNULL, "nullif": kwNULLIF,
		"numanode": kwNUMANODE,
		"object": kwOBJECT, "of": kwOF, "off": kwOFF, "offline": kwOFFLINE,
		"offset": kwOFFSET, "offsets": kwOFFSETS, "old_password": kwOLD_PASSWORD, "on": kwON,
		"only": kwONLY, "open": kwOPEN, "opendatasource": kwOPENDATASOURCE, "openjson": kwOPENJSON,
		"openquery": kwOPENQUERY, "openrowset": kwOPENROWSET, "openxml": kwOPENXML, "optimize": kwOPTIMIZE,
		"option": kwOPTION, "or": kwOR, "order": kwORDER, "out": kwOUT,
		"outer": kwOUTER, "output": kwOUTPUT, "over": kwOVER, "override": kwOVERRIDE,
		"owner": kwOWNER,
		"page": kwPAGE, "parameterization": kwPARAMETERIZATION, "partition": kwPARTITION, "partitions": kwPARTITIONS,
		"password": kwPASSWORD, "path": kwPATH, "pause": kwPAUSE, "percent": kwPERCENT,
		"period": kwPERIOD, "permission_set": kwPERMISSION_SET, "persisted": kwPERSISTED, "pivot": kwPIVOT,
		"plan": kwPLAN, "platform": kwPLATFORM, "poison_message_handling": kwPOISON_MESSAGE_HANDLING, "policy": kwPOLICY,
		"pool": kwPOOL, "population": kwPOPULATION, "preceding": kwPRECEDING, "precision": kwPRECISION,
		"predicate": kwPREDICATE, "predict": kwPREDICT, "primary": kwPRIMARY, "print": kwPRINT,
		"priority": kwPRIORITY, "privileges": kwPRIVILEGES, "proc": kwPROC, "procedure": kwPROCEDURE,
		"procedure_cache": kwPROCEDURE_CACHE, "procedure_name": kwPROCEDURE_NAME, "process": kwPROCESS, "property": kwPROPERTY,
		"provider": kwPROVIDER, "public": kwPUBLIC,
		"query": kwQUERY, "querytraceon": kwQUERYTRACEON, "queue": kwQUEUE,
		"raiserror": kwRAISERROR, "range": kwRANGE, "raw": kwRAW, "read": kwREAD,
		"read_write_filegroups": kwREAD_WRITE_FILEGROUPS, "readonly": kwREADONLY, "readtext": kwREADTEXT, "rebuild": kwREBUILD,
		"receive": kwRECEIVE, "recompile": kwRECOMPILE, "reconfigure": kwRECONFIGURE, "references": kwREFERENCES,
		"regenerate": kwREGENERATE, "related_conversation": kwRELATED_CONVERSATION, "related_conversation_group": kwRELATED_CONVERSATION_GROUP, "remote": kwREMOTE,
		"remove": kwREMOVE, "rename": kwRENAME, "reorganize": kwREORGANIZE, "repeatable": kwREPEATABLE,
		"replica": kwREPLICA, "replication": kwREPLICATION, "resample": kwRESAMPLE, "resource": kwRESOURCE,
		"resource_pool": kwRESOURCE_POOL, "restart": kwRESTART, "restore": kwRESTORE, "restrict": kwRESTRICT,
		"result": kwRESULT, "resume": kwRESUME, "retention": kwRETENTION, "return": kwRETURN,
		"returns": kwRETURNS, "revert": kwREVERT, "revoke": kwREVOKE, "right": kwRIGHT,
		"robust": kwROBUST, "role": kwROLE, "rollback": kwROLLBACK, "rollup": kwROLLUP,
		"root": kwROOT, "round_robin": kwROUND_ROBIN, "route": kwROUTE, "row": kwROW,
		"rowcount": kwROWCOUNT, "rowguidcol": kwROWGUIDCOL, "rows": kwROWS, "rule": kwRULE,
		"sample": kwSAMPLE, "save": kwSAVE, "scheduler": kwSCHEDULER, "schema": kwSCHEMA,
		"schemabinding": kwSCHEMABINDING, "scheme": kwSCHEME, "scoped": kwSCOPED, "search": kwSEARCH,
		"secondary": kwSECONDARY, "seconds": kwSECONDS, "security": kwSECURITY, "securityaudit": kwSECURITYAUDIT,
		"select": kwSELECT, "selective": kwSELECTIVE, "self": kwSELF, "semantickeyphrasetable": kwSEMANTICKEYPHRASETABLE,
		"semanticsimilaritydetailstable": kwSEMANTICSIMILARITYDETAILSTABLE, "semanticsimilaritytable": kwSEMANTICSIMILARITYTABLE, "semijoin": kwSEMIJOIN, "send": kwSEND,
		"sensitivity": kwSENSITIVITY, "sent": kwSENT, "serializable": kwSERIALIZABLE, "server": kwSERVER,
		"service": kwSERVICE, "session": kwSESSION, "session_user": kwSESSION_USER, "set": kwSET,
		"sets": kwSETS, "setuser": kwSETUSER, "shutdown": kwSHUTDOWN, "signature": kwSIGNATURE,
		"size": kwSIZE, "snapshot": kwSNAPSHOT, "softnuma": kwSOFTNUMA, "some": kwSOME,
		"source": kwSOURCE, "sparse": kwSPARSE, "spatial": kwSPATIAL, "specification": kwSPECIFICATION,
		"split": kwSPLIT, "start": kwSTART, "state": kwSTATE, "static": kwSTATIC, "statistical_semantics": kwSTATISTICAL_SEMANTICS,
		"statistics": kwSTATISTICS, "stats": kwSTATS, "status": kwSTATUS, "statusonly": kwSTATUSONLY,
		"stop": kwSTOP, "stoplist": kwSTOPLIST, "stream": kwSTREAM, "streaming": kwSTREAMING,
		"subscription": kwSUBSCRIPTION, "suspend_for_snapshot_backup": kwSUSPEND_FOR_SNAPSHOT_BACKUP, "switch": kwSWITCH, "symmetric": kwSYMMETRIC, "synonym": kwSYNONYM,
		"system": kwSYSTEM, "system_time": kwSYSTEM_TIME, "system_user": kwSYSTEM_USER,
		"table": kwTABLE, "tablesample": kwTABLESAMPLE, "target": kwTARGET, "tb": kwTB,
		"tcp": kwTCP, "tempdb_metadata": kwTEMPDB_METADATA, "textimage_on": kwTEXTIMAGE_ON, "textsize": kwTEXTSIZE,
		"then": kwTHEN, "throw": kwTHROW, "ties": kwTIES, "time": kwTIME,
		"timeout": kwTIMEOUT, "timer": kwTIMER, "to": kwTO, "top": kwTOP,
		"tran": kwTRAN, "transaction": kwTRANSACTION, "trigger": kwTRIGGER, "truncate": kwTRUNCATE,
		"try": kwTRY, "try_cast": kwTRY_CAST, "try_convert": kwTRY_CONVERT, "tsequal": kwTSEQUAL,
		"type": kwTYPE, "type_warning": kwTYPE_WARNING,
		"unbounded": kwUNBOUNDED, "uncommitted": kwUNCOMMITTED, "undefined": kwUNDEFINED, "union": kwUNION,
		"unique": kwUNIQUE, "unknown": kwUNKNOWN, "unlimited": kwUNLIMITED, "unlock": kwUNLOCK,
		"unpivot": kwUNPIVOT, "update": kwUPDATE, "updatetext": kwUPDATETEXT, "url": kwURL,
		"use": kwUSE, "used": kwUSED, "user": kwUSER, "using": kwUSING,
		"validation": kwVALIDATION, "value": kwVALUE, "values": kwVALUES, "varying": kwVARYING,
		"vector": kwVECTOR, "view": kwVIEW, "views": kwVIEWS,
		"waitfor": kwWAITFOR, "when": kwWHEN, "where": kwWHERE, "while": kwWHILE,
		"window": kwWINDOW, "windows": kwWINDOWS, "with": kwWITH, "within": kwWITHIN,
		"without": kwWITHOUT, "without_array_wrapper": kwWITHOUT_ARRAY_WRAPPER, "work": kwWORK, "workload": kwWORKLOAD, "write": kwWRITE, "writetext": kwWRITETEXT,
		"xact_abort": kwXACT_ABORT, "xml": kwXML, "xmldata": kwXMLDATA, "xmlnamespaces": kwXMLNAMESPACES, "xmlschema": kwXMLSCHEMA, "xsinil": kwXSINIL,
		"zone": kwZONE,
		// Newly registered context keywords (cursor, XML, transaction, federation, etc.)
		"absent": kwABSENT, "absolute": kwABSOLUTE, "assembly": kwASSEMBLY,
		"base64": kwBASE64, "binary": kwBINARY, "bucket_count": kwBUCKET_COUNT,
		"delayed_durability": kwDELAYED_DURABILITY, "dynamic": kwDYNAMIC,
		"elements": kwELEMENTS,
		"fast_forward": kwFAST_FORWARD, "filtering": kwFILTERING, "forward_only": kwFORWARD_ONLY,
		"global": kwGLOBAL,
		"include_null_values": kwINCLUDE_NULL_VALUES, "insensitive": kwINSENSITIVE,
		"keyset": kwKEYSET,
		"mark": kwMARK,
		"noreset": kwNORESET,
		"optimistic": kwOPTIMISTIC,
		"read_only": kwREAD_ONLY, "relative": kwRELATIVE, "reset": kwRESET,
		"cache": kwCACHE, "cycle": kwCYCLE, "increment": kwINCREMENT,
		"last": kwLAST, "logon": kwLOGON, "maxvalue": kwMAXVALUE, "minvalue": kwMINVALUE,
		"prior": kwPRIOR, "transfer": kwTRANSFER,
		"scroll": kwSCROLL, "scroll_locks": kwSCROLL_LOCKS, "sequence": kwSEQUENCE,
	}
}

// lookupKeyword returns the token type for a keyword (case-insensitive),
// or tokIDENT if not a keyword.
// KeywordCategory classifies keywords as Core (reserved) or Context (context-sensitive).
type KeywordCategory int

const (
	CoreKeyword    KeywordCategory = iota // SqlScriptDOM registered - cannot be unquoted identifier
	ContextKeyword                        // Context-sensitive - can be unquoted identifier
)

// Keyword holds classification metadata for a registered keyword.
type Keyword struct {
	Name     string
	Token    int
	Category KeywordCategory
}

// keywordClassification maps keyword token types to their classification.
var keywordClassification = map[int]Keyword{
		kwABSENT: {Name: "ABSENT", Token: kwABSENT, Category: ContextKeyword},
		kwABSOLUTE: {Name: "ABSOLUTE", Token: kwABSOLUTE, Category: ContextKeyword},
		kwACCENT_SENSITIVITY: {Name: "ACCENT_SENSITIVITY", Token: kwACCENT_SENSITIVITY, Category: ContextKeyword},
		kwACTION: {Name: "ACTION", Token: kwACTION, Category: ContextKeyword},
		kwACTIVATION: {Name: "ACTIVATION", Token: kwACTIVATION, Category: ContextKeyword},
		kwADD: {Name: "ADD", Token: kwADD, Category: CoreKeyword},
		kwAFFINITY: {Name: "AFFINITY", Token: kwAFFINITY, Category: ContextKeyword},
		kwAFTER: {Name: "AFTER", Token: kwAFTER, Category: ContextKeyword},
		kwAGGREGATE: {Name: "AGGREGATE", Token: kwAGGREGATE, Category: ContextKeyword},
		kwALGORITHM: {Name: "ALGORITHM", Token: kwALGORITHM, Category: ContextKeyword},
		kwALL: {Name: "ALL", Token: kwALL, Category: CoreKeyword},
		kwALL_SPARSE_COLUMNS: {Name: "ALL_SPARSE_COLUMNS", Token: kwALL_SPARSE_COLUMNS, Category: ContextKeyword},
		kwALTER: {Name: "ALTER", Token: kwALTER, Category: CoreKeyword},
		kwALWAYS: {Name: "ALWAYS", Token: kwALWAYS, Category: ContextKeyword},
		kwASSEMBLY: {Name: "ASSEMBLY", Token: kwASSEMBLY, Category: ContextKeyword},
		kwAND: {Name: "AND", Token: kwAND, Category: CoreKeyword},
		kwANY: {Name: "ANY", Token: kwANY, Category: CoreKeyword},
		kwAPPEND: {Name: "APPEND", Token: kwAPPEND, Category: ContextKeyword},
		kwAPPLICATION: {Name: "APPLICATION", Token: kwAPPLICATION, Category: ContextKeyword},
		kwAPPLY: {Name: "APPLY", Token: kwAPPLY, Category: ContextKeyword},
		kwAS: {Name: "AS", Token: kwAS, Category: CoreKeyword},
		kwASC: {Name: "ASC", Token: kwASC, Category: CoreKeyword},
		kwASYMMETRIC: {Name: "ASYMMETRIC", Token: kwASYMMETRIC, Category: ContextKeyword},
		kwAT: {Name: "AT", Token: kwAT, Category: ContextKeyword},
		kwATTACH: {Name: "ATTACH", Token: kwATTACH, Category: ContextKeyword},
		kwATTACH_REBUILD_LOG: {Name: "ATTACH_REBUILD_LOG", Token: kwATTACH_REBUILD_LOG, Category: ContextKeyword},
		kwAUDIT: {Name: "AUDIT", Token: kwAUDIT, Category: ContextKeyword},
		kwAUTHENTICATION: {Name: "AUTHENTICATION", Token: kwAUTHENTICATION, Category: ContextKeyword},
		kwAUTHORIZATION: {Name: "AUTHORIZATION", Token: kwAUTHORIZATION, Category: CoreKeyword},
		kwAUTO: {Name: "AUTO", Token: kwAUTO, Category: ContextKeyword},
		kwAVAILABILITY: {Name: "AVAILABILITY", Token: kwAVAILABILITY, Category: ContextKeyword},
		kwBACKUP: {Name: "BACKUP", Token: kwBACKUP, Category: CoreKeyword},
		kwBASE64: {Name: "BASE64", Token: kwBASE64, Category: ContextKeyword},
		kwBEFORE: {Name: "BEFORE", Token: kwBEFORE, Category: ContextKeyword},
		kwBEGIN: {Name: "BEGIN", Token: kwBEGIN, Category: CoreKeyword},
		kwBETWEEN: {Name: "BETWEEN", Token: kwBETWEEN, Category: CoreKeyword},
		kwBINARY: {Name: "BINARY", Token: kwBINARY, Category: ContextKeyword},
		kwBINDING: {Name: "BINDING", Token: kwBINDING, Category: ContextKeyword},
		kwBLOCK: {Name: "BLOCK", Token: kwBLOCK, Category: ContextKeyword},
		kwBREAK: {Name: "BREAK", Token: kwBREAK, Category: CoreKeyword},
		kwBROKER: {Name: "BROKER", Token: kwBROKER, Category: ContextKeyword},
		kwBROWSE: {Name: "BROWSE", Token: kwBROWSE, Category: CoreKeyword},
		kwBUCKET_COUNT: {Name: "BUCKET_COUNT", Token: kwBUCKET_COUNT, Category: ContextKeyword},
		kwCACHE: {Name: "CACHE", Token: kwCACHE, Category: ContextKeyword},
		kwBUFFER: {Name: "BUFFER", Token: kwBUFFER, Category: ContextKeyword},
		kwBULK: {Name: "BULK", Token: kwBULK, Category: CoreKeyword},
		kwBY: {Name: "BY", Token: kwBY, Category: CoreKeyword},
		kwCALLED: {Name: "CALLED", Token: kwCALLED, Category: ContextKeyword},
		kwCALLER: {Name: "CALLER", Token: kwCALLER, Category: ContextKeyword},
		kwCASCADE: {Name: "CASCADE", Token: kwCASCADE, Category: CoreKeyword},
		kwCASE: {Name: "CASE", Token: kwCASE, Category: CoreKeyword},
		kwCAST: {Name: "CAST", Token: kwCAST, Category: ContextKeyword},
		kwCATALOG: {Name: "CATALOG", Token: kwCATALOG, Category: ContextKeyword},
		kwCATCH: {Name: "CATCH", Token: kwCATCH, Category: ContextKeyword},
		kwCERTIFICATE: {Name: "CERTIFICATE", Token: kwCERTIFICATE, Category: ContextKeyword},
		kwCHANGE_TRACKING: {Name: "CHANGE_TRACKING", Token: kwCHANGE_TRACKING, Category: ContextKeyword},
		kwCHECK: {Name: "CHECK", Token: kwCHECK, Category: CoreKeyword},
		kwCHECKPOINT: {Name: "CHECKPOINT", Token: kwCHECKPOINT, Category: CoreKeyword},
		kwCLASSIFICATION: {Name: "CLASSIFICATION", Token: kwCLASSIFICATION, Category: ContextKeyword},
		kwCLASSIFIER: {Name: "CLASSIFIER", Token: kwCLASSIFIER, Category: ContextKeyword},
		kwCLEANUP: {Name: "CLEANUP", Token: kwCLEANUP, Category: ContextKeyword},
		kwCLEAR: {Name: "CLEAR", Token: kwCLEAR, Category: ContextKeyword},
		kwCLONE: {Name: "CLONE", Token: kwCLONE, Category: ContextKeyword},
		kwCLOSE: {Name: "CLOSE", Token: kwCLOSE, Category: CoreKeyword},
		kwCLUSTER: {Name: "CLUSTER", Token: kwCLUSTER, Category: ContextKeyword},
		kwCLUSTERED: {Name: "CLUSTERED", Token: kwCLUSTERED, Category: CoreKeyword},
		kwCOALESCE: {Name: "COALESCE", Token: kwCOALESCE, Category: CoreKeyword},
		kwCOLLATE: {Name: "COLLATE", Token: kwCOLLATE, Category: CoreKeyword},
		kwCOLLECTION: {Name: "COLLECTION", Token: kwCOLLECTION, Category: ContextKeyword},
		kwCOLUMN: {Name: "COLUMN", Token: kwCOLUMN, Category: CoreKeyword},
		kwCOLUMN_SET: {Name: "COLUMN_SET", Token: kwCOLUMN_SET, Category: ContextKeyword},
		kwCOLUMNSTORE: {Name: "COLUMNSTORE", Token: kwCOLUMNSTORE, Category: ContextKeyword},
		kwCOMMIT: {Name: "COMMIT", Token: kwCOMMIT, Category: CoreKeyword},
		kwCOMMITTED: {Name: "COMMITTED", Token: kwCOMMITTED, Category: ContextKeyword},
		kwCOMPUTE: {Name: "COMPUTE", Token: kwCOMPUTE, Category: CoreKeyword},
		kwCONCAT: {Name: "CONCAT", Token: kwCONCAT, Category: ContextKeyword},
		kwCONFIGURATION: {Name: "CONFIGURATION", Token: kwCONFIGURATION, Category: ContextKeyword},
		kwCONNECTION: {Name: "CONNECTION", Token: kwCONNECTION, Category: ContextKeyword},
		kwCONSTRAINT: {Name: "CONSTRAINT", Token: kwCONSTRAINT, Category: CoreKeyword},
		kwCONTAINMENT: {Name: "CONTAINMENT", Token: kwCONTAINMENT, Category: ContextKeyword},
		kwCONTAINS: {Name: "CONTAINS", Token: kwCONTAINS, Category: CoreKeyword},
		kwCONTAINSTABLE: {Name: "CONTAINSTABLE", Token: kwCONTAINSTABLE, Category: CoreKeyword},
		kwCONTEXT: {Name: "CONTEXT", Token: kwCONTEXT, Category: ContextKeyword},
		kwCONTINUE: {Name: "CONTINUE", Token: kwCONTINUE, Category: CoreKeyword},
		kwCONTRACT: {Name: "CONTRACT", Token: kwCONTRACT, Category: ContextKeyword},
		kwCONVERSATION: {Name: "CONVERSATION", Token: kwCONVERSATION, Category: ContextKeyword},
		kwCONVERT: {Name: "CONVERT", Token: kwCONVERT, Category: CoreKeyword},
		kwCOOKIE: {Name: "COOKIE", Token: kwCOOKIE, Category: ContextKeyword},
		kwCOPY: {Name: "COPY", Token: kwCOPY, Category: ContextKeyword},
		kwCOUNTER: {Name: "COUNTER", Token: kwCOUNTER, Category: ContextKeyword},
		kwCREATE: {Name: "CREATE", Token: kwCREATE, Category: CoreKeyword},
		kwCREDENTIAL: {Name: "CREDENTIAL", Token: kwCREDENTIAL, Category: ContextKeyword},
		kwCROSS: {Name: "CROSS", Token: kwCROSS, Category: CoreKeyword},
		kwCRYPTOGRAPHIC: {Name: "CRYPTOGRAPHIC", Token: kwCRYPTOGRAPHIC, Category: ContextKeyword},
		kwCUBE: {Name: "CUBE", Token: kwCUBE, Category: ContextKeyword},
		kwCURRENT: {Name: "CURRENT", Token: kwCURRENT, Category: CoreKeyword},
		kwCYCLE: {Name: "CYCLE", Token: kwCYCLE, Category: ContextKeyword},
		kwCURRENT_DATE: {Name: "CURRENT_DATE", Token: kwCURRENT_DATE, Category: CoreKeyword},
		kwCURRENT_TIME: {Name: "CURRENT_TIME", Token: kwCURRENT_TIME, Category: CoreKeyword},
		kwCURRENT_TIMESTAMP: {Name: "CURRENT_TIMESTAMP", Token: kwCURRENT_TIMESTAMP, Category: CoreKeyword},
		kwCURRENT_USER: {Name: "CURRENT_USER", Token: kwCURRENT_USER, Category: CoreKeyword},
		kwCURSOR: {Name: "CURSOR", Token: kwCURSOR, Category: CoreKeyword},
		kwDATA: {Name: "DATA", Token: kwDATA, Category: ContextKeyword},
		kwDATA_SOURCE: {Name: "DATA_SOURCE", Token: kwDATA_SOURCE, Category: ContextKeyword},
		kwDATABASE: {Name: "DATABASE", Token: kwDATABASE, Category: CoreKeyword},
		kwDATABASE_SNAPSHOT: {Name: "DATABASE_SNAPSHOT", Token: kwDATABASE_SNAPSHOT, Category: ContextKeyword},
		kwDAYS: {Name: "DAYS", Token: kwDAYS, Category: ContextKeyword},
		kwDBCC: {Name: "DBCC", Token: kwDBCC, Category: CoreKeyword},
		kwDEALLOCATE: {Name: "DEALLOCATE", Token: kwDEALLOCATE, Category: CoreKeyword},
		kwDECLARE: {Name: "DECLARE", Token: kwDECLARE, Category: CoreKeyword},
		kwDECRYPTION: {Name: "DECRYPTION", Token: kwDECRYPTION, Category: ContextKeyword},
		kwDEFAULT: {Name: "DEFAULT", Token: kwDEFAULT, Category: CoreKeyword},
		kwDELAY: {Name: "DELAY", Token: kwDELAY, Category: ContextKeyword},
		kwDELAYED_DURABILITY: {Name: "DELAYED_DURABILITY", Token: kwDELAYED_DURABILITY, Category: ContextKeyword},
		kwDELETE: {Name: "DELETE", Token: kwDELETE, Category: CoreKeyword},
		kwDENY: {Name: "DENY", Token: kwDENY, Category: CoreKeyword},
		kwDEPENDENTS: {Name: "DEPENDENTS", Token: kwDEPENDENTS, Category: ContextKeyword},
		kwDESC: {Name: "DESC", Token: kwDESC, Category: CoreKeyword},
		kwDESCRIPTION: {Name: "DESCRIPTION", Token: kwDESCRIPTION, Category: ContextKeyword},
		kwDIAGNOSTICS: {Name: "DIAGNOSTICS", Token: kwDIAGNOSTICS, Category: ContextKeyword},
		kwDIALOG: {Name: "DIALOG", Token: kwDIALOG, Category: ContextKeyword},
		kwDIRECTORY_NAME: {Name: "DIRECTORY_NAME", Token: kwDIRECTORY_NAME, Category: ContextKeyword},
		kwDISABLE: {Name: "DISABLE", Token: kwDISABLE, Category: ContextKeyword},
		kwDISTINCT: {Name: "DISTINCT", Token: kwDISTINCT, Category: CoreKeyword},
		kwDISTRIBUTED: {Name: "DISTRIBUTED", Token: kwDISTRIBUTED, Category: CoreKeyword},
		kwDISTRIBUTION: {Name: "DISTRIBUTION", Token: kwDISTRIBUTION, Category: ContextKeyword},
		kwDO: {Name: "DO", Token: kwDO, Category: ContextKeyword},
		kwDOUBLE: {Name: "DOUBLE", Token: kwDOUBLE, Category: CoreKeyword},
		kwDROP: {Name: "DROP", Token: kwDROP, Category: CoreKeyword},
		kwDUMP: {Name: "DUMP", Token: kwDUMP, Category: ContextKeyword},
		kwDYNAMIC: {Name: "DYNAMIC", Token: kwDYNAMIC, Category: ContextKeyword},
		kwEDGE: {Name: "EDGE", Token: kwEDGE, Category: ContextKeyword},
		kwELEMENTS: {Name: "ELEMENTS", Token: kwELEMENTS, Category: ContextKeyword},
		kwELSE: {Name: "ELSE", Token: kwELSE, Category: CoreKeyword},
		kwENABLE: {Name: "ENABLE", Token: kwENABLE, Category: ContextKeyword},
		kwENCRYPTED: {Name: "ENCRYPTED", Token: kwENCRYPTED, Category: ContextKeyword},
		kwENCRYPTION: {Name: "ENCRYPTION", Token: kwENCRYPTION, Category: ContextKeyword},
		kwEND: {Name: "END", Token: kwEND, Category: CoreKeyword},
		kwENDPOINT: {Name: "ENDPOINT", Token: kwENDPOINT, Category: ContextKeyword},
		kwERRLVL: {Name: "ERRLVL", Token: kwERRLVL, Category: CoreKeyword},
		kwERROR: {Name: "ERROR", Token: kwERROR, Category: ContextKeyword},
		kwESCAPE: {Name: "ESCAPE", Token: kwESCAPE, Category: CoreKeyword},
		kwEVENT: {Name: "EVENT", Token: kwEVENT, Category: ContextKeyword},
		kwEXCEPT: {Name: "EXCEPT", Token: kwEXCEPT, Category: CoreKeyword},
		kwEXEC: {Name: "EXEC", Token: kwEXEC, Category: CoreKeyword},
		kwEXECUTE: {Name: "EXECUTE", Token: kwEXECUTE, Category: CoreKeyword},
		kwEXISTS: {Name: "EXISTS", Token: kwEXISTS, Category: CoreKeyword},
		kwEXIT: {Name: "EXIT", Token: kwEXIT, Category: CoreKeyword},
		kwEXPAND: {Name: "EXPAND", Token: kwEXPAND, Category: ContextKeyword},
		kwEXTENSION: {Name: "EXTENSION", Token: kwEXTENSION, Category: ContextKeyword},
		kwEXTERNAL: {Name: "EXTERNAL", Token: kwEXTERNAL, Category: CoreKeyword},
		kwFAILOVER: {Name: "FAILOVER", Token: kwFAILOVER, Category: ContextKeyword},
		kwFAN_IN: {Name: "FAN_IN", Token: kwFAN_IN, Category: ContextKeyword},
		kwFAST: {Name: "FAST", Token: kwFAST, Category: ContextKeyword},
		kwFAST_FORWARD: {Name: "FAST_FORWARD", Token: kwFAST_FORWARD, Category: ContextKeyword},
		kwFEDERATION: {Name: "FEDERATION", Token: kwFEDERATION, Category: ContextKeyword},
		kwFETCH: {Name: "FETCH", Token: kwFETCH, Category: CoreKeyword},
		kwFILE: {Name: "FILE", Token: kwFILE, Category: CoreKeyword},
		kwFILEGROUP: {Name: "FILEGROUP", Token: kwFILEGROUP, Category: ContextKeyword},
		kwFILENAME: {Name: "FILENAME", Token: kwFILENAME, Category: ContextKeyword},
		kwFILESTREAM: {Name: "FILESTREAM", Token: kwFILESTREAM, Category: ContextKeyword},
		kwFILESTREAM_ON: {Name: "FILESTREAM_ON", Token: kwFILESTREAM_ON, Category: ContextKeyword},
		kwFILETABLE: {Name: "FILETABLE", Token: kwFILETABLE, Category: ContextKeyword},
		kwFILETABLE_NAMESPACE: {Name: "FILETABLE_NAMESPACE", Token: kwFILETABLE_NAMESPACE, Category: ContextKeyword},
		kwFILLFACTOR: {Name: "FILLFACTOR", Token: kwFILLFACTOR, Category: CoreKeyword},
		kwFILTER: {Name: "FILTER", Token: kwFILTER, Category: ContextKeyword},
		kwFILTERING: {Name: "FILTERING", Token: kwFILTERING, Category: ContextKeyword},
		kwFIRST: {Name: "FIRST", Token: kwFIRST, Category: ContextKeyword},
		kwFOLLOWING: {Name: "FOLLOWING", Token: kwFOLLOWING, Category: ContextKeyword},
		kwFOR: {Name: "FOR", Token: kwFOR, Category: CoreKeyword},
		kwFOR_APPEND: {Name: "FOR_APPEND", Token: kwFOR_APPEND, Category: ContextKeyword},
		kwFORCE: {Name: "FORCE", Token: kwFORCE, Category: ContextKeyword},
		kwFORCE_FAILOVER_ALLOW_DATA_LOSS: {Name: "FORCE_FAILOVER_ALLOW_DATA_LOSS", Token: kwFORCE_FAILOVER_ALLOW_DATA_LOSS, Category: ContextKeyword},
		kwFOREIGN: {Name: "FOREIGN", Token: kwFOREIGN, Category: CoreKeyword},
		kwFORMAT: {Name: "FORMAT", Token: kwFORMAT, Category: ContextKeyword},
		kwFORWARD_ONLY: {Name: "FORWARD_ONLY", Token: kwFORWARD_ONLY, Category: ContextKeyword},
		kwFREETEXT: {Name: "FREETEXT", Token: kwFREETEXT, Category: CoreKeyword},
		kwFREETEXTTABLE: {Name: "FREETEXTTABLE", Token: kwFREETEXTTABLE, Category: CoreKeyword},
		kwFROM: {Name: "FROM", Token: kwFROM, Category: CoreKeyword},
		kwFULL: {Name: "FULL", Token: kwFULL, Category: CoreKeyword},
		kwFULLTEXT: {Name: "FULLTEXT", Token: kwFULLTEXT, Category: ContextKeyword},
		kwFUNCTION: {Name: "FUNCTION", Token: kwFUNCTION, Category: CoreKeyword},
		kwGB: {Name: "GB", Token: kwGB, Category: ContextKeyword},
		kwGENERATED: {Name: "GENERATED", Token: kwGENERATED, Category: ContextKeyword},
		kwGET: {Name: "GET", Token: kwGET, Category: ContextKeyword},
		kwGLOBAL: {Name: "GLOBAL", Token: kwGLOBAL, Category: ContextKeyword},
		kwGO: {Name: "GO", Token: kwGO, Category: ContextKeyword},
		kwGOTO: {Name: "GOTO", Token: kwGOTO, Category: CoreKeyword},
		kwGOVERNOR: {Name: "GOVERNOR", Token: kwGOVERNOR, Category: ContextKeyword},
		kwGRANT: {Name: "GRANT", Token: kwGRANT, Category: CoreKeyword},
		kwGROUP: {Name: "GROUP", Token: kwGROUP, Category: CoreKeyword},
		kwGROUPING: {Name: "GROUPING", Token: kwGROUPING, Category: ContextKeyword},
		kwGROUPS: {Name: "GROUPS", Token: kwGROUPS, Category: ContextKeyword},
		kwHADR: {Name: "HADR", Token: kwHADR, Category: ContextKeyword},
		kwHARDWARE_OFFLOAD: {Name: "HARDWARE_OFFLOAD", Token: kwHARDWARE_OFFLOAD, Category: ContextKeyword},
		kwHASH: {Name: "HASH", Token: kwHASH, Category: ContextKeyword},
		kwHASHED: {Name: "HASHED", Token: kwHASHED, Category: ContextKeyword},
		kwHAVING: {Name: "HAVING", Token: kwHAVING, Category: CoreKeyword},
		kwHEAP: {Name: "HEAP", Token: kwHEAP, Category: ContextKeyword},
		kwHIDDEN: {Name: "HIDDEN", Token: kwHIDDEN, Category: ContextKeyword},
		kwHIGH: {Name: "HIGH", Token: kwHIGH, Category: ContextKeyword},
		kwHINT: {Name: "HINT", Token: kwHINT, Category: ContextKeyword},
		kwHOLDLOCK: {Name: "HOLDLOCK", Token: kwHOLDLOCK, Category: CoreKeyword},
		kwHOURS: {Name: "HOURS", Token: kwHOURS, Category: ContextKeyword},
		kwHTTP: {Name: "HTTP", Token: kwHTTP, Category: ContextKeyword},
		kwIDENTITY: {Name: "IDENTITY", Token: kwIDENTITY, Category: CoreKeyword},
		kwIDENTITY_INSERT: {Name: "IDENTITY_INSERT", Token: kwIDENTITY_INSERT, Category: CoreKeyword},
		kwIDENTITYCOL: {Name: "IDENTITYCOL", Token: kwIDENTITYCOL, Category: CoreKeyword},
		kwIF: {Name: "IF", Token: kwIF, Category: CoreKeyword},
		kwIIF: {Name: "IIF", Token: kwIIF, Category: ContextKeyword},
		kwIMMEDIATE: {Name: "IMMEDIATE", Token: kwIMMEDIATE, Category: ContextKeyword},
		kwIN: {Name: "IN", Token: kwIN, Category: CoreKeyword},
		kwINCLUDE: {Name: "INCLUDE", Token: kwINCLUDE, Category: ContextKeyword},
		kwINCLUDE_NULL_VALUES: {Name: "INCLUDE_NULL_VALUES", Token: kwINCLUDE_NULL_VALUES, Category: ContextKeyword},
		kwINCREMENT: {Name: "INCREMENT", Token: kwINCREMENT, Category: ContextKeyword},
		kwINDEX: {Name: "INDEX", Token: kwINDEX, Category: CoreKeyword},
		kwINNER: {Name: "INNER", Token: kwINNER, Category: CoreKeyword},
		kwINPUT: {Name: "INPUT", Token: kwINPUT, Category: ContextKeyword},
		kwINSENSITIVE: {Name: "INSENSITIVE", Token: kwINSENSITIVE, Category: ContextKeyword},
		kwINSERT: {Name: "INSERT", Token: kwINSERT, Category: CoreKeyword},
		kwINSTEAD: {Name: "INSTEAD", Token: kwINSTEAD, Category: ContextKeyword},
		kwINTERSECT: {Name: "INTERSECT", Token: kwINTERSECT, Category: CoreKeyword},
		kwINTO: {Name: "INTO", Token: kwINTO, Category: CoreKeyword},
		kwIS: {Name: "IS", Token: kwIS, Category: CoreKeyword},
		kwISOLATION: {Name: "ISOLATION", Token: kwISOLATION, Category: ContextKeyword},
		kwJOB: {Name: "JOB", Token: kwJOB, Category: ContextKeyword},
		kwJOIN: {Name: "JOIN", Token: kwJOIN, Category: CoreKeyword},
		kwJSON: {Name: "JSON", Token: kwJSON, Category: ContextKeyword},
		kwKB: {Name: "KB", Token: kwKB, Category: ContextKeyword},
		kwKEEP: {Name: "KEEP", Token: kwKEEP, Category: ContextKeyword},
		kwKEEPFIXED: {Name: "KEEPFIXED", Token: kwKEEPFIXED, Category: ContextKeyword},
		kwKEY:  {Name: "KEY", Token: kwKEY, Category: CoreKeyword},
		kwKEYS: {Name: "KEYS", Token: kwKEYS, Category: ContextKeyword},
		kwKEYSET: {Name: "KEYSET", Token: kwKEYSET, Category: ContextKeyword},
		kwLAST: {Name: "LAST", Token: kwLAST, Category: ContextKeyword},
		kwKILL: {Name: "KILL", Token: kwKILL, Category: CoreKeyword},
		kwLANGUAGE: {Name: "LANGUAGE", Token: kwLANGUAGE, Category: ContextKeyword},
		kwLEFT: {Name: "LEFT", Token: kwLEFT, Category: CoreKeyword},
		kwLEVEL: {Name: "LEVEL", Token: kwLEVEL, Category: ContextKeyword},
		kwLIBRARY: {Name: "LIBRARY", Token: kwLIBRARY, Category: ContextKeyword},
		kwLIFETIME: {Name: "LIFETIME", Token: kwLIFETIME, Category: ContextKeyword},
		kwLIKE: {Name: "LIKE", Token: kwLIKE, Category: CoreKeyword},
		kwLINENO: {Name: "LINENO", Token: kwLINENO, Category: CoreKeyword},
		kwLIST: {Name: "LIST", Token: kwLIST, Category: ContextKeyword},
		kwLISTENER: {Name: "LISTENER", Token: kwLISTENER, Category: ContextKeyword},
		kwLISTENER_IP: {Name: "LISTENER_IP", Token: kwLISTENER_IP, Category: ContextKeyword},
		kwLISTENER_PORT: {Name: "LISTENER_PORT", Token: kwLISTENER_PORT, Category: ContextKeyword},
		kwLOAD: {Name: "LOAD", Token: kwLOAD, Category: ContextKeyword},
		kwLOB_COMPACTION: {Name: "LOB_COMPACTION", Token: kwLOB_COMPACTION, Category: ContextKeyword},
		kwLOCAL: {Name: "LOCAL", Token: kwLOCAL, Category: ContextKeyword},
		kwLOG: {Name: "LOG", Token: kwLOG, Category: ContextKeyword},
		kwLOGIN: {Name: "LOGIN", Token: kwLOGIN, Category: ContextKeyword},
		kwLOGON: {Name: "LOGON", Token: kwLOGON, Category: ContextKeyword},
		kwLOOP: {Name: "LOOP", Token: kwLOOP, Category: ContextKeyword},
		kwLOW: {Name: "LOW", Token: kwLOW, Category: ContextKeyword},
		kwMANUAL: {Name: "MANUAL", Token: kwMANUAL, Category: ContextKeyword},
		kwMANUAL_CUTOVER: {Name: "MANUAL_CUTOVER", Token: kwMANUAL_CUTOVER, Category: ContextKeyword},
		kwMARK: {Name: "MARK", Token: kwMARK, Category: ContextKeyword},
		kwMASKED: {Name: "MASKED", Token: kwMASKED, Category: ContextKeyword},
		kwMASTER: {Name: "MASTER", Token: kwMASTER, Category: ContextKeyword},
		kwMATCHED: {Name: "MATCHED", Token: kwMATCHED, Category: ContextKeyword},
		kwMATERIALIZED: {Name: "MATERIALIZED", Token: kwMATERIALIZED, Category: ContextKeyword},
		kwMAX: {Name: "MAX", Token: kwMAX, Category: ContextKeyword},
		kwMAXVALUE: {Name: "MAXVALUE", Token: kwMAXVALUE, Category: ContextKeyword},
		kwMAX_QUEUE_READERS: {Name: "MAX_QUEUE_READERS", Token: kwMAX_QUEUE_READERS, Category: ContextKeyword},
		kwMAXDOP: {Name: "MAXDOP", Token: kwMAXDOP, Category: ContextKeyword},
		kwMAXRECURSION: {Name: "MAXRECURSION", Token: kwMAXRECURSION, Category: ContextKeyword},
		kwMB: {Name: "MB", Token: kwMB, Category: ContextKeyword},
		kwMEMBER: {Name: "MEMBER", Token: kwMEMBER, Category: ContextKeyword},
		kwMEMORY_OPTIMIZED: {Name: "MEMORY_OPTIMIZED", Token: kwMEMORY_OPTIMIZED, Category: ContextKeyword},
		kwMEMORY_OPTIMIZED_DATA: {Name: "MEMORY_OPTIMIZED_DATA", Token: kwMEMORY_OPTIMIZED_DATA, Category: ContextKeyword},
		kwMERGE: {Name: "MERGE", Token: kwMERGE, Category: CoreKeyword},
		kwMESSAGE: {Name: "MESSAGE", Token: kwMESSAGE, Category: ContextKeyword},
		kwMESSAGE_FORWARD_SIZE: {Name: "MESSAGE_FORWARD_SIZE", Token: kwMESSAGE_FORWARD_SIZE, Category: ContextKeyword},
		kwMESSAGE_FORWARDING: {Name: "MESSAGE_FORWARDING", Token: kwMESSAGE_FORWARDING, Category: ContextKeyword},
		kwMINUTES: {Name: "MINUTES", Token: kwMINUTES, Category: ContextKeyword},
		kwMINVALUE: {Name: "MINVALUE", Token: kwMINVALUE, Category: ContextKeyword},
		kwMIRROR: {Name: "MIRROR", Token: kwMIRROR, Category: ContextKeyword},
		kwMIRRORING: {Name: "MIRRORING", Token: kwMIRRORING, Category: ContextKeyword},
		kwMODE: {Name: "MODE", Token: kwMODE, Category: ContextKeyword},
		kwMODEL: {Name: "MODEL", Token: kwMODEL, Category: ContextKeyword},
		kwMODIFY: {Name: "MODIFY", Token: kwMODIFY, Category: ContextKeyword},
		kwMOVE: {Name: "MOVE", Token: kwMOVE, Category: ContextKeyword},
		kwMUST_CHANGE: {Name: "MUST_CHANGE", Token: kwMUST_CHANGE, Category: ContextKeyword},
		kwNAME: {Name: "NAME", Token: kwNAME, Category: ContextKeyword},
		kwNATIONAL: {Name: "NATIONAL", Token: kwNATIONAL, Category: CoreKeyword},
		kwNATIVE_COMPILATION: {Name: "NATIVE_COMPILATION", Token: kwNATIVE_COMPILATION, Category: ContextKeyword},
		kwNEXT: {Name: "NEXT", Token: kwNEXT, Category: ContextKeyword},
		kwNO: {Name: "NO", Token: kwNO, Category: ContextKeyword},
		kwNOCHECK: {Name: "NOCHECK", Token: kwNOCHECK, Category: CoreKeyword},
		kwNOCOUNT: {Name: "NOCOUNT", Token: kwNOCOUNT, Category: ContextKeyword},
		kwNODE: {Name: "NODE", Token: kwNODE, Category: ContextKeyword},
		kwNOLOCK: {Name: "NOLOCK", Token: kwNOLOCK, Category: ContextKeyword},
		kwNONCLUSTERED: {Name: "NONCLUSTERED", Token: kwNONCLUSTERED, Category: CoreKeyword},
		kwNONE: {Name: "NONE", Token: kwNONE, Category: ContextKeyword},
		kwNORESET: {Name: "NORESET", Token: kwNORESET, Category: ContextKeyword},
		kwNOT: {Name: "NOT", Token: kwNOT, Category: CoreKeyword},
		kwNOTIFICATION: {Name: "NOTIFICATION", Token: kwNOTIFICATION, Category: ContextKeyword},
		kwNOWAIT: {Name: "NOWAIT", Token: kwNOWAIT, Category: ContextKeyword},
		kwNULL: {Name: "NULL", Token: kwNULL, Category: CoreKeyword},
		kwNULLIF: {Name: "NULLIF", Token: kwNULLIF, Category: CoreKeyword},
		kwNUMANODE: {Name: "NUMANODE", Token: kwNUMANODE, Category: ContextKeyword},
		kwOBJECT: {Name: "OBJECT", Token: kwOBJECT, Category: ContextKeyword},
		kwOF: {Name: "OF", Token: kwOF, Category: CoreKeyword},
		kwOFF: {Name: "OFF", Token: kwOFF, Category: CoreKeyword},
		kwOFFLINE: {Name: "OFFLINE", Token: kwOFFLINE, Category: ContextKeyword},
		kwOFFSET: {Name: "OFFSET", Token: kwOFFSET, Category: ContextKeyword},
		kwOFFSETS: {Name: "OFFSETS", Token: kwOFFSETS, Category: CoreKeyword},
		kwOLD_PASSWORD: {Name: "OLD_PASSWORD", Token: kwOLD_PASSWORD, Category: ContextKeyword},
		kwON: {Name: "ON", Token: kwON, Category: CoreKeyword},
		kwONLY: {Name: "ONLY", Token: kwONLY, Category: ContextKeyword},
		kwOPEN: {Name: "OPEN", Token: kwOPEN, Category: CoreKeyword},
		kwOPENDATASOURCE: {Name: "OPENDATASOURCE", Token: kwOPENDATASOURCE, Category: CoreKeyword},
		kwOPENJSON: {Name: "OPENJSON", Token: kwOPENJSON, Category: ContextKeyword},
		kwOPENQUERY: {Name: "OPENQUERY", Token: kwOPENQUERY, Category: CoreKeyword},
		kwOPENROWSET: {Name: "OPENROWSET", Token: kwOPENROWSET, Category: CoreKeyword},
		kwOPENXML: {Name: "OPENXML", Token: kwOPENXML, Category: CoreKeyword},
		kwOPTIMIZE: {Name: "OPTIMIZE", Token: kwOPTIMIZE, Category: ContextKeyword},
		kwOPTIMISTIC: {Name: "OPTIMISTIC", Token: kwOPTIMISTIC, Category: ContextKeyword},
		kwOPTION: {Name: "OPTION", Token: kwOPTION, Category: CoreKeyword},
		kwOR: {Name: "OR", Token: kwOR, Category: CoreKeyword},
		kwORDER: {Name: "ORDER", Token: kwORDER, Category: CoreKeyword},
		kwOUT: {Name: "OUT", Token: kwOUT, Category: ContextKeyword},
		kwOUTER: {Name: "OUTER", Token: kwOUTER, Category: CoreKeyword},
		kwOUTPUT: {Name: "OUTPUT", Token: kwOUTPUT, Category: ContextKeyword},
		kwOVER: {Name: "OVER", Token: kwOVER, Category: CoreKeyword},
		kwOVERRIDE: {Name: "OVERRIDE", Token: kwOVERRIDE, Category: ContextKeyword},
		kwOWNER: {Name: "OWNER", Token: kwOWNER, Category: ContextKeyword},
		kwPAGE: {Name: "PAGE", Token: kwPAGE, Category: ContextKeyword},
		kwPARAMETERIZATION: {Name: "PARAMETERIZATION", Token: kwPARAMETERIZATION, Category: ContextKeyword},
		kwPARTITION: {Name: "PARTITION", Token: kwPARTITION, Category: ContextKeyword},
		kwPARTITIONS: {Name: "PARTITIONS", Token: kwPARTITIONS, Category: ContextKeyword},
		kwPASSWORD: {Name: "PASSWORD", Token: kwPASSWORD, Category: ContextKeyword},
		kwPATH: {Name: "PATH", Token: kwPATH, Category: ContextKeyword},
		kwPAUSE: {Name: "PAUSE", Token: kwPAUSE, Category: ContextKeyword},
		kwPERCENT: {Name: "PERCENT", Token: kwPERCENT, Category: CoreKeyword},
		kwPERIOD: {Name: "PERIOD", Token: kwPERIOD, Category: ContextKeyword},
		kwPERMISSION_SET: {Name: "PERMISSION_SET", Token: kwPERMISSION_SET, Category: ContextKeyword},
		kwPERSISTED: {Name: "PERSISTED", Token: kwPERSISTED, Category: ContextKeyword},
		kwPIVOT: {Name: "PIVOT", Token: kwPIVOT, Category: CoreKeyword},
		kwPLAN: {Name: "PLAN", Token: kwPLAN, Category: CoreKeyword},
		kwPLATFORM: {Name: "PLATFORM", Token: kwPLATFORM, Category: ContextKeyword},
		kwPOISON_MESSAGE_HANDLING: {Name: "POISON_MESSAGE_HANDLING", Token: kwPOISON_MESSAGE_HANDLING, Category: ContextKeyword},
		kwPOLICY: {Name: "POLICY", Token: kwPOLICY, Category: ContextKeyword},
		kwPOOL: {Name: "POOL", Token: kwPOOL, Category: ContextKeyword},
		kwPOPULATION: {Name: "POPULATION", Token: kwPOPULATION, Category: ContextKeyword},
		kwPRECEDING: {Name: "PRECEDING", Token: kwPRECEDING, Category: ContextKeyword},
		kwPRECISION: {Name: "PRECISION", Token: kwPRECISION, Category: ContextKeyword},
		kwPREDICATE: {Name: "PREDICATE", Token: kwPREDICATE, Category: ContextKeyword},
		kwPREDICT: {Name: "PREDICT", Token: kwPREDICT, Category: ContextKeyword},
		kwPRIMARY: {Name: "PRIMARY", Token: kwPRIMARY, Category: CoreKeyword},
		kwPRINT: {Name: "PRINT", Token: kwPRINT, Category: CoreKeyword},
		kwPRIOR: {Name: "PRIOR", Token: kwPRIOR, Category: ContextKeyword},
		kwPRIORITY: {Name: "PRIORITY", Token: kwPRIORITY, Category: ContextKeyword},
		kwPRIVILEGES: {Name: "PRIVILEGES", Token: kwPRIVILEGES, Category: ContextKeyword},
		kwPROC: {Name: "PROC", Token: kwPROC, Category: CoreKeyword},
		kwPROCEDURE: {Name: "PROCEDURE", Token: kwPROCEDURE, Category: CoreKeyword},
		kwPROCEDURE_CACHE: {Name: "PROCEDURE_CACHE", Token: kwPROCEDURE_CACHE, Category: ContextKeyword},
		kwPROCEDURE_NAME: {Name: "PROCEDURE_NAME", Token: kwPROCEDURE_NAME, Category: ContextKeyword},
		kwPROCESS: {Name: "PROCESS", Token: kwPROCESS, Category: ContextKeyword},
		kwPROPERTY: {Name: "PROPERTY", Token: kwPROPERTY, Category: ContextKeyword},
		kwPROVIDER: {Name: "PROVIDER", Token: kwPROVIDER, Category: ContextKeyword},
		kwPUBLIC: {Name: "PUBLIC", Token: kwPUBLIC, Category: CoreKeyword},
		kwQUERY: {Name: "QUERY", Token: kwQUERY, Category: ContextKeyword},
		kwQUERYTRACEON: {Name: "QUERYTRACEON", Token: kwQUERYTRACEON, Category: ContextKeyword},
		kwQUEUE: {Name: "QUEUE", Token: kwQUEUE, Category: ContextKeyword},
		kwRAISERROR: {Name: "RAISERROR", Token: kwRAISERROR, Category: CoreKeyword},
		kwRANGE: {Name: "RANGE", Token: kwRANGE, Category: ContextKeyword},
		kwRAW: {Name: "RAW", Token: kwRAW, Category: ContextKeyword},
		kwREAD: {Name: "READ", Token: kwREAD, Category: CoreKeyword},
		kwREAD_ONLY: {Name: "READ_ONLY", Token: kwREAD_ONLY, Category: ContextKeyword},
		kwREAD_WRITE_FILEGROUPS: {Name: "READ_WRITE_FILEGROUPS", Token: kwREAD_WRITE_FILEGROUPS, Category: ContextKeyword},
		kwREADONLY: {Name: "READONLY", Token: kwREADONLY, Category: ContextKeyword},
		kwREADTEXT: {Name: "READTEXT", Token: kwREADTEXT, Category: CoreKeyword},
		kwREBUILD: {Name: "REBUILD", Token: kwREBUILD, Category: ContextKeyword},
		kwRECEIVE: {Name: "RECEIVE", Token: kwRECEIVE, Category: ContextKeyword},
		kwRECOMPILE: {Name: "RECOMPILE", Token: kwRECOMPILE, Category: ContextKeyword},
		kwRECONFIGURE: {Name: "RECONFIGURE", Token: kwRECONFIGURE, Category: CoreKeyword},
		kwREFERENCES: {Name: "REFERENCES", Token: kwREFERENCES, Category: CoreKeyword},
		kwREGENERATE: {Name: "REGENERATE", Token: kwREGENERATE, Category: ContextKeyword},
		kwRELATED_CONVERSATION: {Name: "RELATED_CONVERSATION", Token: kwRELATED_CONVERSATION, Category: ContextKeyword},
		kwRELATED_CONVERSATION_GROUP: {Name: "RELATED_CONVERSATION_GROUP", Token: kwRELATED_CONVERSATION_GROUP, Category: ContextKeyword},
		kwRELATIVE: {Name: "RELATIVE", Token: kwRELATIVE, Category: ContextKeyword},
		kwREMOTE: {Name: "REMOTE", Token: kwREMOTE, Category: ContextKeyword},
		kwREMOVE: {Name: "REMOVE", Token: kwREMOVE, Category: ContextKeyword},
		kwRENAME: {Name: "RENAME", Token: kwRENAME, Category: ContextKeyword},
		kwREORGANIZE: {Name: "REORGANIZE", Token: kwREORGANIZE, Category: ContextKeyword},
		kwREPEATABLE: {Name: "REPEATABLE", Token: kwREPEATABLE, Category: ContextKeyword},
		kwRESET: {Name: "RESET", Token: kwRESET, Category: ContextKeyword},
		kwREPLICA: {Name: "REPLICA", Token: kwREPLICA, Category: ContextKeyword},
		kwREPLICATION: {Name: "REPLICATION", Token: kwREPLICATION, Category: CoreKeyword},
		kwRESAMPLE: {Name: "RESAMPLE", Token: kwRESAMPLE, Category: ContextKeyword},
		kwRESOURCE: {Name: "RESOURCE", Token: kwRESOURCE, Category: ContextKeyword},
		kwRESOURCE_POOL: {Name: "RESOURCE_POOL", Token: kwRESOURCE_POOL, Category: ContextKeyword},
		kwRESTART: {Name: "RESTART", Token: kwRESTART, Category: ContextKeyword},
		kwRESTORE: {Name: "RESTORE", Token: kwRESTORE, Category: CoreKeyword},
		kwRESTRICT: {Name: "RESTRICT", Token: kwRESTRICT, Category: CoreKeyword},
		kwRESULT: {Name: "RESULT", Token: kwRESULT, Category: ContextKeyword},
		kwRESUME: {Name: "RESUME", Token: kwRESUME, Category: ContextKeyword},
		kwRETENTION: {Name: "RETENTION", Token: kwRETENTION, Category: ContextKeyword},
		kwRETURN: {Name: "RETURN", Token: kwRETURN, Category: CoreKeyword},
		kwRETURNS: {Name: "RETURNS", Token: kwRETURNS, Category: ContextKeyword},
		kwREVERT: {Name: "REVERT", Token: kwREVERT, Category: CoreKeyword},
		kwREVOKE: {Name: "REVOKE", Token: kwREVOKE, Category: CoreKeyword},
		kwRIGHT: {Name: "RIGHT", Token: kwRIGHT, Category: CoreKeyword},
		kwROBUST: {Name: "ROBUST", Token: kwROBUST, Category: ContextKeyword},
		kwROLE: {Name: "ROLE", Token: kwROLE, Category: ContextKeyword},
		kwROLLBACK: {Name: "ROLLBACK", Token: kwROLLBACK, Category: CoreKeyword},
		kwROLLUP: {Name: "ROLLUP", Token: kwROLLUP, Category: ContextKeyword},
		kwROOT: {Name: "ROOT", Token: kwROOT, Category: ContextKeyword},
		kwROUND_ROBIN: {Name: "ROUND_ROBIN", Token: kwROUND_ROBIN, Category: ContextKeyword},
		kwROUTE: {Name: "ROUTE", Token: kwROUTE, Category: ContextKeyword},
		kwROW: {Name: "ROW", Token: kwROW, Category: ContextKeyword},
		kwROWCOUNT: {Name: "ROWCOUNT", Token: kwROWCOUNT, Category: CoreKeyword},
		kwROWGUIDCOL: {Name: "ROWGUIDCOL", Token: kwROWGUIDCOL, Category: CoreKeyword},
		kwROWS: {Name: "ROWS", Token: kwROWS, Category: ContextKeyword},
		kwRULE: {Name: "RULE", Token: kwRULE, Category: CoreKeyword},
		kwSAMPLE: {Name: "SAMPLE", Token: kwSAMPLE, Category: ContextKeyword},
		kwSAVE: {Name: "SAVE", Token: kwSAVE, Category: CoreKeyword},
		kwSCHEDULER: {Name: "SCHEDULER", Token: kwSCHEDULER, Category: ContextKeyword},
		kwSCHEMA: {Name: "SCHEMA", Token: kwSCHEMA, Category: CoreKeyword},
		kwSCHEMABINDING: {Name: "SCHEMABINDING", Token: kwSCHEMABINDING, Category: ContextKeyword},
		kwSCHEME: {Name: "SCHEME", Token: kwSCHEME, Category: ContextKeyword},
		kwSCROLL: {Name: "SCROLL", Token: kwSCROLL, Category: ContextKeyword},
		kwSCROLL_LOCKS: {Name: "SCROLL_LOCKS", Token: kwSCROLL_LOCKS, Category: ContextKeyword},
		kwSCOPED: {Name: "SCOPED", Token: kwSCOPED, Category: ContextKeyword},
		kwSEARCH: {Name: "SEARCH", Token: kwSEARCH, Category: ContextKeyword},
		kwSECONDARY: {Name: "SECONDARY", Token: kwSECONDARY, Category: ContextKeyword},
		kwSECONDS: {Name: "SECONDS", Token: kwSECONDS, Category: ContextKeyword},
		kwSECURITY: {Name: "SECURITY", Token: kwSECURITY, Category: ContextKeyword},
		kwSECURITYAUDIT: {Name: "SECURITYAUDIT", Token: kwSECURITYAUDIT, Category: ContextKeyword},
		kwSELECT: {Name: "SELECT", Token: kwSELECT, Category: CoreKeyword},
		kwSELECTIVE: {Name: "SELECTIVE", Token: kwSELECTIVE, Category: ContextKeyword},
		kwSELF: {Name: "SELF", Token: kwSELF, Category: ContextKeyword},
		kwSEMANTICKEYPHRASETABLE: {Name: "SEMANTICKEYPHRASETABLE", Token: kwSEMANTICKEYPHRASETABLE, Category: CoreKeyword},
		kwSEMANTICSIMILARITYDETAILSTABLE: {Name: "SEMANTICSIMILARITYDETAILSTABLE", Token: kwSEMANTICSIMILARITYDETAILSTABLE, Category: CoreKeyword},
		kwSEMANTICSIMILARITYTABLE: {Name: "SEMANTICSIMILARITYTABLE", Token: kwSEMANTICSIMILARITYTABLE, Category: CoreKeyword},
		kwSEMIJOIN: {Name: "SEMIJOIN", Token: kwSEMIJOIN, Category: ContextKeyword},
		kwSEND: {Name: "SEND", Token: kwSEND, Category: ContextKeyword},
		kwSEQUENCE: {Name: "SEQUENCE", Token: kwSEQUENCE, Category: ContextKeyword},
		kwSENSITIVITY: {Name: "SENSITIVITY", Token: kwSENSITIVITY, Category: ContextKeyword},
		kwSENT: {Name: "SENT", Token: kwSENT, Category: ContextKeyword},
		kwSERIALIZABLE: {Name: "SERIALIZABLE", Token: kwSERIALIZABLE, Category: ContextKeyword},
		kwSERVER: {Name: "SERVER", Token: kwSERVER, Category: ContextKeyword},
		kwSERVICE: {Name: "SERVICE", Token: kwSERVICE, Category: ContextKeyword},
		kwSESSION: {Name: "SESSION", Token: kwSESSION, Category: ContextKeyword},
		kwSESSION_USER: {Name: "SESSION_USER", Token: kwSESSION_USER, Category: CoreKeyword},
		kwSET: {Name: "SET", Token: kwSET, Category: CoreKeyword},
		kwSETS: {Name: "SETS", Token: kwSETS, Category: ContextKeyword},
		kwSETUSER: {Name: "SETUSER", Token: kwSETUSER, Category: CoreKeyword},
		kwSHUTDOWN: {Name: "SHUTDOWN", Token: kwSHUTDOWN, Category: CoreKeyword},
		kwSIGNATURE: {Name: "SIGNATURE", Token: kwSIGNATURE, Category: ContextKeyword},
		kwSIZE: {Name: "SIZE", Token: kwSIZE, Category: ContextKeyword},
		kwSNAPSHOT: {Name: "SNAPSHOT", Token: kwSNAPSHOT, Category: ContextKeyword},
		kwSOFTNUMA: {Name: "SOFTNUMA", Token: kwSOFTNUMA, Category: ContextKeyword},
		kwSOME: {Name: "SOME", Token: kwSOME, Category: CoreKeyword},
		kwSOURCE: {Name: "SOURCE", Token: kwSOURCE, Category: ContextKeyword},
		kwSPARSE: {Name: "SPARSE", Token: kwSPARSE, Category: ContextKeyword},
		kwSPATIAL: {Name: "SPATIAL", Token: kwSPATIAL, Category: ContextKeyword},
		kwSPECIFICATION: {Name: "SPECIFICATION", Token: kwSPECIFICATION, Category: ContextKeyword},
		kwSPLIT: {Name: "SPLIT", Token: kwSPLIT, Category: ContextKeyword},
		kwSTART: {Name: "START", Token: kwSTART, Category: ContextKeyword},
		kwSTATE: {Name: "STATE", Token: kwSTATE, Category: ContextKeyword},
		kwSTATIC: {Name: "STATIC", Token: kwSTATIC, Category: ContextKeyword},
		kwSTATISTICAL_SEMANTICS: {Name: "STATISTICAL_SEMANTICS", Token: kwSTATISTICAL_SEMANTICS, Category: ContextKeyword},
		kwSTATISTICS: {Name: "STATISTICS", Token: kwSTATISTICS, Category: CoreKeyword},
		kwSTATS: {Name: "STATS", Token: kwSTATS, Category: ContextKeyword},
		kwSTATUS: {Name: "STATUS", Token: kwSTATUS, Category: ContextKeyword},
		kwSTATUSONLY: {Name: "STATUSONLY", Token: kwSTATUSONLY, Category: ContextKeyword},
		kwSTOP: {Name: "STOP", Token: kwSTOP, Category: ContextKeyword},
		kwSTOPLIST: {Name: "STOPLIST", Token: kwSTOPLIST, Category: CoreKeyword},
		kwSTREAM: {Name: "STREAM", Token: kwSTREAM, Category: ContextKeyword},
		kwSTREAMING: {Name: "STREAMING", Token: kwSTREAMING, Category: ContextKeyword},
		kwSUBSCRIPTION: {Name: "SUBSCRIPTION", Token: kwSUBSCRIPTION, Category: ContextKeyword},
		kwSUSPEND_FOR_SNAPSHOT_BACKUP: {Name: "SUSPEND_FOR_SNAPSHOT_BACKUP", Token: kwSUSPEND_FOR_SNAPSHOT_BACKUP, Category: ContextKeyword},
		kwSWITCH: {Name: "SWITCH", Token: kwSWITCH, Category: ContextKeyword},
		kwSYMMETRIC: {Name: "SYMMETRIC", Token: kwSYMMETRIC, Category: ContextKeyword},
		kwSYNONYM: {Name: "SYNONYM", Token: kwSYNONYM, Category: ContextKeyword},
		kwSYSTEM: {Name: "SYSTEM", Token: kwSYSTEM, Category: ContextKeyword},
		kwSYSTEM_TIME: {Name: "SYSTEM_TIME", Token: kwSYSTEM_TIME, Category: ContextKeyword},
		kwSYSTEM_USER: {Name: "SYSTEM_USER", Token: kwSYSTEM_USER, Category: CoreKeyword},
		kwTABLE: {Name: "TABLE", Token: kwTABLE, Category: CoreKeyword},
		kwTABLESAMPLE: {Name: "TABLESAMPLE", Token: kwTABLESAMPLE, Category: CoreKeyword},
		kwTARGET: {Name: "TARGET", Token: kwTARGET, Category: ContextKeyword},
		kwTB: {Name: "TB", Token: kwTB, Category: ContextKeyword},
		kwTCP: {Name: "TCP", Token: kwTCP, Category: ContextKeyword},
		kwTEMPDB_METADATA: {Name: "TEMPDB_METADATA", Token: kwTEMPDB_METADATA, Category: ContextKeyword},
		kwTEXTIMAGE_ON: {Name: "TEXTIMAGE_ON", Token: kwTEXTIMAGE_ON, Category: ContextKeyword},
		kwTEXTSIZE: {Name: "TEXTSIZE", Token: kwTEXTSIZE, Category: CoreKeyword},
		kwTHEN: {Name: "THEN", Token: kwTHEN, Category: CoreKeyword},
		kwTHROW: {Name: "THROW", Token: kwTHROW, Category: ContextKeyword},
		kwTIES: {Name: "TIES", Token: kwTIES, Category: ContextKeyword},
		kwTIME: {Name: "TIME", Token: kwTIME, Category: ContextKeyword},
		kwTIMEOUT: {Name: "TIMEOUT", Token: kwTIMEOUT, Category: ContextKeyword},
		kwTIMER: {Name: "TIMER", Token: kwTIMER, Category: ContextKeyword},
		kwTO: {Name: "TO", Token: kwTO, Category: CoreKeyword},
		kwTOP: {Name: "TOP", Token: kwTOP, Category: CoreKeyword},
		kwTRAN: {Name: "TRAN", Token: kwTRAN, Category: CoreKeyword},
		kwTRANSACTION: {Name: "TRANSACTION", Token: kwTRANSACTION, Category: CoreKeyword},
		kwTRANSFER: {Name: "TRANSFER", Token: kwTRANSFER, Category: ContextKeyword},
		kwTRIGGER: {Name: "TRIGGER", Token: kwTRIGGER, Category: CoreKeyword},
		kwTRUNCATE: {Name: "TRUNCATE", Token: kwTRUNCATE, Category: CoreKeyword},
		kwTRY: {Name: "TRY", Token: kwTRY, Category: ContextKeyword},
		kwTRY_CAST: {Name: "TRY_CAST", Token: kwTRY_CAST, Category: ContextKeyword},
		kwTRY_CONVERT: {Name: "TRY_CONVERT", Token: kwTRY_CONVERT, Category: CoreKeyword},
		kwTSEQUAL: {Name: "TSEQUAL", Token: kwTSEQUAL, Category: CoreKeyword},
		kwTYPE: {Name: "TYPE", Token: kwTYPE, Category: ContextKeyword},
		kwTYPE_WARNING: {Name: "TYPE_WARNING", Token: kwTYPE_WARNING, Category: ContextKeyword},
		kwUNBOUNDED: {Name: "UNBOUNDED", Token: kwUNBOUNDED, Category: ContextKeyword},
		kwUNCOMMITTED: {Name: "UNCOMMITTED", Token: kwUNCOMMITTED, Category: ContextKeyword},
		kwUNDEFINED: {Name: "UNDEFINED", Token: kwUNDEFINED, Category: ContextKeyword},
		kwUNION: {Name: "UNION", Token: kwUNION, Category: CoreKeyword},
		kwUNIQUE: {Name: "UNIQUE", Token: kwUNIQUE, Category: CoreKeyword},
		kwUNKNOWN: {Name: "UNKNOWN", Token: kwUNKNOWN, Category: ContextKeyword},
		kwUNLIMITED: {Name: "UNLIMITED", Token: kwUNLIMITED, Category: ContextKeyword},
		kwUNLOCK: {Name: "UNLOCK", Token: kwUNLOCK, Category: ContextKeyword},
		kwUNPIVOT: {Name: "UNPIVOT", Token: kwUNPIVOT, Category: CoreKeyword},
		kwUPDATE: {Name: "UPDATE", Token: kwUPDATE, Category: CoreKeyword},
		kwUPDATETEXT: {Name: "UPDATETEXT", Token: kwUPDATETEXT, Category: CoreKeyword},
		kwURL: {Name: "URL", Token: kwURL, Category: ContextKeyword},
		kwUSE: {Name: "USE", Token: kwUSE, Category: CoreKeyword},
		kwUSED: {Name: "USED", Token: kwUSED, Category: ContextKeyword},
		kwUSER: {Name: "USER", Token: kwUSER, Category: CoreKeyword},
		kwUSING: {Name: "USING", Token: kwUSING, Category: ContextKeyword},
		kwVALIDATION: {Name: "VALIDATION", Token: kwVALIDATION, Category: ContextKeyword},
		kwVALUE: {Name: "VALUE", Token: kwVALUE, Category: ContextKeyword},
		kwVALUES: {Name: "VALUES", Token: kwVALUES, Category: CoreKeyword},
		kwVARYING: {Name: "VARYING", Token: kwVARYING, Category: CoreKeyword},
		kwVECTOR: {Name: "VECTOR", Token: kwVECTOR, Category: ContextKeyword},
		kwVIEW: {Name: "VIEW", Token: kwVIEW, Category: CoreKeyword},
		kwVIEWS: {Name: "VIEWS", Token: kwVIEWS, Category: ContextKeyword},
		kwWAITFOR: {Name: "WAITFOR", Token: kwWAITFOR, Category: CoreKeyword},
		kwWHEN: {Name: "WHEN", Token: kwWHEN, Category: CoreKeyword},
		kwWHERE: {Name: "WHERE", Token: kwWHERE, Category: CoreKeyword},
		kwWHILE: {Name: "WHILE", Token: kwWHILE, Category: CoreKeyword},
		kwWINDOW: {Name: "WINDOW", Token: kwWINDOW, Category: ContextKeyword},
		kwWINDOWS: {Name: "WINDOWS", Token: kwWINDOWS, Category: ContextKeyword},
		kwWITH: {Name: "WITH", Token: kwWITH, Category: CoreKeyword},
		kwWITHIN: {Name: "WITHIN", Token: kwWITHIN, Category: ContextKeyword},
		kwWITHOUT: {Name: "WITHOUT", Token: kwWITHOUT, Category: ContextKeyword},
		kwWITHOUT_ARRAY_WRAPPER: {Name: "WITHOUT_ARRAY_WRAPPER", Token: kwWITHOUT_ARRAY_WRAPPER, Category: ContextKeyword},
		kwWORK: {Name: "WORK", Token: kwWORK, Category: ContextKeyword},
		kwWORKLOAD: {Name: "WORKLOAD", Token: kwWORKLOAD, Category: ContextKeyword},
		kwWRITE: {Name: "WRITE", Token: kwWRITE, Category: ContextKeyword},
		kwWRITETEXT: {Name: "WRITETEXT", Token: kwWRITETEXT, Category: CoreKeyword},
		kwXACT_ABORT: {Name: "XACT_ABORT", Token: kwXACT_ABORT, Category: ContextKeyword},
		kwXML: {Name: "XML", Token: kwXML, Category: ContextKeyword},
		kwXMLDATA: {Name: "XMLDATA", Token: kwXMLDATA, Category: ContextKeyword},
		kwXMLNAMESPACES: {Name: "XMLNAMESPACES", Token: kwXMLNAMESPACES, Category: ContextKeyword},
		kwXMLSCHEMA: {Name: "XMLSCHEMA", Token: kwXMLSCHEMA, Category: ContextKeyword},
		kwXSINIL: {Name: "XSINIL", Token: kwXSINIL, Category: ContextKeyword},
		kwZONE: {Name: "ZONE", Token: kwZONE, Category: ContextKeyword},
}

// lookupKeywordCategory returns the classification category for a keyword token.
func lookupKeywordCategory(token int) KeywordCategory {
	if kw, ok := keywordClassification[token]; ok {
		return kw.Category
	}
	return ContextKeyword // default to context for unknown tokens
}

// isContextKeyword returns true if the token is a context-sensitive keyword.
func isContextKeyword(token int) bool {
	if kw, ok := keywordClassification[token]; ok {
		return kw.Category == ContextKeyword
	}
	return false
}

func lookupKeyword(ident string) int {
	if tok, ok := keywordMap[strings.ToLower(ident)]; ok {
		return tok
	}
	return tokIDENT
}

// Token represents a lexical token.
type Token struct {
	Type int    // Token type
	Str  string // String value for identifiers, string literals, operators
	Ival int64  // Integer value for tokICONST
	Loc  int    // Byte offset in the source text
	End  int    // Exclusive byte offset past token end
}

// Lexer implements a T-SQL lexer.
type Lexer struct {
	input string
	pos   int
	start int
	Err   error
}

// NewLexer creates a new T-SQL lexer.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	tok := l.nextTokenInner()
	tok.End = l.pos
	return tok
}

// nextTokenInner performs the actual lexing. NextToken wraps this to set tok.End.
func (l *Lexer) nextTokenInner() Token {
	l.skipWhitespace()
	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: l.pos}
	}

	l.start = l.pos
	ch := l.input[l.pos]

	// Line comment: --
	if ch == '-' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '-' {
		l.skipLineComment()
		return l.nextTokenInner()
	}

	// Block comment: /* ... */
	if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
		l.skipBlockComment()
		return l.nextTokenInner()
	}

	// N'...' nvarchar string literal
	if (ch == 'N' || ch == 'n') && l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
		l.pos++ // skip N
		return l.lexNString()
	}

	// '...' string literal
	if ch == '\'' {
		return l.lexString()
	}

	// [bracketed identifier]
	if ch == '[' {
		return l.lexBracketedIdent()
	}

	// "quoted identifier"
	if ch == '"' {
		return l.lexQuotedIdent()
	}

	// @variable or @@sysvariable
	if ch == '@' {
		return l.lexVariable()
	}

	// Two-character operators (check before single-char)
	if l.pos+1 < len(l.input) {
		ch2 := l.input[l.pos : l.pos+2]
		switch ch2 {
		case "<>":
			l.pos += 2
			return Token{Type: tokNOTEQUAL, Str: "<>", Loc: l.start}
		case "!=":
			l.pos += 2
			return Token{Type: tokNOTEQUAL, Str: "!=", Loc: l.start}
		case "<=":
			l.pos += 2
			return Token{Type: tokLESSEQUAL, Str: "<=", Loc: l.start}
		case ">=":
			l.pos += 2
			return Token{Type: tokGREATEQUAL, Str: ">=", Loc: l.start}
		case "!<":
			l.pos += 2
			return Token{Type: tokNOTLESS, Str: "!<", Loc: l.start}
		case "!>":
			l.pos += 2
			return Token{Type: tokNOTGREATER, Str: "!>", Loc: l.start}
		case "::":
			l.pos += 2
			return Token{Type: tokCOLONCOLON, Str: "::", Loc: l.start}
		case "+=":
			l.pos += 2
			return Token{Type: tokPLUSEQUAL, Str: "+=", Loc: l.start}
		case "-=":
			l.pos += 2
			return Token{Type: tokMINUSEQUAL, Str: "-=", Loc: l.start}
		case "*=":
			l.pos += 2
			return Token{Type: tokMULEQUAL, Str: "*=", Loc: l.start}
		case "/=":
			l.pos += 2
			return Token{Type: tokDIVEQUAL, Str: "/=", Loc: l.start}
		case "%=":
			l.pos += 2
			return Token{Type: tokMODEQUAL, Str: "%=", Loc: l.start}
		case "&=":
			l.pos += 2
			return Token{Type: tokANDEQUAL, Str: "&=", Loc: l.start}
		case "|=":
			l.pos += 2
			return Token{Type: tokOREQUAL, Str: "|=", Loc: l.start}
		case "^=":
			l.pos += 2
			return Token{Type: tokXOREQUAL, Str: "^=", Loc: l.start}
		}
	}

	// Numbers (digit or .digit)
	if isDigit(ch) || (ch == '.' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1])) {
		return l.lexNumber()
	}

	// Single-character tokens
	if isSingleChar(ch) {
		l.pos++
		return Token{Type: int(ch), Str: string(ch), Loc: l.start}
	}

	// Identifiers and keywords
	if isIdentStart(ch) {
		return l.lexIdent()
	}

	// Unknown character - return as itself
	l.pos++
	return Token{Type: int(ch), Str: string(ch), Loc: l.start}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == '\f' {
			l.pos++
		} else {
			break
		}
	}
}

func (l *Lexer) skipLineComment() {
	l.pos += 2
	for l.pos < len(l.input) {
		if l.input[l.pos] == '\n' {
			l.pos++
			return
		}
		l.pos++
	}
}

func (l *Lexer) skipBlockComment() {
	l.pos += 2
	depth := 1
	for l.pos < len(l.input) && depth > 0 {
		if l.input[l.pos] == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
			depth++
			l.pos += 2
		} else if l.input[l.pos] == '*' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
			depth--
			l.pos += 2
		} else {
			l.pos++
		}
	}
	if depth > 0 {
		l.Err = fmt.Errorf("unterminated block comment")
	}
}

func (l *Lexer) lexString() Token {
	l.pos++ // skip opening quote
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\'' {
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
				buf.WriteByte('\'')
				l.pos += 2
				continue
			}
			l.pos++
			return Token{Type: tokSCONST, Str: buf.String(), Loc: l.start}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	l.Err = fmt.Errorf("unterminated string literal")
	return Token{Type: tokEOF, Loc: l.start}
}

func (l *Lexer) lexNString() Token {
	tok := l.lexString()
	if tok.Type == tokSCONST {
		tok.Type = tokNSCONST
	}
	return tok
}

func (l *Lexer) lexBracketedIdent() Token {
	l.pos++ // skip [
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ']' {
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == ']' {
				buf.WriteByte(']')
				l.pos += 2
				continue
			}
			l.pos++
			return Token{Type: tokIDENT, Str: buf.String(), Loc: l.start}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	l.Err = fmt.Errorf("unterminated bracketed identifier")
	return Token{Type: tokIDENT, Str: buf.String(), Loc: l.start}
}

func (l *Lexer) lexQuotedIdent() Token {
	l.pos++ // skip "
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' {
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '"' {
				buf.WriteByte('"')
				l.pos += 2
				continue
			}
			l.pos++
			str := buf.String()
			if str == "" {
				l.Err = fmt.Errorf("zero-length delimited identifier")
				return Token{Type: tokEOF, Loc: l.start}
			}
			return Token{Type: tokIDENT, Str: str, Loc: l.start}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	l.Err = fmt.Errorf("unterminated quoted identifier")
	return Token{Type: tokIDENT, Str: buf.String(), Loc: l.start}
}

func (l *Lexer) lexVariable() Token {
	l.pos++ // skip first @
	if l.pos < len(l.input) && l.input[l.pos] == '@' {
		// @@sysvariable
		l.pos++
		start := l.pos
		for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
			l.pos++
		}
		name := "@@" + l.input[start:l.pos]
		return Token{Type: tokSYSVARIABLE, Str: name, Loc: l.start}
	}
	// @variable
	start := l.pos
	for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
		l.pos++
	}
	name := "@" + l.input[start:l.pos]
	return Token{Type: tokVARIABLE, Str: name, Loc: l.start}
}

func (l *Lexer) lexNumber() Token {
	start := l.pos
	isFloat := false

	// Handle hex: 0x...
	if l.input[l.pos] == '0' && l.pos+1 < len(l.input) && (l.input[l.pos+1] == 'x' || l.input[l.pos+1] == 'X') {
		l.pos += 2
		for l.pos < len(l.input) && isHexDigit(l.input[l.pos]) {
			l.pos++
		}
		numStr := l.input[start:l.pos]
		val, err := strconv.ParseInt(numStr[2:], 16, 64)
		if err != nil {
			return Token{Type: tokFCONST, Str: numStr, Loc: l.start}
		}
		return Token{Type: tokICONST, Ival: val, Str: numStr, Loc: l.start}
	}

	// Leading dot: .123
	if l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
	}

	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}

	// Decimal point
	if !isFloat && l.pos < len(l.input) && l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
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
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
	}

	numStr := l.input[start:l.pos]
	if isFloat {
		return Token{Type: tokFCONST, Str: numStr, Loc: l.start}
	}

	val, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return Token{Type: tokFCONST, Str: numStr, Loc: l.start}
	}
	return Token{Type: tokICONST, Ival: val, Str: numStr, Loc: l.start}
}

func (l *Lexer) lexIdent() Token {
	start := l.pos
	for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
		l.pos++
	}
	ident := l.input[start:l.pos]

	// Check for multi-word keywords like TRY_CAST, TRY_CONVERT, etc.
	tok := lookupKeyword(ident)
	if tok != tokIDENT {
		return Token{Type: tok, Str: ident, Loc: l.start}
	}

	return Token{Type: tokIDENT, Str: ident, Loc: l.start}
}

// Character classification functions

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isIdentStart(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_' || ch == '#' || ch >= 128
}

func isIdentCont(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9') || ch == '$'
}

func isSingleChar(ch byte) bool {
	return ch == '(' || ch == ')' || ch == ',' || ch == ';' || ch == '.' ||
		ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '%' ||
		ch == '=' || ch == '<' || ch == '>' || ch == '~' || ch == '&' ||
		ch == '|' || ch == '^' || ch == ':' || ch == '!'
}
