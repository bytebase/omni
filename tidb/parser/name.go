package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// keywordCategory classifies MySQL 8.0 keywords into 6 categories matching
// the sql_yacc.yy grammar's identifier context rules.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/keywords.html
// Ref: mysql-server sql/sql_yacc.yy (ident, label_ident, role_ident, lvalue_ident rules)
type keywordCategory int

const (
	kwCatReserved    keywordCategory = iota // Cannot be used as identifiers without quoting
	kwCatUnambiguous                        // ident_keywords_unambiguous — allowed in all identifier contexts
	kwCatAmbiguous1                         // ident_keywords_ambiguous_1_roles_and_labels — NOT allowed as label or role
	kwCatAmbiguous2                         // ident_keywords_ambiguous_2_labels — NOT allowed as label
	kwCatAmbiguous3                         // ident_keywords_ambiguous_3_roles — NOT allowed as role
	kwCatAmbiguous4                         // ident_keywords_ambiguous_4_system_variables — NOT allowed as lvalue
)

// keywordCategories maps keyword token types to their category. Keywords not
// present in this map are not registered keywords (they lex as tokIDENT).
//
// All 65 original reserved keywords are migrated here with kwCatReserved.
// All other registered keywords default to kwCatUnambiguous and will be
// refined to their correct ambiguous category in later sections.
var keywordCategories = map[int]keywordCategory{
	kwSELECT:             kwCatReserved,
	kwINSERT:             kwCatReserved,
	kwUPDATE:             kwCatReserved,
	kwDELETE:             kwCatReserved,
	kwFROM:               kwCatReserved,
	kwWHERE:              kwCatReserved,
	kwCREATE:             kwCatReserved,
	kwDROP:               kwCatReserved,
	kwALTER:              kwCatReserved,
	kwTABLE:              kwCatReserved,
	kwINTO:               kwCatReserved,
	kwVALUES:             kwCatReserved,
	kwSET:                kwCatReserved,
	kwJOIN:               kwCatReserved,
	kwLEFT:               kwCatReserved,
	kwRIGHT:              kwCatReserved,
	kwINNER:              kwCatReserved,
	kwOUTER:              kwCatReserved,
	kwON:                 kwCatReserved,
	kwAND:                kwCatReserved,
	kwOR:                 kwCatReserved,
	kwNOT:                kwCatReserved,
	kwNULL:               kwCatReserved,
	kwTRUE:               kwCatReserved,
	kwFALSE:              kwCatReserved,
	kwIN:                 kwCatReserved,
	kwBETWEEN:            kwCatReserved,
	kwLIKE:               kwCatReserved,
	kwORDER:              kwCatReserved,
	kwGROUP:              kwCatReserved,
	kwBY:                 kwCatReserved,
	kwHAVING:             kwCatReserved,
	kwLIMIT:              kwCatReserved,
	kwAS:                 kwCatReserved,
	kwIS:                 kwCatReserved,
	kwEXISTS_KW:          kwCatReserved,
	kwCASE:               kwCatReserved,
	kwWHEN:               kwCatReserved,
	kwTHEN:               kwCatReserved,
	kwELSE:               kwCatReserved,
	kwEND:                kwCatAmbiguous2, // demoted from reserved — MySQL 8.0 ambiguous_2 (not label)
	kwIF:                 kwCatReserved,
	kwFOR:                kwCatReserved,
	kwWHILE:              kwCatReserved,
	kwINDEX:              kwCatReserved,
	kwKEY:                kwCatReserved,
	kwPRIMARY:            kwCatReserved,
	kwFOREIGN:            kwCatReserved,
	kwREFERENCES:         kwCatReserved,
	kwCONSTRAINT:         kwCatReserved,
	kwUNIQUE:             kwCatReserved,
	kwCHECK:              kwCatReserved,
	kwDEFAULT:            kwCatReserved,
	kwCOLUMN:             kwCatReserved,
	kwADD:                kwCatReserved,
	kwCHANGE:             kwCatReserved,
	kwRENAME:             kwCatReserved,
	kwGRANT:              kwCatReserved,
	kwREVOKE:             kwCatReserved,
	kwALL:                kwCatReserved,
	kwDISTINCT:           kwCatReserved,
	kwUNION:              kwCatReserved,
	kwINTERSECT:          kwCatReserved,
	kwEXCEPT:             kwCatReserved,
	kwACCESSIBLE:         kwCatUnambiguous,
	kwASENSITIVE:         kwCatUnambiguous,
	kwCUBE:               kwCatUnambiguous,
	kwCUME_DIST:          kwCatReserved,
	kwDENSE_RANK:         kwCatReserved,
	kwDUAL:               kwCatReserved,
	kwFIRST_VALUE:        kwCatReserved,
	kwGROUPING:           kwCatUnambiguous,
	kwINSENSITIVE:        kwCatUnambiguous,
	kwLAG:                kwCatReserved,
	kwLAST_VALUE:         kwCatReserved,
	kwLEAD:               kwCatReserved,
	kwNTH_VALUE:          kwCatReserved,
	kwNTILE:              kwCatReserved,
	kwOF:                 kwCatReserved,
	kwOPTIMIZER_COSTS:    kwCatUnambiguous,
	kwPERCENT_RANK:       kwCatReserved,
	kwRANK:               kwCatReserved,
	kwROW_NUMBER:         kwCatReserved,
	kwSENSITIVE:          kwCatUnambiguous,
	kwSPECIFIC:           kwCatUnambiguous,
	kwUSAGE:              kwCatReserved,
	kwVARYING:            kwCatReserved,
	kwDAY_HOUR:           kwCatReserved,
	kwDAY_MICROSECOND:    kwCatReserved,
	kwDAY_MINUTE:         kwCatReserved,
	kwDAY_SECOND:         kwCatReserved,
	kwHOUR_MICROSECOND:   kwCatReserved,
	kwHOUR_MINUTE:        kwCatReserved,
	kwHOUR_SECOND:        kwCatReserved,
	kwMINUTE_MICROSECOND: kwCatReserved,
	kwMINUTE_SECOND:      kwCatReserved,
	kwSECOND_MICROSECOND: kwCatReserved,
	kwYEAR_MONTH:         kwCatReserved,
	kwUTC_DATE:           kwCatReserved,
	kwUTC_TIME:           kwCatReserved,
	kwUTC_TIMESTAMP:      kwCatReserved,
	kwMAXVALUE:           kwCatReserved,
	kwNO_WRITE_TO_BINLOG: kwCatReserved,
	kwIO_AFTER_GTIDS:     kwCatUnambiguous,
	kwIO_BEFORE_GTIDS:    kwCatUnambiguous,
	kwSQLEXCEPTION:       kwCatReserved,
	kwSQLSTATE:           kwCatReserved,
	kwSQLWARNING:         kwCatReserved,
	kwCROSS:              kwCatReserved,
	kwNATURAL:            kwCatReserved,
	kwUSING:              kwCatReserved,
	kwASC:                kwCatReserved,
	kwDESC:               kwCatReserved,
	kwTO:                 kwCatReserved,
	kwDIV:                kwCatReserved,
	kwMOD:                kwCatReserved,
	kwXOR:                kwCatReserved,
	kwREGEXP:             kwCatReserved,
	kwBINARY:             kwCatReserved,
	kwINTERVAL:           kwCatReserved,
	kwMATCH:              kwCatReserved,
	kwCURRENT_DATE:       kwCatReserved,
	kwCURRENT_TIME:       kwCatReserved,
	kwCURRENT_TIMESTAMP:  kwCatReserved,
	kwCURRENT_USER:       kwCatReserved,
	kwDATABASE:           kwCatReserved,
	kwFUNCTION:           kwCatUnambiguous,
	kwPROCEDURE:          kwCatReserved,
	kwTRIGGER:            kwCatReserved,
	kwPARTITION:          kwCatReserved,
	kwRANGE:              kwCatReserved,
	kwROW:                kwCatReserved,
	kwROWS:               kwCatReserved,
	kwOVER:               kwCatReserved,
	kwWINDOW:             kwCatReserved,
	kwFORCE:              kwCatReserved,
	kwCONVERT:            kwCatReserved,
	kwCAST:               kwCatUnambiguous,
	kwWITH:               kwCatReserved,
	kwREPLACE:            kwCatReserved,
	kwIGNORE:             kwCatReserved,
	kwLOAD:               kwCatReserved,
	kwUSE:                kwCatReserved,
	kwKILL:               kwCatReserved,
	kwEXPLAIN:            kwCatReserved,
	kwSPATIAL:            kwCatReserved,
	kwFULLTEXT:           kwCatReserved,
	kwOUTFILE:            kwCatReserved,
	kwGEOMETRY:           kwCatUnambiguous,
	kwPOINT:              kwCatUnambiguous,
	kwLINESTRING:         kwCatUnambiguous,
	kwPOLYGON:            kwCatUnambiguous,
	kwMULTIPOINT:         kwCatUnambiguous,
	kwMULTILINESTRING:    kwCatUnambiguous,
	kwMULTIPOLYGON:       kwCatUnambiguous,
	kwGEOMETRYCOLLECTION: kwCatUnambiguous,
	kwSERIAL:             kwCatUnambiguous,
	kwNATIONAL:           kwCatUnambiguous,
	kwNCHAR:              kwCatUnambiguous,
	kwNVARCHAR:           kwCatUnambiguous,
	kwSIGNED:             kwCatAmbiguous2,
	kwPRECISION:          kwCatReserved,
	kwBOOL:               kwCatUnambiguous,
	kwBOOLEAN:            kwCatUnambiguous,
	kwSRID:               kwCatUnambiguous,
	kwENFORCED:           kwCatUnambiguous,
	kwLESS:               kwCatUnambiguous,
	kwTHAN:               kwCatUnambiguous,
	kwSUBPARTITIONS:      kwCatUnambiguous,
	kwLEAVES:             kwCatUnambiguous,
	kwPARSER:             kwCatUnambiguous,
	kwCOMPRESSION:        kwCatUnambiguous,
	kwINSERT_METHOD:      kwCatUnambiguous,
	kwACTION:             kwCatUnambiguous,
	kwPARTIAL:            kwCatUnambiguous,
	kwFORMAT:             kwCatUnambiguous,
	kwXML:                kwCatUnambiguous,
	kwCONCURRENT:         kwCatUnambiguous,
	kwWORK:               kwCatUnambiguous,
	kwXID:                kwCatUnambiguous,
	kwEXPORT:             kwCatUnambiguous,
	kwUPGRADE:            kwCatUnambiguous,
	kwFAST:               kwCatUnambiguous,
	kwMEDIUM:             kwCatUnambiguous,
	kwCHANGED:            kwCatUnambiguous,
	kwCODE:               kwCatUnambiguous,
	kwEVENTS:             kwCatUnambiguous,
	kwINDEXES:            kwCatUnambiguous,
	kwGRANTS:             kwCatUnambiguous,
	kwTRIGGERS:           kwCatUnambiguous,
	kwSCHEMAS:            kwCatReserved,
	kwPARTITIONS:         kwCatUnambiguous,
	kwHOSTS:              kwCatUnambiguous,
	kwMUTEX:              kwCatUnambiguous,
	kwPROFILE:            kwCatUnambiguous,
	kwREPLICAS:           kwCatUnambiguous,
	kwNAMES:              kwCatUnambiguous,
	kwACCOUNT:            kwCatUnambiguous,
	kwOPTION:             kwCatReserved,
	kwPROXY:              kwCatAmbiguous3,
	kwROUTINE:            kwCatUnambiguous,
	kwEXPIRE:             kwCatUnambiguous,
	kwNEVER:              kwCatUnambiguous,
	kwDAY:                kwCatUnambiguous,
	kwHISTORY:            kwCatUnambiguous,
	kwREUSE:              kwCatUnambiguous,
	kwOPTIONAL:           kwCatUnambiguous,
	kwX509:               kwCatUnambiguous,
	kwISSUER:             kwCatUnambiguous,
	kwSUBJECT:            kwCatUnambiguous,
	kwCIPHER:             kwCatUnambiguous,
	kwSCHEDULE:           kwCatUnambiguous,
	kwCOMPLETION:         kwCatUnambiguous,
	kwPRESERVE:           kwCatUnambiguous,
	kwEVERY:              kwCatUnambiguous,
	kwSTARTS:             kwCatUnambiguous,
	kwENDS:               kwCatUnambiguous,
	kwVALUE:              kwCatUnambiguous,
	kwSTACKED:            kwCatUnambiguous,
	kwUNKNOWN:            kwCatUnambiguous,
	kwWAIT:               kwCatUnambiguous,
	kwACTIVE:             kwCatUnambiguous,
	kwINACTIVE:           kwCatUnambiguous,
	kwATTRIBUTE:          kwCatUnambiguous,
	kwADMIN:              kwCatUnambiguous,
	kwDESCRIPTION:        kwCatUnambiguous,
	kwORGANIZATION:       kwCatUnambiguous,
	kwREFERENCE:          kwCatUnambiguous,
	kwDEFINITION:         kwCatUnambiguous,
	kwNAME:               kwCatUnambiguous,
	kwSYSTEM:             kwCatUnambiguous,
	kwROTATE:             kwCatUnambiguous,
	kwKEYRING:            kwCatUnambiguous,
	kwTLS:                kwCatUnambiguous,
	kwSTREAM:             kwCatUnambiguous,
	kwGENERATE:           kwCatUnambiguous,
	// Section 1.5: Ambiguous category classifications
	// Ambiguous 1 (not label, not role)
	kwEXECUTE: kwCatAmbiguous1,
	// Ambiguous 2 (not label)
	kwBEGIN:     kwCatAmbiguous2,
	kwCOMMIT:    kwCatAmbiguous2,
	kwCONTAINS:  kwCatAmbiguous2,
	kwDO:        kwCatAmbiguous2,
	kwFLUSH:     kwCatAmbiguous2,
	kwFOLLOWS:   kwCatAmbiguous2,
	kwPRECEDES:  kwCatAmbiguous2,
	kwPREPARE:   kwCatAmbiguous2,
	kwREPAIR:    kwCatAmbiguous2,
	kwRESET:     kwCatAmbiguous2,
	kwROLLBACK:  kwCatAmbiguous2,
	kwSAVEPOINT: kwCatAmbiguous2,
	kwSLAVE:     kwCatAmbiguous2,
	kwSTART:     kwCatAmbiguous2,
	kwSTOP:      kwCatAmbiguous2,
	kwTRUNCATE:  kwCatAmbiguous2,
	kwXA:        kwCatAmbiguous2,
	// Ambiguous 3 (not role)
	kwEVENT:       kwCatAmbiguous3,
	kwPROCESS:     kwCatAmbiguous3,
	kwRELOAD:      kwCatAmbiguous3,
	kwREPLICATION: kwCatAmbiguous3,
	// Ambiguous 4 (not lvalue)
	kwGLOBAL:  kwCatAmbiguous4,
	kwSESSION: kwCatAmbiguous4,
	kwLOCAL:   kwCatAmbiguous4,
	// --- 253 missing keyword categories (non-unambiguous only) ---
	kwASCII:               kwCatAmbiguous2,
	kwBIT_AND:             kwCatUnambiguous,
	kwBIT_OR:              kwCatUnambiguous,
	kwBIT_XOR:             kwCatUnambiguous,
	kwBYTE:                kwCatAmbiguous2,
	kwCURDATE:             kwCatUnambiguous,
	kwCURTIME:             kwCatUnambiguous,
	kwDATE_ADD:            kwCatUnambiguous,
	kwDATE_SUB:            kwCatUnambiguous,
	kwEXTERNAL:            kwCatUnambiguous,
	kwFILE:                kwCatAmbiguous3,
	kwFLOAT4:              kwCatReserved,
	kwFLOAT8:              kwCatReserved,
	kwINT1:                kwCatReserved,
	kwINT2:                kwCatReserved,
	kwINT3:                kwCatReserved,
	kwINT4:                kwCatReserved,
	kwINT8:                kwCatReserved,
	kwJSON_ARRAYAGG:       kwCatUnambiguous,
	kwJSON_DUALITY_OBJECT: kwCatUnambiguous,
	kwJSON_OBJECTAGG:      kwCatUnambiguous,
	kwLIBRARY:             kwCatUnambiguous,
	kwLONG:                kwCatReserved,
	kwMANUAL:              kwCatUnambiguous,
	kwMID:                 kwCatUnambiguous,
	kwMIDDLEINT:           kwCatReserved,
	kwNOW:                 kwCatUnambiguous,
	kwPARALLEL:            kwCatUnambiguous,
	kwPERSIST_ONLY:        kwCatAmbiguous4,
	kwQUALIFY:             kwCatUnambiguous,
	kwREAD_WRITE:          kwCatUnambiguous,
	kwSETS:                kwCatUnambiguous,
	kwSTD:                 kwCatUnambiguous,
	kwSTDDEV:              kwCatUnambiguous,
	kwSTDDEV_POP:          kwCatUnambiguous,
	kwSTDDEV_SAMP:         kwCatUnambiguous,
	kwSUBSTR:              kwCatUnambiguous,
	kwSUPER:               kwCatAmbiguous3,
	kwSYSDATE:             kwCatUnambiguous,
	kwTABLESAMPLE:         kwCatReserved,
	kwUNICODE:             kwCatAmbiguous2,
	kwVAR_POP:             kwCatUnambiguous,
	kwVAR_SAMP:            kwCatUnambiguous,
	kwVARCHARACTER:        kwCatReserved,
	kwVARIANCE:            kwCatUnambiguous,
	// --- Section 3.1: Missing reserved classifications ---
	kwANALYZE:             kwCatReserved,
	kwBEFORE:              kwCatUnambiguous,
	kwBIGINT:              kwCatReserved,
	kwBLOB:                kwCatReserved,
	kwBOTH:                kwCatReserved,
	kwCALL:                kwCatReserved,
	kwCASCADE:             kwCatReserved,
	kwCHAR:                kwCatReserved,
	kwCHARACTER:           kwCatReserved,
	kwCOLLATE:             kwCatReserved,
	kwCONDITION:           kwCatUnambiguous,
	kwCONTINUE:            kwCatReserved,
	kwCOUNT:               kwCatUnambiguous,
	kwCURSOR:              kwCatReserved,
	kwDATABASES:           kwCatReserved,
	kwDEC:                 kwCatReserved,
	kwDECIMAL:             kwCatReserved,
	kwDECLARE:             kwCatUnambiguous,
	kwDELAYED:             kwCatReserved,
	kwDESCRIBE:            kwCatReserved,
	kwDETERMINISTIC:       kwCatUnambiguous,
	kwDISTINCTROW:         kwCatReserved,
	kwDOUBLE:              kwCatReserved,
	kwEACH:                kwCatUnambiguous,
	kwELSEIF:              kwCatReserved,
	kwEMPTY:               kwCatUnambiguous,
	kwENCLOSED:            kwCatReserved,
	kwESCAPED:             kwCatReserved,
	kwEXIT:                kwCatReserved,
	kwEXTRACT:             kwCatUnambiguous,
	kwFETCH:               kwCatReserved,
	kwFLOAT:               kwCatReserved,
	kwGENERATED:           kwCatReserved,
	kwGET:                 kwCatUnambiguous,
	kwGROUP_CONCAT:        kwCatUnambiguous,
	kwGROUPS:              kwCatReserved,
	kwHIGH_PRIORITY:       kwCatReserved,
	kwINFILE:              kwCatReserved,
	kwINOUT:               kwCatReserved,
	kwINT:                 kwCatReserved,
	kwINTEGER:             kwCatReserved,
	kwITERATE:             kwCatReserved,
	kwKEYS:                kwCatReserved,
	kwLATERAL:             kwCatUnambiguous,
	kwLEADING:             kwCatReserved,
	kwLEAVE:               kwCatReserved,
	kwLINEAR:              kwCatReserved,
	kwLINES:               kwCatReserved,
	kwLOCALTIME:           kwCatReserved,
	kwLOCALTIMESTAMP:      kwCatReserved,
	kwLOCK:                kwCatReserved,
	kwLONGBLOB:            kwCatReserved,
	kwLONGTEXT:            kwCatReserved,
	kwLOOP:                kwCatUnambiguous,
	kwLOW_PRIORITY:        kwCatReserved,
	kwMAX:                 kwCatUnambiguous,
	kwMEDIUMBLOB:          kwCatReserved,
	kwMEDIUMINT:           kwCatReserved,
	kwMEDIUMTEXT:          kwCatReserved,
	kwMIN:                 kwCatUnambiguous,
	kwMODIFIES:            kwCatUnambiguous,
	kwNUMERIC:             kwCatReserved,
	kwOPTIMIZE:            kwCatReserved,
	kwOPTIONALLY:          kwCatReserved,
	kwOUT:                 kwCatReserved,
	kwPOSITION:            kwCatUnambiguous,
	kwPURGE:               kwCatUnambiguous,
	kwREAD:                kwCatReserved,
	kwREADS:               kwCatUnambiguous,
	kwREAL:                kwCatReserved,
	kwRECURSIVE:           kwCatReserved,
	kwRELEASE:             kwCatReserved,
	kwREPEAT:              kwCatReserved,
	kwREQUIRE:             kwCatReserved,
	kwRESIGNAL:            kwCatUnambiguous,
	kwRESTRICT:            kwCatReserved,
	kwRETURN:              kwCatUnambiguous,
	kwRLIKE:               kwCatReserved,
	kwSCHEMA:              kwCatReserved,
	kwSEPARATOR:           kwCatUnambiguous,
	kwSHOW:                kwCatReserved,
	kwSIGNAL:              kwCatUnambiguous,
	kwSMALLINT:            kwCatReserved,
	kwSQL:                 kwCatReserved,
	kwSQL_BIG_RESULT:      kwCatReserved,
	kwSQL_CALC_FOUND_ROWS: kwCatReserved,
	kwSQL_SMALL_RESULT:    kwCatReserved,
	kwSSL:                 kwCatReserved,
	kwSTARTING:            kwCatReserved,
	kwSTORED:              kwCatReserved,
	kwSTRAIGHT_JOIN:       kwCatReserved,
	kwSUBSTRING:           kwCatUnambiguous,
	kwSUM:                 kwCatUnambiguous,
	kwTERMINATED:          kwCatReserved,
	kwTINYBLOB:            kwCatReserved,
	kwTINYINT:             kwCatReserved,
	kwTINYTEXT:            kwCatReserved,
	kwTRAILING:            kwCatReserved,
	kwTRIM:                kwCatUnambiguous,
	kwUNDO:                kwCatUnambiguous,
	kwUNLOCK:              kwCatReserved,
	kwUNSIGNED:            kwCatReserved,
	kwVARBINARY:           kwCatReserved,
	kwVARCHAR:             kwCatReserved,
	kwVIRTUAL:             kwCatReserved,
	kwWRITE:               kwCatReserved,
	kwZEROFILL:            kwCatReserved,
	// --- Section 3.2: Missing ambiguous classifications ---
	// Ambiguous 2 (not label)
	kwBINLOG:     kwCatAmbiguous2,
	kwCACHE:      kwCatAmbiguous2,
	kwCHARSET:    kwCatAmbiguous2,
	kwCHECKSUM:   kwCatAmbiguous2,
	kwCLONE:      kwCatAmbiguous2,
	kwCOMMENT:    kwCatAmbiguous2,
	kwDEALLOCATE: kwCatAmbiguous2,
	kwHANDLER:    kwCatAmbiguous2,
	kwHELP:       kwCatAmbiguous2,
	kwIMPORT:     kwCatAmbiguous2,
	kwINSTALL:    kwCatAmbiguous2,
	kwLANGUAGE:   kwCatAmbiguous2,
	kwNO:         kwCatAmbiguous2,
	kwUNINSTALL:  kwCatAmbiguous2,
	// Ambiguous 1 (not label, not role)
	kwRESTART:  kwCatAmbiguous1,
	kwSHUTDOWN: kwCatAmbiguous1,
	// Ambiguous 3 (not role)
	kwNONE:     kwCatAmbiguous3,
	kwRESOURCE: kwCatAmbiguous3,
	// Ambiguous 4 (not lvalue)
	kwPERSIST: kwCatAmbiguous4,

	// TiDB-specific keywords.
	kwAUTO_RANDOM:       kwCatUnambiguous,
	kwAUTO_RANDOM_BASE:  kwCatUnambiguous,
	kwSHARD_ROW_ID_BITS: kwCatUnambiguous,
	kwPRE_SPLIT_REGIONS: kwCatUnambiguous,
	kwAUTO_ID_CACHE:     kwCatUnambiguous,
	kwCLUSTERED:         kwCatUnambiguous,
	kwNONCLUSTERED:      kwCatUnambiguous,
	kwTIFLASH:           kwCatUnambiguous,
	kwTTL:               kwCatUnambiguous,
	kwTTL_ENABLE:        kwCatUnambiguous,
	kwTTL_JOB_INTERVAL:  kwCatUnambiguous,
	kwPLACEMENT:         kwCatUnambiguous,
	kwPOLICY:            kwCatUnambiguous,

	// TiDB PLACEMENT POLICY option keywords (all unreserved upstream).
	kwPRIMARY_REGION:       kwCatUnambiguous,
	kwREGIONS:              kwCatUnambiguous,
	kwFOLLOWERS:            kwCatUnambiguous,
	kwVOTERS:               kwCatUnambiguous,
	kwLEARNERS:             kwCatUnambiguous,
	kwCONSTRAINTS:          kwCatUnambiguous,
	kwLEADER_CONSTRAINTS:   kwCatUnambiguous,
	kwFOLLOWER_CONSTRAINTS: kwCatUnambiguous,
	kwVOTER_CONSTRAINTS:    kwCatUnambiguous,
	kwLEARNER_CONSTRAINTS:  kwCatUnambiguous,
	kwSURVIVAL_PREFERENCES: kwCatUnambiguous,

	// TiDB LOCATION LABELS clause tokens (parser.y:7071-7072, unreserved).
	kwLOCATION: kwCatUnambiguous,
	kwLABELS:   kwCatUnambiguous,

	// MASTER_LOG_FILE — legacy MySQL replication alias, unreserved.
	kwMASTER_LOG_FILE: kwCatUnambiguous,

	// TiDB non-transactional DML keywords (BATCH = parser.y:7355;
	// DRY/RUN = parser.y:7397-7398) — all unreserved upstream, usable as
	// identifiers/labels/roles/lvalues.
	kwBATCH: kwCatUnambiguous,
	kwDRY:   kwCatUnambiguous,
	kwRUN:   kwCatUnambiguous,

	// TiDB sequence keywords (parser.y:7077-7085) — all unreserved
	// upstream, so `CREATE SEQUENCE sequence` and a column named cycle
	// stay valid.
	kwSEQUENCE:   kwCatUnambiguous,
	kwINCREMENT:  kwCatUnambiguous,
	kwMINVALUE:   kwCatUnambiguous,
	kwNOMINVALUE: kwCatUnambiguous,
	kwNOMAXVALUE: kwCatUnambiguous,
	kwNOCACHE:    kwCatUnambiguous,
	kwCYCLE:      kwCatUnambiguous,
	kwNOCYCLE:    kwCatUnambiguous,

	// TiDB's own SHOW targets (parser.y:6962 UnReservedKeyword,
	// parser.y:7353 TiDBKeyword) — all unreserved upstream, so a table
	// called config or a column called job stays valid.
	kwBINDINGS:             kwCatUnambiguous,
	kwBINDING_CACHE:        kwCatUnambiguous,
	kwBUILTINS:             kwCatUnambiguous,
	kwBACKUPS:              kwCatUnambiguous,
	kwRESTORES:             kwCatUnambiguous,
	kwCOLUMN_STATS_USAGE:   kwCatUnambiguous,
	kwCONFIG:               kwCatUnambiguous,
	kwDISTRIBUTION:         kwCatUnambiguous,
	kwDISTRIBUTIONS:        kwCatUnambiguous,
	kwHISTOGRAMS_IN_FLIGHT: kwCatUnambiguous,
	kwJOB:                  kwCatUnambiguous,
	kwJOBS:                 kwCatUnambiguous,
	kwNEXT_ROW_ID:          kwCatUnambiguous,
	kwSESSION_STATES:       kwCatUnambiguous,
	kwSTATS_BUCKETS:        kwCatUnambiguous,
	kwSTATS_EXTENDED:       kwCatUnambiguous,
	kwSTATS_HEALTHY:        kwCatUnambiguous,
	kwSTATS_HISTOGRAMS:     kwCatUnambiguous,
	kwSTATS_LOCKED:         kwCatUnambiguous,
	kwSTATS_META:           kwCatUnambiguous,
	kwSTATS_TOPN:           kwCatUnambiguous,

	// TiDB's SQL plan binding keywords, each unreserved on v8.5.5, so each
	// stays usable as a table, column, alias, or index name.
	kwBINDING:  kwCatUnambiguous,
	kwDIGEST:   kwCatUnambiguous,
	kwPLAN:     kwCatUnambiguous,
	kwENABLED:  kwCatUnambiguous,
	kwDISABLED: kwCatUnambiguous,
	// TiDB SPLIT and DISTRIBUTE keywords. SPLIT, DISTRIBUTE, and REGION
	// are TiDBKeyword (parser.y:7390, 7366, 7395); RULE and TIMEOUT are
	// UnReservedKeyword (parser.y:7062, 7079). None is reserved, so
	// `DELETE FROM t split`, `SELECT 1 AS rule`, and a column named
	// timeout all stay valid, as they are on v8.5.5.
	kwSPLIT:      kwCatUnambiguous,
	kwDISTRIBUTE: kwCatUnambiguous,
	kwREGION:     kwCatUnambiguous,
	kwRULE:       kwCatUnambiguous,
	kwTIMEOUT:    kwCatUnambiguous,
	// TiDB's resource control keywords. Upstream splits them between
	// UnReservedKeyword (parser.y:6962) and NotKeywordToken (parser.y:7400),
	// and Identifier (parser.y:6956) takes both, so each of these stays a
	// valid identifier — a column called priority or duration, a table
	// called plan.
	kwRU_PER_SEC:        kwCatUnambiguous,
	kwPRIORITY:          kwCatUnambiguous,
	kwBURSTABLE:         kwCatUnambiguous,
	kwQUERY_LIMIT:       kwCatUnambiguous,
	kwBACKGROUND:        kwCatUnambiguous,
	kwTASK_TYPES:        kwCatUnambiguous,
	kwUTILIZATION_LIMIT: kwCatUnambiguous,
	kwEXEC_ELAPSED:      kwCatUnambiguous,
	kwPROCESSED_KEYS:    kwCatUnambiguous,
	kwRU:                kwCatUnambiguous,
	kwWATCH:             kwCatUnambiguous,
	kwDURATION:          kwCatUnambiguous,
	kwDRYRUN:            kwCatUnambiguous,
	kwCOOLDOWN:          kwCatUnambiguous,
	kwSWITCH_GROUP:      kwCatUnambiguous,
	kwEXACT:             kwCatUnambiguous,
	kwSIMILAR:           kwCatUnambiguous,
	kwUNLIMITED:         kwCatUnambiguous,
	kwLOW:               kwCatUnambiguous,
	kwHIGH:              kwCatUnambiguous,
	kwCALIBRATE:         kwCatUnambiguous,
	kwWORKLOAD:          kwCatUnambiguous,
	kwTPCC:              kwCatUnambiguous,
	kwTPCH_10:           kwCatUnambiguous,
	kwOLTP_READ_WRITE:   kwCatUnambiguous,
	kwOLTP_READ_ONLY:    kwCatUnambiguous,
	kwOLTP_WRITE_ONLY:   kwCatUnambiguous,
	kwSTART_TIME:        kwCatUnambiguous,
	kwEND_TIME:          kwCatUnambiguous,
	kwFLASHBACK:         kwCatUnambiguous,

	// Unreserved on MySQL and in ReservedKeyword on TiDB, parser.y:79, 85 and 289. Each is
	// a 1064 as a bare name on v8.5.5.
	kwARRAY: kwCatReserved,
	kwUNTIL: kwCatReserved,
}

// tidbPlainIdentifiers are words this lexer holds as keywords and TiDB v8.5.5 does not, so
// upstream's scanner hands them over as the bare identifier token. A few grammar rules take
// that token and no keyword: RoleNameString, parser.y:11135, and ExtendedPriv, parser.y:14501.
// `CREATE ROLE loop` parses there, where `CREATE ROLE transaction` is a 1064 though both are
// unreserved. TestRoleNameSweepTiDBOracle holds this list to the server, so a keyword added to
// the lexer that TiDB does not have fails there until it is listed.
var tidbPlainIdentifiers = map[string]bool{
	"absent": true, "accessible": true, "active": true, "aggregate": true,
	"allow_missing_files": true, "asensitive": true,
	"assign_gtids_to_anonymous_transactions": true, "at": true, "authentication": true,
	"auto": true, "auto_refresh": true, "auto_refresh_source": true, "autoextend_size": true,
	"before": true, "bulk": true, "catalog_name": true, "challenge_response": true,
	"changed": true, "channel": true, "class_origin": true, "clone": true, "code": true,
	"column_name": true, "completion": true, "component": true, "concat": true,
	"concurrent": true, "condition": true, "constraint_catalog": true,
	"constraint_name": true, "constraint_schema": true, "contains": true, "count": true,
	"cube": true, "cursor_name": true, "datafile": true, "default_auth": true,
	"definition": true, "description": true, "deterministic": true, "diagnostics": true,
	"duality": true, "dumpfile": true, "each": true, "empty": true, "ends": true,
	"engine_attribute": true, "every": true, "exclude": true, "export": true,
	"extent_size": true, "external": true, "external_format": true, "factor": true,
	"fast": true, "file_block_size": true, "file_format": true, "file_name": true,
	"file_pattern": true, "file_prefix": true, "files": true, "filter": true, "finish": true,
	"follows": true, "generate": true, "geomcollection": true, "geometry": true,
	"geometrycollection": true, "get": true, "get_source_public_key": true, "greatest": true,
	"group_replication": true, "grouping": true, "gtid_only": true, "gtids": true,
	"guided": true, "header": true, "host": true, "ignore_server_ids": true, "inactive": true,
	"initial": true, "initial_size": true, "initiate": true, "insensitive": true,
	"install": true, "io_after_gtids": true, "io_before_gtids": true, "io_thread": true,
	"json_duality_object": true, "json_table": true, "json_value": true, "keyring": true,
	"lateral": true, "least": true, "leaves": true, "library": true, "linestring": true,
	"locks": true, "logfile": true, "loop": true, "manual": true, "master_log_file": true,
	"materialized": true, "max_size": true, "message_text": true, "mid": true,
	"migrate": true, "modifies": true, "multilinestring": true, "multipoint": true,
	"multipolygon": true, "mutex": true, "mysql_errno": true, "name": true, "ndb": true,
	"ndbcluster": true, "nested": true, "network_namespace": true, "new": true,
	"no_wait": true, "nullif": true, "number": true, "oj": true, "old": true, "one": true,
	"optimizer_costs": true, "options": true, "ordinality": true, "organization": true,
	"others": true, "overlay": true, "owner": true, "parallel": true, "parameters": true,
	"parse_tree": true, "path": true, "persist": true, "persist_only": true, "phase": true,
	"placing": true, "plugin": true, "plugin_dir": true, "polygon": true, "port": true,
	"precedes": true, "prev": true, "privilege_checks_user": true, "qualify": true,
	"random": true, "read_only": true, "read_write": true, "reads": true,
	"redo_buffer_size": true, "reference": true, "registration": true, "relational": true,
	"relay": true, "relay_log_file": true, "relay_log_pos": true, "relay_thread": true,
	"relaylog": true, "replicate_do_db": true, "replicate_do_table": true,
	"replicate_ignore_db": true, "replicate_ignore_table": true, "replicate_rewrite_db": true,
	"replicate_wild_do_table": true, "replicate_wild_ignore_table": true,
	"require_row_format": true, "require_table_primary_key_check": true, "resignal": true,
	"retain": true, "return": true, "returned_sqlstate": true, "returning": true,
	"returns": true, "rotate": true, "schema_name": true, "secondary_engine_attribute": true,
	"sensitive": true, "server": true, "session_user": true, "sets": true, "signal": true,
	"socket": true, "soname": true, "sounds": true, "source_auto_position": true,
	"source_bind": true, "source_compression_algorithms": true, "source_connect_retry": true,
	"source_connection_auto_failover": true, "source_delay": true,
	"source_heartbeat_period": true, "source_host": true, "source_log_file": true,
	"source_log_pos": true, "source_password": true, "source_port": true,
	"source_public_key_path": true, "source_retry_count": true, "source_ssl": true,
	"source_ssl_ca": true, "source_ssl_capath": true, "source_ssl_cert": true,
	"source_ssl_cipher": true, "source_ssl_crl": true, "source_ssl_crlpath": true,
	"source_ssl_key": true, "source_ssl_verify_server_cert": true,
	"source_tls_ciphersuites": true, "source_tls_version": true, "source_user": true,
	"source_zstd_compression_level": true, "specific": true, "sql_after_gtids": true,
	"sql_after_mts_gaps": true, "sql_before_gtids": true, "sql_thread": true, "srid": true,
	"st_collect": true, "stacked": true, "starts": true, "stream": true, "strict_load": true,
	"string": true, "subclass_origin": true, "suspend": true, "sysdate": true,
	"system_user": true, "table_name": true, "thread_priority": true, "ties": true,
	"types": true, "undo": true, "undo_buffer_size": true, "undofile": true,
	"uninstall": true, "unregister": true, "upgrade": true, "uri": true, "url": true,
	"use_frm": true, "user_resources": true, "validate": true, "vcpu": true,
	"verify_key_constraints": true, "work": true, "wrapper": true, "xa": true, "xid": true,
	"xml": true, "zone": true,
}

// isPlainIdentifier reports whether the current token is what upstream's scanner would call
// the bare identifier token: an identifier here, or one of tidbPlainIdentifiers.
func (p *Parser) isPlainIdentifier() bool {
	if p.cur.Type == tokIDENT {
		return true
	}
	return p.cur.Type >= 700 && tidbPlainIdentifiers[strings.ToLower(p.cur.Str)]
}

// isReserved returns true if the token type is a reserved keyword that cannot
// be used as an unquoted identifier.
func isReserved(t int) bool {
	cat, ok := keywordCategories[t]
	return ok && cat == kwCatReserved
}

// IsReservedKeyword reports whether s is a reserved keyword, which cannot be
// an unquoted identifier.
func IsReservedKeyword(s string) bool {
	tok, ok := keywords[strings.ToLower(s)]
	return ok && isReserved(tok)
}

// isIdentKeyword returns true if the token type is a non-reserved keyword that
// can be used as an identifier. This covers all 5 non-reserved categories:
// unambiguous, ambiguous_1, ambiguous_2, ambiguous_3, and ambiguous_4.
func isIdentKeyword(t int) bool {
	cat, ok := keywordCategories[t]
	return ok && cat != kwCatReserved
}

// isLabelKeyword returns true if the token type is a non-reserved keyword that
// can be used as a statement label. Includes: unambiguous, ambiguous_3, ambiguous_4.
// Excludes: ambiguous_1 (not label, not role), ambiguous_2 (not label).
func isLabelKeyword(t int) bool {
	cat, ok := keywordCategories[t]
	if !ok {
		return false
	}
	return cat == kwCatUnambiguous || cat == kwCatAmbiguous3 || cat == kwCatAmbiguous4
}

// isRoleKeyword returns true if the token type is a non-reserved keyword that
// can be used as a role name. Includes: unambiguous, ambiguous_2, ambiguous_4.
// Excludes: ambiguous_1 (not label, not role), ambiguous_3 (not role).
func isRoleKeyword(t int) bool {
	cat, ok := keywordCategories[t]
	if !ok {
		return false
	}
	return cat == kwCatUnambiguous || cat == kwCatAmbiguous2 || cat == kwCatAmbiguous4
}

// isLvalueKeyword returns true if the token type is a non-reserved keyword that
// can be used as an lvalue (SET target). Includes: unambiguous, ambiguous_1, ambiguous_2, ambiguous_3.
// Excludes: ambiguous_4 (system variables like GLOBAL, SESSION, LOCAL).
func isLvalueKeyword(t int) bool {
	cat, ok := keywordCategories[t]
	if !ok {
		return false
	}
	return cat == kwCatUnambiguous || cat == kwCatAmbiguous1 || cat == kwCatAmbiguous2 || cat == kwCatAmbiguous3
}

// parseIdent parses an identifier matching MySQL's `ident` grammar rule.
// Accepts tokIDENT plus all 5 non-reserved keyword categories (unambiguous,
// ambiguous_1, ambiguous_2, ambiguous_3, ambiguous_4).
//
// Ref: mysql-server sql/sql_yacc.yy — ident rule
// Ref: https://dev.mysql.com/doc/refman/8.0/en/identifiers.html
func (p *Parser) parseIdent() (string, int, error) {
	if p.cur.Type == tokIDENT {
		// The lexer scans @a whole as an identifier for the expression reader, and
		// a character set introducer such as _utf8 is one word to it. Neither is a
		// name: `CREATE TABLE @a (b int)` and `CREATE TABLE _utf8 (b int)` are 1064s,
		// where backquoted `@a` and `_utf8` are ordinary names.
		switch p.lexer.inputAt(p.cur.Loc) {
		case '@':
			return "", 0, p.syntaxErrorAtCur()
		case '_':
			if isCharsetIntroducerWord(p.cur.Str) {
				return "", 0, p.syntaxErrorAtCur()
			}
		}
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	// Accept non-reserved keyword tokens as identifiers.
	if p.cur.Type >= 700 && !isReserved(p.cur.Type) {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type == tokEOF {
		return "", 0, p.syntaxErrorAtCur()
	}
	return "", 0, &ParseError{
		Message:  "expected identifier",
		Position: p.cur.Loc,
	}
}

// isBareWord reports whether the current token is word, written without
// backquotes. A grammar that names something the keyword table does not hold
// reads it this way: the token arrives as an identifier, and a backquoted
// identifier of the same spelling is not the name the grammar means.
func (p *Parser) isBareWord(word string) bool {
	if p.cur.Type != tokIDENT && p.cur.Type < 700 {
		return false
	}
	if !strings.EqualFold(p.cur.Str, word) {
		return false
	}
	return !strings.HasPrefix(p.inputText(p.cur.Loc, p.cur.End), "`")
}

// parseIdentifier is a thin alias for parseIdent, preserved for gradual migration
// of existing call sites. New code should use parseIdent directly.
func (p *Parser) parseIdentifier() (string, int, error) {
	return p.parseIdent()
}

// parseIdentOrText parses an identifier OR a string literal, matching MySQL's
// `ident_or_text` grammar rule. Used for alias positions where MySQL accepts a
// quoted string as the alias (e.g. SELECT a AS 'x', FROM t AS 'y').
func (p *Parser) parseIdentOrText() (string, int, error) {
	if p.cur.Type == tokSCONST {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	return p.parseIdent()
}

// parseLabelIdent parses an identifier matching MySQL's `label_ident` grammar rule.
// Accepts tokIDENT plus unambiguous, ambiguous_3, and ambiguous_4 keywords.
// Excludes ambiguous_1 (not label, not role) and ambiguous_2 (not label).
//
// Ref: mysql-server sql/sql_yacc.yy — label_ident rule
func (p *Parser) parseLabelIdent() (string, int, error) {
	if p.cur.Type == tokIDENT {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type >= 700 && isLabelKeyword(p.cur.Type) {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type == tokEOF {
		return "", 0, p.syntaxErrorAtCur()
	}
	return "", 0, &ParseError{
		Message:  "expected identifier",
		Position: p.cur.Loc,
	}
}

// parseRoleIdent parses an identifier matching MySQL's `role_ident` grammar rule.
// Accepts tokIDENT plus unambiguous, ambiguous_2, and ambiguous_4 keywords.
// Excludes ambiguous_1 (not label, not role) and ambiguous_3 (not role).
//
// Ref: mysql-server sql/sql_yacc.yy — role_ident rule
func (p *Parser) parseRoleIdent() (string, int, error) {
	if p.cur.Type == tokIDENT {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type >= 700 && isRoleKeyword(p.cur.Type) {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type == tokEOF {
		return "", 0, p.syntaxErrorAtCur()
	}
	return "", 0, &ParseError{
		Message:  "expected identifier",
		Position: p.cur.Loc,
	}
}

// parseLvalueIdent parses an identifier matching MySQL's `lvalue_ident` grammar rule.
// Accepts tokIDENT plus unambiguous, ambiguous_1, ambiguous_2, and ambiguous_3 keywords.
// Excludes ambiguous_4 (system variables like GLOBAL, SESSION, LOCAL).
//
// Ref: mysql-server sql/sql_yacc.yy — lvalue_ident rule
func (p *Parser) parseLvalueIdent() (string, int, error) {
	if p.cur.Type == tokIDENT {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type >= 700 && isLvalueKeyword(p.cur.Type) {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type == tokEOF {
		return "", 0, p.syntaxErrorAtCur()
	}
	return "", 0, &ParseError{
		Message:  "expected identifier",
		Position: p.cur.Loc,
	}
}

// parseKeywordOrIdent parses an identifier or ANY keyword token (including reserved words).
// Use this in contexts where the grammar expects a fixed enum value, action word, or
// option name that may collide with reserved keywords. Examples:
//   - ALGORITHM = DEFAULT/INSTANT/INPLACE/COPY
//   - LOCK = NONE/SHARED/EXCLUSIVE
//   - REQUIRE_TABLE_PRIMARY_KEY_CHECK = ON/OFF/STREAM/GENERATE
//   - ALTER INSTANCE action words: ROTATE/RELOAD/ENABLE/DISABLE
//   - INTERVAL units: DAY/HOUR/MINUTE/SECOND (and compound forms)
//   - EXTRACT units, EXPLAIN FORMAT values, etc.
//
// This matches MySQL's grammar behavior where specific productions explicitly
// list keyword tokens as valid alternatives (e.g., ON_SYM in option values).
func (p *Parser) parseKeywordOrIdent() (string, int, error) {
	if p.cur.Type == tokIDENT {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	// Accept ANY keyword token, including reserved words.
	if p.cur.Type >= 700 {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if p.cur.Type == tokEOF {
		return "", 0, p.syntaxErrorAtCur()
	}
	return "", 0, &ParseError{
		Message:  "expected identifier or keyword",
		Position: p.cur.Loc,
	}
}

// parseColumnRef parses a column reference, which may be qualified:
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/identifier-qualifiers.html
//
//	column_ref:
//	    identifier
//	    | identifier '.' identifier
//	    | identifier '.' identifier '.' identifier
//	    | identifier '.' '*'
//	    | identifier '.' identifier '.' '*'
func (p *Parser) parseColumnRef() (*nodes.ColumnRef, error) {
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	ref := &nodes.ColumnRef{
		Loc:    nodes.Loc{Start: start},
		Column: name,
	}

	// Check for dot-qualification
	if p.cur.Type == '.' {
		p.advance() // consume '.'

		// Check for table.* or schema.table.*
		if p.cur.Type == '*' {
			p.advance()
			ref.Table = name
			ref.Column = ""
			ref.Star = true
			ref.Loc.End = p.prev.End
			return ref, nil
		}

		name2, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}

		// Check for second dot: schema.table.col or schema.table.*
		if p.cur.Type == '.' {
			p.advance() // consume second '.'

			if p.cur.Type == '*' {
				p.advance()
				ref.Schema = name
				ref.Table = name2
				ref.Column = ""
				ref.Star = true
				ref.Loc.End = p.prev.End
				return ref, nil
			}

			name3, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			ref.Schema = name
			ref.Table = name2
			ref.Column = name3
			ref.Loc.End = p.prev.End
			return ref, nil
		}

		// table.col
		ref.Table = name
		ref.Column = name2
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseTableRef parses a table reference (possibly qualified with schema).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/identifier-qualifiers.html
//
//	table_ref:
//	    identifier
//	    | identifier '.' identifier
func (p *Parser) parseTableRef() (*nodes.TableRef, error) {
	start := p.pos()

	if p.matchWildcardDatabase() {
		// Completion: after "*.", offer the same candidates as after "db.".
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
		}
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		return &nodes.TableRef{
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
			Schema: wildcardDatabase,
			Name:   name,
		}, nil
	}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	ref := &nodes.TableRef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// Check for schema.table
	if p.cur.Type == '.' {
		p.advance() // consume '.'

		// Completion: after "db.", offer table_ref qualified with database.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
		}

		if err := refuseTableNameDatabase(name, start); err != nil {
			return nil, err
		}

		name2, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		ref.Schema = name
		ref.Name = name2
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

// wildcardDatabase is the database name TiDB writes as `*`, which stands
// for any database. It reaches a table name through TiDB's cross-database
// SQL plan bindings (`CREATE BINDING USING SELECT * FROM *.t1`), but TiDB
// spells it in the table-name rule itself, so every statement that names a
// table takes it: v8.5.5 accepts `SELECT * FROM *.t1`, `UPDATE *.t1 SET a =
// 1`, `CREATE TABLE *.t (a int)` and `ANALYZE TABLE *.t1`. Only one name
// may follow: `*.*` and `*.t1.a` are both a 1064.
//
// Ref: https://docs.pingcap.com/tidb/v8.5/sql-plan-management/#cross-database-binding
const wildcardDatabase = "*"

// matchWildcardDatabase consumes the `*` and the `.` of a wildcard database
// qualifier and reports whether it did. The `.` is required before the `*`
// is taken, so an expression's multiplication is never read as one.
func (p *Parser) matchWildcardDatabase() bool {
	if p.cur.Type != '*' || p.peekNext().Type != '.' {
		return false
	}
	p.advance() // consume '*'
	p.advance() // consume '.'
	return true
}

// parseTableRefWithAlias parses a table reference with an optional alias.
//
//	table_ref_alias:
//	    table_ref [AS identifier | identifier]
func (p *Parser) parseTableRefWithAlias() (*nodes.TableRef, error) {
	ref, err := p.parseTableRefUpToAlias()
	if err != nil {
		return nil, err
	}
	if err := p.parseTableRefIndexHints(ref); err != nil {
		return nil, err
	}
	return ref, nil
}

// atAsOf reports whether the next tokens are the AS OF that starts an AsOfClause.
//
// TiDB's lexer joins AS and OF into one token wherever AS is followed by the word
// OF (lexer.go:262), so AS OF never introduces an alias named of. The word next to
// a dot is an identifier there and here, which is why `t as of.a` does not count.
func (p *Parser) atAsOf() bool {
	return p.cur.Type == kwAS && p.peekNext().Type == kwOF
}

// parseAsOfClause parses a stale-read clause, with the current token on AS.
//
// Ref: AsOfClause, parser.y:6559 (tag v8.5.5)
//
//	AS OF TIMESTAMP Expression
//
// The expression is any Expression because the grammar does not narrow it: only
// the planner refuses `AS OF TIMESTAMP a`, with error 1054.
func (p *Parser) parseAsOfClause() (*nodes.AsOfClause, error) {
	clause := &nodes.AsOfClause{Loc: nodes.Loc{Start: p.pos()}}
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwOF); err != nil {
		return nil, err
	}

	// Completion: after AS OF, the one word the clause goes on with.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwTIMESTAMP)
		return nil, &ParseError{Message: "collecting"}
	}

	if _, err := p.expect(kwTIMESTAMP); err != nil {
		return nil, err
	}
	ts, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	clause.Timestamp = ts
	clause.Loc.End = p.prev.End
	return clause, nil
}

// parseTableRefUpToAlias parses a table name with its optional PARTITION list and
// alias, and stops before an AS OF, which is not an alias.
func (p *Parser) parseTableRefUpToAlias() (*nodes.TableRef, error) {
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}

	// Optional PARTITION (p0, p1, ...)
	if p.cur.Type == kwPARTITION {
		p.advance()
		parts, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		ref.Partitions = parts
		ref.Loc.End = p.prev.End
	}

	// Optional AS alias. Unlike MySQL, TiDB does NOT accept a string literal as
	// a table alias (ident_or_text applies only to SELECT-list aliases), so this
	// uses parseIdentifier, not parseIdentOrText.
	if p.atAsOf() {
		return ref, nil
	}
	if _, ok := p.match(kwAS); ok {
		alias, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		ref.Alias = alias
		ref.Loc.End = p.prev.End
	} else if p.isIdentToken() {
		// Alias without AS keyword — accepts non-reserved keywords,
		// matching MySQL's grammar: opt_table_alias: [AS] ident
		alias, _, _ := p.parseIdentifier()
		ref.Alias = alias
		ref.Loc.End = p.prev.End
	}

	return ref, nil
}

// parseTableRefIndexHints reads the optional index hints of a table reference:
// USE/FORCE/IGNORE {INDEX|KEY} ...
func (p *Parser) parseTableRefIndexHints(ref *nodes.TableRef) error {
	if p.cur.Type == kwUSE || p.cur.Type == kwFORCE || p.cur.Type == kwIGNORE {
		hints, err := p.parseIndexHints()
		if err != nil {
			return err
		}
		ref.IndexHints = hints
		ref.Loc.End = p.prev.End
	}
	return nil
}

// parseVariableRef parses a user variable (@var) or system variable (@@var).
// The lexer emits user variables as tokIDENT with "@" prefix,
// and system variables as tokAt2 followed by identifier/keyword tokens.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/user-variables.html
// Ref: https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html
//
//	variable_ref:
//	    '@' identifier
//	    | '@@' [GLOBAL | SESSION | LOCAL] '.' identifier ['.' identifier ...]
//	    | '@@' identifier ['.' identifier ...]
func (p *Parser) parseVariableRef() (*nodes.VariableRef, error) {
	// System variable: @@...
	if p.cur.Type == tokAt2 {
		start := p.cur.Loc
		p.advance() // consume @@

		ref := &nodes.VariableRef{
			Loc:    nodes.Loc{Start: start},
			System: true,
		}

		// Check for scope prefix: @@GLOBAL.var, @@SESSION.var, @@LOCAL.var, @@PERSIST.var
		isScopeKeyword := p.cur.Type == kwGLOBAL || p.cur.Type == kwSESSION ||
			p.cur.Type == kwLOCAL || p.cur.Type == kwPERSIST ||
			p.cur.Type == kwPERSIST_ONLY
		if isScopeKeyword && p.cur.Loc == p.prev.End {
			// Look ahead for '.' to distinguish scope from plain variable name
			next := p.peekNext()
			if next.Type == '.' && next.Loc == p.cur.End {
				switch {
				case p.cur.Type == kwGLOBAL:
					ref.Scope = "GLOBAL"
				case p.cur.Type == kwSESSION:
					ref.Scope = "SESSION"
				case p.cur.Type == kwLOCAL:
					ref.Scope = "LOCAL"
				case p.cur.Type == kwPERSIST:
					ref.Scope = "PERSIST"
				default:
					ref.Scope = "PERSIST_ONLY"
				}
				p.advance() // consume scope keyword
				p.advance() // consume '.'
				// Parse the variable name (may be dot-separated: comp.var)
				name, _, err := p.parseSysVarName()
				if err != nil {
					return nil, err
				}
				ref.Name = name
				ref.Loc.End = p.prev.End
				return ref, nil
			}
		}

		// No scope — plain @@variable_name (may be dot-separated)
		name, _, err := p.parseSysVarName()
		if err != nil {
			return nil, err
		}
		ref.Name = name
		ref.Loc.End = p.prev.End
		return ref, nil
	}

	// User variable: @name (lexer emits as tokIDENT with "@" prefix)
	if p.cur.Type == tokIDENT && len(p.cur.Str) > 1 && p.cur.Str[0] == '@' {
		tok := p.cur
		p.advance()
		ref := &nodes.VariableRef{
			Loc:  nodes.Loc{Start: tok.Loc},
			Name: tok.Str[1:],
		}
		ref.Loc.End = p.prev.End
		return ref, nil
	}

	return nil, &ParseError{
		Message:  "expected variable reference",
		Position: p.cur.Loc,
	}
}

// parseSysVarName parses a system variable name, which may contain dots
// (e.g., "validate_password.length" or "comp.var1").
//
// The first part may be a reserved word when it touches what comes before it:
// TiDB's lexer scans the name as part of the @@ token and never looks it up
// (startWithAt, lexer.go:618), so `SELECT @@table` parses and fails later with
// 1193, while `SELECT @@ table` is a 1064. Later parts are always identifiers,
// because the word follows a dot.
func (p *Parser) parseSysVarName() (string, int, error) {
	// Upstream's scanner reads @@ and the name with its dots as one token, so a space
	// anywhere in it ends the token: `SELECT @@ a` and `SELECT @@a . b` are 1064s.
	if p.cur.Loc != p.prev.End {
		return "", 0, p.syntaxErrorAtCur()
	}
	quoted := strings.HasPrefix(p.inputText(p.cur.Loc, p.cur.End), "`")
	name, loc, err := p.parseKeywordOrIdent()
	if err != nil {
		return "", 0, err
	}
	// A backquoted name is the whole token, so @@ then a backquoted a then .b is a 1064.
	if quoted {
		return name, loc, nil
	}
	// Consume additional .ident parts for component-style variable names
	for p.cur.Type == '.' && p.cur.Loc == p.prev.End {
		dot := p.advance() // consume '.'
		// The token takes a dot with nothing after it, so `SELECT @@a. b` reads an alias.
		if p.cur.Loc != dot.End || !(p.cur.Type == tokIDENT || p.cur.Type >= 700) {
			name += "."
			break
		}
		part, _, err := p.parseIdent()
		if err != nil {
			return "", 0, err
		}
		name += "." + part
	}
	return name, loc, nil
}

// isIdentToken returns true if the current token can be used as an identifier.
func (p *Parser) isIdentToken() bool {
	return p.cur.Type == tokIDENT || (p.cur.Type >= 700 && !isReserved(p.cur.Type))
}

// isLabelIdentToken returns true if the current token can be used as a label identifier.
// This matches MySQL's label_ident rule: tokIDENT + unambiguous + ambiguous_3 + ambiguous_4.
// Excludes ambiguous_1 and ambiguous_2 keywords (e.g., BEGIN, COMMIT, END).
func (p *Parser) isLabelIdentToken() bool {
	return p.cur.Type == tokIDENT || (p.cur.Type >= 700 && isLabelKeyword(p.cur.Type))
}

// isLvalueIdentToken returns true if the current token can be used as an lvalue identifier.
// This matches MySQL's lvalue_ident rule: tokIDENT + unambiguous + ambiguous_1 + ambiguous_2 + ambiguous_3.
// Excludes ambiguous_4 keywords (GLOBAL, SESSION, LOCAL).
func (p *Parser) isLvalueIdentToken() bool {
	return p.cur.Type == tokIDENT || (p.cur.Type >= 700 && isLvalueKeyword(p.cur.Type))
}

// isVariableRef returns true if the current token is a variable reference.
func (p *Parser) isVariableRef() bool {
	// System variable: @@...
	if p.cur.Type == tokAt2 {
		return true
	}
	// User variable: @name (lexer emits as tokIDENT with "@" prefix)
	if p.cur.Type != tokIDENT {
		return false
	}
	return len(p.cur.Str) > 0 && p.cur.Str[0] == '@'
}

// indexOf returns the index of the first occurrence of ch in s, or -1.
func indexOf(s string, ch byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ch {
			return i
		}
	}
	return -1
}

// eqFold reports whether s and t are equal under Unicode case-folding (ASCII only).
func eqFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		a, b := s[i], t[i]
		if a >= 'A' && a <= 'Z' {
			a += 'a' - 'A'
		}
		if b >= 'A' && b <= 'Z' {
			b += 'a' - 'A'
		}
		if a != b {
			return false
		}
	}
	return true
}

// refuseTableNameDatabase is TableName's action, which answers 1102, ErrWrongDBName, for a
// database name that is empty or ends in a space.
//
// Ref: TableName and isInCorrectIdentifierName in TiDB v8.5.5 pkg/parser/parser.y and
// misc.go.
func refuseTableNameDatabase(name string, pos int) error {
	if name == "" || strings.HasSuffix(name, " ") {
		return &ParseError{Message: "Incorrect database name '" + name + "'", Position: pos}
	}
	return nil
}
