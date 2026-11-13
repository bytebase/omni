package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
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

// keywordCategories maps keyword token types to their category.
//
// A registered keyword that has no entry here is unambiguous: it is accepted
// in every identifier context (plain identifiers, labels, role names and
// SET targets). Only reserved keywords and the four ambiguous categories need
// an explicit entry. categoryOf() is the single lookup that applies this
// default; TestKeywordClassification pins the MySQL 8.0 keywords whose real
// category is not unambiguous.
var keywordCategories = map[int]keywordCategory{
	kwSKIP:                            kwCatReserved,
	kwSELECT:                          kwCatReserved,
	kwINSERT:                          kwCatReserved,
	kwUPDATE:                          kwCatReserved,
	kwDELETE:                          kwCatReserved,
	kwFROM:                            kwCatReserved,
	kwWHERE:                           kwCatReserved,
	kwCREATE:                          kwCatReserved,
	kwDROP:                            kwCatReserved,
	kwALTER:                           kwCatReserved,
	kwTABLE:                           kwCatReserved,
	kwINTO:                            kwCatReserved,
	kwVALUES:                          kwCatReserved,
	kwSET:                             kwCatReserved,
	kwJOIN:                            kwCatReserved,
	kwLEFT:                            kwCatReserved,
	kwRIGHT:                           kwCatReserved,
	kwINNER:                           kwCatReserved,
	kwOUTER:                           kwCatReserved,
	kwON:                              kwCatReserved,
	kwAND:                             kwCatReserved,
	kwOR:                              kwCatReserved,
	kwNOT:                             kwCatReserved,
	kwNULL:                            kwCatReserved,
	kwTRUE:                            kwCatReserved,
	kwFALSE:                           kwCatReserved,
	kwIN:                              kwCatReserved,
	kwBETWEEN:                         kwCatReserved,
	kwLIKE:                            kwCatReserved,
	kwORDER:                           kwCatReserved,
	kwGROUP:                           kwCatReserved,
	kwBY:                              kwCatReserved,
	kwHAVING:                          kwCatReserved,
	kwLIMIT:                           kwCatReserved,
	kwAS:                              kwCatReserved,
	kwIS:                              kwCatReserved,
	kwEXISTS_KW:                       kwCatReserved,
	kwCASE:                            kwCatReserved,
	kwWHEN:                            kwCatReserved,
	kwTHEN:                            kwCatReserved,
	kwELSE:                            kwCatReserved,
	kwEND:                             kwCatAmbiguous2, // demoted from reserved — MySQL 8.0 ambiguous_2 (not label)
	kwIF:                              kwCatReserved,
	kwFOR:                             kwCatReserved,
	kwWHILE:                           kwCatReserved,
	kwINDEX:                           kwCatReserved,
	kwKEY:                             kwCatReserved,
	kwPRIMARY:                         kwCatReserved,
	kwFOREIGN:                         kwCatReserved,
	kwREFERENCES:                      kwCatReserved,
	kwCONSTRAINT:                      kwCatReserved,
	kwUNIQUE:                          kwCatReserved,
	kwCHECK:                           kwCatReserved,
	kwDEFAULT:                         kwCatReserved,
	kwCOLUMN:                          kwCatReserved,
	kwADD:                             kwCatReserved,
	kwCHANGE:                          kwCatReserved,
	kwRENAME:                          kwCatReserved,
	kwGRANT:                           kwCatReserved,
	kwREVOKE:                          kwCatReserved,
	kwALL:                             kwCatReserved,
	kwDISTINCT:                        kwCatReserved,
	kwUNION:                           kwCatReserved,
	kwINTERSECT:                       kwCatUnambiguous,
	kwEXCEPT:                          kwCatUnambiguous,
	kwACCESSIBLE:                      kwCatReserved,
	kwASENSITIVE:                      kwCatReserved,
	kwCUBE:                            kwCatUnambiguous,
	kwCUME_DIST:                       kwCatUnambiguous,
	kwDENSE_RANK:                      kwCatUnambiguous,
	kwDUAL:                            kwCatReserved,
	kwFIRST_VALUE:                     kwCatUnambiguous,
	kwGROUPING:                        kwCatUnambiguous,
	kwINSENSITIVE:                     kwCatReserved,
	kwLAG:                             kwCatUnambiguous,
	kwLAST_VALUE:                      kwCatUnambiguous,
	kwLEAD:                            kwCatUnambiguous,
	kwNTH_VALUE:                       kwCatUnambiguous,
	kwNTILE:                           kwCatUnambiguous,
	kwOF:                              kwCatUnambiguous,
	kwOPTIMIZER_COSTS:                 kwCatUnambiguous,
	kwPERCENT_RANK:                    kwCatUnambiguous,
	kwRANK:                            kwCatUnambiguous,
	kwROW_NUMBER:                      kwCatUnambiguous,
	kwSENSITIVE:                       kwCatReserved,
	kwSPECIFIC:                        kwCatReserved,
	kwUSAGE:                           kwCatReserved,
	kwVARYING:                         kwCatReserved,
	kwDAY_HOUR:                        kwCatReserved,
	kwDAY_MICROSECOND:                 kwCatReserved,
	kwDAY_MINUTE:                      kwCatReserved,
	kwDAY_SECOND:                      kwCatReserved,
	kwHOUR_MICROSECOND:                kwCatReserved,
	kwHOUR_MINUTE:                     kwCatReserved,
	kwHOUR_SECOND:                     kwCatReserved,
	kwMINUTE_MICROSECOND:              kwCatReserved,
	kwMINUTE_SECOND:                   kwCatReserved,
	kwSECOND_MICROSECOND:              kwCatReserved,
	kwYEAR_MONTH:                      kwCatReserved,
	kwUTC_DATE:                        kwCatReserved,
	kwUTC_TIME:                        kwCatReserved,
	kwUTC_TIMESTAMP:                   kwCatReserved,
	kwMAXVALUE:                        kwCatUnambiguous,
	kwNO_WRITE_TO_BINLOG:              kwCatReserved,
	kwIO_AFTER_GTIDS:                  kwCatReserved,
	kwIO_BEFORE_GTIDS:                 kwCatReserved,
	kwSQLEXCEPTION:                    kwCatReserved,
	kwSQLSTATE:                        kwCatReserved,
	kwSQLWARNING:                      kwCatReserved,
	kwCROSS:                           kwCatReserved,
	kwNATURAL:                         kwCatReserved,
	kwUSING:                           kwCatReserved,
	kwASC:                             kwCatReserved,
	kwDESC:                            kwCatReserved,
	kwTO:                              kwCatReserved,
	kwDIV:                             kwCatReserved,
	kwMOD:                             kwCatReserved,
	kwXOR:                             kwCatReserved,
	kwREGEXP:                          kwCatReserved,
	kwBINARY:                          kwCatReserved,
	kwINTERVAL:                        kwCatReserved,
	kwMATCH:                           kwCatReserved,
	kwCURRENT_DATE:                    kwCatReserved,
	kwCURRENT_TIME:                    kwCatReserved,
	kwCURRENT_TIMESTAMP:               kwCatReserved,
	kwCURRENT_USER:                    kwCatReserved,
	kwDATABASE:                        kwCatReserved,
	kwFUNCTION:                        kwCatUnambiguous,
	kwPROCEDURE:                       kwCatReserved,
	kwTRIGGER:                         kwCatReserved,
	kwPARTITION:                       kwCatReserved,
	kwRANGE:                           kwCatUnambiguous,
	kwROW:                             kwCatUnambiguous,
	kwROWS:                            kwCatUnambiguous,
	kwOVER:                            kwCatUnambiguous,
	kwWINDOW:                          kwCatUnambiguous,
	kwFORCE:                           kwCatReserved,
	kwCONVERT:                         kwCatReserved,
	kwCAST:                            kwCatUnambiguous,
	kwWITH:                            kwCatReserved,
	kwREPLACE:                         kwCatReserved,
	kwIGNORE:                          kwCatReserved,
	kwLOAD:                            kwCatReserved,
	kwUSE:                             kwCatReserved,
	kwKILL:                            kwCatReserved,
	kwEXPLAIN:                         kwCatReserved,
	kwSPATIAL:                         kwCatReserved,
	kwFULLTEXT:                        kwCatReserved,
	kwOUTFILE:                         kwCatReserved,
	kwGEOMETRY:                        kwCatUnambiguous,
	kwPOINT:                           kwCatUnambiguous,
	kwLINESTRING:                      kwCatUnambiguous,
	kwPOLYGON:                         kwCatUnambiguous,
	kwMULTIPOINT:                      kwCatUnambiguous,
	kwMULTILINESTRING:                 kwCatUnambiguous,
	kwMULTIPOLYGON:                    kwCatUnambiguous,
	kwGEOMETRYCOLLECTION:              kwCatUnambiguous,
	kwSERIAL:                          kwCatUnambiguous,
	kwNATIONAL:                        kwCatUnambiguous,
	kwNCHAR:                           kwCatUnambiguous,
	kwNVARCHAR:                        kwCatUnambiguous,
	kwSIGNED:                          kwCatAmbiguous2,
	kwPRECISION:                       kwCatReserved,
	kwBOOL:                            kwCatUnambiguous,
	kwBOOLEAN:                         kwCatUnambiguous,
	kwSRID:                            kwCatUnambiguous,
	kwENFORCED:                        kwCatUnambiguous,
	kwLESS:                            kwCatUnambiguous,
	kwTHAN:                            kwCatUnambiguous,
	kwSUBPARTITIONS:                   kwCatUnambiguous,
	kwLEAVES:                          kwCatUnambiguous,
	kwPARSER:                          kwCatUnambiguous,
	kwCOMPRESSION:                     kwCatUnambiguous,
	kwINSERT_METHOD:                   kwCatUnambiguous,
	kwACTION:                          kwCatUnambiguous,
	kwPARTIAL:                         kwCatUnambiguous,
	kwFORMAT:                          kwCatUnambiguous,
	kwXML:                             kwCatUnambiguous,
	kwCONCURRENT:                      kwCatUnambiguous,
	kwWORK:                            kwCatUnambiguous,
	kwXID:                             kwCatUnambiguous,
	kwEXPORT:                          kwCatUnambiguous,
	kwUPGRADE:                         kwCatUnambiguous,
	kwFAST:                            kwCatUnambiguous,
	kwMEDIUM:                          kwCatUnambiguous,
	kwCHANGED:                         kwCatUnambiguous,
	kwCODE:                            kwCatUnambiguous,
	kwEVENTS:                          kwCatUnambiguous,
	kwINDEXES:                         kwCatUnambiguous,
	kwGRANTS:                          kwCatUnambiguous,
	kwTRIGGERS:                        kwCatUnambiguous,
	kwSCHEMAS:                         kwCatReserved,
	kwPARTITIONS:                      kwCatUnambiguous,
	kwHOSTS:                           kwCatUnambiguous,
	kwMUTEX:                           kwCatUnambiguous,
	kwPROFILE:                         kwCatUnambiguous,
	kwREPLICAS:                        kwCatUnambiguous,
	kwNAMES:                           kwCatUnambiguous,
	kwACCOUNT:                         kwCatUnambiguous,
	kwOPTION:                          kwCatReserved,
	kwPROXY:                           kwCatAmbiguous3,
	kwROUTINE:                         kwCatUnambiguous,
	kwEXPIRE:                          kwCatUnambiguous,
	kwNEVER:                           kwCatUnambiguous,
	kwDAY:                             kwCatUnambiguous,
	kwHISTORY:                         kwCatUnambiguous,
	kwREUSE:                           kwCatUnambiguous,
	kwOPTIONAL:                        kwCatUnambiguous,
	kwX509:                            kwCatUnambiguous,
	kwISSUER:                          kwCatUnambiguous,
	kwSUBJECT:                         kwCatUnambiguous,
	kwCIPHER:                          kwCatUnambiguous,
	kwSCHEDULE:                        kwCatUnambiguous,
	kwCOMPLETION:                      kwCatUnambiguous,
	kwPRESERVE:                        kwCatUnambiguous,
	kwEVERY:                           kwCatUnambiguous,
	kwSTARTS:                          kwCatUnambiguous,
	kwENDS:                            kwCatUnambiguous,
	kwVALUE:                           kwCatUnambiguous,
	kwSTACKED:                         kwCatUnambiguous,
	kwUNKNOWN:                         kwCatUnambiguous,
	kwWAIT:                            kwCatUnambiguous,
	kwACTIVE:                          kwCatUnambiguous,
	kwINACTIVE:                        kwCatUnambiguous,
	kwATTRIBUTE:                       kwCatUnambiguous,
	kwADMIN:                           kwCatUnambiguous,
	kwDESCRIPTION:                     kwCatUnambiguous,
	kwORGANIZATION:                    kwCatUnambiguous,
	kwBLOCK_SIZE:                      kwCatUnambiguous,
	kwENABLE_MACRO_BLOCK_BLOOM_FILTER: kwCatUnambiguous,
	kwPCTFREE:                         kwCatUnambiguous,
	kwREPLICA_NUM:                     kwCatUnambiguous,
	kwTABLET_SIZE:                     kwCatUnambiguous,
	kwUSE_BLOOM_FILTER:                kwCatUnambiguous,
	kwAUTO_INCREMENT_MODE:             kwCatUnambiguous,
	kwDYNAMIC_PARTITION_POLICY:        kwCatUnambiguous,
	kwTEMPLATE:                        kwCatUnambiguous,
	kwREFERENCE:                       kwCatUnambiguous,
	kwDEFINITION:                      kwCatUnambiguous,
	kwNAME:                            kwCatUnambiguous,
	kwSYSTEM:                          kwCatUnambiguous,
	kwROTATE:                          kwCatUnambiguous,
	kwKEYRING:                         kwCatUnambiguous,
	kwTLS:                             kwCatUnambiguous,
	kwSTREAM:                          kwCatUnambiguous,
	kwGENERATE:                        kwCatUnambiguous,
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
	kwREAD_WRITE:          kwCatReserved,
	kwSETS:                kwCatUnambiguous,
	kwSTD:                 kwCatUnambiguous,
	kwSTDDEV:              kwCatUnambiguous,
	kwSTDDEV_POP:          kwCatUnambiguous,
	kwSTDDEV_SAMP:         kwCatUnambiguous,
	kwSUBSTR:              kwCatUnambiguous,
	kwSUPER:               kwCatAmbiguous3,
	kwSYSDATE:             kwCatUnambiguous,
	kwTABLESAMPLE:         kwCatUnambiguous,
	kwUNICODE:             kwCatAmbiguous2,
	kwVAR_POP:             kwCatUnambiguous,
	kwVAR_SAMP:            kwCatUnambiguous,
	kwVARCHARACTER:        kwCatReserved,
	kwVARIANCE:            kwCatUnambiguous,
	// --- Section 3.1: Missing reserved classifications ---
	kwANALYZE:             kwCatReserved,
	kwBEFORE:              kwCatReserved,
	kwBIGINT:              kwCatReserved,
	kwBLOB:                kwCatReserved,
	kwBOTH:                kwCatReserved,
	kwCALL:                kwCatReserved,
	kwCASCADE:             kwCatReserved,
	kwCHAR:                kwCatReserved,
	kwCHARACTER:           kwCatReserved,
	kwCOLLATE:             kwCatReserved,
	kwCONDITION:           kwCatReserved,
	kwCONTINUE:            kwCatReserved,
	kwCOUNT:               kwCatUnambiguous,
	kwCURSOR:              kwCatReserved,
	kwDATABASES:           kwCatReserved,
	kwDEC:                 kwCatUnambiguous,
	kwDECIMAL:             kwCatReserved,
	kwDECLARE:             kwCatReserved,
	kwDELAYED:             kwCatReserved,
	kwDESCRIBE:            kwCatReserved,
	kwDETERMINISTIC:       kwCatReserved,
	kwDISTINCTROW:         kwCatReserved,
	kwDOUBLE:              kwCatReserved,
	kwEACH:                kwCatReserved,
	kwELSEIF:              kwCatReserved,
	kwEMPTY:               kwCatUnambiguous,
	kwENCLOSED:            kwCatReserved,
	kwESCAPED:             kwCatReserved,
	kwEXIT:                kwCatReserved,
	kwEXTRACT:             kwCatUnambiguous,
	kwFETCH:               kwCatReserved,
	kwFLOAT:               kwCatReserved,
	kwGENERATED:           kwCatReserved,
	kwGET:                 kwCatReserved,
	kwGROUP_CONCAT:        kwCatUnambiguous,
	kwGROUPS:              kwCatUnambiguous,
	kwHIGH_PRIORITY:       kwCatReserved,
	kwINFILE:              kwCatReserved,
	kwINOUT:               kwCatReserved,
	kwINT:                 kwCatReserved,
	kwINTEGER:             kwCatReserved,
	kwITERATE:             kwCatReserved,
	kwJSON_TABLE:          kwCatUnambiguous,
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
	kwLOOP:                kwCatReserved,
	kwLOW_PRIORITY:        kwCatReserved,
	kwMAX:                 kwCatUnambiguous,
	kwMEDIUMBLOB:          kwCatReserved,
	kwMEDIUMINT:           kwCatReserved,
	kwMEDIUMTEXT:          kwCatReserved,
	kwMIN:                 kwCatUnambiguous,
	kwMODIFIES:            kwCatReserved,
	kwNUMERIC:             kwCatReserved,
	kwOPTIMIZE:            kwCatReserved,
	kwOPTIONALLY:          kwCatReserved,
	kwOUT:                 kwCatReserved,
	kwPOSITION:            kwCatUnambiguous,
	kwPURGE:               kwCatReserved,
	kwREAD:                kwCatReserved,
	kwREADS:               kwCatReserved,
	kwREAL:                kwCatReserved,
	kwRECURSIVE:           kwCatUnambiguous,
	kwRELEASE:             kwCatReserved,
	kwREPEAT:              kwCatReserved,
	kwREQUIRE:             kwCatReserved,
	kwRESIGNAL:            kwCatReserved,
	kwRESTRICT:            kwCatReserved,
	kwRETURN:              kwCatReserved,
	kwRLIKE:               kwCatReserved,
	kwSCHEMA:              kwCatReserved,
	kwSEPARATOR:           kwCatReserved,
	kwSHOW:                kwCatReserved,
	kwSIGNAL:              kwCatReserved,
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
	kwUNDO:                kwCatReserved,
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
	// --- Aurora MySQL S3 extension keywords (LOAD DATA FROM S3, SELECT ... INTO OUTFILE S3) ---
	// Not MySQL 8.0 keywords; unambiguous so they stay usable as identifiers,
	// labels, role names and SET targets.
	kwCSV:       kwCatUnambiguous,
	kwMANIFEST:  kwCatUnambiguous,
	kwOVERWRITE: kwCatUnambiguous,
	kwPREFIX:    kwCatUnambiguous,
	kwSSE_KMS:   kwCatUnambiguous,
	kwSSE_S3:    kwCatUnambiguous,
	// MySQL 8.0 unambiguous keywords that the Aurora syntax consumes as tokens.
	kwENCRYPTION: kwCatUnambiguous,
	kwHEADER:     kwCatUnambiguous,
	kwS3:         kwCatUnambiguous,
	// OceanBase's own reserved words (see their tokens in lexer.go).
	kwNO_PARALLEL:      kwCatReserved,
	kwNO_REWRITE:       kwCatReserved,
	kwREAD_CONSISTENCY: kwCatReserved,
	kwTABLEGROUP:       kwCatReserved,
}

// categoryOf returns the category of a keyword token. Registered keywords
// without an explicit keywordCategories entry are unambiguous. The helpers
// below guard on t >= 700 themselves: a non-keyword token has no category.
func categoryOf(t int) keywordCategory {
	if cat, ok := keywordCategories[t]; ok {
		return cat
	}
	return kwCatUnambiguous
}

// isReserved returns true if the token type is a reserved keyword that cannot
// be used as an unquoted identifier.
func isReserved(t int) bool {
	cat, ok := keywordCategories[t]
	return ok && cat == kwCatReserved
}

// isIdentKeyword returns true if the token type is a non-reserved keyword that
// can be used as an identifier. This covers all 5 non-reserved categories:
// unambiguous, ambiguous_1, ambiguous_2, ambiguous_3, and ambiguous_4.
func isIdentKeyword(t int) bool {
	return t >= 700 && categoryOf(t) != kwCatReserved
}

// isLabelKeyword returns true if the token type is a non-reserved keyword that
// can be used as a statement label. Includes: unambiguous, ambiguous_3, ambiguous_4.
// Excludes: ambiguous_1 (not label, not role), ambiguous_2 (not label).
func isLabelKeyword(t int) bool {
	return t >= 700 && !plReservedWords[t]
}

// plReservedWords are the words a stored program may not use as a label, measured on OceanBase
// 4.3.5.6. A routine is read by another grammar, whose label_ident (pl_parser_mysql_mode.y:1779)
// is any name or unreserved_keyword, and whose lexer knows only its own tokens, so a word such
// as ADD, reserved in SQL, is a name there.
var plReservedWords = map[int]bool{
	kwALTER: true, kwANALYZE: true, kwBEFORE: true, kwBIGINT: true, kwBINARY: true,
	kwBINLOG: true, kwBLOB: true, kwBY: true, kwCACHE: true, kwCALL: true, kwCASE: true,
	kwCHANGE: true, kwCHAR: true, kwCHARACTER: true, kwCHARSET: true, kwCHECKSUM: true,
	kwCOLLATE: true, kwCOMMIT: true, kwCONDITION: true, kwCONTINUE: true, kwCREATE: true,
	kwCURRENT_USER: true, kwCURSOR: true, kwDATABASE: true, kwDEALLOCATE: true, kwDEC: true,
	kwDECIMAL: true, kwDECLARE: true, kwDEFAULT: true, kwDELETE: true, kwDESC: true,
	kwDESCRIBE: true, kwDETERMINISTIC: true, kwDO: true, kwDOUBLE: true, kwDROP: true,
	kwEACH: true, kwELSE: true, kwELSEIF: true, kwERRORS: true, kwEVENTS: true, kwEXECUTE: true,
	kwEXISTS_KW: true, kwEXIT: true, kwEXPLAIN: true, kwFETCH: true, kwFLOAT: true,
	kwFLOAT4: true, kwFLOAT8: true, kwFLUSH: true, kwFOR: true, kwFROM: true, kwGET: true,
	kwGRANT: true, kwIF: true, kwIN: true, kwINDEX: true, kwINOUT: true, kwINSERT: true,
	kwINSTALL: true, kwINT: true, kwINT1: true, kwINT2: true, kwINT3: true, kwINT4: true,
	kwINT8: true, kwINTEGER: true, kwINTO: true, kwIS: true, kwITERATE: true, kwKILL: true,
	kwLEAVE: true, kwLIMIT: true, kwLOAD: true, kwLOCK: true, kwLOGS: true, kwLONG: true,
	kwLONGBLOB: true, kwLONGTEXT: true, kwLOOP: true, kwMASTER: true, kwMEDIUMBLOB: true,
	kwMEDIUMINT: true, kwMEDIUMTEXT: true, kwMIDDLEINT: true, kwMODIFIES: true, kwNOT: true,
	kwNULL: true, kwNUMERIC: true, kwON: true, kwOPTIMIZE: true, kwOR: true, kwOUT: true,
	kwPLUGIN: true, kwPRECISION: true, kwPREPARE: true, kwPRIVILEGES: true, kwPROCEDURE: true,
	kwQUERY: true, kwREAD: true, kwREADS: true, kwREAL: true, kwRELEASE: true, kwRENAME: true,
	kwREPAIR: true, kwREPEAT: true, kwREPLACE: true, kwRESET: true, kwRESIGNAL: true,
	kwRETURN: true, kwREVOKE: true, kwROLLBACK: true, kwSAVEPOINT: true, kwSELECT: true,
	kwSET: true, kwSHOW: true, kwSIGNAL: true, kwSLAVE: true, kwSMALLINT: true, kwSQL: true,
	kwSQLEXCEPTION: true, kwSQLSTATE: true, kwSQLWARNING: true, kwSTART: true, kwSTOP: true,
	kwTABLE: true, kwTABLES: true, kwTEMPORARY: true, kwTHEN: true, kwTINYBLOB: true,
	kwTINYINT: true, kwTINYTEXT: true, kwTRIGGER: true, kwTRUNCATE: true, kwUNINSTALL: true,
	kwUNIQUE: true, kwUNSIGNED: true, kwUNTIL: true, kwUPDATE: true, kwUSER_RESOURCES: true,
	kwUSING: true, kwVARBINARY: true, kwVARCHAR: true, kwVARCHARACTER: true, kwVARYING: true,
	kwVIEW: true, kwWARNINGS: true, kwWHEN: true, kwWHILE: true, kwWITH: true, kwZEROFILL: true,
}

// isAmbiguousRoleWord reports whether the word is in unreserved_keyword_ambiguous_roles
// (sql_parser_mysql_mode.y), which var_name and a role name both refuse, though it is a name
// everywhere else, as in `SELECT event FROM t`.
func isAmbiguousRoleWord(t int) bool {
	switch t {
	case kwSUPER, kwFILE, kwPROCESS, kwNONE, kwEVENT, kwPROXY, kwRELOAD, kwREPLICATION, kwRESOURCE, kwEXECUTE, kwSHUTDOWN:
		return true
	}
	return false
}

// isRoleKeyword returns true if the token type is a non-reserved keyword that
// can be used as a role name. Includes: unambiguous, ambiguous_2, ambiguous_4.
// Excludes: ambiguous_1 (not label, not role), ambiguous_3 (not role).
func isRoleKeyword(t int) bool {
	if t < 700 {
		return false
	}
	cat := categoryOf(t)
	return cat == kwCatUnambiguous || cat == kwCatAmbiguous2 || cat == kwCatAmbiguous4
}

// isLvalueKeyword returns true if the token type is a non-reserved keyword that
// can be used as an lvalue (SET target). Includes: unambiguous, ambiguous_1, ambiguous_2, ambiguous_3.
// Excludes: ambiguous_4 (system variables like GLOBAL, SESSION, LOCAL).
func isLvalueKeyword(t int) bool {
	if t < 700 {
		return false
	}
	cat := categoryOf(t)
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
		if p.curIsVariableOrIntroducer() {
			return "", 0, p.syntaxErrorAtCur()
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

// curIsVariableOrIntroducer reports whether the identifier token at p.cur is a user
// variable or a character set introducer, neither of which is a name. The lexer scans @a
// whole as an identifier for the expression reader, and _utf8 is one word to it, but
// OceanBase's lexer makes them USER_VARIABLE and _UTF8 tokens, so `CREATE TABLE @a (b
// int)`, `SELECT * FROM @a` and `CREATE TABLE _utf8 (b int)` are 1064s. Backquoted, `@a`
// and `_utf8` are ordinary names. A word that is part of a qualified name is a name too,
// because flex reads `d._utf8` and `_utf8.t` whole as one, the longer match: on either side
// of a dot with nothing between, _utf8 names a table or a database.
//
// Ref: user_variable, _UTF8 and the ID_DOT_ID rules in sql_parser_mysql_mode.l
// (oceanbase v4.3.5_CE_BP6).
func (p *Parser) curIsVariableOrIntroducer() bool {
	// The raw byte, as backtickQuotedAt reads it: inputText expands the executable comments
	// before the token, which every name would pay for.
	switch i := p.cur.Loc - p.lexer.baseOffset; {
	case i < 0 || i >= len(p.lexer.input):
		return false
	case p.lexer.input[i] == '@':
		return true
	case p.lexer.input[i] == '_':
		afterDot := p.prev.Type == '.' && p.prev.End == p.cur.Loc
		beforeDot := p.lexer.dotNameAt(p.cur.End - p.lexer.baseOffset)
		return !afterDot && !beforeDot && isCharsetIntroducer(p.cur.Str)
	}
	return false
}

// parseIdentOrText parses MySQL's `ident_or_text` grammar rule: either an
// identifier (including non-reserved keywords) or a string literal. Used for
// aliases and other positions where MySQL accepts TEXT_STRING_sys equivalently.
//
// Ref: mysql-server sql/sql_yacc.yy — ident_or_text rule
//
//	ident_or_text:
//	    ident
//	  | TEXT_STRING_sys
func (p *Parser) parseIdentOrText() (string, int, error) {
	if p.cur.Type == tokSCONST {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	return p.parseIdent()
}

// parseIdentifier is a thin alias for parseIdent, preserved for gradual migration
// of existing call sites. New code should use parseIdent directly.
func (p *Parser) parseIdentifier() (string, int, error) {
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

// parseColumnRef reads column_ref outside an expression (no_param_column_ref has the
// same forms): a column's name, alone or qualified. With r a relation_name, c a
// column_name and K a word of mysql_reserved_keyword:
//
//	c | r.c | r.* | r.r.c | r.r.*
//	r.K | K.K | r.r.K | r.K.K
//	.r.c | .r.K | .K.K
//
// A reserved word is a name only after a dot, or before a dot and another reserved
// word, and a leading dot stands for no database and takes a table and a column after it.
//
// Ref: column_ref and no_param_column_ref in sql_parser_mysql_mode.y (oceanbase
// v4.3.5_CE_BP6, lines 898 and 11301); TestColumnRefFormsOceanBaseOracle holds the
// forms to the pinned server.
func (p *Parser) parseColumnRef() (*nodes.ColumnRef, error) {
	start := p.pos()
	if p.cur.Type == '.' || p.reservedPairAhead() {
		return p.parseReservedColumnRef()
	}
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
		if reservedNamesAfterDot[p.cur.Type] {
			return p.parseColumnRefReservedTail(start, name)
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

			name3, err := p.parseNameAfterDot()
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

// reservedPairAhead reports whether the current token opens column_ref's K.K: a word of
// mysql_reserved_keyword, a dot, and another such word. Only there is a reserved word a
// table's name without a name before it.
func (p *Parser) reservedPairAhead() bool {
	if !reservedNamesAfterDot[p.cur.Type] || p.peekNext().Type != '.' {
		return false
	}
	s := p.saveScan()
	p.advance()
	p.advance()
	ok := reservedNamesAfterDot[p.cur.Type]
	p.restoreScan(s)
	return ok
}

// parseReservedColumnRef reads the column_ref forms that open with a dot or with a
// reserved word: K.K, .r.c, .r.K and .K.K.
func (p *Parser) parseReservedColumnRef() (*nodes.ColumnRef, error) {
	ref := &nodes.ColumnRef{Loc: nodes.Loc{Start: p.pos()}}
	if p.cur.Type != '.' {
		ref.Table = p.advance().Str
		p.advance() // the dot reservedPairAhead saw
		ref.Column = p.advance().Str
		ref.Loc.End = p.prev.End
		p.reservedColumnRef = ref
		return ref, nil
	}
	p.advance() // the leading dot
	tableReserved := reservedNamesAfterDot[p.cur.Type]
	table, err := p.parseNameAfterDot()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('.'); err != nil {
		return nil, err
	}
	// After a reserved table name the column is reserved too.
	if tableReserved && !reservedNamesAfterDot[p.cur.Type] {
		return nil, p.syntaxErrorAtCur()
	}
	column, err := p.parseNameAfterDot()
	if err != nil {
		return nil, err
	}
	ref.Table, ref.Column = table, column
	ref.Loc.End = p.prev.End
	p.reservedColumnRef = ref
	return ref, nil
}

// parseColumnRefReservedTail reads the rest of r.K and r.K.K, at the reserved word after
// the first dot. A second dot takes another reserved word and nothing else.
func (p *Parser) parseColumnRefReservedTail(start int, first string) (*nodes.ColumnRef, error) {
	ref := &nodes.ColumnRef{Loc: nodes.Loc{Start: start}, Table: first, Column: p.advance().Str}
	if p.cur.Type == '.' {
		p.advance()
		if !reservedNamesAfterDot[p.cur.Type] {
			return nil, p.syntaxErrorAtCur()
		}
		ref.Schema, ref.Table, ref.Column = first, ref.Column, p.advance().Str
	}
	ref.Loc.End = p.prev.End
	p.reservedColumnRef = ref
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

// parseRelationFactor reads relation_factor, the name of a table in most statements:
//
//	relation_factor:        normal_relation_factor | dot_relation_factor
//	normal_relation_factor: relation_name
//	                      | relation_name '.' relation_name
//	                      | relation_name '.' mysql_reserved_keyword
//	dot_relation_factor:    '.' relation_name
//	                      | '.' mysql_reserved_keyword
//
// A leading dot names no database: `.t` is the table t of the current one, and the node
// holds no Schema. After a dot, a word of mysql_reserved_keyword is a name too, so
// `SELECT * FROM d . table` and `EXPLAIN . select` differ: TABLE is in that list and SELECT
// is not. Words joined to their dot were names already, since OceanBase's lexer reads
// `d.table` as one token (ID_DOT_ID).
//
// Not read here: normal_relation_factor's database link (`t@link`) and its three-part
// names (`catalog.d.t`).
//
// Ref: relation_factor, normal_relation_factor and dot_relation_factor in
// sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6, lines 14870-14981). The pinned server
// accepts `EXPLAIN . LINEAR` and the row of each other word (testdata/compat.tsv, rule
// mysql_reserved_keyword); TestRelationFactorFormsOceanBaseOracle holds the other
// statements to it.
func (p *Parser) parseRelationFactor() (*nodes.TableRef, error) {
	ref := &nodes.TableRef{Loc: nodes.Loc{Start: p.pos()}}
	if p.cur.Type == '.' {
		p.advance()
		name, err := p.parseNameAfterDot()
		if err != nil {
			return nil, err
		}
		ref.Name = name
		ref.Loc.End = p.prev.End
		return ref, nil
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	ref.Name = name
	if p.cur.Type == '.' {
		p.advance()

		// Completion: after "db.", offer table_ref qualified with database.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
		}

		name2, err := p.parseNameAfterDot()
		if err != nil {
			return nil, err
		}
		ref.Schema, ref.Name = name, name2
	}
	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseNameAfterDot reads the name after a qualifying dot: relation_name, or a word of
// mysql_reserved_keyword, which the grammar lists beside it there.
func (p *Parser) parseNameAfterDot() (string, error) {
	if reservedNamesAfterDot[p.cur.Type] {
		return p.advance().Str, nil
	}
	name, _, err := p.parseIdent()
	return name, err
}

// reservedNamesAfterDot is mysql_reserved_keyword (sql_parser_mysql_mode.y, oceanbase
// v4.3.5_CE_BP6, line 26758): the reserved words that name a table or a column after a dot.
// It is not every reserved word: the grammar leaves out ALL, DISTINCT, FROM, SELECT, UNIQUE,
// WHEN and ten more, and those stay a 1064 there.
var reservedNamesAfterDot = map[int]bool{
	kwACCESSIBLE: true, kwADD: true, kwALTER: true, kwANALYZE: true, kwAND: true, kwAS: true, kwASC: true,
	kwASENSITIVE: true, kwBEFORE: true, kwBETWEEN: true, kwBIGINT: true, kwBINARY: true, kwBLOB: true,
	kwBY: true, kwCALL: true, kwCASCADE: true, kwCASE: true, kwCHANGE: true, kwCHAR: true,
	kwCHARACTER: true, kwCHECK: true, kwCOLLATE: true, kwCOLUMN: true, kwCONDITION: true,
	kwCONSTRAINT: true, kwCONTINUE: true, kwCONVERT: true, kwCREATE: true, kwCROSS: true,
	kwCURRENT_DATE: true, kwCURRENT_TIME: true, kwCURRENT_TIMESTAMP: true, kwCURRENT_USER: true,
	kwCURSOR: true, kwDATABASE: true, kwDATABASES: true, kwDAY_HOUR: true, kwDAY_MICROSECOND: true,
	kwDAY_MINUTE: true, kwDAY_SECOND: true, kwDECLARE: true, kwDECIMAL: true, kwDEFAULT: true,
	kwDELAYED: true, kwDELETE: true, kwDESC: true, kwDESCRIBE: true, kwDETERMINISTIC: true,
	kwDISTINCTROW: true, kwDIV: true, kwDOUBLE: true, kwDROP: true, kwDUAL: true, kwEACH: true,
	kwELSE: true, kwELSEIF: true, kwENCLOSED: true, kwESCAPED: true, kwEXISTS_KW: true, kwEXIT: true,
	kwEXPLAIN: true, kwFETCH: true, kwFLOAT: true, kwFLOAT4: true, kwFLOAT8: true, kwFOR: true,
	kwFORCE: true, kwFOREIGN: true, kwFULLTEXT: true, kwGENERATED: true, kwGET: true, kwGRANT: true,
	kwGROUP: true, kwHAVING: true, kwHIGH_PRIORITY: true, kwHOUR_MICROSECOND: true, kwHOUR_MINUTE: true,
	kwHOUR_SECOND: true, kwIF: true, kwIGNORE: true, kwIN: true, kwINDEX: true, kwINFILE: true,
	kwINNER: true, kwINOUT: true, kwINSENSITIVE: true, kwINSERT: true, kwINT: true, kwINT1: true,
	kwINT2: true, kwINT3: true, kwINT4: true, kwINT8: true, kwINTEGER: true, kwINTERVAL: true,
	kwINTO: true, kwIO_AFTER_GTIDS: true, kwIO_BEFORE_GTIDS: true, kwIS: true, kwITERATE: true,
	kwJOIN: true, kwKEY: true, kwKEYS: true, kwKILL: true, kwLEAVE: true, kwLEFT: true, kwLIKE: true,
	kwLIMIT: true, kwLINEAR: true, kwLINES: true, kwLOAD: true, kwLOCALTIME: true, kwLOCALTIMESTAMP: true,
	kwLONG: true, kwLONGBLOB: true, kwLONGTEXT: true, kwLOOP: true, kwLOW_PRIORITY: true,
	kwMASTER_BIND: true, kwMASTER_SSL_VERIFY_SERVER_CERT: true, kwMATCH: true, kwMEDIUMBLOB: true,
	kwMEDIUMINT: true, kwMEDIUMTEXT: true, kwMIDDLEINT: true, kwMINUTE_MICROSECOND: true,
	kwMINUTE_SECOND: true, kwMOD: true, kwMODIFIES: true, kwNATURAL: true, kwNOT: true,
	kwNO_WRITE_TO_BINLOG: true, kwNUMERIC: true, kwON: true, kwOPTIMIZE: true, kwOPTION: true,
	kwOPTIONALLY: true, kwOR: true, kwORDER: true, kwOUT: true, kwOUTER: true, kwOUTFILE: true,
	kwPARTITION: true, kwPRECISION: true, kwPRIMARY: true, kwPROCEDURE: true, kwPURGE: true, kwREAD: true,
	kwREADS: true, kwREAD_WRITE: true, kwREAL: true, kwREFERENCES: true, kwREGEXP: true, kwRELEASE: true,
	kwRENAME: true, kwREPEAT: true, kwREPLACE: true, kwREQUIRE: true, kwRESIGNAL: true, kwRESTRICT: true,
	kwRETURN: true, kwREVOKE: true, kwRIGHT: true, kwRLIKE: true, kwSCHEMA: true, kwSCHEMAS: true,
	kwSECOND_MICROSECOND: true, kwSENSITIVE: true, kwSEPARATOR: true, kwSET: true, kwSHOW: true,
	kwSIGNAL: true, kwSMALLINT: true, kwSPATIAL: true, kwSPECIFIC: true, kwSQL: true,
	kwSQLEXCEPTION: true, kwSQLSTATE: true, kwSQLWARNING: true, kwSQL_BIG_RESULT: true,
	kwSQL_SMALL_RESULT: true, kwSSL: true, kwSKIP: true, kwSTARTING: true, kwSTORED: true, kwTABLE: true,
	kwTERMINATED: true, kwTHEN: true, kwTINYBLOB: true, kwTINYINT: true, kwTINYTEXT: true, kwTO: true,
	kwTRIGGER: true, kwUNDO: true, kwUNION: true, kwUNLOCK: true, kwUNSIGNED: true, kwUPDATE: true,
	kwUSAGE: true, kwUSE: true, kwUSING: true, kwUTC_DATE: true, kwUTC_TIME: true, kwUTC_TIMESTAMP: true,
	kwVALUES: true, kwVARBINARY: true, kwVARCHAR: true, kwVARCHARACTER: true, kwVARYING: true,
	kwVIRTUAL: true, kwWHERE: true, kwWHILE: true, kwWITH: true, kwWRITE: true, kwXOR: true,
	kwYEAR_MONTH: true, kwZEROFILL: true,
}

// parseSpName reads sp_name, the name of a stored program with its database or without:
//
//	sp_name: ident '.' ident | ident
//
// Ref: src/pl/parser/pl_parser_mysql_mode.y in oceanbase/oceanbase (v4.3.5_CE_BP6, line
// 737). 4.3.5.6 takes CREATE TRIGGER d.tr, `d`.`tr` and d . tr, and answers 1064 for a
// third part (a.d.tr) and for a string in place of either name.
// TestTriggerNameFormsOceanBaseOracle holds each to the engine.
func (p *Parser) parseSpName() (schema, name string, err error) {
	if name, _, err = p.parseIdent(); err != nil {
		return "", "", err
	}
	if p.cur.Type != '.' {
		return "", name, nil
	}
	p.advance()
	schema = name
	if name, _, err = p.parseIdent(); err != nil {
		return "", "", err
	}
	return schema, name, nil
}

// parseTableRefWithAlias reads tbl_name, a table in a FROM clause: its name and the clauses
// parseTableNameClauses reads after it.
func (p *Parser) parseTableRefWithAlias() (*nodes.TableRef, error) {
	ref, err := p.parseRelationFactor()
	if err != nil {
		return nil, err
	}
	if _, err := p.parseTableNameClauses(ref); err != nil {
		return nil, err
	}
	return ref, nil
}

// parseTableNameClauses reads what tbl_name takes after its relation_factor, in the one
// order the grammar has, and reports whether it read any:
//
//	[use_partition] [sample_clause [seed]] [use_flashback] [[AS] relation_name] [index_hint_list]
//
//	use_partition: PARTITION '(' name_list ')' | PARTITION '(' external_table_partitions ')'
//	sample_clause: SAMPLE [BLOCK | HYBRID] [ALL | BASE | INCR] '(' {INTNUM | DECIMAL_VAL} ')'
//	seed:          SEED '(' INTNUM ')'
//	use_flashback: AS OF SNAPSHOT bit_expr | AS OF PROCTIME '(' ')'
//
// tbl_name lists 66 of the 72 combinations: index hints do not follow a flashback that has
// no sample before it, so `FROM t AS OF SNAPSHOT 1 USE KEY (i)` is a 1064 where `FROM t
// SAMPLE (1) AS OF SNAPSHOT 1 USE KEY (i)` parses.
//
// SAMPLE, SEED and OF are names too. Each is its clause only where the next token says so,
// which is all the engine's one token of lookahead sees: `FROM t sample` and `FROM t AS of`
// are aliases.
//
// Ref: tbl_name, use_partition, sample_clause, seed and use_flashback in
// sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6, lines 14080-14741).
func (p *Parser) parseTableNameClauses(ref *nodes.TableRef) (bool, error) {
	read := false
	if p.cur.Type == kwPARTITION {
		if err := p.parseUsePartition(ref); err != nil {
			return false, err
		}
		read = true
	}
	if p.sampleClauseAhead() {
		sample, err := p.parseSampleClause()
		if err != nil {
			return false, err
		}
		ref.Sample = sample
		read = true
	}
	if p.flashbackAhead() {
		flashback, err := p.parseUseFlashback()
		if err != nil {
			return false, err
		}
		ref.Flashback = flashback
		read = true
	}

	// [AS] relation_name. The alias is a name and never a string, which MySQL's
	// ident_or_text would allow: `FROM t 'a'` and `FROM t AS 'a'` are 1064s.
	if _, ok := p.match(kwAS); ok {
		alias, _, err := p.parseIdent()
		if err != nil {
			return false, err
		}
		ref.Alias = alias
		read = true
	} else if p.isIdentToken() && p.cur.Type != kwEXCEPT && p.cur.Type != kwINTERSECT && p.cur.Type != kwWINDOW {
		alias, _, err := p.parseIdent()
		if err != nil {
			return false, err
		}
		ref.Alias = alias
		read = true
	}

	// Optional index hints: USE/FORCE/IGNORE {INDEX|KEY} ...
	if p.cur.Type == kwUSE || p.cur.Type == kwFORCE || p.cur.Type == kwIGNORE {
		if ref.Flashback != nil && ref.Sample == nil {
			return false, p.syntaxErrorAtCur()
		}
		hints, err := p.parseIndexHints()
		if err != nil {
			return false, err
		}
		ref.IndexHints = hints
		read = true
	}

	if read {
		ref.Loc.End = p.prev.End
	}
	return read, nil
}

// parseUsePartition reads use_partition at PARTITION: a list of partition names, or an
// external table's partition as `column = constant` pairs. A name in parentheses is
// followed by ',' or ')' and a pair's by '=', which tells the two apart.
func (p *Parser) parseUsePartition(ref *nodes.TableRef) error {
	if p.peekNext().Type != '(' {
		p.advance()
		return p.syntaxErrorAtCur()
	}
	s := p.saveScan()
	p.advance() // PARTITION
	p.advance() // (
	_, _, nameErr := p.parseIdent()
	pairs := nameErr == nil && p.cur.Type == '='
	p.restoreScan(s)
	if !pairs {
		p.advance() // PARTITION
		parts, err := p.parseParenIdentList()
		if err != nil {
			return err
		}
		ref.Partitions = parts
		return nil
	}
	p.advance() // PARTITION
	p.advance() // (
	for {
		start := p.pos()
		name, _, err := p.parseIdent()
		if err != nil {
			return err
		}
		if _, err := p.expect('='); err != nil {
			return err
		}
		value, err := p.parseExprConst()
		if err != nil {
			return err
		}
		ref.PartitionValues = append(ref.PartitionValues, &nodes.Assignment{
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
			Column: &nodes.ColumnRef{Loc: nodes.Loc{Start: start, End: start + len(name)}, Column: name},
			Value:  value,
		})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, err := p.expect(')')
	return err
}

// parseExprConst reads expr_const: a literal, a system variable or a parameter marker.
//
// Ref: expr_const in sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6).
func (p *Parser) parseExprConst() (nodes.ExprNode, error) {
	if p.cur.Type == tokAt2 || p.cur.Type == '?' {
		return p.parsePrimaryExpr()
	}
	return p.parseLiteralOnly()
}

// nextIsBareWord reports whether the token after the current one is one of words,
// unquoted.
func (p *Parser) nextIsBareWord(words ...string) bool {
	next := p.peekNext()
	if next.Type != tokIDENT && next.Type < 700 {
		return false
	}
	if p.backtickQuotedAt(next.Loc) {
		return false
	}
	for _, w := range words {
		if strings.EqualFold(next.Str, w) {
			return true
		}
	}
	return false
}

// sampleClauseAhead reports whether SAMPLE at the current token opens a sample_clause
// and is not an alias.
func (p *Parser) sampleClauseAhead() bool {
	if !p.isBareWord("SAMPLE") {
		return false
	}
	return p.peekNext().Type == '(' || p.nextIsBareWord("BLOCK", "HYBRID", "ALL", "BASE", "INCR")
}

// parseSampleClause reads sample_clause and the seed that may follow it.
func (p *Parser) parseSampleClause() (*nodes.TableSample, error) {
	sample := &nodes.TableSample{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // SAMPLE
	if p.isBareWord("BLOCK") || p.isBareWord("HYBRID") {
		sample.Method = strings.ToUpper(p.advance().Str)
	}
	if p.isBareWord("ALL") || p.isBareWord("BASE") || p.isBareWord("INCR") {
		sample.Scope = strings.ToUpper(p.advance().Str)
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// INTNUM or DECIMAL_VAL. A number with an exponent is a DECIMAL_VAL too: the lexer
	// never makes an APPROXNUM (testdata/tokens.tsv), and `SAMPLE (1e1)` parses.
	if p.cur.Type != tokICONST && p.cur.Type != tokFCONST {
		return nil, p.syntaxErrorAtCur()
	}
	sample.Percent = p.advance().Str
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	if p.isBareWord("SEED") && p.peekNext().Type == '(' {
		p.advance() // SEED
		p.advance() // (
		if p.cur.Type != tokICONST {
			return nil, p.syntaxErrorAtCur()
		}
		sample.Seed = p.advance().Str
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}
	sample.Loc.End = p.prev.End
	return sample, nil
}

// flashbackAhead reports whether AS at the current token opens use_flashback: AS OF
// SNAPSHOT or AS OF PROCTIME. OF alone after AS is an alias.
func (p *Parser) flashbackAhead() bool {
	if p.cur.Type != kwAS || p.peekNext().Type != kwOF {
		return false
	}
	s := p.saveScan()
	p.advance() // AS
	p.advance() // OF
	ok := p.isBareWord("SNAPSHOT") || p.isBareWord("PROCTIME")
	p.restoreScan(s)
	return ok
}

// parseUseFlashback reads use_flashback at its AS.
func (p *Parser) parseUseFlashback() (*nodes.TableFlashback, error) {
	flashback := &nodes.TableFlashback{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // AS
	p.advance() // OF
	if p.isBareWord("PROCTIME") {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		flashback.Proctime = true
	} else {
		p.advance() // SNAPSHOT
		snapshot, err := p.parseExprPrec(precBitOr)
		if err != nil {
			return nil, err
		}
		if snapshot == nil {
			return nil, p.syntaxErrorAtCur()
		}
		flashback.Snapshot = snapshot
	}
	flashback.Loc.End = p.prev.End
	return flashback, nil
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

		// The name touches @@: `SELECT @@ a` is a 1064.
		if p.cur.Loc != p.prev.End {
			return nil, p.syntaxErrorAtCur()
		}

		ref := &nodes.VariableRef{
			Loc:    nodes.Loc{Start: start},
			System: true,
		}

		// Check for scope prefix: @@GLOBAL.var, @@SESSION.var, @@LOCAL.var, @@PERSIST.var
		isScopeKeyword := p.cur.Type == kwGLOBAL || p.cur.Type == kwSESSION ||
			p.cur.Type == kwLOCAL || p.cur.Type == kwPERSIST ||
			p.cur.Type == kwPERSIST_ONLY
		if isScopeKeyword {
			// Look ahead for '.' to distinguish scope from plain variable name
			next := p.peekNext()
			if next.Type == '.' {
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
	if p.atQuotedUserVar() {
		at := p.advance()
		name := p.advance()
		return &nodes.VariableRef{
			Loc:  nodes.Loc{Start: at.Loc, End: name.End},
			Name: name.Str,
		}, nil
	}
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
func (p *Parser) parseSysVarName() (string, int, error) {
	name, loc, err := p.parseIdent()
	if err != nil {
		return "", 0, err
	}
	// Consume additional .ident parts for component-style variable names
	for p.cur.Type == '.' {
		p.advance() // consume '.'
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

// atQuotedUserVar reports whether p.cur starts @'name' or @"name". The lexer gives the @ alone,
// since 'user'@'host' takes the same text: a string directly after it is the name, and a
// space makes it something else.
//
// user_variable (sql_parser_mysql_mode.l:68) spells a quoted name from a short list of bytes,
// so `@'a b'` is a 1064 here where MySQL and MariaDB read it.
func (p *Parser) atQuotedUserVar() bool {
	if p.cur.Type != tokIDENT || p.cur.Str != "@" {
		return false
	}
	next := p.peekNext()
	if next.Type != tokSCONST || next.Loc != p.cur.End {
		return false
	}
	for _, b := range []byte(p.inputText(next.Loc, next.End)) {
		switch {
		case b >= 'a' && b <= 'z', b >= 'A' && b <= 'Z', b >= '0' && b <= '9', b >= 0x80:
		case b == '_', b == '.', b == '$', b == '/', b == '%', b == '\'', b == '"', b == '`':
		default:
			return false
		}
	}
	return true
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
