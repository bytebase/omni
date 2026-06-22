package parser

// Token type constants.
const (
	tokEOF     = 0
	tokILLEGAL = -1
)

// Literal token types.
const (
	tokIDENT      = iota + 1000 // unquoted identifier
	tokQUOTED     // double-quoted identifier
	tokSTRING     // single-quoted string
	tokINTEGER    // integer constant
	tokFLOAT      // float constant
	tokUUID       // UUID literal
	tokHEX        // hex literal 0xABCD
	tokCODEBLOCK  // $$...$$ code block
)

// Operator / punctuation token types.
const (
	tokDOT      = iota + 2000 // .
	tokCOMMA                  // ,
	tokSEMI                   // ;
	tokCOLON                  // :
	tokLPAREN                 // (
	tokRPAREN                 // )
	tokLBRACE                 // {
	tokRBRACE                 // }
	tokLBRACK                 // [
	tokRBRACK                 // ]
	tokSTAR                   // *
	tokPLUS                   // +
	tokMINUS                  // -
	tokEQ                     // =
	tokLT                     // <
	tokGT                     // >
	tokLTE                    // <=
	tokGTE                    // >=
	tokMINUSMINUS             // --
)

// Keyword token types.
const (
	tokADD           = iota + 3000
	tokAGGREGATE
	tokALL
	tokALLOW
	tokALTER
	tokAND
	tokANN
	tokAPPLY
	tokAS
	tokASC
	tokASCII
	tokAUTHORIZE
	tokBATCH
	tokBEGIN
	tokBIGINT
	tokBLOB
	tokBOOLEAN
	tokBY
	tokCALLED
	tokCLUSTERING
	tokCOMPACT
	tokCONTAINS
	tokCOUNTER
	tokCREATE
	tokCURRENTDATE
	tokCURRENTTIME
	tokCURRENTTIMESTAMP
	tokCURRENTTIMEUUID
	tokCUSTOM
	tokDATE
	tokDATETIMENOW
	tokDECIMAL
	tokDEFAULT
	tokDELETE
	tokDESC
	tokDESCRIBE
	tokDISTINCT
	tokDOUBLE
	tokDROP
	tokDURABLE_WRITES
	tokDURATION
	tokENTRIES
	tokEXECUTE
	tokEXISTS
	tokFALSE
	tokFILTERING
	tokFINALFUNC
	tokFLOATKW
	tokFROM
	tokFROMJSON
	tokFROZEN
	tokFULL
	tokFUNCTION
	tokFUNCTIONS
	tokGRANT
	tokIF
	tokIN
	tokINDEX
	tokINET
	tokINITCOND
	tokINPUT
	tokINSERT
	tokINT
	tokINTO
	tokIS
	tokJSON
	tokKEY
	tokKEYS
	tokKEYSPACE
	tokKEYSPACES
	tokLANGUAGE
	tokLIMIT
	tokLIST
	tokLOGGED
	tokLOGIN
	tokMAP
	tokMATERIALIZED
	tokMAXTIMEUUID
	tokMINTIMEUUID
	tokMODIFY
	tokNORECURSIVE
	tokNOSUPERUSER
	tokNOT
	tokNOW
	tokNULL
	tokOF
	tokON
	tokOPTIONS
	tokOR
	tokORDER
	tokPASSWORD
	tokPERMISSIONS
	tokPRIMARY
	tokRENAME
	tokREPLACE
	tokREPLICATION
	tokRETURNS
	tokREVOKE
	tokROLE
	tokROLES
	tokSAI
	tokSELECT
	tokSET
	tokSFUNC
	tokSMALLINT
	tokSTATIC
	tokSTORAGE
	tokSTORAGEATTACHEDINDEX
	tokSTYPE
	tokSUPERUSER
	tokTABLE
	tokTEXT
	tokTIME
	tokTIMESTAMP
	tokTIMEUUID
	tokTINYINT
	tokTO
	tokTOJSON
	tokTRIGGER
	tokTRUE
	tokTRUNCATE
	tokTTL
	tokTUPLE
	tokTYPE
	tokUNLOGGED
	tokUNSET
	tokUPDATE
	tokUSE
	tokUSER
	tokUSING
	tokUUID_KW
	tokVALUES
	tokVARCHAR
	tokVARINT
	tokVECTOR
	tokVIEW
	tokWHERE
	tokWITH
)

// Token represents a single lexical token.
type Token struct {
	Type int
	Str  string
	Loc  int // byte offset in source
	End  int // exclusive end byte offset
}
